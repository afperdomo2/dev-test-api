-- +goose Up
--
-- Notificaciones in-app por usuario. La primera (y hoy única) fuente es el job
-- diario de repaso (D-28): la tabla es la fuente del aviso; push solo entrega.
--
-- `creation_day` sostiene la idempotencia del job: una fila por
-- (usuario, tipo, día UTC), aunque el job se ejecute dos veces.

CREATE TABLE notifications (
    id           uuid        PRIMARY KEY,
    user_id      uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type         varchar(30) NOT NULL,
    title        varchar(200) NOT NULL,
    body         text        NOT NULL,
    is_read      boolean     NOT NULL DEFAULT false,
    created_at   timestamptz NOT NULL DEFAULT now(),
    creation_day date        NOT NULL GENERATED ALWAYS AS ((created_at AT TIME ZONE 'UTC')::date) STORED
);

CREATE UNIQUE INDEX notifications_user_type_day_uniq
    ON notifications (user_id, type, creation_day);

CREATE INDEX notifications_user_created_idx ON notifications (user_id, created_at DESC);
CREATE INDEX notifications_user_unread_idx  ON notifications (user_id) WHERE is_read = false;

COMMENT ON TABLE  notifications                IS 'Notificaciones in-app por usuario. Generadas por el job diario y por eventos futuros.';
COMMENT ON COLUMN notifications.id             IS 'Identificador de la notificación (UUID generado en Go).';
COMMENT ON COLUMN notifications.user_id        IS 'Usuario destinatario. Se borra en cascada con el usuario.';
COMMENT ON COLUMN notifications.type           IS 'Tipo de notificación. Hoy solo ''daily_review'' (D-21: string + const Go).';
COMMENT ON COLUMN notifications.title          IS 'Título corto mostrado en la lista.';
COMMENT ON COLUMN notifications.body           IS 'Cuerpo de la notificación.';
COMMENT ON COLUMN notifications.is_read        IS 'Marcada como leída por el usuario.';
COMMENT ON COLUMN notifications.created_at     IS 'Momento de creación (UTC).';
COMMENT ON COLUMN notifications.creation_day   IS 'Día UTC de creación. Columna generada; sostiene la unique de idempotencia. NO se escribe desde la aplicación.';

-- +goose Down
DROP TABLE notifications;
