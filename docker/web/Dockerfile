# syntax=docker/dockerfile:1.7

# ─────────────────────────────────────────────────────────────
# dev-test-api — imagen del frontend (Vue 3 + Vite + Vuetify)
#
# Stages:
#   dev     → Vite dev server con HMR. Lo usa docker/compose.dev.yaml
#   build   → build estático (vue-tsc + vite build)
#   runtime → nginx unprivileged (non-root, :8080)
#
# Build manual:
#   docker build -f docker/web/Dockerfile -t devtest-web .
#
# Nota: se usa Debian slim (no Alpine) porque sass-embedded, rollup y esbuild
# traen binarios nativos; glibc evita variantes musl frágiles.
# ─────────────────────────────────────────────────────────────

FROM node:22-bookworm-slim AS base

ENV CI=1 \
    COREPACK_ENABLE_DOWNLOAD_PROMPT=0 \
    PNPM_HOME=/pnpm \
    PATH=/pnpm:$PATH \
    npm_config_store_dir=/pnpm/store

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@10.30.3 --activate

# ── Dependencias (capa cacheable) ─────────────────────────────
# Se copian todos los manifiestos del workspace para que `--frozen-lockfile`
# valide igual que en local, pero solo se instala el subárbol del web.
FROM base AS deps
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/web/package.json apps/web/package.json
COPY packages/shared/package.json packages/shared/package.json
COPY apps/mobile/package.json apps/mobile/package.json
RUN --mount=type=cache,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter dev-test-frontend...

# ── Desarrollo (Vite + HMR) ──────────────────────────────────
FROM deps AS dev
COPY . .
EXPOSE 3000
CMD ["pnpm", "--filter", "dev-test-frontend", "dev"]

# ── Build estático ───────────────────────────────────────────
FROM deps AS build

# Vacío = same-origin: la SPA llama a /api y nginx lo proxya a api:8080.
# Si el API vive en otro dominio, pasarlo como build arg.
ARG VITE_API_BASE_URL=""
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

COPY packages/shared ./packages/shared
COPY apps/web ./apps/web
RUN pnpm --filter dev-test-frontend build

# ── Runtime (nginx non-root) ─────────────────────────────────
FROM nginxinc/nginx-unprivileged:1.29-alpine-slim AS runtime

COPY docker/web/nginx.conf /etc/nginx/conf.d/default.conf
COPY apps/web/deploy/nginx-security.conf /etc/nginx/snippets/security-headers.conf
COPY --from=build /app/apps/web/dist /usr/share/nginx/html

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8080/healthz || exit 1
