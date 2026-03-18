CREATE TABLE IF NOT EXISTS organisation_tbl (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL,
    name        TEXT,
    slug        TEXT        UNIQUE,
    description TEXT,
    icon_url    TEXT,
    plan        VARCHAR(50) NOT NULL DEFAULT 'free',
    status      VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
