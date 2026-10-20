# SPEC — e2e

| Campo | Valor |
|---|---|
| Module id | `e2e` |
| IDs originales | B5.1 |
| Estado | `ready` (aprobado en D-32) |
| Depende de | backend + web entregados (todo el mapa) |
| Migración goose | **Ninguna** |
| Dependencia nueva | `@playwright/test` (dev) — ⚠️ requiere pasar por `.agents/web/dependencies.md` (D-25) |

## Objective

Tener una prueba E2E que valide el camino crítico de un usuario nuevo **de punta a punta**, contra el backend y la DB reales, para detectar roturas de integración que los tests unitarios (store mockeado, Vitest con `apiClient` mockeado) no pueden ver: contrato de la API, snake/camel de JSON, guards del router, JWT real.

**Flujo objetivo (B5.1)**: `status → setup → login → dashboard → sessions`.

**Éxito**: `make e2e` levanta todo, recorre el flujo y termina en verde desde una DB limpia.

### Por qué necesita una DB limpia

El flujo empieza en **setup** porque `POST /auth/setup` solo funciona **cuando no existe ningún admin** (`auth` module). El guard del router redirige a `/setup` mientras `authStore.needsSetup` sea `true` (`router/index.ts:23-28`). Por tanto la suite debe correr contra una DB **dedicada y vacía**, no contra la de desarrollo.

## Tech stack

- **Playwright** (`@playwright/test`) — Chromium como proyecto por defecto.
- Backend Go + Gin + Postgres (real, vía `make run`).
- Web Vite dev server con el proxy `/api` → `:8080` que ya existe.

⚠️ **`make fe-test` es Vitest y no ejecuta esto.** El E2E es un runner aparte: no se mezcla con la suite unitaria ni entra en `make fe-check`/`make fe-lint`.

## Commands

```bash
# 1) DB dedicada y limpia
DB_NAME=devtest_e2e make migrate-up

# 2) Backend contra esa DB (consume .env; sobreescribir DB_NAME por entorno)
make run

# 3) Suite E2E (levanta Vite automáticamente vía webServer)
make e2e

# Variantes
make e2e-ui       # modo UI/headed para depurar
make e2e-report   # abre el reporte HTML del último run
```

Targets nuevos en el `Makefile` (delegando en `pnpm --filter <web> test:e2e`). Scripts nuevos en `apps/web/package.json`: `test:e2e`, `test:e2e:ui`, `test:e2e:report`.

## Project structure

```
apps/web/
├── playwright.config.ts          # NUEVO
├── e2e/
│   ├── fixtures.ts               # NUEVO: helpers (reset de DB, usuario de prueba)
│   ├── critical-flow.spec.ts     # NUEVO: status → setup → login → dashboard → sessions
│   └── .auth/                    # storageState (gitignored)
└── package.json                  # + scripts test:e2e*

Makefile                          # + e2e, e2e-ui, e2e-report
```

## Contrato de la suite

### `playwright.config.ts`

```ts
export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,     // comparte estado en la DB: sin paralelismo
  workers: 1,
  retries: 0,               // local; en CI (si algún día) 1
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  webServer: {
    command: 'pnpm dev',
    url: 'http://localhost:3000',
    reuseExistingServer: true,   // permite `make fe-dev` ya levantado
    timeout: 120_000,
  },
})
```

El backend **no** lo levanta Playwright: se asume `make run` corriendo en `:8080`. Si `/health` no responde, la suite debe fallar con un mensaje claro ("levanta el backend con `make run`"), no con un timeout críptico.

### Escenario `critical-flow.spec.ts`

1. **status** → `GET /auth/status` confirma que hace falta setup; la app redirige a `/setup`.
2. **setup** → rellenar `SetupPage`, enviar; se crea el primer admin.
3. **login** → el guard lleva a `/login`; login con el admin recién creado; se guarda el JWT (storage del web) y se entra a `/`.
4. **dashboard** → `DashboardPage` renderiza sus cards sin errores de consola ni respuestas 4xx/5xx.
5. **sessions** → navegar a `/sessions`, crear una sesión (o abrir la lista y comprobar el estado vacío), y verificar que la página carga con datos reales del backend.

**Aserciones transversales** (lo que aporta valor frente a los unitarios):
- Ninguna respuesta de `/api/**` con status ≥ 400 durante el flujo (salvo el `401` esperado de `/auth/status` pre-login, si aplica).
- Sin errores en `page.on('console')` de tipo `error` ni `pageerror`.
- El JWT sobrevive a un reload (persistencia real).

### Aislamiento y estado

- `fullyParallel: false` + `workers: 1`: la suite asume una DB que **ella misma** deja configurada.
- La DB se prepara **fuera** de Playwright (`DB_NAME=devtest_e2e make migrate-up` sobre una DB vacía), porque el `setup` requiere ausencia de admin. La suite **no** borra datos: si se re-ejecuta contra una DB ya configurada, el primer paso fallará por diseño; el mensaje debe indicar cómo recrear la DB.
- `storageState` en `e2e/.auth/` (gitignored) para reutilizar la sesión entre specs si se añaden más adelante.

## Code style

- Selectores por **rol y texto accesible** (`getByRole`, `getByLabel`) — los mismos criterios de accesibilidad que exige `ui-ux-pro-max` para la UI. Sin selectores por clase CSS.
- Sin `waitForTimeout`: usar auto-waiting de Playwright + `expect(...).toBeVisible()`.
- TypeScript, siguiendo el `tsconfig` del web.

## Testing strategy

| Nivel | Qué cubre |
|---|---|
| E2E (Playwright) | Integración real: contrato API, guards, JWT, render con datos reales |
| Unit (Vitest) | Ya existente, **no se reemplaza** |
| Backend (`make test`) | Ya existente, **no se reemplaza** |

El E2E **no** busca cobertura: cubre **un** camino crítico. No se duplican aquí validaciones de reglas de negocio (eso es de los unitarios).

## Boundaries

- **Always**: pasar por `.agents/web/dependencies.md` antes de `pnpm add -D @playwright/test` (D-25); documentar en el `README.md`/`Makefile` que el E2E necesita DB dedicada + backend corriendo; selectores accesibles; mantener el E2E **fuera** de `make fe-test`.
- **Ask first**: integrarlo en CI (requiere servicio Postgres en GitHub Actions + secretos → decisión aparte, ver OQ-2); añadir más escenarios; instalar navegadores adicionales (Firefox/WebKit).
- **Never**: reutilizar la DB de desarrollo (rompería el setup y borraría datos reales); meterlo en `make fe-check`/`make fe-lint`/`make fe-test`; `waitForTimeout` arbitrarios; `AutoMigrate`; editar `docs/` a mano; tocar mobile.

## Success Criteria

1. `make e2e` recorre `status → setup → login → dashboard → sessions` desde una DB limpia y termina en verde.
2. Si el backend no está corriendo, la suite falla con un mensaje que explica cómo levantarlo.
3. Si la DB ya tiene admin, la suite falla indicando cómo recrear la DB `devtest_e2e`.
4. El flujo no registra errores de consola ni respuestas `/api/**` ≥ 400.
5. `make fe-test`, `make fe-check` y `make fe-lint` **no** ejecutan el E2E.
6. `@playwright/test` es **devDependency** del web y el gate de `dependencies.md` quedó documentado.
7. `e2e/.auth/` está en `.gitignore` y no se commitea estado de sesión.

## Notas de implementación (2026-10-20)

- **Dependencia**: `@playwright/test@^1.64.1` (Microsoft, MIT) como devDependency + Chromium 1248 (`playwright install chromium`, una vez por máquina). Gate de `.agents/web/dependencies.md`: sin advisories propios; `pnpm audit` conserva los 18 hallazgos preexistentes del workspace.
- **Flujo ajustado a los guards reales**: el admin creado en `setup` no puede entrar a `/sessions` (`requiresNotAdmin`). Tras el dashboard admin, la suite crea un **usuario normal desde la UI de Usuarios**, cierra sesión y valida Sesiones con él (estado vacío contra el backend real).
- **Warm-up de Vite** (`beforeAll`): pide cada `src/**/*.vue|ts` al dev server antes de abrir el navegador. Sin esto la primera visita a una página lazy descubre componentes Vuetify (`vuetify/components/VDataTable`, …), Vite reoptimiza y recarga la página en mitad de la navegación, perdiendo el click. No se cambia `pnpm dev` (contrato de esta spec).
- **Timeouts**: test `120 s` y `expect` `30 s`; el arranque en frío de Vite+Vuetify supera los 30 s por defecto.
- **DB vacía a mano** (OQ-1): recrear con `DROP DATABASE devtest_e2e; CREATE DATABASE devtest_e2e;` + `DB_NAME=devtest_e2e make migrate-up`. Si ya hay admin, el spec falla con ese comando.
- **Artefactos**: `apps/web/e2e/.auth/`, `playwright-report/` y `test-results/` gitignored. `make fe-test`/`fe-check` no ejecutan Playwright (Vitest solo incluye `src/**`).
- **Evidencia**: `make e2e` verde (1 test, 25.3 s) con Vite arrancado por Playwright y `devtest_e2e` recién migrada.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿La suite crea/recrea la DB por sí misma (drop + migrate) o se prepara a mano? | **A mano** (`DB_NAME=devtest_e2e make migrate-up` sobre DB vacía). Automatizarlo implica que el test hace DROP de una DB: un accidente con la variable mal puesta borraría datos de desarrollo. Se documenta el comando en el mensaje de error. |
| OQ-2 | ¿CI? | **Fuera de alcance de B5.1.** Requiere servicio Postgres en el workflow, seed del admin y orquestar dos procesos; es una tarea de infraestructura propia (`backend.yml`/`frontend.yml` no lo cubren). Propuesta: dejarlo local-first y abrir tarea aparte si se quiere CI. |
| OQ-3 | ¿Navegadores? | **Solo Chromium** al principio (instalación más ligera). Firefox/WebKit se añaden solo si aparece un bug específico. |
| OQ-4 | ¿Un solo spec o varios? | **Un spec** (`critical-flow.spec.ts`). Partirlo cuando haya ≥ 3 flujos independientes. |
