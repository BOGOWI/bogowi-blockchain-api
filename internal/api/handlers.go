package api

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// NFT Handlers

// GetNFTBalance returns NFT balance for an address and token ID
func (h *Handler) GetNFTBalance(c *gin.Context) {
	address := c.Param("address")
	tokenID := c.Param("tokenId")

	// Validate Ethereum address
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Ethereum address"})
		return
	}

	// Get NFT balance from SDK
	balance, err := h.SDK.GetNFTBalance(address, tokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"tokenId": tokenID,
		"balance": balance,
	})
}

// MintEventTicketRequest represents the request to mint an event ticket
type MintEventTicketRequest struct {
	To        string `json:"to" binding:"required"`
	EventName string `json:"eventName" binding:"required"`
	EventDate string `json:"eventDate" binding:"required"`
}

// MintEventTicket mints a new event ticket NFT
func (h *Handler) MintEventTicket(c *gin.Context) {
	var req MintEventTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing required fields"})
		return
	}

	// Validate recipient address
	if !common.IsHexAddress(req.To) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid recipient address"})
		return
	}

	// Mint event ticket via SDK
	txHash, err := h.SDK.MintEventTicket(req.To, req.EventName, req.EventDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactionHash": txHash,
	})
}

// MintConservationNFTRequest represents the request to mint a conservation NFT
type MintConservationNFTRequest struct {
	To          string `json:"to" binding:"required"`
	TokenURI    string `json:"tokenURI" binding:"required"`
	Description string `json:"description"`
}

// MintConservationNFT mints a new conservation NFT
func (h *Handler) MintConservationNFT(c *gin.Context) {
	var req MintConservationNFTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing required fields"})
		return
	}

	// Validate recipient address
	if !common.IsHexAddress(req.To) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid recipient address"})
		return
	}

	// Mint conservation NFT via SDK
	txHash, err := h.SDK.MintConservationNFT(req.To, req.TokenURI, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactionHash": txHash,
	})
}

// Rewards Handlers

// GetRewardInfo returns reward information for an address
func (h *Handler) GetRewardInfo(c *gin.Context) {
	address := c.Param("address")

	// Validate Ethereum address
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Ethereum address"})
		return
	}

	// Get reward info from SDK
	info, err := h.SDK.GetRewardInfo(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetAchievementProgress returns achievement progress for an address
func (h *Handler) GetAchievementProgress(c *gin.Context) {
	address := c.Param("address")
	achievementID := c.Param("achievementId")

	// Validate Ethereum address
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Ethereum address"})
		return
	}

	// Get achievement progress from SDK
	progress, err := h.SDK.GetAchievementProgress(address, achievementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// ClaimRewardRequest represents a reward claim request
type ClaimRewardRequest struct {
	Address      string `json:"address" binding:"required"`
	RewardType   string `json:"rewardType" binding:"required"`
	RewardAmount string `json:"rewardAmount" binding:"required"`
}

// ClaimReward processes a reward claim
func (h *Handler) ClaimReward(c *gin.Context) {
	var req ClaimRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing required fields"})
		return
	}

	// Validate user address
	if !common.IsHexAddress(req.Address) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user address"})
		return
	}

	// Claim reward via SDK
	txHash, err := h.SDK.ClaimReward(req.Address, req.RewardType, req.RewardAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactionHash": txHash,
	})
}

// DAO Handlers

// GetDAOInfo returns DAO information
func (h *Handler) GetDAOInfo(c *gin.Context) {
	// For now, return a stub response
	c.JSON(http.StatusOK, gin.H{
		"threshold":        2,
		"signerCount":      5,
		"transactionCount": 0,
	})
}

// GetPendingTransactions returns pending multisig transactions
func (h *Handler) GetPendingTransactions(c *gin.Context) {
	// For now, return a stub response
	c.JSON(http.StatusOK, gin.H{
		"transactions": []interface{}{},
	})
}

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
