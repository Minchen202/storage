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

(Desktop app tasks are unchanged)

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

(Desktop app tasks are unchanged)

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

(Desktop app tasks are unchanged)

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

(Desktop app tasks are unchanged)

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

(Desktop app tasks are unchanged)

---

(Database schema and other sections remain the same)
