CREATE TABLE IF NOT EXISTS user_tbl (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    tenant_id  UUID         NOT NULL,
    name       TEXT,
    email      TEXT,
    password   TEXT,
    salt       TEXT,
    status     BOOLEAN,
    roles      TEXT[]
);
