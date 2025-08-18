package sdk

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRewardBoundContract mocks BoundContract interface for reward tests
type MockRewardBoundContract struct {
	mock.Mock
}

func (m *MockRewardBoundContract) Call(opts *bind.CallOpts, results *[]interface{}, method string, params ...interface{}) error {
	args := m.Called(opts, results, method, params)
	return args.Error(0)
}

func (m *MockRewardBoundContract) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	args := m.Called(opts, method, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

// MockRewardEthClient mocks EthClient interface for reward tests
type MockRewardEthClient struct {
	mock.Mock
}

func (m *MockRewardEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockRewardEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockRewardEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockRewardEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockRewardEthClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	args := m.Called(ctx, call, blockNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockRewardEthClient) Close() {
	m.Called()
}

func TestClaimRewardV2(t *testing.T) {
	tests := []struct {
		name          string
		templateID    string
		recipient     common.Address
		setupMocks    func(*MockRewardBoundContract, *MockRewardEthClient)
		expectError   bool
		errorContains string
	}{
		{
			name:       "successful claim",
			templateID: "welcome_bonus",
			recipient:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				// Mock gas price suggestion
				mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)

				// Mock transaction
				expectedTx := types.NewTransaction(
					1,
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					big.NewInt(0),
					21000,
					big.NewInt(20000000000),
					nil,
				)
				mockContract.On("Transact", mock.Anything, "claimReward", []interface{}{"welcome_bonus"}).
					Return(expectedTx, nil)
			},
			expectError: false,
		},
		{
			name:       "reward distributor not initialized",
			templateID: "welcome_bonus",
			recipient:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				// No setup needed - we'll set rewardDistributor to nil
			},
			expectError:   true,
			errorContains: "reward distributor not initialized",
		},
		{
			name:       "gas price error",
			templateID: "welcome_bonus",
			recipient:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(nil, errors.New("network error"))
			},
			expectError:   true,
			errorContains: "failed to get gas price",
		},
		{
			name:       "transaction error",
			templateID: "welcome_bonus",
			recipient:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)
				mockContract.On("Transact", mock.Anything, "claimReward", []interface{}{"welcome_bonus"}).
					Return(nil, errors.New("insufficient funds"))
			},
			expectError:   true,
			errorContains: "failed to execute claimReward",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContract := new(MockRewardBoundContract)
			mockClient := new(MockRewardEthClient)

			sdk := &BOGOWISDK{
				client:  mockClient,
				chainID: big.NewInt(1),
			}

			// Set up rewardDistributor if not testing nil case
			if tt.name != "reward distributor not initialized" {
				sdk.rewardDistributor = &Contract{Instance: mockContract}
				// Generate a test private key
				privateKey, _ := crypto.GenerateKey()
				sdk.privateKey = privateKey
			}

			tt.setupMocks(mockContract, mockClient)

			tx, err := sdk.ClaimRewardV2(tt.templateID, tt.recipient)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, tx)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tx)
			}

			mockContract.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestClaimCustomReward(t *testing.T) {
	tests := []struct {
		name          string
		recipient     common.Address
		amount        *big.Int
		reason        string
		setupMocks    func(*MockRewardBoundContract, *MockRewardEthClient)
		expectError   bool
		errorContains string
	}{
		{
			name:      "successful custom claim",
			recipient: common.HexToAddress("0x1234567890123456789012345678901234567890"),
			amount:    new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), // 100 BOGO
			reason:    "Bug bounty reward",
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)

				expectedTx := types.NewTransaction(1, common.Address{}, big.NewInt(0), 21000, big.NewInt(20000000000), nil)
				mockContract.On("Transact", mock.Anything, "claimCustomReward",
					[]interface{}{
						common.HexToAddress("0x1234567890123456789012345678901234567890"),
						new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)),
						"Bug bounty reward",
					}).Return(expectedTx, nil)
			},
			expectError: false,
		},
		{
			name:          "amount exceeds maximum",
			recipient:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
			amount:        new(big.Int).Mul(big.NewInt(1001), big.NewInt(1e18)), // 1001 BOGO
			reason:        "Too much",
			setupMocks:    func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {},
			expectError:   true,
			errorContains: "amount exceeds maximum of 1000 BOGO",
		},
		{
			name:      "transaction failure",
			recipient: common.HexToAddress("0x1234567890123456789012345678901234567890"),
			amount:    new(big.Int).Mul(big.NewInt(50), big.NewInt(1e18)),
			reason:    "Test reward",
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)
				mockContract.On("Transact", mock.Anything, "claimCustomReward", mock.Anything).
					Return(nil, errors.New("unauthorized"))
			},
			expectError:   true,
			errorContains: "failed to execute claimCustomReward",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContract := new(MockRewardBoundContract)
			mockClient := new(MockRewardEthClient)

			sdk := &BOGOWISDK{
				client:            mockClient,
				chainID:           big.NewInt(1),
				rewardDistributor: &Contract{Instance: mockContract},
				privateKey:        func() *ecdsa.PrivateKey { key, _ := crypto.GenerateKey(); return key }(),
			}

			tt.setupMocks(mockContract, mockClient)

			tx, err := sdk.ClaimCustomReward(tt.recipient, tt.amount, tt.reason)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, tx)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tx)
			}

			mockContract.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestClaimReferralBonus(t *testing.T) {
	tests := []struct {
		name          string
		referrer      common.Address
		referred      common.Address
		setupMocks    func(*MockRewardBoundContract, *MockRewardEthClient)
		expectError   bool
		errorContains string
	}{
		{
			name:     "successful referral claim",
			referrer: common.HexToAddress("0x1111111111111111111111111111111111111111"),
			referred: common.HexToAddress("0x2222222222222222222222222222222222222222"),
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)

				expectedTx := types.NewTransaction(1, common.Address{}, big.NewInt(0), 21000, big.NewInt(20000000000), nil)
				mockContract.On("Transact", mock.Anything, "claimReferralBonus",
					[]interface{}{common.HexToAddress("0x1111111111111111111111111111111111111111")}).
					Return(expectedTx, nil)
			},
			expectError: false,
		},
		{
			name:     "referral bonus error",
			referrer: common.HexToAddress("0x3333333333333333333333333333333333333333"),
			referred: common.HexToAddress("0x4444444444444444444444444444444444444444"),
			setupMocks: func(mockContract *MockRewardBoundContract, mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)
				mockContract.On("Transact", mock.Anything, "claimReferralBonus", mock.Anything).
					Return(nil, errors.New("already claimed"))
			},
			expectError:   true,
			errorContains: "failed to execute claimReferralBonus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContract := new(MockRewardBoundContract)
			mockClient := new(MockRewardEthClient)

			sdk := &BOGOWISDK{
				client:            mockClient,
				chainID:           big.NewInt(1),
				rewardDistributor: &Contract{Instance: mockContract},
				privateKey:        func() *ecdsa.PrivateKey { key, _ := crypto.GenerateKey(); return key }(),
			}

			tt.setupMocks(mockContract, mockClient)

			tx, err := sdk.ClaimReferralBonus(tt.referrer, tt.referred)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, tx)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tx)
			}

			mockContract.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestCheckRewardEligibility(t *testing.T) {
	tests := []struct {
		name             string
		templateID       string
		wallet           common.Address
		expectedEligible bool
		expectedReason   string
		expectError      bool
	}{
		{
			name:             "welcome bonus eligible",
			templateID:       "welcome_bonus",
			wallet:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedEligible: true,
			expectedReason:   "",
			expectError:      false,
		},
		{
			name:             "founder bonus not whitelisted",
			templateID:       "founder_bonus",
			wallet:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedEligible: false,
			expectedReason:   "Not whitelisted",
			expectError:      false,
		},
		{
			name:             "dao participation cooldown",
			templateID:       "dao_participation",
			wallet:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedEligible: false,
			expectedReason:   "Cooldown period active",
			expectError:      false,
		},
		{
			name:             "unknown template defaults to eligible",
			templateID:       "unknown_template",
			wallet:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedEligible: true,
			expectedReason:   "",
			expectError:      false,
		},
		{
			name:             "reward distributor not initialized",
			templateID:       "any",
			wallet:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedEligible: false,
			expectedReason:   "reward distributor not initialized",
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &BOGOWISDK{}

			// Set up rewardDistributor if not testing nil case
			if tt.name != "reward distributor not initialized" {
				sdk.rewardDistributor = &Contract{Instance: new(MockRewardBoundContract)}
			}

			eligible, reason, err := sdk.CheckRewardEligibility(tt.templateID, tt.wallet)

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, tt.expectedEligible, eligible)
				assert.Equal(t, tt.expectedReason, reason)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedEligible, eligible)
				assert.Equal(t, tt.expectedReason, reason)
			}
		})
	}
}

func TestGetRewardTemplate(t *testing.T) {
	tests := []struct {
		name          string
		templateID    string
		expectError   bool
		errorContains string
		checkTemplate func(*testing.T, *RewardTemplate)
	}{
		{
			name:        "get welcome bonus template",
			templateID:  "welcome_bonus",
			expectError: false,
			checkTemplate: func(t *testing.T, tmpl *RewardTemplate) {
				assert.Equal(t, "welcome_bonus", tmpl.ID)
				assert.Equal(t, new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)), tmpl.FixedAmount)
				assert.Equal(t, big.NewInt(1), tmpl.MaxClaimsPerWallet)
				assert.True(t, tmpl.Active)
				assert.False(t, tmpl.RequiresWhitelist)
			},
		},
		{
			name:        "get founder bonus template",
			templateID:  "founder_bonus",
			expectError: false,
			checkTemplate: func(t *testing.T, tmpl *RewardTemplate) {
				assert.Equal(t, "founder_bonus", tmpl.ID)
				assert.Equal(t, new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), tmpl.FixedAmount)
				assert.Equal(t, big.NewInt(1), tmpl.MaxClaimsPerWallet)
				assert.True(t, tmpl.Active)
				assert.True(t, tmpl.RequiresWhitelist)
			},
		},
		{
			name:          "template not found",
			templateID:    "nonexistent",
			expectError:   true,
			errorContains: "template not found",
		},
		{
			name:          "reward distributor not initialized",
			templateID:    "any",
			expectError:   true,
			errorContains: "reward distributor not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &BOGOWISDK{}

			// Set up rewardDistributor if not testing nil case
			if tt.name != "reward distributor not initialized" {
				sdk.rewardDistributor = &Contract{Instance: new(MockRewardBoundContract)}
			}

			template, err := sdk.GetRewardTemplate(tt.templateID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, template)
			} else {
				require.NoError(t, err)
				require.NotNil(t, template)
				tt.checkTemplate(t, template)
			}
		})
	}
}

func TestGetClaimCount(t *testing.T) {
	wallet := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("successful get claim count", func(t *testing.T) {
		sdk := &BOGOWISDK{
			rewardDistributor: &Contract{Instance: new(MockRewardBoundContract)},
		}

		count, err := sdk.GetClaimCount(wallet, "welcome_bonus")
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(0), count)
	})

	t.Run("reward distributor not initialized", func(t *testing.T) {
		sdk := &BOGOWISDK{}

		count, err := sdk.GetClaimCount(wallet, "welcome_bonus")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reward distributor not initialized")
		assert.Nil(t, count)
	})
}

func TestIsWhitelisted(t *testing.T) {
	wallet := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("check whitelist status", func(t *testing.T) {
		sdk := &BOGOWISDK{
			rewardDistributor: &Contract{Instance: new(MockRewardBoundContract)},
		}

		whitelisted, err := sdk.IsWhitelisted(wallet)
		require.NoError(t, err)
		assert.False(t, whitelisted)
	})

	t.Run("reward distributor not initialized", func(t *testing.T) {
		sdk := &BOGOWISDK{}

		whitelisted, err := sdk.IsWhitelisted(wallet)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reward distributor not initialized")
		assert.False(t, whitelisted)
	})
}

func TestGetRemainingDailyLimit(t *testing.T) {
	t.Run("get remaining daily limit", func(t *testing.T) {
		sdk := &BOGOWISDK{
			rewardDistributor: &Contract{Instance: new(MockRewardBoundContract)},
		}

		limit, err := sdk.GetRemainingDailyLimit()
		require.NoError(t, err)
		expectedLimit := new(big.Int).Mul(big.NewInt(400000), big.NewInt(1e18))
		assert.Equal(t, expectedLimit, limit)
	})

	t.Run("reward distributor not initialized", func(t *testing.T) {
		sdk := &BOGOWISDK{}

		limit, err := sdk.GetRemainingDailyLimit()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reward distributor not initialized")
		assert.Nil(t, limit)
	})
}

func TestGetReferrer(t *testing.T) {
	wallet := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("get referrer address", func(t *testing.T) {
		sdk := &BOGOWISDK{
			rewardDistributor: &Contract{Instance: new(MockRewardBoundContract)},
		}

		referrer, err := sdk.GetReferrer(wallet)
		require.NoError(t, err)
		assert.Equal(t, common.Address{}, referrer)
	})

	t.Run("reward distributor not initialized", func(t *testing.T) {
		sdk := &BOGOWISDK{}

		referrer, err := sdk.GetReferrer(wallet)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reward distributor not initialized")
		assert.Equal(t, common.Address{}, referrer)
	})
}

func TestGetTransactOpts(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockRewardEthClient)
		expectError   bool
		errorContains string
	}{
		{
			name: "successful transact opts creation",
			setupMocks: func(mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)
			},
			expectError: false,
		},
		{
			name: "gas price error",
			setupMocks: func(mockClient *MockRewardEthClient) {
				mockClient.On("SuggestGasPrice", mock.Anything).Return(nil, errors.New("network error"))
			},
			expectError:   true,
			errorContains: "failed to get gas price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockRewardEthClient)

			sdk := &BOGOWISDK{
				client:     mockClient,
				chainID:    big.NewInt(1),
				privateKey: func() *ecdsa.PrivateKey { key, _ := crypto.GenerateKey(); return key }(),
			}

			tt.setupMocks(mockClient)

			opts, err := sdk.getTransactOpts()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, opts)
			} else {
				require.NoError(t, err)
				require.NotNil(t, opts)
				// Verify gas price was increased by 20%
				expectedGasPrice := new(big.Int).Mul(big.NewInt(20000000000), big.NewInt(120))
				expectedGasPrice = new(big.Int).Div(expectedGasPrice, big.NewInt(100))
				assert.Equal(t, expectedGasPrice, opts.GasPrice)
				assert.Equal(t, uint64(300000), opts.GasLimit)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
