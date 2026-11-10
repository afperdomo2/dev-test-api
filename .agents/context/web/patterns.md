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

## Password fields

`PasswordField` (`features/auth/components/PasswordField.vue`) encapsula un `v-text-field` de contraseña y es el patrón a reutilizar (login ya lo usa):

- Arranca con `type="password"`; el botón del `append-inner` es un `v-btn` nativo (`type="button"`, no envía el form) que alterna a `type="text"`.
- El botón lleva `aria-label` dinámico («Mostrar/Ocultar contraseña»): nunca se depende solo del icono.
- Al deshabilitarse el campo (envío en curso) vuelve a ocultarse; la visibilidad no se persiste.
- Autofill: login → `autocomplete="current-password"` (email con `autocomplete="username"`); alta → `new-password`.

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

Constants in `packages/shared/src/constants/index.ts` (re-exported via `@/constants`):
- `ITEMS_PER_PAGE_OPTIONS = [10, 20, 50]`
- `DEFAULT_PER_PAGE = 10`

Use `usePagination()` composable in list pages:

```ts
import { usePagination } from '@/composables/usePagination'

const { page, perPage, reset } = usePagination()
// perPage defaults to DEFAULT_PER_PAGE (10) from shared constants
```

`usePagination` returns reactive refs. Pass `page` and `perPage` to `queryOptions()`. Use `reset()` when filters change or on refresh.

## Shared list components

Two standardized components live in `apps/web/src/components/`. Every list page MUST use them.

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

### For card-based pages (Questions, Progress)

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

### For the sessions list (hybrid row list)

Sessions use a **row list, not a card grid and not a `v-data-table`**: one row per
session (~68px), aligned columns so the eye can scan state → name → progress →
topics → action. Reason: sessions are an actionable history (continue / review),
not a set of comparable scalar values.

```
SessionRow.vue   (presentational: props + emits, no queries/mutations)
  props  session, busy
  emits  open | finish | delete

SessionsListPage.vue
  useQuery(sessionsListOptions(page, perPage, status, sortBy, sortOrder))
  useQuery(sessionsResumeOptions())        → banner "Continuar"
  mutations (finish/delete) live here      → single confirm dialog
  usePagination() + PaginatedFooter :in-table="true" inside the card
```

- Sorting is **server-side**: `sortBy`/`sortOrder` go through the service and into
  the `queryKey`, never a client-side `computed().sort()`.
- The CTA is a real `v-btn :to` and the name is a `RouterLink` — don't make the
  row itself a focusable button (nested interactive elements).
- Row click (`emit('open')`) is only a mouse convenience on top of those links.

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

## Identidad del usuario logueado

El usuario debe poder identificarse en cualquier página **sin abrir menús ni entrar a "Mi perfil"** (necesario cuando hay muchas pestañas abiertas de cuentas distintas).

- **App bar** (`DefaultLayout.vue`): botón de cuenta con **solo el avatar** (inicial). Abre el menú con "Mi perfil" y "Cerrar sesión". El email y el rol no se muestran en la barra ni se repiten en el menú: viven en el sidebar, en el título de la pestaña y en el `aria-label`/`title` del botón.
- **Sidebar** (`#append`): el mismo `UserIdentity` (`variant="sidebar"`) es el **link a "Mi perfil"**: un `v-list-item :to="/profile"` con chevron, hover, estado activo y `aria-label` (`Mi perfil de <email> (<rol>)`). No hay un item "Mi perfil" separado; "Cerrar sesión" queda como acción aparte. El drawer es `permanent`, así que se ve siempre en desktop.
- **Título de pestaña**: `useDocumentTitle()` + `buildDocumentTitle()` → `Rol · email — APP_NAME` (p. ej. `Admin · admin@empresa.com — DevTest`). El rol va **primero** porque el navegador trunca el título por la derecha. Sin usuario → `APP_NAME`. `AuthLayout` fija `APP_NAME` para que al cerrar sesión no quede el usuario anterior.
- **Etiquetas de rol**: `roleLabel()` / `shortRoleLabel()` en `utils/identity.ts`. No hardcodear `'Administrador'` / `'Usuario'` en componentes.
- **Rol = icono + color + texto** (nunca solo color): en el sidebar el chip va `tonal` con `mdi-shield-account` + `error` (admin) o `mdi-account` + `info` (usuario). En la app bar el avatar va blanco con la inicial en `primary` para contrastar con cualquier `VITE_APP_COLOR`.
- **El avatar** está siempre visible en la app bar; el email y el rol se ven en el sidebar (desktop) y en el título de la pestaña.

## Styling

- Vuetify theming: `plugins/vuetify.ts` (`defaultTheme` se inicializa desde el tema persistido; paletas light/dark)
- Theme toggle: `stores/app.store.ts` → `toggleTheme()`. El tema elegido se guarda **por dispositivo** en `localStorage` (`app_theme`) con `getTheme()`/`setTheme()` de `utils/storage.ts`, así sobrevive a un F5 o a reiniciar el dev server. Valores inválidos caen a `light`.
- Font: Roboto via `unplugin-fonts` (configured in `vite.config.ts`)
- Icons: `@mdi/font` (Material Design Icons, imported in `plugins/vuetify.ts`)
- SASS overrides: `apps/web/src/styles/settings.scss`

## Label + `v-progress-linear` rows

Rows that pair a label with a horizontal bar (stats, dashboards, topic lists) must align the bars in a single column:

- Label column with a **fixed `flex-basis`** (`flex: 0 1 220px; min-width: 120px`) — identical in every row, so all bars start at the same X. `shrink` + `min-width` keep the bar usable on mobile.
- Label text with `text-truncate` (ellipsis) **plus `:title`** for the full value. Prefer a single line; don't stack metadata (e.g. the topic's category) under the name — it doubles the row height and eats the width the bar needs.
- Counter **outside** the bar, in a fixed-width right column (`min-width: 32px; text-align: right`): a short bar can't hold a centered number.
- Bar with explicit `class="flex-grow-1"`.

Never use `text-no-wrap` + `max-width` for the label: the text isn't clipped, it overflows and paints over the bar, and the variable width makes the same percentage look like different bar lengths.

## Delete confirmation rule

Any destructive action (delete, bulk remove) MUST show a confirmation dialog before executing. The dialog must:

- Clearly state what will be deleted (`"¿Estás seguro de eliminar <name>?"`)
- Have a **Cancel** button (`variant="text"`) and a **Confirm** button (`color="error"`, loading state)
- Only execute the mutation after the user explicitly confirms
- Show a success/error snackbar after the operation

```html
<v-dialog v-model="deleteDialog" max-width="420">
  <v-card>
    <v-card-title>Eliminar X</v-card-title>
    <v-card-text>
      ¿Estás seguro de eliminar <strong>{{ target?.name }}</strong>?
    </v-card-text>
    <v-card-actions>
      <v-spacer />
      <v-btn variant="text" @click="deleteDialog = false"> Cancelar </v-btn>
      <v-btn color="error" :loading="deleteMut.isPending.value" @click="executeDelete">
        Eliminar
      </v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>
```

This pattern applies to all list pages (Users, Topics, Sessions, etc.).

## v-slot with dots in Vuetify data table

Vuetify `v-data-table` uses dotted slot names like `item.is_admin`. ESLint's `vue/valid-v-slot` treats dots as modifiers. Use bracket syntax:

```html
<!-- ❌ ESLint error -->
<template #item.is_admin="{ item }">

<!-- ✅ Correct -->
<template #[`item.isAdmin`]="{ item }">
```

## Frontend tests (Vitest)

### Stack

- **Runner**: Vitest 3 + jsdom + `@vue/test-utils` + `@testing-library/vue` + `@testing-library/jest-dom`
- **Config**: `vite.config.ts` `test` block (`environment: jsdom`, `globals: true`, `setupFiles: ['src/__tests__/setup.ts']`, `coverage: { provider: 'v8' }`, `server.deps.inline: [/vuetify/]`)
- **Setup**: `src/__tests__/setup.ts` registers Vuetify (`createVuetify({})`) globally, stubs `ResizeObserver`, `visualViewport`, `matchMedia`, `requestAnimationFrame`

### Scripts & Makefile

```bash
pnpm test        # watch
pnpm test:run    # single run
pnpm test:cover  # with coverage v8
make fe-test       # pnpm test:run
make fe-test-cover # pnpm test:cover
```

CI: `.github/workflows/frontend.yml` runs `pnpm test:run --coverage` after type-check.

Hooks (`lefthook.yml`), ambos filtrados por `glob` (no corren si no cambió `apps/web/**` ni `packages/shared/**/*.ts`):
- Pre-commit (`root: apps/web/`): `pnpm vitest related {staged_files}` → solo los tests que importan los archivos staged.
- Pre-push: `pnpm test:run` (suite completo).

### Conventions

- **Co-located** `*.test.ts` next to source (`src/utils/validators.test.ts`, `src/components/ErrorState.test.ts`). No `__tests__` centralization (except `src/__tests__/setup.ts` + `src/__tests__/utils/mount.ts` helpers + `smoke.test.ts`).
- **`globals: true`** — no need to import `describe/it/expect/vi` per file.
- **Deterministic**: `vi.useFakeTimers()` + `await nextTick()` for `useDebounce` (watch is async); `vi.clearAllMocks()` in `beforeEach`.

### Mock strategy

- **Pure utils/composables/stores**: no mount, mock `storage`/`services` via `vi.mock('@/utils/storage', ...)` with `vi.importActual` for partial mocks.
- **api/client**: `axios-mock-adapter` for interceptors (unwrap, 401 `removeToken`, `ApiError` fallback).
- **Services**: `vi.mock('@/api/client', () => ({ default: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }))`.
- **Queries**: `vi.mock('@/api/services/*.service')` and assert `mutationKey`/`queryKey` + `mutationFn`/`queryFn` delegations. For `queryFn` use `(opts as unknown as { queryFn: ... }).queryFn` to bypass `unique symbol` type.
- **Components**: `mount` from `@vue/test-utils` with global Vuetify (from `setup.ts`). Stub `teleport: true` except for `v-dialog` tests where we assert via `vm.dialogTitle`/`vm.isEdit`/`vm.validate()` instead of `wrapper.text()` (teleported content not in `wrapper.text()`). For `v-data-table`/`router-link`/`v-menu`, use `stubs: { RouterLink: true }` and assert via `vm.filteredNav` or `props`.

### Helpers

```ts
// src/__tests__/utils/mount.ts
import { mountWithProviders, createTestQueryClient } from '@/__tests__/utils/mount'

// Mount with Pinia + Vuetify + MemoryRouter + VueQuery (retry: false)
const wrapper = await mountWithProviders(MyComponent, {
  props: { modelValue: true },
  routerRoutes: [{ path: '/', component: { template: '<div/>' } }],
})
```

### Known stubs

- `ResizeObserver`, `visualViewport`, `matchMedia` stubbed in `setup.ts` for Vuetify `VOverlay`/`VDialog`.
- `v-highlight` directive: mock `@devtest/shared` (`getHighlighter`, `normalizeCodeLanguage`) — el registro de lenguajes vive ahí.

### Coverage

Provider `v8`, `include: ['src/**/*.{ts,vue}']`, `exclude: ['src/**/*.types.ts', 'src/main.ts', 'src/__tests__/**']`. Threshold not enforced yet (suggested 60% lines initial).
