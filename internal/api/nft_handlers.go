package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/services/datakyte"
	"bogowi-blockchain-go/internal/services/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// NFTHandler handles NFT ticket-related endpoints
type NFTHandler struct {
	*HandlerV2
	metadataServiceTestnet *datakyte.TicketMetadataService
	metadataServiceMainnet *datakyte.TicketMetadataService
	datakyteConfig         *config.DatakyteConfig
	imageService           *storage.ImageService
}

// NewNFTHandler creates a new NFT handler
func NewNFTHandler(h *HandlerV2) *NFTHandler {
	datakyteConfig := config.GetDatakyteConfig()
	
	// Get contract addresses from config or environment
	testnetContract := h.Config.Testnet.Contracts.BOGOWITickets
	if testnetContract == "" {
		testnetContract = os.Getenv("NFT_TICKETS_TESTNET_CONTRACT")
	}
	
	mainnetContract := h.Config.Mainnet.Contracts.BOGOWITickets
	if mainnetContract == "" {
		mainnetContract = os.Getenv("NFT_TICKETS_MAINNET_CONTRACT")
	}
	
	// If still not configured, use placeholder addresses
	if testnetContract == "" {
		testnetContract = "0x0000000000000000000000000000000000000000" // To be deployed
	}
	if mainnetContract == "" {
		mainnetContract = "0x0000000000000000000000000000000000000000" // To be deployed
	}
	
	// Initialize image service
	bucketName := os.Getenv("TICKET_IMAGE_UPLOAD_BUCKET")
	if bucketName == "" {
		bucketName = "bogowi-tickets"
	}
	cdnBaseURL := os.Getenv("TICKET_IMAGE_BASE_URL")
	if cdnBaseURL == "" {
		cdnBaseURL = "https://storage.bogowi.com/tickets"
	}
	
	imageService, _ := storage.NewImageService(bucketName, cdnBaseURL)
	
	return &NFTHandler{
		HandlerV2:      h,
		datakyteConfig: datakyteConfig,
		imageService:   imageService,
		metadataServiceTestnet: datakyte.NewTicketMetadataService(
			datakyteConfig.TestnetAPIKey,
			testnetContract,
			501, // Camino testnet
		),
		metadataServiceMainnet: datakyte.NewTicketMetadataService(
			datakyteConfig.MainnetAPIKey,
			mainnetContract,
			500, // Camino mainnet
		),
	}
}

// getMetadataService returns the appropriate metadata service for the network
func (h *NFTHandler) getMetadataService(network string) *datakyte.TicketMetadataService {
	if network == "mainnet" {
		return h.metadataServiceMainnet
	}
	return h.metadataServiceTestnet
}

// MintTicketRequest represents a request to mint a new ticket
type MintTicketRequest struct {
	To                 string    `json:"to" binding:"required"`
	BookingID          string    `json:"bookingId" binding:"required"`
	EventID            string    `json:"eventId" binding:"required"`
	ExperienceTitle    string    `json:"experienceTitle" binding:"required"`
	ExperienceType     string    `json:"experienceType" binding:"required"`
	Location           string    `json:"location" binding:"required"`
	Duration           string    `json:"duration" binding:"required"`
	MaxParticipants    int       `json:"maxParticipants"`
	CarbonOffset       int       `json:"carbonOffset"`
	ConservationImpact string    `json:"conservationImpact"`
	ValidUntil         time.Time `json:"validUntil" binding:"required"`
	TransferableAfter  time.Time `json:"transferableAfter" binding:"required"`
	ExpiresAt          time.Time `json:"expiresAt" binding:"required"`
	RewardBasisPoints  uint16    `json:"rewardBasisPoints"`
	RecipientName      string    `json:"recipientName"`
	ProviderName       string    `json:"providerName"`
	ProviderContact    string    `json:"providerContact"`
	ImageURL           string    `json:"imageUrl,omitempty"`
}

// MintTicketResponse represents the response after minting a ticket
type MintTicketResponse struct {
	Success     bool   `json:"success"`
	TokenID     uint64 `json:"tokenId"`
	TxHash      string `json:"txHash"`
	MetadataURI string `json:"metadataUri"`
	DatakyteID  string `json:"datakyteId,omitempty"`
	Message     string `json:"message,omitempty"`
}

// MintTicket mints a new NFT ticket
func (h *NFTHandler) MintTicket(c *gin.Context) {
	network := GetNetworkFromContext(c)
	
	var req MintTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	// Step 1: Mint the NFT on-chain
	tokenID, txHash, err := h.mintOnChain(sdk, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to mint NFT: %v", err),
		})
		return
	}
	
	// Step 2: Create metadata in Datakyte
	metadataService := h.getMetadataService(network)
	
	ticketData := datakyte.BOGOWITicketData{
		TokenID:            tokenID,
		BookingID:          req.BookingID,
		EventID:            req.EventID,
		ExperienceTitle:    req.ExperienceTitle,
		ExperienceType:     req.ExperienceType,
		Location:           req.Location,
		Duration:           req.Duration,
		MaxParticipants:    req.MaxParticipants,
		CarbonOffset:       req.CarbonOffset,
		ConservationImpact: req.ConservationImpact,
		ValidUntil:         req.ValidUntil,
		TransferableAfter:  req.TransferableAfter,
		ExpiresAt:          req.ExpiresAt,
		BOGORewards:        int(req.RewardBasisPoints),
		RecipientAddress:   req.To,
		RecipientName:      req.RecipientName,
		ProviderName:       req.ProviderName,
		ProviderContact:    req.ProviderContact,
		ImageURL:           req.ImageURL,
	}
	
	nft, err := metadataService.CreateTicketMetadata(ticketData)
	if err != nil {
		// Log error but don't fail - NFT is already minted
		fmt.Printf("Warning: Failed to create Datakyte metadata for token %d: %v\n", tokenID, err)
	}
	
	// Prepare response
	response := MintTicketResponse{
		Success:     true,
		TokenID:     tokenID,
		TxHash:      txHash,
		MetadataURI: metadataService.GetMetadataURI(tokenID),
	}
	
	if nft != nil {
		response.DatakyteID = nft.ID
	}
	
	c.JSON(http.StatusOK, response)
}

// mintOnChain handles the blockchain transaction for minting
func (h *NFTHandler) mintOnChain(sdk interface{}, req MintTicketRequest) (uint64, string, error) {
	// TODO: Implement actual SDK call to mint NFT
	// This is a placeholder - needs to be integrated with your SDK
	
	// For now, return mock data
	tokenID := uint64(10001) // This would come from the smart contract
	txHash := "0x" + common.Bytes2Hex([]byte("mock_tx_hash"))
	
	return tokenID, txHash, nil
}

// GetTicketMetadata retrieves metadata for a ticket
func (h *NFTHandler) GetTicketMetadata(c *gin.Context) {
	network := GetNetworkFromContext(c)
	metadataService := h.getMetadataService(network)
	
	tokenIDStr := c.Param("tokenId")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid token ID"})
		return
	}
	
	metadata, err := metadataService.GetTicketMetadata(tokenID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Metadata not found"})
		return
	}
	
	c.JSON(http.StatusOK, metadata)
}

// RedeemTicketRequest represents a request to redeem a ticket
type RedeemTicketRequest struct {
	TokenID   uint64 `json:"tokenId" binding:"required"`
	Redeemer  string `json:"redeemer" binding:"required"`
	Nonce     uint64 `json:"nonce" binding:"required"`
	Deadline  uint64 `json:"deadline" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// RedeemTicket handles ticket redemption
func (h *NFTHandler) RedeemTicket(c *gin.Context) {
	network := GetNetworkFromContext(c)
	
	var req RedeemTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	// TODO: Call smart contract to redeem ticket
	_ = sdk
	
	// Update metadata status in Datakyte
	// Note: We'd need to store the Datakyte NFT ID mapping somewhere
	// For now, this is a placeholder
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket redeemed successfully",
		"tokenId": req.TokenID,
	})
}

// GetUserTickets retrieves all tickets for a user
func (h *NFTHandler) GetUserTickets(c *gin.Context) {
	userAddress := c.Param("address")
	if !common.IsHexAddress(userAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid address"})
		return
	}
	
	// TODO: Query blockchain for user's NFTs
	// Then get metadata from Datakyte for each
	
	c.JSON(http.StatusOK, gin.H{
		"tickets": []interface{}{},
		"message": "Endpoint under development",
	})
}

// BatchMintRequest represents a request to mint multiple tickets
type BatchMintRequest struct {
	Tickets []MintTicketRequest `json:"tickets" binding:"required,dive"`
}

// BatchMintTickets mints multiple tickets at once
func (h *NFTHandler) BatchMintTickets(c *gin.Context) {
	network := GetNetworkFromContext(c)
	
	var req BatchMintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	if len(req.Tickets) > 100 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Batch size exceeds maximum of 100"})
		return
	}
	
	// Get SDK for the network
	sdk, err := h.NetworkHandler.GetSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	// TODO: Implement batch minting
	_ = sdk
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Batch minting %d tickets", len(req.Tickets)),
		"status":  "Under development",
	})
}

// GetPresignedUploadURL generates a presigned URL for image upload
func (h *NFTHandler) GetPresignedUploadURL(c *gin.Context) {
	tokenIDStr := c.Param("tokenId")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid token ID"})
		return
	}
	
	contentType := c.DefaultQuery("contentType", "image/jpeg")
	
	if h.imageService == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Image service not configured"})
		return
	}
	
	uploadURL, err := h.imageService.GeneratePresignedUploadURL(tokenID, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to generate upload URL: %v", err),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"uploadUrl":   uploadURL,
		"tokenId":     tokenID,
		"contentType": contentType,
		"expiresIn":   900, // 15 minutes
	})
}

// UploadTicketImage handles direct image upload
func (h *NFTHandler) UploadTicketImage(c *gin.Context) {
	tokenIDStr := c.Param("tokenId")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid token ID"})
		return
	}
	
	// Get the uploaded file
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No image file provided"})
		return
	}
	defer file.Close()
	
	if h.imageService == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Image service not configured"})
		return
	}
	
	// Validate image
	config := storage.DefaultUploadConfig()
	
	// Create a copy for validation since we'll consume the reader
	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)
	
	if err := h.imageService.ValidateImage(tee, config); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	// Process and upload
	imageURL, err := h.imageService.ProcessAndUploadImage(&buf, tokenID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to process image: %v", err),
		})
		return
	}
	
	// Update metadata with image URL
	network := GetNetworkFromContext(c)
	metadataService := h.getMetadataService(network)
	
	// Get current metadata and update image
	metadata, err := metadataService.GetTicketMetadata(tokenID)
	if err == nil && metadata != nil {
		metadata.Image = imageURL
		// TODO: Update metadata in Datakyte
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"tokenId":  tokenID,
		"imageUrl": imageURL,
		"sizes":    h.imageService.GetImageURLs(tokenID),
		"message":  "Image uploaded successfully",
	})
}

// UpdateTicketStatus updates the status of a ticket in metadata
func (h *NFTHandler) UpdateTicketStatus(c *gin.Context) {
	tokenIDStr := c.Param("tokenId")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid token ID"})
		return
	}
	
	var req struct {
		Status string `json:"status" binding:"required,oneof=Active Redeemed Expired Burned"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	// TODO: Get Datakyte NFT ID from database or mapping
	// For now, this is a placeholder
	datakyteNFTID := fmt.Sprintf("placeholder-%d", tokenID)
	
	network := GetNetworkFromContext(c)
	metadataService := h.getMetadataService(network)
	
	err = metadataService.UpdateTicketStatus(datakyteNFTID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to update status: %v", err),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tokenId": tokenID,
		"status":  req.Status,
	})
}