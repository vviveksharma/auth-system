CREATE TABLE IF NOT EXISTS db_reset_token (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL,
    tenant_id   UUID         NOT NULL,
    otp_hash    VARCHAR(255) NOT NULL,
    otp_type    VARCHAR(20)  NOT NULL DEFAULT 'numeric',
    reset_token VARCHAR(255),
    expires_at  TIMESTAMPTZ  NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    used_at     TIMESTAMPTZ,
    INDEX idx_reset_token_user_id    (user_id),
    INDEX idx_reset_token_tenant_id  (tenant_id),
    INDEX idx_reset_token_expires_at (expires_at),
    INDEX idx_reset_token_is_active  (is_active)
);
