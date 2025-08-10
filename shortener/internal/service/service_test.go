package service

import (
	"testing"
	"time"

	"github.com/rusik69/shortener/internal/db"
)

func TestCreateShortURL(t *testing.T) {
	mockRepo := &MockRepository{}
	svc := NewService(mockRepo)
	testURL := "https://example.com"

	shortCode, err := svc.CreateShortURL(testURL, "")
	if err != nil {
		t.Errorf("CreateShortURL failed: %v", err)
	}

	if len(shortCode) != 8 {
		t.Errorf("Expected short code length 8, got %d", len(shortCode))
	}
}

func TestGetURLStats(t *testing.T) {
	mockRepo := &MockRepository{}
	svc := NewService(mockRepo)
	testCode := "abc12345"

	stats, err := svc.GetURLStats(testCode)
	if err != nil {
		t.Errorf("GetURLStats failed: %v", err)
	}

	if stats.Code != testCode {
		t.Errorf("Expected code %s, got %s", testCode, stats.Code)
	}
}

func TestRedirectURL(t *testing.T) {
	mockRepo := &MockRepository{}
	svc := NewService(mockRepo)
	testCode := "abc12345"
	testURL := "https://example.com"
	testIP := "192.168.1.1"
	testUserAgent := "Mozilla/5.0"

	originalURL, err := svc.RedirectURL(testCode, testIP, testUserAgent)
	if err != nil {
		t.Errorf("RedirectURL failed: %v", err)
	}

	if originalURL != testURL {
		t.Errorf("Expected URL %s, got %s", testURL, originalURL)
	}
}

// MockRepository implements the Repository interface for testing
type MockRepository struct{}

func (m *MockRepository) CreateShortURL(shortCode, originalURL string, userID *int64, expiresAt *time.Time) (*db.ShortURL, error) {
	return &db.ShortURL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		ClickCount:  0,
	}, nil
}

func (m *MockRepository) GetShortURLByCode(shortCode string) (*db.ShortURL, error) {
	return &db.ShortURL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: "https://example.com",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ClickCount:  5,
	}, nil
}

func (m *MockRepository) IncrementClickCount(shortURLID int64) error {
	return nil
}

func (m *MockRepository) GetClicks(shortURLID int64, limit int) ([]db.Click, error) {
	return []db.Click{}, nil
}

func (m *MockRepository) CreateClick(shortURLID int64, userAgent, ipAddress, referrer string) error {
	return nil
}

func (m *MockRepository) GetOrCreateRateLimit(ipAddress string) (*db.RateLimit, error) {
	return &db.RateLimit{}, nil
}

func (m *MockRepository) UpdateRateLimit(rateLimit *db.RateLimit) error {
	return nil
}

