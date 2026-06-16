# dev-test-api — Agent instructions

Fullstack monorepo: **Go + Gin** REST API + **Vue 3 + Vite + Vuetify** SPA frontend.

## Commands

### Backend (root)

| Command | Action |
|---------|--------|
| `make install` | `go mod tidy` — sync Go dependencies |
| `make dev` | Start Go server with live reload (`go tool air`) |
| `make build` | Build Go binary to `./tmp/main` |
| `make run` | `go run main.go` on `:8080` |
| `make swagger` | Regenerate `docs/` from Swagger annotations |
| `make check` | `make fmt` then `make vet` |
| `make setup` | Install lefthook git hooks (run once after clone) |
| `make clean` | Remove `tmp/` and `docs/` |

### Frontend (`frontend/`)

| Command | Action |
|---------|--------|
| `make fe-install` | `pnpm install` |
| `make fe-dev` | Vite dev server on `:3000` (proxies `/api` → `:8080`) |
| `make fe-build` | Production build (`vue-tsc` + `vite build`) |
| `make fe-lint` | ESLint (`--fix`) |
| `make fe-check` | TypeScript check (`vue-tsc --build`) |

Always run `make fe-check` and `make fe-lint` after any frontend change. The `make fe-build` command runs `type-check` automatically before `vite build`.

## Key facts

- **NEVER create commits** — commits are manual. Only stage, diff, suggest.
- **Secrets** — never hardcode keys, passwords, or tokens. Use `.env` for Go, `VITE_*` env vars for frontend. Never commit `.env`.
- **`internal/` is Go-private** — the compiler blocks external imports.
- **`docs/` and `tmp/` are auto-generated and gitignored** — never edit them manually.
- **Air** and **Lefthook** are Go project tools in `go.mod` — invoke via `go tool air`, `go tool lefthook`, or their `make` targets.
- **Pre-commit hooks** (lefthook): `gofmt -l` check + `go vet ./...`. Run `make setup` after clone.
- **CI** (`.github/workflows/ci.yml`): `gofmt -l`, `go vet ./...`, `go build ./...` on push/PR to `main`. Frontend is NOT yet in CI.
- **Swagger UI** at `/swagger/index.html` when server is running.
- **Context7 MCP** — available for Gin, GORM, jwt, swag docs. Always `resolve-library-id` first, then `query-docs`.
- **Vuetify MCP** — available via `opencode.json`. Query component APIs with `vuetify_get_component_api_by_version` or `vuetify_get_feature_guide`.
- **Log icons** — prefix Go `log.*` calls with emoji (❌ errors, ✅ success, 🚀 startup, 🛢️ database, 🌱 seed). Pick the most descriptive icon per context.
- **No tests yet** — no testing framework to break in either backend or frontend.

## Domain rules

Backend and frontend have separate rule files with detailed conventions:

### Backend: `.agents/backend/`

| File | Purpose |
|------|---------|
| `architecture.md` | Module anatomy (dto → store → service → handler → routes), wiring, GORM, Swagger |
| `responses.md` | Envelope format, RFC 9457 errors, error language (title: EN, detail: ES), DTO conventions, JWT context |
| `dependencies.md` | Safety checks before `go get` — security, maintenance, licensing |

### Frontend: `.agents/frontend/`

| File | Purpose |
|------|---------|
| `architecture.md` | Directory ownership, data flow (api → queries → features), feature structure, wiring checklist |
| `api-client.md` | Axios interceptors, envelope unwrap gotcha (`meta` key), error mapping, `useFormErrors` |
| `patterns.md` | TanStack Query patterns, debounce 500ms, Vue/TS conventions, Pinia stores, pagination, v-slot bracket syntax |
