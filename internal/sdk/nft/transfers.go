package nft

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Transfer transfers a ticket to another address
func (c *Client) Transfer(ctx context.Context, to common.Address, tokenID uint64) (*types.Transaction, error) {
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

// SafeTransfer safely transfers a ticket to another address
func (c *Client) SafeTransfer(ctx context.Context, to common.Address, tokenID uint64) (*types.Transaction, error) {
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

// SafeTransferWithData safely transfers a ticket with additional data
func (c *Client) SafeTransferWithData(ctx context.Context, to common.Address, tokenID uint64, data []byte) (*types.Transaction, error) {
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

// Approve approves another address to transfer a specific ticket
func (c *Client) Approve(ctx context.Context, spender common.Address, tokenID uint64) (*types.Transaction, error) {
	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	// Verify ownership
	owner, err := c.GetOwnerOf(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}

	if owner != c.auth.From {
		return nil, fmt.Errorf("sender is not the owner of token %d", tokenID)
	}

	// Execute approval
	tx, err := c.ticketsContract.Approve(
		c.auth,
		spender,
		new(big.Int).SetUint64(tokenID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to approve: %w", err)
	}

	// Wait for confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 0 {
		return tx, fmt.Errorf("approval transaction failed")
	}

	c.ResetAuth()
	return tx, nil
}

// SetApprovalForAll approves or revokes approval for an operator to manage all tickets
func (c *Client) SetApprovalForAll(ctx context.Context, operator common.Address, approved bool) (*types.Transaction, error) {
	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	// Execute approval for all
	tx, err := c.ticketsContract.SetApprovalForAll(
		c.auth,
		operator,
		approved,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set approval for all: %w", err)
	}

	// Wait for confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 0 {
		return tx, fmt.Errorf("set approval for all transaction failed")
	}

	c.ResetAuth()
	return tx, nil
}

// GetApproved returns the approved address for a specific ticket
func (c *Client) GetApproved(ctx context.Context, tokenID uint64) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}

	approved, err := c.ticketsContract.GetApproved(opts, new(big.Int).SetUint64(tokenID))
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get approved: %w", err)
	}

	return approved, nil
}

// IsApprovedForAll checks if an operator is approved to manage all tickets for an owner
func (c *Client) IsApprovedForAll(ctx context.Context, owner, operator common.Address) (bool, error) {
	opts := &bind.CallOpts{Context: ctx}

	approved, err := c.ticketsContract.IsApprovedForAll(opts, owner, operator)
	if err != nil {
		return false, fmt.Errorf("failed to check approval for all: %w", err)
	}

	return approved, nil
}

// TransferFrom transfers a ticket from one address to another (requires approval)
func (c *Client) TransferFrom(ctx context.Context, from, to common.Address, tokenID uint64) (*types.Transaction, error) {
	// Check if ticket is transferable
	transferable, err := c.IsTransferable(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to check transferability: %w", err)
	}
	if !transferable {
		return nil, ErrNotTransferable
	}

	// Verify approval
	approved, err := c.GetApproved(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved: %w", err)
	}

	isApprovedForAll, err := c.IsApprovedForAll(ctx, from, c.auth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to check approval for all: %w", err)
	}

	if approved != c.auth.From && !isApprovedForAll {
		return nil, fmt.Errorf("sender is not approved to transfer token %d", tokenID)
	}

	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	// Execute transfer
	tx, err := c.ticketsContract.TransferFrom(
		c.auth,
		from,
		to,
		new(big.Int).SetUint64(tokenID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer from: %w", err)
	}

	// Wait for confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, fmt.Errorf("transaction failed: %w", err)
	}

	if receipt.Status == 0 {
		return tx, fmt.Errorf("transfer from transaction failed")
	}

	c.ResetAuth()
	return tx, nil
}

// BatchTransfer transfers multiple tickets to the same recipient
func (c *Client) BatchTransfer(ctx context.Context, to common.Address, tokenIDs []uint64) ([]*types.Transaction, error) {
	if len(tokenIDs) == 0 {
		return nil, fmt.Errorf("no token IDs provided")
	}

	txs := make([]*types.Transaction, 0, len(tokenIDs))

	for _, tokenID := range tokenIDs {
		tx, err := c.Transfer(ctx, to, tokenID)
		if err != nil {
			// Return partial results on error
			return txs, fmt.Errorf("failed to transfer token %d: %w", tokenID, err)
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// TransferToMultiple transfers tickets to multiple recipients
func (c *Client) TransferToMultiple(ctx context.Context, transfers map[common.Address][]uint64) ([]*types.Transaction, error) {
	if len(transfers) == 0 {
		return nil, fmt.Errorf("no transfers provided")
	}

	var txs []*types.Transaction

	for recipient, tokenIDs := range transfers {
		for _, tokenID := range tokenIDs {
			tx, err := c.Transfer(ctx, recipient, tokenID)
			if err != nil {
				// Return partial results on error
				return txs, fmt.Errorf("failed to transfer token %d to %s: %w", tokenID, recipient.Hex(), err)
			}
			txs = append(txs, tx)
		}
	}

	return txs, nil
}
