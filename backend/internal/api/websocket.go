package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// In a real application, you'd want a more secure origin check
		return true
	},
}

func (env *Env) WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("WebSocket connection without authentication")
		return
	}

	log.Printf("WebSocket connection established for user %s", userID)

	// In a real application, you would manage the connection here,
	// such as adding it to a pool of active peer connections.

	for {
		// Read messages from the peer (e.g., responses to chunk requests)
		// For now, we'll just keep the connection alive.
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for user %s: %v", userID, err)
			break
		}
	}

	log.Printf("WebSocket connection closed for user %s", userID)
}
