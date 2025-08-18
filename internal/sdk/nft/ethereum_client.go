package nft

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthereumClient defines the interface for Ethereum client operations
// This allows us to mock the client in tests
type EthereumClient interface {
	// Gas and nonce operations
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	
	// Transaction operations
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	
	// Account operations
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	
	// Contract operations
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	
	// Block operations
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	
	// Subscription operations
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	
	// Network operations
	NetworkID(ctx context.Context) (*big.Int, error)
	ChainID(ctx context.Context) (*big.Int, error)
}

// Ensure ethclient.Client implements our interface
// This is a compile-time check
// var _ EthereumClient = (*ethclient.Client)(nil)