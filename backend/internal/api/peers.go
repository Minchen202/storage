package api

import (
	"net/http"
	"p2p-storage/backend/internal/db"

	"github.com/gin-gonic/gin"
)

const UptimeRewardInterval = 600 // 10 minutes

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
