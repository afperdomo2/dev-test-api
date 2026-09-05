#!/usr/bin/env bash
set -Eeuo pipefail

# check-fmt.sh — verifica gofmt con patrones defensivos (bash-defensive-patterns)
# Uso: bash scripts/check-fmt.sh

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

resolve_gofmt() {
  # Windows/WSL: PATH con espacios (Program Files) rompe command -v en bash; probar rutas comunes
  if command -v gofmt &>/dev/null; then echo "gofmt"; return 0; fi
  if command -v gofmt.exe &>/dev/null; then echo "gofmt.exe"; return 0; fi
  for p in "/mnt/c/Program Files/Go/bin/gofmt.exe" "C:/Program Files/Go/bin/gofmt.exe"; do
    if [[ -x "$p" ]]; then echo "$p"; return 0; fi
  done
  if command -v go &>/dev/null || command -v go.exe &>/dev/null; then echo "gofmt"; return 0; fi
  echo "gofmt"
}

check_dependencies() {
  local gofmt_bin
  gofmt_bin="$(resolve_gofmt)"
  if ! "$gofmt_bin" --help &>/dev/null && ! go version &>/dev/null && ! go.exe version &>/dev/null; then
    echo "❌ gofmt/go no encontrado — instala Go o verifica PATH" >&2
    return 1
  fi
}

check_dependencies

# gofmt -l lista archivos desformateados (uno por línea)
GOFMT_BIN="$(resolve_gofmt)"
unformatted="$("$GOFMT_BIN" -l . 2>&1)" || {
  echo "❌ gofmt falló (bin: $GOFMT_BIN)" >&2
  echo "$unformatted" >&2
  exit 1
}

if [[ -n "${unformatted}" ]]; then
  echo "❌ Unformatted files:" >&2
  echo "${unformatted}" >&2
  echo "   Ejecuta: make fmt (go fmt ./...)" >&2
  exit 1
fi

echo "✅ Go fmt OK"
