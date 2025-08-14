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
	// Check required environment variables
	privateKey := os.Getenv("TESTNET_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("TESTNET_PRIVATE_KEY environment variable not set")
	}

	recipientAddr := os.Getenv("RECIPIENT_ADDRESS")
	if recipientAddr == "" {
		recipientAddr = "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7" // Default test address
	}

	// Create SDK client
	client, err := nft.NewClient(nft.ClientConfig{
		PrivateKey:      privateKey,
		Network:         "testnet",
		DatakyteEnabled: true,
	})
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}
	defer client.Close()

	// Check balance
	ctx := context.Background()
	balance, err := client.GetBalance(ctx)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("Account balance: %s CAM\n", balance.String())

	// Prepare mint parameters
	params := nft.MintParams{
		To:                common.HexToAddress(recipientAddr),
		BookingID:         nft.HashBookingID("BOOK-2024-001"),
		EventID:           nft.HashEventID("EVENT-WILDLIFE-SAFARI"),
		UtilityFlags:      0, // Standard ticket
		TransferUnlockAt:  uint64(time.Now().Add(24 * time.Hour).Unix()),
		ExpiresAt:         uint64(time.Now().Add(30 * 24 * time.Hour).Unix()),
		MetadataURI:       "", // Will use Datakyte
		RewardBasisPoints: 500, // 5% rewards
		DatakyteNftID:     "",  // Will be generated
	}

	fmt.Println("\nMinting ticket...")
	fmt.Printf("  To: %s\n", params.To.Hex())
	fmt.Printf("  Transfer unlocks: %s\n", time.Unix(int64(params.TransferUnlockAt), 0).Format(time.RFC3339))
	fmt.Printf("  Expires: %s\n", time.Unix(int64(params.ExpiresAt), 0).Format(time.RFC3339))
	fmt.Printf("  BOGO Rewards: %d basis points\n", params.RewardBasisPoints)

	// Mint the ticket
	tx, tokenID, err := client.MintTicket(ctx, params)
	if err != nil {
		log.Fatalf("Failed to mint ticket: %v", err)
	}

	fmt.Printf("\n✅ Ticket minted successfully!\n")
	fmt.Printf("  Token ID: %d\n", tokenID)
	fmt.Printf("  Transaction: %s\n", tx.Hash().Hex())
	fmt.Printf("  Metadata URL: https://dklnk.to/api/nfts/{contract}/%d/metadata\n", tokenID)
}