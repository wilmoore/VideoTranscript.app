-- Create metrics and performance tracking tables

-- System metrics table for real-time dashboard stats
CREATE TABLE system_metrics (
    id SERIAL PRIMARY KEY,
    metric_name TEXT NOT NULL,
    metric_value NUMERIC NOT NULL,
    metric_type TEXT NOT NULL, -- 'counter', 'gauge', 'histogram'
    tags JSONB DEFAULT '{}',
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Performance metrics for job processing
CREATE TABLE job_performance (
    id SERIAL PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    stage TEXT NOT NULL, -- 'download', 'extract', 'transcribe', 'complete'
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    duration_seconds NUMERIC,
    memory_usage_mb INTEGER,
    cpu_usage_percent NUMERIC,
    success BOOLEAN DEFAULT false,
    error_message TEXT,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- API usage tracking for dashboard metrics
CREATE TABLE api_usage (
    id SERIAL PRIMARY KEY,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    user_id TEXT,
    api_key_id TEXT,
    response_status INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    request_size_bytes INTEGER DEFAULT 0,
    response_size_bytes INTEGER DEFAULT 0,
    user_agent TEXT,
    ip_address INET,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Business metrics for revenue tracking
CREATE TABLE business_metrics (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    total_jobs INTEGER DEFAULT 0,
    completed_jobs INTEGER DEFAULT 0,
    failed_jobs INTEGER DEFAULT 0,
    revenue_usd NUMERIC(10,2) DEFAULT 0.00,
    unique_users INTEGER DEFAULT 0,
    avg_processing_time_seconds NUMERIC DEFAULT 0,
    total_video_duration_seconds BIGINT DEFAULT 0,
    storage_used_bytes BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(date)
);

-- System health metrics
CREATE TABLE system_health (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    status TEXT NOT NULL, -- 'healthy', 'warning', 'critical'
    uptime_seconds BIGINT DEFAULT 0,
    memory_usage_mb INTEGER DEFAULT 0,
    cpu_usage_percent NUMERIC DEFAULT 0,
    disk_usage_percent NUMERIC DEFAULT 0,
    active_connections INTEGER DEFAULT 0,
    queue_size INTEGER DEFAULT 0,
    last_error TEXT,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_system_metrics_name_time ON system_metrics(metric_name, recorded_at DESC);
CREATE INDEX idx_system_metrics_type ON system_metrics(metric_type);
CREATE INDEX idx_job_performance_job_id ON job_performance(job_id);
CREATE INDEX idx_job_performance_stage ON job_performance(stage);
CREATE INDEX idx_job_performance_recorded_at ON job_performance(recorded_at DESC);
CREATE INDEX idx_api_usage_endpoint ON api_usage(endpoint);
CREATE INDEX idx_api_usage_recorded_at ON api_usage(recorded_at DESC);
CREATE INDEX idx_api_usage_user_id ON api_usage(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_business_metrics_date ON business_metrics(date DESC);
CREATE INDEX idx_system_health_service ON system_health(service_name);
CREATE INDEX idx_system_health_recorded_at ON system_health(recorded_at DESC);

-- Views for common dashboard queries
CREATE VIEW current_job_stats AS
SELECT
    COUNT(*) as total_jobs,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_jobs,
    COUNT(CASE WHEN status = 'running' THEN 1 END) as running_jobs,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_jobs,
    COUNT(CASE WHEN status = 'error' THEN 1 END) as failed_jobs,
    AVG(CASE WHEN status = 'completed' AND completed_at IS NOT NULL
        THEN EXTRACT(EPOCH FROM (completed_at - created_at)) END) as avg_processing_time_seconds,
    COUNT(CASE WHEN created_at >= CURRENT_DATE THEN 1 END) as jobs_today,
    COUNT(CASE WHEN created_at >= CURRENT_DATE - INTERVAL '7 days' THEN 1 END) as jobs_this_week
FROM jobs;

CREATE VIEW recent_performance AS
SELECT
    j.id,
    j.status,
    j.progress,
    j.title,
    j.video_id,
    j.created_at,
    j.completed_at,
    EXTRACT(EPOCH FROM (COALESCE(j.completed_at, NOW()) - j.created_at)) as total_duration_seconds,
    jp.stage,
    jp.duration_seconds as stage_duration_seconds
FROM jobs j
LEFT JOIN job_performance jp ON j.id = jp.job_id
WHERE j.created_at >= NOW() - INTERVAL '24 hours'
ORDER BY j.created_at DESC;

-- Function to update business metrics daily
CREATE OR REPLACE FUNCTION update_business_metrics()
RETURNS void AS $$
DECLARE
    target_date DATE := CURRENT_DATE;
    job_stats RECORD;
    revenue_per_job NUMERIC := 2.50;
BEGIN
    -- Get job statistics for the target date
    SELECT
        COUNT(*) as total,
        COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
        COUNT(CASE WHEN status = 'error' THEN 1 END) as failed,
        COUNT(DISTINCT CASE WHEN created_at >= target_date THEN
            EXTRACT('user_id', url) -- This would need proper user tracking
        END) as unique_users,
        AVG(CASE WHEN status = 'completed' AND completed_at IS NOT NULL
            THEN EXTRACT(EPOCH FROM (completed_at - created_at)) END) as avg_processing_time,
        SUM(CASE WHEN status = 'completed'
            THEN COALESCE(EXTRACT(EPOCH FROM duration::INTERVAL), 0)
            ELSE 0 END) as total_video_seconds
    INTO job_stats
    FROM jobs
    WHERE created_at >= target_date AND created_at < target_date + INTERVAL '1 day';

    -- Insert or update business metrics
    INSERT INTO business_metrics (
        date, total_jobs, completed_jobs, failed_jobs,
        revenue_usd, unique_users, avg_processing_time_seconds,
        total_video_duration_seconds, updated_at
    ) VALUES (
        target_date,
        COALESCE(job_stats.total, 0),
        COALESCE(job_stats.completed, 0),
        COALESCE(job_stats.failed, 0),
        COALESCE(job_stats.completed, 0) * revenue_per_job,
        COALESCE(job_stats.unique_users, 0),
        COALESCE(job_stats.avg_processing_time, 0),
        COALESCE(job_stats.total_video_seconds, 0),
        NOW()
    )
    ON CONFLICT (date) DO UPDATE SET
        total_jobs = EXCLUDED.total_jobs,
        completed_jobs = EXCLUDED.completed_jobs,
        failed_jobs = EXCLUDED.failed_jobs,
        revenue_usd = EXCLUDED.revenue_usd,
        unique_users = EXCLUDED.unique_users,
        avg_processing_time_seconds = EXCLUDED.avg_processing_time_seconds,
        total_video_duration_seconds = EXCLUDED.total_video_duration_seconds,
        updated_at = NOW();
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update business metrics when jobs change
CREATE OR REPLACE FUNCTION trigger_update_business_metrics()
RETURNS trigger AS $$
BEGIN
    -- Update business metrics for the affected date
    PERFORM update_business_metrics();
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER jobs_business_metrics_trigger
    AFTER INSERT OR UPDATE OR DELETE ON jobs
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_update_business_metrics();