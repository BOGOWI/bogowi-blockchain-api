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

	// Check for network parameter
	network := c.Query("network")
	if network == "" {
		network = c.GetHeader("X-Network")
	}
	
	// Use network-specific SDK if available and network is specified
	var sdkToUse SDKInterface
	if network != "" && h.NetworkHandler != nil {
		networkSDK, err := h.NetworkHandler.GetSDK(network)
		if err != nil {
			// Fall back to default SDK if network is invalid
			sdkToUse = h.SDK
		} else {
			sdkToUse = networkSDK
		}
	} else {
		// Use default SDK
		sdkToUse = h.SDK
	}

	balance, err := sdkToUse.GetTokenBalance(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}
