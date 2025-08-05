package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Reset global state
	mu.Lock()
	requests = make(map[string][]time.Time)
	mu.Unlock()

	router := gin.New()
	router.Use(RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Test normal request
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestRateLimitExceeded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Reset global state
	mu.Lock()
	requests = make(map[string][]time.Time)
	mu.Unlock()

	router := gin.New()
	router.Use(RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Make requests up to the limit
	for i := 0; i < MaxRequestsPerMinute; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// This should be rate limited
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "Rate limit exceeded")
}

func TestRateLimitDifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Reset global state
	mu.Lock()
	requests = make(map[string][]time.Time)
	mu.Unlock()

	router := gin.New()
	router.Use(RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Test that different IPs have separate limits
	ips := []string{"192.168.1.10:12345", "192.168.1.11:12345", "192.168.1.12:12345"}
	
	for _, ip := range ips {
		for i := 0; i < 50; i++ { // Well under the limit
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = ip
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}

	// All IPs should still work
	for _, ip := range ips {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimitCleanup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Reset global state
	mu.Lock()
	requests = make(map[string][]time.Time)
	// Manually add old requests that should be cleaned up
	old := time.Now().Add(-2 * time.Minute)
	requests["192.168.1.20"] = []time.Time{old, old, old}
	mu.Unlock()

	router := gin.New()
	router.Use(RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Make a request - this should clean up old requests
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.20:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify old requests were cleaned up
	mu.Lock()
	assert.Equal(t, 1, len(requests["192.168.1.20"]))
	mu.Unlock()
}

func TestRateLimitConcurrency(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Reset global state
	mu.Lock()
	requests = make(map[string][]time.Time)
	mu.Unlock()

	router := gin.New()
	router.Use(RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Test concurrent requests don't cause race conditions
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "192.168.1.30:12345"
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				// Should not panic or cause race conditions
			}
		}(i)
	}
	wg.Wait()

	// Verify we can still make requests
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.31:12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
