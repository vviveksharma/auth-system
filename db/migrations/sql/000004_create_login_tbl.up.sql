CREATE TABLE IF NOT EXISTS login_tbl (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID        NOT NULL,
    user_id    UUID        NOT NULL,
    role_id    UUID        NOT NULL,
    role_name  TEXT        NOT NULL,
    jwt_token  TEXT        NOT NULL,
    issued_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked    BOOLEAN     NOT NULL DEFAULT false,
    ip_address VARCHAR(45)
);
