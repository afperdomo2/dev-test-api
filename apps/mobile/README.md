# DevTest Mobile

> App **React Native** + **Expo** (SDK 57) + **TypeScript**

## 📋 Requisitos

| Herramienta | Versión |
|------------|---------|
| Node.js    | >= 20   |
| pnpm       | —       |
| Expo CLI   | via `npx expo` |

## 🚀 Inicio rápido

```bash
# Desde la raíz del workspace
pnpm install
pnpm mobile

# O directamente
cd apps/mobile && npx expo start
```

## 🔧 Variables de entorno

Copia `.env.example` a `.env` y ajusta según tu entorno:

```bash
cp .env.example .env
```

| Variable | Descripción | Default |
|----------|-------------|---------|
| `EXPO_PUBLIC_API_BASE_URL` | URL del backend (sin trailing slash) | `http://localhost:8080` |

**Para dispositivo real** en la misma red WiFi, usa la IP LAN de tu máquina:

```env
EXPO_PUBLIC_API_BASE_URL=http://192.168.1.x:8080
```

> En iOS simulator / Android emulator `localhost` funciona directamente. En dispositivo real necesitas la IP de tu máquina en la LAN.
>
> Con el stack Docker (`make docker-dev`) el API está en `:8082`; con el backend
> nativo (`make dev`) en `:8080`.

## 📁 Estructura

```
apps/mobile/
├── app/                 # Expo Router — rutas + _layout (AuthGate, tabs, login, setup)
├── src/
│   ├── core/
│   │   ├── api/         # client.ts (fetch + Bearer + unwrap envelope) + services/
│   │   ├── auth/        # storage.ts (SecureStore) + auth.store.ts (Zustand)
│   │   ├── components/  # Screen.tsx (primitivo compartido)
│   │   ├── theme/       # tokens.ts (colores, spacing, tipografía)
│   │   └── utils/       # errors.ts (extrae detail de RFC 9457)
│   └── features/
│       ├── auth/        # LoginScreen, SetupScreen, AuthCard, AuthInput, AuthButton
│       ├── home/        # HomeScreen (placeholder)
│       ├── profile/     # ProfileScreen (placeholder)
│       ├── progress/    # ProgressScreen (placeholder)
│       └── themes/      # ThemesScreen (placeholder)
├── metro.config.js      # Config monorepo (watchFolders, nodeModulesPaths)
├── app.json             # Config Expo (plugin expo-router, dark, scheme "devtest")
├── tsconfig.json        # Extiende expo/tsconfig.base + alias "@/*" → src/*
├── .env                 # Variables de entorno local (no commitear)
├── .env.example         # Plantilla de variables de entorno
└── assets/              # Íconos, splash
```

- Navegación: **Expo Router** con **NativeTabs** (4 tabs: Inicio, Temas, Progreso, Perfil).
- Design system (dark, OLED): `.agents/mobile/design-system.md` → tokens en `src/core/theme/tokens.ts`.
- Auth: `expo-secure-store` (JWT) + Zustand (`auth.store.ts`). El `AuthGate` en `app/_layout.tsx` gestiona login/setup/tabs.

## 🔗 Paquete compartido (@devtest/shared)

Los tipos, utilidades de formato y predicados de validación se importan desde `@devtest/shared`:

```ts
import type { User, Session } from '@devtest/shared'
import { formatDateTime, formatScore, isValidEmail } from '@devtest/shared'
```

No duplicar nada que ya exista en shared. Si algo es reutilizable entre web y mobile, pertenece a `packages/shared/`.

Contrato API (envelope, errores RFC 9457, camelCase): ver `.agents/backend/responses.md`.

## 📦 Scripts

| Comando | Descripción |
|---------|-------------|
| `pnpm --filter mobile start` | Expo dev server |
| `pnpm --filter mobile typecheck` | TypeScript check (`tsc --noEmit`) |
| `pnpm --filter mobile android` | Iniciar en Android |
| `pnpm --filter mobile ios` | Iniciar en iOS |

## 📱 APK manual (Android)

El APK se compila a mano con **EAS cloud** (no hay CI de APK). La URL del API se
**compila dentro del APK**: se configura en `eas.json` → `build.preview.env`,
no en `.env`, y debe ser HTTPS.

```bash
make mobile-apk      # preflight + eas build -p android --profile preview
```

Runbook completo, checklist de versión y cómo conectarse a la DB de la nube:
[`docs/mobile-apk.md`](../../docs/mobile-apk.md).

## 📚 Convenciones

- Arquitectura: `.agents/mobile/architecture.md`
- Dependencias: `.agents/mobile/dependencies.md`
