-- =============================================================
-- Migration 000014 DOWN: Remove all gr.* system roles
-- =============================================================

BEGIN;

DELETE FROM route_role_tbl WHERE role_name LIKE 'gr.%';
DELETE FROM role_tbl         WHERE role     LIKE 'gr.%';
DROP FUNCTION IF EXISTS uuid_generate_v7();

COMMIT;
