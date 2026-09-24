package database

import "gorm.io/gorm"

func runDDL(db *gorm.DB) error {
	statements := []string{
		// ===== Tipos (enums) =====
		`DO $$ BEGIN
			CREATE TYPE question_type AS ENUM ('single_choice', 'multiple_choice', 'code_completion', 'true_false');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`ALTER TYPE question_type ADD VALUE IF NOT EXISTS 'true_false';`,

		`DO $$ BEGIN
			CREATE TYPE question_difficulty AS ENUM ('beginner', 'intermediate', 'advanced');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`DO $$ BEGIN
			CREATE TYPE question_source AS ENUM ('ai_generated', 'manual', 'imported');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`DO $$ BEGIN
			CREATE TYPE session_status AS ENUM ('in_progress', 'completed', 'cancelled');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`DO $$ BEGIN
			CREATE TYPE session_mode AS ENUM ('generate', 'review', 'daily');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`ALTER TYPE session_mode ADD VALUE IF NOT EXISTS 'daily';`,

		`DO $$ BEGIN
			CREATE TYPE session_difficulty AS ENUM ('beginner', 'intermediate', 'advanced');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`DO $$ BEGIN
			CREATE TYPE report_status AS ENUM ('pending', 'reviewed');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		// ===== Índices =====
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_topics_slug_created_by
			ON topics (slug, created_by) NULLS NOT DISTINCT;`,

		`DROP INDEX IF EXISTS idx_user_question_progress_next_review;
		CREATE INDEX idx_user_question_progress_next_review
			ON user_question_progress (user_id, next_review_at)
			WHERE next_review_at IS NOT NULL;`,

		`CREATE INDEX IF NOT EXISTS idx_user_question_progress_mastered
			ON user_question_progress (user_id, is_mastered)
			WHERE is_mastered = false;`,

		`CREATE INDEX IF NOT EXISTS idx_sessions_user_status
			ON sessions (user_id, status);`,

		`CREATE INDEX IF NOT EXISTS idx_sessions_user_started
			ON sessions (user_id, started_at DESC);`,

		`CREATE INDEX IF NOT EXISTS idx_session_answers_session
			ON session_answers (session_id);`,

		`CREATE INDEX IF NOT EXISTS idx_tips_topic_id
			ON tips (topic_id);`,

		`CREATE INDEX IF NOT EXISTS idx_reports_target
			ON reports (target_type, target_id);`,

		`CREATE INDEX IF NOT EXISTS idx_reports_status
			ON reports (status);`,

		`CREATE INDEX IF NOT EXISTS idx_questions_type
			ON questions (type);`,

		// ===== Función y trigger de score de sesión =====
		`CREATE OR REPLACE FUNCTION calc_session_score()
		RETURNS TRIGGER AS $$
		BEGIN
			IF NEW.status = 'completed' AND (OLD.status IS DISTINCT FROM 'completed') THEN
				NEW.score = (
					SELECT ROUND(
						(COUNT(*) FILTER (WHERE is_correct = true)::decimal /
						 GREATEST(COUNT(*), 1)) * 100, 1
					)
					FROM session_answers
					WHERE session_id = NEW.id
				);
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;`,

		`DROP TRIGGER IF EXISTS trg_sessions_score ON sessions;
		CREATE TRIGGER trg_sessions_score
			BEFORE UPDATE ON sessions
			FOR EACH ROW
			EXECUTE FUNCTION calc_session_score();`,

		// ===== Documentación: objetos nuevos (tablas/columnas/enums) =====
		`COMMENT ON TABLE tips IS 'Tarjetas tip creadas por administradores (siempre públicas). Leídas por cualquier usuario autenticado y usadas por la sesión diaria.';`,
		`COMMENT ON COLUMN tips.id IS 'Identificador UUID del tip';`,
		`COMMENT ON COLUMN tips.user_id IS 'Administrador que creó el tip';`,
		`COMMENT ON COLUMN tips.topic_id IS 'Tema al que pertenece el tip (un único tema por tip)';`,
		`COMMENT ON COLUMN tips.content IS 'Contenido del tip en texto plano';`,
		`COMMENT ON COLUMN tips.source IS 'Origen del tip: manual o imported (sin contenido generado por IA)';`,
		`COMMENT ON COLUMN tips.created_at IS 'Fecha de creación del tip';`,
		`COMMENT ON COLUMN tips.updated_at IS 'Fecha de última actualización del tip';`,
		`COMMENT ON COLUMN tips.deleted_at IS 'Soft delete (borrado lógico) del tip';`,

		`COMMENT ON TABLE reports IS 'Reportes de contenido (preguntas o tips) hechos por usuarios para revisión del administrador.';`,
		`COMMENT ON COLUMN reports.id IS 'Identificador UUID del reporte';`,
		`COMMENT ON COLUMN reports.reporter_user_id IS 'Usuario que realizó el reporte';`,
		`COMMENT ON COLUMN reports.target_type IS 'Tipo de contenido reportado: question o tip';`,
		`COMMENT ON COLUMN reports.target_id IS 'Identificador del contenido reportado';`,
		`COMMENT ON COLUMN reports.reason IS 'Motivo del reporte en texto libre';`,
		`COMMENT ON COLUMN reports.status IS 'Estado del reporte: pending o reviewed';`,
		`COMMENT ON COLUMN reports.created_at IS 'Fecha de creación del reporte';`,
		`COMMENT ON COLUMN reports.updated_at IS 'Fecha de última actualización del reporte';`,

		`COMMENT ON TABLE user_streaks IS 'Racha diaria de actividad del usuario (días consecutivos con pregunta respondida), racha más larga y freezes disponibles.';`,
		`COMMENT ON COLUMN user_streaks.user_id IS 'Usuario (PK)';`,
		`COMMENT ON COLUMN user_streaks.current_streak IS 'Racha actual en días consecutivos';`,
		`COMMENT ON COLUMN user_streaks.longest_streak IS 'Racha más larga alcanzada en días';`,
		`COMMENT ON COLUMN user_streaks.last_activity_at IS 'Timestamp de la última actividad';`,
		`COMMENT ON COLUMN user_streaks.last_activity_date IS 'Fecha (YYYY-MM-DD, UTC) de la última actividad para comparar días consecutivos';`,
		`COMMENT ON COLUMN user_streaks.freezes_available IS 'Freezes disponibles (congelan un día sin actividad, tope 3)';`,
		`COMMENT ON COLUMN user_streaks.daily_xp_today IS 'XP acumulado en el día actual para la meta diaria';`,
		`COMMENT ON COLUMN user_streaks.daily_xp_date IS 'Fecha (YYYY-MM-DD, UTC) a la que corresponde daily_xp_today';`,
		`COMMENT ON COLUMN user_streaks.created_at IS 'Fecha de creación del registro';`,
		`COMMENT ON COLUMN user_streaks.updated_at IS 'Fecha de última actualización del registro';`,

		`COMMENT ON TABLE user_xps IS 'XP global del usuario (suma acumulada de todas las respuestas correctas).';`,
		`COMMENT ON COLUMN user_xps.user_id IS 'Usuario (PK)';`,
		`COMMENT ON COLUMN user_xps.xp IS 'XP acumulado total';`,
		`COMMENT ON COLUMN user_xps.created_at IS 'Fecha de creación del registro';`,
		`COMMENT ON COLUMN user_xps.updated_at IS 'Fecha de última actualización del registro';`,

		`COMMENT ON TABLE user_topic_xps IS 'XP del usuario por tema (desglose del XP global).';`,
		`COMMENT ON COLUMN user_topic_xps.user_id IS 'Usuario (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN user_topic_xps.topic_id IS 'Tema (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN user_topic_xps.xp IS 'XP acumulado en el tema';`,
		`COMMENT ON COLUMN user_topic_xps.created_at IS 'Fecha de creación del registro';`,
		`COMMENT ON COLUMN user_topic_xps.updated_at IS 'Fecha de última actualización del registro';`,

		`COMMENT ON COLUMN users.daily_goal_xp IS 'Meta diaria de XP configurable por el usuario (presets: casual=25, regular=50, serious=100)';`,
		`COMMENT ON COLUMN sessions.tip_id IS 'Tip del día asociado a una sesión en modo daily (nullable en el resto de modos)';`,

		`COMMENT ON TYPE question_type IS 'Tipo de pregunta: single_choice, multiple_choice, code_completion o true_false';`,
		`COMMENT ON TYPE session_mode IS 'Modo de sesión: generate (IA), review (repaso manual) o daily (sesión diaria inteligente)';`,
		`COMMENT ON TYPE report_status IS 'Estado de un reporte: pending o reviewed';`,

		// ===== Documentación: esquema existente =====
		`COMMENT ON TABLE users IS 'Usuarios de la plataforma (login, permisos y límites diarios)';`,
		`COMMENT ON COLUMN users.id IS 'Identificador UUID del usuario';`,
		`COMMENT ON COLUMN users.email IS 'Email único del usuario (usado para autenticación)';`,
		`COMMENT ON COLUMN users.password_hash IS 'Hash del password (BCrypt); nunca se expone en JSON';`,
		`COMMENT ON COLUMN users.is_admin IS 'Indica si el usuario tiene rol administrador';`,
		`COMMENT ON COLUMN users.daily_import_limit IS 'Cupo diario de preguntas importadas por CSV';`,
		`COMMENT ON COLUMN users.daily_ai_limit IS 'Cupo diario de preguntas generadas por IA';`,
		`COMMENT ON COLUMN users.created_at IS 'Fecha de creación del usuario';`,
		`COMMENT ON COLUMN users.updated_at IS 'Fecha de última actualización del usuario';`,
		`COMMENT ON COLUMN users.deleted_at IS 'Soft delete (borrado lógico) del usuario';`,

		`COMMENT ON TABLE topics IS 'Temas de estudio. Los temas del sistema tienen is_system=true; los creados por usuarios quedan ligados a created_by.';`,
		`COMMENT ON COLUMN topics.id IS 'Identificador UUID del tema';`,
		`COMMENT ON COLUMN topics.slug IS 'Identificador corto único por creador (system o usuario)';`,
		`COMMENT ON COLUMN topics.name IS 'Nombre visible del tema';`,
		`COMMENT ON COLUMN topics.category IS 'Categoría a la que pertenece el tema';`,
		`COMMENT ON COLUMN topics.is_system IS 'Indica si el tema es del sistema (visible para todos)';`,
		`COMMENT ON COLUMN topics.created_by IS 'Usuario creador del tema (NULL si es del sistema)';`,
		`COMMENT ON COLUMN topics.created_at IS 'Fecha de creación del tema';`,
		`COMMENT ON COLUMN topics.updated_at IS 'Fecha de última actualización del tema';`,

		`COMMENT ON TABLE questions IS 'Banco de preguntas. Las de source ai_generated y las públicas (is_public=true) son visibles para todos; el resto es privado de su creador.';`,
		`COMMENT ON COLUMN questions.id IS 'Identificador UUID de la pregunta';`,
		`COMMENT ON COLUMN questions.user_id IS 'Usuario creador de la pregunta';`,
		`COMMENT ON COLUMN questions.type IS 'Tipo de pregunta (ver enum question_type)';`,
		`COMMENT ON COLUMN questions.content IS 'Enunciado de la pregunta';`,
		`COMMENT ON COLUMN questions.explanation IS 'Explicación de la respuesta correcta';`,
		`COMMENT ON COLUMN questions.difficulty IS 'Dificultad: beginner, intermediate o advanced';`,
		`COMMENT ON COLUMN questions.is_public IS 'Indica si la pregunta es pública (visible para todos)';`,
		`COMMENT ON COLUMN questions.source IS 'Origen de la pregunta: ai_generated, manual o imported';`,
		`COMMENT ON COLUMN questions.created_at IS 'Fecha de creación de la pregunta';`,
		`COMMENT ON COLUMN questions.updated_at IS 'Fecha de última actualización de la pregunta';`,
		`COMMENT ON COLUMN questions.deleted_at IS 'Soft delete (borrado lógico) de la pregunta';`,

		`COMMENT ON TABLE question_options IS 'Opciones de respuesta de preguntas de tipo choice (single_choice, multiple_choice, true_false).';`,
		`COMMENT ON COLUMN question_options.id IS 'Identificador UUID de la opción';`,
		`COMMENT ON COLUMN question_options.question_id IS 'Pregunta a la que pertenece la opción';`,
		`COMMENT ON COLUMN question_options.content IS 'Texto de la opción';`,
		`COMMENT ON COLUMN question_options.is_correct IS 'Indica si la opción es la correcta';`,
		`COMMENT ON COLUMN question_options.created_at IS 'Fecha de creación de la opción';`,
		`COMMENT ON COLUMN question_options.updated_at IS 'Fecha de última actualización de la opción';`,

		`COMMENT ON TABLE code_challenges IS 'Desafío de código asociado a preguntas tipo code_completion.';`,
		`COMMENT ON COLUMN code_challenges.id IS 'Identificador UUID del desafío';`,
		`COMMENT ON COLUMN code_challenges.question_id IS 'Pregunta (única) asociada al desafío';`,
		`COMMENT ON COLUMN code_challenges.starter_code IS 'Código inicial provisto al usuario';`,
		`COMMENT ON COLUMN code_challenges.expected_output IS 'Salida esperada para validar la solución';`,
		`COMMENT ON COLUMN code_challenges.language IS 'Lenguaje de programación del desafío';`,
		`COMMENT ON COLUMN code_challenges.test_cases_json IS 'Casos de prueba en formato JSONB';`,
		`COMMENT ON COLUMN code_challenges.created_at IS 'Fecha de creación del desafío';`,
		`COMMENT ON COLUMN code_challenges.updated_at IS 'Fecha de última actualización del desafío';`,

		`COMMENT ON TABLE question_topics IS 'Relación N:M entre preguntas y temas.';`,
		`COMMENT ON COLUMN question_topics.question_id IS 'Pregunta (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN question_topics.topic_id IS 'Tema (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN question_topics.created_at IS 'Fecha de creación de la relación';`,

		`COMMENT ON TABLE user_question_progress IS 'Progreso individual del usuario por pregunta (repeticiones, facilidad, fecha de próximo repaso) basado en SM-2.';`,
		`COMMENT ON COLUMN user_question_progress.user_id IS 'Usuario (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN user_question_progress.question_id IS 'Pregunta (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN user_question_progress.repetitions IS 'Número de repeticiones exitosas consecutivas';`,
		`COMMENT ON COLUMN user_question_progress.ease_factor IS 'Factor de facilidad SM-2 (default 2.5)';`,
		`COMMENT ON COLUMN user_question_progress.interval_days IS 'Intervalo actual entre repasos en días';`,
		`COMMENT ON COLUMN user_question_progress.next_review_at IS 'Fecha del próximo repaso programado';`,
		`COMMENT ON COLUMN user_question_progress.last_reviewed_at IS 'Fecha del último repaso realizado';`,
		`COMMENT ON COLUMN user_question_progress.is_saved IS 'Indica si el usuario guardó la pregunta como favorita';`,
		`COMMENT ON COLUMN user_question_progress.is_mastered IS 'Indica si la pregunta se considera dominada (no vuelve a aparecer en repasos)';`,
		`COMMENT ON COLUMN user_question_progress.created_at IS 'Fecha de creación del progreso';`,
		`COMMENT ON COLUMN user_question_progress.updated_at IS 'Fecha de última actualización del progreso';`,

		`COMMENT ON TABLE sessions IS 'Sesiones de estudio. Modos: generate (IA), review (repaso manual) y daily (sesión diaria inteligente con tip).';`,
		`COMMENT ON COLUMN sessions.id IS 'Identificador UUID de la sesión';`,
		`COMMENT ON COLUMN sessions.user_id IS 'Usuario al que pertenece la sesión';`,
		`COMMENT ON COLUMN sessions.name IS 'Nombre visible de la sesión';`,
		`COMMENT ON COLUMN sessions.status IS 'Estado: in_progress, completed o cancelled';`,
		`COMMENT ON COLUMN sessions.mode IS 'Modo de sesión (ver enum session_mode)';`,
		`COMMENT ON COLUMN sessions.difficulty IS 'Dificultad de las preguntas de la sesión (ver enum session_difficulty)';`,
		`COMMENT ON COLUMN sessions.question_limit IS 'Límite de preguntas de la sesión (NULL = sin límite)';`,
		`COMMENT ON COLUMN sessions.questions_generated IS 'Cantidad de preguntas generadas para la sesión';`,
		`COMMENT ON COLUMN sessions.tip_id IS 'Tip del día asociado a una sesión en modo daily (nullable en el resto de modos)';`,
		`COMMENT ON COLUMN sessions.score IS 'Puntuación final calculada al completar la sesión (0-100, ver trigger trg_sessions_score)';`,
		`COMMENT ON COLUMN sessions.started_at IS 'Fecha de inicio de la sesión';`,
		`COMMENT ON COLUMN sessions.finished_at IS 'Fecha de finalización de la sesión (NULL mientras esté en progreso)';`,
		`COMMENT ON COLUMN sessions.created_at IS 'Fecha de creación de la sesión';`,
		`COMMENT ON COLUMN sessions.updated_at IS 'Fecha de última actualización de la sesión';`,

		`COMMENT ON TABLE session_topics IS 'Relación N:M entre sesiones y temas.';`,
		`COMMENT ON COLUMN session_topics.session_id IS 'Sesión (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN session_topics.topic_id IS 'Tema (parte de la PK compuesta)';`,
		`COMMENT ON COLUMN session_topics.created_at IS 'Fecha de creación de la relación';`,

		`COMMENT ON TABLE session_answers IS 'Respuestas del usuario a preguntas dentro de una sesión (con resultado y tiempo de respuesta).';`,
		`COMMENT ON COLUMN session_answers.id IS 'Identificador UUID de la respuesta';`,
		`COMMENT ON COLUMN session_answers.session_id IS 'Sesión a la que pertenece la respuesta';`,
		`COMMENT ON COLUMN session_answers.user_id IS 'Usuario que respondió';`,
		`COMMENT ON COLUMN session_answers.question_id IS 'Pregunta respondida';`,
		`COMMENT ON COLUMN session_answers.answer_text IS 'Respuesta en texto (para code_completion)';`,
		`COMMENT ON COLUMN session_answers.selected_options IS 'Opciones seleccionadas en formato JSONB (array de UUIDs)';`,
		`COMMENT ON COLUMN session_answers.is_correct IS 'Indica si la respuesta fue correcta';`,
		`COMMENT ON COLUMN session_answers.ai_feedback IS 'Feedback generado por IA sobre la respuesta';`,
		`COMMENT ON COLUMN session_answers.response_time_ms IS 'Tiempo de respuesta en milisegundos';`,
		`COMMENT ON COLUMN session_answers.created_at IS 'Fecha de creación de la respuesta';`,

		`COMMENT ON TYPE question_difficulty IS 'Dificultad de la pregunta: beginner, intermediate o advanced';`,
		`COMMENT ON TYPE question_source IS 'Origen de la pregunta: ai_generated, manual o imported';`,
		`COMMENT ON TYPE session_status IS 'Estado de una sesión: in_progress, completed o cancelled';`,
		`COMMENT ON TYPE session_difficulty IS 'Dificultad de una sesión: beginner, intermediate o advanced';`,

		`COMMENT ON FUNCTION calc_session_score() IS 'Calcula el score (0-100) de la sesión al pasarla a completed, basado en el porcentaje de respuestas correctas';`,
		`COMMENT ON TRIGGER trg_sessions_score ON sessions IS 'Ejecuta calc_session_score() al actualizar una sesión para asignar el score al completarla';`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}

	return nil
}
