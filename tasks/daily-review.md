# Tasks — daily-review (backend + mobile)

> Decisión: D-52 (`specs/DECISIONS.md`) · Estado: **done**

## Problema

Pulsar "Repasar hoy" en mobile abría una sesión `daily` vacía y la pantalla mostraba
"Sesion completada · Respondiste 0 preguntas" sin haber respondido nada. La sesión
quedaba `in_progress` para siempre en el tab Sesiones y cada toque creaba otra.

Causa raíz: el warmup solo aceptaba preguntas `true_false` (la IA solo genera
`single_choice`, D-14 y el seed no traía ninguna), así que sin repasos vencidos
`GET /next` devolvía 404. Además el mobile trataba **cualquier** error de `/next`
como fin de sesión.

## Tareas

- [x] `store.go`: `FindNextDailyQuestion` encadena 3 escalones (vencidas → `true_false` →
      práctica normal). Nueva `findNextPractice` + `practiceCandidatesQuery` (excluye dominadas)
      y `CountPracticeCandidates`. Nueva `FindTodayInProgressDaily`.
      **Hecho:** `internal/modules/sessions/store.go`.
- [x] `service.go`: `gatherDailyData` cuenta también el escalón 3; `CreateDaily` responde **409**
      sin candidatos, **reutiliza** la daily de hoy en progreso (corte medianoche UTC) y recarga
      su tip; `Finish` no otorga recompensa si `len(Answers) == 0`.
      **Hecho:** `internal/modules/sessions/service.go`.
- [x] Swagger: descripción del endpoint + `409`.
      **Hecho:** `handler.go` y `make swagger` (`docs/swagger.{yaml,json}`, `docs/docs.go`).
- [x] Tests backend: guard 409 sin candidatos, reutilización con tip, Finish sin respuestas,
      handler 409.
      **Hecho:** `go test ./internal/modules/sessions/...` en verde.
- [x] Mobile: `SessionStudyPage` distingue 404 con respuestas (`finished`), 404 sin respuestas
      (`empty` con CTA a sesión de estudio / volver, descartando la sesión vacía) y 500/red
      (`error` con reintento). `HomeScreen` interpreta el 409 como estado neutro.
      **Hecho:** `SessionStudyPage.tsx`, `HomeScreen.tsx`, `core/utils/errors.ts` (`getErrorStatus`).
- [x] Tests mobile: 404 con respuestas, 404 sin respuestas, descarte de la sesión vacía, 500 con
      reintento, Home con 409, helper `getErrorStatus`.
      **Hecho:** `pnpm --filter mobile typecheck` + jest relacionados en verde.
- [x] Docs: D-52, este task y `.agents/context/backend/architecture.md`.

## Notas

- El escalón 2 **no** filtra por `is_mastered`: es calentamiento y se acepta repetir.
  El escalón 3 sí lo excluye para no re-preguntar lo ya dominado.
- `DailySessionResponse` **no** cambia de forma: `warmupCount` sigue siendo `true_false`;
  el conteo del escalón 3 es interno.
- Si en el futuro la IA genera `true_false`, el escalón 2 se activa solo, sin tocar código.
