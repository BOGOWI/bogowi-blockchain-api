// +build integration

package api

import (
	"context"
	"os"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewNetworkHandler_Integration tests the actual initialization with real connections
// Run with: go test -tags=integration ./internal/api -run TestNewNetworkHandler_Integration
func TestNewNetworkHandler_Integration(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
		checkFunc func(t *testing.T, handler *NetworkHandler)
	}{
		{
			name: "initialize with testnet NFT SDK",
			config: &config.Config{
				Testnet: config.NetworkConfig{
					RPCUrl:  "http://localhost:8545",
					ChainID: 501,
					Contracts: config.ContractAddresses{
						BOGOWITickets:     "0x3Aa5ebB10DC797CAC828524e59A333d0A371443c",
						BOGOToken:         "0x959922bE3CAee4b8Cd9a407cc3ac1C251C2007B1",
						RewardDistributor: "0x0B306BF915C4d645ff596e518fAf3F9669b97016",
						RoleManager:       "0x9A676e781A523b5d0C0e43731313A708CB607508",
					},
				},
				TestnetPrivateKey: os.Getenv("TESTNET_PRIVATE_KEY"),
			},
			wantErr: false,
			checkFunc: func(t *testing.T, handler *NetworkHandler) {
				assert.NotNil(t, handler.testnetSDK)
				assert.NotNil(t, handler.testnetNFTSDK)
				
				// Test that we can get the NFT SDK
				nftSDK, err := handler.GetNFTSDK("testnet")
				require.NoError(t, err)
				assert.NotNil(t, nftSDK)
				
				// Test that we can get the regular SDK
				sdk, err := handler.GetSDK("testnet")
				require.NoError(t, err)
				assert.NotNil(t, sdk)
			},
		},
		{
			name: "initialize with both testnet and mainnet",
			config: &config.Config{
				Testnet: config.NetworkConfig{
					RPCUrl:  "https://columbus.camino.network/ext/bc/C/rpc",
					ChainID: 501,
					Contracts: config.ContractAddresses{
						BOGOWITickets:     os.Getenv("NFT_TICKETS_TESTNET_CONTRACT"),
						BOGOToken:         os.Getenv("BOGO_TOKEN_TESTNET"),
						RewardDistributor: os.Getenv("REWARD_DISTRIBUTOR_TESTNET"),
					},
				},
				Mainnet: config.NetworkConfig{
					RPCUrl:  "https://api.camino.network/ext/bc/C/rpc",
					ChainID: 500,
					Contracts: config.ContractAddresses{
						BOGOWITickets:     os.Getenv("NFT_TICKETS_MAINNET_CONTRACT"),
						BOGOToken:         os.Getenv("BOGO_TOKEN_MAINNET"),
						RewardDistributor: os.Getenv("REWARD_DISTRIBUTOR_MAINNET"),
					},
				},
				TestnetPrivateKey: os.Getenv("TESTNET_PRIVATE_KEY"),
				MainnetPrivateKey: os.Getenv("MAINNET_PRIVATE_KEY"),
			},
			wantErr: os.Getenv("TESTNET_PRIVATE_KEY") == "" || os.Getenv("MAINNET_PRIVATE_KEY") == "",
			checkFunc: func(t *testing.T, handler *NetworkHandler) {
				if handler != nil {
					// Both SDKs should be initialized if keys are provided
					if os.Getenv("TESTNET_PRIVATE_KEY") != "" {
						assert.NotNil(t, handler.testnetSDK)
						if os.Getenv("NFT_TICKETS_TESTNET_CONTRACT") != "" {
							assert.NotNil(t, handler.testnetNFTSDK)
						}
					}
					
					if os.Getenv("MAINNET_PRIVATE_KEY") != "" {
						assert.NotNil(t, handler.mainnetSDK)
						if os.Getenv("NFT_TICKETS_MAINNET_CONTRACT") != "" {
							assert.NotNil(t, handler.mainnetNFTSDK)
						}
					}
				}
			},
		},
		{
			name: "error when NFT contract configured but private key missing",
			config: &config.Config{
				Testnet: config.NetworkConfig{
					RPCUrl:  "http://localhost:8545",
					ChainID: 501,
					Contracts: config.ContractAddresses{
						BOGOWITickets: "0x3Aa5ebB10DC797CAC828524e59A333d0A371443c",
					},
				},
				TestnetPrivateKey: "", // Missing private key
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := NewNetworkHandler(tt.config)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, handler)
				
				if tt.checkFunc != nil {
					tt.checkFunc(t, handler)
				}
				
				// Test cleanup
				if handler != nil {
					handler.Close()
				}
			}
		})
	}
}

// TestNetworkHandler_NFTOperations_Integration tests actual NFT operations
func TestNetworkHandler_NFTOperations_Integration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}

	// Setup
	cfg := &config.Config{
		Testnet: config.NetworkConfig{
			RPCUrl:  "http://localhost:8545",
			ChainID: 501,
			Contracts: config.ContractAddresses{
				BOGOWITickets:     "0x3Aa5ebB10DC797CAC828524e59A333d0A371443c",
				NFTRegistry:       "0x68B1D87F95878fE05B998F19b66F4baba5De1aed",
				RoleManager:       "0x9A676e781A523b5d0C0e43731313A708CB607508",
				BOGOToken:         "0x959922bE3CAee4b8Cd9a407cc3ac1C251C2007B1",
			},
		},
		TestnetPrivateKey: os.Getenv("TESTNET_PRIVATE_KEY"),
	}

	if cfg.TestnetPrivateKey == "" {
		t.Skip("TESTNET_PRIVATE_KEY not set")
	}

	handler, err := NewNetworkHandler(cfg)
	require.NoError(t, err)
	defer handler.Close()

	// Test NFT SDK operations
	nftSDK, err := handler.GetNFTSDK("testnet")
	require.NoError(t, err)
	assert.NotNil(t, nftSDK)

	// Test that the NFT SDK is properly initialized
	address := nftSDK.GetAddress()
	assert.NotEqual(t, "0x0000000000000000000000000000000000000000", address.Hex())
	
	// Test balance query (should work even with 0 balance)
	ctx := context.Background()
	balance, err := nftSDK.GetBalance(ctx)
	require.NoError(t, err)
	assert.NotNil(t, balance)
}