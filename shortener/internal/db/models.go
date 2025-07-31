package db

import (
	"database/sql"
	"time"
)

// ShortURL represents a shortened URL in the database
type ShortURL struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	UserID      *int64   `json:"user_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ClickCount  int64     `json:"click_count"`
}

// Click represents a click on a shortened URL
type Click struct {
	ID          int64     `json:"id"`
	ShortURLID  int64     `json:"short_url_id"`
	UserAgent   string    `json:"user_agent"`
	IPAddress   string    `json:"ip_address"`
	Referrer    string    `json:"referrer"`
	CreatedAt   time.Time `json:"created_at"`
}

// RateLimit represents a rate limit entry
type RateLimit struct {
	ID        int64     `json:"id"`
	IPAddress string    `json:"ip_address"`
	RequestCount int64   `json:"request_count"`
	ResetAt   time.Time `json:"reset_at"`
	CreatedAt time.Time `json:"created_at"`
}

// CaptchaAttempt represents a captcha attempt
type CaptchaAttempt struct {
	ID        int64     `json:"id"`
	IPAddress string    `json:"ip_address"`
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository interface defines database operations
type Repository interface {
	// ShortURL operations
	CreateShortURL(shortCode, originalURL string, userID *int64, expiresAt *time.Time) (*ShortURL, error)
	GetShortURLByCode(shortCode string) (*ShortURL, error)
	IncrementClickCount(shortURLID int64) error
	GetClicks(shortURLID int64, limit int) ([]Click, error)

	// Click operations
	CreateClick(shortURLID int64, userAgent, ipAddress, referrer string) error

	// Rate limiting operations
	GetOrCreateRateLimit(ipAddress string) (*RateLimit, error)
	UpdateRateLimit(rateLimit *RateLimit) error

	// Captcha operations
	CreateCaptchaAttempt(ipAddress string, success bool) error
	GetRecentCaptchaAttempts(ipAddress string, limit int) ([]CaptchaAttempt, error)
}

// NewRepository creates a new database repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

type repository struct {
	db *sql.DB
}
