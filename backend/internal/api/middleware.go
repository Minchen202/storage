package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter stores a rate limiter for each IP address.
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
}

// NewIPRateLimiter creates a new IPRateLimiter.
func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
	}
}

// AddIP creates a new rate limiter for an IP address.
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Create a new limiter that allows 10 requests per second with a burst of 20.
	limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 20)
	i.ips[ip] = limiter
	return limiter
}

// GetLimiter returns the rate limiter for a given IP address.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}
	return limiter
}

// RateLimitMiddleware is a middleware that limits requests from a single IP.
func RateLimitMiddleware() gin.HandlerFunc {
	limiter := NewIPRateLimiter()
	return func(c *gin.Context) {
		if !limiter.GetLimiter(c.ClientIP()).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
