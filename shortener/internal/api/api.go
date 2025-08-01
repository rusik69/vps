package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/shortener/internal/service"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, svc service.Service) {
	api := r.Group("/api")
	{
		api.POST("/shorten", createShortURL(svc))
		api.GET("/stats/:code", getURLStats(svc))
		api.GET("/:code", redirectURL(svc))
	}
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// createShortURL handles URL shortening requests
func createShortURL(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			URL string `json:"url"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shortURL, err := svc.CreateShortURL(req.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_url": shortURL,
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
		originalURL, err := svc.RedirectURL(code, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, originalURL)
	}
}
