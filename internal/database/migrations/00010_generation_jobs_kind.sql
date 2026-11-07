-- +goose Up
--
-- Discriminador de contenido para los trabajos de generación con IA.
--
-- Un mismo pipeline (cola, worker pool, presupuesto diario) procesa dos tipos
-- de contenido:
--   question → preguntas del banco (tabla questions)
--   tip      → tips/consejos (tabla tips)
--
-- Se reutiliza el pipeline en vez de duplicar worker y presupuesto (D-51).
-- El default 'question' deja los jobs existentes como estaban.

ALTER TABLE generation_jobs
    ADD COLUMN IF NOT EXISTS kind character varying(20) NOT NULL DEFAULT 'question';

ALTER TABLE generation_jobs
    DROP CONSTRAINT IF EXISTS generation_jobs_kind_valid;

ALTER TABLE generation_jobs
    ADD CONSTRAINT generation_jobs_kind_valid CHECK (kind IN ('question', 'tip'));

COMMENT ON COLUMN generation_jobs.kind IS 'Tipo de contenido que genera el trabajo: question (banco de preguntas) o tip (consejos)';

-- Soporta el listado y los contadores por tipo manteniendo el orden FIFO.
CREATE INDEX IF NOT EXISTS idx_generation_jobs_kind_status
    ON public.generation_jobs USING btree (kind, status, created_at);

COMMENT ON INDEX public.idx_generation_jobs_kind_status IS 'Filtro por tipo + estado y soporte del claim FIFO';

-- +goose Down
DROP INDEX IF EXISTS public.idx_generation_jobs_kind_status;

ALTER TABLE generation_jobs
    DROP CONSTRAINT IF EXISTS generation_jobs_kind_valid;

ALTER TABLE generation_jobs
    DROP COLUMN IF EXISTS kind;
