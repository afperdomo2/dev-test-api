# Mobile (Expo/RN) — quick context

Stack: **React Native + Expo SDK 57** + **Expo Router** (file-based) + **NativeTabs** + **TypeScript** + `@devtest/shared`.

## Commands

| Command | Action |
|---------|--------|
| `pnpm --filter mobile start` | Expo dev server |
| `pnpm --filter mobile typecheck` | `tsc --noEmit` (corre tras cada cambio) |
| `pnpm --filter mobile android` / `ios` | Abrir en dispositivo |

## Structure

```
app/            # Expo Router — SOLO rutas y _layout (thin re-exports)
src/features/   # Dominios: home, themes, progress, profile (screens/, components/)
src/core/       # theme (tokens), components (Screen), api (por construir)
```

Rutas actuales: `app/(tabs)/` → Inicio (index), Temas, Progreso, Perfil + `+not-found`.

## Key facts

- **Dark-only** (design system Dark OLED) — `app.json` `userInterfaceStyle: "dark"`.
- **Design tokens** viven en `src/core/theme/tokens.ts` (fuente de verdad: `.agents/mobile/design-system.md`).
- Alias de import `@/*` → `src/*` (tsconfig + `experiments.tsconfigPaths`).
- Tabs con **NativeTabs** (`expo-router/unstable-native-tabs`), labels en español, iconos SF Symbol + Material.
- DTOs/tipos/formatters desde `@devtest/shared` — no duplicar.

## Details

- Arquitectura completa: `.agents/mobile/architecture.md`
- Design system: `.agents/mobile/design-system.md`
- Seguridad de dependencias: `.agents/mobile/dependencies.md`