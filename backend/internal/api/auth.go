package api

import (
	"net/http"
	"p2p-storage/backend/internal/auth"
	"p2p-storage/backend/internal/db"
	"p2p-storage/backend/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (env *Env) RegisterHandler(c *gin.Context) {
	log := logrus.WithFields(logrus.Fields{
		"handler": "RegisterHandler",
		"ip":      c.ClientIP(),
	})
	log.Debug("Register handler started")

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).WithField("content_type", c.ContentType()).Error("Failed to bind register request - invalid JSON or missing fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.WithField("email", req.Email).Debug("Register request parsed successfully")

	log.Debug("Hashing password with bcrypt...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("Failed to hash password - bcrypt error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	log.Debug("Password hashed successfully")

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	log.Debug("Attempting to create user in database...")
	if err := db.CreateUser(env.DB, user); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.WithFields(logrus.Fields{
				"pq_code":    pqErr.Code,
				"pq_name":    pqErr.Code.Name(),
				"pq_message": pqErr.Message,
				"pq_detail":  pqErr.Detail,
			}).Warn("PostgreSQL error during user creation")
			if pqErr.Code.Name() == "unique_violation" {
				c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
				return
			}
		}
		log.WithError(err).Error("Failed to create user in database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User registered successfully")
	c.JSON(http.StatusCreated, user)
}

func (env *Env) LoginHandler(c *gin.Context) {
	log := logrus.WithFields(logrus.Fields{
		"handler": "LoginHandler",
		"ip":      c.ClientIP(),
	})
	log.Debug("Login handler started")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).WithField("content_type", c.ContentType()).Error("Failed to bind login request - invalid JSON or missing fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.WithField("email", req.Email).Debug("Login request parsed successfully")

	log.Debug("Looking up user by email in database...")
	user, err := db.GetUserByEmail(env.DB, req.Email)
	if err != nil {
		log.WithError(err).WithField("email", req.Email).Warn("Failed to get user by email - user not found or database error")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	log.WithField("user_id", user.ID).Debug("User found in database")

	log.Debug("Comparing password hash...")
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.WithError(err).WithField("email", req.Email).Warn("Invalid password - bcrypt comparison failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	log.Debug("Password verified successfully")

	log.Debug("Generating JWT token...")
	token, err := auth.GenerateJWT(user)
	if err != nil {
		log.WithError(err).WithField("user_id", user.ID).Error("Failed to generate JWT - this is a server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	log.Debug("JWT token generated successfully")

	log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User logged in successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}
