# SPEC — mobile-offline

| Campo | Valor |
|---|---|
| Module id | `mobile-offline` |
| Origen | Alcance nuevo mobile (D-26, D-31, D-35) |
| Estado | `ready` |
| Depende de | — (backend: cambio pequeño en `sessions`) |
| Migración goose | `00006_answer_idempotency.sql` (siguiente libre; `notifications` pasa a `00007` y `push` a `00008`) |
| Dependencia nueva | `expo-crypto` (`~57.0.4`, SDK 57) — gate de `.agents/mobile/dependencies.md` OK |

## Objective

Que un fallo de red al enviar una respuesta **no se trague en silencio** y que reintentar sea **seguro** (sin duplicar la respuesta, el XP ni el progreso SRS).

**Historias**:
- Respondo y se va la cobertura: veo "Sin conexión. Inténtalo de nuevo" y un botón **Reintentar**; mi selección sigue ahí.
- Al reintentar, la respuesta se guarda **una vez**, aunque el primer intento sí hubiera llegado al servidor.
- Si el error es del servidor (4xx/5xx), veo el detalle y también puedo reintentar.
- Un doble toque en "Enviar" no crea dos respuestas.

**Éxito**: ningún fallo de envío queda invisible, y ninguna respuesta se duplica.

## Contexto verificado (el bug)

```ts
// apps/mobile/src/features/sessions/screens/SessionStudyPage.tsx — handleAnswer
try {
  const result = await submitAnswer(id, { questionId: question.id, selectedOptions: options })
  ...
  setPhase('result')
} catch {
  setPhase('result')   // ← se traga el error: pasa a "resultado" sin resultado
}
```

Y en el backend, `sessionService.Answer` **solo** rechaza si la sesión está completada; **no** detecta respuestas duplicadas para la misma pregunta (`sessions/service.go:304-341`). Por tanto un reintento a ciegas duplicaría la fila, el XP y el SRS. De ahí D-35.

## Commands

```bash
# Backend (idempotencia)
make migrate-create name=answer_idempotency
make migrate-up && make migrate-status
make check && make test

# Mobile
pnpm --filter mobile typecheck
pnpm --filter mobile test:run
pnpm --filter mobile lint
```

## Parte A — Backend: idempotencia (D-35)

### Migración

```sql
-- +goose Up
ALTER TABLE session_answers ADD COLUMN client_answer_id uuid;
CREATE UNIQUE INDEX session_answers_client_uniq
    ON session_answers (session_id, client_answer_id);
COMMENT ON COLUMN session_answers.client_answer_id
    IS 'UUID generado por el cliente para idempotencia de reintentos. NULL en filas antiguas o clientes que no lo envían.';
-- +goose Down
DROP INDEX IF EXISTS session_answers_client_uniq;
ALTER TABLE session_answers DROP COLUMN IF EXISTS client_answer_id;
```

Postgres trata los `NULL` como distintos en un índice único, así que las filas históricas no chocan.

### Request y modelo

```go
type AnswerRequest struct {
	// ...campos actuales...
	ClientAnswerID *uuid.UUID `json:"clientAnswerId,omitempty"`
}

type SessionAnswer struct {
	// ...campos actuales...
	ClientAnswerID *uuid.UUID `gorm:"type:uuid" json:"clientAnswerId,omitempty"`
}
```

### Store

```go
// CreateAnswer inserta; si choca con (session_id, client_answer_id) no inserta
// y devuelve created=false + la fila existente.
CreateAnswer(answer *models.SessionAnswer) (created bool, err error)
```

Implementación: `clause.OnConflict{Columns: [session_id, client_answer_id], DoNothing: true}`; si `RowsAffected == 0`, `FindByClientAnswerID(sessionID, clientAnswerID)` y devolverla.

### Service

```go
created, err := s.store.CreateAnswer(answer)
if !created {
	// Reintento: devolver la respuesta ya guardada y NO volver a aplicar
	// progreso/gamificación (si no, doble XP y doble SRS).
	return ToAnswerResponse(*existing), nil
}
// ...camino actual: progreso + reward...
```

**Clave**: el `progressService.Answer` (que otorga XP y mueve el SRS) se ejecuta **solo** cuando la fila es nueva. Sin `clientAnswerId`, `created` siempre es `true` (comportamiento actual).

## Parte B — Mobile

### Helper

```ts
// src/core/utils/errors.ts (ampliar)
export function isNetworkError(err: unknown): boolean
// true si el error es el `{ title: 'Network Error' }` que lanza api/client.ts
```

### `SessionStudyPage`

- Al empezar una pregunta, generar **un** `clientAnswerId` (`Crypto.randomUUID()` de `expo-crypto`) y guardarlo en estado junto a la pregunta. **Se reutiliza en los reintentos** y se descarta al pasar a la siguiente pregunta.
- `handleAnswer`:
  - `phase='loading'` (deshabilita el botón: ya protege del doble toque).
  - éxito → como hoy.
  - error de red → `phase='question'`, conservar `selectedOptions`, `submitError='Sin conexión. Revisa tu red e inténtalo de nuevo.'` y mostrar **Reintentar**.
  - error de API → `phase='question'`, `submitError=getErrorDetail(err, 'No se pudo enviar la respuesta')` y **Reintentar** (seguro por la idempotencia).
- "Reintentar" vuelve a llamar a `handleAnswer` con las mismas opciones y el **mismo** `clientAnswerId`.
- El error se limpia al reintentar o al cambiar de pregunta.

### UI del error

- Bloque bajo el `QuestionCard`: texto del error + botón "Reintentar" (touch target ≥ 44 pt, `accessibilityRole="button"`).
- El botón de enviar vuelve a estar habilitado.
- No se avanza de pregunta ni se muestra feedback de acierto/fallo.

## Project structure

```
internal/database/migrations/00006_answer_idempotency.sql
internal/models/session_answer.go            # + ClientAnswerID
internal/modules/sessions/request.go         # + clientAnswerId
internal/modules/sessions/store.go           # CreateAnswer idempotente + FindByClientAnswerID
internal/modules/sessions/service.go         # no duplica progreso/XP
internal/modules/sessions/*_test.go

apps/mobile/src/core/utils/errors.ts          # + isNetworkError
apps/mobile/src/core/utils/errors.test.ts     # + casos
apps/mobile/src/core/api/services/sessions.service.ts  # + clientAnswerId en el body
apps/mobile/src/features/sessions/screens/SessionStudyPage.tsx
apps/mobile/src/features/sessions/screens/SessionStudyPage.test.tsx
```

## Testing strategy

| Nivel | Caso |
|---|---|
| Backend service | misma `(session, clientAnswerId)` dos veces → `created=false`, devuelve la misma respuesta, **`progressService.Answer` se llama una sola vez**; sin `clientAnswerId` → comportamiento actual (dos filas posibles, como hoy) |
| Backend store | índice único funciona; `NULL` permite duplicados históricos |
| Backend migration | `make migrate-up`/`down` limpios contra la DB real |
| Mobile `errors` | `isNetworkError` true para `{title:'Network Error'}`, false para 4xx/5xx y para `Error` genérico |
| Mobile screen | error de red → mensaje + "Reintentar", **no** avanza a resultado; reintento envía el **mismo** `clientAnswerId` (espía del servicio); doble toque envía una sola vez |
| Mobile service | el body incluye `clientAnswerId` cuando se pasa |

## Boundaries

- **Always**: el `clientAnswerId` se genera una vez por pregunta; reintentar es seguro; el error se muestra; los tests del comportamiento actual (de `mobile-tooling`) se actualizan junto con el cambio.
- **Ask first**: persistir respuestas pendientes en disco (X-11); reintentos automáticos con backoff; cambiar la forma de `CreateAnswer` para otros módulos.
- **Never**: reintentar sin idempotencia; duplicar XP/SRS en un reintento; silenciar el error; guardar el `clientAnswerId` en el backend como clave de negocio más allá de la idempotencia.

## Success Criteria

1. Con red, el flujo de responder funciona igual que hoy.
2. Sin red, no se avanza al resultado: se ve el mensaje y **Reintentar**.
3. Reintentar con red guarda la respuesta **una sola vez** (una fila, un XP).
4. Enviar dos veces el mismo `clientAnswerId` (doble toque o reintento) devuelve la misma respuesta y no duplica XP/SRS.
5. Al cambiar de pregunta, el `clientAnswerId` es distinto.
6. `make migrate-up`/`down`, `make test`, `typecheck` y `test:run` en verde.

## Notas de implementación (2026-10-20)

- **Migración `00006`**: al ejecutarse antes que `notifications`/`push`, toma el siguiente número libre; las specs futuras se renumeran a `00007_notifications.sql` y `00008_device_tokens.sql` para no dejar huecos que goose se saltaría.
- **Idempotencia en el store**: `ON CONFLICT (session_id, client_answer_id) DO NOTHING`; `RowsAffected == 0` → `created=false`. Los `NULL` de clientes antiguos no chocan (semántica de Postgres). Verificado contra Postgres real en transacción con rollback.
- **Reintento sin efectos**: el service devuelve la fila guardada (con la explicación de *su* pregunta, no de la del payload) y no llama a `progressService.Answer`; el `Reward` no se repite.
- **Mobile**: el `clientAnswerId` se genera con `expo-crypto` al cargar cada pregunta y vive en un `useRef`; `submittingRef` bloquea el doble toque de forma síncrona. El error se pinta bajo `QuestionCard`, que no se desmonta, así que la selección se conserva.
- **Dependencia**: `expo-crypto@~57.0.4` (oficial, MIT, sin prebuild) instalado con `expo install`; gate de `.agents/mobile/dependencies.md` documentado.
- Tests: backend 5 subtests nuevos; mobile **111 tests / 25 archivos**; `make check`/`test`, `typecheck`, `test:run` y `lint` (0 errores) en verde.
- ⚠️ Queda la prueba manual en **modo avión** con dispositivo/emulador (misma pasada pendiente de MPROG-4).

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Reintento automático (1 intento) o solo manual? | **Manual**: el usuario ve qué pasa. Automático podría amplificar un fallo del servidor. |
| OQ-2 | ¿`clientAnswerId` obligatorio en clientes nuevos? | No: opcional (`omitempty`) para no romper integraciones ni datos históricos. |
| OQ-3 | ¿Reintentar también errores 4xx? | Sí, pero mostrando el detalle; la idempotencia lo hace seguro. Un 409 de sesión completada se mostrará tal cual. |
