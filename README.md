# dev-test-api

> REST API con **Go** + **Gin** · **JWT** · **PostgreSQL** · **Swagger** · live reload con **Air**
>
> Frontend SPA con **Vue 3** + **Vite** + **Vuetify** + **Pinia** + **TanStack Query** → [`apps/web/`](apps/web/)
>
> App mobile con **React Native** + **Expo** → [`apps/mobile/`](apps/mobile/)
>
> Paquete compartido **TypeScript** → [`packages/shared/`](packages/shared/)

## 📋 Requisitos

| Herramienta | Versión |
|------------|---------|
| Go         | 1.26+   |
| PostgreSQL | 14+     |
| Node.js    | 20+     |
| pnpm       | —       |
| make       | —       |
| Docker + Compose | 29+ (opcional, recomendado) |

## 🐳 Docker

Tres entornos listos para usar (guía completa: [`docker/README.md`](docker/README.md)):

| Entorno | Comando | URL | Para qué |
|---|---|---|---|
| **dev** | `make docker-dev` | http://localhost:3030 | Hot reload (Air + Vite HMR) |
| **prod-local** | `make docker-prod` | http://localhost:8088 | Ensayo real en la PC (DB incluida) |
| **Dokploy** | `docker/compose.dokploy.yaml` | `https://tu-dominio` | Deploy en VPS |

```bash
make docker-env         # crea docker/.env (dev)
make docker-dev         # levanta dev con Compose Watch
make docker-env-prod    # crea docker/.env.prod (prod-local)
make docker-prod        # build + up + healthchecks
make docker-smoke       # verifica API, SPA y cabeceras
```

La primera vez: editá `docker/.env.prod`, levantá prod-local, creá el admin con
el `SETUP_TOKEN`, vacialo y redeployá. La DB no se publica; backup/restore con
`make docker-backup` / `make docker-restore FILE=…`.

Mobile no corre en Docker: Metro en el host y APK manual — **Gradle local**
(iteración) o **EAS cloud** (release) — [`docs/mobile-apk.md`](docs/mobile-apk.md).

## 📁 Estructura del proyecto

```
dev-test-api/
├── main.go
├── internal/                  # Módulos de negocio Go
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── modules/
│   └── server/
├── pkg/                       # Paquetes públicos
├── brand/                     # Fuente de iconos web + mobile (make icons)
├── docs/                      # Swagger auto-generado + guías (branching, APK, branding)
├── apps/
│   ├── web/                   # SPA Vue 3 + Vite (ver apps/web/README.md)
│   └── mobile/                # Expo React Native (ver apps/mobile/README.md)
├── packages/
│   └── shared/                # Tipos, utils, constants compartidos (ver packages/shared/)
├── .agents/
│   ├── context/               # Reglas de arquitectura por capa (backend, web, mobile, shared, devops)
│   └── skills/                # Skills de agentes (gitignored, ver skills-lock.json)
├── pnpm-workspace.yaml
├── Makefile
└── README.md
```

Cada módulo en `internal/` contiene todo lo que necesita: modelo, store (DB), servicio (lógica) y handler (HTTP). Esto permite extraer un módulo como microservicio sin desarmar el proyecto.

## 🔐 Flujo de autenticación

### Bootstrap — primer usuario admin

```
POST /api/v1/auth/setup
```

Solo funciona si **no hay usuarios en la base de datos**. Crea el primer usuario con `is_admin: true` y devuelve un JWT.

Si ya existen usuarios, devuelve `409 Conflict`.

### Login

```
POST /api/v1/auth/login
```

Devuelve un JWT que se usa como `Authorization: Bearer <token>` en los endpoints protegidos.

### Usuarios (admin)

```
GET    /api/v1/users          # listar
POST   /api/v1/users          # crear
GET    /api/v1/users/:id      # obtener
DELETE /api/v1/users/:id      # soft-delete
```

## ⚙️ Variables de entorno

Copiá `.env.example` a `.env`:

```bash
cp .env.example .env
```

Los valores de `.env.example` son los defaults del código. En `.env` solo
necesitás:

| Variable | Cuándo |
|---|---|
| `JWT_SECRET` | Siempre (obligatoria) |
| `DB_PASSWORD` | Siempre (obligatoria) |
| `SETUP_TOKEN` | Solo con la DB vacía (mínimo 16 caracteres) |

El resto podés omitirlo salvo que quieras cambiar el default. La lista completa,
con comentarios, está en [`.env.example`](.env.example).

Los compose de Docker fijan defaults propios por entorno (CORS, HSTS,
`TRUSTED_PROXIES`, `JWT_EXPIRY_HOURS`, `GIN_MODE`, `SWAGGER_ENABLED`); ver
[`docker/env/`](docker/env/) y [`docker/README.md`](docker/README.md).

Para web y mobile: [`apps/web/.env.example`](apps/web/.env.example) y
[`apps/mobile/.env.example`](apps/mobile/.env.example).

## 🛠️ Scripts (Makefile)

### Backend

| Comando | Descripción |
|---------|-------------|
| `make install` | Instala dependencias Go (`go mod tidy`) |
| `make dev` | Live reload con Air |
| `make build` | Compila binario en `./tmp/main.exe` |
| `make run` | Corre el servidor con `go run` |
| `make swagger` | Regenera docs Swagger |
| `make clean` | Elimina `tmp/` y `docs/` |
| `make fmt` | `go fmt ./...` |
| `make vet` | `go vet ./...` |
| `make check` | `fmt` + `vet` |
| `make test` | Pruebas unitarias (`go test -short`) |
| `make test-cover` | Pruebas con cobertura |
| `make setup` | Instala git hooks (lefthook) |
| `make migrate-up` | Aplica migraciones pendientes (goose, carga `.env`) |
| `make migrate-down` | Revierte la última migración (`down`) |
| `make migrate-status` | Estado de las migraciones |
| `make migrate-create name=<desc>` | Crea `NNNNN_<desc>.sql` (numeración secuencial) |

### Web (apps/web/)

| Comando | Descripción |
|---------|-------------|
| `make fe-install` | `pnpm install` |
| `make fe-dev` | Servidor de desarrollo (:3000) |
| `make fe-build` | Build de producción |
| `make fe-lint` | ESLint |
| `make fe-check` | TypeScript check |
| `make e2e` | E2E Playwright (requiere DB dedicada + backend) |
| `make e2e-ui` | Playwright en modo UI |
| `make e2e-report` | Abre el reporte del último run |

### Mobile (apps/mobile/)

| Comando | Descripción |
|---------|-------------|
| `make mobile-start` | Expo dev server |
| `make mobile-android` | Abrir en Android |
| `make mobile-ios` | Abrir en iOS |
| `make mobile-web` | Abrir en web (Expo) |
| `make mobile-typecheck` | TypeScript check |
| `make mobile-lint` | Lint |
| `make mobile-apk` | Preflight + APK manual en EAS cloud ([runbook](docs/mobile-apk.md)) |
| `make mobile-apk-local` | APK release local con Gradle (sin EAS, sin cuenta Expo) |
| `make mobile-bump` | Sube `version` (semver) y `versionCode` (+1) antes de compilar |

#### APK local (Android, sin EAS)

```bash
make mobile-apk-local                             # URL: devtest.afperdomo.cloud
make mobile-apk-local API_URL=https://otro.com    # otra URL (se hornea en el bundle)
```

Equivalente manual (requiere `~/Android/Sdk` + JDK 21):

```bash
cd apps/mobile/android
rm -rf app/build/generated/assets app/build/intermediates/assets app/build/intermediates/compressed_assets
ANDROID_HOME=$HOME/Android/Sdk JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64 \
EXPO_PUBLIC_API_BASE_URL=https://devtest.afperdomo.cloud ./gradlew assembleRelease
```

`EXPO_PUBLIC_API_BASE_URL` queda **compilada** en el APK: si cambia el backend,
recompilá. Versioná antes de compilar con `make mobile-bump [VERSION=x.y.z]`
(semver + `versionCode` +1, ver [`docs/mobile-apk.md`](docs/mobile-apk.md) → "Versionado").

**APK resultante** (ruta para pegar en el gestor de archivos):

```
apps/mobile/android/app/build/outputs/apk/release/app-release.apk
```

### Branding (iconos web + mobile)

| Comando | Descripción |
|---------|-------------|
| `make icons` | Regenera favicon/PWA y assets de la app desde `brand/*.svg` ([detalle](docs/branding.md)) |

### Docker

| Comando | Descripción |
|---------|-------------|
| `make docker-env` / `docker-env-prod` | Crea `docker/.env` / `docker/.env.prod` desde los ejemplos |
| `make docker-dev` | Dev con hot reload (Compose Watch) |
| `make docker-dev-up` / `docker-dev-down` / `docker-dev-reset` | Dev en background / bajar / borrar datos |
| `make docker-prod` | Prod-local: build + up + healthchecks (`:8088`) |
| `make docker-prod-down` / `docker-prod-reset` / `docker-prod-logs` | Operación de prod-local |
| `make docker-build` | Construye las imágenes de producción |
| `make docker-config` | Valida los 3 compose |
| `make docker-seed` | Siembra las 50 preguntas en dev |
| `make docker-smoke` | Verifica nginx, API, SPA y cabeceras |
| `make docker-backup` / `docker-restore FILE=…` | Backup / restore de la DB |
| `make e2e-docker` | Playwright contra prod-local |

Equivalentes directos con pnpm:

| Comando | Descripción |
|---------|-------------|
| `pnpm --filter mobile start` | Expo dev server |
| `pnpm --filter mobile typecheck` | TypeScript check |

### Workspace

| Comando | Descripción |
|---------|-------------|
| `pnpm typecheck` | Type-check de todos los packages |
| `pnpm web` | Iniciar dev server web |
| `pnpm mobile` | Iniciar Expo dev server |

## 📖 Documentación Swagger

Los endpoints se documentan mediante anotaciones en el código.

```bash
make swagger
```

Con el servidor corriendo: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 📦 Dependencias

### Backend

| Paquete | Uso |
|---------|-----|
| `gin-gonic/gin` | Framework HTTP |
| `gorm.io/gorm` | ORM |
| `gorm.io/driver/postgres` | Driver PostgreSQL |
| `golang-jwt/jwt/v5` | JWT tokens |
| `golang.org/x/crypto` | bcrypt |
| `joho/godotenv` | Variables de entorno |
| `swaggo/swag` | Docs OpenAPI |
| `swaggo/gin-swagger` | Swagger UI |
| `air-verse/air` | Live reload (project tool) |
| `stretchr/testify` | Pruebas unitarias |

### Shared (@devtest/shared)

Paquete de TypeScript puro compartido entre web y mobile:
- DTOs y tipos de API (`ApiError`, `PaginatedResponse`, etc.)
- Utilidades de formato (`formatDateTime`, `formatScore`)
- Predicados de validación (`isValidEmail`, `isRequired`)
- Constantes (`ITEMS_PER_PAGE_OPTIONS`, `DEFAULT_PER_PAGE`)

Ver `packages/shared/` y `.agents/context/shared/architecture.md` para detalles.

## 🧪 Testing

### Backend

Pruebas unitarias con `testing` stdlib + `testify`. Sin BD (mocks de `Store`).

```bash
make test        # unitarias
make test-cover  # con cobertura (genera coverage.out)
go test -race ./...            # con detector de carreras
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Web

Pruebas con Vitest + Vue Test Utils + Testing Library.

```bash
pnpm --filter dev-test-frontend test:run
```

#### E2E (Playwright)

Un único camino crítico (`status → setup → login → dashboard → sessions`) contra backend + Postgres reales. **No** lo ejecutan `make fe-test`, `make fe-check` ni `make fe-lint`.

La DB debe estar **vacía**: `/auth/setup` solo funciona cuando no existe ningún admin.

```bash
# 1) DB dedicada (una vez; para repetir: DROP DATABASE devtest_e2e; CREATE DATABASE devtest_e2e;)
DB_NAME=devtest_e2e make migrate-up

# 2) Backend contra esa DB
DB_NAME=devtest_e2e make run

# 3) Suite (levanta Vite; la primera vez instala el navegador)
pnpm --filter dev-test-frontend exec playwright install chromium
make e2e
```

En PowerShell, prefija con `$env:DB_NAME='devtest_e2e';` en lugar de `DB_NAME=...`.
Si el backend no responde en `:8080` o la DB ya tiene un admin, la suite falla con el comando exacto para arreglarlo.

Hooks de git (`lefthook.yml`) — todos los jobs se filtran por `glob`, así que un cambio que solo toca docs/markdown no ejecuta nada:

- **pre-commit** (rápido, solo archivos staged): `gofmt` + `go vet` + `go test` **en los paquetes Go afectados**, `shared-typecheck`, web (prettier, eslint, vue-tsc, `vitest related`) y mobile (`jest --findRelatedTests`).
- **pre-push** (completo, solo archivos empujados): backend (fmt/vet/build/test), web (format/lint/type/test/build) y mobile (`test:run`).

## 🔒 Seguridad

El repo escanea secretos con [gitleaks](https://github.com/gitleaks/gitleaks) en cada PR y push a `main` (reglas en `.gitleaks.toml`). Para comprobarlo en local antes de commitear:

```bash
gitleaks detect --source . --redact   # historia completa; debe salir "no leaks found"
gitleaks protect --staged --redact    # solo lo que vas a commitear
```

Instalación: `brew install gitleaks` (macOS), `choco install gitleaks` o `scoop install gitleaks` (Windows), o el binario de la release oficial. `.env`, `*.pem` y `*.key` ya están en `.gitignore`.

Política de reporte, plazos y rotación de credenciales: [`SECURITY.md`](SECURITY.md).

## 🌿 Ramas y releases

- **`develop`**: trabajo diario y backup (no despliega).
- **`main`**: producción, solo por PR desde `develop` y con checks obligatorios.
- **Release semanal**: PR `develop` → `main` el mismo día del merge, tag `vX.Y.Z` después de verificar.

Detalle (hotfix, rollback, migraciones, Dokploy): [`docs/branching.md`](docs/branching.md).

## 🧪 Primer uso

```bash
# 1. Instalar dependencias
make install
pnpm install

# 2. Configurar .env con los datos de PostgreSQL y SETUP_TOKEN
cp .env.example .env
# SETUP_TOKEN: mínimo 16 caracteres (ej: openssl rand -hex 32).
# Con la DB vacía es obligatorio: sin él el servidor no arranca.

# 3. Crear el primer usuario admin
make dev
curl -X POST http://localhost:8080/api/v1/auth/setup \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123456","setupToken":"<tu-setup-token>"}'

# 4. Usar el token devuelto para los endpoints protegidos
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <tu-token>"
```

El formulario de setup web (`/setup`) y mobile piden el mismo token en el campo
**“Token de instalación”**. Si `SETUP_TOKEN` está vacío y ya existen usuarios, el
setup queda deshabilitado (403) y el servidor solo muestra un aviso al arrancar.

## 📝 Respuestas

**Error** — RFC 9457 Problem Detail:

```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 422,
  "detail": "email: required",
  "instance": "/api/v1/users"
}
```

**Éxito** — envelope `{ "data": ... }`:

```json
{
  "data": {
    "token": "...",
    "user": { "id": "...", "email": "...", "is_admin": true }
  }
}
```

## 📝 Notas

- El esquema de DB se gestiona con **migraciones versionadas goose** (`internal/database/migrations/*.sql`, baseline `00001_init.sql`). Se aplican al arrancar el servidor (embedded) o con `make migrate-*`. Cualquier cambio de esquema es una migración nueva, nunca `AutoMigrate`.
- `docs/swagger.yaml`, `docs/swagger.json` y `docs/docs.go` son auto-generados y no se editan manualmente; el resto de `docs/` son guías versionadas. `tmp/` también es auto-generado.
- Air se instaló como project tool (`go get -tool`), no requiere instalación global.
- `internal/` es privado al módulo Go — no puede ser importado desde fuera del proyecto.

## 📄 Licencia

Este proyecto está licenciado bajo la [PolyForm Noncommercial License 1.0.0](LICENSE).

**Resumen**: Puedes usar, copiar, modificar y distribuir este software para fines **no comerciales**. No está permitido usarlo con fines de lucro o comercialización.

Para más detalles, consulta el archivo [LICENSE](LICENSE).
