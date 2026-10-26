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
| `app/(tabs)/_layout.tsx` | `NativeTabs` — 3 triggers: Inicio (`index`), Sesiones (`sesiones`), Perfil (`perfil`) | Yes |
| `src/features/<domain>/` | Feature (domain) modules: `screens/`, `components/` (y `utils/` cuando aplique) | Yes — add features here |
| `src/core/theme/tokens.ts` | Design tokens (colors, spacing, radii, typography) — mirror tipado de `.agents/mobile/design-system.md` | Yes — keep in sync with design-system.md |
| `src/core/components/` | Shared UI primitives (e.g. `Screen.tsx`) | Yes |
| `src/core/api/` | HTTP client + services (`client.ts`, `services/*.ts`) | Yes — add endpoints here |
| `src/core/auth/` | Almacenamiento del token (`expo-secure-store`) y estado de sesión | Yes |
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
- `app/(tabs)/_layout.tsx` renders `<NativeTabs>` with 3 triggers: `index` (Inicio), `sesiones` (Sesiones), `perfil` (Perfil).
- Nested tabs use Stack navigators (e.g. `app/(tabs)/sesiones/` has `_layout.tsx` with `<Stack>` for list → study → summary).
- NativeTabs render no headers — screens show their title via the shared `Screen` primitive or custom header. When a tab needs sub-screens, nest a `<Stack>` group for that tab.
- Tab icons: SF Symbols (`sf`) + Material (`md`). Keep labels in Spanish.

## UX conventions

All mobile UI must follow these rules (source: `pro-rules.md` + `ui-ux-pro-max` skill):

### Icons
- **Always use `MaterialCommunityIcons`** from `@expo/vector-icons` — never emojis for structural/navigation icons.
- Icon sizes: badges `12px`, action rows `20px`, headers `24-48px`.
- One icon style per hierarchy level (don't mix filled/outline at the same level).

### Touch & Interaction
- **Touch targets**: minimum 44pt (iOS) / 48dp (Android). Use `hitSlop` or `minHeight`/`minWidth` to expand.
- **Tap feedback**: every `Pressable` must show `opacity: 0.7-0.85` via `({ pressed }) =>` style pattern.
- **Destructive actions**: always use `Alert.alert()` confirmation before delete/logout.

### Accessibility
- `accessibilityLabel` on every interactive element (cards, buttons, icons).
- `accessibilityRole="button"` on all tappable elements.
- `accessibilityState` for dynamic states (selected, disabled, busy).
- Screen reader focus order must match visual order.

### List rows (sessions)
- **Máximo 3 líneas por card**: nombre / badges (estado + modo + dificultad) / meta (tiempo relativo + progreso + score). Si una card crece a 4–5 líneas, sobra información: muévela al detalle. Los temas no van en la card: el nombre de la sesión ya los insinúa y el resumen los muestra.
- Una sola línea de badges (`flexDirection: 'row'`, sin `flexWrap`) y **sin más de 3 badges**: un cuarto badge con `flexShrink` colapsa a un "punto con borde" ilegible.
- Filtros: **una sola fila con `ScrollView horizontal`**. Nunca `flexWrap`: 4 chips en español no caben y dejan una huérfana en la segunda fila.
  - El `ScrollView` necesita **`flexGrow: 0` + `flexShrink: 0` + `height` explícita**. RN aplica por dentro `flexGrow: 1, flexShrink: 1`; si solo anulas el grow, la lista comprime el contenedor cuando tiene items y los chips se recortan (se ven "cuadrados"). Con `height` fija no depende del contenido.
  - `showsHorizontalScrollIndicator={false}`, gap/padding en `contentContainerStyle` y `paddingRight` corto para que asome el siguiente chip.
- Chips de filtro: **versión compacta** (visual 36pt, `typography.size.xs`, `paddingHorizontal: 12`) con `hitSlop={{ top: 4, bottom: 4 }}` para que el target táctil siga siendo ≥44pt.
- Si hay una sesión `in_progress`, banner de reanudación arriba (nombre + progreso + botón 44pt): es la acción principal de la pantalla.
- Los colores de estado salen de `SESSION_STATUS_COLORS` (shared) traducidos a tokens con un `Record<SessionStatus, string>` — el compilador detecta estados nuevos.


### Layout & Spacing
- Use design tokens from `src/core/theme/tokens.ts` — never hardcode hex values.
- 8dp spacing rhythm (spacing.xs=4, sm=8, md=16, lg=24, xl=32).
- Expo Router tabs handle safe areas automatically — don't add SafeAreaView manually.
- Scroll content must not be hidden behind fixed bars.

### Progress & Feedback
- Visual progress bars for multi-step flows (sessions, quizzes) — not just text counters.
- Loading states: `ActivityIndicator` for initial load, skeleton for content.
- Status = icon + color + text (never color alone for error/success).

### Components
- `@gorhom/bottom-sheet` for action menus (requires `react-native-reanimated` + `react-native-gesture-handler`, both SDK-compatible via `expo install`).
- `react-native-markdown-display` for rendering markdown content (questions, explanations).
- `expo-constants` for app version (`Constants.expoConfig?.version`).

### Pre-delivery checklist
Before delivering any mobile UI, verify against `pro-rules.md`:
- [ ] No emojis as icons
- [ ] All touch targets ≥44pt
- [ ] Tap feedback on every Pressable
- [ ] `accessibilityLabel` on interactive elements
- [ ] Destructive actions have confirmation
- [ ] Status uses icon + color + text
- [ ] Design tokens used (no hardcoded colors)
- [ ] Progress bars for multi-step flows

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

## Push notifications

`src/core/push/` centraliza permisos, obtención del token y registro en el backend. Es **best-effort**: si algo falla, la app sigue funcionando sin push (ver N-09).

`expo-notifications` se carga con `loadNotifications()` (require perezoso), **nunca** con un `import` estático:

- Expo Go en Android (SDK 57) no registra `ExpoTopicSubscriptionModule` y el import estático rompe el arranque antes de renderizar (expo/expo#52175 → N-12).
- En Expo Go Android el push remoto no existe desde SDK 53, así que ahí se omite el módulo por completo.
- En development build / producción se carga normal y el push funciona.

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