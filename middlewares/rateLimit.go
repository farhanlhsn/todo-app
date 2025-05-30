package middlewares

import (
	"net/http"
	"sync"
	"time"
	"todo-app/helpers"

	"github.com/gin-gonic/gin"
)

// RateLimiter stores rate limiting data for each IP
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       *sync.RWMutex
}

// Visitor stores request information for each visitor
type Visitor struct {
	requests  []time.Time
	lastClean time.Time
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		mu:       &sync.RWMutex{},
	}

	// Clean up old visitors every 10 minutes
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			rl.cleanupVisitors()
		}
	}()

	return rl
}

// cleanupVisitors removes inactive visitors to prevent memory leaks
func (rl *RateLimiter) cleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-time.Hour) // Remove visitors inactive for 1 hour
	for ip, visitor := range rl.visitors {
		if visitor.lastClean.Before(cutoff) {
			delete(rl.visitors, ip)
		}
	}
}

// isAllowed checks if the visitor is allowed to make a request
func (rl *RateLimiter) isAllowed(ip string, maxRequests int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	visitor, exists := rl.visitors[ip]
	if !exists {
		visitor = &Visitor{
			requests:  make([]time.Time, 0),
			lastClean: now,
		}
		rl.visitors[ip] = visitor
	}

	// Remove old requests outside the time window
	var validRequests []time.Time
	for _, requestTime := range visitor.requests {
		if requestTime.After(cutoff) {
			validRequests = append(validRequests, requestTime)
		}
	}

	visitor.requests = validRequests
	visitor.lastClean = now

	// Check if under the limit
	if len(visitor.requests) >= maxRequests {
		return false
	}

	// Add current request
	visitor.requests = append(visitor.requests, now)
	return true
}

// Global rate limiter instance
var globalRateLimiter = NewRateLimiter()

// RateLimit creates a rate limiting middleware
// maxRequests: maximum number of requests allowed
// window: time window for the rate limit (e.g., time.Minute)
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !globalRateLimiter.isAllowed(ip, maxRequests, window) {
			c.JSON(http.StatusTooManyRequests, helpers.FormatResponseWithoutData(
				false,
				"Rate limit exceeded. Please try again later.",
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimit is a specific rate limiter for authentication endpoints
func AuthRateLimit() gin.HandlerFunc {
	return RateLimit(10, time.Minute) // 10 requests per minute for auth endpoints
}

// GeneralRateLimit is a general rate limiter for API endpoints
func GeneralRateLimit() gin.HandlerFunc {
	return RateLimit(100, time.Minute) // 100 requests per minute for general endpoints
}
