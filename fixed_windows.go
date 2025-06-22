package rate_limiter

import (
	"sync"
	"time"
)

type fixedWindow struct {
	limit    int
	interval time.Duration
	requests map[string]int
	windows  map[string]time.Time
	mu       sync.Mutex
}

func NewFixedWindow(limit int, interval time.Duration) Limiter {
	return &fixedWindow{
		limit:    limit,
		interval: interval,
		requests: make(map[string]int),
		windows:  make(map[string]time.Time),
	}
}

func (fw *fixedWindow) Allow(userID string) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	windowStart := now.Truncate(fw.interval)

	if fw.windows[userID] != windowStart {
		fw.windows[userID] = windowStart
		fw.requests[userID] = 0
	}

	if fw.requests[userID] < fw.limit {
		fw.requests[userID]++
		return true
	}

	return false
}
