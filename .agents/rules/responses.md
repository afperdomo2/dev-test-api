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
apierr.ErrForbidden("detail", "")         // 403
apierr.ErrUnauthorized("detail", "")      // 401
apierr.ErrInternal("detail", "")          // 500
```

Leave `instance` empty in services — the handler fills it.

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
