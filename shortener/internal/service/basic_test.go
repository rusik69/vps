package service_test

import (
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	if time.Now().IsZero() {
		t.Error("time should not be zero")
	}
}

func TestString(t *testing.T) {
	if "hello" != "hello" {
		t.Error("strings should match")
	}
}
