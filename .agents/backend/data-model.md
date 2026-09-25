# Backend data model — contexto

Catálogo rápido de la estructura de la DB. **Fuente canónica**: `internal/models/*.go` (structs GORM) y **`internal/database/migrations/*.sql`** (esquema real, baseline `00001_init.sql`). Si este archivo difiere del código, leer las migraciones (los `COMMENT ON` viven ahí).

## Tablas (baseline `00001_init.sql`)

| Tabla | Propósito |
|-------|-----------|
| `users` | Usuarios + límites (isAdmin, dailyImportLimit, dailyAiLimit) |
| `topics` | Temas de estudio (slug, name, category, isSystem, createdBy) |
| `questions` | Preguntas (type, content, difficulty, isPublic, source) |
| `question_options` | Opciones de pregunta (content, isCorrect) |
| `code_challenges` | Retos de código (starterCode, language, testCases) |
| `question_topics` | Join N:M Question ↔ Topic |
| `user_question_progress` | Progreso SRS por usuario+pregunta (repetitions, easeFactor, intervalDays, isMastered) |
| `sessions` | Sesión de práctica (status, mode, difficulty, score, tipId) |
| `session_topics` | Join N:M Session ↔ Topic |
| `session_answers` | Respuestas de la sesión (isCorrect, aiFeedback, responseTimeMs) |
| `tips` | Tips del día (content, source, topicId) |
| `reports` | Reportes de moderación (targetType, targetId, reason, status) |
| `user_xps` | XP total por usuario (xp, level) |
| `user_streaks` | Racha (current, longest, freezesAvailable, lastActiveDay) |
| `user_topic_xps` | XP por usuario+tema (xp) |
| `settings` | Clave/valor de configuración administrable (migración `00002_settings.sql`) |
| `generation_jobs` | Lotes de preguntas generados con IA (estado, presupuesto, dueño) (migración `00003_generation_jobs.sql`) |
| `generation_job_topics` | Join N:M GenerationJob ↔ Topic (migración `00003_generation_jobs.sql`) |

## Enums

- `question_type`: single_choice, multiple_choice, code_completion, true_false
- `question_difficulty` / `session_difficulty`: beginner, intermediate, advanced
- `question_source`: ai_generated, manual, imported
- `session_status`: in_progress, completed, cancelled
- `session_mode`: generate, review, daily
- `report_status`: pending, reviewed
- `tips.source` (columna sin enum): manual, imported
- `reports.target_type` (columna sin enum): question, tip
- `generation_jobs.status` (columna varchar + CHECK, no `CREATE TYPE`): pending, running, completed, failed
- `generation_jobs.difficulty` (columna varchar + CHECK, no `CREATE TYPE`): beginner, intermediate, advanced, random

## `settings` (`00002_settings.sql`)

Configuración **tunable** administrable en runtime vía `GET/PUT /api/v1/settings`. Los **secretos no viven aquí** (`AI_API_KEY`, `JWT_SECRET` siguen en `.env`).

- Columnas: `key` (text, PK), `value` (text — el tipo lo define el registry de Go), `description`, `updated_by` (uuid FK `users`, `ON DELETE` restringido; NULL si viene del seed), `created_at`, `updated_at`.
- Seeding idempotente de 5 claves con `ON CONFLICT (key) DO NOTHING` (no pisa valores ya editados): `ai.max_per_job=10`, `ai.max_per_day=200`, `ai.default_quantity=5`, `ai.temperature=0.8`, `ai.worker_concurrency=2`.
- `ai.model` y `ai.request_timeout` **no** se siembran: su default es el valor de `config.AIConfig` (`.env`).
- El registry (`internal/modules/settings/registry.go`) valida 7 claves; `ai.max_per_job` y `ai.default_quantity` tienen techo duro 10. `ai.temperature` es decimal válido y el service lo expone como `*float64` para distinguir "sin definir" de `0`.

## `generation_jobs` (`00003_generation_jobs.sql`)

Un registro = una llamada a la IA que produce hasta 10 preguntas. `quantity` es lo **reservado** al encolar (cobrable contra el tope); `generated` es lo realmente devuelto y persistido en `questions`.

- Columnas: `id` (uuid PK), `requested_by` (uuid FK `users`, sin acción en borrado), `difficulty` (default `intermediate`), `quantity`, `status` (default `pending`), `generated` (default 0), `error` (text, solo si `failed`), `created_at`, `updated_at`, `started_at`, `finished_at`.
- CHECKs: `quantity > 0`, `generated >= 0`, `status IN ('pending','running','completed','failed')`, `difficulty IN ('beginner','intermediate','advanced','random')`.
- Índices: btree en `requested_by`, `created_at` y `status`, más **idx_generation_jobs_pending** (btree en `created_at` `WHERE status = 'pending'`) que soporta el claim FIFO.
- `generation_job_topics`: PK compuesta `(job_id, topic_id)`, ambas FK con `ON DELETE CASCADE` en el job, índice btree en `topic_id`.
- Presupuesto diario **global** (no por usuario): `SUM(quantity)` de jobs **no fallidos** con `created_at` en el día **UTC** en curso; la reserva y el insert ocurren en la misma transacción bajo `pg_advisory_xact_lock`, y el claim toma un `pending` con `FOR UPDATE SKIP LOCKED`.
- Las preguntas resultantes se guardan en `questions` con `source = 'ai_generated'`, `is_public = true` y `created_by = requested_by`.

## Índices / trigger del baseline (`00001_init.sql`)

> Los índices de `settings` y `generation_jobs` están documentados en sus secciones respectivas.

- Único `topics (slug, created_by)` (NULLS NOT DISTINCT)
- `user_question_progress (user_id, next_review_at)` y `(user_id, is_mastered)`
- `sessions (user_id, status)` y `(user_id, started_at DESC)`
- `session_answers (session_id)`
- Trigger `calc_session_score()` → calcula `sessions.score` al pasar a `completed`