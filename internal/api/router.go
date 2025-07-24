package api

import (
	"log"
	"net/http"
	"time"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/sdk"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Handler holds the SDK and configuration
type Handler struct {
	SDK    *sdk.BOGOWISDK
	Config *config.Config
}

// NewRouter creates a new Gin router with all routes configured
func NewRouter(bogoSDK *sdk.BOGOWISDK, cfg *config.Config) *gin.Engine {
	handler := &Handler{
		SDK:    bogoSDK,
		Config: cfg,
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Trust proxy for forwarded headers (required for nginx)
	trustedProxies := []string{"127.0.0.1"}
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		log.Panicf("failed to set trusted proxies: %s", err)
	}

	// Rate limiting middleware
	limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100) // 100 requests per minute
	router.Use(rateLimitMiddleware(limiter))

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
	api.GET("/health", handler.GetHealth)
	api.GET("/gas-price", handler.GetGasPrice)

	// Token endpoints
	token := api.Group("/token")
	token.GET("/balance/:address", handler.GetTokenBalance)
	token.GET("/flavored-balances/:address", handler.GetFlavoredTokenBalances)

	// NFT endpoints
	nft := api.Group("/nft")
	nft.GET("/balance/:address/:tokenId", handler.GetNFTBalance)
	nft.POST("/mint-ticket", handler.MintEventTicket)
	nft.POST("/mint-collectible", handler.MintConservationNFT)

	// Rewards endpoints
	rewards := api.Group("/rewards")
	rewards.GET("/info/:address", handler.GetRewardInfo)
	rewards.GET("/achievement/:address/:achievementId", handler.GetAchievementProgress)
	rewards.POST("/claim", handler.ClaimReward)

	// DAO endpoints
	dao := api.Group("/dao")
	dao.GET("/info", handler.GetDAOInfo)
	dao.GET("/pending-transactions", handler.GetPendingTransactions)

	return router
}

// rateLimitMiddleware implements rate limiting
func rateLimitMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
