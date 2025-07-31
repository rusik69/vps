package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/rusik69/shortener/internal/db"
)

func TestDBOperations(t *testing.T) {
	// Setup test database
	db, err := sql.Open("postgres", "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	repo := db.NewRepository(db)

	// Test ShortURL operations
	t.Run("ShortURL Operations", func(t *testing.T) {
		shortURL := &db.ShortURL{
			ShortCode:   "test123",
			OriginalURL: "https://example.com",
			UserID:      nil,
			ExpiresAt:   nil,
		}

		// Create
		created, err := repo.CreateShortURL(shortURL.ShortCode, shortURL.OriginalURL, nil, nil)
		if err != nil {
			t.Errorf("Failed to create short URL: %v", err)
		}

		// Get
		retrieved, err := repo.GetShortURLByCode(shortURL.ShortCode)
		if err != nil {
			t.Errorf("Failed to get short URL: %v", err)
		}
		if retrieved.ID != created.ID {
			t.Errorf("Retrieved URL ID does not match created URL ID")
		}

		// Increment click count
		err = repo.IncrementClickCount(created.ID)
		if err != nil {
			t.Errorf("Failed to increment click count: %v", err)
		}

		// Get clicks
		clicks, err := repo.GetClicks(created.ID, 10)
		if err != nil {
			t.Errorf("Failed to get clicks: %v", err)
		}
		if len(clicks) != 1 {
			t.Errorf("Expected 1 click, got %d", len(clicks))
		}
	})

	// Test RateLimit operations
	t.Run("RateLimit Operations", func(t *testing.T) {
		rip := "127.0.0.1"

		// Get or create
		rateLimit, err := repo.GetOrCreateRateLimit(ip)
		if err != nil {
			t.Errorf("Failed to get or create rate limit: %v", err)
		}

		// Update
		rateLimit.RequestCount = 5
		err = repo.UpdateRateLimit(rateLimit)
		if err != nil {
			t.Errorf("Failed to update rate limit: %v", err)
		}

		// Verify update
		updated, err := repo.GetOrCreateRateLimit(ip)
		if err != nil {
			t.Errorf("Failed to verify rate limit update: %v", err)
		}
		if updated.RequestCount != 5 {
			t.Errorf("Rate limit count not updated correctly")
		}
	})

	// Test Captcha operations
	t.Run("Captcha Operations", func(t *testing.T) {
		rip := "127.0.0.1"

		// Create attempt
		err := repo.CreateCaptchaAttempt(ip, true)
		if err != nil {
			t.Errorf("Failed to create captcha attempt: %v", err)
		}

		// Get recent attempts
		attempts, err := repo.GetRecentCaptchaAttempts(ip, 5)
		if err != nil {
			t.Errorf("Failed to get captcha attempts: %v", err)
		}
		if len(attempts) != 1 {
			t.Errorf("Expected 1 captcha attempt, got %d", len(attempts))
		}
	})
}
