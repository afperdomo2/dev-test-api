#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Carga .env sin pisar variables ya definidas en el entorno (paridad con godotenv)
load_env() {
  [[ -f .env ]] || return 0
  while IFS='=' read -r key value; do
    [[ "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]] || continue
    if [[ -z "${!key+x}" ]]; then
      export "$key"="$value"
    fi
  done < .env
}
load_env

if [[ -z "${DB_NAME:-}" ]]; then
  echo "❌ DB_NAME vacio. Verifica .env o define DB_* en el entorno." >&2
  exit 1
fi

if [[ "${1:-}" == "create" ]] && [[ -z "${2:-}" ]]; then
  echo "❌ Uso: make migrate-create name=<descripcion>" >&2
  exit 1
fi

MIGRATE_DIR="internal/database/migrations"

# DSN en formato key-value (igual que config.go/.env.example)
# shellcheck disable=SC2154
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=${DB_HOST:-localhost} port=${DB_PORT:-5432} user=${DB_USER:-postgres} password=${DB_PASSWORD:-} dbname=${DB_NAME} sslmode=${DB_SSL_MODE:-disable}"

go tool goose -dir "$MIGRATE_DIR" "$@"