package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/shortener/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api.SetupRoutes(r, nil)

	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create recorder
	rec := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
}
