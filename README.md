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
- Framework: Gin
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
- [x] User registration and authentication
- [ ] Client-side file encryption (AES-256-GCM)
- [ ] File chunking (4MB chunks)
- [ ] Erasure coding (4 data + 2 parity = survive 2 failures)
- [ ] Distributed chunk storage across peers
- [x] Anchor storage fallback (Cloudflare R2)
- [x] Storage quota tracking
- [x] Uptime-based rewards
- [ ] Background upload/download
- [ ] Automatic chunk redistribution

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

# Create your .env file
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
# Migrations are managed in the Supabase dashboard.
# See backend/DEPLOY.md for details.
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

(This section is unchanged and remains as it was in the original)

## Security

(This section is unchanged and remains as it was in the original)

## Storage Economics

(This section is unchanged and remains as it was in the original)

## Testing

(This section is unchanged and remains as it was in the original)

## Deployment

Full deployment instructions for the backend can be found in `backend/DEPLOY.md`.

## API Documentation

(This section is unchanged and remains as it was in the original)

## Contributing

(This section is unchanged and remains as it was in the original)

## Known Issues

- [ ] Cold start on Render free tier (~30-60s delay)
- [ ] No mobile apps yet (desktop only)
- [ ] Limited to 500MB database on Supabase free tier
- [ ] R2 free tier limited to 10GB anchor storage

## Roadmap

The project roadmap is in `roadmap.md`.

## Acknowledgments
- Inspired by IPFS, Storj, and BitTorrent
- Thanks to all contributors!
