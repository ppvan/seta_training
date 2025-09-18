-- 1. Drop the GIN index
DROP INDEX IF EXISTS idx_posts_search;

-- 2. Drop the generated tsvector column
ALTER TABLE posts
DROP COLUMN IF EXISTS search_vector;
