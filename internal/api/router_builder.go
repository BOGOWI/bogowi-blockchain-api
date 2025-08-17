package api

import (
	"net/http"
	"time"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/middleware"
	"bogowi-blockchain-go/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RouterDependencies holds all dependencies needed to build a router
type RouterDependencies struct {
	NetworkHandler *NetworkHandler
	DefaultSDK     SDKInterface
	Config         *config.Config
	Storage        storage.RewardsStorage
	RateLimiter    *rate.Limiter
	AuthMiddleware *middleware.AuthMiddleware
	CORSConfig     *cors.Config
	TrustedProxies []string
	NFTHandlerFunc func(*Handler) *NFTHandler // For dependency injection
}

// DefaultRouterDependencies creates dependencies with sensible defaults
func DefaultRouterDependencies(cfg *config.Config) (*RouterDependencies, error) {
	// Initialize NetworkHandler
	networkHandler, err := NewNetworkHandler(cfg)
	if err != nil {
		return nil, err
	}

	// Get default SDK based on environment
	var defaultSDK SDKInterface
	if cfg.Environment == "development" {
		defaultSDK, _ = networkHandler.GetSDK("testnet")
	} else {
		defaultSDK, _ = networkHandler.GetSDK("mainnet")
	}

	// Create CORS config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	return &RouterDependencies{
		NetworkHandler: networkHandler,
		DefaultSDK:     defaultSDK,
		Config:         cfg,
		Storage:        storage.NewInMemoryRewardsStorage(),
		RateLimiter:    rate.NewLimiter(rate.Every(time.Minute/100), 100),
		AuthMiddleware: middleware.NewAuthMiddleware(cfg.FirebaseProjectID),
		CORSConfig:     &corsConfig,
		TrustedProxies: []string{"127.0.0.1"},
		NFTHandlerFunc: NewNFTHandler,
	}, nil
}

// RouterBuilder helps build routers with testable dependencies
type RouterBuilder struct {
	deps           *RouterDependencies
	engine         *gin.Engine
	handler        *Handler
	skipMiddleware bool // For testing
}

// NewRouterBuilder creates a new router builder
func NewRouterBuilder(deps *RouterDependencies) *RouterBuilder {
	return &RouterBuilder{
		deps:   deps,
		engine: gin.New(),
	}
}

// SkipMiddleware disables middleware for testing
func (rb *RouterBuilder) SkipMiddleware() *RouterBuilder {
	rb.skipMiddleware = true
	return rb
}

// Build constructs the router
func (rb *RouterBuilder) Build() *gin.Engine {
	// Create handler
	rb.handler = &Handler{
		SDK:            rb.deps.DefaultSDK,
		NetworkHandler: rb.deps.NetworkHandler,
		Config:         rb.deps.Config,
		Storage:        rb.deps.Storage,
	}

	// Apply middleware unless skipped (for testing)
	if !rb.skipMiddleware {
		rb.applyMiddleware()
	}

	// Register routes
	rb.registerRoutes()

	return rb.engine
}

// applyMiddleware sets up all middleware
func (rb *RouterBuilder) applyMiddleware() {
	// Basic middleware
	rb.engine.Use(gin.Logger())
	rb.engine.Use(gin.Recovery())

	// CORS
	if rb.deps.CORSConfig != nil {
		rb.engine.Use(cors.New(*rb.deps.CORSConfig))
	}

	// Trusted proxies
	if len(rb.deps.TrustedProxies) > 0 {
		_ = rb.engine.SetTrustedProxies(rb.deps.TrustedProxies)
		// Note: In production, handle this error properly
	}

	// Rate limiting
	if rb.deps.RateLimiter != nil {
		rb.engine.Use(rateLimitMiddleware(rb.deps.RateLimiter))
	}
}

// registerRoutes sets up all API routes
func (rb *RouterBuilder) registerRoutes() {
	// Documentation routes
	rb.registerDocsRoutes()

	// API routes
	api := rb.engine.Group("/api")

	// System endpoints
	rb.registerSystemRoutes(api)

	// Token endpoints
	rb.registerTokenRoutes(api)

	// Rewards endpoints
	rb.registerRewardRoutes(api)

	// NFT routes
	if rb.deps.NFTHandlerFunc != nil {
		RegisterNFTRoutes(api, rb.handler)
		RegisterPublicNFTRoutes(rb.engine)
	}
}

// registerDocsRoutes sets up documentation endpoints
func (rb *RouterBuilder) registerDocsRoutes() {
	// Redoc documentation
	rb.engine.GET("/docs", func(c *gin.Context) {
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

	// OpenAPI spec
	rb.engine.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("./openapi.yaml")
	})

	// Root redirect
	rb.engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})
}

// registerSystemRoutes sets up system endpoints
func (rb *RouterBuilder) registerSystemRoutes(api *gin.RouterGroup) {
	api.GET("/health", rb.handler.GetHealth)
	api.GET("/gas-price", rb.handler.GetGasPrice)
}

// registerTokenRoutes sets up token endpoints
func (rb *RouterBuilder) registerTokenRoutes(api *gin.RouterGroup) {
	token := api.Group("/token")
	token.GET("/balance/:address", rb.handler.GetTokenBalance)
	token.POST("/transfer", rb.handler.TransferBOGOTokens)
}

// registerRewardRoutes sets up reward endpoints
func (rb *RouterBuilder) registerRewardRoutes(api *gin.RouterGroup) {
	rewards := api.Group("/rewards")

	// Public endpoints
	rewards.GET("/templates", rb.handler.GetRewardTemplates)
	rewards.GET("/templates/:id", rb.handler.GetRewardTemplate)

	// Authenticated endpoints
	if rb.deps.AuthMiddleware != nil {
		auth := AuthMiddleware(rb.deps.AuthMiddleware)
		rewards.GET("/eligibility", auth, rb.handler.CheckRewardEligibility)
		rewards.GET("/history", auth, rb.handler.GetRewardHistory)
		rewards.POST("/claim", auth, rb.handler.ClaimReward)
		rewards.POST("/claim-v2", auth, rb.handler.ClaimRewardV2) // Backward compatibility
		rewards.POST("/claim-referral", auth, rb.handler.ClaimReferralBonus)
	}

	// Backend-only endpoint
	rewards.POST("/claim-custom", rb.handler.ClaimCustomReward)
}

// NewRouterWithBuilder creates a router using the builder pattern (backward compatible)
func NewRouterWithBuilder(networkHandler *NetworkHandler, defaultSDK SDKInterface, cfg *config.Config) *gin.Engine {
	deps := &RouterDependencies{
		NetworkHandler: networkHandler,
		DefaultSDK:     defaultSDK,
		Config:         cfg,
		Storage:        storage.NewInMemoryRewardsStorage(),
		RateLimiter:    rate.NewLimiter(rate.Every(time.Minute/100), 100),
		AuthMiddleware: middleware.NewAuthMiddleware(cfg.FirebaseProjectID),
		TrustedProxies: []string{"127.0.0.1"},
		NFTHandlerFunc: NewNFTHandler,
	}

	// Create CORS config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	deps.CORSConfig = &corsConfig

	builder := NewRouterBuilder(deps)
	return builder.Build()
}
