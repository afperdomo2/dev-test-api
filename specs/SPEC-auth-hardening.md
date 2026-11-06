# SPEC — auth-hardening

| Campo | Valor |
|---|---|
| Module id | `auth-hardening` |
| Origen | Revisión de seguridad 2026-10-21 (H3, H4, L3) |
| Estado | `ready` |
| Depende de | `http-hardening` (`SetTrustedProxies` + body limit) |
| Consumido por | — (consume `client-hardening` para las pantallas) |
| Migración goose | Ninguna |
| Dependencia nueva | `golang.org/x/time` (rate limiter oficial de Go) — *Ask First* |
| Cambio de contrato | Sí: `/auth/setup` exige token; `/auth/status` añade `setupRequiresToken` |

## Objective

`POST /api/v1/auth/setup` crea el primer admin si la tabla está vacía
(`auth/service.go:40-52`) y es **pública**: en una instancia recién desplegada,
cualquiera que llegue primero se queda con el rol admin. Además, `POST /auth/login`
no tiene ningún límite de intentos y filtra por timing si un email existe.

**Historias**

- Como operador, quiero que solo quien conozca `SETUP_TOKEN` pueda inicializar la
  instancia.
- Como operador, quiero que el login resista fuerza bruta básica sin bloquear a
  usuarios legítimos.
- Como usuario, no quiero que un atacante pueda saber qué emails están
  registrados midiendo tiempos de respuesta.

**Éxito**

- Sin `SETUP_TOKEN` (o con uno incorrecto), `/auth/setup` responde 403 en **todos**
  los modos, y no se crea ningún usuario.
- Con el token correcto y DB vacía, el setup funciona igual que hoy.
- El intento 6 de login desde la misma IP en una ventana de un minuto responde 429
  con `Retry-After`.
- Email inexistente y contraseña incorrecta devuelven el mismo body y consumen
  tiempo comparable (bcrypt en ambos caminos).

## Tech stack

Go + Gin. `golang.org/x/time/rate` para el token bucket. Sin Redis: limiter en
memoria (single instance, ver Open Questions).

## Commands

```bash
go get golang.org/x/time@latest && go mod tidy
make check && make test && make build && make swagger

# prueba manual del gate
SETUP_TOKEN=$(openssl rand -hex 16) make run &
curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:8080/api/v1/auth/setup \
  -H 'Content-Type: application/json' -d '{"email":"a@b.com","password":"password123"}'          # 403
curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:8080/api/v1/auth/setup \
  -H "X-Setup-Token: $SETUP_TOKEN" -H 'Content-Type: application/json' \
  -d '{"email":"a@b.com","password":"password123","setupToken":"'$SETUP_TOKEN'"}'               # 201

# rate limit
for i in $(seq 1 8); do
  curl -s -o /dev/null -w '%{http_code} ' -X POST localhost:8080/api/v1/auth/login \
    -H 'Content-Type: application/json' -d '{"email":"a@b.com","password":"x"}'
done; echo   # 401 x5 y luego 429
```

## Project structure

```
internal/config/config.go             # + SetupToken, RateLimitAuth*
internal/modules/auth/request.go      # SetupRequest + setupToken
internal/modules/auth/response.go     # StatusResponse + setupRequiresToken
internal/modules/auth/handler.go      # validación de token + wiring del limiter
internal/modules/auth/routes.go       # limiter en setup/login
internal/modules/auth/service.go      # dummy bcrypt en not-found
internal/modules/auth/*_test.go
internal/middleware/ratelimit.go      # NUEVO
internal/middleware/ratelimit_test.go # NUEVO
internal/server/server.go             # wiring (limiter + token)
.env.example                          # + SETUP_TOKEN, AUTH_RATE_LIMIT_*
packages/shared/src/types/auth.types.ts  # setupToken + setupRequiresToken
apps/web/src/features/auth/pages/SetupPage.vue
apps/mobile/src/features/auth/screens/SetupScreen.tsx
apps/web/e2e/fixtures.ts + critical-flow.spec.ts
docs/                                 # make swagger
```

## Contrato

### 1. `SETUP_TOKEN` (fail closed)

Config:

```go
type Config struct {
    // …
    SetupToken string // SETUP_TOKEN
    RateLimitAuthPerMinute int // AUTH_RATE_LIMIT_PER_MINUTE (default 10; 0 = deshabilitado)
    RateLimitAuthBurst     int // AUTH_RATE_LIMIT_BURST (default 5)
}
```

Reglas:

- `SETUP_TOKEN` vacío ⇒ `/auth/setup` **siempre** 403, en debug y release.
- `SETUP_TOKEN` presente ⇒ se exige coincidencia; si no coincide o falta ⇒ 403
  con el mismo mensaje y body (no distinguir "token ausente" de "token inválido").
- Si viene informado con **menos de 16 caracteres** ⇒ `log.Fatal` en
  `config.Load` (evita `SETUP_TOKEN=123`). Vacío no es fatal aquí: lo decide el
  guard de arranque.
- La comparación es en tiempo constante sobre hashes SHA-256 (evita filtrar la
  longitud del token):

```go
func setupTokenMatches(got, want string) bool {
    if want == "" {
        return false
    }
    a, b := sha256.Sum256([]byte(got)), sha256.Sum256([]byte(want))
    return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}
```

- El token viaja en el **body** (`setupToken`), no en query string (el logger de
  Gin registra el query).
- No se loguea nunca.

**Guard de arranque** (en `server.Run`, con la DB ya conectada):

```go
count, err := userStore.Count()
if err != nil { log.Fatalf("❌ No se pudo comprobar el estado inicial: %v", err) }

if count == 0 && cfg.SetupToken == "" {
    log.Fatal("❌ SETUP_TOKEN es obligatorio con la DB vacía: sin él no se puede crear el primer admin")
}
if count > 0 && cfg.SetupToken == "" {
    log.Println("⚠️ SETUP_TOKEN vacío: /auth/setup deshabilitado (el sistema ya está inicializado)")
}
```

Así una instancia nueva **no arranca en un estado inutilizable**, y una ya
inicializada no queda acoplada a un secreto de un solo uso.

- La carrera `Count()` + `Create()` se **acepta** ahora que el endpoint está
  detrás de un secreto: dos requests concurrentes con el token correcto son un
  escenario del propio operador. Se documenta en `SECURITY.md`.

### 2. Cambios de DTO (backend + shared)

```go
// internal/modules/auth/request.go
type SetupRequest struct {
    Email      string `json:"email"      binding:"required,email"`
    Password   string `json:"password"   binding:"required,min=8,max=72"`
    SetupToken string `json:"setupToken" binding:"omitempty,max=256"`
}

// internal/modules/auth/response.go
type StatusResponse struct {
    Initialized        bool `json:"initialized"`
    SetupRequiresToken bool `json:"setupRequiresToken"`
}
```

`GET /auth/status`:

- `initialized`: igual que hoy.
- `setupRequiresToken`: `cfg.SetupToken != ""`.
  - `true` ⇒ hay que mostrar el campo de token.
  - `false` ⇒ el setup está deshabilitado: los clientes muestran un aviso y no
    permiten enviar el formulario.

```ts
// packages/shared/src/types/auth.types.ts
export interface SetupRequest { email: string; password: string; setupToken?: string }
export interface StatusResponse { initialized: boolean; setupRequiresToken: boolean }
```

`make swagger` actualiza `docs/`.

### 3. Rate limiting

```go
// internal/middleware/ratelimit.go
type Limit struct {
    PerMinute int // 0 = deshabilitado
    Burst     int
}

func RateLimit(l Limit, key func(*gin.Context) string) gin.HandlerFunc
```

- Implementación: `map[string]*rate.Limiter` + `sync.Mutex` + `lastSeen`.
  Barrido perezoso cuando el mapa supera 10 000 claves (elimina las no usadas en
  > 10 min). Sin goroutine de fondo.
- Clave por defecto: `c.ClientIP()` (fiables gracias a `TRUSTED_PROXIES`).
- Respuesta al exceder:

```go
c.Header("Retry-After", "60")
response.Problem(c, apierr.ErrTooManyRequests(
    "Demasiados intentos. Espera un minuto e inténtalo de nuevo.", c.Request.URL.Path))
```

- Wiring (solo en los endpoints que crean estado):

```go
auth.RegisterRoutes(api, authHandler,
    middleware.RateLimit(middleware.Limit{
        PerMinute: cfg.RateLimitAuthPerMinute,
        Burst:     cfg.RateLimitAuthBurst,
    }, nil), // nil = ClientIP
)
```

y en `auth/routes.go`:

```go
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, loginLimit gin.HandlerFunc) {
    auth := rg.Group("/auth")
    auth.GET("/status", handler.Status)                    // sin límite: lo llama el guard del SPA
    auth.POST("/setup", loginLimit, handler.Setup)
    auth.POST("/login", loginLimit, handler.Login)
}
```

- `AUTH_RATE_LIMIT_PER_MINUTE=0` deshabilita el limiter (útil en tests/E2E); el
  default es 10/min con burst 5.
- La cuota diaria de import ya existe y no se duplica aquí.

### 4. Anti-timing en login

```go
// auth/service.go
var dummyHash = func() []byte {
    h, err := bcrypt.GenerateFromPassword([]byte("timing-defense-not-a-real-password"), bcrypt.DefaultCost)
    if err != nil { panic(err) }
    return h
}()

func (s *authService) Login(email, password string) (*AuthResponse, error) {
    user, err := s.store.FindByEmail(email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            _ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password)) // mismo costo
            return nil, apierr.ErrUnauthorized("Email o contraseña inválidos", "")
        }
        return nil, apierr.ErrInternal("Error al buscar el usuario", "")
    }
    // …igual que hoy
}
```

El `dummyHash` se calcula una vez al arrancar (~60 ms) en una `sync.Once` si se
quiere evitar coste en tests; cualquiera de las dos formas es válida.

### 5. Clientes (web + mobile)

- **Web** (`SetupPage.vue`): el guard ya consulta `getStatus()`; si
  `setupRequiresToken` es `true`, muestra un campo **"Token de instalación"**
  (`type="password"`, `autocomplete="off"`) con validación "obligatorio"; si es
  `false`, muestra un `v-alert` de "Configuración inicial deshabilitada" y
  deshabilita el submit.
- **Mobile** (`SetupScreen.tsx`): mismo comportamiento con `AuthInput` y un aviso
  en `AuthCard`.
- Ambos envían `setupToken` dentro del body (sin cambios en `apiClient`).
- **E2E** (`apps/web/e2e/fixtures.ts`): añadir
  `setupToken: process.env.E2E_SETUP_TOKEN ?? 'e2e-setup-token'` y rellenar el
  campo en `critical-flow.spec.ts` cuando esté visible. Documentar en el fixture:
  el backend debe arrancar con `SETUP_TOKEN=e2e-setup-token` (o el mismo valor que
  `E2E_SETUP_TOKEN`).

## Testing strategy

Backend:

| # | Caso | Esperado |
|---|---|---|
| T1 | setup sin `SetupToken` en config | 403, `users.Count` no cambia |
| T2 | setup con token incorrecto | 403, mismo body que T1 |
| T3 | setup con token correcto + DB vacía | 201 + JWT admin |
| T4 | setup con token correcto + DB inicializada | 409 (sin cambios) |
| T5 | `GET /auth/status` con token configurado | `setupRequiresToken: true` |
| T6 | `GET /auth/status` sin token | `setupRequiresToken: false` |
| T7 | 5 logins fallidos + 6.º | 401 ×5, 429 con `Retry-After: 60` |
| T8 | rate limit con `PerMinute=0` | nunca 429 |
| T9 | `/auth/status` no consume cuota | 20 llamadas seguidas → 200 |
| T10 | login email inexistente vs password mala | mismo status/body |

Frontend/mobile:

| # | Caso | Esperado |
|---|---|---|
| T11 | `SetupPage` con `setupRequiresToken: true` | muestra el campo y bloquea submit vacío |
| T12 | `SetupPage` con `setupRequiresToken: false` | alerta + submit deshabilitado |
| T13 | `SetupScreen` (mobile) equivalentes | idem |
| T14 | E2E del flujo crítico | verde con `SETUP_TOKEN` alineado |

## Boundaries

- **Always**: `make check` + `make test` + `make swagger`; comparar el token en
  tiempo constante; mantener `/auth/status` sin rate limit; actualizar shared y
  los dos clientes en el mismo PR (el contrato es uno).
- **Ask first**: añadir `golang.org/x/time` (si prefieres cero dependencias, se
  implementa un limiter de ~40 líneas con `time.Now()`); añadir lockout por
  cuenta en DB; mover el limiter a Redis.
- **Never**: aceptar el token por query string, loguearlo, permitir setup con
  `SETUP_TOKEN` vacío "solo en producción", usar `==` para comparar el token,
  romper el guard del SPA limitando `/auth/status`.

## Success criteria

1. Sin `SETUP_TOKEN`, `POST /auth/setup` es 403 en debug y en release y no crea
   usuarios.
2. Con token correcto, el flujo de setup web, mobile y E2E sigue funcionando.
3. El 6.º login por minuto desde una IP es 429 con `Retry-After`.
4. `/auth/status` responde 200 de forma sostenida (no lo afecta el limiter).
5. Email inexistente y contraseña incorrecta son indistinguibles en
   status/body/mensaje.
6. `pnpm typecheck` (shared) + `make fe-check` + `pnpm --filter mobile typecheck`
   verdes.
7. `docs/swagger.yaml` documenta `setupToken`, `setupRequiresToken`, 403 y 429.

## Open questions

- ¿`golang.org/x/time` o limiter propio sin dependencias? (propuesto: `x/time`,
  es el oficial y pesa ~50 KB).
- ¿Se quiere lockout por cuenta (persistido) además del límite por IP? (fuera de
  alcance de esta spec).
- ¿El despliegue tendrá más de una instancia? Si sí, el limiter en memoria pasa a
  ser un problema y hay que moverlo a Redis (nueva spec).
- ¿Se rota `SETUP_TOKEN` tras inicializar? Recomendado: sí, y el endpoint seguirá
  devolviendo 409 porque ya hay admin.
