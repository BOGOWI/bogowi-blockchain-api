package datakyte

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// BaseURL is the Datakyte API endpoint
	BaseURL = "https://dklnk.to"
	
	// API endpoints
	EndpointPing           = "/ping"
	EndpointValidateKey    = "/auth/validate-api-key"
	EndpointCreateNFT      = "/api/nfts"
	EndpointListNFTs       = "/api/nfts"
	EndpointGetNFT         = "/api/nfts/%s"
	EndpointUpdateMetadata = "/api/nfts/%s/metadata"
	EndpointGetMetadata    = "/api/nfts/%s/%s/metadata" // contractAddress/tokenId
	EndpointGetVersions    = "/api/nfts/%s/versions"
	EndpointRestoreVersion = "/api/nfts/%s/versions/%d/restore"
)

// Client represents a Datakyte API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Datakyte API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Response represents a standard API response
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *ErrorResponse  `json:"error,omitempty"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
}

// NFTMetadata represents the metadata structure for Datakyte NFTs
type NFTMetadata struct {
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	Image        string                   `json:"image"`
	ExternalURL  string                   `json:"external_url,omitempty"`
	AnimationURL string                   `json:"animation_url,omitempty"`
	Attributes   []NFTAttribute           `json:"attributes"`
	Properties   map[string]interface{}   `json:"properties,omitempty"`
}

// NFTAttribute represents a single NFT attribute
type NFTAttribute struct {
	TraitType   string      `json:"trait_type"`
	Value       interface{} `json:"value"`
	DisplayType string      `json:"display_type,omitempty"`
}

// CreateNFTRequest represents the request to create a new NFT
type CreateNFTRequest struct {
	TokenID         string                 `json:"tokenId"`
	ContractAddress string                 `json:"contractAddress"`
	ChainID         int                    `json:"chainId"`
	Metadata        NFTMetadata            `json:"metadata"`
	UserID          string                 `json:"userId,omitempty"`
	CustomData      map[string]interface{} `json:"customData,omitempty"`
}

// NFT represents a Datakyte NFT
type NFT struct {
	ID              string                 `json:"id"`
	TokenID         string                 `json:"tokenId"`
	ContractAddress string                 `json:"contractAddress"`
	ChainID         int                    `json:"chainId"`
	Metadata        NFTMetadata            `json:"metadata"`
	UserID          string                 `json:"userId"`
	CustomData      map[string]interface{} `json:"customData,omitempty"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	Deleted         bool                   `json:"deleted"`
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, endpoint string, body interface{}) (*Response, error) {
	url := c.baseURL + endpoint
	
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}
	
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var response Response
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	// Check for error response
	if !response.Success && response.Error != nil {
		return nil, fmt.Errorf("API error %d: %s - %s", 
			response.Error.StatusCode, 
			response.Error.Name, 
			response.Error.Message)
	}
	
	return &response, nil
}

// Ping checks if the API is available
func (c *Client) Ping() error {
	_, err := c.doRequest("GET", EndpointPing, nil)
	return err
}

// ValidateAPIKey validates the current API key
func (c *Client) ValidateAPIKey() (*User, error) {
	resp, err := c.doRequest("GET", EndpointValidateKey, nil)
	if err != nil {
		return nil, err
	}
	
	var user User
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}
	
	return &user, nil
}

// CreateNFT creates a new NFT with metadata
func (c *Client) CreateNFT(nft CreateNFTRequest) (*NFT, error) {
	resp, err := c.doRequest("POST", EndpointCreateNFT, nft)
	if err != nil {
		return nil, err
	}
	
	var createdNFT NFT
	if err := json.Unmarshal(resp.Data, &createdNFT); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NFT data: %w", err)
	}
	
	return &createdNFT, nil
}

// GetNFT retrieves an NFT by ID
func (c *Client) GetNFT(nftID string) (*NFT, error) {
	endpoint := fmt.Sprintf(EndpointGetNFT, nftID)
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	
	var nft NFT
	if err := json.Unmarshal(resp.Data, &nft); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NFT data: %w", err)
	}
	
	return &nft, nil
}

// UpdateMetadata updates the metadata for an existing NFT
func (c *Client) UpdateMetadata(nftID string, metadata NFTMetadata) (*NFT, error) {
	endpoint := fmt.Sprintf(EndpointUpdateMetadata, nftID)
	resp, err := c.doRequest("POST", endpoint, metadata)
	if err != nil {
		return nil, err
	}
	
	var updatedNFT NFT
	if err := json.Unmarshal(resp.Data, &updatedNFT); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NFT data: %w", err)
	}
	
	return &updatedNFT, nil
}

// GetMetadataByTokenID retrieves metadata by contract address and token ID (ERC-721 compliant)
func (c *Client) GetMetadataByTokenID(contractAddress, tokenID string) (*NFTMetadata, error) {
	endpoint := fmt.Sprintf(EndpointGetMetadata, contractAddress, tokenID)
	// This endpoint is public, so we don't need authentication
	url := c.baseURL + endpoint
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get metadata: status %d", resp.StatusCode)
	}
	
	var metadata NFTMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}
	
	return &metadata, nil
}

// ListNFTs lists the current user's NFTs
func (c *Client) ListNFTs() ([]NFT, error) {
	resp, err := c.doRequest("GET", EndpointListNFTs, nil)
	if err != nil {
		return nil, err
	}
	
	var nfts []NFT
	if err := json.Unmarshal(resp.Data, &nfts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NFTs data: %w", err)
	}
	
	return nfts, nil
}

// User represents a Datakyte user
type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
}