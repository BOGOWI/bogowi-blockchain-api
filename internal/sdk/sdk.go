package sdk

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"bogowi-blockchain-go/internal/config"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BOGOWISDK represents the main SDK for interacting with BOGOWI blockchain contracts
type BOGOWISDK struct {
	client    *ethclient.Client
	auth      *bind.TransactOpts
	chainID   *big.Int
	contracts *ContractInstances
	config    *config.Config
}

// ContractInstances holds all initialized contract instances
type ContractInstances struct {
	BOGOTokenV2       *Contract
	ConservationNFT   *Contract
	CommercialNFT     *Contract
	RewardDistributor *Contract
	MultisigTreasury  *Contract
	OceanBOGO         *Contract
	EarthBOGO         *Contract
	WildlifeBOGO      *Contract
}

// Contract represents a generic contract with ABI and address
type Contract struct {
	Address  common.Address
	ABI      abi.ABI
	Instance *bind.BoundContract
}

// TokenBalance represents a token balance response
type TokenBalance struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

// FlavoredBalances represents flavored token balances
type FlavoredBalances struct {
	Address  string            `json:"address"`
	Balances map[string]string `json:"balances"`
}

// DAOInfo represents DAO information
type DAOInfo struct {
	Threshold        int `json:"threshold"`
	SignerCount      int `json:"signerCount"`
	TransactionCount int `json:"transactionCount"`
}

// NewBOGOWISDK creates a new BOGOWI SDK instance
func NewBOGOWISDK(cfg *config.Config) (*BOGOWISDK, error) {
	// Connect to Ethereum client
	client, err := ethclient.Dial(cfg.RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get chain ID
	chainID := big.NewInt(cfg.ChainID)

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	sdk := &BOGOWISDK{
		client:  client,
		auth:    auth,
		chainID: chainID,
		config:  cfg,
	}

	// Initialize contracts
	if err := sdk.initializeContracts(); err != nil {
		return nil, fmt.Errorf("failed to initialize contracts: %w", err)
	}

	return sdk, nil
}

// initializeContracts initializes all contract instances
func (s *BOGOWISDK) initializeContracts() error {
	s.contracts = &ContractInstances{}

	// Initialize BOGOTokenV2
	if s.config.Contracts.BOGOTokenV2 != "" {
		contract, err := s.initializeContract(s.config.Contracts.BOGOTokenV2, ERC20ABI)
		if err != nil {
			return fmt.Errorf("failed to initialize BOGOTokenV2: %w", err)
		}
		s.contracts.BOGOTokenV2 = contract
	}

	// Initialize ConservationNFT
	if s.config.Contracts.ConservationNFT != "" {
		contract, err := s.initializeContract(s.config.Contracts.ConservationNFT, ERC721ABI)
		if err != nil {
			return fmt.Errorf("failed to initialize ConservationNFT: %w", err)
		}
		s.contracts.ConservationNFT = contract
	}

	// Initialize CommercialNFT
	if s.config.Contracts.CommercialNFT != "" {
		contract, err := s.initializeContract(s.config.Contracts.CommercialNFT, ERC721ABI)
		if err != nil {
			return fmt.Errorf("failed to initialize CommercialNFT: %w", err)
		}
		s.contracts.CommercialNFT = contract
	}

	// Initialize flavored tokens
	if s.config.Contracts.OceanBOGO != "" {
		contract, err := s.initializeContract(s.config.Contracts.OceanBOGO, ERC20ABI)
		if err != nil {
			return fmt.Errorf("failed to initialize OceanBOGO: %w", err)
		}
		s.contracts.OceanBOGO = contract
	}

	if s.config.Contracts.EarthBOGO != "" {
		contract, err := s.initializeContract(s.config.Contracts.EarthBOGO, ERC20ABI)
		if err != nil {
			return fmt.Errorf("failed to initialize EarthBOGO: %w", err)
		}
		s.contracts.EarthBOGO = contract
	}

	if s.config.Contracts.WildlifeBOGO != "" {
		contract, err := s.initializeContract(s.config.Contracts.WildlifeBOGO, ERC20ABI)
		if err != nil {
			return fmt.Errorf("failed to initialize WildlifeBOGO: %w", err)
		}
		s.contracts.WildlifeBOGO = contract
	}

	return nil
}

// initializeContract creates a contract instance with the given address and ABI
func (s *BOGOWISDK) initializeContract(address, abiJSON string) (*Contract, error) {
	if address == "" {
		return nil, fmt.Errorf("contract address is empty")
	}

	contractAddress := common.HexToAddress(address)

	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	instance := bind.NewBoundContract(contractAddress, contractABI, s.client, s.client, s.client)

	return &Contract{
		Address:  contractAddress,
		ABI:      contractABI,
		Instance: instance,
	}, nil
}

// GetTokenBalance gets the BOGO token balance for an address
func (s *BOGOWISDK) GetTokenBalance(address string) (*TokenBalance, error) {
	if s.contracts.BOGOTokenV2 == nil {
		return nil, fmt.Errorf("BOGOTokenV2 contract not initialized")
	}

	addr := common.HexToAddress(address)
	var balance *big.Int

	err := s.contracts.BOGOTokenV2.Instance.Call(
		&bind.CallOpts{Context: context.Background()},
		&[]interface{}{&balance},
		"balanceOf",
		addr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	// Convert wei to ether (18 decimals)
	decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	balanceEther := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		new(big.Float).SetInt(decimals),
	)

	return &TokenBalance{
		Address: address,
		Balance: balanceEther.String(),
	}, nil
}

// GetFlavoredTokenBalances gets balances for all flavored tokens
func (s *BOGOWISDK) GetFlavoredTokenBalances(address string) (*FlavoredBalances, error) {
	addr := common.HexToAddress(address)
	balances := make(map[string]string)

	// Get Ocean BOGO balance
	if s.contracts.OceanBOGO != nil {
		balance, err := s.getTokenBalanceFromContract(s.contracts.OceanBOGO, addr)
		if err != nil {
			balances["ocean"] = "0"
		} else {
			balances["ocean"] = balance
		}
	}

	// Get Earth BOGO balance
	if s.contracts.EarthBOGO != nil {
		balance, err := s.getTokenBalanceFromContract(s.contracts.EarthBOGO, addr)
		if err != nil {
			balances["earth"] = "0"
		} else {
			balances["earth"] = balance
		}
	}

	// Get Wildlife BOGO balance
	if s.contracts.WildlifeBOGO != nil {
		balance, err := s.getTokenBalanceFromContract(s.contracts.WildlifeBOGO, addr)
		if err != nil {
			balances["wildlife"] = "0"
		} else {
			balances["wildlife"] = balance
		}
	}

	return &FlavoredBalances{
		Address:  address,
		Balances: balances,
	}, nil
}

// getTokenBalanceFromContract gets balance from a specific token contract
func (s *BOGOWISDK) getTokenBalanceFromContract(contract *Contract, address common.Address) (string, error) {
	var balance *big.Int

	err := contract.Instance.Call(
		&bind.CallOpts{Context: context.Background()},
		&[]interface{}{&balance},
		"balanceOf",
		address,
	)
	if err != nil {
		return "0", err
	}

	// Convert wei to ether (18 decimals)
	decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	balanceEther := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		new(big.Float).SetInt(decimals),
	)
	return balanceEther.String(), nil
}

// GetGasPrice gets the current gas price
func (s *BOGOWISDK) GetGasPrice() (string, error) {
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// Convert to Gwei
	gwei := new(big.Float).Quo(new(big.Float).SetInt(gasPrice), new(big.Float).SetInt(big.NewInt(1000000000)))
	return fmt.Sprintf("%.2f gwei", gwei), nil
}

// GetPublicKey returns the public key associated with the private key
func (s *BOGOWISDK) GetPublicKey() (string, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(s.config.PrivateKey, "0x"))
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex(), nil
}

// Close closes the SDK and cleans up resources
func (s *BOGOWISDK) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
