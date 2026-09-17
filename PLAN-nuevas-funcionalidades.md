# PLAN — Nuevas funcionalidades (V/F, Tips, Visibilidad, Sesión diaria, Gamificación)

> Backend: Go + Gin + GORM (Postgres). Frontend: Vue 3 + Vite + Vuetify + TanStack Query + Pinia.
> Objetivo: ir ejecutando poco a poco, marcando cada tarea del checklist al completarla.
> Convención: modelos → AutoMigrate en `internal/database/database.go` · `make swagger` tras cada endpoint ·
> `make check` (backend) · `make fe-check` + `make fe-lint` (frontend) · tests con store mockeado.

## Estado general
- [ ] Fase 0 — Fundaciones (modelos + migraciones)
- [ ] Fase 1 — Preguntas Verdadero/Falso
- [ ] Fase 2 — Visibilidad (private / pending_review / published)
- [ ] Fase 3 — Tarjetas de tip
- [ ] Fase 4 — Moderación (reportes)
- [ ] Fase 5 — Gamificación (XP, racha, meta diaria, feedback inmediato)
- [ ] Fase 6 — Sesión diaria inteligente ("Repasar hoy")
- [ ] Fase 7 — Sync web/shared + docs
- [ ] Fase 8 — Verificación final

---

## Decisiones bloqueadas (ya confirmadas)

| Tema | Decisión |
|---|---|
| V/F | Extender `Question` con `type="true_false"` (2 opciones Verdadero/Falso, 1 correcta). Reusa pipeline create/import/answer. |
| Tips ↔ temas | Un solo tema por tip: `Tip.TopicID` FK indexada. |
| Import UI (web) | Una página unificada con selector de tipo "preguntas / V-F / tips"; cambia plantilla CSV + endpoint destino. |
| Gamificación | Un solo módulo `internal/modules/gamification/`. |
| Visibilidad | Reemplazar `IsPublic bool` por `Visibility string` (`private` / `pending_review` / `published`). |
| Sesión diaria | Modo `daily` nuevo en `Session` + selección on-the-fly (sin tabla de cola persistente). |

---

## Convenciones del repo a respetar

- Modelos: `internal/models/<singular_snake>.go`, struct `PascalCase`, tabla pluralizada por GORM, IDs `uuid.UUID` + hook `BeforeCreate`, JSON `camelCase`.
- Módulos: `internal/modules/<plural>/{request,response,mapper,store,service,handler,routes}.go`.
  - `Store` interface + `gormStore` + `NewStore(db)`; `Service` interface + struct + `NewService(deps...)`.
- Errores: `apierr.ErrXxx(title EN, detail ES, instance)`; en el handler fijar `e.Instance = c.Request.URL.Path`.
- Respuestas: `response.Success` / `response.Paginated` / `response.Problem` (envelope `{ data }`, RFC 9457).
- `Question.Type` es `string` (los `CREATE TYPE ... AS ENUM` de `ddl.go` son decorativos; la validación real es Go `binding:oneof`).
- No existe hoy patrón de consultas paralelas: solo goroutines fire-and-forget (`go s.generateBatch(...)` en `sessions/service.go`).
- Logs con emoji (❌✅🚀🛢️🌱🤖).

---

## 1. Preguntas Verdadero/Falso

Extender `Question` con `type="true_false"`. Se guardan 2 opciones (`Verdadero`/`Falso`) con una marcada correcta.

Cambios:
- `questions/request.go`: `CreateQuestionRequest.Type` → añadir `true_false` al `binding:oneof`.
- `questions/service.go` (`Create`): validar que `true_false` traiga exactamente 2 opciones con 1 correcta.
- `questions/import.go`: `validImportTypes` + `true_false`; en `validateImportRow`, para `true_false` exigir 2 opciones y 1 correcta.
- `sessions/service.go` (`evaluateCorrectness`): case `true_false` (misma lógica que `single_choice`).
- `ai/generator.go` (`parseAIResponse`) y `ai/prompts.go`: aceptar/generar `true_false` (opcional, fase posterior).
- `database/ddl.go`: enum `question_type` + `'true_false'`.

---

## 2. Tarjetas de tip

Modelo nuevo `internal/models/tip.go`:
```go
type Tip struct {
    ID         uuid.UUID      // pk
    UserID     uuid.UUID      // not null, index
    User       *User
    TopicID    uuid.UUID      // not null, index  (un solo tema)
    Topic      *Topic
    Content    string         // text not null
    Source     string         // ai_generated|manual|imported (default manual)
    Visibility string         // private|pending_review|published (default private)
    CreatedAt, UpdatedAt, DeletedAt
}
```

Módulo nuevo `internal/modules/tips/` (espejo simplificado de `questions`): request, response, mapper, store, service, handler, routes, import.

Endpoints:
| Método | Ruta | Propósito |
|---|---|---|
| GET | `/api/v1/tips` | Listar (paginado, filtro `topicId`) |
| GET | `/api/v1/tips/:id` | Detalle |
| POST | `/api/v1/tips` | Crear manual |
| PUT | `/api/v1/tips/:id` | Actualizar (dueño) |
| DELETE | `/api/v1/tips/:id` | Eliminar (dueño) |
| POST | `/api/v1/tips/import` | Importar CSV/texto |
| GET | `/api/v1/tips/import-quota` | Cupo diario (reusa `DailyImportLimit`) |
| POST | `/api/v1/tips/:id/publish` | Publicar (ver §3) |

Dependencias del service: `tipStore`, `topicStore` (`FindBySlugAndUser`), `userStore` (`FindByID`), `aiGenerator`.

Import: reutilizar `parseCSV` / `resolveTopicIDs` de `questions` (o extraer a helper compartido). Encabezado: `content,topic,difficulty?,source?` (definir junto a la UI unificada, Fase 7).

---

## 3. Visibilidad + moderación

Reemplazar `IsPublic bool` por `Visibility string` en `questions` y `tips`. Valores: `private`, `pending_review`, `published`.

Reglas:
- Creador no-admin → `private` por defecto.
- Admin / `source='ai_generated'` → `published` (banco compartido).
- `publish` → validación automática por IA: coherente → `published`; dudoso → `pending_review`.
- Filtro unificado: `(visibility='published' OR user_id=? OR source='ai_generated')`.

Lugares a actualizar:
- `questions/store.go`: `FindPage` (~L51), `Stats` (~L172–207).
- `sessions/store.go`: `FindNextQuestion` (~L108), `CountAvailableQuestions` (~L135).
- `ai/generator.go`: `existingContent` (~L147) — solo contenido `published`/`ai_generated`.
- `questions/service.go` (`Import` ~L378, `Create`) y `mapper.go`/DTOs: `isPublic` → `visibility`.
- **Topics quedan fuera**: ya modelan visibilidad con `is_system` + `created_by`.

Validación por IA — nuevo método en `ai` service: `Validate(content string) (verdict, reason string, err error)` (approved/pending). Fallback sin IA: `published` para admin, `pending_review` para no-admin.

Modelo reporte `internal/models/report.go`:
```go
type Report struct {
    ID             uuid.UUID // pk
    ReporterUserID uuid.UUID // index
    TargetType     string    // question|tip
    TargetID       uuid.UUID // index
    Reason         string    // text
    Status         string    // pending|reviewed (default pending)
    CreatedAt, UpdatedAt
}
```

Módulo nuevo `internal/modules/moderation/`:
| Método | Ruta | Propósito |
|---|---|---|
| POST | `/api/v1/reports` | Reportar contenido (`targetType`, `targetId`, `reason`) |
| GET | `/api/v1/reports` (admin) | Listar reportes pendientes |
| PUT | `/api/v1/reports/:id/review` (admin) | Marcar revisado (opcional: des-publicar target) |

`publish` vive en `questions` y `tips` (cada uno posee su store); `reports` vive en `moderation`.

---

## 4. Sesión diaria inteligente ("Repasar hoy")

- `Session.Mode` admite `daily`; añadir `TipID *uuid.UUID` (nullable) a `Session`.
- `database/ddl.go`: enum `session_mode` + `'daily'`.

Endpoint: `POST /api/v1/sessions/daily` → `CreateDaily(userID)`:
1. En paralelo (goroutines + `sync.WaitGroup`) consultar: (a) repasos vencidos (`progress.FindUpcoming`), (b) tip del día (ligado a temas con errores recientes), (c) candidatos `true_false`.
2. Crear `Session{Mode:"daily"}` con `TipID`, `Difficulty` derivada.
3. Lanzar `go generateBatch(sess, 1..2)` (IA en background, igual que hoy) para no bloquear la respuesta.
4. Responder: sesión + tip + contadores (`dueCount`, `warmupCount`).

`NextQuestion` modo `daily` — nuevo método `FindNextDailyQuestion(userID, answeredIDs)`:
1. siguiente repaso vencido (`next_review_at <= now AND is_mastered=false`, no respondida);
2. si no hay, siguiente `true_false` visible no respondida;
3. si no hay, fallback IA (`aiGenerator`, como modo `generate`) o 404.

La sesión manual por tema (`POST /sessions` con `topicIds`) sigue existiendo intacta.

**Patrón de concurrencia**: nuevo — helper pequeño con goroutines + `sync.WaitGroup` + colección de errores para las 3 lecturas independientes. `*gorm.DB` es seguro para uso concurrente (pool). No usar `errgroup` (evitar dependencia nueva).

---

## 5. Racha (streak)

Modelo `internal/models/user_streak.go`:
```go
type UserStreak struct {
    UserID           uuid.UUID // pk
    User             *User
    CurrentStreak    int       // default 0
    LongestStreak    int       // default 0
    LastActivityAt   *time.Time
    LastActivityDate string    // "2006-01-02" (UTC) para comparar días
    FreezesAvailable int       // default 0
    CreatedAt, UpdatedAt
}
```
- Se actualiza en el mismo momento de responder/completar (dentro de `progress.Answer`, ver §8).
- Continuidad: hoy → no cambia; ayer → `CurrentStreak++`; > ayer → se rompe (salvo freeze).
- Freeze (ganado por actividad): +1 freeze cada 7 días consecutivos de racha (tope, ej. 3). `POST /gamification/streak/freeze` consume 1 freeze para cubrir hoy sin actividad.

---

## 6. XP y nivel por tema

Modelos:
```go
type UserXp struct     { UserID uuid.UUID pk; Xp int default 0; UpdatedAt }                          // global
type UserTopicXp struct{ UserID uuid.UUID pk; TopicID uuid.UUID pk; Xp int default 0; UpdatedAt }   // por tema
```
- XP por respuesta correcta según dificultad: `beginner=5`, `intermediate=10`, `advanced=15` (constantes en shared).
- Nivel **derivado** (no almacenado): `level = 1 + floor(sqrt(xp/100))`.
- Se premia global + por cada tema de la pregunta (unión `question_topics`).

---

## 7. Meta diaria configurable

- Añadir `DailyGoalXp int not null default 50` a `User`. Presets cliente: casual=25 / regular=50 / serious=100.
- `GET /gamification/daily-goal` / `PUT /gamification/daily-goal` (o reutilizar `/profile`).
- Progreso diario: sumar XP del día (campo `DailyXpToday` + `DailyXpDate` en `UserStreak`, o desde counters).

---

## 8. Feedback inmediato (XP/racha en la respuesta)

Centralizar el "reward" en `progress.Answer` (único chokepoint de respuesta), que llama a `gamification.RegisterAnswer(userID, questionID, isCorrect)` y devuelve:
```
Reward { xpGained, totalXp, level, streak: { current, longest, freezes }, dailyGoal: { current, target, reached } }
```
- `POST /sessions/:id/answer` → `SessionAnswerResponse` + campo `reward`.
- `POST /sessions/:id/finish` → bonus de sesión + `reward`.
- `POST /progress/:id/answer` → `ProgressResponse` + `reward`.

Evita doble conteo: el XP se otorga solo dentro de `progress.Answer`; `sessions.Answer` ya lo llama y solo propaga el `reward`.

---

## Módulo `gamification`

`internal/modules/gamification/` — store, service, handler, routes, response, request.

Service:
- `RegisterAnswer(userID, questionID, isCorrect) (*Reward, error)` — XP global+tema, streak, progreso meta diaria.
- `CompleteSession(userID) (*Reward, error)` — bonus.
- `Freeze(userID) error`
- `GetOverview(userID) (*OverviewResponse, error)` — streak + XP/level global + meta diaria + XP por tema.
- `GetDailyGoal / SetDailyGoal(userID, xp int)`

Endpoints:
| Método | Ruta | Propósito |
|---|---|---|
| GET | `/api/v1/gamification/overview` | Racha + XP global/nivel + meta diaria + XP por tema |
| PUT | `/api/v1/gamification/daily-goal` | Fijar meta diaria |
| POST | `/api/v1/gamification/streak/freeze` | Congelar un día |

Dependencias: `gamificationStore` (usa `*gorm.DB`; consulta `question_topics` directo → sin ciclos de import).

---

## Cableado en `server.go`

```go
gamificationStore := gamification.NewStore(db)
gamificationService := gamification.NewService(gamificationStore)

progressService := progress.NewService(progressStore, gamificationService)                             // +gamification
sessionService  := sessions.NewService(sessionStore, progressService, aiGenerator, gamificationService) // +gamification

tipService         := tips.NewService(tipStore, topicStore, userStore, aiGenerator)                    // nuevo
moderationService  := moderation.NewService(reportStore)                                               // nuevo
```
Rutas: `tips`, `gamification` bajo `protected`; `moderation` bajo `protected` (admin para GET/PUT reports).

---

## Migraciones

`database.go` `AutoMigrate` añadir: `models.Tip{}`, `models.Report{}`, `models.UserStreak{}`, `models.UserXp{}`, `models.UserTopicXp{}`.

`ddl.go` añadir/actualizar:
- enum `question_type` + `'true_false'`.
- enum `session_mode` + `'daily'`.
- enums nuevos (decorativos): `content_visibility` (`private,pending_review,published`), `report_status` (`pending,reviewed`).
- Backfill: `UPDATE questions SET visibility='published' WHERE is_public OR source='ai_generated'; ... ELSE 'private'`.
- Columnas: `questions.visibility`, `tips.visibility`, `users.daily_goal_xp`, `sessions.tip_id`.
- Índices: `tips(topic_id)`, `reports(target_type,target_id)`, `reports(status)`, `questions(type, visibility)`, PKs compuestos `user_xp`/`user_topic_xp`/`user_streak`.

---

## Checklist ordenado por dependencias

### Fase 0 — Fundaciones (sin dependencias)
- [ ] **0.1** Añadir modelos: `Tip`, `Report`, `UserStreak`, `UserXp`, `UserTopicXp`; campo `Visibility` + `true_false` en `Question`; `DailyGoalXp` en `User`; `TipID` + `Mode="daily"` en `Session`.
- [ ] **0.2** `database.go`: registrar modelos en `AutoMigrate`.
- [ ] **0.3** `ddl.go`: enums (`true_false`, `daily`, `content_visibility`, `report_status`), backfill de `visibility`, columnas e índices.

### Fase 1 — Verdadero/Falso
- [ ] **1.1** `questions/request.go` + `questions/service.go`: aceptar/validar `true_false`.
- [ ] **1.2** `questions/import.go`: tipo `true_false` (2 opciones, 1 correcta).
- [ ] **1.3** `sessions/service.go` `evaluateCorrectness`: case `true_false`.
- [ ] **1.4** Shared TS: `QuestionType` + `'true_false'` + maps/icons.

### Fase 2 — Visibilidad (prerequisito de Tips y Moderación)
- [ ] **2.1** Reemplazar filtros de visibilidad en `questions/store.go`, `sessions/store.go`, `ai/generator.go`.
- [ ] **2.2** `questions/service.go` + `mapper.go`/DTOs: `Visibility`.
- [ ] **2.3** `ai.Generator.Validate(content)` (revisión de coherencia).

### Fase 3 — Tips
- [ ] **3.1** Módulo `tips` (store/service/handler/routes/request/response/mapper).
- [ ] **3.2** Import CSV/texto + `import-quota` para tips.
- [ ] **3.3** `publish` en tips (usa `ai.Validate`).
- [ ] **3.4** (opcional) `ai.Generator.GenerateTip`.

### Fase 4 — Moderación
- [ ] **4.1** Módulo `moderation` (report + admin review).
- [ ] **4.2** `publish` en `questions` (usa `ai.Validate`).

### Fase 5 — Gamificación
- [ ] **5.1** Módulo `gamification` (store + service: XP global/tema, streak, freeze, meta diaria, overview).
- [ ] **5.2** `progress.NewService(store, gamificationService)` + retornar `Reward`.
- [ ] **5.3** `sessions.Answer`/`Finish`: incluir `reward` en la respuesta.
- [ ] **5.4** Rutas `gamification` en `server.go` + wiring.
- [ ] **5.5** Tests `service_test.go` (store mockeado) + `handler_test.go`.

### Fase 6 — Sesión diaria
- [ ] **6.1** `sessions`: `CreateDaily` + `POST /sessions/daily`.
- [ ] **6.2** Store: `FindNextDailyQuestion`, `FindTipOfTheDay` (temas con errores recientes), `FindWarmupCandidates`.
- [ ] **6.3** Helper de concurrencia (goroutines + `sync.WaitGroup`) para el ensamblado.
- [ ] **6.4** `NextQuestion` manejo de modo `daily`.

### Fase 7 — Sync web/shared + docs
- [ ] **7.1** Shared: `tip.types.ts`, `moderation.types.ts`, `gamification.types.ts`, `SessionMode`+`daily`, `Reward`, constants de import unificado.
- [ ] **7.2** Web: servicios `tips/gamification/moderation`, queries, páginas (tips list, selector de tipo en import unificado, overview en Progress/Dashboard, review de reportes admin).
- [ ] **7.3** `make swagger` + actualizar `.agents/backend/data-model.md` y `architecture.md` (fuentes de verdad).

### Fase 8 — Verificación
- [ ] **8.1** `make fmt` + `make vet`.
- [ ] **8.2** Tests unitarios por módulo nuevo (`make test`).
- [ ] **8.3** `make fe-check` + `make fe-lint`.

---

## Flujo de trabajo
1. Marcar tarea como `in_progress` antes de empezar.
2. Implementar + verificar (`make check` / `make fe-check` + `make fe-lint`; tests).
3. Marcar completada y seguir con la siguiente.
4. Al terminar, commit manual (no automático).

---

## Notas
- **Consolidación con `TASKS.md`**: ese archivo cubre gamificación con otro diseño (tabla única `UserStats`, badges, catálogo `const`). Reconciliar antes de ejecutar la Fase 5:
  - `TASKS.md` Fase 1 propone `UserStats` (una tabla) + `UserBadge`; este plan propone `UserStreak` + `UserXp` + `UserTopicXp` (XP global y por tema) sin badges.
  - `TASKS.md` Fase 2 (temas débiles) y Fase 4 (reportes) son complementarias a este plan.
  - Decidir si mantener ambos archivos o fusionar; este archivo es temporal hasta decidir.
- Logros/gamificación: si se adoptan badges, hacer catálogo en código (`const`), no tabla.
- Dependencias nuevas → pasar por `.agents/backend/dependencies.md` antes de `go get` / `pnpm add`.
- Docs/`tmp/` son auto-generados y gitignored — nunca editarlos a mano.