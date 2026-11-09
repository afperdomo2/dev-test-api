# Tasks — daily-review (backend + mobile)

> Decisión: D-52 (`specs/DECISIONS.md`) · Estado: **done**

## Problema

Pulsar "Repasar hoy" en mobile abría una sesión `daily` vacía y la pantalla mostraba
"Sesion completada · Respondiste 0 preguntas" sin haber respondido nada. La sesión
quedaba `in_progress` para siempre en el tab Sesiones y cada toque creaba otra.

Causa raíz: el warmup solo aceptaba preguntas `true_false` y el seed no traía ninguna (en ese
momento la IA solo generaba `single_choice`, D-14; hoy D-54 permite los 3 tipos), así que sin
repasos vencidos `GET /next` devolvía 404. Además el mobile trataba **cualquier** error de `/next`
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

## Cupos por sesión (D-55)

Problema: con D-54 la IA puede generar `true_false` en masa y, como el calentamiento no
tenía tope, **el repaso entero eran verdadero/falso** y nunca llegaban las de selección.

- [x] `dailyTrueFalseLimit = 3`: `NextQuestion` cuenta las `true_false` ya respondidas en la
      sesión (`CountSessionAnswersByType`) y, al llegar a 3, pasa `skipTrueFalse` para saltar
      también las vencidas de ese tipo. Las pendientes siguen vencidas para el día siguiente.
- [x] `dailyQuestionLimit = 15`: `CreateDaily` fija `QuestionLimit = 15`; el mobile ya muestra
      "N / 15", la barra de progreso y termina en "Sesión completada".
- [x] Al reutilizar una daily de hoy creada antes del cambio (sin tope), `ensureDailyLimit`
      le aplica `QuestionLimit = 15` (best effort).
- [x] Tests: cupo disponible/agotado, error del contador, tope aplicado al crear, tope aplicado
      al reutilizar y no reescrito si ya existía.
      **Hecho:** `go test -count=1 ./internal/modules/sessions/...` en verde.
- [x] Docs: D-55 y `.agents/context/backend/architecture.md`.
- [x] Mobile: **sin cambios** (el `questionLimit` ya viaja en `SessionResponse` y la UI lo usa).

## Contrato de pregunta y `code_completion` (D-56)

- [x] `toNextQuestionItem` emite `options` y `topics` como arrays **siempre** (`[]`, nunca
      omitidos ni `null`); se quitó `omitempty` de `Options`.
      **Hecho:** `sessions/mapper.go`, `sessions/response.go` + `mapper_test.go` (JSON con `[]`).
- [x] `code_completion` excluida de la selección de preguntas: `FindNextQuestion`,
      `findNextDueReview`, `practiceCandidatesQuery`, `CountAvailableQuestions` y
      `CountDueReviews` (con `JOIN questions` y filtro de visibilidad alineado).
      **Hecho:** `sessions/store.go`; el crash ocurría al llegar al escalón 3 (D-55).
- [x] Mobile: `QuestionCard` usa `question.options ?? []` y `question.topics ?? []`.
      **Hecho:** test "no rompe si el backend omite options/topics".
- [x] Docs: D-56. Sin migración. `make swagger` no cambia (el spec no modela `omitempty`).

## Notas

- El escalón 2 **no** filtra por `is_mastered`: es calentamiento y se acepta repetir.
  El escalón 3 sí lo excluye para no re-preguntar lo ya dominado.
- `DailySessionResponse` **no** cambia de forma: `warmupCount` sigue siendo `true_false`;
  el conteo del escalón 3 es interno.
- Si en el futuro la IA genera `true_false`, el escalón 2 se activa solo, sin tocar código.
