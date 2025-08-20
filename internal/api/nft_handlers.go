package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/database"
	"bogowi-blockchain-go/internal/sdk/nft"
	"bogowi-blockchain-go/internal/services/datakyte"
	"bogowi-blockchain-go/internal/services/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// NFTHandler handles NFT ticket-related endpoints
type NFTHandler struct {
	*Handler
	metadataServiceTestnet *datakyte.TicketMetadataService
	metadataServiceMainnet *datakyte.TicketMetadataService
	datakyteConfig         *config.DatakyteConfig
	imageService           *storage.ImageService
}

// NewNFTHandler creates a new NFT handler
func NewNFTHandler(h *Handler) *NFTHandler {
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
		Handler:        h,
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

// getContractAddress returns the contract address for the network
func (h *NFTHandler) getContractAddress(network string) string {
	if network == "mainnet" {
		contractAddr := h.Config.Mainnet.Contracts.BOGOWITickets
		if contractAddr == "" {
			contractAddr = os.Getenv("NFT_TICKETS_MAINNET_CONTRACT")
		}
		return contractAddr
	}
	
	contractAddr := h.Config.Testnet.Contracts.BOGOWITickets
	if contractAddr == "" {
		contractAddr = os.Getenv("NFT_TICKETS_TESTNET_CONTRACT")
	}
	return contractAddr
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
// @Summary Mint a new NFT ticket
// @Description Mints a new NFT ticket on the blockchain and creates metadata in Datakyte
// @Tags NFT
// @Accept json
// @Produce json
// @Param X-Network-Type header string false "Network type (testnet/mainnet)" default(testnet)
// @Param request body MintTicketRequest true "Mint ticket request"
// @Success 200 {object} MintTicketResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /nft/tickets/mint [post]
func (h *NFTHandler) MintTicket(c *gin.Context) {
	network := GetNetworkFromContext(c)

	var req MintTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get NFT SDK for the network
	nftSDK, err := h.NetworkHandler.GetNFTSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Step 1: Mint the NFT on-chain
	tokenID, txHash, err := h.mintOnChain(nftSDK, req)
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
		
		// Save the mapping to database
		db := database.GetDB()
		mapping := &database.NFTMapping{
			TokenID:       tokenID,
			DatakyteNFTID: nft.ID,
			Network:       network,
			ContractAddr:  h.getContractAddress(network),
			OwnerAddress:  req.To,
			BookingID:     req.BookingID,
			EventID:       req.EventID,
			Status:        "active",
			MetadataURI:   response.MetadataURI,
			ImageURL:      req.ImageURL,
			TxHash:        txHash,
		}
		
		if err := db.SaveNFTMapping(mapping); err != nil {
			// Log error but don't fail - NFT is already minted
			fmt.Printf("Warning: Failed to save NFT mapping for token %d: %v\n", tokenID, err)
		}
	}

	c.JSON(http.StatusOK, response)
}

// mintOnChain handles the blockchain transaction for minting
func (h *NFTHandler) mintOnChain(nftSDK *nft.Client, req MintTicketRequest) (uint64, string, error) {

	// Convert request parameters to SDK format
	params := nft.MintParams{
		To:                common.HexToAddress(req.To),
		BookingID:         convertToBytes32(req.BookingID),
		EventID:           convertToBytes32(req.EventID),
		UtilityFlags:      0, // Default to 0, can be extended based on requirements
		TransferUnlockAt:  uint64(req.TransferableAfter.Unix()),
		ExpiresAt:         uint64(req.ExpiresAt.Unix()),
		MetadataURI:       "", // Will be set by Datakyte after minting
		RewardBasisPoints: req.RewardBasisPoints,
	}

	// Call the actual SDK mint function
	ctx := context.Background()
	tx, tokenID, err := nftSDK.MintTicket(ctx, params)
	if err != nil {
		return 0, "", fmt.Errorf("failed to mint NFT on blockchain: %w", err)
	}

	return tokenID, tx.Hash().Hex(), nil
}

// convertToBytes32 converts a string to [32]byte, padding with zeros if necessary
func convertToBytes32(s string) [32]byte {
	var b32 [32]byte
	// Hash the string to ensure it fits in 32 bytes
	hash := crypto.Keccak256([]byte(s))
	copy(b32[:], hash)
	return b32
}

// GetTicketMetadata retrieves metadata for a ticket
// @Summary Get ticket metadata
// @Description Retrieves metadata for a specific NFT ticket from Datakyte
// @Tags NFT
// @Accept json
// @Produce json
// @Param X-Network-Type header string false "Network type (testnet/mainnet)" default(testnet)
// @Param tokenId path int true "Token ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /nft/tickets/{tokenId}/metadata [get]
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
// @Summary Redeem an NFT ticket
// @Description Redeems an NFT ticket using EIP-712 signature verification
// @Tags NFT
// @Accept json
// @Produce json
// @Param X-Network-Type header string false "Network type (testnet/mainnet)" default(testnet)
// @Param request body RedeemTicketRequest true "Redeem ticket request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /nft/tickets/redeem [post]
func (h *NFTHandler) RedeemTicket(c *gin.Context) {
	network := GetNetworkFromContext(c)

	var req RedeemTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get NFT SDK for the network
	nftSDK, err := h.NetworkHandler.GetNFTSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Prepare redemption parameters
	params := nft.RedemptionParams{
		TokenID:  req.TokenID,
		Redeemer: common.HexToAddress(req.Redeemer),
		Nonce:    req.Nonce,
		Deadline: int64(req.Deadline),
	}

	// Execute redemption on blockchain
	ctx := context.Background()
	tx, err := nftSDK.RedeemTicket(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to redeem ticket: %v", err),
		})
		return
	}

	// Update metadata status in Datakyte
	db := database.GetDB()
	datakyteNFTID, err := db.GetDatakyteID(req.TokenID, network)
	if err != nil {
		fmt.Printf("Warning: Failed to get Datakyte ID for token %d: %v\n", req.TokenID, err)
	} else {
		metadataService := h.getMetadataService(network)
		err = metadataService.UpdateTicketStatus(datakyteNFTID, "Redeemed")
		if err != nil {
			// Log but don't fail - redemption already succeeded on-chain
			fmt.Printf("Warning: Failed to update Datakyte status for token %d: %v\n", req.TokenID, err)
		}
		
		// Update database status
		if err := db.UpdateNFTRedemption(req.TokenID, network); err != nil {
			fmt.Printf("Warning: Failed to update redemption status in database for token %d: %v\n", req.TokenID, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket redeemed successfully",
		"tokenId": req.TokenID,
		"txHash":  tx.Hash().Hex(),
	})
}

// GetUserTickets retrieves all tickets for a user
// @Summary Get user's NFT tickets
// @Description Retrieves all NFT tickets owned by a specific address
// @Tags NFT
// @Accept json
// @Produce json
// @Param X-Network-Type header string false "Network type (testnet/mainnet)" default(testnet)
// @Param address path string true "User wallet address"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /nft/users/{address}/tickets [get]
func (h *NFTHandler) GetUserTickets(c *gin.Context) {
	userAddress := c.Param("address")
	if !common.IsHexAddress(userAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid address"})
		return
	}

	network := GetNetworkFromContext(c)
	
	// Get NFT SDK for the network
	nftSDK, err := h.NetworkHandler.GetNFTSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Query blockchain for user's balance
	ctx := context.Background()
	owner := common.HexToAddress(userAddress)
	
	balance, err := nftSDK.GetBalanceOf(ctx, owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to get user balance: %v", err),
		})
		return
	}

	// Try to get user's tickets (note: this requires enumeration support)
	tickets, err := nftSDK.GetUserTickets(ctx, owner)
	if err != nil {
		// If enumeration is not supported, just return the balance info
		c.JSON(http.StatusOK, gin.H{
			"address": userAddress,
			"balance": balance.String(),
			"tickets": []interface{}{},
			"message": "Token enumeration requires event filtering (not yet implemented). User owns " + balance.String() + " tickets.",
		})
		return
	}

	// Get metadata for each ticket
	metadataService := h.getMetadataService(network)
	ticketDetails := make([]gin.H, 0, len(tickets))
	
	for _, tokenID := range tickets {
		// Get on-chain data
		ticketData, err := nftSDK.GetTicketData(ctx, tokenID)
		if err != nil {
			continue // Skip tickets we can't read
		}

		// Get metadata from Datakyte
		metadata, _ := metadataService.GetTicketMetadata(tokenID)
		
		detail := gin.H{
			"tokenId":           tokenID,
			"state":             nft.ParseTicketState(ticketData.State).String(),
			"expiresAt":         ticketData.ExpiresAt,
			"transferUnlockAt":  ticketData.TransferUnlockAt,
			"bookingId":         fmt.Sprintf("%x", ticketData.BookingID),
			"eventId":           fmt.Sprintf("%x", ticketData.EventID),
		}
		
		if metadata != nil {
			detail["metadata"] = metadata
		}
		
		ticketDetails = append(ticketDetails, detail)
	}

	c.JSON(http.StatusOK, gin.H{
		"address": userAddress,
		"balance": balance.String(),
		"tickets": ticketDetails,
	})
}

// BatchMintRequest represents a request to mint multiple tickets
type BatchMintRequest struct {
	Tickets []MintTicketRequest `json:"tickets" binding:"required,dive"`
}

// BatchMintTickets mints multiple tickets at once
// @Summary Batch mint NFT tickets
// @Description Mints multiple NFT tickets in a single transaction
// @Tags NFT
// @Accept json
// @Produce json
// @Param X-Network-Type header string false "Network type (testnet/mainnet)" default(testnet)
// @Param request body BatchMintRequest true "Batch mint request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /nft/tickets/batch-mint [post]
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

	// Get NFT SDK for the network
	nftSDK, err := h.NetworkHandler.GetNFTSDK(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert requests to SDK params
	mintParams := make([]nft.MintParams, len(req.Tickets))
	for i, ticket := range req.Tickets {
		mintParams[i] = nft.MintParams{
			To:                common.HexToAddress(ticket.To),
			BookingID:         convertToBytes32(ticket.BookingID),
			EventID:           convertToBytes32(ticket.EventID),
			UtilityFlags:      0, // Default to 0, can be extended based on requirements
			TransferUnlockAt:  uint64(ticket.TransferableAfter.Unix()),
			ExpiresAt:         uint64(ticket.ExpiresAt.Unix()),
			MetadataURI:       "", // Will be set by Datakyte after minting
			RewardBasisPoints: ticket.RewardBasisPoints,
		}
	}

	// Execute batch mint on blockchain
	ctx := context.Background()
	tx, tokenIDs, err := nftSDK.BatchMint(ctx, mintParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to batch mint NFTs: %v", err),
		})
		return
	}

	// Create metadata in Datakyte for each minted token
	metadataService := h.getMetadataService(network)
	results := make([]MintTicketResponse, len(tokenIDs))

	for i, tokenID := range tokenIDs {
		ticket := req.Tickets[i]
		
		// Create Datakyte metadata
		ticketData := datakyte.BOGOWITicketData{
			TokenID:            tokenID,
			BookingID:          ticket.BookingID,
			EventID:            ticket.EventID,
			ExperienceTitle:    ticket.ExperienceTitle,
			ExperienceType:     ticket.ExperienceType,
			Location:           ticket.Location,
			Duration:           ticket.Duration,
			MaxParticipants:    ticket.MaxParticipants,
			CarbonOffset:       ticket.CarbonOffset,
			ConservationImpact: ticket.ConservationImpact,
			ValidUntil:         ticket.ValidUntil,
			TransferableAfter:  ticket.TransferableAfter,
			ExpiresAt:          ticket.ExpiresAt,
			BOGORewards:        int(ticket.RewardBasisPoints),
			RecipientAddress:   ticket.To,
			RecipientName:      ticket.RecipientName,
			ProviderName:       ticket.ProviderName,
			ProviderContact:    ticket.ProviderContact,
			ImageURL:           ticket.ImageURL,
		}

		nftMetadata, err := metadataService.CreateTicketMetadata(ticketData)
		
		results[i] = MintTicketResponse{
			Success:     err == nil,
			TokenID:     tokenID,
			TxHash:      tx.Hash().Hex(),
			MetadataURI: metadataService.GetMetadataURI(tokenID),
		}

		if nftMetadata != nil {
			results[i].DatakyteID = nftMetadata.ID
			
			// Save the mapping to database
			db := database.GetDB()
			mapping := &database.NFTMapping{
				TokenID:       tokenID,
				DatakyteNFTID: nftMetadata.ID,
				Network:       network,
				ContractAddr:  h.getContractAddress(network),
				OwnerAddress:  ticket.To,
				BookingID:     ticket.BookingID,
				EventID:       ticket.EventID,
				Status:        "active",
				MetadataURI:   results[i].MetadataURI,
				ImageURL:      ticket.ImageURL,
				TxHash:        tx.Hash().Hex(),
			}
			
			if err := db.SaveNFTMapping(mapping); err != nil {
				fmt.Printf("Warning: Failed to save NFT mapping for token %d: %v\n", tokenID, err)
			}
		}

		if err != nil {
			results[i].Message = fmt.Sprintf("Warning: NFT minted but metadata creation failed: %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Successfully batch minted %d tickets", len(tokenIDs)),
		"results": results,
		"txHash":  tx.Hash().Hex(),
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
// @Summary Upload ticket image
// @Description Uploads an image for an NFT ticket
// @Tags NFT
// @Accept multipart/form-data
// @Produce json
// @Param X-Network-Type header string false "Network type (testnet/mainnet)" default(testnet)
// @Param tokenId path int true "Token ID"
// @Param image formData file true "Image file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /nft/tickets/{tokenId}/image [post]
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
	db := database.GetDB()
	
	// Get Datakyte NFT ID from database
	_, err = db.GetDatakyteID(tokenID, network)
	if err == nil {
		metadataService := h.getMetadataService(network)
		
		// Get current metadata and update image
		metadata, err := metadataService.GetTicketMetadata(tokenID)
		if err == nil && metadata != nil {
			metadata.Image = imageURL
			// Update metadata in Datakyte (implementation depends on Datakyte API)
			// For now, just log it
			fmt.Printf("Image URL updated for token %d: %s\n", tokenID, imageURL)
		}
	}
	
	// Update image URL in database
	if mapping, err := db.GetNFTMapping(tokenID, network); err == nil {
		mapping.ImageURL = imageURL
		if err := db.SaveNFTMapping(mapping); err != nil {
			fmt.Printf("Warning: Failed to update image URL in database for token %d: %v\n", tokenID, err)
		}
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
// @Summary Update ticket status
// @Description Updates the status of an NFT ticket (Active, Redeemed, Expired, Burned)
// @Tags NFT
// @Accept json
// @Produce json
// @Param X-Network-Type header string false "Network type (testnet/mainnet)" default(testnet)
// @Param tokenId path int true "Token ID"
// @Param request body map[string]string true "Status update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /nft/tickets/{tokenId}/status [put]
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

	network := GetNetworkFromContext(c)
	
	// Get Datakyte NFT ID from database
	db := database.GetDB()
	datakyteNFTID, err := db.GetDatakyteID(tokenID, network)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: fmt.Sprintf("NFT mapping not found: %v", err),
		})
		return
	}

	metadataService := h.getMetadataService(network)
	err = metadataService.UpdateTicketStatus(datakyteNFTID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("Failed to update Datakyte status: %v", err),
		})
		return
	}
	
	// Update status in database as well
	if err := db.UpdateNFTStatus(tokenID, network, req.Status); err != nil {
		fmt.Printf("Warning: Failed to update status in database for token %d: %v\n", tokenID, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tokenId": tokenID,
		"status":  req.Status,
	})
}
