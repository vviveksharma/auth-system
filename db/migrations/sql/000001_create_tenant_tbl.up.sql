CREATE TABLE IF NOT EXISTS tenant_tbl (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT,
    email      TEXT,
    salt       TEXT,
    campany    TEXT,
    password   TEXT,
    status     TEXT         NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    INDEX idx_tenant_tbl_status (status)
);
