package rate_limiter

import (
	"sync"
	"time"
)

type tokenBucket struct {
	capacity   int
	tokens     map[string]float64
	lastAccess map[string]time.Time
	rate       float64 // tokens/second
	mu         sync.Mutex
}

func NewTokenBucket(capacity int, rate float64) Limiter {
	return &tokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     make(map[string]float64),
		lastAccess: make(map[string]time.Time),
	}
}

func (tb *tokenBucket) Allow(userID string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	last := tb.lastAccess[userID]
	elapsed := now.Sub(last).Seconds()

	tb.tokens[userID] += elapsed * tb.rate
	if tb.tokens[userID] > float64(tb.capacity) {
		tb.tokens[userID] = float64(tb.capacity)
	}

	if tb.tokens[userID] >= 1 {
		tb.tokens[userID] -= 1
		tb.lastAccess[userID] = now
		return true
	}

	tb.lastAccess[userID] = now
	return false
}
