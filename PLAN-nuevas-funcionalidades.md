# PLAN — Nuevas funcionalidades (V/F, Tips, Moderación, Sesión diaria, Gamificación)

> Backend: Go + Gin + GORM (Postgres). Frontend: Vue 3 + Vite + Vuetify + TanStack Query + Pinia.
> Objetivo: ir ejecutando poco a poco, marcando cada tarea del checklist al completarla.
> Convención: modelos → AutoMigrate en `internal/database/database.go` · `make swagger` tras cada endpoint ·
> `make check` (backend) · `make fe-check` + `make fe-lint` (frontend) · tests con store mockeado.

## Estado general
- [ ] Fase 0 — Fundaciones (modelos + migraciones)
- [ ] Fase 1 — Preguntas Verdadero/Falso
- [ ] Fase 2 — Tarjetas de tip (solo admin, públicas)
- [ ] Fase 3 — Moderación (reportes)
- [ ] Fase 4 — Gamificación (XP, racha, meta diaria, feedback inmediato)
- [ ] Fase 5 — Sesión diaria inteligente ("Repasar hoy")
- [ ] Fase 6 — Sync web/shared + docs
- [ ] Fase 7 — Verificación final

---

## Decisiones bloqueadas (confirmadas)

| Tema | Decisión |
|---|---|
| V/F | Extender `Question` con `type="true_false"` (2 opciones Verdadero/Falso, 1 correcta). Reusa pipeline create/import/answer. |
| Tips ↔ temas | Un solo tema por tip: `Tip.TopicID` FK indexada. |
| Tips creación | **Solo admin**; siempre públicos. **Solo manual + import CSV** (sin IA). Usuario normal solo lectura. |
| Visibilidad | **SIN cambios**: mantener `is_public` + `source` como está. Solo IA y admin son públicos; el usuario normal queda privado. Sin `pending_review`, sin acción "publicar", sin refactor `isPublic → visibility`. |
| Moderación | Solo **reportar contenido público + revisión admin** (sin "publish"). |
| Import UI (web) | Una página unificada con selector "preguntas / V-F / tips"; tips solo visible en la UI para admin. |
| Gamificación | Un solo módulo `internal/modules/gamification/`. |
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
- Documentación de esquema: toda tabla, columna, enum, función y trigger (nuevo o existente) debe llevar su `COMMENT ON ...` en `ddl.go`. Las tablas actuales se documentan en la tarea **0.4**; los objetos nuevos, en **0.3**.

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

## 2. Tarjetas de tip (solo admin, públicas)

Modelo `internal/models/tip.go`:
```go
type Tip struct {
    ID         uuid.UUID // pk
    UserID     uuid.UUID // admin que lo crea (not null, index)
    User       *User
    TopicID    uuid.UUID // not null, index  (un solo tema)
    Topic      *Topic
    Content    string    // text not null
    Source     string    // manual|imported (default manual) — sin ai_generated
    CreatedAt, UpdatedAt, DeletedAt
}
```
- **Solo admin** crea/edita/elimina/importa (grupo `admin` en `server.go`).
- **Siempre públicos** (sin columna `visibility`): cualquier usuario autenticado los lee (inicio + sesión diaria).
- **Sin cupo diario** de importación (es admin). Sin `publish`. Sin `GenerateTip` (sin IA).

Endpoints:
| Método | Ruta | Acceso | Propósito |
|---|---|---|---|
| GET | `/api/v1/tips` | autenticado | Listar (paginado, filtro `topicId`) |
| GET | `/api/v1/tips/:id` | autenticado | Detalle |
| POST | `/api/v1/tips` | admin | Crear manual |
| PUT | `/api/v1/tips/:id` | admin | Actualizar |
| DELETE | `/api/v1/tips/:id` | admin | Eliminar |
| POST | `/api/v1/tips/import` | admin | Importar CSV/texto |

- Método público del service: `GetTipOfTheDay(userID)` → usado por la sesión diaria (elige un tip de tema afín a errores recientes o aleatorio).
- Dependencias del service: `tipStore`, `topicStore` (`FindBySlugAndUser` para import).
- Import: reutilizar el parser CSV de `questions` (columna `content` + `topic`; sin `options`). Encabezado a definir con la UI unificada (Fase 6).

---

## 3. Moderación (solo reportes)

**Visibilidad de preguntas: sin cambios.** Se mantiene `is_public` + `source`:
- `source='ai_generated'` y `is_public=true` (admin) → visibles para todos.
- Usuario normal (`manual`/`imported` con `is_public=false`) → privado para su creador.
- Sin `pending_review`, sin acción "publicar", sin refactor de visibilidad.

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
| Método | Ruta | Acceso | Propósito |
|---|---|---|---|
| POST | `/api/v1/reports` | autenticado | Reportar contenido (`targetType`, `targetId`, `reason`) |
| GET | `/api/v1/reports` | admin | Listar reportes pendientes |
| PUT | `/api/v1/reports/:id/review` | admin | Marcar revisado; opcional des-publicar target |

Nota sobre "des-publicar": solo aplica a contenido de admin (`is_public=true` → `false`). El contenido `ai_generated` no se puede "des-publicar" con este mecanismo (retiro = decisión futura, fuera de alcance).

---

## 4. Sesión diaria inteligente ("Repasar hoy")

- `Session.Mode` admite `daily`; añadir `TipID *uuid.UUID` (nullable) a `Session`.
- `database/ddl.go`: enum `session_mode` + `'daily'`.

Endpoint: `POST /api/v1/sessions/daily` → `CreateDaily(userID)`:
1. En paralelo (goroutines + `sync.WaitGroup`) consultar: (a) repasos vencidos (`progress.FindUpcoming`), (b) tip del día (**solo tips admin**, vía `tipService.GetTipOfTheDay`), (c) candidatos `true_false`.
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

tipStore := tips.NewStore(db)
tipService := tips.NewService(tipStore, topicStore)                                                    // nuevo (admin-only; sin aiGenerator)

sessionService := sessions.NewService(sessionStore, progressService, aiGenerator, gamificationService, tipService) // +gamification +tipService

moderationService := moderation.NewService(reportStore)                                                // nuevo
```
Rutas:
- `tips`: lectura (GET) bajo `protected`; CRUD + import bajo grupo `admin`.
- `gamification`: bajo `protected`.
- `moderation`: `POST /reports` bajo `protected`; `GET/PUT /reports` bajo `admin`.

---

## Migraciones

`database.go` `AutoMigrate` añadir: `models.Tip{}`, `models.Report{}`, `models.UserStreak{}`, `models.UserXp{}`, `models.UserTopicXp{}`.

`ddl.go` añadir/actualizar:
- enum `question_type` + `'true_false'`.
- enum `session_mode` + `'daily'`.
- enum nuevo (decorativo): `report_status` (`pending,reviewed`).
- Columnas: `users.daily_goal_xp`, `sessions.tip_id`.
- Índices: `tips(topic_id)`, `reports(target_type,target_id)`, `reports(status)`, `questions(type)`, PKs compuestos `user_xp`/`user_topic_xp`/`user_streak`.
- **Comentarios de documentación** (`COMMENT ON`) para los objetos nuevos: tablas `tips`/`reports`/`user_streak`/`user_xp`/`user_topic_xp`, columnas `users.daily_goal_xp` y `sessions.tip_id`, y enums `true_false`/`daily`/`report_status`. Si en fases posteriores se crean triggers/funciones, documentarlos igual.
- **No**: `questions.visibility`, `tips.visibility`, enum `content_visibility`, backfill de visibilidad.

---

## Checklist ordenado por dependencias

### Fase 0 — Fundaciones (sin dependencias)
- [ ] **0.1** Añadir modelos: `Tip`, `Report`, `UserStreak`, `UserXp`, `UserTopicXp`; `DailyGoalXp` en `User`; `TipID` + `Mode="daily"` en `Session`.
- [ ] **0.2** `database.go`: registrar modelos en `AutoMigrate`.
- [ ] **0.3** `ddl.go`: enums (`true_false`, `daily`, `report_status`), columnas (`daily_goal_xp`, `tip_id`), índices, y **comentarios `COMMENT ON`** de los objetos nuevos (tablas/columnas/enums; y triggers/funciones si se crean).
- [ ] **0.4** Documentar el esquema existente con `COMMENT ON` en `ddl.go`: tablas actuales (`users`, `topics`, `questions`, `question_options`, `code_challenges`, `question_topics`, `user_question_progress`, `sessions`, `session_topics`, `session_answers`), enums existentes (`question_type`, `question_difficulty`, `question_source`, `session_status`, `session_mode`, `session_difficulty`), función `calc_session_score()` y trigger `trg_sessions_score`.

### Fase 1 — Verdadero/Falso
- [ ] **1.1** `questions/request.go` + `questions/service.go`: aceptar/validar `true_false`.
- [ ] **1.2** `questions/import.go`: tipo `true_false` (2 opciones, 1 correcta).
- [ ] **1.3** `sessions/service.go` `evaluateCorrectness`: case `true_false`.
- [ ] **1.4** Shared TS: `QuestionType` + `'true_false'` + maps/icons.

### Fase 2 — Tips (solo admin)
- [ ] **2.1** Módulo `tips` (store/service/handler/routes/request/response/mapper) con CRUD admin + lectura pública.
- [ ] **2.2** Import CSV/texto (admin, sin quota).
- [ ] **2.3** `GetTipOfTheDay(userID)` para la sesión diaria.
- [ ] **2.4** Tests `service_test.go` (store mockeado) + `handler_test.go`.

### Fase 3 — Moderación
- [ ] **3.1** Módulo `moderation` (report + admin review).
- [ ] **3.2** Rutas (POST `protected`; GET/PUT admin) + wiring en `server.go`.
- [ ] **3.3** Tests.

### Fase 4 — Gamificación
- [ ] **4.1** Módulo `gamification` (store + service: XP global/tema, streak, freeze, meta diaria, overview).
- [ ] **4.2** `progress.NewService(store, gamificationService)` + retornar `Reward`.
- [ ] **4.3** `sessions.Answer`/`Finish`: incluir `reward` en la respuesta.
- [ ] **4.4** Rutas `gamification` en `server.go` + wiring.
- [ ] **4.5** Tests `service_test.go` (store mockeado) + `handler_test.go`.

### Fase 5 — Sesión diaria
- [ ] **5.1** `sessions`: `CreateDaily` + `POST /sessions/daily`.
- [ ] **5.2** Store: `FindNextDailyQuestion`, `FindTipOfTheDay` (temas con errores recientes), `FindWarmupCandidates`.
- [ ] **5.3** Helper de concurrencia (goroutines + `sync.WaitGroup`) para el ensamblado.
- [ ] **5.4** `NextQuestion` manejo de modo `daily`.

### Fase 6 — Sync web/shared + docs
- [ ] **6.1** Shared: `true_false` en `question.types.ts`, `tip.types.ts`, `moderation.types.ts`, `gamification.types.ts`, `SessionMode`+`daily`, `Reward`, constants de import unificado.
- [ ] **6.2** Web: servicios `tips/gamification/moderation`, queries, páginas (tips admin, selector de tipo en import unificado, overview en Progress/Dashboard, review de reportes admin).
- [ ] **6.3** `make swagger` + actualizar `.agents/backend/data-model.md` y `architecture.md` (fuentes de verdad).

### Fase 7 — Verificación
- [ ] **7.1** `make fmt` + `make vet`.
- [ ] **7.2** Tests unitarios por módulo nuevo (`make test`).
- [ ] **7.3** `make fe-check` + `make fe-lint`.

---

## Flujo de trabajo
1. Marcar tarea como `in_progress` antes de empezar.
2. Implementar + verificar (`make check` / `make fe-check` + `make fe-lint`; tests).
3. Marcar completada y seguir con la siguiente.
4. Al terminar, commit manual (no automático).

---

## Notas
- **Consolidación con `TASKS.md`**: ese archivo cubre gamificación con otro diseño (tabla única `UserStats`, badges, catálogo `const`). Reconciliar antes de ejecutar la Fase 4:
  - `TASKS.md` Fase 1 propone `UserStats` (una tabla) + `UserBadge`; este plan propone `UserStreak` + `UserXp` + `UserTopicXp` (XP global y por tema) sin badges.
  - `TASKS.md` Fase 2 (temas débiles) y Fase 4 (reportes) son complementarias a este plan.
  - Decidir si mantener ambos archivos o fusionar; este archivo es temporal hasta decidir.
- Logros/gamificación: si se adoptan badges, hacer catálogo en código (`const`), no tabla.
- "Des-publicar" por reporte solo aplica a contenido admin (`is_public=true`). Contenido `ai_generated` requiere mecanismo futuro (fuera de alcance).
- Dependencias nuevas → pasar por `.agents/backend/dependencies.md` antes de `go get` / `pnpm add`.
- Docs/`tmp/` son auto-generados y gitignored — nunca editarlos a mano.