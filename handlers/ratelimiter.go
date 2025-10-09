package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// RateLimiter stores rate limiting data for each IP
type RateLimiter struct {
	mu       sync.RWMutex
	visitors map[string]*Visitor
	limit    int           // requests per window
	window   time.Duration // time window
}

// Visitor represents a visitor with their request count and window start time
type Visitor struct {
	count     int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		limit:    limit,
		window:   window,
	}

	// Clean up old visitors every 10 minutes
	go rl.cleanupVisitors()

	return rl
}

// cleanupVisitors removes visitors that haven't made requests recently
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastReset) > rl.window*2 {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if the IP is allowed to make a request
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	visitor, exists := rl.visitors[ip]

	if !exists {
		// First request from this IP
		rl.visitors[ip] = &Visitor{
			count:     1,
			lastReset: now,
		}
		return true
	}

	// Check if window has expired
	if now.Sub(visitor.lastReset) > rl.window {
		// Reset the window
		visitor.count = 1
		visitor.lastReset = now
		return true
	}

	// Check if under limit
	if visitor.count < rl.limit {
		visitor.count++
		return true
	}

	// Rate limit exceeded
	return false
}

// GetRemoteIP extracts the real IP address from the request
func GetRemoteIP(r *http.Request) string {
	// Check for forwarded IP (behind proxy/load balancer)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to remote address
	return r.RemoteAddr
}

// RateLimitMiddleware creates middleware for rate limiting
func RateLimitMiddleware(rl *RateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := GetRemoteIP(r)

			if !rl.Allow(ip) {
				// Rate limit exceeded
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				response := map[string]interface{}{
					"success": false,
					"error":   "Rate limit exceeded. Please try again later.",
					"message": "",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			// Continue to next handler
			next(w, r)
		}
	}
}
