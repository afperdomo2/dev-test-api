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
- [x] **RM-8** — Tests web + checks
  - Acceptance: tests de service, query y página (render, badge Débil, vacío); los 3 checks y la suite en verde
  - Verify: `make fe-format` → `make fe-check` → `make fe-lint` → `make fe-test`
  - Files: `apps/web/src/api/services/progress.service.test.ts`, `apps/web/src/queries/progress.queries.test.ts`, test de `ProgressPage`
  - Hecho: **+12 tests** (1 service, 1 query, **6 de componente** en `TopicMasteryList.test.ts` — nuevo, no estaba previsto — y 4 de página) → suite **303 tests / 51 archivos** · `fe-format`/`fe-check`/`fe-lint`/`fe-test` = 0 · `shared typecheck` = 0
  - Notas de implementación: (a) el `v-card` de Vuetify pinta un **loader interno** con `role="progressbar"` sin `aria-label`, así que los tests seleccionan `[role="progressbar"][aria-label]` para no confundirlo con las barras reales; (b) los tests respetan `noUncheckedIndexedAccess` (accesos con `?.`) y la regla que prohíbe `!` (helper `tabByLabel` que lanza si no encuentra la pestaña)

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

- [x] **IN-1** — Helpers de XP en `gamification` + esqueleto del módulo
  - Acceptance: `gamification.XpForDifficulty(difficulty) int` y `gamification.SessionBonusXp() int` exportados y usados internamente por `gamification` (una sola fuente del XP, N-01); módulo `insights` creado con la anatomía estándar
  - Verify: `make check` + `make test` (los tests existentes de gamification siguen verdes)
  - Files: `internal/modules/gamification/service.go`, `internal/modules/insights/{request,response,store,service,handler,routes}.go`
  - Hecho: helpers exportados + `xp_test.go` (5/10/15, desconocida → 0, bonus 20)
- [x] **IN-2** — Agregaciones del store + `Summary`
  - Acceptance: `CountSessions`, `CountAnswers` (total y correctas), `SumResponseTimeMs`, `TopTopics(limit)`, `DailyActivity`, `CountCorrectByDifficulty`, `CountCompletedSessions`; filtran por `created_at` en la ventana y sesión `completed` (`models.SessionStatusCompleted`); el service rellena los días sin actividad con ceros; `accuracy` a 4 decimales y sin división por cero; `weakTopics` viene de `progress.Service.TopicMastery` filtrado por `IsWeak`; `xpGained` con los helpers de IN-1; `progress` y `gamification` se consultan en paralelo con `sync.WaitGroup` + `sync.Mutex` (**no** `errgroup`, D-13)
  - Verify: `make test`; contraste de los agregados contra la DB real con transacción con rollback
  - Files: `internal/modules/insights/store.go`, `internal/modules/insights/service.go`
  - Hecho: **8 consultas en paralelo** (`WaitGroup` + mutex solo para `errs`; cada goroutine escribe su propia variable) · `CountCompletedSessions` **no se implementó**: `CountSessions` ya devuelve el mismo dato y se usa para el bonus de XP · ventana y días en **UTC** (determinista, sin depender de la TZ de Postgres) · `topTopics` desempata por accuracy en SQL
- [x] **IN-3** — Endpoint + Swagger + wiring
  - Acceptance: `GET /api/v1/insights/summary?period=week|month` (default `week`, `oneof`), registrado en `protected` y cableado en `server.go` después de `progress` y `gamification`
  - Verify: `make swagger` + `make check`; `period=year` → `422` RFC 9457; `docs/swagger.yaml` con el path y el parámetro
  - Files: `internal/modules/insights/handler.go`, `internal/modules/insights/routes.go`, `internal/server/server.go`
  - Hecho: handler con `getUserID` del módulo, `e.Instance` y `response.Problem` · `make swagger` OK · wiring tras `progress`
- [x] **IN-4** — Tests backend
  - Acceptance: ventana 7 vs 30 días; usuario sin actividad → ceros con `dailyActivity` completo; `xpGained` = correctas×dificultad + 20×sesiones; `weakTopics` filtrado desde el mock de `progress`; error de store → `apierr`; `XpForDifficulty` (3 dificultades + desconocida → 0) y `SessionBonusXp() == 20`
  - Verify: `make test`
  - Files: `internal/modules/insights/service_test.go`, `internal/modules/gamification/service_test.go`
  - Hecho: 8 subtests de service (ventana, accuracy 0.6667, XP 75, weak, zero-fill, racha, gamification caída no rompe, error de store) + 4 de handler (200/422/401/500) + helpers de XP
- [x] **IN-5** — Shared + service + query
  - Acceptance: `insight.types.ts` (`InsightPeriod`, `TopicInsightItem`, `DailyActivityItem`, `InsightSummary`, `PERIOD_OPTIONS`) + barrel; `TopicMasteryItem` **importado** de `progress.types.ts` (no re-exportado desde dos sitios); re-export fino; `getInsightSummary(period)`; `insightSummaryOptions(period)` con el periodo en la queryKey
  - Verify: `pnpm --filter @devtest/shared typecheck` + `make fe-check`
  - Files: `packages/shared/src/types/insight.types.ts`, `packages/shared/src/index.ts`, `apps/web/src/types/insight.types.ts`, `apps/web/src/api/services/insights.service.ts`, `apps/web/src/queries/insights.queries.ts`
  - Hecho: `TopicMasteryItem` se importa (no se re-exporta) · barrel con `export *`
- [x] **IN-6** — Página Informes + ruta
  - Acceptance: selector semana/mes; KPI cards (sesiones, preguntas, acierto, tiempo legible, XP, racha); listas "Más practicados" y "Temas débiles" con barras y badge con texto; actividad diaria con barras; **sin** librería de gráficas; ruta `/insights` con `requiresNotAdmin`
  - Verify: manual en `make fe-dev` + `make fe-check`
  - Files: `apps/web/src/features/insights/pages/InsightsPage.vue`, `apps/web/src/features/insights/components/*.vue`, `apps/web/src/router/routes.ts`
  - Hecho: una sola página (sin componentes extra) con `v-btn-toggle`, KPIs, barras y actividad diaria en divs · sin librería de gráficas · ruta `/insights` (`requiresNotAdmin`)
- [x] **IN-7** — Nav (título "Estudio" + item Informes)
  - Acceptance: la sección del nav que agrupa Sesiones/Progreso pasa a titularse **Estudio** (⚠️ hoy está sin título, `DefaultLayout.vue:26-31`) y contiene `{ title: 'Informes', to: '/insights', icon: 'mdi-file-chart', userOnly: true }`
  - Verify: `make fe-check` + `make fe-test` con `DefaultLayout.test.ts` actualizado
  - Files: `apps/web/src/layouts/DefaultLayout.vue`, `apps/web/src/layouts/DefaultLayout.test.ts`
  - Hecho: sección **Estudio** con Sesiones, Progreso e Informes · tests actualizados (Informes visible para usuario, oculto para admin; orden de la sección)
- [x] **IN-8** — Tests web + checks
  - Acceptance: tests de service, query y página (cambio de periodo re-consulta, KPIs, listas vacías, barras); checks y suite en verde
  - Verify: `make fe-format` → `make fe-check` → `make fe-lint` → `make fe-test`
  - Files: `apps/web/src/api/services/insights.service.test.ts`, `apps/web/src/queries/insights.queries.test.ts`, test de `InsightsPage`
  - Hecho: 8 tests nuevos (service 1, query 1, página 6) → suite **311 tests / 54 archivos** · `fe-format`/`fe-check`/`fe-lint`/`fe-test` = 0 · backend `make check` + `make test` = 0

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

## Módulo `push` — `specs/SPEC-push.md`

- [ ] **PUSH-1** — Modelo + migración `00007_device_tokens.sql`
  - Acceptance: tabla con `token` único global, `platform` (`ios|android`), `last_seen_at`, FK a `users` con `ON DELETE CASCADE`, `COMMENT ON` de tabla y columnas; `Down` con `DROP TABLE`
  - Verify: `make migrate-up` + `make migrate-status` + `make migrate-down`; insertar dos veces el mismo token → falla el segundo (índice único)
  - Files: `internal/database/migrations/00007_device_tokens.sql`, `internal/models/device_token.go`
- [ ] **PUSH-2** — Store + service (registro y baja)
  - Acceptance: `Store` con upsert por token, delete por token, listado por usuario; el upsert **reasigna** el token si ya existía con otro usuario; `Unregister` ajeno no borra; errores a `apierr`
  - Verify: `make check` + `make test`
  - Files: `internal/modules/push/{request,response,store,service}.go`
- [ ] **PUSH-3** — Sender de Expo + `Pusher` + hook en notifications
  - Acceptance: `ExpoSender` hace POST a `EXPO_PUSH_URL` en lotes de 100 con timeout; `DeviceNotRegistered` borra el token; `Pusher` implementa la interfaz de `notifications`; `RunDailyJob` llama a `Push` **después** de insertar y no falla si Expo falla
  - Verify: `make test` (fake de Expo con `httptest`) + prueba real con un token falso → log de error sin romper el job
  - Files: `internal/modules/push/sender.go`, `internal/modules/push/service.go`, `internal/modules/notifications/service.go`, `internal/config/config.go`
- [ ] **PUSH-4** — Handler + rutas + wiring + Swagger
  - Acceptance: `POST /push/tokens` y `POST /push/tokens/unregister` (protegidos, idempotentes); wiring en `server.go` con `push` **antes** de `notifications`
  - Verify: `make swagger` + `make check`; arranque real sin conflictos; `docs/swagger.yaml` con los 2 paths
  - Files: `internal/modules/push/{handler,routes}.go`, `internal/server/server.go`
- [ ] **PUSH-5** — Tests backend
  - Acceptance: sender (lotes, timeout, `DeviceNotRegistered`), service (upsert, reasignación, baja), handler (401/422/200); store real en transacción con rollback
  - Verify: `make test`
  - Files: `internal/modules/push/{sender,service,handler}_test.go`

---

## Módulo `mobile-tooling` — `specs/SPEC-mobile-tooling.md`

- [x] **MTOOL-1** — ESLint en mobile
  - Acceptance: `eslint.config.js` con `eslint-config-expo`; script `lint` real (deja de ser un `echo`); gate de `.agents/mobile/dependencies.md` documentado
  - Verify: `pnpm --filter mobile lint` ✅ sin errores
  - Files: `apps/mobile/eslint.config.js`, `apps/mobile/package.json`
  - Hecho: `eslint@9` + `eslint-config-expo@58` (MIT, peer `eslint >=8.10`, gate OK) · **0 errores, 11 warnings** · la regla nueva `react-hooks/set-state-in-effect` se dejó en `warn` (marca el patrón estándar de fetch inicial de `SessionsListPage`; como error bloquearía código válido) · se **eliminó un `useEffect` redundante** en `QuestionCard` (el padre ya monta con `key={question.id}`) · `types/eslint-config-expo.d.ts` para el typecheck
- [x] **MTOOL-2** — Configuración de Jest + mocks nativos
  - Acceptance: `jest.config.js` con `jest-expo`, `setupFilesAfterEnv`; `jest.setup.ts` con mocks de `expo-secure-store`, `expo-router`, `@expo/vector-icons`; matchers de RTL
  - Verify: `pnpm --filter mobile test:run` (los 2 tests existentes siguen verdes)
  - Files: `apps/mobile/jest.config.js`, `apps/mobile/jest.setup.ts`, `apps/mobile/package.json`
  - Hecho: Jest no resolvía el workspace → `moduleNameMapper` para `@devtest/shared` (main es TS) y `@/` · `testTimeout: 30000` (el primer render supera los 5 s por defecto) · el bloque `"jest"` de `package.json` se movió al config
- [x] **MTOOL-3** — Tests del cliente HTTP y servicios
  - Acceptance: `client.ts` (9 casos de la spec, incluido `Network Error` y 401) + los 5 servicios (URL/método/payload) + `auth/storage`
  - Verify: `pnpm --filter mobile test:run`; romper `unwrap()` a propósito → falla un test
  - Files: `apps/mobile/src/core/api/client.test.ts`, `apps/mobile/src/core/api/services/*.test.ts`, `apps/mobile/src/core/auth/storage.test.ts`
  - Hecho: **44 tests** (cliente 9, storage 4, auth 4, sessions 6, questions 3, progress 3, gamification 3, + los 2 de utils)
- [x] **MTOOL-4** — Tests de pantallas críticas
  - Acceptance: `HomeScreen` (racha, meta, CTA diaria, freeze con confirmación), `QuestionCard` (selección única/múltiple, deshabilitado), `SessionStudyPage` (responder y avanzar), `LoginScreen`/`SetupScreen`
  - Verify: `pnpm --filter mobile test:run` + `typecheck` + `lint`; romper `handleAnswer` a propósito → falla un test
  - Files: `apps/mobile/src/features/**/*.test.tsx`
  - Hecho: **20 tests** → total **64 tests / 14 archivos** · se añadió `accessibilityLabel` a `AuthButton` (faltaba) y `testID` a `AuthInput` (RTL v14 ya no expone `UNSAFE_getAllByType`) · los tests de `Login/Setup` validan formularios sin tocar el backend

---

## Módulo `mobile-daily-tip` — `specs/SPEC-mobile-daily-tip.md`

- [x] **MTIP-1** — Store en memoria del tip
  - Acceptance: `setPendingTip`, `getPendingTip` (no consume), `clearPendingTip`; otra sesión devuelve `null`
  - Verify: `pnpm --filter mobile test:run`
  - Files: `apps/mobile/src/features/sessions/dailyTipStore.ts`, `dailyTipStore.test.ts`
  - Hecho: 5 tests (guardar/leer, otra sesión, no consume, sin tip, clear)
- [x] **MTIP-2** — `DailyTipCard`
  - Acceptance: título, tema, contenido con `MarkdownText`, cerrar con `accessibilityLabel`, solo tokens de diseño
  - Verify: `pnpm --filter mobile test:run` + `lint`
  - Files: `apps/mobile/src/features/sessions/components/DailyTipCard.tsx`, `.test.tsx`
  - Hecho: botón de cerrar de 44×44 pt con `accessibilityLabel="Cerrar tip del día"` · icono + color del tema (`colors.accent`), sin hex
- [x] **MTIP-3** — Integrar en Home y en la sesión diaria
  - Acceptance: `handleDailySession` guarda `res.tip`; la pantalla de sesión pinta la tarjeta solo en fase `question`; se limpia al salir
  - Verify: tests de `HomeScreen` y `SessionStudyPage` + prueba manual con sesión diaria
  - Files: `apps/mobile/src/features/home/screens/HomeScreen.tsx`, `apps/mobile/src/features/sessions/screens/SessionStudyPage.tsx` (+ tests)
  - Hecho: `setPendingTip` antes de navegar · la tarjeta solo se pinta en la **primera** pregunta (`questionNumber === 1`), no en `result`/`finished`/`error` · `clearPendingTip()` al desmontar · tests: se guarda en Home (1), se muestra/oculta/cierra/no aparece sin tip (4)

---

## Módulo `mobile-progress` — `specs/SPEC-mobile-progress.md`

- [x] **MPROG-1** — Servicio `getTopicMastery`
  - Acceptance: llama `GET /progress/topics-mastery` y devuelve `TopicMasteryItem[]` (tipo de shared)
  - Verify: `pnpm --filter mobile test:run` + `typecheck`
  - Files: `apps/mobile/src/core/api/services/progress.service.ts`, `progress.service.test.ts`
  - Hecho: función añadida + test (72 tests en total tras el módulo)
- [x] **MPROG-2** — 4ª pestaña + hub
  - Acceptance: `NativeTabs` con 4 triggers (Inicio, Sesiones, Progreso, Perfil); `progreso/_layout.tsx` con Stack; hub con la fila **Temas** (las otras dos llegan en sus olas)
  - Verify: prueba manual en dispositivo (iOS/Android): las 3 pestañas previas siguen funcionando
  - Files: `apps/mobile/app/(tabs)/_layout.tsx`, `apps/mobile/app/(tabs)/progreso/{_layout,index}.tsx`, `apps/mobile/src/features/progress/screens/ProgressHubScreen.tsx`
  - Hecho: 4 triggers · Stack con headers nativos por pantalla (`index` = "Progreso", `temas` = "Temas", así hay botón atrás) · hub con `ProgressHubRow` (touch target 64 pt, `accessibilityRole="button"`)
- [x] **MPROG-3** — Pantalla de Temas
  - Acceptance: lista con barra (`View` + %), `correctCount/attempts`, chip "Débil" con texto, `accessibilityLabel` con el porcentaje, estados vacío/error/carga, pull-to-refresh
  - Verify: tests de pantalla + `lint` + prueba manual contra el backend real
  - Files: `apps/mobile/app/(tabs)/progreso/temas.tsx`, `apps/mobile/src/features/progress/screens/TopicMasteryScreen.tsx`, `components/{TopicMasteryRow,ProgressHubRow}.tsx` (+ tests)
  - Hecho: barra con `accessibilityRole="progressbar"` + `accessibilityValue {min,max,now}` (en vez de `aria-*`) · chip "Débil" con icono + texto · estado vacío, error con "Reintentar" y `RefreshControl` · tests de hub, fila y pantalla (incluido el reintento)
- [ ] **MPROG-4** — Verificación en dispositivo
  - Acceptance: 4 pestañas correctas, navegación hub→Temas con back, datos reales coinciden con la web, touch targets ≥ 44 pt
  - Verify: manual en dev build/dispositivo; `typecheck` + `test:run` + `lint`
  - Files: — (verificación)
  - ⚠️ **Pendiente**: requiere dispositivo/emulador. El código y los tests están verdes; falta la pasada visual (barra inferior con 4 tabs y header nativo con atrás).

---

## Módulo `mobile-weekly` — `specs/SPEC-mobile-weekly.md`

- [x] **MWEEK-1** — Servicio de insights
  - Acceptance: `getInsightSummary(period)` con `params.period`; tipos de shared
  - Verify: `pnpm --filter mobile test:run` + `typecheck`
  - Files: `apps/mobile/src/core/api/services/insights.service.ts`, `.test.ts`
  - Hecho: función + test (URL y `period`)
- [x] **MWEEK-2** — Pantalla "Mi semana" + fila en el hub
  - Acceptance: selector Semana/Mes, 6 KPIs, más practicados, temas débiles (o sección oculta), actividad diaria en barras, estados vacío/error
  - Verify: `typecheck` + `lint` + prueba manual con datos reales
  - Files: `apps/mobile/app/(tabs)/progreso/semana.tsx`, `features/progress/screens/WeeklySummaryScreen.tsx`, `components/{KpiCard,PeriodToggle,DailyActivityBar}.tsx`, `ProgressHubScreen.tsx`
  - Hecho: ruta `semana` con header "Mi semana" · fila nueva en el hub · `TopicMasteryRow` se **generalizó** (recibe campos, no el DTO) para reutilizarla en informes sin `correctCount` ni chip · barras diarias con altura proporcional y `accessibilityLabel` por día · estados vacío/error/reintento
- [x] **MWEEK-3** — Tests de la pantalla
  - Acceptance: pinta KPIs, cambia de periodo y re-consulta, vacío, error con reintento, oculta temas débiles vacíos; `DailyActivityBar` con altura y a11y correctas
  - Verify: `pnpm --filter mobile test:run`
  - Files: `apps/mobile/src/features/progress/screens/WeeklySummaryScreen.test.tsx`, `components/*.test.tsx`
  - Hecho: **11 tests** nuevos (pantalla 4, toggle 1, barra 3, hub +2, servicio 1) → suite mobile **95 tests / 23 archivos** · `typecheck` + `lint` + `pnpm -r typecheck` en verde

---

## Módulo `mobile-badges` — `specs/SPEC-mobile-badges.md`

- [ ] **MBADGE-1** — Servicio `getBadges`
  - Acceptance: llama `GET /gamification/badges`; tipos `Badge`/`BadgesResponse` de shared
  - Verify: `pnpm --filter mobile test:run` + `typecheck`
  - Files: `apps/mobile/src/core/api/services/gamification.service.ts`, `.test.ts`
- [ ] **MBADGE-2** — Pantalla de Logros + fila en el hub
  - Acceptance: contador `earned/total`, cuadrícula de 2 columnas, conseguido = icono + "Conseguido", pendiente = `current/target` + barra, sin `earnedAt`, estados vacío/error
  - Verify: `typecheck` + `lint` + prueba manual
  - Files: `apps/mobile/app/(tabs)/progreso/logros.tsx`, `features/progress/screens/AchievementsScreen.tsx`, `components/BadgeCard.tsx`, `ProgressHubScreen.tsx`
- [ ] **MBADGE-3** — Tests
  - Acceptance: 8 badges, contador, ambos estados, a11y del `BadgeCard`, mapeo de icono con fallback
  - Verify: `pnpm --filter mobile test:run`
  - Files: `apps/mobile/src/features/progress/screens/AchievementsScreen.test.tsx`, `components/BadgeCard.test.tsx`

---

## Módulo `mobile-push` — `specs/SPEC-mobile-push.md`

- [ ] **MPUSH-1** — Dependencias + `app.json` + canal Android
  - Acceptance: `expo-notifications`, `expo-device`, `expo-constants` instalados (gate documentado); plugin en `app.json`; `extra.eas.projectId` añadido; canal `default` creado antes del token
  - Verify: `pnpm --filter mobile typecheck`; comprobar que la app arranca en dev build
  - Files: `apps/mobile/app.json`, `apps/mobile/package.json`
- [ ] **MPUSH-2** — Módulo de push + servicio HTTP + tests
  - Acceptance: `ensurePermission`, `getDevicePushToken` (null en emulador/sin projectId), `configureNotificationHandling`; `registerForPush`/`unregisterFromPush` que **nunca lanzan**; `push.service.ts` con register/unregister
  - Verify: `pnpm --filter mobile test:run` (mocks de `expo-notifications`/`expo-device`/`expo-constants`)
  - Files: `apps/mobile/src/core/push/{notifications,register}.ts` (+ tests), `apps/mobile/src/core/api/services/push.service.ts` (+ test)
- [ ] **MPUSH-3** — Wiring + navegación al tocar + prueba real
  - Acceptance: registro una vez tras el bootstrap autenticado; baja en logout; tocar una notificación (y arranque en frío) abre Inicio
  - Verify: prueba en **dev build** con un envío real desde el backend (`push`); sin permiso la app funciona igual
  - Files: bootstrap de auth, pantalla de Perfil/logout, `app/_layout.tsx`

---

## Módulo `mobile-offline` — `specs/SPEC-mobile-offline.md`

- [ ] **MOFF-1** — Backend: migración + modelo + request
  - Acceptance: `session_answers.client_answer_id` + índice único `(session_id, client_answer_id)` con `COMMENT ON`; `ClientAnswerID *uuid.UUID` en `AnswerRequest` y en el modelo (opcional)
  - Verify: `make migrate-up`/`down`; insertar dos filas con el mismo par → falla; con `NULL` → permitido
  - Files: `internal/database/migrations/NNNNN_answer_idempotency.sql`, `internal/models/session_answer.go`, `internal/modules/sessions/request.go`
- [ ] **MOFF-2** — Backend: store + service idempotentes
  - Acceptance: `CreateAnswer` devuelve `created bool` con `ON CONFLICT DO NOTHING`; si la respuesta ya existía, el service devuelve la guardada y **no** ejecuta `progressService.Answer` (sin doble XP/SRS); sin `clientAnswerId` se comporta como hoy
  - Verify: `make test` + `make check`; contraste real en transacción con rollback (dos envíos → 1 fila, 1 XP)
  - Files: `internal/modules/sessions/store.go`, `internal/modules/sessions/service.go`, `internal/modules/sessions/service_test.go`
- [ ] **MOFF-3** — Mobile: `isNetworkError` + `clientAnswerId`
  - Acceptance: helper en `core/utils/errors.ts`; `submitAnswer` acepta y envía `clientAnswerId`; se genera **una vez por pregunta** con `expo-crypto` y se reutiliza al reintentar
  - Verify: `pnpm --filter mobile test:run` + `typecheck`; gate de dependencias de `expo-crypto`
  - Files: `apps/mobile/src/core/utils/errors.ts` (+ test), `apps/mobile/src/core/api/services/sessions.service.ts` (+ test)
- [ ] **MOFF-4** — Mobile: UI de reintento
  - Acceptance: error de red → no avanza, conserva la selección, mensaje + "Reintentar"; error de API → detalle + "Reintentar"; doble toque no duplica; el `clientAnswerId` cambia al cambiar de pregunta
  - Verify: tests de `SessionStudyPage` (espía del `clientAnswerId`) + prueba manual en modo avión
  - Files: `apps/mobile/src/features/sessions/screens/SessionStudyPage.tsx`, `.test.tsx`

---

## Verificación final

- [ ] **VF-1** — Backend: `make fmt` + `make vet` + `make test` (+ `make build`)
- [ ] **VF-2** — Web: `make fe-format` + `make fe-check` + `make fe-lint` + `make fe-test`
- [ ] **VF-3** — Mobile: `pnpm --filter mobile typecheck` + `pnpm --filter mobile test:run` + `pnpm --filter mobile lint`
- [ ] **VF-4** — Shared: `pnpm --filter @devtest/shared typecheck`
- [ ] **VF-5** — Esquema: `make migrate-status` sin migraciones pendientes; `docs/swagger.yaml` regenerado y actualizado
