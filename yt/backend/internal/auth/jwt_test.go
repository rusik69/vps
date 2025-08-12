package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTManager(t *testing.T) {
	secretKey := "test-secret"
	issuer := "test-issuer"
	duration := 1 * time.Hour

	manager := NewJWTManager(secretKey, issuer, duration)

	assert.Equal(t, []byte(secretKey), manager.secretKey)
	assert.Equal(t, issuer, manager.issuer)
	assert.Equal(t, duration, manager.duration)
}

func TestGenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-issuer", 1*time.Hour)
	userID := 123

	token, err := manager.GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-issuer", 1*time.Hour)
	userID := 123

	// Generate token
	token, err := manager.GenerateToken(userID)
	require.NoError(t, err)

	// Validate token
	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "test-issuer", claims.Issuer)
}

func TestValidateTokenInvalid(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-issuer", 1*time.Hour)

	// Test invalid token
	claims, err := manager.ValidateToken("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateTokenExpired(t *testing.T) {
	// Create manager with very short duration
	manager := NewJWTManager("test-secret", "test-issuer", 1*time.Nanosecond)
	userID := 123

	// Generate token
	token, err := manager.GenerateToken(userID)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(1 * time.Millisecond)

	// Validate expired token
	claims, err := manager.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateTokenWrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret1", "test-issuer", 1*time.Hour)
	manager2 := NewJWTManager("secret2", "test-issuer", 1*time.Hour)
	userID := 123

	// Generate token with first manager
	token, err := manager1.GenerateToken(userID)
	require.NoError(t, err)

	// Try to validate with second manager (different secret)
	claims, err := manager2.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestAuthMiddleware(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-issuer", 1*time.Hour)
	userID := 123

	// Generate valid token
	token, err := manager.GenerateToken(userID)
	require.NoError(t, err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from context
		contextUserID, err := GetUserIDFromRequest(r)
		assert.NoError(t, err)
		assert.Equal(t, userID, contextUserID)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with auth middleware
	protectedHandler := manager.AuthMiddleware(testHandler)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "Missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization header required\n",
		},
		{
			name:           "Invalid auth header format",
			authHeader:     "InvalidFormat " + token,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format\n",
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token: token is malformed: token contains an invalid number of segments\n",
		},
		{
			name:           "Missing Bearer prefix",
			authHeader:     token,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			protectedHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	userID := 123

	// Test with valid context
	ctx := context.WithValue(context.Background(), UserIDKey, userID)
	retrievedUserID, err := GetUserIDFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, userID, retrievedUserID)

	// Test with empty context
	emptyCtx := context.Background()
	_, err = GetUserIDFromContext(emptyCtx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID not found in context")

	// Test with wrong type in context
	wrongTypeCtx := context.WithValue(context.Background(), UserIDKey, "not-an-int")
	_, err = GetUserIDFromContext(wrongTypeCtx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID not found in context")
}

func TestGetUserIDFromRequest(t *testing.T) {
	userID := 123

	// Test with valid request context
	ctx := context.WithValue(context.Background(), UserIDKey, userID)
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(ctx)

	retrievedUserID, err := GetUserIDFromRequest(req)
	assert.NoError(t, err)
	assert.Equal(t, userID, retrievedUserID)

	// Test with empty request context
	emptyReq := httptest.NewRequest("GET", "/test", nil)
	_, err = GetUserIDFromRequest(emptyReq)
	assert.Error(t, err)
}

func TestParseUserID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{
			name:        "Valid integer",
			input:       "123",
			expected:    123,
			expectError: false,
		},
		{
			name:        "Zero",
			input:       "0",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Negative integer",
			input:       "-1",
			expected:    -1,
			expectError: false,
		},
		{
			name:        "Invalid string",
			input:       "abc",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Float",
			input:       "123.45",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseUserID(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}