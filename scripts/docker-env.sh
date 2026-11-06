#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Crea el archivo de entorno del stack Docker sin sobrescribir uno existente.
# Uso: bash scripts/docker-env.sh [dev|prod]

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

KIND="${1:-dev}"

case "$KIND" in
  dev)
    SRC="docker/env/.env.dev.example"
    DEST="docker/.env"
    ;;
  prod)
    SRC="docker/env/.env.prod.example"
    DEST="docker/.env.prod"
    ;;
  *)
    echo "❌ Uso: bash scripts/docker-env.sh [dev|prod]" >&2
    exit 1
    ;;
esac

[[ -f "$SRC" ]] || { echo "❌ No existe $SRC" >&2; exit 1; }

if [[ -f "$DEST" ]]; then
  echo "ℹ️  $DEST ya existe (no se sobrescribe)"
  exit 0
fi

cp "$SRC" "$DEST"
echo "✅ $DEST creado desde $SRC"
if [[ "$KIND" == "prod" ]]; then
  echo "⚠️  Editá los valores CHANGE_ME antes de levantar prod-local."
fi
