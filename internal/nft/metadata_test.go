package nft

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetadataGenerator(t *testing.T) {
	baseImageURL := "https://images.bogowi.com"
	baseExternalURL := "https://bogowi.com"
	ipfsGateway := "https://ipfs.io"

	generator := NewMetadataGenerator(baseImageURL, baseExternalURL, ipfsGateway)

	assert.NotNil(t, generator)
	assert.Equal(t, baseImageURL, generator.baseImageURL)
	assert.Equal(t, baseExternalURL, generator.baseExternalURL)
	assert.Equal(t, ipfsGateway, generator.ipfsGateway)
}

func TestMetadataGenerator_GenerateTicketMetadata(t *testing.T) {
	generator := NewMetadataGenerator(
		"https://images.bogowi.com",
		"https://bogowi.com",
		"https://ipfs.io",
	)

	validUntil := time.Now().Add(30 * 24 * time.Hour)
	transferableAfter := time.Now().Add(7 * 24 * time.Hour)

	tests := []struct {
		name               string
		tokenID            uint64
		bookingID          string
		eventID            string
		experienceType     string
		location           string
		duration           string
		maxParticipants    int
		carbonOffset       int
		conservationImpact string
		bogoRewards        int
		imageHash          string
		validateFunc       func(*testing.T, *Metadata)
	}{
		{
			name:               "basic ticket metadata",
			tokenID:            1,
			bookingID:          "BOOK-123456789",
			eventID:            "EVENT-456",
			experienceType:     "Wildlife Safari",
			location:           "Kenya",
			duration:           "3 days",
			maxParticipants:    10,
			carbonOffset:       50,
			conservationImpact: "Supporting local wildlife conservation",
			bogoRewards:        500,
			imageHash:          "",
			validateFunc: func(t *testing.T, m *Metadata) {
				assert.Equal(t, "BOGOWI Eco-Experience #1", m.Name)
				assert.Contains(t, m.Description, "Wildlife Safari")
				assert.Contains(t, m.Description, "Kenya")
				assert.Equal(t, "https://images.bogowi.com/1.png", m.Image)
				assert.Equal(t, "https://bogowi.com/experiences/1", m.ExternalURL)

				// Check attributes
				assert.Len(t, m.Attributes, 12)

				// Find specific attributes
				var foundType, foundLocation, foundStatus bool
				for _, attr := range m.Attributes {
					switch attr.TraitType {
					case "Experience Type":
						assert.Equal(t, "Wildlife Safari", attr.Value)
						foundType = true
					case "Location":
						assert.Equal(t, "Kenya", attr.Value)
						foundLocation = true
					case "Status":
						assert.Equal(t, "Active", attr.Value)
						foundStatus = true
					case "Carbon Offset":
						assert.Equal(t, "50 kg CO2", attr.Value)
					case "BOGO Rewards":
						assert.Equal(t, 500, attr.Value)
						assert.Equal(t, "boost_percentage", attr.DisplayType)
					}
				}
				assert.True(t, foundType, "Experience Type attribute not found")
				assert.True(t, foundLocation, "Location attribute not found")
				assert.True(t, foundStatus, "Status attribute not found")

				// Check properties
				assert.NotNil(t, m.Properties)
				assert.NotNil(t, m.Properties.BookingDetails)
				assert.Equal(t, "BOOK-123456789", m.Properties.BookingDetails.BookingID)
				assert.Equal(t, "EVENT-456", m.Properties.BookingDetails.EventID)

				assert.NotNil(t, m.Properties.Redemption)
				assert.Equal(t, "bogowi://redeem/1/BOOK-123456789", m.Properties.Redemption.QRCode)
				assert.Equal(t, "BWX-1-BOOK-123", m.Properties.Redemption.RedemptionCode)

				assert.NotNil(t, m.Properties.Rewards)
				assert.Equal(t, 500, m.Properties.Rewards.BOGOTokens)
				assert.Equal(t, 50, m.Properties.Rewards.CarbonCredits)
				assert.Equal(t, 5000, m.Properties.Rewards.LoyaltyPoints)
			},
		},
		{
			name:               "ticket with IPFS image",
			tokenID:            2,
			bookingID:          "BOOK-987654321",
			eventID:            "EVENT-789",
			experienceType:     "Marine Conservation",
			location:           "Maldives",
			duration:           "1 week",
			maxParticipants:    8,
			carbonOffset:       100,
			conservationImpact: "Coral reef restoration",
			bogoRewards:        1000,
			imageHash:          "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
			validateFunc: func(t *testing.T, m *Metadata) {
				assert.Equal(t, "BOGOWI Eco-Experience #2", m.Name)
				assert.Equal(t, "https://ipfs.io/ipfs/QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco", m.Image)
				assert.Contains(t, m.Description, "Marine Conservation")
				assert.Contains(t, m.Description, "Maldives")

				// Check rewards calculation
				assert.Equal(t, 1000, m.Properties.Rewards.BOGOTokens)
				assert.Equal(t, 10000, m.Properties.Rewards.LoyaltyPoints)
			},
		},
		{
			name:               "short booking ID handling",
			tokenID:            3,
			bookingID:          "BK-123", // Short booking ID (less than 8 chars)
			eventID:            "EV-1",
			experienceType:     "Forest Trek",
			location:           "Amazon",
			duration:           "2 days",
			maxParticipants:    5,
			carbonOffset:       30,
			conservationImpact: "Rainforest protection",
			bogoRewards:        300,
			imageHash:          "",
			validateFunc: func(t *testing.T, m *Metadata) {
				// Should handle short booking ID gracefully
				assert.Equal(t, "BWX-3-BK-123", m.Properties.Redemption.RedemptionCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := generator.GenerateTicketMetadata(
				tt.tokenID,
				tt.bookingID,
				tt.eventID,
				tt.experienceType,
				tt.location,
				tt.duration,
				tt.maxParticipants,
				tt.carbonOffset,
				tt.conservationImpact,
				validUntil,
				transferableAfter,
				tt.bogoRewards,
				tt.imageHash,
			)

			assert.NotNil(t, metadata)
			tt.validateFunc(t, metadata)

			// Common validations
			assert.NotEmpty(t, metadata.Name)
			assert.NotEmpty(t, metadata.Description)
			assert.NotEmpty(t, metadata.Image)
			assert.NotEmpty(t, metadata.ExternalURL)
			assert.NotEmpty(t, metadata.Attributes)
		})
	}
}

func TestMetadata_ToJSON(t *testing.T) {
	metadata := &Metadata{
		Name:        "Test NFT",
		Description: "Test Description",
		Image:       "https://example.com/image.png",
		ExternalURL: "https://example.com/nft/1",
		Attributes: []MetadataAttribute{
			{
				TraitType: "Color",
				Value:     "Blue",
			},
			{
				TraitType:   "Level",
				Value:       5,
				DisplayType: "number",
			},
		},
		Properties: &Properties{
			BookingDetails: &BookingDetails{
				BookingID: "BOOK-123",
				EventID:   "EVENT-456",
				Provider:  "Test Provider",
			},
		},
	}

	jsonStr, err := metadata.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonStr)

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &parsed)
	assert.NoError(t, err)

	// Check structure
	assert.Equal(t, "Test NFT", parsed["name"])
	assert.Equal(t, "Test Description", parsed["description"])

	// Check attributes
	attrs, ok := parsed["attributes"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, attrs, 2)

	// Check indentation (should be pretty-printed)
	assert.Contains(t, jsonStr, "\n  ")
}

func TestMetadata_UpdateStatus(t *testing.T) {
	tests := []struct {
		name      string
		metadata  *Metadata
		newStatus string
		validate  func(*testing.T, *Metadata)
	}{
		{
			name: "update existing status",
			metadata: &Metadata{
				Attributes: []MetadataAttribute{
					{TraitType: "Color", Value: "Blue"},
					{TraitType: "Status", Value: "Active"},
					{TraitType: "Level", Value: 5},
				},
			},
			newStatus: "Redeemed",
			validate: func(t *testing.T, m *Metadata) {
				assert.Len(t, m.Attributes, 3)
				found := false
				for _, attr := range m.Attributes {
					if attr.TraitType == "Status" {
						assert.Equal(t, "Redeemed", attr.Value)
						found = true
						break
					}
				}
				assert.True(t, found, "Status attribute not found")
			},
		},
		{
			name: "add status when missing",
			metadata: &Metadata{
				Attributes: []MetadataAttribute{
					{TraitType: "Color", Value: "Blue"},
					{TraitType: "Level", Value: 5},
				},
			},
			newStatus: "Expired",
			validate: func(t *testing.T, m *Metadata) {
				assert.Len(t, m.Attributes, 3)
				found := false
				for _, attr := range m.Attributes {
					if attr.TraitType == "Status" {
						assert.Equal(t, "Expired", attr.Value)
						found = true
						break
					}
				}
				assert.True(t, found, "Status attribute not added")
			},
		},
		{
			name: "empty attributes list",
			metadata: &Metadata{
				Attributes: []MetadataAttribute{},
			},
			newStatus: "Burned",
			validate: func(t *testing.T, m *Metadata) {
				assert.Len(t, m.Attributes, 1)
				assert.Equal(t, "Status", m.Attributes[0].TraitType)
				assert.Equal(t, "Burned", m.Attributes[0].Value)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metadata.UpdateStatus(tt.newStatus)
			tt.validate(t, tt.metadata)
		})
	}
}

func TestMetadataAttribute_JSONSerialization(t *testing.T) {
	tests := []struct {
		name      string
		attribute MetadataAttribute
		validate  func(*testing.T, string)
	}{
		{
			name: "basic attribute",
			attribute: MetadataAttribute{
				TraitType: "Color",
				Value:     "Blue",
			},
			validate: func(t *testing.T, jsonStr string) {
				assert.Contains(t, jsonStr, `"trait_type":"Color"`)
				assert.Contains(t, jsonStr, `"value":"Blue"`)
				assert.NotContains(t, jsonStr, "display_type")
			},
		},
		{
			name: "attribute with display type",
			attribute: MetadataAttribute{
				TraitType:   "Level",
				Value:       10,
				DisplayType: "number",
			},
			validate: func(t *testing.T, jsonStr string) {
				assert.Contains(t, jsonStr, `"trait_type":"Level"`)
				assert.Contains(t, jsonStr, `"value":10`)
				assert.Contains(t, jsonStr, `"display_type":"number"`)
			},
		},
		{
			name: "attribute with date display type",
			attribute: MetadataAttribute{
				TraitType:   "Created",
				Value:       1234567890,
				DisplayType: "date",
			},
			validate: func(t *testing.T, jsonStr string) {
				assert.Contains(t, jsonStr, `"display_type":"date"`)
				assert.Contains(t, jsonStr, `"value":1234567890`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.attribute)
			require.NoError(t, err)
			jsonStr := string(data)
			tt.validate(t, jsonStr)

			// Verify it can be unmarshalled back
			var parsed MetadataAttribute
			err = json.Unmarshal(data, &parsed)
			assert.NoError(t, err)
			assert.Equal(t, tt.attribute.TraitType, parsed.TraitType)
		})
	}
}

func TestProperties_JSONSerialization(t *testing.T) {
	props := &Properties{
		BookingDetails: &BookingDetails{
			BookingID:       "BOOK-123",
			EventID:         "EVENT-456",
			Provider:        "Test Provider",
			ProviderContact: "contact@provider.com",
		},
		Redemption: &RedemptionInfo{
			QRCode:         "bogowi://redeem/1/BOOK-123",
			RedemptionCode: "BWX-1-BOOK-123",
			Instructions:   "Present at venue",
		},
		Rewards: &RewardInfo{
			BOGOTokens:    100,
			CarbonCredits: 50,
			LoyaltyPoints: 1000,
		},
	}

	data, err := json.Marshal(props)
	require.NoError(t, err)

	jsonStr := string(data)

	// Check all fields are present
	assert.Contains(t, jsonStr, "booking_details")
	assert.Contains(t, jsonStr, "BOOK-123")
	assert.Contains(t, jsonStr, "EVENT-456")
	assert.Contains(t, jsonStr, "Test Provider")

	assert.Contains(t, jsonStr, "redemption")
	assert.Contains(t, jsonStr, "bogowi://redeem")
	assert.Contains(t, jsonStr, "BWX-1-BOOK-123")

	assert.Contains(t, jsonStr, "rewards")
	assert.Contains(t, jsonStr, "bogo_tokens")
	assert.Contains(t, jsonStr, "carbon_credits")
	assert.Contains(t, jsonStr, "loyalty_points")

	// Verify it can be unmarshalled back
	var parsed Properties
	err = json.Unmarshal(data, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, props.BookingDetails.BookingID, parsed.BookingDetails.BookingID)
	assert.Equal(t, props.Rewards.BOGOTokens, parsed.Rewards.BOGOTokens)
}

func TestMetadata_CompleteJSONRoundTrip(t *testing.T) {
	generator := NewMetadataGenerator(
		"https://images.bogowi.com",
		"https://bogowi.com",
		"https://ipfs.io",
	)

	original := generator.GenerateTicketMetadata(
		1,
		"BOOK-123456789",
		"EVENT-456",
		"Wildlife Safari",
		"Kenya",
		"3 days",
		10,
		50,
		"Supporting conservation",
		time.Now().Add(30*24*time.Hour),
		time.Now().Add(7*24*time.Hour),
		500,
		"",
	)

	// Convert to JSON
	jsonStr, err := original.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonStr)

	// Parse back
	var parsed Metadata
	err = json.Unmarshal([]byte(jsonStr), &parsed)
	require.NoError(t, err)

	// Verify key fields
	assert.Equal(t, original.Name, parsed.Name)
	assert.Equal(t, original.Description, parsed.Description)
	assert.Equal(t, original.Image, parsed.Image)
	assert.Equal(t, original.ExternalURL, parsed.ExternalURL)
	assert.Len(t, parsed.Attributes, len(original.Attributes))

	// Verify properties
	assert.NotNil(t, parsed.Properties)
	assert.Equal(t, original.Properties.BookingDetails.BookingID, parsed.Properties.BookingDetails.BookingID)
	assert.Equal(t, original.Properties.Rewards.BOGOTokens, parsed.Properties.Rewards.BOGOTokens)
}

func TestMetadataGenerator_EdgeCases(t *testing.T) {
	generator := NewMetadataGenerator(
		"https://images.bogowi.com",
		"https://bogowi.com",
		"https://ipfs.io",
	)

	t.Run("zero values", func(t *testing.T) {
		metadata := generator.GenerateTicketMetadata(
			0,           // zero token ID
			"",          // empty booking ID
			"",          // empty event ID
			"",          // empty experience type
			"",          // empty location
			"",          // empty duration
			0,           // zero participants
			0,           // zero carbon offset
			"",          // empty conservation impact
			time.Time{}, // zero time
			time.Time{}, // zero time
			0,           // zero rewards
			"",          // empty image hash
		)

		assert.NotNil(t, metadata)
		assert.Equal(t, "BOGOWI Eco-Experience #0", metadata.Name)
		assert.Contains(t, metadata.Description, "experience in")
		assert.Equal(t, "https://images.bogowi.com/0.png", metadata.Image)
	})

	t.Run("very long strings", func(t *testing.T) {
		longString := strings.Repeat("A", 1000)
		metadata := generator.GenerateTicketMetadata(
			999999,
			longString,
			longString,
			longString,
			longString,
			longString,
			9999,
			9999,
			longString,
			time.Now(),
			time.Now(),
			999999,
			"",
		)

		assert.NotNil(t, metadata)
		// Should not panic with long strings
		jsonStr, err := metadata.ToJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, longString)
	})

	t.Run("special characters in strings", func(t *testing.T) {
		specialChars := `Special "quotes" & 'apostrophes' <tags> \backslash/ 中文`
		metadata := generator.GenerateTicketMetadata(
			1,
			specialChars,
			"EVENT-1",
			specialChars,
			specialChars,
			"1 day",
			5,
			10,
			specialChars,
			time.Now(),
			time.Now(),
			100,
			"",
		)

		assert.NotNil(t, metadata)
		jsonStr, err := metadata.ToJSON()
		assert.NoError(t, err)

		// JSON should properly escape special characters
		assert.Contains(t, jsonStr, `\"quotes\"`)

		// Should be able to parse back
		var parsed Metadata
		err = json.Unmarshal([]byte(jsonStr), &parsed)
		assert.NoError(t, err)
	})
}

func TestBookingDetails(t *testing.T) {
	details := &BookingDetails{
		BookingID:       "BOOK-123",
		EventID:         "EVENT-456",
		Provider:        "Test Provider",
		ProviderContact: "contact@test.com",
	}

	data, err := json.Marshal(details)
	require.NoError(t, err)

	var parsed BookingDetails
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, details.BookingID, parsed.BookingID)
	assert.Equal(t, details.EventID, parsed.EventID)
	assert.Equal(t, details.Provider, parsed.Provider)
	assert.Equal(t, details.ProviderContact, parsed.ProviderContact)
}

func TestRedemptionInfo(t *testing.T) {
	info := &RedemptionInfo{
		QRCode:         "bogowi://redeem/1/BOOK-123",
		RedemptionCode: "BWX-1-BOOK",
		Instructions:   "Present this at the venue",
	}

	data, err := json.Marshal(info)
	require.NoError(t, err)

	var parsed RedemptionInfo
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, info.QRCode, parsed.QRCode)
	assert.Equal(t, info.RedemptionCode, parsed.RedemptionCode)
	assert.Equal(t, info.Instructions, parsed.Instructions)
}

func TestRewardInfo(t *testing.T) {
	rewards := &RewardInfo{
		BOGOTokens:    500,
		CarbonCredits: 50,
		LoyaltyPoints: 5000,
	}

	data, err := json.Marshal(rewards)
	require.NoError(t, err)

	var parsed RewardInfo
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, rewards.BOGOTokens, parsed.BOGOTokens)
	assert.Equal(t, rewards.CarbonCredits, parsed.CarbonCredits)
	assert.Equal(t, rewards.LoyaltyPoints, parsed.LoyaltyPoints)
}

func BenchmarkGenerateTicketMetadata(b *testing.B) {
	generator := NewMetadataGenerator(
		"https://images.bogowi.com",
		"https://bogowi.com",
		"https://ipfs.io",
	)

	validUntil := time.Now().Add(30 * 24 * time.Hour)
	transferableAfter := time.Now().Add(7 * 24 * time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generator.GenerateTicketMetadata(
			uint64(i),
			fmt.Sprintf("BOOK-%d", i),
			fmt.Sprintf("EVENT-%d", i),
			"Wildlife Safari",
			"Kenya",
			"3 days",
			10,
			50,
			"Conservation impact",
			validUntil,
			transferableAfter,
			500,
			"",
		)
	}
}

func BenchmarkMetadataToJSON(b *testing.B) {
	metadata := &Metadata{
		Name:        "Test NFT",
		Description: "Test Description",
		Image:       "https://example.com/image.png",
		Attributes: []MetadataAttribute{
			{TraitType: "Color", Value: "Blue"},
			{TraitType: "Level", Value: 5},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = metadata.ToJSON()
	}
}
