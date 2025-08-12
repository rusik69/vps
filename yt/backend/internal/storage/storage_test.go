package storage

import (
	"os"
	"testing"
	"time"

	"github.com/rusik69/vps/yt/backend/internal/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *Storage {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/youtube_clone_test?sslmode=disable"
	}

	storage, err := New(databaseURL)
	require.NoError(t, err)

	// Clean up tables before each test
	cleanupTables(t, storage)

	return storage
}

func cleanupTables(t *testing.T, storage *Storage) {
	// Clean up in reverse order due to foreign key constraints
	_, err := storage.db.Exec("DELETE FROM yt_videos")
	require.NoError(t, err)
	_, err = storage.db.Exec("DELETE FROM yt_users")
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := storage.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestCreateUserDuplicate(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	user1 := &models.User{
		Username: "testuser",
		Email:    "test1@example.com",
		Password: "hashedpassword",
	}

	user2 := &models.User{
		Username: "testuser", // Same username
		Email:    "test2@example.com",
		Password: "hashedpassword",
	}

	err := storage.CreateUser(user1)
	assert.NoError(t, err)

	err = storage.CreateUser(user2)
	assert.Error(t, err) // Should fail due to duplicate username
}

func TestGetUserByUsername(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Test user not found
	user, err := storage.GetUserByUsername("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, user)

	// Create and retrieve user
	originalUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err = storage.CreateUser(originalUser)
	require.NoError(t, err)

	retrievedUser, err := storage.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, originalUser.Username, retrievedUser.Username)
	assert.Equal(t, originalUser.Email, retrievedUser.Email)
	assert.Equal(t, originalUser.Password, retrievedUser.Password)
}

func TestGetUserByID(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Test user not found
	user, err := storage.GetUserByID(999)
	assert.NoError(t, err)
	assert.Nil(t, user)

	// Create and retrieve user
	originalUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err = storage.CreateUser(originalUser)
	require.NoError(t, err)

	retrievedUser, err := storage.GetUserByID(originalUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, originalUser.ID, retrievedUser.ID)
	assert.Equal(t, originalUser.Username, retrievedUser.Username)
}

func TestCreateVideo(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Create user first
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := storage.CreateUser(user)
	require.NoError(t, err)

	video := &models.Video{
		Title:        "Test Video",
		Description:  "Test Description",
		URL:          "https://example.com/video.mp4",
		ThumbnailURL: "https://example.com/thumb.jpg",
		UserID:       user.ID,
	}

	err = storage.CreateVideo(video)
	assert.NoError(t, err)
	assert.NotZero(t, video.ID)
	assert.NotZero(t, video.CreatedAt)
	assert.NotZero(t, video.UpdatedAt)
	assert.Equal(t, 0, video.Views) // Should default to 0
}

func TestGetAllVideos(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Create user
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := storage.CreateUser(user)
	require.NoError(t, err)

	// Test empty result
	videos, err := storage.GetAllVideos()
	assert.NoError(t, err)
	assert.Empty(t, videos)

	// Create videos
	video1 := &models.Video{
		Title:       "Video 1",
		Description: "Description 1",
		URL:         "https://example.com/video1.mp4",
		UserID:      user.ID,
	}
	video2 := &models.Video{
		Title:       "Video 2",
		Description: "Description 2",
		URL:         "https://example.com/video2.mp4",
		UserID:      user.ID,
	}

	err = storage.CreateVideo(video1)
	require.NoError(t, err)
	
	// Add small delay to ensure different created_at times
	time.Sleep(1 * time.Millisecond)
	
	err = storage.CreateVideo(video2)
	require.NoError(t, err)

	// Retrieve all videos
	videos, err = storage.GetAllVideos()
	assert.NoError(t, err)
	assert.Len(t, videos, 2)
	
	// Should be ordered by created_at DESC (newest first)
	assert.Equal(t, "Video 2", videos[0].Title)
	assert.Equal(t, "Video 1", videos[1].Title)
	assert.Equal(t, user.Username, videos[0].Username)
}

func TestGetVideoByID(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Test video not found
	video, err := storage.GetVideoByID(999)
	assert.NoError(t, err)
	assert.Nil(t, video)

	// Create user and video
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err = storage.CreateUser(user)
	require.NoError(t, err)

	originalVideo := &models.Video{
		Title:        "Test Video",
		Description:  "Test Description",
		URL:          "https://example.com/video.mp4",
		ThumbnailURL: "https://example.com/thumb.jpg",
		UserID:       user.ID,
	}
	err = storage.CreateVideo(originalVideo)
	require.NoError(t, err)

	// Retrieve video
	retrievedVideo, err := storage.GetVideoByID(originalVideo.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedVideo)
	assert.Equal(t, originalVideo.Title, retrievedVideo.Title)
	assert.Equal(t, originalVideo.Description, retrievedVideo.Description)
	assert.Equal(t, originalVideo.URL, retrievedVideo.URL)
	assert.Equal(t, user.Username, retrievedVideo.Username)
}

func TestGetVideosByUserID(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Create users
	user1 := &models.User{
		Username: "user1",
		Email:    "user1@example.com",
		Password: "hashedpassword",
	}
	user2 := &models.User{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "hashedpassword",
	}
	err := storage.CreateUser(user1)
	require.NoError(t, err)
	err = storage.CreateUser(user2)
	require.NoError(t, err)

	// Create videos for different users
	video1 := &models.Video{
		Title:  "User1 Video 1",
		URL:    "https://example.com/video1.mp4",
		UserID: user1.ID,
	}
	video2 := &models.Video{
		Title:  "User1 Video 2",
		URL:    "https://example.com/video2.mp4",
		UserID: user1.ID,
	}
	video3 := &models.Video{
		Title:  "User2 Video 1",
		URL:    "https://example.com/video3.mp4",
		UserID: user2.ID,
	}

	err = storage.CreateVideo(video1)
	require.NoError(t, err)
	err = storage.CreateVideo(video2)
	require.NoError(t, err)
	err = storage.CreateVideo(video3)
	require.NoError(t, err)

	// Get videos for user1
	user1Videos, err := storage.GetVideosByUserID(user1.ID)
	assert.NoError(t, err)
	assert.Len(t, user1Videos, 2)

	// Get videos for user2
	user2Videos, err := storage.GetVideosByUserID(user2.ID)
	assert.NoError(t, err)
	assert.Len(t, user2Videos, 1)
	assert.Equal(t, "User2 Video 1", user2Videos[0].Title)
}

func TestUpdateVideo(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Create user and video
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := storage.CreateUser(user)
	require.NoError(t, err)

	video := &models.Video{
		Title:       "Original Title",
		Description: "Original Description",
		URL:         "https://example.com/original.mp4",
		UserID:      user.ID,
	}
	err = storage.CreateVideo(video)
	require.NoError(t, err)

	// Update video
	video.Title = "Updated Title"
	video.Description = "Updated Description"
	video.URL = "https://example.com/updated.mp4"
	video.ThumbnailURL = "https://example.com/updated-thumb.jpg"

	err = storage.UpdateVideo(video)
	assert.NoError(t, err)

	// Verify update
	updatedVideo, err := storage.GetVideoByID(video.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedVideo.Title)
	assert.Equal(t, "Updated Description", updatedVideo.Description)
	assert.Equal(t, "https://example.com/updated.mp4", updatedVideo.URL)
	assert.Equal(t, "https://example.com/updated-thumb.jpg", updatedVideo.ThumbnailURL)
}

func TestDeleteVideo(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Create user and video
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := storage.CreateUser(user)
	require.NoError(t, err)

	video := &models.Video{
		Title:  "Test Video",
		URL:    "https://example.com/video.mp4",
		UserID: user.ID,
	}
	err = storage.CreateVideo(video)
	require.NoError(t, err)

	// Delete video
	err = storage.DeleteVideo(video.ID, user.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedVideo, err := storage.GetVideoByID(video.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedVideo)
}

func TestDeleteVideoUnauthorized(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Create users
	user1 := &models.User{
		Username: "user1",
		Email:    "user1@example.com",
		Password: "hashedpassword",
	}
	user2 := &models.User{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "hashedpassword",
	}
	err := storage.CreateUser(user1)
	require.NoError(t, err)
	err = storage.CreateUser(user2)
	require.NoError(t, err)

	// Create video owned by user1
	video := &models.Video{
		Title:  "User1 Video",
		URL:    "https://example.com/video.mp4",
		UserID: user1.ID,
	}
	err = storage.CreateVideo(video)
	require.NoError(t, err)

	// Try to delete with user2 (should not delete)
	err = storage.DeleteVideo(video.ID, user2.ID)
	assert.NoError(t, err) // No error but no rows affected

	// Verify video still exists
	existingVideo, err := storage.GetVideoByID(video.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existingVideo)
}

func TestIncrementVideoViews(t *testing.T) {
	storage := setupTestDB(t)
	defer storage.Close()

	// Create user and video
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := storage.CreateUser(user)
	require.NoError(t, err)

	video := &models.Video{
		Title:  "Test Video",
		URL:    "https://example.com/video.mp4",
		UserID: user.ID,
	}
	err = storage.CreateVideo(video)
	require.NoError(t, err)

	// Verify initial view count
	retrievedVideo, err := storage.GetVideoByID(video.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, retrievedVideo.Views)

	// Increment views
	err = storage.IncrementVideoViews(video.ID)
	assert.NoError(t, err)

	err = storage.IncrementVideoViews(video.ID)
	assert.NoError(t, err)

	// Verify updated view count
	updatedVideo, err := storage.GetVideoByID(video.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, updatedVideo.Views)
}