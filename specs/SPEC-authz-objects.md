# SPEC — authz-objects

| Campo | Valor |
|---|---|
| Module id | `authz-objects` |
| Origen | Revisión de seguridad 2026-10-21 (H1, H2) |
| Estado | `ready` |
| Depende de | — |
| Consumido por | — (interno) |
| Migración goose | Ninguna |
| Dependencia nueva | Ninguna |
| Cambio de contrato | Sí: códigos de error en `sessions` y `questions` (ver "Contract changes") |

## Objective

Hoy cualquier usuario autenticado puede operar sobre objetos de otro usuario
porque el handler recibe el `id` pero el service nunca compara el dueño con el
`sub` del JWT:

- `GET /api/v1/sessions/:id` — `sessions/handler.go:140` → `service.GetByID(id)`
- `GET /api/v1/sessions/:id/summary` — `handler.go:307` → `service.Summary(id)`
- `PUT /api/v1/sessions/:id/finish` — `handler.go:205` → `service.Finish(id)`
- `GET /api/v1/sessions/:id/next` — `handler.go:235` → `service.NextQuestion(id)`
- `POST /api/v1/sessions/:id/answer` — `handler.go:267` → `service.Answer(sessionID, userID, …)` sin comprobar `sess.UserID`
- `GET /api/v1/questions/:id` — `questions/handler.go:87` → `service.GetByID(id)`, sin visibilidad, y el mapper (`questions/mapper.go:43-52`) devuelve `isCorrect` + `explanation`

**Historias**

- Como usuario, no quiero que otro usuario pueda leer mis sesiones, mi feedback
  ni las respuestas correctas de mis preguntas privadas.
- Como usuario, no quiero que otro pueda finalizar mi sesión y otorgarme XP.
- Como administrador, quiero seguir viendo y gestionando todo.

**Éxito**

- Con un JWT válido de otro usuario, los 5 endpoints de sesión responden 404 y no
  filtran ni mutan nada.
- Con un JWT de otro usuario, `GET /questions/:id` de una pregunta privada
  responde 404; las públicas y las `ai_generated` siguen respondiendo 200.
- No cambia ningún caso legítimo (dueño, admin, preguntas públicas).

## Tech stack

Go + Gin + GORM + Postgres. Solo cambios de backend.

## Commands

```bash
make fmt && make vet          # o make check
make test                     # tests unitarios (mocks de Store)
make test-cover
make swagger                  # regenera docs/ por los códigos de error
make build
```

## Project structure

```
internal/modules/sessions/
  service.go            # + userID en GetByID/Summary/Finish/NextQuestion/Answer
  handler.go            # pasa userID desde getUserID(c)
  service_test.go       # ownership (tabla)
  handler_test.go       # HTTP 404 con JWT de otro usuario
internal/modules/questions/
  store.go              # + FindVisibleByID(id, userID, isAdmin)
  service.go            # GetByID(id, userID, isAdmin)
  handler.go            # usa getUserRoleAndID(c) (ya existe, handler.go:340)
  service_test.go       # visibilidad
  handler_test.go       # 404 para privada ajena
docs/                   # generado por make swagger
```

## Contrato

### Decisiones

| ID | Decisión | Motivo |
|---|---|---|
| D-AZ-1 | El fallo de autorización por objeto devuelve **404**, no 403 | No revela la existencia del recurso ajeno (evita oráculo de enumeración). 403 se reserva para permisos de rol/política. |
| D-AZ-2 | La comprobación vive en el **service**, no en el handler | Testeable con mocks; imposible de olvidar al añadir rutas. |
| D-AZ-3 | `sessions.Delete` cambia su 403 por **404** en el caso "no es tuya"; el 403 de la regla de 24 h se mantiene | Consistencia con D-AZ-1. La regla de 24 h sí es una política legítima. |
| D-AZ-4 | Las preguntas visibles = `is_public = true OR source = 'ai_generated' OR user_id = caller` (admin: todas) | Reutiliza exactamente la regla de `questions/store.go:50-52` (listado). Una sola definición. |

### Cambios de firma

```go
// internal/modules/sessions/service.go
GetByID(id, userID uuid.UUID) (*SessionResponse, error)
Summary(id, userID uuid.UUID) (*SessionSummaryResponse, error)
Finish(id, userID uuid.UUID) (*SessionResponse, error)
NextQuestion(id, userID uuid.UUID) (*NextQuestionResponse, error)
Answer(sessionID, userID uuid.UUID, input AnswerRequest) (*SessionAnswerResponse, error)

// internal/modules/questions/store.go
FindVisibleByID(id, userID uuid.UUID, isAdmin bool) (*models.Question, error)

// internal/modules/questions/service.go
GetByID(id, userID uuid.UUID, isAdmin bool) (*QuestionResponse, error)
```

Handler de sesiones (patrón repetido, con `getUserID` ya existente en
`sessions/handler.go:325`):

```go
userID, apiErr := getUserID(c)
if apiErr != nil {
    apiErr.Instance = c.Request.URL.Path
    response.Problem(c, apiErr)
    return
}

session, err := h.service.GetByID(id, userID)
```

Service de sesiones (patrón único, factorizado):

```go
func findOwnedSession(store Store, id, userID uuid.UUID) (*models.Session, error) {
    sess, err := store.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apierr.ErrNotFound("Sesion", "")
        }
        return nil, apierr.ErrInternal("Error al obtener la sesion", "")
    }
    if sess.UserID != userID {
        // D-AZ-1: no revelar que existe
        return nil, apierr.ErrNotFound("Sesion", "")
    }
    return sess, nil
}
```

Se usa en `GetByID`, `Summary`, `Finish`, `NextQuestion` y antes de la lógica de
idempotencia en `Answer` (`00006_answer_idempotency`): la comprobación de dueño
va **primero**, para que un tercero no pueda sondear la idempotencia de otro.

`Delete` mantiene su estructura actual (`service.go:425-445`) y solo cambia
`ErrForbidden` → `ErrNotFound` en la comparación de dueño.

Store de preguntas:

```go
func (s *gormStore) FindVisibleByID(id, userID uuid.UUID, isAdmin bool) (*models.Question, error) {
    q := s.db.Preload("Options").Preload("CodeChallenge").Preload("Topics")
    if isAdmin {
        q = q.Where("id = ?", id)
    } else {
        q = q.Where("id = ? AND (is_public = ? OR source = ? OR user_id = ?)",
            id, true, "ai_generated", userID)
    }
    var question models.Question
    if err := q.First(&question).Error; err != nil {
        return nil, err
    }
    return &question, nil
}
```

### Contract changes (Swagger)

| Endpoint | Antes | Después |
|---|---|---|
| `GET /sessions/{id}` | 401, 404 | 401, 404 (404 también si es ajena) |
| `GET /sessions/{id}/summary` | 404 | 404 |
| `PUT /sessions/{id}/finish` | 401, 404, 409 | 401, 404, 409 |
| `GET /sessions/{id}/next` | 401, 404, 409 | 401, 404, 409 |
| `POST /sessions/{id}/answer` | 400, 401, 404, 409 | 400, 401, 404, 409 |
| `DELETE /sessions/{id}` | 403 (dueño), 403 (24 h) | **404 (dueño)**, 403 (24 h) |
| `GET /questions/{id}` | 404 | 404 (privada ajena) |

`make swagger` regenera `docs/swagger.yaml|json|docs.go`; no se editan a mano.

### Fuera del contrato

- `NextQuestion` ya oculta `isCorrect` (`sessions/mapper.go:toNextQuestionItem`).
  No se toca.
- `ToSessionDetailResponse` sigue devolviendo `isCorrect` **después** de
  responder: es intencional (pantalla de revisión del dueño).

## Testing strategy

Tests de service con `Store` mockeado (`testify`) y tests de handler a nivel HTTP
con router + `middleware.Auth` + JWT firmado (mismo estilo que RM-4).

| # | Caso | Esperado |
|---|---|---|
| T1 | `GetByID(id, owner)` | 200 con detalle |
| T2 | `GetByID(id, otro)` | 404, sin `userID` ni respuestas en el body |
| T3 | `GetByID(id inexistente, cualquiera)` | 404 idéntico a T2 (mismo body) |
| T4 | `Finish(id, otro)` | 404 y **no** se llama `gamification.CompleteSession` (spy) |
| T5 | `Finish(id, owner)` | 200, `reward` presente |
| T6 | `Answer(id, otro, …)` | 404 y `store.FindAnswerBy…`/`Create` no invocados |
| T7 | `Answer(id, owner, repetido)` | idempotencia actual intacta (409/200 según `00006`) |
| T8 | `Summary(id, otro)` | 404 sin contadores |
| T9 | `NextQuestion(id, otro)` | 404 sin contenido de pregunta |
| T10 | `Delete(id, otro)` | 404 (cambio D-AZ-3) |
| T11 | `Delete(id, owner, >24 h)` | 403 (regla de política) |
| T12 | `Questions.GetByID(privada_ajena, user, false)` | 404 |
| T13 | `Questions.GetByID(pública, user, false)` | 200 |
| T14 | `Questions.GetByID(ai_generated, user, false)` | 200 |
| T15 | `Questions.GetByID(privada_ajena, admin, true)` | 200 |
| T16 | `Questions.List` con `userID` | sin cambios (regresión) |

Cobertura objetivo: los 5 métodos de sesión y `questions.GetByID` quedan al 100 %
de ramas (owner / no owner / not found).

## Boundaries

- **Always**: `make check` + `make test` + `make swagger` antes de cerrar;
  tests nuevos por cada rama de autorización; 404 con el mismo body que un
  recurso inexistente.
- **Ask first**: cambiar `Delete` a 403 si prefieres mantener el contrato previo;
  tocar la regla de visibilidad (p. ej. permitir preguntas de temas compartidos).
- **Never**: comprobar propiedad en el handler, usar `fmt.Sprintf` para el WHERE,
  devolver 403/404 con detalles distintos que permitan distinguir "no existe" de
  "es de otro", añadir un `Preload` que reintroduzca `UserID` ajeno en la
  respuesta.

## Success criteria

1. Los 5 endpoints de sesión devuelven 404 con el JWT de otro usuario y no
   producen escrituras (verificado en T4/T6 con spy).
2. `GET /questions/:id` respeta la misma visibilidad que el listado (T12–T15).
3. `make test` verde y `make swagger` sin diffs inesperados (solo los códigos de
   error documentados).
4. Los tests E2E existentes siguen verdes (el dueño no cambia de comportamiento).
5. No hay migraciones ni cambios de esquema.

## Open questions

- ¿Prefieres mantener 403 en `Delete` para no cambiar el contrato publicado?
  (propuesto: 404 por D-AZ-1 y D-AZ-3).
- ¿Un admin debe poder finalizar la sesión de otro usuario para soporte?
  (propuesto: no; el admin gestiona usuarios y contenido, no sesiones).
