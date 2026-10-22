# Plan técnico — Iniciativa `security-hardening`

> Referencias: `specs/MAP-security.md` (módulos y orden), `specs/SPEC-<id>.md` (contrato de cada módulo).
> Tareas concretas: `tasks/security-todo.md`.
> Decisiones ya tomadas: S-1 (backend en Internet), S-2 (P0+P1+P2), S-3 (`SETUP_TOKEN`), S-4 (no reescribir historia).

## 1. Grafo de dependencias

```
secrets-hygiene ──┐
                  ├──► deps-supply-chain
                  │
authz-objects ────┼── (independiente)

http-hardening ───┼──► auth-hardening
                  └──► client-hardening
```

- `auth-hardening` **depende de** `http-hardening`: necesita `SetTrustedProxies`
  correcto (si no, el rate limit por IP es evadible con `X-Forwarded-For`) y el
  body limit global para que el login no sea un vector de memoria.
- `deps-supply-chain` **depende de** `secrets-hygiene`: ambos editan
  `.github/workflows/`; el job `secrets` y `permissions: contents: read` se
  definen una sola vez.
- `client-hardening` **depende de** `http-hardening`: documenta y verifica las
  cabeceras que la API ya emite.
- `authz-objects` no comparte archivos con nadie: puede ir en paralelo desde el
  minuto uno.

## 2. Orden de implementación (olas)

| Ola | Módulos | Por qué |
|---|---|---|
| 0 | `secrets-hygiene` · `authz-objects` · `http-hardening` | Sin dependencias entre sí y cierran los riesgos de mayor impacto (IDOR) y de publicación (secretos) |
| 1 | `auth-hardening` · `deps-supply-chain` · `client-hardening` | Consumen contratos de la ola 0 |

**Paralelizable dentro de cada ola**: en la ola 0 los tres módulos tocan
conjuntos de archivos disjuntos (`specs/`, `.github/`, `sessions|questions`,
`middleware|server`). En la ola 1, `deps-supply-chain` y `client-hardening` no
comparten archivos con `auth-hardening` salvo `package.json`/lock
(`client-hardening` no toca dependencias).

## 3. Componentes y touch points

| Componente | Backend | Web | Mobile | CI/Infra |
|---|---|---|---|---|
| `secrets-hygiene` | — | — | — | `.gitignore`, `.gitleaks.toml`, `SECURITY.md`, `workflows/secrets.yml`, settings de GitHub |
| `authz-objects` | `sessions/{service,handler}.go`, `questions/{store,service,handler}.go` | — | — | `docs/` (swagger) |
| `http-hardening` | `config.go`, `middleware/{security,bodylimit}.go`, `server.go`, `pkg/apierr`, `questions/handler.go` | — | — | `.env.example` |
| `auth-hardening` | `auth/{request,response,handler,routes,service}.go`, `middleware/ratelimit.go`, `server.go` | `SetupPage.vue` | `auth.store.ts`, `SetupScreen.tsx` | `docs/`, `e2e/` |
| `deps-supply-chain` | `go.mod`, `go.sum` | `package.json` (vite) | `package.json`, `MarkdownText.tsx` | `dependabot.yml`, `Makefile`, workflows |
| `client-hardening` | — | `index.html`, `deploy/`, READMEs | `config.ts` | `.env.example` |

**Reglas transversales**

- Backend: `.agents/backend/architecture.md` (request → response → mapper → store →
  service → handler → routes).
- Contrato HTTP: `.agents/backend/responses.md` (envelope, RFC 9457, camelCase,
  `title` EN / `detail` ES).
- Web: `.agents/web/architecture.md` (services → queries → features).
- Mobile: `.agents/mobile/architecture.md` (rutas thin + `src/features` + `src/core`).
- Cualquier cambio de contrato API ⇒ `make swagger` y actualizar
  `packages/shared/src/types/` en el mismo PR.

## 4. Riesgos y mitigaciones

| Riesgo | Prob. | Impacto | Mitigación |
|---|---|---|---|
| `SETUP_TOKEN` rompe el setup local y el E2E existente | **Alta** | Medio | `AH-8` alinea `E2E_SETUP_TOKEN` con `SETUP_TOKEN`; el fixture documenta el acoplamiento; `SetupPage`/`SetupScreen` muestran un aviso claro cuando falta |
| Cambiar `Delete` de 403 a 404 rompe un test/cliente existente | Media | Bajo | `AZ-3` incluye el test del caso 24 h (sigue 403) y `make swagger`; el E2E no borra sesiones ajenas |
| El rate limiter por IP bloquea a varios usuarios tras un NAT | Media | Medio | Burst 5 y 10/min; `AUTH_RATE_LIMIT_PER_MINUTE=0` para desactivar; documentado en `SECURITY.md` |
| Un `WriteTimeout` corta un import grande | Baja | Medio | Límite de 200 preguntas/día ya existente; `HH-5` verifica con el máximo real |
| Forzar `markdown-it` 14 rompe el render en mobile | Media | Alto | `DS-3` es un spike con fallback (truncado a 8 000 chars) y test de humo; nunca se fuerza el override sin test |
| Un `pnpm up` global rompe Expo | Media | Alto | Prohibido por spec; solo `pnpm.overrides` de transitivas + `expo install --check` |
| `govulncheck@latest` en CI cambia de comportamiento sin avisar | Baja | Bajo | Se fija la versión del módulo `golang.org/x/vuln` en el comando (`@vX.Y.Z`) al implementar `DS-5` |
| Se silencia una CVE sin justificar | Media | Alto | `ignoreCves` exige entrada con motivo + fecha en `SECURITY.md`; revisión trimestral |
| El hook de gitleaks local no existe en todas las máquinas | Media | Bajo | Solo CI es obligatorio; el hook local se documenta como opcional (`SH-2`) |
| `sec-fetch`/CSP rompe Vuetify en producción | Media | Medio | `CH-2` deja `style-src 'unsafe-inline'` y se verifica en staging (`CH-4`) |

## 5. Checkpoints de verificación

| Checkpoint | Cuándo | Criterio |
|---|---|---|
| V1 | fin de `secrets-hygiene` | `git check-ignore` correcto, `gitleaks detect` limpio, canario detectado, CI verde, secret scanning activo |
| V2 | fin de `authz-objects` | 5 endpoints de sesión → 404 con JWT ajeno y sin escrituras; `questions.Get` respeta visibilidad; `make test` + `make swagger` |
| V3 | fin de `http-hardening` | `curl -I` con 5 cabeceras + HSTS; Swagger 404 en release; POST de 20 MB → 413 sin OOM; `X-Forwarded-For` ignorado sin `TRUSTED_PROXIES` |
| V4 | fin de `auth-hardening` | Setup sin token → 403; con token → 201; 6.º login → 429; `/auth/status` sostenido; web + mobile + E2E verdes |
| V5 | fin de `deps-supply-chain` | `govulncheck` 0 alcanzables; `pnpm audit --prod --audit-level high` 0; Dependabot activo; actions pinneadas |
| V6 | fin de `client-hardening` | Mobile falla al arrancar con `http://` en producción; cabeceras en staging; SPA sin errores de CSP |

## 6. Reglas de ejecución

1. Una tarea por sesión; marcar `[~]` al empezar y `[x]` al terminar (más
   `Hecho:` con lo verificado, como en `tasks/todo.md`).
2. Antes de cerrar una tarea: ejecutar todos los comandos de su `Verify:`.
3. Si el contrato cambia respecto a la spec, **actualizar la spec primero** y
   anotarlo en `specs/DECISIONS.md` si es una decisión nueva.
4. Commits **manuales** (el agente nunca commitea), uno por tarea o por grupo
   acoplado del mismo módulo.
5. No mezclar módulos en un mismo commit salvo que compartan archivo
   (`server.go`, `package.json`).
6. **Doc sync**: si el cambio consolida una convención, contrato, comando o
   variable, actualizar su fuente de verdad (`.agents/`, `specs/DECISIONS.md`,
   `README.md`) y la referencia en `AGENTS.md`. Una tarea no está terminada si la
   documentación quedó con el estado anterior.
7. Antes de proponer trabajo nuevo, comprobar `specs/`, `tasks/` y `git log`: si
   ya está hecho o decidido, no se rehace.

## 7. Mapeo spec → tareas

| Spec | Tareas |
|---|---|
| `SPEC-secrets-hygiene.md` | `SH-1` … `SH-4` |
| `SPEC-authz-objects.md` | `AZ-1` … `AZ-5` |
| `SPEC-http-hardening.md` | `HH-1` … `HH-5` |
| `SPEC-auth-hardening.md` | `AH-1` … `AH-8` |
| `SPEC-deps-supply-chain.md` | `DS-1` … `DS-6` |
| `SPEC-client-hardening.md` | `CH-1` … `CH-4` |

**Total: 32 tareas.** Ninguna requiere migración goose. Una sola dependencia
nueva (`golang.org/x/time`, sujeta a aprobación en `AH-2`).
