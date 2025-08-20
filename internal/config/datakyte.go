package config

import (
	"os"
)

// DatakyteConfig holds Datakyte API configuration
type DatakyteConfig struct {
	TestnetAPIKey string
	MainnetAPIKey string
	BaseURL       string
}

// GetDatakyteConfig returns the Datakyte configuration
func GetDatakyteConfig() *DatakyteConfig {
	return &DatakyteConfig{
		TestnetAPIKey: getEnvOrDefault("DATAKYTE_TESTNET_API_KEY", "dk_d707e26c919e72ab2bb3b81897566c393f4e2eba54d07ff680d765ee03d6cc5d"),
		MainnetAPIKey: getEnvOrDefault("DATAKYTE_MAINNET_API_KEY", "dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d"),
		BaseURL:       getEnvOrDefault("DATAKYTE_BASE_URL", "https://dklnk.to"),
	}
}

// GetAPIKeyForNetwork returns the appropriate API key based on network
func (c *DatakyteConfig) GetAPIKeyForNetwork(network string) string {
	switch network {
	case "mainnet":
		return c.MainnetAPIKey
	case "testnet":
		return c.TestnetAPIKey
	default:
		// Default to testnet for safety
		return c.TestnetAPIKey
	}
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
