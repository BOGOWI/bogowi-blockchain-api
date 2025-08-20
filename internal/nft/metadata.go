package nft

import (
	"encoding/json"
	"fmt"
	"time"
)

// MetadataAttribute represents a single attribute in NFT metadata
type MetadataAttribute struct {
	TraitType   string      `json:"trait_type"`
	Value       interface{} `json:"value"`
	DisplayType string      `json:"display_type,omitempty"`
}

// BookingDetails contains booking-specific information
type BookingDetails struct {
	BookingID       string `json:"booking_id"`
	EventID         string `json:"event_id"`
	Provider        string `json:"provider"`
	ProviderContact string `json:"provider_contact"`
}

// RedemptionInfo contains redemption-specific information
type RedemptionInfo struct {
	QRCode         string `json:"qr_code"`
	RedemptionCode string `json:"redemption_code"`
	Instructions   string `json:"instructions"`
}

// RewardInfo contains reward-specific information
type RewardInfo struct {
	BOGOTokens    int `json:"bogo_tokens"`
	CarbonCredits int `json:"carbon_credits"`
	LoyaltyPoints int `json:"loyalty_points"`
}

// Properties contains additional structured data
type Properties struct {
	BookingDetails *BookingDetails `json:"booking_details,omitempty"`
	Redemption     *RedemptionInfo `json:"redemption,omitempty"`
	Rewards        *RewardInfo     `json:"rewards,omitempty"`
}

// Metadata represents the complete NFT metadata structure
type Metadata struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Image        string              `json:"image"`
	AnimationURL string              `json:"animation_url,omitempty"`
	ExternalURL  string              `json:"external_url,omitempty"`
	Attributes   []MetadataAttribute `json:"attributes"`
	Properties   *Properties         `json:"properties,omitempty"`
}

// MetadataGenerator handles creation of NFT metadata
type MetadataGenerator struct {
	baseImageURL    string
	baseExternalURL string
	ipfsGateway     string
}

// NewMetadataGenerator creates a new metadata generator
func NewMetadataGenerator(baseImageURL, baseExternalURL, ipfsGateway string) *MetadataGenerator {
	return &MetadataGenerator{
		baseImageURL:    baseImageURL,
		baseExternalURL: baseExternalURL,
		ipfsGateway:     ipfsGateway,
	}
}

// GenerateTicketMetadata creates metadata for a BOGOWI ticket NFT
func (g *MetadataGenerator) GenerateTicketMetadata(
	tokenID uint64,
	bookingID string,
	eventID string,
	experienceType string,
	location string,
	duration string,
	maxParticipants int,
	carbonOffset int,
	conservationImpact string,
	validUntil time.Time,
	transferableAfter time.Time,
	bogoRewards int,
	imageHash string,
) *Metadata {
	// Generate QR code data
	qrData := fmt.Sprintf("bogowi://redeem/%d/%s", tokenID, bookingID)

	// Create attributes
	attributes := []MetadataAttribute{
		{
			TraitType: "Experience Type",
			Value:     experienceType,
		},
		{
			TraitType: "Location",
			Value:     location,
		},
		{
			TraitType: "Duration",
			Value:     duration,
		},
		{
			TraitType: "Max Participants",
			Value:     maxParticipants,
		},
		{
			TraitType: "Carbon Offset",
			Value:     fmt.Sprintf("%d kg CO2", carbonOffset),
		},
		{
			TraitType: "Conservation Impact",
			Value:     conservationImpact,
		},
		{
			TraitType:   "Booking ID",
			Value:       bookingID,
			DisplayType: "string",
		},
		{
			TraitType:   "Event ID",
			Value:       eventID,
			DisplayType: "string",
		},
		{
			TraitType:   "Valid Until",
			Value:       validUntil.Unix(),
			DisplayType: "date",
		},
		{
			TraitType:   "Transferable After",
			Value:       transferableAfter.Unix(),
			DisplayType: "date",
		},
		{
			TraitType:   "BOGO Rewards",
			Value:       bogoRewards,
			DisplayType: "boost_percentage",
		},
		{
			TraitType: "Status",
			Value:     "Active",
		},
	}

	// Build image URL (could be IPFS or regular URL)
	var imageURL string
	if imageHash != "" {
		imageURL = fmt.Sprintf("%s/ipfs/%s", g.ipfsGateway, imageHash)
	} else {
		imageURL = fmt.Sprintf("%s/%d.png", g.baseImageURL, tokenID)
	}

	// Create metadata
	metadata := &Metadata{
		Name:        fmt.Sprintf("BOGOWI Eco-Experience #%d", tokenID),
		Description: fmt.Sprintf("%s experience in %s. %s", experienceType, location, conservationImpact),
		Image:       imageURL,
		ExternalURL: fmt.Sprintf("%s/experiences/%d", g.baseExternalURL, tokenID),
		Attributes:  attributes,
		Properties: &Properties{
			BookingDetails: &BookingDetails{
				BookingID: bookingID,
				EventID:   eventID,
				Provider:  "BOGOWI Partner Network",
			},
			Redemption: &RedemptionInfo{
				QRCode:         qrData,
				RedemptionCode: fmt.Sprintf("BWX-%d-%s", tokenID, truncateBookingID(bookingID, 8)),
				Instructions:   "Present this QR code at the venue to redeem your experience",
			},
			Rewards: &RewardInfo{
				BOGOTokens:    bogoRewards,
				CarbonCredits: carbonOffset,
				LoyaltyPoints: bogoRewards * 10,
			},
		},
	}

	return metadata
}

// ToJSON converts metadata to JSON string
func (m *Metadata) ToJSON() (string, error) {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UpdateStatus updates the status attribute of the metadata
func (m *Metadata) UpdateStatus(newStatus string) {
	for i, attr := range m.Attributes {
		if attr.TraitType == "Status" {
			m.Attributes[i].Value = newStatus
			return
		}
	}
	// If status attribute doesn't exist, add it
	m.Attributes = append(m.Attributes, MetadataAttribute{
		TraitType: "Status",
		Value:     newStatus,
	})
}

// truncateBookingID safely truncates a booking ID to a maximum length
func truncateBookingID(bookingID string, maxLen int) string {
	if len(bookingID) <= maxLen {
		return bookingID
	}
	return bookingID[:maxLen]
}
