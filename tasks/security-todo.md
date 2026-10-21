# Tareas — Iniciativa `security-hardening`

> Ordenadas por dependencia (ver `tasks/security-plan.md`). Specs: `specs/SPEC-<id>.md`, mapa: `specs/MAP-security.md`.
>
> **Leyenda**: `[ ]` pendiente → `[~]` en progreso → `[x]` hecha → `[-]` descartada (con motivo).
> **Commits**: manuales, uno por tarea o por grupo acoplado del mismo módulo.
> **Regla**: al terminar, añadir una línea `Hecho:` con lo verificado (como en `tasks/todo.md`).

---

## Módulo `secrets-hygiene` — `specs/SPEC-secrets-hygiene.md`

- [ ] **SH-1** — `.gitignore`: variantes de `.env` y material criptográfico
  - Acceptance: el bloque de la spec está añadido en la raíz; `.env`, `.env.local`, `.env.production`, `apps/mobile/.env.preview`, `*.pem`, `*.key`, `*.p12`, `*.jks`, `*.keystore` y `.npmrc` quedan ignorados; `.env.example` (raíz, web y mobile) **no** queda ignorado; ningún archivo ya trackeado pasa a estar ignorado
  - Verify: `git check-ignore -v .env.local .env.production apps/mobile/.env.preview cert.pem server.key` imprime reglas; `git check-ignore -v .env.example apps/web/.env.example` no imprime nada; `git status --porcelain` limpio
  - Files: `.gitignore`

- [ ] **SH-2** — `.gitleaks.toml` + canario local
  - Acceptance: extiende las reglas por defecto; allowlist solo de `*.env.example`, `testdata/`, `db/seeds/`, `docs/swagger.*`, `*.test.ts(x)` y `e2e/fixtures.ts`; no hay paths de producción (`internal/`, `apps/*/src/`) silenciados
  - Verify: `gitleaks detect --source . --redact --no-banner` → exit 0; `printf 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n' > "$TMPDIR/canary.txt"` + `gitleaks detect --no-git --source "$TMPDIR" --redact` → exit 1 y finding `aws-access-token` sin imprimir el token; borrar el canario
  - Files: `.gitleaks.toml`

- [ ] **SH-3** — Workflow `Secrets` + `SECURITY.md`
  - Acceptance: `.github/workflows/secrets.yml` corre en push a `main` y en PR; `permissions: contents: read`; `fetch-depth: 0`; versión de gitleaks pineada con `sha256sum -c`; `SECURITY.md` con las 6 secciones de la spec (reporte, plazos, alcance, rotación, allowlist, alcance del módulo)
  - Verify: YAML parsea (`act -l` o `yq '.' .github/workflows/secrets.yml`); `gitleaks detect` local limpio; revisión manual de `SECURITY.md` contra los 6 puntos
  - Files: `.github/workflows/secrets.yml`, `SECURITY.md`

- [ ] **SH-4** — GitHub: Secret Scanning + Push Protection (manual)
  - Acceptance: Secret scanning y Push protection activos (Settings → Code security); un push de canario a una rama desechable es **rechazado**; el checklist firmado con fecha queda en `SECURITY.md`
  - Verify: `git switch -c test/canary && echo 'AKIAIOSFODNN7EXAMPLE' > canary.txt && git add . && git commit -m test && git push -u origin test/canary` → rechazado por push protection; borrar rama y archivo; `git push origin --delete test/canary`
  - Files: `SECURITY.md` (checklist)

---

## Módulo `authz-objects` — `specs/SPEC-authz-objects.md`

- [ ] **AZ-1** — Ownership en `sessions.Service` (5 métodos)
  - Acceptance: helper `findOwnedSession` (D-AZ-2) usado por `GetByID`, `Summary`, `Finish`, `NextQuestion` y `Answer`; `Answer` valida dueño **antes** de la idempotencia; no-owner → `ErrNotFound` (D-AZ-1); el modo dueño/admin no cambia
  - Verify: `make check && make test` — tests T1–T9 de la spec (incluye spy de que `CompleteSession` no se llama en T4 y que no hay `Create` en T6)
  - Files: `internal/modules/sessions/service.go`, `internal/modules/sessions/service_test.go`

- [ ] **AZ-2** — Handlers de sesiones pasan `userID` + tests HTTP
  - Acceptance: los 5 handlers usan `getUserID(c)` y responden 401 sin token y 404 con JWT ajeno; `e.Instance = c.Request.URL.Path` en todos; ningún handler lee `userID` de query/path
  - Verify: `make test`; test HTTP con router + `middleware.Auth` + JWT firmado (T2, T4, T6, T8, T9); `curl` manual con dos usuarios sembrados
  - Files: `internal/modules/sessions/handler.go`, `internal/modules/sessions/handler_test.go`

- [ ] **AZ-3** — `DELETE /sessions/:id`: 403 → 404 por ownership
  - Acceptance: dueño ajeno → 404 (D-AZ-3); regla de 24 h sigue devolviendo 403; swagger actualizado; sin cambios en el resto del contrato
  - Verify: `make test` (T10, T11); `make swagger` y `git diff docs/` muestra solo el cambio documentado; E2E web verde
  - Files: `internal/modules/sessions/service.go`, `internal/modules/sessions/service_test.go`, `docs/`

- [ ] **AZ-4** — Visibilidad en `questions` (store + service + handler)
  - Acceptance: `Store.FindVisibleByID(id, userID, isAdmin)` con la regla D-AZ-4 (`is_public OR source='ai_generated' OR user_id = caller`; admin ve todo); `Service.GetByID(id, userID, isAdmin)`; handler usa `getUserRoleAndID(c)`; privada ajena → 404; el listado no cambia
  - Verify: `make check && make test` (T12–T16 en `service_test.go`); `make build`
  - Files: `internal/modules/questions/store.go`, `internal/modules/questions/service.go`, `internal/modules/questions/handler.go`, `internal/modules/questions/service_test.go`

- [ ] **AZ-5** — Tests HTTP de preguntas + Swagger
  - Acceptance: 404 para privada ajena a nivel HTTP (no solo service); 200 para pública, `ai_generated`, dueño y admin; `@Failure 404` reflejado; `docs/` regenerado sin cambios ajenos
  - Verify: `make test`; `make swagger`; `git diff --stat docs/`
  - Files: `internal/modules/questions/handler_test.go`, `docs/`

---

## Módulo `http-hardening` — `specs/SPEC-http-hardening.md`

- [ ] **HH-1** — Configuración HTTP (`MAX_BODY_BYTES`, `TRUSTED_PROXIES`, `HSTS_ENABLED`, `SWAGGER_ENABLED`)
  - Acceptance: `Config` expone los 4 campos con los defaults de la spec; `SWAGGER_ENABLED` vacío ⇒ `GinMode != "release"`; `TRUSTED_PROXIES` se parsea a `[]string` sin entradas vacías; `.env.example` documenta los 4 con comentarios
  - Verify: `make check && make test`; `SWAGGER_ENABLED= GIN_MODE=release make run` (log de arranque) y `SWAGGER_ENABLED=true GIN_MODE=release make run`
  - Files: `internal/config/config.go`, `.env.example`

- [ ] **HH-2** — Middleware `SecurityHeaders`
  - Acceptance: emite `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `Permissions-Policy` y CSP de API; `Strict-Transport-Security` solo con `hsts=true`; no toca respuestas de CORS/preflight
  - Verify: `make test` (T1, T2 de la spec)
  - Files: `internal/middleware/security.go`, `internal/middleware/security_test.go`

- [ ] **HH-3** — `BodyLimit` + `apierr.ErrPayloadTooLarge` (413)
  - Acceptance: `BodyLimit` envuelve el body en `http.MaxBytesReader`; `IsPayloadTooLarge` detecta `*http.MaxBytesError`; `ErrPayloadTooLarge` con título `Payload Too Large` y status 413; sin cambios en `ErrTooManyRequests`
  - Verify: `make test` (T3, T4, T5); `make check`
  - Files: `internal/middleware/bodylimit.go`, `internal/middleware/bodylimit_test.go`, `pkg/apierr/errors.go`, `pkg/apierr/errors_test.go`

- [ ] **HH-4** — `server.go`: timeouts, proxies, wiring y gate de Swagger
  - Acceptance: `http.Server` con `ReadHeaderTimeout 5s`, `ReadTimeout 30s`, `WriteTimeout 60s`, `IdleTimeout 120s`, `MaxHeaderBytes 1 MiB` como constantes; `SetTrustedProxies(nil)` si la lista está vacía y `log.Fatalf` si es inválida; `SecurityHeaders` + `BodyLimit` registrados antes de CORS; `/swagger/*any` solo si `SwaggerEnabled`
  - Verify: `make build`; `curl -sI localhost:8080/health` (5 cabeceras); `GIN_MODE=release SWAGGER_ENABLED=false` → `/swagger/index.html` 404; `TRUSTED_PROXIES=` + `curl -H 'X-Forwarded-For: 1.2.3.4'` y comprobar IP real en el log
  - Files: `internal/server/server.go`

- [ ] **HH-5** — Import con límite de 8 MiB → 413 + Swagger
  - Acceptance: `questions.Import` rechaza `fh.Size > 8 MiB` y bodies que excedan `MAX_BODY_BYTES` con 413; copia con `io.LimitReader` (nunca `io.Copy` ilimitado); `@Failure 413` en swagger
  - Verify: `make test` (T6, T7); `make swagger`; POST de 20 MB con `curl` → 413 y RSS estable
  - Files: `internal/modules/questions/handler.go`, `internal/modules/questions/handler_test.go`, `docs/`

---

## Módulo `auth-hardening` — `specs/SPEC-auth-hardening.md`

- [ ] **AH-1** — Configuración: `SETUP_TOKEN` + límites de auth
  - Acceptance: `Config.SetupToken`, `RateLimitAuthPerMinute` (10) y `RateLimitAuthBurst` (5); `log.Println` de aviso cuando `SETUP_TOKEN` está vacío (no fatal); `.env.example` documenta las 3 variables; `AUTH_RATE_LIMIT_PER_MINUTE=0` deshabilita
  - Verify: `make check`; arranque con y sin `SETUP_TOKEN` y revisar el log
  - Files: `internal/config/config.go`, `.env.example`

- [ ] **AH-2** — Middleware `RateLimit` (token bucket)
  - Acceptance: `RateLimit(Limit{PerMinute, Burst}, keyFn)` con `map[string]*rate.Limiter` + mutex; barrido perezoso >10 000 claves y >10 min; `PerMinute=0` ⇒ pass-through; al exceder 429 + `Retry-After: 60` + `ErrTooManyRequests`; dependencia `golang.org/x/time` cubierta por D-39 (revisar `.agents/backend/dependencies.md` antes de instalar)
  - Verify: `go get golang.org/x/time && go mod tidy`; `make test` (T7, T8); `make check`
  - Files: `internal/middleware/ratelimit.go`, `internal/middleware/ratelimit_test.go`, `go.mod`, `go.sum`

- [ ] **AH-3** — `SETUP_TOKEN` en el backend + `setupRequiresToken`
  - Acceptance: `SetupRequest.SetupToken` (`omitempty,max=256`); comparación en tiempo constante sobre SHA-256 de la spec; token vacío/incorrecto → 403 con body idéntico; no se loguea; `StatusResponse.SetupRequiresToken` refleja si está configurado; `GET /auth/status` no se ve afectado
  - Verify: `make test` (T1–T6, T10); `make check`
  - Files: `internal/modules/auth/request.go`, `internal/modules/auth/response.go`, `internal/modules/auth/handler.go`, `internal/modules/auth/handler_test.go`

- [ ] **AH-4** — Rutas: limiter solo en setup/login + Swagger
  - Acceptance: `RegisterRoutes` recibe el `gin.HandlerFunc` del limiter; `/auth/status` queda sin límite; setup/login lo aplican antes del handler; `docs/` documenta 403 y 429; wiring en `server.go`
  - Verify: `make test` (T9: 20 llamadas a status → 200); `make swagger`; bucle de 8 logins → `401×5` y `429`
  - Files: `internal/modules/auth/routes.go`, `internal/server/server.go`, `docs/`

- [ ] **AH-5** — Anti-timing en login (bcrypt dummy)
  - Acceptance: email inexistente ejecuta `bcrypt.CompareHashAndPassword` contra un hash fijo antes de devolver 401; mismo status/body que contraseña incorrecta; el hash dummy se calcula una sola vez
  - Verify: `make test` (T10); `make check`; medición manual de 20 intentos por caso (diferencia < 2×)
  - Files: `internal/modules/auth/service.go`, `internal/modules/auth/service_test.go`

- [ ] **AH-6** — Shared + `SetupPage` web
  - Acceptance: `SetupRequest.setupToken?` y `StatusResponse.setupRequiresToken` en shared (el re-export de `apps/web/src/types` es fino y no cambia); `SetupPage` muestra campo "Token de instalación" cuando `setupRequiresToken`, lo valida y lo envía en el body; si es `false`, muestra alerta de "Configuración inicial deshabilitada" y deshabilita el submit
  - Verify: `pnpm --filter @devtest/shared typecheck`; `make fe-check`; `make fe-test` (T11, T12); prueba manual contra backend con y sin `SETUP_TOKEN`
  - Files: `packages/shared/src/types/auth.types.ts`, `apps/web/src/features/auth/pages/SetupPage.vue`, `apps/web/src/features/auth/pages/SetupPage.test.ts`

- [ ] **AH-7** — Mobile: `auth.store` + `SetupScreen`
  - Acceptance: `setup(email, password, setupToken)` en el store; `bootstrap` guarda `setupRequiresToken`; `SetupScreen` pide el token cuando corresponde, muestra aviso y bloquea submit cuando está deshabilitado; usa `AuthInput` y tokens de diseño (sin hex)
  - Verify: `pnpm --filter mobile typecheck`; `pnpm --filter mobile test:run` (T13); prueba manual en Expo
  - Files: `apps/mobile/src/core/auth/auth.store.ts`, `apps/mobile/src/features/auth/screens/SetupScreen.tsx`, `apps/mobile/src/features/auth/screens/SetupScreen.test.tsx`

- [ ] **AH-8** — E2E + documentación del acoplamiento
  - Acceptance: fixture `setupToken` (`E2E_SETUP_TOKEN ?? 'e2e-setup-token'`); `critical-flow.spec.ts` rellena el campo cuando está visible; comentario/README indicando que el backend debe arrancar con el mismo `SETUP_TOKEN`; `apps/web/README.md` documenta el flujo de primera instalación
  - Verify: backend con `SETUP_TOKEN=e2e-setup-token` + `pnpm --filter dev-test-frontend test:e2e`
  - Files: `apps/web/e2e/fixtures.ts`, `apps/web/e2e/critical-flow.spec.ts`, `apps/web/README.md`

---

## Módulo `deps-supply-chain` — `specs/SPEC-deps-supply-chain.md`

- [ ] **DS-1** — Go 1.27.3 + `x/net` + `quic-go`
  - Acceptance: directiva `go 1.27.3`; `golang.org/x/net v0.60.0`; `quic-go v0.59.1`; `go mod tidy` sin ruido; `govulncheck ./...` sin vulnerabilidades alcanzables
  - Verify: `go get golang.org/x/net@v0.60.0 github.com/quic-go/quic-go@v0.59.1 && go mod tidy`; `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`; `make check && make test && make build`
  - Files: `go.mod`, `go.sum`

- [ ] **DS-2** — Overrides pnpm de transitivas vulnerables
  - Acceptance: `pnpm.overrides` con `follow-redirects@<1.16.1 → 1.16.1` y `shell-quote@<1.11.0 → 1.11.0`; `pnpm.auditConfig.ignoreCves` con las 3 CVEs sin parche (braces, node-forge, sprintf-js); **sin** `pnpm up` global; lock regenerado
  - Verify: `pnpm install`; `pnpm why follow-redirects` → 1.16.1; `pnpm why shell-quote` → ≥1.11.0; `pnpm audit --prod --audit-level high` → 0; `pnpm -r typecheck`
  - Files: `package.json`, `pnpm-lock.yaml`

- [ ] **DS-3** — Spike `markdown-it` en mobile + mitigación
  - Acceptance: se documenta si existe `react-native-markdown-display` compatible con `markdown-it >= 14.3.3`; si sí, bump + test de humo; si no, **no** se fuerza el override y `MarkdownText` trunca a 8 000 caracteres con `MAX_MARKDOWN_CHARS`; decisión registrada
  - Verify: `pnpm --filter mobile test:run` (humo: negrita, fence, lista, enlace); `pnpm --filter mobile typecheck`; `npx expo install --check`; `pnpm audit --prod`
  - Files: `apps/mobile/package.json`, `apps/mobile/src/core/components/MarkdownText.tsx`, `apps/mobile/src/core/components/MarkdownText.test.tsx` (nuevo), `pnpm-lock.yaml`, `specs/DECISIONS.md`

- [ ] **DS-4** — Vite 6.4.4+
  - Acceptance: `apps/web` resuelve `vite >= 6.4.4` dentro del major 6; build y tests verdes; sin cambios de configuración
  - Verify: `pnpm --filter dev-test-frontend why vite`; `make fe-build`; `make fe-test`
  - Files: `apps/web/package.json`, `pnpm-lock.yaml`

- [ ] **DS-5** — CI: `permissions`, `govulncheck`, audit y pinning por SHA
  - Acceptance: ambos workflows con `permissions: contents: read`; `govulncheck` en backend con versión fijada; `pnpm audit --prod --audit-level high` en frontend; las 4 actions pinneadas por SHA con comentario de versión
  - Verify: PR de prueba verde; revertir temporalmente un override → job rojo; `yq`/`act -l` para validar YAML
  - Files: `.github/workflows/backend.yml`, `.github/workflows/frontend.yml`

- [ ] **DS-6** — Dependabot + `make vuln`/`make audit` + docs
  - Acceptance: `.github/dependabot.yml` con gomod/npm/actions semanales, grupos y `expo*` ignorado; targets `vuln` y `audit` en `Makefile`; filas añadidas a `AGENTS.md`
  - Verify: `make vuln`; `make audit`; `yq '.' .github/dependabot.yml`; revisión de `AGENTS.md`
  - Files: `.github/dependabot.yml`, `Makefile`, `AGENTS.md`

---

## Módulo `client-hardening` — `specs/SPEC-client-hardening.md`

- [ ] **CH-1** — Mobile: HTTPS obligatorio fuera de `__DEV__`
  - Acceptance: `readApiBaseUrl()` lanza si `!__DEV__` y la URL no empieza por `https://` (incluido el default `http://localhost:8080`); el mensaje sanea credenciales `user:pass@`; en `__DEV__` todo sigue igual; `.env.example` de mobile documenta la regla
  - Verify: `pnpm --filter mobile typecheck`; `pnpm --filter mobile test:run` (T1–T5 con `jest.isolateModules` + `global.__DEV__`)
  - Files: `apps/mobile/src/core/config.ts`, `apps/mobile/src/core/config.test.ts` (nuevo), `apps/mobile/.env.example`

- [ ] **CH-2** — Web: plantillas de cabeceras + `index.html`
  - Acceptance: `apps/web/deploy/nginx-security.conf` y `apps/web/deploy/_headers` con CSP, HSTS, `nosniff`, `frame-ancestors`, `Referrer-Policy` y `Permissions-Policy`; `connect-src` usa `<API_ORIGIN>` (nunca `*`); `deploy/README.md` explica sustitución y verificación; `index.html` añade `meta referrer`
  - Verify: `get-content apps/web/index.html` contiene el meta; revisar que ambas plantillas tienen las 6 cabeceras; `make fe-build` sigue verde
  - Files: `apps/web/deploy/nginx-security.conf` (nuevo), `apps/web/deploy/_headers` (nuevo), `apps/web/deploy/README.md` (nuevo), `apps/web/index.html`

- [ ] **CH-3** — Documentar política de token y TTL
  - Acceptance: `apps/web/README.md` explica el riesgo de `localStorage`, las mitigaciones (CSP, 401 wipe, TTL) y cómo verificar las cabeceras; `.agents/web/architecture.md` gana una sección de seguridad del cliente; `.env.example` documenta el rango recomendado de `JWT_EXPIRY_HOURS` (8–12)
  - Verify: revisión de los 3 documentos; `git diff` sin cambios de código
  - Files: `apps/web/README.md`, `.agents/web/architecture.md`, `.env.example`

- [ ] **CH-4** — Verificación en staging (manual)
  - Acceptance: con el build desplegado, `curl -sI` muestra las 6 cabeceras con el origen correcto; la SPA carga sin errores de CSP en consola; una llamada a la API funciona (CORS + `connect-src`); mobile apunta a HTTPS
  - Verify: `curl -sI https://<dominio>/ | grep -Ei 'content-security-policy|strict-transport|nosniff|frame|referrer|permissions-policy'`; abrir DevTools → Console (0 errores CSP); login + sesión de estudio en staging
  - Files: — (verificación)

---

## Orden recomendado de ejecución

```
Ola 0: SH-1 → SH-2 → SH-3 → SH-4
       AZ-1 → AZ-2 → AZ-3 → AZ-4 → AZ-5
       HH-1 → HH-2 → HH-3 → HH-4 → HH-5
Ola 1: AH-1 → AH-2 → AH-3 → AH-4 → AH-5 → AH-6 → AH-7 → AH-8
       DS-1 → DS-2 → DS-3 → DS-4 → DS-5 → DS-6
       CH-1 → CH-2 → CH-3 → CH-4
```

`CH-4` se ejecuta cuando exista un despliegue; si no hay staging todavía, se deja
`[ ]` con una nota y se cierra en el primer deploy.
