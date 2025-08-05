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
	"github.com/rusik69/shortener/internal/service"
)

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Use in-memory SQLite for testing
	testDB, err := sql.Open("pgx", "postgresql://postgres:postgres@localhost:5432/url_shortener_test?sslmode=disable")
	if err != nil {
		t.Skipf("Skipping integration tests: %v", err)
	}

	// Run migrations
	if err := db.MigrateDatabase(testDB); err != nil {
		t.Skipf("Skipping integration tests: failed to migrate: %v", err)
	}

	cleanup := func() {
		// Clean up test data
		testDB.Exec("DELETE FROM clicks")
		testDB.Exec("DELETE FROM short_urls")
		testDB.Exec("DELETE FROM users")
		testDB.Close()
	}

	return testDB, cleanup
}

func setupTestRouter(testDB *sql.DB) *gin.Engine {
	repo := db.NewRepository(testDB)
	svc := service.NewService(repo)
	
	router := gin.New()
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
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

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
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
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
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	shortCode := response["short_code"].(string)

	// Now test the redirect
	req, _ = http.NewRequest("GET", "/"+shortCode, nil)
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
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	shortCode := response["short_code"].(string)

	// Test getting stats
	req, _ = http.NewRequest("GET", "/api/stats/"+shortCode, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var statsResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &statsResponse)

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

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

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
	jsonBody, _ := json.Marshal(reqBody)

	// Make requests up to the rate limit
	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "127.0.0.1:12345" // Set consistent IP
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// The first requests should succeed, later ones should be rate limited
		if i < 99 {
			if w.Code != http.StatusCreated {
				t.Logf("Request %d: Expected status %d, got %d", i, http.StatusCreated, w.Code)
			}
		} else {
			// The 100th request should be rate limited
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Expected rate limiting on request %d, got status %d", i, w.Code)
			}
		}
	}
}

func TestIntegrationNotFound(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(testDB)

	// Test accessing non-existent short code
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	// Test getting stats for non-existent code
	req, _ = http.NewRequest("GET", "/api/stats/nonexistent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
