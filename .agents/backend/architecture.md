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

## Concurrency

### Connection pool (`internal/database/database.go`)

The underlying `*sql.DB` is configured with safe limits right after `gorm.Open()`:

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
sqlDB.SetConnMaxIdleTime(1 * time.Minute)
```

Adjust these values based on production load, but never leave them at the Go defaults (unlimited).

### Graceful shutdown (`internal/server/server.go`)

The server runs in a goroutine. The main goroutine blocks on a signal channel (`os.Signal`). On `SIGINT`/`SIGTERM`, `http.Server.Shutdown()` drains in-flight requests with a 10-second timeout.

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
go func() { srv.ListenAndServe() }()
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

### Parallel independent queries (`golang.org/x/sync/errgroup`)

When a service method makes **two or more store calls that don't depend on each other**, run them in parallel with `errgroup.Group`:

```go
import "golang.org/x/sync/errgroup"

func (s *someService) DoWork(id uuid.UUID) (*Result, error) {
    var itemA *models.ItemA
    var itemB []models.ItemB

    g := new(errgroup.Group)

    g.Go(func() error {
        var err error
        itemA, err = s.store.FindA(id)
        if err != nil {
            if err == gorm.ErrRecordNotFound {
                return apierr.ErrNotFound("ItemA", "")
            }
            return apierr.ErrInternal("Error al obtener itemA", "")
        }
        return nil
    })

    g.Go(func() error {
        var err error
        itemB, err = s.store.FindB(id)
        if err != nil {
            return apierr.ErrInternal("Error al obtener itemB", "")
        }
        return nil
    })

    if err := g.Wait(); err != nil {
        return nil, err
    }

    // Both itemA and itemB are ready — proceed with dependent logic.
}
```

**Rules:**
- Declare result variables outside the closures so they're visible after `g.Wait()`.
- Convert GORM errors to `*apierr.APIError` inside each closure — the service must never return raw GORM errors.
- `g.Wait()` returns the **first** error encountered. Subsequent goroutines are cancelled via the errgroup's internal context.
- Only parallelize queries that are **truly independent** — if query B needs the result of query A, keep them sequential.
- This pattern requires the connection pool configured above so goroutines can borrow separate connections.

**Where this already applies:**
- `sessions/service.go` → `NextQuestion()` — `FindByID` and `FindAnsweredQuestionIDs` run in parallel.

### Fire-and-forget (for non-critical side effects)

When a side effect doesn't affect the response (e.g., logging, analytics, non-blocking progress updates), use a bare goroutine:

```go
go func() {
    s.analyticsService.Track(userID, event)
}()
```

Do **not** use this for operations where the caller needs the result or error. Prefer `errgroup` when you need to wait for completion.

## Swagger

- General API info annotations live in `main.go` (title, version, host, securityDefinitions)
- Endpoint annotations live on handler methods in each module's `handler.go`
- After ANY change to annotations → run `make swagger`
- The blank import `_ "github.com/felipe/dev-test-api/docs"` in `main.go` registers the generated docs
