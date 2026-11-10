#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Crea el archivo de entorno del stack Docker sin sobrescribir uno existente y
# completa con secretos aleatorios las variables OBLIGATORIAS que el .example
# deja vacías (DB_PASSWORD, JWT_SECRET, SETUP_TOKEN).
#
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

# openssl es lo habitual; /dev/urandom cubre Git Bash y contenedores mínimos.
rand_hex() {
  local bytes="${1:-32}"
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex "$bytes"
  elif [[ -r /dev/urandom ]]; then
    od -An -tx1 -N "$bytes" /dev/urandom | tr -d ' \n'
  else
    echo "❌ No hay openssl ni /dev/urandom para generar secretos" >&2
    return 1
  fi
}

# Reemplaza la línea `CLAVE=` (vacía) por `CLAVE=<aleatorio>`.
fill_secret() {
  local key="$1" bytes="$2" value tmp
  value="$(rand_hex "$bytes")"

  if ! grep -qx "${key}=" "$DEST"; then
    echo "❌ $DEST no tiene la línea vacía '${key}=' (¿cambió el .example?)" >&2
    return 1
  fi

  tmp="$(mktemp)"
  awk -v k="$key" -v v="$value" '$0 == k "=" { print k "=" v; next } { print }' \
    "$DEST" >"$tmp"
  mv "$tmp" "$DEST"
}

read_env() {
  local key="$1"
  sed -n "s/^${key}=//p" "$DEST" | head -n 1
}

cp "$SRC" "$DEST"
# Si falla la generación, no dejar un entorno a medias (quedaría "ya existe").
trap 'echo "❌ Error generando secretos: se descarta $DEST" >&2; rm -f "$DEST"' ERR

fill_secret DB_PASSWORD 16
fill_secret JWT_SECRET 32
fill_secret SETUP_TOKEN 16

# Secretos completos: volver al manejo de errores normal.
trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR
echo "✅ $DEST creado desde $SRC"

if [[ "$KIND" == "prod" ]]; then
  echo "🔐 Secretos generados (no los reutilices en otro entorno):"
else
  echo "🔐 Secretos generados para desarrollo local:"
fi
echo "   DB_USER=$(read_env DB_USER)"
echo "   DB_PASSWORD=$(read_env DB_PASSWORD)"
echo "   SETUP_TOKEN=$(read_env SETUP_TOKEN)"

if [[ "$KIND" == "prod" ]]; then
  echo "⚠️  Usá el SETUP_TOKEN de arriba para crear el primer admin; después vacialo en $DEST."
fi
