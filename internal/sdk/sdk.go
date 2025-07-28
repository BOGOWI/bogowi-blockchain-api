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
	client            EthClient
	auth              *bind.TransactOpts
	chainID           *big.Int
	contracts         *ContractInstances
	config            *config.Config
	privateKey        *ecdsa.PrivateKey
	rewardDistributor *Contract
}

// ContractInstances holds all initialized contract instances
type ContractInstances struct {
	BOGOTokenV2       *Contract
	ConservationNFT   *Contract
	CommercialNFT     *Contract
	RewardDistributor *Contract
	MultisigTreasury  *Contract
}

// Contract represents a generic contract with ABI and address
type Contract struct {
	Address  common.Address
	ABI      abi.ABI
	Instance BoundContract
}

// TokenBalance represents a token balance response
type TokenBalance struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
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
		client:     client,
		auth:       auth,
		chainID:    chainID,
		config:     cfg,
		privateKey: privateKey,
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

	// Check if client is a real ethclient.Client to create BoundContract
	if ethClient, ok := s.client.(*ethclient.Client); ok {
		instance := bind.NewBoundContract(contractAddress, contractABI, ethClient, ethClient, ethClient)
		return &Contract{
			Address:  contractAddress,
			ABI:      contractABI,
			Instance: &BoundContractWrapper{instance},
		}, nil
	}

	return nil, fmt.Errorf("client is not an ethclient.Client")
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

// TransferBOGOTokens transfers BOGO tokens to a recipient
func (s *BOGOWISDK) TransferBOGOTokens(to string, amount string) (string, error) {
	if !common.IsHexAddress(to) {
		return "", fmt.Errorf("invalid recipient address")
	}

	if s.contracts.BOGOTokenV2 == nil {
		return "", fmt.Errorf("BOGO token contract not initialized")
	}

	// Parse amount (assuming it's in ether units, convert to wei)
	amountFloat, ok := new(big.Float).SetString(amount)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// Convert to wei (multiply by 10^18)
	decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	amountWei := new(big.Int)
	amountFloat.Mul(amountFloat, new(big.Float).SetInt(decimals))
	amountFloat.Int(amountWei)

	// Prepare transaction
	toAddress := common.HexToAddress(to)

	// Get current nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// Set transaction options
	s.auth.Nonce = big.NewInt(int64(nonce))
	s.auth.GasPrice = gasPrice
	s.auth.GasLimit = uint64(100000) // Standard gas limit for ERC20 transfer

	// Execute transfer
	tx, err := s.contracts.BOGOTokenV2.Instance.Transact(s.auth, "transfer", toAddress, amountWei)
	if err != nil {
		return "", fmt.Errorf("failed to execute transfer: %w", err)
	}

	return tx.Hash().Hex(), nil
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
