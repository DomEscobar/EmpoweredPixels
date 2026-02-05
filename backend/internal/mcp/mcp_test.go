package mcp

import (
	"context"
	"testing"
	"time"
)

func TestFairnessFilter(t *testing.T) {
	// Limit: 2 actions per 1 second
	filter := NewFairnessFilter(2, 1*time.Second)
	key := "test-ai-key-123456"

	// 1st action: OK
	if ok, _ := filter.Allow(key); !ok {
		t.Error("Should allow 1st action")
	}

	// 2nd action: OK
	if ok, _ := filter.Allow(key); !ok {
		t.Error("Should allow 2nd action")
	}

	// 3rd action: Should fail
	if ok, _ := filter.Allow(key); ok {
		t.Error("Should block 3rd action within same window")
	}

	// Wait for window to pass
	time.Sleep(1100 * time.Millisecond)

	// After wait: OK
	if ok, _ := filter.Allow(key); !ok {
		t.Error("Should allow action after window expired")
	}
}

func TestAuthentication(t *testing.T) {
	filter := NewFairnessFilter(10, 10*time.Second)
	handler := NewMCPHandler(filter)

	// Short key: Fail
	if _, err := handler.Authenticate(context.Background(), "short"); err == nil {
		t.Error("Should fail for short API key")
	}

	// Long key: OK
	if _, err := handler.Authenticate(context.Background(), "very-long-valid-api-key-12345"); err != nil {
		t.Errorf("Should authenticate valid key, got: %v", err)
	}
}
