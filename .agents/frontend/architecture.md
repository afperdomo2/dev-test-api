# Frontend architecture

## Directory ownership

| Path | Purpose | Editable? |
|------|---------|-----------|
| `frontend/src/api/client.ts` | Axios instance + interceptors (JWT, envelope unwrap, RFC 9457 errors) | Rarely |
| `frontend/src/api/services/*.service.ts` | Pure HTTP functions per domain | Yes — add new endpoints here |
| `frontend/src/queries/*.queries.ts` | TanStack Query: `queryOptions()` + mutation definitions | Yes — per domain |
| `frontend/src/features/<domain>/` | Business domains: pages + components | Yes — add features here |
| `frontend/src/stores/` | Pinia stores (auth, app) | Yes — add global state |
| `frontend/src/router/` | Vue Router: routes.ts + index.ts with guards | Yes — add routes here |
| `frontend/src/composables/` | Reusable composables (useDebounce, usePagination, useFormErrors) | Yes |
| `frontend/src/types/` | TypeScript DTOs + constant maps | Yes — add types per domain |
| `frontend/src/utils/` | Pure utilities (storage, format, validators) | Yes |
| `frontend/src/plugins/` | Vuetify + VueQueryPlugin setup | Rarely |
| `frontend/src/components/` | Shared components (ErrorState) | Yes |
| `frontend/dist/` | Build output | **NEVER** — use `pnpm build` |

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

1. Create types in `types/<name>.types.ts` (add constant maps if needed: labels, colors, options arrays)
2. Create service in `api/services/<name>.service.ts` (pure functions, return typed promises)
3. Create queries in `queries/<name>.queries.ts` (export `queryOptions()` for reads, mutation objects for writes)
4. Create pages + components in `features/<name>/`
5. Register routes in `router/routes.ts` (add `meta: { requiresAuth, requiresAdmin? }`)
6. Run `make fe-check` and `make fe-lint`
