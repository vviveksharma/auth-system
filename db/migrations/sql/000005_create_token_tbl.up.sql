CREATE TABLE IF NOT EXISTS token_tbl (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL,
    name            TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at    TIMESTAMPTZ,
    usage_count     BIGINT      NOT NULL DEFAULT 0,
    expires_at      TIMESTAMPTZ NOT NULL,
    is_active       BOOLEAN,
    application_key BOOLEAN,
    revoked_at      TIMESTAMPTZ
);
