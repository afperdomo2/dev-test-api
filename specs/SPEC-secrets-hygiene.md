# SPEC — secrets-hygiene

| Campo | Valor |
|---|---|
| Module id | `secrets-hygiene` |
| Origen | Revisión de seguridad 2026-10-21 (M5 + §1) |
| Estado | `ready` |
| Depende de | — |
| Consumido por | `deps-supply-chain` (reutiliza el job `secrets`) |
| Migración goose | Ninguna |
| Dependencia nueva | Ninguna de runtime. `gitleaks` es un binario de CI/local, **no** entra en `go.mod` |

## Objective

El repositorio ya es público y la auditoría confirmó que **no hay secretos en la
historia** (175 commits, árbol remoto incluido). Este módulo es **preventivo**:
que un `git add` accidental o un commit futuro no pueda publicar una credencial,
y que si ocurre, se detecte antes de salir del portátil.

**Historias**

- Como mantenedor, quiero que `git status` nunca me ofrezca `.env.local`,
  `*.pem` o `*.key` para commitear.
- Como mantenedor, quiero que un push con una llave falle en CI antes de
  mergearse.
- Como visitante del repo, quiero un canal claro para reportar una
  vulnerabilidad sin abrir un issue público con el secreto dentro.

**Éxito**

- `git check-ignore` marca como ignoradas todas las variantes de `.env` salvo
  los `.env.example`.
- Un canario con un token falso hace fallar `gitleaks detect` y el job de CI.
- `SECURITY.md` existe y GitHub tiene Secret Scanning + Push Protection activos.

## Tech stack

Git + GitHub Actions + [gitleaks](https://github.com/gitleaks/gitleaks) v8
(binario). Sin dependencias de runtime.

## Commands

```bash
# verificación local (requiere gitleaks instalado)
gitleaks detect --source . --redact --no-banner

# canario: debe fallar (exit 1) y NO debe estar dentro del repo.
# La clave se genera en runtime para no dejar material tipo PEM en el repo.
# (El token de documentación `AKIA…EXAMPLE` NO dispara la regla por defecto
# de gitleaks 8.30.1, por eso no se usa como canario.)
canary="$(mktemp -d)"
node -e "const c=require('crypto');require('fs').writeFileSync(process.argv[1],c.generateKeyPairSync('rsa',{modulusLength:2048,privateKeyEncoding:{type:'pkcs1',format:'pem'},publicKeyEncoding:{type:'spki',format:'pem'}}).privateKey)" "$canary/id_rsa.txt"
gitleaks detect --no-git --source "$canary" --config .gitleaks.toml --redact --no-banner   # exit 1
rm -rf "$canary"

# comprobación de ignores
git check-ignore -v .env.local .env.production apps/mobile/.env.preview cert.pem server.key

# CI local equivalente
act -j secrets   # opcional
```

## Project structure

```
.gitignore                       # MODIFICADO: variantes de .env + material criptográfico
.gitleaks.toml                   # NUEVO: reglas + allowlist del repo
SECURITY.md                      # NUEVO: canal de reporte y política
README.md                        # MODIFICADO: sección "Seguridad" con el uso local
.github/workflows/secrets.yml    # NUEVO: job de detección en push/PR
specs/SPEC-secrets-hygiene.md    # este archivo
```

## Contrato

### 1. `.gitignore` (raíz)

Añadir el bloque **después** de la sección actual de Node/pnpm y antes de
`# IA` (el orden importa: las negaciones van al final de su grupo):

```gitignore
# Secretos y material sensible (cualquier nivel)
.env
.env.*
!.env.example
!**/.env.example
*.pem
*.key
*.p12
*.jks
*.keystore
.npmrc
coverage/
```

Notas de implementación:

- `apps/web/.gitignore` ya tiene `.env` + `!.env.example`; no se toca (el patrón
  raíz aplica a todos los niveles igualmente).
- `apps/mobile/.gitignore` solo cubre `.env*.local`; queda cubierto por la raíz.
- `*.key` no debe romper el build: verificar con `git ls-files | rg '\.key$'`
  (hoy devuelve vacío).
- No añadir `.agents/skills/` (ya está).

### 2. `.gitleaks.toml`

```toml
[extend]
useDefault = true

[[allowlists]]
description = "Ficheros de ejemplo, tests y datos de seed del repo"
paths = [
  '''\.env\.example$''',
  '''apps/[^/]+/\.env\.example$''',
  '''(^|/)testdata/''',
  '''(^|/)db/seeds/''',
  '''(^|/)docs/swagger\.(json|yaml)$''',
  '''(^|/).*\.test\.(ts|tsx)$''',
  '''apps/web/e2e/fixtures\.ts$''',
]
regexes = [
  '''(?i)(cambia-esto|placeholder|dummy|example|changeme|your[-_]?key)''',
]
stopwords = ["seedDummyHash", "ExponentPushToken"]
```

Regla: **no** añadir a la allowlist rutas de producción (`internal/`, `apps/*/src/`).
Si un falso positivo aparece ahí, se corrige el patrón, no se silencia el path.

### 3. `SECURITY.md`

Secciones obligatorias:

1. Versiones soportadas (solo `main`).
2. Cómo reportar: email de contacto o GitHub Security Advisories
   ("Report a vulnerability"), **nunca** un issue público.
3. Qué incluir: versión/commit, pasos, impacto, PoC sin datos reales.
4. Plazo de respuesta: acuse ≤ 72 h, triage ≤ 7 días, divulgación coordinada ≤ 90 días.
5. Rotación de credenciales: qué hacer si se commitea una llave
   (revocar → rotar → purgar historia → forzar push → auditar uso).
6. Alcance: qué se considera vulnerabilidad (auth, IDOR, inyección, DoS) y qué no.

### 4. CI — `.github/workflows/secrets.yml`

```yaml
name: Secrets
on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: read

jobs:
  gitleaks:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5
        with:
          fetch-depth: 0        # historia completa: detecta secretos en commits, no solo en el diff
      - name: Install gitleaks
        run: |
          VERSION=8.28.0        # ← fijar a la última estable al implementar
          curl -sSfL "https://github.com/gitleaks/gitleaks/releases/download/v${VERSION}/gitleaks_${VERSION}_linux_x64.tar.gz" \
            | sudo tar -xz -C /usr/local/bin gitleaks
          gitleaks version
      - name: Detect
        run: gitleaks detect --source . --redact --no-banner --exit-code 1
```

Reglas:

- Versión **pineada** y, si es posible, `sha256sum -c` del tarball (la tarea de
  implementación deja el hash verificado).
- `fetch-depth: 0` obligatorio (si no, solo escanea el commit actual).
- `permissions: contents: read` explícito.
- No se usa `gitleaks/gitleaks-action` para evitar la dependencia de licencia en
  organizaciones.

### 5. Hook local (opcional, documentado)

`lefthook.yml` **no** se modifica por defecto (el binario no está garantizado en
todas las máquinas y rompería el pre-commit). El uso local se documenta en
`README.md` → sección **Seguridad** (y `SECURITY.md` enlaza allí, para no
duplicar los comandos):

```bash
gitleaks detect --source . --redact   # historia completa
gitleaks protect --staged --redact    # solo lo staged
```

Si más adelante se quiere automatizar, se añade el job con `only:`/`glob:` y un
guard `command -v gitleaks || exit 0`.

### 6. Configuración de GitHub (manual, checklist)

Se aplica vía API (`gh api -X PATCH repos/{owner}/{repo} -f
"security_and_analysis[...]"`) y se verifica leyendo el mismo objeto.

- [x] **Secret scanning**: On
- [x] **Push protection**: On
- [x] **Dependabot alerts** y **security updates**: On (los consume `deps-supply-chain`)
- [ ] **Non-provider patterns** / **validity checks**: el PATCH responde 200 pero
      no aplica; requieren GitHub Secret Protection (plan de pago) y no están
      disponibles para este repo público de usuario.

**Alcance real de Push Protection**

GitHub bloquea por defecto solo **patrones de proveedor** con formato válido. Los
**patrones genéricos** (claves privadas, cadenas de conexión) no se bloquean
salvo con Secret Protection. Verificado con 4 canarios (RSA PKCS#1, EC, RSA
PKCS#8 y una API key de Google falsa): ninguno fue bloqueado y no se crearon
alertas.

Consecuencia: el canario con clave PEM **no** es una prueba válida en este repo y
la cobertura de genéricos queda en manos de gitleaks (CI + `protect --staged`). No
se debe "arreglar" esto cambiando la configuración de GitHub: es una limitación
del plan.

## Testing strategy

No hay tests unitarios: el módulo es configuración. La verificación es por
comandos y por el canario.

| Caso | Comando | Esperado |
|---|---|---|
| `.env.local` ignorado | `git check-ignore -v .env.local` | imprime la regla |
| `.env.example` **no** ignorado | `git check-ignore -v .env.example` | exit 1 (no ignorado) |
| `.env.production` ignorado | idem | imprime la regla |
| `*.pem` / `*.key` ignorados | idem | imprime la regla |
| Canario detectado | `gitleaks detect --no-git --source "$canary"` (clave PEM ficticia) | exit 1 + finding de clave privada |
| Allowlist acotada | mismo secreto en `db/seeds/` (esperado 0) y en `internal/` (esperado 1) | solo se silencia el path permitido |
| Repo limpio | `gitleaks detect --source . --redact` | exit 0 |
| CI | PR de prueba | job `gitleaks` verde |

## Boundaries

- **Always**: ejecutar `gitleaks detect` antes de un push que toque `.env.example`
  o CI; mantener `permissions: contents: read`; documentar en `SECURITY.md`
  cualquier allowlist nueva con su motivo.
- **Ask first**: añadir `gitleaks` a `go.mod` (`tool`), añadir un hook que
  bloquee commits, cambiar la política de divulgación.
- **Never**: commitear `.env` reales (ni "temporalmente"), silenciar una regla
  con un path de producción, poner secretos en issues/PRs, usar `--no-verify`.

## Success criteria

1. `git add -A` sobre un árbol con `.env.local`, `cert.pem` y `server.key` no los
   stagea.
2. `gitleaks detect` sobre el repo completo devuelve 0 halls.
3. El canario PEM falla con exit 1 y `--redact` (no imprime la clave).
4. `.github/workflows/secrets.yml` corre en PR y en push a `main` con
   `permissions: contents: read` y `fetch-depth: 0`.
5. `SECURITY.md` cubre reporte, plazos y rotación.
6. Secret scanning + push protection verificados con un push de prueba rechazado.

## Open questions

- ¿Se instala `gitleaks` en el runner vía tarball (propuesto) o se usa la acción
  oficial asumiendo la dependencia de licencia? → propuesto: tarball.
- ¿Se añade el hook local a `lefthook.yml` en una segunda iteración? → por
  ahora solo documentado.
- ¿Se cambia el email del autor en la historia? → **No** (decisión S-4).
