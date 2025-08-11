package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Never serialize password
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Video represents a video in the system
type Video struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	ThumbnailURL string   `json:"thumbnail_url" db:"thumbnail_url"`
	UserID      int       `json:"user_id" db:"user_id"`
	Username    string    `json:"username" db:"username"` // For JOIN queries
	Views       int       `json:"views" db:"views"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateVideoRequest represents the request payload for creating a video
type CreateVideoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}