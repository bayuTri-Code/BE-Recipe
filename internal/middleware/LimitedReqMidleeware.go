package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string, maxRequests int, durationSeconds int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(
			rate.Every(time.Duration(durationSeconds)*time.Second/time.Duration(maxRequests)),
			maxRequests,
		)
		visitors[ip] = limiter
	}

	return limiter
}

func RateLimiter(maxRequests int, durationSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitor(ip, maxRequests, durationSeconds)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
