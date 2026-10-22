# Response & error rules

## Success responses

All success responses use the `{ "data": ... }` envelope. **Do not** return raw JSON at the top level.

```go
// single resource
response.Success(c, http.StatusCreated, user.ToResponse())

// collection with metadata
response.Paginated(c, http.StatusOK, items, response.Meta{Total: 42, Page: 1, PerPage: 20})
```

Import path: `github.com/felipe/dev-test-api/pkg/response`

## Error responses (RFC 9457)

All errors return a Problem Detail object:

```json
{
  "type":     "about:blank",
  "title":    "Validation Error",
  "status":   422,
  "detail":   "email: required; password: min length 8",
  "instance": "/api/v1/auth/login"
}
```

### How to use

**In services**: return `*apierr.APIError` (not raw errors). Use the predefined constructors:

```go
// from pkg/apierr
apierr.ErrNotFound("user", "")           // 404
apierr.ErrValidation("detail...", "")    // 422
apierr.ErrConflict("title", "detail", "") // 409
apierr.ErrForbidden("detail", "")         // 403 (solo política, ver abajo)
apierr.ErrUnauthorized("detail", "")      // 401
apierr.ErrPayloadTooLarge("detail...", "") // 413
apierr.ErrInternal("detail", "")          // 500
```

Leave `instance` empty in services — the handler fills it.

## Object-level authorization (404, no 403)

Cuando una operación por ID no encuentra el recurso **o** el recurso pertenece a
otro usuario, el service devuelve `ErrNotFound` (**404**), nunca `ErrForbidden`.
Así no se revela la existencia de recursos ajenos (decisión D-37 en
`specs/DECISIONS.md`).

- La comprobación vive en el **service**, no en el handler.
- El handler solo pasa `userID` / `isAdmin` extraídos del JWT.
- **Excepción**: las reglas de *política* (p. ej. "solo sesiones de las últimas
  24 h") sí usan `ErrForbidden` (403), porque el recurso sí es del usuario.
- Los listados filtran por `user_id` (o visibilidad) en el store; nunca se
  devuelve una página con recursos ajenos.

Si añades un endpoint por ID, aplica esta regla y añade el test de "otro
usuario → 404".

**In handlers**: convert and send:

```go
user, err := h.service.GetByID(id)
if err != nil {
    e := err.(*apierr.APIError)
    e.Instance = c.Request.URL.Path   // ← fill instance here
    response.Problem(c, e)
    return
}
```

**Validation errors**: use the shorthand:

```go
if err := c.ShouldBindJSON(&req); err != nil {
    response.ValidationError(c, err.Error(), c.Request.URL.Path)
    return
}
```

## DTO conventions

- `PasswordHash` field uses `json:"-"` — never serialized
- `DeletedAt` uses `json:"-"` — soft delete is invisible to clients
- Model → response conversion via `ToResponse()` method, never expose GORM model directly
- Request DTOs use `binding:"required,email"` / `binding:"required,min=8,max=72"` tags

## Auth context

JWT claims are stored in the Gin context by `middleware.Auth()`:

```go
claims, exists := c.Get("user_claims")  // type *jwt.MapClaims
```

Helper functions in `internal/auth/service.go`:
- `auth.GetUserID(&claims)` → `(uuid.UUID, bool)`
- `auth.IsAdmin(&claims)` → `bool`

## Error language

Error responses follow a bilingual convention: `title` is always in **English**, `detail` is always in **Spanish** with the first letter capitalized.

**Rules:**

| Field   | Language | Example |
|---------|----------|---------|
| `title` | English  | `"Not Found"`, `"Conflict"`, `"Unauthorized"` |
| `detail` | Spanish (first letter caps) | `"Usuario no encontrado"`, `"Email o contraseña inválidos"` |

**Template definitions** (`pkg/apierr/errors.go`):
- `ErrNotFound(entity, instance)` → title `"Not Found"`, detail `entity + " no encontrado"`
- `ErrValidation(detail, instance)` → title `"Validation Error"`, detail passed verbatim
- `ErrConflict(title, detail, instance)` → both `title` and `detail` passed verbatim (title in English, detail in Spanish)
- `ErrForbidden(detail, instance)` → title `"Forbidden"`, detail verbatim
- `ErrUnauthorized(detail, instance)` → title `"Unauthorized"`, detail verbatim
- `ErrPayloadTooLarge(detail, instance)` → title `"Payload Too Large"`, detail verbatim (413)
- `ErrInternal(detail, instance)` → title `"Internal Server Error"`, detail verbatim

**When adding new errors:**
- Pass the `title` in English (even for `ErrConflict` which accepts a custom title)
- Pass the `detail` in Spanish, starting with a capital letter
- Entity names follow the same rule: `apierr.ErrNotFound("Usuario", "")` → `"Usuario no encontrado"`
- Keep Go-internal errors (like `fmt.Errorf` in middleware) in English — they never reach the client
