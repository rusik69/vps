package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rusik69/vps/yt/backend/internal/auth"
)

const (
	maxUploadSize = 500 << 20 // 500MB
	uploadDir     = "/app/uploads"
)

// UploadVideo handles video file uploads
func (h *Handler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, fileHeader, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close uploaded file: %v", err)
		}
	}()

	// Validate file type
	if !isValidVideoFile(fileHeader.Filename) {
		http.Error(w, "Invalid file type. Only video files are allowed.", http.StatusBadRequest)
		return
	}

	// Create upload directory if it doesn't exist
	uploadPath := getUploadDir()
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	filename := generateUniqueFilename(userID, fileHeader.Filename)
	filePath := filepath.Join(uploadPath, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("Failed to close destination file: %v", err)
		}
	}()

	// Copy the uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		// Clean up the file if copy fails
		if removeErr := os.Remove(filePath); removeErr != nil {
			log.Printf("Failed to remove file after copy failure: %v", removeErr)
		}
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Return the file URL
	fileURL := fmt.Sprintf("/videos/%s", filename)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprintf(w, `{"url": "%s", "filename": "%s"}`, fileURL, filename); err != nil {
		log.Printf("Failed to write upload response: %v", err)
	}
}

// ServeVideo serves uploaded video files
func (h *Handler) ServeVideo(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	path := strings.TrimPrefix(r.URL.Path, "/videos/")
	if path == "" || strings.Contains(path, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Get full file path
	filePath := filepath.Join(getUploadDir(), path)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", getContentType(path))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// Helper functions

func isValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".mp4", ".webm", ".ogg", ".avi", ".mov", ".wmv", ".flv", ".mkv"}
	
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func generateUniqueFilename(userID int, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	base := strings.TrimSuffix(filepath.Base(originalFilename), ext)
	
	// Clean filename of potentially dangerous characters
	base = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, base)
	
	// Generate unique filename with user ID and timestamp
	timestamp := fmt.Sprintf("%d", getUserTimestamp())
	return fmt.Sprintf("%d_%s_%s%s", userID, timestamp, base, ext)
}

func getUploadDir() string {
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		return dir
	}
	return uploadDir
}

func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogg":
		return "video/ogg"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".flv":
		return "video/x-flv"
	case ".mkv":
		return "video/x-matroska"
	default:
		return "application/octet-stream"
	}
}

func getUserTimestamp() int64 {
	return time.Now().Unix()
}