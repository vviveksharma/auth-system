CREATE TABLE IF NOT EXISTS reset_creds_tbl (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID        NOT NULL,
    user_id    UUID        NOT NULL,
    active     BOOLEAN,
    code_hash  TEXT        NOT NULL,
    salt       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    used_at    TIMESTAMPTZ,
    INDEX idx_reset_creds_code_hash (code_hash),
    INDEX idx_reset_creds_salt      (salt)
);
