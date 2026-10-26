# Decisiones — log vivo

> Rescatado de `ROADMAP.md` (sección "Decisiones bloqueadas" + "Notas") al migrar el backlog a SDD. **Estas decisiones no son tareas**: son guardarraíles que restringen cualquier implementación futura.
>
> Reglas: una decisión no se borra nunca; si cambia, se marca `SUPERSEDED por D-xx`. Si una decisión contradice a `.agents/*.md`, gana este archivo y hay que corregir el doc (ver D-13).

## Cómo leer las etiquetas

- **Bloqueada** — no se discute; si algo la viola, el cambio está mal.
- **SUPERSEDED** — reemplazada por otra decisión posterior.
- **Abierta** — pendiente de decisión; su spec tiene la Open Question.

---

## Decisiones de producto / alcance

| ID | Estado | Decisión |
|---|---|---|
| D-01 | Bloqueada | **V/F**: `Question.Type = "true_false"` (2 opciones, 1 correcta). Reusa el pipeline create/import/answer. |
| D-02 | Bloqueada | **Tips ↔ temas**: un solo tema por tip (`Tip.TopicID`, FK indexada). |
| D-03 | Bloqueada | **Creación de tips**: solo admin y siempre públicos. Backend solo manual + import CSV; **nunca** genera tips con IA. El prompt helper de `/questions/import` es texto para copiar/pegar en una IA externa. Usuario normal: solo lectura. |
| D-04 | Bloqueada | **Visibilidad sin cambios**: se mantiene `is_public` + `source`. Solo IA y admin son públicos; el usuario normal queda privado. **Sin** `pending_review`, **sin** acción "publicar", **sin** refactor `isPublic → visibility`. |
| D-05 | Bloqueada | **Moderación**: solo reportar contenido público + revisión admin. Sin "publish". |
| D-06 | Bloqueada | **Import UI (web)**: página unificada con selector "preguntas / V-F / tips"; tips solo visible para admin. |
| D-07 | Bloqueada | **Gamificación**: un solo módulo `internal/modules/gamification/`. **Sin tabla `UserStats`**: se usa `UserStreak` + `UserXp` + `UserTopicXp`. XP por dificultad `beginner=5 / intermediate=10 / advanced=15` + `sessionBonusXp=20` al terminar sesión (`gamification/service.go:13-24`). |
| D-08 | Bloqueada | **Definición única de "tema débil"**: vive en `review-mastery`. `isWeak = accuracy < 0.6 && attempts >= 5` (constantes en Go). Ningún otro módulo la recalcula (`insights` filtra, `badges` consume). |
| D-09 | Bloqueada | **Definición única de "due review"**: `sessions.Store.CountDueReviews` (`sessions/store.go:210`) — `next_review_at <= now() AND is_mastered = false`. No se reimplementa en `progress` ni en `notifications`. |
| D-10 | Bloqueada | **Accuracy**: respuestas (`SessionAnswer`) de **sesiones `completed`**, contando todos los reintentos. `attempts` = respuestas totales al tema (no preguntas distintas). ⚠️ El estado correcto es `completed`, **no** `finished` (el `ROADMAP.md` decía "finished"): `SessionStatus = 'in_progress' \| 'completed' \| 'cancelled'` (`packages/shared/src/types/session.types.ts:5`). |
| D-11 | Bloqueada | **Sesión diaria**: modo `daily` en `Session` + selección on-the-fly, **sin** tabla de cola persistente. Es **mobile-first**: la web **no** tiene botón "Repasar hoy". |
| D-12 | Bloqueada | **Insights** (no "reports"): `GET /insights/summary`, para no chocar con `/reports` (moderación). |
| D-13 | Bloqueada | **Concurrencia**: goroutines + `sync.WaitGroup`. **Sin `errgroup`** (no añadir dependencia). El patrón para trabajo de fondo es `internal/modules/generation/worker.go` (pool + `context` cancelable + `WaitGroup` + reconciliación); para llamadas independientes dentro de un service, `sync.WaitGroup` + `sync.Mutex` (patrón real: `sessions/service.go` → `gatherDailyData()`). Nota: `golang.org/x/sync` figura como **indirect** en `go.mod` (arrastrado por otra dependencia) — eso **no** autoriza usarlo en código propio. ✅ El doc `.agents/backend/architecture.md` quedó corregido el 2026-10-17: documentaba `errgroup` como canónico y afirmaba en falso que ya se usaba en `sessions/service.go` (N-07). |
| D-14 | Bloqueada | **IA**: solo por lotes y solo admins (`/generation`), y solo genera preguntas `single_choice`. Tope **global** diario `ai.max_per_day`. El modo `generate` de sesiones y la cuota de IA por usuario **se retiraron** (fases 10 y 11). **No reintroducir IA por usuario en sesiones.** |
| D-15 | Bloqueada | **Cuotas**: `ImportQuota` / `DailyImportLimit` (import CSV, con exención de admin) **se conserva**. La cuota de IA por usuario **no existe**. |
| D-16 | Bloqueada | **Logros**: catálogo en código (`const`), **sin tabla** `UserBadge` ni `user_badges`. Se derivan on-the-fly. Sin `earnedAt` (no hay persistencia). |
| D-17 | Bloqueada | **Scheduler de notificaciones**: ticker stdlib (`time.Timer` calculando el próximo corte) + goroutine. **Sin `robfig/cron`**. |
| D-18 | SUPERSEDED por D-26 | ~~**Alcance de clientes del backlog**: backend + web. Mobile fuera de alcance.~~ |
| D-19 | Bloqueada | **Responses sin paginar** en `topics-mastery` (el set de temas es acotado). El resto de listas nuevos sí pagina. |
| D-26 | Bloqueada | **Mobile entra en alcance** (reemplaza D-18). Se abre el backlog mobile con 8 módulos: `push`, `mobile-tooling`, `mobile-daily-tip`, `mobile-progress`, `mobile-weekly`, `mobile-badges`, `mobile-offline`, `mobile-push`. Ver `specs/MAP.md`. |
| D-27 | Bloqueada | **Push con Expo Push**: `expo-notifications` en el cliente + Expo Push API en el backend. **Sin** Firebase ni APNs directos (ni `google-services.json` ni claves `.p8`). |
| D-28 | Bloqueada | **Un solo dominio de avisos**: `notifications` sigue siendo la fuente (tabla + endpoints + job diario). `push` es **solo entrega**: guarda tokens y empuja lo que `notifications` crea, vía una interfaz `Pusher`. Por eso `push` **depende de** `notifications`. |
| D-29 | Bloqueada | **Logros en las 3 capas**: backend (`catalog.go` + criterios + endpoint), web (`/achievements`) y mobile (dentro de Progreso). Deja de ser opt-in. |
| D-30 | Bloqueada | **El tip del día se muestra al empezar la sesión diaria**. El backend solo lo entrega en `POST /sessions/daily` (no existe `GET` de tip), así que no se pide aparte ni se muestra en Home. |
| D-31 | Bloqueada | **Offline = reintento con aviso**: si el envío de una respuesta falla por red, se conserva la selección, se avisa y se ofrece "Reintentar". **Sin** cola persistente en disco (ver X-11). |
| D-32 | Bloqueada | **El E2E de Playwright se mantiene** en el plan (un flujo, opt-in ya aprobado). |
| D-33 | Bloqueada | **Navegación mobile**: `NativeTabs` pasa a **4 pestañas** — Inicio, Sesiones, **Progreso**, Perfil. Dentro de Progreso viven **Temas**, **Mi semana** y **Logros**. |
| D-34 | Bloqueada | **`mobile-tooling` incluye tests retroactivos**: además de lint, se cubren con tests el código **ya existente** (cliente HTTP, servicios y pantallas críticas), no solo lo nuevo. Objetivo: que los cambios futuros rompan menos. |
| D-35 | Bloqueada | **Idempotencia de respuestas**: `POST /sessions/{id}/answer` acepta un `clientAnswerId` opcional (UUID generado por el cliente) con índice único `(session_id, client_answer_id)`. Reintentar por falta de red —o un doble toque— devuelve la respuesta ya guardada **sin** duplicar la fila, el XP ni el SRS. Sin el campo, el comportamiento es el actual (retrocompatible). Requisito de `mobile-offline`. |
| D-36 | Bloqueada | **Priorización: notificaciones y push al final.** `notifications` → `push` → `mobile-push` se ejecutan después del resto del backlog (Temas, tip, insights/weekly, badges, offline, e2e). Consecuencia asumida: el aviso diario en el teléfono llega más tarde; nada más depende de ellos. |
| D-37 | Bloqueada | **Autorización por objeto = 404.** Cuando un recurso no existe **o** pertenece a otro usuario, el service devuelve `ErrNotFound` (nunca 403 ni detalles distintos). Evita el oráculo de existencia. `DELETE /sessions/{id}` cambia su 403 de ownership a 404; el 403 de la regla de 24 h se mantiene. Ver `specs/SPEC-authz-objects.md` (D-AZ-1, D-AZ-3). |
| D-38 | Bloqueada | **Setup con token de entorno.** `POST /auth/setup` exige `SETUP_TOKEN` comparado en tiempo constante (SHA-256); si la variable está vacía el endpoint responde 403 en **todos** los modos. El arranque **falla** si la DB está vacía y falta el token (guard en `server.Run`) y exige **≥ 16 caracteres** cuando está informado; con la DB ya inicializada solo avisa. `GET /auth/status` expone `setupRequiresToken` para que web y mobile muestren u oculten el campo. Ver `specs/SPEC-auth-hardening.md`. |
| D-39 | Bloqueada | **Rate limit en memoria por IP.** Token bucket (`golang.org/x/time/rate`) con `map` + mutex y barrido perezoso; 10/min con burst 5 en `/auth/setup` y `/auth/login`; `/auth/status` exento. Válido solo para **una instancia**: si el despliegue escala horizontalmente habrá que moverlo a Redis (nueva spec). |
| D-40 | Bloqueada | **Vigilancia de dependencias.** `govulncheck` y `pnpm audit --prod --audit-level high` son gates de CI. Solo se permiten `pnpm.overrides` de transitivas, nunca `pnpm up` global (rompe el pinning de Expo). Cada `ignoreCves` exige motivo y fecha de revisión en `SECURITY.md`. |
| D-41 | Bloqueada | **HTTPS obligatorio en mobile de producción.** `EXPO_PUBLIC_API_BASE_URL` debe empezar por `https://` cuando `__DEV__ === false`; si no, el módulo de config lanza al importarse y la app no envía ninguna petición. `http://localhost` sigue permitido en desarrollo. |
| D-42 | Bloqueada | **El JWT se queda en `localStorage` (web).** Migrar a cookie `httpOnly` exige CSRF y cambios de contrato: fuera de alcance. Mitigaciones obligatorias: CSP estricta en el borde, `removeToken()` en 401 y TTL recomendada de 8–12 h. Un XSS sigue siendo el riesgo residual aceptado. |

## Convenciones del repo rescatadas (no estaban en `.agents/`)

| ID | Convención |
|---|---|
| D-20 | `Question.Type` es **`string`**: los `CREATE TYPE ... AS ENUM` de las migraciones son **decorativos**. La validación real es Go `binding:"oneof=..."`. No "arreglar" la DB para usar los enums. |
| D-21 | Enums de dominio se modelan como `string` + constante Go + `CHECK`/`binding` explícito. |
| D-22 | Logs con emoji (`❌` error, `✅` éxito, `🚀` arranque, `🛢️` DB, `🌱` seed, `🤖` IA). |
| D-23 | Todo objeto de esquema (tabla, columna, enum, función, trigger) lleva su `COMMENT ON` **en la migración** que lo crea. Nunca se edita una migración ya aplicada. |
| D-24 | Cambio de esquema = migración goose nueva (`make migrate-create name=<desc>`). `Up` idempotente cuando sea seguro; `Down` reversible o documentado como destructivo. |
| D-25 | Dependencia nueva (Go o pnpm) → pasar por `.agents/backend/dependencies.md` o `.agents/web/dependencies.md` **antes** de instalar. |

## Decisiones descartadas (no reintroducir)

| ID | Descartado | Motivo |
|---|---|---|
| X-01 | Tabla `UserStats` | Sustituida por `UserStreak` + `UserXp` + `UserTopicXp` (D-07). |
| X-02 | Tabla `UserBadge` / `user_badges` | Contradice D-16. El `ROADMAP.md` se contradecía a sí mismo en B4.3; gana D-16. |
| X-03 | `pending_review` + acción "publicar" | D-04. |
| X-04 | Refactor `isPublic → visibility` | D-04. |
| X-05 | `golang.org/x/sync/errgroup` | D-13. |
| X-06 | `github.com/robfig/cron/v3` | D-17. |
| X-07 | Lotes de IA de 20 preguntas / `AI_MAX_PER_JOB=50` | Un job = una llamada, máx 10; tope global (D-14). |
| X-08 | Tabla de cola persistente para la sesión diaria | D-11. |
| X-09 | IA por usuario en sesiones (`mode=generate`) | D-14. |
| X-10 | Cuota de IA por usuario (`users.daily_ai_limit`) | D-15. |
| X-11 | Cola offline persistente (guardar respuestas en disco y sincronizar al volver la red) | D-31: exige idempotencia en el backend para no duplicar respuestas; se resuelve con reintento explícito. |

## Notas de implementación relevantes para el backlog

| ID | Nota |
|---|---|
| N-01 | **`xpGained` por periodo no es derivable de las tablas de gamificación**: `UserXp` y `UserTopicXp` son acumuladores de una fila (`models/user_xp.go:12`). `UserStreak.DailyXpToday` solo guarda el día actual (`gamification/service.go:276-281`). Se deriva de `SessionAnswer` (dificultad → 5/10/15) + `20 ×` sesiones `finished` en la ventana. Ver Open Question en `SPEC-insights.md`. |
| N-02 | **`SessionAnswer.ResponseTimeMs`** (`models/session_answer.go:22`) es la fuente de tiempo real de estudio y del badge `fast_solver`. Preferirla a `FinishedAt - StartedAt` (incluye tiempo muerto). |
| N-03 | "Des-publicar" por reporte solo aplica a contenido con `is_public=true` de admin. Des-publicar contenido `ai_generated` queda **fuera de alcance** (requiere mecanismo futuro). |
| N-04 | `.agents/shared/architecture.md` citaba `DEFAULT_DAILY_AI_LIMIT` (eliminado en la fase 11). ✅ Corregido el 2026-10-17 (N-07). No reintroducir el símbolo. |
| N-05 | El estado de sesión es un **literal mágico** en Go: `"completed"` aparece en `sessions/service.go:220,225,257,313` sin constante. `"cancelled"` figura en la doc de Swagger del filtro pero ningún código lo asigna. Aplica D-21: al tocar estos estados desde `review-mastery` o `insights`, definir constantes en `internal/models/session.go` (`SessionStatusInProgress`, `SessionStatusCompleted`, `SessionStatusCancelled`) en vez de ampliar los literales. |
| N-06 | Las migraciones **ya aplicadas** `00004_sessions_mode_no_generate.sql` y `00005_drop_users_daily_ai_limit.sql` citan en sus comentarios ("Ver ROADMAP Fase 10/11") un archivo que **ya no existe**. No se editan (D-23/D-24: nunca tocar una migración aplicada). El contexto vivo de esas fases está en D-14, X-09 y X-10, y el histórico en el commit que eliminó `ROADMAP.md`. |
| N-07 | **Limpieza de docs desactualizados** (2026-10-17, commit propio tras la migración a SDD): (a) `.agents/backend/architecture.md` — la sección de concurrencia documentaba `errgroup` como canónico y afirmaba en falso que se usaba en `NextQuestion()`; reescrita con el patrón real (`sync.WaitGroup` + `sync.Mutex`, `gatherDailyData()`), con la prohibición de D-13 y aclarando que `NextQuestion()` es secuencial. También: la tabla de módulos decía que `questions` hacía "generación IA" (está en `generation` desde la fase 9) y la nota del worker apuntaba al id de backlog `B2.5`, inexistente tras borrar el `ROADMAP.md` (ahora apunta a `SPEC-notifications.md` / `NT-5`). (b) `AGENTS.md` — afirmaba que `docs/` está gitignored; en realidad está **versionado** (solo `tmp/` se ignora). (c) `.agents/shared/architecture.md` — el listado de archivos omitía `tip/moderation/gamification/generation.types.ts`, el ejemplo de re-export citaba `DEFAULT_DAILY_AI_LIMIT` (eliminado en la fase 11) y las descripciones de `utils/` y `constants/` estaban incompletas. (d) `.agents/backend/data-model.md` — `users` listaba `dailyAiLimit` (columna dropeada en la fase 11) y omitía `dailyGoalXp`; `user_xps` listaba una columna `level` que **no existe** (el nivel se deriva del XP); `user_streaks` usaba nombres de columna inventados y omitía los contadores de XP diario; `session_mode` no advertía del `CHECK` de la migración `00004`. (e) `.agents/web/architecture.md` — citaba un `utils/rules.ts` inexistente y una lista incompleta de componentes compartidos. (f) `.agents/mobile/architecture.md` y `apps/mobile/AGENTS.md` — el doc se contradecía consigo mismo sobre los tabs (decía 4: "Inicio, Temas, Progreso, Perfil"; son **3**) y documentaba `src/core/hooks/`, que no existe (ahora documenta `src/core/auth/`); faltaban las rutas `login`/`setup` y el dominio `auth`. (g) `.agents/backend/data-model.md` — decía que la IA guarda las preguntas con `created_by = requested_by`; la columna real es `user_id` (verificado contra `information_schema` y `generation.Store.SaveQuestions(userID, ...)`). |
| N-08 | **Estado real de los tests mobile** (2026-10-19): solo existen `src/core/utils/dates.test.ts` y `src/core/utils/errors.test.ts`. La infraestructura **sí** está instalada (`jest-expo`, `@testing-library/react-native` v14, `@types/jest`) y hay scripts `test` / `test:run` / `test:related`. **No hay tests** de servicios, cliente HTTP ni pantallas, y **no hay ESLint** (`"lint": "echo 'TODO: add eslint'"`). Esto es lo que cubre `mobile-tooling` (D-34). |
| N-09 | **Push en Expo requiere dev build**: desde SDK 53 las notificaciones remotas **no funcionan en Expo Go** (solo en un development build / EAS). Además `apps/mobile/app.json` **no** tiene `expo.extra.eas.projectId`, que `getExpoPushTokenAsync()` necesita. Ambas cosas se resuelven en `mobile-push`; la app debe seguir funcionando (sin push) si el token no se puede obtener. |
| N-10 | **`@testing-library/react-native` v14 es async**: `render`, `rerender`, `unmount` y todas las variantes de `fireEvent` devuelven `Promise`, así que hay que `await`-earlas. `screen` solo es válido después de `await render`. Tampoco existen las queries `UNSAFE_*`: para `TextInput` se usa `testID`. Detalle en `specs/SPEC-mobile-tooling.md`. |
| N-11 | **Jest no resuelve alias de Metro**: hay que mapear `@devtest/shared` (cuyo `main` apunta a TypeScript) y `@/` en `moduleNameMapper` de `apps/mobile/jest.config.js`. El primer render de cada archivo tarda varios segundos → `testTimeout: 30000`. |
| N-12 | **`expo-notifications` debe cargarse en diferido** (2026-10-25): en Expo Go para Android (SDK 57) el nativo `ExpoTopicSubscriptionModule` no está registrado y `expo-notifications` lo resuelve al importarse, así que un `import` estático deja la app en pantalla roja antes de renderizar (`Cannot find native module 'ExpoTopicSubscriptionModule'` → `Route "./.../_layout.tsx" is missing the required default export` → `Cannot read property 'ErrorBoundary' of undefined`; ver expo/expo#52175). `src/core/push/notifications.ts` usa `loadNotifications()` con `require()` perezoso y se salta el módulo cuando `Platform.OS === 'android' && Constants.executionEnvironment === ExecutionEnvironment.StoreClient` (ahí el push remoto tampoco existe desde SDK 53, N-09). Nunca volver al `import * as Notifications from 'expo-notifications'` estático. |
