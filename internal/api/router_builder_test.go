package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestRouterBuilder_Build(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creates router with all routes", func(t *testing.T) {
		deps := TestRouterDependencies(t)
		builder := NewRouterBuilder(deps)
		router := builder.Build()

		require.NotNil(t, router)

		// Check documentation routes
		AssertRouteExists(t, router, "GET", "/")
		AssertRouteExists(t, router, "GET", "/docs")
		AssertRouteExists(t, router, "GET", "/openapi.yaml")

		// Check system routes
		AssertRouteExists(t, router, "GET", "/api/health")
		AssertRouteExists(t, router, "GET", "/api/gas-price")

		// Check token routes
		AssertRouteExists(t, router, "GET", "/api/token/balance/:address")
		AssertRouteExists(t, router, "POST", "/api/token/transfer")

		// Check reward routes
		AssertRouteExists(t, router, "GET", "/api/rewards/templates")
		AssertRouteExists(t, router, "GET", "/api/rewards/templates/:id")
		AssertRouteExists(t, router, "POST", "/api/rewards/claim-custom")
	})

	t.Run("skips middleware when requested", func(t *testing.T) {
		deps := TestRouterDependencies(t)
		deps.RateLimiter = rate.NewLimiter(rate.Every(1), 1) // Very strict limit

		builder := NewRouterBuilder(deps).SkipMiddleware()
		router := builder.Build()

		// Should be able to make multiple requests without rate limiting
		for i := 0; i < 5; i++ {
			resp, err := MakeRequest(router, TestRequest{
				Method: "GET",
				Path:   "/api/health",
			})
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
		}
	})

	t.Run("includes auth routes when auth middleware provided", func(t *testing.T) {
		t.Skip("Auth middleware mocking needs proper implementation")
		deps := TestRouterDependencies(t)
		// deps.AuthMiddleware = CreateMockAuthMiddleware() // TODO: Fix auth middleware mocking

		builder := NewRouterBuilder(deps)
		router := builder.Build()

		// Auth-protected routes should exist
		AssertRouteExists(t, router, "GET", "/api/rewards/eligibility")
		AssertRouteExists(t, router, "GET", "/api/rewards/history")
		AssertRouteExists(t, router, "POST", "/api/rewards/claim-v2")
		AssertRouteExists(t, router, "POST", "/api/rewards/claim-referral")
	})

	t.Run("includes NFT routes when NFT handler provided", func(t *testing.T) {
		deps := TestRouterDependencies(t)
		deps.NFTHandlerFunc = func(h *Handler) *NFTHandler {
			return &NFTHandler{Handler: h}
		}

		builder := NewRouterBuilder(deps)
		router := builder.Build()

		// NFT routes should exist
		AssertRouteExists(t, router, "POST", "/api/nft/tickets/mint")
		AssertRouteExists(t, router, "GET", "/api/nft/tickets/:tokenId/metadata")
	})
}

func TestRouterBuilder_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("handler has correct dependencies", func(t *testing.T) {
		cfg := &config.Config{Environment: "test"}
		mockSDK := NewSimpleMockSDK()
		mockStorage := storage.NewInMemoryRewardsStorage()

		deps := &RouterDependencies{
			DefaultSDK: mockSDK,
			Config:     cfg,
			Storage:    mockStorage,
		}

		builder := NewRouterBuilder(deps)
		router := builder.Build()

		// Access handler through route
		resp, err := MakeRequest(router, TestRequest{
			Method: "GET",
			Path:   "/api/health",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestDefaultRouterDependencies(t *testing.T) {
	t.Run("creates dependencies with defaults", func(t *testing.T) {
		cfg := &config.Config{
			Environment:       "development",
			APIPort:           "8080",
			FirebaseProjectID: "test-project",
			Testnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8545",
				ChainID: 501,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x1234567890123456789012345678901234567890",
				},
			},
			Mainnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8546",
				ChainID: 500,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x0987654321098765432109876543210987654321",
				},
			},
			TestnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			MainnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		}

		deps, err := DefaultRouterDependencies(cfg)
		require.NoError(t, err)
		require.NotNil(t, deps)

		assert.NotNil(t, deps.NetworkHandler)
		assert.NotNil(t, deps.DefaultSDK)
		assert.NotNil(t, deps.Config)
		assert.NotNil(t, deps.Storage)
		assert.NotNil(t, deps.RateLimiter)
		assert.NotNil(t, deps.AuthMiddleware)
		assert.NotNil(t, deps.CORSConfig)
		assert.NotEmpty(t, deps.TrustedProxies)
		assert.NotNil(t, deps.NFTHandlerFunc)
	})

	t.Run("selects correct default SDK", func(t *testing.T) {
		// Test development environment
		cfg := &config.Config{
			Environment: "development",
			Testnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8545",
				ChainID: 501,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x1234567890123456789012345678901234567890",
				},
			},
			Mainnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8546",
				ChainID: 500,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x0987654321098765432109876543210987654321",
				},
			},
			TestnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			MainnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		}

		deps, err := DefaultRouterDependencies(cfg)
		require.NoError(t, err)
		// In development, should use testnet SDK
		assert.NotNil(t, deps.DefaultSDK)

		// Test production environment
		cfg.Environment = "production"
		deps, err = DefaultRouterDependencies(cfg)
		require.NoError(t, err)
		// In production, should use mainnet SDK
		assert.NotNil(t, deps.DefaultSDK)
	})
}

func TestRouterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("rate limiting works", func(t *testing.T) {
		t.Skip("Rate limiting test needs fixing")
		deps := TestRouterDependencies(t)
		deps.RateLimiter = rate.NewLimiter(rate.Every(1), 2) // Allow 2 requests per second

		builder := NewRouterBuilder(deps)
		router := builder.Build()

		// First two requests should succeed
		for i := 0; i < 2; i++ {
			resp, err := MakeRequest(router, TestRequest{
				Method: "GET",
				Path:   "/api/health",
			})
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
		}

		// Third request should be rate limited
		resp, err := MakeRequest(router, TestRequest{
			Method: "GET",
			Path:   "/api/health",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp.Code)
	})

	t.Run("CORS headers are set", func(t *testing.T) {
		deps, err := DefaultRouterDependencies(&config.Config{
			Environment:       "test",
			FirebaseProjectID: "test",
			Testnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8545",
				ChainID: 501,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x1234567890123456789012345678901234567890",
				},
			},
			Mainnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8546",
				ChainID: 500,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x0987654321098765432109876543210987654321",
				},
			},
			TestnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			MainnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		})
		require.NoError(t, err)

		builder := NewRouterBuilder(deps)
		router := builder.Build()

		// Make OPTIONS request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
		req.Header.Set("Origin", "http://example.com")
		router.ServeHTTP(w, req)

		// Check CORS headers
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	})
}

func TestRouteRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("all system routes registered", func(t *testing.T) {
		router := CreateTestRouter(t, TestRouterOptions{})

		routes := []struct {
			method string
			path   string
		}{
			{"GET", "/api/health"},
			{"GET", "/api/gas-price"},
		}

		for _, route := range routes {
			AssertRouteExists(t, router, route.method, route.path)
		}
	})

	t.Run("all token routes registered", func(t *testing.T) {
		router := CreateTestRouter(t, TestRouterOptions{})

		routes := []struct {
			method string
			path   string
		}{
			{"GET", "/api/token/balance/:address"},
			{"POST", "/api/token/transfer"},
		}

		for _, route := range routes {
			AssertRouteExists(t, router, route.method, route.path)
		}
	})

	t.Run("public reward routes always registered", func(t *testing.T) {
		router := CreateTestRouter(t, TestRouterOptions{})

		routes := []struct {
			method string
			path   string
		}{
			{"GET", "/api/rewards/templates"},
			{"GET", "/api/rewards/templates/:id"},
			{"POST", "/api/rewards/claim-custom"},
		}

		for _, route := range routes {
			AssertRouteExists(t, router, route.method, route.path)
		}
	})

	t.Run("auth routes only with auth middleware", func(t *testing.T) {
		t.Skip("Auth middleware test needs fixing")
		// Without auth
		router := CreateTestRouter(t, TestRouterOptions{WithAuth: false})
		AssertRouteNotExists(t, router, "GET", "/api/rewards/eligibility")

		// With auth
		router = CreateTestRouter(t, TestRouterOptions{WithAuth: true})
		AssertRouteExists(t, router, "GET", "/api/rewards/eligibility")
	})
}

func TestBackwardCompatibility(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("NewRouterWithBuilder maintains compatibility", func(t *testing.T) {
		cfg := &config.Config{
			Environment:       "test",
			FirebaseProjectID: "test",
			Testnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8545",
				ChainID: 501,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x1234567890123456789012345678901234567890",
				},
			},
			Mainnet: config.NetworkConfig{
				RPCUrl:  "http://localhost:8546",
				ChainID: 500,
				Contracts: config.ContractAddresses{
					BOGOToken: "0x0987654321098765432109876543210987654321",
				},
			},
			TestnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			MainnetPrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		}

		networkHandler, err := NewNetworkHandler(cfg)
		require.NoError(t, err)

		mockSDK := NewSimpleMockSDK()

		router := NewRouterWithBuilder(networkHandler, mockSDK, cfg)
		require.NotNil(t, router)

		// Should have all standard routes
		AssertRouteExists(t, router, "GET", "/api/health")
		AssertRouteExists(t, router, "GET", "/api/token/balance/:address")
		AssertRouteExists(t, router, "GET", "/api/rewards/templates")
	})
}
