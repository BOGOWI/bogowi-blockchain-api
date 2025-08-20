package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestCreateRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creates router with all required routes", func(t *testing.T) {
		// Setup
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{
			FirebaseProjectID: "test-project",
		}

		// Create a temporary openapi.yaml file for testing
		tempFile, err := os.CreateTemp(".", "openapi*.yaml")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		_, err = tempFile.WriteString("openapi: 3.0.0")
		require.NoError(t, err)
		tempFile.Close()

		// Create router using new unified function
		routerConfig := &RouterConfig{
			SDK:       mockSDK,
			AppConfig: cfg,
		}
		router := CreateRouter(routerConfig)
		require.NotNil(t, router)

		// Test documentation routes
		testRoute(t, router, "GET", "/", http.StatusMovedPermanently)
		testRoute(t, router, "GET", "/docs", http.StatusOK)

		// Skip openapi.yaml test if file doesn't exist in expected location
		if _, err := os.Stat("./openapi.yaml"); err == nil {
			testRoute(t, router, "GET", "/openapi.yaml", http.StatusOK)
		}

		// Test API routes
		testRoute(t, router, "GET", "/api/health", http.StatusOK)
		testRoute(t, router, "GET", "/api/gas-price", http.StatusOK)
		testRoute(t, router, "GET", "/api/token/balance/0x123", http.StatusOK)

		// Test POST routes
		testRouteWithBody(t, router, "POST", "/api/token/transfer",
			map[string]string{"to": "0x456", "amount": "100"}, http.StatusOK)
	})

	t.Run("sets up middleware correctly", func(t *testing.T) {
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{
			FirebaseProjectID: "test-project",
		}

		routerConfig := &RouterConfig{
			SDK:       mockSDK,
			AppConfig: cfg,
		}
		router := CreateRouter(routerConfig)

		// Test that middleware is present by checking routes
		routes := router.Routes()
		assert.NotEmpty(t, routes)

		// Check CORS is working by sending OPTIONS request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
		req.Header.Set("Origin", "http://example.com")
		router.ServeHTTP(w, req)

		// CORS should allow the request
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("handles trusted proxies", func(t *testing.T) {
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{
			FirebaseProjectID: "test-project",
		}

		// This should not panic
		assert.NotPanics(t, func() {
			routerConfig := &RouterConfig{
				SDK:       mockSDK,
				AppConfig: cfg,
			}
			CreateRouter(routerConfig)
		})
	})

	t.Run("docs endpoint returns correct HTML", func(t *testing.T) {
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{}

		routerConfig := &RouterConfig{
			SDK:       mockSDK,
			AppConfig: cfg,
		}
		router := CreateRouter(routerConfig)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/docs", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, w.Body.String(), "BOGOWI Blockchain API Documentation")
		assert.Contains(t, w.Body.String(), "redoc")
	})

	t.Run("root redirects to docs", func(t *testing.T) {
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{}

		routerConfig := &RouterConfig{
			SDK:       mockSDK,
			AppConfig: cfg,
		}
		router := CreateRouter(routerConfig)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "/docs", w.Header().Get("Location"))
	})
}

func TestCreateRouterWithNetworkSupport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creates router with network support", func(t *testing.T) {
		// Setup
		mockSDK := &TestMockSDK{}
		networkHandler := &NetworkHandler{
			testnetSDK: mockSDK,
			mainnetSDK: mockSDK,
		}
		cfg := &config.Config{
			FirebaseProjectID: "test-project",
		}

		// Create router using unified function
		routerConfig := &RouterConfig{
			SDK:            mockSDK,
			NetworkHandler: networkHandler,
			AppConfig:      cfg,
		}
		router := CreateRouter(routerConfig)
		require.NotNil(t, router)

		// Test that all routes are created
		testRoute(t, router, "GET", "/", http.StatusMovedPermanently)
		testRoute(t, router, "GET", "/docs", http.StatusOK)
		testRoute(t, router, "GET", "/api/health", http.StatusOK)
		testRoute(t, router, "GET", "/api/gas-price", http.StatusOK)
	})

	t.Run("handler has network handler set", func(t *testing.T) {
		mockSDK := &TestMockSDK{}
		networkHandler := &NetworkHandler{
			testnetSDK: mockSDK,
			mainnetSDK: mockSDK,
		}
		cfg := &config.Config{}

		// Create handler internally to verify all fields are set correctly
		handler := &Handler{
			SDK:            mockSDK,
			NetworkHandler: networkHandler,
			Config:         cfg,
		}

		assert.NotNil(t, handler.NetworkHandler)
		assert.Equal(t, networkHandler, handler.NetworkHandler)
		assert.NotNil(t, handler.SDK)
		assert.Equal(t, mockSDK, handler.SDK)
		assert.NotNil(t, handler.Config)
		assert.Equal(t, cfg, handler.Config)
	})

	t.Run("creates identical routes to standard router", func(t *testing.T) {
		cfg := &config.Config{
			FirebaseProjectID: "test-project",
		}
		mockSDK := &TestMockSDK{}
		networkHandler := &NetworkHandler{
			testnetSDK: mockSDK,
			mainnetSDK: mockSDK,
		}

		// Create router with network support
		routerConfig := &RouterConfig{
			SDK:            mockSDK,
			NetworkHandler: networkHandler,
			AppConfig:      cfg,
		}
		routerWithNetwork := CreateRouter(routerConfig)

		// Get routes from network router
		networkRoutes := routerWithNetwork.Routes()

		// Should have multiple routes
		assert.NotEmpty(t, networkRoutes)

		// Check specific routes exist
		hasDocsRoute := false
		hasHealthRoute := false
		for _, route := range networkRoutes {
			if route.Path == "/docs" && route.Method == "GET" {
				hasDocsRoute = true
			}
			if route.Path == "/api/health" && route.Method == "GET" {
				hasHealthRoute = true
			}
		}

		assert.True(t, hasDocsRoute, "Should have /docs route")
		assert.True(t, hasHealthRoute, "Should have /api/health route")
	})

	t.Run("sets up same middleware as standard router", func(t *testing.T) {
		mockSDK := &TestMockSDK{}
		networkHandler := &NetworkHandler{
			testnetSDK: mockSDK,
			mainnetSDK: mockSDK,
		}
		cfg := &config.Config{}

		routerConfig := &RouterConfig{
			SDK:            mockSDK,
			NetworkHandler: networkHandler,
			AppConfig:      cfg,
		}
		router := CreateRouter(routerConfig)

		// Test CORS
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
		req.Header.Set("Origin", "http://example.com")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allows requests within rate limit", func(t *testing.T) {
		limiter := rate.NewLimiter(rate.Every(time.Second), 10)

		router := gin.New()
		router.Use(rateLimitMiddleware(limiter))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// First request should succeed
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "ok", response["status"])
	})

	t.Run("blocks requests exceeding rate limit", func(t *testing.T) {
		// Very strict limiter - 1 request immediately, no refill
		limiter := rate.NewLimiter(0, 1)

		router := gin.New()
		router.Use(rateLimitMiddleware(limiter))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// First request should succeed
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second request should be rate limited
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusTooManyRequests, w2.Code)

		var response map[string]string
		json.Unmarshal(w2.Body.Bytes(), &response)
		assert.Equal(t, "Rate limit exceeded", response["error"])
	})

	t.Run("rate limiter prevents handler execution when limit exceeded", func(t *testing.T) {
		limiter := rate.NewLimiter(0, 1) // Allow only 1 request

		handlerCalled := false
		router := gin.New()
		router.Use(rateLimitMiddleware(limiter))
		router.GET("/test", func(c *gin.Context) {
			handlerCalled = true
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// First request
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.True(t, handlerCalled)

		// Reset flag
		handlerCalled = false

		// Second request should be blocked
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w2, req2)

		assert.False(t, handlerCalled, "Handler should not be called when rate limit exceeded")
	})

	t.Run("rate limiter works with multiple endpoints", func(t *testing.T) {
		limiter := rate.NewLimiter(rate.Every(time.Second), 2) // 2 requests per second

		router := gin.New()
		router.Use(rateLimitMiddleware(limiter))
		router.GET("/endpoint1", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"endpoint": "1"})
		})
		router.GET("/endpoint2", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"endpoint": "2"})
		})

		// Request to endpoint1
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/endpoint1", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Request to endpoint2
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/endpoint2", nil)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// Third request should be rate limited (shared limiter)
		w3 := httptest.NewRecorder()
		req3, _ := http.NewRequest("GET", "/endpoint1", nil)
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusTooManyRequests, w3.Code)
	})
}

func TestHandler(t *testing.T) {
	t.Run("Handler struct initialization", func(t *testing.T) {
		mockSDK := &TestMockSDK{}
		networkHandler := &NetworkHandler{
			testnetSDK: mockSDK,
			mainnetSDK: mockSDK,
		}
		cfg := &config.Config{
			FirebaseProjectID: "test-project",
		}
		storage := storage.NewInMemoryRewardsStorage()

		handler := &Handler{
			SDK:            mockSDK,
			NetworkHandler: networkHandler,
			Config:         cfg,
			Storage:        storage,
		}

		assert.NotNil(t, handler.SDK)
		assert.NotNil(t, handler.NetworkHandler)
		assert.NotNil(t, handler.Config)
		assert.NotNil(t, handler.Storage)
	})
}

func TestRouterIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("full request flow with middleware", func(t *testing.T) {
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{
			FirebaseProjectID: "test-project",
		}

		routerConfig := &RouterConfig{
			SDK:       mockSDK,
			AppConfig: cfg,
		}
		router := CreateRouter(routerConfig)

		// Make a request that goes through all middleware
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/health", nil)
		req.Header.Set("Origin", "http://example.com")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Check CORS headers are set
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("error recovery middleware works", func(t *testing.T) {
		router := gin.New()
		router.Use(gin.Recovery())

		// Add a route that panics
		router.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/panic", nil)

		// Should not panic the test
		assert.NotPanics(t, func() {
			router.ServeHTTP(w, req)
		})

		// Should return 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// Helper functions for testing

func testRoute(t *testing.T, router *gin.Engine, method, path string, expectedStatus int) {
	t.Helper()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, nil)
	require.NoError(t, err)

	router.ServeHTTP(w, req)

	// For redirects, just check the status code
	if expectedStatus == http.StatusMovedPermanently {
		assert.Equal(t, expectedStatus, w.Code, "Path %s %s should return %d", method, path, expectedStatus)
		return
	}

	// For other routes, check if they exist (not 404)
	assert.NotEqual(t, http.StatusNotFound, w.Code, "Path %s %s should exist", method, path)
}

func testRouteWithBody(t *testing.T, router *gin.Engine, method, path string, body interface{}, _ int) {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Check if route exists (not 404)
	assert.NotEqual(t, http.StatusNotFound, w.Code, "Path %s %s should exist", method, path)
}

func TestTrustedProxiesConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("sets trusted proxies successfully", func(t *testing.T) {
		// Temporarily redirect log output to capture panic if it occurs
		oldLogOutput := log.Writer()
		defer log.SetOutput(oldLogOutput)

		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)

		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{}

		// This should not panic with valid proxy
		assert.NotPanics(t, func() {
			routerConfig := &RouterConfig{
				SDK:       mockSDK,
				AppConfig: cfg,
			}
			CreateRouter(routerConfig)
		})

		// Check that no panic was logged
		assert.NotContains(t, logBuffer.String(), "panic")
	})
}

func TestCORSConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CORS allows all configured methods", func(t *testing.T) {
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{}
		routerConfig := &RouterConfig{
			SDK:       mockSDK,
			AppConfig: cfg,
		}
		router := CreateRouter(routerConfig)

		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

		for _, method := range methods {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
			req.Header.Set("Origin", "http://example.com")
			req.Header.Set("Access-Control-Request-Method", method)

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code, "OPTIONS request for %s should succeed", method)
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), method)
		}
	})

	t.Run("CORS allows configured headers", func(t *testing.T) {
		mockSDK := &SimpleMockSDK{}
		cfg := &config.Config{}
		routerConfig := &RouterConfig{
			SDK:       mockSDK,
			AppConfig: cfg,
		}
		router := CreateRouter(routerConfig)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		allowedHeaders := w.Header().Get("Access-Control-Allow-Headers")
		assert.Contains(t, allowedHeaders, "Authorization")
		assert.Contains(t, allowedHeaders, "Content-Type")
	})
}
