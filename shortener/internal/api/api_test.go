package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create a mock service
	mockService := &MockService{}

	// Create Gin router
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create recorder
	rec := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ok")
}

func TestShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockService{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test valid URL shortening
	requestBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "short_code")
	assert.Contains(t, rec.Body.String(), "abc12345")
}

func TestShortenURLInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockServiceWithErrors{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test invalid URL
	requestBody := map[string]string{
		"url": "not-a-valid-url",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid URL format")
}

func TestShortenURLMissingURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockService{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test missing URL
	requestBody := map[string]string{}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request format")
}

func TestGetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockService{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test valid stats retrieval
	req, err := http.NewRequest("GET", "/api/stats/abc12345", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "abc12345")
	assert.Contains(t, rec.Body.String(), "https://example.com")
	assert.Contains(t, rec.Body.String(), "5")
}

func TestGetStatsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockServiceWithErrors{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test stats for non-existent URL
	req, err := http.NewRequest("GET", "/api/stats/notfound", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "URL not found")
}

func TestRedirectURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockService{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test valid redirect
	req, err := http.NewRequest("GET", "/abc12345", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
}

func TestRedirectURLNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockServiceWithErrors{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test redirect for non-existent URL
	req, err := http.NewRequest("GET", "/notfound", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestShortenURLServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockServiceWithErrors{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test service error during URL shortening
	requestBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "service error")
}

func TestInvalidJSONRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockService{}
	router := gin.Default()
	SetupRoutes(router, mockService)

	// Test invalid JSON
	req, err := http.NewRequest("POST", "/api/shorten", strings.NewReader("invalid json"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// MockService implements the Service interface for testing
type MockService struct{}

func (m *MockService) CreateShortURL(originalURL, customCode string) (string, error) {
	if customCode != "" {
		return customCode, nil
	}
	return "abc12345", nil
}

func (m *MockService) GetURLStats(code string) (service.URLStats, error) {
	return service.URLStats{
		Code:        code,
		OriginalURL: "https://example.com",
		Clicks:      5,
		LastAccess:  time.Now(),
	}, nil
}

func (m *MockService) RedirectURL(code, ip, userAgent string) (string, error) {
	return "https://example.com", nil
}

// MockServiceWithErrors implements the Service interface for error testing
type MockServiceWithErrors struct{}

func (m *MockServiceWithErrors) CreateShortURL(originalURL, customCode string) (string, error) {
	if originalURL == "https://example.com" {
		return "", errors.New("service error")
	}
	return "", service.ErrInvalidURL
}

func (m *MockServiceWithErrors) GetURLStats(code string) (service.URLStats, error) {
	return service.URLStats{}, errors.New("URL not found")
}

func (m *MockServiceWithErrors) RedirectURL(code, ip, userAgent string) (string, error) {
	return "", errors.New("URL not found")
}
