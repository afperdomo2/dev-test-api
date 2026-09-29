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
- [x] Fase 6b — Sync web/shared: UI web (tips admin, import unificado, gamificación, moderación admin)
- [x] Fase 7 — Mobile foundation (auth + API client + sesión) + prompt sesión diaria
- [x] Fase 8 — Verificación final
- [x] Fase 9 — Generación de preguntas con IA (solo admin, por lotes + worker pool)
- [x] Fase 10 — Retirada de la IA en sesiones (modo `generate`)
- [x] Fase 11 — Retirada de la cuota de IA por usuario
- [ ] Backlog — Repaso inteligente (temas débiles), Notificaciones in-app, Insights, Logros, E2E

---

## Decisiones bloqueadas (confirmadas)

| Tema | Decisión |
|---|---|
| V/F | Extender `Question` con `type="true_false"` (2 opciones Verdadero/Falso, 1 correcta). Reusa pipeline create/import/answer. |
| Tips ↔ temas | Un solo tema por tip: `Tip.TopicID` FK indexada. |
| Tips creación | **Solo admin**; siempre públicos. En el **backend** solo manual + import CSV (sin generación de IA). El admin puede usar el prompt helper de `/questions/import` para armar el CSV de tips en una IA externa. Usuario normal solo lectura. |
| Visibilidad | **SIN cambios**: mantener `is_public` + `source` como está. Solo IA y admin son públicos; el usuario normal queda privado. Sin `pending_review`, sin acción "publicar", sin refactor `isPublic → visibility`. |
| Moderación | Solo **reportar contenido público + revisión admin** (sin "publish"). |
| Import UI (web) | Una página unificada con selector "preguntas / V-F / tips"; tips solo visible en la UI para admin. |
| Gamificación | Un solo módulo `internal/modules/gamification/` · **sin tabla `UserStats`** (se usa `UserStreak` + `UserXp` + `UserTopicXp`) · XP por dificultad `beginner=5 / intermediate=10 / advanced=15`. |
| Logros (badges) | **Opt-in futuro** en backlog. Si se adoptan: catálogo en código (`const`), **sin tabla** `UserBadge`. |
| Sesión diaria | Modo `daily` nuevo en `Session` + selección on-the-fly (sin tabla de cola persistente). |
| Informes | **Renombrado a "Insights"** (`/insights/summary`) para no chocar con moderación `/reports` (reportes de contenido). |
| Concurrencia | goroutines + `sync.WaitGroup`; **sin `errgroup`** (evitar dependencia nueva). |
| IA (alcance) | La IA se genera **solo por lotes y solo por admins** (Fase 9). El **backend** (`/generation`) solo genera preguntas (single_choice); **nunca** genera tips. El generador de prompt de `/questions/import` es un helper de texto externo (copy-paste), independiente de la IA del backend. El modo `generate` de las sesiones y la cuota de IA por usuario se retiraron en las **Fases 10 y 11** (hecho). |
| Cuotas | `ImportQuota` (import CSV, con exención de admin) se conserva; la **cuota de IA por usuario se eliminó** (hecho, Fase 11) — el tope es global `ai.max_per_day`. |

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

### Fase 6b — Sync web/shared: UI web
> Orden de ejecución: **6.6 → 6.4 → 6.7 → 6.5**.

- [x] **6.4** Página admin de **tips** (`features/tips/`): `TipsListPage` (filtro por tema, `v-data-table`, delete con confirmación) + `TipFormDialog` + `TipImportPanel` reutilizable. Ruta `/tips` (`requiresAdmin`) + nav item *Tips* en sección Administración. Componente compartido `components/TopicSelect.vue` (autocomplete de un solo tema).
- [x] **6.5** **Import unificado**: selector "Preguntas / Verdadero-Falso / Tips" en `QuestionImportPage` (tips solo para admin). Para tips se reusa `TipImportPanel` (encabezado `content,topic`); para V-F, hint y placeholder propios (`type=true_false`, 2 opciones, 1 correcta). Sin cambios backend.
- [x] **6.6** **Overview de gamificación**: `features/gamification/components/GamificationPanel.vue` (XP total, nivel con barra, racha + freeze, meta diaria, XP por tema) + `DailyGoalDialog` (meta 10–200 XP) + `TopicXpList`. Integrado como **pestaña "Gamificación"** en `ProgressPage` y como **card** en `Dashboard` (solo no-admin). Controles: `setDailyGoalMutation`, `freezeStreakMutation`.
- [x] **6.7** **Moderación admin**: `features/moderation/pages/ReportsListPage.vue` (filtro por status, paginación) + `ReportReviewDialog` (opción de des-publicar solo para `question`; en tips el backend la ignora y se informa en la UI). Ruta `/reports` (`requiresAdmin`) + nav item *Moderación*.
- [x] **6.8** Tests: `tips/gamification/moderation` service + queries, `TipFormDialog`, y nav de `DefaultLayout`. `make fe-format` + `make fe-check` + `make fe-lint` + `make fe-test` (272 tests) + `pnpm --filter @devtest/shared typecheck` + `pnpm --filter mobile typecheck` en verde.
- [x] **6.9** **Prompt de IA por tipo de contenido** en `/questions/import`: unificar "Preguntas" y "Verdadero/Falso" en un solo tab "Preguntas" con selector de tipo (`single_choice` / `multiple_choice` / `true_false`); prompt, placeholder y plantilla CSV adaptan al tipo seleccionado. Añadir generador de prompt para tips (formato `content,topic`). Mover "Descargar plantilla" a la columna Importar. Añadir guía 1-2-3 del flujo copy-paste.

### Fase 8 — Verificación final
- [x] **8.1** `make fmt` + `make vet`.
- [x] **8.2** Tests unitarios por módulo nuevo (`make test`).
- [x] **8.3** `make fe-check` + `make fe-lint` (+ `make fe-test`: 272 tests).
- [x] **8.4** `pnpm --filter @devtest/shared typecheck` + `pnpm --filter mobile typecheck`.

### Fase 9 — Generación de preguntas con IA (solo admin)
> Suplanta a la Fase 9 original de "importación masiva con IA". Decisiones: **un job = una llamada a la IA** con máximo **10** preguntas (no lotes de 20, no `AI_MAX_PER_JOB=50`); el presupuesto es **global** por día, no por usuario; el worker pool es **reutilizable** para futuros schedulers (ver nota en B2.5).
> **Alcance**: esta fase **solo añade** el batch admin. No toca la IA legacy de las sesiones ni la cuota por usuario — se retiran en las Fases 10 y 11.

- [x] **9.1** goose `00002_settings.sql` (tabla `settings` + seeds `ai.*`) y `00003_generation_jobs.sql` (`generation_jobs` + `generation_job_topics`, CHECKs e índices); modelos `Setting`, `GenerationJob`, `GenerationJobTopic`.
- [x] **9.2** Módulo `internal/modules/settings/`: registry tipado de 7 claves (`ai.max_per_job`, `ai.max_per_day`, `ai.default_quantity`, `ai.model`, `ai.temperature`, `ai.request_timeout`, `ai.worker_concurrency`), defaults, validación cruzada, auditoría `updated_by`; `GET/PUT /settings` (admin). `ai.model` y `ai.request_timeout` **no** se siembran: caen a `config.AIConfig` (`.env`).
- [x] **9.3** `internal/services/ai/`: `Settings` + `SettingsProvider` (releído por llamada), timeout por request vía `context`, `GenerateQuestions` con prompt batch en español y parser de array (limpia code fences, valida cada elemento, devuelve los válidos — respuesta parcial permitida). `Temperature *float64` para aceptar el valor válido `0`.
- [x] **9.4** Módulo `internal/modules/generation/`: store con reserva diaria y claim FIFO `FOR UPDATE SKIP LOCKED` en la misma transacción, `recoverRunning` al arranque, persistencia transaccional de preguntas (`source=ai_generated`, `is_public=true`, dueño = admin solicitante); service con validación de topics/dificultad/cantidad y estadísticas; handler + rutas `POST/GET /generation`, `GET /generation/stats`, `GET /generation/{id}` (todas admin).
- [x] **9.5** `internal/modules/generation/worker.go`: pool con concurrencia configurable en caliente desde `ai.worker_concurrency` (reconciliación periódica), apagado ordenado y recuperación de jobs `running` huérfanos.
- [x] **9.6** Wiring en `server.go` (settings → generator → generation) + shutdown del pool antes de cerrar el HTTP server.
- [x] **9.7** Tests unitarios: `settings/service_test.go`, `generation/service_test.go`, `ai/batch_test.go`. `ai.NewGenerator` acepta `SettingsProvider` (nil = config estática) y `internal/modules/sessions/service_test.go` se adaptó.
- [x] **9.8** Migraciones aplicadas y verificadas contra la DB real: seeds, claim FIFO, carga de topics, reserva diaria, `recoverRunning` y CHECKs (todas las pruebas en transacciones con rollback).
- [x] **9.9** `make swagger` (`docs/swagger.yaml` incluye los 4 paths nuevos) + shared/web: `generation.types.ts` (barrel + re-export fino), `generation.service.ts`, `generation.queries.ts` (polling adaptativo 3 s con jobs en vuelo / 30 s en reposo), `GenerationPage` (encolar + stats + historial), `GenerationSettingsPage`, rutas `/generation` y `/generation/settings`, nav *Generar con IA* / *Ajustes de IA* y card de consumo global en el dashboard.
- [x] **9.10** Tests frontend: `generation.service.test.ts`, `generation.queries.test.ts` y cobertura de las entradas de nav en `DefaultLayout.test.ts` → 287 tests en verde.

---

## Fases pendientes

### Fase 7 — Mobile foundation + sesión diaria
> La sesión diaria es mobile-first (prompt al entrar), pero la app mobile no tiene infraestructura aún. Dependencias previas obligatorias.

- [x] **7.1** **Auth mobile**: pantalla de login/setup + `expo-secure-store` para JWT.
- [x] **7.2** **API client mobile** (`src/core/api/client.ts`): fetch + Bearer token + unwrap envelope + RFC 9457 errors + query params.
- [x] **7.3** **Servicios mobile** (`src/core/api/services/`): sessions, questions, progress, gamification (usando `@devtest/shared`).
- [x] **7.4** **Flujo de sesión** en mobile: crear sesión → `NextQuestion` → responder → feedback con `Reward`.
- [x] **7.5** **Sesión diaria (prompt al entrar)**: al abrir la app, preguntar si desea hacer el repaso diario. `POST /sessions/daily` → tip del día + contadores → flujo de estudio. Mostrar racha y meta diaria en Home.

> Nota: la **Fase 8** se ejecutó al cerrar la 6b (6b era la última fase de implementación abierta). Volver a correrla al final de la Fase 7, cuando haya código nuevo de backend/mobile.

### Fase 10 — Retirada de la IA en las sesiones (modo `generate`) ✅
> **Completada (2026-09-28)**. Cerró el trabajo de la **Fase 9**: hoy solo hay UN camino de IA — el batch admin (`/api/v1/generation`). Se eliminó el legacy por sesión (modo `generate`) y la opción "Generar preguntas" del diálogo de Nueva sesión; el modo queda fijo en `review`. Migración `00004_sessions_mode_no_generate.sql` aplicada (la sesión histórica `generate` pasó a `review`, default cambiado, CHECK añadido, columna `questions_generated` eliminada).
> **Numeración**: se llama Fase 10 (no "Fase 2") para no chocar con la Fase 2 histórica de Tips.

**Estado antes de la fase (contexto histórico)** — el modo `generate` estaba **plenamente visible en la web**:

- **Backend**
  - `internal/modules/sessions/service.go` → campo `aiGenerator` (L40), `NewService` (L44-45), `generateBatch()` (L433), `const maxGeneratedPerSession = 10` (L22), y las **3 llamadas**: `Create` (L120), `NextQuestion` (L300-302, con mapeo 429 de `ai.DailyLimitError` en L303), `Answer` (L384, L393).
  - `internal/services/ai/generator.go` → `GenerateQuestion()` (L92, 1 pregunta por llamada), `checkDailyAiLimit()` (L70), `evaluateAiQuota()` (L60), `DailyLimitError` (L31).
  - `internal/modules/sessions/request.go:20` → `Mode` con `binding:"required,oneof=generate review"`.
  - `internal/models/session.go` → `Mode` `default:generate` (L16), `QuestionsGenerated` (L19); `sessions/service.go:483` lo mapea a la response.
- **Shared** — `packages/shared/src/types/session.types.ts`: `SessionMode` incluye `'generate'` (L6), `SESSION_MODES` incluye "Generar preguntas" (L48), y `questionsGenerated` en `SessionResponse` (L91) + `SessionSummaryResponse` (L139).
- **Frontend — sí hay UI de creación de sesión** (`features/sessions/pages/` + `features/sessions/components/`):
  - `SessionsListPage.vue` → botón "Nueva sesión" → diálogo con radio de `SESSION_MODES` (L236), iconos/descripciones (L251, L255), **`getAiQuota()`** (L5, L66-71) con alertas de cupo (L264-319), `aiQuotaExceeds` (L85-94), `createSessionMutation` (L4, L54), default del form `mode: 'review'` (L45).
  - `SessionCard.vue` → `isGenerating` con `questionsGenerated < 2` (L41-45), `SESSION_MODE_LABELS` (L11, L166), `MODE_ICONS` (L57).
  - `SessionStudyPage.vue` → `isGenerating` (L29-34), `watch` que redirige a `/sessions` si la IA aún no generó (L42-48), chequeo en el resumen (L385).

> ✅ **"Generar preguntas" hoy SÍ funciona**: `POST /sessions` con `mode=generate` → `go s.generateBatch(session, 2)` → IA legacy. Si la IA no está configurada en `.env` (`AI_API_KEY`/`AI_API_URL`), el log es silencioso ("AI no configurada — omitiendo generación") y la sesión queda **sin preguntas**; al estudiar responde "No hay más preguntas disponibles".

#### Backend
- [x] **10.1** `sessions/service.go`: quitar `aiGenerator *ai.Generator` del struct (L40) y de `NewService` (L44-45); eliminar `generateBatch()` (L433).
- [x] **10.2** `sessions/service.go`: `Create()` ya no lanza la goroutine `go s.generateBatch(session, 2)` (L120).
- [x] **10.3** `sessions/service.go`: `NextQuestion()` ya no cae a `aiGenerator.GenerateQuestion()` (L300-302) ni mapea `ai.DailyLimitError` → 429 (L303) → devuelve 404 directo.
- [x] **10.4** `sessions/service.go`: `Answer()` ya no dispara `go s.generateBatch(sess, 1)` (L384, L393).
- [x] **10.5** `ai/generator.go`: borrar `GenerateQuestion()` (L92), `checkDailyAiLimit()` (L70), `evaluateAiQuota()` (L60) y `DailyLimitError` (L31) — ⚠️ `evaluateAiQuota`/`DailyLimitError` quedan **sin consumidor tras esta fase** (el batch nunca los usó), así que se borran aquí y no en la 11. **Conservar** `GeneratedQuestion`, `stripCodeFence`, `validateGeneratedQuestion` (los usa el batch).
- [x] **10.6** `ai/prompts.go`: borrar el `systemPrompt` legacy y `buildUserPrompt()`; conservar `batchSystemPrompt` + `buildBatchUserPrompt`.
- [x] **10.7** `sessions/request.go:20`: `Mode` → `binding:"oneof=review"`. ✅ **`daily` NO se acepta aquí** (verificado): la sesión diaria se crea por su endpoint dedicado con `CreateDaily`, que fija el modo internamente.
- [x] **10.8** `server.go`: `sessions.NewService(...)` sin el `*ai.Generator` (el generator solo se inyecta en `generation`).
- [x] **10.9** `models/session.go`: `Mode` L16 → `default:review`; quitar el campo `QuestionsGenerated` (L19) + su mapper (`sessions/service.go:483`) + `maxGeneratedPerSession` (L22).
- [x] **10.10** Migración goose `NNNNN_sessions_mode_no_generate.sql`. ✅ **`sessions.mode` es `character varying`, NO el enum** — verificado en la BD real: `information_schema.columns.data_type = 'character varying'` y `column_default = 'generate'::character varying`. El enum `public.session_mode` (`00001_init.sql:114`) existe pero **la columna no lo usa** (solo se menciona en el `COMMENT ON COLUMN` de L740), coherente con la convención del repo de que los `CREATE TYPE ... AS ENUM` son decorativos. Por tanto es trivial: `ALTER TABLE sessions ALTER COLUMN mode SET DEFAULT 'review'` + `ADD CONSTRAINT ... CHECK (mode IN ('review','daily')) NOT VALID` seguido de `VALIDATE CONSTRAINT` (sin lock de escritura) + `COMMENT ON TYPE` actualizado.
  - ✅ **Decisión tomada**: en la BD de desarrollo hay **1 sesión en modo `generate`** → la migración hace `UPDATE sessions SET mode = 'review' WHERE mode = 'generate'` (conserva el histórico y la hace jugable).
- [x] **10.11** Tests: `sessions/service_test.go` (quitar casos de generate), `ai/generator_test.go` (quitar `TestEvaluateAiQuota` L14 y `TestCheckDailyAiLimit_FailClosedOnDBError` L92).
- [x] **10.12** `make swagger` + `make check` + `make test`.

#### Shared / Frontend
- [x] **10.13** `session.types.ts`: `SessionMode` → `'review' | 'daily'` (L6); quitar `generate` de `SESSION_MODES` (L48), `SESSION_MODE_LABELS`, `SESSION_MODE_ICONS`, `SESSION_MODE_DESCRIPTIONS`; quitar `questionsGenerated` de `SessionResponse` (L91) y `SessionSummaryResponse` (L139).
- [x] **10.14** `SessionsListPage.vue`: ✅ **decisión tomada** — quitar por completo la opción "Generar preguntas": el radio queda con una sola opción (`review`), por lo que **se elimina el `v-radio-group`/`SESSION_MODES` y el modo queda fijo**; borrar `aiQuota` (L66-71), `aiQuotaPercent`, `aiQuotaColor`, `aiQuotaExceeds` (L73-94), los 2 `v-alert` de cupo (L263-319) y el `watch` de invalidación (L96-103); quitar el import de `getAiQuota` (L5).
- [x] **10.15** `SessionCard.vue`: borrar `isGenerating` (L41-45) y su uso en el template; quitar `generate` de `MODE_ICONS` (L57) y de `SESSION_MODE_LABELS` (import L11, uso L166).
- [x] **10.16** `SessionStudyPage.vue`: borrar `isGenerating` (L29-34), el `watch` de redirección por IA (L42-48) y el chequeo de `questionsGenerated` (L385).
- [x] **10.17** Tests web (`SessionCard`, `SessionStudyPage`, `SessionsListPage`, `sessions.queries.test.ts`) + `make fe-format` + `make fe-check` + `make fe-lint` + `pnpm typecheck`.

> ⚠️ **OJO — no romper `daily`**: `mode="daily"` (Fase 5) es independiente de `generate`; la sesión diaria NO usa IA y debe seguir funcionando tras esta fase. `SessionMode` mantiene `'daily'` en el tipo, pero **el diálogo de creación ya no debe ofrecerlo** (se crea sola desde la Fase 5).

### Fase 11 — Retirada de la cuota de IA por usuario ✅
> **Completada (2026-09-28)**. Cerró el alcance de la Fase 9: la IA es **exclusiva de admins** con presupuesto **global** (`ai.max_per_day`). Se eliminaron la ruta `GET /questions/ai-quota`, `GetAiQuota`, `CountAiGeneratedSince`, el campo `dailyAiLimit` de la API/UI de usuarios y la columna `users.daily_ai_limit` (migración `00005` aplicada). `ImportQuota` se conserva intacto.
> ⚠️ **Alcance real**: `DailyAiLimit` no era solo un contador interno — estaba **expuesto en la API de usuarios** y en la **UI de administración de usuarios** (UserFormDialog, UserTable, Dashboard), por eso la fase tocó backend, shared y web.
> Nota de secuencia: al borrarse `checkDailyAiLimit` en la **Fase 10** (10.5), la lógica de cuota ya no tiene consumidor de IA; aquí solo queda retirar la **superficie de administración y visualization** de la cuota.

**Estado antes de la fase (contexto histórico)** — estaba vivo:
- `internal/models/user.go` → `DefaultDailyAiLimit = 20` (L12), campo `DailyAiLimit` `default:20` (L22).
- `internal/modules/users/` → `request.go` (L15 create, L22 update: `binding:"omitempty,min=1,max=50"`), `response.go` (L14, L24), `mapper.go` (L11, L23), `handler.go` (L79), `service.go` (L13, L28, L44-46, L54, L108-109).
- `internal/modules/questions/service.go` → `GetAiQuota()` (L309, L315); `handler.go` + `routes.go` → `GET /questions/ai-quota`.
- `internal/modules/questions/store.go` → `CountAiGeneratedSince()` (el conteo de uso diario).
- `packages/shared/src/types/user.types.ts` → `dailyAiLimit` en las 3 interfaces (L9, L19, L26); `question.types.ts` → tipo `AiQuota`.
- `apps/web`:
  - `features/users/components/UserFormDialog.vue` → input "Cupo IA" (L41, L62, L71, L123, L136, **L231**).
  - `features/users/components/UserTable.vue` → columna "Cupo IA" (L22, L69, **L88**).
  - `features/dashboard/pages/DashboardPage.vue` → **tarjeta de consumo de IA** con barra de progreso (L5, L27-31, L121-163) para no-admins.
  - `features/sessions/pages/SessionsListPage.vue` → alertas de cupo (se van en la 10).
  - `queries/questions.queries.ts` → `aiQuotaOptions()` (L59) — **helper huérfano**, nadie lo consume (los componentes usan `getAiQuota` directo).

#### Backend
- [x] **11.1** `questions/service.go`: borrar `GetAiQuota()` (L309, L315); `handler.go` + `routes.go`: borrar el handler y la ruta `GET /questions/ai-quota`.
- [x] **11.2** `questions/store.go`: borrar `CountAiGeneratedSince()` (⚠️ verificar antes que no lo use otro servicio).
- [x] **11.3** `users/request.go`: quitar `DailyAiLimit` de `CreateUserRequest` y `UpdateUserRequest` (L15, L22).
- [x] **11.4** `users/response.go` + `mapper.go`: quitar `DailyAiLimit` (L14, L24, L11, L23).
- [x] **11.5** `users/service.go`: quitar el parámetro `dailyAiLimit *int` de `Create` (L13, L28) y toda su asignación (L44-46, L54, L108-109); ajustar la interfaz (L13) y el llamador de `handler.go:79`.
- [x] **11.6** `models/user.go`: quitar `DefaultDailyAiLimit` (L12) y el campo `DailyAiLimit` (L22). ✅ **decisión tomada** — la columna `users.daily_ai_limit` (default=20, 2 usuarios en la BD de desarrollo) se **dropea** con migración goose; el dato ya no significa nada y evita lecturas por error.
- [x] **11.7** Tests: `users/handler_test.go` (L23 la declara en el mock), `questions/service_test.go` (L115), `import_test.go` (L243).
- [x] **11.8** `make swagger` + `make check` + `make test`.

#### Shared / Frontend
- [x] **11.9** `user.types.ts`: quitar `dailyAiLimit` de las 3 interfaces (L9, L19, L26).
- [x] **11.10** `question.types.ts`: quitar `AiQuota` + su export en el barrel `index.ts`.
- [x] **11.11** `UserFormDialog.vue`: quitar el input "Cupo IA" (L231) y su carga/envío (L41, L62, L71, L123, L136).
- [x] **11.12** `UserTable.vue`: quitar la columna "Cupo IA" (L22, L69, L88) y su `v-chip` (⚠️ revisar si el template queda desbalanceado).
- [x] **11.13** `DashboardPage.vue`: quitar la tarjeta de consumo de IA (L121-163) y su `useQuery` (L27-31) + el import de `getAiQuota` (L5). ⚠️ **Los admins ya ven el consumo global** (tarjeta de la Fase 9) — no perder ese bloque.
- [x] **11.14** `questions.queries.ts`: borrar `aiQuotaOptions()` (L59, ya huérfano).
- [x] **11.15** Tests web: `DashboardPage.test.ts` (L17 mockea `getAiQuota`), tests de users + `make fe-format` + `make fe-check` + `make fe-lint` + `pnpm typecheck`.

> **Nota de alcance**: `ImportQuota` / `DailyImportLimit` (import CSV) es **independiente** y se conserva, incluida la exención para admins. Esta fase solo toca la **cuota de IA**.


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
- [ ] **B2.5** Scheduler en `server.go` (goroutine + `github.com/robfig/cron/v3` — validar con `dependencies.md` antes de `go get`). **Reutilizar el patrón de pool de `internal/modules/generation/worker.go`** (goroutines + `sync.WaitGroup` + `context` cancelable + reconciliación de concurrencia), no `errgroup`. Revisar si se extrae a un worker genérico compartido.
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
- **Historia de la IA**: la IA en sesiones (modo `generate`, 1 pregunta por llamada) existía desde antes del batch. La **Fase 9** introdujo el batch admin sin tocar el legacy, y se acordó retirarlo después en fases separadas para no hacer un cambio grande de una sola vez. **Fase 10** = quitar el modo `generate`; **Fase 11** = quitar la cuota de IA por usuario (deja de tener consumidor tras la 10). No reintroducir IA por usuario en sesiones.
- Notificaciones in-app: sin secretos ni env vars nuevas; scheduler sin estado compartido.
- Dependencias nuevas → pasar por `.agents/backend/dependencies.md` antes de `go get` / `pnpm add`.
- Docs/`tmp/` son auto-generados y gitignored — nunca editarlos a mano.