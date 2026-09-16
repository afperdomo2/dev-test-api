# DevTest Frontend

> SPA con **Vue 3** + **Vite** + **Vuetify** + **Pinia** + **TanStack Query** + **TypeScript**

## 📋 Requisitos

| Herramienta | Versión |
|------------|---------|
| Node.js    | >= 20   |
| pnpm       | —       |

## 🚀 Inicio rápido

```bash
# Desde la raíz del workspace
pnpm install
pnpm web              # Servidor de desarrollo (:3000)

# O desde esta carpeta
make fe-dev
```

## 📁 Estructura

```
apps/web/
├── src/
│   ├── api/                     # Capa HTTP (Axios + interceptores)
│   │   ├── client.ts            # Axios instance: JWT, unwrap envelope, errores RFC 9457
│   │   └── services/            # Funciones puras que llaman a cada endpoint
│   │       ├── auth.service.ts
│   │       ├── users.service.ts
│   │       ├── topics.service.ts
│   │       ├── questions.service.ts
│   │       ├── sessions.service.ts
│   │       └── progress.service.ts
│   │
│   ├── queries/                 # TanStack Query (caché, loading, refetch)
│   │   ├── auth.queries.ts
│   │   ├── users.queries.ts
│   │   └── ...
│   │
│   ├── composables/             # Composables reutilizables
│   │   ├── useDebounce.ts       # Debounce 500ms para búsquedas
│   │   ├── usePagination.ts     # page/perPage/total
│   │   └── useFormErrors.ts     # Mapeo errores API → campos
│   │
│   ├── layouts/                 # Layouts de la aplicación
│   │   ├── AuthLayout.vue       # Layout para login/setup
│   │   └── DefaultLayout.vue    # AppBar + NavigationDrawer + main
│   │
│   ├── features/                # Funcionalidades por dominio
│   │   ├── auth/                # Login + setup
│   │   ├── dashboard/           # Dashboard post-login
│   │   ├── users/               # CRUD usuarios (admin)
│   │   ├── questions/           # Listado + detalle + filtros
│   │   ├── topics/              # Listado + CRUD (admin)
│   │   ├── sessions/            # Sesiones de estudio
│   │   └── progress/            # Progreso SM-2
│   │
│   ├── router/                  # Vue Router + guards
│   │   ├── index.ts
│   │   └── routes.ts
│   │
│   ├── stores/                  # Pinia stores
│   │   ├── auth.store.ts        # JWT, user, login/logout
│   │   └── app.store.ts         # Tema, sidebar, snackbar
│   │
│   ├── types/                   # Re-exports desde @devtest/shared
│   │   ├── api.types.ts         # → ApiError, ApiResponse, PaginatedResponse
│   │   ├── auth.types.ts        # → LoginRequest, AuthResponse
│   │   └── ...
│   │
│   ├── utils/                   # Utilidades
│   │   ├── format.ts            # Re-export desde @devtest/shared
│   │   ├── validators.ts        # Reglas Vuetify (predicados desde shared)
│   │   └── storage.ts           # localStorage helpers (token)
│   │
│   └── plugins/                 # Plugins
│       ├── vuetify.ts           # createVuetify + tema + íconos (MDI)
│       └── query.ts             # VueQueryPlugin + defaults
│
├── .env                         # VITE_API_BASE_URL
├── .eslintrc.json               # ESLint (no any, strict TS, Vue rules)
├── .prettierrc                  # Prettier
├── vite.config.ts               # Vite + Vue + Vuetify + proxy /api
└── tsconfig.json
```

## 🔗 Paquete compartido (@devtest/shared)

Los `types/*`, `constants/*` y `utils/format.ts` son **re-exports** desde `packages/shared/`. Los predicados de validación (`isValidEmail`, `isRequired`) también viven en shared; las reglas Vuetify (`requiredRule`, `emailRule`) se definen en `utils/validators.ts` importando desde shared.

Al agregar un tipo nuevo: agregar la definición en `packages/shared/src/types/` y luego el re-export en `apps/web/src/types/`.

Contrato API (envelope, errores RFC 9457, camelCase): ver `.agents/backend/responses.md`.

## 🔗 Conexión con la API

El cliente HTTP (`api/client.ts`) gestiona automáticamente:

- **JWT** — agrega `Authorization: Bearer <token>` a cada request
- **Envelope** — desenvuelve `{ data: ... }` en respuestas exitosas
- **Errores RFC 9457** — transforma errores del backend en `ApiError` tipado
- **401** — limpia el token al recibir un 401

## 🧠 TanStack Query

Cada módulo del dominio exporta:

- **`queryOptions`** — definiciones reutilizables con `queryKey`, `queryFn`, `staleTime`
- **Mutaciones** — `mutationKey` + `mutationFn` para POST/PUT/DELETE

StaleTime por defecto: **60s**. Compartir `queryKey` entre componentes reutiliza el caché automáticamente.

## 🔐 Autenticación

- **Setup**: `POST /api/v1/auth/setup` — primer usuario admin (solo si DB vacía)
- **Login**: `POST /api/v1/auth/login` — devuelve JWT + user
- **Router guard**: `beforeEach` redirige a `/login` si `meta.requiresAuth` y no hay token
- **Admin guard**: verifica `meta.requiresAdmin` y `is_admin` en claims

## 🎨 Vuetify

- Tema **light/dark** con toggle en AppBar
- Íconos **Material Design Icons** (MDI)
- Fuente: **Roboto** via `unplugin-fonts`

## 📏 ESLint

Reglas destacadas:
- `@typescript-eslint/no-explicit-any` → **error**
- `@typescript-eslint/no-non-null-assertion` → **error**
- `@typescript-eslint/consistent-type-imports` → type imports
- `vue/component-api-style` → script-setup
- `vue/block-order` → script → template → style
- `no-console` → warn

Comandos: `pnpm lint`, `make fe-lint`

## ⚙️ Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `VITE_API_BASE_URL` | `http://localhost:8080` | URL base de la API |

El servidor de desarrollo usa el proxy en `vite.config.ts` para `/api` → `localhost:8080`.

## 📦 Scripts

| Comando | Descripción |
|---------|-------------|
| `pnpm dev` | Servidor de desarrollo (:3000) |
| `pnpm build` | Build de producción |
| `pnpm type-check` | TypeScript check |
| `pnpm lint` | ESLint + fix |
| `pnpm format` | Prettier |
