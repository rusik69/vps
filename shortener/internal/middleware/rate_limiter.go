package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/vps/url-shortener/internal/db"
)

const (
	// Rate limits
	MaxRequestsPerMinute = 100
	MaxRequestsPerHour   = 1000
)

// RateLimiter implements rate limiting middleware
type RateLimiter struct {
	db db.Repository
}

// NewRateLimiter creates a new rate limiter
type RateLimiterConfig struct {
	DB db.Repository
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		db: config.DB,
	}
}

// RateLimiterMiddleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		// Check minute rate limit
		if err := rl.checkRateLimit(ip, time.Minute, MaxRequestsPerMinute); err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Too many requests. Please wait %s before trying again.", err),
			})
			return
		}

		// Check hour rate limit
		if err := rl.checkRateLimit(ip, time.Hour, MaxRequestsPerHour); err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Too many requests. Please wait %s before trying again.", err),
			})
			return
		}

		c.Next()
	}
}

// checkRateLimit checks the rate limit for a given IP and time window
func (rl *RateLimiter) checkRateLimit(ip string, window time.Duration, maxRequests int) error {
	// Get or create rate limit entry
	rateLimit, err := rl.db.GetOrCreateRateLimit(ip)
	if err != nil {
		return fmt.Errorf("failed to get rate limit: %v", err)
	}

	// Check if we need to reset the counter
	if time.Now().After(rateLimit.ResetAt) {
		rateLimit.RequestCount = 0
		rateLimit.ResetAt = time.Now().Add(window)
		if err := rl.db.UpdateRateLimit(rateLimit); err != nil {
			return fmt.Errorf("failed to update rate limit: %v", err)
		}
	}

	// Check if we've exceeded the limit
	if rateLimit.RequestCount >= int64(maxRequests) {
		return fmt.Sprintf("%v", time.Until(rateLimit.ResetAt))
	}

	// Increment the request count
	rateLimit.RequestCount++
	if err := rl.db.UpdateRateLimit(rateLimit); err != nil {
		return fmt.Errorf("failed to update rate limit: %v", err)
	}

	return nil
}
