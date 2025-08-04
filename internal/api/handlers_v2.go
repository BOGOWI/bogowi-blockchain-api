package api

import (
	"fmt"
	"net/http"

	"bogowi-blockchain-go/internal/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// GetHealthV2 returns health status with network info
func (h *HandlerV2) GetHealthV2(c *gin.Context) {
	network := GetNetworkFromContext(c)

	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Get public key from SDK
	publicKey, err := sdk.GetPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get public key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"network":   network,
		"publicKey": publicKey,
		"version":   "2.0.0",
		"contracts": getContractAddresses(h.Config, network),
	})
}

// GetNetworkInfo returns information about the current network
func (h *HandlerV2) GetNetworkInfo(c *gin.Context) {
	network := GetNetworkFromContext(c)

	c.JSON(http.StatusOK, gin.H{
		"network":   network,
		"contracts": getContractAddresses(h.Config, network),
		"rpc_url":   getRPCUrl(h.Config, network),
		"chain_id":  getChainID(h.Config, network),
	})
}

// GetTokenBalanceV2 returns token balance for an address
func (h *HandlerV2) GetTokenBalanceV2(c *gin.Context) {
	network := GetNetworkFromContext(c)
	address := c.Param("address")

	// Validate Ethereum address
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Ethereum address"})
		return
	}

	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Get token balance
	balance, err := sdk.GetTokenBalance(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get token balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"network": network,
		"address": balance.Address,
		"balance": balance.Balance,
	})
}

// GetGasPriceV2 returns current gas price
func (h *HandlerV2) GetGasPriceV2(c *gin.Context) {
	network := GetNetworkFromContext(c)

	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	gasPrice, err := sdk.GetGasPrice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get gas price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"network":  network,
		"gasPrice": gasPrice,
	})
}

// TransferBOGOTokensV2 transfers BOGO tokens
func (h *HandlerV2) TransferBOGOTokensV2(c *gin.Context) {
	network := GetNetworkFromContext(c)

	var req struct {
		To     string `json:"to" binding:"required"`
		Amount string `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Validate address
	if !common.IsHexAddress(req.To) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid recipient address"})
		return
	}

	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Transfer tokens
	txHash, err := sdk.TransferBOGOTokens(req.To, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Transfer failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"network":         network,
		"transactionHash": txHash,
		"to":              req.To,
		"amount":          req.Amount,
	})
}

// GetRewardTemplatesV2 returns available reward templates
func (h *HandlerV2) GetRewardTemplatesV2(c *gin.Context) {
	network := GetNetworkFromContext(c)

	// TODO: Implement reward templates logic
	c.JSON(http.StatusOK, gin.H{
		"network":   network,
		"templates": []interface{}{},
	})
}

// GetRewardTemplateV2 returns a specific reward template
func (h *HandlerV2) GetRewardTemplateV2(c *gin.Context) {
	network := GetNetworkFromContext(c)
	templateID := c.Param("id")

	// TODO: Implement reward template retrieval logic
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: fmt.Sprintf("Template '%s' not found on %s", templateID, network),
	})
}

// CheckRewardEligibilityV2 checks if user is eligible for rewards
func (h *HandlerV2) CheckRewardEligibilityV2(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid") // From auth middleware

	// TODO: Implement eligibility check logic
	c.JSON(http.StatusOK, gin.H{
		"network":  network,
		"eligible": false,
		"userID":   userID,
	})
}

// GetRewardHistoryV2 returns user's reward history
func (h *HandlerV2) GetRewardHistoryV2(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid") // From auth middleware

	// TODO: Implement reward history logic
	c.JSON(http.StatusOK, gin.H{
		"network": network,
		"userID":  userID,
		"rewards": []interface{}{},
	})
}

// ClaimRewardV3 processes reward claims
func (h *HandlerV2) ClaimRewardV3(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid") // From auth middleware

	var req struct {
		RewardType string `json:"rewardType" binding:"required"`
		Amount     string `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// TODO: Implement reward claim logic
	c.JSON(http.StatusOK, gin.H{
		"network":    network,
		"userID":     userID,
		"status":     "pending",
		"rewardType": req.RewardType,
	})
}

// ClaimReferralV3 processes referral rewards
func (h *HandlerV2) ClaimReferralV3(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid") // From auth middleware

	var req struct {
		ReferralCode string `json:"referralCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// TODO: Implement referral claim logic
	c.JSON(http.StatusOK, gin.H{
		"network":      network,
		"userID":       userID,
		"status":       "pending",
		"referralCode": req.ReferralCode,
	})
}

// ClaimCustomRewardV3 processes custom rewards (backend only)
func (h *HandlerV2) ClaimCustomRewardV3(c *gin.Context) {
	network := GetNetworkFromContext(c)

	var req struct {
		RecipientAddress string `json:"recipientAddress" binding:"required"`
		Amount           string `json:"amount" binding:"required"`
		RewardType       string `json:"rewardType" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Validate address
	if !common.IsHexAddress(req.RecipientAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid recipient address"})
		return
	}

	// TODO: Implement custom reward distribution logic
	c.JSON(http.StatusOK, gin.H{
		"network":    network,
		"status":     "pending",
		"recipient":  req.RecipientAddress,
		"amount":     req.Amount,
		"rewardType": req.RewardType,
	})
}

// Helper functions
func getContractAddresses(cfg *config.Config, network string) map[string]string {
	var contracts config.ContractAddresses
	if network == "mainnet" {
		contracts = cfg.Mainnet.Contracts
	} else {
		contracts = cfg.Testnet.Contracts
	}

	return map[string]string{
		"roleManager":       contracts.RoleManager,
		"bogoToken":         contracts.BOGOToken,
		"rewardDistributor": contracts.RewardDistributor,
	}
}

func getRPCUrl(cfg *config.Config, network string) string {
	if network == "mainnet" {
		return cfg.Mainnet.RPCUrl
	}
	return cfg.Testnet.RPCUrl
}

func getChainID(cfg *config.Config, network string) int64 {
	if network == "mainnet" {
		return cfg.Mainnet.ChainID
	}
	return cfg.Testnet.ChainID
}
