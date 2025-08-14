package nft

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"bogowi-blockchain-go/internal/sdk/contracts"
	"bogowi-blockchain-go/internal/services/datakyte"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

// MintTicket mints a new NFT ticket
func (c *Client) MintTicket(ctx context.Context, params MintParams) (*types.Transaction, uint64, error) {
	// Prepare contract parameters
	contractParams := contracts.IBOGOWITicketsMintParams{
		To:                params.To,
		BookingId:         params.BookingID,
		EventId:           params.EventID,
		UtilityFlags:      params.UtilityFlags,
		TransferUnlockAt:  params.TransferUnlockAt,
		ExpiresAt:         params.ExpiresAt,
		MetadataURI:       params.MetadataURI,
		RewardBasisPoints: new(big.Int).SetUint64(uint64(params.RewardBasisPoints)),
	}

	// Estimate gas
	opts := &bind.TransactOpts{
		From:    c.auth.From,
		Signer:  c.auth.Signer,
		Context: ctx,
		NoSend:  true, // Don't send, just estimate
	}

	// Create a copy of auth for gas estimation
	_, err := c.ticketsContract.MintTicket(opts, contractParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Set gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Update auth with gas settings
	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	// Send the actual transaction
	tx, err := c.ticketsContract.MintTicket(c.auth, contractParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to mint ticket: %w", err)
	}

	// Wait for transaction confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, 0, fmt.Errorf("transaction failed: %w", err)
	}

	// Extract token ID from events
	tokenID, err := c.extractTokenIDFromReceipt(receipt)
	if err != nil {
		return tx, 0, fmt.Errorf("failed to extract token ID: %w", err)
	}

	// Create Datakyte metadata if service is available
	if c.datakyteService != nil && params.DatakyteNftID != "" {
		go c.syncDatakyteMetadata(tokenID, params)
	}

	// Reset auth for next transaction
	c.ResetAuth()

	return tx, tokenID, nil
}

// BatchMint mints multiple tickets in a single transaction
func (c *Client) BatchMint(ctx context.Context, params []MintParams) (*types.Transaction, []uint64, error) {
	if len(params) == 0 {
		return nil, nil, fmt.Errorf("empty batch")
	}

	if len(params) > 100 {
		return nil, nil, fmt.Errorf("batch size exceeds maximum of 100")
	}

	// Convert to contract format
	contractParams := make([]contracts.IBOGOWITicketsMintParams, len(params))
	for i, p := range params {
		contractParams[i] = contracts.IBOGOWITicketsMintParams{
			To:                p.To,
			BookingId:         p.BookingID,
			EventId:           p.EventID,
			UtilityFlags:      p.UtilityFlags,
			TransferUnlockAt:  p.TransferUnlockAt,
			ExpiresAt:         p.ExpiresAt,
			MetadataURI:       p.MetadataURI,
			RewardBasisPoints: new(big.Int).SetUint64(uint64(p.RewardBasisPoints)),
		}
	}

	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Update auth
	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx
	c.auth.GasLimit = uint64(150000 * len(params)) // Estimate 150k gas per mint

	// Send transaction
	tx, err := c.ticketsContract.MintBatch(c.auth, contractParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to batch mint: %w", err)
	}

	// Wait for confirmation
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return tx, nil, fmt.Errorf("transaction failed: %w", err)
	}

	// Extract all token IDs
	tokenIDs, err := c.extractTokenIDsFromBatchReceipt(receipt, len(params))
	if err != nil {
		return tx, nil, fmt.Errorf("failed to extract token IDs: %w", err)
	}

	// Sync Datakyte metadata for all tokens
	if c.datakyteService != nil {
		for i, tokenID := range tokenIDs {
			if params[i].DatakyteNftID != "" {
				go c.syncDatakyteMetadata(tokenID, params[i])
			}
		}
	}

	c.ResetAuth()
	return tx, tokenIDs, nil
}

// extractTokenIDFromReceipt extracts the token ID from a mint transaction receipt
func (c *Client) extractTokenIDFromReceipt(receipt *types.Receipt) (uint64, error) {
	// Parse TicketMinted event
	for _, log := range receipt.Logs {
		event, err := c.ticketsContract.ParseTicketMinted(*log)
		if err == nil {
			return event.TokenId.Uint64(), nil
		}
	}
	return 0, fmt.Errorf("TicketMinted event not found")
}

// extractTokenIDsFromBatchReceipt extracts multiple token IDs from batch mint receipt
func (c *Client) extractTokenIDsFromBatchReceipt(receipt *types.Receipt, expectedCount int) ([]uint64, error) {
	tokenIDs := make([]uint64, 0, expectedCount)
	
	for _, log := range receipt.Logs {
		event, err := c.ticketsContract.ParseTicketMinted(*log)
		if err == nil {
			tokenIDs = append(tokenIDs, event.TokenId.Uint64())
		}
	}

	if len(tokenIDs) != expectedCount {
		return nil, fmt.Errorf("expected %d tokens, got %d", expectedCount, len(tokenIDs))
	}

	return tokenIDs, nil
}

// syncDatakyteMetadata syncs ticket metadata with Datakyte
func (c *Client) syncDatakyteMetadata(tokenID uint64, params MintParams) {
	if c.datakyteService == nil {
		return
	}

	// Create Datakyte metadata
	ticketData := datakyte.BOGOWITicketData{
		TokenID:           tokenID,
		BookingID:         string(params.BookingID[:]),
		EventID:           string(params.EventID[:]),
		ExperienceTitle:   "BOGOWI Experience", // These would come from additional params
		ExperienceType:    "Eco-Adventure",
		Location:          "TBD",
		Duration:          "1 Day",
		ValidUntil:        time.Unix(int64(params.ExpiresAt), 0),
		TransferableAfter: time.Unix(int64(params.TransferUnlockAt), 0),
		ExpiresAt:         time.Unix(int64(params.ExpiresAt), 0),
		BOGORewards:       int(params.RewardBasisPoints),
		RecipientAddress:  params.To.Hex(),
	}

	_, err := c.datakyteService.CreateTicketMetadata(ticketData)
	if err != nil {
		// Log error but don't fail the transaction
		fmt.Printf("Warning: Failed to sync Datakyte metadata for token %d: %v\n", tokenID, err)
	}
}

// SetBaseURI updates the base URI for Datakyte metadata
func (c *Client) SetBaseURI(ctx context.Context, baseURI string) (*types.Transaction, error) {
	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	tx, err := c.ticketsContract.SetBaseURI(c.auth, baseURI)
	if err != nil {
		return nil, fmt.Errorf("failed to set base URI: %w", err)
	}

	c.ResetAuth()
	return tx, nil
}

// ExpireTicket marks a ticket as expired
func (c *Client) ExpireTicket(ctx context.Context, tokenID uint64) (*types.Transaction, error) {
	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	tx, err := c.ticketsContract.ExpireTicket(c.auth, new(big.Int).SetUint64(tokenID))
	if err != nil {
		return nil, fmt.Errorf("failed to expire ticket: %w", err)
	}

	// Update Datakyte status if available
	if c.datakyteService != nil {
		go func() {
			// Note: Would need to store Datakyte NFT ID mapping
			fmt.Printf("TODO: Update Datakyte status for expired ticket %d\n", tokenID)
		}()
	}

	c.ResetAuth()
	return tx, nil
}

// UpdateTransferUnlock updates the transfer unlock time for a ticket
func (c *Client) UpdateTransferUnlock(ctx context.Context, tokenID uint64, newUnlockTime uint64) (*types.Transaction, error) {
	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	tx, err := c.ticketsContract.UpdateTransferUnlock(
		c.auth, 
		new(big.Int).SetUint64(tokenID), 
		newUnlockTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update transfer unlock: %w", err)
	}

	c.ResetAuth()
	return tx, nil
}

// Burn burns a ticket NFT
func (c *Client) Burn(ctx context.Context, tokenID uint64) (*types.Transaction, error) {
	// Get gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	c.auth.GasPrice = gasPrice
	c.auth.Context = ctx

	tx, err := c.ticketsContract.Burn(c.auth, new(big.Int).SetUint64(tokenID))
	if err != nil {
		return nil, fmt.Errorf("failed to burn ticket: %w", err)
	}

	c.ResetAuth()
	return tx, nil
}