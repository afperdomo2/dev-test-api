#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Binario de gh. En Git Bash de Windows no siempre está en PATH; se puede
# apuntar explícitamente con:
#   GH_BIN="/c/Program Files/GitHub CLI/gh.exe" bash scripts/setup-branch-protection.sh
GH_BIN="${GH_BIN:-gh}"

if ! command -v "$GH_BIN" >/dev/null 2>&1; then
  for candidate in \
    "/c/Program Files/GitHub CLI/gh.exe" \
    "/c/Program Files (x86)/GitHub CLI/gh.exe"; do
    if [[ -x "$candidate" ]]; then
      GH_BIN="$candidate"
      break
    fi
  done
fi

# Aplica la protección de ramas acordada:
#   main    → PR + checks obligatorios + lineal + sin force-push/borrado + SIN bypass
#   develop → sin force-push ni borrado (sin PR obligatorio)
#
# Es idempotente: borra y recrea los rulesets por nombre.
# Uso: scripts/setup-branch-protection.sh
# Requiere: gh autenticado con permiso de administración del repo.

# Contextos de los checks obligatorios. Deben coincidir EXACTAMENTE con el
# nombre del job en el workflow (ver .github/workflows/*.yml).
MAIN_CHECKS=("backend" "shared" "web" "gitleaks")

check_dependencies() {
  local -a missing=()
  for cmd in "$GH_BIN"; do
    command -v "$cmd" &>/dev/null || missing+=("$cmd")
  done
  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "❌ Faltan dependencias: ${missing[*]}" >&2
    exit 1
  fi
}

check_dependencies

REPO="$("$GH_BIN" repo view --json nameWithOwner --jq '.nameWithOwner')"
echo "🔒 Configurando protección para $REPO"

delete_ruleset_if_exists() {
  local name="$1" id
  id="$("$GH_BIN" api "repos/$REPO/rulesets" --jq ".[] | select(.name == \"$name\") | .id" || true)"
  if [[ -n "$id" ]]; then
    "$GH_BIN" api --method DELETE "repos/$REPO/rulesets/$id" >/dev/null
    echo "🗑️  Ruleset '$name' anterior eliminado (id $id)"
  fi
}

create_main_ruleset() {
  local checks_json=""
  local check
  for check in "${MAIN_CHECKS[@]}"; do
    [[ -n "$checks_json" ]] && checks_json+=","
    checks_json+="{\"context\":\"$check\"}"
  done

  "$GH_BIN" api --method POST "repos/$REPO/rulesets" --input - >/dev/null <<JSON
{
  "name": "main-protection",
  "target": "branch",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "include": ["refs/heads/main"],
      "exclude": []
    }
  },
  "bypass_actors": [],
  "rules": [
    { "type": "deletion" },
    { "type": "non_fast_forward" },
    { "type": "required_linear_history" },
    {
      "type": "pull_request",
      "parameters": {
        "required_approving_review_count": 0,
        "dismiss_stale_reviews_on_push": false,
        "require_code_owner_review": false,
        "require_last_push_approval": false,
        "required_review_thread_resolution": false
      }
    },
    {
      "type": "required_status_checks",
      "parameters": {
        "strict_required_status_checks_policy": true,
        "required_status_checks": [ ${checks_json} ]
      }
    }
  ]
}
JSON
  echo "✅ main protegida: PR + ${#MAIN_CHECKS[@]} checks + lineal, sin bypass"
}

create_develop_ruleset() {
  "$GH_BIN" api --method POST "repos/$REPO/rulesets" --input - >/dev/null <<'JSON'
{
  "name": "develop-protection",
  "target": "branch",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "include": ["refs/heads/develop"],
      "exclude": []
    }
  },
  "bypass_actors": [],
  "rules": [
    { "type": "deletion" },
    { "type": "non_fast_forward" }
  ]
}
JSON
  echo "✅ develop protegida: sin force-push ni borrado"
}

delete_ruleset_if_exists "main-protection"
delete_ruleset_if_exists "develop-protection"
create_main_ruleset
create_develop_ruleset

echo
echo "📋 Rulesets activos:"
"$GH_BIN" api "repos/$REPO/rulesets" --jq '.[] | "  - \(.name) (\(.target), \(.enforcement))"'
