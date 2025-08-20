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

// RouterConfig contains all dependencies needed to create a router
type RouterConfig struct {
	SDK            SDKInterface           // Required: SDK for blockchain operations
	NetworkHandler *NetworkHandler        // Optional: for network switching support
	AppConfig      *config.Config         // Required: application configuration
	Storage        storage.RewardsStorage // Optional: defaults to in-memory storage
}

// CreateRouter creates a new Gin router with all routes configured
// This replaces both NewRouter and NewRouterWithNetworkSupport
func CreateRouter(cfg *RouterConfig) *gin.Engine {
	// Validate required fields
	if cfg.SDK == nil {
		log.Panic("SDK is required")
	}
	if cfg.AppConfig == nil {
		log.Panic("AppConfig is required")
	}

	// Set default storage if not provided
	if cfg.Storage == nil {
		cfg.Storage = storage.NewInMemoryRewardsStorage()
	}

	// Create handler with all dependencies
	handler := &Handler{
		SDK:            cfg.SDK,
		NetworkHandler: cfg.NetworkHandler,
		Config:         cfg.AppConfig,
		Storage:        cfg.Storage,
	}

	router := gin.New()

	// Setup middleware
	setupMiddleware(router)

	// Setup documentation routes
	setupDocumentationRoutes(router)

	// Setup API routes
	setupAPIRoutes(router, handler, cfg.AppConfig)

	return router
}

// setupMiddleware configures all middleware for the router
func setupMiddleware(router *gin.Engine) {
	// Basic middleware
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
		log.Printf("Warning: failed to set trusted proxies: %s", err)
		// Don't panic, just log the warning
	}

	// Rate limiting middleware
	limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100) // 100 requests per minute
	router.Use(rateLimitMiddleware(limiter))
}

// setupDocumentationRoutes configures documentation endpoints
func setupDocumentationRoutes(router *gin.Engine) {
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
}

// setupAPIRoutes configures all API endpoints
func setupAPIRoutes(router *gin.Engine, handler *Handler, cfg *config.Config) {
	api := router.Group("/api")

	// System endpoints
	api.GET("/health", handler.GetHealth)
	api.GET("/gas-price", handler.GetGasPrice)

	// Token endpoints
	setupTokenRoutes(api, handler)

	// Rewards endpoints
	setupRewardRoutes(api, handler, cfg)

	// NFT routes
	RegisterNFTRoutes(api, handler)

	// Register public NFT routes (no auth required)
	RegisterPublicNFTRoutes(router)
}

// setupTokenRoutes configures token-related endpoints
func setupTokenRoutes(api *gin.RouterGroup, handler *Handler) {
	token := api.Group("/token")
	token.GET("/balance/:address", handler.GetTokenBalance)
	token.POST("/transfer", handler.TransferBOGOTokens)
}

// setupRewardRoutes configures reward-related endpoints
func setupRewardRoutes(api *gin.RouterGroup, handler *Handler, cfg *config.Config) {
	rewards := api.Group("/rewards")

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.FirebaseProjectID)

	// Public reward endpoints
	rewards.GET("/templates", handler.GetRewardTemplates)
	rewards.GET("/templates/:id", handler.GetRewardTemplate)

	// Authenticated reward endpoints
	rewards.GET("/eligibility", AuthMiddleware(authMiddleware), handler.CheckRewardEligibility)
	rewards.GET("/history", AuthMiddleware(authMiddleware), handler.GetRewardHistory)

	// Main reward endpoints
	rewards.POST("/claim", AuthMiddleware(authMiddleware), handler.ClaimReward)
	rewards.POST("/claim-referral", AuthMiddleware(authMiddleware), handler.ClaimReferralBonus)
	rewards.POST("/claim-custom", handler.ClaimCustomReward)

	// Backward compatibility endpoint (DEPRECATED)
	rewards.POST("/claim-v2", AuthMiddleware(authMiddleware), handler.ClaimRewardV2)
}
