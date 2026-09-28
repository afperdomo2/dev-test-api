-- +goose Up
--
-- Retirada de la cuota de IA por usuario. Ver ROADMAP Fase 11.
--
-- Tras la Fase 10 (retirada del modo 'generate' de las sesiones), la IA se
-- genera exclusivamente por el batch admin con presupuesto global
-- (ai.max_per_day). La columna users.daily_ai_limit queda sin consumidor y
-- se elimina.

ALTER TABLE public.users DROP COLUMN IF EXISTS daily_ai_limit;

-- +goose Down
ALTER TABLE public.users ADD COLUMN IF NOT EXISTS daily_ai_limit bigint DEFAULT 20 NOT NULL;
COMMENT ON COLUMN public.users.daily_ai_limit IS 'Cupo diario de preguntas generadas por IA (obsoleto: la IA ahora se genera solo por el batch admin con presupuesto global)';