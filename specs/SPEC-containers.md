# SPEC — Contenedores (dev, prod-local y Dokploy)

> Estado: implementado el stack; deploy en Dokploy y build del APK quedan como
> tareas manuales (`tasks/containers-todo.md`).
> Decisiones: D-44…D-47 en `specs/DECISIONS.md`.
> Fuente de verdad operativa: `docker/README.md` y `.agents/devops/containers.md`.

## Objetivo

Correr el proyecto completo (API + web + Postgres) en contenedores con tres
propósitos distintos, sin depender de instalar Go/Node/Postgres en cada máquina:

1. **dev**: hot reload real (Air + Vite HMR) con un solo comando.
2. **prod-local**: ensayar el despliegue real en la PC (imágenes release, DB en
   el compose, un solo origen, healthchecks, backups).
3. **dokploy**: stack listo para desplegar en un VPS con Dokploy + Traefik.

Mobile queda fuera de Docker: Metro corre en el host y el APK se compila a mano
con EAS cloud.

## Decisiones que implementa

| ID | Decisión |
|---|---|
| D-44 | Tres entornos Compose independientes (dev / prod-local / dokploy), con Compose Watch en dev |
| D-45 | Un solo origen: nginx sirve la SPA y proxya `/api` → `api:8080` |
| D-46 | Imágenes multi-stage non-root (distroless / nginx unprivileged); DB dentro del compose de Dokploy con volumen nombrado |
| D-47 | Mobile fuera de Docker; `EXPO_PUBLIC_API_BASE_URL` es build-time, HTTPS y exige rebuild al cambiar de backend |

## Entregables

### Infraestructura

- `docker/backend/Dockerfile` (+ `Dockerfile.dockerignore`, `air.toml`):
  stages `dev` / `build` / `runtime`; binario estático; distroless nonroot;
  healthcheck.
- `docker/web/Dockerfile` (+ `.dockerignore`): Node 22 + pnpm 10.30.3 (build),
  nginx unprivileged (runtime); `VITE_API_BASE_URL` como build arg.
- `docker/web/nginx.conf`: SPA fallback, proxy `/api` + `/health` + `/swagger`,
  cache de assets, gzip, `client_max_body_size 9m`, CSP/HSTS.
- `apps/web/deploy/nginx-security.conf`: cabeceras de la SPA
  (SPEC-client-hardening CH-2/CH-3).
- `docker/compose.{dev,prod,dokploy}.yaml` + `docker/env/.env.*.example`.
- `.dockerignore` raíz + ignores por Dockerfile.

### Código

- `apps/web/vite.config.ts`: `VITE_PROXY_TARGET` y `host: true`.
- `apps/web/index.html`: meta `referrer`.
- `apps/web/playwright.config.ts`: `PLAYWRIGHT_BASE_URL` (E2E contra compose).
- `apps/mobile/app.json`: `android.package` + `versionCode`.
- `apps/mobile/eas.json`: perfil `preview` (APK) + `production` (AAB).

### Operación

- Make: `docker-env`, `docker-env-prod`, `docker-dev`, `docker-dev-up/down/reset`,
  `docker-prod`, `docker-prod-down/reset/logs`, `docker-build`, `docker-config`,
  `docker-seed`, `docker-smoke`, `docker-backup`, `docker-restore`,
  `e2e-docker`, `mobile-apk`.
- Scripts (strict mode): `docker-env.sh`, `docker-seed.sh`, `docker-smoke.sh`,
  `docker-backup.sh`, `docker-restore.sh`, `e2e-docker.sh`, `preflight-apk.sh`.
- CI: build de las dos imágenes + validación de los compose.

### Documentación

- `docker/README.md` (runbook completo).
- `.agents/devops/containers.md` (invariantes para agentes).
- `docs/mobile-apk.md` (runbook del APK).
- `AGENTS.md`, `README.md`, `specs/DECISIONS.md`.

## Testing matrix

| # | Caso | Cómo | Esperado |
|---|---|---|---|
| T1 | Los 3 compose parsean | `make docker-config` | exit 0, sin warnings de interpolación |
| T2 | Imagen API | `docker build …` | build OK; `id` → uid 65532; ~60–65 MB |
| T3 | Imagen web | `docker build …` | build OK; nginx nonroot; ~30 MB |
| T4 | Hot reload Go | `make docker-dev` + editar `internal/**` | Air recompila y reinicia |
| T5 | Hot reload web/shared | editar `apps/web/src/**` y `packages/shared/src/**` | Vite HMR sin restart |
| T6 | Arranque prod-local | `make docker-prod` | 3 servicios healthy; `:8088` responde |
| T7 | Contrato HTTP | `make docker-smoke` | health, auth/status, SPA y 5 headers |
| T8 | Persistencia | reiniciar stack | los datos siguen en `db-data` |
| T9 | Backup/restore | `make docker-backup` + `make docker-restore` | restore deja la DB usable |
| T10 | Migraciones | arrancar API con DB vacía | migraciones aplicadas y admin setup disponible |
| T11 | Seed | `make docker-seed` | 50 preguntas; idempotente al repetir |
| T12 | E2E | `make docker-prod-reset && make docker-prod && make e2e-docker` | suite verde |
| T13 | Dokploy | deploy real | dominio HTTPS, `/health` OK, volumen persistente |
| T14 | APK | `make mobile-apk` | APK instala, login y request real |
| T15 | Cambio de URL | editar `eas.json` | preflight falla con placeholder / pasa con HTTPS real |

## Boundaries

**Always**

- `make docker-config` tras tocar cualquier compose/Dockerfile.
- Replicar cambios de servicios en los 3 compose.
- Mantener non-root, DB no publicada, `TRUSTED_PROXIES` y `TZ`.
- Documentar el cambio en `docker/README.md` y/o `.agents/devops/containers.md`.

**Ask first**

- Agregar un cuarto entorno (staging).
- Publicar la DB en algún entorno.
- Cambiar la topología a dos dominios o a una DB gestionada en Dokploy.
- Dockerizar Metro/Android SDK.
- Introducir un registry de imágenes.

**Never**

- Secretos en Dockerfiles, compose versionado o build args.
- `down -v` en prod-local/Dokploy.
- `container_name` fijo.
- Bajar de non-root a root "para que funcione".
- `connect-src *` o desactivar CSP en producción.
- Asumir que el APK toma la URL del shell/EAS distinto de `eas.json`.

## Success criteria

1. `make docker-dev` levanta los 3 servicios y el hot reload funciona en Go,
   Vue y shared sin bind mounts.
2. `make docker-prod` levanta el stack con imágenes release y `:8088` sirve SPA +
   API con CSP/HSTS/nosniff.
3. `make docker-smoke` pasa.
4. `make docker-backup` + `make docker-restore` recuperan la DB.
5. El compose de Dokploy despliega con un dominio HTTPS, DB en volumen nombrado
   y `SETUP_TOKEN` retirado tras el primer admin.
6. El APK se compila a mano con la URL HTTPS correcta y se instala en Android.
7. Un agente puede responder "cómo levanto dev / prod / Dokploy / APK" leyendo
   `.agents/devops/containers.md` sin mirar el código.

## Open questions

- ¿Staging en Dokploy (segunda rama) o se mantiene solo prod? (hoy: no).
- ¿Se migra la DB de Dokploy a un recurso gestionado cuando haya backups
  automáticos verificados? (hoy: DB en compose).
- ¿Se automatiza el APK en CI/EAS Workflows? (hoy: manual).
- ¿Se añade `trivy`/`hadolint` como gate de CI? (hoy: opcional, documentado).
