package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHealth returns the API status and configured smart contract addresses
func (h *Handler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"contracts": h.Config.Testnet.Contracts,
	})
}

// GetGasPrice returns the current gas price
func (h *Handler) GetGasPrice(c *gin.Context) {
	gasPrice, err := h.SDK.GetGasPrice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get gas price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gasPrice": gasPrice,
	})
}
