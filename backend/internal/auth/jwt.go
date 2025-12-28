package auth

import (
	"os"
	"p2p-storage/backend/pkg/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var jwtKey []byte

func init() {
	log := logrus.WithField("component", "jwt_init")
	log.Debug("Initializing JWT configuration...")

	key := os.Getenv("JWT_SECRET")
	if key == "" {
		log.Fatal("JWT_SECRET environment variable is not set - server cannot start without this")
	}

	log.WithField("key_length", len(key)).Debug("JWT_SECRET loaded successfully")
	jwtKey = []byte(key)
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User) (string, error) {
	log := logrus.WithFields(logrus.Fields{
		"function": "GenerateJWT",
		"user_id":  user.ID,
	})
	log.Debug("Generating JWT token...")

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		log.WithError(err).Error("Failed to sign JWT token")
		return "", err
	}

	log.WithFields(logrus.Fields{
		"expires_at":   expirationTime.Format(time.RFC3339),
		"token_length": len(tokenString),
	}).Debug("JWT token generated successfully")

	return tokenString, nil
}
