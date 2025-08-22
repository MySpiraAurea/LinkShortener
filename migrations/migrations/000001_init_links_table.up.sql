CREATE TABLE IF NOT EXISTS links (
    short_id VARCHAR(6) PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);