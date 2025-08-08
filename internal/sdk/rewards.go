package sdk

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// RewardTemplate represents a reward template from the contract
type RewardTemplate struct {
	ID                 string
	FixedAmount        *big.Int
	MaxAmount          *big.Int
	CooldownPeriod     *big.Int
	MaxClaimsPerWallet *big.Int
	RequiresWhitelist  bool
	Active             bool
}

// ClaimRewardV2 claims a fixed reward from a template
func (s *BOGOWISDK) ClaimRewardV2(templateID string, recipient common.Address) (*types.Transaction, error) {
	if s.rewardDistributor == nil {
		return nil, fmt.Errorf("reward distributor not initialized")
	}

	opts, err := s.getTransactOpts()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction options: %v", err)
	}

	// Call the contract method using the bound contract instance
	// The method signature is: claimReward(string templateId)
	tx, err := s.rewardDistributor.Instance.Transact(opts, "claimReward", templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute claimReward: %v", err)
	}

	return tx, nil
}

// ClaimCustomReward claims a custom amount reward (backend only)
func (s *BOGOWISDK) ClaimCustomReward(recipient common.Address, amount *big.Int, reason string) (*types.Transaction, error) {
	if s.rewardDistributor == nil {
		return nil, fmt.Errorf("reward distributor not initialized")
	}

	// Validate amount (max 1000 BOGO)
	maxAmount := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	if amount.Cmp(maxAmount) > 0 {
		return nil, fmt.Errorf("amount exceeds maximum of 1000 BOGO")
	}

	// Get transaction options
	opts, err := s.getTransactOpts()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction options: %v", err)
	}

	// Call the contract method using the bound contract instance
	// The method signature is: claimCustomReward(address recipient, uint256 amount, string reason)
	tx, err := s.rewardDistributor.Instance.Transact(opts, "claimCustomReward", recipient, amount, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to execute claimCustomReward: %v", err)
	}

	return tx, nil
}

// ClaimReferralBonus claims a referral bonus
func (s *BOGOWISDK) ClaimReferralBonus(referrer common.Address, referred common.Address) (*types.Transaction, error) {
	if s.rewardDistributor == nil {
		return nil, fmt.Errorf("reward distributor not initialized")
	}

	opts, err := s.getTransactOpts()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction options: %v", err)
	}

	// Call the contract method using the bound contract instance
	// The method signature is: claimReferralBonus(address referrer)
	tx, err := s.rewardDistributor.Instance.Transact(opts, "claimReferralBonus", referrer)
	if err != nil {
		return nil, fmt.Errorf("failed to execute claimReferralBonus: %v", err)
	}

	return tx, nil
}

// CheckRewardEligibility checks if a wallet is eligible for a reward
func (s *BOGOWISDK) CheckRewardEligibility(templateID string, wallet common.Address) (bool, string, error) {
	if s.rewardDistributor == nil {
		return false, "reward distributor not initialized", fmt.Errorf("reward distributor not initialized")
	}

	// For now, return mock eligibility
	// TODO: Call actual contract view method when ABI is available
	mockTemplates := map[string]struct {
		eligible bool
		reason   string
	}{
		"welcome_bonus":     {true, ""},
		"founder_bonus":     {false, "Not whitelisted"},
		"first_nft_mint":    {true, ""},
		"dao_participation": {false, "Cooldown period active"},
	}

	if tmpl, exists := mockTemplates[templateID]; exists {
		return tmpl.eligible, tmpl.reason, nil
	}

	return true, "", nil
}

// GetReferrer gets the referrer address for a wallet
func (s *BOGOWISDK) GetReferrer(wallet common.Address) (common.Address, error) {
	if s.rewardDistributor == nil {
		return common.Address{}, fmt.Errorf("reward distributor not initialized")
	}

	// For now, return zero address (not referred)
	// TODO: Call actual contract view method when ABI is available
	return common.Address{}, nil
}

// GetRewardTemplate gets details for a specific template
func (s *BOGOWISDK) GetRewardTemplate(templateID string) (*RewardTemplate, error) {
	if s.rewardDistributor == nil {
		return nil, fmt.Errorf("reward distributor not initialized")
	}

	// For now, return mock template data
	// TODO: Call actual contract view method when ABI is available
	templates := map[string]*RewardTemplate{
		"welcome_bonus": {
			ID:                 "welcome_bonus",
			FixedAmount:        new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)),
			MaxClaimsPerWallet: big.NewInt(1),
			Active:             true,
		},
		"founder_bonus": {
			ID:                 "founder_bonus",
			FixedAmount:        new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)),
			MaxClaimsPerWallet: big.NewInt(1),
			RequiresWhitelist:  true,
			Active:             true,
		},
	}

	if tmpl, exists := templates[templateID]; exists {
		return tmpl, nil
	}

	return nil, fmt.Errorf("template not found")
}

// GetClaimCount gets the number of times a wallet has claimed a template
func (s *BOGOWISDK) GetClaimCount(wallet common.Address, templateID string) (*big.Int, error) {
	if s.rewardDistributor == nil {
		return nil, fmt.Errorf("reward distributor not initialized")
	}

	// For now, return 0
	// TODO: Call actual contract view method when ABI is available
	return big.NewInt(0), nil
}

// IsWhitelisted checks if a wallet is whitelisted for founder bonus
func (s *BOGOWISDK) IsWhitelisted(wallet common.Address) (bool, error) {
	if s.rewardDistributor == nil {
		return false, fmt.Errorf("reward distributor not initialized")
	}

	// For now, return false
	// TODO: Call actual contract view method when ABI is available
	return false, nil
}

// GetRemainingDailyLimit gets the remaining daily distribution limit
func (s *BOGOWISDK) GetRemainingDailyLimit() (*big.Int, error) {
	if s.rewardDistributor == nil {
		return nil, fmt.Errorf("reward distributor not initialized")
	}

	// For now, return mock value (400k BOGO remaining)
	// TODO: Call actual contract view method when ABI is available
	return new(big.Int).Mul(big.NewInt(400000), big.NewInt(1e18)), nil
}


// Helper method to get transaction options
func (s *BOGOWISDK) getTransactOpts() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return nil, err
	}

	// Get suggested gas price from the network
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Add 20% buffer to suggested price to ensure transaction goes through
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(300000) // This could also be estimated dynamically

	return auth, nil
}
