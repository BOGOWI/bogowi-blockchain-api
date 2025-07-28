package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bogowi-blockchain-go/internal/sdk"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTokenRouter() (*gin.Engine, *MockSDK) {
	gin.SetMode(gin.TestMode)
	
	mockSDK := new(MockSDK)
	handler := &Handler{
		SDK: mockSDK,
	}
	
	router := gin.New()
	api := router.Group("/api")
	token := api.Group("/token")
	token.GET("/balance/:address", handler.GetTokenBalance)
	token.POST("/transfer", handler.TransferBOGOTokens)
	
	return router, mockSDK
}

func TestGetTokenBalance(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		mockBalance    *sdk.TokenBalance
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful balance retrieval",
			address: "0x742d35Cc6634C0532925a3b844Bc9e7595f8E97D",
			mockBalance: &sdk.TokenBalance{
				Address: "0x742d35Cc6634C0532925a3b844Bc9e7595f8E97D",
				Balance: "1000.5",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ethereum address",
			address:        "invalid-address",
			mockBalance:    nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid Ethereum address",
		},
		{
			name:           "sdk error",
			address:        "0x742d35Cc6634C0532925a3b844Bc9e7595f8E97D",
			mockBalance:    nil,
			mockError:      fmt.Errorf("connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "connection failed",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockSDK := setupTokenRouter()
			
			if tt.mockBalance != nil || tt.mockError != nil {
				mockSDK.On("GetTokenBalance", tt.address).Return(tt.mockBalance, tt.mockError)
			}
			
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/token/balance/"+tt.address, nil)
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedStatus == http.StatusOK {
				var response sdk.TokenBalance
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.mockBalance.Address, response.Address)
				assert.Equal(t, tt.mockBalance.Balance, response.Balance)
			} else {
				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}
			
			mockSDK.AssertExpectations(t)
		})
	}
}

func TestTransferBOGOTokens(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockTxHash     string
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful transfer",
			requestBody: TransferBOGOTokensRequest{
				To:     "0x742d35Cc6634C0532925a3b844Bc9e7595f8E97D",
				Amount: "100.5",
			},
			mockTxHash:     "0x1234567890abcdef",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid recipient address",
			requestBody: TransferBOGOTokensRequest{
				To:     "invalid-address",
				Amount: "100.5",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid recipient address",
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"to": "0x742d35Cc6634C0532925a3b844Bc9e7595f8E97D",
				// amount is missing
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'TransferBOGOTokensRequest.Amount' Error:Field validation for 'Amount' failed on the 'required' tag",
		},
		{
			name: "sdk transfer error",
			requestBody: TransferBOGOTokensRequest{
				To:     "0x742d35Cc6634C0532925a3b844Bc9e7595f8E97D",
				Amount: "100.5",
			},
			mockError:      fmt.Errorf("insufficient funds"),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "insufficient funds",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockSDK := setupTokenRouter()
			
			// Setup mock only for valid address cases that will reach SDK
			if req, ok := tt.requestBody.(TransferBOGOTokensRequest); ok && len(req.To) > 0 && req.To[0:2] == "0x" {
				mockSDK.On("TransferBOGOTokens", req.To, req.Amount).Return(tt.mockTxHash, tt.mockError)
			}
			
			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/token/transfer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "Transfer initiated successfully", response["message"])
				assert.Equal(t, tt.mockTxHash, response["transaction"])
			} else {
				assert.Equal(t, tt.expectedError, response["error"])
			}
			
			mockSDK.AssertExpectations(t)
		})
	}
}