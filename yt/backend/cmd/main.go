package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rusik69/vps/yt/backend/internal/auth"
	"github.com/rusik69/vps/yt/backend/internal/handlers"
	"github.com/rusik69/vps/yt/backend/internal/storage"
)

func main() {
	// Environment variables with defaults
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/youtube_clone?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	port := getEnv("PORT", "8080")

	// Initialize storage
	store, err := storage.New(databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("Failed to close storage: %v", err)
		}
	}()

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(jwtSecret, "youtube-clone", 24*time.Hour)

	// Initialize handlers
	handler := handlers.New(store, jwtManager)

	// Setup routes
	router := mux.NewRouter()

	// Apply CORS to all routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	// Public routes
	router.HandleFunc("/api/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/api/videos", handler.GetVideos).Methods("GET")
	router.HandleFunc("/api/videos/{id:[0-9]+}", handler.GetVideo).Methods("GET")

	// Protected routes
	router.HandleFunc("/api/videos", jwtManager.AuthMiddleware(handler.CreateVideo)).Methods("POST")
	router.HandleFunc("/api/videos/{id:[0-9]+}", jwtManager.AuthMiddleware(handler.UpdateVideo)).Methods("PUT")
	router.HandleFunc("/api/videos/{id:[0-9]+}", jwtManager.AuthMiddleware(handler.DeleteVideo)).Methods("DELETE")
	router.HandleFunc("/api/my-videos", jwtManager.AuthMiddleware(handler.GetMyVideos)).Methods("GET")
	router.HandleFunc("/api/upload", jwtManager.AuthMiddleware(handler.UploadVideo)).Methods("POST")

	// Video file serving
	router.PathPrefix("/videos/").HandlerFunc(handler.ServeVideo).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().UTC().Format(time.RFC3339)); err != nil {
			log.Printf("Failed to write health response: %v", err)
		}
	}).Methods("GET")

	fmt.Printf("Starting YouTube Clone API server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}