package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NetworkMiddleware extracts and validates the network parameter
func NetworkMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		network := c.Query("network")
		if network == "" {
			network = c.GetHeader("X-Network")
		}
		if network == "" {
			network = "testnet" // Default to testnet
		}

		// Validate network
		if network != "testnet" && network != "mainnet" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid network. Use 'testnet' or 'mainnet'"})
			c.Abort()
			return
		}

		// Store in context
		c.Set("network", network)
		c.Next()
	}
}

// GetNetworkFromContext retrieves the network from gin context
func GetNetworkFromContext(c *gin.Context) string {
	network, _ := c.Get("network")
	if network != nil {
		return network.(string)
	}
	return "testnet" // Default to testnet if not set
}
