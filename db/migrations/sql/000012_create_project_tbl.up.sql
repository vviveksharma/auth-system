CREATE TABLE project_tbl (
    id UUID PRIMARY KEY,
    org_id UUID NOT NULL REFERENCES organisation_tbl(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    environment VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,

    -- Constraints
    CONSTRAINT projects_name_length CHECK (char_length(name) >= 2 AND char_length(name) <= 255),
    CONSTRAINT projects_environment_valid CHECK (
        environment IS NULL OR
        environment IN ('production', 'staging', 'development', 'testing')
    )
);

-- Indexes for performance
CREATE INDEX idx_projects_org ON project_tbl(org_id);
CREATE INDEX idx_projects_tenant ON project_tbl(tenant_id);
CREATE INDEX idx_projects_org_tenant ON project_tbl(org_id, tenant_id);
CREATE INDEX idx_projects_org_name ON project_tbl(org_id, name);
CREATE INDEX idx_projects_environment ON project_tbl(environment) WHERE environment IS NOT NULL;
CREATE INDEX idx_projects_created ON project_tbl(created_at DESC);

-- Unique constraint: project name must be unique within organization
CREATE UNIQUE INDEX idx_projects_org_name_unique ON project_tbl(org_id, LOWER(name));
