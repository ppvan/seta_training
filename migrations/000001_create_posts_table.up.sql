-- Create posts table with basic fields (idempotent)
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    tags TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Example usage queries:
-- Find posts containing a specific tag
-- SELECT * FROM posts WHERE 'postgresql' = ANY(tags);

-- Find posts containing any of multiple tags
-- SELECT * FROM posts WHERE tags && ARRAY['postgresql', 'database'];

-- Find posts containing all specified tags
-- SELECT * FROM posts WHERE tags @> ARRAY['postgresql', 'performance'];