package sdk

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockEthClient is a mock Ethereum client
type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockEthClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	args := m.Called(ctx, call, blockNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEthClient) Close() {
	m.Called()
}

// MockBoundContract is a mock bound contract
type MockBoundContract struct {
	mock.Mock
}

func (m *MockBoundContract) Call(opts *bind.CallOpts, results *[]interface{}, method string, params ...interface{}) error {
	args := m.Called(opts, results, method, params)
	// Set the balance result if provided
	if method == "balanceOf" && len(*results) > 0 {
		if args.Get(0) != nil {
			// Set the value through the pointer
			balancePtr := (*results)[0].(**big.Int)
			*balancePtr = args.Get(0).(*big.Int)
		}
	}
	return args.Error(1)
}

func (m *MockBoundContract) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	args := m.Called(opts, method, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func TestNewBOGOWISDK(t *testing.T) {
	// Test creating new SDK instance
	cfg := &config.Config{
		RPCUrl:     "https://columbus.camino.network/ext/bc/C/rpc",
		ChainID:    501,
		PrivateKey: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Contracts:  config.ContractAddresses{},
	}

	sdk, err := NewBOGOWISDK(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, sdk)
}

func TestSDKWithInvalidPrivateKey(t *testing.T) {
	// Test creating SDK with invalid private key
	cfg := &config.Config{
		RPCUrl:     "https://columbus.camino.network/ext/bc/C/rpc",
		ChainID:    501,
		PrivateKey: "invalid",
		Contracts:  config.ContractAddresses{},
	}

	sdk, err := NewBOGOWISDK(cfg)
	assert.Error(t, err)
	assert.Nil(t, sdk)
}

func TestGetTokenBalance(t *testing.T) {
	cfg := &config.Config{
		PrivateKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		RPCUrl:     "http://localhost:8545",
		ChainID:    1,
		Contracts: config.ContractAddresses{
			BOGOTokenV2: "0x1234567890123456789012345678901234567890",
		},
	}

	sdk := &BOGOWISDK{
		config: cfg,
		contracts: &ContractInstances{
			BOGOTokenV2: &Contract{
				Instance: nil, // Will be set in test
			},
		},
	}

	tests := []struct {
		name        string
		address     string
		mockBalance *big.Int
		mockError   error
		wantBalance string
		wantError   bool
	}{
		{
			name:        "successful balance query",
			address:     "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			mockBalance: big.NewInt(1000000000000000000), // 1 token
			mockError:   nil,
			wantBalance: "1",
			wantError:   false,
		},
		{
			name:        "balance query error",
			address:     "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			mockBalance: nil,
			mockError:   errors.New("connection error"),
			wantBalance: "",
			wantError:   true,
		},
		{
			name:        "zero balance",
			address:     "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			mockBalance: big.NewInt(0),
			mockError:   nil,
			wantBalance: "0",
			wantError:   false,
		},
		{
			name:        "large balance",
			address:     "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			mockBalance: new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1000000000000000000)), // 1M tokens
			mockError:   nil,
			wantBalance: "1000000",
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContract := new(MockBoundContract)
			sdk.contracts.BOGOTokenV2.Instance = mockContract

			if tt.mockError != nil {
				mockContract.On("Call", mock.Anything, mock.Anything, "balanceOf", mock.Anything).
					Return(nil, tt.mockError)
			} else {
				mockContract.On("Call", mock.Anything, mock.Anything, "balanceOf", mock.Anything).
					Return(tt.mockBalance, nil)
			}

			balance, err := sdk.GetTokenBalance(tt.address)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, balance)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, balance)
				assert.Equal(t, tt.address, balance.Address)
				assert.Equal(t, tt.wantBalance, balance.Balance)
			}

			mockContract.AssertExpectations(t)
		})
	}
}

func TestGetTokenBalance_NoContract(t *testing.T) {
	sdk := &BOGOWISDK{
		contracts: &ContractInstances{},
	}

	balance, err := sdk.GetTokenBalance("0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BOGOTokenV2 contract not initialized")
	assert.Nil(t, balance)
}

func TestGetGasPrice(t *testing.T) {
	mockClient := new(MockEthClient)
	sdk := &BOGOWISDK{
		client: mockClient,
	}

	tests := []struct {
		name          string
		mockGasPrice  *big.Int
		mockError     error
		expectedPrice string
		wantError     bool
	}{
		{
			name:          "normal gas price",
			mockGasPrice:  big.NewInt(25000000000), // 25 Gwei
			mockError:     nil,
			expectedPrice: "25.00 gwei",
			wantError:     false,
		},
		{
			name:          "high gas price",
			mockGasPrice:  big.NewInt(150000000000), // 150 Gwei
			mockError:     nil,
			expectedPrice: "150.00 gwei",
			wantError:     false,
		},
		{
			name:          "low gas price",
			mockGasPrice:  big.NewInt(1000000000), // 1 Gwei
			mockError:     nil,
			expectedPrice: "1.00 gwei",
			wantError:     false,
		},
		{
			name:          "gas price error",
			mockGasPrice:  nil,
			mockError:     errors.New("network error"),
			expectedPrice: "",
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("SuggestGasPrice", mock.Anything).
				Return(tt.mockGasPrice, tt.mockError).Once()

			price, err := sdk.GetGasPrice()

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, price)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPrice, price)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestTransferBOGOTokens(t *testing.T) {
	mockClient := new(MockEthClient)
	mockContract := new(MockBoundContract)
	
	sdk := &BOGOWISDK{
		client: mockClient,
		auth:   &bind.TransactOpts{From: common.HexToAddress("0x1234567890123456789012345678901234567890")},
		contracts: &ContractInstances{
			BOGOTokenV2: &Contract{
				Instance: mockContract,
			},
		},
	}

	tests := []struct {
		name         string
		to           string
		amount       string
		mockNonce    uint64
		mockGasPrice *big.Int
		mockTx       *types.Transaction
		mockError    error
		wantTxHash   string
		wantError    bool
		errorMsg     string
	}{
		{
			name:         "successful transfer",
			to:           "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			amount:       "100.5",
			mockNonce:    5,
			mockGasPrice: big.NewInt(25000000000),
			mockTx:       types.NewTransaction(5, common.Address{}, big.NewInt(0), 100000, big.NewInt(25000000000), nil),
			mockError:    nil,
			wantTxHash:   "0x",
			wantError:    false,
		},
		{
			name:         "invalid recipient address",
			to:           "invalid-address",
			amount:       "100",
			wantError:    true,
			errorMsg:     "invalid recipient address",
		},
		{
			name:         "invalid amount format",
			to:           "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			amount:       "abc",
			wantError:    true,
			errorMsg:     "invalid amount format",
		},
		{
			name:         "transaction error",
			to:           "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			amount:       "100",
			mockNonce:    5,
			mockGasPrice: big.NewInt(25000000000),
			mockError:    errors.New("insufficient funds"),
			wantError:    true,
			errorMsg:     "failed to execute transfer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockClient.ExpectedCalls = nil
			mockContract.ExpectedCalls = nil

			if tt.to != "invalid-address" && tt.amount != "abc" {
				mockClient.On("PendingNonceAt", mock.Anything, sdk.auth.From).
					Return(tt.mockNonce, nil).Once()
				mockClient.On("SuggestGasPrice", mock.Anything).
					Return(tt.mockGasPrice, nil).Once()
				
				if tt.mockError != nil {
					mockContract.On("Transact", mock.Anything, "transfer", mock.Anything, mock.Anything).
						Return(nil, tt.mockError).Once()
				} else {
					mockContract.On("Transact", mock.Anything, "transfer", mock.Anything, mock.Anything).
						Return(tt.mockTx, nil).Once()
				}
			}

			txHash, err := sdk.TransferBOGOTokens(tt.to, tt.amount)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, txHash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, txHash)
			}

			mockClient.AssertExpectations(t)
			mockContract.AssertExpectations(t)
		})
	}
}

func TestTransferBOGOTokens_NoContract(t *testing.T) {
	sdk := &BOGOWISDK{
		contracts: &ContractInstances{},
	}

	txHash, err := sdk.TransferBOGOTokens("0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2", "100")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BOGO token contract not initialized")
	assert.Empty(t, txHash)
}

func TestGetPublicKey(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		wantError  bool
	}{
		{
			name:       "valid private key",
			privateKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantError:  false,
		},
		{
			name:       "valid private key with 0x prefix",
			privateKey: "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantError:  false,
		},
		{
			name:       "invalid private key",
			privateKey: "invalid",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &BOGOWISDK{
				config: &config.Config{
					PrivateKey: tt.privateKey,
				},
			}

			pubKey, err := sdk.GetPublicKey()

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, pubKey)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, pubKey)
				assert.True(t, common.IsHexAddress(pubKey))
			}
		})
	}
}

func TestClose(t *testing.T) {
	mockClient := new(MockEthClient)
	sdk := &BOGOWISDK{
		client: mockClient,
	}

	mockClient.On("Close").Once()
	
	sdk.Close()
	
	mockClient.AssertExpectations(t)
}

func TestClose_NilClient(t *testing.T) {
	sdk := &BOGOWISDK{
		client: nil,
	}
	
	// Should not panic
	assert.NotPanics(t, func() {
		sdk.Close()
	})
}

func TestInitializeContract(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		abiJSON   string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty address",
			address:   "",
			abiJSON:   ERC20ABI,
			wantError: true,
			errorMsg:  "contract address is empty",
		},
		{
			name:      "invalid ABI",
			address:   "0x1234567890123456789012345678901234567890",
			abiJSON:   "invalid json",
			wantError: true,
			errorMsg:  "failed to parse ABI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &BOGOWISDK{
				client: new(MockEthClient),
			}
			
			contract, err := sdk.initializeContract(tt.address, tt.abiJSON)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, contract)
			}
		})
	}
}
