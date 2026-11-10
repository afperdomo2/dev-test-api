#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Copia el APK local compilado a una carpeta de fácil acceso (Descargas).
#
#   scripts/copy-mobile-apk.sh [ruta-apk] [destino]
#
# Sin argumentos usa el APK de release de Gradle
# (apps/mobile/android/app/build/outputs/apk/release/app-release.apk) y el
# directorio de descargas del usuario. El nombre incluye app, `version` y
# `versionCode` para no pisar builds anteriores:
#
#   ~/Descargas/DevTest-1.1.4-11.apk
#
# NO mueve el original: `adb install -r` sigue apuntando a la ruta de build.
# Detalle: docs/mobile-apk.md → "Copiar el APK a Descargas".

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_JSON="$ROOT_DIR/apps/mobile/app.json"
APK_DEFAULT="$ROOT_DIR/apps/mobile/android/app/build/outputs/apk/release/app-release.apk"

fail() {
  echo "❌ $*" >&2
  exit 1
}

command -v node >/dev/null 2>&1 || fail "node no está disponible en el PATH"
[[ -f "$APP_JSON" ]] || fail "No existe $APP_JSON"

APK_PATH="${1:-$APK_DEFAULT}"
[[ -f "$APK_PATH" ]] \
  || fail "No existe el APK: $APK_PATH (¿compilaste con 'make mobile-apk-local'?)"

# Directorio de descargas: XDG → ~/Downloads → ~/Descargas → ~
resolve_downloads() {
  local dir candidate
  if command -v xdg-user-dir >/dev/null 2>&1; then
    dir="$(xdg-user-dir DOWNLOAD 2>/dev/null || true)"
    if [[ -n "$dir" && "$dir" != "$HOME" ]]; then
      printf '%s\n' "$dir"
      return
    fi
  fi
  for candidate in "$HOME/Downloads" "$HOME/Descargas"; do
    if [[ -d "$candidate" ]]; then
      printf '%s\n' "$candidate"
      return
    fi
  done
  printf '%s\n' "$HOME"
}

DEST_DIR="${2:-${APK_DEST_DIR:-}}"
[[ -n "$DEST_DIR" ]] || DEST_DIR="$(resolve_downloads)"

read_json() {
  node -e "const fs=require('fs');const j=JSON.parse(fs.readFileSync(process.argv[1],'utf8'));console.log($1)" \
    "$APP_JSON"
}

APP_NAME="$(read_json "(j.expo&&j.expo.name)||'DevTest'")"
VERSION="$(read_json "(j.expo&&j.expo.version)||'0.0.0'")"
VERSION_CODE="$(read_json "(j.expo&&j.expo.android&&j.expo.android.versionCode)||0")"

SAFE_NAME="${APP_NAME// /-}"
DEST="$DEST_DIR/${SAFE_NAME}-${VERSION}-${VERSION_CODE}.apk"

mkdir -p "$DEST_DIR"
cp -f "$APK_PATH" "$DEST"

SIZE="$(du -h "$DEST" | cut -f1)"
echo "✅ APK copiado: $DEST"
echo "   tamaño:    $SIZE"
echo "   original:  $APK_PATH"

# Aviso temprano: el bundle horneado no debería conservar placeholders ni dev.
if command -v unzip >/dev/null 2>&1; then
  if unzip -p "$DEST" assets/index.android.bundle 2>/dev/null \
    | grep -aqE 'https?://(app\.example\.com|localhost|10\.0\.2\.2|127\.0\.0\.1)'; then
    echo "⚠️  El bundle contiene localhost o el placeholder app.example.com: recompilá con API_URL real."
  fi
fi
