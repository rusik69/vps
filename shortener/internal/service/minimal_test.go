package service

import (
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	if time.Now().IsZero() {
		t.Error("time should not be zero")
	}
}

func TestServiceInterface(t *testing.T) {
	// Just test that the interface compiles
	var _ Service = (*service)(nil)
}
