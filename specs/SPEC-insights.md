# SPEC — insights

| Campo | Valor |
|---|---|
| Module id | `insights` |
| IDs originales | B3.1–B3.8 |
| Estado | `ready` |
| Depende de | `review-mastery` (temas débiles), `gamification` (XP, racha), `sessions` |
| Migración goose | **Ninguna** — solo lectura |
| Dependencia nueva | **Ninguna** (sin librería de gráficas) |

## Objective

Que el usuario vea un **informe de su actividad** de la última semana o mes en una sola pantalla: cuánto practicó, con qué acierto, en qué temas, cuánto XP ganó y cómo se movió su racha.

**Historias**:
- Como usuario, entro en **Informes**, alterno semana/mes y veo los KPIs del periodo (sesiones, preguntas, acierto, tiempo, XP, racha).
- Como usuario, veo mis temas más practicados y, si los tengo, mis **temas débiles** — los mismos que marca la pestaña Temas de Progreso, no una lista distinta.
- Como usuario, veo la actividad diaria del periodo como barras.

**Éxito**: un solo endpoint devuelve el resumen completo, y `weakTopics` es **exactamente** lo que devuelve `review-mastery` filtrado por `isWeak` (D-08) — cero duplicación de la heurística.

### Asimetría deliberada `topTopics` vs `weakTopics` (no "arreglar")

- `topTopics` → métrica **del periodo**: los temas con más respuestas en la ventana. Es una medida de *actividad*, no de dominio.
- `weakTopics` → métrica **histórica**: los temas débiles según la definición única de `review-mastery` (all-time, `accuracy < 0.6 && attempts >= 5`).

Si ambos se calcularan en la ventana, habría **dos** definiciones de "débil" (viola D-08) y un tema débil con pocos intentos recientes desaparecería del informe. Queda documentado aquí para que no se "corrija".

## Tech stack

Go + Gin + GORM + Postgres · Vue 3 + Vuetify + TanStack Query · `@devtest/shared`.

## Commands

```bash
make check && make test && make swagger && make build
make fe-format && make fe-check && make fe-lint && make fe-test
pnpm --filter @devtest/shared typecheck
```

## Project structure

```
internal/modules/gamification/service.go      # + XpForDifficulty() / SessionBonusXp() exportados (fuente única del XP)
internal/modules/insights/                    # NUEVO
├── request.go      # SummaryParams
├── response.go     # InsightSummaryResponse, TopicInsightItem, DailyActivityItem
├── store.go        # Store + gormStore (agregaciones)
├── service.go      # Service + Summary()
├── handler.go      # Summary + swagger
├── routes.go       # RegisterRoutes
└── service_test.go

internal/server/server.go                     # wiring (después de progress y gamification)

packages/shared/src/types/insight.types.ts    # NUEVO
packages/shared/src/index.ts                  # barrel
apps/web/src/types/insight.types.ts           # re-export fino (NUEVO)
apps/web/src/api/services/insights.service.ts # NUEVO
apps/web/src/queries/insights.queries.ts      # NUEVO
apps/web/src/features/insights/pages/InsightsPage.vue       # NUEVO
apps/web/src/features/insights/components/*.vue             # NUEVO (KPIs, barras)
apps/web/src/layouts/DefaultLayout.vue        # nav Informes + título de sección
apps/web/src/router/routes.ts                 # + /insights
```

## Contrato

### Endpoint

```
GET /api/v1/insights/summary?period=week|month     (protected — requiere JWT)
```

- `period` opcional, `binding:"omitempty,oneof=week month"`, default **`week`**.
- `week` = últimos **7** días (hoy incluido); `month` = últimos **30** días. **Ventana móvil**, no semana/mes natural: evita bordes de calendario y no necesita configurar el primer día de la semana.
- TZ: la del servidor para el cálculo de días; el detalle se serializa en RFC 3339.

### Response

```go
type TopicInsightItem struct {
	TopicID   uuid.UUID `json:"topicId"`
	TopicName string    `json:"topicName"`
	Attempts  int       `json:"attempts"`
	Accuracy  float64   `json:"accuracy"`
}

type DailyActivityItem struct {
	Date              string `json:"date"` // YYYY-MM-DD (hora local del servidor)
	QuestionsAnswered int    `json:"questionsAnswered"`
}

type InsightSummaryResponse struct {
	Period            string                      `json:"period"`
	From              time.Time                   `json:"from"`
	To                time.Time                   `json:"to"`
	Sessions          int                         `json:"sessions"`
	QuestionsAnswered int                         `json:"questionsAnswered"`
	Accuracy          float64                     `json:"accuracy"`
	TimeSpentMs       int64                       `json:"timeSpentMs"`
	TopTopics         []TopicInsightItem          `json:"topTopics"`
	WeakTopics        []progress.TopicMasteryItem `json:"weakTopics"`
	XpGained          int                         `json:"xpGained"`
	StreakDays        int                         `json:"streakDays"`
	DailyActivity     []DailyActivityItem         `json:"dailyActivity"`
}
```

`WeakTopics` **reutiliza** `progress.TopicMasteryItem` — hay precedente de DTOs cruzados en el repo (`progress/response.go:6-7` importa `gamification` y `questions`). El handler responde `response.Success(c, http.StatusOK, resp)` → `{ "data": {...} }`.

### Definiciones (fijadas, no negociables al implementar)

| Campo | Definición |
|---|---|
| `Sessions` | `sessions` del usuario con `status = models.SessionStatusCompleted` y `created_at` en la ventana |
| `QuestionsAnswered` | `COUNT(session_answers)` del usuario, `created_at` en la ventana, cuya sesión está `completed` |
| `Accuracy` | `correct / QuestionsAnswered`, redondeado a 4 decimales; `0` si no hay respuestas |
| `TimeSpentMs` | `SUM(session_answers.response_time_ms)` (N-02), no `finished_at - started_at` (evita tiempo muerto) |
| `TopTopics` | top 5 por `attempts` en la ventana, `JOIN question_topics`; desempate por `accuracy` DESC y nombre |
| `WeakTopics` | `progress.Service.TopicMastery(userID)` filtrado por `IsWeak` (D-08) |
| `XpGained` | `Σ XP por dificultad de las respuestas correctas de la ventana + sessionBonusXp × sesiones completed de la ventana` |
| `StreakDays` | `gamification.Service.GetOverview(userID).Streak.Current` |
| `DailyActivity` | un punto por día de la ventana, **con ceros**, en orden ascendente |

La ventana filtra por `created_at` (no `finished_at`) en ambas tablas: es un valor no nulo y evita asimetrías de medianoche.

### Fuente única del XP (evita duplicar la tabla de XP)

`gamification/service.go:20-24` tiene `xpByDifficulty` **sin exportar** y `sessionBonusXp=20` (`:14`). Duplicar esos números en SQL sería una segunda fuente de verdad. La spec exporta dos helpers en `gamification`:

```go
// gamification/service.go
func XpForDifficulty(difficulty string) int
func SessionBonusXp() int
```

Y el store de `insights` devuelve los conteos crudos por dificultad:

```go
// insights/store.go
CountCorrectByDifficulty(userID uuid.UUID, from time.Time) (map[string]int, error)
CountCompletedSessions(userID uuid.UUID, from time.Time) (int, error)
```

```go
// insights/service.go
xp := 0
for difficulty, n := range correctByDifficulty {
	xp += n * gamification.XpForDifficulty(difficulty)
}
xp += sessions * gamification.SessionBonusXp()
```

### Store (agregaciones)

```go
type Store interface {
	CountSessions(userID uuid.UUID, from time.Time) (int, error)
	CountAnswers(userID uuid.UUID, from time.Time) (int, int, error) // total, correctas
	SumResponseTimeMs(userID uuid.UUID, from time.Time) (int64, error)
	TopTopics(userID uuid.UUID, from time.Time, limit int) ([]TopicInsightItem, error)
	DailyActivity(userID uuid.UUID, from, to time.Time) (map[string]int, error)
	CountCorrectByDifficulty(userID uuid.UUID, from time.Time) (map[string]int, error)
	CountCompletedSessions(userID uuid.UUID, from time.Time) (int, error)
}
```

`DailyActivity` devuelve el mapa de días con datos; el **service** rellena con ceros los días sin actividad hasta completar la ventana (no hacerlo en SQL con `generate_series`: mezcla lógica de presentación con la query).

### Service

```go
func (s *insightsService) Summary(userID uuid.UUID, period string) (*InsightSummaryResponse, error)
```

- `progress.TopicMastery` y `gamification.GetOverview` son **independientes** entre sí y de las agregaciones → van en paralelo con `sync.WaitGroup` + `sync.Mutex` (D-13, patrón de `sessions/service.go:128-129`), **no** `errgroup`.
- Cualquier error de store → `*apierr.APIError` con mensaje en español.
- Si no hay actividad: respuesta válida con ceros, `topTopics: []`, `weakTopics` según mastery, `dailyActivity` con la ventana completa en ceros.

### Shared + Web

```ts
// packages/shared/src/types/insight.types.ts
export type InsightPeriod = 'week' | 'month'

export interface TopicInsightItem {
  topicId: string
  topicName: string
  attempts: number
  accuracy: number
}

export interface DailyActivityItem {
  date: string
  questionsAnswered: number
}

export interface InsightSummary {
  period: InsightPeriod
  from: string
  to: string
  sessions: number
  questionsAnswered: number
  accuracy: number
  timeSpentMs: number
  topTopics: Array<TopicInsightItem>
  weakTopics: Array<TopicMasteryItem> // importado de progress.types (no redefinir)
  xpGained: number
  streakDays: number
  dailyActivity: Array<DailyActivityItem>
}

export const PERIOD_OPTIONS: Array<{ title: string; value: InsightPeriod }> = [
  { title: 'Semana', value: 'week' },
  { title: 'Mes', value: 'month' },
]
```

> `TopicMasteryItem` se **importa** de `progress.types.ts`, no se re-exporta por el barrel desde dos sitios (regla anti-duplicación de `.agents/shared/architecture.md`).

```ts
// apps/web/src/api/services/insights.service.ts
export async function getInsightSummary(period: InsightPeriod): Promise<InsightSummary> {
  const res = await apiClient.get<InsightSummary>('/api/v1/insights/summary', { params: { period } })
  return res.data
}

// apps/web/src/queries/insights.queries.ts
export function insightSummaryOptions(period: () => InsightPeriod) {
  return queryOptions({
    queryKey: ['insights', 'summary', period],
    queryFn: () => insightsService.getInsightSummary(period()),
    staleTime: 60 * 1000,
  })
}
```

**UI** (`InsightsPage.vue`):
- Selector semana/mes (`v-btn-toggle` con `PERIOD_OPTIONS`).
- Grid de KPI cards: sesiones, preguntas, acierto (`formatScore` de shared), tiempo (`formatDateTime`/helper local de duración), XP ganado, racha.
- Dos listas con barras (`v-progress-linear`): **Más practicados** y **Temas débiles** (badge "Débil" con texto, no solo color). Vacío → `v-alert` informativo.
- **Actividad diaria** con barras CSS/`v-progress-linear` horizontales. **Sin librería de gráficas**: no se añade dependencia.
- Accesibilidad: `aria-label` en las barras con el valor; el tiempo se muestra en texto legible (`12 min`), no solo milisegundos.

**Nav (B3.7)**: ⚠️ **la sección "Estudio" no existe** en `DefaultLayout.vue` — el grupo de Sesiones/Progreso (`:26-31`) está **sin título**. Se añade `title: 'Estudio'` a esa sección y dentro el item `{ title: 'Informes', to: '/insights', icon: 'mdi-file-chart', userOnly: true }`. Así el sidebar queda coherente (Contenido y Administración ya están tituladas).

**Ruta**: `{ path: 'insights', name: 'Insights', component: () => import('@/features/insights/pages/InsightsPage.vue'), meta: { requiresAuth: true, requiresNotAdmin: true } }` — alineado con el `userOnly` del nav (los admins no ven Sesiones/Progreso).

## Code style

- Módulo con la anatomía estándar (`request/response/store/service/handler/routes`); **sin `mapper.go`** si no hay transformación (los DTOs se construyen en el service).
- Errores: `apierr.ErrValidation` para `period` inválido (aunque el `binding` ya lo corta en el handler con `response.ValidationError`).
- Handler fija `e.Instance = c.Request.URL.Path`.
- Redondeo de `accuracy` a 4 decimales (igual que `review-mastery`) para tests estables.
- Mobile fuera de alcance (D-18).

## Testing strategy

| Nivel | Archivo | Qué cubre |
|---|---|---|
| Unit service | `internal/modules/insights/service_test.go` | tabla por `period` (7 vs 30 días); sin actividad → ceros con `dailyActivity` relleno; `xpGained` = correctas×dificultad + 20×sesiones; `weakTopics` filtrado desde el mock de `progress`; error de store → `apierr`; `accuracy` con 0 respuestas no divide por cero |
| Unit plain | `gamification` | `XpForDifficulty` para las 3 dificultades y una desconocida (→ `0`), `SessionBonusXp() == 20` |
| Unit plain | — | Las agregaciones SQL se verifican contra la DB real en transacción con rollback (nota en la tarea) |
| Web service | `insights.service.test.ts` | URL + `params.period` + desempaquetado |
| Web query | `insights.queries.test.ts` | queryKey depende del periodo |
| Web page | `InsightsPage.test.ts` | cambio de periodo re-consulta; KPIs; listas vacías; barras de actividad |

Mock de `progress.Service` y de `gamification.Service` (interfaces ya inyectables). Sin DB en tests.

## Boundaries

- **Always**: `make check` + `make test` + `make swagger`; `make fe-*` + `make fe-test`; reutilizar `progress.Service.TopicMastery` y los helpers de XP de `gamification` en vez de recalcular; `sync.WaitGroup` para las llamadas independientes.
- **Ask first**: añadir librería de gráficas; cambiar la ventana a semana/mes natural o configurar el primer día de la semana; añadir export CSV/PDF del informe; añadir índices para las agregaciones (requiere migración).
- **Never**: `errgroup` (D-13); reimplementar "tema débil" (D-08); duplicar `xpByDifficulty`/`sessionBonusXp` (N-01); usar `finished_at - started_at` para el tiempo (N-02); `AutoMigrate`; editar `docs/` a mano; tocar mobile.

## Success Criteria

1. `GET /api/v1/insights/summary` sin `period` responde con `period: "week"`; con `period=month`, 30 días de ventana.
2. `period=year` → `422` con el formato RFC 9457 del repo.
3. Con actividad conocida sembrada, los 11 campos cuadran con el cálculo manual (incluido `timeSpentMs` = suma de `responseTimeMs`).
4. Usuario sin actividad: `200`, ceros, `dailyActivity` con exactamente 7/30 elementos, `topTopics: []`.
5. `weakTopics` coincide **item por item** con lo que devuelve `/progress/topics-mastery` filtrado por `isWeak` (mismo umbral, mismos valores).
6. `xpGained` = `5/10/15` por dificultad de las correctas + `20` por sesión completed de la ventana.
7. `docs/swagger.yaml` incluye `GET /api/v1/insights/summary` con el parámetro `period`.
8. La página Informes aparece en el nav (sección ahora titulada **Estudio**) solo para no-admin, la ruta `/insights` está protegida con `requiresNotAdmin`, y `make fe-check` + `make fe-lint` + `make fe-test` pasan.
9. `pnpm --filter @devtest/shared typecheck` pasa.
10. Ninguna dependencia nueva en `go.mod` ni en los `package.json`.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | `xpGained` derivado ¿es aceptable frente a un ledger `xp_events`? | **Derivar** (opción recomendada). Exacto hoy; el riesgo es que si `sessionBonusXp` o `xpByDifficulty` cambian, los periodos **anteriores** se recalculan con las reglas nuevas. Un ledger sería un cambio de esquema fuera del alcance de B3. Se documenta el límite en el código con un comentario. |
| OQ-2 | ¿Titular la sección del nav como "Estudio" o dejar el item en la sección sin título? | **Titular "Estudio"**: es lo que pedía B3.7 y alinea el sidebar con Contenido/Administración. |
| OQ-3 | ¿`topTopics` top 3 o top 5? | **5**, con `limit` como constante del service (fácil de ajustar). |
| OQ-4 | ¿Tiempo total del periodo en minutos/horas en la UI? | El backend devuelve `timeSpentMs`; la UI formatea (`Xh Ym` / `Y min`). No se añade un campo derivado a la API. |
