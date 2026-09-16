# Web architecture

## Directory ownership

| Path | Purpose | Editable? |
|------|---------|-----------|
| `apps/web/src/api/client.ts` | Axios instance + interceptors (JWT, envelope unwrap, RFC 9457 errors) | Rarely |
| `apps/web/src/api/services/*.service.ts` | Pure HTTP functions per domain | Yes — add new endpoints here |
| `apps/web/src/queries/*.queries.ts` | TanStack Query: `queryOptions()` + mutation definitions | Yes — per domain |
| `apps/web/src/features/<domain>/` | Business domains: pages + components | Yes — add features here |
| `apps/web/src/stores/` | Pinia stores (auth, app) | Yes — add global state |
| `apps/web/src/router/` | Vue Router: routes.ts + index.ts with guards | Yes — add routes here |
| `apps/web/src/composables/` | Reusable composables (useDebounce, usePagination, useFormErrors) | Yes |
| `apps/web/src/types/` | Thin re-exports from `@devtest/shared` | Yes — add new re-exports when shared adds types |
| `apps/web/src/utils/` | `storage.ts` (localStorage), `rules.ts` (Vuetify validation) | Yes |
| `apps/web/src/utils/format.ts` | Re-export from `@devtest/shared` | Keep in sync with shared |
| `apps/web/src/utils/validators.ts` | Vuetify rule factories (import predicates from `@devtest/shared`) | Yes |
| `apps/web/src/components/` | Shared components (ErrorState, ListPageHeader, PaginatedFooter) | Yes — add reusable UI here |
| `apps/web/src/constants/` | Re-exports from `@devtest/shared` | Keep in sync with shared |
| `apps/web/src/plugins/` | Vuetify + VueQueryPlugin setup | Rarely |
| `apps/web/dist/` | Build output | **NEVER** — use `pnpm build` |

## Shared package

`types/*`, `constants/*`, `utils/format.ts`, and `utils/validators.ts` (predicates) live in `@devtest/shared` (`packages/shared/`). The files in web are **thin re-exports** that preserve existing `@/types/...` imports.

When adding a new type or constant:
1. Add the source definition to `packages/shared/src/types/` or `packages/shared/src/constants/`
2. Add the re-export to `apps/web/src/types/<name>.types.ts` or `apps/web/src/constants/index.ts`

Full details: see `.agents/shared/architecture.md`.

## Data flow

```
Page → useQuery(queries.*Options()) → service function → apiClient → backend
```

Three layers, strictly separated:

| Layer | Directory | Depends on |
|-------|-----------|------------|
| API services | `api/services/` | `api/client.ts` (Axios) |
| Queries | `queries/` | `api/services/` (TanStack Query wraps them) |
| Pages/components | `features/` | `queries/`, `stores/`, `composables/` |

`api/services/` must NOT import TanStack Query. `queries/` must NOT import Vue components.

## API contract

The envelope `{ data: ... }` and RFC 9457 error format are defined by the backend. See `.agents/backend/responses.md` for the canonical contract. The web `api/client.ts` implements the consuming side (Axios interceptors that unwrap the envelope and transform errors).

## Feature structure

Each feature domain in `features/<name>/` follows this pattern:

```
features/<name>/
├── pages/
│   └── <Name>Page.vue      # Full page: useQuery/useMutation, layout, routing
└── components/
    └── <Name>Component.vue  # Reusable: props + emits, no routing logic
```

Pages own TanStack Query calls. Components are presentational — they receive data via props and emit events. Components never import query/mutation definitions directly.

## Wiring: new feature checklist

1. Create types in `types/<name>.types.ts` (add constant maps if needed: labels, colors, options arrays) — re-export from `@devtest/shared` when the DTO comes from the shared package
2. Create service in `api/services/<name>.service.ts` (pure functions, return typed promises)
3. Create queries in `queries/<name>.queries.ts` (export `queryOptions()` for reads, mutation objects for writes)
4. Create pages + components in `features/<name>/`
5. Register routes in `router/routes.ts` (add `meta: { requiresAuth, requiresAdmin? }`)
6. Run `make fe-check` and `make fe-lint`
