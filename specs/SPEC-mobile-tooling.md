# SPEC — mobile-tooling

| Campo | Valor |
|---|---|
| Module id | `mobile-tooling` |
| Origen | Alcance nuevo mobile (D-26, D-34) |
| Estado | `ready` |
| Depende de | — (puede ir en la Ola 1, en paralelo) |
| Consumido por | todos los módulos `mobile-*` (les da red de seguridad) |
| Migración goose | Ninguna |
| Dependencia nueva | `eslint` + `eslint-config-expo` (dev) — **pasa por `.agents/mobile/dependencies.md`** |

## Objective

Tener **lint** y una base real de **tests** en mobile. Hoy solo hay dos tests de utilidades; el resto del código (cliente HTTP, servicios, pantallas) no tiene red. El objetivo es que los cambios de los próximos módulos no rompan lo que ya funciona (D-34).

**Historias**:
- Como desarrollador, si rompo el unwrap del envelope o el manejo de errores del cliente, un test falla.
- Como desarrollador, si rompo el flujo de responder una pregunta, un test falla.
- Como desarrollador, `pnpm --filter mobile lint` me dice si el código cumple el estilo (hoy es un `echo 'TODO'`).

**Éxito**: `lint` y `test:run` en verde, cubriendo **servicios y pantallas críticas existentes** — no solo lo nuevo.

## Estado actual (verificado, N-08)

```
src/core/utils/dates.test.ts      ✅ existe
src/core/utils/errors.test.ts     ✅ existe
resto del código                  ❌ sin tests
lint                              ❌ "echo 'TODO: add eslint'"
infra                             ✅ jest-expo + @testing-library/react-native v14 + @types/jest
```

## Tech stack

Expo SDK 57 · Jest (`jest-expo`) · `@testing-library/react-native` v14 · ESLint (`eslint-config-expo`).

## Commands

```bash
pnpm --filter mobile test:run        # todos los tests
pnpm --filter mobile test:related    # solo lo afectado por archivos cambiados (hooks)
pnpm --filter mobile lint            # NUEVO
pnpm --filter mobile typecheck
```

## Project structure

```
apps/mobile/
├── eslint.config.js               # NUEVO (flat config, expo)
├── jest.config.js                 # NUEVO (o ampliar el bloque "jest" de package.json)
├── jest.setup.ts                  # NUEVO (matchers de RTL + mocks globales)
├── package.json                   # + lint, devDeps
└── src/
    ├── core/api/client.test.ts                    # NUEVO
    ├── core/api/services/auth.service.test.ts     # NUEVO
    ├── core/api/services/sessions.service.test.ts # NUEVO
    ├── core/api/services/questions.service.test.ts# NUEVO
    ├── core/api/services/progress.service.test.ts # NUEVO
    ├── core/api/services/gamification.service.test.ts # NUEVO
    ├── core/auth/storage.test.ts                  # NUEVO (SecureStore mockeado)
    ├── core/utils/errors.test.ts                  # ya existe (ampliar si falta)
    ├── core/utils/dates.test.ts                   # ya existe
    ├── features/home/screens/HomeScreen.test.tsx  # NUEVO
    ├── features/sessions/components/QuestionCard.test.tsx # NUEVO
    ├── features/sessions/screens/SessionStudyPage.test.tsx # NUEVO
    └── features/auth/screens/*.test.tsx           # NUEVO (login/setup)
```

## Contrato de los tests (qué se cubre, no cómo)

### `client.ts` — prioridad máxima

| Caso | Aserción |
|---|---|
| Token presente | `Authorization: Bearer <token>` en el `fetch` |
| Sin token | no envía la cabecera |
| Envelope `{ data }` | devuelve `data` desempaquetado |
| Envelope con `meta` | **no** desempaqueta (listas paginadas) |
| Error RFC 9457 | el error lanzado conserva `title/status/detail/instance` |
| Error sin JSON | construye un `ApiError` de fallback |
| 401 | borra el token y dispara `onUnauthorized` |
| Fallo de red (`fetch` rechaza) | lanza `{ title: 'Network Error', status: 500 }` — **contrato que usará `mobile-offline`** |
| Query params | `undefined` se omite; el resto se serializa |

### Servicios

Cada servicio: URL, método y payload correctos (mockeando `apiClient`). Especial atención a:
- `sessions.service.submitAnswer` (lo usa el flujo que rompe `mobile-offline`).
- `progress.service` (base de `mobile-progress`).
- `auth.service` (login/setup/status).

### Pantallas críticas

| Pantalla | Casos mínimos |
|---|---|
| `HomeScreen` | pinta racha y meta; el CTA llama `createDailySession` y navega; "usar freeze" pide confirmación con `Alert.alert`; error de red muestra mensaje |
| `QuestionCard` | selección única vs múltiple; deshabilitado al enviar; feedback con texto (no solo color) |
| `SessionStudyPage` | responder llama `submitAnswer`; avanza a resultado; **hoy un fallo de red se traga** (se documenta como comportamiento actual y se ajustará en `mobile-offline`) |
| `LoginScreen` / `SetupScreen` | validación básica y envío |

## Code style

- Tests **co-located**, nombre `*.test.ts(x)`, en español los `describe`/`it` (consistente con web).
- Sin `waitForTimeout` arbitrarios: `findBy*` / `waitFor`.
- Mocks de módulos nativos (`expo-secure-store`, `expo-router`, `@expo/vector-icons`) en `jest.setup.ts`, no repetidos por archivo.
- ESLint flat config con `eslint-config-expo`; reglas extra alineadas con web donde aplique (`no-unused-vars`, `eqeqeq`).

### ⚠️ RTL v14: `render` y `fireEvent` son **async**

`@testing-library/react-native` v14 cambió la API (verificado al implementar):

```ts
await render(<Component />)          // devuelve el RenderResult
await fireEvent.press(screen.getByLabelText('…'))
await fireEvent.changeText(screen.getByTestId('login-email'), 'x')
```

- `screen.*` solo es válido **después** de `await render`.
- Los imports no necesitan `extend-expect`: la librería registra sus matchers al importarse.
- No existen las queries `UNSAFE_*` en v14: para localizar `TextInput` se usa `testID` (por eso `AuthInput` lo acepta).
- `render`/`fireEvent` son lentos en frío → `testTimeout: 30000` en `jest.config.js`.

## Testing strategy

No se persigue un número de cobertura global. Se exige:
1. **Todos** los servicios y el cliente HTTP con tests.
2. **Las 4 pantallas críticas** de la tabla.
3. Los tests existentes (`dates`, `errors`) siguen verdes.

## Boundaries

- **Always**: co-located tests; mocks nativos centralizados; `typecheck` + `test:run` + `lint` en verde; pasar por `.agents/mobile/dependencies.md` antes de instalar ESLint.
- **Ask first**: añadir `@testing-library/jest-native` (deprecado), cambiar el preset de Jest, activar umbrales de cobertura en CI, añadir tests E2E mobile (Detox/Maestro — fuera de alcance).
- **Never**: tests que dependan de red real o de un dispositivo; `jest --watch` en hooks; desactivar reglas de ESLint globalmente para "que pase"; duplicar helpers de test en cada archivo.

## Success Criteria

1. `pnpm --filter mobile lint` existe y pasa (deja de ser un `echo`).
2. `pnpm --filter mobile test:run` pasa con **≥ 30 tests** y **≥ 12 archivos** (cobertura de cliente + 5 servicios + 4 pantallas).
3. Los tests del cliente cubren explícitamente los 9 casos de la tabla, incluido `Network Error`.
4. Un cambio deliberado en `unwrap()` de `client.ts` hace fallar al menos un test.
5. Un cambio deliberado en `handleAnswer` de `SessionStudyPage` hace fallar al menos un test.
6. `pnpm --filter mobile typecheck` sigue en verde.
7. `eslint` es devDependency y el gate de dependencias quedó documentado.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿`jest.config.js` o el bloque `"jest"` de `package.json`? | **`jest.config.js`**: permite `setupFilesAfterEnv` y `transformIgnorePatterns` sin ensuciar `package.json`. |
| OQ-2 | ¿Cobertura mínima en CI? | **No** por ahora: exigir suites concretas es más útil que un % global. Se puede añadir después. |
| OQ-3 | ¿Tests de los hooks (si aparecen)? | Sí, cuando existan; hoy no hay `src/core/hooks/`. |
