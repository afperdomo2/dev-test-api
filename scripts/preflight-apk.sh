#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Preflight del APK manual (EAS cloud).
#
# Valida lo que más se olvida antes de compilar:
#   1. La URL del API en apps/mobile/eas.json es HTTPS y no es el placeholder.
#   2. El backend responde /health.
#   3. El proxy del API responde /api/v1/auth/status.
#   4. app.json tiene android.package y versionCode.
#
# EAS cloud usa el `env` de eas.json, NO las variables exportadas en tu shell.
# Si cambiás de backend, editá eas.json y volvé a compilar: la URL queda
# compilada dentro del APK.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_JSON="$ROOT_DIR/apps/mobile/app.json"
EAS_JSON="$ROOT_DIR/apps/mobile/eas.json"

fail() {
  echo "❌ $*" >&2
  exit 1
}

pass() {
  echo "✅ $*"
}

command -v node >/dev/null 2>&1 || fail "node no está disponible en el PATH"
command -v curl >/dev/null 2>&1 || fail "curl no está disponible en el PATH"
[[ -f "$APP_JSON" ]] || fail "No existe $APP_JSON"
[[ -f "$EAS_JSON" ]] || fail "No existe $EAS_JSON (¿corriste eas init?)"

read_json() {
  local file="$1"
  local expr="$2"
  node -e "const fs=require('fs');const j=JSON.parse(fs.readFileSync(process.argv[1],'utf8'));console.log($expr)" "$file"
}

URL_EAS="$(read_json "$EAS_JSON" "(j.build&&j.build.preview&&j.build.preview.env&&j.build.preview.env.EXPO_PUBLIC_API_BASE_URL)||''")"
URL_EAS="${URL_EAS%/}"

if [[ -n "${EXPO_PUBLIC_API_BASE_URL:-}" && "${EXPO_PUBLIC_API_BASE_URL%/}" != "$URL_EAS" ]]; then
  fail "EXPO_PUBLIC_API_BASE_URL (${EXPO_PUBLIC_API_BASE_URL}) no coincide con eas.json ($URL_EAS). EAS cloud compila el valor de eas.json: editá ese archivo y volvé a intentar."
fi

URL="$URL_EAS"

[[ -n "$URL" ]] || fail "eas.json → build.preview.env.EXPO_PUBLIC_API_BASE_URL está vacío"

if [[ "$URL" == *example.com* ]]; then
  fail "La URL del APK sigue siendo el placeholder ($URL). Editá apps/mobile/eas.json → build.preview.env.EXPO_PUBLIC_API_BASE_URL con el dominio real."
fi

if [[ "$URL" != https://* ]]; then
  if [[ "${ALLOW_HTTP:-0}" != "1" ]]; then
    fail "La URL debe ser https:// (D-41). Para una prueba local explícita: ALLOW_HTTP=1 $0"
  fi
  echo "⚠️  ALLOW_HTTP=1: se permite $URL solo para pruebas locales"
fi

PACKAGE="$(read_json "$APP_JSON" "(j.expo&&j.expo.android&&j.expo.android.package)||''")"
VERSION_CODE="$(read_json "$APP_JSON" "(j.expo&&j.expo.android&&j.expo.android.versionCode)||0")"
VERSION_NAME="$(read_json "$APP_JSON" "(j.expo&&j.expo.version)||''")"

[[ -n "$PACKAGE" ]] || fail "app.json no tiene expo.android.package"
[[ "$VERSION_CODE" =~ ^[0-9]+$ ]] || fail "app.json → expo.android.versionCode no es numérico"
(( VERSION_CODE >= 1 )) || fail "app.json → expo.android.versionCode debe ser >= 1"

pass "eas.json → $URL"
pass "app.json → $PACKAGE v$VERSION_NAME (versionCode $VERSION_CODE)"

echo "🔎 Verificando backend..."
curl -fsS --max-time 15 "$URL/health" | grep -q '"status":"ok"' \
  || fail "El backend no responde /health en $URL"
pass "/health responde ok"

curl -fsS --max-time 15 "$URL/api/v1/auth/status" | grep -q 'setupRequiresToken' \
  || fail "El proxy no responde /api/v1/auth/status en $URL"
pass "/api/v1/auth/status responde"

echo
echo "📱 Siguiente paso:"
echo "   cd apps/mobile && npx --yes eas-cli@latest build -p android --profile preview"
echo "   (o: make mobile-apk)"
echo
echo "⚠️  Si cambiás el backend, actualizá eas.json y volvé a compilar: la URL"
echo "    queda compilada en el APK."
