-- 1. Add a tsvector column (can be generated or manually updated)
ALTER TABLE posts
ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content, ''))
    ) STORED;

-- 2. Create a GIN index on that column
CREATE INDEX IF NOT EXISTS idx_posts_search
    ON posts
    USING GIN (search_vector);
