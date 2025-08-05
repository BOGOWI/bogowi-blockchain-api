package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Test loading config with environment variables
	os.Setenv("NODE_ENV", "test")
	os.Setenv("API_PORT", "8080")
	os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.Environment)
	assert.Equal(t, "8080", cfg.APIPort)

	// Cleanup
	os.Unsetenv("NODE_ENV")
	os.Unsetenv("API_PORT")
	os.Unsetenv("TESTNET_PRIVATE_KEY")
}

func TestConfigDefaults(t *testing.T) {
	// Test loading config with defaults
	os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Environment)
	assert.Equal(t, "development", cfg.Environment) // default value
	assert.Equal(t, "3001", cfg.APIPort)            // default value

	// Cleanup
	os.Unsetenv("TESTNET_PRIVATE_KEY")
}
