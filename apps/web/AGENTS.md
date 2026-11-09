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

- Contracto API (envelope `{ data }`, errores RFC 9457, camelCase): `.agents/context/backend/responses.md`.
- **Identidad visible**: avatar en la app bar, email + rol en el sidebar (el bloque es el link a "Mi perfil") y `document.title` = `Rol · email — APP_NAME`. El rol se lee sin abrir menús. Detalle: `.agents/context/web/patterns.md` → "Identidad del usuario logueado".
- Design system UI: Vuetify (componentes + `VITE_*` env). Ver `.agents/context/web/architecture.md`.
- Identidad de build: `VITE_APP_NAME` (default `DevTest`) y `VITE_APP_COLOR` (presets en `src/config/app.ts`); se compilan en el bundle y exigen rebuild.
- **Iconos/PWA**: favicons y manifest en `apps/web/public/` (`favicon.svg|ico`, `apple-touch-icon.png`, `icon-192|512.png`, `maskable-512.png`, `site.webmanifest`). Son generados: no se editan a mano, se regeneran con `make icons` desde `brand/`. Detalle: `docs/branding.md`.

## Details

- Arquitectura: `.agents/context/web/architecture.md`
- API client: `.agents/context/web/api-client.md`
- Patrones: `.agents/context/web/patterns.md`
- Seguridad de dependencias: `.agents/context/web/dependencies.md`