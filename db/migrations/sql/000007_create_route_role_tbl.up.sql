CREATE TABLE IF NOT EXISTS route_role_tbl (
    id          UUID   PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name   TEXT,
    tenant_id   UUID   NOT NULL,
    role_id     UUID,
    permissions JSONB,
    routes      TEXT[]
);
