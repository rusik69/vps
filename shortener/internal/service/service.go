package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rusik69/shortener/internal/db"
)

// Service represents the main service interface
//go:generate mockery --name Service --output mocks

type Service interface {
	CreateShortURL(originalURL string) (string, error)
	GetURLStats(code string) (db.URLStats, error)
	RedirectURL(code, ip, userAgent string) (string, error)
}

type service struct {
	db *sql.DB
}

// NewService creates a new service instance
func NewService(db *sql.DB) *service {
	return &service{db: db}
}

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

// CreateShortURL creates a new shortened URL
func (s *service) CreateShortURL(originalURL string) (string, error) {
	// Generate a unique short code
	shortCode := generateShortCode()

	// Insert into database
	_, err := s.db.Exec(
		"INSERT INTO shortener.urls (original_url, short_code) VALUES ($1, $2)",
		originalURL,
		shortCode,
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:8080/%s", shortCode), nil
}

// GetURLStats retrieves statistics for a URL
func (s *service) GetURLStats(code string) (db.URLStats, error) {
	var stats db.URLStats
	row := s.db.QueryRow(
		"SELECT click_count, created_at, metadata FROM shortener.urls WHERE short_code = $1",
		code,
	)

	if err := row.Scan(&stats.ClickCount, &stats.CreatedAt, &stats.Metadata); err != nil {
		return db.URLStats{}, err
	}

	return stats, nil
}

// RedirectURL handles URL redirection and tracks analytics
func (s *service) RedirectURL(code, ip, userAgent string) (string, error) {
	var originalURL string
	row := s.db.QueryRow(
		"SELECT original_url, click_count FROM shortener.urls WHERE short_code = $1",
		code,
	)

	if err := row.Scan(&originalURL, &stats.ClickCount); err != nil {
		return "", err
	}

	// Update click count
	_, err := s.db.Exec(
		"UPDATE shortener.urls SET click_count = click_count + 1 WHERE short_code = $1",
		code,
	)
	if err != nil {
		return "", err
	}

	// Record analytics
	_, err = s.db.Exec(
		"INSERT INTO shortener.analytics (url_id, ip_address, user_agent) VALUES ((SELECT id FROM shortener.urls WHERE short_code = $1), $2, $3)",
		code,
		ip,
		userAgent,
	)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

// generateShortCode generates a unique short code
func generateShortCode() string {
	return uuid.New().String()[:8] // Simplified for example
}
