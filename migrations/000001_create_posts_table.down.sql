DROP INDEX IF EXISTS idx_posts_tags_gin;

-- Drop posts table (idempotent)
DROP TABLE IF EXISTS posts;