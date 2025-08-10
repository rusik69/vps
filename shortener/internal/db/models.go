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
	ID           int64     `json:"id"`
	IPAddress    string    `json:"ip_address"`
	RequestCount int64     `json:"request_count"`
	ResetAt      time.Time `json:"reset_at"`
	CreatedAt    time.Time `json:"created_at"`
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
}

// NewRepository creates a new database repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

type repository struct {
	db *sql.DB
}

func (r *repository) CreateShortURL(shortCode, originalURL string, userID *int64, expiresAt *time.Time) (*ShortURL, error) {
	now := time.Now()
	var id int64
	err := r.db.QueryRow(
		"INSERT INTO short_urls (short_code, original_url, user_id, expires_at, created_at, updated_at, click_count) VALUES ($1, $2, $3, $4, $5, $6, 0) RETURNING id",
		shortCode, originalURL, userID, expiresAt, now, now,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &ShortURL{
		ID:          id,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		ExpiresAt:   expiresAt,
		ClickCount:  0,
	}, nil
}

func (r *repository) GetShortURLByCode(shortCode string) (*ShortURL, error) {
	var url ShortURL
	err := r.db.QueryRow(
		"SELECT id, short_code, original_url, user_id, created_at, updated_at, expires_at, click_count FROM short_urls WHERE short_code = $1",
		shortCode,
	).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.CreatedAt, &url.UpdatedAt, &url.ExpiresAt, &url.ClickCount)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *repository) IncrementClickCount(shortURLID int64) error {
	_, err := r.db.Exec(
		"UPDATE short_urls SET click_count = click_count + 1, updated_at = $1 WHERE id = $2",
		time.Now(), shortURLID,
	)
	return err
}

func (r *repository) GetClicks(shortURLID int64, limit int) ([]Click, error) {
	rows, err := r.db.Query(
		"SELECT id, short_url_id, user_agent, ip_address, referrer, created_at FROM clicks WHERE short_url_id = $1 ORDER BY created_at DESC LIMIT $2",
		shortURLID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var clicks []Click
	for rows.Next() {
		var click Click
		err := rows.Scan(&click.ID, &click.ShortURLID, &click.UserAgent, &click.IPAddress, &click.Referrer, &click.CreatedAt)
		if err != nil {
			return nil, err
		}
		clicks = append(clicks, click)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clicks, nil
}

func (r *repository) CreateClick(shortURLID int64, userAgent, ipAddress, referrer string) error {
	_, err := r.db.Exec(
		"INSERT INTO clicks (short_url_id, user_agent, ip_address, referrer, created_at) VALUES ($1, $2, $3, $4, $5)",
		shortURLID, userAgent, ipAddress, referrer, time.Now(),
	)
	return err
}

func (r *repository) GetOrCreateRateLimit(ipAddress string) (*RateLimit, error) {
	// Handle invalid IP addresses
	if ipAddress == "" || ipAddress == "::" || ipAddress == "::1" {
		ipAddress = "127.0.0.1" // Default to localhost for rate limiting
	}
	
	var rateLimit RateLimit
	err := r.db.QueryRow(
		"SELECT id, ip_address, request_count, reset_at, created_at FROM rate_limits WHERE ip_address = $1",
		ipAddress,
	).Scan(&rateLimit.ID, &rateLimit.IPAddress, &rateLimit.RequestCount, &rateLimit.ResetAt, &rateLimit.CreatedAt)

	if err == sql.ErrNoRows {
		// Create new rate limit entry
		now := time.Now()
		resetAt := now.Add(time.Hour)
		err = r.db.QueryRow(
			"INSERT INTO rate_limits (ip_address, request_count, reset_at, created_at) VALUES ($1, 1, $2, $3) RETURNING id",
			ipAddress, resetAt, now,
		).Scan(&rateLimit.ID)
		if err != nil {
			return nil, err
		}
		rateLimit.IPAddress = ipAddress
		rateLimit.RequestCount = 1
		rateLimit.ResetAt = resetAt
		rateLimit.CreatedAt = now
	} else if err != nil {
		return nil, err
	}

	return &rateLimit, nil
}

func (r *repository) UpdateRateLimit(rateLimit *RateLimit) error {
	_, err := r.db.Exec(
		"UPDATE rate_limits SET request_count = $1, reset_at = $2 WHERE id = $3",
		rateLimit.RequestCount, rateLimit.ResetAt, rateLimit.ID,
	)
	return err
}

