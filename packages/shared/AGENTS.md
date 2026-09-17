# Shared (@devtest/shared) — quick context

Paquete TypeScript puro (sin runtime de React/Vue) compartido entre **web** y **mobile**.

## Commands

| Command | Action |
|---------|--------|
| `pnpm --filter @devtest/shared typecheck` | Type-check |

## What lives here

- `types/` — DTOs y tipos de API (`ApiError`, `PaginatedResponse`, `Topic`, `Question`, `Progress`, …)
- `constants/` — `ITEMS_PER_PAGE_OPTIONS`, `DEFAULT_PER_PAGE`, `QUESTION_*` maps
- `utils/` — `format.ts` (formatDateTime, formatScore), `validators.ts` (isValidEmail, isRequired)

Barrel export: `src/index.ts`.

## Rules

- **Nunca duplicar**: si algo es reutilizable entre web y mobile, pertenece aquí.
- Web re-exporta los tipos en `apps/web/src/types/*.types.ts`; mobile importa vía `@devtest/shared` directo.
- Flujo de sync: `docs/swagger.yaml` → `packages/shared/src/types/` → apps.

## Details

- Arquitectura y sync: `.agents/shared/architecture.md`