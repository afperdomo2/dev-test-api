#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Genera los iconos de web y mobile a partir de los másters de `brand/`.
#
#   brand/icon.svg  → tile redondeado (favicon y PWA "any")
#   brand/mark.svg  → monograma suelto (adaptive icon, maskable, apple-touch)
#
# Requiere Google Chrome/Chromium (render de SVG) y Python 3 + Pillow (ICO y
# normalización de PNG). Los resultados se commitean: este script es de uso
# manual (`make icons`), no corre en CI ni en los git hooks.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BRAND_DIR="$ROOT_DIR/brand"
WEB_PUBLIC="$ROOT_DIR/apps/web/public"
MOBILE_ASSETS="$ROOT_DIR/apps/mobile/assets"
BRAND_COLOR="#0F766E"

die() {
  echo "❌ $*" >&2
  exit 1
}

find_chrome() {
  local bin
  for bin in \
    google-chrome \
    google-chrome-stable \
    chromium \
    chromium-browser \
    "/c/Program Files/Google/Chrome/Application/chrome.exe" \
    "/c/Program Files (x86)/Google/Chrome/Application/chrome.exe"; do
    if command -v "$bin" >/dev/null 2>&1; then
      command -v "$bin"
      return 0
    fi
  done
  return 1
}

find_python() {
  local bin
  for bin in python3 python; do
    if command -v "$bin" >/dev/null 2>&1; then
      command -v "$bin"
      return 0
    fi
  done
  return 1
}

[[ -f "$BRAND_DIR/icon.svg" ]] || die "No existe $BRAND_DIR/icon.svg"
[[ -f "$BRAND_DIR/mark.svg" ]] || die "No existe $BRAND_DIR/mark.svg"

CHROME_BIN="$(find_chrome)" \
  || die "No encontré Chrome/Chromium. Instalalo o generá los PNG con rsvg-convert/ImageMagick."

PYTHON_BIN="$(find_python)" || die "Falta python3/python."
"$PYTHON_BIN" -c 'import PIL' 2>/dev/null || die "Falta Pillow: python3 -m pip install pillow"

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

ICON_SVG="$(cat "$BRAND_DIR/icon.svg")"
MARK_SVG="$(cat "$BRAND_DIR/mark.svg")"

# Renderiza un SVG a PNG cuadrado con Chrome headless.
#
#   render <modo> <px> <salida>
#
# Modos:
#   tile       tile redondeado, esquinas transparentes
#   full       monograma sobre fondo de marca (opaco, full-bleed)
#   full-safe  igual que full, pero con margen para iconos maskable
#   safe       monograma transparente con margen (adaptive icon Android)
#   mark       monograma transparente al tamaño natural
render() {
  local mode="$1"
  local size="$2"
  local out="$3"
  local svg bg scale

  case "$mode" in
    tile)
      svg="$ICON_SVG"
      bg="transparent"
      scale="1"
      ;;
    full)
      svg="$MARK_SVG"
      bg="$BRAND_COLOR"
      scale="1"
      ;;
    full-safe)
      svg="$MARK_SVG"
      bg="$BRAND_COLOR"
      # Zona segura maskable: círculo del 80% del ancho.
      scale="0.95"
      ;;
    safe)
      svg="$MARK_SVG"
      bg="transparent"
      # Zona segura Android: círculo de 66dp sobre un lienzo de 108dp.
      scale="0.77"
      ;;
    mark)
      svg="$MARK_SVG"
      bg="transparent"
      scale="1"
      ;;
    *)
      die "Modo de render desconocido: $mode"
      ;;
  esac

  local html="$TMP_DIR/render-${mode}-${size}.html"
  cat >"$html" <<HTML
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
html, body { margin: 0; padding: 0; width: 100%; height: 100%; overflow: hidden; background: ${bg}; }
.wrap { position: fixed; inset: 0; display: grid; place-items: center; }
svg { width: 100%; height: 100%; display: block; transform: scale(${scale}); transform-origin: center; }
</style>
</head>
<body>
<div class="wrap">${svg}</div>
</body>
</html>
HTML

  "$CHROME_BIN" \
    --headless \
    --disable-gpu \
    --disable-dev-shm-usage \
    --hide-scrollbars \
    --force-device-scale-factor=1 \
    --default-background-color=00000000 \
    --window-size="${size},${size}" \
    --screenshot="$out" \
    "file://$html" >/dev/null 2>&1 \
    || die "Chrome falló al renderizar $mode a ${size}px"

  [[ -s "$out" ]] || die "Chrome no generó la imagen $out"
}

echo "🎨 Renderizando iconos desde brand/…"

# Web
render tile 256 "$TMP_DIR/favicon-256.png"
render full 180 "$TMP_DIR/apple-touch-icon.png"
render full 192 "$TMP_DIR/icon-192.png"
render full 512 "$TMP_DIR/icon-512.png"
render full-safe 512 "$TMP_DIR/maskable-512.png"

# Mobile
render full 1024 "$TMP_DIR/icon.png"
render safe 1024 "$TMP_DIR/android-icon-foreground.png"
render safe 432 "$TMP_DIR/android-icon-monochrome.png"
render tile 48 "$TMP_DIR/mobile-favicon.png"
render mark 1024 "$TMP_DIR/splash-icon.png"

# Normaliza los PNG (quita metadatos de Chrome para que el script sea
# idempotente) y arma el .ico multi-tamaño desde el render de 256 px.
"$PYTHON_BIN" - "$TMP_DIR" <<'PY'
import sys
from pathlib import Path

from PIL import Image

tmp = Path(sys.argv[1])

for png in sorted(tmp.glob("*.png")):
    with Image.open(png) as img:
        img.load()
        # `copy()` desacopla el archivo antes de reescribirlo (Windows bloquea
        # el path mientras el handle de lectura sigue abierto).
        normalized = img.copy()

    normalized.save(png, format="PNG", optimize=True)

with Image.open(tmp / "favicon-256.png") as img:
    img.convert("RGBA").save(
        tmp / "favicon.ico",
        format="ICO",
        sizes=[(16, 16), (24, 24), (32, 32), (48, 48)],
    )
PY

mkdir -p "$WEB_PUBLIC" "$MOBILE_ASSETS"

# Web
cp "$BRAND_DIR/icon.svg" "$WEB_PUBLIC/favicon.svg"
cp "$TMP_DIR/favicon.ico" "$WEB_PUBLIC/favicon.ico"
cp "$TMP_DIR/apple-touch-icon.png" "$WEB_PUBLIC/apple-touch-icon.png"
cp "$TMP_DIR/icon-192.png" "$WEB_PUBLIC/icon-192.png"
cp "$TMP_DIR/icon-512.png" "$WEB_PUBLIC/icon-512.png"
cp "$TMP_DIR/maskable-512.png" "$WEB_PUBLIC/maskable-512.png"

# Mobile
cp "$TMP_DIR/icon.png" "$MOBILE_ASSETS/icon.png"
cp "$TMP_DIR/android-icon-foreground.png" "$MOBILE_ASSETS/android-icon-foreground.png"
cp "$TMP_DIR/android-icon-monochrome.png" "$MOBILE_ASSETS/android-icon-monochrome.png"
cp "$TMP_DIR/mobile-favicon.png" "$MOBILE_ASSETS/favicon.png"
cp "$TMP_DIR/splash-icon.png" "$MOBILE_ASSETS/splash-icon.png"

echo "✅ Web  → apps/web/public/{favicon.svg,favicon.ico,apple-touch-icon.png,icon-192.png,icon-512.png,maskable-512.png}"
echo "✅ Mobile → apps/mobile/assets/{icon.png,android-icon-foreground.png,android-icon-monochrome.png,favicon.png,splash-icon.png}"
echo "   Siguiente: make mobile-apk-local APP_NAME=\"DevTest\" (hace prebuild)"
