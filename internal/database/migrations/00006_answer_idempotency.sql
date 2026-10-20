-- +goose Up
--
-- Idempotencia de respuestas (D-35): el cliente puede reintentar el envío sin
-- duplicar la fila, el XP ni el progreso SRS.
--
-- `client_answer_id` lo genera el cliente (UUID) una vez por pregunta y se
-- reutiliza en los reintentos. Las filas históricas y los clientes que no lo
-- envían quedan en NULL; Postgres considera los NULL distintos en un índice
-- único, así que no chocan entre sí.

ALTER TABLE session_answers ADD COLUMN client_answer_id uuid;

CREATE UNIQUE INDEX session_answers_client_uniq
    ON session_answers (session_id, client_answer_id);

COMMENT ON COLUMN session_answers.client_answer_id
    IS 'UUID generado por el cliente para idempotencia de reintentos. NULL en filas antiguas o clientes que no lo envían.';

-- +goose Down
DROP INDEX IF EXISTS session_answers_client_uniq;
ALTER TABLE session_answers DROP COLUMN IF EXISTS client_answer_id;
