#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Siembra las 50 preguntas de desarrollo en la DB del stack dev.
# Requiere el stack levantado (make docker-dev / make docker-dev-up).

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILE="docker/compose.dev.yaml"
COMPOSE_ENV_FILE="docker/.env"
SEED_FILE="db/seeds/seed_questions.sql"

[[ -f "$COMPOSE_ENV_FILE" ]] || { echo "❌ No existe $COMPOSE_ENV_FILE (¿corriste make docker-env?)" >&2; exit 1; }
[[ -f "$SEED_FILE" ]] || { echo "❌ No existe $SEED_FILE" >&2; exit 1; }

echo "🌱 Ejecutando $SEED_FILE ..."
docker compose -f "$COMPOSE_FILE" --env-file "$COMPOSE_ENV_FILE" exec -T db sh -c \
  'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -q' < "$SEED_FILE"

echo "✅ Seed completado"
