package api

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/models"

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

	// Get templates from storage
	templates, err := h.Storage.GetAllRewardTemplates(c.Request.Context(), network, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"network":   network,
		"templates": templates,
	})
}

// GetRewardTemplateV2 returns a specific reward template
func (h *HandlerV2) GetRewardTemplateV2(c *gin.Context) {
	network := GetNetworkFromContext(c)
	templateID := c.Param("id")

	// Get template from storage
	template, err := h.Storage.GetRewardTemplate(c.Request.Context(), templateID, network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve template"})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: fmt.Sprintf("Template '%s' not found on %s", templateID, network),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CheckRewardEligibilityV2 checks if user is eligible for rewards
func (h *HandlerV2) CheckRewardEligibilityV2(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid") // From auth middleware
	templateID := c.Query("templateId")

	if templateID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "templateId is required"})
		return
	}

	// Check eligibility from storage
	eligibility, err := h.Storage.GetUserEligibility(c.Request.Context(), userID, templateID, network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to check eligibility"})
		return
	}

	if eligibility == nil {
		// No previous eligibility record, create default
		eligibility = &models.UserRewardEligibility{
			UserID:      userID,
			TemplateID:  templateID,
			IsEligible:  true,
			Reason:      "",
			ClaimCount:  0,
			Network:     network,
			LastChecked: time.Now(),
		}

		// Save eligibility
		h.Storage.SaveUserEligibility(c.Request.Context(), eligibility)
	}

	c.JSON(http.StatusOK, gin.H{
		"network":        network,
		"eligible":       eligibility.IsEligible,
		"userID":         userID,
		"templateId":     templateID,
		"reason":         eligibility.Reason,
		"claimCount":     eligibility.ClaimCount,
		"nextEligibleAt": eligibility.NextEligibleAt,
	})
}

// GetRewardHistoryV2 returns user's reward history
func (h *HandlerV2) GetRewardHistoryV2(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid")    // From auth middleware
	wallet := c.GetString("wallet") // From auth middleware

	if wallet == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Wallet not found in context"})
		return
	}

	// Get reward claims from storage
	rewardClaims, err := h.Storage.GetRewardClaimsByWallet(c.Request.Context(), wallet, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve reward history"})
		return
	}

	// Filter by network if needed
	var filteredClaims []gin.H
	for _, claim := range rewardClaims {
		if claim.Network == network {
			filteredClaims = append(filteredClaims, gin.H{
				"id":         claim.ID,
				"templateId": claim.TemplateID,
				"amount":     claim.Amount,
				"status":     claim.Status,
				"txHash":     claim.TxHash,
				"claimedAt":  claim.ClaimedAt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"network": network,
		"userID":  userID,
		"wallet":  wallet,
		"rewards": filteredClaims,
		"total":   len(filteredClaims),
	})
}

// ClaimRewardV3 processes reward claims
func (h *HandlerV2) ClaimRewardV3(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid")    // From auth middleware
	wallet := c.GetString("wallet") // From auth middleware

	var req struct {
		RewardType string `json:"rewardType" binding:"required"`
		Amount     string `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	if wallet == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Wallet not found in context"})
		return
	}

	// Check if template exists
	template, err := h.Storage.GetRewardTemplate(c.Request.Context(), req.RewardType, network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve template"})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Template '%s' not found", req.RewardType)})
		return
	}

	// Check eligibility
	eligibility, _ := h.Storage.GetUserEligibility(c.Request.Context(), userID, req.RewardType, network)
	if eligibility != nil && !eligibility.IsEligible {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: eligibility.Reason})
		return
	}

	// Determine amount
	amount := template.FixedAmount
	if req.Amount != "" {
		amount = req.Amount
	}

	// Create claim record
	claimRecord := &models.RewardClaim{
		WalletAddress: wallet,
		TemplateID:    req.RewardType,
		Amount:        amount,
		Status:        "pending",
		ClaimedAt:     time.Now(),
		Network:       network,
	}

	err = h.Storage.CreateRewardClaim(c.Request.Context(), claimRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create claim record"})
		return
	}

	// Get SDK for network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		h.Storage.UpdateRewardClaimStatus(c.Request.Context(), claimRecord.ID, "failed", "")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Network %s not available", network)})
		return
	}

	// Process claim through SDK
	walletAddr := common.HexToAddress(wallet)
	tx, err := sdk.ClaimRewardV2(req.RewardType, walletAddr)
	if err != nil {
		h.Storage.UpdateRewardClaimStatus(c.Request.Context(), claimRecord.ID, "failed", "")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error claiming reward: %v", err)})
		return
	}

	// Update claim status
	h.Storage.UpdateRewardClaimStatus(c.Request.Context(), claimRecord.ID, "completed", tx.Hash().Hex())

	// Update eligibility
	if eligibility == nil {
		eligibility = &models.UserRewardEligibility{
			UserID:     userID,
			TemplateID: req.RewardType,
			Network:    network,
		}
	}
	eligibility.ClaimCount++
	eligibility.LastChecked = time.Now()

	// Check if max claims reached
	if uint64(eligibility.ClaimCount) >= template.MaxClaimsPerWallet {
		eligibility.IsEligible = false
		eligibility.Reason = "Maximum claims reached"
	} else if template.CooldownPeriod > 0 {
		eligibility.NextEligibleAt = time.Now().Add(time.Duration(template.CooldownPeriod) * time.Second)
	}

	h.Storage.SaveUserEligibility(c.Request.Context(), eligibility)

	c.JSON(http.StatusOK, gin.H{
		"network":    network,
		"userID":     userID,
		"status":     "completed",
		"rewardType": req.RewardType,
		"txHash":     tx.Hash().Hex(),
		"claimId":    claimRecord.ID,
	})
}

// ClaimReferralV3 processes referral rewards
func (h *HandlerV2) ClaimReferralV3(c *gin.Context) {
	network := GetNetworkFromContext(c)
	userID := c.GetString("uid")    // From auth middleware
	wallet := c.GetString("wallet") // From auth middleware

	var req struct {
		ReferralCode string `json:"referralCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	if wallet == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Wallet not found in context"})
		return
	}

	// For demo purposes, treating referral code as the referrer's address
	// In production, this would map to actual referral codes
	if !common.IsHexAddress(req.ReferralCode) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid referral code format"})
		return
	}

	// Create referral claim record
	referralClaim := &models.ReferralClaim{
		ReferrerAddress: req.ReferralCode,
		ReferredAddress: wallet,
		ReferralCode:    req.ReferralCode,
		BonusAmount:     "5000000000000000000", // Default 5 BOGO
		Status:          "pending",
		ClaimedAt:       time.Now(),
		Network:         network,
	}

	err := h.Storage.CreateReferralClaim(c.Request.Context(), referralClaim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create referral claim"})
		return
	}

	// Get SDK for network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		h.Storage.UpdateReferralClaimStatus(c.Request.Context(), referralClaim.ID, "failed", "")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Network %s not available", network)})
		return
	}

	// Process referral through SDK
	referrerAddr := common.HexToAddress(req.ReferralCode)
	referredAddr := common.HexToAddress(wallet)

	tx, err := sdk.ClaimReferralBonus(referrerAddr, referredAddr)
	if err != nil {
		h.Storage.UpdateReferralClaimStatus(c.Request.Context(), referralClaim.ID, "failed", "")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error claiming referral: %v", err)})
		return
	}

	// Update status to completed
	h.Storage.UpdateReferralClaimStatus(c.Request.Context(), referralClaim.ID, "completed", tx.Hash().Hex())

	c.JSON(http.StatusOK, gin.H{
		"network":      network,
		"userID":       userID,
		"status":       "completed",
		"referralCode": req.ReferralCode,
		"txHash":       tx.Hash().Hex(),
		"claimId":      referralClaim.ID,
	})
}

// ClaimCustomRewardV3 processes custom rewards (backend only)
func (h *HandlerV2) ClaimCustomRewardV3(c *gin.Context) {
	network := GetNetworkFromContext(c)

	// Check backend authentication based on network
	authHeader := c.GetHeader("X-Backend-Auth")
	var expectedSecret string
	if network == "testnet" {
		expectedSecret = h.Config.DevBackendSecret
	} else {
		expectedSecret = h.Config.BackendSecret
	}

	if authHeader != expectedSecret {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

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

	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Parse amount from string to big.Int
	amount := new(big.Int)
	_, success := amount.SetString(req.Amount, 10)
	if !success {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid amount format"})
		return
	}

	// Call smart contract to distribute reward
	tx, err := sdk.ClaimCustomReward(common.HexToAddress(req.RecipientAddress), amount, req.RewardType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to claim reward: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"network":         network,
		"status":          "success",
		"transactionHash": tx.Hash().Hex(),
		"recipient":       req.RecipientAddress,
		"amount":          req.Amount,
		"rewardType":      req.RewardType,
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
