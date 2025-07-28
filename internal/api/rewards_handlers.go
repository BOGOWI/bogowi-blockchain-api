package api

import (
	"fmt"
	"math/big"
	"net/http"

	"bogowi-blockchain-go/internal/middleware"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type ClaimRewardRequestV2 struct {
	TemplateID string `json:"templateId" binding:"required"`
}

type ClaimCustomRewardRequestV2 struct {
	Wallet string `json:"wallet" binding:"required"`
	Amount string `json:"amount" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

type ClaimReferralRequestV2 struct {
	ReferrerAddress string `json:"referrerAddress" binding:"required"`
}

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
		"welcome_bonus":      {"id": "welcome_bonus", "fixedAmount": "10", "maxClaimsPerWallet": 1, "active": true},
		"founder_bonus":      {"id": "founder_bonus", "fixedAmount": "100", "maxClaimsPerWallet": 1, "requiresWhitelist": true, "active": true},
		"referral_bonus":     {"id": "referral_bonus", "fixedAmount": "20", "active": true},
		"first_nft_mint":     {"id": "first_nft_mint", "fixedAmount": "25", "maxClaimsPerWallet": 1, "active": true},
		"dao_participation":  {"id": "dao_participation", "fixedAmount": "15", "cooldownPeriod": 2592000, "active": true},
		"attraction_tier_1":  {"id": "attraction_tier_1", "fixedAmount": "10", "active": true},
		"attraction_tier_2":  {"id": "attraction_tier_2", "fixedAmount": "20", "active": true},
		"attraction_tier_3":  {"id": "attraction_tier_3", "fixedAmount": "40", "active": true},
		"attraction_tier_4":  {"id": "attraction_tier_4", "fixedAmount": "50", "active": true},
		"custom_reward":      {"id": "custom_reward", "maxAmount": "1000", "active": true},
	}

	template, exists := templates[templateID]
	if !exists {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// ClaimRewardV2 handles reward claiming with JWT auth
func (h *Handler) ClaimRewardV2(c *gin.Context) {
	// Get wallet from JWT context (set by middleware)
	wallet, exists := c.Get("wallet")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	var req ClaimRewardRequestV2
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

	// Claim reward
	tx, err := h.SDK.ClaimRewardV2(req.TemplateID, walletAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error claiming reward: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"txHash":  tx.Hash().Hex(),
		"message": fmt.Sprintf("Successfully claimed %s", req.TemplateID),
	})
}

// ClaimReferralV2 handles referral bonus claims
func (h *Handler) ClaimReferralV2(c *gin.Context) {
	wallet, exists := c.Get("wallet")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	var req ClaimReferralRequestV2
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	if !common.IsHexAddress(req.ReferrerAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid referrer address"})
		return
	}

	referrerAddr := common.HexToAddress(req.ReferrerAddress)
	referredAddr := common.HexToAddress(wallet.(string))

	// Check if already referred
	existingReferrer, err := h.SDK.GetReferrer(referredAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error checking referral status: %v", err)})
		return
	}

	if existingReferrer != (common.Address{}) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Already referred"})
		return
	}

	// Claim referral bonus
	tx, err := h.SDK.ClaimReferralBonus(referrerAddr, referredAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error claiming referral: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"txHash":  tx.Hash().Hex(),
		"message": "Referral bonus claimed successfully",
	})
}

// ClaimCustomRewardV2 handles custom reward claims (backend only)
func (h *Handler) ClaimCustomRewardV2(c *gin.Context) {
	// This endpoint requires backend service authentication
	authHeader := c.GetHeader("X-Backend-Auth")
	if authHeader != h.Config.BackendSecret { // Add BackendSecret to config
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	var req ClaimCustomRewardRequestV2
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	if !common.IsHexAddress(req.Wallet) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid wallet address"})
		return
	}

	// Validate amount
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid amount format"})
		return
	}

	// Convert to wei (multiply by 10^18)
	weiAmount := new(big.Int).Mul(amount, big.NewInt(1e18))

	// Check max amount (1000 BOGO)
	maxAmount := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	if weiAmount.Cmp(maxAmount) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Amount exceeds maximum (1000 BOGO)"})
		return
	}

	walletAddr := common.HexToAddress(req.Wallet)

	// Claim custom reward
	tx, err := h.SDK.ClaimCustomReward(walletAddr, weiAmount, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error claiming custom reward: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"txHash":  tx.Hash().Hex(),
		"amount":  req.Amount,
		"reason":  req.Reason,
	})
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

	// TODO: Implement claim history from database
	c.JSON(http.StatusOK, gin.H{
		"wallet": wallet,
		"claims": []interface{}{},
		"message": "Claim history not yet implemented",
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