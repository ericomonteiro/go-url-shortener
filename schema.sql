-- Create links table
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    redirect_code VARCHAR(50) NOT NULL,
    destiny_url VARCHAR(2000) NOT NULL,
    clicks INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on redirect_code for faster lookups
CREATE INDEX idx_redirect_code ON links(redirect_code);
