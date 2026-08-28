# TASKS — Frontend Tests (Vitest)

> Stack: Vue 3 + Vite 6 + Vuetify 3 + Pinia + Vue Router + TanStack Query
> Runner: Vitest + jsdom + @vue/test-utils + @testing-library/vue
> Ver plan detallado en conversación previa.

## Estado general
- [x] Fase 0 — Infraestructura Vitest
- [ ] Fase 1 — Tests puros (utils / composables / stores / api client)
- [ ] Fase 2 — Services + Queries
- [ ] Fase 3 — Componentes base + features
- [ ] Fase 4 — Router + layouts + integración
- [ ] Fase 5 — CI / Makefile / DX + E2E (opcional)

---

### Fase 0 — Infraestructura Vitest
- [x] **0.1** Instalar devDeps: `vitest`, `jsdom`, `@vue/test-utils`, `@testing-library/vue`, `@testing-library/jest-dom`, `msw`/`axios-mock-adapter`, `@pinia/testing`, `@vitest/coverage-v8`, `jsdom` stubs
- [x] **0.2** Configurar `frontend/vite.config.ts` — bloque `test` (environment `jsdom`, `globals:true`, `include`, `setupFiles`, `coverage` v8, alias `@`)
- [x] **0.3** Crear `frontend/src/__tests__/setup.ts` (jest-dom, Vuetify plugin global, `ResizeObserver`/`matchMedia` stubs, `vuetify/styles` mock si hace falta)
- [x] **0.4** Ajustar `frontend/env.d.ts` / `tsconfig.json` para tipos `vitest/globals` + `jsdom`
- [x] **0.5** Añadir scripts `frontend/package.json`: `test`, `test:run`, `test:cover`, `test:ui`
- [x] **0.6** Verificar `pnpm test:run` pasa en vacío (smoke test dummy) + `pnpm type-check` verde

### Fase 1 — Tests puros (sin mount, alto ROI)
- [ ] **1.1** `src/utils/validators.test.ts` — `isValidEmail`, `isValidPassword`, `isRequired`, `isMinLength/isMaxLength`, `requiredRule/emailRule/passwordRule/confirmPasswordRule/maxLengthRule`, `validateRules`
- [ ] **1.2** `src/utils/format.test.ts` — `formatDateTime`, `formatDate`, `formatScore`
- [ ] **1.3** `src/utils/storage.test.ts` — `getToken/setToken/removeToken` con `localStorage` mock + `TOKEN_KEY`
- [ ] **1.4** `src/composables/useFormErrors.test.ts` — `extractFieldErrors` (split `;` `,` `:`, lowercase keys, detail vacío, error sin detail)
- [ ] **1.5** `src/composables/useDebounce.test.ts` — valor inicial, debounce 500ms con `vi.useFakeTimers()`, cambio rápido resetea timer, delay custom
- [ ] **1.6** `src/composables/usePagination.test.ts` — estado inicial, cambios de page
- [ ] **1.7** `src/stores/auth.store.test.ts` — `setSession/clearSession/setUser`, `isLoggedIn/isAdmin`, `initSession` éxito/401, `checkStatus` initialized true/false/error (mock `users.service` + `auth.service` + `storage`)
- [ ] **1.8** `src/stores/app.store.test.ts` — estado y acciones del app store
- [ ] **1.9** `src/api/client.test.ts` — interceptor request añade `Bearer`, no añade sin token, unwrap `{data}` sin `meta`, no unwrap con `meta`, 401 → `removeToken`, `ApiError` fallback 500, rechazo `Promise.reject`

### Fase 2 — Services + Queries (mock apiClient)
- [ ] **2.1** `src/api/services/auth.service.test.ts` — `login`/`setup`/`getStatus` llaman `POST/GET` a `/api/v1/auth/*` y retornan `res.data`
- [ ] **2.2** `src/api/services/users.service.test.ts` + `topics` + `sessions` + `questions` + `progress` — CRUD y query params
- [ ] **2.3** `src/queries/auth.queries.test.ts` — `loginMutation`/`setupMutation` keys y `mutationFn` delega a service
- [ ] **2.4** `src/queries/users.queries.test.ts` + `topics` + `sessions` + `questions` + `progress` — `queryKey`, `queryFn`, `enabled` conditional

### Fase 3 — Componentes (con mount)
- [ ] **3.1** Helper `src/__tests__/utils/mount.ts` — `mountWithProviders` (Pinia + Vuetify + Router + Vue Query, stubs `RouterLink`)
- [ ] **3.2** `src/components/ErrorState.test.ts` + `PaginatedFooter.test.ts` + `ListPageHeader.test.ts` + `CodeContent.test.ts`
- [ ] **3.3** `src/features/auth/components/LoginForm.test.ts` — validación, submit, loading, `extractFieldErrors` integration
- [ ] **3.4** `src/features/users/components/UserTable.test.ts` + `UserFormDialog.test.ts`
- [ ] **3.5** `src/features/questions/components/QuestionTable.test.ts` + `QuestionFormDialog.test.ts` + `QuestionFilters.test.ts`
- [ ] **3.6** `src/features/topics/components/TopicFormDialog.test.ts`
- [ ] **3.7** `src/directives/highlight.test.ts`

### Fase 4 — Router + layouts + integración
- [ ] **4.1** `src/router/index.test.ts` — guards: `needsSetup→/setup`, `needsSetup false + /setup→/login`, `requiresAuth→/login?redirect`, `requiresAdmin`, `requiresNotAdmin`, `already logged + /login→/` (tabla `it.each`)
- [ ] **4.2** `src/layouts/DefaultLayout.test.ts` + `AuthLayout.test.ts`
- [ ] **4.3** Smoke `src/features/auth/pages/LoginPage.test.ts` + `SetupPage.test.ts` + `src/features/dashboard/pages/DashboardPage.test.ts`

### Fase 5 — CI / DX / E2E opcional
- [ ] **5.1** Makefile raíz — targets `fe-test` (`pnpm test:run`), `fe-test-cover`, actualizar `fe-check` si hace falta
- [ ] **5.2** `.github/workflows/frontend.yml` — step `Test` con `pnpm test:run --coverage`
- [ ] **5.3** `lefthook.yml` — hook pre-commit/push `pnpm test:run` (o `test:cover` con threshold)
- [ ] **5.4** Documentar convención en `.agents/frontend/patterns.md` (co-localizado `.test.ts`, `globals:true`, `vi.mock` vs `msw`)
- [ ] **5.5** (Opcional) Playwright E2E — flujo `status → setup → login → dashboard → sessions` con `msw` o backend real

---

## Flujo de trabajo
1. Marcar tarea como `in_progress` antes de empezar.
2. Implementar + verificar con `pnpm test:run --coverage` y `pnpm type-check`.
3. Marcar como completada en este archivo y continuar con la siguiente.
4. Al terminar todas, crear commit manual (no automático).

## Notas
- Convención: tests co-localizados `*.test.ts` junto al source (elegido por tamaño del repo).
- `globals: true` para no importar `describe/it/expect` en cada archivo.
- `jsdom` como environment; fallback `happy-dom` si Vuetify da problemas.
- Coverage provider `v8`, threshold inicial sugerido `60%` líneas.
