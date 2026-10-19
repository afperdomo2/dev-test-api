# SPEC — review-mastery

| Campo | Valor |
|---|---|
| Module id | `review-mastery` |
| IDs originales | B1.1–B1.8 |
| Estado | `ready` |
| Depende de | `progress` (módulo existente), `sessions` (dato existente) |
| Consumido por | `insights` (filtra `isWeak`), `badges` (`topics_mastered_5`) |
| Migración goose | **Ninguna** — funcionalidad de solo lectura |

## Objective

Que el usuario vea su **dominio por tema** (accuracy real, no XP) y detecte de un vistazo en qué temas es débil, para saber qué repasar.

**Historias**:
- Como usuario, en `ProgressPage` abro la pestaña **Temas** y veo una barra de dominio por tema, ordenada con los temas débiles primero.
- Como usuario, un tema con accuracy bajo y suficiente historial aparece marcado con el badge **"Débil"**.
- Como usuario nuevo (sin respuestas), la pestaña muestra un estado vacío, no errores.

**Éxito**: el endpoint devuelve, por tema con respuestas, `attempts`, `correctCount`, `accuracy` e `isWeak`, con una **sola definición** de "débil" en todo el sistema (D-08).

### Aclaración de nombres (resuelve hallazgo D)

"Mastery" aquí significa **accuracy agregado por tema**. Es distinto del `UserQuestionProgress.IsMastered` del SRS (`progress/service.go:149`), que es por pregunta. Por eso:

- método de store: `FindTopicAccuracy` (no `FindTopicMastery`, que colisionaba conceptualmente),
- campo de respuesta: `isWeak` (no `isMastered`),
- el DTO **no** expone `isMastered`.

## Tech stack

Go + Gin + GORM + Postgres (backend) · Vue 3 + Vuetify + TanStack Query (web) · TypeScript puro (`@devtest/shared`). Sin dependencias nuevas.

## Commands

```bash
# Backend
make check          # make fmt + make vet
make test           # go test -short -count=1 (store mockeado, sin DB)
make swagger        # regenera docs/ tras las anotaciones nuevas
make build

# Web + shared
make fe-format      # prettier (primero, para que el check no falle)
make fe-check       # vue-tsc --build
make fe-lint        # eslint --fix
make fe-test        # vitest
pnpm --filter @devtest/shared typecheck
```

## Project structure

```
internal/models/session.go    # + consts SessionStatusInProgress/Completed/Cancelled (N-05)
internal/modules/progress/
├── store.go        # + FindTopicAccuracy + TopicAccuracyRow
├── service.go      # + TopicMastery en la interface + consts de umbral
├── response.go     # + TopicMasteryItem
├── handler.go      # + TopicMastery (anotaciones swagger)
├── routes.go       # + GET /progress/topics-mastery
└── service_test.go # + casos

packages/shared/src/types/progress.types.ts   # + TopicMasteryItem
packages/shared/src/index.ts                  # barrel (si re-exporta explícito)
apps/web/src/types/progress.types.ts          # re-export fino
apps/web/src/api/services/progress.service.ts # + getTopicMastery
apps/web/src/queries/progress.queries.ts      # + topicMasteryOptions
apps/web/src/features/progress/components/TopicMasteryList.vue  # NUEVO
apps/web/src/features/progress/pages/ProgressPage.vue           # 4ª pestaña
```

## Contrato

### Store (`store.go`)

```go
// Antes de tocar la query: D-21 exige que los valores de estado de sesión dejen de
// ser literales mágicos. Se añaden constantes en internal/models/session.go
// (N-05) y se usan aquí y en las agregaciones de insights.
type TopicAccuracyRow struct {
	TopicID      uuid.UUID `gorm:"column:topic_id"`
	TopicName    string    `gorm:"column:topic_name"`
	Attempts     int       `gorm:"column:attempts"`
	CorrectCount int       `gorm:"column:correct_count"`
}

// FindTopicAccuracy agrega las respuestas del usuario por tema.
// Una pregunta con N temas cuenta para los N temas.
FindTopicAccuracy(userID uuid.UUID) ([]TopicAccuracyRow, error)
```

```go
func (s *gormStore) FindTopicAccuracy(userID uuid.UUID) ([]TopicAccuracyRow, error) {
	var rows []TopicAccuracyRow
	err := s.db.
		Table("session_answers AS sa").
		Select(`t.id AS topic_id, t.name AS topic_name,
		        COUNT(sa.id) AS attempts,
		        COUNT(sa.id) FILTER (WHERE sa.is_correct) AS correct_count`).
		Joins("JOIN sessions AS s ON s.id = sa.session_id AND s.status = ?", models.SessionStatusCompleted).
		Joins("JOIN question_topics AS qt ON qt.question_id = sa.question_id").
		Joins("JOIN topics AS t ON t.id = qt.topic_id").
		Where("sa.user_id = ?", userID).
		Group("t.id, t.name").
		Scan(&rows).Error
	return rows, err
}
```

> Solo sesiones `completed` (D-10). Los temas sin respuestas **no aparecen**: el `JOIN` los excluye (no hace falta `LEFT JOIN`).

### Service (`service.go`)

```go
const (
	weakAccuracyThreshold = 0.6 // D-08: accuracy < 60%
	minAttempts           = 5   // D-08: al menos 5 respuestas
)

type Service interface {
	// ...existentes...
	TopicMastery(userID uuid.UUID) ([]TopicMasteryItem, error)
}

func (s *progressService) TopicMastery(userID uuid.UUID) ([]TopicMasteryItem, error) {
	rows, err := s.store.FindTopicAccuracy(userID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el dominio por tema", "")
	}

	items := make([]TopicMasteryItem, 0, len(rows))
	for _, r := range rows {
		accuracy := 0.0
		if r.Attempts > 0 {
			accuracy = math.Round(float64(r.CorrectCount)/float64(r.Attempts)*10000) / 10000
		}
		items = append(items, TopicMasteryItem{
			TopicID:      r.TopicID,
			TopicName:    r.TopicName,
			Attempts:     r.Attempts,
			CorrectCount: r.CorrectCount,
			Accuracy:     accuracy,
			IsWeak:       accuracy < weakAccuracyThreshold && r.Attempts >= minAttempts,
		})
	}

	// Orden determinista: débiles primero, luego accuracy ascendente, luego nombre.
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].IsWeak != items[j].IsWeak {
			return items[i].IsWeak
		}
		if items[i].Accuracy != items[j].Accuracy {
			return items[i].Accuracy < items[j].Accuracy
		}
		return items[i].TopicName < items[j].TopicName
	})

	return items, nil
}
```

`accuracy` se redondea a 4 decimales para que los tests sean estables y el JSON no arrastre ruido de coma flotante.

### Response (`response.go`)

```go
type TopicMasteryItem struct {
	TopicID      uuid.UUID `json:"topicId"`
	TopicName    string    `json:"topicName"`
	Attempts     int       `json:"attempts"`
	CorrectCount int       `json:"correctCount"`
	Accuracy     float64   `json:"accuracy"`
	IsWeak       bool      `json:"isWeak"`
}
```

> **Desviación de B1.3**: no se crea `TopicMasteryResponse` (wrapper). `response.Success(c, status, data)` ya envuelve en `{ "data": ... }` (`pkg/response/response.go:17`); un wrapper propio produciría `{"data":{"topics":[...]}}`.

### Endpoint

```
GET /api/v1/progress/topics-mastery     (protected — requiere JWT)
200 { "data": [ TopicMasteryItem, ... ] }
401 apierr.APIError
```

Sin paginación (D-19) y sin query params. Handler sigue el patrón del módulo (`getUserID(c)`, `e.Instance = c.Request.URL.Path`, `response.Problem`).

```go
// @Summary      Dominio por tema
// @Description  Accuracy agregado por tema de las sesiones finalizadas + detección de temas débiles
// @Tags         progress
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   TopicMasteryItem
// @Failure      401  {object}  apierr.APIError
// @Router       /api/v1/progress/topics-mastery [get]
```

### Shared + Web

```ts
// packages/shared/src/types/progress.types.ts
export interface TopicMasteryItem {
  topicId: string
  topicName: string
  attempts: number
  correctCount: number
  accuracy: number
  isWeak: boolean
}
```

```ts
// apps/web/src/api/services/progress.service.ts
export async function getTopicMastery(): Promise<TopicMasteryItem[]> {
  const res = await apiClient.get<TopicMasteryItem[]>('/api/v1/progress/topics-mastery')
  return res.data
}
```

```ts
// apps/web/src/queries/progress.queries.ts
export function topicMasteryOptions() {
  return queryOptions({
    queryKey: ['progress', 'topics-mastery'],
    queryFn: () => progressService.getTopicMastery(),
    staleTime: 30 * 1000,
  })
}
```

**UI**: `ProgressPage.vue` extiende `activeTab` a `'upcoming' | 'saved' | 'gamification' | 'topics'` y añade la pestaña **Temas** con `TopicMasteryList.vue` (presentacional: recibe `items` por props; **no emite eventos** — no hay acción definida para un tema). Nunca importa queries. Cada fila: nombre del tema, `v-progress-linear` con `accuracy * 100`, texto `correctCount/attempts`, y `v-chip` "Débil" cuando `isWeak`. Estado vacío con `v-alert` cuando `items.length === 0`.

**Accesibilidad (skill `ui-ux-pro-max`)**: el chip "Débil" no puede ser solo color — lleva texto; `aria-label` en la barra con el porcentaje; sin emojis como icono estructural.

## Code style

Convenciones obligatorias (detalle en `specs/DECISIONS.md` y `.agents/backend/architecture.md`):

- Store devuelve errores GORM crudos; el **service** los convierte a `*apierr.APIError`.
- Service fija `instance` a `""`; el **handler** lo rellena con `c.Request.URL.Path`.
- Nombres: `TopicMastery` / `TopicMasteryItem` / `FindTopicAccuracy` / `topicMasteryOptions` (camelCase en JSON y TS, PascalCase en Go exportado).
- Logs con emoji solo si se añade log (no es necesario aquí).

## Testing strategy

| Nivel | Archivo | Qué cubre |
|---|---|---|
| Unit service | `internal/modules/progress/service_test.go` | tabla: sin filas → lista vacía; redondeo de accuracy; `isWeak` en el borde exacto (`0.59/5` débil, `0.60` **no** débil, `0.4/4` **no** débil por attempts); orden determinista; error de store → `apierr` interno |
| Unit plain | — | El SQL se valida contra la DB real de forma manual (no hay test de integración en el repo) |
| Web service | `apps/web/src/api/services/progress.service.test.ts` | `getTopicMastery` llama la URL correcta y desempaqueta `{ data }` |
| Web query | `apps/web/src/queries/progress.queries.test.ts` | `topicMasteryOptions` con `queryKey` estable y `queryFn` correcto |
| Web page | `apps/web/src/features/progress/pages/ProgressPage.test.ts` (si existe; si no, se crea) | pestaña Temas renderiza items, badge Débil y estado vacío |

Store mockeado (patrón `mockProgressStore` ya existente en `service_test.go`). Sin DB en tests.

## Boundaries

- **Always**: `make check` + `make test` tras el backend; `make fe-format` → `make fe-check` → `make fe-lint` + `make fe-test` tras el frontend; `make swagger` tras tocar anotaciones; mapear los errores GORM a `apierr`.
- **Ask first**: añadir índices (p. ej. `question_topics(topic_id)`) si el agregado resulta lento en producción — requiere migración goose nueva; añadir campos al DTO más allá de los acordados.
- **Never**: `AutoMigrate`; editar `docs/` a mano; duplicar la definición de "débil" fuera de este módulo; paginar este endpoint (D-19); tocar mobile; exponer `isMastered` del SRS en este DTO.

## Success Criteria

1. `GET /api/v1/progress/topics-mastery` devuelve `{ "data": [TopicMasteryItem] }` con `200` para un usuario con respuestas, y `{ "data": [] }` para un usuario nuevo.
2. Un tema con 4 respuestas y 0 aciertos **no** es `isWeak`; con 5 respuestas y 2 aciertos (0.4) **sí**; con 5 respuestas y 3 aciertos (0.6) **no**.
3. La lista siempre ordena: débiles primero, luego accuracy ascendente, luego nombre.
4. `accuracy` es un número con 4 decimales máximo.
5. `docs/swagger.yaml` incluye el path nuevo tras `make swagger`.
6. La pestaña **Temas** en `ProgressPage` muestra barras + badge "Débil" + estado vacío, y `make fe-check` + `make fe-lint` + `make fe-test` pasan.
7. `pnpm --filter @devtest/shared typecheck` pasa.
8. `progress.Service` expone `TopicMastery` de forma que `insights` pueda consumirlo sin redefinir "débil".
9. Los literales `"completed"` / `"in_progress"` no se amplían: `internal/models/session.go` define las constantes y el código nuevo las usa (N-05, D-21).

## Open Questions

Ninguna bloqueante. Pendiente solo al implementar: si el agregado necesita índice en `question_topics` (ver Boundaries → Ask first).
