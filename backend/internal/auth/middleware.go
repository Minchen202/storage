package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logrus.WithFields(logrus.Fields{
			"middleware": "AuthMiddleware",
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"ip":         c.ClientIP(),
		})
		log.Debug("Auth middleware processing request")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		log.Debug("Authorization header present")

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Warn("Bearer prefix missing in Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}
		log.WithField("token_length", len(tokenString)).Debug("Bearer token extracted")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.WithField("alg", token.Header["alg"]).Warn("Unexpected signing method")
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil {
			log.WithError(err).Warn("JWT parsing failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			log.Warn("JWT token is not valid (may be expired or malformed)")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		log.WithField("user_id", claims.UserID).Debug("JWT validated successfully, user authenticated")
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
