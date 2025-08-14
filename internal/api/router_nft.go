package api

import (
	"github.com/gin-gonic/gin"
)

// RegisterNFTRoutes registers all NFT-related routes
func RegisterNFTRoutes(router *gin.RouterGroup, handler *HandlerV2) {
	// Create NFT handler
	nftHandler := NewNFTHandler(handler)
	
	// NFT ticket endpoints
	nft := router.Group("/nft")
	{
		// Ticket management
		tickets := nft.Group("/tickets")
		{
			// Minting
			tickets.POST("/mint", nftHandler.MintTicket)
			tickets.POST("/batch-mint", nftHandler.BatchMintTickets)
			
			// Images
			tickets.GET("/:tokenId/upload-url", nftHandler.GetPresignedUploadURL)
			tickets.POST("/:tokenId/image", nftHandler.UploadTicketImage)
			
			// Metadata
			tickets.GET("/:tokenId/metadata", nftHandler.GetTicketMetadata)
			tickets.PATCH("/:tokenId/status", nftHandler.UpdateTicketStatus)
			
			// Redemption
			tickets.POST("/:tokenId/redeem", nftHandler.RedeemTicket)
			
			// User queries
			tickets.GET("/user/:address", nftHandler.GetUserTickets)
		}
		
		// Public metadata endpoint (ERC-721 compliant)
		// This should be accessible without authentication for marketplaces
		nft.GET("/metadata/:contractAddress/:tokenId", func(c *gin.Context) {
			// Redirect to Datakyte endpoint
			contractAddress := c.Param("contractAddress")
			tokenId := c.Param("tokenId")
			c.Redirect(301, "https://dklnk.to/api/nfts/"+contractAddress+"/"+tokenId+"/metadata")
		})
	}
}

// RegisterPublicNFTRoutes registers public NFT routes (no auth required)
func RegisterPublicNFTRoutes(router *gin.Engine) {
	// OpenSea/marketplace compatible metadata endpoint
	router.GET("/api/nft/:contractAddress/:tokenId", func(c *gin.Context) {
		contractAddress := c.Param("contractAddress")
		tokenId := c.Param("tokenId")
		
		// Redirect to Datakyte's public endpoint
		c.Redirect(301, "https://dklnk.to/api/nfts/"+contractAddress+"/"+tokenId+"/metadata")
	})
	
	// Contract-level metadata (for collection info)
	router.GET("/api/nft/contract/:contractAddress", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "BOGOWI Eco-Experience Tickets",
			"description": "NFT tickets for sustainable travel experiences that support conservation efforts worldwide",
			"image": "https://storage.bogowi.com/collection/bogowi-tickets.png",
			"external_link": "https://bogowi.com",
			"seller_fee_basis_points": 500, // 5% royalty
			"fee_recipient": "0x...", // Conservation DAO address
		})
	})
}