# PLAN — Nuevas funcionalidades (V/F, Tips, Moderación, Sesión diaria, Gamificación)

> Backend: Go + Gin + GORM (Postgres). Frontend: Vue 3 + Vite + Vuetify + TanStack Query + Pinia.
> Objetivo: ir ejecutando poco a poco, marcando cada tarea del checklist al completarla.
> Convención: modelos + migración en `internal/database/migrations/` · `make swagger` tras cada endpoint ·
> `make check` (backend) · `make fe-check` + `make fe-lint` (frontend) · tests con store mockeado.

> ⚠️ **Esquema de DB**: desde septiembre 2026 se gestiona con migraciones versionadas **goose** (`internal/database/migrations/`), que corren al arranque y son la única fuente de verdad. Ver **`TASKS.md` Fase 0**. Las menciones a `AutoMigrate`/`ddl.go` en este documento son **históricas** (fases ya completadas) y quedan reemplazadas; no implementar cambios de esquema por esa vía.

## Estado general
- [x] Fase 0 — Fundaciones (modelos + migraciones)
- [x] Fase 1 — Preguntas Verdadero/Falso
- [x] Fase 2 — Tarjetas de tip (solo admin, públicas)
- [x] Fase 3 — Moderación (reportes)
- [x] Fase 4 — Gamificación (XP, racha, meta diaria, feedback inmediato)
- [x] Fase 5 — Sesión diaria inteligente ("Repasar hoy")
- [x] Fase 6a — Sync web/shared: capa de datos (tipos/servicios/queries)
- [ ] Fase 6b — Sync web/shared: UI web (tips admin, import unificado, gamificación, reportes)
- [ ] Fase 7 — Mobile foundation (auth + API client + sesión) + prompt sesión diaria
- [ ] Fase 8 — Verificación final

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
- `Question.Type` es `string` (los `CREATE TYPE ... AS ENUM` creados en las migraciones son decorativos; la validación real es Go `binding:oneof`).
- No existe hoy patrón de consultas paralelas: solo goroutines fire-and-forget (`go s.generateBatch(...)` en `sessions/service.go`).
- Logs con emoji (❌✅🚀🛢️🌱🤖).
- Documentación de esquema: toda tabla, columna, enum, función y trigger (nuevo o existente) debe llevar su `COMMENT ON ...` en las **migraciones** (`internal/database/migrations/`). Las tablas actuales se documentan en el baseline `00001_init.sql`; los objetos nuevos, en su migración correspondiente.

---

## 1. Preguntas Verdadero/Falso

Extender `Question` con `type="true_false"`. Se guardan 2 opciones (`Verdadero`/`Falso`) con una marcada correcta.

Cambios:
- `questions/request.go`: `CreateQuestionRequest.Type` → añadir `true_false` al `binding:oneof`.
- `questions/service.go` (`Create`): validar que `true_false` traiga exactamente 2 opciones con 1 correcta.
- `questions/import.go`: `validImportTypes` + `true_false`; en `validateImportRow`, para `true_false` exigir 2 opciones y 1 correcta.
- `sessions/service.go` (`evaluateCorrectness`): case `true_false` (misma lógica que `single_choice`).
- `ai/generator.go` (`parseAIResponse`) y `ai/prompts.go`: aceptar/generar `true_false` (opcional, fase posterior).
- Migración `internal/database/migrations/`: enum `question_type` + `'true_false'` (ya incluido en el baseline `00001_init.sql`).

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
- Migración `internal/database/migrations/`: enum `session_mode` + `'daily'` (ya incluido en el baseline `00001_init.sql`).

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

> **Nota (2026-09):** esta sección registra **históricamente** cómo se materializó el esquema de estas fases (vía `AutoMigrate` + `ddl.go`). **Sustituida** por el sistema de migraciones versionadas goose — ver **`TASKS.md` Fase 0**. En el nuevo modelo estos cambios viven en el baseline `00001_init.sql` (o en migraciones `NNNNN_*.sql` posteriores), no en `database.go`/`ddl.go`.

Registro histórico de los objetos añadidos: `models.Tip{}`, `models.Report{}`, `models.UserStreak{}`, `models.UserXp{}`, `models.UserTopicXp{}`; enums `true_false`/`daily`/`report_status`; columnas `users.daily_goal_xp`, `sessions.tip_id`; índices `tips(topic_id)`, `reports(target_type,target_id)`, `reports(status)`, `questions(type)`, PKs compuestos `user_xps`/`user_topic_xps`/`user_streaks`; `COMMENT ON` de los objetos nuevos.
- **No** (mantenido): `questions.visibility`, `tips.visibility`, enum `content_visibility`, backfill de visibilidad.

---

## Checklist ordenado por dependencias

### Fase 0 — Fundaciones (sin dependencias)
- [x] **0.1** Añadir modelos: `Tip`, `Report`, `UserStreak`, `UserXp`, `UserTopicXp`; `DailyGoalXp` en `User`; `TipID` + `Mode="daily"` en `Session`.
- [x] **0.2** (histórico) Registrar los modelos nuevos en `AutoMigrate` — hoy: contenido del baseline `00001_init.sql` de goose.
- [x] **0.3** (histórico) `ddl.go`: enums (`true_false`, `daily`, `report_status`), columnas (`daily_goal_xp`, `tip_id`), índices, y **comentarios `COMMENT ON`** de los objetos nuevos — hoy: absorbido en las migraciones de `internal/database/migrations/`.
- [x] **0.4** (histórico) Documentar el esquema existente con `COMMENT ON` en `ddl.go`: tablas actuales, enums existentes, función `calc_session_score()` y trigger `trg_sessions_score` — hoy: `COMMENT ON` incluidos en el baseline `00001_init.sql`.

### Fase 1 — Verdadero/Falso
- [x] **1.1** `questions/request.go` + `questions/service.go`: aceptar/validar `true_false`.
- [x] **1.2** `questions/import.go`: tipo `true_false` (2 opciones, 1 correcta).
- [x] **1.3** `sessions/service.go` `evaluateCorrectness`: case `true_false`.
- [x] **1.4** Shared TS: `QuestionType` + `'true_false'` + maps/icons.

### Fase 2 — Tips (solo admin)
- [x] **2.1** Módulo `tips` (store/service/handler/routes/request/response/mapper) con CRUD admin + lectura pública.
- [x] **2.2** Import CSV/texto (admin, sin quota).
- [x] **2.3** `GetTipOfTheDay(userID)` para la sesión diaria.
- [x] **2.4** Tests `service_test.go` (store mockeado) + `handler_test.go`.

### Fase 3 — Moderación
- [x] **3.1** Módulo `moderation` (report + admin review).
- [x] **3.2** Rutas (POST `protected`; GET/PUT admin) + wiring en `server.go`.
- [x] **3.3** Tests.

### Fase 4 — Gamificación
- [x] **4.1** Módulo `gamification` (store + service: XP global/tema, streak, freeze, meta diaria, overview).
- [x] **4.2** `progress.NewService(store, gamificationService)` + retornar `Reward`.
- [x] **4.3** `sessions.Answer`/`Finish`: incluir `reward` en la respuesta.
- [x] **4.4** Rutas `gamification` en `server.go` + wiring.
- [x] **4.5** Tests `service_test.go` (store mockeado) + `handler_test.go`.

### Fase 5 — Sesión diaria
- [x] **5.1** `sessions`: `CreateDaily` + `POST /sessions/daily`.
- [x] **5.2** Store: `FindNextDailyQuestion`, `CountDueReviews`, `CountWarmupCandidates` (el tip del día viene de `tipService.GetTipOfTheDay`, no del store).
- [x] **5.3** Helper de concurrencia (goroutines + `sync.WaitGroup`) para el ensamblado.
- [x] **5.4** `NextQuestion` manejo de modo `daily`.

### Fase 6a — Sync web/shared: capa de datos (tipos/servicios/queries)
- [x] **6.1** Shared: `tip.types.ts`, `moderation.types.ts`, `gamification.types.ts` (nuevos), `SessionMode` + `daily` + `Reward` + `DailySessionResponse` en `session.types.ts`, barrel `index.ts`. `ImportResult`/`ImportRowError` re-exportados desde `question.types.ts` (evita colisión). Fix: `'daily'` quitado de `SESSION_MODES` (no va en el form de crear sesión; `CreateSessionRequest.Mode` solo acepta `generate|review`).
- [x] **6.2** Web: re-exports finos `tip/moderation/gamification` + `session.types`; servicios `tips/gamification/moderation` y `createDailySession()` en `sessions.service.ts`; queries `tips/gamification/moderation` + `createDailySessionMutation`.
- [x] **6.3** `make swagger` + actualizar `.agents/backend/data-model.md` y `architecture.md` (fuentes de verdad).

### Fase 6b — Sync web/shared: UI web (pendiente)
- [ ] **6.4** Página admin de **tips** (`features/tips/`): listado + CRUD + import CSV. Ruta `/tips` (admin) + nav item.
- [ ] **6.5** **Import unificado**: selector "preguntas / V-F / tips" en la página de importación. Tips solo visible para admin. Reusa `importTips`/`importTipsFromText`. Encabezado CSV distinto para tips (`content,topic`).
- [ ] **6.6** **Overview de gamificación** en Progress/Dashboard: XP total, nivel, racha, meta diaria, XP por tema. Controles: `setDailyGoalMutation`, `freezeStreakMutation`.
- [ ] **6.7** **Moderación admin**: listado de reportes con filtro por status, revisión (unpublish). Ruta `/reports` (admin) + nav item.

### Fase 7 — Mobile foundation + sesión diaria (pendiente)
> La sesión diaria es mobile-first (prompt al entrar), pero la app mobile no tiene infraestructura aún. Dependencias previas obligatorias.

- [ ] **7.1** **Auth mobile**: pantalla de login/setup + `expo-secure-store` para JWT.
- [ ] **7.2** **API client mobile** (`src/core/api/client.ts`): fetch + Bearer token + unwrap envelope + RFC 9457 errors.
- [ ] **7.3** **Servicios mobile** (`src/core/api/services/`): sessions, questions, progress, gamification (usando `@devtest/shared`).
- [ ] **7.4** **Flujo de sesión** en mobile: crear sesión → `NextQuestion` → responder → feedback con `Reward`.
- [ ] **7.5** **Sesión diaria (prompt al entrar)**: al abrir la app, preguntar si desea hacer el repaso diario. `POST /sessions/daily` → tip del día + contadores → flujo de estudio. Mostrar racha y meta diaria en Home/Progreso.

### Fase 8 — Verificación final
- [ ] **8.1** `make fmt` + `make vet`.
- [ ] **8.2** Tests unitarios por módulo nuevo (`make test`).
- [ ] **8.3** `make fe-check` + `make fe-lint`.
- [ ] **8.4** `pnpm --filter @devtest/shared typecheck` + `pnpm --filter mobile typecheck`.

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