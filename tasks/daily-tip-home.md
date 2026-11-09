# Tasks — daily-tip-home

> Decisión: D-53 (`specs/DECISIONS.md`) · Estado: **done**
> Reemplaza el comportamiento de `specs/SPEC-mobile-daily-tip.md` (histórico).

## Problema

El tip del día se pintaba dentro de la pantalla de estudio, arriba de la primera pregunta.
Como el tip era un `RANDOM()` global sobre `tips` (sin relación con el usuario, el día ni la
pregunta), un usuario podía leer un tip de Angular y responder una pregunta de Go. Además el
tip "reaparecía" al reutilizar la sesión diaria.

## Comportamiento nuevo

- El tip vive en **Inicio**, debajo del CTA "Sesión de estudio".
- `GET /tips/today` lo elige **determinista por fecha UTC**: mismo tip todo el día, rota al siguiente.
- Tiene botón de cerrar; el cierre se recuerda **por dispositivo y por día** (SecureStore).
  Al día siguiente vuelve a mostrarse.
- Si no hay tips (404) o falla la petición, no se muestra nada y la Home sigue igual.
- La pantalla de estudio ya no muestra tip.

## Tareas

- [x] Backend: `tips.Store.FindOfTheDay(day)` (`ORDER BY md5(id || fecha)`) reemplaza a `FindRandom`.
      **Hecho:** `internal/modules/tips/store.go`.
- [x] Backend: `GetTipOfTheDay()` sin `userID` (los tips son globales) + handler `GetOfTheDay` + ruta
      `GET /tips/today` (antes de `/:id`).
      **Hecho:** `tips/service.go`, `tips/handler.go`, `tips/routes.go`.
- [x] Backend: quitar `tip` de `POST /sessions/daily` (`DailySessionResponse`) y la dependencia
      `tips.Service` de `sessions`.
      **Hecho:** `sessions/response.go`, `sessions/service.go`, `internal/server/server.go`.
- [x] Shared: quitar `tip?: Tip` de `DailySessionResponse`.
      **Hecho:** `packages/shared/src/types/session.types.ts`.
- [x] Tests backend: servicio (tip del día), handler (200/404), `CreateDaily` sin tip.
      **Hecho:** `go test ./internal/modules/tips/... ./internal/modules/sessions/...` en verde.
- [x] Mobile: `getTipOfTheDay()` + tarjeta en Home + cierre persistido por día.
      **Hecho:** `core/api/services/tips.service.ts`, `core/components/DailyTipCard.tsx` (movido),
      `features/home/dailyTipDismissal.ts`, `HomeScreen.tsx`.
- [x] Mobile: quitar el tip de la pantalla de estudio y borrar `dailyTipStore`.
      **Hecho:** `SessionStudyPage.tsx`; `features/sessions/dailyTipStore{,.test}.ts` eliminados.
- [x] Tests mobile: tips.service, dismissal, tarjeta, Home (muestra/cierra/ya cerrado/sin tip),
      estudio sin tip.
      **Hecho:** jest relacionados en verde + `pnpm --filter mobile typecheck`.
- [x] Docs: D-53, este task, `.agents/context/{backend,mobile}/architecture.md`, spec histórico marcado.
      **Hecho:** `make swagger` (nuevo endpoint) + `docs/swagger.*`.

## Notas

- La preferencia de cierre es del **dispositivo**, no de la cuenta: otro usuario en el mismo
  teléfono hereda el cierre del día. Aceptable para una tarjeta de lectura; si molesta, la clave
  pasa a incluir el `userId`.
- `sessions.tip_id` sigue en la DB (no se toca una migración aplicada) pero ya no se escribe.
- Un tip manual/importado es global: `GET /tips` siempre listó todos, así que no hay scoping por usuario.
