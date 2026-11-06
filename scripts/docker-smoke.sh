#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Smoke test del stack prod-local (o de cualquier stack con un solo origen).
#
# Uso:
#   bash scripts/docker-smoke.sh
#   BASE_URL=http://localhost:8088 bash scripts/docker-smoke.sh
#
# Verifica:
#   - nginx vivo (/healthz)
#   - API proxeado (/health)
#   - contrato de auth (/api/v1/auth/status)
#   - SPA servida (index.html)
#   - cabeceras de seguridad de la SPA

BASE_URL="${BASE_URL:-http://localhost:8088}"
BASE_URL="${BASE_URL%/}"

pass() { echo "✅ $*"; }
fail() { echo "❌ $*" >&2; exit 1; }

command -v curl >/dev/null 2>&1 || fail "curl no está disponible en el PATH"

echo "🔎 Smoke test contra $BASE_URL"

curl -fsS --max-time 10 "$BASE_URL/healthz" >/dev/null \
  || fail "nginx no responde /healthz"
pass "nginx /healthz"

curl -fsS --max-time 10 "$BASE_URL/health" | grep -q '"status":"ok"' \
  || fail "el API no responde /health a través de nginx"
pass "API /health"

curl -fsS --max-time 10 "$BASE_URL/api/v1/auth/status" | grep -q 'setupRequiresToken' \
  || fail "el proxy no responde /api/v1/auth/status"
pass "API /api/v1/auth/status"

curl -fsS --max-time 10 "$BASE_URL/" | grep -q 'id="app"' \
  || fail "nginx no sirve la SPA (index.html)"
pass "SPA index.html"

HEADERS="$(curl -fsSI --max-time 10 "$BASE_URL/")"
for header in content-security-policy x-content-type-options x-frame-options referrer-policy permissions-policy; do
  echo "$HEADERS" | grep -qi "^${header}:" \
    || fail "falta la cabecera $header en la SPA"
done
pass "cabeceras de seguridad de la SPA"

echo
echo "✅ Smoke test OK — $BASE_URL"
