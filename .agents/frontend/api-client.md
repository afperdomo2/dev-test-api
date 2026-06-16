# API client behavior

## Axios interceptors (`api/client.ts`)

The client automatically handles two concerns:

### Request interceptor

Attaches `Authorization: Bearer <token>` from `localStorage` to every request. Tokens are managed via `utils/storage.ts` (`getToken`, `setToken`, `removeToken`).

### Response interceptor — processing order

The interceptor applies two transformations in sequence on every response:

```
1. Envelope unwrap  →  { data: ... } extracted if not paginated
2. snake→camel      →  all keys deep-transformed via snakeToCamel()
```

**Critical gotcha — paginated responses**: the unwrap step checks for `data` without `meta`. Paginated responses pass through step 1 unchanged, then get camelized in step 2 (including `meta.per_page` → `meta.perPage`).

| Backend returns | After unwrap | After camelize |
|---|---|---|
| `{ data: { token: "...", user: { is_admin: true } } }` | `{ token: "...", user: { is_admin: true } }` | `{ token: "...", user: { isAdmin: true } }` |
| `{ data: [...], meta: { total: 1, per_page: 20 } }` | `{ data: [...], meta: { total: 1, per_page: 20 } }` | `{ data: [...], meta: { total: 1, perPage: 20 } }` |

**Never** add a `meta` field to non-paginated backend responses.

### Error interceptor

- Transforms RFC 9457 errors into typed `ApiError` objects (`{ type, title, status, detail, instance }`)
- On 401: auto-removes the token from localStorage (forcing re-login)
- The rejected promise carries the `ApiError` object, not the raw Axios error

## Naming convention: snake_case vs camelCase

The backend (Go) uses `snake_case` in JSON (`is_admin`, `created_at`). The frontend (TypeScript/Vue) MUST use `camelCase` (`isAdmin`, `createdAt`).

### Rule

| Direction | Convention | Where |
|-----------|-----------|-------|
| **Response data** (API → frontend) | `camelCase` | All `types/*.types.ts` interfaces, all template bindings, all headers keys, all slot names |
| **Request bodies** (frontend → API) | `snake_case` | POST/PUT data, query params, form submissions |

### How it works

The `snakeToCamel()` utility in `utils/transform.ts` runs inside the response interceptor (`api/client.ts`). It recursively converts all object keys from `snake_case` to `camelCase`. This happens automatically — services and components receive already-camelized data.

Request bodies are NOT transformed. They are sent as-is to the backend, which expects `snake_case` JSON tags.

### When adding a new type

1. Define the **response type** with `camelCase` properties matching the camelized output of `snakeToCamel()`
2. Define the **request type** with `snake_case` properties matching the backend's JSON tags
3. Data table header `key` values use `camelCase` (they match response object keys)
4. Template slot names for data table columns use `camelCase`: `#[`item.isAdmin`]`, `#[`item.createdAt`]`

Example:
```ts
// Response type — camelCase
export interface User {
  id: string
  isAdmin: boolean      // backend: is_admin
  createdAt: string     // backend: created_at
}

// Request type — snake_case (sent to backend as-is)
export interface CreateUserRequest {
  email: string
  password: string
  is_admin?: boolean    // backend expects this key name
}

// Data table header — key matches response type
{ title: 'Rol', key: 'isAdmin' }

// Data table slot — name matches header key
<template #[`item.isAdmin`]="{ item }">
```

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
