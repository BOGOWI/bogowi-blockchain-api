package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/sdk"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test config
	cfg := &config.Config{
		Environment:      "test",
		APIPort:          "3001",
		BackendSecret:    "test-secret",
		DevBackendSecret: "test-dev-secret",
		Testnet: config.NetworkConfig{
			RPCUrl:    "https://columbus.camino.network/ext/bc/C/rpc",
			ChainID:   501,
			Contracts: config.ContractAddresses{},
		},
		Mainnet: config.NetworkConfig{
			RPCUrl:    "https://api.camino.network/ext/bc/C/rpc",
			ChainID:   500,
			Contracts: config.ContractAddresses{},
		},
	}

	// Create test handler
	handler := &Handler{
		SDK:    &sdk.BOGOWISDK{},
		Config: cfg,
	}

	router := gin.New()
	router.GET("/health", handler.GetHealth)

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
}
