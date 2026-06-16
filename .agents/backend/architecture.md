# Architecture rules

## Directory ownership

| Path | Purpose | Editable? |
|------|---------|-----------|
| `internal/models/` | GORM structs (pure DB schema) | Yes — add new models here |
| `internal/modules/<module>/` | Business modules (auth, users, questions) | Yes — add features here |
| `internal/server/server.go` | Gin engine + middleware + DI wiring + Run() | Yes — register new modules here |
| `main.go` | Server entrypoint (config.Load + db.Connect + server.Run) | Rarely |
| `internal/config/config.go` | Env var loading | Rarely — add new env vars only |
| `internal/database/database.go` | GORM connection + AutoMigrate | Rarely — register new models here |
| `internal/middleware/` | Cross-cutting: JWT auth, logging | Rarely |
| `pkg/` | Shared utilities: response, apierr | Rarely |
| `docs/` | Swagger auto-generated | **NEVER** — use `make swagger` |
| `tmp/` | Air build output | **NEVER** |

## Module anatomy

Every business module follows this structure:

### `dto.go`
- Request DTOs (with `binding:"required,email"` validation tags)
- Response DTOs (with `json:"..."` tags)
- Helper functions like `ToUserResponse(models.User) UserResponse` for safe serialization (hides password hash, etc.)

### `internal/models/<model>.go`
- GORM model struct (with `gorm:"..."` tags, `BeforeCreate` hook for UUID)
- No DTOs, no business logic — pure DB schema

### `store.go`
- Exported `Store` **interface** defining all DB operations
- Unexported `gormStore` struct implementing it
- `NewStore(db *gorm.DB) Store` constructor
- Methods return raw GORM errors — the service layer converts them to API errors

### `service.go`
- Exported `Service` **interface** with business operations
- Unexported struct implementing it
- Constructor takes `Store` interface (not concrete type) — enables testing
- All errors returned as `*apierr.APIError` (never raw GORM errors)
- `gorm.ErrRecordNotFound` → `apierr.ErrNotFound`
- Service sets `instance` to `""` — the handler fills it with the request path

### `handler.go`
- Exported `Handler` struct with methods matching Gin signature
- Constructor takes `Service` interface
- Parses request, calls service, sends response
- Fills `err.Instance` with `c.Request.URL.Path` before calling `response.Problem()`
- Uses `response.ValidationError()` for binding failures
- Has Swagger annotations (`@Summary`, `@Tags`, `@Router`, etc.) directly above each method

### `routes.go`
- Single exported function: `RegisterRoutes(rg *gin.RouterGroup, handler *Handler)`
- Groups routes under a sub-path (e.g., `/users`, `/auth`)
- Protected routes: called inside a `gin.RouterGroup` that already has `middleware.Auth()`
- Public routes (auth): called on a group without auth middleware

## Wiring in internal/server/server.go

Pattern for adding a new module:
```go
xStore    := xmod.NewStore(db)
xService  := xmod.NewService(xStore)
xHandler  := xmod.NewHandler(xService)
// ...
xmod.RegisterRoutes(protected, xHandler)
```

Always create stores first, then services, then handlers. Modules that depend on other modules pass the interface (e.g., `auth.NewService(userStore, ...)`).

## GORM AutoMigrate

All models must be registered in `internal/database/database.go` → `db.AutoMigrate()`. Tables are created/updated on every server start.

## Swagger

- General API info annotations live in `main.go` (title, version, host, securityDefinitions)
- Endpoint annotations live on handler methods in each module's `handler.go`
- After ANY change to annotations → run `make swagger`
- The blank import `_ "github.com/felipe/dev-test-api/docs"` in `main.go` registers the generated docs
