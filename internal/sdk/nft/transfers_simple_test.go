package nft

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApprove(t *testing.T) {
	mockContract := new(MockTicketsContract)

	// Test just the contract interaction
	spender := common.HexToAddress("0x9876543210987654321098765432109876543210")

	t.Run("contract approve call", func(t *testing.T) {
		mockContract.On("Approve", mock.Anything, spender, big.NewInt(1)).Return(nil, nil).Once()

		// Direct contract call test
		_, err := mockContract.Approve(&bind.TransactOpts{}, spender, big.NewInt(1))
		assert.NoError(t, err)

		mockContract.AssertExpectations(t)
	})

	t.Run("contract approve error", func(t *testing.T) {
		mockContract.On("Approve", mock.Anything, spender, big.NewInt(1)).Return(nil, errors.New("approval failed")).Once()

		_, err := mockContract.Approve(&bind.TransactOpts{}, spender, big.NewInt(1))
		assert.Error(t, err)

		mockContract.AssertExpectations(t)
	})
}

func TestSetApprovalForAll(t *testing.T) {
	mockContract := new(MockTicketsContract)
	operator := common.HexToAddress("0x9876543210987654321098765432109876543210")

	t.Run("set approval for all", func(t *testing.T) {
		mockContract.On("SetApprovalForAll", mock.Anything, operator, true).Return(nil, nil).Once()

		_, err := mockContract.SetApprovalForAll(&bind.TransactOpts{}, operator, true)
		assert.NoError(t, err)

		mockContract.AssertExpectations(t)
	})

	t.Run("revoke approval for all", func(t *testing.T) {
		mockContract.On("SetApprovalForAll", mock.Anything, operator, false).Return(nil, nil).Once()

		_, err := mockContract.SetApprovalForAll(&bind.TransactOpts{}, operator, false)
		assert.NoError(t, err)

		mockContract.AssertExpectations(t)
	})
}

func TestGetApproved(t *testing.T) {
	mockContract := new(MockTicketsContract)
	approved := common.HexToAddress("0x9876543210987654321098765432109876543210")

	t.Run("get approved address", func(t *testing.T) {
		mockContract.On("GetApproved", mock.Anything, big.NewInt(1)).Return(approved, nil).Once()

		result, err := mockContract.GetApproved(&bind.CallOpts{}, big.NewInt(1))
		assert.NoError(t, err)
		assert.Equal(t, approved, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("error getting approved", func(t *testing.T) {
		mockContract.On("GetApproved", mock.Anything, big.NewInt(1)).Return(common.Address{}, errors.New("contract error")).Once()

		_, err := mockContract.GetApproved(&bind.CallOpts{}, big.NewInt(1))
		assert.Error(t, err)

		mockContract.AssertExpectations(t)
	})
}

func TestIsApprovedForAll(t *testing.T) {
	mockContract := new(MockTicketsContract)
	owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
	operator := common.HexToAddress("0x9876543210987654321098765432109876543210")

	t.Run("is approved", func(t *testing.T) {
		mockContract.On("IsApprovedForAll", mock.Anything, owner, operator).Return(true, nil).Once()

		result, err := mockContract.IsApprovedForAll(&bind.CallOpts{}, owner, operator)
		assert.NoError(t, err)
		assert.True(t, result)

		mockContract.AssertExpectations(t)
	})

	t.Run("not approved", func(t *testing.T) {
		mockContract.On("IsApprovedForAll", mock.Anything, owner, operator).Return(false, nil).Once()

		result, err := mockContract.IsApprovedForAll(&bind.CallOpts{}, owner, operator)
		assert.NoError(t, err)
		assert.False(t, result)

		mockContract.AssertExpectations(t)
	})
}
