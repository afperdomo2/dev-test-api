# Mobile architecture

## Stack

- **React Native + Expo** (SDK 57, TypeScript)
- **Expo Router** (file-based routing, `app/` directory)
- **NativeTabs** (`expo-router/unstable-native-tabs`) — native tab bar (iOS liquid glass / Android Material 3)
- **Metro bundler** (monorepo-aware via `metro.config.js`)
- **Shared package**: `@devtest/shared` (workspace)
- **Design system**: source of truth en `.agents/mobile/design-system.md` (generado con `ui-ux-pro-max`) → espejo tipado en `src/core/theme/tokens.ts`

## Directory ownership

| Path | Purpose | Editable? |
|------|---------|-----------|
| `app/` | Expo Router — **only** route files and `_layout.tsx` (thin re-exports). Never co-locate components/types/utils here | Yes — add routes |
| `app/_layout.tsx` | Root layout (Stack, headerShown false, theme) | Rarely |
| `app/(tabs)/_layout.tsx` | `NativeTabs` — 4 triggers: Inicio, Temas, Progreso, Perfil | Yes |
| `src/features/<domain>/` | Feature (domain) modules: `screens/`, `components/`, `hooks/` | Yes — add features here |
| `src/core/theme/tokens.ts` | Design tokens (colors, spacing, radii, typography) — mirror tipado de `.agents/mobile/design-system.md` | Yes — keep in sync with design-system.md |
| `src/core/components/` | Shared UI primitives (e.g. `Screen.tsx`) | Yes |
| `src/core/api/` | HTTP client + services (`client.ts`, `services/*.ts`) | Yes — add endpoints here |
| `src/core/hooks/` | Shared hooks (color scheme, etc.) | Yes |
| `src/core/utils/` | Shared utilities | Yes |
| `apps/mobile/package.json` | `"main": "expo-router/entry"` | Rarely |
| `apps/mobile/app.json` | Expo config (plugins, scheme, userInterfaceStyle: "dark") | Yes |
| `apps/mobile/tsconfig.json` | Extends `expo/tsconfig.base`, `@/*` → `src/*` | Rarely |
| `.agents/mobile/design-system.md` | Design system (Dark OLED, paleta, tipografía, anti-patrones) — source of truth | Yes |

## Feature structure

Feature-driven (domain-oriented), mirroring the web app:

```
src/features/<name>/
├── screens/
│   └── <Name>Screen.tsx     # Page: owns data fetching, layout
└── components/
    └── <Name>Component.tsx  # Presentational: props + events, no logic
```

Screens own data fetching. Components are presentational. Shared/reusable UI lives in `src/core/components/`; shared data access in `src/core/api/services/`.

## Data flow

```
Screen → apiClient (fetch + SecureStore JWT) → backend
```

| Layer | Depends on |
|-------|-----------|
| `src/core/api/client.ts` | `expo-secure-store` (token), `@devtest/shared` (ApiError) |
| `src/core/api/services/*.ts` | `src/core/api/client.ts`, `@devtest/shared` (DTOs) |
| Screens | Services + `@devtest/shared` (format, validators) |

## Routing & navigation

- Routes live in `app/`. Each file exports a default component (thin re-export of the feature screen).
- `app/(tabs)/_layout.tsx` renders `<NativeTabs>` with 4 static triggers: `index` (Inicio), `temas`, `progreso`, `perfil`.
- NativeTabs render no headers — screens show their title via the shared `Screen` primitive (`src/core/components/Screen.tsx`). When a tab needs sub-screens, nest a `<Stack>` group for that tab.
- Tab icons: SF Symbols (`sf`) + Material (`md`). Keep labels in Spanish.

## Design system

Source of truth: `.agents/mobile/design-system.md`. Mirror it as typed tokens in `src/core/theme/tokens.ts`:

```ts
import { colors, spacing, radius } from '@/core/theme/tokens'
```

- Style: Dark Mode (OLED). `app.json` sets `userInterfaceStyle: "dark"` — dark-only.
- Palette: background `#0F172A`, primary green `#22C55E`, accent amber `#D97706`, destructive red `#DC2626`. Difficulty gradient (green/amber/red) already matches `DIFFICULTY_COLORS` in `@devtest/shared`.
- Typography: system font for UI (apple-design: default to platform system font); monospace (JetBrains Mono / SF Mono) for code blocks.
- La skill `ui-ux-pro-max` puede proponer ajustes (opcional, no es el archivo canónico):
  `python .agents/skills/ui-ux-pro-max/scripts/search.py "coding challenge leetcode practice developer tool" --design-system`

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

## Path aliases

`tsconfig.json` maps `@/*` → `src/*`. Prefer aliases over relative imports. Expo resolves tsconfig paths (ensure `tsconfigPaths` is enabled if not default).

## Monorepo

`metro.config.js` adds the monorepo root to `watchFolders` and `nodeModulesPaths` so Metro resolves `@devtest/shared` from the workspace symlink. No build step required — Metro transpiles `.ts` source directly.

## Commands

| Command | Action |
|---------|--------|
| `pnpm --filter mobile start` | Start Expo dev server |
| `pnpm --filter mobile typecheck` | `tsc --noEmit` |
| `pnpm --filter mobile android` | Start on Android |
| `pnpm --filter mobile ios` | Start on iOS |