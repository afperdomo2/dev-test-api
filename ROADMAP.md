# ROADMAP — Funcionalidades (V/F, Tips, Moderación, Sesión diaria, Gamificación, Mobile)

> Backend: Go + Gin + GORM (Postgres). Frontend: Vue 3 + Vite + Vuetify + TanStack Query + Pinia. Mobile: React Native + Expo.
> Objetivo: ir ejecutando poco a poco, marcando cada tarea del checklist al completarla.
> Convención: modelos + migración en `internal/database/migrations/` · `make swagger` tras cada endpoint ·
> `make check` (backend) · `make fe-check` + `make fe-lint` (frontend) · tests con store mockeado ·
> `pnpm --filter mobile typecheck`.

> ⚠️ **Esquema de DB**: desde septiembre 2026 se gestiona con migraciones versionadas **goose** (`internal/database/migrations/`), que corren al arranque y son la única fuente de verdad. Las menciones a `AutoMigrate`/`ddl.go` en este documento son **históricas** (fases ya completadas) y quedan reemplazadas; no implementar cambios de esquema por esa vía.

## Estado general
- [x] Fase 0 — Fundaciones (goose migraciones + modelos)
- [x] Fase 1 — Preguntas Verdadero/Falso
- [x] Fase 2 — Tarjetas de tip (solo admin, públicas)
- [x] Fase 3 — Moderación (reportes)
- [x] Fase 4 — Gamificación (XP, racha, meta diaria, feedback inmediato)
- [x] Fase 5 — Sesión diaria inteligente ("Repasar hoy")
- [x] Fase 6a — Sync web/shared: capa de datos (tipos/servicios/queries)
- [ ] Fase 6b — Sync web/shared: UI web (tips admin, import unificado, gamificación, moderación admin)
- [ ] Fase 7 — Mobile foundation (auth + API client + sesión) + prompt sesión diaria
- [ ] Fase 8 — Verificación final
- [ ] Backlog — Repaso inteligente (temas débiles), Notificaciones in-app, Informes, Logros, E2E

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
| Gamificación | Un solo módulo `internal/modules/gamification/` · **sin tabla `UserStats`** (se usa `UserStreak` + `UserXp` + `UserTopicXp`) · XP por dificultad `beginner=5 / intermediate=10 / advanced=15`. |
| Logros (badges) | **Opt-in futuro** en backlog. Si se adoptan: catálogo en código (`const`), **sin tabla** `UserBadge`. |
| Sesión diaria | Modo `daily` nuevo en `Session` + selección on-the-fly (sin tabla de cola persistente). |
| Informes | **Renombrado a "Insights"** (`/insights/summary`) para no chocar con moderación `/reports` (reportes de contenido). |
| Concurrencia | goroutines + `sync.WaitGroup`; **sin `errgroup`** (evitar dependencia nueva). |

---

## Convenciones del repo a respetar

- Modelos: `internal/models/<singular_snake>.go`, struct `PascalCase`, tabla pluralizada por GORM, IDs `uuid.UUID` + hook `BeforeCreate`, JSON `camelCase`.
- Módulos: `internal/modules/<plural>/{request,response,mapper,store,service,handler,routes}.go`.
  - `Store` interface + `gormStore` + `NewStore(db)`; `Service` interface + struct + `NewService(deps...)`.
- Errores: `apierr.ErrXxx(title EN, detail ES, instance)`; en el handler fijar `e.Instance = c.Request.URL.Path`.
- Respuestas: `response.Success` / `response.Paginated` / `response.Problem` (envelope `{ data }`, RFC 9457).
- `Question.Type` es `string` (los `CREATE TYPE ... AS ENUM` creados en las migraciones son decorativos; la validación real es Go `binding:oneof`).
- Logs con emoji (❌✅🚀🛢️🌱🤖).
- Documentación de esquema: toda tabla, columna, enum, función y trigger (nuevo o existente) debe llevar su `COMMENT ON ...` en las **migraciones** (`internal/database/migrations/`). Las tablas actuales se documentan en el baseline `00001_init.sql`; los objetos nuevos, en su migración correspondiente.
- Esquema: desde el sistema goose, la fuente de verdad son las migraciones; **nunca editar una ya aplicada**. Añadir un struct en `internal/models/` ya no cambia la DB: hay que escribir una migración `NNNNN_descripcion.sql`.
- Cambio de esquema (tabla/columna/enum/índice/trigger/comentario) = migración nueva; `Up` idempotente cuando sea seguro; `Down` reversible (o documentado si es destructivo).
- Sesión diaria es **mobile-first** (prompt al entrar en la app). Web **no** tiene botón "Repasar hoy".
- Dependencias nuevas → pasar por `.agents/backend/dependencies.md` (Go) o `.agents/web/dependencies.md` (pnpm) antes de instalar.
- Docs/`tmp/` son auto-generados y gitignored — nunca editarlos a mano.

---

## Fases completadas

### Fase 0 — Fundaciones (goose migraciones + modelos)
- [x] **0.1** Modelos: `Tip`, `Report`, `UserStreak`, `UserXp`, `UserTopicXp`; `DailyGoalXp` en `User`; `TipID` + `Mode="daily"` en `Session`.
- [x] **0.2** goose: `go get github.com/pressly/goose/v3` + `tool .../cmd/goose`; baseline `internal/database/migrations/00001_init.sql` (pg_dump `--schema-only --no-owner --no-privileges`, secciones `-- +goose Up`/`Down`).
- [x] **0.3** `internal/database/migrate.go` (`//go:embed` + `goose.NewProvider` + `WithSessionLocker` + `provider.Up`).
- [x] **0.4** `database.go`: `AutoMigrate` → `runMigrations`; eliminar `runDDL` y `internal/database/ddl.go`; mantener `seedDefaultTopics`.
- [x] **0.5** Makefile targets `migrate-up`/`down`/`status`/`create` + `scripts/migrate.sh`/`migrate.ps1`.
- [x] **0.6** Docs (`architecture.md`, `data-model.md`, `AGENTS.md`, `README.md`).
- [x] **0.7** Verificación: `make check` + `make build` + `make test` (DB limpia y existente).

### Fase 1 — Verdadero/Falso
- [x] **1.1** `questions/request.go` + `questions/service.go`: aceptar/validar `true_false` (2 opciones, 1 correcta).
- [x] **1.2** `questions/import.go`: tipo `true_false`.
- [x] **1.3** `sessions/service.go` `evaluateCorrectness`: case `true_false`.
- [x] **1.4** Shared TS: `QuestionType` + `'true_false'` + maps/icons.

### Fase 2 — Tips (solo admin)
- [x] **2.1** Módulo `tips` (CRUD admin + lectura pública; `GetTipOfTheDay(userID)`).
- [x] **2.2** Import CSV/texto (admin, sin quota).
- [x] **2.3** Tests `service_test.go` (store mockeado) + `handler_test.go`.

### Fase 3 — Moderación
- [x] **3.1** Módulo `moderation` (report + admin review). Rutas: `POST /reports` (protected); `GET /reports` + `PUT /reports/:id/review` (admin).
- [x] **3.2** Wiring en `server.go` + Tests.

### Fase 4 — Gamificación
- [x] **4.1** Módulo `gamification` (store + service: XP global/tema, streak, freeze, meta diaria, overview).
- [x] **4.2** `progress.NewService(store, gamificationService)` + retornar `Reward` (centralizado en `progress.Answer`, evita doble conteo).
- [x] **4.3** `sessions.Answer`/`Finish`: incluir `reward` en la respuesta.
- [x] **4.4** Rutas + wiring en `server.go` + Tests.

### Fase 5 — Sesión diaria
- [x] **5.1** `sessions`: `CreateDaily` + `POST /sessions/daily`.
- [x] **5.2** Store: `FindNextDailyQuestion`, `CountDueReviews`, `CountWarmupCandidates`; tip del día vía `tipService.GetTipOfTheDay`.
- [x] **5.3** Helper de concurrencia (goroutines + `sync.WaitGroup`).
- [x] **5.4** `NextQuestion` manejo de modo `daily`.

### Fase 6a — Sync web/shared: capa de datos
- [x] **6.1** Shared: `tip.types.ts`, `moderation.types.ts`, `gamification.types.ts` (nuevos), `SessionMode` + `daily` + `Reward` + `DailySessionResponse` en `session.types.ts`, barrel `index.ts`. `ImportResult`/`ImportRowError` re-exportados desde `question.types.ts` (evita colisión). Fix: `'daily'` quitado de `SESSION_MODES` (no va en el form de crear sesión; `CreateSessionRequest.Mode` solo acepta `generate|review`).
- [x] **6.2** Web: re-exports finos `tip/moderation/gamification` + `session.types`; servicios `tips/gamification/moderation` y `createDailySession()` en `sessions.service.ts`; queries `tips/gamification/moderation` + `createDailySessionMutation`.
- [x] **6.3** `make swagger` + actualizar `.agents/backend/data-model.md` y `architecture.md` (fuentes de verdad).

---

## Fases pendientes

### Fase 6b — Sync web/shared: UI web
> Orden sugerido de ejecución: **6.6 → 6.4 → 6.7 → 6.5**.

- [ ] **6.4** Página admin de **tips** (`features/tips/`): listado + CRUD + import CSV. Ruta `/tips` (admin) + nav item.
- [ ] **6.5** **Import unificado**: selector "preguntas / V-F / tips" en la página de importación. Tips solo visible para admin. Reusa `importTips`/`importTipsFromText`. Encabezado CSV distinto para tips (`content,topic`).
- [ ] **6.6** **Overview de gamificación** en Progress/Dashboard: XP total, nivel, racha, meta diaria, XP por tema. Controles: `setDailyGoalMutation`, `freezeStreakMutation`.
- [ ] **6.7** **Moderación admin**: listado de reportes con filtro por status, revisión (unpublish). Ruta `/reports` (admin) + nav item.

### Fase 7 — Mobile foundation + sesión diaria
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

## Backlog (absorbido de TASKS.md)

> Diseño de gamificación de TASKS.md Fase 1 (tabla `UserStats` + `UserBadge`) fue **descartado** en favor del diseño implementado (`UserStreak` + `UserXp` + `UserTopicXp`, sin badges). Las siguientes fases permanecen como trabajo futuro.

### B1 — Repaso inteligente (temas débiles)
#### Backend
- [ ] **B1.1** `progress/store.go` → `FindTopicMastery(userID)` (agregación `SessionAnswer` JOIN `question_topics` GROUP BY tema) + `CountDue(userID)`.
- [ ] **B1.2** `progress/service.go` → `TopicMastery()` con `isWeak = accuracy < 60%` y `attempts >= 5`.
- [ ] **B1.3** `response.go` → `TopicMasteryResponse` / `TopicMasteryItem`.
- [ ] **B1.4** Endpoint `GET /progress/topics-mastery` + anotaciones Swagger.
- [ ] **B1.5** Tests service/store.
#### Frontend
- [ ] **B1.6** `types/progress.types.ts` ampliado + service/query.
- [ ] **B1.7** `ProgressPage.vue` — nueva pestaña "Temas" con barras de dominio + badge "Débil".
- [ ] **B1.8** Tests + `make fe-check` + `make fe-lint`.

### B2 — Notificaciones in-app
#### Backend
- [ ] **B2.1** Modelo `internal/models/notification.go` → `Notification` (ID, `UserID`, `Type`, `Title`, `Body`, `IsRead`, `CreatedAt`); migración `NNNNN_notification.sql`.
- [ ] **B2.2** Módulo `internal/modules/notifications/` (store/service/handler/routes).
- [ ] **B2.3** Endpoints: `GET /notifications`, `POST /notifications/:id/read`, `POST /notifications/read-all`, `GET /notifications/unread-count`.
- [ ] **B2.4** Generador: job diario crea "Hoy tienes N preguntas para repasar" (desde `next_review_at <= now`).
- [ ] **B2.5** Scheduler en `server.go` (goroutine + `github.com/robfig/cron/v3` — validar con `dependencies.md` antes de `go get`).
- [ ] **B2.6** Tests.
#### Frontend
- [ ] **B2.7** `types/notification.types.ts` + service + queries.
- [ ] **B2.8** Campana en `v-app-bar` con badge no-leídos + `features/notifications/pages/NotificationsPage.vue` + ruta.
- [ ] **B2.9** Tests + `make fe-check` + `make fe-lint`.

### B3 — Insights (informes semanal/mensual)
> Renombrado desde "reports" para no chocar con el módulo de moderación `moderation` (que ya expone `/reports`).
#### Backend
- [ ] **B3.1** Módulo `internal/modules/insights/` (store/service/handler/routes).
- [ ] **B3.2** Endpoint `GET /insights/summary?period=week|month` → `{ sessions, questionsAnswered, accuracy, timeSpentMs, topTopics[], weakTopics[], xpGained, streakDays, dailyActivity[] }`.
- [ ] **B3.3** Agregaciones desde `Session` + `SessionAnswer` + `question_topics` + tablas de gamificación.
- [ ] **B3.4** Tests.
#### Frontend
- [ ] **B3.5** `types/insight.types.ts` + service + queries.
- [ ] **B3.6** `features/insights/pages/InsightsPage.vue` (selector semana/mes, KPIs, barras por tema) + ruta.
- [ ] **B3.7** Nav: item "Informes" (`mdi-file-chart`, `userOnly`) en sección **Estudio**.
- [ ] **B3.8** Tests + `make fe-check` + `make fe-lint`.

### B4 — Logros (badges) [opt-in]
> Solo si se decide adoptar. Catálogo en código (`const`), **sin tabla** `UserBadge`.
- [ ] **B4.1** `catalog.go` — catálogo de logros como `const`: `first_session`, `streak_3`, `streak_7`, `streak_30`, `perfect_session`, `answers_100`, `topics_mastered_5`, `fast_solver`.
- [ ] **B4.2** Evaluación de criterios en `RecordAnswer`/`RecordSessionFinished` (en el service de gamificación).
- [ ] **B4.3** `GET /gamification/badges` + `user_badges` como tabla (solo si se quiere persistir ganados; decisión pendiente).
- [ ] **B4.4** `features/gamification/pages/AchievementsPage.vue` (grid logros ganados/bloqueados) + ruta.
- [ ] **B4.5** Tests + `make fe-check` + `make fe-lint`.

### B5 — E2E (Playwright) [opcional]
- [ ] **B5.1** Playwright E2E — flujo `status → setup → login → dashboard → sessions`.

---

## Flujo de trabajo
1. Marcar tarea como `in_progress` antes de empezar.
2. Implementar + verificar (`make check` / `make fe-check` + `make fe-lint`; tests).
3. Marcar completada y seguir con la siguiente.
4. Al terminar, commit manual (no automático).

---

## Notas
- **Historia**: este archivo unifica `PLAN-nuevas-funcionalidades.md` y `TASKS.md` (ambos eliminados). El diseño de gamificación de TASKS (tabla única `UserStats` + `UserBadge`) fue descartado por el diseño implementado en PLAN. `TASKS.md` Fase 2 (temas débiles), Fase 3 (notificaciones) y Fase 4 (reportes) pasan al Backlog (informes renombrados a "insights").
- "Des-publicar" por reporte solo aplica a contenido admin (`is_public=true`). Contenido `ai_generated` requiere mecanismo futuro (fuera de alcance).
- Notificaciones in-app: sin secretos ni env vars nuevas; scheduler sin estado compartido.
- Dependencias nuevas → pasar por `.agents/backend/dependencies.md` antes de `go get` / `pnpm add`.
- Docs/`tmp/` son auto-generados y gitignored — nunca editarlos a mano.