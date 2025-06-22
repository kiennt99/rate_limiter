package rate_limiter

import (
	"testing"
	"time"
)

func TestSlidingWindow_AllowsWithinLimit(t *testing.T) {
	limiter := NewSlidingWindow(3, 2*time.Second) // Max 3 requests per 2 seconds
	userID := "user1"

	// Send 3 requests quickly — all should be allowed
	for i := 0; i < 3; i++ {
		if !limiter.Allow(userID) {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	// 4th request should be rejected
	if limiter.Allow(userID) {
		t.Error("4th request should have been rejected")
	}
}

func TestSlidingWindow_ExpiresOldRequests(t *testing.T) {
	limiter := NewSlidingWindow(2, 1*time.Second)
	userID := "user2"

	// Send 2 requests
	if !limiter.Allow(userID) {
		t.Error("First request should be allowed")
	}
	if !limiter.Allow(userID) {
		t.Error("Second request should be allowed")
	}

	// Wait for 1.1 seconds for previous requests to expire
	time.Sleep(1100 * time.Millisecond)

	// New request should be allowed again
	if !limiter.Allow(userID) {
		t.Error("Request should be allowed after sliding window expires")
	}
}

func TestSlidingWindow_DifferentUsersIndependent(t *testing.T) {
	limiter := NewSlidingWindow(1, 1*time.Second)

	// Allow 1 request for user A
	if !limiter.Allow("userA") {
		t.Error("userA first request should be allowed")
	}

	// Allow 1 request for user B
	if !limiter.Allow("userB") {
		t.Error("userB first request should be allowed")
	}

	// Both should now be rate limited independently
	if limiter.Allow("userA") {
		t.Error("userA second request should be rejected")
	}
	if limiter.Allow("userB") {
		t.Error("userB second request should be rejected")
	}
}

func TestSlidingWindow_WindowMovesForward(t *testing.T) {
	limiter := NewSlidingWindow(2, 2*time.Second)
	userID := "user3"

	// Send 2 requests
	if !limiter.Allow(userID) {
		t.Error("First request should be allowed")
	}
	time.Sleep(1 * time.Second)
	if !limiter.Allow(userID) {
		t.Error("Second request should be allowed")
	}

	// Sleep until first request is out of window
	time.Sleep(1100 * time.Millisecond)

	if !limiter.Allow(userID) {
		t.Error("Third request should be allowed as oldest one expired")
	}
}
