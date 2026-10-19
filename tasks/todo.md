# Tareas — Backlog funcional (B1–B5)

> Ordenadas por dependencia (ver `tasks/plan.md`). Un módulo se implementa entero antes de pasar al siguiente, salvo la rama paralela `notifications` ↔ `review-mastery`/`insights`.
>
> **Leyenda**: `[ ]` pendiente · `[~]` en progreso · `[x]` hecha · `[-]` descartada (con motivo).
> **Commits**: manuales, uno por tarea o por grupo acoplado del mismo módulo.

---

## Módulo `review-mastery` — `specs/SPEC-review-mastery.md`

- [x] **RM-1** — Constantes de estado de sesión + agregación por tema
  - Acceptance: `models.SessionStatusInProgress/Completed/Cancelled` definidas en `internal/models/session.go` y usadas en el código nuevo (N-05, D-21); `Store.FindTopicAccuracy(userID)` devuelve `TopicAccuracyRow{TopicID, TopicName, Attempts, CorrectCount}` agregando `session_answers` (sesión `completed`) ⋈ `question_topics` ⋈ `topics`; los temas sin respuestas no aparecen
  - Verify: `make check` + `make build`; contraste de la query contra la DB real con una transacción con rollback (`session_answers` sembrados), confirmando un caso con pregunta multi-tema
  - Files: `internal/models/session.go`, `internal/modules/progress/store.go`, `internal/modules/progress/service_test.go` (stub del mock: añadir un método a la interfaz `Store` obliga a implementarlo)
  - Hecho: `make check` ✅ · `make build` ✅ · `make test` ✅ · 7/7 comprobaciones contra Postgres real en transacción revertida (multi-tema, `in_progress` excluida, aislamiento entre usuarios, usuario sin respuestas)
- [x] **RM-2** — `TopicMastery` en el service
  - Acceptance: `weakAccuracyThreshold = 0.6` y `minAttempts = 5` como const del paquete; `accuracy` redondeada a 4 decimales; `isWeak = accuracy < 0.6 && attempts >= 5`; orden débil → accuracy asc → nombre; errores de store convertidos a `*apierr.APIError` con `instance` vacío
  - Verify: `make test`
  - Files: `internal/modules/progress/service.go`, `internal/modules/progress/service_test.go` (stub del mock de `Store`)
  - Hecho: devuelve slice vacío (no `nil`) para que el envelope sea `{"data":[]}` · stub añadido a `fakeProgressAdapter` (`sessions`) y `mockProgressService` (`progress/handler_test.go`), que también implementan la interfaz · `make check` ✅ · `make test` ✅
- [x] **RM-3** — DTO `TopicMasteryItem`
  - Acceptance: `TopicMasteryItem{TopicID, TopicName, Attempts, CorrectCount, Accuracy, IsWeak}` con JSON camelCase; **sin** wrapper `TopicMasteryResponse`; sin exponer `isMastered` del SRS
  - Verify: `make build`
  - Files: `internal/modules/progress/response.go`
  - Hecho: necesario antes de RM-2 para poder compilar (el service lo devuelve)
- [x] **RM-4** — Endpoint + Swagger
  - Acceptance: `GET /api/v1/progress/topics-mastery` (protegido) responde `200 {"data":[...]}`; handler con `getUserID`, `e.Instance = c.Request.URL.Path` y `response.Problem`; anotaciones Swagger del estilo del módulo
  - Verify: `make swagger` + `make check`; `docs/swagger.yaml` contiene el path; llamada real contra la DB con un usuario sembrado
  - Files: `internal/modules/progress/handler.go`, `internal/modules/progress/routes.go`, `docs/` (generado)
  - Hecho: ruta estática registrada **antes** de `/:question_id` (conviven sin conflicto en Gin 1.12) · `make swagger` ✅ (solo añadidos: 151 líneas, 0 borrados) · `make check` ✅ · verificación HTTP real con router + `middleware.Auth` + JWT firmado: 13/13 checks (401 sin token, 200 con token, envelope `{"data":[…]}` no `null`, claves camelCase sin `isMastered`, invariantes de accuracy/isWeak/orden, coincidencia con el service directo y `/:question_id` sin quedar tapada)
- [x] **RM-5** — Tests backend
  - Acceptance: store mockeado; casos: lista vacía; redondeo; borde exacto de `isWeak` (5 resp./3 aciertos = 0.6 → **no** débil; 5/2 = 0.4 → débil; 4/0 → **no** débil por `attempts`); orden determinista; error de store → `apierr`
  - Verify: `make test`
  - Files: `internal/modules/progress/service_test.go`
  - Hecho: ⏩ **adelantada antes de RM-4** — verifica justo la lógica que entrega RM-2, así que se escribió junto a él. 5 subtests en verde (vacío/no-`nil`, error de store, redondeo a 0.3333, tabla de 7 bordes de `isWeak`, orden `Zeta→Alfa→Delta→Gamma→Beta`)
- [x] **RM-6** — Shared + service + query
  - Acceptance: `TopicMasteryItem` en `packages/shared/src/types/progress.types.ts` + barrel; re-export fino en `apps/web/src/types/progress.types.ts`; `getTopicMastery()` desempaqueta `{ data }`; `topicMasteryOptions()` con `queryKey: ['progress','topics-mastery']` y `staleTime: 30s`
  - Verify: `pnpm --filter @devtest/shared typecheck` + `make fe-check`
  - Files: `packages/shared/src/types/progress.types.ts`, `packages/shared/src/index.ts`, `apps/web/src/types/progress.types.ts`, `apps/web/src/api/services/progress.service.ts`, `apps/web/src/queries/progress.queries.ts`
  - Hecho: el barrel **no** necesitó cambios (`export * from './types/progress.types'`) · `pnpm --filter @devtest/shared typecheck` ✅ · `make fe-check` ✅ · `pnpm --filter mobile typecheck` ✅ · tests existentes de `progress` (10) ✅
  - ⚠️ **Incidente de entorno** (no del código): `pnpm --filter @devtest/shared typecheck` fallaba con `TS2688: Cannot find type definition file for 'react'`. Causa: `node_modules/@types/react` era un directorio huérfano del 2026-09-29 (solo `ts5.0/`, sin `package.json`), que TypeScript descubría como type library implícita. El paquete real vive en `.pnpm/@types+react@19.2.18` y `apps/mobile` resuelve por junction ahí. Se eliminó el directorio huérfano (solo `node_modules`, reversible con `pnpm install`). Si reaparece, el diagnóstico es el mismo.
- [x] **RM-7** — Pestaña "Temas"
  - Acceptance: `activeTab` incluye `'topics'`; `TopicMasteryList.vue` es presentacional (props, sin queries) con `v-progress-linear`, texto `correct/total`, `v-chip` "Débil" con texto (no solo color), `aria-label` con el porcentaje y estado vacío
  - Verify: manual en `make fe-dev` + `make fe-check`
  - Files: `apps/web/src/features/progress/components/TopicMasteryList.vue` (nuevo), `apps/web/src/features/progress/pages/ProgressPage.vue`
  - Hecho: 4ª pestaña **Temas** con rama `v-else-if` **antes** de la tarjeta de preguntas (si no, caería en la rejilla paginada) · `handleRefresh` invalida también `['progress','topics-mastery']` · el componente recibe `items` + `loading` (skeleton dentro de la tarjeta) · chip "Débil" con **texto + icono** (nunca solo color) y barra con `role="progressbar"` + `aria-label`/`aria-valuenow` · sin emojis estructurales (`mdi-*`)
  - ⚠️ **Desviación del spec**: no emite `select`. `SPEC-review-mastery.md` mencionaba "emite `select`" pero no existe acción definida para un tema (no hay drill-down ni filtro). Emitir sin consumidor sería código muerto; queda para cuando exista la acción.
  - Verificado: `make fe-format` ✅ · `make fe-check` ✅ · `make fe-lint` ✅ · `make fe-test` ✅ **291 tests / 49 archivos**. La comprobación visual real (barras, chip, vacío) la cubre el test montado de RM-8 — no la sustituye una revisión a ojo.
- [ ] **RM-8** — Tests web + checks
  - Acceptance: tests de service, query y página (render, badge Débil, vacío); los 3 checks y la suite en verde
  - Verify: `make fe-format` → `make fe-check` → `make fe-lint` → `make fe-test`
  - Files: `apps/web/src/api/services/progress.service.test.ts`, `apps/web/src/queries/progress.queries.test.ts`, test de `ProgressPage`

---

## Módulo `notifications` — `specs/SPEC-notifications.md`

- [ ] **NT-1** — Modelo + migración `00006_notifications.sql`
  - Acceptance: tabla con `creation_day` generada, índice único `(user_id, type, creation_day)`, índices de listado y no-leídas, `COMMENT ON` de tabla y de **todas** las columnas (D-23); `Down` con `DROP TABLE`
  - Verify: `make migrate-up` → `make migrate-status` → `make migrate-down` (la columna generada debe ser aceptada por Postgres; si no, aplicar el fallback documentado en la spec); comprobar el índice único intentando insertar un duplicado
  - Files: `internal/database/migrations/00006_notifications.sql`, `internal/models/notification.go`
- [ ] **NT-2** — Store + service (CRUD)
  - Acceptance: `Store` con listado paginado (`created_at DESC`), `FindByIDAndUser`, `MarkRead`, `MarkAllRead`, `CountUnread`, `InsertBatch` (`OnConflict{DoNothing:true}`); service con `apierr` y `404` para id ajeno; `mapper.ToNotificationResponse` exportado
  - Verify: `make check` + `make test`
  - Files: `internal/modules/notifications/{request,response,mapper,store,service}.go`
- [ ] **NT-3** — Handler + rutas + Swagger
  - Acceptance: los 4 endpoints (`GET /notifications`, `GET /notifications/unread-count`, `POST /notifications/:id/read`, `POST /notifications/read-all`) registrados en `protected`; hermanos static+param sin conflicto (ya verificado en Gin 1.12); anotaciones Swagger
  - Verify: `make swagger` + `make check`; arranque real del servidor sin pánico de rutas; `docs/swagger.yaml` con los 4 paths
  - Files: `internal/modules/notifications/handler.go`, `internal/modules/notifications/routes.go`
- [ ] **NT-4** — Job diario (fan-out idempotente)
  - Acceptance: `DueReviewCounts()` reutiliza la definición de repaso vencido (D-09); `RunDailyJob` no inserta para usuarios con `n = 0`; el texto es `Tienes N preguntas para repasar`; la idempotencia la garantiza el índice único
  - Verify: `make test`; doble ejecución real contra la DB → segunda vez inserta `0`
  - Files: `internal/modules/notifications/store.go`, `internal/modules/notifications/service.go`
- [ ] **NT-5** — Scheduler + wiring
  - Acceptance: `Scheduler` con `time.Timer` + `nextRun` a las 08:00 locales, sin dependencias nuevas (D-17); arranca en su goroutine antes del HTTP server y se detiene **antes** de cerrarlo, con espera acotada (patrón de `server.go:144-168`); logs con emoji
  - Verify: `make build`; arrancar y parar el servidor dos veces → sin goroutines colgadas ni pánico; el log de cierre aparece
  - Files: `internal/modules/notifications/scheduler.go`, `internal/server/server.go`
- [ ] **NT-6** — Tests backend
  - Acceptance: `RunDailyJob` (0 usuarios / N usuarios / error de store); `nextRun` (antes y después de las 08:00, cambio de día); `Start` retorna al cancelar el context; `mark read` idempotente; `404` en id ajeno
  - Verify: `make test`
  - Files: `internal/modules/notifications/service_test.go`, `internal/modules/notifications/scheduler_test.go`
- [ ] **NT-7** — Shared + service + queries
  - Acceptance: `notification.types.ts` (`NotificationType`, `Notification`, `UnreadCount`) + barrel; re-export fino; 4 funciones de servicio; `unreadCountOptions()` con `refetchInterval: 60s` y `refetchOnWindowFocus`; mutaciones que invalidan `['notifications']`
  - Verify: `pnpm --filter @devtest/shared typecheck` + `make fe-check`
  - Files: `packages/shared/src/types/notification.types.ts`, `packages/shared/src/index.ts`, `apps/web/src/types/notification.types.ts`, `apps/web/src/api/services/notifications.service.ts`, `apps/web/src/queries/notifications.queries.ts`
- [ ] **NT-8** — Campana + página + ruta
  - Acceptance: campana en `v-app-bar` con `v-badge` (oculto en 0) y `aria-label`; `NotificationsPage` con lista paginada, "marcar todas como leídas", no-leído con texto+color, estado vacío; ruta `/notifications` (`requiresAuth`)
  - Verify: manual en `make fe-dev` + `make fe-check`
  - Files: `apps/web/src/features/notifications/pages/NotificationsPage.vue`, `apps/web/src/layouts/DefaultLayout.vue`, `apps/web/src/router/routes.ts`
- [ ] **NT-9** — Tests web + checks
  - Acceptance: tests de service, queries, campana en `DefaultLayout` y página; checks y suite en verde
  - Verify: `make fe-format` → `make fe-check` → `make fe-lint` → `make fe-test`
  - Files: `apps/web/src/api/services/notifications.service.test.ts`, `apps/web/src/queries/notifications.queries.test.ts`, `apps/web/src/layouts/DefaultLayout.test.ts`, test de `NotificationsPage`

---

## Módulo `insights` — `specs/SPEC-insights.md`

- [ ] **IN-1** — Helpers de XP en `gamification` + esqueleto del módulo
  - Acceptance: `gamification.XpForDifficulty(difficulty) int` y `gamification.SessionBonusXp() int` exportados y usados internamente por `gamification` (una sola fuente del XP, N-01); módulo `insights` creado con la anatomía estándar
  - Verify: `make check` + `make test` (los tests existentes de gamification siguen verdes)
  - Files: `internal/modules/gamification/service.go`, `internal/modules/insights/{request,response,store,service,handler,routes}.go`
- [ ] **IN-2** — Agregaciones del store + `Summary`
  - Acceptance: `CountSessions`, `CountAnswers` (total y correctas), `SumResponseTimeMs`, `TopTopics(limit)`, `DailyActivity`, `CountCorrectByDifficulty`, `CountCompletedSessions`; filtran por `created_at` en la ventana y sesión `completed` (`models.SessionStatusCompleted`); el service rellena los días sin actividad con ceros; `accuracy` a 4 decimales y sin división por cero; `weakTopics` viene de `progress.Service.TopicMastery` filtrado por `IsWeak`; `xpGained` con los helpers de IN-1; `progress` y `gamification` se consultan en paralelo con `sync.WaitGroup` + `sync.Mutex` (**no** `errgroup`, D-13)
  - Verify: `make test`; contraste de los agregados contra la DB real con transacción con rollback
  - Files: `internal/modules/insights/store.go`, `internal/modules/insights/service.go`
- [ ] **IN-3** — Endpoint + Swagger + wiring
  - Acceptance: `GET /api/v1/insights/summary?period=week|month` (default `week`, `oneof`), registrado en `protected` y cableado en `server.go` después de `progress` y `gamification`
  - Verify: `make swagger` + `make check`; `period=year` → `422` RFC 9457; `docs/swagger.yaml` con el path y el parámetro
  - Files: `internal/modules/insights/handler.go`, `internal/modules/insights/routes.go`, `internal/server/server.go`
- [ ] **IN-4** — Tests backend
  - Acceptance: ventana 7 vs 30 días; usuario sin actividad → ceros con `dailyActivity` completo; `xpGained` = correctas×dificultad + 20×sesiones; `weakTopics` filtrado desde el mock de `progress`; error de store → `apierr`; `XpForDifficulty` (3 dificultades + desconocida → 0) y `SessionBonusXp() == 20`
  - Verify: `make test`
  - Files: `internal/modules/insights/service_test.go`, `internal/modules/gamification/service_test.go`
- [ ] **IN-5** — Shared + service + query
  - Acceptance: `insight.types.ts` (`InsightPeriod`, `TopicInsightItem`, `DailyActivityItem`, `InsightSummary`, `PERIOD_OPTIONS`) + barrel; `TopicMasteryItem` **importado** de `progress.types.ts` (no re-exportado desde dos sitios); re-export fino; `getInsightSummary(period)`; `insightSummaryOptions(period)` con el periodo en la queryKey
  - Verify: `pnpm --filter @devtest/shared typecheck` + `make fe-check`
  - Files: `packages/shared/src/types/insight.types.ts`, `packages/shared/src/index.ts`, `apps/web/src/types/insight.types.ts`, `apps/web/src/api/services/insights.service.ts`, `apps/web/src/queries/insights.queries.ts`
- [ ] **IN-6** — Página Informes + ruta
  - Acceptance: selector semana/mes; KPI cards (sesiones, preguntas, acierto, tiempo legible, XP, racha); listas "Más practicados" y "Temas débiles" con barras y badge con texto; actividad diaria con barras; **sin** librería de gráficas; ruta `/insights` con `requiresNotAdmin`
  - Verify: manual en `make fe-dev` + `make fe-check`
  - Files: `apps/web/src/features/insights/pages/InsightsPage.vue`, `apps/web/src/features/insights/components/*.vue`, `apps/web/src/router/routes.ts`
- [ ] **IN-7** — Nav (título "Estudio" + item Informes)
  - Acceptance: la sección del nav que agrupa Sesiones/Progreso pasa a titularse **Estudio** (⚠️ hoy está sin título, `DefaultLayout.vue:26-31`) y contiene `{ title: 'Informes', to: '/insights', icon: 'mdi-file-chart', userOnly: true }`
  - Verify: `make fe-check` + `make fe-test` con `DefaultLayout.test.ts` actualizado
  - Files: `apps/web/src/layouts/DefaultLayout.vue`, `apps/web/src/layouts/DefaultLayout.test.ts`
- [ ] **IN-8** — Tests web + checks
  - Acceptance: tests de service, query y página (cambio de periodo re-consulta, KPIs, listas vacías, barras); checks y suite en verde
  - Verify: `make fe-format` → `make fe-check` → `make fe-lint` → `make fe-test`
  - Files: `apps/web/src/api/services/insights.service.test.ts`, `apps/web/src/queries/insights.queries.test.ts`, test de `InsightsPage`

---

## Módulo `badges` (opt-in) — `specs/SPEC-badges.md`

> ⚠️ No empezar sin go-ahead explícito del usuario.

- [ ] **BD-1** — Catálogo `const`
  - Acceptance: `BadgeID` + `badgeCatalog` con los 8 logros (id, nombre, descripción, icono `mdi-*`, `Target`); sin lógica; sin tabla
  - Verify: `make check` + `make test`
  - Files: `internal/modules/gamification/catalog.go`
- [ ] **BD-2** — Evaluación de criterios
  - Acceptance: cada badge con su criterio (tabla de `SPEC-badges.md`), incluido `topics_mastered_5` según OQ-1; umbrales en el borde exacto; agregaciones en paralelo con `sync.WaitGroup` (**no** `errgroup`); ≤ 6 queries; `earned` y `current/target`
  - Verify: `make test` (tabla por badge, incluidos los bordes)
  - Files: `internal/modules/gamification/badges.go`, `internal/modules/gamification/store.go`
- [ ] **BD-3** — Endpoint `GET /gamification/badges`
  - Acceptance: `Service.GetBadges()` en la interface; respuesta `{ badges, earned, total }` con orden estable y **sin** `earnedAt`; ruta registrada en `Protected` existente
  - Verify: `make swagger` + `make check`; `docs/swagger.yaml` con el path
  - Files: `internal/modules/gamification/{service,response,handler,routes}.go`
- [ ] **BD-4** — Página de logros
  - Acceptance: grid ganados/bloqueados con icono + texto (nunca solo color), barras `current/target` en los bloqueados, estado vacío; ruta `/achievements` (`requiresNotAdmin`); nav "Logros" (`userOnly`) en la sección Estudio
  - Verify: `make fe-check` + `make fe-test`
  - Files: `apps/web/src/features/gamification/pages/AchievementsPage.vue`, `apps/web/src/router/routes.ts`, `apps/web/src/layouts/DefaultLayout.vue`, `packages/shared/src/types/gamification.types.ts`, `apps/web/src/api/services/gamification.service.ts`, `apps/web/src/queries/gamification.queries.ts`
- [ ] **BD-5** — Tests + checks
  - Acceptance: umbrales en el borde (99/100 respuestas, `Longest` 2/3, sesión de 4 respuestas no cuenta, media 10 001 ms no cumple); tests web de service/query/página; checks verdes
  - Verify: `make test` + `make fe-format` → `make fe-check` → `make fe-lint` → `make fe-test`
  - Files: `internal/modules/gamification/badges_test.go`, `apps/web/src/features/gamification/pages/AchievementsPage.test.ts`

---

## Módulo `e2e` (opt-in) — `specs/SPEC-e2e.md`

> ⚠️ No empezar sin go-ahead explícito del usuario. Requiere el gate de `.agents/web/dependencies.md`.

- [ ] **E2-1** — Suite Playwright del camino crítico
  - Acceptance: `@playwright/test` como devDependency (gate de dependencias documentado); `playwright.config.ts` con `workers: 1`, `webServer` para Vite y `baseURL :3000`; `critical-flow.spec.ts` recorre `status → setup → login → dashboard → sessions`; sin `waitForTimeout`; selectores por rol/texto; targets `make e2e` / `e2e-ui` / `e2e-report`; `e2e/.auth/` gitignored; el E2E **no** entra en `make fe-test`
  - Verify: `make e2e` verde desde una DB limpia (`DB_NAME=devtest_e2e make migrate-up` + `make run`); falla con mensaje claro sin backend o con DB ya configurada
  - Files: `apps/web/playwright.config.ts`, `apps/web/e2e/critical-flow.spec.ts`, `apps/web/e2e/fixtures.ts`, `apps/web/package.json`, `Makefile`, `.gitignore`

---

## Verificación final (re-run de la Fase 8 del `ROADMAP`)

- [ ] **VF-1** — Backend: `make fmt` + `make vet` + `make test` (+ `make build`)
- [ ] **VF-2** — Frontend: `make fe-format` + `make fe-check` + `make fe-lint` + `make fe-test`
- [ ] **VF-3** — Tipos compartidos y mobile: `pnpm -r typecheck`
- [ ] **VF-4** — Esquema: `make migrate-status` sin migraciones pendientes; `docs/swagger.yaml` regenerado y actualizado
