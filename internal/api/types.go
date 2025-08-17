package api

import (
	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

// Handler holds the SDK and configuration
type Handler struct {
	SDK            SDKInterface
	NetworkHandler *NetworkHandler
	Config         *config.Config
	Storage        storage.RewardsStorage
}

// ErrorResponse is the standard error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse is the standard success response structure
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// rateLimitMiddleware creates a middleware for rate limiting
func rateLimitMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
