-- +goose Up
--
-- Retirada del modo 'generate' (IA legacy por sesión). Ver ROADMAP Fase 10.
--
-- La columna `sessions.mode` es varchar (NO usa el enum public.session_mode,
-- que es decorativo según la convención del repo). Por tanto no hace falta
-- tocar el tipo: solo se cambia el default, se migra el dato histórico y se
-- añade un CHECK que impida volver a usar 'generate'.
--
-- `questions_generated` contaba las preguntas que la IA legacy generaba por
-- sesión; sin el modo generate deja de tener consumidor (Go ya no la mapea)
-- y se elimina.

UPDATE public.sessions SET mode = 'review' WHERE mode = 'generate';

ALTER TABLE public.sessions ALTER COLUMN mode SET DEFAULT 'review'::character varying;

-- CHECK sin bloquear la tabla en escritura durante la validación.
ALTER TABLE public.sessions ADD CONSTRAINT sessions_mode_valid CHECK (mode IN ('review', 'daily')) NOT VALID;
ALTER TABLE public.sessions VALIDATE CONSTRAINT sessions_mode_valid;

ALTER TABLE public.sessions DROP COLUMN questions_generated;

COMMENT ON TYPE public.session_mode IS 'Modo de sesión: review (repaso manual) o daily (sesión diaria inteligente). El valor generate (IA por sesión) se retiró en la Fase 10 del ROADMAP y queda solo por compatibilidad histórica del tipo.';
COMMENT ON COLUMN public.sessions.mode IS 'Modo de sesión: review (repaso manual) o daily (sesión diaria inteligente). El modo generate (IA por sesión) se retiró en la Fase 10.';

-- +goose Down
ALTER TABLE public.sessions DROP CONSTRAINT IF EXISTS sessions_mode_valid;
ALTER TABLE public.sessions ADD COLUMN IF NOT EXISTS questions_generated bigint DEFAULT 0 NOT NULL;
ALTER TABLE public.sessions ALTER COLUMN mode SET DEFAULT 'generate'::character varying;

COMMENT ON TYPE public.session_mode IS 'Modo de sesión: generate (IA), review (repaso manual) o daily (sesión diaria inteligente)';
COMMENT ON COLUMN public.sessions.mode IS 'Modo de sesión (ver enum session_mode)';