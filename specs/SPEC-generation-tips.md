# SPEC — generation-tips

| Campo | Valor |
|---|---|
| Module id | `generation-tips` |
| Estado | Implementado |
| Decisión | D-51 (presupuesto compartido; pipeline reutilizado con `kind`) |
| Fecha | 2026-11-07 |

## Contexto

`/generation` solo generaba preguntas con IA. Para tips existía únicamente el prompt copy-paste de `QuestionImportPage` (`buildTipsImportPrompt`), que obliga al admin a salir a ChatGPT/Claude y volver a importar el CSV.

Objetivo: generar tips con IA desde el mismo panel de administración, sin duplicar cola, worker pool, presupuesto diario ni la UI de historial.

## Alcance

**Backend**

- `generation_jobs.kind` (`question` | `tip`, default `question`) + `CHECK` + índice `(kind, status, created_at)` (migración `00010_generation_jobs_kind.sql`).
- `ai.Generator.GenerateTips` + `TipBatchRequest` + `tipsSystemPrompt` (JSON `[{content, topic}]`, `topic` = slug de un tema del job).
- `ProcessJob` ramifica por `kind`: `GenerateQuestions`/`SaveQuestions` o `GenerateTips`/`SaveTips`.
- `Store.SaveTips` resuelve `slug`/nombre → `topic_id` contra los temas del job y descarta tips con tema desconocido; `source = 'ai_generated'`.
- `Store.ExistingTipContentByTopics` para evitar duplicados en el prompt.
- Contrato: `kind` opcional en `POST /generation`; `?kind=` en `GET /generation` y `GET /generation/stats` (los contadores de cola se filtran, el presupuesto es global).
- `difficulty` es opcional en el request y se ignora para tips (la columna queda en `intermediate` para respetar el `CHECK`).

**Web**

- `/generation` con tabs **Preguntas** / **Tips** (`?type=` persiste el tab).
- `GenerationPanel` reutiliza `GenerationEnqueueForm`, `GenerationStatsCards` y la tabla con un prop `kind`.
- En Tips se oculta Dificultad y la tabla no muestra la columna Dificultad.
- La tarjeta de stats aclara que el tope diario es compartido.

**Shared**

- `GenerationKind`; `kind` en `GenerationJob` y `CreateGenerationJobRequest`; `difficulty` pasa a opcional.

## Fuera de alcance

- Presupuesto separado por tipo (`ai.tips.max_per_day`, etc.). Se descartó por ahora: un tip cuenta como una unidad del cupo global.
- Revisión/moderación previa de tips generados (se publican al instante, igual que el import actual). Nota N-03: moderación no puede des-publicar tips.
- Mobile: no está expuesto en la app.

## Archivos

```
internal/database/migrations/00010_generation_jobs_kind.sql   NUEVO
internal/models/generation_job.go                             kind + constantes
internal/services/ai/tips.go                                  NUEVO: GenerateTips + parse/validación
internal/services/ai/tips_test.go                             NUEVO
internal/services/ai/prompts.go                               tipsSystemPrompt
internal/modules/generation/{service,store,request,response,mapper,handler,worker}.go
apps/web/src/features/generation/components/GenerationPanel.vue             NUEVO
apps/web/src/features/generation/{pages/GenerationPage,components/GenerationEnqueueForm,components/GenerationStatsCards}.vue
apps/web/src/{api/services,queries}/generation.*.ts
packages/shared/src/types/generation.types.ts
```

## Verificación

- `go test -short -count=1 ./internal/modules/generation/... ./internal/services/ai/...`
- `pnpm --filter @devtest/shared typecheck`
- `apps/web`: `vue-tsc --build`, eslint de los archivos tocados, `vitest run` de `generation.queries.test.ts` y `generation.service.test.ts`.
- `make swagger` (cambian `POST /generation`, `GET /generation`, `GET /generation/stats` y `CreateJobRequest`).

## Criterios de aceptación

1. Encolar un job `kind=tip` genera tips vía IA, los guarda en `tips` con `source='ai_generated'` y completa el job con la cantidad real guardada.
2. Un job de tips no toca `questions`; uno de preguntas no toca `tips`.
3. El historial y los contadores de cada tab solo muestran su `kind`; la barra de consumo es la misma en ambos (presupuesto compartido).
4. Si la IA devuelve un slug desconocido, ese tip se descarta y el resto se guarda.
5. `POST /generation` sin `kind` sigue creando preguntas (retrocompatible).
