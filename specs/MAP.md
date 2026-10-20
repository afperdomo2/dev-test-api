# Capability Map — Backlog funcional

> **Origen**: migrado desde `ROADMAP.md` el 2026-10-17. Las **fases 0–11 ya completadas no se migran** (quedan en el historial de git). `ROADMAP.md` se elimina en el mismo commit que estos artefactos.
>
> **Este archivo es el índice**. No describe implementación: para eso están `specs/SPEC-<id>.md` (contrato) y `tasks/` (plan + tareas). Las decisiones transversales viven en `specs/DECISIONS.md`.

## Módulos

| Module id | Responsabilidad | Depende de | IDs originales | Estado |
|---|---|---|---|---|
| `review-mastery` | Dominio (accuracy) **por tema** y detección de temas débiles | `progress` (existente), `sessions` (existente) | B1.1–B1.8 | `done` |
| `notifications` | Notificaciones in-app + job diario de repaso | `sessions` (due count ya existente) | B2.1–B2.9 | `ready` |
| `insights` | Agregados de actividad (semana/mes): KPIs, temas, racha, XP | `review-mastery`, `gamification`, `sessions` | B3.1–B3.8 | `ready` |
| `badges` | Catálogo de logros `const`, evaluados on-the-fly | `gamification`, `review-mastery`, `insights` | B4.1–B4.5 | `opt-in` |
| `e2e` | Flujo E2E con Playwright | todos los entregados | B5.1 | `opt-in` |

**Build order**: `review-mastery` → `notifications`, `insights` → `badges` → `e2e`

Sin ciclos. `notifications` e `insights` son independientes entre sí y pueden ejecutarse en paralelo una vez cerrado `review-mastery`; `badges` consume ambos; `e2e` va al final por definición (prueba el sistema integrado).

## Interfaces en las fronteras

| Proveedor | Contrato expuesto | Consumidor |
|---|---|---|
| `review-mastery` | `progress.Service.TopicMastery() ([]TopicMasteryItem, error)` — incluye la definición única de `isWeak` | `insights` (filtra `isWeak`), `badges` (`topics_mastered_5`) |
| `sessions` | `Store.CountDueReviews(userID) (int64, error)` — **ya existe** (`sessions/store.go:210`), no se reimplementa | `notifications` (job diario) |
| `gamification` | `Service.GetOverview()` (XP, streak, `topicXp`), `RegisterAnswer`, `CompleteSession` | `insights`, `badges` |
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
| B2.1 | `NT-1` | Modelo + migración `00006` |
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

- **Mobile** (Expo/RN): ningún módulo de este mapa incluye tareas mobile. La sesión diaria es mobile-first, pero estas funcionalidades son backend + web.
- **Cambios de esquema** en `review-mastery` e `insights`: son de solo lectura.
- **Retención/purga** de notificaciones y **`xp_events` ledger**: ver Open Questions de cada spec.
