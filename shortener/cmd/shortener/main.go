package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/shortener/internal/api"
	"github.com/rusik69/shortener/internal/middleware"
	"github.com/rusik69/shortener/internal/service"
)

func main() {
	// Initialize dependencies
	db, err := service.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	service := service.NewService(db)

	// Create router
	router := gin.Default()

	// Middleware
	router.Use(middleware.RateLimiter())
	router.Use(middleware.CaptchaMiddleware())

	// API Routes
	api.SetupRoutes(router, service)

	// Web Routes
	setupWebRoutes(router)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupWebRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/stats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "stats.html", nil)
	})
}
