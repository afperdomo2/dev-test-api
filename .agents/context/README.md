# `.agents/context/` — reglas de arquitectura por capa

Documentación detallada que los `AGENTS.md` referencian pero **no duplican**. Es
la fuente de verdad de convenciones, contratos y patrones de cada ecosistema;
si una regla cambia, se cambia acá (y se ajusta el índice en `AGENTS.md`).

No agregar estos archivos a `opencode.json` → `instructions`: son ~130 KB y se
cargan de forma diferida desde los `AGENTS.md` (ver D-50 en `specs/DECISIONS.md`).
Las skills viven aparte, en `.agents/skills/` (gitignored, registradas en
`skills-lock.json`).

## Backend (Go, repo root)

| File | Purpose |
|------|---------|
| `backend/architecture.md` | Módulos, anatomía (request → response → mapper → store → service → handler → routes), wiring, migraciones goose, worker pool, Swagger, middlewares/timeouts |
| `backend/responses.md` | **SOURCE OF TRUTH** del contrato HTTP: envelope, errores RFC 9457, idioma de errores (title EN, detail ES), camelCase, DTOs, contexto JWT |
| `backend/data-model.md` | Contexto de las 13 tablas + enums + índices (canónico: `internal/models/*.go` y `internal/database/migrations/*.sql`) |
| `backend/dependencies.md` | Chequeos de seguridad antes de `go get` |

## Web (`apps/web/`)

| File | Purpose |
|------|---------|
| `web/architecture.md` | Ownership de directorios, data flow (api → queries → features), estructura de feature, checklist de wiring, re-exports de shared |
| `web/api-client.md` | Interceptores Axios, unwrap del envelope, patrones de service, `useFormErrors` (contrato en `backend/responses.md`) |
| `web/patterns.md` | TanStack Query, debounce 500ms, convenciones Vue/TS, Pinia, paginación, sintaxis `v-slot` |
| `web/dependencies.md` | Chequeos de seguridad antes de `pnpm add` |

## Mobile (`apps/mobile/`)

| File | Purpose |
|------|---------|
| `mobile/architecture.md` | Expo Router, feature-driven (`src/features` + `src/core`), NativeTabs, tokens, metro monorepo, SecureStore, HTTP client, UX conventions |
| `mobile/design-system.md` | **SOURCE OF TRUTH** del diseño mobile: Dark Mode (OLED), paleta, tipografía, anti-patrones |
| `mobile/dependencies.md` | Compatibilidad Expo y chequeo de módulos nativos antes de `pnpm add` |

## Shared (`packages/shared/`)

| File | Purpose |
|------|---------|
| `shared/architecture.md` | Qué va en shared vs app, barrel exports, re-export pattern en web, dedup de `CodeChallenge`, flujo de sync (swagger → shared → apps) |

## DevOps (`docker/`)

| File | Purpose |
|------|---------|
| `devops/containers.md` | **SOURCE OF TRUTH** de Docker: 3 entornos, mapa de archivos, invariantes (non-root, DB no publicada, same-origin, TZ, `TRUSTED_PROXIES`), Compose Watch, Dokploy, troubleshooting |
