# SPEC — notifications

| Campo | Valor |
|---|---|
| Module id | `notifications` |
| IDs originales | B2.1–B2.9 |
| Estado | `ready` |
| Depende de | `sessions` (due count existente: `sessions/store.go:210`) |
| Migración goose | `00006_notifications.sql` |
| Dependencia nueva | **Ninguna** (ticker stdlib — D-17) |

## Objective

Que el usuario se entere **dentro de la app** de que tiene repaso pendiente hoy, sin depender de email ni push, y que tenga un centro donde revisar y marcar como leídas sus notificaciones.

**Historias**:
- Como usuario con repasos vencidos, recibo una vez al día una notificación "Tienes N preguntas para repasar" y veo un badge con el número de no leídas en la barra superior.
- Como usuario, abro la campana, veo mi historial paginado, marco una como leída o todas a la vez.
- Como usuario sin repasos pendientes, **no** recibo notificación ese día.

**Éxito**: el job diario es **idempotente** (correrlo dos veces no duplica) y no genera filas para usuarios sin repasos pendientes.

## Tech stack

Go + Gin + GORM + Postgres · Vue 3 + Vuetify + TanStack Query · `@devtest/shared`. Sin dependencias nuevas.

## Commands

```bash
make migrate-create name=notifications   # crea 00006_notifications.sql
make migrate-up                          # aplica
make migrate-status
make migrate-down                        # revertir en desarrollo
make check && make test && make swagger
make fe-format && make fe-check && make fe-lint && make fe-test
pnpm --filter @devtest/shared typecheck
```

## Project structure

```
internal/models/notification.go                     # NUEVO
internal/database/migrations/00006_notifications.sql # NUEVO
internal/modules/notifications/
├── request.go      # ListNotificationsParams (embebe common.PaginationParams)
├── response.go     # NotificationResponse, UnreadCountResponse, MarkAllReadResponse
├── mapper.go       # ToNotificationResponse
├── store.go        # Store + gormStore
├── service.go      # Service + RunDailyJob
├── scheduler.go    # NUEVO: ticker stdlib
├── handler.go      # 4 handlers + swagger
├── routes.go       # RegisterRoutes
├── service_test.go
└── scheduler_test.go

internal/server/server.go                           # wiring + arranque/parada del scheduler

packages/shared/src/types/notification.types.ts      # NUEVO
packages/shared/src/index.ts                         # barrel
apps/web/src/types/notification.types.ts             # re-export fino (NUEVO)
apps/web/src/api/services/notifications.service.ts   # NUEVO
apps/web/src/queries/notifications.queries.ts        # NUEVO
apps/web/src/features/notifications/pages/NotificationsPage.vue  # NUEVO
apps/web/src/layouts/DefaultLayout.vue               # campana + badge
apps/web/src/router/routes.ts                        # + /notifications
```

## Contrato

### Migración `00006_notifications.sql`

```sql
-- +goose Up
CREATE TABLE notifications (
    id           uuid        PRIMARY KEY,
    user_id      uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type         varchar(30) NOT NULL,
    title        varchar(200) NOT NULL,
    body         text        NOT NULL,
    is_read      boolean     NOT NULL DEFAULT false,
    created_at   timestamptz NOT NULL DEFAULT now(),
    creation_day date        NOT NULL GENERATED ALWAYS AS ((created_at AT TIME ZONE 'UTC')::date) STORED
);

-- Idempotencia del job diario: una notificación por (usuario, tipo, día UTC).
CREATE UNIQUE INDEX notifications_user_type_day_uniq
    ON notifications (user_id, type, creation_day);

CREATE INDEX notifications_user_created_idx ON notifications (user_id, created_at DESC);
CREATE INDEX notifications_user_unread_idx  ON notifications (user_id) WHERE is_read = false;

-- D-23: toda tabla y columna lleva su COMMENT ON.
COMMENT ON TABLE  notifications                IS 'Notificaciones in-app por usuario. Generadas por el job diario y por eventos futuros.';
COMMENT ON COLUMN notifications.type           IS 'Tipo de notificación. Hoy solo ''daily_review'' (D-21: string + const Go).';
COMMENT ON COLUMN notifications.title          IS 'Título corto mostrado en la lista.';
COMMENT ON COLUMN notifications.body           IS 'Cuerpo de la notificación.';
COMMENT ON COLUMN notifications.is_read        IS 'Marcada como leída por el usuario.';
COMMENT ON COLUMN notifications.creation_day   IS 'Día UTC de creación. Columna generada; sostiene la unique de idempotencia. NO se escribe desde la aplicación.';

-- +goose Down
DROP TABLE notifications;
```

⚠️ **Verificación obligatoria antes de dar por buena la migración**: una columna `GENERATED ... STORED` exige una expresión `IMMUTABLE`; `AT TIME ZONE 'UTC'` debe serlo. Confirmar con `make migrate-up` contra la DB real. Si Postgres la rechaza, el fallback es declarar `creation_day date NOT NULL` y que la aplicación la escriba (mismo índice único, idempotencia equivalente).

### Modelo

```go
type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	User      *User     `gorm:"foreignKey:UserID" json:"-"`
	Type      string    `gorm:"not null;size:30" json:"type"`
	Title     string    `gorm:"not null;size:200" json:"title"`
	Body      string    `gorm:"not null" json:"body"`
	IsRead    bool      `gorm:"not null;default:false" json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
```

`creation_day` **no** se mapea: la escribe Postgres. GORM no la incluye en el `INSERT`.

### Tipos

```go
// D-21: string + const, como el resto del repo.
const TypeDailyReview = "daily_review"
```

### Endpoints (todos bajo `protected`, cualquier usuario autenticado)

| Método | Ruta | Respuesta |
|---|---|---|
| `GET` | `/api/v1/notifications?page&perPage` | `Paginated` → `NotificationResponse[]` + `Meta` |
| `GET` | `/api/v1/notifications/unread-count` | `{ "data": { "count": 3 } }` |
| `POST` | `/api/v1/notifications/:id/read` | `{ "data": NotificationResponse }` — **idempotente** |
| `POST` | `/api/v1/notifications/read-all` | `{ "data": { "updated": 3 } }` |

`GET /notifications` ordena por `created_at DESC` (no expone `sortBy`: el orden es el del producto).

✅ **Hermanos static+param verificados**: `POST /notifications/read-all` convive con `POST /notifications/:id/read` sin conflicto de wildcard en Gin 1.12.0 (comprobado ejecutando el router real). Se mantienen las rutas de B2.3 tal cual.

Marcar como leída una notificación ajena o inexistente → `404 NotFoundResponse` (no revelar existencia ajena). `read-all` sin filas → `{ "updated": 0 }`, no error.

### Scheduler (`scheduler.go`) — sin dependencias nuevas

```go
const dailyHour = 8 // hora local del servidor

type Scheduler struct {
	service Service
}

func NewScheduler(service Service) *Scheduler { return &Scheduler{service: service} }

// Start bloquea hasta que ctx se cancele. Reutiliza el patrón de
// internal/modules/generation/worker.go (D-13): goroutine de larga vida +
// context cancelable, sin errgroup.
func (s *Scheduler) Start(ctx context.Context) {
	for {
		timer := time.NewTimer(time.Until(nextRun(time.Now())))
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			if err := s.service.RunDailyJob(ctx); err != nil {
				log.Printf("❌ Job diario de notificaciones: %v", err)
			}
		}
	}
}

// nextRun devuelve hoy a las dailyHour:00 locales, o mañana si ya pasó.
func nextRun(now time.Time) time.Time { /* ... */ }
```

Zona horaria: **del servidor** (`time.Local`), hora fija 08:00 configurable como constante. No se añade clave en `settings` (fuera de alcance de B2).

### Job diario (fan-out idempotente)

```go
// store
DueReviewCounts() ([]DueReviewUser, error)              // user_id + n (COUNT de repasos vencidos)
InsertBatch(items []models.Notification) (int64, error) // OnConflict{DoNothing: true}, devuelve insertadas
```

```go
func (s *service) RunDailyJob(ctx context.Context) error {
	due, err := s.store.DueReviewCounts()
	if err != nil {
		return apierr.ErrInternal("Error al calcular los repasos pendientes", "")
	}
	if len(due) == 0 {
		return nil
	}

	items := make([]models.Notification, 0, len(due))
	for _, d := range due {
		items = append(items, models.Notification{
			UserID: d.UserID,
			Type:   TypeDailyReview,
			Title:  "Repaso de hoy",
			Body:   fmt.Sprintf("Tienes %d preguntas para repasar", d.N),
		})
	}

	inserted, err := s.store.InsertBatch(items)
	if err != nil {
		return apierr.ErrInternal("Error al crear las notificaciones", "")
	}
	log.Printf("✅ Job diario de notificaciones: %d nuevas (%d usuarios con repaso pendiente)", inserted, len(due))
	return nil
}
```

- `DueReviewCounts` reutiliza la definición de repaso vencido (D-09): `next_review_at <= now() AND is_mastered = false`, agrupado por `user_id`.
- Solo usuarios con `n > 0` → sin ruido.
- Dos queries + un bulk insert con UUID generados en Go (usa el hook `BeforeCreate`; evita depender de `gen_random_uuid()` y de la versión de Postgres).
- La idempotencia la **garantiza el índice único**, no el código.

### Wiring en `server.go`

```go
notificationStore   := notifications.NewStore(db)
notificationService := notifications.NewService(notificationStore)
notificationHandler := notifications.NewHandler(notificationService)
notificationScheduler := notifications.NewScheduler(notificationService)
```

Registro: `notifications.RegisterRoutes(protected, notificationHandler)` (dentro del grupo con `middleware.Auth`).

Arranque/parada, espejo del worker de generación (`server.go:144-168`):

```go
schedCtx, stopScheduler := context.WithCancel(context.Background())
schedDone := make(chan struct{})
go func() {
	defer close(schedDone)
	notificationScheduler.Start(schedCtx)
}()
// ...
stopScheduler()
select {
case <-schedDone:
case <-time.After(5 * time.Second):
	log.Println("⚠️ El scheduler de notificaciones no terminó a tiempo")
}
```

### Shared + Web

```ts
// packages/shared/src/types/notification.types.ts
export type NotificationType = 'daily_review'

export interface Notification {
  id: string
  userId: string
  type: NotificationType
  title: string
  body: string
  isRead: boolean
  createdAt: string
}

export interface UnreadCount {
  count: number
}
```

```ts
// apps/web/src/api/services/notifications.service.ts
export async function getNotifications(page: number, perPage: number): Promise<PaginatedResponse<Notification>>
export async function getUnreadCount(): Promise<UnreadCount>
export async function markNotificationRead(id: string): Promise<Notification>
export async function markAllNotificationsRead(): Promise<{ updated: number }>
```

```ts
// apps/web/src/queries/notifications.queries.ts
export function notificationsOptions(page: () => number, perPage: () => number) { /* key: ['notifications','list',page,perPage] */ }
export function unreadCountOptions() {
  return queryOptions({
    queryKey: ['notifications', 'unread-count'],
    queryFn: () => notificationsService.getUnreadCount(),
    refetchInterval: 60 * 1000,
    refetchOnWindowFocus: true,
  })
}
export function markReadMutation() { /* invalida ['notifications'] */ }
export function markAllReadMutation() { /* invalida ['notifications'] */ }
```

**UI**:
- `DefaultLayout.vue`: `v-btn icon="mdi-bell-outline"` con `v-badge` (`:content="unreadCount"`, `:model-value="unreadCount > 0"`) entre el toggle de tema y el menú de usuario; navega a `/notifications`. `aria-label="Notificaciones"`.
- `NotificationsPage.vue` (`useQuery` + `useMutation`, `v-list` de items con título/cuerpo/fecha relativa, botón "Marcar todas como leídas", paginación con `PaginatedFooter`, estado vacío, item no leído con indicador **texto + color**, no solo color).
- Ruta `/notifications` en `routes.ts` con `meta: { requiresAuth: true }`.
- El badge **no** consulta para admins/users por separado: es para todos.

## Code style

- `mapper.go` con `ToNotificationResponse(models.Notification) NotificationResponse` (exportado).
- El service fija `instance` a `""`; el handler lo rellena con `c.Request.URL.Path`.
- Store: errores GORM crudos; service los traduce a `*apierr.APIError`.
- Logs con emoji (`✅` job ok, `❌` error, `🚀`/`🛑` si aplica).
- Mobile fuera de alcance (D-18): no se toca `apps/mobile/`.

## Testing strategy

| Nivel | Archivo | Qué cubre |
|---|---|---|
| Unit service | `internal/modules/notifications/service_test.go` | `RunDailyJob` con 0 usuarios → no inserta; con N → construye `N` items con el texto correcto; error de store → `apierr`; lista paginada; `mark read` idempotente; `404` en id ajeno |
| Unit scheduler | `internal/modules/notifications/scheduler_test.go` | `nextRun` (hoy antes/después de las 08:00, cambio de día); `Start` retorna al cancelar el context |
| Unit plain | — | El índice único y la columna generada se verifican contra la DB real (transacción con rollback → nota en la tarea) |
| Web service | `notifications.service.test.ts` | 4 funciones: URL, método y desempaquetado de `{ data }` / `PaginatedResponse` |
| Web query | `notifications.queries.test.ts` | queryKeys estables; mutaciones invalidan `['notifications']` |
| Web layout | `DefaultLayout.test.ts` | la campana aparece; el badge muestra el contador y se oculta en 0 |
| Web page | `NotificationsPage.test.ts` | render de lista, estado vacío, "marcar todas" |

Store mockeado. Sin DB en tests.

## Boundaries

- **Always**: migración goose nueva para cualquier cambio de esquema, con `COMMENT ON` (D-23/D-24); `make check` + `make test` + `make swagger`; `make fe-*` + `make fe-test`; idempotencia garantizada por índice único.
- **Ask first**: cambiar `dailyHour` o la zona horaria del job; añadir un segundo `Type`; añadir retención/purga; exponer `sortBy` en el listado.
- **Never**: `robfig/cron` u otra dependencia de scheduling (D-17); `errgroup` (D-13); notificar a usuarios sin repasos pendientes; escribir `creation_day` desde Go; editar `docs/` a mano; tocar mobile.

## Success Criteria

1. `make migrate-up` aplica `00006_notifications.sql` y `make migrate-down` la revierte limpiamente.
2. Ejecutar `RunDailyJob` dos veces el mismo día inserta filas la primera vez y **0** la segunda (idempotencia por índice único).
3. Un usuario sin repasos vencidos no recibe ninguna fila.
4. El texto es `Tienes N preguntas para repasar` con el `N` real del usuario.
5. `GET /notifications` devuelve `Paginated` ordenado por `createdAt` descendente; `GET /notifications/unread-count` devuelve el número de no leídas.
6. `POST /notifications/:id/read` sobre una notificación ajena devuelve `404`; sobre la propia, es idempotente.
7. Arrancar y parar el servidor no deja goroutines colgadas y loguea el cierre del scheduler.
8. La campana con badge aparece en la barra superior y navega a `/notifications`; `make fe-check` + `make fe-lint` + `make fe-test` pasan.
9. `docs/swagger.yaml` incluye los 4 paths nuevos.
10. No se añadió ninguna dependencia a `go.mod`.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Política de retención (purgar >90 días)? | **Fuera de alcance** de B2. La tabla crece lento (1 fila/usuario/día con repasos). Si molesta, se abre una tarea aparte con índice + job de purga. |
| OQ-2 | ¿Hora y zona del job? | 08:00 `time.Local` del servidor, constante en `scheduler.go`. Sin clave en `settings`. |
| OQ-3 | ¿El badge sondea cada 60 s o solo al enfocar? | `refetchInterval: 60s` + `refetchOnWindowFocus`. Ajustable sin coste. |
| OQ-4 | ¿`read-all` devuelve el número actualizado? | Sí: `{ "updated": n }` — permite feedback inmediato en la UI. |
