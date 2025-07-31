package service

import "time"

// URLStats represents URL statistics
type URLStats struct {
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int       `json:"clicks"`
	LastAccess  time.Time `json:"last_access"`
}
