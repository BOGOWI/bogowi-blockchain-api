package datakyte

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	assert.NotNil(t, client)
	assert.Equal(t, apiKey, client.apiKey)
	assert.Equal(t, BaseURL, client.baseURL)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 30*time.Second, client.httpClient.Timeout)
}

func TestClient_Ping(t *testing.T) {
	tests := []struct {
		name       string
		response   Response
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful ping",
			response: Response{
				Success: true,
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failed ping",
			response: Response{
				Success: false,
				Error: &ErrorResponse{
					StatusCode: 500,
					Name:       "InternalServerError",
					Message:    "Service unavailable",
				},
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, EndpointPing, r.URL.Path)
				assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL

			err := client.Ping()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_ValidateAPIKey(t *testing.T) {
	expectedUser := User{
		ID:          "user-123",
		Email:       "test@example.com",
		Name:        "Test User",
		Permissions: []string{"read", "write"},
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name       string
		response   Response
		statusCode int
		wantUser   *User
		wantErr    bool
	}{
		{
			name: "valid API key",
			response: Response{
				Success: true,
				Data:    mustMarshal(expectedUser),
			},
			statusCode: http.StatusOK,
			wantUser:   &expectedUser,
			wantErr:    false,
		},
		{
			name: "invalid API key",
			response: Response{
				Success: false,
				Error: &ErrorResponse{
					StatusCode: 401,
					Name:       "Unauthorized",
					Message:    "Invalid API key",
				},
			},
			statusCode: http.StatusUnauthorized,
			wantUser:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, EndpointValidateKey, r.URL.Path)

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL

			user, err := client.ValidateAPIKey()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Email, user.Email)
			}
		})
	}
}

func TestClient_CreateNFT(t *testing.T) {
	nftRequest := CreateNFTRequest{
		TokenID:         "1",
		ContractAddress: "0x123",
		ChainID:         501,
		Metadata: NFTMetadata{
			Name:        "Test NFT",
			Description: "Test Description",
			Image:       "https://example.com/image.png",
			Attributes: []NFTAttribute{
				{TraitType: "Color", Value: "Blue"},
			},
		},
	}

	expectedNFT := NFT{
		ID:              "nft-123",
		TokenID:         nftRequest.TokenID,
		ContractAddress: nftRequest.ContractAddress,
		ChainID:         nftRequest.ChainID,
		Metadata:        nftRequest.Metadata,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	tests := []struct {
		name       string
		response   Response
		statusCode int
		wantNFT    *NFT
		wantErr    bool
	}{
		{
			name: "successful creation",
			response: Response{
				Success: true,
				Data:    mustMarshal(expectedNFT),
			},
			statusCode: http.StatusCreated,
			wantNFT:    &expectedNFT,
			wantErr:    false,
		},
		{
			name: "duplicate NFT",
			response: Response{
				Success: false,
				Error: &ErrorResponse{
					StatusCode: 409,
					Name:       "Conflict",
					Message:    "NFT already exists",
				},
			},
			statusCode: http.StatusConflict,
			wantNFT:    nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, EndpointCreateNFT, r.URL.Path)

				var reqBody CreateNFTRequest
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				require.NoError(t, err)
				assert.Equal(t, nftRequest.TokenID, reqBody.TokenID)

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL

			nft, err := client.CreateNFT(nftRequest)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, nft)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nft)
				assert.Equal(t, tt.wantNFT.ID, nft.ID)
				assert.Equal(t, tt.wantNFT.TokenID, nft.TokenID)
			}
		})
	}
}

func TestClient_GetNFT(t *testing.T) {
	expectedNFT := NFT{
		ID:              "nft-123",
		TokenID:         "1",
		ContractAddress: "0x123",
		ChainID:         501,
		Metadata: NFTMetadata{
			Name:        "Test NFT",
			Description: "Test Description",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		nftID      string
		response   Response
		statusCode int
		wantNFT    *NFT
		wantErr    bool
	}{
		{
			name:  "existing NFT",
			nftID: "nft-123",
			response: Response{
				Success: true,
				Data:    mustMarshal(expectedNFT),
			},
			statusCode: http.StatusOK,
			wantNFT:    &expectedNFT,
			wantErr:    false,
		},
		{
			name:  "non-existent NFT",
			nftID: "nft-999",
			response: Response{
				Success: false,
				Error: &ErrorResponse{
					StatusCode: 404,
					Name:       "NotFound",
					Message:    "NFT not found",
				},
			},
			statusCode: http.StatusNotFound,
			wantNFT:    nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Contains(t, r.URL.Path, tt.nftID)

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL

			nft, err := client.GetNFT(tt.nftID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, nft)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nft)
				assert.Equal(t, tt.wantNFT.ID, nft.ID)
			}
		})
	}
}

func TestClient_UpdateMetadata(t *testing.T) {
	newMetadata := NFTMetadata{
		Name:        "Updated NFT",
		Description: "Updated Description",
		Image:       "https://example.com/new-image.png",
	}

	updatedNFT := NFT{
		ID:              "nft-123",
		TokenID:         "1",
		ContractAddress: "0x123",
		ChainID:         501,
		Metadata:        newMetadata,
		UpdatedAt:       time.Now(),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "nft-123")
		assert.Contains(t, r.URL.Path, "metadata")

		var reqBody NFTMetadata
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)
		assert.Equal(t, newMetadata.Name, reqBody.Name)

		response := Response{
			Success: true,
			Data:    mustMarshal(updatedNFT),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL

	nft, err := client.UpdateMetadata("nft-123", newMetadata)
	assert.NoError(t, err)
	assert.NotNil(t, nft)
	assert.Equal(t, updatedNFT.ID, nft.ID)
	assert.Equal(t, newMetadata.Name, nft.Metadata.Name)
}

func TestClient_GetMetadataByTokenID(t *testing.T) {
	expectedMetadata := NFTMetadata{
		Name:        "Test NFT",
		Description: "Test Description",
		Image:       "https://example.com/image.png",
		Attributes: []NFTAttribute{
			{TraitType: "Color", Value: "Blue"},
		},
	}

	tests := []struct {
		name            string
		contractAddress string
		tokenID         string
		statusCode      int
		wantMetadata    *NFTMetadata
		wantErr         bool
	}{
		{
			name:            "successful retrieval",
			contractAddress: "0x123",
			tokenID:         "1",
			statusCode:      http.StatusOK,
			wantMetadata:    &expectedMetadata,
			wantErr:         false,
		},
		{
			name:            "not found",
			contractAddress: "0x999",
			tokenID:         "999",
			statusCode:      http.StatusNotFound,
			wantMetadata:    nil,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Contains(t, r.URL.Path, tt.contractAddress)
				assert.Contains(t, r.URL.Path, tt.tokenID)

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(expectedMetadata)
				}
			}))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL

			metadata, err := client.GetMetadataByTokenID(tt.contractAddress, tt.tokenID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, metadata)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, metadata)
				assert.Equal(t, tt.wantMetadata.Name, metadata.Name)
			}
		})
	}
}

func TestClient_ListNFTs(t *testing.T) {
	expectedNFTs := []NFT{
		{
			ID:              "nft-1",
			TokenID:         "1",
			ContractAddress: "0x123",
			ChainID:         501,
		},
		{
			ID:              "nft-2",
			TokenID:         "2",
			ContractAddress: "0x123",
			ChainID:         501,
		},
	}

	tests := []struct {
		name       string
		response   Response
		statusCode int
		wantNFTs   []NFT
		wantErr    bool
	}{
		{
			name: "successful list",
			response: Response{
				Success: true,
				Data:    mustMarshal(expectedNFTs),
			},
			statusCode: http.StatusOK,
			wantNFTs:   expectedNFTs,
			wantErr:    false,
		},
		{
			name: "empty list",
			response: Response{
				Success: true,
				Data:    mustMarshal([]NFT{}),
			},
			statusCode: http.StatusOK,
			wantNFTs:   []NFT{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, EndpointListNFTs, r.URL.Path)

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL

			nfts, err := client.ListNFTs()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantNFTs), len(nfts))
			}
		})
	}
}

func TestClient_doRequest_NetworkError(t *testing.T) {
	client := NewClient("test-api-key")
	client.baseURL = "http://invalid-url-that-does-not-exist.local"

	_, err := client.doRequest("GET", "/test", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestClient_doRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL

	_, err := client.doRequest("GET", "/test", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
}

// Helper function to marshal data for tests
func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return json.RawMessage(data)
}
