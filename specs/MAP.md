# Capability Map — Backlog funcional

> **Origen**: migrado desde `ROADMAP.md` el 2026-10-17. Las **fases 0–11 ya completadas no se migran** (quedan en el historial de git). `ROADMAP.md` se elimina en el mismo commit que estos artefactos.
>
> **Este archivo es el índice**. No describe implementación: para eso están `specs/SPEC-<id>.md` (contrato) y `tasks/` (plan + tareas). Las decisiones transversales viven en `specs/DECISIONS.md`.
>
> **Iniciativa separada**: el endurecimiento de seguridad tiene su propio mapa en `specs/MAP-security.md` y sus tareas en `tasks/security-todo.md`. No forma parte de este backlog funcional.

## Módulos

| Module id | Responsabilidad | Depende de | IDs originales | Estado |
|---|---|---|---|---|
| `review-mastery` | Dominio (accuracy) **por tema** y detección de temas débiles | `progress` (existente), `sessions` (existente) | B1.1–B1.8 | `done` |
| `notifications` | Notificaciones in-app + job diario de repaso | `sessions` (due count ya existente) | B2.1–B2.9 | `done` |
| `insights` | Agregados de actividad (semana/mes): KPIs, temas, racha, XP | `review-mastery`, `gamification`, `sessions` | B3.1–B3.8 | `done` |
| `badges` | Catálogo de logros `const`, evaluados on-the-fly | `gamification`, `review-mastery`, `insights` | B4.1–B4.5 | `done` |
| `e2e` | Flujo E2E con Playwright | todos los entregados | B5.1 | `done` |
| `push` | Entrega push: tokens de dispositivo + Expo Push API | `notifications` | nuevo (mobile) | `done` |
| `mobile-tooling` | ESLint + tests del código mobile existente y nuevo | — | nuevo (mobile) | `done` |
| `mobile-daily-tip` | Mostrar el tip del día al empezar la sesión diaria | — (el backend ya lo envía) | nuevo (mobile) | `done` |
| `mobile-progress` | 4ª pestaña **Progreso** + pantalla de Temas | `review-mastery` ✅ | nuevo (mobile) | `done` |
| `mobile-weekly` | "Mi semana": KPIs del periodo dentro de Progreso | `insights` ✅ | nuevo (mobile) | `done` |
| `mobile-badges` | Logros dentro de Progreso | `badges` (backend) | nuevo (mobile) | `done` |
| `mobile-progress-ux` | Progreso como dashboard (anillos, actividad, logros) + Inicio acción-first | `mobile-progress` ✅, `mobile-weekly` ✅, `mobile-badges` ✅ | nuevo (mobile) | `done` |
| `mobile-push` | Registrar el dispositivo y recibir/abrir push | `push` | nuevo (mobile) | `blocked` (EAS/dev build) |
| `mobile-offline` | Reintento con aviso al fallar el envío de una respuesta | — | nuevo (mobile) | `done` |

**Build order** (olas por dependencia; lo que no bloquea va en paralelo):

```
Ola 0  review-mastery ✅ · mobile-tooling ✅ · mobile-progress ✅
Ola 1  mobile-daily-tip
Ola 2  insights  →  mobile-weekly
Ola 3  badges (backend → web)  →  mobile-badges
Ola 4  mobile-offline  ·  e2e
Ola 5  notifications  →  push  →  mobile-push   (al final: priorización del producto, D-36)
```

Sin ciclos. `mobile-progress` crea la pestaña que `mobile-weekly` y `mobile-badges` reutilizan. `push` cuelga de `notifications` (D-28): la tabla de avisos es la fuente y el dispatcher solo entrega; por eso ambos quedan al final junto con `mobile-push`.

## Interfaces en las fronteras

| Proveedor | Contrato expuesto | Consumidor |
|---|---|---|
| `review-mastery` | `progress.Service.TopicMastery() ([]TopicMasteryItem, error)` — incluye la definición única de `isWeak` | `insights` (filtra `isWeak`), `badges` (`topics_mastered_5`) |
| `sessions` | `Store.CountDueReviews(userID) (int64, error)` — **ya existe** (`sessions/store.go:210`), no se reimplementa | `notifications` (job diario) |
| `gamification` | `Service.GetOverview()` (XP, streak, `topicXp`), `RegisterAnswer`, `CompleteSession` | `insights`, `badges` |
| `insights` | `insights.Service.Summary(userID, period)` + `GET /insights/summary` | `mobile-weekly` |
| `badges` | `GET /gamification/badges` | `mobile-badges` |
| `notifications` | `Service.RunDailyJob()` + `Pusher` interface (`Push(ctx, []Notification) error`) | `push` |
| `push` | `POST /push/tokens`, `POST /push/tokens/unregister`, `ExpoSender` | `mobile-push` |
| `insights` | `insights.Service.Summary(userID, period)` | `badges` (`answers_100` si se decide derivar de ahí) |

**Regla de reutilización**: `isWeak`, "tema débil" y "due review" tienen **una sola definición** en el sistema. Ningún módulo los recalcula (ver `specs/DECISIONS.md` D-08, D-09).

## Trazabilidad `ROADMAP.Bx.y → spec → tarea`

| ID original | Tarea | Nota |
|---|---|---|
| B1.1 | `RM-1` | ⚠️ `CountDue(userID)` **no se implementa**: duplica `sessions.CountDueReviews`. `FindTopicMastery` se renombra a `FindTopicAccuracy` (colisiona con el `IsMastered` del SRS). Incluye las constantes `SessionStatus*` (N-05) |
| B1.2 | `RM-2` | Umbrales como constantes en Go |
| B1.3 | `RM-3` | ⚠️ `TopicMasteryResponse` (wrapper) **no se crea**: `response.Success` ya envuelve en `{ data }` |
| B1.4 | `RM-4` | `GET /progress/topics-mastery` |
| B1.5 | `RM-5` | Tests service + store |
| B1.6 | `RM-6` | Shared types + service + query |
| B1.7 | `RM-7` | 4ª pestaña "Temas" en `ProgressPage` |
| B1.8 | `RM-8` | Tests web + checks |
| B2.1 | `NT-1` | Modelo + migración `00007` |
| B2.2 | `NT-2` | Módulo `notifications` |
| B2.3 | `NT-3` | 4 endpoints |
| B2.4 | `NT-4` | Job diario (fan-out idempotente) |
| B2.5 | `NT-5` | ⚠️ **sin `robfig/cron`**: ticker stdlib |
| B2.6 | `NT-6` | Tests |
| B2.7 | `NT-7` | Shared types + service + queries |
| B2.8 | `NT-8` | Campana + página + ruta |
| B2.9 | `NT-9` | Tests web + checks |
| B3.1 | `IN-1` | Módulo `insights` + helpers de XP exportados en `gamification` |
| B3.2 | `IN-3` | `GET /insights/summary` (handler/rutas/swagger) |
| B3.3 | `IN-2` | Agregaciones + `Summary()` (incl. `xpGained` derivado) |
| B3.4 | `IN-4` | Tests |
| B3.5 | `IN-5` | Shared types + service + queries |
| B3.6 | `IN-6` | Página + ruta |
| B3.7 | `IN-7` | Nav "Informes" en sección Estudio |
| B3.8 | `IN-8` | Tests web + checks |
| B4.1 | `BD-1` | Catálogo `const` |
| B4.2 | `BD-2` | ⚠️ **no** se evalúa en `RecordAnswer`/`RecordSessionFinished` (derivado, no persistido): la evaluación vive en `gamification/badges.go` |
| B4.3 | `BD-3` | ⚠️ **sin tabla `user_badges`** (decisión bloqueada). `GET /gamification/badges` deriva |
| B4.4 | `BD-4` | Página de logros + ruta |
| B4.5 | `BD-5` | Tests + checks |
| B5.1 | `E2-1` | Playwright |

**Cobertura**: 31/31 items con tarea asignada. 5 items tienen desviación documentada (marcados ⚠️), todas resueltas en las specs y en `DECISIONS.md`.

## Flujo de trabajo

1. Marcar la tarea como `in_progress` en `tasks/todo.md` (`- [ ]` → `- [~]`).
2. Implementar + verificar con los comandos que la tarea declara en `Verify:`.
3. Marcar completada (`- [~]` → `- [x]`).
4. Commit **manual** (el agente nunca commitea).

**Leyenda de estado**: `[ ]` pendiente · `[~]` en progreso · `[x]` hecha · `[-]` descartada (con motivo en la nota).

Los gates de SDD (revisión de la spec antes de implementar, y del plan antes de las tareas) se aplican por módulo, en el build order de arriba.

## Fuera de alcance

- **Retención/purga** de notificaciones y **`xp_events` ledger**: ver Open Questions de cada spec.
- **Cola offline persistente** (guardar respuestas en disco y sincronizar): descartada en D-31 / X-11.
- **Push con Firebase/APNs directos**: descartado en D-27 (se usa Expo Push).
- **E2E de la app mobile** (Detox/Maestro): no se incluye; el E2E del plan es solo web (Playwright).

## Backlog mobile (sin IDs `Bx.y`)

Los módulos mobile son alcance nuevo (D-26); no tienen IDs heredados del `ROADMAP`. Su trazabilidad es el propio module id:

| Module id | Spec | Tareas |
|---|---|---|
| `push` | `SPEC-push.md` | `PUSH-1…` |
| `mobile-tooling` | `SPEC-mobile-tooling.md` | `MTOOL-1…` |
| `mobile-daily-tip` | `SPEC-mobile-daily-tip.md` | `MTIP-1…` |
| `mobile-progress` | `SPEC-mobile-progress.md` | `MPROG-1…` |
| `mobile-weekly` | `SPEC-mobile-weekly.md` | `MWEEK-1…` |
| `mobile-badges` | `SPEC-mobile-badges.md` | `MBADGE-1…` |
| `mobile-progress-ux` | `SPEC-mobile-progress-ux.md` | `MUX-1…` |
| `mobile-push` | `SPEC-mobile-push.md` | `MPUSH-1…` |
| `mobile-offline` | `SPEC-mobile-offline.md` | `MOFF-1…` |
