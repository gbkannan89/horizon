package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/errors"
)

type rateLimitEntry struct {
	count    int
	resetAt time.Time
}

var (
	rateLimits   = sync.Map{}
	rateLimitMtx sync.Mutex
)

// RateLimiter returns middleware that limits requests per IP.
func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		rateLimitMtx.Lock()
		val, _ := rateLimits.LoadOrStore(ip, &rateLimitEntry{count: 0, resetAt: now.Add(window)})
		entry := val.(*rateLimitEntry)

		if now.After(entry.resetAt) {
			entry.count = 0
			entry.resetAt = now.Add(window)
		}

		entry.count++
		remaining := maxRequests - entry.count
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", itoa(remaining))
		c.Header("X-RateLimit-Reset", itoa(int(entry.resetAt.Unix())))

		if entry.count > maxRequests {
			rateLimitMtx.Unlock()
			c.Header("Retry-After", itoa(int(window.Seconds())))
			errors.Respond(c, 429, "RATE_LIMIT_ERROR", "TOO_MANY_REQUESTS", "Rate limit exceeded. Try again later.")
			c.Abort()
			return
		}
		rateLimitMtx.Unlock()
		c.Next()
	}
}

func itoa(n int) string {
	if n == 0 { return "0" }
	s := ""
	for n > 0 { s = string(rune('0'+n%10)) + s; n /= 10 }
	return s
}
