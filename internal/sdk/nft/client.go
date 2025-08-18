package nft

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"time"

	"bogowi-blockchain-go/internal/sdk/contracts"
	"bogowi-blockchain-go/internal/services/datakyte"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client is the main SDK client for NFT ticket operations
type Client struct {
	ethClient          *ethclient.Client
	ticketsContract    TicketsContractInterface
	ticketsAddress     common.Address
	roleManager        *contracts.RoleManager
	roleManagerAddress common.Address
	auth               *bind.TransactOpts
	chainID            *big.Int
	network            string
	config             *ClientConfig
	networkConfig      *NetworkConfig
	datakyteService    *datakyte.TicketMetadataService
	privateKey         *ecdsa.PrivateKey
}

// NewClient creates a new NFT SDK client
func NewClient(config ClientConfig) (*Client, error) {
	// Get network configuration
	networkConfig, err := GetNetworkConfig(config.Network)
	if err != nil {
		return nil, fmt.Errorf("failed to get network config: %w", err)
	}

	// Use custom RPC if provided
	rpcURL := networkConfig.RPCURL
	if config.CustomRPCURL != "" {
		rpcURL = config.CustomRPCURL
	}

	// Connect to Ethereum client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ethClient, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Camino network: %w", err)
	}

	// Verify chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	if chainID.Cmp(networkConfig.ChainID) != 0 {
		return nil, fmt.Errorf("chain ID mismatch: expected %s, got %s",
			networkConfig.ChainID.String(), chainID.String())
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Create auth transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set default gas settings
	if config.GasMultiplier == 0 {
		config.GasMultiplier = 1.2
	}

	// Initialize client
	client := &Client{
		ethClient:     ethClient,
		auth:          auth,
		chainID:       chainID,
		network:       config.Network,
		config:        &config,
		networkConfig: networkConfig,
		privateKey:    privateKey,
	}

	// Load contract addresses from environment or config
	if err := client.loadContracts(); err != nil {
		return nil, fmt.Errorf("failed to load contracts: %w", err)
	}

	// Initialize Datakyte service if enabled
	if config.DatakyteEnabled {
		client.initDatakyteService()
	}

	return client, nil
}

// loadContracts loads the smart contract instances
func (c *Client) loadContracts() error {
	// Get contract addresses from environment
	ticketsAddr := os.Getenv(fmt.Sprintf("NFT_TICKETS_%s_CONTRACT", c.network))
	if ticketsAddr == "" {
		// Use default addresses if available
		if c.network == "testnet" {
			ticketsAddr = os.Getenv("NFT_TICKETS_TESTNET_CONTRACT")
		} else {
			ticketsAddr = os.Getenv("NFT_TICKETS_MAINNET_CONTRACT")
		}
	}

	if ticketsAddr == "" {
		return fmt.Errorf("tickets contract address not configured for %s", c.network)
	}

	// Load BOGOWITickets contract
	c.ticketsAddress = common.HexToAddress(ticketsAddr)
	ticketsContract, err := contracts.NewBOGOWITickets(c.ticketsAddress, c.ethClient)
	if err != nil {
		return fmt.Errorf("failed to load tickets contract: %w", err)
	}
	c.ticketsContract = NewContractAdapter(ticketsContract)

	// Load RoleManager if configured
	roleManagerAddr := os.Getenv(fmt.Sprintf("ROLE_MANAGER_%s", c.network))
	if roleManagerAddr != "" {
		c.roleManagerAddress = common.HexToAddress(roleManagerAddr)
		roleManager, err := contracts.NewRoleManager(c.roleManagerAddress, c.ethClient)
		if err != nil {
			return fmt.Errorf("failed to load role manager: %w", err)
		}
		c.roleManager = roleManager
	}

	return nil
}

// initDatakyteService initializes the Datakyte metadata service
func (c *Client) initDatakyteService() {
	apiKey := os.Getenv(fmt.Sprintf("DATAKYTE_%s_API_KEY", c.network))
	if apiKey == "" {
		return
	}

	contractAddr := c.ticketsAddress.Hex()
	chainID := int(c.chainID.Int64())

	c.datakyteService = datakyte.NewTicketMetadataService(apiKey, contractAddr, chainID)
}

// GetAddress returns the client's Ethereum address
func (c *Client) GetAddress() common.Address {
	return c.auth.From
}

// GetBalance returns the CAM balance of the client
func (c *Client) GetBalance(ctx context.Context) (*big.Int, error) {
	return c.ethClient.BalanceAt(ctx, c.auth.From, nil)
}

// GetTicketsContract returns the tickets contract instance
func (c *Client) GetTicketsContract() TicketsContractInterface {
	return c.ticketsContract
}

// GetRoleManager returns the role manager contract instance
func (c *Client) GetRoleManager() *contracts.RoleManager {
	return c.roleManager
}

// EstimateGas estimates gas for a transaction
func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	gas, err := c.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Apply multiplier for safety
	gasWithBuffer := float64(gas) * c.config.GasMultiplier
	return uint64(gasWithBuffer), nil
}

// WaitForTransaction waits for a transaction to be confirmed
func (c *Client) WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	// Create a ticker for checking transaction status
	ticker := time.NewTicker(c.networkConfig.BlockTime)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("transaction timeout")
		case <-ticker.C:
			receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
			if err != nil {
				continue // Transaction not yet mined
			}

			// Check if we need to wait for more confirmations
			if c.networkConfig.ConfirmationWait > 0 {
				currentBlock, err := c.ethClient.BlockNumber(ctx)
				if err != nil {
					continue
				}

				confirmations := currentBlock - receipt.BlockNumber.Uint64()
				if confirmations < uint64(c.networkConfig.ConfirmationWait) {
					continue
				}
			}

			// Check transaction status
			if receipt.Status == 0 {
				return receipt, fmt.Errorf("transaction failed")
			}

			return receipt, nil
		}
	}
}

// UpdateAuth updates transaction auth settings
func (c *Client) UpdateAuth(gasLimit uint64, gasPrice *big.Int, nonce *big.Int) {
	if gasLimit > 0 {
		c.auth.GasLimit = gasLimit
	}
	if gasPrice != nil {
		c.auth.GasPrice = gasPrice
	}
	if nonce != nil {
		c.auth.Nonce = nonce
	}
}

// ResetAuth resets auth to default values
func (c *Client) ResetAuth() {
	c.auth.GasLimit = 0
	c.auth.GasPrice = nil
	c.auth.Nonce = nil
}

// GetNonce returns the current nonce for the client's address
func (c *Client) GetNonce(ctx context.Context) (uint64, error) {
	return c.ethClient.PendingNonceAt(ctx, c.auth.From)
}

// SuggestGasPrice returns a suggested gas price
func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// Apply max gas price limit if configured
	if c.config.MaxGasPrice != nil && gasPrice.Cmp(c.config.MaxGasPrice) > 0 {
		return c.config.MaxGasPrice, nil
	}

	return gasPrice, nil
}

// Close closes the client connections
func (c *Client) Close() {
	if c.ethClient != nil {
		c.ethClient.Close()
	}
}
