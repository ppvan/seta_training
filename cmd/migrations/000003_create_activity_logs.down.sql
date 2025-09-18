-- Reverse SQL: Drop everything (idempotent)
DROP INDEX IF EXISTS idx_activity_logs_post_logged;
DROP INDEX IF EXISTS idx_activity_logs_logged_at;
DROP INDEX IF EXISTS idx_activity_logs_post_id;
DROP TABLE IF EXISTS activity_logs;