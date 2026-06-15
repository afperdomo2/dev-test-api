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

The project follows a **vertical-slice module** pattern. Each module in `internal/` owns everything it needs:

```
internal/<module>/
├── model.go      # GORM struct + request/response DTOs
├── store.go      # interface Store + gormStore{} implementation
├── service.go    # business logic, depends on Store interface
├── handler.go    # Gin handlers, depend on Service interface
└── routes.go     # RegisterRoutes(rg *gin.RouterGroup, ...)
```

**Layers** → `Handler → Service → Store → DB`. Each layer only knows the interface of the one below it.

**Cross-module dependencies**: modules MAY import each other's `Store` interface and model types. For example, `auth` imports `users.Store` and `users.User`. This is intentional — no circular deps yet.

**Wiring**: dependency injection is **manual** in `main.go`. No DI framework. Create stores → services → handlers → register routes.

**Response format**:
- Success: `{ "data": ... }` via `response.Success()` / `response.Paginated()`
- Errors: RFC 9457 via `response.Problem()` with typed errors from `pkg/apierr`

## Adding a new endpoint

1. Add handler function with Swagger annotations in the module's `handler.go`
2. Run **`make swagger`** to regenerate `docs/`
3. If the handler is in a new module: create the 5 files following the module pattern above, then wire it in `main.go`

## Rules files

For detailed conventions, see:
- `.agents/rules/responses.md` — error handling, response helpers, validation errors
- `.agents/rules/architecture.md` — module anatomy, store pattern, service pattern, wiring
