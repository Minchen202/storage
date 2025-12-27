package api

import (
	"database/sql"
	"p2p-storage/backend/internal/storage"
)

type Env struct {
	DB       *sql.DB
	R2Client *storage.R2Client
}
