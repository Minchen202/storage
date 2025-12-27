package peer

import (
	"database/sql"
	"fmt"
	"p2p-storage/backend/internal/db"
	"p2p-storage/backend/pkg/models"
)

// SelectPeersForChunk selects the best peers for storing a chunk.
func SelectPeersForChunk(database *sql.DB, chunkHash string, requiredPeers int, excludeUserID string) ([]models.User, error) {
	peers, err := db.GetOnlinePeers(database, requiredPeers, excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get online peers: %w", err)
	}

	if len(peers) < requiredPeers {
		return nil, fmt.Errorf("not enough online peers available")
	}

	return peers, nil
}
