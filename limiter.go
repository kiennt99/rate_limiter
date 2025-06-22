package rate_limiter

type Limiter interface {
	Allow(userID string) bool
}
