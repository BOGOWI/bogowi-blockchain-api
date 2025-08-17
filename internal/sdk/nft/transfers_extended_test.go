package nft

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransferErrorHandling(t *testing.T) {
	ctx := context.Background()
	mockContract := new(MockTicketsContract)
	mockClient := new(MockEthClient)

	client := &TestClient{
		ticketsContract: mockContract,
		ethClient:       mockClient,
		auth: &bind.TransactOpts{
			From:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Nonce: big.NewInt(1),
		},
	}

	to := common.HexToAddress("0x9876543210987654321098765432109876543210")
	tokenID := uint64(1)

	t.Run("IsTransferable error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(false, assert.AnError).Once()

		_, err := client.Transfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check transferability")

		mockContract.AssertExpectations(t)
	})

	t.Run("SuggestGasPrice error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(nil, assert.AnError).Once()

		_, err := client.Transfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get gas price")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("GetOwnerOf error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(common.Address{}, assert.AnError).Once()

		_, err := client.Transfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get current owner")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("TransferFrom error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("TransferFrom", mock.Anything, client.auth.From, to, big.NewInt(1)).Return(nil, assert.AnError).Once()

		_, err := client.Transfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to transfer ticket")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("WaitForTransaction error", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("TransferFrom", mock.Anything, client.auth.From, to, big.NewInt(1)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(nil, assert.AnError).Once()

		result, err := client.Transfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction failed")
		assert.Equal(t, tx, result) // Transaction is still returned even on error

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}

func TestSafeTransferErrorHandling(t *testing.T) {
	ctx := context.Background()
	mockContract := new(MockTicketsContract)
	mockClient := new(MockEthClient)

	client := &TestClient{
		ticketsContract: mockContract,
		ethClient:       mockClient,
		auth: &bind.TransactOpts{
			From:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Nonce: big.NewInt(1),
		},
	}

	to := common.HexToAddress("0x9876543210987654321098765432109876543210")
	tokenID := uint64(1)

	t.Run("not transferable", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(false, nil).Once()

		_, err := client.SafeTransfer(ctx, to, tokenID)
		assert.ErrorIs(t, err, ErrNotTransferable)

		mockContract.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		otherOwner := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(otherOwner, nil).Once()

		_, err := client.SafeTransfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sender is not the owner")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("SafeTransferFrom error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom", mock.Anything, client.auth.From, to, big.NewInt(1)).Return(nil, assert.AnError).Once()

		_, err := client.SafeTransfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to safe transfer ticket")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("transaction failed", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 0}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom", mock.Anything, client.auth.From, to, big.NewInt(1)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		_, err := client.SafeTransfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "safe transfer transaction failed")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}

func TestSafeTransferWithDataErrorHandling(t *testing.T) {
	ctx := context.Background()
	mockContract := new(MockTicketsContract)
	mockClient := new(MockEthClient)

	client := &TestClient{
		ticketsContract: mockContract,
		ethClient:       mockClient,
		auth: &bind.TransactOpts{
			From:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Nonce: big.NewInt(1),
		},
	}

	to := common.HexToAddress("0x9876543210987654321098765432109876543210")
	tokenID := uint64(1)
	data := []byte("test data")

	t.Run("IsTransferable error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(false, assert.AnError).Once()

		_, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check transferability")

		mockContract.AssertExpectations(t)
	})

	t.Run("not transferable", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(false, nil).Once()

		_, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.ErrorIs(t, err, ErrNotTransferable)

		mockContract.AssertExpectations(t)
	})

	t.Run("SuggestGasPrice error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(nil, assert.AnError).Once()

		_, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get gas price")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("GetOwnerOf error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(common.Address{}, assert.AnError).Once()

		_, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get current owner")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		otherOwner := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(otherOwner, nil).Once()

		_, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sender is not the owner")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("SafeTransferFrom0 error", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom0", mock.Anything, client.auth.From, to, big.NewInt(1), data).Return(nil, assert.AnError).Once()

		_, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to safe transfer with data")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("transaction failed", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 0}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom0", mock.Anything, client.auth.From, to, big.NewInt(1), data).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		_, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "safe transfer with data transaction failed")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("nil data", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom0", mock.Anything, client.auth.From, to, big.NewInt(1), []byte(nil)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.SafeTransferWithData(ctx, to, tokenID, nil)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("empty data", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}
		emptyData := []byte{}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom0", mock.Anything, client.auth.From, to, big.NewInt(1), emptyData).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.SafeTransferWithData(ctx, to, tokenID, emptyData)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}

func TestApprovalMethods(t *testing.T) {
	mockContract := new(MockTicketsContract)

	t.Run("Approve with zero address", func(t *testing.T) {
		mockContract.On("Approve", mock.Anything, common.Address{}, big.NewInt(1)).Return(nil, nil).Once()

		_, err := mockContract.Approve(&bind.TransactOpts{}, common.Address{}, big.NewInt(1))
		assert.NoError(t, err)

		mockContract.AssertExpectations(t)
	})

	t.Run("GetApproved with zero return", func(t *testing.T) {
		mockContract.On("GetApproved", mock.Anything, big.NewInt(999)).Return(common.Address{}, nil).Once()

		result, err := mockContract.GetApproved(&bind.CallOpts{}, big.NewInt(999))
		assert.NoError(t, err)
		assert.Equal(t, common.Address{}, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("IsApprovedForAll same owner and operator", func(t *testing.T) {
		addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
		mockContract.On("IsApprovedForAll", mock.Anything, addr, addr).Return(true, nil).Once()

		result, err := mockContract.IsApprovedForAll(&bind.CallOpts{}, addr, addr)
		assert.NoError(t, err)
		assert.True(t, result)

		mockContract.AssertExpectations(t)
	})
}

func TestTransferBoundaryConditions(t *testing.T) {
	ctx := context.Background()
	mockContract := new(MockTicketsContract)
	mockClient := new(MockEthClient)

	client := &TestClient{
		ticketsContract: mockContract,
		ethClient:       mockClient,
		auth: &bind.TransactOpts{
			From:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Nonce: big.NewInt(1),
		},
	}

	t.Run("transfer to same address", func(t *testing.T) {
		from := client.auth.From
		tokenID := uint64(1)
		tx := types.NewTransaction(1, from, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(from, nil).Once()
		mockContract.On("TransferFrom", mock.Anything, from, from, big.NewInt(1)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.Transfer(ctx, from, tokenID)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("transfer to zero address", func(t *testing.T) {
		zeroAddr := common.Address{}
		tokenID := uint64(1)
		tx := types.NewTransaction(1, zeroAddr, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("TransferFrom", mock.Anything, client.auth.From, zeroAddr, big.NewInt(1)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.Transfer(ctx, zeroAddr, tokenID)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("transfer with max uint64 token ID", func(t *testing.T) {
		to := common.HexToAddress("0x9876543210987654321098765432109876543210")
		tokenID := uint64(^uint64(0)) // Max uint64
		tokenIDBig := new(big.Int).SetUint64(tokenID)
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}

		mockContract.On("IsTransferable", mock.Anything, tokenIDBig).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, tokenIDBig).Return(client.auth.From, nil).Once()
		mockContract.On("TransferFrom", mock.Anything, client.auth.From, to, tokenIDBig).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.Transfer(ctx, to, tokenID)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}
