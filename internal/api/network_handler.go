package api

import (
	"fmt"
	"sync"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/sdk"
	"bogowi-blockchain-go/internal/sdk/nft"
)

// NetworkHandler manages SDK instances for both testnet and mainnet
type NetworkHandler struct {
	testnetSDK    SDKInterface
	mainnetSDK    SDKInterface
	testnetNFTSDK *nft.Client
	mainnetNFTSDK *nft.Client
	config        *config.Config
	mu            sync.RWMutex
}

// NewNetworkHandler creates a new network-aware handler
func NewNetworkHandler(cfg *config.Config) (*NetworkHandler, error) {
	handler := &NetworkHandler{
		config: cfg,
	}

	// Initialize testnet SDK
	if cfg.Testnet.Contracts.BOGOToken != "" || cfg.Testnet.Contracts.RewardDistributor != "" {
		testnetPrivateKey := cfg.TestnetPrivateKey
		if testnetPrivateKey == "" {
			return nil, fmt.Errorf("TESTNET_PRIVATE_KEY is required for testnet operations")
		}
		testnetSDK, err := sdk.NewBOGOWISDK(&cfg.Testnet, testnetPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize testnet SDK: %w", err)
		}
		handler.testnetSDK = testnetSDK
	}

	// Initialize mainnet SDK
	if cfg.Mainnet.Contracts.BOGOToken != "" || cfg.Mainnet.Contracts.RewardDistributor != "" {
		mainnetPrivateKey := cfg.MainnetPrivateKey
		if mainnetPrivateKey == "" {
			return nil, fmt.Errorf("MAINNET_PRIVATE_KEY is required for mainnet operations")
		}
		mainnetSDK, err := sdk.NewBOGOWISDK(&cfg.Mainnet, mainnetPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mainnet SDK: %w", err)
		}
		handler.mainnetSDK = mainnetSDK
	}

	// Initialize NFT SDKs if BOGOWITickets contract is configured
	if cfg.Testnet.Contracts.BOGOWITickets != "" {
		testnetPrivateKey := cfg.TestnetPrivateKey
		if testnetPrivateKey == "" {
			return nil, fmt.Errorf("TESTNET_PRIVATE_KEY is required for NFT operations")
		}
		testnetNFTConfig := nft.ClientConfig{
			PrivateKey:      testnetPrivateKey,
			Network:         "testnet",
			CustomRPCURL:    cfg.Testnet.RPCUrl,
			DatakyteEnabled: true,
		}
		testnetNFTSDK, err := nft.NewClient(testnetNFTConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize testnet NFT SDK: %w", err)
		}
		handler.testnetNFTSDK = testnetNFTSDK
	}

	if cfg.Mainnet.Contracts.BOGOWITickets != "" {
		mainnetPrivateKey := cfg.MainnetPrivateKey
		if mainnetPrivateKey == "" {
			return nil, fmt.Errorf("MAINNET_PRIVATE_KEY is required for NFT operations")
		}
		mainnetNFTConfig := nft.ClientConfig{
			PrivateKey:      mainnetPrivateKey,
			Network:         "mainnet",
			CustomRPCURL:    cfg.Mainnet.RPCUrl,
			DatakyteEnabled: true,
		}
		mainnetNFTSDK, err := nft.NewClient(mainnetNFTConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mainnet NFT SDK: %w", err)
		}
		handler.mainnetNFTSDK = mainnetNFTSDK
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

// GetNFTSDK returns the appropriate NFT SDK based on the network parameter
func (h *NetworkHandler) GetNFTSDK(network string) (*nft.Client, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	switch network {
	case "testnet", "columbus":
		if h.testnetNFTSDK == nil {
			return nil, fmt.Errorf("testnet NFT SDK not initialized")
		}
		return h.testnetNFTSDK, nil
	case "mainnet", "camino":
		if h.mainnetNFTSDK == nil {
			return nil, fmt.Errorf("mainnet NFT SDK not initialized")
		}
		return h.mainnetNFTSDK, nil
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
	if h.testnetNFTSDK != nil {
		h.testnetNFTSDK.Close()
	}
	if h.mainnetNFTSDK != nil {
		h.mainnetNFTSDK.Close()
	}
}
