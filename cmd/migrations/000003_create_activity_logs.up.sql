-- Create activity_logs table (idempotent)
CREATE TABLE IF NOT EXISTS activity_logs (
    id SERIAL PRIMARY KEY,
    action VARCHAR(255) NOT NULL,
    post_id INTEGER,
    logged_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE activity_logs ADD CONSTRAINT fk_activity_logs_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE SET NULL;