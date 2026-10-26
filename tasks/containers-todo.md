# Tasks — Contenedores

> Spec: `specs/SPEC-containers.md` · Guía: `docker/README.md` ·
> Invariantes: `.agents/devops/containers.md`

## Fase 0 — Decisiones y documentación

- [x] **CT-0.1** Decisiones D-44…D-47 en `specs/DECISIONS.md`
  - Hecho: 3 entornos; same-origin; DB en compose + non-root; mobile fuera de Docker con URL build-time.
- [x] **CT-0.2** Spec + tasks
  - Hecho: `specs/SPEC-containers.md`, este archivo y enlaces desde `AGENTS.md` (domains + documentation map).

## Fase 1 — Imágenes

- [x] **CT-1.1** Dockerfile del backend con stages `dev`/`build`/`runtime`
  - Acceptance: Air en dev, binario estático, distroless nonroot, healthcheck.
  - Hecho: `docker/backend/Dockerfile` + `air.toml` + `Dockerfile.dockerignore`; build verificado; `uid=65532(nonroot)`; runtime ~61 MB, dev ~384 MB (módulos en volumen, no en la imagen).
- [x] **CT-1.2** Dockerfile del web con stages `dev`/`build`/`runtime`
  - Acceptance: Vite HMR en dev, build con `VITE_API_BASE_URL`, nginx non-root.
  - Hecho: `docker/web/Dockerfile` + `Dockerfile.dockerignore` + `nginx.conf`; build verificado.
- [x] **CT-1.3** `.dockerignore` por Dockerfile + raíz
  - Hecho: contexto del API sin `apps/`/`node_modules`; contexto web sin `internal/`/`docs/`; mobile solo `package.json`.
- [x] **CT-1.4** Cabeceras de seguridad en el borde
  - Acceptance: CSP/HSTS/nosniff/XFO/Referrer/Permissions en la SPA.
  - Hecho: `apps/web/deploy/nginx-security.conf` incluido por `docker/web/nginx.conf`; Swagger con CSP relajada solo en su location.

## Fase 2 — Dev con hot reload

- [x] **CT-2.1** `docker/compose.dev.yaml` con Compose Watch
  - Acceptance: sin bind mounts; Air/Vite reaccionan a cambios de Go, Vue y shared.
  - Hecho: `develop.watch` para `main.go`, `internal/`, `pkg/`, `docs/`, `apps/web/src`, `packages/shared/src`; `rebuild` en manifests. Verificado: Air rebuild con cambio en Go y Vite reload con cambios en web/shared (`docker compose watch --no-up`).
- [x] **CT-2.2** Vite configurable para Docker
  - Hecho: `VITE_PROXY_TARGET` (default `localhost:8080`, Docker `api:8080`) y `host: true`.
- [x] **CT-2.3** Make targets + env de dev
  - Hecho: `docker-env`, `docker-dev`, `docker-dev-up/down/reset`, `docker-seed`; `docker/env/.env.dev.example`; puertos 3030/8082/5435; volúmenes `go-mod-cache`, `go-cache`, `vite-cache`.

## Fase 3 — Prod-local

- [x] **CT-3.1** `docker/compose.prod.yaml` same-origin
  - Hecho: db + api (read-only, cap_drop) + web (`:8088`), healthchecks, `stop_grace_period`, logs rotados, DB sin puerto.
- [x] **CT-3.2** Nginx proxy + cache + límites
  - Hecho: `/api`, `/health`, `/swagger`; `client_max_body_size 9m`; assets 1y; `index.html` sin cache; DNS runtime (`resolver 127.0.0.11`).
- [x] **CT-3.3** Smoke, seed, backup y restore
  - Hecho: `scripts/docker-smoke.sh`, `docker-seed.sh`, `docker-backup.sh`, `docker-restore.sh` (strict mode) + targets.
- [x] **CT-3.4** E2E contra el stack
  - Hecho: `PLAYWRIGHT_BASE_URL` en `playwright.config.ts` (omite `webServer`), `scripts/e2e-docker.sh`, `make e2e-docker`.
- [x] **CT-3.5** Ejecutar el ciclo completo en la PC
  - Acceptance: `make docker-prod && make docker-smoke && make docker-backup` + restore verificado.
  - Hecho: `make docker-prod` (3 healthy en `:8088`), `make docker-smoke` verde, setup + login vía nginx OK, backup y restore con tabla `restore_probe` verificados.

## Fase 4 — Dokploy

- [x] **CT-4.1** `docker/compose.dokploy.yaml`
  - Hecho: web en `dokploy-network`, api/db en red interna, `expose` (sin `ports`), volumen `db-data`, env desde UI, `APP_ORIGIN` obligatorio.
- [x] **CT-4.2** Runbook de deploy
  - Hecho: `docker/README.md` sección Dokploy (DNS, compose path, envs, dominio, setup token, backups, update/rollback).
- [ ] **CT-4.3** Primer deploy real
  - Acceptance: dominio HTTPS, `/health` OK, admin creado, `SETUP_TOKEN` retirado, volume backup probado.
  - Pendiente: requiere VPS/Dokploy del humano.

## Fase 5 — APK

- [x] **CT-5.1** Config EAS + identidad Android
  - Hecho: `apps/mobile/eas.json` (preview APK / production AAB), `android.package` + `versionCode` en `app.json`.
- [x] **CT-5.2** Preflight + runbook
  - Hecho: `scripts/preflight-apk.sh` (HTTPS, placeholder, health, auth/status, versionCode) + `docs/mobile-apk.md` + `make mobile-apk`.
- [ ] **CT-5.3** Primer APK instalado
  - Acceptance: `eas init`, URL real en `eas.json`, APK instalado, login + request real.
  - Pendiente: cuenta Expo y backend desplegado.

## Fase 6 — CI

- [x] **CT-6.1** Workflow de imágenes + compose
  - Hecho: `.github/workflows/containers.yml` (buildx de api/web + `compose config` de los 3 entornos).
- [ ] **CT-6.2** Endurecer CI (opcional)
  - Acceptance: `trivy image` y `hadolint` como gate; actualización automatizada de tags base.
  - Pendiente: decidir en una spec aparte.
