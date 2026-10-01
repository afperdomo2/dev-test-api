---
description: Corre los tests/checks de los paquetes afectados
---
Archivos cambiados:
!`git status --short`

Ejecuta solo los checks de las áreas afectadas (no hay jest en el repo):
- Backend (raíz, Go): `make test` (unit, sin DB)
- Web (apps/web, vitest): `pnpm web:test` + `make fe-check` + `make fe-lint`
- Mobile (apps/mobile): sin runner de tests → `pnpm mobile:typecheck`
- Shared (packages/shared): `pnpm --filter @devtest/shared typecheck`

Si algo falla, analiza la causa y propón el arreglo.
