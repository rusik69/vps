package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/vps/url-shortener/internal/api"
	"github.com/rusik69/vps/url-shortener/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestAPIEndpoints(t *testing.T) {
	// Setup test database and repository
	db, err := sql.Open("postgres", "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	repo := db.NewRepository(db)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api.InitRoutes(r)

	// Test Shorten URL endpoint
	t.Run("Shorten URL", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(`{
			"url": "https://example.com",
			"captcha": "test123"
		}`))
		req.Header.Set("Content-Type", "application/json")
		
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response, "shortUrl")
	})

	// Test Get URL endpoint
	t.Run("Get URL", func(t *testing.T) {
		// First create a short URL
		shortCode := "test123"
		_, err := repo.CreateShortURL(shortCode, "https://example.com", nil, nil)
		if err != nil {
			t.Fatalf("Failed to create test URL: %v", err)
		}

		req, _ := http.NewRequest("GET", "/api/urls/test123", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
	})

	// Test Rate Limiting
	t.Run("Rate Limiting", func(t *testing.T) {
		// Send multiple requests to trigger rate limiting
		for i := 0; i < 105; i++ {
			req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(`{
				"url": "https://example.com",
				"captcha": "test123"
			}`))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Real-IP", "127.0.0.1")
			
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if i >= 100 {
				assert.Equal(t, http.StatusTooManyRequests, rec.Code)
			} else {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		}
	})

	// Test Captcha
	t.Run("Captcha", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/captcha", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response, "image")
	})
}
