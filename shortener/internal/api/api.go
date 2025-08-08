package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/shortener/internal/middleware"
	"github.com/rusik69/shortener/internal/service"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, svc service.Service) {
	// Trust all proxies
	r.ForwardedByClientIP = true
	
	// Apply rate limiting middleware
	r.Use(middleware.RateLimitMiddleware())
	
	// Serve static files and load templates (skip in test mode)
	if gin.Mode() != gin.TestMode {
		r.Static("/static", "./web")
		r.LoadHTMLGlob("web/*.html")
	}
	
	// Web frontend
	r.GET("/", func(c *gin.Context) {
		if gin.Mode() == gin.TestMode {
			c.JSON(http.StatusOK, gin.H{"message": "URL Shortener"})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "URL Shortener",
			})
		}
	})
	
	// API routes
	api := r.Group("/api")
	{
		api.POST("/shorten", createShortURL(svc))
		api.GET("/stats/:code", getURLStats(svc))
	}
	
	// Redirect route (not under /api to keep URLs short)
	r.GET("/:code", redirectURL(svc))
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
}

// CreateURLRequest represents the request payload for URL shortening
type CreateURLRequest struct {
	URL string `json:"url" binding:"required"`
}

// CreateURLResponse represents the response for URL shortening
type CreateURLResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
	FullURL   string `json:"full_url"`
}

// createShortURL handles URL shortening requests
func createShortURL(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateURLRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		shortCode, err := svc.CreateShortURL(req.URL)
		if err != nil {
			if err == service.ErrInvalidURL {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid URL format",
					"details": err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to create short URL",
					"details": err.Error(),
				})
			}
			return
		}

		// Build full URL
		host := c.Request.Host
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		fullURL := scheme + "://" + host + "/" + shortCode

		c.JSON(http.StatusCreated, CreateURLResponse{
			ShortURL:  shortCode,
			ShortCode: shortCode,
			FullURL:   fullURL,
		})
	}
}

// getURLStats retrieves statistics for a short URL
func getURLStats(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		stats, err := svc.GetURLStats(code)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

// redirectURL handles URL redirection
func redirectURL(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		
		// Validate code format
		if len(code) == 0 || len(code) > 10 {
			if gin.Mode() == gin.TestMode {
				c.JSON(http.StatusNotFound, gin.H{"error": "Invalid short code"})
			} else {
				c.HTML(http.StatusNotFound, "error.html", gin.H{
					"error": "Invalid short code",
				})
			}
			return
		}

		// Get IP address - pass as-is to service layer for proper handling
		ip := c.ClientIP()
		
		originalURL, err := svc.RedirectURL(code, ip, c.Request.UserAgent())
		if err != nil {
			if gin.Mode() == gin.TestMode {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			} else {
				c.HTML(http.StatusNotFound, "error.html", gin.H{
					"error": "Short URL not found",
				})
			}
			return
		}

		c.Redirect(http.StatusMovedPermanently, originalURL)
	}
}
