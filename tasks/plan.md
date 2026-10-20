# Plan técnico — Backlog funcional (web + mobile)

> Referencias: `specs/MAP.md` (contratos y orden), `specs/SPEC-<id>.md` (detalle), `specs/DECISIONS.md` (guardarraíles). Las tareas concretas están en `tasks/todo.md`.

## 1. Grafo de dependencias

```
review-mastery ✅ ──────────────► mobile-progress ──► mobile-weekly ◄── insights
                                 │                    └► mobile-badges ◄── badges
                                 │
notifications ──► push ──► mobile-push

(sin dependencias)  mobile-tooling · mobile-daily-tip · mobile-offline · e2e
```

- `review-mastery` ✅ ya está hecho y desbloquea `mobile-progress`. La pestaña que crea `mobile-progress` es la que reutilizan `mobile-weekly` y `mobile-badges` (por eso va antes).
- `push` **cuelga de `notifications`** (D-28): la tabla de avisos es la fuente y `push` solo entrega.
- `mobile-tooling`, `mobile-daily-tip`, `mobile-offline` y `e2e` no dependen de nada y pueden ir en paralelo.
- `mobile-offline` toca backend (idempotencia, D-35) y mobile: es el único módulo mixto del backlog mobile.

## 2. Orden de implementación (olas)

| Ola | Módulos | Por qué |
|---|---|---|
| 0 | `review-mastery` ✅ · `mobile-tooling` ✅ · `mobile-progress` ✅ | Hechos. El endpoint de temas y la pestaña Progreso ya existen |
| 1 | `mobile-daily-tip` | Barato y cierra un hueco: el tip se descarta hoy |
| 2 | `insights` → `mobile-weekly` | Backend primero; la pantalla después |
| 3 | `badges` (backend → web) → `mobile-badges` | Igual: endpoint antes que cliente |
| 4 | `mobile-offline` · `e2e` | Robustez y verificación |
| 5 | `notifications` → `push` → `mobile-push` | **Al final por priorización (D-36)**. `push` cuelga de `notifications` (D-28) |

**Paralelizable dentro de cada ola**: en la Ola 1 se puede trabajar `notifications` (backend) mientras `mobile-tooling` y `mobile-daily-tip` van por otro lado; solo comparten `server.go`/`DefaultLayout` en puntos triviales.

## 3. Componentes y touch points

| Componente | Backend | Web | Mobile |
|---|---|---|---|
| `notifications` | `models/notification.go`, `modules/notifications/*`, `server.go`, migración `00007` | campana, `NotificationsPage`, ruta | — |
| `insights` | `modules/insights/*`, helpers XP en `gamification` | `InsightsPage`, nav | — |
| `badges` | `gamification/{catalog,badges}.go`, endpoint | `AchievementsPage`, ruta, nav | — |
| `push` | `models/device_token.go`, `modules/push/*`, migración `00008` | — | — |
| `mobile-tooling` | — | — | eslint, jest, tests |
| `mobile-daily-tip` | — | — | store + card + Home + Study |
| `mobile-progress` | — | — | tab + hub + Temas |
| `mobile-weekly` | — | — | pantalla + servicio |
| `mobile-badges` | — | — | pantalla + servicio |
| `mobile-push` | — | — | `expo-notifications` + registro |
| `mobile-offline` | idempotencia en `sessions` (migración) | — | reintento + UUID |
| `e2e` | — | Playwright | — |

**Reglas transversales**: backend sigue la anatomía de `.agents/backend/architecture.md`; web el checklist de `.agents/web/architecture.md`; mobile la estructura `app/` (rutas thin) + `src/features/` + `src/core/` de `.agents/mobile/architecture.md`.

## 4. Riesgos y mitigaciones

| Riesgo | Prob. | Impacto | Mitigación |
|---|---|---|---|
| **Push no funciona en Expo Go** (SDK 53+) | **Alta** | Alto | N-09: probar en dev build; la app degrada sin push; documentarlo en el README de mobile |
| **Falta `eas.projectId`** en `app.json` | Alta | Medio | `MPUSH-1` lo añade; si falta, `getDevicePushToken` devuelve `null` y avisa por consola |
| Reintento de respuesta duplica XP/SRS | Media | **Alto** | D-35: `clientAnswerId` + índice único; el service no re-aplica progreso si la fila ya existía |
| El job diario se bloquea esperando a Expo | Media | Medio | Push best-effort con timeout; el aviso in-app se crea antes |
| Tokens push muertos acumulados | Media | Bajo | `DeviceNotRegistered` borra el token; purga por antigüedad queda como OQ |
| `errgroup` reintroducido por copiar docs viejos | Baja | Medio | D-13 + doc corregido (N-07); las specs dicen "`sync.WaitGroup`, no `errgroup`" |
| Se rompe la barra de pestañas al pasar a 4 | Media | Medio | `MPROG-4` incluye verificación manual en dispositivo (iOS/Android) |
| Tests retroactivos de mobile fallan por mocks nativos | Media | Bajo | `jest.setup.ts` centraliza mocks de `expo-secure-store`, `expo-router`, `expo-notifications` |
| Se "arregla" la asimetría `topTopics` (ventana) vs `weakTopics` (histórico) | Media | Alto | Documentado en `SPEC-insights.md`; test de consistencia |
| `noUncheckedIndexedAccess` / `no-non-null-assertion` rompen tests nuevos | Alta | Bajo | Ya nos pasó en RM-8: usar `?.` y helpers que lancen, no `!` |
| `v-card` inyecta un `role="progressbar"` oculto que confunde tests | Media | Bajo | Selectores `[role="progressbar"][aria-label]` (documentado en RM-8) |
| Dependencias nuevas (`expo-notifications`, `expo-crypto`, `eslint-config-expo`) | Alta | Bajo | Gate obligatorio en `.agents/mobile/dependencies.md` antes de instalar |

## 5. Checkpoints de verificación

| Checkpoint | Cuándo | Criterio |
|---|---|---|
| V1 | `notifications` | `migrate-up/down` limpios; doble ejecución del job → 0 duplicados; 4 endpoints; swagger |
| V2 | `mobile-tooling` | `lint` deja de ser un TODO; ≥ 12 archivos de test; cliente HTTP cubierto |
| V3 | `mobile-progress` | 4 pestañas en dispositivo; Temas coincide con el endpoint; vacío y error |
| V4 | `insights` + `mobile-weekly` | 11 campos cuadran; `weakTopics` idéntico a `mastery.filter(isWeak)`; semana/mes |
| V5 | `badges` + `mobile-badges` | Umbrales en el borde; sin tabla; web y mobile muestran lo mismo |
| V6 | `push` + `mobile-push` | Token registrado una vez; sin projectId degrada; baja en logout; aviso recibido en dev build |
| V7 | `mobile-offline` | Reintento no duplica (1 fila, 1 XP); sin red no avanza a resultado |
| V8 | `e2e` | `make e2e` verde desde DB limpia |
| V9 | **Verificación final** | `make check` + `make test` · `make fe-*` · `pnpm -r typecheck` · `pnpm --filter mobile test:run` |

## 6. Qué NO se hace

- **Cola offline persistente** (X-11) y **push con Firebase/APNs directos** (D-27).
- **Retención de notificaciones**, **ledger `xp_events`** y **persistencia de logros**: Open Questions de sus specs.
- **E2E de la app mobile** (Detox/Maestro): fuera de alcance.
- **Refactors colaterales**: no se arreglan más docs desactualizados ni literales mágicos fuera de lo que tocan estas specs.

## 7. Estrategia de commits

Un commit por módulo (o por tarea, si el módulo es grande). El mensaje referencia `specs/SPEC-<id>.md` y el id de tarea. **Los commits son manuales**: el agente deja el working tree listo.
