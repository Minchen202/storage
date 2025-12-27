package models

import (
	"time"
)

type File struct {
	ID                      string    `json:"id"`
	UserID                  string    `json:"user_id"`
	Filename                string    `json:"filename"`
	FileSize                int64     `json:"file_size"`
	ChunkCount              int       `json:"chunk_count"`
	EncryptionKeyEncrypted  string    `json:"-"` // Do not expose
	CreatedAt               time.Time `json:"created_at"`
	LastAccessed            time.Time `json:"last_accessed"`
}
