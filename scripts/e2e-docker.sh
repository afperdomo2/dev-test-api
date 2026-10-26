#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# E2E contra el stack prod-local ya levantado (no arranca Vite).
# Requiere una DB limpia: la suite hace setup del primer admin.
#
# Uso: bash scripts/e2e-docker.sh
#      E2E_BASE_URL=http://localhost:9000 bash scripts/e2e-docker.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

BASE_URL="${1:-${E2E_BASE_URL:-http://localhost:8088}}"
BASE_URL="${BASE_URL%/}"

echo "🧪 E2E contra $BASE_URL"
PLAYWRIGHT_BASE_URL="$BASE_URL" pnpm --filter dev-test-frontend test:e2e
