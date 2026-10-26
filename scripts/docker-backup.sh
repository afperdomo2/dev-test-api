#!/usr/bin/env bash
set -Eeuo pipefail

trap 'cleanup; echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Backup lógico de la DB del stack prod-local (o de cualquier compose con un
# servicio `db`). Genera un .sql.gz con DROP ... IF EXISTS para poder restaurar
# sobre una DB existente.
#
# Uso:
#   bash scripts/docker-backup.sh
#   bash scripts/docker-backup.sh docker/backups/mi-backup.sql.gz
#   COMPOSE_FILE=docker/compose.dokploy.yaml bash scripts/docker-backup.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILE="${COMPOSE_FILE:-docker/compose.prod.yaml}"
COMPOSE_ENV_FILE="${COMPOSE_ENV_FILE:-docker/.env.prod}"
BACKUP_DIR="${BACKUP_DIR:-docker/backups}"
STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="${1:-$BACKUP_DIR/devtest-$STAMP.sql.gz}"
TMP=""

cleanup() {
  if [[ -n "${TMP:-}" && -f "$TMP" ]]; then
    rm -f "$TMP"
  fi
}

command -v docker >/dev/null 2>&1 || { echo "❌ docker no está disponible" >&2; exit 1; }
[[ -f "$COMPOSE_FILE" ]] || { echo "❌ No existe $COMPOSE_FILE" >&2; exit 1; }
[[ -f "$COMPOSE_ENV_FILE" ]] || { echo "❌ No existe $COMPOSE_ENV_FILE (¿corriste make docker-env-prod?)" >&2; exit 1; }

mkdir -p "$(dirname "$OUT")"
TMP="$(mktemp "${OUT}.partial.XXXXXX")"

echo "🗄️  Backup de la DB → $OUT"
docker compose -f "$COMPOSE_FILE" --env-file "$COMPOSE_ENV_FILE" exec -T db sh -c \
  'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" --clean --if-exists' | gzip -9 > "$TMP"

mv "$TMP" "$OUT"
TMP=""

echo "✅ Backup listo: $OUT ($(du -h "$OUT" | cut -f1))"
