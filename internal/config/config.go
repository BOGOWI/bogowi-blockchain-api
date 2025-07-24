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

	// Blockchain
	RPCUrl     string `json:"rpc_url"`
	ChainID    int64  `json:"chain_id"`
	PrivateKey string `json:"private_key"`

	// Contract Addresses
	Contracts ContractAddresses `json:"contracts"`

	// Auth
	SwaggerUsername string `json:"swagger_username"`
	SwaggerPassword string `json:"swagger_password"`
}

// ContractAddresses holds all smart contract addresses
type ContractAddresses struct {
	BOGOTokenV2       string `json:"bogo_token_v2"`
	ConservationNFT   string `json:"conservation_nft"`
	CommercialNFT     string `json:"commercial_nft"`
	RewardDistributor string `json:"reward_distributor"`
	MultisigTreasury  string `json:"multisig_treasury"`
	OceanBOGO         string `json:"ocean_bogo"`
	EarthBOGO         string `json:"earth_bogo"`
	WildlifeBOGO      string `json:"wildlife_bogo"`
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
		RPCUrl:      getEnv("RPC_URL", "https://columbus.camino.network/ext/bc/C/rpc"),
		ChainID:     501,
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

	cfg.Contracts = ContractAddresses{
		BOGOTokenV2:       getEnv("BOGO_TOKEN_V2_ADDRESS", ""),
		ConservationNFT:   getEnv("CONSERVATION_NFT_ADDRESS", ""),
		CommercialNFT:     getEnv("COMMERCIAL_NFT_ADDRESS", ""),
		RewardDistributor: getEnv("REWARD_DISTRIBUTOR_V2_ADDRESS", ""),
		MultisigTreasury:  getEnv("MULTISIG_ADDRESS", ""),
		OceanBOGO:         getEnv("OCEAN_BOGO_ADDRESS", ""),
		EarthBOGO:         getEnv("EARTH_BOGO_ADDRESS", ""),
		WildlifeBOGO:      getEnv("WILDLIFE_BOGO_ADDRESS", ""),
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
		"BOGO_TOKEN_V2_ADDRESS",
		"CONSERVATION_NFT_ADDRESS",
		"COMMERCIAL_NFT_ADDRESS",
		"REWARD_DISTRIBUTOR_V2_ADDRESS",
		"MULTISIG_ADDRESS",
		"OCEAN_BOGO_ADDRESS",
		"EARTH_BOGO_ADDRESS",
		"WILDLIFE_BOGO_ADDRESS",
		"SWAGGER_USERNAME",
		"SWAGGER_PASSWORD",
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
		os.Setenv(*param.Name, *param.Value)
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
