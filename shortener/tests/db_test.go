package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rusik69/shortener/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryCreateShortURL(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test successful URL creation
	mock.ExpectQuery("INSERT INTO short_urls").
		WithArgs("abc12345", "https://example.com", nil, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	url, err := repo.CreateShortURL("abc12345", "https://example.com", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "abc12345", url.ShortCode)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetShortURLByCode(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test successful URL retrieval
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "short_code", "original_url", "user_id", "created_at", "updated_at", "expires_at", "click_count"}).
		AddRow(1, "abc12345", "https://example.com", nil, now, now, nil, 5)

	mock.ExpectQuery("SELECT (.+) FROM short_urls WHERE short_code = \\$1").
		WithArgs("abc12345").
		WillReturnRows(rows)

	url, err := repo.GetShortURLByCode("abc12345")
	assert.NoError(t, err)
	assert.Equal(t, "abc12345", url.ShortCode)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	assert.Equal(t, int64(5), url.ClickCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetShortURLByCodeNotFound(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test URL not found
	mock.ExpectQuery("SELECT (.+) FROM short_urls WHERE short_code = \\$1").
		WithArgs("notfound").
		WillReturnError(sql.ErrNoRows)

	url, err := repo.GetShortURLByCode("notfound")
	assert.Error(t, err)
	assert.Nil(t, url)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryIncrementClickCount(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test successful click increment
	mock.ExpectExec("UPDATE short_urls SET click_count = click_count \\+ 1").
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.IncrementClickCount(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryCreateClick(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test successful click recording
	mock.ExpectExec("INSERT INTO clicks").
		WithArgs(int64(1), "Mozilla/5.0", "192.168.1.1", "https://google.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateClick(1, "Mozilla/5.0", "192.168.1.1", "https://google.com")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetOrCreateRateLimit(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test existing rate limit retrieval
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "ip_address", "request_count", "reset_at", "created_at"}).
		AddRow(1, "192.168.1.1", 10, now.Add(time.Hour), now)

	mock.ExpectQuery("SELECT (.+) FROM rate_limits WHERE ip_address = \\$1").
		WithArgs("192.168.1.1").
		WillReturnRows(rows)

	rateLimit, err := repo.GetOrCreateRateLimit("192.168.1.1")
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", rateLimit.IPAddress)
	assert.Equal(t, int64(10), rateLimit.RequestCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetOrCreateRateLimitNew(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test new rate limit creation
	mock.ExpectQuery("SELECT (.+) FROM rate_limits WHERE ip_address = \\$1").
		WithArgs("192.168.1.2").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("INSERT INTO rate_limits").
		WithArgs("192.168.1.2", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	rateLimit, err := repo.GetOrCreateRateLimit("192.168.1.2")
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.2", rateLimit.IPAddress)
	assert.Equal(t, int64(1), rateLimit.RequestCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryUpdateRateLimit(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test successful rate limit update
	mock.ExpectExec("UPDATE rate_limits SET request_count = \\$1, reset_at = \\$2 WHERE id = \\$3").
		WithArgs(int64(5), sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rateLimit := &db.RateLimit{
		ID:           1,
		IPAddress:    "192.168.1.1",
		RequestCount: 5,
		ResetAt:      time.Now().Add(time.Hour),
	}

	err = repo.UpdateRateLimit(rateLimit)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryCreateCaptchaAttempt(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test successful captcha attempt creation
	mock.ExpectExec("INSERT INTO captcha_attempts").
		WithArgs("192.168.1.1", true, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateCaptchaAttempt("192.168.1.1", true)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetClicks(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = database.Close() }()

	repo := db.NewRepository(database)

	// Test successful clicks retrieval
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "short_url_id", "user_agent", "ip_address", "referrer", "created_at"}).
		AddRow(1, 1, "Mozilla/5.0", "192.168.1.1", "https://google.com", now).
		AddRow(2, 1, "Chrome/91.0", "192.168.1.2", "https://facebook.com", now)

	mock.ExpectQuery("SELECT (.+) FROM clicks WHERE short_url_id = \\$1 ORDER BY created_at DESC LIMIT \\$2").
		WithArgs(int64(1), 10).
		WillReturnRows(rows)

	clicks, err := repo.GetClicks(1, 10)
	assert.NoError(t, err)
	assert.Len(t, clicks, 2)
	assert.Equal(t, "Mozilla/5.0", clicks[0].UserAgent)
	assert.Equal(t, "Chrome/91.0", clicks[1].UserAgent)
	assert.NoError(t, mock.ExpectationsWereMet())
}
