# Web (Vue 3) — quick context

Stack: **Vue 3 + Vite + Vuetify 3** + **Pinia** + **TanStack Query** + **TypeScript** + `@devtest/shared`.

## Commands

| Command | Action |
|---------|--------|
| `make fe-dev` | Vite dev server en `:3000` (proxy `/api` → `:8080`) |
| `make fe-check` | `vue-tsc --build` (tras cada cambio) |
| `make fe-lint` | ESLint (`--fix`) |
| `make fe-build` | Build producción (solo CI/deploy) |

## Structure

```
src/api/services/    # Funciones HTTP por dominio
src/queries/         # TanStack Query (queryOptions + mutations)
src/features/<domain>/  # pages/ (usan queries) + components/ (presentacionales)
src/stores/          # Pinia (auth, app)
src/router/          # rutas + guards
src/types/           # re-exports finos de @devtest/shared
src/composables/     # useDebounce, usePagination, useFormErrors
src/components/      # componentes compartidos
```

## Data flow

```
Page → useQuery(queries.*Options()) → service → apiClient → backend
```

`api/services/` NO importa TanStack Query. `queries/` NO importa componentes Vue.

## Key facts

- Contracto API (envelope `{ data }`, errores RFC 9457, camelCase): `.agents/backend/responses.md`.
- Design system UI: Vuetify (componentes + `VITE_*` env). Ver `.agents/web/architecture.md`.

## Details

- Arquitectura: `.agents/web/architecture.md`
- API client: `.agents/web/api-client.md`
- Patrones: `.agents/web/patterns.md`
- Seguridad de dependencias: `.agents/web/dependencies.md`