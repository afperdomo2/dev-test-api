-- +goose Up
--
-- Tokens de push de Expo por dispositivo. Un token identifica un dispositivo:
-- no puede estar asociado a dos usuarios a la vez (si alguien re-inicia sesión
-- en el mismo teléfono, el upsert lo reasigna).

CREATE TABLE device_tokens (
    id           uuid         PRIMARY KEY,
    user_id      uuid         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token        varchar(255) NOT NULL,
    platform     varchar(10)  NOT NULL,
    created_at   timestamptz  NOT NULL DEFAULT now(),
    updated_at   timestamptz  NOT NULL DEFAULT now(),
    last_seen_at timestamptz  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX device_tokens_token_uniq ON device_tokens (token);
CREATE INDEX device_tokens_user_idx ON device_tokens (user_id);

COMMENT ON TABLE  device_tokens              IS 'Tokens de push (Expo) por dispositivo. Un token pertenece a un solo usuario.';
COMMENT ON COLUMN device_tokens.id           IS 'Identificador del registro (UUID generado en Go).';
COMMENT ON COLUMN device_tokens.user_id      IS 'Usuario dueño del dispositivo. Se borra en cascada con el usuario.';
COMMENT ON COLUMN device_tokens.token        IS 'ExponentPushToken[...] devuelto por expo-notifications. Único global.';
COMMENT ON COLUMN device_tokens.platform     IS 'Plataforma del dispositivo: ios o android.';
COMMENT ON COLUMN device_tokens.created_at   IS 'Primera vez que se registró este token.';
COMMENT ON COLUMN device_tokens.updated_at   IS 'Última modificación del registro (upsert).';
COMMENT ON COLUMN device_tokens.last_seen_at IS 'Última vez que la app registró/refrescó el token. Sirve para detectar tokens abandonados.';

-- +goose Down
DROP TABLE device_tokens;
