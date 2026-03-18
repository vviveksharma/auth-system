-- Drop indexes first
DROP INDEX IF EXISTS idx_provider_daily_provider;
DROP INDEX IF EXISTS idx_provider_daily_org;
DROP INDEX IF EXISTS idx_project_monthly_project;
DROP INDEX IF EXISTS idx_project_daily_org;
DROP INDEX IF EXISTS idx_project_daily_project;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS provider_daily_stats;
DROP TABLE IF EXISTS project_monthly_stats;
DROP TABLE IF EXISTS project_daily_stats;
