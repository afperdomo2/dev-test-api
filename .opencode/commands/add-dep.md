---
description: Añade una dependencia tras pasar las reglas de seguridad del ecosistema
---
Dependencia: $ARGUMENTS

1. Identifica el ecosistema (backend Go / web / mobile). Si es ambiguo, pregunta
   antes de continuar.
2. Lee `.agents/<ecosistema>/dependencies.md` y ejecuta TODOS los checks:
   reputación del origen, mantenimiento, registro de seguridad/CVEs, licencia,
   calidad/API, supply chain. Investiga con websearch/webfetch (repo, última
   release, advisories).
3. Si algún check falla → NO instales; explica el riesgo y propone alternativas.
4. Solo si todo pasa, instala:
   - backend: `go get <module>` + `go mod tidy`
   - web: `pnpm add <pkg>` (o `-D` para dev)
   - mobile: `pnpm --filter mobile add <pkg>`
5. Verifica tras instalar: `make check` (backend) o typechecks del workspace.
