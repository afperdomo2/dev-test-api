#!/usr/bin/env bash
set -Eeuo pipefail

# go-vet-affected.sh — ejecuta go vet solo en los paquetes afectados por archivos staged
# Uso: bash scripts/go-vet-affected.sh <file1.go> [file2.go ...]

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

resolve_go() {
  if command -v go &>/dev/null; then echo "go"; return 0; fi
  if command -v go.exe &>/dev/null; then echo "go.exe"; return 0; fi
  for p in "/mnt/c/Program Files/Go/bin/go.exe" "C:/Program Files/Go/bin/go.exe" "$(cygpath -u "$LOCALAPPDATA/Programs/Go/bin/go.exe" 2>/dev/null || true)"; do
    if [[ -x "$p" ]]; then echo "$p"; return 0; fi
  done
  echo "go"
}

GO_BIN="$(resolve_go)"

files=()
for f in "$@"; do
  if [[ "$f" == *.go && -f "$f" ]]; then
    files+=("$(dirname -- "$f")")
  fi
done

if [[ ${#files[@]} -eq 0 ]]; then
  echo "✅ No Go files staged — skipping vet"
  exit 0
fi

mapfile -t dirs < <(printf '%s\n' "${files[@]}" | sort -u)

packages=()
for dir in "${dirs[@]}"; do
  if [[ "$dir" == "." ]]; then
    packages+=(".")
  else
    packages+=("./$dir")
  fi
done

echo "🔍 Vetting packages: ${packages[*]}"
"$GO_BIN" vet "${packages[@]}"