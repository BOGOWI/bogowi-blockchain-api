package nft

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestClient_GetTicketData tests the actual Client.GetTicketData method
func TestClient_GetTicketData(t *testing.T) {
	mockContract := new(MockTicketsContract)
	client := &Client{
		ticketsContract: mockContract,
	}

	ctx := context.Background()
	tokenID := uint64(123)

	expectedContractData := TicketDataContract{
		BookingId:                  [32]byte{1, 2, 3},
		EventId:                    [32]byte{4, 5, 6},
		TransferUnlockAt:           1700000000,
		ExpiresAt:                  1800000000,
		UtilityFlags:               0x0FFF,
		State:                      1,
		NonTransferableAfterRedeem: true,
		BurnOnRedeem:               false,
	}

	t.Run("success", func(t *testing.T) {
		mockContract.On("GetTicketData", mock.Anything, big.NewInt(int64(tokenID))).
			Return(expectedContractData, nil).Once()

		result, err := client.GetTicketData(ctx, tokenID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedContractData.BookingId, result.BookingID)
		assert.Equal(t, expectedContractData.EventId, result.EventID)
		assert.Equal(t, expectedContractData.TransferUnlockAt, result.TransferUnlockAt)
		assert.Equal(t, expectedContractData.ExpiresAt, result.ExpiresAt)
		assert.Equal(t, expectedContractData.UtilityFlags, result.UtilityFlags)
		assert.Equal(t, expectedContractData.State, result.State)
		assert.Equal(t, expectedContractData.NonTransferableAfterRedeem, result.NonTransferableAfterRedeem)
		assert.Equal(t, expectedContractData.BurnOnRedeem, result.BurnOnRedeem)

		mockContract.AssertExpectations(t)
	})

	t.Run("contract error", func(t *testing.T) {
		mockContract.On("GetTicketData", mock.Anything, big.NewInt(int64(tokenID))).
			Return(TicketDataContract{}, assert.AnError).Once()

		result, err := client.GetTicketData(ctx, tokenID)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get ticket data")

		mockContract.AssertExpectations(t)
	})
}

// TestClient_IsTransferable tests the actual Client.IsTransferable method
func TestClient_IsTransferable(t *testing.T) {
	mockContract := new(MockTicketsContract)
	client := &Client{
		ticketsContract: mockContract,
	}

	ctx := context.Background()
	tokenID := uint64(456)

	t.Run("transferable", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(int64(tokenID))).
			Return(true, nil).Once()

		result, err := client.IsTransferable(ctx, tokenID)

		require.NoError(t, err)
		assert.True(t, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("not transferable", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(int64(tokenID))).
			Return(false, nil).Once()

		result, err := client.IsTransferable(ctx, tokenID)

		require.NoError(t, err)
		assert.False(t, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(int64(tokenID))).
			Return(false, assert.AnError).Once()

		result, err := client.IsTransferable(ctx, tokenID)

		require.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), "failed to check transferability")

		mockContract.AssertExpectations(t)
	})
}

// TestClient_GetOwnerOf tests the actual Client.GetOwnerOf method
func TestClient_GetOwnerOf(t *testing.T) {
	mockContract := new(MockTicketsContract)
	client := &Client{
		ticketsContract: mockContract,
	}

	ctx := context.Background()
	tokenID := uint64(789)
	expectedOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("success", func(t *testing.T) {
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(int64(tokenID))).
			Return(expectedOwner, nil).Once()

		result, err := client.GetOwnerOf(ctx, tokenID)

		require.NoError(t, err)
		assert.Equal(t, expectedOwner, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(int64(tokenID))).
			Return(common.Address{}, assert.AnError).Once()

		result, err := client.GetOwnerOf(ctx, tokenID)

		require.Error(t, err)
		assert.Equal(t, common.Address{}, result)
		assert.Contains(t, err.Error(), "failed to get owner")

		mockContract.AssertExpectations(t)
	})
}

// TestClient_GetTokenURI tests the actual Client.GetTokenURI method
func TestClient_GetTokenURI(t *testing.T) {
	mockContract := new(MockTicketsContract)
	client := &Client{
		ticketsContract: mockContract,
	}

	ctx := context.Background()
	tokenID := uint64(999)
	expectedURI := "https://api.example.com/metadata/999"

	t.Run("success", func(t *testing.T) {
		mockContract.On("TokenURI", mock.Anything, big.NewInt(int64(tokenID))).
			Return(expectedURI, nil).Once()

		result, err := client.GetTokenURI(ctx, tokenID)

		require.NoError(t, err)
		assert.Equal(t, expectedURI, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockContract.On("TokenURI", mock.Anything, big.NewInt(int64(tokenID))).
			Return("", assert.AnError).Once()

		result, err := client.GetTokenURI(ctx, tokenID)

		require.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to get token URI")

		mockContract.AssertExpectations(t)
	})
}

// TestClient_GetBalanceOf tests the actual Client.GetBalanceOf method
func TestClient_GetBalanceOf(t *testing.T) {
	mockContract := new(MockTicketsContract)
	client := &Client{
		ticketsContract: mockContract,
	}

	ctx := context.Background()
	owner := common.HexToAddress("0xABCDEF1234567890123456789012345678901234")
	expectedBalance := big.NewInt(5)

	t.Run("success", func(t *testing.T) {
		mockContract.On("BalanceOf", mock.Anything, owner).
			Return(expectedBalance, nil).Once()

		result, err := client.GetBalanceOf(ctx, owner)

		require.NoError(t, err)
		assert.Equal(t, expectedBalance, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("zero balance", func(t *testing.T) {
		mockContract.On("BalanceOf", mock.Anything, owner).
			Return(big.NewInt(0), nil).Once()

		result, err := client.GetBalanceOf(ctx, owner)

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(0), result)

		mockContract.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockContract.On("BalanceOf", mock.Anything, owner).
			Return((*big.Int)(nil), assert.AnError).Once()

		result, err := client.GetBalanceOf(ctx, owner)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get balance")

		mockContract.AssertExpectations(t)
	})
}
