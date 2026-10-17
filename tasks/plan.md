# Plan técnico — Backlog funcional (B1–B5)

> Referencias: `specs/MAP.md` (contratos y orden), `specs/SPEC-<id>.md` (detalle), `specs/DECISIONS.md` (guardarraíles). Las tareas concretas están en `tasks/todo.md`.

## 1. Grafo de dependencias

```
                        ┌──────────────────┐
                        │ review-mastery   │  (sin migración)
                        │ progress + consts│
                        └────────┬─────────┘
                                 │
              ┌──────────────────┴──────────────────┐
              ▼                                     ▼
   ┌────────────────────┐                ┌────────────────────┐
   │ notifications      │                │ insights           │
   │ + migración 00006  │                │ (sin migración)    │
   │ + scheduler        │                │ reusa mastery+gamif│
   └────────────────────┘                └─────────┬──────────┘
                                                   │
                                     ┌─────────────┴─────────────┐
                                     ▼                           ▼
                            ┌────────────────┐          ┌─────────────────┐
                            │ badges (opt-in)│          │ e2e (opt-in)    │
                            │ reusa ambos    │          │ flujo completo  │
                            └────────────────┘          └─────────────────┘
```

- `review-mastery` es **bloqueante** para `insights` (definición única de `isWeak`, D-08) y para `badges` (`topics_mastered_5`).
- `notifications` **no** depende de `review-mastery`: reutiliza `sessions.CountDueReviews`, que ya existe (D-09). Puede ejecutarse en **paralelo** con `review-mastery` e `insights`.
- `e2e` va al final: prueba el sistema integrado; ejecutarlo antes solo produce falsos negativos.

**Caminos paralelizables**: `notifications` ↔ (`review-mastery` → `insights`). Ambas ramas solo se tocan en `internal/server/server.go` (wiring) y en `DefaultLayout.vue` (nav/campana) — dos puntos de conflicto, triviales.

## 2. Orden de implementación

| # | Módulo | Por qué en esta posición |
|---|---|---|
| 1 | `review-mastery` | Desbloquea `insights` y `badges`; sin migración (riesgo mínimo); define `SessionStatus*` (N-05) que usan los siguientes |
| 2 | `notifications` | Independiente; es el único con **migración + job de fondo** → conviene aislarlo para validar el esquema con calma |
| 3 | `insights` | Consume `mastery` y `gamification`; la parte delicada es `xpGained` (N-01) |
| 4 | `badges` (opt-in) | Consume `mastery` + `gamification`; catálogo puro |
| 5 | `e2e` (opt-in) | Valida todo lo anterior de punta a punta |

## 3. Componentes y dónde vive cada cosa

| Componente | Touch points backend | Touch points frontend |
|---|---|---|
| Constantes de estado de sesión | `internal/models/session.go` | — |
| Dominio por tema | `internal/modules/progress/*` | `ProgressPage` + `TopicMasteryList` |
| Notificaciones | `internal/models/notification.go`, `modules/notifications/*`, `server.go`, migración | campana en `DefaultLayout`, `NotificationsPage`, `router/routes.ts` |
| Informes | `internal/modules/insights/*`, helpers en `gamification`, `server.go` | `InsightsPage`, nav "Estudio", `router/routes.ts` |
| Logros | `modules/gamification/{catalog,badges,response,handler,routes}.go` | `AchievementsPage`, nav |
| E2E | — | `apps/web/e2e/`, `playwright.config.ts`, `Makefile` |

**Regla transversal**: cada módulo nuevo en backend sigue la anatomía de `.agents/backend/architecture.md` (store → service → handler → routes) y se registra en `server.go` en orden de dependencias (stores → services → handlers). En frontend, el checklist de `.agents/web/architecture.md`: types → service → query → página/componente → ruta.

## 4. Riesgos y mitigaciones

| Riesgo | Prob. | Impacto | Mitigación |
|---|---|---|---|
| Columna `GENERATED ... STORED` rechazada por Postgres (expresión no `IMMUTABLE`) | Baja | Medio | Fallback documentado en `SPEC-notifications.md`: `creation_day` como columna normal escrita por la app, mismo índice único |
| Job diario duplica notificaciones | Media | Alto | Idempotencia **en la DB** (índice único `(user_id, type, creation_day)`), no en el código; test de doble ejecución |
| `xpGained` se desvía del XP real tras cambiar `sessionBonusXp`/`xpByDifficulty` | Media | Bajo | Helpers exportados en `gamification` (fuente única); comentario en el código; OQ-1 de insights documenta el límite |
| Agregaciones lentas sin índice en `question_topics` | Baja | Medio | Medir con datos reales antes de crear migración; el índice va como "Ask first" en la spec |
| `errgroup` reintroducido por copiar un doc desactualizado | Baja | Medio | D-13 lo prohíbe; `.agents/backend/architecture.md` quedó corregido (N-07) y documenta el patrón `sync.WaitGroup` + `sync.Mutex` real; las 3 specs que paralelizan dicen "`sync.WaitGroup`, no `errgroup`" |
| Se "arregla" la asimetría `topTopics` (ventana) vs `weakTopics` (histórico) y se rompe D-08 | Media | Alto | Documentado en `SPEC-insights.md` con el porqué; test de consistencia `weakTopics == mastery.filter(isWeak)` |
| Navegación/aparición del nav rompe `DefaultLayout.test.ts` | Alta | Bajo | Los tests de nav se actualizan en la tarea del módulo correspondiente (ya previsto en los specs) |
| `tsc`/`vue-tsc` falla por tipos compartidos cruzados (`TopicMasteryItem` en `insight.types.ts`) | Media | Bajo | `TopicMasteryItem` se **importa**, no se re-exporta desde dos sitios (regla anti-duplicación) |
| E2E accidenta una DB de desarrollo | Baja | Alto | DB dedicada `devtest_e2e`; la suite no hace DROP; `workers: 1` |

## 5. Checkpoints de verificación

| Checkpoint | Cuándo | Criterio |
|---|---|---|
| V1 | Al cerrar `review-mastery` | `make check` + `make test` + `make swagger` verdes; `/progress/topics-mastery` responde en la DB real; `make fe-check` + `make fe-lint` + `make fe-test` verdes |
| V2 | Al cerrar `notifications` | `make migrate-up` + `make migrate-down` limpios; doble ejecución del job → 0 duplicados; arranque y parada del servidor sin goroutines colgadas |
| V3 | Al cerrar `insights` | Los 11 campos cuadran con el cálculo manual; `weakTopics` idéntico a `mastery.filter(isWeak)`; `dailyActivity` con 7/30 puntos |
| V4 | Al cerrar `badges` | Los 8 badges con umbrales en el borde exacto; sin tabla nueva |
| V5 | Al cerrar E2E | `make e2e` verde desde DB limpia; `make fe-test` no lo ejecuta |
| V6 | **Verificación final** (re-run de la Fase 8) | `make fmt` + `make vet` + `make test` · `make fe-check` + `make fe-lint` + `make fe-test` · `pnpm -r typecheck` |

## 6. Qué NO se hace (fuera de alcance, ya decidido)

- **Mobile**: ningún módulo toca `apps/mobile/` (D-18).
- **Retención de notificaciones**, **ledger `xp_events`**, **persistencia de logros**, **CI del E2E**: cada uno está como Open Question en su spec; ninguno se implementa sin decisión explícita.
- **Refactors colaterales**: no se arreglan los docs desactualizados de `.agents/` (D-13, N-04) ni los literales mágicos fuera de lo que tocan estas specs. Se anotan, no se arreglan.

## 7. Estrategia de commits

Un commit por tarea (o por puñado de tareas acopladas del mismo módulo), con el mensaje referenciando `specs/SPEC-<id>.md` y el id de tarea. **Los commits son manuales**: el agente deja el working tree listo y el usuario commitea.
