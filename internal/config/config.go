package config

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Environment string `json:"environment"`
	APIPort     string `json:"api_port"`

	// Private Key (same for both networks)
	PrivateKey string `json:"private_key"`

	// Network-specific configurations
	Testnet NetworkConfig `json:"testnet"`
	Mainnet NetworkConfig `json:"mainnet"`

	// Auth
	SwaggerUsername   string `json:"swagger_username"`
	SwaggerPassword   string `json:"swagger_password"`
	FirebaseProjectID string `json:"firebase_project_id"`
	BackendSecret     string `json:"backend_secret"`
	DevBackendSecret  string `json:"dev_backend_secret"`
}

// NetworkConfig holds network-specific configuration
type NetworkConfig struct {
	RPCUrl    string            `json:"rpc_url"`
	ChainID   int64             `json:"chain_id"`
	Contracts ContractAddresses `json:"contracts"`
}

// ContractAddresses holds all smart contract addresses
type ContractAddresses struct {
	// V1 Contracts
	RoleManager       string `json:"role_manager"`
	BOGOToken         string `json:"bogo_token"`
	RewardDistributor string `json:"reward_distributor"`

	// Legacy contracts (to be removed after migration)
	BOGOTokenV2      string `json:"bogo_token_v2"`
	ConservationNFT  string `json:"conservation_nft"`
	CommercialNFT    string `json:"commercial_nft"`
	MultisigTreasury string `json:"multisig_treasury"`
}

// Load loads configuration from environment variables and AWS SSM
func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}

	cfg := &Config{
		Environment: getEnv("NODE_ENV", "development"),
		APIPort:     getEnv("API_PORT", "3001"),

		// Configure testnet
		Testnet: NetworkConfig{
			RPCUrl:  "https://columbus.camino.network/ext/bc/C/rpc",
			ChainID: 501,
		},

		// Configure mainnet
		Mainnet: NetworkConfig{
			RPCUrl:  "https://api.camino.network/ext/bc/C/rpc",
			ChainID: 500,
		},
	}

	// Load secrets from AWS SSM in production
	if cfg.Environment == "production" && getEnv("PRIVATE_KEY", "") == "" {
		log.Println("Loading secrets from AWS Systems Manager...")
		if err := loadSecretsFromSSM(cfg); err != nil {
			return nil, fmt.Errorf("failed to load secrets from AWS SSM: %w", err)
		}
	} else {
		log.Println("Using local environment variables for configuration")
		loadFromEnv(cfg)
	}

	// Validate required fields
	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY is required")
	}

	return cfg, nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	cfg.PrivateKey = getEnv("API_PRIVATE_KEY", getEnv("PRIVATE_KEY", ""))
	cfg.SwaggerUsername = getEnv("SWAGGER_USERNAME", "")
	cfg.SwaggerPassword = getEnv("SWAGGER_PASSWORD", "")
	cfg.FirebaseProjectID = getEnv("FIREBASE_PROJECT_ID", "")
	cfg.BackendSecret = getEnv("BACKEND_SECRET", "backend-secret-key")
	cfg.DevBackendSecret = getEnv("DEV_BACKEND_SECRET", cfg.BackendSecret) // Default to main secret if not set

	// Load testnet contracts
	cfg.Testnet.Contracts = ContractAddresses{
		// V1 Contracts - Testnet
		RoleManager:       getEnv("TESTNET_ROLE_MANAGER_ADDRESS", ""),
		BOGOToken:         getEnv("TESTNET_BOGO_TOKEN_ADDRESS", ""),
		RewardDistributor: getEnv("TESTNET_REWARD_DISTRIBUTOR_ADDRESS", ""),

		// Legacy contracts (for backward compatibility)
		BOGOTokenV2:      getEnv("TESTNET_BOGO_TOKEN_V2_ADDRESS", ""),
		ConservationNFT:  getEnv("TESTNET_CONSERVATION_NFT_ADDRESS", ""),
		CommercialNFT:    getEnv("TESTNET_COMMERCIAL_NFT_ADDRESS", ""),
		MultisigTreasury: getEnv("TESTNET_MULTISIG_ADDRESS", ""),
	}

	// Load mainnet contracts
	cfg.Mainnet.Contracts = ContractAddresses{
		// V1 Contracts - Mainnet
		RoleManager:       getEnv("MAINNET_ROLE_MANAGER_ADDRESS", ""),
		BOGOToken:         getEnv("MAINNET_BOGO_TOKEN_ADDRESS", ""),
		RewardDistributor: getEnv("MAINNET_REWARD_DISTRIBUTOR_ADDRESS", ""),

		// Legacy contracts (for backward compatibility)
		BOGOTokenV2:      getEnv("MAINNET_BOGO_TOKEN_V2_ADDRESS", ""),
		ConservationNFT:  getEnv("MAINNET_CONSERVATION_NFT_ADDRESS", ""),
		CommercialNFT:    getEnv("MAINNET_COMMERCIAL_NFT_ADDRESS", ""),
		MultisigTreasury: getEnv("MAINNET_MULTISIG_ADDRESS", ""),
	}
}

// loadSecretsFromSSM loads secrets from AWS Systems Manager Parameter Store
func loadSecretsFromSSM(cfg *Config) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(getEnv("AWS_REGION", "us-east-1")),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	ssmClient := ssm.New(sess)

	// Parameter names to load
	paramNames := []string{
		"PRIVATE_KEY",
		"API_PRIVATE_KEY",
		// V1 Contracts
		"ROLE_MANAGER_ADDRESS",
		"BOGO_TOKEN_ADDRESS",
		"REWARD_DISTRIBUTOR_ADDRESS",
		// Legacy contracts
		"BOGO_TOKEN_V2_ADDRESS",
		"CONSERVATION_NFT_ADDRESS",
		"COMMERCIAL_NFT_ADDRESS",
		"REWARD_DISTRIBUTOR_V2_ADDRESS",
		"MULTISIG_ADDRESS",
		// Auth and other configs
		"SWAGGER_USERNAME",
		"SWAGGER_PASSWORD",
		"FIREBASE_PROJECT_ID",
		"BACKEND_SECRET",
		"BACKEND_WALLET_ADDRESS",
	}

	// Get parameters
	input := &ssm.GetParametersInput{
		Names:          aws.StringSlice(paramNames),
		WithDecryption: aws.Bool(true),
	}

	result, err := ssmClient.GetParameters(input)
	if err != nil {
		return fmt.Errorf("failed to get parameters from SSM: %w", err)
	}

	// Set parameters as environment variables
	for _, param := range result.Parameters {
		if err := os.Setenv(*param.Name, *param.Value); err != nil {
			log.Printf("Warning: Failed to set environment variable %s: %v", *param.Name, err)
		}
	}

	// Log missing parameters
	if len(result.InvalidParameters) > 0 {
		log.Printf("Warning: Missing SSM parameters: %v", result.InvalidParameters)
	}

	// Load from environment now that SSM values are set
	loadFromEnv(cfg)

	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
