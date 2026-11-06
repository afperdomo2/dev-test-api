# SPEC — profile

| Campo | Valor |
|---|---|
| Module id | `profile` |
| Origen | Petición de producto 2026-11-05 |
| Estado | `ready` |
| Depende de | `auth`, `users` |
| Consumido por | web y mobile |
| Migración goose | Ninguna |
| Dependencia nueva | Ninguna |
| Cambio de contrato | Sí: nuevo `PUT /api/v1/profile/password` |

## Objective

El usuario autenticado puede **ver sus datos** (email, rol, fecha de alta, ID y cupo
de importación) y **cambiar su propia contraseña** sin depender de un
administrador. Hoy solo existe `PUT /api/v1/users/:id`, que es admin-only.

**Historias**

- Como usuario, quiero ver mis datos en una página propia.
- Como usuario, quiero cambiar mi contraseña sin pedírselo a un admin.
- Como usuario, quiero seguir dentro de la app después de cambiarla.

**Éxito**

- `GET /api/v1/profile` sigue devolviendo el `UserResponse` actual (sin cambios).
- `PUT /api/v1/profile/password` responde 200 con `{"data":{"message":"..."}}`
  cuando `currentPassword` coincide y `newPassword` es válida y distinta.
- La contraseña actual incorrecta responde **422**, no 401, para no disparar el
  logout automático del cliente web.
- La sesión JWT sigue siendo válida después del cambio.
- Web muestra `/profile` con el resumen del usuario y el formulario. Mobile añade
  la sección "Seguridad" a `ProfileScreen`.

## Contrato

```http
PUT /api/v1/profile/password
Authorization: Bearer <jwt>
Content-Type: application/json

{ "currentPassword": "actual123", "newPassword": "nueva12345" }
```

| Caso | Status | `detail` |
|---|---|---|
| Éxito | 200 | — (`data.message`) |
| Body inválido (requerido / 8–72) | 422 | mensaje de Gin |
| Contraseña actual incorrecta | 422 | `La contraseña actual es incorrecta` |
| Nueva igual a la actual | 422 | `La nueva contraseña debe ser diferente a la actual` |
| Usuario inexistente (soft-deleted) | 404 | `Usuario no encontrado` |
| Sin JWT | 401 | `Usuario no autenticado` |

## Alcance

**Backend**

- `internal/modules/users/request.go`: `ChangePasswordRequest`.
- `internal/modules/users/service.go`: `ChangePassword` (bcrypt compare + hash +
  `store.Update`).
- `internal/modules/users/handler.go`: handler `ChangePassword` con `getUserID`.
- `internal/modules/users/routes.go`: ruta protegida `/profile/password`.
- `docs/` regenerado con `make swagger`.

**Shared / web**

- `ChangePasswordRequest` en `packages/shared/src/types/auth.types.ts` y re-export
  en `apps/web/src/types/auth.types.ts`.
- `changePassword()` en `api/services/users.service.ts` y
  `changePasswordMutation()` en `queries/users.queries.ts`.
- `components/UserSummaryCard.vue`: card de usuario extraído del Dashboard y
  reutilizado en `/profile`.
- `features/profile/pages/ProfilePage.vue` y
  `features/profile/components/ChangePasswordCard.vue`.
- Ruta `/profile` (`requiresAuth`) y entradas "Mi perfil" en sidebar y menú de
  avatar.

**Mobile**

- `changePassword()` en `core/api/services/auth.service.ts`.
- `features/profile/components/ChangePasswordSection.tsx` y su integración en
  `ProfileScreen` dentro de la sección "Seguridad".

**Fuera de alcance**

- Revocar JWT existentes al cambiar la contraseña (ver D-49).
- Editar email o `dailyImportLimit` desde el perfil.
- Recuperación de contraseña por email.
- Rate limiting específico del endpoint (requiere JWT válido; el limiter actual
  solo cubre `/auth/setup` y `/auth/login`, D-39).

## Verify

- `make fmt`, `make vet`, `go test ./internal/modules/users/...`, `make swagger`.
- `pnpm --filter @devtest/shared typecheck`.
- `make fe-check`, Prettier/ESLint del set tocado y Vitest de los archivos
  nuevos/modificados.
- `pnpm --filter mobile typecheck` y Jest de los tests nuevos.
- `make staged-check` antes del commit manual.
