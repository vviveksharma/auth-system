CREATE TABLE IF NOT EXISTS message_tbl (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email       VARCHAR(255) NOT NULL,
    tenant_id        UUID         NOT NULL,
    "current_role"   VARCHAR(100) NOT NULL,
    "requested_role" VARCHAR(100) NOT NULL,
    status           VARCHAR(50)  NOT NULL DEFAULT 'pending',
    request_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    action           BOOLEAN      NOT NULL DEFAULT false
);
