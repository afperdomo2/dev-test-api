# Containers (Docker) — fuente de verdad

Guía operativa completa: **`docker/README.md`**. Este archivo es el contexto que
un agente debe leer antes de tocar Docker, compose o el build del APK. Si hay
conflicto entre este archivo y `docker/README.md`, gana este (y hay que corregir
el README).

## Entornos (3, no agregar más sin spec)

| Entorno | Compose | Env | Propósito |
|---|---|---|---|
| dev | `docker/compose.dev.yaml` | `docker/.env` (auto en project dir) | Hot reload con Compose Watch |
| prod-local | `docker/compose.prod.yaml` | `docker/.env.prod` (`--env-file`) | Ensayo realista en la PC |
| dokploy | `docker/compose.dokploy.yaml` | UI de Dokploy | Deploy real (VPS + Traefik) |

No hay staging. CI solo construye imágenes y valida compose; no levanta stacks.

## Mapa de archivos

```
docker/
├── backend/
│  ├── Dockerfile              # stages: base → deps → dev / build → runtime
│   ├── Dockerfile.dockerignore # contexto mínimo del API
│   └── air.toml               # Air para Linux (el de la raíz es Windows)
├── web/
│  ├── Dockerfile              # deps → dev / build → runtime (nginx)
│  ├── Dockerfile.dockerignore
│  └── nginx.conf              # SPA + proxy /api + CSP
├── compose.dev.yaml
├── compose.prod.yaml
├── compose.dokploy.yaml
├── env/.env.{dev,prod,dokploy}.example
├── backups/                   # gitignored (dumps locales)
└── README.md
apps/web/deploy/nginx-security.conf   # CSP/headers de la SPA
docs/mobile-apk.md                    # runbook del APK
scripts/{docker-env,docker-seed,docker-smoke,docker-backup,docker-restore,e2e-docker,preflight-apk}.sh
```

## Comandos

| Comando | Qué hace |
|---|---|
| `make docker-env` | Copia `.env.dev.example` → `docker/.env` (no sobrescribe) |
| `make docker-env-prod` | Copia `.env.prod.example` → `docker/.env.prod` |
| `make docker-dev` | `docker compose watch` (hot reload; Ctrl+C sale) |
| `make docker-dev-up` | Levanta dev en background con `--wait` |
| `make docker-dev-reset` | `down -v` del stack dev (borra datos) |
| `make docker-prod` | Build + up + healthchecks de prod-local |
| `make docker-prod-reset` | `down -v` del stack prod-local (borra datos) |
| `make docker-config` | Valida los 3 compose con sus env de ejemplo |
| `make docker-seed` | Siembra `db/seeds/seed_questions.sql` en dev |
| `make docker-smoke` | Verifica nginx, /health, auth/status, SPA y headers |
| `make docker-backup` | `pg_dump` comprimido de prod-local |
| `make docker-restore FILE=…` | Restore destructivo (pide confirmación) |
| `make e2e-docker` | Playwright contra prod-local (DB limpia) |
| `make mobile-apk` | Preflight + `eas build` APK (manual, cloud) |

## Invariantes (no romper)

1. **Un solo origen**: la SPA y `/api` comparten dominio vía nginx. No agregar
   CORS de otro origen ni `connect-src *`; si se separa el API, actualizar
   `apps/web/deploy/nginx-security.conf` y `CORS_ALLOWED_ORIGINS`.
2. **Imágenes non-root**: API `USER 65532` (distroless), nginx unprivileged
   (uid 101). No cambiar a `root` para "arreglar" permisos.
3. **DB nunca publicada** en prod-local ni Dokploy. En dev sí (5435) para
   inspección local. No agregar `ports` a `db` en prod/dokploy.
4. **Volumen nombrado** `db-data`. Nunca `docker compose down -v` en prod ni
   Dokploy. Los backups de Dokploy solo funcionan con volúmenes nombrados.
5. **Sin `container_name`**: rompe logs/métricas de Dokploy y el escalado.
6. **Sin secretos en `ARG`/`ENV` de build**. `VITE_*` y `EXPO_PUBLIC_*` son
   públicos por definición: solo URLs.
7. **`EXPO_PUBLIC_API_BASE_URL` es build-time**. Si cambia la URL del backend hay
   que recompilar el APK. EAS cloud usa `apps/mobile/eas.json`, no el shell.
   Debe ser `https://` fuera de desarrollo (D-41).
8. **`TRUSTED_PROXIES`** debe incluir el CIDR del proxy (Docker `172.16.0.0/12`,
   Dokploy `10.0.0.0/8`). Si queda vacío, el rate-limit de auth cuenta todo como
   una sola IP.
9. **`client_max_body_size 9m`** en nginx: el backend acepta imports de hasta 8 MiB.
10. **`TZ`** define las 08:00 del scheduler. Default de los compose:
    `America/Bogota`. No asumir UTC.
11. **Dockerfiles/compose en LF** (`.gitattributes`). Un CRLF en un `RUN ... \`
    o en un YAML rompe el build en Linux.
12. **`docs/` se compila dentro del binario** (`main.go` importa `docs`). No lo
    excluyas del contexto del backend ni de la imagen.

## Verificación al tocar Docker

```bash
make docker-config                 # los 3 compose parsean
make docker-build                  # imágenes prod-local
make docker-prod                   # up + healthchecks
make docker-smoke                  # contrato HTTP + headers
make docker-backup                 # backup real
```

Si el cambio toca `vite.config.ts`, `nginx.conf` o el proxy, además levantar dev
y comprobar HMR:

```bash
make docker-dev
# editar internal/... (Air), apps/web/src/... (Vite), packages/shared/src/...
```

## Detalles de las imágenes

- **API** (`docker/backend/Dockerfile`): Go 1.27 Alpine (build) + distroless
  static Debian 13 nonroot (runtime). Binario estático, `-trimpath`,
  `-ldflags="-s -w"`. Busybox solo para `wget` del `HEALTHCHECK`. Imagen
  ~60–65 MB (el binario son ~37 MB por Gin+GORM+goose+Swagger). Sin shell real,
  sin package manager. En dev no lleva módulos: los volúmenes `go-mod-cache` y
  `go-cache` (más `vite-cache`) evitan repetir descarga/compilación.
- **Web** (`docker/web/Dockerfile`): Node 22 bookworm-slim (build, glibc por
  `sass-embedded`/rollup/esbuild) + nginx unprivileged Alpine **slim** (runtime).
  Vite build con `VITE_API_BASE_URL` vacío = same-origin. Imagen ~30 MB.
- **pnpm**: corepack pinneado a `10.30.3` (el `packageManager` del repo). Store
  en `/pnpm/store` con cache mount.
- **Contexto**: cada Dockerfile tiene su `Dockerfile.dockerignore`; BuildKit lo
  usa en lugar del `.dockerignore` raíz.

## Compose Watch (dev)

- No hay bind mounts: `watch` sincroniza al contenedor. Esto evita inotify roto
  en Docker Desktop/Windows.
- API: sync de `main.go`, `internal/`, `pkg/`, `docs/`; `rebuild` en
  `go.mod`/`go.sum`; `sync+restart` en `air.toml`.
- Web: sync de `apps/web/src`, `index.html`, `packages/shared/src`;
  `sync+restart` en `vite.config.ts`/`tsconfig.json`; `rebuild` en manifests/lock.
- Si el watch no dispara: reiniciar `make docker-dev`; si persiste,
  `make docker-dev-reset`. Verificar que el archivo esté en un path del
  `develop.watch`.

## Dokploy (resumen)

- Compose Path: `docker/compose.dokploy.yaml`.
- Dominio: servicio `web`, puerto `8080` (UI de Dokploy; Traefik inyecta labels).
- `web` está en `dokploy-network`; `api`/`db` en la red interna del compose.
- Variables: pegar `docker/env/.env.dokploy.example` con secretos reales.
  `APP_ORIGIN` es obligatorio (https, sin `/` final).
- Primer deploy: `SETUP_TOKEN` obligatorio → crear admin → borrar token →
  redeploy.
- Backups: Volume Backup de Dokploy sobre `db-data` o `make docker-backup`
  apuntando al stack. Probar el restore, no solo el dump.

## Troubleshooting rápido

| Síntoma | Causa / fix |
|---|---|
| `port is already allocated` | Puertos nativos 3000/8080/8081/5432/5433; Docker usa 3030/8082/5435/8088 |
| API reinicia en loop al arrancar | DB no lista o falta `SETUP_TOKEN` con DB vacía |
| 429 global en login | `TRUSTED_PROXIES` vacío o sin el CIDR del proxy |
| Import falla con 413 | nginx sin `client_max_body_size 9m` |
| APK no conecta | URL de `eas.json` incorrecta o backend sin HTTPS |
| Swagger roto vía nginx | CSP estricta; usar dev (`:8082`) o desactivar Swagger |
| Datos perdidos | `down -v`. En prod/Dokploy nunca usarlo |
| Watch no recompila | Reiniciar `make docker-dev` / `docker-dev-reset` |

## Regla de doc sync

Todo cambio en `docker/` debe actualizar, en el mismo cambio:

1. `docker/README.md` si cambia un comando, variable o topología.
2. Este archivo si cambia una invariante o el mapa de archivos.
3. `AGENTS.md` solo si cambia una regla que el agente deba aplicar.
4. `specs/DECISIONS.md` si se toma una decisión nueva (D-44+).
