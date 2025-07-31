package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortURL(t *testing.T) {
	// Setup test service
	service := &service.MockService{}
	service.On("CreateShortURL", "https://example.com").Return("http://localhost:8080/test-code", nil)

	// Setup router
	r := gin.Default()
	SetupRoutes(r, service)

	// Create request
	req, err := http.NewRequest("POST", "/api/shorten", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create recorder
	rec := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{
		"short_url": "http://localhost:8080/test-code"
	}`, rec.Body.String())

	// Verify mock expectations
	service.AssertExpectations(t)
}

func TestGetURLStats(t *testing.T) {
	// Setup test service
	service := &service.MockService{}
	stats := service.URLStats{
		ClickCount: 10,
		CreatedAt:  time.Now(),
		Metadata:   "{}",
	}
	service.On("GetURLStats", "test-code").Return(stats, nil)

	// Setup router
	r := gin.Default()
	SetupRoutes(r, service)

	// Create request
	req, err := http.NewRequest("GET", "/api/stats/test-code", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create recorder
	rec := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{
		"click_count": 10,
		"created_at": "`+stats.CreatedAt.Format(time.RFC3339)+`",
		"metadata": "{}"
	}`, rec.Body.String())

	// Verify mock expectations
	service.AssertExpectations(t)
}

func TestRedirectURL(t *testing.T) {
	// Setup test service
	service := &service.MockService{}
	service.On("RedirectURL", "test-code", "127.0.0.1", "test-agent").Return("https://example.com", nil)

	// Setup router
	r := gin.Default()
	SetupRoutes(r, service)

	// Create request
	req, err := http.NewRequest("GET", "/test-code", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create recorder
	rec := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "https://example.com", rec.Header().Get("Location"))

	// Verify mock expectations
	service.AssertExpectations(t)
}
