package api

import (
	"math/big"

	"bogowi-blockchain-go/internal/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SimpleMockSDK is a simple mock implementation of SDKInterface for testing
type SimpleMockSDK struct {
	// Control behavior
	ShouldFail  bool
	FailMessage string

	// Return values
	Balance         *sdk.TokenBalance
	GasPrice        string
	TransactionHash string
	PublicKey       string

	// Track calls
	Calls []string
}

// NewSimpleMockSDK creates a new mock SDK with defaults
func NewSimpleMockSDK() *SimpleMockSDK {
	return &SimpleMockSDK{
		Balance: &sdk.TokenBalance{
			Address: "0x1234567890123456789012345678901234567890",
			Balance: "1000000000000000000", // 1 token
		},
		GasPrice:        "1000000000", // 1 gwei
		TransactionHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		PublicKey:       "0x04abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}
}

// GetTokenBalance implements SDKInterface
func (m *SimpleMockSDK) GetTokenBalance(address string) (*sdk.TokenBalance, error) {
	m.Calls = append(m.Calls, "GetTokenBalance")
	if m.ShouldFail {
		return nil, &MockError{Message: m.FailMessage}
	}
	return m.Balance, nil
}

// GetGasPrice implements SDKInterface
func (m *SimpleMockSDK) GetGasPrice() (string, error) {
	m.Calls = append(m.Calls, "GetGasPrice")
	if m.ShouldFail {
		return "", &MockError{Message: m.FailMessage}
	}
	return m.GasPrice, nil
}

// TransferBOGOTokens implements SDKInterface
func (m *SimpleMockSDK) TransferBOGOTokens(to string, amount string) (string, error) {
	m.Calls = append(m.Calls, "TransferBOGOTokens")
	if m.ShouldFail {
		return "", &MockError{Message: m.FailMessage}
	}
	return m.TransactionHash, nil
}

// GetPublicKey implements SDKInterface
func (m *SimpleMockSDK) GetPublicKey() (string, error) {
	m.Calls = append(m.Calls, "GetPublicKey")
	if m.ShouldFail {
		return "", &MockError{Message: m.FailMessage}
	}
	return m.PublicKey, nil
}

// Close implements SDKInterface
func (m *SimpleMockSDK) Close() {
	m.Calls = append(m.Calls, "Close")
}

// CheckRewardEligibility implements SDKInterface
func (m *SimpleMockSDK) CheckRewardEligibility(templateID string, wallet common.Address) (bool, string, error) {
	m.Calls = append(m.Calls, "CheckRewardEligibility")
	if m.ShouldFail {
		return false, "", &MockError{Message: m.FailMessage}
	}
	return true, "Eligible for reward", nil
}

// ClaimRewardV2 implements SDKInterface
func (m *SimpleMockSDK) ClaimRewardV2(templateID string, recipient common.Address) (*types.Transaction, error) {
	m.Calls = append(m.Calls, "ClaimRewardV2")
	if m.ShouldFail {
		return nil, &MockError{Message: m.FailMessage}
	}
	// Create a mock transaction
	tx := types.NewTransaction(
		1,                               // nonce
		recipient,                       // to
		big.NewInt(1000000000000000000), // value (1 token)
		21000,                           // gas limit
		big.NewInt(1000000000),          // gas price (1 gwei)
		nil,                             // data
	)
	return tx, nil
}

// ClaimCustomReward implements SDKInterface
func (m *SimpleMockSDK) ClaimCustomReward(recipient common.Address, amount *big.Int, reason string) (*types.Transaction, error) {
	m.Calls = append(m.Calls, "ClaimCustomReward")
	if m.ShouldFail {
		return nil, &MockError{Message: m.FailMessage}
	}
	// Create a mock transaction
	tx := types.NewTransaction(
		1,                      // nonce
		recipient,              // to
		amount,                 // value
		21000,                  // gas limit
		big.NewInt(1000000000), // gas price (1 gwei)
		nil,                    // data
	)
	return tx, nil
}

// ClaimReferralBonus implements SDKInterface
func (m *SimpleMockSDK) ClaimReferralBonus(referrer common.Address, referred common.Address) (*types.Transaction, error) {
	m.Calls = append(m.Calls, "ClaimReferralBonus")
	if m.ShouldFail {
		return nil, &MockError{Message: m.FailMessage}
	}
	// Create a mock transaction
	tx := types.NewTransaction(
		1,                              // nonce
		referrer,                       // to
		big.NewInt(500000000000000000), // value (0.5 token)
		21000,                          // gas limit
		big.NewInt(1000000000),         // gas price (1 gwei)
		nil,                            // data
	)
	return tx, nil
}

// GetReferrer implements SDKInterface
func (m *SimpleMockSDK) GetReferrer(wallet common.Address) (common.Address, error) {
	m.Calls = append(m.Calls, "GetReferrer")
	if m.ShouldFail {
		return common.Address{}, &MockError{Message: m.FailMessage}
	}
	// Return zero address (no referrer)
	return common.Address{}, nil
}

// GetRewardTemplate implements SDKInterface
func (m *SimpleMockSDK) GetRewardTemplate(templateID string) (*sdk.RewardTemplate, error) {
	m.Calls = append(m.Calls, "GetRewardTemplate")
	if m.ShouldFail {
		return nil, &MockError{Message: m.FailMessage}
	}
	// 1 token = 10^18, 10 tokens = 10 * 10^18
	fixedAmount := new(big.Int)
	fixedAmount.SetString("1000000000000000000", 10)
	maxAmount := new(big.Int)
	maxAmount.SetString("10000000000000000000", 10)

	return &sdk.RewardTemplate{
		ID:                 templateID,
		FixedAmount:        fixedAmount,
		MaxAmount:          maxAmount,
		CooldownPeriod:     big.NewInt(86400), // 1 day
		MaxClaimsPerWallet: big.NewInt(100),
		RequiresWhitelist:  false,
		Active:             true,
	}, nil
}

// GetClaimCount implements SDKInterface
func (m *SimpleMockSDK) GetClaimCount(wallet common.Address, templateID string) (*big.Int, error) {
	m.Calls = append(m.Calls, "GetClaimCount")
	if m.ShouldFail {
		return nil, &MockError{Message: m.FailMessage}
	}
	return big.NewInt(0), nil
}

// IsWhitelisted implements SDKInterface
func (m *SimpleMockSDK) IsWhitelisted(wallet common.Address) (bool, error) {
	m.Calls = append(m.Calls, "IsWhitelisted")
	if m.ShouldFail {
		return false, &MockError{Message: m.FailMessage}
	}
	return true, nil
}

// GetRemainingDailyLimit implements SDKInterface
func (m *SimpleMockSDK) GetRemainingDailyLimit() (*big.Int, error) {
	m.Calls = append(m.Calls, "GetRemainingDailyLimit")
	if m.ShouldFail {
		return nil, &MockError{Message: m.FailMessage}
	}
	// 1000 tokens = 1000 * 10^18
	limit := new(big.Int)
	limit.SetString("1000000000000000000000", 10)
	return limit, nil
}

// MockError is a simple error type for testing
type MockError struct {
	Message string
}

func (e *MockError) Error() string {
	return e.Message
}
