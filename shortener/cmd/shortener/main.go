package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rusik69/shortener/internal/api"
	"github.com/rusik69/shortener/internal/service"
)

// InitDatabase initializes the database connection
func InitDatabase() (*sql.DB, error) {
	dsn := "postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}

func main() {
	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Create service
	service := service.NewService(db)

	// Create Gin router
	router := gin.Default()

	// Setup routes
	api.SetupRoutes(router, service)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
