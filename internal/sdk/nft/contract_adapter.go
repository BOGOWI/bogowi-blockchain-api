package nft

import (
	"math/big"

	"bogowi-blockchain-go/internal/sdk/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ContractAdapter wraps the actual BOGOWITickets contract to implement TicketsContractInterface
type ContractAdapter struct {
	contract *contracts.BOGOWITickets
}

// NewContractAdapter creates a new adapter for the BOGOWITickets contract
func NewContractAdapter(contract *contracts.BOGOWITickets) TicketsContractInterface {
	return &ContractAdapter{contract: contract}
}

func (a *ContractAdapter) TransferFrom(opts *bind.TransactOpts, from, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return a.contract.TransferFrom(opts, from, to, tokenId)
}

func (a *ContractAdapter) SafeTransferFrom(opts *bind.TransactOpts, from, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return a.contract.SafeTransferFrom(opts, from, to, tokenId)
}

func (a *ContractAdapter) SafeTransferFrom0(opts *bind.TransactOpts, from, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return a.contract.SafeTransferFrom0(opts, from, to, tokenId, data)
}

func (a *ContractAdapter) Approve(opts *bind.TransactOpts, spender common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return a.contract.Approve(opts, spender, tokenId)
}

func (a *ContractAdapter) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return a.contract.SetApprovalForAll(opts, operator, approved)
}

func (a *ContractAdapter) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	return a.contract.GetApproved(opts, tokenId)
}

func (a *ContractAdapter) IsApprovedForAll(opts *bind.CallOpts, owner, operator common.Address) (bool, error) {
	return a.contract.IsApprovedForAll(opts, owner, operator)
}

func (a *ContractAdapter) IsTransferable(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	return a.contract.IsTransferable(opts, tokenId)
}

func (a *ContractAdapter) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	return a.contract.OwnerOf(opts, tokenId)
}

func (a *ContractAdapter) GetTicketData(opts *bind.CallOpts, tokenId *big.Int) (TicketDataContract, error) {
	data, err := a.contract.GetTicketData(opts, tokenId)
	if err != nil {
		return TicketDataContract{}, err
	}

	// Convert contracts.IBOGOWITicketsTicketData to TicketDataContract
	return TicketDataContract{
		BookingId:                  data.BookingId,
		EventId:                    data.EventId,
		TransferUnlockAt:           data.TransferUnlockAt,
		ExpiresAt:                  data.ExpiresAt,
		UtilityFlags:               data.UtilityFlags,
		State:                      data.State,
		NonTransferableAfterRedeem: data.NonTransferableAfterRedeem,
		BurnOnRedeem:               data.BurnOnRedeem,
	}, nil
}

func (a *ContractAdapter) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	return a.contract.TokenURI(opts, tokenId)
}

func (a *ContractAdapter) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	return a.contract.BalanceOf(opts, owner)
}

func (a *ContractAdapter) MintTicket(opts *bind.TransactOpts, to common.Address, bookingId [32]byte, eventId [32]byte, utilityFlags uint32, transferUnlockAt uint64, expiresAt uint64, metadataURI string, rewardBasisPoints uint16) (*types.Transaction, error) {
	params := contracts.IBOGOWITicketsMintParams{
		To:                to,
		BookingId:         bookingId,
		EventId:           eventId,
		UtilityFlags:      utilityFlags,
		TransferUnlockAt:  transferUnlockAt,
		ExpiresAt:         expiresAt,
		MetadataURI:       metadataURI,
		RewardBasisPoints: big.NewInt(int64(rewardBasisPoints)),
	}
	return a.contract.MintTicket(opts, params)
}

func (a *ContractAdapter) MintBatch(opts *bind.TransactOpts, tos []common.Address, bookingIds [][32]byte, eventIds [][32]byte, utilityFlags []uint32, transferUnlockAts []uint64, expiresAts []uint64, metadataURIs []string, rewardBasisPoints []uint16) (*types.Transaction, error) {
	// Convert to array of MintParams
	params := make([]contracts.IBOGOWITicketsMintParams, len(tos))
	for i := range tos {
		params[i] = contracts.IBOGOWITicketsMintParams{
			To:                tos[i],
			BookingId:         bookingIds[i],
			EventId:           eventIds[i],
			UtilityFlags:      utilityFlags[i],
			TransferUnlockAt:  transferUnlockAts[i],
			ExpiresAt:         expiresAts[i],
			MetadataURI:       metadataURIs[i],
			RewardBasisPoints: big.NewInt(int64(rewardBasisPoints[i])),
		}
	}
	return a.contract.MintBatch(opts, params)
}

func (a *ContractAdapter) SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error) {
	return a.contract.SetBaseURI(opts, newBaseURI)
}

func (a *ContractAdapter) ParseTicketMinted(log types.Log) (*TicketMintedEvent, error) {
	event, err := a.contract.ParseTicketMinted(log)
	if err != nil {
		return nil, err
	}

	return &TicketMintedEvent{
		TokenId:   event.TokenId,
		To:        event.Buyer,
		BookingId: event.BookingIdHash,
		EventId:   event.EventIdHash,
		Raw:       event.Raw,
	}, nil
}

func (a *ContractAdapter) ExpireTicket(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return a.contract.ExpireTicket(opts, tokenId)
}

func (a *ContractAdapter) RedeemTicket(opts *bind.TransactOpts, redemptionData RedemptionDataContract) (*types.Transaction, error) {
	// Convert to contract type
	contractData := contracts.IBOGOWITicketsRedemptionData{
		TokenId:   redemptionData.TokenId,
		Redeemer:  redemptionData.Redeemer,
		Nonce:     redemptionData.Nonce,
		Deadline:  redemptionData.Deadline,
		ChainId:   redemptionData.ChainId,
		Signature: redemptionData.Signature,
	}
	return a.contract.RedeemTicket(opts, contractData)
}

func (a *ContractAdapter) UpdateTransferUnlock(opts *bind.TransactOpts, tokenId *big.Int, newUnlockTime uint64) (*types.Transaction, error) {
	return a.contract.UpdateTransferUnlock(opts, tokenId, newUnlockTime)
}

func (a *ContractAdapter) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return a.contract.Burn(opts, tokenId)
}
