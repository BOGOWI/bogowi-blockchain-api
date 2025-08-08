package api

import (
	"math/big"

	"bogowi-blockchain-go/internal/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SDKInterface defines the methods we use from the SDK
type SDKInterface interface {
	GetTokenBalance(address string) (*sdk.TokenBalance, error)
	GetGasPrice() (string, error)
	TransferBOGOTokens(to string, amount string) (string, error)
	GetPublicKey() (string, error)
	Close()

	// New reward system methods
	CheckRewardEligibility(templateID string, wallet common.Address) (bool, string, error)
	ClaimRewardV2(templateID string, recipient common.Address) (*types.Transaction, error)
	ClaimCustomReward(recipient common.Address, amount *big.Int, reason string) (*types.Transaction, error)
	ClaimReferralBonus(referrer common.Address, referred common.Address) (*types.Transaction, error)
	GetReferrer(wallet common.Address) (common.Address, error)
	GetRewardTemplate(templateID string) (*sdk.RewardTemplate, error)
	GetClaimCount(wallet common.Address, templateID string) (*big.Int, error)
	IsWhitelisted(wallet common.Address) (bool, error)
	GetRemainingDailyLimit() (*big.Int, error)
}
