package sdk

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthClient is an interface for Ethereum client operations
type EthClient interface {
	ChainID(ctx context.Context) (*big.Int, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	Close()
}

// BoundContract is an interface for bound contract operations
type BoundContract interface {
	Call(opts *bind.CallOpts, results *[]interface{}, method string, params ...interface{}) error
	Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error)
}

// BoundContractWrapper wraps *bind.BoundContract to implement BoundContract interface
type BoundContractWrapper struct {
	*bind.BoundContract
}

// Call implements BoundContract
func (w *BoundContractWrapper) Call(opts *bind.CallOpts, results *[]interface{}, method string, params ...interface{}) error {
	return w.BoundContract.Call(opts, results, method, params...)
}

// Transact implements BoundContract
func (w *BoundContractWrapper) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return w.BoundContract.Transact(opts, method, params...)
}