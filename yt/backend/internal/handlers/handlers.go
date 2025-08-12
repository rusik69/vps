package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/rusik69/vps/yt/backend/internal/auth"
	"github.com/rusik69/vps/yt/backend/internal/models"
	"github.com/rusik69/vps/yt/backend/internal/storage"
)

type Handler struct {
	storage    storage.StorageInterface
	jwtManager *auth.JWTManager
}

func New(storage storage.StorageInterface, jwtManager *auth.JWTManager) *Handler {
	return &Handler{
		storage:    storage,
		jwtManager: jwtManager,
	}
}

// CORS middleware
func (h *Handler) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}

// Auth handlers
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, err := h.storage.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.storage.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user by username
	user, err := h.storage.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Video handlers
func (h *Handler) GetVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.storage.GetAllVideos()
	if err != nil {
		http.Error(w, "Failed to fetch videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// Increment views
	h.storage.IncrementVideoViews(id)

	video, err := h.storage.GetVideoByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
		return
	}
	if video == nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

func (h *Handler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	video := &models.Video{
		Title:        req.Title,
		Description:  req.Description,
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		UserID:       userID,
	}

	if err := h.storage.CreateVideo(video); err != nil {
		http.Error(w, "Failed to create video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

func (h *Handler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	var req models.CreateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	video := &models.Video{
		ID:           id,
		Title:        req.Title,
		Description:  req.Description,
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		UserID:       userID,
	}

	if err := h.storage.UpdateVideo(video); err != nil {
		http.Error(w, "Failed to update video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

func (h *Handler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteVideo(id, userID); err != nil {
		http.Error(w, "Failed to delete video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetMyVideos(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	videos, err := h.storage.GetVideosByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}