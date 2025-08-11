package storage

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rusik69/vps/yt/backend/internal/models"
)

type Storage struct {
	db *sqlx.DB
}

func New(databaseURL string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

// User methods
func (s *Storage) CreateUser(user *models.User) error {
	query := `
		INSERT INTO yt_users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	
	return s.db.QueryRow(query, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, created_at, updated_at FROM yt_users WHERE username = $1`
	
	err := s.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

func (s *Storage) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, created_at, updated_at FROM yt_users WHERE id = $1`
	
	err := s.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

// Video methods
func (s *Storage) CreateVideo(video *models.Video) error {
	query := `
		INSERT INTO yt_videos (title, description, url, thumbnail_url, user_id, views, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 0, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	
	return s.db.QueryRow(query, video.Title, video.Description, video.URL, video.ThumbnailURL, video.UserID).
		Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt)
}

func (s *Storage) GetAllVideos() ([]models.Video, error) {
	var videos []models.Video
	query := `
		SELECT v.id, v.title, v.description, v.url, v.thumbnail_url, v.user_id, u.username, v.views, v.created_at, v.updated_at
		FROM yt_videos v
		JOIN yt_users u ON v.user_id = u.id
		ORDER BY v.created_at DESC`
	
	err := s.db.Select(&videos, query)
	return videos, err
}

func (s *Storage) GetVideoByID(id int) (*models.Video, error) {
	var video models.Video
	query := `
		SELECT v.id, v.title, v.description, v.url, v.thumbnail_url, v.user_id, u.username, v.views, v.created_at, v.updated_at
		FROM yt_videos v
		JOIN yt_users u ON v.user_id = u.id
		WHERE v.id = $1`
	
	err := s.db.Get(&video, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &video, nil
}

func (s *Storage) GetVideosByUserID(userID int) ([]models.Video, error) {
	var videos []models.Video
	query := `
		SELECT v.id, v.title, v.description, v.url, v.thumbnail_url, v.user_id, u.username, v.views, v.created_at, v.updated_at
		FROM yt_videos v
		JOIN yt_users u ON v.user_id = u.id
		WHERE v.user_id = $1
		ORDER BY v.created_at DESC`
	
	err := s.db.Select(&videos, query, userID)
	return videos, err
}

func (s *Storage) UpdateVideo(video *models.Video) error {
	query := `
		UPDATE yt_videos 
		SET title = $1, description = $2, url = $3, thumbnail_url = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6`
	
	_, err := s.db.Exec(query, video.Title, video.Description, video.URL, video.ThumbnailURL, video.ID, video.UserID)
	return err
}

func (s *Storage) DeleteVideo(id, userID int) error {
	query := `DELETE FROM yt_videos WHERE id = $1 AND user_id = $2`
	_, err := s.db.Exec(query, id, userID)
	return err
}

func (s *Storage) IncrementVideoViews(id int) error {
	query := `UPDATE yt_videos SET views = views + 1 WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}