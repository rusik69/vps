package storage

import "github.com/rusik69/vps/yt/backend/internal/models"

// StorageInterface defines the contract for storage operations
type StorageInterface interface {
	// User operations
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)

	// Video operations
	CreateVideo(video *models.Video) error
	GetAllVideos() ([]models.Video, error)
	GetVideoByID(id int) (*models.Video, error)
	GetVideosByUserID(userID int) ([]models.Video, error)
	UpdateVideo(video *models.Video) error
	DeleteVideo(id, userID int) error
	IncrementVideoViews(id int) error

	// Cleanup
	Close() error
}