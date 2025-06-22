package rate_limiter

import (
	"sync"
	"time"
)

type slidingWindow struct {
	limit    int
	interval time.Duration
	requests map[string][]time.Time
	mu       sync.Mutex
}

func NewSlidingWindow(limit int, interval time.Duration) Limiter {
	return &slidingWindow{
		limit:    limit,
		interval: interval,
		requests: make(map[string][]time.Time),
	}
}

func (sw *slidingWindow) Allow(userID string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-sw.interval)

	var recent []time.Time
	for _, t := range sw.requests[userID] {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}

	if len(recent) < sw.limit {
		sw.requests[userID] = append(recent, now)
		return true
	}

	sw.requests[userID] = recent
	return false
}
