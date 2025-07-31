package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)



// Service defines the interface for URL shortening operations
type Service interface {
	CreateShortURL(originalURL string) (string, error)
	GetURLStats(code string) (URLStats, error)
	RedirectURL(code, ip, userAgent string) (string, error)
}

type service struct {
	db *sql.DB
}

// NewService creates a new service instance
func NewService(db *sql.DB) Service {
	return &service{db: db}
}

func (s *service) CreateShortURL(originalURL string) (string, error) {
	shortCode := uuid.New().String()[:8]
	
	_, err := s.db.Exec(
		"INSERT INTO urls (short_code, original_url, created_at) VALUES ($1, $2, $3)",
		shortCode, originalURL, time.Now(),
	)
	if err != nil {
		return "", err
	}
	
	return shortCode, nil
}

func (s *service) GetURLStats(code string) (URLStats, error) {
	var stats URLStats
	row := s.db.QueryRow(
		"SELECT short_code, original_url, click_count, created_at FROM urls WHERE short_code = $1",
		code,
	)

	if err := row.Scan(&stats.Code, &stats.OriginalURL, &stats.Clicks, &stats.LastAccess); err != nil {
		return URLStats{}, err
	}

	return stats, nil
}

func (s *service) RedirectURL(code, ip, userAgent string) (string, error) {
	var originalURL string
	row := s.db.QueryRow(
		"SELECT original_url FROM urls WHERE short_code = $1",
		code,
	)

	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}

	// Record analytics
	_, err := s.db.Exec(
		"INSERT INTO analytics (url_id, ip_address, user_agent, accessed_at) VALUES ((SELECT id FROM urls WHERE short_code = $1), $2, $3, $4)",
		code,
		ip,
		userAgent,
		time.Now(),
	)
	if err != nil {
		// Don't fail the redirect if analytics fails
		fmt.Printf("Failed to record analytics: %v\n", err)
	}

	return originalURL, nil
}



// generateShortCode generates a unique short code
func generateShortCode() string {
	return uuid.New().String()[:8]
}
