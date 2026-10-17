# Decisiones — log vivo

> Rescatado de `ROADMAP.md` (sección "Decisiones bloqueadas" + "Notas") al migrar el backlog a SDD. **Estas decisiones no son tareas**: son guardarraíles que restringen cualquier implementación futura.
>
> Reglas: una decisión no se borra nunca; si cambia, se marca `SUPERSEDED por D-xx`. Si una decisión contradice a `.agents/*.md`, gana este archivo y hay que corregir el doc (ver D-13).

## Cómo leer las etiquetas

- **Bloqueada** — no se discute; si algo la viola, el cambio está mal.
- **SUPERSEDED** — reemplazada por otra decisión posterior.
- **Abierta** — pendiente de decisión; su spec tiene la Open Question.

---

## Decisiones de producto / alcance

| ID | Estado | Decisión |
|---|---|---|
| D-01 | Bloqueada | **V/F**: `Question.Type = "true_false"` (2 opciones, 1 correcta). Reusa el pipeline create/import/answer. |
| D-02 | Bloqueada | **Tips ↔ temas**: un solo tema por tip (`Tip.TopicID`, FK indexada). |
| D-03 | Bloqueada | **Creación de tips**: solo admin y siempre públicos. Backend solo manual + import CSV; **nunca** genera tips con IA. El prompt helper de `/questions/import` es texto para copiar/pegar en una IA externa. Usuario normal: solo lectura. |
| D-04 | Bloqueada | **Visibilidad sin cambios**: se mantiene `is_public` + `source`. Solo IA y admin son públicos; el usuario normal queda privado. **Sin** `pending_review`, **sin** acción "publicar", **sin** refactor `isPublic → visibility`. |
| D-05 | Bloqueada | **Moderación**: solo reportar contenido público + revisión admin. Sin "publish". |
| D-06 | Bloqueada | **Import UI (web)**: página unificada con selector "preguntas / V-F / tips"; tips solo visible para admin. |
| D-07 | Bloqueada | **Gamificación**: un solo módulo `internal/modules/gamification/`. **Sin tabla `UserStats`**: se usa `UserStreak` + `UserXp` + `UserTopicXp`. XP por dificultad `beginner=5 / intermediate=10 / advanced=15` + `sessionBonusXp=20` al terminar sesión (`gamification/service.go:13-24`). |
| D-08 | Bloqueada | **Definición única de "tema débil"**: vive en `review-mastery`. `isWeak = accuracy < 0.6 && attempts >= 5` (constantes en Go). Ningún otro módulo la recalcula (`insights` filtra, `badges` consume). |
| D-09 | Bloqueada | **Definición única de "due review"**: `sessions.Store.CountDueReviews` (`sessions/store.go:210`) — `next_review_at <= now() AND is_mastered = false`. No se reimplementa en `progress` ni en `notifications`. |
| D-10 | Bloqueada | **Accuracy**: respuestas (`SessionAnswer`) de **sesiones `completed`**, contando todos los reintentos. `attempts` = respuestas totales al tema (no preguntas distintas). ⚠️ El estado correcto es `completed`, **no** `finished` (el `ROADMAP.md` decía "finished"): `SessionStatus = 'in_progress' \| 'completed' \| 'cancelled'` (`packages/shared/src/types/session.types.ts:5`). |
| D-11 | Bloqueada | **Sesión diaria**: modo `daily` en `Session` + selección on-the-fly, **sin** tabla de cola persistente. Es **mobile-first**: la web **no** tiene botón "Repasar hoy". |
| D-12 | Bloqueada | **Insights** (no "reports"): `GET /insights/summary`, para no chocar con `/reports` (moderación). |
| D-13 | Bloqueada | **Concurrencia**: goroutines + `sync.WaitGroup`. **Sin `errgroup`** (no añadir dependencia). El patrón para trabajo de fondo es `internal/modules/generation/worker.go` (pool + `context` cancelable + `WaitGroup` + reconciliación). ⚠️ **`.agents/backend/architecture.md:144-207` contradice esto**: documenta `errgroup` como canónico y afirma que ya se usa en `sessions/service.go`, lo cual es **falso** (usa `sync.WaitGroup` + `sync.Mutex`, `sessions/service.go:128-129`). Este archivo gana. |
| D-14 | Bloqueada | **IA**: solo por lotes y solo admins (`/generation`), y solo genera preguntas `single_choice`. Tope **global** diario `ai.max_per_day`. El modo `generate` de sesiones y la cuota de IA por usuario **se retiraron** (fases 10 y 11). **No reintroducir IA por usuario en sesiones.** |
| D-15 | Bloqueada | **Cuotas**: `ImportQuota` / `DailyImportLimit` (import CSV, con exención de admin) **se conserva**. La cuota de IA por usuario **no existe**. |
| D-16 | Bloqueada | **Logros**: catálogo en código (`const`), **sin tabla** `UserBadge` ni `user_badges`. Se derivan on-the-fly. Sin `earnedAt` (no hay persistencia). |
| D-17 | Bloqueada | **Scheduler de notificaciones**: ticker stdlib (`time.Timer` calculando el próximo corte) + goroutine. **Sin `robfig/cron`**. |
| D-18 | Bloqueada | **Alcance de clientes del backlog**: backend + web. Mobile fuera de alcance. |
| D-19 | Bloqueada | **Responses sin paginar** en `topics-mastery` (el set de temas es acotado). El resto de listas nuevos sí pagina. |

## Convenciones del repo rescatadas (no estaban en `.agents/`)

| ID | Convención |
|---|---|
| D-20 | `Question.Type` es **`string`**: los `CREATE TYPE ... AS ENUM` de las migraciones son **decorativos**. La validación real es Go `binding:"oneof=..."`. No "arreglar" la DB para usar los enums. |
| D-21 | Enums de dominio se modelan como `string` + constante Go + `CHECK`/`binding` explícito. |
| D-22 | Logs con emoji (`❌` error, `✅` éxito, `🚀` arranque, `🛢️` DB, `🌱` seed, `🤖` IA). |
| D-23 | Todo objeto de esquema (tabla, columna, enum, función, trigger) lleva su `COMMENT ON` **en la migración** que lo crea. Nunca se edita una migración ya aplicada. |
| D-24 | Cambio de esquema = migración goose nueva (`make migrate-create name=<desc>`). `Up` idempotente cuando sea seguro; `Down` reversible o documentado como destructivo. |
| D-25 | Dependencia nueva (Go o pnpm) → pasar por `.agents/backend/dependencies.md` o `.agents/web/dependencies.md` **antes** de instalar. |

## Decisiones descartadas (no reintroducir)

| ID | Descartado | Motivo |
|---|---|---|
| X-01 | Tabla `UserStats` | Sustituida por `UserStreak` + `UserXp` + `UserTopicXp` (D-07). |
| X-02 | Tabla `UserBadge` / `user_badges` | Contradice D-16. El `ROADMAP.md` se contradecía a sí mismo en B4.3; gana D-16. |
| X-03 | `pending_review` + acción "publicar" | D-04. |
| X-04 | Refactor `isPublic → visibility` | D-04. |
| X-05 | `golang.org/x/sync/errgroup` | D-13. |
| X-06 | `github.com/robfig/cron/v3` | D-17. |
| X-07 | Lotes de IA de 20 preguntas / `AI_MAX_PER_JOB=50` | Un job = una llamada, máx 10; tope global (D-14). |
| X-08 | Tabla de cola persistente para la sesión diaria | D-11. |
| X-09 | IA por usuario en sesiones (`mode=generate`) | D-14. |
| X-10 | Cuota de IA por usuario (`users.daily_ai_limit`) | D-15. |

## Notas de implementación relevantes para el backlog

| ID | Nota |
|---|---|
| N-01 | **`xpGained` por periodo no es derivable de las tablas de gamificación**: `UserXp` y `UserTopicXp` son acumuladores de una fila (`models/user_xp.go:12`). `UserStreak.DailyXpToday` solo guarda el día actual (`gamification/service.go:276-281`). Se deriva de `SessionAnswer` (dificultad → 5/10/15) + `20 ×` sesiones `finished` en la ventana. Ver Open Question en `SPEC-insights.md`. |
| N-02 | **`SessionAnswer.ResponseTimeMs`** (`models/session_answer.go:22`) es la fuente de tiempo real de estudio y del badge `fast_solver`. Preferirla a `FinishedAt - StartedAt` (incluye tiempo muerto). |
| N-03 | "Des-publicar" por reporte solo aplica a contenido con `is_public=true` de admin. Des-publicar contenido `ai_generated` queda **fuera de alcance** (requiere mecanismo futuro). |
| N-04 | `.agents/shared/architecture.md:60` sigue mencionando `DEFAULT_DAILY_AI_LIMIT`, eliminado en la fase 11. Referencia obsoleta, no reintroducir el símbolo. |
| N-05 | El estado de sesión es un **literal mágico** en Go: `"completed"` aparece en `sessions/service.go:220,225,257,313` sin constante. `"cancelled"` figura en la doc de Swagger del filtro pero ningún código lo asigna. Aplica D-21: al tocar estos estados desde `review-mastery` o `insights`, definir constantes en `internal/models/session.go` (`SessionStatusInProgress`, `SessionStatusCompleted`, `SessionStatusCancelled`) en vez de ampliar los literales. |
| N-06 | Las migraciones **ya aplicadas** `00004_sessions_mode_no_generate.sql` y `00005_drop_users_daily_ai_limit.sql` citan en sus comentarios ("Ver ROADMAP Fase 10/11") un archivo que **ya no existe**. No se editan (D-23/D-24: nunca tocar una migración aplicada). El contexto vivo de esas fases está en D-14, X-09 y X-10, y el histórico en el commit que eliminó `ROADMAP.md`. |
