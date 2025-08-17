package api

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"bogowi-blockchain-go/internal/middleware"
	"bogowi-blockchain-go/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// Request structs for reward endpoints
type ClaimRewardRequest struct {
	TemplateID string `json:"templateId" binding:"required"`
}

type ClaimCustomRewardRequest struct {
	Wallet           string `json:"wallet,omitempty"`
	RecipientAddress string `json:"recipientAddress,omitempty"`
	Amount           string `json:"amount" binding:"required"`
	Reason           string `json:"reason,omitempty"`
	RewardType       string `json:"rewardType,omitempty"`
}

type ClaimReferralRequest struct {
	ReferrerAddress string `json:"referrerAddress" binding:"required"`
}

// Type aliases for backward compatibility with tests
type ClaimRewardRequestV2 = ClaimRewardRequest
type ClaimCustomRewardRequestV2 = ClaimCustomRewardRequest
type ClaimReferralRequestV2 = ClaimReferralRequest

// GetRewardTemplates returns all available reward templates
func (h *Handler) GetRewardTemplates(c *gin.Context) {
	templates := []gin.H{
		{"id": "welcome_bonus", "fixedAmount": "10", "maxClaimsPerWallet": 1, "active": true},
		{"id": "founder_bonus", "fixedAmount": "100", "maxClaimsPerWallet": 1, "requiresWhitelist": true, "active": true},
		{"id": "referral_bonus", "fixedAmount": "20", "active": true},
		{"id": "first_nft_mint", "fixedAmount": "25", "maxClaimsPerWallet": 1, "active": true},
		{"id": "dao_participation", "fixedAmount": "15", "cooldownPeriod": 2592000, "active": true},
		{"id": "attraction_tier_1", "fixedAmount": "10", "active": true},
		{"id": "attraction_tier_2", "fixedAmount": "20", "active": true},
		{"id": "attraction_tier_3", "fixedAmount": "40", "active": true},
		{"id": "attraction_tier_4", "fixedAmount": "50", "active": true},
		{"id": "custom_reward", "maxAmount": "1000", "active": true},
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
	})
}

// GetRewardTemplate returns a specific reward template
func (h *Handler) GetRewardTemplate(c *gin.Context) {
	templateID := c.Param("id")

	templates := map[string]gin.H{
		"welcome_bonus":     {"id": "welcome_bonus", "fixedAmount": "10", "maxClaimsPerWallet": 1, "active": true},
		"founder_bonus":     {"id": "founder_bonus", "fixedAmount": "100", "maxClaimsPerWallet": 1, "requiresWhitelist": true, "active": true},
		"referral_bonus":    {"id": "referral_bonus", "fixedAmount": "20", "active": true},
		"first_nft_mint":    {"id": "first_nft_mint", "fixedAmount": "25", "maxClaimsPerWallet": 1, "active": true},
		"dao_participation": {"id": "dao_participation", "fixedAmount": "15", "cooldownPeriod": 2592000, "active": true},
		"attraction_tier_1": {"id": "attraction_tier_1", "fixedAmount": "10", "active": true},
		"attraction_tier_2": {"id": "attraction_tier_2", "fixedAmount": "20", "active": true},
		"attraction_tier_3": {"id": "attraction_tier_3", "fixedAmount": "40", "active": true},
		"attraction_tier_4": {"id": "attraction_tier_4", "fixedAmount": "50", "active": true},
		"custom_reward":     {"id": "custom_reward", "maxAmount": "1000", "active": true},
	}

	template, exists := templates[templateID]
	if !exists {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// ClaimReward handles reward claiming with JWT auth
func (h *Handler) ClaimReward(c *gin.Context) {
	// Get authenticated user's wallet from context
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
		return
	}

	var req ClaimRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get wallet address from claims
	firebaseClaims := claims.(*middleware.FirebaseClaims)
	wallet := firebaseClaims.WalletAddress
	if wallet == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Wallet address not found in token"})
		return
	}

	// Convert wallet to address
	if !common.IsHexAddress(wallet) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid wallet address"})
		return
	}
	walletAddr := common.HexToAddress(wallet)

	// Check eligibility
	eligible, message, err := h.SDK.CheckRewardEligibility(req.TemplateID, walletAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to check eligibility: %v", err)})
		return
	}

	if !eligible {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: message})
		return
	}

	// Claim the reward using the clean interface method
	tx, err := h.SDK.ClaimRewardV2(req.TemplateID, walletAddr) // TODO: SDK method needs renaming too
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to claim reward: %v", err)})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"transactionHash": tx.Hash().Hex(),
		"wallet":          wallet,
		"templateId":      req.TemplateID,
	})
}

// ClaimReferralBonus handles referral bonus claims
func (h *Handler) ClaimReferralBonus(c *gin.Context) {
	// Get authenticated user's wallet from context
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
		return
	}

	var req ClaimReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Validate addresses
	if !common.IsHexAddress(req.ReferrerAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid referrer address"})
		return
	}

	// Get referred wallet from claims
	firebaseClaims := claims.(*middleware.FirebaseClaims)
	referredWallet := firebaseClaims.WalletAddress
	if referredWallet == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Wallet address not found in token"})
		return
	}

	if !common.IsHexAddress(referredWallet) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid referred wallet address"})
		return
	}

	referrerAddr := common.HexToAddress(req.ReferrerAddress)
	referredAddr := common.HexToAddress(referredWallet)

	// Claim referral bonus
	tx, err := h.SDK.ClaimReferralBonus(referrerAddr, referredAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to claim referral bonus: %v", err)})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"transactionHash": tx.Hash().Hex(),
		"referrer":        req.ReferrerAddress,
		"referred":        referredWallet,
	})
}

// ClaimCustomReward handles custom reward claims with optional network support
func (h *Handler) ClaimCustomReward(c *gin.Context) {
	// Determine which SDK to use
	var sdk SDKInterface
	if h.NetworkHandler != nil {
		// Get network from query or header
		network := c.Query("network")
		if network == "" {
			network = c.GetHeader("X-Network")
		}
		if network == "" {
			network = "testnet" // default
		}

		// Get SDK for the specified network
		var err error
		sdk, err = h.NetworkHandler.GetSDK(network)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Failed to get SDK for network %s: %v", network, err)})
			return
		}
	} else {
		// Use default SDK
		sdk = h.SDK
	}

	// Authenticate backend request
	if !h.authenticateBackendRequest(c) {
		return
	}

	var req ClaimCustomRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Determine recipient address
	recipientAddress := req.RecipientAddress
	if recipientAddress == "" && req.Wallet != "" {
		recipientAddress = req.Wallet
	}
	if recipientAddress == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Recipient address is required"})
		return
	}

	// Validate address
	if !common.IsHexAddress(recipientAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid recipient address"})
		return
	}

	// Parse amount
	amount := new(big.Int)
	_, ok := amount.SetString(req.Amount, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid amount format"})
		return
	}

	// Check max amount (1000 BOGO = 1000 * 10^18 wei)
	maxAmount := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	if amount.Cmp(maxAmount) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Amount exceeds maximum (1000 BOGO)"})
		return
	}

	// Use reason or rewardType (for backward compatibility)
	reason := req.Reason
	if reason == "" {
		reason = req.RewardType
	}
	if reason == "" {
		reason = "custom_reward"
	}

	recipientAddr := common.HexToAddress(recipientAddress)

	// Claim custom reward
	tx, err := sdk.ClaimCustomReward(recipientAddr, amount, reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to claim custom reward: %v", err)})
		return
	}

	// Return success response
	response := gin.H{
		"success":         true,
		"transactionHash": tx.Hash().Hex(),
		"recipient":       recipientAddress,
		"amount":          req.Amount,
		"reason":          reason,
	}

	// Add network info if using network handler
	if h.NetworkHandler != nil {
		network := c.Query("network")
		if network == "" {
			network = c.GetHeader("X-Network")
		}
		if network == "" {
			network = "testnet"
		}
		response["network"] = network
	}

	c.JSON(http.StatusOK, response)
}

// authenticateBackendRequest checks if the request is from a trusted backend
func (h *Handler) authenticateBackendRequest(c *gin.Context) bool {
	// Support both Authorization and X-Backend-Auth headers for backward compatibility
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		authHeader = c.GetHeader("X-Backend-Auth")
	}
	network := c.Query("network")
	if network == "" {
		network = c.GetHeader("X-Network")
	}
	if network == "" {
		network = "testnet"
	}

	var expectedSecret string
	if network == "testnet" {
		expectedSecret = h.Config.DevBackendSecret
		if expectedSecret == "" {
			expectedSecret = h.Config.BackendSecret // fallback
		}
	} else {
		expectedSecret = h.Config.BackendSecret
	}

	if authHeader != expectedSecret {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid backend authentication"})
		return false
	}

	return true
}

// CheckRewardEligibility checks what rewards a user can claim
func (h *Handler) CheckRewardEligibility(c *gin.Context) {
	wallet, exists := c.Get("wallet")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	walletAddr := common.HexToAddress(wallet.(string))
	templateID := c.Query("templateId")

	var eligibilities []gin.H

	if templateID != "" {
		// Check specific template
		eligible, reason, err := h.SDK.CheckRewardEligibility(templateID, walletAddr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error checking eligibility: %v", err)})
			return
		}

		eligibilities = append(eligibilities, gin.H{
			"templateId": templateID,
			"eligible":   eligible,
			"reason":     reason,
		})
	} else {
		// Check all templates
		templates := []string{
			"welcome_bonus", "founder_bonus", "first_nft_mint",
			"dao_participation", "attraction_tier_1", "attraction_tier_2",
			"attraction_tier_3", "attraction_tier_4",
		}

		for _, tmpl := range templates {
			eligible, reason, _ := h.SDK.CheckRewardEligibility(tmpl, walletAddr)
			eligibilities = append(eligibilities, gin.H{
				"templateId": tmpl,
				"eligible":   eligible,
				"reason":     reason,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"eligibilities": eligibilities,
	})
}

// GetRewardHistory returns the user's reward claim history
func (h *Handler) GetRewardHistory(c *gin.Context) {
	wallet, exists := c.Get("wallet")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	walletAddr := wallet.(string)

	// Get reward claims from storage
	rewardClaims, err := h.Storage.GetRewardClaimsByWallet(c.Request.Context(), walletAddr, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve claim history"})
		return
	}

	// Get referral claims from storage
	referralClaims, err := h.Storage.GetReferralClaimsByWallet(c.Request.Context(), walletAddr, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve referral history"})
		return
	}

	// Combine and format the claims
	var allClaims []gin.H

	for _, claim := range rewardClaims {
		allClaims = append(allClaims, gin.H{
			"type":       "reward",
			"templateId": claim.TemplateID,
			"amount":     claim.Amount,
			"status":     claim.Status,
			"txHash":     claim.TxHash,
			"claimedAt":  claim.ClaimedAt,
			"network":    claim.Network,
		})
	}

	for _, claim := range referralClaims {
		allClaims = append(allClaims, gin.H{
			"type":            "referral",
			"referrerAddress": claim.ReferrerAddress,
			"bonusAmount":     claim.BonusAmount,
			"status":          claim.Status,
			"txHash":          claim.TxHash,
			"claimedAt":       claim.ClaimedAt,
			"network":         claim.Network,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet": walletAddr,
		"claims": allClaims,
		"total":  len(allClaims),
	})
}

// AuthMiddleware wraps the JWT authentication for Gin
func AuthMiddleware(auth *middleware.AuthMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the Gin context to work with the auth middleware
		auth.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get wallet from context and set it in Gin context
			wallet, err := middleware.GetWalletFromContext(r.Context())
			if err != nil {
				c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
				c.Abort()
				return
			}

			c.Set("wallet", wallet)
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}

// ======= BACKWARD COMPATIBILITY - V2 Methods =======
// These are kept for backward compatibility but internally call the new methods

// ClaimRewardV2 - DEPRECATED: Use ClaimReward instead
func (h *Handler) ClaimRewardV2(c *gin.Context) {
	// Get wallet from JWT context (set by middleware)
	wallet, exists := c.Get("wallet")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	var req ClaimRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	walletAddr := common.HexToAddress(wallet.(string))

	// Check eligibility first
	eligible, reason, err := h.SDK.CheckRewardEligibility(req.TemplateID, walletAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error checking eligibility: %v", err)})
		return
	}

	if !eligible {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: reason})
		return
	}

	// Get template info for amount
	template, err := h.SDK.GetRewardTemplate(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error getting template: %v", err)})
		return
	}

	// Store claim record with pending status
	claimRecord := &models.RewardClaim{
		WalletAddress: wallet.(string),
		TemplateID:    req.TemplateID,
		Amount:        template.FixedAmount.String(),
		Status:        "pending",
		ClaimedAt:     time.Now(),
		Network:       "camino",
	}

	err = h.Storage.CreateRewardClaim(c.Request.Context(), claimRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to record claim"})
		return
	}

	// Claim reward
	tx, err := h.SDK.ClaimRewardV2(req.TemplateID, walletAddr)
	if err != nil {
		// Update claim status to failed
		h.Storage.UpdateRewardClaimStatus(c.Request.Context(), claimRecord.ID, "failed", "")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error claiming reward: %v", err)})
		return
	}

	// Update claim status to completed
	h.Storage.UpdateRewardClaimStatus(c.Request.Context(), claimRecord.ID, "completed", tx.Hash().Hex())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"txHash":  tx.Hash().Hex(),
		"message": fmt.Sprintf("Successfully claimed %s", req.TemplateID),
		"claimId": claimRecord.ID,
	})
}

// ClaimCustomRewardUnified - DEPRECATED: Use ClaimCustomReward instead
func (h *Handler) ClaimCustomRewardUnified(c *gin.Context) {
	h.ClaimCustomReward(c)
}

// ClaimCustomRewardV2 - DEPRECATED: Use ClaimCustomReward instead
func (h *Handler) ClaimCustomRewardV2(c *gin.Context) {
	h.ClaimCustomReward(c)
}

// ClaimReferralV2 - DEPRECATED: Use ClaimReferralBonus instead
func (h *Handler) ClaimReferralV2(c *gin.Context) {
	h.ClaimReferralBonus(c)
}

// ClaimCustomRewardV2WithNetwork - DEPRECATED: Use ClaimCustomReward instead
func (h *Handler) ClaimCustomRewardV2WithNetwork(c *gin.Context) {
	h.ClaimCustomReward(c)
}
