package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Example 1: Generate a redemption QR code
	fmt.Println("Generating redemption QR code...")
	
	tokenID := uint64(10001) // The ticket token ID
	redeemer := client.GetAddress() // Redeemer address (could be different)
	
	qrData, err := client.GenerateRedemptionQR(tokenID, redeemer)
	if err != nil {
		log.Fatalf("Failed to generate QR code: %v", err)
	}
	
	fmt.Printf("QR Code Data: %s\n", qrData)
	fmt.Println("This QR code is valid for 5 minutes")
	fmt.Println()

	// Example 2: Redeem a ticket with explicit parameters
	fmt.Println("Redeeming ticket...")
	
	params := nft.RedemptionParams{
		TokenID:  tokenID,
		Redeemer: redeemer,
		Nonce:    uint64(time.Now().Unix()),
		Deadline: time.Now().Add(5 * time.Minute).Unix(),
	}

	ctx := context.Background()
	tx, err := client.RedeemTicket(ctx, params)
	if err != nil {
		log.Fatalf("Failed to redeem ticket: %v", err)
	}

	fmt.Printf("Redemption transaction sent: %s\n", tx.Hash().Hex())
	fmt.Println("Ticket successfully redeemed!")

	// Example 3: Demonstrate signature verification would be done on-chain
	fmt.Println("\nNote: Signature verification happens on-chain during redemption")
	fmt.Println("The smart contract verifies the EIP-712 signature automatically")

	// Example 4: Redeem with a different redeemer address
	fmt.Println("\nRedeeming for a different user...")
	
	differentRedeemer := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7")
	
	params2 := nft.RedemptionParams{
		TokenID:  10002, // Different ticket
		Redeemer: differentRedeemer,
		Nonce:    uint64(time.Now().Unix()),
		Deadline: time.Now().Add(5 * time.Minute).Unix(),
	}

	tx2, err := client.RedeemTicket(ctx, params2)
	if err != nil {
		fmt.Printf("Failed to redeem for different user: %v\n", err)
		// This might fail if the signer doesn't have permission
	} else {
		fmt.Printf("Redemption for different user sent: %s\n", tx2.Hash().Hex())
	}

	fmt.Println("\nRedemption examples completed!")
}