# SPEC — deps-supply-chain

| Campo | Valor |
|---|---|
| Module id | `deps-supply-chain` |
| Origen | Revisión de seguridad 2026-10-21 (M3, L6) |
| Estado | `ready` |
| Depende de | `secrets-hygiene` (comparte `.github/workflows/`) |
| Consumido por | — |
| Migración goose | Ninguna |
| Dependencia nueva | Ninguna directa (solo bumps y overrides) |
| Cambio de contrato | No (salvo que el bump de Vite cambie el build) |

## Objective

`govulncheck ./...` reporta **14 vulnerabilidades alcanzables** (toolchain Go
1.27.0, `golang.org/x/net v0.58.0`, `quic-go v0.59.0`) y `pnpm audit --prod`
reporta **21** (1 crítica, 9 altas, 9 medias, 2 bajas), concentradas en
`markdown-it@10` (mobile runtime), `follow-redirects` (axios) y tooling
(Expo CLI, Vite, devtools).

Hoy nada de eso se vigila en CI: no hay `govulncheck`, no hay gate de auditoría y
las actions están referenciadas por tag mayor.

**Historias**

- Como mantenedor, quiero que una CVE publicada aparezca como PR automático, no
  como sorpresa.
- Como operador, quiero que el build falle si entra una vulnerabilidad alta en
  una dependencia de producción.
- Como revisor, quiero que las actions estén pineadas por SHA.

**Éxito**

- `govulncheck ./...` → 0 vulnerabilidades alcanzables.
- `pnpm audit --prod` → 0 high/critical fuera de la allowlist documentada.
- `Dependabot` abre PRs semanales agrupados para Go, pnpm y Actions.
- CI falla ante una regresión (verificado con un canary temporal).

## Tech stack

Go 1.27.3, pnpm 10, GitHub Actions. Sin dependencias nuevas de runtime.

## Commands

```bash
# Go
go get golang.org/x/net@v0.60.0 github.com/quic-go/quic-go@v0.59.1
go mod tidy
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
make check && make test && make build

# Node
pnpm install
pnpm audit --prod --audit-level high
pnpm -r typecheck && pnpm -r lint
pnpm --filter dev-test-frontend test:run
pnpm --filter mobile test:run
pnpm --filter mobile exec expo install --check   # compatibilidad Expo (requiere red)

# CI local
act -j build
```

## Project structure

```
go.mod / go.sum                       # toolchain + x/net + quic-go
package.json                          # pnpm.overrides + pnpm.auditConfig.ignoreCves
pnpm-lock.yaml                        # regenerado
.github/workflows/backend.yml         # permissions + govulncheck
.github/workflows/frontend.yml        # permissions + audit
.github/dependabot.yml                # NUEVO
Makefile                              # + vuln, audit
AGENTS.md                             # tabla de comandos
```

## Contrato

### 1. Go

| Elemento | Antes | Después |
|---|---|---|
| directiva `go` | `1.26.6` | `1.27.3` |
| `golang.org/x/net` | `v0.58.0` | `v0.60.0` |
| `github.com/quic-go/quic-go` | `v0.59.0` | `v0.59.1` |

- `go get …` + `go mod tidy`; **no** editar `go.sum` a mano.
- `go-version-file: go.mod` en CI instala automáticamente 1.27.3.
- Aceptación: `govulncheck ./...` sin "Your code is affected by…".

### 2. pnpm — overrides mínimos (no `pnpm up` global)

`package.json` (raíz):

```json
{
  "pnpm": {
    "overrides": {
      "follow-redirects@<1.16.1": "1.16.1",
      "shell-quote@<1.11.0": "1.11.0"
    },
    "auditConfig": {
      "ignoreCves": [
        "GHSA-vfj7-8cjw-p6xm",
        "GHSA-86w9-cpqp-85rv",
        "GHSA-hp3w-g68c-fv3c"
      ]
    }
  }
}
```

Reglas:

- **No** se ejecuta `pnpm up` global: rompería el pinning de Expo
  (`expo ~57.0.23`). Solo overrides de transitivas vulnerables.
- Cada `ignoreCves` lleva un comentario con motivo y **fecha de revisión** en
  `SECURITY.md`:
  - `braces@3.0.3` y `node-forge@1.4.0`: sin versión parcheada publicada;
    alcanzables solo desde `expo/@expo/cli` en tiempo de build.
  - `sprintf-js`: solo en la cadena de `jest` (test).
- `markdown-it@10` (mobile runtime) se trata aparte (ver §3).

### 3. `markdown-it@10` en mobile — spike con fallback

`apps/mobile` usa `react-native-markdown-display@^7.0.2` → `markdown-it@10.0.0`
con 5 advisories (1 alta de complejidad cuadrática). Dos caminos, en este orden:

1. **Bump del wrapper**: comprobar si `react-native-markdown-display` publicó una
   versión con `markdown-it >= 14.3.3`. Si existe y es compatible con Expo 57, se
   actualiza y se elimina el override de `markdown-it`.
2. **Fallback si no hay versión compatible**:
   - Añadir `pnpm.overrides { "markdown-it@<14.3.3": "14.3.3" }` **solo** si el
     render de `MarkdownText` sigue funcionando (test de humo con negrita, código,
     listas y enlaces).
   - Si el override rompe la API del wrapper: **no** forzar. Mitigar en código:
     truncar el contenido a `MAX_MARKDOWN_CHARS = 8000` en
     `apps/mobile/src/core/components/MarkdownText.tsx` antes de renderizar
     (el contenido viene de IA/preguntas; un texto patológico deja de colgar el
     hilo de UI), y documentar el riesgo aceptado en `SECURITY.md`.

La decisión queda registrada en `specs/DECISIONS.md` al implementar.

### 4. Vite

`apps/web` usa `vite ^6.2.6`, resuelto a 6.4.3. Subir a `^6.4.4` (mismo major)
y verificar `pnpm build`. Las advisories de Vite afectan al **dev server** y al
preview (no al bundle de producción), pero el gate de auditoría las cuenta.

### 5. CI

`backend.yml`:

```yaml
permissions:
  contents: read
# …
- name: Vulnerabilities
  run: go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

`frontend.yml`:

```yaml
permissions:
  contents: read
# …
- name: Audit
  run: pnpm audit --prod --audit-level high
```

Pinning de actions por SHA (obligatorio en este módulo):

| Action | Versión actual | Acción |
|---|---|---|
| `actions/checkout` | `@v5` | fijar SHA de v5 y comentar la versión |
| `actions/setup-go` | `@v5` | idem |
| `actions/setup-node` | `@v4` | idem |
| `pnpm/action-setup` | `@v5` | idem |

> El SHA se obtiene de la release oficial; se deja como
> `uses: owner/action@<sha> # vX.Y.Z`.

### 6. Dependabot (`.github/dependabot.yml`)

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule: { interval: weekly, day: monday }
    groups:
      go-deps:
        patterns: ["*"]
  - package-ecosystem: npm
    directory: /
    schedule: { interval: weekly, day: monday }
    groups:
      js-deps:
        patterns: ["*"]
    ignore:
      - dependency-name: expo
      - dependency-name: react-native
      - dependency-name: "expo-*"
  - package-ecosystem: github-actions
    directory: /
    schedule: { interval: weekly, day: monday }
```

Los paquetes de Expo se ignoran porque se actualizan con
`expo install`, no con Dependabot (evita romper la matriz del SDK).

### 7. Makefile + AGENTS.md

```make
vuln:   ## govulncheck en todo el módulo
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

audit:  ## auditoría de dependencias JS de producción
	pnpm audit --prod
```

Añadir ambas filas a la tabla de comandos de `AGENTS.md`.

## Testing strategy

| # | Verificación | Comando | Esperado |
|---|---|---|---|
| T1 | Go sin vulns alcanzables | `govulncheck ./...` | exit 0 |
| T2 | Backend intacto | `make check && make test` | verde |
| T3 | Auditoría JS | `pnpm audit --prod --audit-level high` | exit 0 |
| T4 | Overrides aplicados | `pnpm why follow-redirects` | 1.16.1 |
| T5 | Build web | `pnpm --filter dev-test-frontend build` | verde |
| T6 | Mobile intacta | `pnpm --filter mobile typecheck && test:run` | verde |
| T7 | Markdown mobile | test de humo de `MarkdownText` (negrita, fence, lista, enlace) | render correcto |
| T8 | CI gate | PR con override revertido temporalmente | job rojo |

## Boundaries

- **Always**: cambios de dependencia con `go mod tidy` / `pnpm install` y lock
  actualizado; justificar cada `ignoreCves` con fecha de revisión; correr los
  tests de las 3 apps.
- **Ask first**: subir de major `vite`/`expo`/`react`; añadir overrides que
  cambien el runtime de mobile; ignorar una CVE alta de runtime.
- **Never**: `pnpm up` global, `go get -u ./...`, borrar tests que fallan tras un
  bump, ignorar una CVE sin entrada en `SECURITY.md`, usar `npm audit fix --force`.

## Success criteria

1. `govulncheck ./...` = 0 alcanzables y los 3 comandos de build/test verdes.
2. `pnpm audit --prod --audit-level high` = 0, con solo las 3 CVEs allowlisteadas
   documentadas.
3. `.github/dependabot.yml` activo (GitHub muestra los PRs la semana siguiente).
4. Los 4 workflows usan `permissions: contents: read` y actions por SHA.
5. `make vuln` y `make audit` documentados en `AGENTS.md`.
6. La app mobile sigue renderizando markdown (test de humo) y el build de web
   pasa con Vite 6.4.4.

## Open questions

- ¿Hay versión de `react-native-markdown-display` con `markdown-it >= 14`? Se
  resuelve en el spike; si no, se aplica el fallback de truncado.
- ¿Se permite a Dependabot tocar `expo-*`? (propuesto: no).
- ¿Se añade `trivy` para escanear la imagen Docker? No hay Dockerfile en el repo
  hoy; fuera de alcance.
