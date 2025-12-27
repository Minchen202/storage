-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    storage_contributed BIGINT DEFAULT 10737418240, -- 10GB in bytes
    storage_quota BIGINT DEFAULT 10737418240, -- 10GB in bytes
    total_uptime_seconds BIGINT DEFAULT 0,
    last_online TIMESTAMP
);

-- Files table
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    filename VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    chunk_count INTEGER NOT NULL,
    encryption_key_encrypted TEXT NOT NULL, -- encrypted with user master key
    created_at TIMESTAMP DEFAULT NOW(),
    last_accessed TIMESTAMP
);

-- Chunks table
CREATE TABLE chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    chunk_hash VARCHAR(64) UNIQUE NOT NULL, -- SHA-256
    chunk_index INTEGER NOT NULL,
    chunk_size INTEGER NOT NULL,
    is_parity BOOLEAN DEFAULT FALSE,
    stored_on_anchor BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Chunk locations (which peers have which chunks)
CREATE TABLE chunk_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chunk_id UUID REFERENCES chunks(id) ON DELETE CASCADE,
    peer_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    verified_at TIMESTAMP,
    UNIQUE(chunk_id, peer_user_id)
);

-- Peer sessions (for tracking online peers)
CREATE TABLE peer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    last_heartbeat TIMESTAMP DEFAULT NOW(),
    is_online BOOLEAN DEFAULT TRUE
);

-- Indexes for performance
CREATE INDEX idx_chunks_file_id ON chunks(file_id);
CREATE INDEX idx_chunk_locations_chunk_id ON chunk_locations(chunk_id);
CREATE INDEX idx_chunk_locations_peer ON chunk_locations(peer_user_id);
CREATE INDEX idx_peer_sessions_online ON peer_sessions(is_online, last_heartbeat);
