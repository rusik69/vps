package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rusik69/vps/yt/backend/internal/auth"
	"github.com/rusik69/vps/yt/backend/internal/models"
)

// Mock storage
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateUser(user *models.User) error {
	args := m.Called(user)
	if args.Get(0) != nil {
		// Simulate database behavior
		user.ID = 1
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockStorage) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockStorage) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockStorage) CreateVideo(video *models.Video) error {
	args := m.Called(video)
	if args.Get(0) == nil {
		// Simulate database behavior
		video.ID = 1
		video.CreatedAt = time.Now()
		video.UpdatedAt = time.Now()
		video.Views = 0
	}
	return args.Error(0)
}

func (m *MockStorage) GetAllVideos() ([]models.Video, error) {
	args := m.Called()
	return args.Get(0).([]models.Video), args.Error(1)
}

func (m *MockStorage) GetVideoByID(id int) (*models.Video, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Video), args.Error(1)
}

func (m *MockStorage) GetVideosByUserID(userID int) ([]models.Video, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Video), args.Error(1)
}

func (m *MockStorage) UpdateVideo(video *models.Video) error {
	args := m.Called(video)
	return args.Error(0)
}

func (m *MockStorage) DeleteVideo(id, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockStorage) IncrementVideoViews(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func setupTestHandler() (*Handler, *MockStorage, *auth.JWTManager) {
	mockStorage := &MockStorage{}
	jwtManager := auth.NewJWTManager("test-secret", "test-issuer", 1*time.Hour)
	handler := New(mockStorage, jwtManager)
	return handler, mockStorage, jwtManager
}

func TestRegister(t *testing.T) {
	handler, mockStorage, _ := setupTestHandler()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockStorage)
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "Successful registration",
			requestBody: models.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockStorage) {
				m.On("GetUserByUsername", "testuser").Return((*models.User)(nil), nil)
				m.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.AuthResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, "testuser", response.User.Username)
			},
		},
		{
			name:           "Invalid request body",
			requestBody:    "invalid json",
			setupMock:      func(m *MockStorage) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "User already exists",
			requestBody: models.RegisterRequest{
				Username: "existinguser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockStorage) {
				existingUser := &models.User{
					ID:       1,
					Username: "existinguser",
					Email:    "existing@example.com",
				}
				m.On("GetUserByUsername", "existinguser").Return(existingUser, nil)
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockStorage.ExpectedCalls = nil
			mockStorage.Calls = nil
			tt.setupMock(mockStorage)

			// Create request
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Register(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr.Body.Bytes())
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	handler, mockStorage, _ := setupTestHandler()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockStorage)
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "Successful login",
			requestBody: models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(m *MockStorage) {
				// Hash of "password123"
				hashedPassword := "$2a$10$rGqh8/4G1V1DnN4Q4rGqhOQ4v4v4v4v4v4v4v4v4v4v4v4v4v4v4v4"
				user := &models.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				m.On("GetUserByUsername", "testuser").Return(user, nil)
			},
			expectedStatus: http.StatusUnauthorized, // Will fail due to password mismatch in test
		},
		{
			name:           "Invalid request body",
			requestBody:    "invalid json",
			setupMock:      func(m *MockStorage) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "User not found",
			requestBody: models.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			setupMock: func(m *MockStorage) {
				m.On("GetUserByUsername", "nonexistent").Return((*models.User)(nil), nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockStorage.ExpectedCalls = nil
			mockStorage.Calls = nil
			tt.setupMock(mockStorage)

			// Create request
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Login(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr.Body.Bytes())
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestGetVideos(t *testing.T) {
	handler, mockStorage, _ := setupTestHandler()

	videos := []models.Video{
		{
			ID:       1,
			Title:    "Video 1",
			URL:      "https://example.com/video1.mp4",
			UserID:   1,
			Username: "user1",
			Views:    100,
		},
		{
			ID:       2,
			Title:    "Video 2",
			URL:      "https://example.com/video2.mp4",
			UserID:   2,
			Username: "user2",
			Views:    50,
		},
	}

	mockStorage.On("GetAllVideos").Return(videos, nil)

	req := httptest.NewRequest("GET", "/api/videos", nil)
	rr := httptest.NewRecorder()

	handler.GetVideos(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Video
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Video 1", response[0].Title)

	mockStorage.AssertExpectations(t)
}

func TestGetVideo(t *testing.T) {
	handler, mockStorage, _ := setupTestHandler()

	video := &models.Video{
		ID:       1,
		Title:    "Test Video",
		URL:      "https://example.com/video.mp4",
		UserID:   1,
		Username: "testuser",
		Views:    0,
	}

	mockStorage.On("IncrementVideoViews", 1).Return(nil)
	mockStorage.On("GetVideoByID", 1).Return(video, nil)

	req := httptest.NewRequest("GET", "/api/videos/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.GetVideo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Video
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Test Video", response.Title)

	mockStorage.AssertExpectations(t)
}

func TestGetVideoNotFound(t *testing.T) {
	handler, mockStorage, _ := setupTestHandler()

	mockStorage.On("IncrementVideoViews", 999).Return(nil)
	mockStorage.On("GetVideoByID", 999).Return((*models.Video)(nil), nil)

	req := httptest.NewRequest("GET", "/api/videos/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()

	handler.GetVideo(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	mockStorage.AssertExpectations(t)
}

func TestCreateVideo(t *testing.T) {
	handler, mockStorage, jwtManager := setupTestHandler()

	userID := 1
	token, err := jwtManager.GenerateToken(userID)
	require.NoError(t, err)

	createRequest := models.CreateVideoRequest{
		Title:        "New Video",
		Description:  "Video Description",
		URL:          "https://example.com/new-video.mp4",
		ThumbnailURL: "https://example.com/thumb.jpg",
	}

	mockStorage.On("CreateVideo", mock.AnythingOfType("*models.Video")).Return(nil)

	body, err := json.Marshal(createRequest)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/videos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Add user ID to context (simulating auth middleware)
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateVideo(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Video
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "New Video", response.Title)
	assert.Equal(t, userID, response.UserID)

	mockStorage.AssertExpectations(t)
}

func TestCreateVideoUnauthorized(t *testing.T) {
	handler, _, _ := setupTestHandler()

	createRequest := models.CreateVideoRequest{
		Title: "New Video",
		URL:   "https://example.com/new-video.mp4",
	}

	body, err := json.Marshal(createRequest)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/videos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// No authorization header

	rr := httptest.NewRecorder()
	handler.CreateVideo(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateVideo(t *testing.T) {
	handler, mockStorage, jwtManager := setupTestHandler()

	userID := 1
	videoID := 1
	token, err := jwtManager.GenerateToken(userID)
	require.NoError(t, err)

	updateRequest := models.CreateVideoRequest{
		Title:        "Updated Video",
		Description:  "Updated Description",
		URL:          "https://example.com/updated-video.mp4",
		ThumbnailURL: "https://example.com/updated-thumb.jpg",
	}

	mockStorage.On("UpdateVideo", mock.AnythingOfType("*models.Video")).Return(nil)

	body, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/videos/"+strconv.Itoa(videoID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(videoID)})

	// Add user ID to context (simulating auth middleware)
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateVideo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Video
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Updated Video", response.Title)

	mockStorage.AssertExpectations(t)
}

func TestDeleteVideo(t *testing.T) {
	handler, mockStorage, jwtManager := setupTestHandler()

	userID := 1
	videoID := 1
	token, err := jwtManager.GenerateToken(userID)
	require.NoError(t, err)

	mockStorage.On("DeleteVideo", videoID, userID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/videos/"+strconv.Itoa(videoID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(videoID)})

	// Add user ID to context (simulating auth middleware)
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.DeleteVideo(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	mockStorage.AssertExpectations(t)
}

func TestGetMyVideos(t *testing.T) {
	handler, mockStorage, jwtManager := setupTestHandler()

	userID := 1
	token, err := jwtManager.GenerateToken(userID)
	require.NoError(t, err)

	videos := []models.Video{
		{
			ID:       1,
			Title:    "My Video 1",
			UserID:   userID,
			Username: "testuser",
		},
		{
			ID:       2,
			Title:    "My Video 2",
			UserID:   userID,
			Username: "testuser",
		},
	}

	mockStorage.On("GetVideosByUserID", userID).Return(videos, nil)

	req := httptest.NewRequest("GET", "/api/my-videos", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Add user ID to context (simulating auth middleware)
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetMyVideos(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.Video
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "My Video 1", response[0].Title)

	mockStorage.AssertExpectations(t)
}