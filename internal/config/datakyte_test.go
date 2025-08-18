package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatakyteConfig(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *DatakyteConfig
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			expectedConfig: &DatakyteConfig{
				TestnetAPIKey: "dk_d707e26c919e72ab2bb3b81897566c393f4e2eba54d07ff680d765ee03d6cc5d",
				MainnetAPIKey: "dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d",
				BaseURL:       "https://dklnk.to",
			},
		},
		{
			name: "custom environment variables",
			envVars: map[string]string{
				"DATAKYTE_TESTNET_API_KEY": "custom_testnet_key",
				"DATAKYTE_MAINNET_API_KEY": "custom_mainnet_key",
				"DATAKYTE_BASE_URL":        "https://custom.datakyte.com",
			},
			expectedConfig: &DatakyteConfig{
				TestnetAPIKey: "custom_testnet_key",
				MainnetAPIKey: "custom_mainnet_key",
				BaseURL:       "https://custom.datakyte.com",
			},
		},
		{
			name: "partial custom configuration",
			envVars: map[string]string{
				"DATAKYTE_TESTNET_API_KEY": "only_testnet_custom",
			},
			expectedConfig: &DatakyteConfig{
				TestnetAPIKey: "only_testnet_custom",
				MainnetAPIKey: "dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d",
				BaseURL:       "https://dklnk.to",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			// Get configuration
			config := GetDatakyteConfig()

			// Assert
			assert.Equal(t, tt.expectedConfig.TestnetAPIKey, config.TestnetAPIKey)
			assert.Equal(t, tt.expectedConfig.MainnetAPIKey, config.MainnetAPIKey)
			assert.Equal(t, tt.expectedConfig.BaseURL, config.BaseURL)
		})
	}
}

func TestGetAPIKeyForNetwork(t *testing.T) {
	config := &DatakyteConfig{
		TestnetAPIKey: "testnet_key_123",
		MainnetAPIKey: "mainnet_key_456",
		BaseURL:       "https://api.datakyte.com",
	}

	tests := []struct {
		name        string
		network     string
		expectedKey string
	}{
		{
			name:        "mainnet network",
			network:     "mainnet",
			expectedKey: "mainnet_key_456",
		},
		{
			name:        "testnet network",
			network:     "testnet",
			expectedKey: "testnet_key_123",
		},
		{
			name:        "unknown network defaults to testnet",
			network:     "unknown",
			expectedKey: "testnet_key_123",
		},
		{
			name:        "empty network defaults to testnet",
			network:     "",
			expectedKey: "testnet_key_123",
		},
		{
			name:        "development network defaults to testnet",
			network:     "development",
			expectedKey: "testnet_key_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey := config.GetAPIKeyForNetwork(tt.network)
			assert.Equal(t, tt.expectedKey, apiKey)
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "environment variable exists",
			envKey:       "TEST_ENV_VAR",
			envValue:     "env_value",
			defaultValue: "default_value",
			expected:     "env_value",
		},
		{
			name:         "environment variable does not exist",
			envKey:       "NON_EXISTENT_VAR",
			envValue:     "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
		{
			name:         "empty environment variable returns default",
			envKey:       "EMPTY_VAR",
			envValue:     "",
			defaultValue: "fallback",
			expected:     "fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			result := getEnvOrDefault(tt.envKey, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatakyteConfigIntegration(t *testing.T) {
	// Test the full flow of getting config and using it
	t.Run("full integration test", func(t *testing.T) {
		// Set custom environment variables
		os.Setenv("DATAKYTE_TESTNET_API_KEY", "integration_testnet")
		os.Setenv("DATAKYTE_MAINNET_API_KEY", "integration_mainnet")
		os.Setenv("DATAKYTE_BASE_URL", "https://integration.test")
		defer func() {
			os.Unsetenv("DATAKYTE_TESTNET_API_KEY")
			os.Unsetenv("DATAKYTE_MAINNET_API_KEY")
			os.Unsetenv("DATAKYTE_BASE_URL")
		}()

		// Get configuration
		config := GetDatakyteConfig()

		// Verify configuration
		assert.NotNil(t, config)
		assert.Equal(t, "integration_testnet", config.TestnetAPIKey)
		assert.Equal(t, "integration_mainnet", config.MainnetAPIKey)
		assert.Equal(t, "https://integration.test", config.BaseURL)

		// Test network-specific API key retrieval
		assert.Equal(t, "integration_mainnet", config.GetAPIKeyForNetwork("mainnet"))
		assert.Equal(t, "integration_testnet", config.GetAPIKeyForNetwork("testnet"))
		assert.Equal(t, "integration_testnet", config.GetAPIKeyForNetwork("staging"))
	})
}
