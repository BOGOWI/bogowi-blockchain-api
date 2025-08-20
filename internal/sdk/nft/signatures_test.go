package nft

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEIP712Domain(t *testing.T) {
	chainID := big.NewInt(1)
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	domain := GetEIP712Domain(chainID, contractAddress)

	assert.Equal(t, "BOGOWITickets", domain.Name)
	assert.Equal(t, "1", domain.Version)
	assert.Equal(t, chainID, domain.ChainId)
	assert.Equal(t, contractAddress, domain.VerifyingContract)
}

func TestGenerateAndVerifyRedemptionSignature(t *testing.T) {
	// Generate a test private key
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Test parameters
	tokenID := big.NewInt(123)
	redeemer := common.HexToAddress("0x9876543210987654321098765432109876543210")
	nonce := big.NewInt(1)
	deadline := big.NewInt(1700000000)
	chainID := big.NewInt(31337) // Local test chain
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("generate valid signature", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
		)

		assert.NoError(t, err)
		assert.NotNil(t, signature)
		assert.Len(t, signature, 65)
		// Check that v value is adjusted (should be 27 or 28)
		assert.True(t, signature[64] == 27 || signature[64] == 28)
	})

	t.Run("verify valid signature", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
		)
		require.NoError(t, err)

		valid, err := VerifyRedemptionSignature(
			signature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signerAddress,
		)

		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("verify invalid signature - wrong signer", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
		)
		require.NoError(t, err)

		wrongSigner := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		valid, err := VerifyRedemptionSignature(
			signature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			wrongSigner,
		)

		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("verify invalid signature - wrong token ID", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
		)
		require.NoError(t, err)

		wrongTokenID := big.NewInt(999)
		valid, err := VerifyRedemptionSignature(
			signature,
			wrongTokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signerAddress,
		)

		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("verify invalid signature - wrong chain ID", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
		)
		require.NoError(t, err)

		wrongChainID := big.NewInt(1)
		valid, err := VerifyRedemptionSignature(
			signature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			wrongChainID,
			contractAddress,
			signerAddress,
		)

		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("verify invalid signature length", func(t *testing.T) {
		invalidSignature := []byte("too short")

		valid, err := VerifyRedemptionSignature(
			invalidSignature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signerAddress,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid signature length")
		assert.False(t, valid)
	})

	t.Run("verify corrupted signature", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
		)
		require.NoError(t, err)

		// Corrupt the signature
		signature[10] = ^signature[10]

		valid, err := VerifyRedemptionSignature(
			signature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signerAddress,
		)

		// May return error or just false depending on corruption
		if err != nil {
			assert.Contains(t, err.Error(), "failed to recover public key")
		}
		assert.False(t, valid)
	})
}

func TestGenerateRedemptionQRCode(t *testing.T) {
	tokenID := big.NewInt(123)
	redeemer := common.HexToAddress("0x9876543210987654321098765432109876543210")
	nonce := big.NewInt(1)
	deadline := big.NewInt(1700000000)
	signature := make([]byte, 65)
	baseURL := "https://example.com"

	t.Run("generate QR code URL", func(t *testing.T) {
		qrCode := GenerateRedemptionQRCode(
			tokenID,
			redeemer,
			nonce,
			deadline,
			signature,
			baseURL,
		)

		expected := fmt.Sprintf(
			"%s/redeem?tokenId=%s&redeemer=%s&nonce=%s&deadline=%s&sig=%x",
			baseURL,
			tokenID.String(),
			redeemer.Hex(),
			nonce.String(),
			deadline.String(),
			signature,
		)

		assert.Equal(t, expected, qrCode)
	})

	t.Run("generate QR code with different values", func(t *testing.T) {
		tokenID := big.NewInt(999999)
		redeemer := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		nonce := big.NewInt(42)
		deadline := big.NewInt(2000000000)
		signature := []byte{0x01, 0x02, 0x03}
		baseURL := "https://bogowi.com"

		qrCode := GenerateRedemptionQRCode(
			tokenID,
			redeemer,
			nonce,
			deadline,
			signature,
			baseURL,
		)

		assert.Contains(t, qrCode, baseURL)
		assert.Contains(t, qrCode, "tokenId=999999")
		assert.Contains(t, qrCode, redeemer.Hex())
		assert.Contains(t, qrCode, "nonce=42")
		assert.Contains(t, qrCode, "deadline=2000000000")
		assert.Contains(t, qrCode, "sig=010203")
	})

	t.Run("generate QR code with empty base URL", func(t *testing.T) {
		qrCode := GenerateRedemptionQRCode(
			tokenID,
			redeemer,
			nonce,
			deadline,
			signature,
			"",
		)

		assert.Equal(t, "/redeem?tokenId=123&redeemer=0x9876543210987654321098765432109876543210&nonce=1&deadline=1700000000&sig=0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", qrCode)
	})
}

func TestParseRedemptionQRCode(t *testing.T) {
	t.Run("not implemented", func(t *testing.T) {
		qrData := "https://example.com/redeem?tokenId=123&redeemer=0x123&nonce=1&deadline=1700000000&sig=0x123"

		result, err := ParseRedemptionQRCode(qrData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "QR code parsing not yet implemented")
		assert.Nil(t, result)
	})
}

func TestEIP712SignatureEdgeCases(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	t.Run("nil private key", func(t *testing.T) {
		_, err := GenerateRedemptionSignature(
			nil,
			big.NewInt(1),
			common.HexToAddress("0x123"),
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
			common.HexToAddress("0x456"),
		)

		assert.Error(t, err)
	})

	t.Run("zero token ID", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			big.NewInt(0),
			common.HexToAddress("0x123"),
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
			common.HexToAddress("0x456"),
		)

		assert.NoError(t, err)
		assert.NotNil(t, signature)
	})

	t.Run("very large values", func(t *testing.T) {
		largeValue := new(big.Int)
		largeValue.SetString("999999999999999999999999999999999999999999", 10)

		signature, err := GenerateRedemptionSignature(
			privateKey,
			largeValue,
			common.HexToAddress("0x123"),
			largeValue,
			largeValue,
			big.NewInt(1),
			common.HexToAddress("0x456"),
		)

		assert.NoError(t, err)
		assert.NotNil(t, signature)
	})

	t.Run("zero address", func(t *testing.T) {
		signature, err := GenerateRedemptionSignature(
			privateKey,
			big.NewInt(1),
			common.Address{},
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
			common.Address{},
		)

		assert.NoError(t, err)
		assert.NotNil(t, signature)
	})
}

func TestSignatureVValueHandling(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Generate multiple signatures to test v value handling
	for i := 0; i < 10; i++ {
		nonce := big.NewInt(int64(i))

		signature, err := GenerateRedemptionSignature(
			privateKey,
			big.NewInt(1),
			common.HexToAddress("0x123"),
			nonce,
			big.NewInt(1700000000),
			big.NewInt(1),
			common.HexToAddress("0x456"),
		)

		require.NoError(t, err)

		// V value should always be 27 or 28 after adjustment
		assert.True(t, signature[64] == 27 || signature[64] == 28,
			"V value should be 27 or 28, got %d", signature[64])

		// Verify the signature works
		valid, err := VerifyRedemptionSignature(
			signature,
			big.NewInt(1),
			common.HexToAddress("0x123"),
			nonce,
			big.NewInt(1700000000),
			big.NewInt(1),
			common.HexToAddress("0x456"),
			signerAddress,
		)

		assert.NoError(t, err)
		assert.True(t, valid)
	}
}

func BenchmarkGenerateRedemptionSignature(b *testing.B) {
	privateKey, _ := crypto.GenerateKey()
	tokenID := big.NewInt(123)
	redeemer := common.HexToAddress("0x9876543210987654321098765432109876543210")
	nonce := big.NewInt(1)
	deadline := big.NewInt(1700000000)
	chainID := big.NewInt(1)
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateRedemptionSignature(
			privateKey,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
		)
	}
}

func BenchmarkVerifyRedemptionSignature(b *testing.B) {
	privateKey, _ := crypto.GenerateKey()
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	tokenID := big.NewInt(123)
	redeemer := common.HexToAddress("0x9876543210987654321098765432109876543210")
	nonce := big.NewInt(1)
	deadline := big.NewInt(1700000000)
	chainID := big.NewInt(1)
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	signature, _ := GenerateRedemptionSignature(
		privateKey,
		tokenID,
		redeemer,
		nonce,
		deadline,
		chainID,
		contractAddress,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = VerifyRedemptionSignature(
			signature,
			tokenID,
			redeemer,
			nonce,
			deadline,
			chainID,
			contractAddress,
			signerAddress,
		)
	}
}
