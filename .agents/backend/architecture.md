# Architecture rules

## Directory ownership

| Path | Purpose | Editable? |
|------|---------|-----------|
| `internal/models/` | GORM structs (pure DB schema) | Yes — add new models here |
| `internal/modules/<module>/` | Business modules (auth, users, questions) | Yes — add features here |
| `internal/server/server.go` | Gin engine + middleware + DI wiring + Run() | Yes — register new modules here |
| `main.go` | Server entrypoint (config.Load + db.Connect + server.Run) | Rarely |
| `internal/config/config.go` | Env var loading | Rarely — add new env vars only |
| `internal/database/database.go` | GORM connection + goose migrations + seed | Rarely |
| `internal/middleware/` | Cross-cutting: JWT auth, logging | Rarely |
| `pkg/` | Shared utilities: response, apierr | Rarely |
| `docs/` | Swagger auto-generated | **NEVER** — use `make swagger` |
| `tmp/` | Air build output | **NEVER** |

## Modules (`internal/modules/`)

| Module | Purpose |
|--------|---------|
| `auth` | Auth: `POST /auth/setup` (primer admin), `POST /auth/login`, `GET /auth/status` |
| `users` | CRUD usuarios (admin) |
| `topics` | CRUD temas |
| `questions` | CRUD preguntas + importación CSV (la generación con IA vive en `generation`) |
| `sessions` | Sesiones de práctica (quiz) + sesión diaria `POST /sessions/daily` |
| `progress` | Progreso SRS del usuario |
| `tips` | Tips del día (admin CRUD + import CSV) + `GetTipOfTheDay` consumido por la sesión diaria |
| `moderation` | Reportes de contenido: `POST /reports` (usuario), listado/revisión (admin) |
| `gamification` | Gamificación: `GET /overview`, `PUT /daily-goal`, `POST /streak/freeze` |
| `opencode` | Admin: `GET /opencode/usage` |
| `settings` | Configuración tunable en runtime: `GET/PUT /settings` (admin) |
| `generation` | Jobs de generación IA por lotes (admin): `POST/GET /generation`, `GET /generation/stats`, `GET /generation/{id}` + worker pool en `worker.go` |

Servicios compartidos: `internal/services/ai` (generación IA: cliente + batch), `internal/services/opencode` (client).

Data model (tablas DB): ver `data-model.md`.

## Module anatomy

Every business module follows this structure:

### `request.go`
- Request DTOs (with `binding:"required,email"` validation tags)
- `sortConfig` variable (`common.SortConfig`) for allowed sort fields and default order
- `ListXxxParams` struct (if the module has a paginated list endpoint) that embeds `common.PaginationParams` and adds module-specific filter fields (e.g., `Search`, `MyOnly`, `Type`, `Difficulty`, `TopicIDs`). The handler parses query params into this struct, then passes it to service → store.

### `response.go`
- Response DTOs (with `json:"..."` tags)
- Sub-response types nested inside responses (e.g., `OptionResponse`, `TopicInfo`)
- Both full detail (`XxxResponse`) and list (`XxxListResponse`) variants when they differ

### `mapper.go`
- Helper functions like `ToUserResponse(models.User) UserResponse` for safe serialization (hides password hash, etc.)
- Always exported (`ToXxxResponse`, not `toXxxResponse`)
- Modules with no mapper functions (e.g., `auth`, `progress`) omit this file

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

## Goose migrations

Schema changes are **versioned migrations**, not AutoMigrate. The database is migrated on every server start and via `make migrate-*`.

- Source of truth for the schema: `internal/database/migrations/*.sql` (base `00001_init.sql`).
- Migration runner: `internal/database/migrate.go` — embeds `migrations/*.sql`, runs `goose` provider with a Postgres session locker (safe under concurrent startups).
- Wired in `internal/database/database.go` → `runMigrations(sqlDB)` before `seedDefaultTopics(db)`.
- **Rule:** DDL/DML de esquema va SIEMPRE como migración nueva (`make migrate-create name=<desc>`), nunca como `db.AutoMigrate()` ni edición del baseline una vez aplicado.
- CLI: `make migrate-up|down|status|create name=<desc>` (carga `.env`, construye la DSN key-value). Scripts: `scripts/migrate.sh` (Unix/Git Bash/WSL) y `scripts/migrate.ps1` (Windows).
- Convenciones: funciones SQL embebidas con `-- +goose StatementBegin/End`; `Down` en orden topológico inverso; los `COMMENT ON` del modelo viven en las migraciones.

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

### Parallel independent queries (`sync.WaitGroup` + `sync.Mutex`)

When a service method makes **two or more store calls that don't depend on each other**, run them in parallel with `sync.WaitGroup` and collect results/errors under a `sync.Mutex`.

**Do not use `golang.org/x/sync/errgroup`**: es una dependencia nueva y el repo la evita deliberadamente (ver `specs/DECISIONS.md` → D-13). El patrón válido es el de `sessions/service.go` → `gatherDailyData()`:

```go
type dailyData struct {
	dueCount int64
	tip      *tips.TipResponse
	errs     []error
}

func (s *someService) gatherData(userID uuid.UUID) dailyData {
	var (
		wg  sync.WaitGroup
		mu  sync.Mutex
		res dailyData
	)

	// collect envuelve cada llamada: marca el WaitGroup y acumula el error.
	collect := func(fn func() error) {
		defer wg.Done()
		if err := fn(); err != nil {
			mu.Lock()
			res.errs = append(res.errs, err)
			mu.Unlock()
		}
	}

	wg.Add(2)
	go collect(func() error {
		n, err := s.store.CountDueReviews(userID)
		mu.Lock()
		res.dueCount = n
		mu.Unlock()
		return err
	})
	go collect(func() error {
		tip, err := s.tipService.GetTipOfTheDay(userID)
		if err != nil {
			// Un 404 esperado no es un error.
			var apiErr *apierr.APIError
			if errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound {
				return nil
			}
			return err
		}
		mu.Lock()
		res.tip = tip
		mu.Unlock()
		return nil
	})
	wg.Wait()

	return res
}
```

**Rules:**
- `wg.Done()` siempre con `defer` dentro del closure `collect`.
- **Toda** escritura en el struct compartido va bajo `mu.Lock()` / `mu.Unlock()` — incluidas las asignaciones de resultado, no solo los errores.
- Los errores se acumulan en un slice (`errs []error`) y **el llamador** decide: `if len(res.errs) > 0 { return nil, apierr.ErrInternal(...) }`. No hay cancelación temprana: todas las goroutines terminan.
- Convierte los errores GORM a `*apierr.APIError` antes de guardarlos — el service nunca devuelve errores GORM crudos.
- Solo paraleliza llamadas **realmente independientes**; si la query B necesita el resultado de la A, van secuenciales.
- El patrón requiere el connection pool configurado arriba para que las goroutines tomen conexiones distintas.

**Where this already applies:**
- `sessions/service.go` → `gatherDailyData()` (lo usa `CreateDaily()`) — `CountDueReviews`, `CountWarmupCandidates` y `tipService.GetTipOfTheDay` corren en paralelo; un `404` del tip se traduce en `nil` sin error.

> `NextQuestion()` es **secuencial**: sus lecturas van en serie. No lo "optimices" con concurrencia sin medir antes.

### Fire-and-forget (for non-critical side effects)

When a side effect doesn't affect the response (e.g., logging, analytics, non-blocking progress updates), use a bare goroutine:

```go
go func() {
    s.analyticsService.Track(userID, event)
}()
```

Do **not** use this for operations where the caller needs the result or error. Cuando haya que esperar el resultado, usa el patrón `sync.WaitGroup` + `sync.Mutex` de arriba.

### Background worker pool (`internal/modules/generation/worker.go`)

Patrón canónico para trabajo en segundo plano. **No** usa `errgroup`: son goroutines de larga vida, no un lote acotado (y `errgroup` está prohibido, D-13).

- `Pool` con `N` goroutines `run()` + `sync.WaitGroup`; cada una hace su propio bucle `claim → process → sleep`.
- Cancelación por `context`: cada iteración comprueba `ctx.Err()` y encola cualquier job que ya estaba reclamado para que el shutdown no lo deje huérfano en `running`.
- **Concurrencia ajustable en caliente** desde `settings` (`ai.worker_concurrency`): cada iteración llama a `reconcile(ctx)`, que compara el objetivo con la actual y hace `spawn` o cierra el canal `quit` del worker sobrante.
- Al arrancar, `recoverRunning` devuelve a `pending` los jobs que quedaron en `running` de un proceso anterior.
- El pool se detiene **antes** que el HTTP server en `Run()` (`stopPool` + espera), para no perder trabajo a mitad de una llamada a la IA.
- Cuando se implemente el scheduler de notificaciones (`specs/SPEC-notifications.md`, tarea `NT-5`), reutilizar este patrón (o extraerlo a un worker genérico) en lugar de inventar otro.

## Swagger

- General API info annotations live in `main.go` (title, version, host, securityDefinitions)
- Endpoint annotations live on handler methods in each module's `handler.go`
- After ANY change to annotations → run `make swagger`
- The blank import `_ "github.com/felipe/dev-test-api/docs"` in `main.go` registers the generated docs
