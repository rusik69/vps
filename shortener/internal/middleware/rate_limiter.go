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

var (
	window   = time.Minute
	maxReq   = MaxRequestsPerMinute
	requests = make(map[string][]time.Time)
	mu       sync.RWMutex
)

// RateLimitMiddleware limits requests per IP
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		// Handle cases where IP is empty or invalid
		if ip == "" || ip == "::" || ip == "::1" {
			ip = "127.0.0.1"
		}
		
		mu.Lock()
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
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		
		// Add current request
		requests[ip] = append(requests[ip], now)
		mu.Unlock()
		
		c.Next()
	}
}

// ResetRateLimiter resets the rate limiter state for testing purposes.
func ResetRateLimiter() {
	mu.Lock()
	requests = make(map[string][]time.Time)
	mu.Unlock()
}
