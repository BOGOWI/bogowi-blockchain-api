package nft

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// MintParams represents parameters for minting a new ticket
type MintParams struct {
	To                common.Address
	BookingID         [32]byte
	EventID           [32]byte
	UtilityFlags      uint32
	TransferUnlockAt  uint64
	ExpiresAt         uint64
	MetadataURI       string
	RewardBasisPoints uint16
	DatakyteNftID     string // Optional Datakyte linking
}

// TicketData represents on-chain ticket data
type TicketData struct {
	BookingID                  [32]byte
	EventID                    [32]byte
	TransferUnlockAt           uint64
	ExpiresAt                  uint64
	UtilityFlags               uint32
	State                      uint8 // 0=ISSUED, 1=REDEEMED, 2=EXPIRED
	NonTransferableAfterRedeem bool
	BurnOnRedeem               bool
}

// TicketState represents the lifecycle state of a ticket
type TicketState uint8

const (
	TicketStateIssued   TicketState = 0
	TicketStateRedeemed TicketState = 1
	TicketStateExpired  TicketState = 2
)

// RedemptionData represents data needed for ticket redemption
type RedemptionData struct {
	TokenID   *big.Int
	Redeemer  common.Address
	Nonce     *big.Int
	Deadline  *big.Int
	ChainID   *big.Int
	Signature []byte
}

// RedemptionParams represents parameters for redeeming a ticket
type RedemptionParams struct {
	TokenID  uint64
	Redeemer common.Address
	Nonce    uint64
	Deadline int64
}

// EventFilter represents parameters for filtering events
type EventFilter struct {
	FromBlock *big.Int
	ToBlock   *big.Int
	TokenIDs  []*big.Int
	Addresses []common.Address
}

// TransactionResult represents the result of a blockchain transaction
type TransactionResult struct {
	TxHash    common.Hash
	GasUsed   uint64
	BlockHash common.Hash
	BlockNum  uint64
	Status    uint64
}

// TokenMetadata represents NFT metadata
type TokenMetadata struct {
	TokenID            uint64
	Name               string
	Description        string
	Image              string
	ExternalURL        string
	Attributes         []MetadataAttribute
	DatakyteID         string
	ConservationImpact string
}

// MetadataAttribute represents an NFT metadata attribute
type MetadataAttribute struct {
	TraitType   string      `json:"trait_type"`
	Value       interface{} `json:"value"`
	DisplayType string      `json:"display_type,omitempty"`
}

// NetworkConfig represents network-specific configuration
type NetworkConfig struct {
	ChainID          *big.Int
	RPCURL           string
	TicketsContract  common.Address
	RoleManager      common.Address
	DatakyteAPIKey   string
	DatakyteBaseURL  string
	BlockTime        time.Duration
	ConfirmationWait int
}

// ClientConfig represents SDK client configuration
type ClientConfig struct {
	PrivateKey      string
	Network         string // "testnet" or "mainnet"
	CustomRPCURL    string // Optional custom RPC
	GasMultiplier   float64
	MaxGasPrice     *big.Int
	RequestTimeout  time.Duration
	RetryAttempts   int
	DatakyteEnabled bool
}

// Error types
type SDKError struct {
	Code    string
	Message string
	Details interface{}
}

func (e SDKError) Error() string {
	return e.Message
}

// Common errors
var (
	ErrInvalidNetwork      = SDKError{Code: "INVALID_NETWORK", Message: "Invalid network specified"}
	ErrContractNotDeployed = SDKError{Code: "CONTRACT_NOT_DEPLOYED", Message: "Contract not deployed on this network"}
	ErrInvalidSignature    = SDKError{Code: "INVALID_SIGNATURE", Message: "Invalid signature"}
	ErrTicketNotFound      = SDKError{Code: "TICKET_NOT_FOUND", Message: "Ticket not found"}
	ErrTicketExpired       = SDKError{Code: "TICKET_EXPIRED", Message: "Ticket has expired"}
	ErrTicketRedeemed      = SDKError{Code: "TICKET_REDEEMED", Message: "Ticket already redeemed"}
	ErrNotTransferable     = SDKError{Code: "NOT_TRANSFERABLE", Message: "Ticket is not transferable"}
	ErrInsufficientGas     = SDKError{Code: "INSUFFICIENT_GAS", Message: "Insufficient gas for transaction"}
	ErrDatakyteSyncFailed  = SDKError{Code: "DATAKYTE_SYNC_FAILED", Message: "Failed to sync with Datakyte"}
)

// Helper functions

// HashBookingID converts a booking ID string to bytes32
func HashBookingID(bookingID string) [32]byte {
	return common.BytesToHash([]byte(bookingID))
}

// HashEventID converts an event ID string to bytes32
func HashEventID(eventID string) [32]byte {
	return common.BytesToHash([]byte(eventID))
}

// ParseTicketState converts uint8 to TicketState
func ParseTicketState(state uint8) TicketState {
	return TicketState(state)
}

// GetNetworkConfig returns configuration for specified network
func GetNetworkConfig(network string) (*NetworkConfig, error) {
	switch network {
	case "testnet":
		return &NetworkConfig{
			ChainID:          big.NewInt(501),
			RPCURL:           "https://columbus.camino.network/ext/bc/C/rpc",
			DatakyteBaseURL:  "https://api.datakyte.io/testnet",
			BlockTime:        2 * time.Second,
			ConfirmationWait: 3,
		}, nil
	case "mainnet":
		return &NetworkConfig{
			ChainID:          big.NewInt(500),
			RPCURL:           "https://api.camino.network/ext/bc/C/rpc",
			DatakyteBaseURL:  "https://api.datakyte.io/mainnet",
			BlockTime:        2 * time.Second,
			ConfirmationWait: 6,
		}, nil
	default:
		return nil, ErrInvalidNetwork
	}
}
