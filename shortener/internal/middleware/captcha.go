package middleware

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rusik69/shortener/internal/db"
)

var (
	captchaStore = make(map[string]string)
	captchaMutex sync.RWMutex
)

// CaptchaMiddleware handles captcha validation
func CaptchaMiddleware(repo db.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For now, captcha middleware is simplified - captcha validation is handled at service layer
			next.ServeHTTP(w, r)
		})
	}
}

// GenerateCaptcha generates a new captcha
func GenerateCaptcha(repo db.Repository) (string, []byte, error) {
	captchaID := strconv.FormatInt(time.Now().UnixNano(), 10)
	captchaText := generateRandomString(6)

	// Store in database - for now, we'll just create the attempt with success=true
	// In a real implementation, this would be handled by the service layer
	err := repo.CreateCaptchaAttempt("127.0.0.1", true)
	if err != nil {
		return "", nil, err
	}

	// Generate image
	img := generateCaptchaImage(captchaText)
	
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", nil, err
	}

	return captchaID, buf.Bytes(), nil
}

// Helper functions
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func generateCaptchaImage(text string) *image.RGBA {
	width, height := 200, 80
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{240, 240, 240, 255})
		}
	}
	
	// Add text
	chars := []byte(text)
	for i := range chars {
		color := color.RGBA{
			R: uint8(rand.Intn(100) + 100),
			G: uint8(rand.Intn(100) + 100),
			B: uint8(rand.Intn(100) + 100),
			A: 255,
		}
		// Simple text rendering - in real app, use proper font rendering
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if (i*20+x) < width && y < height {
					img.Set(i*20+x, height/2+y, color)
				}
			}
		}
	}
	
	return img
}

func line(img *image.RGBA, x1, y1, x2, y2 int, color color.Color) {
	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 {
		for y := y1; y <= y2; y++ {
			img.Set(x1, y, color)
		}
		return
	}

	if dy == 0 {
		for x := x1; x <= x2; x++ {
			img.Set(x, y1, color)
		}
		return
	}

	dx1 := abs(dx)
	dy1 := abs(dy)

	if dx1 >= dy1 {
		xstep := 1
		if dx < 0 {
			xstep = -1
		}
		y := y1
		d := dy1 - (dx1 / 2)
		for x := x1; x != x2; x += xstep {
			img.Set(x, y, color)
			if d >= 0 {
				y += 1
				d -= dx1
			}
			d += dy1
		}
	} else {
		ystep := 1
		if dy < 0 {
			ystep = -1
		}
		x := x1
		d := dx1 - (dy1 / 2)
		for y := y1; y != y2; y += ystep {
			img.Set(x, y, color)
			if d >= 0 {
				x += 1
				d -= dy1
			}
			d += dx1
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
