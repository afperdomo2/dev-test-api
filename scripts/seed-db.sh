#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# ── Dependency check (bash-defensive-patterns) ─────────────────────
check_dependencies() {
  local -a missing=()
  for cmd in docker; do
    if ! command -v "$cmd" &>/dev/null; then
      missing+=("$cmd")
    fi
  done
  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "❌ Faltan dependencias: ${missing[*]}" >&2
    return 1
  fi
}
check_dependencies

# Carga .env si existe (ignora lineas vacias y comentarios)
if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

# Configuracion de la BD (valores por defecto coinciden con .env.example)
DB_CONTAINER="${DB_CONTAINER:-dev-postgres}"
DB_USER="${DB_USER:-devuser}"
DB_NAME="${DB_NAME:-dev_test_api}"
SEED_FILE="${1:-db/seeds/seed_questions.sql}"

if [[ ! -f "$SEED_FILE" ]]; then
  echo "❌ No se encontro el archivo seed: $SEED_FILE" >&2
  exit 1
fi

if ! docker ps --format '{{.Names}}' | grep -qx "$DB_CONTAINER"; then
  echo "❌ Contenedor no encontrado o no esta corriendo: $DB_CONTAINER" >&2
  echo "   Contenedores activos:" >&2
  docker ps --format "  - {{.Names}} ({{.Image}})"
  echo "   Ajusta DB_CONTAINER en .env o pasa el nombre como variable de entorno." >&2
  exit 1
fi

echo "🌱 Ejecutando seed: $SEED_FILE -> $DB_CONTAINER ($DB_USER/$DB_NAME)..."

# Pasa el SQL por stdin al psql del contenedor (no requiere psql en el host)
# Usa input redirection + pipefail para detectar fallos del cat o del psql
if docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -q < "$SEED_FILE"; then
  echo "✅ Seed ejecutado correctamente"
else
  echo "❌ Error al ejecutar el seed" >&2
  exit 1
fi

echo "📊 Verificacion rapida:"
docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT 'questions' AS tabla, count(*) FROM questions UNION ALL SELECT 'question_options', count(*) FROM question_options UNION ALL SELECT 'question_topics', count(*) FROM question_topics UNION ALL SELECT 'code_challenges', count(*) FROM code_challenges ORDER BY tabla;"
