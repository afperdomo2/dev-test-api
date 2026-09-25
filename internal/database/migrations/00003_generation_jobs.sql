-- +goose Up
--
-- Trabajos de generación de preguntas con IA (batch, disparados por admins).
--
-- El procesamiento lo hace un pool de workers (internal/modules/generation/worker.go):
--   pending → running → completed | failed
--
-- `quantity` es la cantidad RESERVADA al encolar el trabajo; `generated` es la
-- cantidad real que la IA devolvió y se guardó en el banco de preguntas.
-- El presupuesto diario se computa como SUM(quantity) de los trabajos NO fallidos
-- creados hoy, de modo que encolar dos veces en paralelo no sobrepasa el tope.
--
-- `status` es varchar (no un CREATE TYPE): sigue la convención del repo donde la
-- validación real se hace en Go (binding:oneof) y los enums de Postgres son decorativos.

CREATE TABLE IF NOT EXISTS generation_jobs (
    id uuid NOT NULL,
    requested_by uuid NOT NULL,
    difficulty character varying(20) NOT NULL DEFAULT 'intermediate',
    quantity integer NOT NULL,
    status character varying(20) NOT NULL DEFAULT 'pending',
    generated integer NOT NULL DEFAULT 0,
    error text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    CONSTRAINT generation_jobs_pkey PRIMARY KEY (id),
    CONSTRAINT generation_jobs_quantity_positive CHECK (quantity > 0),
    CONSTRAINT generation_jobs_generated_non_negative CHECK (generated >= 0),
    CONSTRAINT generation_jobs_status_valid CHECK (status IN ('pending', 'running', 'completed', 'failed')),
    CONSTRAINT generation_jobs_difficulty_valid CHECK (difficulty IN ('beginner', 'intermediate', 'advanced', 'random')),
    CONSTRAINT fk_generation_jobs_requester FOREIGN KEY (requested_by) REFERENCES public.users(id)
);

COMMENT ON TABLE generation_jobs IS 'Trabajos de generación de preguntas con IA en lotes, disparados por administradores y procesados por un pool de workers.';
COMMENT ON COLUMN generation_jobs.id IS 'Identificador UUID del trabajo (generado en Go por BeforeCreate)';
COMMENT ON COLUMN generation_jobs.requested_by IS 'Administrador que encoló el trabajo y que será dueño de las preguntas generadas';
COMMENT ON COLUMN generation_jobs.difficulty IS 'Dificultad solicitada: beginner, intermediate, advanced o random';
COMMENT ON COLUMN generation_jobs.quantity IS 'Cantidad de preguntas reservadas al encolar (objetivo del trabajo, cobrable contra el tope diario)';
COMMENT ON COLUMN generation_jobs.status IS 'Estado del trabajo: pending, running, completed o failed';
COMMENT ON COLUMN generation_jobs.generated IS 'Cantidad de preguntas realmente generadas y guardadas en el banco';
COMMENT ON COLUMN generation_jobs.error IS 'Motivo del fallo cuando status = failed (NULL en el resto de estados)';
COMMENT ON COLUMN generation_jobs.created_at IS 'Fecha de encolado (define a qué día se le cobra contra el tope diario)';
COMMENT ON COLUMN generation_jobs.updated_at IS 'Fecha de última actualización del trabajo';
COMMENT ON COLUMN generation_jobs.started_at IS 'Fecha en que un worker reclamó el trabajo (NULL mientras está pending)';
COMMENT ON COLUMN generation_jobs.finished_at IS 'Fecha de finalización (completed o failed; NULL en pending/running)';

CREATE INDEX IF NOT EXISTS idx_generation_jobs_requested_by ON public.generation_jobs USING btree (requested_by);
CREATE INDEX IF NOT EXISTS idx_generation_jobs_created_at ON public.generation_jobs USING btree (created_at);
CREATE INDEX IF NOT EXISTS idx_generation_jobs_status ON public.generation_jobs USING btree (status);
-- Soporta el claim FIFO del worker (WHERE status = 'pending' ORDER BY created_at)
CREATE INDEX IF NOT EXISTS idx_generation_jobs_pending ON public.generation_jobs USING btree (created_at) WHERE status = 'pending';

CREATE TABLE IF NOT EXISTS generation_job_topics (
    job_id uuid NOT NULL,
    topic_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT generation_job_topics_pkey PRIMARY KEY (job_id, topic_id),
    CONSTRAINT fk_generation_job_topics_job FOREIGN KEY (job_id) REFERENCES public.generation_jobs(id) ON DELETE CASCADE,
    CONSTRAINT fk_generation_job_topics_topic FOREIGN KEY (topic_id) REFERENCES public.topics(id)
);

COMMENT ON TABLE generation_job_topics IS 'Relación N:M entre trabajos de generación con IA y los temas requested.';
COMMENT ON COLUMN generation_job_topics.job_id IS 'Trabajo de generación (parte de la PK compuesta)';
COMMENT ON COLUMN generation_job_topics.topic_id IS 'Tema solicitado para la generación (parte de la PK compuesta)';
COMMENT ON COLUMN generation_job_topics.created_at IS 'Fecha de creación de la relación';

CREATE INDEX IF NOT EXISTS idx_generation_job_topics_topic_id ON public.generation_job_topics USING btree (topic_id);

-- +goose Down
DROP TABLE IF EXISTS generation_job_topics;
DROP TABLE IF EXISTS generation_jobs;
