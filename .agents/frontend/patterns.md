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

Constants in `constants/index.ts`:
- `ITEMS_PER_PAGE_OPTIONS = [10, 20, 50]`
- `DEFAULT_PER_PAGE = 20`

Use `usePagination()` composable in list pages:

```ts
import { usePagination } from '@/composables/usePagination'

const { page, perPage, reset } = usePagination()
// perPage defaults to DEFAULT_PER_PAGE (20) from constants
```

`usePagination` returns reactive refs. Pass `page` and `perPage` to `queryOptions()`. Use `reset()` when filters change or on refresh.

## Shared list components

Two standardized components live in `frontend/src/components/`. Every list page MUST use them.

### `ListPageHeader.vue`

Props: `title` (string), `createLabel` (string), `createTo?` (route), `showCreate?` (boolean, default true).
Emits: `refresh`, `create`.

Renders: page title on left, refresh button + create button on right. Refresh button has a 1-second cooldown (`REFRESH_COOLDOWN_MS` from constants). Show/hide create button with `showCreate` (condition on `authStore.isAdmin` for admin-only create).

### `PaginatedFooter.vue`

Props: `page`, `perPage`, `total` (all numbers). Emits: `update:page`, `update:perPage`.

Renders: per-page selector (10/20/50) + "Mostrando X–Y de Z" on left, pagination with `show-first-last-page` (◀◀ ◀ pages ▶ ▶▶) on right. Only renders pagination controls when `total > perPage` (handled internally).

## List page standard pattern

Every list page follows this structure:

```
<template>
  <v-container>
    <ListPageHeader />
    <!-- filters (optional) -->
    <!-- skeleton loaders (loading) -->
    <!-- data table OR card grid (data) -->
    <!-- empty state card (no data) -->
    <PaginatedFooter />
    <!-- dialogs (create, delete) -->
  </v-container>
</template>
```

### For data table pages (Users, Topics)

Use `v-data-table` with a `#bottom` slot containing `<PaginatedFooter>`:

```html
<v-data-table ...>
  <template #[`item.xxx`]="{ item }">...</template>
  <template #bottom>
    <PaginatedFooter
      :page="page" :per-page="perPage" :total="totalX"
      @update:page="page = $event"
      @update:per-page="perPage = $event"
    />
  </template>
</v-data-table>
```

The `#bottom` slot replaces the built-in data table footer.

### For card-based pages (Questions, Sessions, Progress)

Place `<PaginatedFooter>` after the card grid:

```html
<v-row v-if="items.length">...</v-row>
<v-card v-else><!-- empty state --></v-card>
<PaginatedFooter
  :page="page" :per-page="perPage" :total="totalX"
  class="mt-4"
  @update:page="page = $event"
  @update:per-page="perPage = $event"
/>
```

### Empty state standard

Empty states show ONLY text + icon. No inline create buttons (the create button is already in `ListPageHeader`):

```html
<v-card v-else>
  <v-card-text class="text-center py-8">
    <v-icon size="48" color="grey-lighten-1" class="mb-2"> mdi-xxx </v-icon>
    <p class="text-body-1 text-medium-emphasis">No hay X aún</p>
  </v-card-text>
</v-card>
```

### Refresh pattern

Every list page implements `handleRefresh()` as:

```ts
function handleRefresh() {
  resetPagination()                        // back to page 1
  queryClient.invalidateQueries({ queryKey: ['domain', 'list'] })  // bust cache
  refetch()                                // force immediate refetch
}
```

Cooldown (1s) is handled internally by `ListPageHeader`. No extra logic needed in the page.

### Full wiring example (new module)

```ts
// features/xxx/pages/XxxListPage.vue
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import { usePagination } from '@/composables/usePagination'
import { useAuthStore } from '@/stores/auth.store'

const authStore = useAuthStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const { data, isLoading, refetch } = useQuery(xxxListOptions(() => page.value, () => perPage.value))

const items = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['xxx', 'list'] })
  refetch()
}
```

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
