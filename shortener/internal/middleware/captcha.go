package middleware

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusik69/vps/url-shortener/internal/db"
)

const (
	// Captcha settings
	CaptchaWidth  = 160
	CaptchaHeight = 60
	CaptchaChars  = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	CaptchaLength = 6
)

// CaptchaMiddleware implements captcha protection middleware
type CaptchaMiddleware struct {
	db db.Repository
}

// NewCaptchaMiddleware creates a new captcha middleware
type CaptchaMiddlewareConfig struct {
	DB db.Repository
}

// NewCaptchaMiddleware creates a new captcha middleware instance
func NewCaptchaMiddleware(config CaptchaMiddlewareConfig) *CaptchaMiddleware {
	return &CaptchaMiddleware{
		db: config.DB,
	}
}

// CaptchaMiddleware returns a Gin middleware for captcha protection
func (cm *CaptchaMiddleware) CaptchaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if captcha is required based on recent attempts
		ip := c.ClientIP()
		attempts, err := cm.db.GetRecentCaptchaAttempts(ip, 5)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check captcha attempts",
			})
			return
		}

		// If there are recent failed attempts, require captcha
		if len(attempts) > 0 && attempts[0].Success == false {
			c.Set("captcha_required", true)
		}

		c.Next()
	}
}

// GenerateCaptcha generates a new captcha image and answer
func (cm *CaptchaMiddleware) GenerateCaptcha() (string, string, error) {
	// Generate random captcha text
	captchaText := generateRandomString(CaptchaChars, CaptchaLength)

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, CaptchaWidth, CaptchaHeight))
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	// Fill background
	for x := 0; x < CaptchaWidth; x++ {
		for y := 0; y < CaptchaHeight; y++ {
			img.Set(x, y, white)
		}
	}

	// Add noise
	for i := 0; i < 50; i++ {
		x := rand.Intn(CaptchaWidth)
		y := rand.Intn(CaptchaHeight)
		img.Set(x, y, black)
	}

	// Add lines
	for i := 0; i < 3; i++ {
		x1 := rand.Intn(CaptchaWidth)
		y1 := rand.Intn(CaptchaHeight)
		x2 := rand.Intn(CaptchaWidth)
		y2 := rand.Intn(CaptchaHeight)
		line(img, x1, y1, x2, y2, black)
	}

	// Add text
	font := loadFont() // You'll need to implement this or use an existing font
	// Draw text with some rotation and offset
	// This is a simplified version - you might want to use a proper text drawing library
	for i, char := range captchaText {
		x := (i * CaptchaWidth / len(captchaText)) + 10
		y := CaptchaHeight/2 + rand.Intn(10)
		// Draw character with some rotation
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", "", err
	}

	return captchaText, buf.String(), nil
}

// VerifyCaptcha verifies a captcha attempt
func (cm *CaptchaMiddleware) VerifyCaptcha(ip string, answer string, correct bool) error {
	return cm.db.CreateCaptchaAttempt(ip, correct)
}

// Helper functions
func generateRandomString(chars string, length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
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
