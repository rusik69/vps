package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rusik69/shortener/internal/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [migrate|seed|reset]")
		os.Exit(1)
	}

	command := os.Args[1]

	// Initialize database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"
	}

	dbConn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	switch command {
	case "migrate":
		if err := db.MigrateDatabase(dbConn); err != nil {
			log.Fatal("Migration failed:", err)
		}
		fmt.Println("✅ Database migration completed successfully")

	case "seed":
		if err := db.SeedDatabase(dbConn); err != nil {
			log.Fatal("Seeding failed:", err)
		}
		fmt.Println("✅ Database seeding completed successfully")

	case "reset":
		if err := resetDatabase(dbConn); err != nil {
			log.Fatal("Reset failed:", err)
		}
		fmt.Println("✅ Database reset completed successfully")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: migrate, seed, reset")
		os.Exit(1)
	}
}

func resetDatabase(dbConn *sql.DB) error {
	log.Println("Dropping all tables...")

	// Drop tables in reverse order due to foreign key constraints
	tables := []string{
		"captcha_attempts",
		"rate_limits",
		"clicks",
		"short_urls",
		"users",
	}

	for _, table := range tables {
		_, err := dbConn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %v", table, err)
		}
	}

	log.Println("Running migrations...")
	return db.MigrateDatabase(dbConn)
}
