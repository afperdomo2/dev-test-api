---
description: Sincroniza modelos de Go con los tipos TypeScript (shared → apps)
---
Archivos Go cambiados:
!`git diff HEAD --name-only -- '*.go'`

El backend (raíz: `internal/`, `pkg/`) es la fuente de verdad del contrato API.
Sigue el flujo de sync de `.agents/shared/architecture.md`:
1. Si cambiaron shapes de endpoints: regenera con `make swagger` y verifica `docs/swagger.yaml`.
2. Actualiza `packages/shared/src/types/` (fuente de verdad TS, camelCase en JSON).
3. Verifica los re-exports en `apps/web/src/types/` y los imports de mobile (`@devtest/shared`).
4. Detecta rupturas: `pnpm --filter @devtest/shared typecheck` y `pnpm typecheck`.

Señala cualquier endpoint/shape que haya quedado sin sincronizar.
