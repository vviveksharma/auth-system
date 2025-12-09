-- Test data seeding script for integration tests
-- This creates a test tenant and application key for automated testing

-- Insert test tenant (use fixed UUID for consistency)
-- Password: TestPass123! (bcrypt hash)
INSERT INTO tenant_tbl (id, name, email, campany, password, salt, status, created_at, updated_at)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479'::UUID,
    'Test Tenant',
    'test@integration.test',
    'Test Integration Company',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',  -- TestPass123!
    'integration_salt',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Insert test application key for API endpoints (application_key=true)
INSERT INTO token_tbl (id, tenant_id, name, is_active, application_key, expires_at, created_at, updated_at, usage_count)
VALUES (
    'a1b2c3d4-e5f6-4789-a012-b3c4d5e6f7a8'::UUID,
    'f47ac10b-58cc-4372-a567-0e02b2c3d479'::UUID,
    'Test Application Key',
    true,
    true,  -- This is an application key (used in query params)
    NOW() + INTERVAL '365 days',
    NOW(),
    NOW(),
    0
) ON CONFLICT (id) DO NOTHING;

-- Insert test login token for tenant authorization (application_key=false)
INSERT INTO token_tbl (id, tenant_id, name, is_active, application_key, expires_at, created_at, updated_at, usage_count)
VALUES (
    'b2c3d4e5-f6a7-4890-b123-c4d5e6f7a8b9'::UUID,
    'f47ac10b-58cc-4372-a567-0e02b2c3d479'::UUID,
    'Test Login Token',
    true,
    false,  -- This is a login token (used in Authorization header)
    NOW() + INTERVAL '7 days',
    NOW(),
    NOW(),
    0
) ON CONFLICT (id) DO NOTHING;

-- Insert default roles if they don't exist
INSERT INTO role_tbl (id, tenant_id, role, display_name, role_type, status, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::UUID, 'admin', 'Administrator', 'default', true, NOW(), NOW()),
    (gen_random_uuid(), 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::UUID, 'user', 'User', 'default', true, NOW(), NOW()),
    (gen_random_uuid(), 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::UUID, 'moderator', 'Moderator', 'default', true, NOW(), NOW())
ON CONFLICT DO NOTHING;
