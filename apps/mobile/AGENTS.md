# Mobile (Expo/RN) — quick context

Stack: **React Native + Expo SDK 57** + **Expo Router** (file-based) + **NativeTabs** + **TypeScript** + `@devtest/shared`.

## Commands

| Command | Action |
|---------|--------|
| `pnpm --filter mobile start` | Expo dev server |
| `pnpm --filter mobile typecheck` | `tsc --noEmit` (corre tras cada cambio) |
| `pnpm --filter mobile test:run` | Unit tests (`jest-expo`, co-located `*.test.ts`) |
| `pnpm --filter mobile android` / `ios` | Abrir en dispositivo |

## Structure

```
app/            # Expo Router — SOLO rutas y _layout (thin re-exports)
src/features/   # Dominios: auth, home, sessions, profile (screens/, components/)
src/core/       # theme (tokens), components (Screen, MarkdownText), api (client + services), auth, utils
```

Rutas actuales: `app/login.tsx` y `app/setup.tsx` (auth) + `app/(tabs)/` → Inicio (index), Sesiones (con sub-rutas), Perfil, y `+not-found`.

## Key facts

- **Dark-only** (design system Dark OLED) — `app.json` `userInterfaceStyle: "dark"`.
- **Design tokens** viven en `src/core/theme/tokens.ts` (fuente de verdad: `.agents/mobile/design-system.md`).
- Alias de import `@/*` → `src/*` (tsconfig + `experiments.tsconfigPaths`).
- Tabs con **NativeTabs** (`expo-router/unstable-native-tabs`), labels en español, iconos SF Symbol + Material.
- DTOs/tipos/formatters desde `@devtest/shared` — no duplicar.
- **Iconos**: `MaterialCommunityIcons` de `@expo/vector-icons` — nunca emojis para iconos estructurales.
- **Bottom sheet**: `@gorhom/bottom-sheet` para menús de acciones (requiere reanimated + gesture-handler del SDK).
- **Markdown**: `react-native-markdown-display` para contenido de preguntas/explicaciones.

## UX conventions

Ver `.agents/mobile/architecture.md` sección "UX conventions" para reglas completas. Resumen:
- Touch targets ≥44pt/48dp · Tap feedback en todo Pressable · `accessibilityLabel` en controles interactivos
- Confirmación `Alert.alert()` antes de acciones destructivas · Progress bars para flujos multi-paso
- Status = icono + color + texto (nunca solo color) · Usar design tokens (nunca hex hardcodeados)

## Details

- Arquitectura completa: `.agents/mobile/architecture.md`
- Design system: `.agents/mobile/design-system.md`
- Seguridad de dependencias: `.agents/mobile/dependencies.md`
