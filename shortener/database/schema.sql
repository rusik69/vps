-- URL Shortener Database Schema

-- Create shortener users table
CREATE TABLE IF NOT EXISTS shortener_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create short URLs table
CREATE TABLE IF NOT EXISTS shortener_urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id INTEGER REFERENCES shortener_users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    click_count INTEGER DEFAULT 0
);

-- Create clicks table
CREATE TABLE IF NOT EXISTS shortener_clicks (
    id SERIAL PRIMARY KEY,
    short_url_id INTEGER REFERENCES shortener_urls(id),
    user_agent TEXT,
    ip_address INET,
    referrer TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create rate limits table
CREATE TABLE IF NOT EXISTS shortener_rate_limits (
    id SERIAL PRIMARY KEY,
    ip_address INET NOT NULL,
    request_count INTEGER DEFAULT 0,
    reset_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_shortener_urls_short_code ON shortener_urls(short_code);
CREATE INDEX IF NOT EXISTS idx_shortener_clicks_short_url_id ON shortener_clicks(short_url_id);
CREATE INDEX IF NOT EXISTS idx_shortener_rate_limits_ip_address ON shortener_rate_limits(ip_address);
CREATE INDEX IF NOT EXISTS idx_shortener_users_username ON shortener_users(username);

-- Update updated_at timestamp automatically
CREATE OR REPLACE FUNCTION update_shortener_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_shortener_users_updated_at BEFORE UPDATE ON shortener_users 
    FOR EACH ROW EXECUTE FUNCTION update_shortener_updated_at_column();

CREATE TRIGGER update_shortener_urls_updated_at BEFORE UPDATE ON shortener_urls 
    FOR EACH ROW EXECUTE FUNCTION update_shortener_updated_at_column();