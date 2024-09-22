package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Define a struct to hold the limiter and last seen time for each client
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Create a map to store our rate limiters for each client
var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

// Create a new rate limiter for a client
func getClientLimiter(clientIP string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// Check if the client already has a rate limiter
	if _, exists := clients[clientIP]; !exists {
		// If not, create a new rate limiter: allow 10 requests per minute
		clients[clientIP] = &client{
			limiter: rate.NewLimiter(rate.Every(6*time.Second), 10),
		}
	}

	// Update the last seen time
	clients[clientIP].lastSeen = time.Now()

	return clients[clientIP].limiter
}

// RateLimiterMiddleware is a Gin middleware for rate limiting
func RateLimiterMiddleware() gin.HandlerFunc {
	go cleanupOldEntries() // Start the cleanup goroutine

	return func(c *gin.Context) {
		// Get the client's IP address
		clientIP := c.ClientIP()

		// Get (or create) the rate limiter for this client
		limiter := getClientLimiter(clientIP)

		// Check if this request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Periodically clean up old entries
func cleanupOldEntries() {
	for {
		time.Sleep(5 * time.Minute)

		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 10*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
