package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortURL(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	service := NewService(db)

	// Mock database behavior
	mock.ExpectExec("INSERT INTO shortener.urls").
		WithArgs("https://example.com", "test-code").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test case
	shortURL, err := service.CreateShortURL("https://example.com")
	assert.NoError(t, err)
	assert.Contains(t, shortURL, "/test-code")

	// Verify mock expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetURLStats(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	service := NewService(db)

	// Mock database behavior
	mock.ExpectQuery("SELECT click_count, created_at, metadata").
		WithArgs("test-code").
		WillReturnRows(sqlmock.NewRows([]string{"click_count", "created_at", "metadata")).
		AddRow(10, time.Now(), "{}"))

	// Test case
	stats, err := service.GetURLStats("test-code")
	assert.NoError(t, err)
	assert.Equal(t, 10, stats.ClickCount)

	// Verify mock expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRedirectURL(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	service := NewService(db)

	// Mock database behavior
	mock.ExpectQuery("SELECT original_url, click_count").
		WithArgs("test-code").
		WillReturnRows(sqlmock.NewRows([]string{"original_url", "click_count")).
		AddRow("https://example.com", 0))

	mock.ExpectExec("UPDATE shortener.urls").
		WithArgs("test-code").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO shortener.analytics").
		WithArgs(sqlmock.AnyArg(), "127.0.0.1", "test-agent").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test case
	originalURL, err := service.RedirectURL("test-code", "127.0.0.1", "test-agent")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)

	// Verify mock expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
