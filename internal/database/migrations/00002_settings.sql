-- +goose Up
--
-- Tabla de configuraciones de la aplicación (key/value).
--
-- Contiene únicamente valores TUNABLE (límites, modelo, temperatura, timeout,
-- concurrencia de workers). Los SECRETOS (AI_API_KEY, JWT_SECRET) NO se
-- almacenan aquí: viven exclusivamente en variables de entorno (.env).
--
-- Los valores ausentes caen al default declarado en el registro de Go
-- (internal/modules/settings/registry.go). Por eso `ai.model` y
-- `ai.request_timeout` NO se siembran aquí: su default es el valor de .env.

CREATE TABLE IF NOT EXISTS settings (
    key text NOT NULL,
    value text NOT NULL,
    description text,
    updated_by uuid,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT settings_pkey PRIMARY KEY (key),
    CONSTRAINT fk_settings_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id)
);

COMMENT ON TABLE settings IS 'Configuraciones tunables de la aplicación, editables en runtime por el administrador. Los secretos NO se almacenan aquí (viven en .env).';
COMMENT ON COLUMN settings.key IS 'Clave del ajuste (PK). Prefijo por dominio, p.ej. ai.max_per_job';
COMMENT ON COLUMN settings.value IS 'Valor del ajuste siempre almacenado como texto; el tipo lo define el registro de Go';
COMMENT ON COLUMN settings.description IS 'Descripción legible del ajuste';
COMMENT ON COLUMN settings.updated_by IS 'Administrador que editó el ajuste por última vez (NULL si viene del seed)';
COMMENT ON COLUMN settings.created_at IS 'Fecha de creación del ajuste';
COMMENT ON COLUMN settings.updated_at IS 'Fecha de última actualización del ajuste';

-- Semilla de los ajustes con default puro (no dependen del entorno).
-- Idempotente: no sobrescribe valores ya editados por un administrador.
INSERT INTO settings (key, value, description) VALUES
    ('ai.max_per_job', '10', 'Máximo de preguntas por solicitud de generación con IA'),
    ('ai.max_per_day', '200', 'Máximo de preguntas reservadas por día para generación con IA (tope global)'),
    ('ai.default_quantity', '5', 'Cantidad de preguntas preseleccionada por defecto al encolar un trabajo de generación'),
    ('ai.temperature', '0.8', 'Temperatura del modelo de IA (0.0 - 2.0)'),
    ('ai.worker_concurrency', '2', 'Número de workers que procesan trabajos de generación en paralelo')
ON CONFLICT (key) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS settings;
