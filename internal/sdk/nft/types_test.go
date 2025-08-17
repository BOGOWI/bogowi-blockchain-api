package nft

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestHashBookingID(t *testing.T) {
	tests := []struct {
		name      string
		bookingID string
		wantEmpty bool
	}{
		{
			name:      "valid booking ID",
			bookingID: "BOOK-12345",
			wantEmpty: false,
		},
		{
			name:      "empty booking ID",
			bookingID: "",
			wantEmpty: false, // Still returns a hash, just of empty string
		},
		{
			name:      "long booking ID",
			bookingID: "BOOK-1234567890-ABCDEFGHIJKLMNOP",
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashBookingID(tt.bookingID)

			// Check that we got a 32-byte array
			if len(got) != 32 {
				t.Errorf("HashBookingID() returned array of length %d, want 32", len(got))
			}

			// Check if result is empty (all zeros)
			isEmpty := true
			for _, b := range got {
				if b != 0 {
					isEmpty = false
					break
				}
			}

			if tt.bookingID != "" && isEmpty {
				t.Errorf("HashBookingID() returned empty hash for non-empty input")
			}
		})
	}
}

func TestHashEventID(t *testing.T) {
	tests := []struct {
		name    string
		eventID string
	}{
		{
			name:    "valid event ID",
			eventID: "EVENT-2024-001",
		},
		{
			name:    "empty event ID",
			eventID: "",
		},
		{
			name:    "numeric event ID",
			eventID: "123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashEventID(tt.eventID)

			// Check that we got a 32-byte array
			if len(got) != 32 {
				t.Errorf("HashEventID() returned array of length %d, want 32", len(got))
			}
		})
	}
}

func TestParseTicketState(t *testing.T) {
	tests := []struct {
		name  string
		state uint8
		want  TicketState
	}{
		{
			name:  "issued state",
			state: 0,
			want:  TicketStateIssued,
		},
		{
			name:  "redeemed state",
			state: 1,
			want:  TicketStateRedeemed,
		},
		{
			name:  "expired state",
			state: 2,
			want:  TicketStateExpired,
		},
		{
			name:  "unknown state",
			state: 99,
			want:  TicketState(99),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTicketState(tt.state); got != tt.want {
				t.Errorf("ParseTicketState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketState_String(t *testing.T) {
	tests := []struct {
		name  string
		state TicketState
		want  string
	}{
		{
			name:  "issued state",
			state: TicketStateIssued,
			want:  "Issued",
		},
		{
			name:  "redeemed state",
			state: TicketStateRedeemed,
			want:  "Redeemed",
		},
		{
			name:  "expired state",
			state: TicketStateExpired,
			want:  "Expired",
		},
		{
			name:  "unknown state",
			state: TicketState(99),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("TicketState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkConfig(t *testing.T) {
	tests := []struct {
		name      string
		network   string
		wantChain *big.Int
		wantErr   bool
	}{
		{
			name:      "testnet config",
			network:   "testnet",
			wantChain: big.NewInt(501),
			wantErr:   false,
		},
		{
			name:      "mainnet config",
			network:   "mainnet",
			wantChain: big.NewInt(500),
			wantErr:   false,
		},
		{
			name:      "invalid network",
			network:   "invalid",
			wantChain: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNetworkConfig(tt.network)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("GetNetworkConfig() returned nil config")
					return
				}

				if got.ChainID.Cmp(tt.wantChain) != 0 {
					t.Errorf("GetNetworkConfig() ChainID = %v, want %v", got.ChainID, tt.wantChain)
				}

				if got.RPCURL == "" {
					t.Errorf("GetNetworkConfig() returned empty RPC URL")
				}
			}
		})
	}
}

func TestSDKError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  SDKError
		want string
	}{
		{
			name: "invalid network error",
			err:  ErrInvalidNetwork,
			want: "Invalid network specified",
		},
		{
			name: "ticket expired error",
			err:  ErrTicketExpired,
			want: "Ticket has expired",
		},
		{
			name: "custom error",
			err: SDKError{
				Code:    "CUSTOM_ERROR",
				Message: "Custom error message",
			},
			want: "Custom error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("SDKError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedemptionParams(t *testing.T) {
	// Test that RedemptionParams struct is properly defined
	params := RedemptionParams{
		TokenID:  10001,
		Redeemer: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7"),
		Nonce:    12345,
		Deadline: 1735689600,
	}

	if params.TokenID != 10001 {
		t.Errorf("TokenID = %v, want %v", params.TokenID, 10001)
	}

	if params.Nonce != 12345 {
		t.Errorf("Nonce = %v, want %v", params.Nonce, 12345)
	}

	if params.Deadline != 1735689600 {
		t.Errorf("Deadline = %v, want %v", params.Deadline, 1735689600)
	}
}
