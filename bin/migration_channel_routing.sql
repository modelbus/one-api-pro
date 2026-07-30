-- Migration: Add channel routing fields
-- This migration adds fields for channel concurrency, cooldown, RPM, and error tracking

ALTER TABLE channels ADD COLUMN IF NOT EXISTS max_concurrency int DEFAULT 0;
ALTER TABLE channels ADD COLUMN IF NOT EXISTS cooldown_seconds int DEFAULT 60;
ALTER TABLE channels ADD COLUMN IF NOT EXISTS rpm int DEFAULT 0;
ALTER TABLE channels ADD COLUMN IF NOT EXISTS last_error varchar(512) DEFAULT '';
ALTER TABLE channels ADD COLUMN IF NOT EXISTS last_error_time bigint DEFAULT 0;

-- Remove error_config column (was JSON field, replaced by global ErrorNext config)
-- ALTER TABLE channels DROP COLUMN IF EXISTS error_config;

ALTER TABLE logs ADD COLUMN IF NOT EXISTS session_key varchar(128) DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_logs_session_key ON logs(session_key);