package service

import (
	"database/sql"
	"time"

	"github.com/rusik69/shortener/internal/db"
)

// AnalyticsService handles URL analytics
//go:generate mockery --name AnalyticsService --output mocks

type AnalyticsService interface {
	TrackURLAccess(code, ip, userAgent string) error
	GetURLStats(code string) (db.URLStats, error)
	GetAnalytics(code string, startDate, endDate time.Time) ([]db.URLAccess, error)
}

type analyticsService struct {
	db *sql.DB
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *sql.DB) *analyticsService {
	return &analyticsService{db: db}
}

// TrackURLAccess records a URL access event
func (s *analyticsService) TrackURLAccess(code, ip, userAgent string) error {
	_, err := s.db.Exec(
		"INSERT INTO shortener.analytics (url_id, ip_address, user_agent) VALUES ((SELECT id FROM shortener.urls WHERE short_code = $1), $2, $3)",
		code,
		ip,
		userAgent,
	)
	return err
}

// GetURLStats retrieves statistics for a URL
func (s *analyticsService) GetURLStats(code string) (db.URLStats, error) {
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

// GetAnalytics retrieves detailed access analytics for a URL
func (s *analyticsService) GetAnalytics(code string, startDate, endDate time.Time) ([]db.URLAccess, error) {
	rows, err := s.db.Query(
		"SELECT accessed_at, ip_address, user_agent, referrer, country_code FROM shortener.analytics WHERE url_id = (SELECT id FROM shortener.urls WHERE short_code = $1) AND accessed_at BETWEEN $2 AND $3 ORDER BY accessed_at DESC",
		code,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accesses []db.URLAccess
	for rows.Next() {
		var access db.URLAccess
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
