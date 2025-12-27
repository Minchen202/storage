package api

import (
	"net/http"
	"p2p-storage/backend/internal/db"

	"github.com/gin-gonic/gin"
)

const UptimeRewardInterval = 600 // 10 minutes in seconds

func (env *Env) HeartbeatHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := db.CreateOrUpdatePeerSession(env.DB, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update peer session"})
		return
	}

	if err := db.UpdateUptime(env.DB, userID.(string), UptimeRewardInterval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update uptime"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Heartbeat received"})
}

func (env *Env) DiscoverPeersHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	peers, err := db.GetOnlinePeers(env.DB, 50, userID.(string)) // Return up to 50 peers
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to discover peers"})
		return
	}

	c.JSON(http.StatusOK, peers)
}
