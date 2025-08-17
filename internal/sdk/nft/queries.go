package nft

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// GetTicketData retrieves the on-chain data for a ticket
func (c *Client) GetTicketData(ctx context.Context, tokenID uint64) (*TicketData, error) {
	opts := &bind.CallOpts{Context: ctx}
	
	data, err := c.ticketsContract.GetTicketData(opts, new(big.Int).SetUint64(tokenID))
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket data: %w", err)
	}

	return &TicketData{
		BookingID:                  data.BookingId,
		EventID:                    data.EventId,
		TransferUnlockAt:           data.TransferUnlockAt,
		ExpiresAt:                  data.ExpiresAt,
		UtilityFlags:               data.UtilityFlags,
		State:                      data.State,
		NonTransferableAfterRedeem: data.NonTransferableAfterRedeem,
		BurnOnRedeem:               data.BurnOnRedeem,
	}, nil
}

// IsTransferable checks if a ticket can be transferred
func (c *Client) IsTransferable(ctx context.Context, tokenID uint64) (bool, error) {
	opts := &bind.CallOpts{Context: ctx}
	
	transferable, err := c.ticketsContract.IsTransferable(opts, new(big.Int).SetUint64(tokenID))
	if err != nil {
		return false, fmt.Errorf("failed to check transferability: %w", err)
	}

	return transferable, nil
}

// GetOwnerOf returns the owner of a ticket
func (c *Client) GetOwnerOf(ctx context.Context, tokenID uint64) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}
	
	owner, err := c.ticketsContract.OwnerOf(opts, new(big.Int).SetUint64(tokenID))
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get owner: %w", err)
	}

	return owner, nil
}

// GetTokenURI returns the metadata URI for a ticket
func (c *Client) GetTokenURI(ctx context.Context, tokenID uint64) (string, error) {
	opts := &bind.CallOpts{Context: ctx}
	
	uri, err := c.ticketsContract.TokenURI(opts, new(big.Int).SetUint64(tokenID))
	if err != nil {
		return "", fmt.Errorf("failed to get token URI: %w", err)
	}

	return uri, nil
}

// GetBalanceOf returns the number of tickets owned by an address
func (c *Client) GetBalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	
	balance, err := c.ticketsContract.BalanceOf(opts, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

// GetTotalSupply returns the total number of tickets minted
// Note: BOGOWITickets doesn't have a totalSupply method, would need event filtering to track
func (c *Client) GetTotalSupply(ctx context.Context) (*big.Int, error) {
	// BOGOWITickets contract doesn't expose totalSupply
	// This would need to be tracked via events or a separate counter
	return nil, fmt.Errorf("getTotalSupply requires event filtering to track total minted tokens (implement in Phase 7)")
}

// GetUserTickets returns all ticket IDs owned by a user
// Note: This requires either ERC721Enumerable support or event filtering
func (c *Client) GetUserTickets(ctx context.Context, owner common.Address) ([]uint64, error) {
	// First check if we have enumeration support by trying to get balance
	balance, err := c.GetBalanceOf(ctx, owner)
	if err != nil {
		return nil, err
	}

	if balance.Int64() == 0 {
		return []uint64{}, nil
	}

	// Since TokenOfOwnerByIndex might not be available, we need to use events
	// This is a placeholder that indicates event filtering is needed
	return nil, fmt.Errorf("getUserTickets requires event filtering to enumerate tokens (implement in Phase 7). Balance: %s tokens", balance.String())
}

// GetTicketsByEvent returns all tickets for a specific event
func (c *Client) GetTicketsByEvent(ctx context.Context, eventID [32]byte) ([]uint64, error) {
	// This would typically be done by filtering events
	// For now, return an error indicating it needs event filtering
	return nil, fmt.Errorf("event-based queries require event filtering (Phase 7)")
}

// GetActiveTickets returns all active (non-expired, non-redeemed) tickets for an owner
// Note: This currently requires event filtering for full functionality
func (c *Client) GetActiveTickets(ctx context.Context, owner common.Address) ([]uint64, error) {
	// Check balance first
	balance, err := c.GetBalanceOf(ctx, owner)
	if err != nil {
		return nil, err
	}

	if balance.Int64() == 0 {
		return []uint64{}, nil
	}

	// Without enumeration, we need event filtering
	return nil, fmt.Errorf("getActiveTickets requires event filtering (implement in Phase 7). User has %s tickets", balance.String())
}

// GetRedeemedTickets returns all redeemed tickets for an owner
// Note: This currently requires event filtering for full functionality
func (c *Client) GetRedeemedTickets(ctx context.Context, owner common.Address) ([]uint64, error) {
	// Without enumeration, we need event filtering
	return nil, fmt.Errorf("getRedeemedTickets requires event filtering (implement in Phase 7)")
}

// GetExpiredTickets returns all expired tickets for an owner
// Note: This currently requires event filtering for full functionality
func (c *Client) GetExpiredTickets(ctx context.Context, owner common.Address) ([]uint64, error) {
	// Without enumeration, we need event filtering
	return nil, fmt.Errorf("getExpiredTickets requires event filtering (implement in Phase 7)")
}

// GetTransferableTickets returns all tickets that can currently be transferred
// Note: This currently requires event filtering for full functionality
func (c *Client) GetTransferableTickets(ctx context.Context, owner common.Address) ([]uint64, error) {
	// Check balance first
	balance, err := c.GetBalanceOf(ctx, owner)
	if err != nil {
		return nil, err
	}

	if balance.Int64() == 0 {
		return []uint64{}, nil
	}

	// Without enumeration, we need event filtering
	return nil, fmt.Errorf("getTransferableTickets requires event filtering (implement in Phase 7). User has %s tickets", balance.String())
}

// GetTicketMetadata retrieves full metadata including Datakyte data
func (c *Client) GetTicketMetadata(ctx context.Context, tokenID uint64) (*TokenMetadata, error) {
	// Get on-chain URI
	uri, err := c.GetTokenURI(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	// Get on-chain data
	data, err := c.GetTicketData(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	metadata := &TokenMetadata{
		TokenID:     tokenID,
		ExternalURL: uri,
	}

	// Add attributes from on-chain data
	metadata.Attributes = []MetadataAttribute{
		{TraitType: "State", Value: ParseTicketState(data.State).String()},
		{TraitType: "Expires At", Value: data.ExpiresAt, DisplayType: "date"},
		{TraitType: "Transfer Unlock", Value: data.TransferUnlockAt, DisplayType: "date"},
		{TraitType: "Utility Flags", Value: data.UtilityFlags},
	}

	// Get Datakyte metadata if available
	if c.datakyteService != nil {
		datakyteMetadata, err := c.datakyteService.GetTicketMetadata(tokenID)
		if err == nil && datakyteMetadata != nil {
			metadata.Name = datakyteMetadata.Name
			metadata.Description = datakyteMetadata.Description
			metadata.Image = datakyteMetadata.Image
			
			// Extract conservation impact from attributes
			for _, attr := range datakyteMetadata.Attributes {
				if attr.TraitType == "Conservation Impact" {
					if impact, ok := attr.Value.(string); ok {
						metadata.ConservationImpact = impact
					}
				}
				
				// Add all attributes to metadata
				metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
					TraitType:   attr.TraitType,
					Value:       attr.Value,
					DisplayType: attr.DisplayType,
				})
			}
		}
	}

	return metadata, nil
}

// HasRole checks if an address has a specific role
func (c *Client) HasRole(ctx context.Context, role [32]byte, account common.Address) (bool, error) {
	if c.roleManager == nil {
		return false, fmt.Errorf("role manager not configured")
	}

	opts := &bind.CallOpts{Context: ctx}
	hasRole, err := c.roleManager.HasRole(opts, role, account)
	if err != nil {
		return false, fmt.Errorf("failed to check role: %w", err)
	}

	return hasRole, nil
}

// GetRedemptionNonce gets the redemption nonce for an address
// Note: This would need to be tracked via events or stored separately
func (c *Client) GetRedemptionNonce(ctx context.Context, address common.Address) (*big.Int, error) {
	// The contract doesn't expose redemptionNonces publicly
	// In practice, nonces would be tracked off-chain or via events
	// For now, we can use the current timestamp as a simple nonce
	return big.NewInt(time.Now().Unix()), nil
}

// String returns the string representation of a ticket state
func (s TicketState) String() string {
	switch s {
	case TicketStateIssued:
		return "Issued"
	case TicketStateRedeemed:
		return "Redeemed"
	case TicketStateExpired:
		return "Expired"
	default:
		return "Unknown"
	}
}