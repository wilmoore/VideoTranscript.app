-- Expand jobs table to match dashboard Job struct
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS video_id TEXT;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS title TEXT;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS progress INTEGER DEFAULT 0;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS update_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS log_file TEXT;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS output_dir TEXT;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS duration TEXT;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS file_count INTEGER DEFAULT 0;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS file_size TEXT;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS stage TEXT;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS stage_progress JSONB;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS category_class TEXT DEFAULT 'entertainment';
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS category_icon TEXT DEFAULT '🎬';
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS status_text TEXT;

-- Add indexes for dashboard queries
CREATE INDEX IF NOT EXISTS idx_jobs_video_id ON jobs(video_id);
CREATE INDEX IF NOT EXISTS idx_jobs_progress ON jobs(progress);
CREATE INDEX IF NOT EXISTS idx_jobs_update_time ON jobs(update_time);

-- Update existing columns with better defaults
UPDATE jobs SET
    video_id = COALESCE(video_id, 'unknown'),
    title = COALESCE(title, 'Loading...'),
    progress = COALESCE(progress, 0),
    start_time = COALESCE(start_time, created_at),
    update_time = COALESCE(update_time, created_at),
    log_file = COALESCE(log_file, 'logs/' || id || '.log'),
    output_dir = COALESCE(output_dir, 'transcripts/' || COALESCE(video_id, 'unknown')),
    duration = COALESCE(duration, '00:00'),
    file_count = COALESCE(file_count, 0),
    file_size = COALESCE(file_size, '0 KB'),
    stage = COALESCE(stage, ''),
    category_class = COALESCE(category_class, 'entertainment'),
    category_icon = COALESCE(category_icon, '🎬'),
    status_text = CASE
        WHEN status = 'completed' THEN 'Transcription complete'
        WHEN status = 'failed' THEN 'Processing failed'
        WHEN status = 'queued' THEN 'Queued for processing'
        WHEN status = 'running' THEN 'Processing in progress'
        ELSE status
    END
WHERE video_id IS NULL OR title IS NULL OR start_time IS NULL;