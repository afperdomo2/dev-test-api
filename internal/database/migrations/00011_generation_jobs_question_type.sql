-- +goose Up
--
-- Tipo de pregunta pedido a la IA en un trabajo de generación.
--
-- Hasta ahora todos los trabajos generaban preguntas single_choice (D-14).
-- Se agrega el tipo para que el admin pueda pedir selección única, selección
-- múltiple, verdadero/falso o una mezcla de los tres (D-54).
--
-- `code_completion` queda fuera a propósito: no se corrige automáticamente en
-- ningún cliente (web siempre la marca incorrecta y mobile no la puede
-- responder), así que generarla en masa solo ensuciaría el banco.
--
-- El default 'single_choice' mantiene veraces los trabajos históricos: todos
-- los anteriores pidieron ese tipo.

ALTER TABLE generation_jobs
    ADD COLUMN IF NOT EXISTS question_type character varying(20) NOT NULL DEFAULT 'single_choice';

ALTER TABLE generation_jobs
    DROP CONSTRAINT IF EXISTS generation_jobs_question_type_valid;

ALTER TABLE generation_jobs
    ADD CONSTRAINT generation_jobs_question_type_valid
    CHECK (question_type IN ('mixed', 'single_choice', 'multiple_choice', 'true_false'));

COMMENT ON COLUMN generation_jobs.question_type IS 'Tipo de pregunta solicitado: mixed (mezcla de 3), single_choice, multiple_choice o true_false. Ignorado en trabajos de tipo tip';

-- +goose Down
ALTER TABLE generation_jobs
    DROP CONSTRAINT IF EXISTS generation_jobs_question_type_valid;

ALTER TABLE generation_jobs
    DROP COLUMN IF EXISTS question_type;
