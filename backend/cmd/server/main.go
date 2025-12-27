package main

import (
	"database/sql"
	"log"
	"os"
	"p2p-storage/backend/internal/api"
	"p2p-storage/backend/internal/auth"
	"p2p-storage/backend/internal/storage"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r2Client, err := storage.NewR2Client()
	if err != nil {
		log.Fatal(err)
	}

	env := &api.Env{DB: db, R2Client: r2Client}

	r := gin.Default()

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", env.RegisterHandler)
		authRoutes.POST("/login", env.LoginHandler)
	}

	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(auth.AuthMiddleware())
	{
		// Profile endpoint
		protectedRoutes.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{
				"message": "ok",
				"userID":  userID,
			})
		})

		// Peer management endpoints
		protectedRoutes.POST("/peer/heartbeat", env.HeartbeatHandler)

		// File management endpoints
		protectedRoutes.POST("/files", env.CreateFileHandler)
		protectedRoutes.GET("/files", env.GetFilesHandler)

		// Chunk management endpoints
		protectedRoutes.POST("/chunks/upload/:hash", env.UploadChunkHandler)
		protectedRoutes.GET("/chunks/:hash", env.DownloadChunkHandler)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.Run(":8080")
}
