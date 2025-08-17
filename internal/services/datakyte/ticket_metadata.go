package datakyte

import (
	"fmt"
	"time"
)

// TicketMetadataService handles BOGOWI ticket metadata using Datakyte
type TicketMetadataService struct {
	client          *Client
	contractAddress string
	chainID         int
	baseImageURL    string
}

// NewTicketMetadataService creates a new ticket metadata service
func NewTicketMetadataService(apiKey, contractAddress string, chainID int) *TicketMetadataService {
	return &TicketMetadataService{
		client:          NewClient(apiKey),
		contractAddress: contractAddress,
		chainID:         chainID,
		baseImageURL:    "https://storage.bogowi.com/tickets", // Can be configured
	}
}

// BOGOWITicketData contains all information for a BOGOWI ticket
type BOGOWITicketData struct {
	TokenID            uint64
	BookingID          string
	EventID            string
	ExperienceTitle    string
	ExperienceType     string
	Location           string
	Duration           string
	MaxParticipants    int
	CarbonOffset       int // kg CO2
	ConservationImpact string
	ValidUntil         time.Time
	TransferableAfter  time.Time
	ExpiresAt          time.Time
	BOGORewards        int // basis points
	RecipientAddress   string
	RecipientName      string
	ProviderName       string
	ProviderContact    string
	ImageURL           string // Optional custom image URL
}

// CreateTicketMetadata creates metadata for a new BOGOWI ticket
func (s *TicketMetadataService) CreateTicketMetadata(data BOGOWITicketData) (*NFT, error) {
	// Generate metadata
	metadata := s.generateMetadata(data)

	// Create NFT request
	request := CreateNFTRequest{
		TokenID:         fmt.Sprintf("%d", data.TokenID),
		ContractAddress: s.contractAddress,
		ChainID:         s.chainID,
		Metadata:        metadata,
		CustomData: map[string]interface{}{
			"bookingId":        data.BookingID,
			"eventId":          data.EventID,
			"recipientAddress": data.RecipientAddress,
			"provider":         data.ProviderName,
			"createdAt":        time.Now().Format(time.RFC3339),
		},
	}

	// Create NFT in Datakyte
	nft, err := s.client.CreateNFT(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create NFT metadata: %w", err)
	}

	return nft, nil
}

// generateMetadata creates the NFTMetadata structure for a ticket
func (s *TicketMetadataService) generateMetadata(data BOGOWITicketData) NFTMetadata {
	// Determine image URL
	imageURL := data.ImageURL
	if imageURL == "" {
		imageURL = fmt.Sprintf("%s/%d.png", s.baseImageURL, data.TokenID)
	}

	// Create attributes
	attributes := []NFTAttribute{
		{
			TraitType: "Experience Type",
			Value:     data.ExperienceType,
		},
		{
			TraitType: "Location",
			Value:     data.Location,
		},
		{
			TraitType: "Duration",
			Value:     data.Duration,
		},
		{
			TraitType: "Max Participants",
			Value:     data.MaxParticipants,
		},
		{
			TraitType: "Carbon Offset",
			Value:     fmt.Sprintf("%d kg CO2", data.CarbonOffset),
		},
		{
			TraitType: "Conservation Impact",
			Value:     data.ConservationImpact,
		},
		{
			TraitType:   "Valid Until",
			Value:       data.ValidUntil.Unix(),
			DisplayType: "date",
		},
		{
			TraitType:   "Transferable After",
			Value:       data.TransferableAfter.Unix(),
			DisplayType: "date",
		},
		{
			TraitType:   "Expires At",
			Value:       data.ExpiresAt.Unix(),
			DisplayType: "date",
		},
		{
			TraitType:   "BOGO Rewards",
			Value:       data.BOGORewards / 100, // Convert basis points to percentage
			DisplayType: "boost_percentage",
		},
		{
			TraitType: "Status",
			Value:     "Active",
		},
		{
			TraitType: "Network",
			Value:     s.getNetworkName(),
		},
	}

	// Create properties
	properties := map[string]interface{}{
		"booking_details": map[string]interface{}{
			"booking_id":       data.BookingID,
			"event_id":         data.EventID,
			"provider":         data.ProviderName,
			"provider_contact": data.ProviderContact,
		},
		"redemption": map[string]interface{}{
			"qr_code":         fmt.Sprintf("bogowi://redeem/%d/%s", data.TokenID, data.BookingID),
			"redemption_code": fmt.Sprintf("BWX-%d-%s", data.TokenID, truncateString(data.BookingID, 8)),
			"instructions":    "Present this QR code at the venue to redeem your experience",
		},
		"rewards": map[string]interface{}{
			"bogo_tokens":    data.BOGORewards,
			"carbon_credits": data.CarbonOffset,
			"loyalty_points": data.BOGORewards * 10,
		},
	}

	return NFTMetadata{
		Name:         fmt.Sprintf("BOGOWI Eco-Experience #%d", data.TokenID),
		Description:  fmt.Sprintf("%s in %s. Duration: %s. %s", data.ExperienceType, data.Location, data.Duration, data.ConservationImpact),
		Image:        imageURL,
		ExternalURL:  fmt.Sprintf("https://bogowi.com/experiences/%d", data.TokenID),
		AnimationURL: "", // Could add animated QR code or video
		Attributes:   attributes,
		Properties:   properties,
	}
}

// UpdateTicketStatus updates the status of a ticket (e.g., after redemption)
func (s *TicketMetadataService) UpdateTicketStatus(nftID string, status string) error {
	// Get current metadata
	nft, err := s.client.GetNFT(nftID)
	if err != nil {
		return fmt.Errorf("failed to get NFT: %w", err)
	}

	// Update status attribute
	for i, attr := range nft.Metadata.Attributes {
		if attr.TraitType == "Status" {
			nft.Metadata.Attributes[i].Value = status
			break
		}
	}

	// Update metadata
	_, err = s.client.UpdateMetadata(nftID, nft.Metadata)
	if err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	return nil
}

// GetTicketMetadata retrieves metadata for a ticket by token ID
func (s *TicketMetadataService) GetTicketMetadata(tokenID uint64) (*NFTMetadata, error) {
	metadata, err := s.client.GetMetadataByTokenID(s.contractAddress, fmt.Sprintf("%d", tokenID))
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	return metadata, nil
}

// GetMetadataURI returns the URI for accessing metadata (for smart contract)
func (s *TicketMetadataService) GetMetadataURI(tokenID uint64) string {
	// This returns the Datakyte endpoint that can be used in tokenURI
	return fmt.Sprintf("https://dklnk.to/api/nfts/%s/%d/metadata", s.contractAddress, tokenID)
}

// BatchCreateTickets creates metadata for multiple tickets
func (s *TicketMetadataService) BatchCreateTickets(tickets []BOGOWITicketData) ([]*NFT, error) {
	results := make([]*NFT, 0, len(tickets))

	for _, ticket := range tickets {
		nft, err := s.CreateTicketMetadata(ticket)
		if err != nil {
			// Log error but continue with other tickets
			fmt.Printf("Failed to create metadata for token %d: %v\n", ticket.TokenID, err)
			continue
		}
		results = append(results, nft)
	}

	return results, nil
}

// getNetworkName returns the network name based on chain ID
func (s *TicketMetadataService) getNetworkName() string {
	switch s.chainID {
	case 500:
		return "Camino Mainnet"
	case 501:
		return "Camino Testnet"
	default:
		return fmt.Sprintf("Chain %d", s.chainID)
	}
}

// TicketStatus represents possible ticket statuses
const (
	StatusActive   = "Active"
	StatusRedeemed = "Redeemed"
	StatusExpired  = "Expired"
	StatusBurned   = "Burned"
	StatusPending  = "Pending"
)

// truncateString safely truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
