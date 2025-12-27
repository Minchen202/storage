package db

import (
	"database/sql"
	"time"
)

func CreateOrUpdatePeerSession(db *sql.DB, userID string) error {
	query := `
		INSERT INTO peer_sessions (user_id, last_heartbeat, is_online)
		VALUES ($1, $2, TRUE)
		ON CONFLICT (user_id) DO UPDATE SET
			last_heartbeat = $2,
			is_online = TRUE`

	_, err := db.Exec(query, userID, time.Now())
	return err
}

func UpdateUptime(db *sql.DB, userID string, secondsToAdd int64) error {
	query := `
		UPDATE users
		SET total_uptime_seconds = total_uptime_seconds + $1
		WHERE id = $2`

	_, err := db.Exec(query, secondsToAdd, userID)
	return err
}
