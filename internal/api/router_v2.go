package api

import (
	"log"
	"net/http"
	"time"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/middleware"
	"bogowi-blockchain-go/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// HandlerV2 holds the network handler and configuration
type HandlerV2 struct {
	NetworkHandler *NetworkHandler
	Config         *config.Config
	Storage        storage.RewardsStorage
}

// NetworkMiddleware extracts and validates the network parameter
func NetworkMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		network := c.Query("network")
		if network == "" {
			network = c.GetHeader("X-Network")
		}
		if network == "" {
			network = "testnet" // Default to testnet
		}

		// Validate network
		if network != "testnet" && network != "mainnet" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid network. Use 'testnet' or 'mainnet'"})
			c.Abort()
			return
		}

		// Store in context
		c.Set("network", network)
		c.Next()
	}
}

// GetNetworkFromContext retrieves the network from gin context
func GetNetworkFromContext(c *gin.Context) string {
	network, _ := c.Get("network")
	return network.(string)
}

// NewRouterV2 creates a new Gin router with network support
func NewRouterV2(cfg *config.Config) (*gin.Engine, error) {
	// Create network handler
	networkHandler, err := NewNetworkHandler(cfg)
	if err != nil {
		return nil, err
	}

	handler := &HandlerV2{
		NetworkHandler: networkHandler,
		Config:         cfg,
		Storage:        storage.NewInMemoryRewardsStorage(),
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Network"}
	router.Use(cors.New(corsConfig))

	// Trust proxy for forwarded headers (required for nginx)
	trustedProxies := []string{"127.0.0.1"}
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		log.Panicf("failed to set trusted proxies: %s", err)
	}

	// Rate limiting middleware
	limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100) // 100 requests per minute
	router.Use(rateLimitMiddleware(limiter))

	// Network middleware - applies to all API routes
	router.Use(NetworkMiddleware())

	// Redoc documentation endpoint
	router.GET("/docs", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>BOGOWI Blockchain API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <redoc spec-url='/openapi.yaml'></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@2.1.5/bundles/redoc.standalone.js"></script>
</body>
</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// Serve OpenAPI spec
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("./openapi.yaml")
	})

	// Root redirect to docs
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})

	// API routes
	api := router.Group("/api")

	// System endpoints
	api.GET("/health", handler.GetHealthV2)
	api.GET("/gas-price", handler.GetGasPriceV2)
	api.GET("/network-info", handler.GetNetworkInfo)

	// Token endpoints
	token := api.Group("/token")
	token.GET("/balance/:address", handler.GetTokenBalanceV2)
	token.POST("/transfer", handler.TransferBOGOTokensV2)

	// Rewards endpoints
	rewards := api.Group("/rewards")

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.FirebaseProjectID)

	// Public reward endpoints
	rewards.GET("/templates", handler.GetRewardTemplatesV2)
	rewards.GET("/templates/:id", handler.GetRewardTemplateV2)

	// Authenticated reward endpoints
	rewards.GET("/eligibility", AuthMiddleware(authMiddleware), handler.CheckRewardEligibilityV2)
	rewards.GET("/history", AuthMiddleware(authMiddleware), handler.GetRewardHistoryV2)
	rewards.POST("/claim-v2", AuthMiddleware(authMiddleware), handler.ClaimRewardV3)
	rewards.POST("/claim-referral", AuthMiddleware(authMiddleware), handler.ClaimReferralV3)

	// Backend-only endpoint
	rewards.POST("/claim-custom", handler.ClaimCustomRewardV3)

	return router, nil
}
