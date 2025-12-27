package db

import (
	"database/sql"
	"p2p-storage/backend/pkg/models"
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

func GetOnlinePeers(db *sql.DB, count int, excludeUserID string) ([]models.User, error) {
	query := `
		SELECT u.id, u.email
		FROM users u
		JOIN peer_sessions ps ON u.id = ps.user_id
		WHERE ps.is_online = TRUE
		  AND u.id != $1
		  AND u.storage_quota > u.storage_contributed -- A simple check for available space
		ORDER BY random()
		LIMIT $2`

	rows, err := db.Query(query, excludeUserID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []models.User
	for rows.Next() {
		var peer models.User
		if err := rows.Scan(&peer.ID, &peer.Email); err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}

	return peers, nil
}
