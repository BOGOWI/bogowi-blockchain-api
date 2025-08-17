package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"bogowi-blockchain-go/internal/sdk/nft"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// Load configuration from environment
	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY environment variable not set")
	}

	// Configure SDK client
	config := nft.ClientConfig{
		PrivateKey:      privateKey,
		Network:         "testnet", // or "mainnet"
		DatakyteEnabled: true,
	}

	// Create SDK client
	client, err := nft.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Get ticket data
	fmt.Println("=== Getting Ticket Data ===")
	tokenID := uint64(10001)
	
	ticketData, err := client.GetTicketData(ctx, tokenID)
	if err != nil {
		fmt.Printf("Failed to get ticket data: %v\n", err)
	} else {
		fmt.Printf("Ticket #%d Data:\n", tokenID)
		fmt.Printf("  State: %s\n", nft.ParseTicketState(ticketData.State))
		fmt.Printf("  Expires At: %d\n", ticketData.ExpiresAt)
		fmt.Printf("  Transfer Unlock: %d\n", ticketData.TransferUnlockAt)
		fmt.Printf("  Utility Flags: %d\n", ticketData.UtilityFlags)
	}

	// Example 2: Check if ticket is transferable
	fmt.Println("\n=== Checking Transferability ===")
	transferable, err := client.IsTransferable(ctx, tokenID)
	if err != nil {
		fmt.Printf("Failed to check transferability: %v\n", err)
	} else {
		if transferable {
			fmt.Printf("✓ Ticket #%d is transferable\n", tokenID)
		} else {
			fmt.Printf("✗ Ticket #%d is NOT transferable\n", tokenID)
		}
	}

	// Example 3: Get owner of a ticket
	fmt.Println("\n=== Getting Ticket Owner ===")
	owner, err := client.GetOwnerOf(ctx, tokenID)
	if err != nil {
		fmt.Printf("Failed to get owner: %v\n", err)
	} else {
		fmt.Printf("Ticket #%d is owned by: %s\n", tokenID, owner.Hex())
	}

	// Example 4: Get user's balance
	fmt.Println("\n=== Getting User Balance ===")
	userAddress := client.GetAddress()
	balance, err := client.GetBalanceOf(ctx, userAddress)
	if err != nil {
		fmt.Printf("Failed to get balance: %v\n", err)
	} else {
		fmt.Printf("User %s owns %s tickets\n", userAddress.Hex(), balance.String())
	}

	// Example 5: Get all user's tickets
	fmt.Println("\n=== Getting User's Tickets ===")
	tickets, err := client.GetUserTickets(ctx, userAddress)
	if err != nil {
		fmt.Printf("Note: %v\n", err)
		fmt.Println("(Contract may not support enumeration)")
	} else {
		fmt.Printf("User owns tickets: %v\n", tickets)
	}

	// Example 6: Get active tickets
	fmt.Println("\n=== Getting Active Tickets ===")
	activeTickets, err := client.GetActiveTickets(ctx, userAddress)
	if err != nil {
		fmt.Printf("Failed to get active tickets: %v\n", err)
	} else {
		fmt.Printf("Active tickets: %v\n", activeTickets)
	}

	// Example 7: Get ticket metadata (including Datakyte)
	fmt.Println("\n=== Getting Full Ticket Metadata ===")
	metadata, err := client.GetTicketMetadata(ctx, tokenID)
	if err != nil {
		fmt.Printf("Failed to get metadata: %v\n", err)
	} else {
		fmt.Printf("Ticket #%d Metadata:\n", tokenID)
		fmt.Printf("  Name: %s\n", metadata.Name)
		fmt.Printf("  Description: %s\n", metadata.Description)
		fmt.Printf("  Image: %s\n", metadata.Image)
		fmt.Printf("  Datakyte ID: %s\n", metadata.DatakyteID)
		fmt.Printf("  Conservation Impact: %s\n", metadata.ConservationImpact)
		fmt.Printf("  Attributes: %d\n", len(metadata.Attributes))
		for _, attr := range metadata.Attributes {
			fmt.Printf("    - %s: %v\n", attr.TraitType, attr.Value)
		}
	}

	// Example 8: Get token URI
	fmt.Println("\n=== Getting Token URI ===")
	uri, err := client.GetTokenURI(ctx, tokenID)
	if err != nil {
		fmt.Printf("Failed to get token URI: %v\n", err)
	} else {
		fmt.Printf("Token URI: %s\n", uri)
	}

	// Example 9: Check role (if role manager is configured)
	fmt.Println("\n=== Checking Roles ===")
	minterRole := [32]byte{} // MINTER_ROLE hash
	copy(minterRole[:], []byte("MINTER_ROLE"))
	
	hasMinterRole, err := client.HasRole(ctx, minterRole, userAddress)
	if err != nil {
		fmt.Printf("Role check failed: %v\n", err)
	} else {
		if hasMinterRole {
			fmt.Printf("✓ User has MINTER_ROLE\n")
		} else {
			fmt.Printf("✗ User does not have MINTER_ROLE\n")
		}
	}

	// Example 10: Get redemption nonce
	fmt.Println("\n=== Getting Redemption Nonce ===")
	nonce, err := client.GetRedemptionNonce(ctx, userAddress)
	if err != nil {
		fmt.Printf("Failed to get redemption nonce: %v\n", err)
	} else {
		fmt.Printf("Current redemption nonce: %s\n", nonce.String())
	}

	// Example 11: Query another user's tickets
	fmt.Println("\n=== Querying Another User's Tickets ===")
	otherUser := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7")
	
	otherBalance, err := client.GetBalanceOf(ctx, otherUser)
	if err != nil {
		fmt.Printf("Failed to get other user's balance: %v\n", err)
	} else {
		fmt.Printf("User %s owns %s tickets\n", otherUser.Hex(), otherBalance.String())
	}

	fmt.Println("\n=== Query Examples Completed ===")
}