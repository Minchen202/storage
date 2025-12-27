package db

import (
	"database/sql"
	"p2p-storage/backend/pkg/models"
)

func CreateChunk(db *sql.DB, chunk *models.Chunk) error {
	query := `
		INSERT INTO chunks (file_id, chunk_hash, chunk_index, chunk_size, is_parity, stored_on_anchor)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return db.QueryRow(
		query,
		chunk.FileID,
		chunk.ChunkHash,
		chunk.ChunkIndex,
		chunk.ChunkSize,
		chunk.IsParity,
		chunk.StoredOnAnchor,
	).Scan(
		&chunk.ID,
		&chunk.CreatedAt,
	)
}

func AssignChunkToPeer(db *sql.DB, chunkID, peerID string) error {
	query := `
		INSERT INTO chunk_locations (chunk_id, peer_user_id)
		VALUES ($1, $2)`

	_, err := db.Exec(query, chunkID, peerID)
	return err
}

func GetChunkLocations(db *sql.DB, chunkID string) ([]string, error) {
	query := `SELECT peer_user_id FROM chunk_locations WHERE chunk_id = $1`
	rows, err := db.Query(query, chunkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peerIDs []string
	for rows.Next() {
		var peerID string
		if err := rows.Scan(&peerID); err != nil {
			return nil, err
		}
		peerIDs = append(peerIDs, peerID)
	}
	return peerIDs, nil
}
