package rate_limiter

import (
	"testing"
	"time"
)

func TestFixedWindow_AllowsWithinLimit(t *testing.T) {
	limiter := NewFixedWindow(3, 2*time.Second)
	userID := "user1"

	for i := 0; i < 3; i++ {
		if !limiter.Allow(userID) {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	if limiter.Allow(userID) {
		t.Error("4th request should be rejected in the same window")
	}
}

func TestFixedWindow_ResetsAfterInterval(t *testing.T) {
	limiter := NewFixedWindow(2, 1*time.Second)
	userID := "user2"

	// Use 2 tokens
	if !limiter.Allow(userID) {
		t.Error("1st request should be allowed")
	}
	if !limiter.Allow(userID) {
		t.Error("2nd request should be allowed")
	}

	// Wait for new window
	time.Sleep(1100 * time.Millisecond)

	// Should reset and allow again
	if !limiter.Allow(userID) {
		t.Error("Request should be allowed after window reset")
	}
}

func TestFixedWindow_DifferentUsers(t *testing.T) {
	limiter := NewFixedWindow(1, 1*time.Second)

	if !limiter.Allow("userA") {
		t.Error("userA first request should be allowed")
	}
	if !limiter.Allow("userB") {
		t.Error("userB first request should be allowed")
	}

	if limiter.Allow("userA") {
		t.Error("userA second request should be rejected")
	}
	if limiter.Allow("userB") {
		t.Error("userB second request should be rejected")
	}
}
