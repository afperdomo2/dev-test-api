# Mobile architecture

## Stack

- **React Native + Expo** (SDK 57, TypeScript)
- **Metro bundler** (monorepo-aware via `metro.config.js`)
- **Shared package**: `@devtest/shared` (workspace)

## Directory ownership

| Path | Purpose | Editable? |
|------|---------|-----------|
| `apps/mobile/App.tsx` | App entrypoint (minimal scaffold) | Yes |
| `apps/mobile/src/` | Feature screens, components, navigation, API client | Yes — add features here |
| `apps/mobile/metro.config.js` | Monorepo config (watchFolders, nodeModulesPaths) | Rarely |
| `apps/mobile/app.json` | Expo config (splash, icons, scheme) | Yes |
| `apps/mobile/tsconfig.json` | Extends `expo/tsconfig.base`, strict | Rarely |

## Data flow

```
Screen → apiClient (fetch + SecureStore JWT) → backend
```

| Layer | Depends on |
|-------|-----------|
| `src/api/client.ts` | `expo-secure-store` (token), `@devtest/shared` (ApiError) |
| `src/api/services/*.ts` | `src/api/client.ts`, `@devtest/shared` (DTOs) |
| Screens | Services + `@devtest/shared` (format, validators) |

## Shared package usage

All DTOs, type definitions, format utilities, and validation predicates come from `@devtest/shared`:

```ts
import type { User, Session } from '@devtest/shared'
import { formatDateTime, isValidEmail } from '@devtest/shared'
```

**Never** duplicate types or utilities in `apps/mobile/src/` that exist in shared. If something is reusable across web and mobile, it belongs in shared.

## Token storage

Use `expo-secure-store` (not `localStorage`):

```ts
import * as SecureStore from 'expo-secure-store'

SecureStore.setItemAsync('auth_token', token)
SecureStore.getItemAsync('auth_token')
SecureStore.deleteItemAsync('auth_token')
```

## HTTP client

The mobile app uses `fetch` (not axios). The client must:

1. Attach `Authorization: Bearer <token>` from SecureStore
2. Unwrap the `{ data: ... }` envelope (same contract as web)
3. Transform RFC 9457 errors into `ApiError` (from `@devtest/shared`)
4. On 401: delete token and redirect to login

The envelope/error contract is defined in `.agents/backend/responses.md` — the mobile client implements the consuming side.

## Monorepo

`metro.config.js` adds the monorepo root to `watchFolders` and `nodeModulesPaths` so Metro resolves `@devtest/shared` from the workspace symlink. No build step required — Metro transpiles `.ts` source directly.

## Commands

| Command | Action |
|---------|--------|
| `pnpm --filter mobile start` | Start Expo dev server |
| `pnpm --filter mobile typecheck` | `tsc --noEmit` |
| `pnpm --filter mobile android` | Start on Android |
| `pnpm --filter mobile ios` | Start on iOS |
