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

## Key facts

- **`internal/` is Go-private** — the compiler enforces that no external module can import `internal/*` packages. This protects module boundaries.
- **`docs/` and `tmp/` are auto-generated and gitignored** — never edit them manually.
- **Air is a project tool** declared in `go.mod` (`tool github.com/air-verse/air`). Invoke via `go tool air` or `make dev`. Not a global install.
- **No tests, no linter, no CI yet** — no testing framework or CI config to break.
- **Swagger UI** at `/swagger/index.html`.
- **Context7 MCP** is available — use it to fetch current docs for Gin, GORM, jwt, swag, etc. rather than relying on training data. Always resolve library ID first, then query docs with the user's full question.

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

## Adding a new endpoint

1. Add handler function with Swagger annotations in the module's `handler.go`
2. Run **`make swagger`** to regenerate `docs/`
3. If the handler is in a new module: create the files following the module pattern above, then wire it in `internal/server/server.go`

## Rules files

For detailed conventions, see:
- `.agents/rules/responses.md` — error handling, response helpers, validation errors
- `.agents/rules/architecture.md` — module anatomy, store pattern, service pattern, wiring
