package sdk

import (
	"testing"

	"bogowi-blockchain-go/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNewBOGOWISDK(t *testing.T) {
	// Test creating new SDK instance
	cfg := &config.Config{
		RPCUrl:     "https://columbus.camino.network/ext/bc/C/rpc",
		ChainID:    501,
		PrivateKey: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Contracts:  config.ContractAddresses{},
	}

	sdk, err := NewBOGOWISDK(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, sdk)
}

func TestSDKWithInvalidPrivateKey(t *testing.T) {
	// Test creating SDK with invalid private key
	cfg := &config.Config{
		RPCUrl:     "https://columbus.camino.network/ext/bc/C/rpc",
		ChainID:    501,
		PrivateKey: "invalid",
		Contracts:  config.ContractAddresses{},
	}

	sdk, err := NewBOGOWISDK(cfg)
	assert.Error(t, err)
	assert.Nil(t, sdk)
}
