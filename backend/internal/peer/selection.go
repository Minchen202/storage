package peer

import (
	"fmt"
	"p2p-storage/backend/pkg/models"
)

// SelectPeersForChunk simulates selecting the best peers for storing a chunk.
// In a real application, this would involve checking peer availability, storage space, and location.
func SelectPeersForChunk(chunkHash string, requiredPeers int) ([]models.User, error) {
	fmt.Printf("Simulating peer selection for chunk %s\n", chunkHash)

	// In a real application, this would query the database for online, available peers.
	// For now, we'll return a dummy list of peers.
	if requiredPeers > 5 {
		return nil, fmt.Errorf("cannot select more than 5 dummy peers")
	}

	dummyPeers := make([]models.User, requiredPeers)
	for i := 0; i < requiredPeers; i++ {
		dummyPeers[i] = models.User{
			ID:    fmt.Sprintf("dummy-peer-%d", i+1),
			Email: fmt.Sprintf("peer%d@example.com", i+1),
		}
	}

	return dummyPeers, nil
}
