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

## 📁 Estructura

```
apps/mobile/
├── app/               # Expo Router — solo rutas y _layout (app/(tabs)/, +not-found)
├── src/
│   ├── features/      # Dominios: home, themes, progress, profile (screens/, components/)
│   └── core/          # theme (tokens), components (Screen), api (por construir)
├── metro.config.js    # Config monorepo (watchFolders, nodeModulesPaths)
├── app.json           # Config Expo (plugin expo-router, dark, scheme "devtest")
├── tsconfig.json      # Extiende expo/tsconfig.base + alias "@/*" → src/*
└── assets/            # Íconos, splash
```

- Navegación: **Expo Router** con **NativeTabs** (4 tabs: Inicio, Temas, Progreso, Perfil).
- Design system (dark, OLED): `.agents/mobile/design-system.md` → tokens en `src/core/theme/tokens.ts`.

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

## 📚 Convenciones

- Arquitectura: `.agents/mobile/architecture.md`
- Dependencias: `.agents/mobile/dependencies.md`
