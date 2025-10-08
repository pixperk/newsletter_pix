-- Create subscribers table
CREATE TABLE IF NOT EXISTS subscribers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create email verification table
CREATE TABLE IF NOT EXISTS email_verifications (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    verification_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_subscribers_email ON subscribers(email);
CREATE INDEX IF NOT EXISTS idx_email_verifications_email ON email_verifications(email);
CREATE INDEX IF NOT EXISTS idx_email_verifications_token ON email_verifications(verification_token);

-- Insert some sample data (optional)
-- INSERT INTO subscribers (email) VALUES ('test@example.com') ON CONFLICT DO NOTHING;