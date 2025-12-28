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
	// Setup structured logging with DEBUG level for troubleshooting
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	log.Info("=== P2P Storage Backend Starting ===")

	// Log environment variable status (not values for security)
	envVars := []string{"DATABASE_URL", "JWT_SECRET", "R2_ACCOUNT_ID", "R2_ACCESS_KEY_ID", "R2_SECRET_ACCESS_KEY", "R2_BUCKET_NAME", "PORT"}
	for _, env := range envVars {
		if os.Getenv(env) != "" {
			log.WithField("env_var", env).Debug("Environment variable is set")
		} else {
			log.WithField("env_var", env).Warn("Environment variable is NOT set")
		}
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	log.Debug("Connecting to database...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.WithError(err).Fatal("Failed to open database connection")
	}
	defer db.Close()

	// Test database connection
	log.Debug("Testing database connection with ping...")
	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("Failed to ping database - check DATABASE_URL")
	}
	log.Info("Database connection established successfully")

	log.Debug("Initializing R2 storage client...")
	r2Client, err := storage.NewR2Client()
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize R2 client")
	}
	log.Info("R2 storage client initialized successfully")

	env := &api.Env{DB: db, R2Client: r2Client}

	r := gin.New()
	// Custom recovery middleware with logging
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(logrus.Fields{
					"error":  err,
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
					"ip":     c.ClientIP(),
				}).Error("PANIC RECOVERED - this may indicate a bug")
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			}
		}()
		c.Next()
	})
	// Custom logger middleware with comprehensive request/response logging
	r.Use(func(c *gin.Context) {
		start := time.Now()

		// Log incoming request
		log.WithFields(logrus.Fields{
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"ip":           c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"content_type": c.ContentType(),
			"has_auth":     c.GetHeader("Authorization") != "",
		}).Debug("Incoming request")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Log response with appropriate level
		fields := logrus.Fields{
			"status":     status,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency_ms": latency.Milliseconds(),
			"user_agent": c.Request.UserAgent(),
		}

		if status >= 500 {
			log.WithFields(fields).Error("Request completed with server error")
		} else if status >= 400 {
			log.WithFields(fields).Warn("Request completed with client error")
		} else {
			log.WithFields(fields).Info("Request completed successfully")
		}
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
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				return
			}
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.WithField("port", port).Info("=== Server starting on port " + port + " ===")
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.WithError(err).Fatal("Server failed to start")
	}
}
