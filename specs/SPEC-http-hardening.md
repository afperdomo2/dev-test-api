# SPEC — http-hardening

| Campo | Valor |
|---|---|
| Module id | `http-hardening` |
| Origen | Revisión de seguridad 2026-10-21 (M1, M2, M4, L2) |
| Estado | `ready` |
| Depende de | — |
| Consumido por | `auth-hardening` (proxies confiables + body limit), `client-hardening` (headers) |
| Migración goose | Ninguna |
| Dependencia nueva | Ninguna |
| Cambio de contrato | Sí: 413 en import; `/swagger` deja de existir en release por defecto |

## Objective

El servidor no tiene defensas básicas de transporte: `http.Server` se crea solo
con `Addr` (`server.go:158-161`), no hay límite de tamaño de request, Swagger se
sirve sin autenticación en release (`server.go:125`), Gin confía en todos los
proxies y no se emite ninguna cabecera de seguridad.

**Historias**

- Como operador, quiero que un request de 5 GB no tumbe el proceso.
- Como operador, quiero que un cliente lento no agote conexiones.
- Como operador, quiero que la documentación interna no sea pública en producción.
- Como usuario, quiero que el navegador aplique `nosniff`/`frame-ancestors` sobre
  las respuestas de la API.

**Éxito**

- Un body mayor que el límite se rechaza sin reservar memoria proporcional.
- `GIN_MODE=release` sin `SWAGGER_ENABLED=true` devuelve 404 en `/swagger/index.html`.
- Todas las respuestas llevan `X-Content-Type-Options`, `X-Frame-Options`,
  `Referrer-Policy`, `Content-Security-Policy` y (si aplica) `Strict-Transport-Security`.
- `c.ClientIP()` no se puede falsear con `X-Forwarded-For` cuando no hay proxy.

## Tech stack

Go `net/http`, Gin, `gin-contrib/cors` (sin cambios). Sin dependencias nuevas.

## Commands

```bash
make check && make test && make build
make swagger

# verificación manual de cabeceras
curl -sI http://localhost:8080/health | grep -Ei 'x-content-type|frame|referrer|content-security|strict-transport'

# swagger gate
GIN_MODE=release SWAGGER_ENABLED=false make run &
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/swagger/index.html   # 404
GIN_MODE=release SWAGGER_ENABLED=true make run &
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/swagger/index.html   # 200

# body limit (debe ser 413/422 sin OOM)
head -c 20000000 /dev/zero | curl -s -o /dev/null -w '%{http_code}\n' \
  -X POST http://localhost:8080/api/v1/questions/import \
  -H "Authorization: Bearer $TOKEN" -F "file=@-;filename=big.csv"   # 413
```

## Project structure

```
internal/config/config.go              # + SwaggerEnabled, HSTS, TrustedProxies, límites
internal/middleware/security.go        # NUEVO: SecurityHeaders
internal/middleware/bodylimit.go       # NUEVO: BodyLimit + IsPayloadTooLarge
internal/middleware/security_test.go   # NUEVO
internal/middleware/bodylimit_test.go  # NUEVO
internal/server/server.go              # timeouts, SetTrustedProxies, gate Swagger, wiring
internal/modules/questions/handler.go # import: límite + 413
pkg/apierr/errors.go                   # + ErrPayloadTooLarge (413)
.env.example                           # + TRUSTED_PROXIES, SWAGGER_ENABLED, HSTS_ENABLED, MAX_BODY_BYTES
docs/                                  # generado si cambia el contrato de import
```

## Contrato

### 1. Configuración (`.env.example`)

```dotenv
# HTTP hardening
MAX_BODY_BYTES=1048576          # 1 MiB para JSON (login, respuestas, etc.)
TRUSTED_PROXIES=                # IPs/CIDR del reverse proxy, separadas por coma. Vacío = no confiar en X-Forwarded-For
HSTS_ENABLED=true               # solo tiene efecto si el borde sirve HTTPS
SWAGGER_ENABLED=                # vacío = true en debug, false en release
```

`config.Load()`:

- `SwaggerEnabled`: si la env está vacía, `GinMode != "release"`.
- `TrustedProxies`: `[]string` (split + TrimSpace, descartando vacíos).
- El resto con defaults de arriba; `MAX_BODY_BYTES` se valida > 0 (si no, default).

### 2. `http.Server` (server.go)

```go
srv := &http.Server{
    Addr:              ":" + cfg.Port,
    Handler:           r,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       30 * time.Second,
    WriteTimeout:      60 * time.Second,
    IdleTimeout:       120 * time.Second,
    MaxHeaderBytes:    1 << 20,
}
```

Notas:

- `WriteTimeout: 60s` es suficiente: la generación con IA es asíncrona (worker);
  el import máximo son 200 preguntas (`UpdateUserRequest.DailyImportLimit`).
- Los valores se declaran como constantes del paquete `server` para poder
  testearlos por reflexión si se desea.

### 3. Proxies confiables

```go
if len(cfg.TrustedProxies) == 0 {
    if err := r.SetTrustedProxies(nil); err != nil {
        log.Printf("⚠️ No se pudo desactivar la confianza en proxies: %v", err)
    }
} else if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
    log.Fatalf("❌ TRUSTED_PROXIES inválido: %v", err)
}
```

`log.Fatalf` cuando el operador configura algo inválido: mejor fallar al arrancar
que arrancar con IPs falseables.

### 4. Body limit

```go
// internal/middleware/bodylimit.go
func BodyLimit(max int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Body != nil {
            c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max)
        }
        c.Next()
    }
}

func IsPayloadTooLarge(err error) bool {
    var mbe *http.MaxBytesError
    return errors.As(err, &mbe)
}
```

Wiring global, antes de las rutas:

```go
r.Use(middleware.BodyLimit(cfg.MaxBodyBytes))
```

Consecuencia esperada: `ShouldBindJSON` devuelve error y los handlers responden
422 con el detalle del binding (el body está acotado, que es el objetivo de
seguridad). **Excepción**: `questions.Import` debe responder **413**.

### 5. Import con límite explícito (questions/handler.go)

```go
const maxImportBytes int64 = 8 << 20 // 8 MiB

if fh, err := c.FormFile("file"); err == nil {
    if fh.Size > maxImportBytes {
        response.Problem(c, apierr.ErrPayloadTooLarge(
            "El archivo supera el máximo de 8 MiB", c.Request.URL.Path))
        return
    }
    f, err := fh.Open()
    if err != nil { /* 422 como hoy */ }
    defer f.Close()

    buf := new(bytes.Buffer)
    buf.Grow(int(fh.Size))
    n, err := io.Copy(buf, io.LimitReader(f, maxImportBytes+1))
    if err != nil || n > maxImportBytes {
        response.Problem(c, apierr.ErrPayloadTooLarge(
            "El archivo supera el máximo de 8 MiB", c.Request.URL.Path))
        return
    }
    reader = buf
}
```

Además: si `c.FormFile` falla por `MaxBytesError` (body global), devolver 413 en
vez de caer al `PostForm("content")` vacío:

```go
if fh, err := c.FormFile("file"); err != nil && middleware.IsPayloadTooLarge(err) {
    response.Problem(c, apierr.ErrPayloadTooLarge("Request demasiado grande", c.Request.URL.Path))
    return
} else if err == nil { /* … */ }
```

`pkg/apierr/errors.go`:

```go
func ErrPayloadTooLarge(detail, instance string) *APIError {
    return NewAPIError(http.StatusRequestEntityTooLarge, "Payload Too Large", detail, instance)
}
```

Swagger de `POST /questions/import`: añadir `@Failure 413`.

### 6. Cabeceras de seguridad

```go
// internal/middleware/security.go
func SecurityHeaders(hsts bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        h := c.Writer.Header()
        h.Set("X-Content-Type-Options", "nosniff")
        h.Set("X-Frame-Options", "DENY")
        h.Set("Referrer-Policy", "no-referrer")
        h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        // API JSON: no debe cargar ni ejecutar nada
        h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
        if hsts {
            h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }
        c.Next()
    }
}
```

Se registra **antes** de CORS y rutas:

```go
r.Use(middleware.SecurityHeaders(cfg.HSTSEnabled))
r.Use(middleware.Logger())
r.Use(middleware.CORS(cfg.Cors.AllowedOrigins))
```

> La CSP de la SPA **no** es responsabilidad de este módulo (la sirve el hosting
> estático). Ver `SPEC-client-hardening.md`.

### 7. Gate de Swagger

```go
if cfg.SwaggerEnabled {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
} else {
    log.Println("🔒 Swagger deshabilitado (release)")
}
```

Sin ruta registrada → 404 de Gin (sin filtrar que existe).

## Testing strategy

| # | Test | Archivo | Esperado |
|---|---|---|---|
| T1 | `SecurityHeaders(true)` | `security_test.go` | 6 cabeceras presentes, incluida HSTS |
| T2 | `SecurityHeaders(false)` | idem | sin HSTS, resto igual |
| T3 | `BodyLimit(1024)` + `io.ReadAll(body)` | `bodylimit_test.go` | error `*http.MaxBytesError` |
| T4 | `IsPayloadTooLarge(err)` | idem | true/false correctos |
| T5 | Handler que hace `ShouldBindJSON` sobre body 2× límite | idem | 422 (no panic, no OOM) |
| T6 | `questions.Import` con `fh.Size > 8 MiB` | `questions/handler_test.go` | 413 |
| T7 | `questions.Import` con body global excedido | idem | 413 |
| T8 | `config.Load` con `SWAGGER_ENABLED` vacío + `GIN_MODE=release` | `config` (si hay tests) o manual | `false` |
| T9 | `TRUSTED_PROXIES` inválido | manual | el proceso no arranca |

Verificación manual obligatoria (no automatizable sin levantar el server real):
cabeceras con `curl -I`, 404 de Swagger en release, y que `curl` con un body de
20 MB sobre `/questions/import` devuelve 413 sin que el RSS se dispare.

## Boundaries

- **Always**: `make check` + `make test`; cualquier endpoint nuevo hereda el body
  limit global; documentar los límites en `.env.example`.
- **Ask first**: bajar `MAX_BODY_BYTES` por debajo de 1 MiB (rompe payloads
  legítimos), añadir `TRUSTED_PROXIES` por defecto, cambiar `WriteTimeout`.
- **Never**: usar `c.Request.Body` sin el middleware en un endpoint nuevo,
  `SetTrustedProxies([]string{"0.0.0.0/0"})`, servir Swagger en release sin
  decisión explícita, poner CSP de la SPA en la API.

## Success criteria

1. `curl -I /health` devuelve las 5 cabeceras (+ HSTS cuando `HSTS_ENABLED=true`).
2. `GIN_MODE=release` + `SWAGGER_ENABLED` vacío → `/swagger/index.html` = 404.
3. Un POST de 20 MB a `/questions/import` devuelve 413 y el proceso no supera
   ~50 MB de RSS incremental (verificado con `docker stats` o `ps`).
4. `TRUSTED_PROXIES` vacío → `c.ClientIP()` devuelve la IP de la conexión aunque
   se envíe `X-Forwarded-For: 1.2.3.4`.
5. `make test` verde; sin migraciones.

## Open questions

- ¿Qué borde TLS hay en producción (nginx, Caddy, Cloudflare)? De ello depende
  `TRUSTED_PROXIES` y si `HSTS_ENABLED` debe quedar en `true`.
- ¿Se quiere un `Server-Timing`/request-id? Fuera de alcance.
