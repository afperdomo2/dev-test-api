# API client behavior

## API contract

The envelope `{ data: ... }`, RFC 9457 error format, and camelCase convention are defined by the backend. **See `.agents/backend/responses.md` for the canonical contract.** This file documents only the web-specific consuming side (Axios interceptors, localStorage, service patterns).

## Axios interceptors (`api/client.ts`)

The client automatically handles two concerns:

### Request interceptor

Attaches `Authorization: Bearer <token>` from `localStorage` to every request. Tokens are managed via `utils/storage.ts` (`getToken`, `setToken`, `removeToken`).

### Response interceptor — envelope unwrap

The interceptor unwraps the API envelope on every response:

```
{ data: ... } → extracted value (if no meta key)
```

**Critical gotcha — paginated responses**: the unwrap only triggers when the response has `data` but NOT `meta`. Paginated responses `{ data: [...], meta: { total, page, perPage } }` pass through unchanged so components can access `data.value?.meta?.total`.

| Backend returns | After interceptor |
|---|---|
| `{ data: { token: "...", user: { isAdmin: true } } }` | `{ token: "...", user: { isAdmin: true } }` |
| `{ data: [...], meta: { total: 1, perPage: 20 } }` | `{ data: [...], meta: { total: 1, perPage: 20 } }` |

**Never** add a `meta` field to non-paginated backend responses.

### Error interceptor

- Transforms RFC 9457 errors into typed `ApiError` objects (`{ type, title, status, detail, instance }`)
- On 401: auto-removes the token from localStorage (forcing re-login)
- The rejected promise carries the `ApiError` object, not the raw Axios error

## Naming convention

Both the backend (Go) and web/mobile (TypeScript) use **camelCase** for all JSON keys. No runtime conversion exists. See `.agents/backend/responses.md` for the Go struct tag rules. When adding a new type:

1. Read `docs/swagger.yaml` to see the exact JSON shape
2. Define both **response** and **request** types with `camelCase` properties matching the backend's JSON tags
3. Data table header `key` values use `camelCase` (they match response object keys)
4. Template slot names for data table columns use `camelCase`: `#[`item.isAdmin`]`, `#[`item.createdAt`]`

## Service functions

Service files export plain async functions. They receive typed inputs and return typed promises:

```ts
export async function listUsers(page: number, perPage: number): Promise<PaginatedResponse<User>> {
  const res = await apiClient.get<PaginatedResponse<User>>('/api/v1/users', {
    params: { page, perPage },
  })
  return res.data
}
```

Always use `res.data` (the interceptor already unwrapped the envelope). The type annotation on `apiClient.get<T>()` is the type of the unwrapped data.

Query param keys use camelCase (`perPage`, `sortBy`, `sortOrder`, filter names) — the backend parses them from `c.Query("perPage")` / `c.Query("sortBy")`.

## Form error mapping

Backend validation errors use the format: `"field: message; field2: message2"` in the `detail` field.

`useFormErrors()` composable parses this into `Record<string, string>` mapping field names to messages:

```ts
const { extractFieldErrors } = useFormErrors()
// ...
catch (err) {
  serverErrors.value = extractFieldErrors(err)
}
```
