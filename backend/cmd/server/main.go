package main

import (
	"database/sql"
	"os"
	"p2p-storage/backend/internal/api"
	"p2p-storage/backend/internal/auth"
	"p2p-storage/backend/internal/storage"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// Setup structured logging
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

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

	r := gin.New()
	r.Use(gin.Recovery())
	// Custom logger middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
		}).Info("request")
	})
	// Apply rate limiting to all requests
	r.Use(api.RateLimitMiddleware())

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", env.RegisterHandler)
		authRoutes.POST("/login", env.LoginHandler)
	}

	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(auth.AuthMiddleware())
	{
		// User endpoints
		protectedRoutes.GET("/user/stats", env.GetUserStatsHandler)
		protectedRoutes.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{
				"message": "ok",
				"userID":  userID,
			})
		})


		// Peer management endpoints
		protectedRoutes.POST("/peer/heartbeat", env.HeartbeatHandler)
		protectedRoutes.GET("/peers/discover", env.DiscoverPeersHandler)
		protectedRoutes.GET("/peer/connect", env.WebSocketHandler) // WebSocket endpoint

		// File management endpoints
		protectedRoutes.POST("/files", env.CreateFileHandler)
		protectedRoutes.GET("/files", env.GetFilesHandler)
		protectedRoutes.DELETE("/files/:id", env.DeleteFileHandler)

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
