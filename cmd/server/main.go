package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"rate_limiter"
	"time"
)

const (
	fixedWindows   = "fixed"
	slidingWindows = "sliding"
	bucketToken    = "token"
)

func main() {
	algo := flag.String("rate-limiter", "fixed", "Rate rate_limiter algorithm: fixed | sliding | token")
	limit := flag.Int("limit", 5, "Maximum requests per second per user")
	flag.Parse()

	var l rate_limiter.Limiter
	interval := time.Second

	switch *algo {
	case fixedWindows:
		l = rate_limiter.NewFixedWindow(*limit, interval)
	case slidingWindows:
		l = rate_limiter.NewSlidingWindow(*limit, interval)
	case bucketToken:
		l = rate_limiter.NewTokenBucket(*limit, float64(*limit)) // tokens/sec = limit
	default:
		log.Fatalf("Unknown rate rate_limiter: %s", *algo)
	}

	rateLimitMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("X-User-ID")

			if !l.Allow(userID) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request allowed")
	})

	fmt.Printf("Starting server on :8080 with %s rate_limiter (limit = %d/sec)\n", *algo, *limit)
	http.Handle("/", rateLimitMiddleware(handler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
