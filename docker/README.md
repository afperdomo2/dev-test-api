# Docker — entornos, operación y despliegue

Esta carpeta contiene todo lo necesario para correr el proyecto en contenedores:

| Entorno | Archivo | Env | Para qué |
|---|---|---|---|
| **dev** | `docker/compose.dev.yaml` | `docker/.env` | Día a día con hot reload (Air + Vite HMR) |
| **prod-local** | `docker/compose.prod.yaml` | `docker/.env.prod` | Ensayar el despliegue real en tu PC |
| **dokploy** | `docker/compose.dokploy.yaml` | UI de Dokploy | Deploy en un VPS con Traefik |

Los tres stacks son **independientes** (no son overlays): podés tener dev y
prod-local levantados a la vez porque usan puertos y archivos de entorno
distintos. La contrapartida es que un cambio de servicios/healthchecks hay que
replicarlo en los tres YAML (ver [Mantenimiento](#mantenimiento)).

> Fuente de verdad para agentes: `.agents/devops/containers.md`.
> Runbook del APK: `docs/mobile-apk.md`.

## Requisitos

- Docker Desktop con WSL2 (probado con Docker 29 / Compose v5).
- `make` y `bash` (Git Bash en Windows). Los targets llaman a `scripts/*.sh`.
- ~4 GB de RAM libres para dev; ~2 GB para prod-local.

## Preparación (una vez)

```bash
make docker-env        # crea docker/.env desde .env.dev.example
make docker-env-prod   # crea docker/.env.prod desde .env.prod.example
```

Editá `docker/.env.prod` y reemplazá los `CHANGE_ME`. Generá secretos reales:

```bash
openssl rand -hex 32   # JWT_SECRET
openssl rand -hex 16   # DB_PASSWORD / SETUP_TOKEN
```

`docker/.env` y `docker/.env.prod` están gitignorados. Nunca commitearlos.

---

## 1) Desarrollo (`make docker-dev`)

```bash
make docker-dev
```

Equivale a `docker compose -f docker/compose.dev.yaml --env-file docker/.env watch`.
Es un comando **en primer plano**: muestra logs y sale con Ctrl+C.

| Servicio | URL desde el host | Notas |
|---|---|---|
| web (Vite) | http://localhost:3030 | HMR; `/api` proxya a `api:8080` |
| API | http://localhost:8082 | Swagger en `/swagger/index.html` |
| PostgreSQL | localhost:5435 | `devtest` / `devtest` / `dev_test_api` |

### Hot reload

No hay bind mounts. `docker compose watch` sincroniza las fuentes al contenedor:

- **Go**: `main.go`, `internal/`, `pkg/`, `docs/` → Air recompila.
- **Vue**: `apps/web/src`, `index.html`, `packages/shared/src` → Vite HMR.
- `go.mod`/`go.sum`/`package.json`/lockfile → rebuild de la imagen.
- `vite.config.ts`, `tsconfig.json`, `air.toml` → restart del contenedor.

Esto evita el problema clásico de eventos de archivo con bind mounts en Docker
Desktop/Windows. Si el watch deja de disparar: reiniciá `make docker-dev`; si
persiste, `make docker-dev-reset` (borra datos y cachés) y volvé a levantar.

> El **primer** `make docker-dev` descarga los módulos de Go en el volumen
> `go-mod-cache` y Air compila el proyecto: puede tardar varios minutos. La
> caché de compilación queda en `go-cache` y la de Vite en `vite-cache`, así que
> los arranques siguientes son de segundos. `make docker-dev-reset` borra los
> tres volúmenes y vuelve a empezar de cero.

### Datos y seed

```bash
make docker-seed          # 50 preguntas de db/seeds/seed_questions.sql
make docker-dev-reset     # borra contenedores + volumen db-data (destructivo)
```

La DB arranca vacía. El primer admin se crea desde la web (`/setup`) con el
`SETUP_TOKEN` de `docker/.env`. Con la DB vacía el token es **obligatorio**.

### Sin watch

```bash
make docker-dev-up        # segundo plano + healthchecks
make docker-dev-down
```

---

## 2) Prod-local (`make docker-prod`)

Construye las imágenes release, levanta el stack y espera los healthchecks.

```bash
make docker-prod
make docker-smoke          # valida contrato HTTP y cabeceras
```

- **Un solo origen**: http://localhost:8088 sirve la SPA y proxya `/api`,
  `/health` y `/swagger` hacia `api:8080`.
- La DB **no** se publica al host.
- `api` corre read-only, non-root, sin capabilities; `web` non-root con
  `cap_drop: ALL`.
- `GIN_MODE=release`, Swagger apagado por defecto (`SWAGGER_ENABLED=false`).

Flujo recomendado la primera vez:

1. `make docker-env-prod` y editar `docker/.env.prod`.
2. `make docker-prod`.
3. Abrir http://localhost:8088, crear el admin con el `SETUP_TOKEN`.
4. Vaciar `SETUP_TOKEN` en `docker/.env.prod` y `make docker-prod` de nuevo.
5. `make docker-smoke`.

### Backup y restore

```bash
make docker-backup                              # docker/backups/devtest-<fecha>.sql.gz
make docker-restore FILE=docker/backups/xxx.sql.gz
```

El dump usa `--clean --if-exists`: el restore **reemplaza** los objetos de la DB
(y pide confirmación). Probá el restore antes de confiar en el backup.

### E2E contra prod-local

```bash
make docker-prod-reset   # DB limpia (la suite hace setup del primer admin)
make docker-prod
make e2e-docker
```

Playwright usa `PLAYWRIGHT_BASE_URL` y **no** arranca Vite cuando está definida.

---

## 3) Dokploy

### Topología

Un solo stack con `db` + `api` + `web`:

- `web` se une a `dokploy-network` (externa) para que Traefik lo alcance.
- `api` y `db` viven en la red interna del compose; la DB no se expone.
- Un único dominio: `https://app.tu-dominio.com`.
- El volumen `db-data` es nombrado (requisito para los backups de Dokploy).

### Pasos

1. **DNS**: creá un registro `A` de `app.tu-dominio.com` → IP del VPS.
2. En Dokploy: **Project → Create Service → Docker Compose**.
3. **Compose Path**: `docker/compose.dokploy.yaml`
   (el contexto de build `..` resuelve a la raíz del repo).
4. **Environment**: pegá las variables de `docker/env/.env.dokploy.example` con
   secretos reales. `APP_ORIGIN` es obligatorio (`https://…`, sin `/` final).
5. **Domains**: agregá el dominio al servicio **`web`**, puerto **`8080`**.
   Dejá que Dokploy inyecte las labels de Traefik (no las escribas a mano).
6. **Deploy**. El primer arranque exige `SETUP_TOKEN` (DB vacía); si falta, el
   API sale con error y el contenedor reinicia.
7. Crear el admin en `https://app.tu-dominio.com/setup`.
8. Borrar `SETUP_TOKEN` en Dokploy y **redeployar**.
9. Verificar: `https://app.tu-dominio.com/health` → `{"status":"ok"}`.
10. Activar **Volume Backup** sobre `db-data` y probar un restore.

### Variables que no pueden faltar

| Variable | Ejemplo | Nota |
|---|---|---|
| `APP_ORIGIN` | `https://app.tu-dominio.com` | CORS + CSP; sin `/` final |
| `DB_PASSWORD` | secreto fuerte | La usan `db` y `api` |
| `JWT_SECRET` | `openssl rand -hex 32` | Firma de tokens |
| `SETUP_TOKEN` | ≥16 chars | Solo hasta crear el admin |
| `TRUSTED_PROXIES` | `172.16.0.0/12,10.0.0.0/8` | Sin esto el rate-limit es global |
| `TZ` | `America/Bogota` | Define las 08:00 del scheduler |

Opcionales: `AI_API_URL`, `AI_API_KEY`, `AI_MODEL`, `EXPO_ACCESS_TOKEN`,
`SWAGGER_ENABLED`, `VITE_API_BASE_URL` (vacío = same-origin).

### Actualizar / rollback

- **Update**: push a la rama configurada → Deploy (o botón Deploy).
- Los cambios en `db-data` no se tocan al redeployar.
- **Nunca** ejecutes `docker compose down -v` en el VPS: borra la DB.
- Las migraciones goose corren solas al arrancar `api` (con lock de sesión, apto
  para más de una instancia).

---

## Variables y secretos

| Archivo | Se usa en | Se versiona |
|---|---|---|
| `docker/.env` | dev (auto, project dir) | No |
| `docker/.env.prod` | prod-local (`--env-file`) | No |
| `deploy` UI | Dokploy | No (se pega en la UI) |
| `docker/env/*.example` | plantillas + `make docker-config` | Sí |

Reglas:

- **Secretos**: `DB_PASSWORD`, `JWT_SECRET`, `SETUP_TOKEN`, `AI_API_KEY`,
  `EXPO_ACCESS_TOKEN`. Nunca en `ARG`/`ENV` de una imagen.
- **Públicos**: `VITE_*` y `EXPO_PUBLIC_*` se compilan en el bundle. Solo URLs.
- `DB_*` y `POSTGRES_*` salen del mismo valor para que no se desincronicen.
- `docker compose config` valida los tres YAML con los `.example`, sin secretos
  reales: `make docker-config`.

---

## Puertos

| Host | Contenedor | Entorno | Por qué |
|---|---|---|---|
| 3030 | web:3000 | dev | 3000 suele estar ocupado por Vite nativo |
| 8082 | api:8080 | dev | 8080/8081 suelen estar ocupados por Air u otro proceso |
| 5435 | db:5432 | dev | 5432/5433/5434 suelen estar ocupados por Postgres local |
| 8088 | web:8080 | prod-local | Único punto de entrada |
| — | db:5432 | prod/dokploy | Solo red interna |

Se cambian con `WEB_PORT_HOST`, `API_PORT_HOST`, `DB_PORT_HOST` en el `.env`
correspondiente.

---

## Conectarse a la base de datos

### Dentro del servidor (recomendado)

```bash
docker compose -f docker/compose.prod.yaml --env-file docker/.env.prod exec db \
  psql -U devtest -d dev_test_api
```

En Dokploy, usá la terminal del servicio `db` o `docker compose exec`.

### Desde tu PC a la DB de la nube

La DB **no** debe publicarse. En orden de preferencia:

1. **Túnel SSH** (no expone nada al mundo):

   ```bash
   ssh -L 5433:localhost:5432 root@TU_VPS
   # Luego conectá DBeaver/psql a localhost:5433
   ```

   Si el contenedor no publica 5432, el túnel debe apuntar al contenedor:
   `ssh -L 5433:<ip-contenedor>:5432 root@TU_VPS` (o publicá temporalmente el
   puerto con `docker compose port`).

2. **External Port de Dokploy** con firewall/allowlist y solo durante una
   sesión de mantenimiento. Nunca lo dejes abierto.

3. **Backup + restore local** si solo necesitás datos:

   ```bash
   make docker-backup
   # copiar el .sql.gz y restaurarlo en un Postgres local
   ```

> El APK **no** se conecta a Postgres. Solo habla con el API por HTTPS. Ver
> `docs/mobile-apk.md`.

---

## Imágenes

### Backend (`docker/backend/Dockerfile`)

| Stage | Base | Uso |
|---|---|---|
| `dev` | `golang:1.27-alpine` + Air | `compose.dev.yaml` |
| `build` | `golang:1.27-alpine` | compila el binario |
| `runtime` | `gcr.io/distroless/static-debian13:nonroot` | prod/dokploy |

- `CGO_ENABLED=0`, `-trimpath`, `-ldflags="-s -w -buildid="`.
- `USER 65532`, sin shell/package manager (busybox solo para el healthcheck).
- Distroless incluye CA certs + tzdata.
- Tamaño: **~60–65 MB** (el binario pesa ~37 MB por Gin + GORM + goose + Swagger).

### Web (`docker/web/Dockerfile`)

| Stage | Base | Uso |
|---|---|---|
| `dev` | `node:22-bookworm-slim` | `compose.dev.yaml` (Vite) |
| `build` | idem | `vue-tsc` + `vite build` |
| `runtime` | `nginxinc/nginx-unprivileged:1.29-alpine-slim` | prod/dokploy |

- pnpm `10.30.3` vía corepack; instalación filtrada al web.
- Debian (no Alpine) por binarios nativos (`sass-embedded`, rollup, esbuild).
- nginx non-root en `:8080`, CSP/HSTS/nosniff, gzip, cache de assets,
  `client_max_body_size 9m`, proxy a `api:8080`.
- `VITE_API_BASE_URL` vacío = same-origin (build arg).
- Tamaño: **~30 MB** (la variante `-slim` evita ~60 MB de módulos que no usamos).

### Build manual / escaneo

```bash
docker build -f docker/backend/Dockerfile -t devtest-api .
docker build -f docker/web/Dockerfile -t devtest-web .
docker run --rm --entrypoint /busybox devtest-api id   # uid=65532(nonroot)
```

Opcional (recomendado en CI): `trivy image devtest-api` y `hadolint` sobre los
Dockerfiles. Los tags de base están pinneados a versión; actualizarlos es una
tarea periódica (Renovate/Dependabot).

---

## Operación diaria

```bash
make docker-prod-logs                 # logs del stack prod-local
docker compose -f docker/compose.prod.yaml --env-file docker/.env.prod ps
docker compose -f docker/compose.prod.yaml --env-file docker/.env.prod restart api
```

- **Healthchecks**: `db` (`pg_isready`), `api` (`/health`), `web` (`/healthz`).
- `api` tiene `stop_grace_period: 45s` para el shutdown graceful (scheduler +
  worker pool).
- `restart: unless-stopped` en todos los servicios.
- Logs rotados a 10 MB × 3 por contenedor.

---

## Troubleshooting

| Síntoma | Causa / solución |
|---|---|
| `port is already allocated` | Hay un Vite/Air/Postgres nativo. Cambiá `*_PORT_HOST` o pará los procesos locales |
| `api` reinicia en loop | DB no healthy o falta `SETUP_TOKEN` con DB vacía. `make docker-prod-logs` |
| 429 en `/auth/login` | `TRUSTED_PROXIES` sin el CIDR del proxy: todo comparte bucket |
| Import falla 413 | nginx sin `client_max_body_size 9m` (no lo bajes) |
| 502 en `/api` | `api` no healthy; revisar logs y `depends_on` |
| Swagger roto vía nginx | CSP estricta. Usalo en dev (`:8082`) o dejá `SWAGGER_ENABLED=false` |
| SPA muestra 404 al refrescar | `try_files … /index.html` roto en `nginx.conf` |
| Datos desaparecieron | Se ejecutó `down -v`. Restaurar backup |
| Watch no recompila | Reiniciar `make docker-dev`; si persiste, `make docker-dev-reset` (borra cachés) |
| `permission denied` al hablar con Docker desde WSL | Docker Desktop → Settings → Resources → File sharing/WSL integration |
| Tras cambiar de dominio la app no llama al API | Rebuild de web; en mobile, recompilar APK con la nueva URL |
| Hora de notificaciones corrida | `TZ` del contenedor (default `America/Bogota`) |

---

## Checklist de seguridad

- [ ] Secretos solo en `.env`/UI de Dokploy; nada en imágenes ni en git.
- [ ] Contenedores non-root (`65532`, `101`).
- [ ] `cap_drop: ALL` + `no-new-privileges` en api/web.
- [ ] DB sin puerto público en prod/dokploy.
- [ ] `TRUSTED_PROXIES` con el CIDR real.
- [ ] `HSTS_ENABLED=true` solo detrás de HTTPS (Dokploy).
- [ ] CSP/HSTS/nosniff verificados con `make docker-smoke`.
- [ ] Imágenes pinneadas y escaneadas antes de una release.
- [ ] Backup probado (dump **y** restore).
- [ ] `SETUP_TOKEN` eliminado después del primer admin.

---

## Mantenimiento

Los tres compose comparten servicios casi idénticos. Si cambiás algo de
`api`/`web`/`db` (env, healthcheck, volúmenes, security), revisá **los tres** y
actualizá:

1. `docker/compose.dev.yaml`, `docker/compose.prod.yaml`, `docker/compose.dokploy.yaml`.
2. `docker/README.md` (si cambia un comando/variable/topología).
3. `.agents/devops/containers.md` (si cambia una invariante).
4. `AGENTS.md` solo si cambia una regla que el agente deba aplicar.
5. `specs/DECISIONS.md` si es una decisión nueva.

Validación mínima antes de cerrar un cambio:

```bash
make docker-config
make docker-prod
make docker-smoke
```

## APK

El APK se compila fuera de Docker, a mano, con EAS cloud. La URL del backend se
**compila dentro del APK**. Runbook completo: [`docs/mobile-apk.md`](../docs/mobile-apk.md).
