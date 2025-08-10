package service

import (
	"errors"
	"strings"
	"time"
)

// MockService is a mock implementation of Service for testing
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

func (m *MockService) CreateShortURL(originalURL, customCode string) (string, error) {
	if originalURL == "" || !strings.HasPrefix(originalURL, "http") {
		return "", ErrInvalidURL
	}
	
	var code string
	if customCode != "" {
		if _, exists := m.urls[customCode]; exists {
			return "", errors.New("custom code already exists")
		}
		code = customCode
	} else {
		code = "abc123"
	}
	
	m.urls[code] = originalURL
	m.stats[code] = URLStats{
		Code:        code,
		OriginalURL: originalURL,
		Clicks:      0,
		LastAccess:  time.Now(),
	}
	return code, nil
}

func (m *MockService) GetURLStats(code string) (URLStats, error) {
	if stats, exists := m.stats[code]; exists {
		return stats, nil
	}
	return URLStats{}, errors.New("URL not found")
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
	return "", errors.New("URL not found")
}


// CheckRateLimit checks if the IP is rate limited
func (m *MockService) CheckRateLimit(ip string) (bool, error) {
	return false, nil
}

// IncrementRateLimit increments the rate limit counter for an IP
func (m *MockService) IncrementRateLimit(ip string) error {
	return nil
}
