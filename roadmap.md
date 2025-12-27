# Development Plan: P2P Storage System

## Project Overview
Hybrid P2P file storage with centralized coordination. Users contribute storage to get storage. Chunks encrypted client-side, distributed across peers with erasure coding for redundancy.

---

## Tech Stack

### Backend
- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL (Supabase free tier - 500MB)
- **Object Storage:** Cloudflare R2 (10GB free)
- **Hosting:** Render (free tier)
- **Real-time:** WebSockets (gorilla/websocket)

### Desktop App
- **Language:** C# .NET 8+
- **UI:** Avalonia (cross-platform)
- **Networking:** HttpClient + WebSocket libraries
- **Local DB:** SQLite
- **Crypto:** System.Security.Cryptography

### Infrastructure
- **Auth:** JWT tokens
- **API:** RESTful + WebSocket endpoints
- **Encryption:** AES-256-GCM (client-side)
- **Erasure Coding:** reed-solomon library

---

## Phase 1: Foundation (Weeks 1-2)

### Backend Tasks

**Week 1: Core API Setup**
- [x] Set up Go project structure
- [x] Initialize Render deployment config (`render.yaml`)
- [x] Set up Supabase PostgreSQL connection
- [x] Set up Cloudflare R2 SDK integration
- [x] Implement health check endpoint (`/health`)
- [x] Create database schema
- [x] Implement JWT authentication middleware
- [x] Create user registration endpoint (`POST /api/auth/register`)
- [x] Create user login endpoint (`POST /api/auth/login`)
- [x] Create user profile endpoint (`GET /api/user/profile`)

**Week 2: Storage Management**
- [x] Implement storage quota tracking (in schema)
- [x] Create uptime tracking system (ping endpoint)
- [ ] Build contribution calculation logic (partially done)
- [x] Create endpoint to report peer status (`POST /api/peer/heartbeat`)
- [ ] Create endpoint to get storage stats (`GET /api/user/stats`)

### Desktop App Tasks

**Week 1: Project Setup**
- [ ] Create .NET 8 project (Avalonia or WPF)
- [ ] Set up project structure (Models, Services, Views)
- [ ] Create basic UI wireframes (login, dashboard, settings)
- [ ] Implement local SQLite database setup
- [ ] Create encryption/decryption service (AES-256)

**Week 2: Authentication**
- [ ] Build login screen UI
- [ ] Build registration screen UI
- [ ] Implement API client service (HttpClient wrapper)
- [ ] Connect login/register to backend API
- [ ] Store JWT token securely (Windows: DPAPI, cross-platform: keychain)
- [ ] Implement auto-login with stored token

---

## Phase 2: File Upload System (Weeks 3-4)

### Backend Tasks

**Week 3: Chunk Management**
- [x] Create chunk metadata table
- [x] Implement chunk upload endpoint (`POST /api/chunks/upload`)
- [x] Implement file metadata endpoint (`POST /api/files/create`)
- [x] Build peer selection algorithm
- [ ] Create chunk distribution logic

**Week 4: Peer Distribution**
- [ ] Implement chunk assignment to peers (`POST /api/chunks/assign`)
- [x] Create endpoint for peers to download assigned chunks (`GET /api/chunks/:hash`)
- [ ] Build chunk verification system (challenge-response)
- [ ] Implement queue system for offline peers
- [ ] Add background job for chunk redistribution

### Desktop App Tasks

**Week 3: File Chunking**
- [ ] Build file picker UI
- [ ] Implement file encryption (AES-256-GCM with random key)
- [ ] Implement file chunking (4MB chunks)
- [ ] Add erasure coding library (reed-solomon)
  - 4 data chunks + 2 parity = 6 total per 16MB
- [ ] Calculate SHA-256 hash for each chunk
- [ ] Store encryption keys in local DB (encrypted with master key)

**Week 4: Upload Logic**
- [ ] Build upload progress UI (progress bar, status)
- [ ] Implement chunk upload to backend
  - Upload all chunks to server
  - Handle upload failures/retries
- [ ] Create background upload queue
- [ ] Add upload notification system
- [ ] Store file metadata locally

---

## Phase 3: Peer Storage (Weeks 5-6)

### Backend Tasks

**Week 5: Peer Coordination**
- [ ] Build WebSocket server for real-time peer communication
- [x] Implement peer online/offline tracking (via heartbeats)
- [ ] Create peer discovery endpoint (`GET /api/peers/discover`)
- [ ] Build chunk request routing (server requests from peers)
- [ ] Add timeout handling for unresponsive peers

**Week 6: Anchor Storage**
- [x] Implement fallback to anchor (R2) for unavailable chunks
- [ ] Add caching logic (last 30 days uploads stay on anchor)
- [ ] Create background job to clean old anchor chunks
- [ ] Implement redundancy monitoring
- [ ] Add automatic redistribution when peers drop below 5

### Desktop App Tasks

**Week 5: Peer Storage Service**
- [ ] Create local chunk storage directory structure
- [ ] Implement WebSocket connection to backend
- [ ] Build chunk storage handler (receive chunks from server)
- [ ] Add local storage quota management
- [ ] Create settings UI for storage contribution
- [ ] Implement chunk verification responses

**Week 6: Background Services**
- [ ] Build heartbeat service (ping server every 10 min)
- [ ] Implement chunk serving (respond to server requests)
- [ ] Add uptime tracking display in UI
- [ ] Create storage usage visualization
- [ ] Implement auto-start on boot option

---

## Phase 4: File Download System (Weeks 7-8)

### Backend Tasks

**Week 7: Download Orchestration**
- [x] Create file list endpoint (`GET /api/files`)
- [ ] Implement chunk location lookup
- [ ] Build chunk retrieval from peers
- [ ] Create download stream endpoint (`GET /api/files/:id/download`)
- [ ] Add chunk reassembly logic

### Desktop App Tasks

**Week 7: Download UI**
- [ ] Build file browser UI (list uploaded files)
- [ ] Create download button/action
- [ ] Implement chunk download from backend
- [ ] Add progress tracking per file
- [ ] Build chunk reassembly logic

**Week 8: Decryption & Saving**
- [ ] Implement erasure code reconstruction (if chunks missing)
- [ ] Add file decryption after download
- [ ] Create "Save As" dialog
- [ ] Implement download queue
- [ ] Add download history

---

## Phase 5: Polish & Testing (Weeks 9-10)

### Backend Tasks

**Week 9: Monitoring & Logging**
- [ ] Add structured logging
- [ ] Implement error tracking
- [ ] Create admin dashboard endpoint
- [ ] Add database query optimization
- [ ] Implement rate limiting

**Week 10: Security & Testing**
- [ ] Security audit
- [ ] Add CORS configuration
- [ ] Implement request validation
- [ ] Write integration tests
- [ ] Load testing

### Desktop App Tasks

**Week 9: UX Improvements**
- [ ] Add error messages and user feedback
- [ ] Implement retry logic for failed operations
- [ ] Create system tray integration
- [ ] Add notifications for completed uploads/downloads
- [ ] Build settings panel (storage, auto-start, etc.)

**Week 10: Testing & Packaging**
- [ ] End-to-end testing (upload/download cycles)
- [ ] Test with multiple peers simultaneously
- [ ] Test network failure scenarios
- [ ] Create installer (Windows: MSI, macOS: DMG, Linux: AppImage)
- [ ] Write user documentation

---

## Database Schema

### PostgreSQL (Supabase)

```sql
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
```

### SQLite (Desktop App)

```sql
-- Local user info
CREATE TABLE local_user (
    id INTEGER PRIMARY KEY,
    user_uuid TEXT NOT NULL,
    email TEXT NOT NULL,
    jwt_token TEXT,
    master_key_encrypted BLOB -- encrypted with OS keychain
);

-- Local stored chunks (that we're storing for others)
CREATE TABLE stored_chunks (
    id INTEGER PRIMARY KEY,
    chunk_hash TEXT UNIQUE NOT NULL,
    file_path TEXT NOT NULL,
    chunk_size INTEGER NOT NULL,
    stored_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Our uploaded files (cached from server)
CREATE TABLE my_files (
    id INTEGER PRIMARY KEY,
    file_uuid TEXT UNIQUE NOT NULL,
    filename TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    chunk_hashes TEXT NOT NULL, -- JSON array
    encryption_key BLOB NOT NULL, -- encrypted
    uploaded_at TIMESTAMP
);

-- Settings
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT
);
```

---

## API Endpoints Reference

### Authentication
- `POST /api/auth/register` - Create account
- `POST /api/auth/login` - Login, get JWT
- `GET /api/auth/refresh` - Refresh JWT token

### User Management
- `GET /api/user/profile` - Get user info
- `GET /api/user/stats` - Storage quota, uptime, etc.
- `PUT /api/user/settings` - Update storage contribution

### File Operations
- `POST /api/files/create` - Register new file upload
- `GET /api/files` - List user's files
- `GET /api/files/:id` - Get file metadata
- `DELETE /api/files/:id` - Delete file
- `GET /api/files/:id/download` - Download file (streams chunks)

### Chunk Operations
- `POST /api/chunks/upload` - Upload chunk to server
- `GET /api/chunks/:hash` - Download chunk (for peers)
- `POST /api/chunks/verify` - Respond to verification challenge

### Peer Coordination
- `POST /api/peer/heartbeat` - Keep-alive ping
- `GET /api/peers/discover` - Get list of online peers
- `WS /api/peer/connect` - WebSocket for real-time coordination

---

## Key Implementation Details

### Erasure Coding Example (Go)

```go
import "github.com/klauspost/reedsolomon"

// Create encoder: 4 data shards + 2 parity shards
enc, _ := reedsolomon.New(4, 2)

// Split file into 4 chunks
dataShards := [][]byte{chunk1, chunk2, chunk3, chunk4}

// Generate 2 parity chunks
parityShards := make([][]byte, 2)
allShards := append(dataShards, parityShards...)
enc.Encode(allShards)

// Now you have 6 chunks total (4 data + 2 parity)
// Can lose any 2 and still reconstruct
```

### Client-Side Encryption (C#)

```csharp
using System.Security.Cryptography;

// Generate random key for this file
byte[] key = new byte[32]; // 256-bit
RandomNumberGenerator.Fill(key);

byte[] iv = new byte[16]; // 128-bit
RandomNumberGenerator.Fill(iv);

// Encrypt file
using var aes = Aes.Create();
aes.Key = key;
aes.IV = iv;
aes.Mode = CipherMode.GCM;

using var encryptor = aes.CreateEncryptor();
byte[] encryptedData = encryptor.TransformFinalBlock(fileData, 0, fileData.Length);

// Store key encrypted with user's master key
// Master key derived from password using PBKDF2
```

### Peer Selection Algorithm (Go)

```go
func SelectPeersForChunk(chunkID string, count int) ([]Peer, error) {
    // Get online peers with available storage
    peers := GetOnlinePeers()
    
    // Filter: has storage space, not already storing this chunk
    available := FilterPeers(peers, func(p Peer) bool {
        return p.AvailableStorage > ChunkSize && 
               !p.HasChunk(chunkID)
    })
    
    // Shuffle for randomness
    rand.Shuffle(len(available), func(i, j int) {
        available[i], available[j] = available[j], available[i]
    })
    
    // Take first N peers
    if len(available) < count {
        // Not enough peers, use anchor storage
        return available, ErrNotEnoughPeers
    }
    
    return available[:count], nil
}
```

---

## Deployment Checklist

### Render Backend Setup
1. [ ] Create Render account
2. [ ] Create new Web Service (Docker or Go)
3. [ ] Connect GitHub repo (auto-deploy on push)
4. [ ] Add environment variables:
   - `DATABASE_URL` (from Supabase)
   - `R2_ACCOUNT_ID`, `R2_ACCESS_KEY`, `R2_SECRET_KEY`
   - `JWT_SECRET`
5. [ ] Configure health check endpoint
6. [ ] Set up custom domain (optional)

### Supabase Setup
1. [ ] Create Supabase project
2. [ ] Run database migrations (schema above)
3. [ ] Copy connection string to Render env vars
4. [ ] Enable Row Level Security (optional for beta)

### Cloudflare R2 Setup
1. [ ] Create Cloudflare account
2. [ ] Create R2 bucket: `peer-storage-anchor`
3. [ ] Generate API tokens
4. [ ] Configure CORS for your domain

### Desktop App Distribution
1. [ ] Code sign application (Windows/macOS)
2. [ ] Create installer packages
3. [ ] Host installers on GitHub Releases
4. [ ] Create download page

---

## Testing Strategy

### Unit Tests
- Encryption/decryption correctness
- Chunking produces correct sizes
- Erasure coding reconstruction
- JWT token validation

### Integration Tests
- Full upload/download cycle
- Peer selection algorithm
- Chunk redistribution when peer goes offline
- Storage quota enforcement

### Load Tests
- 100 concurrent uploads
- 1000 chunks stored across 50 peers
- Server handles 500 heartbeat requests/min

### Manual Testing Scenarios
1. Upload 100MB file with 3 peers online
2. Download file while 2/5 peers are offline
3. Leave app running for 4 days (verify uptime bonus)
4. Fill storage quota, verify new uploads fail
5. Kill a peer mid-upload, verify redistribution

---

## Milestones & Deliverables

**Week 2:** Auth + basic API works
**Week 4:** Can upload files (to server only)
**Week 6:** Peers can store chunks
**Week 8:** Can download files from peers
**Week 10:** Beta-ready application

---

## Post-Launch Improvements (Future)

- [ ] Web interface (view files from browser)
- [ ] File sharing (generate public links)
- [ ] Compression before encryption
- [ ] Deduplication (same file uploaded by multiple users)
- [ ] Mobile apps (iOS/Android)
- [ ] Admin dashboard (monitor network health)
- [ ] Automatic bandwidth throttling
- [ ] Priority peers (pay for better performance)
- [ ] Folder sync (like Dropbox)

---

## Questions to Resolve

1. **Master key derivation:** Use password-based key (PBKDF2) or force users to save a recovery key?
2. **Peer authentication:** Should peers authenticate each other or trust the server completely?
3. **Chunk size:** 4MB good? Test with different sizes?
4. **Redistribution trigger:** Immediately when dropping below 5 copies, or wait 24h grace period?
5. **Free tier limits:** What happens when user hits Render bandwidth limit mid-month?
