package api

import (
	"net/http"
	"p2p-storage/backend/internal/db"

	"github.com/gin-gonic/gin"
)

func (env *Env) GetUserStatsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := db.GetUserStats(env.DB, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
