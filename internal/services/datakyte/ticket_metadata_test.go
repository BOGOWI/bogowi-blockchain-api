package datakyte

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTicketMetadataService(t *testing.T) {
	apiKey := "test-api-key"
	contractAddress := "0x123"
	chainID := 501

	service := NewTicketMetadataService(apiKey, contractAddress, chainID)

	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
	assert.Equal(t, contractAddress, service.contractAddress)
	assert.Equal(t, chainID, service.chainID)
	assert.Equal(t, "https://storage.bogowi.com/tickets", service.baseImageURL)
}

func TestTicketMetadataService_CreateTicketMetadata(t *testing.T) {
	validUntil := time.Now().Add(30 * 24 * time.Hour)
	transferableAfter := time.Now().Add(7 * 24 * time.Hour)
	expiresAt := time.Now().Add(60 * 24 * time.Hour)

	ticketData := BOGOWITicketData{
		TokenID:            1,
		BookingID:          "BOOK-123",
		EventID:            "EVENT-456",
		ExperienceTitle:    "Rainforest Trek",
		ExperienceType:     "Nature Tour",
		Location:           "Costa Rica",
		Duration:           "4 hours",
		MaxParticipants:    10,
		CarbonOffset:       50,
		ConservationImpact: "Supporting local wildlife conservation",
		ValidUntil:         validUntil,
		TransferableAfter:  transferableAfter,
		ExpiresAt:          expiresAt,
		BOGORewards:        500, // 5%
		RecipientAddress:   "0xRecipient",
		RecipientName:      "John Doe",
		ProviderName:       "EcoTours CR",
		ProviderContact:    "contact@ecotours.cr",
	}

	expectedNFT := NFT{
		ID:              "nft-123",
		TokenID:         "1",
		ContractAddress: "0x123",
		ChainID:         501,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/nfts", r.URL.Path)

		var reqBody CreateNFTRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify request structure
		assert.Equal(t, "1", reqBody.TokenID)
		assert.Equal(t, "0x123", reqBody.ContractAddress)
		assert.Equal(t, 501, reqBody.ChainID)
		assert.Equal(t, "BOGOWI Eco-Experience #1", reqBody.Metadata.Name)
		assert.Contains(t, reqBody.Metadata.Description, "Nature Tour")
		assert.Contains(t, reqBody.Metadata.Description, "Costa Rica")

		// Check attributes
		assert.Greater(t, len(reqBody.Metadata.Attributes), 0)
		
		// Check custom data
		assert.Equal(t, "BOOK-123", reqBody.CustomData["bookingId"])
		assert.Equal(t, "EVENT-456", reqBody.CustomData["eventId"])

		// Return success response
		response := Response{
			Success: true,
			Data:    mustMarshal(expectedNFT),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service := NewTicketMetadataService("test-api-key", "0x123", 501)
	service.client.baseURL = server.URL

	nft, err := service.CreateTicketMetadata(ticketData)
	assert.NoError(t, err)
	assert.NotNil(t, nft)
	assert.Equal(t, expectedNFT.ID, nft.ID)
}

func TestTicketMetadataService_generateMetadata(t *testing.T) {
	service := NewTicketMetadataService("test-api-key", "0x123", 501)

	validUntil := time.Now().Add(30 * 24 * time.Hour)
	transferableAfter := time.Now().Add(7 * 24 * time.Hour)
	expiresAt := time.Now().Add(60 * 24 * time.Hour)

	tests := []struct {
		name     string
		data     BOGOWITicketData
		validate func(t *testing.T, metadata NFTMetadata)
	}{
		{
			name: "basic ticket with default image",
			data: BOGOWITicketData{
				TokenID:            1,
				BookingID:          "BOOK-123",
				EventID:            "EVENT-456",
				ExperienceTitle:    "Rainforest Trek",
				ExperienceType:     "Nature Tour",
				Location:           "Costa Rica",
				Duration:           "4 hours",
				MaxParticipants:    10,
				CarbonOffset:       50,
				ConservationImpact: "Supporting local wildlife conservation",
				ValidUntil:         validUntil,
				TransferableAfter:  transferableAfter,
				ExpiresAt:          expiresAt,
				BOGORewards:        500, // 5%
				RecipientAddress:   "0xRecipient",
				ProviderName:       "EcoTours CR",
			},
			validate: func(t *testing.T, metadata NFTMetadata) {
				assert.Equal(t, "BOGOWI Eco-Experience #1", metadata.Name)
				assert.Contains(t, metadata.Description, "Nature Tour")
				assert.Equal(t, "https://storage.bogowi.com/tickets/1.png", metadata.Image)
				assert.Equal(t, "https://bogowi.com/experiences/1", metadata.ExternalURL)
				
				// Check attributes
				found := false
				for _, attr := range metadata.Attributes {
					if attr.TraitType == "Experience Type" {
						assert.Equal(t, "Nature Tour", attr.Value)
						found = true
						break
					}
				}
				assert.True(t, found, "Experience Type attribute not found")
				
				// Check BOGO rewards conversion
				found = false
				for _, attr := range metadata.Attributes {
					if attr.TraitType == "BOGO Rewards" {
						assert.Equal(t, 5, attr.Value) // 500 basis points = 5%
						found = true
						break
					}
				}
				assert.True(t, found, "BOGO Rewards attribute not found")
			},
		},
		{
			name: "ticket with custom image",
			data: BOGOWITicketData{
				TokenID:            2,
				BookingID:          "BOOK-456",
				EventID:            "EVENT-789",
				ExperienceType:     "Safari",
				Location:           "Kenya",
				Duration:           "2 days",
				MaxParticipants:    8,
				CarbonOffset:       100,
				ConservationImpact: "Wildlife protection",
				ValidUntil:         validUntil,
				TransferableAfter:  transferableAfter,
				ExpiresAt:          expiresAt,
				BOGORewards:        1000, // 10%
				ImageURL:           "https://custom.com/image.jpg",
			},
			validate: func(t *testing.T, metadata NFTMetadata) {
				assert.Equal(t, "BOGOWI Eco-Experience #2", metadata.Name)
				assert.Equal(t, "https://custom.com/image.jpg", metadata.Image)
				
				// Check carbon offset formatting
				found := false
				for _, attr := range metadata.Attributes {
					if attr.TraitType == "Carbon Offset" {
						assert.Equal(t, "100 kg CO2", attr.Value)
						found = true
						break
					}
				}
				assert.True(t, found, "Carbon Offset attribute not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := service.generateMetadata(tt.data)
			tt.validate(t, metadata)
			
			// Common validations
			assert.NotEmpty(t, metadata.Attributes)
			assert.NotNil(t, metadata.Properties)
			
			// Check properties structure
			bookingDetails, ok := metadata.Properties["booking_details"].(map[string]interface{})
			assert.True(t, ok, "booking_details should be a map")
			assert.Equal(t, tt.data.BookingID, bookingDetails["booking_id"])
			
			redemption, ok := metadata.Properties["redemption"].(map[string]interface{})
			assert.True(t, ok, "redemption should be a map")
			assert.Contains(t, redemption["qr_code"], fmt.Sprintf("%d", tt.data.TokenID))
		})
	}
}

func TestTicketMetadataService_UpdateTicketStatus(t *testing.T) {
	currentNFT := NFT{
		ID:              "nft-123",
		TokenID:         "1",
		ContractAddress: "0x123",
		ChainID:         501,
		Metadata: NFTMetadata{
			Name: "Test NFT",
			Attributes: []NFTAttribute{
				{TraitType: "Status", Value: "Active"},
				{TraitType: "Location", Value: "Costa Rica"},
			},
		},
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		
		if callCount == 1 {
			// First call: GET NFT
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "nft-123")
			
			response := Response{
				Success: true,
				Data:    mustMarshal(currentNFT),
			}
			json.NewEncoder(w).Encode(response)
		} else {
			// Second call: Update metadata
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "metadata")
			
			var metadata NFTMetadata
			err := json.NewDecoder(r.Body).Decode(&metadata)
			require.NoError(t, err)
			
			// Check that status was updated
			found := false
			for _, attr := range metadata.Attributes {
				if attr.TraitType == "Status" {
					assert.Equal(t, "Redeemed", attr.Value)
					found = true
					break
				}
			}
			assert.True(t, found, "Status attribute not found or not updated")
			
			updatedNFT := currentNFT
			updatedNFT.Metadata = metadata
			
			response := Response{
				Success: true,
				Data:    mustMarshal(updatedNFT),
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	service := NewTicketMetadataService("test-api-key", "0x123", 501)
	service.client.baseURL = server.URL

	err := service.UpdateTicketStatus("nft-123", StatusRedeemed)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount, "Expected 2 API calls")
}

func TestTicketMetadataService_GetTicketMetadata(t *testing.T) {
	expectedMetadata := NFTMetadata{
		Name:        "BOGOWI Eco-Experience #1",
		Description: "Nature Tour in Costa Rica",
		Image:       "https://storage.bogowi.com/tickets/1.png",
		Attributes: []NFTAttribute{
			{TraitType: "Experience Type", Value: "Nature Tour"},
			{TraitType: "Status", Value: "Active"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "0x123")
		assert.Contains(t, r.URL.Path, "1")
		assert.Contains(t, r.URL.Path, "metadata")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedMetadata)
	}))
	defer server.Close()

	service := NewTicketMetadataService("test-api-key", "0x123", 501)
	service.client.baseURL = server.URL

	metadata, err := service.GetTicketMetadata(1)
	assert.NoError(t, err)
	assert.NotNil(t, metadata)
	assert.Equal(t, expectedMetadata.Name, metadata.Name)
	assert.Equal(t, len(expectedMetadata.Attributes), len(metadata.Attributes))
}

func TestTicketMetadataService_GetMetadataURI(t *testing.T) {
	service := NewTicketMetadataService("test-api-key", "0x123", 501)

	tests := []struct {
		tokenID     uint64
		expectedURI string
	}{
		{
			tokenID:     1,
			expectedURI: "https://dklnk.to/api/nfts/0x123/1/metadata",
		},
		{
			tokenID:     999,
			expectedURI: "https://dklnk.to/api/nfts/0x123/999/metadata",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("token_%d", tt.tokenID), func(t *testing.T) {
			uri := service.GetMetadataURI(tt.tokenID)
			assert.Equal(t, tt.expectedURI, uri)
		})
	}
}

func TestTicketMetadataService_BatchCreateTickets(t *testing.T) {
	tickets := []BOGOWITicketData{
		{
			TokenID:            1,
			BookingID:          "BOOK-1",
			EventID:            "EVENT-1",
			ExperienceType:     "Nature Tour",
			Location:           "Costa Rica",
			Duration:           "4 hours",
			ValidUntil:         time.Now().Add(30 * 24 * time.Hour),
			TransferableAfter:  time.Now().Add(7 * 24 * time.Hour),
			ExpiresAt:          time.Now().Add(60 * 24 * time.Hour),
		},
		{
			TokenID:            2,
			BookingID:          "BOOK-2",
			EventID:            "EVENT-2",
			ExperienceType:     "Safari",
			Location:           "Kenya",
			Duration:           "2 days",
			ValidUntil:         time.Now().Add(30 * 24 * time.Hour),
			TransferableAfter:  time.Now().Add(7 * 24 * time.Hour),
			ExpiresAt:          time.Now().Add(60 * 24 * time.Hour),
		},
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		
		var reqBody CreateNFTRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Create response based on token ID
		nft := NFT{
			ID:              fmt.Sprintf("nft-%s", reqBody.TokenID),
			TokenID:         reqBody.TokenID,
			ContractAddress: reqBody.ContractAddress,
			ChainID:         reqBody.ChainID,
			Metadata:        reqBody.Metadata,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		response := Response{
			Success: true,
			Data:    mustMarshal(nft),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service := NewTicketMetadataService("test-api-key", "0x123", 501)
	service.client.baseURL = server.URL

	results, err := service.BatchCreateTickets(tickets)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 2, callCount)
	
	// Verify results
	assert.Equal(t, "nft-1", results[0].ID)
	assert.Equal(t, "nft-2", results[1].ID)
}

func TestTicketMetadataService_getNetworkName(t *testing.T) {
	tests := []struct {
		chainID      int
		expectedName string
	}{
		{500, "Camino Mainnet"},
		{501, "Camino Testnet"},
		{1, "Chain 1"},
		{137, "Chain 137"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("chain_%d", tt.chainID), func(t *testing.T) {
			service := NewTicketMetadataService("test-api-key", "0x123", tt.chainID)
			name := service.getNetworkName()
			assert.Equal(t, tt.expectedName, name)
		})
	}
}

func TestTicketStatusConstants(t *testing.T) {
	// Verify all status constants are defined
	assert.Equal(t, "Active", StatusActive)
	assert.Equal(t, "Redeemed", StatusRedeemed)
	assert.Equal(t, "Expired", StatusExpired)
	assert.Equal(t, "Burned", StatusBurned)
	assert.Equal(t, "Pending", StatusPending)
}