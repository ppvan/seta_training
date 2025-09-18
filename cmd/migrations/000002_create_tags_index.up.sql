
-- Create GIN index on tags array for optimized tag searches (idempotent)
CREATE INDEX IF NOT EXISTS idx_posts_tags_gin ON posts USING GIN (tags);