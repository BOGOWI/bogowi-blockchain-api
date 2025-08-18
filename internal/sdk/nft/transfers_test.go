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

// MockTicketsContract for testing
type MockTicketsContract struct {
	mock.Mock
}

func (m *MockTicketsContract) TransferFrom(opts *bind.TransactOpts, from, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	args := m.Called(opts, from, to, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) SafeTransferFrom(opts *bind.TransactOpts, from, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	args := m.Called(opts, from, to, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) SafeTransferFrom0(opts *bind.TransactOpts, from, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	args := m.Called(opts, from, to, tokenId, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) Approve(opts *bind.TransactOpts, spender common.Address, tokenId *big.Int) (*types.Transaction, error) {
	args := m.Called(opts, spender, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	args := m.Called(opts, operator, approved)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	args := m.Called(opts, tokenId)
	return args.Get(0).(common.Address), args.Error(1)
}

func (m *MockTicketsContract) IsApprovedForAll(opts *bind.CallOpts, owner, operator common.Address) (bool, error) {
	args := m.Called(opts, owner, operator)
	return args.Bool(0), args.Error(1)
}

func (m *MockTicketsContract) IsTransferable(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	args := m.Called(opts, tokenId)
	return args.Bool(0), args.Error(1)
}

func (m *MockTicketsContract) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	args := m.Called(opts, tokenId)
	return args.Get(0).(common.Address), args.Error(1)
}

func (m *MockTicketsContract) GetTicketData(opts *bind.CallOpts, tokenId *big.Int) (TicketDataContract, error) {
	args := m.Called(opts, tokenId)
	return args.Get(0).(TicketDataContract), args.Error(1)
}

func (m *MockTicketsContract) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	args := m.Called(opts, tokenId)
	return args.String(0), args.Error(1)
}

func (m *MockTicketsContract) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	args := m.Called(opts, owner)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockTicketsContract) MintTicket(opts *bind.TransactOpts, to common.Address, bookingId [32]byte, eventId [32]byte, utilityFlags uint32, transferUnlockAt uint64, expiresAt uint64, metadataURI string, rewardBasisPoints uint16) (*types.Transaction, error) {
	args := m.Called(opts, to, bookingId, eventId, utilityFlags, transferUnlockAt, expiresAt, metadataURI, rewardBasisPoints)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) MintBatch(opts *bind.TransactOpts, tos []common.Address, bookingIds [][32]byte, eventIds [][32]byte, utilityFlags []uint32, transferUnlockAts []uint64, expiresAts []uint64, metadataURIs []string, rewardBasisPoints []uint16) (*types.Transaction, error) {
	args := m.Called(opts, tos, bookingIds, eventIds, utilityFlags, transferUnlockAts, expiresAts, metadataURIs, rewardBasisPoints)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error) {
	args := m.Called(opts, newBaseURI)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) ParseTicketMinted(log types.Log) (*TicketMintedEvent, error) {
	args := m.Called(log)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TicketMintedEvent), args.Error(1)
}

func (m *MockTicketsContract) ExpireTicket(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	args := m.Called(opts, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) RedeemTicket(opts *bind.TransactOpts, redemptionData RedemptionDataContract) (*types.Transaction, error) {
	args := m.Called(opts, redemptionData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) UpdateTransferUnlock(opts *bind.TransactOpts, tokenId *big.Int, newUnlockTime uint64) (*types.Transaction, error) {
	args := m.Called(opts, tokenId, newUnlockTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockTicketsContract) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	args := m.Called(opts, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

// MockEthClient for testing
type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := m.Called(ctx, txHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Receipt), args.Error(1)
}

func TestTransfer(t *testing.T) {
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

	t.Run("successful transfer", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("TransferFrom", mock.Anything, client.auth.From, to, big.NewInt(1)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.Transfer(ctx, to, tokenID)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("not transferable", func(t *testing.T) {
		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(false, nil).Once()

		_, err := client.Transfer(ctx, to, tokenID)
		assert.ErrorIs(t, err, ErrNotTransferable)

		mockContract.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		otherOwner := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(otherOwner, nil).Once()

		_, err := client.Transfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sender is not the owner")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("transaction failed", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 0}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("TransferFrom", mock.Anything, client.auth.From, to, big.NewInt(1)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		_, err := client.Transfer(ctx, to, tokenID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transfer transaction failed")

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}

func TestSafeTransfer(t *testing.T) {
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

	t.Run("successful safe transfer", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom", mock.Anything, client.auth.From, to, big.NewInt(1)).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.SafeTransfer(ctx, to, tokenID)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}

func TestSafeTransferWithData(t *testing.T) {
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

	t.Run("successful safe transfer with data", func(t *testing.T) {
		tx := types.NewTransaction(1, to, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
		receipt := &types.Receipt{Status: 1}

		mockContract.On("IsTransferable", mock.Anything, big.NewInt(1)).Return(true, nil).Once()
		mockClient.On("SuggestGasPrice", ctx).Return(big.NewInt(1000000000), nil).Once()
		mockContract.On("OwnerOf", mock.Anything, big.NewInt(1)).Return(client.auth.From, nil).Once()
		mockContract.On("SafeTransferFrom0", mock.Anything, client.auth.From, to, big.NewInt(1), data).Return(tx, nil).Once()
		mockClient.On("TransactionReceipt", ctx, tx.Hash()).Return(receipt, nil).Once()

		result, err := client.SafeTransferWithData(ctx, to, tokenID, data)
		assert.NoError(t, err)
		assert.Equal(t, tx, result)

		mockContract.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}
