package db

import (
	"database/sql"
	"p2p-storage/backend/pkg/models"
)

func CreateFile(db *sql.DB, file *models.File) error {
	query := `
		INSERT INTO files (user_id, filename, file_size, chunk_count, encryption_key_encrypted)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, last_accessed`

	return db.QueryRow(
		query,
		file.UserID,
		file.Filename,
		file.FileSize,
		file.ChunkCount,
		file.EncryptionKeyEncrypted,
	).Scan(
		&file.ID,
		&file.CreatedAt,
		&file.LastAccessed,
	)
}

func GetFilesByUser(db *sql.DB, userID string) ([]models.File, error) {
	query := `SELECT id, user_id, filename, file_size, chunk_count, created_at, last_accessed FROM files WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(
			&file.ID,
			&file.UserID,
			&file.Filename,
			&file.FileSize,
			&file.ChunkCount,
			&file.CreatedAt,
			&file.LastAccessed,
		); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
