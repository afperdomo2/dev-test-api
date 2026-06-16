# API client behavior

## Axios interceptors (`api/client.ts`)

The client automatically handles two concerns:

### Request interceptor

Attaches `Authorization: Bearer <token>` from `localStorage` to every request. Tokens are managed via `utils/storage.ts` (`getToken`, `setToken`, `removeToken`).

### Response interceptor — envelope unwrap (critical gotcha)

The backend wraps all success responses in `{ "data": ... }`. The interceptor unwraps this envelope automatically, **except** for paginated responses:

```ts
// interceptor logic — simplified
if ('data' in response.data && !('meta' in response.data)) {
  response.data = response.data.data  // unwrap single/collection responses
}
// Paginated responses ({ data: [...], meta: {...} }) pass through unchanged
```

| Backend returns | Interceptor → service receives |
|---|---|
| `{ data: { token: "...", user: {...} } }` | `{ token: "...", user: {...} }` |
| `{ data: { id: "...", email: "..." } }` | `{ id: "...", email: "..." }` |
| `{ data: [...], meta: { total: 1, ... } }` | `{ data: [...], meta: { total: 1, ... } }` |

**Never** add a `meta` field to non-paginated backend responses — the unwrap logic would skip it and the frontend would receive the raw envelope instead of the unwrapped data.

### Error interceptor

- Transforms RFC 9457 errors into typed `ApiError` objects (`{ type, title, status, detail, instance }`)
- On 401: auto-removes the token from localStorage (forcing re-login)
- The rejected promise carries the `ApiError` object, not the raw Axios error

## Service functions

Service files export plain async functions. They receive typed inputs and return typed promises. Example:

```ts
// api/services/users.service.ts
export async function listUsers(page: number, perPage: number): Promise<PaginatedResponse<User>> {
  const res = await apiClient.get<PaginatedResponse<User>>('/api/v1/users', {
    params: { page, per_page: perPage },
  })
  return res.data
}
```

Always use `res.data` (the interceptor already unwrapped the envelope). The type annotation on `apiClient.get<T>()` is the type of the unwrapped data.

## Form error mapping

Backend validation errors use the format: `"field: message; field2: message2"` in the `detail` field.

`useFormErrors()` composable parses this into `Record<string, string>` mapping field names to messages. Use it in form components to display server-side validation errors on the correct fields:

```ts
const { extractFieldErrors } = useFormErrors()
// ...
catch (err) {
  serverErrors.value = extractFieldErrors(err)
}
```
