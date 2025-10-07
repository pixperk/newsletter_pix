-- Create subscribers table
CREATE TABLE IF NOT EXISTS subscribers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_subscribers_email ON subscribers(email);

-- Insert some sample data (optional)
-- INSERT INTO subscribers (email) VALUES ('test@example.com') ON CONFLICT DO NOTHING;