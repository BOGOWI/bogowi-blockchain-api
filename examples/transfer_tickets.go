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
	
	// Example recipient addresses
	recipient1 := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7")
	recipient2 := common.HexToAddress("0x5aAeb6053f3E94C9b9A09f33669435E7Ef1BeAed")
	operator := common.HexToAddress("0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359")

	// Example 1: Simple transfer
	fmt.Println("=== Simple Transfer ===")
	tokenID := uint64(10001)
	
	// Check if ticket is transferable first
	transferable, err := client.IsTransferable(ctx, tokenID)
	if err != nil {
		fmt.Printf("Failed to check transferability: %v\n", err)
	} else if !transferable {
		fmt.Printf("Ticket #%d is not transferable yet\n", tokenID)
	} else {
		tx, err := client.Transfer(ctx, recipient1, tokenID)
		if err != nil {
			fmt.Printf("Failed to transfer: %v\n", err)
		} else {
			fmt.Printf("✓ Transferred ticket #%d to %s\n", tokenID, recipient1.Hex())
			fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
		}
	}

	// Example 2: Safe transfer
	fmt.Println("\n=== Safe Transfer ===")
	tokenID2 := uint64(10002)
	
	tx, err := client.SafeTransfer(ctx, recipient2, tokenID2)
	if err != nil {
		fmt.Printf("Failed to safe transfer: %v\n", err)
	} else {
		fmt.Printf("✓ Safely transferred ticket #%d to %s\n", tokenID2, recipient2.Hex())
		fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
	}

	// Example 3: Safe transfer with data
	fmt.Println("\n=== Safe Transfer with Data ===")
	tokenID3 := uint64(10003)
	data := []byte("Custom transfer data")
	
	tx, err = client.SafeTransferWithData(ctx, recipient1, tokenID3, data)
	if err != nil {
		fmt.Printf("Failed to safe transfer with data: %v\n", err)
	} else {
		fmt.Printf("✓ Safely transferred ticket #%d with data\n", tokenID3)
		fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
	}

	// Example 4: Approve another address to transfer
	fmt.Println("\n=== Approval ===")
	tokenID4 := uint64(10004)
	
	tx, err = client.Approve(ctx, operator, tokenID4)
	if err != nil {
		fmt.Printf("Failed to approve: %v\n", err)
	} else {
		fmt.Printf("✓ Approved %s to transfer ticket #%d\n", operator.Hex(), tokenID4)
		fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
	}

	// Check approval
	approved, err := client.GetApproved(ctx, tokenID4)
	if err != nil {
		fmt.Printf("Failed to get approved: %v\n", err)
	} else {
		fmt.Printf("  Approved address: %s\n", approved.Hex())
	}

	// Example 5: Set approval for all tickets
	fmt.Println("\n=== Approval for All ===")
	
	tx, err = client.SetApprovalForAll(ctx, operator, true)
	if err != nil {
		fmt.Printf("Failed to set approval for all: %v\n", err)
	} else {
		fmt.Printf("✓ Approved %s for all tickets\n", operator.Hex())
		fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
	}

	// Check approval for all
	isApprovedForAll, err := client.IsApprovedForAll(ctx, client.GetAddress(), operator)
	if err != nil {
		fmt.Printf("Failed to check approval for all: %v\n", err)
	} else {
		if isApprovedForAll {
			fmt.Printf("  ✓ Operator is approved for all tickets\n")
		} else {
			fmt.Printf("  ✗ Operator is not approved for all tickets\n")
		}
	}

	// Example 6: Revoke approval for all
	fmt.Println("\n=== Revoke Approval for All ===")
	
	tx, err = client.SetApprovalForAll(ctx, operator, false)
	if err != nil {
		fmt.Printf("Failed to revoke approval for all: %v\n", err)
	} else {
		fmt.Printf("✓ Revoked approval for all tickets from %s\n", operator.Hex())
		fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
	}

	// Example 7: Batch transfer multiple tickets
	fmt.Println("\n=== Batch Transfer ===")
	tokenIDs := []uint64{10005, 10006, 10007}
	
	txs, err := client.BatchTransfer(ctx, recipient1, tokenIDs)
	if err != nil {
		fmt.Printf("Batch transfer error: %v\n", err)
		fmt.Printf("Successful transfers: %d\n", len(txs))
	} else {
		fmt.Printf("✓ Batch transferred %d tickets to %s\n", len(tokenIDs), recipient1.Hex())
		for i, tx := range txs {
			fmt.Printf("  Token %d: %s\n", tokenIDs[i], tx.Hash().Hex())
		}
	}

	// Example 8: Transfer to multiple recipients
	fmt.Println("\n=== Transfer to Multiple Recipients ===")
	transfers := map[common.Address][]uint64{
		recipient1: {10008, 10009},
		recipient2: {10010, 10011, 10012},
	}
	
	txs, err = client.TransferToMultiple(ctx, transfers)
	if err != nil {
		fmt.Printf("Multiple transfer error: %v\n", err)
		fmt.Printf("Successful transfers: %d\n", len(txs))
	} else {
		fmt.Printf("✓ Transferred tickets to multiple recipients\n")
		fmt.Printf("  Total transfers: %d\n", len(txs))
	}

	// Example 9: Transfer from (as approved operator)
	fmt.Println("\n=== Transfer From (as Operator) ===")
	// This would typically be done by the approved operator
	// For demo purposes, we'll show the structure
	
	fromAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	tokenID5 := uint64(10013)
	
	// Note: This would fail unless the client's address is approved
	tx, err = client.TransferFrom(ctx, fromAddress, recipient1, tokenID5)
	if err != nil {
		fmt.Printf("Transfer from failed (expected): %v\n", err)
	} else {
		fmt.Printf("✓ Transferred ticket #%d from %s to %s\n", tokenID5, fromAddress.Hex(), recipient1.Hex())
		fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
	}

	fmt.Println("\n=== Transfer Examples Completed ===")
}