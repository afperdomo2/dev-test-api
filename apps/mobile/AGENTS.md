# Mobile (Expo/RN) — quick context

Stack: **React Native + Expo SDK 57** + **Expo Router** (file-based) + **NativeTabs** + **TypeScript** + `@devtest/shared`.

## Commands

| Command | Action |
|---------|--------|
| `pnpm --filter mobile start` | Expo dev server |
| `pnpm --filter mobile typecheck` | `tsc --noEmit` (corre tras cada cambio) |
| `pnpm --filter mobile test:run` | Unit tests (`jest-expo`, co-located `*.test.ts`) |
| `pnpm --filter mobile android` / `ios` | Abrir en dispositivo |
| `make mobile-apk-local` | APK Android local con Gradle (`API_URL=…`, `APP_NAME=…` opcional); ruta y detalle en `docs/mobile-apk.md` |
| `make mobile-bump [VERSION=x.y.z]` | Sube `version` (semver) y `versionCode` (+1) y sincroniza el proyecto nativo |

## Structure

```
app/            # Expo Router — SOLO rutas y _layout (thin re-exports)
src/features/   # Dominios: auth, home, sessions, profile (screens/, components/)
src/core/       # theme (tokens), components (Screen, MarkdownText), api (client + services), auth, utils
```

Rutas actuales: `app/login.tsx` y `app/setup.tsx` (auth) + `app/(tabs)/` → Inicio (index), Sesiones (lista, `nueva`, detalle y resumen), **Progreso** (dashboard + `temas`, `actividad`, `logros`) y Perfil, más `+not-found`.

## Key facts

- **Dark-only** (design system Dark OLED) — `app.json` `userInterfaceStyle: "dark"`.
- **Design tokens** viven en `src/core/theme/tokens.ts` (fuente de verdad: `.agents/context/mobile/design-system.md`).
- Alias de import `@/*` → `src/*` (tsconfig + `experiments.tsconfigPaths`).
- Tabs con **NativeTabs** (`expo-router/unstable-native-tabs`), labels en español, iconos SF Symbol + Material.
- DTOs/tipos/formatters desde `@devtest/shared` — no duplicar.
- **Iconos**: `MaterialCommunityIcons` de `@expo/vector-icons` — nunca emojis para iconos estructurales.
- **`expo-notifications` se carga en diferido** (`loadNotifications()` en `src/core/push/notifications.ts`): un import estático rompe Expo Go en Android (N-12).
- **Bottom sheet**: `@gorhom/bottom-sheet` para menús de acciones (requiere reanimated + gesture-handler del SDK).
- **Markdown**: `react-native-markdown-display` para contenido de preguntas/explicaciones.
- **Nombre por entorno**: `EXPO_PUBLIC_APP_NAME` (build-time, default `DevTest`). En local: `make mobile-apk-local APP_NAME="DevTest Pruebas"` (ejecuta prebuild). `app.config.ts` sobreescribe `expo.name`; `app.json` sigue siendo la fuente de `version`/`versionCode`.
- **Iconos**: `assets/icon.png` + capas adaptive (foreground/monochrome) se generan desde `brand/` con `make icons`. Cambiar assets exige `expo prebuild` + bump + rebuild. Detalle: `docs/branding.md`.

## Versionado

`app.json` es la fuente única: `expo.version` (semver) + `expo.android.versionCode`.

- **Todo cambio mobile que entregue APK nuevo se versiona** — no se cierra un
  cambio dejando dos APKs con el mismo `versionCode`.
- MAJOR si rompe el contrato con el backend, MINOR si agrega funcionalidad
  visible, PATCH si es fix/config/estilo/textos.
- Antes de compilar: `make mobile-bump` (mismo version) o
  `make mobile-bump VERSION=x.y.z`. Detalle: `docs/mobile-apk.md` → "Versionado".

## UX conventions

Ver `.agents/context/mobile/architecture.md` sección "UX conventions" para reglas completas. Resumen:
- Touch targets ≥44pt/48dp · Tap feedback en todo Pressable · `accessibilityLabel` en controles interactivos
- Confirmación `Alert.alert()` antes de acciones destructivas · Progress bars para flujos multi-paso
- Status = icono + color + texto (nunca solo color) · Usar design tokens (nunca hex hardcodeados)

## Details

- Arquitectura completa: `.agents/context/mobile/architecture.md`
- Design system: `.agents/context/mobile/design-system.md`
- Seguridad de dependencias: `.agents/context/mobile/dependencies.md`
