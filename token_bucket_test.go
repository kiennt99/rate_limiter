package rate_limiter

import (
	"testing"
	"time"
)

func TestTokenBucket_AllowsWithinCapacity(t *testing.T) {
	limiter := NewTokenBucket(5, 1.0) // 5 token capacity, 1 token/sec
	userID := "user1"

	for i := 0; i < 5; i++ {
		if !limiter.Allow(userID) {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	// 6th request should be rejected immediately
	if limiter.Allow(userID) {
		t.Error("6th request should have been rejected due to no tokens left")
	}
}

func TestTokenBucket_RefillsTokensOverTime(t *testing.T) {
	limiter := NewTokenBucket(2, 1.0) // 2 token cap, 1 token/sec
	userID := "user2"

	// Use all tokens
	limiter.Allow(userID)
	limiter.Allow(userID)

	// Should be blocked now
	if limiter.Allow(userID) {
		t.Error("Request should have been rejected")
	}

	// Wait 2.1 seconds to accumulate 2 new tokens
	time.Sleep(2100 * time.Millisecond)

	// Should be allowed again
	if !limiter.Allow(userID) {
		t.Error("Request should have been allowed after token refill")
	}
	if !limiter.Allow(userID) {
		t.Error("Second request after refill should also be allowed")
	}
}

func TestTokenBucket_CapacityDoesNotOverflow(t *testing.T) {
	limiter := NewTokenBucket(3, 1.0) // max 3 tokens
	userID := "user3"

	// Wait 5 seconds — should still only have 3 tokens max
	time.Sleep(5 * time.Second)

	// Try to send 4 requests, 4th should fail
	for i := 0; i < 3; i++ {
		if !limiter.Allow(userID) {
			t.Errorf("Request %d should be allowed after refill", i+1)
		}
	}
	if limiter.Allow(userID) {
		t.Error("4th request should be rejected, capacity is only 3")
	}
}
