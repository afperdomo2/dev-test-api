# Frontend patterns & conventions

## TanStack Query

### Query definitions (`queries/*.queries.ts`)

Export `queryOptions()` functions for reads and mutation objects for writes:

```ts
// queries/users.queries.ts
export function usersListOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['users', 'list', page, perPage],
    queryFn: () => usersService.listUsers(page(), perPage()),
    staleTime: 30 * 1000,  // 30s for lists
  })
}

export function createUserMutation() {
  return {
    mutationKey: ['users', 'create'],
    mutationFn: usersService.createUser,
  }
}
```

- **Query keys** follow `[domain, action, ...params]` pattern
- **staleTime defaults**: 60s for single items, 30s for lists, 0 for real-time data (e.g., next question in a session)
- **Mutation params**: use destructured objects for mutations with multiple args (e.g., `{ id, data }`)

### Page usage

```ts
// features/users/pages/UsersListPage.vue
const { page, perPage } = usePagination()
const { data, isLoading } = useQuery(usersListOptions(() => page.value, () => perPage.value))
const deleteMut = useMutation(deleteUserMutation())

// Access paginated data
const items = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)

// Invalidate after mutation
await deleteMut.mutateAsync(id)
queryClient.invalidateQueries({ queryKey: ['users', 'list'] })
```

Components that use the same `queryKey` share the cache automatically. No prop drilling needed for data.

### Reactivity in queryOptions

Query and mutation options receive **getter functions** (`() => value`, not raw `value` or `Ref`). This ensures TanStack Query re-evaluates when the underlying ref changes.

## Debounce

Use `useDebounce(value, 500)` for search inputs:

```ts
const { value: searchText, debouncedValue: debouncedSearch } = useDebounce('', 500)

watch(debouncedSearch, (val) => {
  // Trigger search/filter only after 500ms of inactivity
})
```

Returns `{ value: Ref<T>, debouncedValue: Ref<T> }`. Bind `value` to the input, react to `debouncedValue`.

## Vue component conventions

ESLint enforces these rules automatically:
- **Script setup**: `<script setup lang="ts">` always, no Options API
- **Block order**: `script` → `template` → `style`
- **Template casing**: PascalCase for components, hyphenated for attributes/props
- **Type imports**: `import type { X }` for type-only imports

## TypeScript rules

| Rule | Setting |
|------|---------|
| `no-explicit-any` | error |
| `no-non-null-assertion` | error |
| `consistent-type-imports` | error (prefer `type-imports`) |
| `array-type` | error (prefer `Array<T>`, not `T[]`) |

## Pinia stores

Use Composition API syntax (`defineStore('id', () => { ... })`):

```ts
// stores/auth.store.ts
export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const isLoggedIn = computed(() => !!token.value)
  return { token, isLoggedIn }
})
```

Stores are global (app-level or cross-feature). Feature-specific state belongs in `useQuery` cache or local component state.

## Pagination

Use `usePagination()` composable in list pages:

```ts
const { page, perPage, total, reset } = usePagination()
```

Returns reactive refs. Pass `page` and `perPage` to `queryOptions()`. Use `reset()` when filters change (not between pages of the same filtered view).

## Styling

- Vuetify theming: `plugins/vuetify.ts` (light/dark themes with color palette)
- Theme toggle: `stores/app.store.ts` → `toggleTheme()`
- Font: Roboto via `unplugin-fonts` (configured in `vite.config.ts`)
- Icons: `@mdi/font` (Material Design Icons, imported in `plugins/vuetify.ts`)
- SASS overrides: `frontend/src/styles/settings.scss`

## v-slot with dots in Vuetify data table

Vuetify `v-data-table` uses dotted slot names like `item.is_admin`. ESLint's `vue/valid-v-slot` treats dots as modifiers. Use bracket syntax:

```html
<!-- ❌ ESLint error -->
<template #item.is_admin="{ item }">

<!-- ✅ Correct -->
<template #[`item.is_admin`]="{ item }">
```
