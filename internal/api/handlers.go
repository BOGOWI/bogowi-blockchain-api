package api

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// TransferBOGOTokensRequest represents a BOGO token transfer request
type TransferBOGOTokensRequest struct {
	To     string `json:"to" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

// TransferBOGOTokens transfers BOGO tokens to a recipient
func (h *Handler) TransferBOGOTokens(c *gin.Context) {
	var req TransferBOGOTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Validate recipient address
	if !common.IsHexAddress(req.To) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid recipient address"})
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

	// Execute the transfer
	txHash, err := networkSDK.TransferBOGOTokens(req.To, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transfer initiated successfully",
		"transaction": txHash,
		"to":          req.To,
		"amount":      req.Amount,
	})
}
