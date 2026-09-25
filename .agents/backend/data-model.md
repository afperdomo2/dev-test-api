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

## Enums

- `question_type`: single_choice, multiple_choice, code_completion, true_false
- `question_difficulty` / `session_difficulty`: beginner, intermediate, advanced
- `question_source`: ai_generated, manual, imported
- `session_status`: in_progress, completed, cancelled
- `session_mode`: generate, review, daily
- `report_status`: pending, reviewed
- `tips.source` (columna sin enum): manual, imported
- `reports.target_type` (columna sin enum): question, tip

## Índices / trigger (migraciones, `00001_init.sql`)

- Único `topics (slug, created_by)` (NULLS NOT DISTINCT)
- `user_question_progress (user_id, next_review_at)` y `(user_id, is_mastered)`
- `sessions (user_id, status)` y `(user_id, started_at DESC)`
- `session_answers (session_id)`
- Trigger `calc_session_score()` → calcula `sessions.score` al pasar a `completed`