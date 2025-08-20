package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	// Create temporary database for testing
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDB(dbPath)
	require.NoError(t, err)
	defer db.Close()

	t.Run("SaveAndRetrieveMapping", func(t *testing.T) {
		mapping := &NFTMapping{
			TokenID:       10001,
			DatakyteNFTID: "dk_nft_test123",
			Network:       "testnet",
			ContractAddr:  "0x1234567890abcdef",
			OwnerAddress:  "0xabcdef1234567890",
			BookingID:     "BOOK-001",
			EventID:       "EVENT-001",
			Status:        "active",
			MetadataURI:   "https://metadata.example.com/10001",
			ImageURL:      "https://image.example.com/10001.png",
			TxHash:        "0xdeadbeef",
		}

		// Save mapping
		err := db.SaveNFTMapping(mapping)
		assert.NoError(t, err)

		// Retrieve by token ID
		datakyteID, err := db.GetDatakyteID(10001, "testnet")
		assert.NoError(t, err)
		assert.Equal(t, "dk_nft_test123", datakyteID)

		// Get full mapping
		retrieved, err := db.GetNFTMapping(10001, "testnet")
		assert.NoError(t, err)
		assert.Equal(t, mapping.TokenID, retrieved.TokenID)
		assert.Equal(t, mapping.DatakyteNFTID, retrieved.DatakyteNFTID)
		assert.Equal(t, mapping.BookingID, retrieved.BookingID)
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		// Create a mapping
		mapping := &NFTMapping{
			TokenID:       10002,
			DatakyteNFTID: "dk_nft_test456",
			Network:       "testnet",
			ContractAddr:  "0x1234567890abcdef",
			OwnerAddress:  "0xabcdef1234567890",
			BookingID:     "BOOK-002",
			EventID:       "EVENT-002",
			Status:        "active",
			TxHash:        "0xbeefbeef",
		}

		err := db.SaveNFTMapping(mapping)
		require.NoError(t, err)

		// Update status
		err = db.UpdateNFTStatus(10002, "testnet", "expired")
		assert.NoError(t, err)

		// Verify update
		retrieved, err := db.GetNFTMapping(10002, "testnet")
		assert.NoError(t, err)
		assert.Equal(t, "expired", retrieved.Status)
	})

	t.Run("UpdateRedemption", func(t *testing.T) {
		// Create a mapping
		mapping := &NFTMapping{
			TokenID:       10003,
			DatakyteNFTID: "dk_nft_test789",
			Network:       "testnet",
			ContractAddr:  "0x1234567890abcdef",
			OwnerAddress:  "0xabcdef1234567890",
			Status:        "active",
			TxHash:        "0xcafecafe",
		}

		err := db.SaveNFTMapping(mapping)
		require.NoError(t, err)

		// Mark as redeemed
		err = db.UpdateNFTRedemption(10003, "testnet")
		assert.NoError(t, err)

		// Verify update
		retrieved, err := db.GetNFTMapping(10003, "testnet")
		assert.NoError(t, err)
		assert.Equal(t, "redeemed", retrieved.Status)
		assert.NotNil(t, retrieved.RedeemedAt)
	})

	t.Run("GetUserNFTs", func(t *testing.T) {
		ownerAddr := "0xuser123"

		// Create multiple mappings for the same user
		for i := uint64(20001); i <= 20003; i++ {
			mapping := &NFTMapping{
				TokenID:       i,
				DatakyteNFTID: fmt.Sprintf("dk_nft_%d", i),
				Network:       "testnet",
				ContractAddr:  "0x1234567890abcdef",
				OwnerAddress:  ownerAddr,
				Status:        "active",
				TxHash:        fmt.Sprintf("0xhash%d", i),
			}
			err := db.SaveNFTMapping(mapping)
			require.NoError(t, err)
		}

		// Get user's NFTs
		nfts, err := db.GetUserNFTs(ownerAddr, "testnet")
		assert.NoError(t, err)
		assert.Len(t, nfts, 3)
		// Verify all tokens are present
		tokenIDs := []uint64{nfts[0].TokenID, nfts[1].TokenID, nfts[2].TokenID}
		assert.Contains(t, tokenIDs, uint64(20001))
		assert.Contains(t, tokenIDs, uint64(20002))
		assert.Contains(t, tokenIDs, uint64(20003))
	})

	t.Run("NonExistentMapping", func(t *testing.T) {
		// Try to get non-existent mapping
		_, err := db.GetDatakyteID(99999, "testnet")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no mapping found")
	})
}

func TestSingleton(t *testing.T) {
	// Set a custom path for testing
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Get singleton instance
	db1 := GetDB()
	db2 := GetDB()

	// Should be the same instance
	assert.Same(t, db1, db2)

	// Cleanup
	db1.Close()
}
