# P2P Distributed Storage System

A hybrid peer-to-peer file storage system where users contribute storage to earn storage. Files are encrypted client-side, split into chunks, and distributed across multiple peers with erasure coding for redundancy.

## Core Concept

- **Contribute storage to get storage**: Minimum 10GB contribution = 10GB storage allocation
- **Earn more by being online**: 4 days cumulative uptime = +1GB storage quota
- **Privacy-first**: End-to-end encryption, files encrypted before leaving your device
- **Resilient**: Files stored across 5+ peers with erasure coding (survive 2 peer failures)
- **No single point of failure**: Distributed storage with anchor servers as backup

## Architecture

### Components

```
┌─────────────────┐         ┌──────────────────┐
│  Desktop App    │────────▶│  Backend API     │
│  (C# .NET)      │  HTTPS  │  (Go)            │
│                 │◀────────│                  │
└─────────────────┘         └──────────────────┘
        │                            │
        │ Stores chunks              │
        │ for others                 │ Coordinates
        │                            │ peer storage
        ▼                            ▼
┌─────────────────┐         ┌──────────────────┐
│  Local Storage  │         │  PostgreSQL      │
│  (Encrypted     │         │  (Supabase)      │
│   Chunks)       │         │  + Cloudflare R2 │
└─────────────────┘         └──────────────────┘
```

### Tech Stack

**Backend**
- Language: Go 1.21+
- Framework: Gin or Fiber
- Database: PostgreSQL (Supabase free tier)
- Object Storage: Cloudflare R2
- Hosting: Render (free tier)
- Auth: JWT tokens

**Desktop Application**
- Language: C# .NET 8+
- UI Framework: Avalonia (cross-platform)
- Local DB: SQLite
- Encryption: AES-256-GCM
- Erasure Coding: Reed-Solomon

## Features

### Current (MVP)
- User registration and authentication
- Client-side file encryption (AES-256-GCM)
- File chunking (4MB chunks)
- Erasure coding (4 data + 2 parity = survive 2 failures)
- Distributed chunk storage across peers
- Anchor storage fallback (Cloudflare R2)
- Storage quota tracking
- Uptime-based rewards
- Background upload/download
- Automatic chunk redistribution

### Planned (Post-MVP)
- Web interface
- File sharing (public links)
- Compression before encryption
- Deduplication
- Mobile apps
- Admin dashboard

## Prerequisites

### Backend Development
- Go 1.21 or higher
- PostgreSQL 14+ (or Supabase account)
- Cloudflare account (R2 storage)
- Render account (deployment)

### Desktop App Development
- .NET 8 SDK
- Visual Studio 2022 / Rider / VS Code
- Windows/macOS/Linux for testing

## Setup

### 1. Backend Setup

```bash
# Clone repository
git clone https://github.com/yourusername/p2p-storage.git
cd p2p-storage/backend

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env
```

**Environment Variables (.env)**
```env
# Database
DATABASE_URL=postgresql://user:pass@host:5432/dbname

# Cloudflare R2
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key
R2_SECRET_ACCESS_KEY=your_secret_key
R2_BUCKET_NAME=peer-storage-anchor

# JWT
JWT_SECRET=your-super-secret-key-change-this

# Server
PORT=8080
ENVIRONMENT=development
```

**Run Database Migrations**
```bash
# Using golang-migrate
migrate -path ./migrations -database "$DATABASE_URL" up

# Or manually run SQL from docs/schema.sql
psql $DATABASE_URL < docs/schema.sql
```

**Run Backend**
```bash
# Development
go run cmd/server/main.go

# Production build
go build -o bin/server cmd/server/main.go
./bin/server
```

### 2. Desktop App Setup

```bash
cd ../desktop-app

# Restore dependencies
dotnet restore

# Update appsettings.json with your backend URL
nano appsettings.json
```

**appsettings.json**
```json
{
  "ApiBaseUrl": "http://localhost:8080",
  "LocalStoragePath": "~/.p2p-storage",
  "DefaultStorageContribution": 10737418240,
  "ChunkSize": 4194304
}
```

**Run Desktop App**
```bash
# Development
dotnet run

# Publish for distribution
dotnet publish -c Release -r win-x64 --self-contained
dotnet publish -c Release -r osx-x64 --self-contained
dotnet publish -c Release -r linux-x64 --self-contained
```

## Project Structure

```
p2p-storage/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # Entry point
│   ├── internal/
│   │   ├── api/                     # HTTP handlers
│   │   │   ├── auth.go
│   │   │   ├── files.go
│   │   │   ├── chunks.go
│   │   │   └── peers.go
│   │   ├── db/                      # Database layer
│   │   │   ├── postgres.go
│   │   │   └── queries.go
│   │   ├── storage/                 # R2 integration
│   │   │   └── r2.go
│   │   ├── peer/                    # Peer coordination
│   │   │   ├── selection.go
│   │   │   └── redistribution.go
│   │   └── auth/                    # JWT middleware
│   │       └── jwt.go
│   ├── pkg/
│   │   └── models/                  # Data models
│   ├── migrations/                  # SQL migrations
│   ├── go.mod
│   └── go.sum
├── desktop-app/
│   ├── Models/                      # Data models
│   ├── Services/                    # Business logic
│   │   ├── ApiClient.cs
│   │   ├── EncryptionService.cs
│   │   ├── ChunkService.cs
│   │   └── PeerService.cs
│   ├── ViewModels/                  # MVVM view models
│   ├── Views/                       # UI views
│   │   ├── LoginView.axaml
│   │   ├── DashboardView.axaml
│   │   └── SettingsView.axaml
│   ├── Storage/                     # Local storage
│   │   └── LocalDatabase.cs
│   ├── App.axaml
│   └── Program.cs
├── docs/
│   ├── api.md                       # API documentation
│   ├── schema.sql                   # Database schema
│   ├── architecture.md              # System architecture
│   └── security.md                  # Security considerations
└── README.md
```

## Security

### Encryption
- **Client-side encryption**: Files encrypted before upload using AES-256-GCM
- **Unique keys per file**: Each file gets a random 256-bit encryption key
- **Master key**: User's encryption keys are encrypted with a master key derived from password (PBKDF2)
- **Zero-knowledge**: Server never sees unencrypted data or encryption keys

### Authentication
- JWT tokens with 24-hour expiration
- Secure password hashing (bcrypt, cost factor 12)
- Token refresh mechanism

### Data Integrity
- SHA-256 hashing for chunk verification
- Periodic chunk verification (challenge-response)
- Redundancy ensures data survives peer failures

### Privacy
- Peers don't know whose data they're storing
- Chunk filenames are hashes (no metadata leakage)
- All traffic over HTTPS/WSS

## Storage Economics

### Base System
- Minimum contribution: **10GB**
- Base storage allocation: **10GB**
- Additional contribution = additional allocation (1:1 ratio)

### Uptime Rewards
- 4 days cumulative uptime = **+1GB** storage quota
- Uptime tracked via heartbeat (every 10 minutes)
- Rewards accumulate over time

### Example
```
User contributes: 20GB
Base allocation: 20GB
30 days uptime: +7.5GB bonus
Total storage: 27.5GB
```

### Redundancy Overhead
- Files stored with 5x redundancy (5 copies across peers)
- With erasure coding (4+2): ~1.5x overhead instead
- Network efficiency: 70-80% of contributed storage is usable

## Testing

### Backend Tests
```bash
cd backend

# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# With coverage
go test -cover ./...
```

### Desktop App Tests
```bash
cd desktop-app

# Run all tests
dotnet test

# With coverage
dotnet test /p:CollectCoverage=true
```

### Manual Testing Checklist
- [ ] Register new user
- [ ] Upload file (< 10MB)
- [ ] Upload file (> 100MB)
- [ ] Download file immediately after upload
- [ ] Download file after closing/reopening app
- [ ] Verify storage quota updates
- [ ] Test with 3+ peers online
- [ ] Test with 1 peer online (should use anchor)
- [ ] Leave app running 24 hours (verify uptime tracking)
- [ ] Fill storage quota, verify rejection

## Deployment

### Backend (Render)

1. Create `render.yaml` in repository root:
```yaml
services:
  - type: web
    name: p2p-storage-api
    env: go
    buildCommand: go build -o bin/server cmd/server/main.go
    startCommand: ./bin/server
    envVars:
      - key: DATABASE_URL
        sync: false
      - key: JWT_SECRET
        generateValue: true
      - key: R2_ACCOUNT_ID
        sync: false
      - key: R2_ACCESS_KEY_ID
        sync: false
      - key: R2_SECRET_ACCESS_KEY
        sync: false
```

2. Connect GitHub repository to Render
3. Add environment variables in Render dashboard
4. Deploy automatically on git push

### Database (Supabase)

1. Create project at supabase.com
2. Run migrations from `docs/schema.sql`
3. Copy connection string to Render environment

### Object Storage (Cloudflare R2)

1. Create R2 bucket: `peer-storage-anchor`
2. Generate API tokens
3. Add credentials to Render environment

### Desktop App Distribution

**Windows (MSI Installer)**
```bash
dotnet publish -c Release -r win-x64 --self-contained
# Use WiX Toolset to create MSI
```

**macOS (DMG)**
```bash
dotnet publish -c Release -r osx-x64 --self-contained
# Use create-dmg or similar tool
```

**Linux (AppImage)**
```bash
dotnet publish -c Release -r linux-x64 --self-contained
# Use appimagetool
```

Host installers on GitHub Releases or your own server.

## API Documentation

### Quick Reference

**Authentication**
- `POST /api/auth/register` - Create account
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh token

**Files**
- `POST /api/files` - Create file metadata
- `GET /api/files` - List user files
- `GET /api/files/:id` - Get file details
- `GET /api/files/:id/download` - Download file
- `DELETE /api/files/:id` - Delete file

**Chunks**
- `POST /api/chunks/upload` - Upload chunk
- `GET /api/chunks/:hash` - Download chunk (for peers)

**Peers**
- `POST /api/peer/heartbeat` - Keep-alive ping
- `WS /api/peer/connect` - WebSocket connection

## Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style
- **Go**: Follow official Go style guide, use `gofmt`
- **C#**: Follow Microsoft C# conventions, use `dotnet format`

## Known Issues

- [ ] Cold start on Render free tier (~30-60s delay)
- [ ] No mobile apps yet (desktop only)
- [ ] Limited to 500MB database on Supabase free tier
- [ ] R2 free tier limited to 10GB anchor storage

## Roadmap

### Version 0.1 (MVP) - Current
- [ ] Basic file upload/download
- [ ] Peer storage coordination
- [ ] Uptime tracking
- [ ] Quota management

### Version 0.2 (Q1 2026)
- [ ] Web interface
- [ ] File sharing links
- [ ] Improved error handling
- [ ] Admin dashboard

### Version 0.3 (Q2 2026)
- [ ] Mobile apps (iOS/Android)
- [ ] Compression before encryption
- [ ] Deduplication
- [ ] Performance optimizations

### Version 1.0 (Q3 2026)
- [ ] Production-ready stability
- [ ] Automated testing suite
- [ ] Documentation complete
- [ ] Marketing website

## Acknowledgments
- Inspired by IPFS, Storj, and BitTorrent
- Thanks to all contributors!
