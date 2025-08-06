package api

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// GetTokenBalance returns the BOGO token balance for a specific address
// @Summary Get BOGO token balance
// @Description Returns the balance of BOGO tokens for a given address
// @Tags Tokens
// @Param address path string true "Wallet address"
// @Param network query string false "Network (testnet or mainnet)"
// @Success 200 {object} sdk.TokenBalance
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /token/balance/{address} [get]
func (h *Handler) GetTokenBalance(c *gin.Context) {
	address := c.Param("address")

	// Validate Ethereum address
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Ethereum address"})
		return
	}

	// Get network parameter (required)
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

	balance, err := networkSDK.GetTokenBalance(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}
