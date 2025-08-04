package api

import (
	"fmt"
	"sync"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/sdk"
)

// NetworkHandler manages SDK instances for both testnet and mainnet
type NetworkHandler struct {
	testnetSDK  SDKInterface
	mainnetSDK  SDKInterface
	config      *config.Config
	mu          sync.RWMutex
}

// NewNetworkHandler creates a new network-aware handler
func NewNetworkHandler(cfg *config.Config) (*NetworkHandler, error) {
	handler := &NetworkHandler{
		config: cfg,
	}

	// Initialize testnet SDK
	if cfg.Testnet.Contracts.BOGOToken != "" || cfg.Testnet.Contracts.RewardDistributor != "" {
		testnetSDK, err := sdk.NewBOGOWISDK(&cfg.Testnet, cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize testnet SDK: %w", err)
		}
		handler.testnetSDK = testnetSDK
	}

	// Initialize mainnet SDK
	if cfg.Mainnet.Contracts.BOGOToken != "" || cfg.Mainnet.Contracts.RewardDistributor != "" {
		mainnetSDK, err := sdk.NewBOGOWISDK(&cfg.Mainnet, cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mainnet SDK: %w", err)
		}
		handler.mainnetSDK = mainnetSDK
	}

	return handler, nil
}

// GetSDK returns the appropriate SDK based on the network parameter
func (h *NetworkHandler) GetSDK(network string) (SDKInterface, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	switch network {
	case "testnet", "columbus":
		if h.testnetSDK == nil {
			return nil, fmt.Errorf("testnet SDK not initialized")
		}
		return h.testnetSDK, nil
	case "mainnet", "camino":
		if h.mainnetSDK == nil {
			return nil, fmt.Errorf("mainnet SDK not initialized")
		}
		return h.mainnetSDK, nil
	default:
		return nil, fmt.Errorf("invalid network: %s (use 'testnet' or 'mainnet')", network)
	}
}

// Close closes all SDK connections
func (h *NetworkHandler) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.testnetSDK != nil {
		h.testnetSDK.Close()
	}
	if h.mainnetSDK != nil {
		h.mainnetSDK.Close()
	}
}