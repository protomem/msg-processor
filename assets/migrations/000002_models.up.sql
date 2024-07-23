BEGIN;

CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    message TEXT NOT NULL                   CHECK (length(message) > 0),
    status  TEXT NOT NULL DEFAULT 'created' CHECK (status = 'created' OR status = 'processing' OR status = 'completed')
);

COMMIT;