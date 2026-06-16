# dev-test-api — Agent instructions

Modular Go + Gin API with JWT auth, PostgreSQL/GORM, Swagger docs, and Air live reload.

## Commands

| Command | Action |
|---------|--------|
| `make install` | `go mod tidy` — sync dependencies |
| `make dev` | Start with live reload (`go tool air`) |
| `make build` | Build binary to `./tmp/main` |
| `make run` | `go run main.go` on `:8080` |
| `make swagger` | Regenerate `docs/` from annotations |
| `make clean` | Remove `tmp/` and `docs/` |
| `make fmt` | `go fmt ./...` — format all Go files |
| `make vet` | `go vet ./...` — static analysis |
| `make check` | Run `fmt` then `vet` |
| `make setup` | Install lefthook git hooks (run once after clone) |

## Key facts

- **NEVER create commits** — commits are a manual user action. Only stage, diff, and suggest changes; the user decides when and what to commit.
- **`internal/` is Go-private** — the compiler enforces that no external module can import `internal/*` packages. This protects module boundaries.
- **`docs/` and `tmp/` are auto-generated and gitignored** — never edit them manually.
- **Air** and **Lefthook** are project tools declared in `go.mod` (`tool github.com/air-verse/air`, `tool github.com/evilmartians/lefthook/v2`). Invoke via `go tool <name>` or their `make` targets. Not a global install.
- **Pre-commit hooks** via lefthook: runs `gofmt -l` check and `go vet ./...` on every commit. Run `make setup` after clone to install.
- **CI**: GitHub Actions workflow in `.github/workflows/ci.yml` — runs `gofmt -l`, `go vet`, and `go build` on push to `main` and pull requests.
- **No tests yet** — no testing framework to break.
- **Secrets** — never hardcode API keys, JWT secrets, DB credentials, or any sensitive values. Always use environment variables via `.env` (loaded by `godotenv` in `internal/config/config.go`). If you see a hardcoded secret while coding, extract it to a new env var in `config.go`. Never commit `.env` files.
- **Log icons** — always prefix `log.Println`, `log.Printf`, and `log.Fatal`/`log.Fatalf` calls with an emoji icon that best represents the context of the message. Choose the most descriptive icon for each case — don't reuse a generic one just because it was used elsewhere. The list below are examples already in the codebase, not a closed set: ❌ errors, ✅ success, 🚀 startup, 🛢️ database, 🌱 seed.
- **Swagger UI** at `/swagger/index.html`.
- **Context7 MCP** is available — use it to fetch current docs for Gin, GORM, jwt, swag, etc. rather than relying on training data. Always resolve library ID first, then query docs with the user's full question.
- **Dependency safety** — before installing any new library or package, validate it is trustworthy, secure, and maintained (see `.agents/rules/dependencies.md`). Never `go get` without running the checks first.

## Architecture

The project follows a **vertical-slice module** pattern. Business modules live in `internal/modules/`, GORM models live in `internal/models/`, and server setup lives in `internal/server/`.

```
internal/
├── config/                 # Env var loading → Config struct
│   └── config.go
├── database/               # GORM connection + AutoMigrate
│   └── database.go
├── middleware/              # Cross-cutting: JWT auth, logging
│   ├── auth.go
│   └── logger.go
├── models/                 # GORM structs (pure DB schema)
│   ├── user.go
│   └── question.go
├── modules/                # Business modules
│   ├── <module>/
│   │   ├── dto.go          # Request + response DTOs
│   │   ├── store.go        # interface Store + gormStore{} implementation
│   │   ├── service.go      # business logic, depends on Store interface
│   │   ├── handler.go      # Gin handlers, depend on Service interface
│   │   └── routes.go       # RegisterRoutes(rg *gin.RouterGroup, ...)
│   ├── auth/
│   ├── questions/
│   └── users/
└── server/
    └── server.go           # Gin engine + middleware + DI wiring + Run()
```

**Layers** → `Handler → Service → Store → DB`. Each layer only knows the interface of the one below it.

**Cross-module dependencies**: modules MAY import each other's `Store` interface and DTO types. For example, `auth` imports `users.Store` and `users.UserResponse`. Models in `internal/models/` import nothing from modules.

**Wiring**: dependency injection is **manual** in `internal/server/server.go`. No DI framework. Create stores → services → handlers → register routes.

**Response format**:
- Success: `{ "data": ... }` via `response.Success()` / `response.Paginated()`
- Errors: RFC 9457 via `response.Problem()` with typed errors from `pkg/apierr`
- **Language**: `title` in **English**, `detail` in **Spanish** with first letter capitalized (see `.agents/rules/responses.md` → *Error language*)

## Pagination

List endpoints that return many records MUST use pagination via `response.Paginated()`.

**Query params** — read from `c.DefaultQuery()` in the handler:

| Param | Default | Description |
|-------|---------|-------------|
| `page` | `1` | Page number |
| `per_page` | `20` | Items per page |

**Store** — keep both methods:

```go
FindAll() ([]Model, error)                               // sin paginación
FindPage(page, perPage int) ([]Model, int64, error)      // paginado, retorna items + total
```

**Service** — expose paginated List for collection endpoints:

```go
List(page, perPage int) ([]Response, int64, error)
```

**Handler** — read query params, call service, respond with `response.Paginated()`:

```go
page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

items, total, err := h.service.List(page, perPage)
response.Paginated(c, http.StatusOK, items, response.Meta{
    Total: total, Page: page, PerPage: perPage,
})
```

**Swagger** — add query params to paginated handlers:

```go
// @Param        page      query  int  false  "Número de página (default: 1)"
// @Param        per_page  query  int  false  "Elementos por página (default: 20)"
```

## Adding a new endpoint

1. Add handler function with Swagger annotations in the module's `handler.go`
2. Run **`make swagger`** to regenerate `docs/`
3. If the handler is in a new module: create the files following the module pattern above, then wire it in `internal/server/server.go`

## Rules files

For detailed conventions, see:
- `.agents/rules/responses.md` — error handling, response helpers, validation errors
- `.agents/rules/architecture.md` — module anatomy, store pattern, service pattern, wiring
- `.agents/rules/dependencies.md` — dependency safety: security, maintenance, licensing checks before `go get`
