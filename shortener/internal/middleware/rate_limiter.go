package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MaxRequestsPerMinute = 100
	MaxRequestsPerHour   = 1000
)

// RateLimiter implements rate limiting middleware
type RateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	window   time.Duration
	maxReq   int
}

// RateLimitMiddleware limits requests per IP
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		mu.Lock()
		defer mu.Unlock()
		
		now := time.Now()
		
		// Clean old requests
		if times, exists := requests[ip]; exists {
			validTimes := make([]time.Time, 0)
			for _, t := range times {
				if now.Sub(t) < window {
					validTimes = append(validTimes, t)
				}
			}
			requests[ip] = validTimes
		}
		
		// Check rate limit
		if len(requests[ip]) >= maxReq {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		
		// Add current request
		requests[ip] = append(requests[ip], now)
		
		c.Next()
	}
}

var (
	window = time.Minute
	maxReq = 10
	requests = make(map[string][]time.Time)
	mu       sync.RWMutex
)
