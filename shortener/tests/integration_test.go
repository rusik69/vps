//go:build integration

package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rusik69/shortener/internal/api"
	"github.com/rusik69/shortener/internal/db"
	"github.com/rusik69/shortener/internal/middleware"
	"github.com/rusik69/shortener/internal/service"
)

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Connect to PostgreSQL container for testing
	testDB, err := sql.Open("pgx", "postgresql://postgres:postgres@db:5432/url_shortener_test?sslmode=disable")
	if err != nil {
		t.Skipf("Skipping integration tests: %v", err)
	}

	// Run migrations
	if err := db.MigrateDatabase(testDB); err != nil {
		t.Skipf("Skipping integration tests: failed to migrate: %v", err)
	}

	cleanup := func() {
		// Clean up test data - ignore errors as these are cleanup operations
		_, _ = testDB.Exec("DELETE FROM clicks")
		_, _ = testDB.Exec("DELETE FROM short_urls")
		_, _ = testDB.Exec("DELETE FROM users")
		if err := testDB.Close(); err != nil {
			t.Logf("Error closing test database: %v", err)
		}
	}

	return testDB, cleanup
}

func setupTestRouter(testDB *sql.DB) *gin.Engine {
	repo := db.NewRepository(testDB)
	svc := service.NewService(repo)

	// Reset rate limiter before each test
	middleware.ResetRateLimiter()

	router := gin.Default()
	api.SetupRoutes(router, svc)

	return router
}

func TestIntegrationCreateShortURL(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	// Test valid URL shortening
	reqBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["short_code"] == nil {
		t.Error("Expected short_code in response")
	}

	if response["full_url"] == nil {
		t.Error("Expected full_url in response")
	}
}

func TestIntegrationInvalidURL(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	// Test invalid URL
	reqBody := map[string]string{
		"url": "not-a-valid-url",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestIntegrationRedirectFlow(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	// First, create a short URL
	reqBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	shortCode := response["short_code"].(string)

	// Now test the redirect
	req, err = http.NewRequest("GET", "/"+shortCode, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", http.StatusMovedPermanently, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("Expected redirect to https://example.com, got %s", location)
	}
}

func TestIntegrationGetStats(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	// First, create a short URL
	reqBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	shortCode := response["short_code"].(string)

	// Test getting stats
	req, err = http.NewRequest("GET", "/api/stats/"+shortCode, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var statsResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &statsResponse); err != nil {
		t.Fatalf("Failed to unmarshal stats response: %v", err)
	}

	if statsResponse["code"] != shortCode {
		t.Errorf("Expected code %s, got %v", shortCode, statsResponse["code"])
	}

	if statsResponse["original_url"] != "https://example.com" {
		t.Errorf("Expected original_url https://example.com, got %v", statsResponse["original_url"])
	}
}

func TestIntegrationHealthCheck(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal health response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status ok, got %v", response["status"])
	}
}

func TestIntegrationRateLimiting(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	// Make multiple requests quickly to trigger rate limiting
	reqBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Make requests up to the rate limit
	for i := 0; i < 101; i++ {
		req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "127.0.0.1") // Set consistent IP

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// The first 100 requests should succeed, the 101st should be rate limited
		if i < 100 {
			if w.Code != http.StatusCreated {
				t.Errorf("Request %d: Expected status %d, got %d", i+1, http.StatusCreated, w.Code)
			}
		} else {
			// The 101st request should be rate limited
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Expected rate limiting on request %d, got status %d", i+1, w.Code)
			}
		}
	}
}

func TestIntegrationNotFound(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	// Test accessing non-existent short code
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	// Test getting stats for non-existent code
	req, err = http.NewRequest("GET", "/api/stats/nonexistent", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
