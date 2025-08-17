package nft

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// EIP712Domain represents the domain separator for EIP-712
type EIP712Domain struct {
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
}

// RedeemTicketData represents the typed data for ticket redemption
type RedeemTicketData struct {
	TokenID  *big.Int
	Redeemer common.Address
	Nonce    *big.Int
	Deadline *big.Int
}

// GetEIP712Domain returns the domain separator for the network
func GetEIP712Domain(chainID *big.Int, contractAddress common.Address) EIP712Domain {
	return EIP712Domain{
		Name:              "BOGOWITickets",
		Version:           "1",
		ChainId:           chainID,
		VerifyingContract: contractAddress,
	}
}

// GenerateRedemptionSignature generates an EIP-712 signature for ticket redemption
func GenerateRedemptionSignature(
	privateKey *ecdsa.PrivateKey,
	tokenID *big.Int,
	redeemer common.Address,
	nonce *big.Int,
	deadline *big.Int,
	chainID *big.Int,
	contractAddress common.Address,
) ([]byte, error) {
	// Create the EIP-712 typed data
	domain := apitypes.TypedDataDomain{
		Name:              "BOGOWITickets",
		Version:           "1",
		ChainId:           (*math.HexOrDecimal256)(chainID),
		VerifyingContract: contractAddress.Hex(),
	}

	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"RedeemTicket": {
			{Name: "tokenId", Type: "uint256"},
			{Name: "redeemer", Type: "address"},
			{Name: "nonce", Type: "uint256"},
			{Name: "deadline", Type: "uint256"},
		},
	}

	message := apitypes.TypedDataMessage{
		"tokenId":  (*math.HexOrDecimal256)(tokenID),
		"redeemer": redeemer.Hex(),
		"nonce":    (*math.HexOrDecimal256)(nonce),
		"deadline": (*math.HexOrDecimal256)(deadline),
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: "RedeemTicket",
		Domain:      domain,
		Message:     message,
	}

	// Hash the typed data
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to hash domain: %w", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to hash message: %w", err)
	}

	// Compute the digest
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	digest := crypto.Keccak256(rawData)

	// Sign the digest
	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	// Adjust v value for Ethereum (27 or 28)
	if signature[64] < 27 {
		signature[64] += 27
	}

	return signature, nil
}

// VerifyRedemptionSignature verifies an EIP-712 redemption signature
func VerifyRedemptionSignature(
	signature []byte,
	tokenID *big.Int,
	redeemer common.Address,
	nonce *big.Int,
	deadline *big.Int,
	chainID *big.Int,
	contractAddress common.Address,
	expectedSigner common.Address,
) (bool, error) {
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length: expected 65, got %d", len(signature))
	}

	// Create the same typed data structure
	domain := apitypes.TypedDataDomain{
		Name:              "BOGOWITickets",
		Version:           "1",
		ChainId:           (*math.HexOrDecimal256)(chainID),
		VerifyingContract: contractAddress.Hex(),
	}

	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"RedeemTicket": {
			{Name: "tokenId", Type: "uint256"},
			{Name: "redeemer", Type: "address"},
			{Name: "nonce", Type: "uint256"},
			{Name: "deadline", Type: "uint256"},
		},
	}

	message := apitypes.TypedDataMessage{
		"tokenId":  (*math.HexOrDecimal256)(tokenID),
		"redeemer": redeemer.Hex(),
		"nonce":    (*math.HexOrDecimal256)(nonce),
		"deadline": (*math.HexOrDecimal256)(deadline),
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: "RedeemTicket",
		Domain:      domain,
		Message:     message,
	}

	// Hash the typed data
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return false, fmt.Errorf("failed to hash domain: %w", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return false, fmt.Errorf("failed to hash message: %w", err)
	}

	// Compute the digest
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	digest := crypto.Keccak256(rawData)

	// Recover the signer
	sig := make([]byte, 65)
	copy(sig, signature)
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	pubKey, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr == expectedSigner, nil
}

// GenerateRedemptionQRCode generates QR code data for ticket redemption
func GenerateRedemptionQRCode(
	tokenID *big.Int,
	redeemer common.Address,
	nonce *big.Int,
	deadline *big.Int,
	signature []byte,
	baseURL string,
) string {
	// Format: baseURL/redeem?tokenId=X&redeemer=0x...&nonce=Y&deadline=Z&sig=0x...
	return fmt.Sprintf(
		"%s/redeem?tokenId=%s&redeemer=%s&nonce=%s&deadline=%s&sig=%x",
		baseURL,
		tokenID.String(),
		redeemer.Hex(),
		nonce.String(),
		deadline.String(),
		signature,
	)
}

// ParseRedemptionQRCode parses QR code data for ticket redemption
func ParseRedemptionQRCode(qrData string) (*RedemptionData, error) {
	// This is a simplified parser - in production, use proper URL parsing
	// Expected format: baseURL/redeem?tokenId=X&redeemer=0x...&nonce=Y&deadline=Z&sig=0x...
	
	// TODO: Implement proper URL parsing
	return nil, fmt.Errorf("QR code parsing not yet implemented")
}