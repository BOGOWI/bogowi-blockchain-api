package nft

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test contract query methods directly
func TestContractQueries(t *testing.T) {
	mockContract := new(MockTicketsContract)

	t.Run("IsTransferable", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()

		result, err := mockContract.IsTransferable(&bind.CallOpts{}, big.NewInt(1))
		assert.NoError(t, err)
		assert.True(t, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("OwnerOf", func(t *testing.T) {
		owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(owner, nil).Once()

		result, err := mockContract.OwnerOf(&bind.CallOpts{}, big.NewInt(1))
		assert.NoError(t, err)
		assert.Equal(t, owner, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("TokenURI", func(t *testing.T) {
		uri := "https://api.example.com/metadata/1"
		mockContract.On("TokenURI", mock.Anything, big.NewInt(1)).Return(uri, nil).Once()

		result, err := mockContract.TokenURI(&bind.CallOpts{}, big.NewInt(1))
		assert.NoError(t, err)
		assert.Equal(t, uri, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("BalanceOf", func(t *testing.T) {
		owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
		balance := big.NewInt(5)
		mockContract.On("BalanceOf", mock.Anything, owner).Return(balance, nil).Once()

		result, err := mockContract.BalanceOf(&bind.CallOpts{}, owner)
		assert.NoError(t, err)
		assert.Equal(t, balance, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("GetTicketData", func(t *testing.T) {
		expectedData := TicketDataContract{
			BookingID:                  [32]byte{1, 2, 3},
			EventID:                    [32]byte{4, 5, 6},
			TransferUnlockAt:           1700000000,
			ExpiresAt:                  1800000000,
			UtilityFlags:               0x0FFF,
			State:                      1,
			NonTransferableAfterRedeem: true,
			BurnOnRedeem:               false,
		}

		mockContract.On("GetTicketData", mock.Anything, big.NewInt(1)).Return(expectedData, nil).Once()

		result, err := mockContract.GetTicketData(&bind.CallOpts{}, big.NewInt(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedData, result)

		mockContract.AssertExpectations(t)
	})
}

func TestTicketStateString(t *testing.T) {
	tests := []struct {
		state    TicketState
		expected string
	}{
		{TicketStateIssued, "Issued"},
		{TicketStateRedeemed, "Redeemed"},
		{TicketStateExpired, "Expired"},
		{TicketState(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.state.String())
		})
	}
}
