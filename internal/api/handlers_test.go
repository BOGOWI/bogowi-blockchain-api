package api

import (
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
		Environment: "test",
		APIPort:     "3001",
		Contracts:   config.ContractAddresses{},
	}

	// Create test handler
	handler := &Handler{
		SDK:    &sdk.BOGOWISDK{},
		Config: cfg,
	}

	router := gin.New()
	router.GET("/health", handler.GetHealth)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
}

func TestNFTBalanceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test handler
	handler := &Handler{
		SDK:    &sdk.BOGOWISDK{},
		Config: &config.Config{},
	}

	router := gin.New()
	router.GET("/nft/balance/:address/:tokenId", handler.GetNFTBalance)

	req, _ := http.NewRequest("GET", "/nft/balance/0x742d35Cc6634C0532925a3b8D84d9C74D938f1f1/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "balance")
}
