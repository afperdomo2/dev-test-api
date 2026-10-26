# syntax=docker/dockerfile:1.7

# ─────────────────────────────────────────────────────────────
# dev-test-api — imagen del backend (Go + Gin)
#
# Stages:
#   dev     → Go + Air. El module cache vive en volúmenes (no en la imagen)
#   build   → compila el binario estático (cache mounts de BuildKit)
#   runtime → distroless nonroot (imagen final, ~60–65 MB)
#
# Build manual:
#   docker build -f docker/backend/Dockerfile -t devtest-api .
# ─────────────────────────────────────────────────────────────

FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS base

# BuildKit los define automáticamente; se declaran para que los stages
# descendientes puedan usarlos al compilar multi-arquitectura.
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOTOOLCHAIN=local \
    GOFLAGS=-mod=readonly

WORKDIR /src

# ── Desarrollo con Air ───────────────────────────────────────
# Sin módulos ni caché dentro de la imagen: docker/compose.dev.yaml monta
# `go-mod-cache` y `go-cache`. El primer arranque descarga+compila; los
# siguientes son de segundos.
FROM base AS dev
COPY . .
EXPOSE 8080
CMD ["sh", "-c", "go mod download && go mod download github.com/air-verse/air && exec go tool air -c docker/backend/air.toml"]

# ── Build de producción ──────────────────────────────────────
FROM base AS build
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -buildid=" -o /out/api .

# ── Runtime (distroless) ─────────────────────────────────────
# Distroless static ya trae ca-certificates, tzdata y el usuario nonroot.
# Busybox aporta únicamente wget/sh para el HEALTHCHECK (distroless no tiene shell).
FROM busybox:1.37.0-musl AS busybox

FROM gcr.io/distroless/static-debian13:nonroot AS runtime

WORKDIR /app

COPY --from=busybox /bin/busybox /busybox
COPY --from=build /out/api /app/api

ENV GIN_MODE=release \
    PORT=8080 \
    TZ=UTC

EXPOSE 8080
USER 65532:65532

HEALTHCHECK --interval=30s --timeout=3s --start-period=25s --retries=3 \
    CMD ["/busybox", "wget", "-q", "-O", "/dev/null", "http://127.0.0.1:8080/health"]

ENTRYPOINT ["/app/api"]
