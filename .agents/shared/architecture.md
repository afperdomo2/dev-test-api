# Shared package architecture

## Purpose

`@devtest/shared` (`packages/shared/`) contains **pure TypeScript** shared across web and mobile. It holds DTOs, type definitions, format utilities, validation predicates, and constants.

## What goes in shared

| ✅ Belongs in shared | ❌ Does NOT belong in shared |
|---|---|
| TypeScript interfaces (DTOs, API contracts) | Vue/Nuxt components or composables |
| Pure validation predicates (`isValidEmail`, `isRequired`) | Vuetify validation rules (`requiredRule`, `emailRule`) |
| Format functions (`formatDateTime`, `formatScore`) | Anything using `localStorage`, `import.meta.env`, DOM APIs |
| Constants (`ITEMS_PER_PAGE_OPTIONS`, `DEFAULT_PER_PAGE`) | Anything using `document`, `window`, `navigator` |
| Type-only re-exports from Swagger (`ApiError`, `PaginatedResponse`) | Platform-specific storage (SecureStore, AsyncStorage) |

## File structure

```
packages/shared/src/
├── types/        # DTOs + API types (mirrors docs/swagger.yaml)
│   ├── api.types.ts        # ApiError, ApiResponse, PaginatedResponse
│   ├── auth.types.ts       # LoginRequest, AuthResponse
│   ├── user.types.ts       # User, CreateUserRequest
│   ├── question.types.ts   # Question, QuestionType, QuestionStats
│   ├── session.types.ts    # Session, SessionAnswer
│   ├── topic.types.ts      # Topic, CreateTopicRequest
│   ├── progress.types.ts   # Progress, UpcomingQuestion
│   └── opencode.types.ts   # OpenCodeUsage
├── utils/
│   ├── format.ts           # formatDateTime, formatScore, formatTimeUntil
│   └── validators.ts       # isValidEmail, isRequired, isMinLength, isMaxLength
├── constants/
│   ├── index.ts            # ITEMS_PER_PAGE_OPTIONS, DEFAULT_PER_PAGE, REFRESH_COOLDOWN_MS
│   └── questionImport.ts   # buildImportPrompt, IMPORT_CSV_HEADER
└── index.ts                # Barrel export (re-exports all)
```

## Package config

- **`package.json`**: `"main": "./src/index.ts"`, `"types": "./src/index.ts"` — exports TypeScript source directly (no build step)
- Both Vite (web) and Metro (mobile) transpile `.ts` from `node_modules` workspace symlinks
- **`tsconfig.json`**: `"composite": true`, `"outDir": "dist"` — used by `tsc --noEmit` for type-checking only

## Barrel exports

`src/index.ts` re-exports everything from all modules. Both consumers import from the package root:

```ts
import type { User, Session, ApiError } from '@devtest/shared'
import { formatDateTime, isValidEmail } from '@devtest/shared'
```

## Re-export pattern in web

`apps/web/src/types/*.ts` contains **thin re-exports** from `@devtest/shared`:

```ts
// apps/web/src/types/user.types.ts
export { DEFAULT_DAILY_IMPORT_LIMIT, DEFAULT_DAILY_AI_LIMIT } from '@devtest/shared'
export type { User, CreateUserRequest, UpdateUserRequest } from '@devtest/shared'
```

This preserves all existing `@/types/...` imports in web without changes. **Keep these re-exports in sync** — if you add a new type to shared, add it to the corresponding web re-export file.

## CodeChallenge deduplication

`CodeChallenge` is defined once in `question.types.ts` and imported (not re-exported) in `session.types.ts`. The barrel `index.ts` exports it from `question.types.ts` only. Do not re-export the same type from multiple modules.

## Sync: Swagger → shared → apps

The backend is the source of truth for API contracts (`docs/swagger.yaml`).

1. **Backend changes** (`internal/models/`, `response.go`, `request.go`) → run `make swagger`
2. **Shared types** (`packages/shared/src/types/`) → update to match new swagger shapes
3. **Web re-exports** (`apps/web/src/types/`) → add new types to re-export files
4. **Mobile** → automatically gets new types via `@devtest/shared` import

**Key rule**: if a field is added to a Go struct's JSON tags, it must appear in `packages/shared/src/types/*.types.ts`. The re-exports in web propagate it to `@/types/...`. Mobile imports from `@devtest/shared` directly.

## Tests

Tests live in `packages/shared/src/utils/` (co-located with source). Run via `pnpm --filter @devtest/shared typecheck` (type-check only; test runner TBD).

Test files are excluded from `tsconfig.json` (no vitest globals in typecheck scope).
