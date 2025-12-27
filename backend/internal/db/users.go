package db

import (
	"database/sql"
	"p2p-storage/backend/pkg/models"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func CreateUser(db *sql.DB, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, storage_contributed, storage_quota, total_uptime_seconds, last_online`

	err := db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.StorageContributed,
		&user.StorageQuota,
		&user.TotalUptimeSeconds,
		&user.LastOnline,
	)

	return err
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at, storage_contributed, storage_quota, total_uptime_seconds, last_online FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.StorageContributed,
		&user.StorageQuota,
		&user.TotalUptimeSeconds,
		&user.LastOnline,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
