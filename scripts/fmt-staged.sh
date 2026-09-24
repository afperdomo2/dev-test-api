#!/usr/bin/env bash
set -Eeuo pipefail

# fmt-staged.sh — formatea con gofmt solo los archivos .go staged
# Uso: bash scripts/fmt-staged.sh <file1.go> [file2.go ...]

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

GOFMT_BIN="$(resolve_gofmt)"

files=()
for f in "$@"; do
  if [[ "$f" == *.go && -f "$f" ]]; then
    files+=("$f")
  fi
done

if [[ ${#files[@]} -eq 0 ]]; then
  echo "✅ No Go files staged — skipping format"
  exit 0
fi

"$GOFMT_BIN" -w "${files[@]}"
echo "✅ gofmt applied to ${#files[@]} file(s)"