# Capability Map — `security-hardening`

| Campo | Valor |
|---|---|
| Iniciativa | `security-hardening` |
| Origen | Revisión de seguridad 2026-10-21 |
| Estado | `ready` (aprobado el 2026-10-21) |
| Alcance | P0 + P1 + P2 (decisión S-2) |
| Entorno objetivo | Backend expuesto a Internet (decisión S-1) |
| Migraciones goose | **Ninguna** |
| Dependencia nueva | `golang.org/x/time` (rate limit) — *Ask First* |

## Módulos

| Module id | Responsabilidad | Depende de | Hallazgos que cierra |
|---|---|---|---|
| `secrets-hygiene` | Que ningún secreto pueda entrar al repo ni pasar desapercibido: `.gitignore`, gitleaks local/CI, `SECURITY.md`, secret scanning + push protection | — | M5, §1 |
| `authz-objects` | Que un usuario solo pueda leer/escribir sus sesiones y las preguntas visibles para él | — | H1, H2 |
| `http-hardening` | Que el servidor resista abuso trivial y no exponga superficie innecesaria: body limits, timeouts, cabeceras, Swagger, proxies | — | M1, M2, M4, L2 |
| `auth-hardening` | Que nadie pueda tomar el rol admin en una instancia nueva ni hacer fuerza bruta al login | `http-hardening` | H3, H4, L3 |
| `deps-supply-chain` | Cerrar vulnerabilidades conocidas de dependencias y automatizar la vigilancia | `secrets-hygiene` | M3, L6 |
| `client-hardening` | Que los clientes no viajen en claro y que el despliegue web tenga política de seguridad | `http-hardening` | M4 (web), L1, L7 |

## Build order

```
Ola 0   secrets-hygiene ∥ authz-objects ∥ http-hardening
Ola 1   auth-hardening ∥ deps-supply-chain ∥ client-hardening
```

Sin ciclos. En la ola 0 ningún módulo comparte archivos: `http-hardening` es el
único que toca `internal/server/server.go` y `auth-hardening` lo reutiliza en la
ola 1 (necesita `SetTrustedProxies` y el body limit ya resueltos).

`deps-supply-chain` edita `.github/workflows/*` y por eso va después de
`secrets-hygiene` (que crea `secrets.yml` y fija `permissions`), para no
provocar conflictos de merge en los mismos archivos.

## Interfaces en las fronteras

| Proveedor | Contrato expuesto | Consumidor |
|---|---|---|
| `http-hardening` | `middleware.SecurityHeaders()`, `middleware.BodyLimit(maxBytes)`; `server.go` con `SetTrustedProxies` y timeouts | `auth-hardening`, `client-hardening` |
| `http-hardening` | `apierr.ErrPayloadTooLarge(detail, instance)` (413) | `auth-hardening` (no lo usa, pero comparte el paquete `apierr`) |
| `authz-objects` | `sessions.Service.{GetByID,Summary,Finish,NextQuestion}(..., userID)` y `questions.Service.GetByID(id, userID, isAdmin)` | — (interno) |
| `auth-hardening` | `StatusResponse.setupRequiresToken` + `SetupRequest.setupToken` | `client-hardening` (pantallas de setup web/mobile) |
| `secrets-hygiene` | `.gitleaks.toml` + job `secrets` en CI | `deps-supply-chain` |

**Regla de consistencia**: la autorización por objeto se resuelve **siempre** en
la capa `service` (no en el handler), y el mensaje de error es
`ErrNotFound` cuando el objeto no existe *o* no pertenece al usuario. Ningún
módulo reimplementa la comprobación.

## Trazabilidad hallazgo → módulo → tarea (IDs tentativos)

| Hallazgo | Sev. | Módulo | Tareas |
|---|---|---|---|
| H1 IDOR sesiones | Alta | `authz-objects` | `AZ-1`, `AZ-2`, `AZ-3` |
| H2 Fuga de answer key en preguntas | Alta | `authz-objects` | `AZ-4`, `AZ-5` |
| H3 `/auth/setup` público | Alta | `auth-hardening` | `AH-1`, `AH-2`, `AH-3` |
| H4 Sin rate limit en login | Alta | `auth-hardening` | `AH-4`, `AH-5` |
| M1 Body/timeouts | Media | `http-hardening` | `HH-1`, `HH-2`, `HH-3` |
| M2 Swagger público | Media | `http-hardening` | `HH-4` |
| M4 Headers + CSP | Media | `http-hardening`, `client-hardening` | `HH-5`, `CH-2`, `CH-3` |
| M3 Dependencias vulnerables | Media | `deps-supply-chain` | `DS-1`, `DS-2`, `DS-3` |
| M5 `.gitignore` incompleto | Media | `secrets-hygiene` | `SH-1` |
| L1 JWT sin revocación / TTL | Baja | `client-hardening` | `CH-4` |
| L2 Proxies confiables | Baja | `http-hardening` | `HH-6` |
| L3 Timing login | Baja | `auth-hardening` | `AH-6` |
| L6 CI sin `permissions` / pinning | Baja | `deps-supply-chain` | `DS-4` |
| L7 Mobile HTTP | Baja | `client-hardening` | `CH-1` |
| §1 Higiene de secretos | Media | `secrets-hygiene` | `SH-2`, `SH-3`, `SH-4` |

## Riesgos del mapa

| Riesgo | Mitigación |
|---|---|
| `auth-hardening` cambia el contrato de setup y obliga a tocar web **y** mobile (2 clientes + E2E) | La spec define el contrato una vez; los clientes solo añaden un campo y un aviso. Si se quiere reducir alcance, se puede cortar el módulo a backend y dejar la UI en una spec aparte |
| `deps-supply-chain` puede romper Expo al forzar overrides de `markdown-it` | Spike explícito con fallback: limitar el tamaño de contenido renderizado y aceptar el riesgo documentado |
| `http-hardening` con `WriteTimeout` puede cortar imports grandes | Límite de import ya existe (≤ 200 preguntas/día); verificar con el import máximo |
| El rate limiter en memoria no sirve con >1 instancia | Documentado en Open Questions; el despliegue actual es single-instance |

## Fuera de alcance

- Reescritura de historia de Git (decisión S-4).
- Migrar de `localStorage` a cookie `httpOnly` (requiere CSRF + cambios de contrato).
- WAF, mTLS, secret manager (Vault/SSM).
- Lockout por cuenta persistido en DB (se resuelve con rate limit por IP).
- Pentest automatizado / DAST.
