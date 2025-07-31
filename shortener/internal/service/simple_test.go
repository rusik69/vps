package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		assert.True(t, true)
		assert.Equal(t, 1, 1)
	})

	t.Run("time test", func(t *testing.T) {
		now := time.Now()
		assert.NotZero(t, now)
	})
}
