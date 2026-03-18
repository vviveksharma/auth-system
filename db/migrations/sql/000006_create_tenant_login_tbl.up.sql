CREATE TABLE IF NOT EXISTS tenant_login_tbl (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT,
    tenant_id  UUID        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active  BOOLEAN,
    ip_address TEXT
);
