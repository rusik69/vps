package service

import (
	"database/sql"
	"log"
	"time"
)

//go:generate mockery --name AnalyticsService --output mocks

// URLAccess represents a URL access record
type URLAccess struct {
	AccessedAt  time.Time `json:"accessed_at"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Referrer    string    `json:"referrer"`
	CountryCode string    `json:"country_code"`
}



// AnalyticsService handles URL analytics
type AnalyticsService interface {
	TrackURLAccess(code, ip, userAgent string) error
	GetURLStats(code string) (URLStats, error)
	GetAnalytics(code string, startDate, endDate time.Time) ([]URLAccess, error)
}

// analyticsService implements AnalyticsService
type analyticsService struct {
	db *sql.DB
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *sql.DB) AnalyticsService {
	return &analyticsService{db: db}
}

// TrackURLAccess records a URL access event
func (s *analyticsService) TrackURLAccess(code, ip, userAgent string) error {
	_, err := s.db.Exec(
		"INSERT INTO analytics (code, ip, user_agent, timestamp) VALUES (?, ?, ?, ?)",
		code, ip, userAgent, time.Now(),
	)
	return err
}

// GetURLStats retrieves statistics for a URL
func (s *analyticsService) GetURLStats(code string) (URLStats, error) {
	var stats URLStats
	row := s.db.QueryRow(
		"SELECT code, original_url, clicks, last_access FROM urls WHERE code = ?",
		code,
	)

	if err := row.Scan(&stats.Code, &stats.OriginalURL, &stats.Clicks, &stats.LastAccess); err != nil {
		return URLStats{}, err
	}

	return stats, nil
}

// GetAnalytics retrieves detailed access analytics for a URL
func (s *analyticsService) GetAnalytics(code string, startDate, endDate time.Time) ([]URLAccess, error) {
	rows, err := s.db.Query(
		"SELECT accessed_at, ip_address, user_agent, referrer, country_code FROM analytics WHERE code = ? AND accessed_at BETWEEN ? AND ? ORDER BY accessed_at DESC",
		code, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var accesses []URLAccess
	for rows.Next() {
		var access URLAccess
		if err := rows.Scan(
			&access.AccessedAt,
			&access.IPAddress,
			&access.UserAgent,
			&access.Referrer,
			&access.CountryCode,
		); err != nil {
			return nil, err
		}
		accesses = append(accesses, access)
	}

	return accesses, nil
}
