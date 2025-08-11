-- YouTube Clone Database Schema

-- YouTube Users table
CREATE TABLE IF NOT EXISTS yt_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- YouTube Videos table
CREATE TABLE IF NOT EXISTS yt_videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    user_id INTEGER REFERENCES yt_users(id) ON DELETE CASCADE,
    views INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_yt_videos_user_id ON yt_videos(user_id);
CREATE INDEX IF NOT EXISTS idx_yt_videos_created_at ON yt_videos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_yt_users_username ON yt_users(username);
CREATE INDEX IF NOT EXISTS idx_yt_users_email ON yt_users(email);

-- Update updated_at timestamp automatically
CREATE OR REPLACE FUNCTION update_yt_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_yt_users_updated_at BEFORE UPDATE ON yt_users 
    FOR EACH ROW EXECUTE FUNCTION update_yt_updated_at_column();

CREATE TRIGGER update_yt_videos_updated_at BEFORE UPDATE ON yt_videos 
    FOR EACH ROW EXECUTE FUNCTION update_yt_updated_at_column();