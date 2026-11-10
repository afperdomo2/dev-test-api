-- +goose Up
--
-- Nuevos defaults de generación con IA: 20 preguntas por trabajo y techo
-- configurable de 50 (antes 10). El registry de Go ya declara los nuevos
-- defaults, pero las filas sembradas por 00002_settings.sql los pisan porque
-- `resolve` prioriza el valor persistido.
--
-- Solo se reescribe el valor si sigue siendo el seed viejo: una
-- personalización del administrador se respeta.

UPDATE settings
SET value = '20', updated_at = now()
WHERE key = 'ai.max_per_job' AND value = '10';

UPDATE settings
SET value = '20', updated_at = now()
WHERE key = 'ai.default_quantity' AND value = '5';

-- +goose Down
-- Revierte únicamente si el valor sigue en el nuevo default.
UPDATE settings
SET value = '10', updated_at = now()
WHERE key = 'ai.max_per_job' AND value = '20';

UPDATE settings
SET value = '5', updated_at = now()
WHERE key = 'ai.default_quantity' AND value = '20';
