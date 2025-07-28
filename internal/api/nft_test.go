package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNFTBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		address    string
		tokenId    string
		mockError  error
		statusCode int
		errorMsg   string
	}{
		{
			name:       "successful NFT balance query",
			address:    "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			tokenId:    "1",
			mockError:  nil,
			statusCode: http.StatusOK,
		},
		{
			name:       "invalid address",
			address:    "invalid",
			tokenId:    "1",
			mockError:  nil,
			statusCode: http.StatusBadRequest,
			errorMsg:   "Invalid Ethereum address",
		},
		{
			name:       "SDK error",
			address:    "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			tokenId:    "1",
			mockError:  errors.New("network error"),
			statusCode: http.StatusInternalServerError,
			errorMsg:   "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSDK := new(MockSDK)
			handler := &Handler{SDK: mockSDK}

			router := gin.New()
			router.GET("/api/nft/balance/:address/:tokenId", handler.GetNFTBalance)

			if tt.mockError == nil && tt.address != "invalid" {
				mockSDK.On("GetNFTBalance", tt.address, tt.tokenId).Return("1", nil)
			} else if tt.mockError != nil {
				mockSDK.On("GetNFTBalance", tt.address, tt.tokenId).Return("", tt.mockError)
			}

			req, _ := http.NewRequest("GET", "/api/nft/balance/"+tt.address+"/"+tt.tokenId, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if tt.errorMsg != "" {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.errorMsg)
			}

			mockSDK.AssertExpectations(t)
		})
	}
}

func TestMintEventTicket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    map[string]interface{}
		mockError  error
		mockTxHash string
		statusCode int
		errorMsg   string
	}{
		{
			name: "successful ticket mint",
			request: map[string]interface{}{
				"to":        "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
				"eventName": "Conservation Summit 2024",
				"eventDate": "2024-12-01",
			},
			mockError:  nil,
			mockTxHash: "0x1234567890abcdef",
			statusCode: http.StatusOK,
		},
		{
			name: "missing required fields",
			request: map[string]interface{}{
				"to": "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
			},
			statusCode: http.StatusBadRequest,
			errorMsg:   "Missing required fields",
		},
		{
			name: "SDK error",
			request: map[string]interface{}{
				"to":        "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
				"eventName": "Conservation Summit 2024",
				"eventDate": "2024-12-01",
			},
			mockError:  errors.New("insufficient funds"),
			statusCode: http.StatusInternalServerError,
			errorMsg:   "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSDK := new(MockSDK)
			handler := &Handler{SDK: mockSDK}

			router := gin.New()
			router.POST("/api/nft/mint-ticket", handler.MintEventTicket)

			body, _ := json.Marshal(tt.request)

			if tt.mockError == nil && tt.request["eventName"] != nil && tt.request["eventDate"] != nil {
				mockSDK.On("MintEventTicket",
					tt.request["to"].(string),
					tt.request["eventName"].(string),
					tt.request["eventDate"].(string),
				).Return(tt.mockTxHash, nil)
			} else if tt.mockError != nil {
				mockSDK.On("MintEventTicket",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return("", tt.mockError)
			}

			req, _ := http.NewRequest("POST", "/api/nft/mint-ticket", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if tt.errorMsg != "" {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.errorMsg)
			} else if tt.statusCode == http.StatusOK {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.mockTxHash, response["transactionHash"])
			}

			mockSDK.AssertExpectations(t)
		})
	}
}

func TestMintConservationNFT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    map[string]interface{}
		mockError  error
		mockTxHash string
		statusCode int
		errorMsg   string
	}{
		{
			name: "successful NFT mint",
			request: map[string]interface{}{
				"to":          "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
				"tokenURI":    "https://ipfs.io/ipfs/QmXx...",
				"description": "Rare Bogotá Wildlife NFT",
			},
			mockError:  nil,
			mockTxHash: "0xabcdef1234567890",
			statusCode: http.StatusOK,
		},
		{
			name: "missing tokenURI",
			request: map[string]interface{}{
				"to":          "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
				"description": "Rare Bogotá Wildlife NFT",
			},
			statusCode: http.StatusBadRequest,
			errorMsg:   "Missing required fields",
		},
		{
			name: "SDK error",
			request: map[string]interface{}{
				"to":          "0x742d35Cc6634C0532925a3b844Bc9e7595f8f8E2",
				"tokenURI":    "https://ipfs.io/ipfs/QmXx...",
				"description": "Rare Bogotá Wildlife NFT",
			},
			mockError:  errors.New("contract paused"),
			statusCode: http.StatusInternalServerError,
			errorMsg:   "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSDK := new(MockSDK)
			handler := &Handler{SDK: mockSDK}

			router := gin.New()
			router.POST("/api/nft/mint-collectible", handler.MintConservationNFT)

			body, _ := json.Marshal(tt.request)

			if tt.mockError == nil && tt.request["tokenURI"] != nil {
				description := ""
				if tt.request["description"] != nil {
					description = tt.request["description"].(string)
				}
				mockSDK.On("MintConservationNFT",
					tt.request["to"].(string),
					tt.request["tokenURI"].(string),
					description,
				).Return(tt.mockTxHash, nil)
			} else if tt.mockError != nil {
				mockSDK.On("MintConservationNFT",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return("", tt.mockError)
			}

			req, _ := http.NewRequest("POST", "/api/nft/mint-collectible", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if tt.errorMsg != "" {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.errorMsg)
			} else if tt.statusCode == http.StatusOK {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.mockTxHash, response["transactionHash"])
			}

			mockSDK.AssertExpectations(t)
		})
	}
}
