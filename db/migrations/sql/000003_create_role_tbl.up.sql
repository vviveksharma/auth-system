CREATE TABLE IF NOT EXISTS role_tbl (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    role         TEXT,
    display_name TEXT,
    description  TEXT,
    role_id      UUID,
    tenant_id    UUID        NOT NULL,
    role_type    TEXT,
    status       BOOLEAN,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
