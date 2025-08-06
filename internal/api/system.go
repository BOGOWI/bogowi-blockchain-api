package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHealth returns the API status and configured smart contract addresses
func (h *Handler) GetHealth(c *gin.Context) {
	// Get network parameter
	network := c.Query("network")
	if network == "" {
		network = c.GetHeader("X-Network")
	}
	if network == "" {
		network = "mainnet" // Default to mainnet if not specified
	}

	// Get contracts for the specified network
	var contracts interface{}
	if network == "testnet" || network == "columbus" {
		contracts = h.Config.Testnet.Contracts
	} else {
		contracts = h.Config.Mainnet.Contracts
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"network":   network,
		"contracts": contracts,
	})
}

// GetGasPrice returns the current gas price
func (h *Handler) GetGasPrice(c *gin.Context) {
	// Get network parameter
	network := c.Query("network")
	if network == "" {
		network = c.GetHeader("X-Network")
	}
	if network == "" {
		network = "mainnet" // Default to mainnet if not specified
	}

	// Get network-specific SDK
	if h.NetworkHandler == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Network handler not initialized"})
		return
	}

	networkSDK, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid network: " + network + ". Use 'testnet' or 'mainnet'"})
		return
	}

	gasPrice, err := networkSDK.GetGasPrice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get gas price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"network":  network,
		"gasPrice": gasPrice,
	})
}
