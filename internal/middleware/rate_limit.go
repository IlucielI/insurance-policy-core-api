package middleware

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	
	// Rate limits
	publicRate  rate.Limit  // requests per second for public endpoints
	publicBurst int         // burst size for public
	authRate    rate.Limit  // requests per second for authenticated
	authBurst   int         // burst size for authenticated
	
	// Cleanup
	cleanupInterval time.Duration
}

// NewRateLimiter creates a new rate limiter
// publicRPM: requests per minute for public endpoints
// authRPM: requests per minute for authenticated users
func NewRateLimiter(publicRPM, authRPM int) *RateLimiter {
	rl := &RateLimiter{
		visitors:        make(map[string]*visitor),
		publicRate:      rate.Limit(float64(publicRPM) / 60.0),  // convert to per-second
		publicBurst:     publicRPM / 10,                          // 10% burst
		authRate:        rate.Limit(float64(authRPM) / 60.0),
		authBurst:       authRPM / 10,
		cleanupInterval: 3 * time.Minute,
	}
	
	// Start cleanup goroutine
	go rl.cleanupVisitors()
	
	return rl
}

// getVisitor returns the rate limiter for a given key
func (rl *RateLimiter) getVisitor(key string, isAuthenticated bool) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	v, exists := rl.visitors[key]
	if !exists {
		var limiter *rate.Limiter
		if isAuthenticated {
			limiter = rate.NewLimiter(rl.authRate, rl.authBurst)
		} else {
			limiter = rate.NewLimiter(rl.publicRate, rl.publicBurst)
		}
		
		rl.visitors[key] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}
	
	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupVisitors removes old entries from the map
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.cleanupInterval)
		
		rl.mu.Lock()
		for key, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanupInterval {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware returns a Fiber middleware for rate limiting
func (rl *RateLimiter) RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Determine if user is authenticated
		isAuthenticated := false
		var key string
		
		// Check if user has valid auth token
		claims, ok := c.Locals(string(UserContextKey)).(*UserClaims)
		if ok && claims != nil {
			isAuthenticated = true
			key = fmt.Sprintf("user:%s", claims.UserID)
		} else {
			// Use IP address for unauthenticated requests
			key = fmt.Sprintf("ip:%s", c.IP())
		}
		
		// Get limiter for this visitor
		limiter := rl.getVisitor(key, isAuthenticated)
		
		// Check if request is allowed
		if !limiter.Allow() {
			// Calculate retry-after in seconds
			reservation := limiter.Reserve()
			if !reservation.OK() {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
				})
			}
			
			delay := reservation.Delay()
			reservation.Cancel()
			
			retryAfter := int(delay.Seconds()) + 1
			
			// Log rate limit violation
			log.Printf("⚠️  Rate limit exceeded - Key: %s, Auth: %v, IP: %s, Path: %s",
				key, isAuthenticated, c.IP(), c.Path())
			
			c.Set("Retry-After", fmt.Sprintf("%d", retryAfter))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"retry_after": retryAfter,
				"message":     fmt.Sprintf("Too many requests. Please retry after %d seconds", retryAfter),
			})
		}
		
		return c.Next()
	}
}

// Global rate limiter instance
var globalRateLimiter *RateLimiter
var once sync.Once

// InitRateLimiter initializes the global rate limiter (call once at startup)
func InitRateLimiter(publicRPM, authRPM int) *RateLimiter {
	once.Do(func() {
		globalRateLimiter = NewRateLimiter(publicRPM, authRPM)
		log.Printf("✅ Rate limiter initialized - Public: %d req/min, Auth: %d req/min", publicRPM, authRPM)
	})
	return globalRateLimiter
}

// RateLimit returns the middleware using the global rate limiter
func RateLimit() fiber.Handler {
	if globalRateLimiter == nil {
		log.Fatal("Rate limiter not initialized. Call InitRateLimiter first.")
	}
	return globalRateLimiter.RateLimitMiddleware()
}
