#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "❌ Error en ${BASH_SOURCE[0]} línea $LINENO" >&2' ERR

# Valida los tres compose y verifica que cada .example esté SINCRONIZADO con su
# compose:
#   1. toda variable interpolada en el compose existe en el .example;
#   2. ninguna clave del .example sobra respecto al compose;
#   3. las requeridas (`${VAR:?}`) están VACÍAS (las completa el humano);
#   4. las que tienen default (`${VAR:-x}`) valen exactamente `x` si están
#      escritas;
#   5. no queda ningún placeholder `CHANGE_ME`.
#
# Las requeridas se validan con dummies porque `docker compose config` no
# interpola `${VAR:?}` sin valor. Los dummies NO son secretos ni se escriben en
# ningún archivo.
#
# Uso: bash scripts/docker-config.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

STACKS=(
  "dev:docker/compose.dev.yaml:docker/env/.env.dev.example"
  "prod-local:docker/compose.prod.yaml:docker/env/.env.prod.example"
  "dokploy:docker/compose.dokploy.yaml:docker/env/.env.dokploy.example"
)

DUMMY_ENV=(
  "DB_PASSWORD=dummy-validation-only"
  "JWT_SECRET=dummy-validation-only"
  "APP_ORIGIN=https://dummy.invalid"
)

# AWK_REFERENCIAS imprime "VAR<TAB>tipo<TAB>default" por cada ${VAR...} del
# compose. Ignora las referencias escapadas al contenedor ($${VAR}).
read -r -d '' AWK_REFS <<'AWK' || true
{
  line = $0
  gsub(/\$\$\{[^}]*\}/, "", line)

  while ((start = index(line, "${")) > 0) {
    line = substr(line, start + 2)
    close_at = index(line, "}")
    if (close_at == 0) break
    body = substr(line, 1, close_at - 1)
    line = substr(line, close_at + 1)

    op_at = match(body, /:-|:\?/)
    if (op_at > 0) {
      print substr(body, 1, op_at - 1) "\t" (substr(body, op_at, 2) == ":-" ? "default" : "required") "\t" substr(body, op_at + 2)
    } else {
      print body "\tsin-default\t"
    }
  }
}
AWK

errors=0
error() {
  echo "   ❌ $*" >&2
  errors=$((errors + 1))
}

check_stack() {
  local name="$1" compose="$2" envfile="$3"
  local line key val ref

  echo "🔎 $name: $compose ↔ $envfile"

  # 1) El compose interpola bien (requeridas con dummies).
  env "${DUMMY_ENV[@]}" docker compose -f "$compose" --env-file "$envfile" config -q ||
    error "$compose no interpola"

  # 2) Referencias del compose.
  declare -A REF_KIND=() REF_DEFAULT=()
  while IFS=$'\t' read -r key ref val; do
    [[ -n "$key" ]] || continue
    if [[ -z "${REF_KIND[$key]:-}" ]]; then
      REF_KIND[$key]="$ref"
      REF_DEFAULT[$key]="$val"
    elif [[ "${REF_KIND[$key]}" != "$ref" || "${REF_DEFAULT[$key]:-}" != "$val" ]]; then
      error "$compose interpola '$key' de dos formas distintas"
    fi
  done < <(awk "$AWK_REFS" "$compose")

  # 3) Claves del .example.
  declare -A ENV_VAL=() ENV_OVERRIDE=()
  while IFS= read -r line || [[ -n "$line" ]]; do
    line="${line%$'\r'}"
    [[ "$line" == *=* ]] || continue
    [[ "$line" =~ ^[[:space:]]*# ]] && continue
    key="${line%%=*}"
    key="${key//[[:space:]]/}"
    val="${line#*=}"
    # Un valor distinto del default solo se acepta con el marcador `# override`.
    if [[ "$val" == *" #"* ]]; then
      [[ "${val#* #}" == *override* ]] && ENV_OVERRIDE[$key]=1
      val="${val%% #*}"
    fi
    ENV_VAL[$key]="$val"
  done < "$envfile"

  # 4) Variables del compose que faltan en el .example.
  for key in "${!REF_KIND[@]}"; do
    if [[ -z "${ENV_VAL[$key]+set}" ]]; then
      error "$envfile no documenta '$key' (lo usa $compose)"
    fi
  done

  # 5) Claves del .example que el compose no usa.
  for key in "${!ENV_VAL[@]}"; do
    if [[ -z "${REF_KIND[$key]:-}" ]]; then
      error "$envfile define '$key', que $compose no interpola"
    fi
  done

  # 6) Valores: requeridas vacías, defaults exactos, sin placeholders.
  for key in "${!REF_KIND[@]}"; do
    [[ -n "${ENV_VAL[$key]+set}" ]] || continue
    val="${ENV_VAL[$key]}"
    case "${REF_KIND[$key]}" in
      required | sin-default)
        [[ -z "$val" ]] ||
          error "$key debe quedar VACÍA en $envfile (no tiene default; el humano pone el valor)"
        ;;
      default)
        if [[ -z "${ENV_OVERRIDE[$key]:-}" && -n "$val" && "$val" != "${REF_DEFAULT[$key]}" ]]; then
          error "$key=${val} no es el default del compose ('${REF_DEFAULT[$key]}'); vaciá el valor o marcá la línea con '# override'"
        fi
        ;;
    esac
  done

  if grep -q "CHANGE_ME" "$envfile"; then
    error "$envfile todavía tiene placeholders CHANGE_ME; las requeridas van vacías"
  fi
}

for stack in "${STACKS[@]}"; do
  IFS=':' read -r name compose envfile <<<"$stack"
  check_stack "$name" "$compose" "$envfile"
done

if ((errors > 0)); then
  echo "❌ $errors problema(s) de configuración" >&2
  exit 1
fi

echo "✅ Compose válidos y .example sincronizados con su compose"
