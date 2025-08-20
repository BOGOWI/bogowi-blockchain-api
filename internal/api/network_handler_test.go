package api

import (
	"errors"
	"math/big"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/sdk"
	"bogowi-blockchain-go/internal/sdk/nft"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestMockSDK is a test-specific mock SDK for network handler tests
type TestMockSDK struct {
	mock.Mock
	CloseFunc func() error
}

func (m *TestMockSDK) GetTokenBalance(address string) (*sdk.TokenBalance, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.TokenBalance), args.Error(1)
}

func (m *TestMockSDK) GetGasPrice() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *TestMockSDK) TransferBOGOTokens(to string, amount string) (string, error) {
	args := m.Called(to, amount)
	return args.String(0), args.Error(1)
}

func (m *TestMockSDK) GetPublicKey() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *TestMockSDK) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	} else {
		m.Called()
	}
}

func (m *TestMockSDK) CheckRewardEligibility(templateID string, wallet common.Address) (bool, string, error) {
	args := m.Called(templateID, wallet)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *TestMockSDK) ClaimRewardV2(templateID string, recipient common.Address) (*types.Transaction, error) {
	args := m.Called(templateID, recipient)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *TestMockSDK) ClaimCustomReward(recipient common.Address, amount *big.Int, reason string) (*types.Transaction, error) {
	args := m.Called(recipient, amount, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *TestMockSDK) ClaimReferralBonus(referrer common.Address, referred common.Address) (*types.Transaction, error) {
	args := m.Called(referrer, referred)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *TestMockSDK) GetReferrer(wallet common.Address) (common.Address, error) {
	args := m.Called(wallet)
	return args.Get(0).(common.Address), args.Error(1)
}

func (m *TestMockSDK) GetRewardTemplate(templateID string) (*sdk.RewardTemplate, error) {
	args := m.Called(templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.RewardTemplate), args.Error(1)
}

func (m *TestMockSDK) GetClaimCount(wallet common.Address, templateID string) (*big.Int, error) {
	args := m.Called(wallet, templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *TestMockSDK) IsWhitelisted(wallet common.Address) (bool, error) {
	args := m.Called(wallet)
	return args.Bool(0), args.Error(1)
}

func (m *TestMockSDK) GetRemainingDailyLimit() (*big.Int, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func TestNewNetworkHandler(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful initialization with no networks configured",
			config: &config.Config{
				Testnet: config.NetworkConfig{
					Contracts: config.ContractAddresses{},
				},
				Mainnet: config.NetworkConfig{
					Contracts: config.ContractAddresses{},
				},
			},
			wantErr: false,
		},
		{
			name: "error when testnet configured but private key missing",
			config: &config.Config{
				Testnet: config.NetworkConfig{
					Contracts: config.ContractAddresses{
						BOGOToken: "0x123",
					},
				},
				TestnetPrivateKey: "",
			},
			wantErr: true,
			errMsg:  "TESTNET_PRIVATE_KEY is required for testnet operations",
		},
		{
			name: "error when mainnet configured but private key missing",
			config: &config.Config{
				Mainnet: config.NetworkConfig{
					Contracts: config.ContractAddresses{
						BOGOToken: "0x456",
					},
				},
				MainnetPrivateKey: "",
			},
			wantErr: true,
			errMsg:  "MAINNET_PRIVATE_KEY is required for mainnet operations",
		},
		{
			name: "error when mainnet configured with RewardDistributor but private key missing",
			config: &config.Config{
				Mainnet: config.NetworkConfig{
					Contracts: config.ContractAddresses{
						RewardDistributor: "0x789",
					},
				},
				MainnetPrivateKey: "",
			},
			wantErr: true,
			errMsg:  "MAINNET_PRIVATE_KEY is required for mainnet operations",
		},
		{
			name: "error when testnet configured with RewardDistributor but private key missing",
			config: &config.Config{
				Testnet: config.NetworkConfig{
					Contracts: config.ContractAddresses{
						RewardDistributor: "0xabc",
					},
				},
				TestnetPrivateKey: "",
			},
			wantErr: true,
			errMsg:  "TESTNET_PRIVATE_KEY is required for testnet operations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := NewNetworkHandler(tt.config)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, handler)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, handler)
				assert.Equal(t, tt.config, handler.config)
			}
		})
	}
}

func TestNetworkHandler_GetSDK(t *testing.T) {
	// Create a mock SDK for testing
	mockTestnetSDK := &TestMockSDK{}
	mockMainnetSDK := &TestMockSDK{}

	tests := []struct {
		name        string
		handler     *NetworkHandler
		network     string
		wantErr     bool
		errMsg      string
		expectedSDK SDKInterface
	}{
		{
			name: "get testnet SDK with 'testnet' parameter",
			handler: &NetworkHandler{
				testnetSDK: mockTestnetSDK,
				mainnetSDK: mockMainnetSDK,
			},
			network:     "testnet",
			wantErr:     false,
			expectedSDK: mockTestnetSDK,
		},
		{
			name: "get testnet SDK with 'columbus' parameter",
			handler: &NetworkHandler{
				testnetSDK: mockTestnetSDK,
				mainnetSDK: mockMainnetSDK,
			},
			network:     "columbus",
			wantErr:     false,
			expectedSDK: mockTestnetSDK,
		},
		{
			name: "get mainnet SDK with 'mainnet' parameter",
			handler: &NetworkHandler{
				testnetSDK: mockTestnetSDK,
				mainnetSDK: mockMainnetSDK,
			},
			network:     "mainnet",
			wantErr:     false,
			expectedSDK: mockMainnetSDK,
		},
		{
			name: "get mainnet SDK with 'camino' parameter",
			handler: &NetworkHandler{
				testnetSDK: mockTestnetSDK,
				mainnetSDK: mockMainnetSDK,
			},
			network:     "camino",
			wantErr:     false,
			expectedSDK: mockMainnetSDK,
		},
		{
			name: "error when testnet SDK not initialized",
			handler: &NetworkHandler{
				testnetSDK: nil,
				mainnetSDK: mockMainnetSDK,
			},
			network: "testnet",
			wantErr: true,
			errMsg:  "testnet SDK not initialized",
		},
		{
			name: "error when mainnet SDK not initialized",
			handler: &NetworkHandler{
				testnetSDK: mockTestnetSDK,
				mainnetSDK: nil,
			},
			network: "mainnet",
			wantErr: true,
			errMsg:  "mainnet SDK not initialized",
		},
		{
			name: "error with invalid network parameter",
			handler: &NetworkHandler{
				testnetSDK: mockTestnetSDK,
				mainnetSDK: mockMainnetSDK,
			},
			network: "invalid",
			wantErr: true,
			errMsg:  "invalid network: invalid (use 'testnet' or 'mainnet')",
		},
		{
			name: "error with empty network parameter",
			handler: &NetworkHandler{
				testnetSDK: mockTestnetSDK,
				mainnetSDK: mockMainnetSDK,
			},
			network: "",
			wantErr: true,
			errMsg:  "invalid network:  (use 'testnet' or 'mainnet')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk, err := tt.handler.GetSDK(tt.network)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, sdk)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSDK, sdk)
			}
		})
	}
}

func TestNetworkHandler_GetSDK_Concurrent(t *testing.T) {
	// Test concurrent access to GetSDK
	mockTestnetSDK := &TestMockSDK{}
	mockMainnetSDK := &TestMockSDK{}

	handler := &NetworkHandler{
		testnetSDK: mockTestnetSDK,
		mainnetSDK: mockMainnetSDK,
	}

	// Run multiple goroutines accessing GetSDK concurrently
	done := make(chan bool, 20)
	for i := 0; i < 10; i++ {
		go func() {
			sdk, err := handler.GetSDK("testnet")
			assert.NoError(t, err)
			assert.NotNil(t, sdk)
			done <- true
		}()

		go func() {
			sdk, err := handler.GetSDK("mainnet")
			assert.NoError(t, err)
			assert.NotNil(t, sdk)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestNetworkHandler_GetNFTSDK(t *testing.T) {
	// Create mock NFT SDKs for testing
	mockTestnetNFTSDK := &nft.Client{}
	mockMainnetNFTSDK := &nft.Client{}

	tests := []struct {
		name           string
		handler        *NetworkHandler
		network        string
		wantErr        bool
		errMsg         string
		expectedNFTSDK *nft.Client
	}{
		{
			name: "get testnet NFT SDK with 'testnet' parameter",
			handler: &NetworkHandler{
				testnetNFTSDK: mockTestnetNFTSDK,
				mainnetNFTSDK: mockMainnetNFTSDK,
			},
			network:        "testnet",
			wantErr:        false,
			expectedNFTSDK: mockTestnetNFTSDK,
		},
		{
			name: "get testnet NFT SDK with 'columbus' parameter",
			handler: &NetworkHandler{
				testnetNFTSDK: mockTestnetNFTSDK,
				mainnetNFTSDK: mockMainnetNFTSDK,
			},
			network:        "columbus",
			wantErr:        false,
			expectedNFTSDK: mockTestnetNFTSDK,
		},
		{
			name: "get mainnet NFT SDK with 'mainnet' parameter",
			handler: &NetworkHandler{
				testnetNFTSDK: mockTestnetNFTSDK,
				mainnetNFTSDK: mockMainnetNFTSDK,
			},
			network:        "mainnet",
			wantErr:        false,
			expectedNFTSDK: mockMainnetNFTSDK,
		},
		{
			name: "get mainnet NFT SDK with 'camino' parameter",
			handler: &NetworkHandler{
				testnetNFTSDK: mockTestnetNFTSDK,
				mainnetNFTSDK: mockMainnetNFTSDK,
			},
			network:        "camino",
			wantErr:        false,
			expectedNFTSDK: mockMainnetNFTSDK,
		},
		{
			name: "error when testnet NFT SDK not initialized",
			handler: &NetworkHandler{
				testnetNFTSDK: nil,
				mainnetNFTSDK: mockMainnetNFTSDK,
			},
			network: "testnet",
			wantErr: true,
			errMsg:  "testnet NFT SDK not initialized",
		},
		{
			name: "error when mainnet NFT SDK not initialized",
			handler: &NetworkHandler{
				testnetNFTSDK: mockTestnetNFTSDK,
				mainnetNFTSDK: nil,
			},
			network: "mainnet",
			wantErr: true,
			errMsg:  "mainnet NFT SDK not initialized",
		},
		{
			name: "error with invalid network parameter",
			handler: &NetworkHandler{
				testnetNFTSDK: mockTestnetNFTSDK,
				mainnetNFTSDK: mockMainnetNFTSDK,
			},
			network: "invalid",
			wantErr: true,
			errMsg:  "invalid network: invalid (use 'testnet' or 'mainnet')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nftSDK, err := tt.handler.GetNFTSDK(tt.network)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, nftSDK)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedNFTSDK, nftSDK)
			}
		})
	}
}

func TestNetworkHandler_GetNFTSDK_Concurrent(t *testing.T) {
	// Test concurrent access to GetNFTSDK
	mockTestnetNFTSDK := &nft.Client{}
	mockMainnetNFTSDK := &nft.Client{}

	handler := &NetworkHandler{
		testnetNFTSDK: mockTestnetNFTSDK,
		mainnetNFTSDK: mockMainnetNFTSDK,
	}

	// Run multiple goroutines accessing GetNFTSDK concurrently
	done := make(chan bool, 20)
	for i := 0; i < 10; i++ {
		go func() {
			nftSDK, err := handler.GetNFTSDK("testnet")
			assert.NoError(t, err)
			assert.NotNil(t, nftSDK)
			done <- true
		}()

		go func() {
			nftSDK, err := handler.GetNFTSDK("mainnet")
			assert.NoError(t, err)
			assert.NotNil(t, nftSDK)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestNetworkHandler_Close(t *testing.T) {
	tests := []struct {
		name               string
		handler            *NetworkHandler
		expectTestnetClose bool
		expectMainnetClose bool
		expectTestnetNFTClose bool
		expectMainnetNFTClose bool
	}{
		{
			name: "close all SDKs including NFT SDKs",
			handler: &NetworkHandler{
				testnetSDK: &TestMockSDK{},
				mainnetSDK: &TestMockSDK{},
				testnetNFTSDK: &nft.Client{},
				mainnetNFTSDK: &nft.Client{},
			},
			expectTestnetClose: true,
			expectMainnetClose: true,
			expectTestnetNFTClose: true,
			expectMainnetNFTClose: true,
		},
		{
			name: "close only testnet SDK",
			handler: &NetworkHandler{
				testnetSDK: &TestMockSDK{},
				mainnetSDK: nil,
			},
			expectTestnetClose: true,
			expectMainnetClose: false,
		},
		{
			name: "close only mainnet SDK",
			handler: &NetworkHandler{
				testnetSDK: nil,
				mainnetSDK: &TestMockSDK{},
			},
			expectTestnetClose: false,
			expectMainnetClose: true,
		},
		{
			name: "close with no SDKs initialized",
			handler: &NetworkHandler{
				testnetSDK: nil,
				mainnetSDK: nil,
			},
			expectTestnetClose: false,
			expectMainnetClose: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Track if Close was called on each SDK
			var testnetClosed, mainnetClosed bool

			if tt.handler.testnetSDK != nil {
				mockTestnet := tt.handler.testnetSDK.(*TestMockSDK)
				mockTestnet.CloseFunc = func() error {
					testnetClosed = true
					return nil
				}
			}

			if tt.handler.mainnetSDK != nil {
				mockMainnet := tt.handler.mainnetSDK.(*TestMockSDK)
				mockMainnet.CloseFunc = func() error {
					mainnetClosed = true
					return nil
				}
			}

			// Call Close
			tt.handler.Close()

			// Verify expectations
			assert.Equal(t, tt.expectTestnetClose, testnetClosed, "testnet SDK Close() called")
			assert.Equal(t, tt.expectMainnetClose, mainnetClosed, "mainnet SDK Close() called")
		})
	}
}

func TestNetworkHandler_Close_Concurrent(t *testing.T) {
	// Test that Close is thread-safe
	mockTestnetSDK := &TestMockSDK{}
	mockMainnetSDK := &TestMockSDK{}

	handler := &NetworkHandler{
		testnetSDK: mockTestnetSDK,
		mainnetSDK: mockMainnetSDK,
	}

	closeCount := 0
	mockTestnetSDK.CloseFunc = func() error {
		closeCount++
		return nil
	}
	mockMainnetSDK.CloseFunc = func() error {
		closeCount++
		return nil
	}

	// Run Close concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			handler.Close()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Close should be called multiple times safely
	// Each SDK's Close should be called at least once
	assert.GreaterOrEqual(t, closeCount, 2)
}

func TestNetworkHandler_Close_WithError(t *testing.T) {
	// Test Close when SDKs return errors
	mockTestnetSDK := &TestMockSDK{}
	mockMainnetSDK := &TestMockSDK{}

	handler := &NetworkHandler{
		testnetSDK: mockTestnetSDK,
		mainnetSDK: mockMainnetSDK,
	}

	testnetErr := errors.New("testnet close error")
	mainnetErr := errors.New("mainnet close error")

	testnetClosed := false
	mainnetClosed := false

	mockTestnetSDK.CloseFunc = func() error {
		testnetClosed = true
		return testnetErr
	}
	mockMainnetSDK.CloseFunc = func() error {
		mainnetClosed = true
		return mainnetErr
	}

	// Close should not panic even if SDKs return errors
	assert.NotPanics(t, func() {
		handler.Close()
	})

	// Both SDKs should still be closed despite errors
	assert.True(t, testnetClosed)
	assert.True(t, mainnetClosed)
}
