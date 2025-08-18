package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestLoadConfigWithMainnetKey(t *testing.T) {
	// Test loading config with mainnet key only
	os.Setenv("MAINNET_PRIVATE_KEY", "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", cfg.MainnetPrivateKey)

	// Cleanup
	os.Unsetenv("MAINNET_PRIVATE_KEY")
}

func TestLoadConfigWithBothKeys(t *testing.T) {
	// Test loading config with both testnet and mainnet keys
	os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	os.Setenv("MAINNET_PRIVATE_KEY", "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", cfg.TestnetPrivateKey)
	assert.Equal(t, "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", cfg.MainnetPrivateKey)

	// Cleanup
	os.Unsetenv("TESTNET_PRIVATE_KEY")
	os.Unsetenv("MAINNET_PRIVATE_KEY")
}

func TestLoadConfigWithLegacyPrivateKey(t *testing.T) {
	// Test loading config with legacy PRIVATE_KEY
	os.Setenv("PRIVATE_KEY", "0xlegacy1234567890abcdef1234567890abcdef1234567890abcdef1234567890")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "0xlegacy1234567890abcdef1234567890abcdef1234567890abcdef1234567890", cfg.TestnetPrivateKey)
	assert.Equal(t, "0xlegacy1234567890abcdef1234567890abcdef1234567890abcdef1234567890", cfg.PrivateKey)

	// Cleanup
	os.Unsetenv("PRIVATE_KEY")
}

func TestLoadConfigWithAPIPrivateKey(t *testing.T) {
	// Test loading config with API_PRIVATE_KEY (higher priority than PRIVATE_KEY)
	os.Setenv("API_PRIVATE_KEY", "0xapi1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	os.Setenv("PRIVATE_KEY", "0xlegacy1234567890abcdef1234567890abcdef1234567890abcdef1234567890")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "0xapi1234567890abcdef1234567890abcdef1234567890abcdef1234567890", cfg.TestnetPrivateKey)
	assert.Equal(t, "0xapi1234567890abcdef1234567890abcdef1234567890abcdef1234567890", cfg.PrivateKey)

	// Cleanup
	os.Unsetenv("API_PRIVATE_KEY")
	os.Unsetenv("PRIVATE_KEY")
}

func TestLoadConfigNoPrivateKeys(t *testing.T) {
	// Test that config fails when no private keys are provided
	// Make sure no keys are set
	os.Unsetenv("TESTNET_PRIVATE_KEY")
	os.Unsetenv("MAINNET_PRIVATE_KEY")
	os.Unsetenv("PRIVATE_KEY")
	os.Unsetenv("API_PRIVATE_KEY")

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "At least one private key")
}

func TestLoadConfigWithAuthSettings(t *testing.T) {
	// Test loading auth-related configuration
	os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	os.Setenv("SWAGGER_USERNAME", "admin")
	os.Setenv("SWAGGER_PASSWORD", "secret123")
	os.Setenv("FIREBASE_PROJECT_ID", "my-firebase-project")
	os.Setenv("BACKEND_SECRET", "my-backend-secret")
	os.Setenv("DEV_BACKEND_SECRET", "my-dev-backend-secret")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "admin", cfg.SwaggerUsername)
	assert.Equal(t, "secret123", cfg.SwaggerPassword)
	assert.Equal(t, "my-firebase-project", cfg.FirebaseProjectID)
	assert.Equal(t, "my-backend-secret", cfg.BackendSecret)
	assert.Equal(t, "my-dev-backend-secret", cfg.DevBackendSecret)

	// Cleanup
	os.Unsetenv("TESTNET_PRIVATE_KEY")
	os.Unsetenv("SWAGGER_USERNAME")
	os.Unsetenv("SWAGGER_PASSWORD")
	os.Unsetenv("FIREBASE_PROJECT_ID")
	os.Unsetenv("BACKEND_SECRET")
	os.Unsetenv("DEV_BACKEND_SECRET")
}

func TestLoadConfigDevBackendSecretDefault(t *testing.T) {
	// Test that DEV_BACKEND_SECRET defaults to BACKEND_SECRET if not set
	os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	os.Setenv("BACKEND_SECRET", "main-secret")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "main-secret", cfg.BackendSecret)
	assert.Equal(t, "main-secret", cfg.DevBackendSecret) // Should default to BackendSecret

	// Cleanup
	os.Unsetenv("TESTNET_PRIVATE_KEY")
	os.Unsetenv("BACKEND_SECRET")
}

func TestLoadConfigWithContractAddresses(t *testing.T) {
	// Test loading contract addresses
	os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	// Testnet contracts
	os.Setenv("TESTNET_ROLE_MANAGER_ADDRESS", "0xTestnetRoleManager")
	os.Setenv("TESTNET_BOGO_TOKEN_ADDRESS", "0xTestnetBOGO")
	os.Setenv("TESTNET_REWARD_DISTRIBUTOR_ADDRESS", "0xTestnetReward")
	os.Setenv("NFT_TICKETS_TESTNET_CONTRACT", "0xTestnetTickets")

	// Mainnet contracts
	os.Setenv("MAINNET_ROLE_MANAGER_ADDRESS", "0xMainnetRoleManager")
	os.Setenv("MAINNET_BOGO_TOKEN_ADDRESS", "0xMainnetBOGO")
	os.Setenv("MAINNET_REWARD_DISTRIBUTOR_ADDRESS", "0xMainnetReward")
	os.Setenv("NFT_TICKETS_MAINNET_CONTRACT", "0xMainnetTickets")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check testnet contracts
	assert.Equal(t, "0xTestnetRoleManager", cfg.Testnet.Contracts.RoleManager)
	assert.Equal(t, "0xTestnetBOGO", cfg.Testnet.Contracts.BOGOToken)
	assert.Equal(t, "0xTestnetReward", cfg.Testnet.Contracts.RewardDistributor)
	assert.Equal(t, "0xTestnetTickets", cfg.Testnet.Contracts.BOGOWITickets)

	// Check mainnet contracts
	assert.Equal(t, "0xMainnetRoleManager", cfg.Mainnet.Contracts.RoleManager)
	assert.Equal(t, "0xMainnetBOGO", cfg.Mainnet.Contracts.BOGOToken)
	assert.Equal(t, "0xMainnetReward", cfg.Mainnet.Contracts.RewardDistributor)
	assert.Equal(t, "0xMainnetTickets", cfg.Mainnet.Contracts.BOGOWITickets)

	// Cleanup
	os.Unsetenv("TESTNET_PRIVATE_KEY")
	os.Unsetenv("TESTNET_ROLE_MANAGER_ADDRESS")
	os.Unsetenv("TESTNET_BOGO_TOKEN_ADDRESS")
	os.Unsetenv("TESTNET_REWARD_DISTRIBUTOR_ADDRESS")
	os.Unsetenv("NFT_TICKETS_TESTNET_CONTRACT")
	os.Unsetenv("MAINNET_ROLE_MANAGER_ADDRESS")
	os.Unsetenv("MAINNET_BOGO_TOKEN_ADDRESS")
	os.Unsetenv("MAINNET_REWARD_DISTRIBUTOR_ADDRESS")
	os.Unsetenv("NFT_TICKETS_MAINNET_CONTRACT")
}

func TestLoadConfigBackwardCompatibility(t *testing.T) {
	t.Run("development environment uses simple names for testnet", func(t *testing.T) {
		os.Setenv("NODE_ENV", "development")
		os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		os.Setenv("ROLE_MANAGER_ADDRESS", "0xSimpleRoleManager")
		os.Setenv("BOGO_TOKEN_ADDRESS", "0xSimpleBOGO")
		os.Setenv("REWARD_DISTRIBUTOR_ADDRESS", "0xSimpleReward")

		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		// In development, simple names override testnet
		assert.Equal(t, "0xSimpleRoleManager", cfg.Testnet.Contracts.RoleManager)
		assert.Equal(t, "0xSimpleBOGO", cfg.Testnet.Contracts.BOGOToken)
		assert.Equal(t, "0xSimpleReward", cfg.Testnet.Contracts.RewardDistributor)

		// Cleanup
		os.Unsetenv("NODE_ENV")
		os.Unsetenv("TESTNET_PRIVATE_KEY")
		os.Unsetenv("ROLE_MANAGER_ADDRESS")
		os.Unsetenv("BOGO_TOKEN_ADDRESS")
		os.Unsetenv("REWARD_DISTRIBUTOR_ADDRESS")
	})

	t.Run("production environment uses simple names for mainnet", func(t *testing.T) {
		os.Setenv("NODE_ENV", "production")
		os.Setenv("MAINNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		os.Setenv("ROLE_MANAGER_ADDRESS", "0xSimpleRoleManager")
		os.Setenv("BOGO_TOKEN_ADDRESS", "0xSimpleBOGO")
		os.Setenv("REWARD_DISTRIBUTOR_ADDRESS", "0xSimpleReward")

		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		// In production, simple names override mainnet
		assert.Equal(t, "0xSimpleRoleManager", cfg.Mainnet.Contracts.RoleManager)
		assert.Equal(t, "0xSimpleBOGO", cfg.Mainnet.Contracts.BOGOToken)
		assert.Equal(t, "0xSimpleReward", cfg.Mainnet.Contracts.RewardDistributor)

		// Cleanup
		os.Unsetenv("NODE_ENV")
		os.Unsetenv("MAINNET_PRIVATE_KEY")
		os.Unsetenv("ROLE_MANAGER_ADDRESS")
		os.Unsetenv("BOGO_TOKEN_ADDRESS")
		os.Unsetenv("REWARD_DISTRIBUTOR_ADDRESS")
	})
}

func TestNetworkConfig(t *testing.T) {
	os.Setenv("TESTNET_PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Test testnet configuration
	assert.Equal(t, "https://columbus.camino.network/ext/bc/C/rpc", cfg.Testnet.RPCUrl)
	assert.Equal(t, int64(501), cfg.Testnet.ChainID)

	// Test mainnet configuration
	assert.Equal(t, "https://api.camino.network/ext/bc/C/rpc", cfg.Mainnet.RPCUrl)
	assert.Equal(t, int64(500), cfg.Mainnet.ChainID)

	// Cleanup
	os.Unsetenv("TESTNET_PRIVATE_KEY")
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{
			name:         "existing environment variable",
			key:          "TEST_VAR",
			value:        "test_value",
			defaultValue: "default",
			expected:     "test_value",
		},
		{
			name:         "non-existing environment variable",
			key:          "NON_EXISTENT",
			value:        "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
