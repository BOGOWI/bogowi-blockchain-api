package nft

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TicketsContractInterface defines the interface for BOGOWITickets contract interactions
type TicketsContractInterface interface {
	TransferFrom(opts *bind.TransactOpts, from, to common.Address, tokenID *big.Int) (*types.Transaction, error)
	SafeTransferFrom(opts *bind.TransactOpts, from, to common.Address, tokenID *big.Int) (*types.Transaction, error)
	SafeTransferFrom0(opts *bind.TransactOpts, from, to common.Address, tokenID *big.Int, data []byte) (*types.Transaction, error)
	Approve(opts *bind.TransactOpts, spender common.Address, tokenID *big.Int) (*types.Transaction, error)
	SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error)
	GetApproved(opts *bind.CallOpts, tokenID *big.Int) (common.Address, error)
	IsApprovedForAll(opts *bind.CallOpts, owner, operator common.Address) (bool, error)
	IsTransferable(opts *bind.CallOpts, tokenID *big.Int) (bool, error)
	OwnerOf(opts *bind.CallOpts, tokenID *big.Int) (common.Address, error)
	GetTicketData(opts *bind.CallOpts, tokenID *big.Int) (TicketDataContract, error)
	TokenURI(opts *bind.CallOpts, tokenID *big.Int) (string, error)
	BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error)
	MintTicket(opts *bind.TransactOpts, to common.Address, bookingID [32]byte, eventID [32]byte, utilityFlags uint32, transferUnlockAt uint64, expiresAt uint64, metadataURI string, rewardBasisPoints uint16) (*types.Transaction, error)
	MintBatch(opts *bind.TransactOpts, tos []common.Address, bookingIDs [][32]byte, eventIDs [][32]byte, utilityFlags []uint32, transferUnlockAts []uint64, expiresAts []uint64, metadataURIs []string, rewardBasisPoints []uint16) (*types.Transaction, error)
	SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error)
	ParseTicketMinted(log types.Log) (*TicketMintedEvent, error)
	ExpireTicket(opts *bind.TransactOpts, tokenID *big.Int) (*types.Transaction, error)
	RedeemTicket(opts *bind.TransactOpts, redemptionData RedemptionDataContract) (*types.Transaction, error)
	UpdateTransferUnlock(opts *bind.TransactOpts, tokenID *big.Int, newUnlockTime uint64) (*types.Transaction, error)
	Burn(opts *bind.TransactOpts, tokenID *big.Int) (*types.Transaction, error)
}

// TicketDataContract represents the contract return type for GetTicketData
type TicketDataContract struct {
	BookingID                  [32]byte
	EventID                    [32]byte
	TransferUnlockAt           uint64
	ExpiresAt                  uint64
	UtilityFlags               uint32
	State                      uint8
	NonTransferableAfterRedeem bool
	BurnOnRedeem               bool
}

// TicketMintedEvent represents the TicketMinted event
type TicketMintedEvent struct {
	TokenID   *big.Int
	To        common.Address
	BookingID [32]byte
	EventID   [32]byte
	Raw       types.Log
}

// RedemptionDataContract represents the contract parameter for RedeemTicket
type RedemptionDataContract struct {
	TokenID   *big.Int
	Redeemer  common.Address
	Nonce     *big.Int
	Deadline  *big.Int
	ChainID   *big.Int
	Signature []byte
}

// EthClientInterface defines the interface for Ethereum client interactions
type EthClientInterface interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// TestClient is a test version of Client that uses interfaces
type TestClient struct {
	ticketsContract TicketsContractInterface
	ethClient       EthClientInterface
	auth            *bind.TransactOpts
	roleManager     interface{}
	datakyteService interface{}
}

// Convert TestClient methods to work with interfaces
func (c *TestClient) Transfer(ctx context.Context, to common.Address, tokenID uint64) (*types.Transaction, error) {
	// Check if ticket is transferable first
	transferable, err := c.IsTransferable(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to check transferability: %w", err)
	}
	if !transferable {
		return nil, ErrNotTransferable
	}

	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	// Get the current owner to use as 'from' address
	owner, err := c.GetOwnerOf(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current owner: %w", err)
	}

	// Ensure the transaction sender is the owner
	if owner != c.auth.From {
		return nil, fmt.Errorf("sender is not the owner of token %d", tokenID)
	}

	// Execute transfer
	tx, err := c.ticketsContract.TransferFrom(
		c.auth,
		owner,
		to,
		new(big.Int).SetUint64(tokenID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer ticket: %w", err)
	}

	// Wait for confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 0 {
		return tx, fmt.Errorf("transfer transaction failed")
	}

	c.ResetAuth()
	return tx, nil
}

func (c *TestClient) SafeTransfer(ctx context.Context, to common.Address, tokenID uint64) (*types.Transaction, error) {
	// Check if ticket is transferable
	transferable, err := c.IsTransferable(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to check transferability: %w", err)
	}
	if !transferable {
		return nil, ErrNotTransferable
	}

	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	// Get the current owner
	owner, err := c.GetOwnerOf(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current owner: %w", err)
	}

	if owner != c.auth.From {
		return nil, fmt.Errorf("sender is not the owner of token %d", tokenID)
	}

	// Execute safe transfer
	tx, err := c.ticketsContract.SafeTransferFrom(
		c.auth,
		owner,
		to,
		new(big.Int).SetUint64(tokenID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to safe transfer ticket: %w", err)
	}

	// Wait for confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 0 {
		return tx, fmt.Errorf("safe transfer transaction failed")
	}

	c.ResetAuth()
	return tx, nil
}

func (c *TestClient) SafeTransferWithData(ctx context.Context, to common.Address, tokenID uint64, data []byte) (*types.Transaction, error) {
	// Check if ticket is transferable
	transferable, err := c.IsTransferable(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to check transferability: %w", err)
	}
	if !transferable {
		return nil, ErrNotTransferable
	}

	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	// Get the current owner
	owner, err := c.GetOwnerOf(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current owner: %w", err)
	}

	if owner != c.auth.From {
		return nil, fmt.Errorf("sender is not the owner of token %d", tokenID)
	}

	// Execute safe transfer with data
	tx, err := c.ticketsContract.SafeTransferFrom0(
		c.auth,
		owner,
		to,
		new(big.Int).SetUint64(tokenID),
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to safe transfer with data: %w", err)
	}

	// Wait for confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 0 {
		return tx, fmt.Errorf("safe transfer with data transaction failed")
	}

	c.ResetAuth()
	return tx, nil
}

func (c *TestClient) IsTransferable(ctx context.Context, tokenID uint64) (bool, error) {
	opts := &bind.CallOpts{Context: ctx}
	return c.ticketsContract.IsTransferable(opts, new(big.Int).SetUint64(tokenID))
}

func (c *TestClient) GetOwnerOf(ctx context.Context, tokenID uint64) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}
	return c.ticketsContract.OwnerOf(opts, new(big.Int).SetUint64(tokenID))
}

func (c *TestClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.ethClient.SuggestGasPrice(ctx)
}

func (c *TestClient) WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.ethClient.TransactionReceipt(ctx, txHash)
}

func (c *TestClient) ResetAuth() {
	// Reset auth state after transaction
}
