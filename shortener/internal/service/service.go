package service

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/rusik69/shortener/internal/db"
)



// Service defines the interface for URL shortening operations
type Service interface {
	CreateShortURL(originalURL string) (string, error)
	GetURLStats(code string) (URLStats, error)
	RedirectURL(code, ip, userAgent string) (string, error)
}

type service struct {
	repo db.Repository
}

// NewService creates a new service instance
func NewService(repo db.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateShortURL(originalURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", ErrInvalidURL
	}
	shortCode := uuid.New().String()[:8]

	_, err = s.repo.CreateShortURL(shortCode, originalURL, nil, nil)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func (s *service) GetURLStats(code string) (URLStats, error) {
	shortURL, err := s.repo.GetShortURLByCode(code)
	if err != nil {
		return URLStats{}, err
	}

	stats := URLStats{
		Code:        shortURL.ShortCode,
		OriginalURL: shortURL.OriginalURL,
		Clicks:      int(shortURL.ClickCount),
		LastAccess:  shortURL.CreatedAt,
	}

	return stats, nil
}

func (s *service) RedirectURL(code, ip, userAgent string) (string, error) {
	shortURL, err := s.repo.GetShortURLByCode(code)
	if err != nil {
		return "", err
	}

	// Increment click count
	err = s.repo.IncrementClickCount(shortURL.ID)
	if err != nil {
		fmt.Printf("Failed to increment click count: %v\n", err)
	}

	// Record analytics - ensure we have a valid IP
	if ip == "" || ip == "::" || ip == "::1" {
		ip = "127.0.0.1" // Use localhost for invalid IPs
	}
	err = s.repo.CreateClick(shortURL.ID, userAgent, ip, "")
	if err != nil {
		// Don't fail the redirect if analytics fails
		fmt.Printf("Failed to record analytics: %v\n", err)
	}

	return shortURL.OriginalURL, nil
}
