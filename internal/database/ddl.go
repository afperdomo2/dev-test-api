package database

import "gorm.io/gorm"

func runDDL(db *gorm.DB) error {
	statements := []string{
		`DO $$ BEGIN
			CREATE TYPE question_type AS ENUM ('single_choice', 'multiple_choice', 'code_completion');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

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
			CREATE TYPE session_mode AS ENUM ('generate', 'review');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`DO $$ BEGIN
			CREATE TYPE session_difficulty AS ENUM ('beginner', 'intermediate', 'advanced');
		EXCEPTION WHEN duplicate_object THEN null; END $$;`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_topics_slug_created_by
			ON topics (slug, created_by) NULLS NOT DISTINCT;`,

		`CREATE INDEX IF NOT EXISTS idx_user_question_progress_next_review
			ON user_question_progress (user_id, next_review_at)
			WHERE is_saved = true;`,

		`CREATE INDEX IF NOT EXISTS idx_user_question_progress_mastered
			ON user_question_progress (user_id, is_mastered)
			WHERE is_mastered = false;`,

		`CREATE INDEX IF NOT EXISTS idx_sessions_user_status
			ON sessions (user_id, status);`,

		`CREATE INDEX IF NOT EXISTS idx_sessions_user_started
			ON sessions (user_id, started_at DESC);`,

		`CREATE INDEX IF NOT EXISTS idx_session_answers_session
			ON session_answers (session_id);`,

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
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}

	return nil
}
