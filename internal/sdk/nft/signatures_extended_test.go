package nft

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRedemptionSignatureErrorCases(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	chainID := big.NewInt(31337)
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	redeemer := common.HexToAddress("0x9876543210987654321098765432109876543210")

	t.Run("nil token ID", func(t *testing.T) {
		_, err := GenerateRedemptionSignature(
			privateKey,
			nil,
			redeemer,
			big.NewInt(1),
			big.NewInt(1700000000),
			chainID,
			contractAddress,
		)
		// Should error - nil values not handled
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid integer value")
	})

	t.Run("nil nonce", func(t *testing.T) {
		_, err := GenerateRedemptionSignature(
			privateKey,
			big.NewInt(1),
			redeemer,
			nil,
			big.NewInt(1700000000),
			chainID,
			contractAddress,
		)
		// Should error - nil values not handled
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid integer value")
	})

	t.Run("nil deadline", func(t *testing.T) {
		_, err := GenerateRedemptionSignature(
			privateKey,
			big.NewInt(1),
			redeemer,
			big.NewInt(1),
			nil,
			chainID,
			contractAddress,
		)
		// Should error - nil values not handled
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid integer value")
	})

	t.Run("nil chain ID", func(t *testing.T) {
		_, err := GenerateRedemptionSignature(
			privateKey,
			big.NewInt(1),
			redeemer,
			big.NewInt(1),
			big.NewInt(1700000000),
			nil,
			contractAddress,
		)
		// Should error - nil values not handled
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid integer value")
	})

	t.Run("negative values", func(t *testing.T) {
		_, err := GenerateRedemptionSignature(
			privateKey,
			big.NewInt(-1),
			redeemer,
			big.NewInt(-999),
			big.NewInt(-1700000000),
			big.NewInt(-1),
			contractAddress,
		)
		// Should error - negative values not allowed for unsigned types
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid negative value")
	})
}

func TestVerifyRedemptionSignatureEdgeCases(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	signer := crypto.PubkeyToAddress(privateKey.PublicKey)

	tokenID := big.NewInt(123)
	redeemer := common.HexToAddress("0x9876543210987654321098765432109876543210")
	nonce := big.NewInt(1)
	deadline := big.NewInt(1700000000)
	chainID := big.NewInt(31337)
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	validSignature, err := GenerateRedemptionSignature(
		privateKey,
		tokenID,
		redeemer,
		nonce,
		deadline,
		chainID,
		contractAddress,
	)
	require.NoError(t, err)

	t.Run("empty signature", func(t *testing.T) {
		valid, err := VerifyRedemptionSignature(
			[]byte{},
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signer,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid signature length")
		assert.False(t, valid)
	})

	t.Run("signature too short", func(t *testing.T) {
		valid, err := VerifyRedemptionSignature(
			make([]byte, 64),
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signer,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid signature length")
		assert.False(t, valid)
	})

	t.Run("signature too long", func(t *testing.T) {
		valid, err := VerifyRedemptionSignature(
			make([]byte, 66),
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signer,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid signature length")
		assert.False(t, valid)
	})

	t.Run("invalid signature format", func(t *testing.T) {
		// Create a signature with invalid recovery ID
		invalidSig := make([]byte, 65)
		copy(invalidSig, validSignature)
		invalidSig[64] = 255 // Invalid recovery ID

		valid, err := VerifyRedemptionSignature(
			invalidSig,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signer,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to recover public key")
		assert.False(t, valid)
	})

	t.Run("nil token ID in verification", func(t *testing.T) {
		valid, err := VerifyRedemptionSignature(
			validSignature,
			nil, // nil token ID
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signer,
		)
		// Should error - nil values not handled
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid integer value")
		assert.False(t, valid)
	})

	t.Run("zero address as signer", func(t *testing.T) {
		valid, err := VerifyRedemptionSignature(
			validSignature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			common.Address{}, // Zero address
		)
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("different contract address", func(t *testing.T) {
		valid, err := VerifyRedemptionSignature(
			validSignature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			signer,
		)
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("v value already adjusted", func(t *testing.T) {
		// Test with v value already < 27
		sigWithLowV := make([]byte, 65)
		copy(sigWithLowV, validSignature)
		if sigWithLowV[64] >= 27 {
			sigWithLowV[64] -= 27
		}

		valid, err := VerifyRedemptionSignature(
			sigWithLowV,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signer,
		)
		assert.NoError(t, err)
		assert.True(t, valid)
	})
}

func TestGenerateRedemptionQRCodeEdgeCases(t *testing.T) {
	t.Run("nil values", func(t *testing.T) {
		qrCode := GenerateRedemptionQRCode(
			nil,
			common.Address{},
			nil,
			nil,
			nil,
			"https://example.com",
		)
		assert.Contains(t, qrCode, "https://example.com/redeem?")
		assert.Contains(t, qrCode, "tokenId=<nil>")
		assert.Contains(t, qrCode, "nonce=<nil>")
		assert.Contains(t, qrCode, "deadline=<nil>")
	})

	t.Run("very large numbers", func(t *testing.T) {
		largeNumber := new(big.Int)
		largeNumber.SetString("999999999999999999999999999999999999999999", 10)

		qrCode := GenerateRedemptionQRCode(
			largeNumber,
			common.HexToAddress("0x123"),
			largeNumber,
			largeNumber,
			make([]byte, 65),
			"https://example.com",
		)
		assert.Contains(t, qrCode, "999999999999999999999999999999999999999999")
	})

	t.Run("special characters in base URL", func(t *testing.T) {
		qrCode := GenerateRedemptionQRCode(
			big.NewInt(1),
			common.HexToAddress("0x123"),
			big.NewInt(1),
			big.NewInt(1),
			make([]byte, 65),
			"https://example.com/path?existing=param",
		)
		assert.Contains(t, qrCode, "https://example.com/path?existing=param/redeem?")
	})

	t.Run("empty signature", func(t *testing.T) {
		qrCode := GenerateRedemptionQRCode(
			big.NewInt(1),
			common.HexToAddress("0x123"),
			big.NewInt(1),
			big.NewInt(1),
			[]byte{},
			"https://example.com",
		)
		assert.Contains(t, qrCode, "sig=")
	})

	t.Run("single byte signature", func(t *testing.T) {
		qrCode := GenerateRedemptionQRCode(
			big.NewInt(1),
			common.HexToAddress("0x123"),
			big.NewInt(1),
			big.NewInt(1),
			[]byte{0xFF},
			"https://example.com",
		)
		assert.Contains(t, qrCode, "sig=ff")
	})
}

func TestParseRedemptionQRCodeImplementation(t *testing.T) {
	t.Run("various invalid formats", func(t *testing.T) {
		testCases := []string{
			"",
			"not-a-url",
			"https://example.com",
			"https://example.com/wrong-path",
			"https://example.com/redeem",
			"https://example.com/redeem?incomplete=params",
		}

		for _, tc := range testCases {
			result, err := ParseRedemptionQRCode(tc)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "QR code parsing not yet implemented")
			assert.Nil(t, result)
		}
	})
}

func TestEIP712DomainConsistency(t *testing.T) {
	testCases := []struct {
		name            string
		chainID         *big.Int
		contractAddress common.Address
	}{
		{
			name:            "mainnet",
			chainID:         big.NewInt(1),
			contractAddress: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
		{
			name:            "testnet",
			chainID:         big.NewInt(5),
			contractAddress: common.HexToAddress("0xABCDEF1234567890123456789012345678901234"),
		},
		{
			name:            "local",
			chainID:         big.NewInt(31337),
			contractAddress: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			domain := GetEIP712Domain(tc.chainID, tc.contractAddress)

			assert.Equal(t, "BOGOWITickets", domain.Name)
			assert.Equal(t, "1", domain.Version)
			assert.Equal(t, tc.chainID, domain.ChainId)
			assert.Equal(t, tc.contractAddress, domain.VerifyingContract)
		})
	}
}

func TestSignatureRoundTrip(t *testing.T) {
	// Test that a signature generated can be verified
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	signer := crypto.PubkeyToAddress(privateKey.PublicKey)

	testCases := []struct {
		name            string
		tokenID         *big.Int
		redeemer        common.Address
		nonce           *big.Int
		deadline        *big.Int
		chainID         *big.Int
		contractAddress common.Address
	}{
		{
			name:            "standard values",
			tokenID:         big.NewInt(1),
			redeemer:        common.HexToAddress("0x123"),
			nonce:           big.NewInt(0),
			deadline:        big.NewInt(1700000000),
			chainID:         big.NewInt(1),
			contractAddress: common.HexToAddress("0x456"),
		},
		{
			name:            "large values",
			tokenID:         big.NewInt(1000000),
			redeemer:        common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
			nonce:           big.NewInt(999999),
			deadline:        big.NewInt(9999999999),
			chainID:         big.NewInt(137),
			contractAddress: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		},
		{
			name:            "zero values",
			tokenID:         big.NewInt(0),
			redeemer:        common.Address{},
			nonce:           big.NewInt(0),
			deadline:        big.NewInt(0),
			chainID:         big.NewInt(0),
			contractAddress: common.Address{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate signature
			signature, err := GenerateRedemptionSignature(
				privateKey,
				tc.tokenID,
				tc.redeemer,
				tc.nonce,
				tc.deadline,
				tc.chainID,
				tc.contractAddress,
			)
			require.NoError(t, err)
			require.Len(t, signature, 65)

			// Verify signature
			valid, err := VerifyRedemptionSignature(
				signature,
				tc.tokenID,
				tc.redeemer,
				tc.nonce,
				tc.deadline,
				tc.chainID,
				tc.contractAddress,
				signer,
			)
			assert.NoError(t, err)
			assert.True(t, valid)

			// Verify with wrong signer fails
			wrongSigner := common.HexToAddress("0x9999999999999999999999999999999999999999")
			valid, err = VerifyRedemptionSignature(
				signature,
				tc.tokenID,
				tc.redeemer,
				tc.nonce,
				tc.deadline,
				tc.chainID,
				tc.contractAddress,
				wrongSigner,
			)
			assert.NoError(t, err)
			assert.False(t, valid)
		})
	}
}
