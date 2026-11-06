#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Restore de un backup generado por scripts/docker-backup.sh.
#
# Uso:
#   bash scripts/docker-restore.sh docker/backups/devtest-YYYYMMDD-HHMMSS.sql.gz
#   FORCE=1 bash scripts/docker-restore.sh <archivo>   # sin confirmación
#
# El dump incluye DROP ... IF EXISTS: reemplaza los objetos de la DB del stack.
# No borra el volumen ni la base: la DB debe existir.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILE="${COMPOSE_FILE:-docker/compose.prod.yaml}"
COMPOSE_ENV_FILE="${COMPOSE_ENV_FILE:-docker/.env.prod}"
FILE="${1:-}"

usage() {
  echo "Uso: bash scripts/docker-restore.sh <archivo.sql.gz>" >&2
  exit 1
}

[[ -n "$FILE" ]] || usage
[[ -f "$FILE" ]] || { echo "❌ No existe el archivo $FILE" >&2; exit 1; }
[[ -f "$COMPOSE_FILE" ]] || { echo "❌ No existe $COMPOSE_FILE" >&2; exit 1; }
[[ -f "$COMPOSE_ENV_FILE" ]] || { echo "❌ No existe $COMPOSE_ENV_FILE (¿corriste make docker-env-prod?)" >&2; exit 1; }

command -v docker >/dev/null 2>&1 || { echo "❌ docker no está disponible" >&2; exit 1; }

if [[ "${FORCE:-0}" != "1" ]]; then
  echo "⚠️  Esto reemplaza los objetos de la DB de $COMPOSE_FILE con $FILE."
  read -r -p "¿Continuar? [y/N] " answer
  [[ "$answer" =~ ^[Yy]$ ]] || { echo "Cancelado."; exit 0; }
fi

echo "♻️  Restaurando $FILE ..."
gunzip -c "$FILE" | docker compose -f "$COMPOSE_FILE" --env-file "$COMPOSE_ENV_FILE" exec -T db sh -c \
  'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -q'

echo "✅ Restore completado"
