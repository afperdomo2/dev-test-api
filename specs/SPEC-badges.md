# SPEC — badges (opt-in)

| Campo | Valor |
|---|---|
| Module id | `badges` |
| IDs originales | B4.1–B4.5 |
| Estado | **`opt-in`** — no se implementa sin un go-ahead explícito del usuario |
| Depende de | `gamification`, `review-mastery`, `insights` |
| Migración goose | **Ninguna** — catálogo `const` + evaluación derivada (D-16) |
| Dependencia nueva | **Ninguna** |

## Objective

Dar reconocimiento a hitos de uso (primera sesión, rachas, sesiones perfectas, volumen de respuestas, temas dominados, rapidez) en una página de logros ganados/bloqueados.

**Historias**:
- Como usuario, abro **Logros** y veo los 8 logros, cuáles tengo y cuáles no, con mi progreso hacia los bloqueados.
- Como usuario con racha de 7 días, veo el logro `streak_7` desbloqueado.

**Éxito**: el catálogo vive en código (D-16), no hay tabla nueva, no hay `earnedAt` prometido, y el endpoint se resuelve con un número acotado de agregaciones.

### Desviación de B4.2 (explícita)

`B4.2` pedía evaluar los criterios dentro de `RecordAnswer`/`RecordSessionFinished`. **No se hace**: sin tabla de logros ganados, evaluar en el camino caliente no persiste nada y solo añade latencia al responder y al terminar sesión. La evaluación vive en `gamification/badges.go` y se ejecuta **bajo demanda** en el endpoint (D-16). Es una consecuencia obligada de la decisión bloqueada "sin tabla `UserBadge`".

`B4.3` mencionaba "`user_badges` como tabla (decisión pendiente)" y se contradecía con la decisión bloqueada del propio `ROADMAP.md`; **gana D-16 / X-02: sin tabla**.

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
internal/modules/gamification/
├── catalog.go      # NUEVO: BadgeID + catálogo const (id, nombre, descripción, icono, color)
├── badges.go       # NUEVO: evaluación de criterios (agregaciones)
├── store.go        # + agregaciones necesarias
├── service.go      # + GetBadges() en la interface
├── response.go     # + BadgeItem / BadgesResponse
├── handler.go      # + Badges (swagger)
└── routes.go       # + GET /gamification/badges

packages/shared/src/types/gamification.types.ts   # + BadgeId, Badge, BadgesResponse, BADGES_*
apps/web/src/api/services/gamification.service.ts # + getBadges
apps/web/src/queries/gamification.queries.ts      # + badgesOptions
apps/web/src/features/gamification/pages/AchievementsPage.vue  # NUEVO
apps/web/src/router/routes.ts                     # + /achievements
apps/web/src/layouts/DefaultLayout.vue            # nav "Logros" (userOnly)
```

## Contrato

### Catálogo (`catalog.go`)

```go
type BadgeID string

const (
	BadgeFirstSession    BadgeID = "first_session"
	BadgeStreak3         BadgeID = "streak_3"
	BadgeStreak7         BadgeID = "streak_7"
	BadgeStreak30        BadgeID = "streak_30"
	BadgePerfectSession  BadgeID = "perfect_session"
	BadgeAnswers100      BadgeID = "answers_100"
	BadgeTopicsMastered5 BadgeID = "topics_mastered_5"
	BadgeFastSolver      BadgeID = "fast_solver"
)

// Catálogo inmutable: el orden es el que se muestra en la UI.
var badgeCatalog = []BadgeDefinition{
	{ID: BadgeFirstSession, Name: "Primera sesión", Description: "Completa tu primera sesión", Icon: "mdi-flag-checkered", Target: 1},
	{ID: BadgeStreak3, Name: "Racha de 3", Description: "3 días seguidos de estudio", Icon: "mdi-fire", Target: 3},
	{ID: BadgeStreak7, Name: "Racha de 7", Description: "7 días seguidos de estudio", Icon: "mdi-fire", Target: 7},
	{ID: BadgeStreak30, Name: "Racha de 30", Description: "30 días seguidos de estudio", Icon: "mdi-fire", Target: 30},
	{ID: BadgePerfectSession, Name: "Sesión perfecta", Description: "Termina una sesión con el 100% de acierto (mín. 5 preguntas)", Icon: "mdi-star", Target: 1},
	{ID: BadgeAnswers100, Name: "100 respuestas", Description: "Responde 100 preguntas", Icon: "mdi-check-all", Target: 100},
	{ID: BadgeTopicsMastered5, Name: "5 temas dominados", Description: "Domina 5 temas", Icon: "mdi-school", Target: 5},
	{ID: BadgeFastSolver, Name: "Resolvedor veloz", Description: "Sesión perfecta con menos de 10 s por pregunta", Icon: "mdi-lightning-bolt", Target: 1},
}
```

### Criterios (todos derivados de datos existentes)

| Badge | Criterio | Fuente |
|---|---|---|
| `first_session` | ≥ 1 sesión `completed` | `sessions` |
| `streak_3` / `streak_7` / `streak_30` | `UserStreak.Longest >= 3 / 7 / 30` | `user_streaks` |
| `perfect_session` | ≥ 1 sesión `completed` con **≥ 5** respuestas y **100 %** de acierto | `sessions` + `session_answers` |
| `answers_100` | `COUNT(session_answers)` ≥ 100 | `session_answers` |
| `topics_mastered_5` | ≥ 5 temas dominados (ver OQ-1) | `user_question_progress` / `TOPIC accuracy` |
| `fast_solver` | ≥ 1 sesión `completed` con ≥ 5 respuestas, 100 % de acierto y media de `response_time_ms` ≤ 10 000 | `session_answers` |

Se usa `Longest` (no `Current`): un logro es "lo has conseguido alguna vez", no "lo tienes hoy".

### Endpoint

```
GET /api/v1/gamification/badges     (protected — requiere JWT)
200 { "data": { "badges": [ BadgeItem, ... ], "earned": 3, "total": 8 } }
```

```go
type BadgeItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Earned      bool   `json:"earned"`
	Current     int    `json:"current"` // progreso hacia el objetivo
	Target      int    `json:"target"`
}

type BadgesResponse struct {
	Badges []BadgeItem `json:"badges"`
	Earned int         `json:"earned"`
	Total  int         `json:"total"`
}
```

⚠️ **No hay `earnedAt`**: sin persistencia no existe la fecha. La UI no puede prometer "conseguido el día X" ni ordenar por fecha de obtención. Documentado como límite de D-16, no como bug.

`Current` permite "12/100 respuestas" en los bloqueados. Para badges booleanos (`perfect_session`, `fast_solver`) `Current` es `0` o `Target`.

### Evaluación (`badges.go`)

- Agregaciones independientes → en paralelo con `sync.WaitGroup` + `sync.Mutex` (D-13).
- ≤ 6 queries por llamada; si alguna crece, se consolidan en una sola query con `FILTER`.
- `topics_mastered_5` reutiliza la agregación de `review-mastery` (`progress.Store.FindTopicAccuracy`) — no se reimplementa el cálculo por tema.

### Shared + Web

```ts
// packages/shared/src/types/gamification.types.ts
export type BadgeId =
  | 'first_session' | 'streak_3' | 'streak_7' | 'streak_30'
  | 'perfect_session' | 'answers_100' | 'topics_mastered_5' | 'fast_solver'

export interface Badge {
  id: BadgeId
  name: string
  description: string
  icon: string
  earned: boolean
  current: number
  target: number
}

export interface BadgesResponse {
  badges: Array<Badge>
  earned: number
  total: number
}
```

**UI** (`AchievementsPage.vue`): grid de cards; ganado = icono `mdi` + color + texto "Conseguido"; bloqueado = atenuado + `v-progress-linear` `current/target` + texto `12/100`. Estado de ganado/bloqueado **nunca solo por color** (texto + icono). Botón de refresco manual. `staleTime` alto (5 min) porque el estado cambia poco.

**Ruta**: `/achievements` con `meta: { requiresAuth: true, requiresNotAdmin: true }`. **Nav**: item `{ title: 'Logros', to: '/achievements', icon: 'mdi-trophy-outline', userOnly: true }` en la sección **Estudio** (la misma que titula `SPEC-insights.md`; si `insights` no está implementado aún, se añade el título igualmente).

## Code style

- `catalog.go`: datos como `var` de solo lectura + `const` para los ids. Nada de lógica.
- `badges.go`: funciones puras de evaluación que reciben los agregados ya cargados — testeables sin mocks de store complejos.
- Handler/service/store siguen la anatomía estándar; errores a `*apierr.APIError`.
- Mobile fuera de alcance (D-18).

## Testing strategy

| Nivel | Archivo | Qué cubre |
|---|---|---|
| Unit puro | `internal/modules/gamification/badges_test.go` | tabla por badge: umbral exacto en el borde (`Longest=2` → `streak_3` false, `=3` → true), `<>` de `answers_100` (99/100), `perfect_session` con 4 respuestas (no cuenta) y con 5 al 100 % (sí), `fast_solver` con media 10 001 ms (no) y 9 999 ms (sí) |
| Unit service | `internal/modules/gamification/service_test.go` | `GetBadges` compone el catálogo completo, cuenta `earned`/`total`, orden estable |
| Unit plain | — | Contraste de los agregados contra la DB real en transacción con rollback (nota en la tarea) |
| Web service/query | `gamification.service.test.ts` / `gamification.queries.test.ts` | `getBadges` + `badgesOptions` |
| Web page | `AchievementsPage.test.ts` | grid con ganados/bloqueados, progreso, estado vacío |

## Boundaries

- **Always**: catálogo en código (D-16); reutilizar agregaciones existentes; `make check` + `make test` + `make swagger`; `make fe-*` + `make fe-test`.
- **Ask first**: adoptar el módulo (es `opt-in`); cualquier persistencia de logros ganados (implica migración y **supersede D-16**, requiere ADR explícito); añadir un badge nuevo; cambiar los umbrales.
- **Never**: crear la tabla `user_badges` sin superseder D-16 (X-02); evaluar en `RecordAnswer`/`RecordSessionFinished` (B4.2 desviado); prometer `earnedAt`; `errgroup` (D-13); editar `docs/` a mano; tocar mobile.

## Success Criteria

1. `GET /api/v1/gamification/badges` devuelve los **8** badges del catálogo, con `earned`/`total` correctos y orden estable.
2. Los umbrales se cumplen en el borde exacto (ver tabla de tests).
3. Usuario nuevo: `earned: 0`, todos con `current: 0`.
4. La página de logros muestra ganados y bloqueados con progreso, sin depender solo del color, y `make fe-check` + `make fe-lint` + `make fe-test` pasan.
5. `docs/swagger.yaml` incluye el path nuevo.
6. No se añade tabla ni dependencia.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | Definición de "tema dominado" para `topics_mastered_5` | **Basada en el SRS existente**: un tema está dominado si el usuario tiene **≥ 3** preguntas respondidas en él y **todas** tienen `UserQuestionProgress.is_mastered = true`. Reutiliza la semántica ya existente de "mastered" en vez de inventar un umbral de accuracy nuevo (que colisionaría con el nombre de `review-mastery`). Alternativa descartada: `accuracy ≥ 0.8 && attempts ≥ 5`, que introduce un tercer umbral numérico. |
| OQ-2 | Umbrales de `fast_solver` | **≥ 5 respuestas, 100 % acierto, media ≤ 10 000 ms** (nombre del badge en la UI: "Resolvedor veloz"). Ajustable en `catalog.go`/`badges.go` sin cambio de contrato. |
| OQ-3 | ¿`current` en badges booleanos? | `0` o `Target` (sin progreso parcial). Alternativa: ocultar la barra cuando `Target == 1`. |
| OQ-4 | ¿Se persiste algún día? | **No** mientras D-16 esté vigente. Si el producto pide "conseguido el <fecha>", hay que superseder D-16 con un ADR y añadir migración — decisión explícita del usuario, no del implementador. |
