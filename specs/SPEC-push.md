# SPEC — push

| Campo | Valor |
|---|---|
| Module id | `push` |
| Origen | Alcance nuevo mobile (D-26) |
| Estado | `ready` |
| Depende de | `notifications` (D-28: la tabla de avisos es la fuente) |
| Consumido por | `mobile-push` |
| Migración goose | `00008_device_tokens.sql` |
| Dependencia nueva | **Ninguna** (HTTP con `net/http` contra la Expo Push API) |

## Objective

Que un usuario reciba en el teléfono el mismo aviso que crea el job diario de `notifications` ("Tienes N preguntas para repasar"), sin abrir la app.

**Historias**:
- Al abrir la app con sesión iniciada, el dispositivo queda registrado y puede recibir avisos.
- Cuando el job diario crea las notificaciones, se envían a los dispositivos de esos usuarios.
- Si Expo responde que un token ya no es válido, se borra: no se reintenta para siempre.

**Éxito**: el job diario sigue funcionando aunque Expo esté caído (el aviso in-app se crea igual) y un token muerto deja de intentarse.

## Tech stack

Go + Gin + GORM + Postgres · Expo Push API (`https://exp.host/--/api/v2/push/send`) · sin SDK nuevo.

## Commands

```bash
make migrate-create name=device_tokens   # 00008_device_tokens.sql
make migrate-up && make migrate-status
make check && make test && make swagger
make build
```

## Project structure

```
internal/models/device_token.go                      # NUEVO
internal/database/migrations/00008_device_tokens.sql # NUEVO
internal/modules/push/
├── request.go      # RegisterTokenRequest, UnregisterTokenRequest
├── response.go     # (reutiliza response.Success; sin DTO propio si no hace falta)
├── store.go        # Store + gormStore (upsert, delete, list por usuario)
├── service.go      # Service (Register/Unregister) + Pusher
├── sender.go       # ExpoSender: HTTP a la Expo Push API
├── handler.go      # register / unregister + swagger
├── routes.go       # POST /push/tokens, POST /push/tokens/unregister
└── *_test.go

internal/server/server.go            # wiring: push antes de notifications
internal/modules/notifications/service.go  # llama al Pusher tras crear los avisos
```

## Contrato

### Migración `00008_device_tokens.sql`

```sql
-- +goose Up
CREATE TABLE device_tokens (
    id           uuid         PRIMARY KEY,
    user_id      uuid         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token        varchar(255) NOT NULL,
    platform     varchar(10)  NOT NULL,
    created_at   timestamptz  NOT NULL DEFAULT now(),
    updated_at   timestamptz  NOT NULL DEFAULT now(),
    last_seen_at timestamptz  NOT NULL DEFAULT now()
);

-- Un token identifica un dispositivo: no puede estar en dos usuarios.
CREATE UNIQUE INDEX device_tokens_token_uniq ON device_tokens (token);
CREATE INDEX device_tokens_user_idx ON device_tokens (user_id);

COMMENT ON TABLE  device_tokens             IS 'Tokens de push (Expo) por dispositivo. Un token pertenece a un solo usuario.';
COMMENT ON COLUMN device_tokens.token       IS 'ExponentPushToken[...] devuelto por expo-notifications. Único global.';
COMMENT ON COLUMN device_tokens.platform    IS 'Plataforma del dispositivo: ios o android.';
COMMENT ON COLUMN device_tokens.last_seen_at IS 'Última vez que la app registró/refrescó el token. Sirve para detectar tokens abandonados.';
-- +goose Down
DROP TABLE device_tokens;
```

`platform` validado con `binding:"required,oneof=ios android"` (D-20: columna varchar, validación en Go).

### Endpoints (protegidos, cualquier usuario autenticado)

```
POST /api/v1/push/tokens              body: { token, platform }   -> 200 { data: { registered: true } }
POST /api/v1/push/tokens/unregister   body: { token }             -> 200 { data: { removed: true } }
```

- `POST /push/tokens` es **idempotente** y hace *upsert* por `token`: si el token ya existía con otro usuario (dispositivo compartido / re-login), se reasigna al usuario actual.
- `POST /push/tokens/unregister` se llama al cerrar sesión. Borrar un token ajeno no da error (responde `removed: false`): no se filtra información.
- `DELETE` con body es incómodo en algunos clientes; por eso el unregister es `POST` (coherente con `/notifications/read-all`).

### Service y entrega

```go
// Interfaz que consume notifications (D-28). push la implementa.
type Pusher interface {
	Push(ctx context.Context, items []models.Notification) error
}

type Service interface {
	Register(userID uuid.UUID, token, platform string) error
	Unregister(userID uuid.UUID, token string) error
}
```

- `Service` **no** conoce Expo: eso vive en `sender.go`.
- `Pusher` implementa la interfaz de `notifications`. Si en el futuro aparece otro canal (email), se compone igual.
- `Push` agrupa por token y envía en lotes de **máx. 100** (límite de Expo), con `context` con timeout (10 s por lote).
- Respuesta de Expo: si un *ticket* trae `details.error == "DeviceNotRegistered"`, se borra ese token del store y se loguea `❌ Token push inválido eliminado`.
- Cualquier fallo de red o HTTP se **loguea y se descarta**: el aviso in-app ya está persistido, y el push es best-effort (no se reintenta; el día siguiente vuelve a intentarse).

### Configuración

| Variable | Default | Uso |
|---|---|---|
| `EXPO_PUSH_URL` | `https://exp.host/--/api/v2/push/send` | Endpoint (permite apuntar a un fake en tests) |
| `EXPO_ACCESS_TOKEN` | vacío | Opcional; si está, se envía como `Authorization: Bearer` |

Ambas en `internal/config` (`PushConfig`). Sin secretos obligatorios.

### Integración con `notifications`

`notifications.Service.RunDailyJob()` recibe un `Pusher` (puede ser `nil`) y, **después** de insertar los avisos, llama a `Push(ctx, inserted)`. Si es `nil`, no hace nada (modo web-only). El aviso in-app nunca depende del resultado del push.

Orden de wiring en `server.go`: `pushStore → pushService → notificationStore → notificationService(pusher: pushService)`.

## Code style

- Módulo con la anatomía estándar; handlers con `getUserID` propio del módulo (o reutilizando el patrón existente) y `e.Instance = c.Request.URL.Path`.
- Errores de validación: `response.ValidationError`; upsert/borrado → `apierr.ErrInternal("Error al registrar el dispositivo", "")`.
- Logs con emoji: `✅` lote enviado, `❌` fallo de envío, `🗑️` token eliminado.
- Sin `errgroup` (D-13); si hay que esperar varios lotes, `sync.WaitGroup` + `sync.Mutex`.

## Testing strategy

| Nivel | Archivo | Qué cubre |
|---|---|---|
| Unit sender | `push/sender_test.go` | Con un `httptest.Server` como fake de Expo: cuerpo correcto (array de mensajes con `to/title/body/data`), lotes de 100, timeout, `DeviceNotRegistered` dispara borrado, HTTP 500 no rompe |
| Unit service | `push/service_test.go` | upsert, reasignación de token a otro usuario, unregister idempotente, unregister ajeno no borra |
| Unit handler | `push/handler_test.go` | 401 sin token, 422 con `platform` inválida, 200 en alta y baja |
| Integración | — | Contra la DB real en transacción con rollback: única del token y `ON DELETE CASCADE` |

Store mockeado en unitarios. Sin llamadas reales a Expo.

## Boundaries

- **Always**: migración goose nueva con `COMMENT ON`; upsert idempotente por token; push best-effort (nunca bloquea ni falla el aviso in-app); `Pusher` opcional (`nil` deshabilitado).
- **Ask first**: añadir un tercer `platform`; persistir el histórico de envíos; añadir reintentos con backoff; cambiar el límite de lote.
- **Never**: credenciales de Firebase/APNs (D-27); bloquear el job diario esperando a Expo; guardar tokens en `users`; devolver si un token pertenece a otro usuario; `AutoMigrate`; editar `docs/` a mano.

## Success Criteria

1. `make migrate-up` aplica `00007` y `make migrate-down` la revierte.
2. Registrar dos veces el mismo token deja **una** fila (upsert) y actualiza `last_seen_at`.
3. Registrar con el mismo token desde otro usuario lo reasigna (no duplica).
4. `POST /push/tokens` con `platform=web` → `422` RFC 9457.
5. Con el fake de Expo caído, `RunDailyJob` **sí** crea los avisos in-app y loguea el fallo del push.
6. Un ticket `DeviceNotRegistered` elimina el token.
7. `docs/swagger.yaml` incluye los 2 paths.
8. Sin dependencias nuevas en `go.mod`.

## Notas de implementación (2026-10-20)

- **`Unregister` devuelve `(bool, error)`** (la spec declaraba solo `error`): la respuesta necesita distinguir `removed: false` de `removed: true` sin revelar si el token es ajeno.
- **Upsert**: `ON CONFLICT (token) DO UPDATE` reasigna `user_id`/`platform` y refresca `last_seen_at`; el `id` y `created_at` se conservan. Verificado contra la DB real, junto con `ON DELETE CASCADE`, en transacción con rollback.
- **Sender**: lotes de 100, timeout de 10 s por lote y `Authorization: Bearer` opcional. Los tickets `DeviceNotRegistered` se borran (`🗑️`); los errores HTTP/red se propagan al job, que solo loguea `⚠️` y conserva el aviso in-app.
- **Wiring**: `pushStore → pushService → notificationStore → notificationService(pusher: pushService)`, y `push.RegisterRoutes` antes de `notifications`.
- **Config**: `PushConfig{ExpoURL, AccessToken}` con defaults (`exp.host`); `EXPO_*` documentadas en `.env.example`.
- Tests: **17 nuevos** en `push`; `make check`/`test`/`build` en verde; 2 rutas responden 401 sin token y `docs/swagger.yaml` incluye ambos paths.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Push también al crear notificaciones que no sean del job diario? | Sí, por eso es un `Pusher` genérico; hoy solo hay `daily_review`. |
| OQ-2 | ¿Purga de tokens con `last_seen_at` viejo (> 90 días)? | Fuera de alcance; `DeviceNotRegistered` ya limpia los muertos. Si molesta, tarea aparte. |
| OQ-3 | ¿Envío síncrono dentro del job o worker propio? | **Síncrono con timeout** por lote: el job corre de madrugada y el volumen es bajo. Si crece, se extrae a worker con el patrón de `generation/worker.go`. |
