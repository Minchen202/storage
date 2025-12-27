package models

import (
	"time"
)

type Chunk struct {
	ID             string    `json:"id"`
	FileID         string    `json:"file_id"`
	ChunkHash      string    `json:"chunk_hash"`
	ChunkIndex     int       `json:"chunk_index"`
	ChunkSize      int       `json:"chunk_size"`
	IsParity       bool      `json:"is_parity"`
	StoredOnAnchor bool      `json:"stored_on_anchor"`
	CreatedAt      time.Time `json:"created_at"`
}
