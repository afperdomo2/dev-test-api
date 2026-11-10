#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Sube la versión de la app mobile y la sincroniza con el proyecto nativo.
#
#   scripts/bump-mobile-version.sh            # versionCode +1, mismo version
#   scripts/bump-mobile-version.sh 1.1.0      # versionCode +1 + version nueva
#
# Fuente única: apps/mobile/app.json (`expo.version` + `expo.android.versionCode`).
# Refleja `versionCode`/`versionName` en apps/mobile/android/app/build.gradle, que
# es generado por prebuild y NO lee app.json. Detalle: docs/mobile-apk.md → "Versionado".

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_JSON="$ROOT_DIR/apps/mobile/app.json"
BUILD_GRADLE="$ROOT_DIR/apps/mobile/android/app/build.gradle"

fail() {
  echo "❌ $*" >&2
  exit 1
}

[[ -f "$APP_JSON" ]] || fail "No existe $APP_JSON"
command -v node >/dev/null 2>&1 || fail "node no está disponible en el PATH"

read_json() {
  node -e "const fs=require('fs');const j=JSON.parse(fs.readFileSync(process.argv[1],'utf8'));console.log($1)" "$APP_JSON"
}

APP_NAME="$(read_json "(j.expo&&j.expo.name)||''")"
CURRENT_VERSION="$(read_json "(j.expo&&j.expo.version)||''")"
CURRENT_CODE="$(read_json "(j.expo&&j.expo.android&&j.expo.android.versionCode)||0")"

[[ -n "$CURRENT_VERSION" ]] || fail "app.json no tiene expo.version"
[[ "$CURRENT_CODE" =~ ^[0-9]+$ ]] || fail "app.json → expo.android.versionCode no es numérico"

NEW_VERSION="${1:-$CURRENT_VERSION}"

[[ "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] \
  || fail "VERSION debe ser semver x.y.z (recibido: '$NEW_VERSION')"

LOWEST="$(printf '%s\n%s\n' "$CURRENT_VERSION" "$NEW_VERSION" | sort -V | head -n1)"
[[ "$LOWEST" == "$CURRENT_VERSION" ]] \
  || fail "La versión no puede bajar ($CURRENT_VERSION → $NEW_VERSION)"

NEW_CODE=$((CURRENT_CODE + 1))

# Reemplazo con archivo temporal: portable entre GNU sed y BSD sed (sin -i).
patch_file() {
  local file="$1"
  shift
  local tmp
  tmp="$(mktemp)"
  sed -E "$@" "$file" > "$tmp"
  mv "$tmp" "$file"
}

patch_file "$APP_JSON" \
  -e "s/(\"version\":[[:space:]]*\")[^\"]*(\")/\1${NEW_VERSION}\2/" \
  -e "s/(\"versionCode\":[[:space:]]*)[0-9]+/\1${NEW_CODE}/"

grep -Eq "\"version\":[[:space:]]*\"${NEW_VERSION}\"" "$APP_JSON" \
  || fail "No se pudo actualizar expo.version en $APP_JSON"
grep -Eq "\"versionCode\":[[:space:]]*${NEW_CODE}" "$APP_JSON" \
  || fail "No se pudo actualizar expo.android.versionCode en $APP_JSON"

if [[ -f "$BUILD_GRADLE" ]]; then
  patch_file "$BUILD_GRADLE" \
    -e "s/(versionCode[[:space:]]+)[0-9]+/\1${NEW_CODE}/" \
    -e "s/(versionName[[:space:]]+\")[^\"]*(\")/\1${NEW_VERSION}\2/"

  grep -Eq "versionCode[[:space:]]+${NEW_CODE}" "$BUILD_GRADLE" \
    || fail "No se pudo actualizar versionCode en $BUILD_GRADLE"
  grep -Eq "versionName[[:space:]]+\"${NEW_VERSION}\"" "$BUILD_GRADLE" \
    || fail "No se pudo actualizar versionName en $BUILD_GRADLE"
else
  echo "⚠️  No existe apps/mobile/android/: solo se actualizó app.json"
  echo "    (el próximo 'expo prebuild' genera el proyecto nativo con esos valores)"
fi

echo "✅ ${APP_NAME:-mobile} ${NEW_VERSION} (versionCode ${NEW_CODE})"
echo "   Siguiente: make mobile-apk-local"
