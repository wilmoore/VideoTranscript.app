-- Drop metrics and performance tracking tables

-- Drop triggers and functions first
DROP TRIGGER IF EXISTS jobs_business_metrics_trigger ON jobs;
DROP FUNCTION IF EXISTS trigger_update_business_metrics();
DROP FUNCTION IF EXISTS update_business_metrics();

-- Drop views
DROP VIEW IF EXISTS recent_performance;
DROP VIEW IF EXISTS current_job_stats;

-- Drop tables in reverse order
DROP TABLE IF EXISTS system_health;
DROP TABLE IF EXISTS business_metrics;
DROP TABLE IF EXISTS api_usage;
DROP TABLE IF EXISTS job_performance;
DROP TABLE IF EXISTS system_metrics;