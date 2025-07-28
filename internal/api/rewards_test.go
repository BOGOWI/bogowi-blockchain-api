package api

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRewardTemplates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		BackendSecret: "test-secret",
	}
	
	handler := &Handler{
		SDK:    &MockSDK{},
		Config: cfg,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/rewards/templates", nil)

	handler.GetRewardTemplates(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	templates, ok := response["templates"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 10, len(templates))
}

func TestGetRewardTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		BackendSecret: "test-secret",
	}
	
	handler := &Handler{
		SDK:    &MockSDK{},
		Config: cfg,
	}

	tests := []struct {
		name       string
		templateID string
		wantStatus int
	}{
		{
			name:       "Valid template",
			templateID: "welcome_bonus",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid template",
			templateID: "invalid_template",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/api/rewards/templates/"+tt.templateID, nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.templateID}}

			handler.GetRewardTemplate(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestClaimRewardV2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name       string
		wallet     string
		request    ClaimRewardRequestV2
		setupMock  func(m *MockSDK)
		wantStatus int
	}{
		{
			name:   "Successful claim",
			wallet: "0x1234567890123456789012345678901234567890",
			request: ClaimRewardRequestV2{
				TemplateID: "welcome_bonus",
			},
			setupMock: func(m *MockSDK) {
				walletAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
				m.On("CheckRewardEligibility", "welcome_bonus", walletAddr).Return(true, "", nil)
				// Return a properly initialized transaction
				tx := types.NewTransaction(0, common.HexToAddress("0x0"), big.NewInt(0), 0, big.NewInt(0), nil)
				m.On("ClaimRewardV2", "welcome_bonus", walletAddr).Return(tx, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "Not eligible",
			wallet: "0x1234567890123456789012345678901234567890",
			request: ClaimRewardRequestV2{
				TemplateID: "welcome_bonus",
			},
			setupMock: func(m *MockSDK) {
				walletAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
				m.On("CheckRewardEligibility", "welcome_bonus", walletAddr).Return(false, "Already claimed", nil)
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSDK := new(MockSDK)
			if tt.setupMock != nil {
				tt.setupMock(mockSDK)
			}

			cfg := &config.Config{
				BackendSecret: "test-secret",
			}
			
			handler := &Handler{
				SDK:    mockSDK,
				Config: cfg,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			body, _ := json.Marshal(tt.request)
			c.Request, _ = http.NewRequest("POST", "/api/rewards/claim-v2", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("wallet", tt.wallet)

			handler.ClaimRewardV2(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			mockSDK.AssertExpectations(t)
		})
	}
}

func TestClaimCustomRewardV2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name       string
		authHeader string
		request    ClaimCustomRewardRequestV2
		setupMock  func(m *MockSDK)
		wantStatus int
	}{
		{
			name:       "Valid request",
			authHeader: "test-secret",
			request: ClaimCustomRewardRequestV2{
				Wallet: "0x1234567890123456789012345678901234567890",
				Amount: "500",
				Reason: "contest_winner",
			},
			setupMock: func(m *MockSDK) {
				walletAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
				amount := new(big.Int).Mul(big.NewInt(500), big.NewInt(1e18))
				tx := types.NewTransaction(0, common.HexToAddress("0x0"), big.NewInt(0), 0, big.NewInt(0), nil)
				m.On("ClaimCustomReward", walletAddr, amount, "contest_winner").Return(tx, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid auth",
			authHeader: "wrong-secret",
			request: ClaimCustomRewardRequestV2{
				Wallet: "0x1234567890123456789012345678901234567890",
				Amount: "500",
				Reason: "contest_winner",
			},
			setupMock: func(m *MockSDK) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Amount too high",
			authHeader: "test-secret",
			request: ClaimCustomRewardRequestV2{
				Wallet: "0x1234567890123456789012345678901234567890",
				Amount: "1500",
				Reason: "contest_winner",
			},
			setupMock: func(m *MockSDK) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSDK := new(MockSDK)
			if tt.setupMock != nil {
				tt.setupMock(mockSDK)
			}

			cfg := &config.Config{
				BackendSecret: "test-secret",
			}
			
			handler := &Handler{
				SDK:    mockSDK,
				Config: cfg,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			body, _ := json.Marshal(tt.request)
			c.Request, _ = http.NewRequest("POST", "/api/rewards/claim-custom", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Header.Set("X-Backend-Auth", tt.authHeader)

			handler.ClaimCustomRewardV2(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			mockSDK.AssertExpectations(t)
		})
	}
}

func TestCheckRewardEligibility(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		BackendSecret: "test-secret",
	}

	tests := []struct {
		name       string
		wallet     string
		templateID string
		setupMock  func(m *MockSDK)
		wantStatus int
	}{
		{
			name:       "Check specific template",
			wallet:     "0x1234567890123456789012345678901234567890",
			templateID: "welcome_bonus",
			setupMock: func(m *MockSDK) {
				walletAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
				m.On("CheckRewardEligibility", "welcome_bonus", walletAddr).Return(true, "", nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Check all templates",
			wallet:     "0x1234567890123456789012345678901234567890",
			templateID: "",
			setupMock: func(m *MockSDK) {
				walletAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
				templates := []string{
					"welcome_bonus", "founder_bonus", "first_nft_mint",
					"dao_participation", "attraction_tier_1", "attraction_tier_2",
					"attraction_tier_3", "attraction_tier_4",
				}
				for _, tmpl := range templates {
					m.On("CheckRewardEligibility", tmpl, walletAddr).Return(true, "", nil).Maybe()
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSDK := new(MockSDK)
			if tt.setupMock != nil {
				tt.setupMock(mockSDK)
			}

			handler := &Handler{
				SDK:    mockSDK,
				Config: cfg,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			url := "/api/rewards/eligibility"
			if tt.templateID != "" {
				url += "?templateId=" + tt.templateID
			}
			
			c.Request, _ = http.NewRequest("GET", url, nil)
			c.Set("wallet", tt.wallet)

			handler.CheckRewardEligibility(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			
			eligibilities, ok := response["eligibilities"].([]interface{})
			assert.True(t, ok)
			
			if tt.templateID != "" {
				assert.Equal(t, 1, len(eligibilities))
			} else {
				assert.True(t, len(eligibilities) > 1)
			}
			
			mockSDK.AssertExpectations(t)
		})
	}
}