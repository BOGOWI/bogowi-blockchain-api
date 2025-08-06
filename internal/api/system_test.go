package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSDK is a mock implementation of the SDK
type MockSDK struct {
	mock.Mock
}

func (m *MockSDK) GetTokenBalance(address string) (*sdk.TokenBalance, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.TokenBalance), args.Error(1)
}

func (m *MockSDK) GetGasPrice() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockSDK) TransferBOGOTokens(to string, amount string) (string, error) {
	args := m.Called(to, amount)
	return args.String(0), args.Error(1)
}

func (m *MockSDK) GetPublicKey() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockSDK) Close() {
	m.Called()
}

func (m *MockSDK) GetNFTBalance(address string, tokenId string) (string, error) {
	args := m.Called(address, tokenId)
	return args.String(0), args.Error(1)
}

func (m *MockSDK) MintEventTicket(to string, eventName string, eventDate string) (string, error) {
	args := m.Called(to, eventName, eventDate)
	return args.String(0), args.Error(1)
}

func (m *MockSDK) MintConservationNFT(to string, tokenURI string, description string) (string, error) {
	args := m.Called(to, tokenURI, description)
	return args.String(0), args.Error(1)
}

func (m *MockSDK) GetRewardInfo(address string) (map[string]interface{}, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockSDK) GetAchievementProgress(address string, achievementId string) (map[string]interface{}, error) {
	args := m.Called(address, achievementId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockSDK) ClaimReward(address string, rewardType string, rewardAmount string) (string, error) {
	args := m.Called(address, rewardType, rewardAmount)
	return args.String(0), args.Error(1)
}

// New reward system methods
func (m *MockSDK) CheckRewardEligibility(templateID string, wallet common.Address) (bool, string, error) {
	args := m.Called(templateID, wallet)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockSDK) ClaimRewardV2(templateID string, recipient common.Address) (*types.Transaction, error) {
	args := m.Called(templateID, recipient)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockSDK) ClaimCustomReward(recipient common.Address, amount *big.Int, reason string) (*types.Transaction, error) {
	args := m.Called(recipient, amount, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockSDK) ClaimReferralBonus(referrer common.Address, referred common.Address) (*types.Transaction, error) {
	args := m.Called(referrer, referred)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockSDK) GetReferrer(wallet common.Address) (common.Address, error) {
	args := m.Called(wallet)
	return args.Get(0).(common.Address), args.Error(1)
}

func (m *MockSDK) GetRewardTemplate(templateID string) (*sdk.RewardTemplate, error) {
	args := m.Called(templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.RewardTemplate), args.Error(1)
}

func (m *MockSDK) GetClaimCount(wallet common.Address, templateID string) (*big.Int, error) {
	args := m.Called(wallet, templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockSDK) IsWhitelisted(wallet common.Address) (bool, error) {
	args := m.Called(wallet)
	return args.Bool(0), args.Error(1)
}

func (m *MockSDK) GetRemainingDailyLimit() (*big.Int, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockSDK, *config.Config) {
	gin.SetMode(gin.TestMode)

	mockSDK := new(MockSDK)
	cfg := &config.Config{
		Environment:      "test",
		APIPort:          "3001",
		BackendSecret:    "test-secret",
		DevBackendSecret: "test-dev-secret",
		Testnet: config.NetworkConfig{
			RPCUrl:  "https://columbus.camino.network/ext/bc/C/rpc",
			ChainID: 501,
			Contracts: config.ContractAddresses{
				BOGOToken:         "0x123",
				ConservationNFT:   "0x456",
				CommercialNFT:     "0x789",
				RewardDistributor: "0xabc",
				MultisigTreasury:  "0xdef",
			},
		},
		Mainnet: config.NetworkConfig{
			RPCUrl:  "https://api.camino.network/ext/bc/C/rpc",
			ChainID: 500,
			Contracts: config.ContractAddresses{
				BOGOToken:         "0x123",
				ConservationNFT:   "0x456",
				CommercialNFT:     "0x789",
				RewardDistributor: "0xabc",
				MultisigTreasury:  "0xdef",
			},
		},
	}

	// Create a mock NetworkHandler
	networkHandler := &NetworkHandler{
		testnetSDK: mockSDK,
		mainnetSDK: mockSDK,
		config:     cfg,
	}

	handler := &Handler{
		SDK:            mockSDK,
		NetworkHandler: networkHandler,
		Config:         cfg,
	}

	router := gin.New()
	api := router.Group("/api")
	api.GET("/health", handler.GetHealth)
	api.GET("/gas-price", handler.GetGasPrice)

	return router, mockSDK, cfg
}

func TestGetHealth(t *testing.T) {
	router, _, cfg := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])

	contracts, ok := response["contracts"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, cfg.Testnet.Contracts.BOGOToken, contracts["bogo_token"])
}

func TestGetGasPrice(t *testing.T) {
	tests := []struct {
		name           string
		mockGasPrice   string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful gas price retrieval",
			mockGasPrice:   "25.00 gwei",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"gasPrice": "25.00 gwei",
			},
		},
		{
			name:           "gas price retrieval error",
			mockGasPrice:   "",
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to get gas price",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockSDK, _ := setupTestRouter()

			mockSDK.On("GetGasPrice").Return(tt.mockGasPrice, tt.mockError)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/gas-price", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody["gasPrice"], response["gasPrice"])
			} else {
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}

			mockSDK.AssertExpectations(t)
		})
	}
}
