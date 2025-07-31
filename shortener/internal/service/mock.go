package service

import (
	"errors"
	"strings"
	"time"
)

// MockService is a mock implementation of Service for testing
var ErrInvalidURL = errors.New("invalid URL")
var ErrURLNotFound = errors.New("URL not found")

type MockService struct {
	urls map[string]string
	stats map[string]URLStats
}

func NewMockService() *MockService {
	return &MockService{
		urls: make(map[string]string),
		stats: make(map[string]URLStats),
	}
}

func (m *MockService) CreateShortURL(originalURL string) (string, error) {
	if originalURL == "" || !strings.HasPrefix(originalURL, "http") {
		return "", ErrInvalidURL
	}
	code := "abc123"
	m.urls[code] = originalURL
	m.stats[code] = URLStats{
		Code:        code,
		OriginalURL: originalURL,
		Clicks:      0,
		LastAccess:  time.Now(),
	}
	return "http://localhost:8080/" + code, nil
}

func (m *MockService) GetURLStats(code string) (URLStats, error) {
	if stats, exists := m.stats[code]; exists {
		return stats, nil
	}
	return URLStats{}, ErrURLNotFound
}

func (m *MockService) RedirectURL(code, ip, userAgent string) (string, error) {
	if originalURL, exists := m.urls[code]; exists {
		if stats, exists := m.stats[code]; exists {
			stats.Clicks++
			stats.LastAccess = time.Now()
			m.stats[code] = stats
		}
		return originalURL, nil
	}
	return "", ErrURLNotFound
}

func (m *MockService) ValidateCaptcha(captcha string) bool {
	return captcha == "test123"
}

// CheckRateLimit checks if the IP is rate limited
func (m *MockService) CheckRateLimit(ip string) (bool, error) {
	return false, nil
}

// IncrementRateLimit increments the rate limit counter for an IP
func (m *MockService) IncrementRateLimit(ip string) error {
	return nil
}
