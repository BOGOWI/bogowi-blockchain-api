// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IBOGOWITicketsMintParams is an auto generated low-level Go binding around an user-defined struct.
type IBOGOWITicketsMintParams struct {
	To                common.Address
	BookingId         [32]byte
	EventId           [32]byte
	UtilityFlags      uint32
	TransferUnlockAt  uint64
	ExpiresAt         uint64
	MetadataURI       string
	RewardBasisPoints *big.Int
}

// IBOGOWITicketsRedemptionData is an auto generated low-level Go binding around an user-defined struct.
type IBOGOWITicketsRedemptionData struct {
	TokenId   *big.Int
	Redeemer  common.Address
	Nonce     *big.Int
	Deadline  *big.Int
	ChainId   *big.Int
	Signature []byte
}

// IBOGOWITicketsTicketData is an auto generated low-level Go binding around an user-defined struct.
type IBOGOWITicketsTicketData struct {
	BookingId                  [32]byte
	EventId                    [32]byte
	TransferUnlockAt           uint64
	ExpiresAt                  uint64
	UtilityFlags               uint32
	State                      uint8
	NonTransferableAfterRedeem bool
	BurnOnRedeem               bool
}

// BOGOWITicketsMetaData contains all meta data concerning the BOGOWITickets contract.
var BOGOWITicketsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_roleManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_conservationDAO\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"ERC2981InvalidDefaultRoyalty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC2981InvalidDefaultRoyaltyReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"ERC2981InvalidTokenRoyalty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC2981InvalidTokenRoyaltyReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RoleManagerNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"StringsInsufficientHexLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"UnauthorizedRole\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newBaseURI\",\"type\":\"string\"}],\"name\":\"BaseURIUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"BatchMintStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"datakyteNftId\",\"type\":\"string\"}],\"name\":\"DatakyteMetadataLinked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"NonceUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"roleManagerAddress\",\"type\":\"address\"}],\"name\":\"RoleManagerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"feeBasisPoints\",\"type\":\"uint96\"}],\"name\":\"RoyaltyInfoUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"TicketBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"TicketExpired\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"bookingIdHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"eventIdHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardBasisPoints\",\"type\":\"uint256\"}],\"name\":\"TicketMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"redeemedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"TicketRedeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newUnlockTime\",\"type\":\"uint64\"}],\"name\":\"TransferUnlockUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CAMINO_MAINNET_CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CAMINO_TESTNET_CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COMMIT_DELAY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ROYALTY_BPS\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GAS_PER_MINT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIAL_TOKEN_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_BATCH_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REDEMPTION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"conservationDAO\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"expireTicket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRoleManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getTicketData\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"bookingId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"eventId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"transferUnlockAt\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"expiresAt\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"utilityFlags\",\"type\":\"uint32\"},{\"internalType\":\"enumIBOGOWITickets.TicketState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"nonTransferableAfterRedeem\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"burnOnRedeem\",\"type\":\"bool\"}],\"internalType\":\"structIBOGOWITickets.TicketData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"isExpired\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"isRedeemed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"isTransferable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"bookingId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"eventId\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"utilityFlags\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"transferUnlockAt\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"expiresAt\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"metadataURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"rewardBasisPoints\",\"type\":\"uint256\"}],\"internalType\":\"structIBOGOWITickets.MintParams[]\",\"name\":\"params\",\"type\":\"tuple[]\"}],\"name\":\"mintBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"bookingId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"eventId\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"utilityFlags\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"transferUnlockAt\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"expiresAt\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"metadataURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"rewardBasisPoints\",\"type\":\"uint256\"}],\"internalType\":\"structIBOGOWITickets.MintParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"mintTicket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structIBOGOWITickets.RedemptionData\",\"name\":\"redemptionData\",\"type\":\"tuple\"}],\"name\":\"redeemTicket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roleManager\",\"outputs\":[{\"internalType\":\"contractIRoleManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newBaseURI\",\"type\":\"string\"}],\"name\":\"setBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"feeBasisPoints\",\"type\":\"uint96\"}],\"name\":\"setRoyaltyInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"newUnlockTime\",\"type\":\"uint64\"}],\"name\":\"updateTransferUnlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structIBOGOWITickets.RedemptionData\",\"name\":\"redemptionData\",\"type\":\"tuple\"}],\"name\":\"verifyRedemptionSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BOGOWITicketsABI is the input ABI used to generate the binding from.
// Deprecated: Use BOGOWITicketsMetaData.ABI instead.
var BOGOWITicketsABI = BOGOWITicketsMetaData.ABI

// BOGOWITickets is an auto generated Go binding around an Ethereum contract.
type BOGOWITickets struct {
	BOGOWITicketsCaller     // Read-only binding to the contract
	BOGOWITicketsTransactor // Write-only binding to the contract
	BOGOWITicketsFilterer   // Log filterer for contract events
}

// BOGOWITicketsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BOGOWITicketsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BOGOWITicketsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BOGOWITicketsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BOGOWITicketsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BOGOWITicketsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BOGOWITicketsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BOGOWITicketsSession struct {
	Contract     *BOGOWITickets    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BOGOWITicketsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BOGOWITicketsCallerSession struct {
	Contract *BOGOWITicketsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BOGOWITicketsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BOGOWITicketsTransactorSession struct {
	Contract     *BOGOWITicketsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BOGOWITicketsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BOGOWITicketsRaw struct {
	Contract *BOGOWITickets // Generic contract binding to access the raw methods on
}

// BOGOWITicketsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BOGOWITicketsCallerRaw struct {
	Contract *BOGOWITicketsCaller // Generic read-only contract binding to access the raw methods on
}

// BOGOWITicketsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BOGOWITicketsTransactorRaw struct {
	Contract *BOGOWITicketsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBOGOWITickets creates a new instance of BOGOWITickets, bound to a specific deployed contract.
func NewBOGOWITickets(address common.Address, backend bind.ContractBackend) (*BOGOWITickets, error) {
	contract, err := bindBOGOWITickets(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BOGOWITickets{BOGOWITicketsCaller: BOGOWITicketsCaller{contract: contract}, BOGOWITicketsTransactor: BOGOWITicketsTransactor{contract: contract}, BOGOWITicketsFilterer: BOGOWITicketsFilterer{contract: contract}}, nil
}

// NewBOGOWITicketsCaller creates a new read-only instance of BOGOWITickets, bound to a specific deployed contract.
func NewBOGOWITicketsCaller(address common.Address, caller bind.ContractCaller) (*BOGOWITicketsCaller, error) {
	contract, err := bindBOGOWITickets(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsCaller{contract: contract}, nil
}

// NewBOGOWITicketsTransactor creates a new write-only instance of BOGOWITickets, bound to a specific deployed contract.
func NewBOGOWITicketsTransactor(address common.Address, transactor bind.ContractTransactor) (*BOGOWITicketsTransactor, error) {
	contract, err := bindBOGOWITickets(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsTransactor{contract: contract}, nil
}

// NewBOGOWITicketsFilterer creates a new log filterer instance of BOGOWITickets, bound to a specific deployed contract.
func NewBOGOWITicketsFilterer(address common.Address, filterer bind.ContractFilterer) (*BOGOWITicketsFilterer, error) {
	contract, err := bindBOGOWITickets(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsFilterer{contract: contract}, nil
}

// bindBOGOWITickets binds a generic wrapper to an already deployed contract.
func bindBOGOWITickets(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BOGOWITicketsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BOGOWITickets *BOGOWITicketsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BOGOWITickets.Contract.BOGOWITicketsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BOGOWITickets *BOGOWITicketsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.BOGOWITicketsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BOGOWITickets *BOGOWITicketsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.BOGOWITicketsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BOGOWITickets *BOGOWITicketsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BOGOWITickets.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BOGOWITickets *BOGOWITicketsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BOGOWITickets *BOGOWITicketsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.contract.Transact(opts, method, params...)
}

// CAMINOMAINNETCHAINID is a free data retrieval call binding the contract method 0xb2c98fc7.
//
// Solidity: function CAMINO_MAINNET_CHAIN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCaller) CAMINOMAINNETCHAINID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "CAMINO_MAINNET_CHAIN_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CAMINOMAINNETCHAINID is a free data retrieval call binding the contract method 0xb2c98fc7.
//
// Solidity: function CAMINO_MAINNET_CHAIN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) CAMINOMAINNETCHAINID() (*big.Int, error) {
	return _BOGOWITickets.Contract.CAMINOMAINNETCHAINID(&_BOGOWITickets.CallOpts)
}

// CAMINOMAINNETCHAINID is a free data retrieval call binding the contract method 0xb2c98fc7.
//
// Solidity: function CAMINO_MAINNET_CHAIN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCallerSession) CAMINOMAINNETCHAINID() (*big.Int, error) {
	return _BOGOWITickets.Contract.CAMINOMAINNETCHAINID(&_BOGOWITickets.CallOpts)
}

// CAMINOTESTNETCHAINID is a free data retrieval call binding the contract method 0x8cd5a48c.
//
// Solidity: function CAMINO_TESTNET_CHAIN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCaller) CAMINOTESTNETCHAINID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "CAMINO_TESTNET_CHAIN_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CAMINOTESTNETCHAINID is a free data retrieval call binding the contract method 0x8cd5a48c.
//
// Solidity: function CAMINO_TESTNET_CHAIN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) CAMINOTESTNETCHAINID() (*big.Int, error) {
	return _BOGOWITickets.Contract.CAMINOTESTNETCHAINID(&_BOGOWITickets.CallOpts)
}

// CAMINOTESTNETCHAINID is a free data retrieval call binding the contract method 0x8cd5a48c.
//
// Solidity: function CAMINO_TESTNET_CHAIN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCallerSession) CAMINOTESTNETCHAINID() (*big.Int, error) {
	return _BOGOWITickets.Contract.CAMINOTESTNETCHAINID(&_BOGOWITickets.CallOpts)
}

// COMMITDELAY is a free data retrieval call binding the contract method 0x50857fb2.
//
// Solidity: function COMMIT_DELAY() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCaller) COMMITDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "COMMIT_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COMMITDELAY is a free data retrieval call binding the contract method 0x50857fb2.
//
// Solidity: function COMMIT_DELAY() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) COMMITDELAY() (*big.Int, error) {
	return _BOGOWITickets.Contract.COMMITDELAY(&_BOGOWITickets.CallOpts)
}

// COMMITDELAY is a free data retrieval call binding the contract method 0x50857fb2.
//
// Solidity: function COMMIT_DELAY() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCallerSession) COMMITDELAY() (*big.Int, error) {
	return _BOGOWITickets.Contract.COMMITDELAY(&_BOGOWITickets.CallOpts)
}

// DEFAULTROYALTYBPS is a free data retrieval call binding the contract method 0x5b1ab434.
//
// Solidity: function DEFAULT_ROYALTY_BPS() view returns(uint96)
func (_BOGOWITickets *BOGOWITicketsCaller) DEFAULTROYALTYBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "DEFAULT_ROYALTY_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTROYALTYBPS is a free data retrieval call binding the contract method 0x5b1ab434.
//
// Solidity: function DEFAULT_ROYALTY_BPS() view returns(uint96)
func (_BOGOWITickets *BOGOWITicketsSession) DEFAULTROYALTYBPS() (*big.Int, error) {
	return _BOGOWITickets.Contract.DEFAULTROYALTYBPS(&_BOGOWITickets.CallOpts)
}

// DEFAULTROYALTYBPS is a free data retrieval call binding the contract method 0x5b1ab434.
//
// Solidity: function DEFAULT_ROYALTY_BPS() view returns(uint96)
func (_BOGOWITickets *BOGOWITicketsCallerSession) DEFAULTROYALTYBPS() (*big.Int, error) {
	return _BOGOWITickets.Contract.DEFAULTROYALTYBPS(&_BOGOWITickets.CallOpts)
}

// GASPERMINT is a free data retrieval call binding the contract method 0xbf17374e.
//
// Solidity: function GAS_PER_MINT() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCaller) GASPERMINT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "GAS_PER_MINT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GASPERMINT is a free data retrieval call binding the contract method 0xbf17374e.
//
// Solidity: function GAS_PER_MINT() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) GASPERMINT() (*big.Int, error) {
	return _BOGOWITickets.Contract.GASPERMINT(&_BOGOWITickets.CallOpts)
}

// GASPERMINT is a free data retrieval call binding the contract method 0xbf17374e.
//
// Solidity: function GAS_PER_MINT() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCallerSession) GASPERMINT() (*big.Int, error) {
	return _BOGOWITickets.Contract.GASPERMINT(&_BOGOWITickets.CallOpts)
}

// INITIALTOKENID is a free data retrieval call binding the contract method 0x6b2036fb.
//
// Solidity: function INITIAL_TOKEN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCaller) INITIALTOKENID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "INITIAL_TOKEN_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALTOKENID is a free data retrieval call binding the contract method 0x6b2036fb.
//
// Solidity: function INITIAL_TOKEN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) INITIALTOKENID() (*big.Int, error) {
	return _BOGOWITickets.Contract.INITIALTOKENID(&_BOGOWITickets.CallOpts)
}

// INITIALTOKENID is a free data retrieval call binding the contract method 0x6b2036fb.
//
// Solidity: function INITIAL_TOKEN_ID() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCallerSession) INITIALTOKENID() (*big.Int, error) {
	return _BOGOWITickets.Contract.INITIALTOKENID(&_BOGOWITickets.CallOpts)
}

// MAXBATCHSIZE is a free data retrieval call binding the contract method 0xcfdbf254.
//
// Solidity: function MAX_BATCH_SIZE() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCaller) MAXBATCHSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "MAX_BATCH_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBATCHSIZE is a free data retrieval call binding the contract method 0xcfdbf254.
//
// Solidity: function MAX_BATCH_SIZE() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) MAXBATCHSIZE() (*big.Int, error) {
	return _BOGOWITickets.Contract.MAXBATCHSIZE(&_BOGOWITickets.CallOpts)
}

// MAXBATCHSIZE is a free data retrieval call binding the contract method 0xcfdbf254.
//
// Solidity: function MAX_BATCH_SIZE() view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCallerSession) MAXBATCHSIZE() (*big.Int, error) {
	return _BOGOWITickets.Contract.MAXBATCHSIZE(&_BOGOWITickets.CallOpts)
}

// REDEMPTIONTYPEHASH is a free data retrieval call binding the contract method 0x0010322b.
//
// Solidity: function REDEMPTION_TYPEHASH() view returns(bytes32)
func (_BOGOWITickets *BOGOWITicketsCaller) REDEMPTIONTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "REDEMPTION_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// REDEMPTIONTYPEHASH is a free data retrieval call binding the contract method 0x0010322b.
//
// Solidity: function REDEMPTION_TYPEHASH() view returns(bytes32)
func (_BOGOWITickets *BOGOWITicketsSession) REDEMPTIONTYPEHASH() ([32]byte, error) {
	return _BOGOWITickets.Contract.REDEMPTIONTYPEHASH(&_BOGOWITickets.CallOpts)
}

// REDEMPTIONTYPEHASH is a free data retrieval call binding the contract method 0x0010322b.
//
// Solidity: function REDEMPTION_TYPEHASH() view returns(bytes32)
func (_BOGOWITickets *BOGOWITicketsCallerSession) REDEMPTIONTYPEHASH() ([32]byte, error) {
	return _BOGOWITickets.Contract.REDEMPTIONTYPEHASH(&_BOGOWITickets.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _BOGOWITickets.Contract.BalanceOf(&_BOGOWITickets.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BOGOWITickets *BOGOWITicketsCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _BOGOWITickets.Contract.BalanceOf(&_BOGOWITickets.CallOpts, owner)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_BOGOWITickets *BOGOWITicketsCaller) BaseURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "baseURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_BOGOWITickets *BOGOWITicketsSession) BaseURI() (string, error) {
	return _BOGOWITickets.Contract.BaseURI(&_BOGOWITickets.CallOpts)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_BOGOWITickets *BOGOWITicketsCallerSession) BaseURI() (string, error) {
	return _BOGOWITickets.Contract.BaseURI(&_BOGOWITickets.CallOpts)
}

// ConservationDAO is a free data retrieval call binding the contract method 0xd5f2bfbf.
//
// Solidity: function conservationDAO() view returns(address)
func (_BOGOWITickets *BOGOWITicketsCaller) ConservationDAO(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "conservationDAO")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ConservationDAO is a free data retrieval call binding the contract method 0xd5f2bfbf.
//
// Solidity: function conservationDAO() view returns(address)
func (_BOGOWITickets *BOGOWITicketsSession) ConservationDAO() (common.Address, error) {
	return _BOGOWITickets.Contract.ConservationDAO(&_BOGOWITickets.CallOpts)
}

// ConservationDAO is a free data retrieval call binding the contract method 0xd5f2bfbf.
//
// Solidity: function conservationDAO() view returns(address)
func (_BOGOWITickets *BOGOWITicketsCallerSession) ConservationDAO() (common.Address, error) {
	return _BOGOWITickets.Contract.ConservationDAO(&_BOGOWITickets.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_BOGOWITickets *BOGOWITicketsCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_BOGOWITickets *BOGOWITicketsSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _BOGOWITickets.Contract.Eip712Domain(&_BOGOWITickets.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_BOGOWITickets *BOGOWITicketsCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _BOGOWITickets.Contract.Eip712Domain(&_BOGOWITickets.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BOGOWITickets *BOGOWITicketsCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BOGOWITickets *BOGOWITicketsSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _BOGOWITickets.Contract.GetApproved(&_BOGOWITickets.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BOGOWITickets *BOGOWITicketsCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _BOGOWITickets.Contract.GetApproved(&_BOGOWITickets.CallOpts, tokenId)
}

// GetRoleManager is a free data retrieval call binding the contract method 0x51331ad7.
//
// Solidity: function getRoleManager() view returns(address)
func (_BOGOWITickets *BOGOWITicketsCaller) GetRoleManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "getRoleManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleManager is a free data retrieval call binding the contract method 0x51331ad7.
//
// Solidity: function getRoleManager() view returns(address)
func (_BOGOWITickets *BOGOWITicketsSession) GetRoleManager() (common.Address, error) {
	return _BOGOWITickets.Contract.GetRoleManager(&_BOGOWITickets.CallOpts)
}

// GetRoleManager is a free data retrieval call binding the contract method 0x51331ad7.
//
// Solidity: function getRoleManager() view returns(address)
func (_BOGOWITickets *BOGOWITicketsCallerSession) GetRoleManager() (common.Address, error) {
	return _BOGOWITickets.Contract.GetRoleManager(&_BOGOWITickets.CallOpts)
}

// GetTicketData is a free data retrieval call binding the contract method 0xd9895e6a.
//
// Solidity: function getTicketData(uint256 tokenId) view returns((bytes32,bytes32,uint64,uint64,uint32,uint8,bool,bool))
func (_BOGOWITickets *BOGOWITicketsCaller) GetTicketData(opts *bind.CallOpts, tokenId *big.Int) (IBOGOWITicketsTicketData, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "getTicketData", tokenId)

	if err != nil {
		return *new(IBOGOWITicketsTicketData), err
	}

	out0 := *abi.ConvertType(out[0], new(IBOGOWITicketsTicketData)).(*IBOGOWITicketsTicketData)

	return out0, err

}

// GetTicketData is a free data retrieval call binding the contract method 0xd9895e6a.
//
// Solidity: function getTicketData(uint256 tokenId) view returns((bytes32,bytes32,uint64,uint64,uint32,uint8,bool,bool))
func (_BOGOWITickets *BOGOWITicketsSession) GetTicketData(tokenId *big.Int) (IBOGOWITicketsTicketData, error) {
	return _BOGOWITickets.Contract.GetTicketData(&_BOGOWITickets.CallOpts, tokenId)
}

// GetTicketData is a free data retrieval call binding the contract method 0xd9895e6a.
//
// Solidity: function getTicketData(uint256 tokenId) view returns((bytes32,bytes32,uint64,uint64,uint32,uint8,bool,bool))
func (_BOGOWITickets *BOGOWITicketsCallerSession) GetTicketData(tokenId *big.Int) (IBOGOWITicketsTicketData, error) {
	return _BOGOWITickets.Contract.GetTicketData(&_BOGOWITickets.CallOpts, tokenId)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BOGOWITickets.Contract.HasRole(&_BOGOWITickets.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BOGOWITickets.Contract.HasRole(&_BOGOWITickets.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _BOGOWITickets.Contract.IsApprovedForAll(&_BOGOWITickets.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _BOGOWITickets.Contract.IsApprovedForAll(&_BOGOWITickets.CallOpts, owner, operator)
}

// IsExpired is a free data retrieval call binding the contract method 0xd9548e53.
//
// Solidity: function isExpired(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) IsExpired(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "isExpired", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExpired is a free data retrieval call binding the contract method 0xd9548e53.
//
// Solidity: function isExpired(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) IsExpired(tokenId *big.Int) (bool, error) {
	return _BOGOWITickets.Contract.IsExpired(&_BOGOWITickets.CallOpts, tokenId)
}

// IsExpired is a free data retrieval call binding the contract method 0xd9548e53.
//
// Solidity: function isExpired(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) IsExpired(tokenId *big.Int) (bool, error) {
	return _BOGOWITickets.Contract.IsExpired(&_BOGOWITickets.CallOpts, tokenId)
}

// IsRedeemed is a free data retrieval call binding the contract method 0x32d33cd0.
//
// Solidity: function isRedeemed(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) IsRedeemed(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "isRedeemed", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRedeemed is a free data retrieval call binding the contract method 0x32d33cd0.
//
// Solidity: function isRedeemed(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) IsRedeemed(tokenId *big.Int) (bool, error) {
	return _BOGOWITickets.Contract.IsRedeemed(&_BOGOWITickets.CallOpts, tokenId)
}

// IsRedeemed is a free data retrieval call binding the contract method 0x32d33cd0.
//
// Solidity: function isRedeemed(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) IsRedeemed(tokenId *big.Int) (bool, error) {
	return _BOGOWITickets.Contract.IsRedeemed(&_BOGOWITickets.CallOpts, tokenId)
}

// IsTransferable is a free data retrieval call binding the contract method 0xb2564569.
//
// Solidity: function isTransferable(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) IsTransferable(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "isTransferable", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTransferable is a free data retrieval call binding the contract method 0xb2564569.
//
// Solidity: function isTransferable(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) IsTransferable(tokenId *big.Int) (bool, error) {
	return _BOGOWITickets.Contract.IsTransferable(&_BOGOWITickets.CallOpts, tokenId)
}

// IsTransferable is a free data retrieval call binding the contract method 0xb2564569.
//
// Solidity: function isTransferable(uint256 tokenId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) IsTransferable(tokenId *big.Int) (bool, error) {
	return _BOGOWITickets.Contract.IsTransferable(&_BOGOWITickets.CallOpts, tokenId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BOGOWITickets *BOGOWITicketsCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BOGOWITickets *BOGOWITicketsSession) Name() (string, error) {
	return _BOGOWITickets.Contract.Name(&_BOGOWITickets.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BOGOWITickets *BOGOWITicketsCallerSession) Name() (string, error) {
	return _BOGOWITickets.Contract.Name(&_BOGOWITickets.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BOGOWITickets *BOGOWITicketsCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BOGOWITickets *BOGOWITicketsSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _BOGOWITickets.Contract.OwnerOf(&_BOGOWITickets.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BOGOWITickets *BOGOWITicketsCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _BOGOWITickets.Contract.OwnerOf(&_BOGOWITickets.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) Paused() (bool, error) {
	return _BOGOWITickets.Contract.Paused(&_BOGOWITickets.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) Paused() (bool, error) {
	return _BOGOWITickets.Contract.Paused(&_BOGOWITickets.CallOpts)
}

// RoleManager is a free data retrieval call binding the contract method 0x00435da5.
//
// Solidity: function roleManager() view returns(address)
func (_BOGOWITickets *BOGOWITicketsCaller) RoleManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "roleManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RoleManager is a free data retrieval call binding the contract method 0x00435da5.
//
// Solidity: function roleManager() view returns(address)
func (_BOGOWITickets *BOGOWITicketsSession) RoleManager() (common.Address, error) {
	return _BOGOWITickets.Contract.RoleManager(&_BOGOWITickets.CallOpts)
}

// RoleManager is a free data retrieval call binding the contract method 0x00435da5.
//
// Solidity: function roleManager() view returns(address)
func (_BOGOWITickets *BOGOWITicketsCallerSession) RoleManager() (common.Address, error) {
	return _BOGOWITickets.Contract.RoleManager(&_BOGOWITickets.CallOpts)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 tokenId, uint256 salePrice) view returns(address receiver, uint256 amount)
func (_BOGOWITickets *BOGOWITicketsCaller) RoyaltyInfo(opts *bind.CallOpts, tokenId *big.Int, salePrice *big.Int) (struct {
	Receiver common.Address
	Amount   *big.Int
}, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "royaltyInfo", tokenId, salePrice)

	outstruct := new(struct {
		Receiver common.Address
		Amount   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receiver = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 tokenId, uint256 salePrice) view returns(address receiver, uint256 amount)
func (_BOGOWITickets *BOGOWITicketsSession) RoyaltyInfo(tokenId *big.Int, salePrice *big.Int) (struct {
	Receiver common.Address
	Amount   *big.Int
}, error) {
	return _BOGOWITickets.Contract.RoyaltyInfo(&_BOGOWITickets.CallOpts, tokenId, salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 tokenId, uint256 salePrice) view returns(address receiver, uint256 amount)
func (_BOGOWITickets *BOGOWITicketsCallerSession) RoyaltyInfo(tokenId *big.Int, salePrice *big.Int) (struct {
	Receiver common.Address
	Amount   *big.Int
}, error) {
	return _BOGOWITickets.Contract.RoyaltyInfo(&_BOGOWITickets.CallOpts, tokenId, salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BOGOWITickets.Contract.SupportsInterface(&_BOGOWITickets.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BOGOWITickets.Contract.SupportsInterface(&_BOGOWITickets.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BOGOWITickets *BOGOWITicketsCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BOGOWITickets *BOGOWITicketsSession) Symbol() (string, error) {
	return _BOGOWITickets.Contract.Symbol(&_BOGOWITickets.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BOGOWITickets *BOGOWITicketsCallerSession) Symbol() (string, error) {
	return _BOGOWITickets.Contract.Symbol(&_BOGOWITickets.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BOGOWITickets *BOGOWITicketsCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BOGOWITickets *BOGOWITicketsSession) TokenURI(tokenId *big.Int) (string, error) {
	return _BOGOWITickets.Contract.TokenURI(&_BOGOWITickets.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BOGOWITickets *BOGOWITicketsCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _BOGOWITickets.Contract.TokenURI(&_BOGOWITickets.CallOpts, tokenId)
}

// VerifyRedemptionSignature is a free data retrieval call binding the contract method 0x27f1fab2.
//
// Solidity: function verifyRedemptionSignature((uint256,address,uint256,uint256,uint256,bytes) redemptionData) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCaller) VerifyRedemptionSignature(opts *bind.CallOpts, redemptionData IBOGOWITicketsRedemptionData) (bool, error) {
	var out []interface{}
	err := _BOGOWITickets.contract.Call(opts, &out, "verifyRedemptionSignature", redemptionData)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyRedemptionSignature is a free data retrieval call binding the contract method 0x27f1fab2.
//
// Solidity: function verifyRedemptionSignature((uint256,address,uint256,uint256,uint256,bytes) redemptionData) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsSession) VerifyRedemptionSignature(redemptionData IBOGOWITicketsRedemptionData) (bool, error) {
	return _BOGOWITickets.Contract.VerifyRedemptionSignature(&_BOGOWITickets.CallOpts, redemptionData)
}

// VerifyRedemptionSignature is a free data retrieval call binding the contract method 0x27f1fab2.
//
// Solidity: function verifyRedemptionSignature((uint256,address,uint256,uint256,uint256,bytes) redemptionData) view returns(bool)
func (_BOGOWITickets *BOGOWITicketsCallerSession) VerifyRedemptionSignature(redemptionData IBOGOWITicketsRedemptionData) (bool, error) {
	return _BOGOWITickets.Contract.VerifyRedemptionSignature(&_BOGOWITickets.CallOpts, redemptionData)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Approve(&_BOGOWITickets.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Approve(&_BOGOWITickets.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Burn(&_BOGOWITickets.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Burn(&_BOGOWITickets.TransactOpts, tokenId)
}

// ExpireTicket is a paid mutator transaction binding the contract method 0xc90f131b.
//
// Solidity: function expireTicket(uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) ExpireTicket(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "expireTicket", tokenId)
}

// ExpireTicket is a paid mutator transaction binding the contract method 0xc90f131b.
//
// Solidity: function expireTicket(uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsSession) ExpireTicket(tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.ExpireTicket(&_BOGOWITickets.TransactOpts, tokenId)
}

// ExpireTicket is a paid mutator transaction binding the contract method 0xc90f131b.
//
// Solidity: function expireTicket(uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) ExpireTicket(tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.ExpireTicket(&_BOGOWITickets.TransactOpts, tokenId)
}

// MintBatch is a paid mutator transaction binding the contract method 0xadfee456.
//
// Solidity: function mintBatch((address,bytes32,bytes32,uint32,uint64,uint64,string,uint256)[] params) returns(uint256[])
func (_BOGOWITickets *BOGOWITicketsTransactor) MintBatch(opts *bind.TransactOpts, params []IBOGOWITicketsMintParams) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "mintBatch", params)
}

// MintBatch is a paid mutator transaction binding the contract method 0xadfee456.
//
// Solidity: function mintBatch((address,bytes32,bytes32,uint32,uint64,uint64,string,uint256)[] params) returns(uint256[])
func (_BOGOWITickets *BOGOWITicketsSession) MintBatch(params []IBOGOWITicketsMintParams) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.MintBatch(&_BOGOWITickets.TransactOpts, params)
}

// MintBatch is a paid mutator transaction binding the contract method 0xadfee456.
//
// Solidity: function mintBatch((address,bytes32,bytes32,uint32,uint64,uint64,string,uint256)[] params) returns(uint256[])
func (_BOGOWITickets *BOGOWITicketsTransactorSession) MintBatch(params []IBOGOWITicketsMintParams) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.MintBatch(&_BOGOWITickets.TransactOpts, params)
}

// MintTicket is a paid mutator transaction binding the contract method 0xbc975404.
//
// Solidity: function mintTicket((address,bytes32,bytes32,uint32,uint64,uint64,string,uint256) params) returns(uint256)
func (_BOGOWITickets *BOGOWITicketsTransactor) MintTicket(opts *bind.TransactOpts, params IBOGOWITicketsMintParams) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "mintTicket", params)
}

// MintTicket is a paid mutator transaction binding the contract method 0xbc975404.
//
// Solidity: function mintTicket((address,bytes32,bytes32,uint32,uint64,uint64,string,uint256) params) returns(uint256)
func (_BOGOWITickets *BOGOWITicketsSession) MintTicket(params IBOGOWITicketsMintParams) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.MintTicket(&_BOGOWITickets.TransactOpts, params)
}

// MintTicket is a paid mutator transaction binding the contract method 0xbc975404.
//
// Solidity: function mintTicket((address,bytes32,bytes32,uint32,uint64,uint64,string,uint256) params) returns(uint256)
func (_BOGOWITickets *BOGOWITicketsTransactorSession) MintTicket(params IBOGOWITicketsMintParams) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.MintTicket(&_BOGOWITickets.TransactOpts, params)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_BOGOWITickets *BOGOWITicketsSession) Pause() (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Pause(&_BOGOWITickets.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) Pause() (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Pause(&_BOGOWITickets.TransactOpts)
}

// RedeemTicket is a paid mutator transaction binding the contract method 0xa3eb1c40.
//
// Solidity: function redeemTicket((uint256,address,uint256,uint256,uint256,bytes) redemptionData) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) RedeemTicket(opts *bind.TransactOpts, redemptionData IBOGOWITicketsRedemptionData) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "redeemTicket", redemptionData)
}

// RedeemTicket is a paid mutator transaction binding the contract method 0xa3eb1c40.
//
// Solidity: function redeemTicket((uint256,address,uint256,uint256,uint256,bytes) redemptionData) returns()
func (_BOGOWITickets *BOGOWITicketsSession) RedeemTicket(redemptionData IBOGOWITicketsRedemptionData) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.RedeemTicket(&_BOGOWITickets.TransactOpts, redemptionData)
}

// RedeemTicket is a paid mutator transaction binding the contract method 0xa3eb1c40.
//
// Solidity: function redeemTicket((uint256,address,uint256,uint256,uint256,bytes) redemptionData) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) RedeemTicket(redemptionData IBOGOWITicketsRedemptionData) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.RedeemTicket(&_BOGOWITickets.TransactOpts, redemptionData)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SafeTransferFrom(&_BOGOWITickets.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SafeTransferFrom(&_BOGOWITickets.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BOGOWITickets *BOGOWITicketsSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SafeTransferFrom0(&_BOGOWITickets.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SafeTransferFrom0(&_BOGOWITickets.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BOGOWITickets *BOGOWITicketsSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SetApprovalForAll(&_BOGOWITickets.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SetApprovalForAll(&_BOGOWITickets.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "setBaseURI", newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_BOGOWITickets *BOGOWITicketsSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SetBaseURI(&_BOGOWITickets.TransactOpts, newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SetBaseURI(&_BOGOWITickets.TransactOpts, newBaseURI)
}

// SetRoyaltyInfo is a paid mutator transaction binding the contract method 0x02fa7c47.
//
// Solidity: function setRoyaltyInfo(address receiver, uint96 feeBasisPoints) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) SetRoyaltyInfo(opts *bind.TransactOpts, receiver common.Address, feeBasisPoints *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "setRoyaltyInfo", receiver, feeBasisPoints)
}

// SetRoyaltyInfo is a paid mutator transaction binding the contract method 0x02fa7c47.
//
// Solidity: function setRoyaltyInfo(address receiver, uint96 feeBasisPoints) returns()
func (_BOGOWITickets *BOGOWITicketsSession) SetRoyaltyInfo(receiver common.Address, feeBasisPoints *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SetRoyaltyInfo(&_BOGOWITickets.TransactOpts, receiver, feeBasisPoints)
}

// SetRoyaltyInfo is a paid mutator transaction binding the contract method 0x02fa7c47.
//
// Solidity: function setRoyaltyInfo(address receiver, uint96 feeBasisPoints) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) SetRoyaltyInfo(receiver common.Address, feeBasisPoints *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.SetRoyaltyInfo(&_BOGOWITickets.TransactOpts, receiver, feeBasisPoints)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.TransferFrom(&_BOGOWITickets.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.TransferFrom(&_BOGOWITickets.TransactOpts, from, to, tokenId)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_BOGOWITickets *BOGOWITicketsSession) Unpause() (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Unpause(&_BOGOWITickets.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) Unpause() (*types.Transaction, error) {
	return _BOGOWITickets.Contract.Unpause(&_BOGOWITickets.TransactOpts)
}

// UpdateTransferUnlock is a paid mutator transaction binding the contract method 0x4327bd8d.
//
// Solidity: function updateTransferUnlock(uint256 tokenId, uint64 newUnlockTime) returns()
func (_BOGOWITickets *BOGOWITicketsTransactor) UpdateTransferUnlock(opts *bind.TransactOpts, tokenId *big.Int, newUnlockTime uint64) (*types.Transaction, error) {
	return _BOGOWITickets.contract.Transact(opts, "updateTransferUnlock", tokenId, newUnlockTime)
}

// UpdateTransferUnlock is a paid mutator transaction binding the contract method 0x4327bd8d.
//
// Solidity: function updateTransferUnlock(uint256 tokenId, uint64 newUnlockTime) returns()
func (_BOGOWITickets *BOGOWITicketsSession) UpdateTransferUnlock(tokenId *big.Int, newUnlockTime uint64) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.UpdateTransferUnlock(&_BOGOWITickets.TransactOpts, tokenId, newUnlockTime)
}

// UpdateTransferUnlock is a paid mutator transaction binding the contract method 0x4327bd8d.
//
// Solidity: function updateTransferUnlock(uint256 tokenId, uint64 newUnlockTime) returns()
func (_BOGOWITickets *BOGOWITicketsTransactorSession) UpdateTransferUnlock(tokenId *big.Int, newUnlockTime uint64) (*types.Transaction, error) {
	return _BOGOWITickets.Contract.UpdateTransferUnlock(&_BOGOWITickets.TransactOpts, tokenId, newUnlockTime)
}

// BOGOWITicketsApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the BOGOWITickets contract.
type BOGOWITicketsApprovalIterator struct {
	Event *BOGOWITicketsApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsApproval represents a Approval event raised by the BOGOWITickets contract.
type BOGOWITicketsApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*BOGOWITicketsApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsApprovalIterator{contract: _BOGOWITickets.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsApproval)
				if err := _BOGOWITickets.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseApproval(log types.Log) (*BOGOWITicketsApproval, error) {
	event := new(BOGOWITicketsApproval)
	if err := _BOGOWITickets.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the BOGOWITickets contract.
type BOGOWITicketsApprovalForAllIterator struct {
	Event *BOGOWITicketsApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsApprovalForAll represents a ApprovalForAll event raised by the BOGOWITickets contract.
type BOGOWITicketsApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*BOGOWITicketsApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsApprovalForAllIterator{contract: _BOGOWITickets.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsApprovalForAll)
				if err := _BOGOWITickets.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseApprovalForAll(log types.Log) (*BOGOWITicketsApprovalForAll, error) {
	event := new(BOGOWITicketsApprovalForAll)
	if err := _BOGOWITickets.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsBaseURIUpdatedIterator is returned from FilterBaseURIUpdated and is used to iterate over the raw logs and unpacked data for BaseURIUpdated events raised by the BOGOWITickets contract.
type BOGOWITicketsBaseURIUpdatedIterator struct {
	Event *BOGOWITicketsBaseURIUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsBaseURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsBaseURIUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsBaseURIUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsBaseURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsBaseURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsBaseURIUpdated represents a BaseURIUpdated event raised by the BOGOWITickets contract.
type BOGOWITicketsBaseURIUpdated struct {
	NewBaseURI string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBaseURIUpdated is a free log retrieval operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterBaseURIUpdated(opts *bind.FilterOpts) (*BOGOWITicketsBaseURIUpdatedIterator, error) {

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsBaseURIUpdatedIterator{contract: _BOGOWITickets.contract, event: "BaseURIUpdated", logs: logs, sub: sub}, nil
}

// WatchBaseURIUpdated is a free log subscription operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchBaseURIUpdated(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsBaseURIUpdated) (event.Subscription, error) {

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsBaseURIUpdated)
				if err := _BOGOWITickets.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBaseURIUpdated is a log parse operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseBaseURIUpdated(log types.Log) (*BOGOWITicketsBaseURIUpdated, error) {
	event := new(BOGOWITicketsBaseURIUpdated)
	if err := _BOGOWITickets.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsBatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the BOGOWITickets contract.
type BOGOWITicketsBatchMetadataUpdateIterator struct {
	Event *BOGOWITicketsBatchMetadataUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsBatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsBatchMetadataUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsBatchMetadataUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsBatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsBatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsBatchMetadataUpdate represents a BatchMetadataUpdate event raised by the BOGOWITickets contract.
type BOGOWITicketsBatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*BOGOWITicketsBatchMetadataUpdateIterator, error) {

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsBatchMetadataUpdateIterator{contract: _BOGOWITickets.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsBatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsBatchMetadataUpdate)
				if err := _BOGOWITickets.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseBatchMetadataUpdate(log types.Log) (*BOGOWITicketsBatchMetadataUpdate, error) {
	event := new(BOGOWITicketsBatchMetadataUpdate)
	if err := _BOGOWITickets.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsBatchMintStartedIterator is returned from FilterBatchMintStarted and is used to iterate over the raw logs and unpacked data for BatchMintStarted events raised by the BOGOWITickets contract.
type BOGOWITicketsBatchMintStartedIterator struct {
	Event *BOGOWITicketsBatchMintStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsBatchMintStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsBatchMintStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsBatchMintStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsBatchMintStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsBatchMintStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsBatchMintStarted represents a BatchMintStarted event raised by the BOGOWITickets contract.
type BOGOWITicketsBatchMintStarted struct {
	BatchSize *big.Int
	Minter    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBatchMintStarted is a free log retrieval operation binding the contract event 0x80674971f5df274976a70fc833db39f6d22faf9e5b78396a26daa21e7531fcda.
//
// Solidity: event BatchMintStarted(uint256 batchSize, address indexed minter)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterBatchMintStarted(opts *bind.FilterOpts, minter []common.Address) (*BOGOWITicketsBatchMintStartedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "BatchMintStarted", minterRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsBatchMintStartedIterator{contract: _BOGOWITickets.contract, event: "BatchMintStarted", logs: logs, sub: sub}, nil
}

// WatchBatchMintStarted is a free log subscription operation binding the contract event 0x80674971f5df274976a70fc833db39f6d22faf9e5b78396a26daa21e7531fcda.
//
// Solidity: event BatchMintStarted(uint256 batchSize, address indexed minter)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchBatchMintStarted(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsBatchMintStarted, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "BatchMintStarted", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsBatchMintStarted)
				if err := _BOGOWITickets.contract.UnpackLog(event, "BatchMintStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatchMintStarted is a log parse operation binding the contract event 0x80674971f5df274976a70fc833db39f6d22faf9e5b78396a26daa21e7531fcda.
//
// Solidity: event BatchMintStarted(uint256 batchSize, address indexed minter)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseBatchMintStarted(log types.Log) (*BOGOWITicketsBatchMintStarted, error) {
	event := new(BOGOWITicketsBatchMintStarted)
	if err := _BOGOWITickets.contract.UnpackLog(event, "BatchMintStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsDatakyteMetadataLinkedIterator is returned from FilterDatakyteMetadataLinked and is used to iterate over the raw logs and unpacked data for DatakyteMetadataLinked events raised by the BOGOWITickets contract.
type BOGOWITicketsDatakyteMetadataLinkedIterator struct {
	Event *BOGOWITicketsDatakyteMetadataLinked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsDatakyteMetadataLinkedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsDatakyteMetadataLinked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsDatakyteMetadataLinked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsDatakyteMetadataLinkedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsDatakyteMetadataLinkedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsDatakyteMetadataLinked represents a DatakyteMetadataLinked event raised by the BOGOWITickets contract.
type BOGOWITicketsDatakyteMetadataLinked struct {
	TokenId       *big.Int
	DatakyteNftId string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDatakyteMetadataLinked is a free log retrieval operation binding the contract event 0xc114fb138057fc73fb3e245bbf32ed1412c8953965c0a664277ef2e8f47c6830.
//
// Solidity: event DatakyteMetadataLinked(uint256 indexed tokenId, string datakyteNftId)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterDatakyteMetadataLinked(opts *bind.FilterOpts, tokenId []*big.Int) (*BOGOWITicketsDatakyteMetadataLinkedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "DatakyteMetadataLinked", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsDatakyteMetadataLinkedIterator{contract: _BOGOWITickets.contract, event: "DatakyteMetadataLinked", logs: logs, sub: sub}, nil
}

// WatchDatakyteMetadataLinked is a free log subscription operation binding the contract event 0xc114fb138057fc73fb3e245bbf32ed1412c8953965c0a664277ef2e8f47c6830.
//
// Solidity: event DatakyteMetadataLinked(uint256 indexed tokenId, string datakyteNftId)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchDatakyteMetadataLinked(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsDatakyteMetadataLinked, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "DatakyteMetadataLinked", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsDatakyteMetadataLinked)
				if err := _BOGOWITickets.contract.UnpackLog(event, "DatakyteMetadataLinked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDatakyteMetadataLinked is a log parse operation binding the contract event 0xc114fb138057fc73fb3e245bbf32ed1412c8953965c0a664277ef2e8f47c6830.
//
// Solidity: event DatakyteMetadataLinked(uint256 indexed tokenId, string datakyteNftId)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseDatakyteMetadataLinked(log types.Log) (*BOGOWITicketsDatakyteMetadataLinked, error) {
	event := new(BOGOWITicketsDatakyteMetadataLinked)
	if err := _BOGOWITickets.contract.UnpackLog(event, "DatakyteMetadataLinked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the BOGOWITickets contract.
type BOGOWITicketsEIP712DomainChangedIterator struct {
	Event *BOGOWITicketsEIP712DomainChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsEIP712DomainChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsEIP712DomainChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsEIP712DomainChanged represents a EIP712DomainChanged event raised by the BOGOWITickets contract.
type BOGOWITicketsEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*BOGOWITicketsEIP712DomainChangedIterator, error) {

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsEIP712DomainChangedIterator{contract: _BOGOWITickets.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsEIP712DomainChanged)
				if err := _BOGOWITickets.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseEIP712DomainChanged(log types.Log) (*BOGOWITicketsEIP712DomainChanged, error) {
	event := new(BOGOWITicketsEIP712DomainChanged)
	if err := _BOGOWITickets.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the BOGOWITickets contract.
type BOGOWITicketsMetadataUpdateIterator struct {
	Event *BOGOWITicketsMetadataUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsMetadataUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsMetadataUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsMetadataUpdate represents a MetadataUpdate event raised by the BOGOWITickets contract.
type BOGOWITicketsMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*BOGOWITicketsMetadataUpdateIterator, error) {

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsMetadataUpdateIterator{contract: _BOGOWITickets.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsMetadataUpdate)
				if err := _BOGOWITickets.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseMetadataUpdate(log types.Log) (*BOGOWITicketsMetadataUpdate, error) {
	event := new(BOGOWITicketsMetadataUpdate)
	if err := _BOGOWITickets.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsNonceUsedIterator is returned from FilterNonceUsed and is used to iterate over the raw logs and unpacked data for NonceUsed events raised by the BOGOWITickets contract.
type BOGOWITicketsNonceUsedIterator struct {
	Event *BOGOWITicketsNonceUsed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsNonceUsedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsNonceUsed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsNonceUsed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsNonceUsedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsNonceUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsNonceUsed represents a NonceUsed event raised by the BOGOWITickets contract.
type BOGOWITicketsNonceUsed struct {
	Nonce *big.Int
	User  common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNonceUsed is a free log retrieval operation binding the contract event 0xd116bceaa03ee2cf1d4bb4c0cf26498cd3a2a20781cc967b85632c7c060cfe40.
//
// Solidity: event NonceUsed(uint256 indexed nonce, address indexed user)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterNonceUsed(opts *bind.FilterOpts, nonce []*big.Int, user []common.Address) (*BOGOWITicketsNonceUsedIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "NonceUsed", nonceRule, userRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsNonceUsedIterator{contract: _BOGOWITickets.contract, event: "NonceUsed", logs: logs, sub: sub}, nil
}

// WatchNonceUsed is a free log subscription operation binding the contract event 0xd116bceaa03ee2cf1d4bb4c0cf26498cd3a2a20781cc967b85632c7c060cfe40.
//
// Solidity: event NonceUsed(uint256 indexed nonce, address indexed user)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchNonceUsed(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsNonceUsed, nonce []*big.Int, user []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "NonceUsed", nonceRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsNonceUsed)
				if err := _BOGOWITickets.contract.UnpackLog(event, "NonceUsed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNonceUsed is a log parse operation binding the contract event 0xd116bceaa03ee2cf1d4bb4c0cf26498cd3a2a20781cc967b85632c7c060cfe40.
//
// Solidity: event NonceUsed(uint256 indexed nonce, address indexed user)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseNonceUsed(log types.Log) (*BOGOWITicketsNonceUsed, error) {
	event := new(BOGOWITicketsNonceUsed)
	if err := _BOGOWITickets.contract.UnpackLog(event, "NonceUsed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the BOGOWITickets contract.
type BOGOWITicketsPausedIterator struct {
	Event *BOGOWITicketsPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsPaused represents a Paused event raised by the BOGOWITickets contract.
type BOGOWITicketsPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterPaused(opts *bind.FilterOpts) (*BOGOWITicketsPausedIterator, error) {

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsPausedIterator{contract: _BOGOWITickets.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsPaused) (event.Subscription, error) {

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsPaused)
				if err := _BOGOWITickets.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParsePaused(log types.Log) (*BOGOWITicketsPaused, error) {
	event := new(BOGOWITicketsPaused)
	if err := _BOGOWITickets.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsRoleManagerSetIterator is returned from FilterRoleManagerSet and is used to iterate over the raw logs and unpacked data for RoleManagerSet events raised by the BOGOWITickets contract.
type BOGOWITicketsRoleManagerSetIterator struct {
	Event *BOGOWITicketsRoleManagerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsRoleManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsRoleManagerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsRoleManagerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsRoleManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsRoleManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsRoleManagerSet represents a RoleManagerSet event raised by the BOGOWITickets contract.
type BOGOWITicketsRoleManagerSet struct {
	RoleManagerAddress common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRoleManagerSet is a free log retrieval operation binding the contract event 0x765235f6b1f9df25a0fa901c365a8db93771de0abb8f48ffed12959c5c4d59b9.
//
// Solidity: event RoleManagerSet(address indexed roleManagerAddress)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterRoleManagerSet(opts *bind.FilterOpts, roleManagerAddress []common.Address) (*BOGOWITicketsRoleManagerSetIterator, error) {

	var roleManagerAddressRule []interface{}
	for _, roleManagerAddressItem := range roleManagerAddress {
		roleManagerAddressRule = append(roleManagerAddressRule, roleManagerAddressItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "RoleManagerSet", roleManagerAddressRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsRoleManagerSetIterator{contract: _BOGOWITickets.contract, event: "RoleManagerSet", logs: logs, sub: sub}, nil
}

// WatchRoleManagerSet is a free log subscription operation binding the contract event 0x765235f6b1f9df25a0fa901c365a8db93771de0abb8f48ffed12959c5c4d59b9.
//
// Solidity: event RoleManagerSet(address indexed roleManagerAddress)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchRoleManagerSet(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsRoleManagerSet, roleManagerAddress []common.Address) (event.Subscription, error) {

	var roleManagerAddressRule []interface{}
	for _, roleManagerAddressItem := range roleManagerAddress {
		roleManagerAddressRule = append(roleManagerAddressRule, roleManagerAddressItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "RoleManagerSet", roleManagerAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsRoleManagerSet)
				if err := _BOGOWITickets.contract.UnpackLog(event, "RoleManagerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleManagerSet is a log parse operation binding the contract event 0x765235f6b1f9df25a0fa901c365a8db93771de0abb8f48ffed12959c5c4d59b9.
//
// Solidity: event RoleManagerSet(address indexed roleManagerAddress)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseRoleManagerSet(log types.Log) (*BOGOWITicketsRoleManagerSet, error) {
	event := new(BOGOWITicketsRoleManagerSet)
	if err := _BOGOWITickets.contract.UnpackLog(event, "RoleManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsRoyaltyInfoUpdatedIterator is returned from FilterRoyaltyInfoUpdated and is used to iterate over the raw logs and unpacked data for RoyaltyInfoUpdated events raised by the BOGOWITickets contract.
type BOGOWITicketsRoyaltyInfoUpdatedIterator struct {
	Event *BOGOWITicketsRoyaltyInfoUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsRoyaltyInfoUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsRoyaltyInfoUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsRoyaltyInfoUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsRoyaltyInfoUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsRoyaltyInfoUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsRoyaltyInfoUpdated represents a RoyaltyInfoUpdated event raised by the BOGOWITickets contract.
type BOGOWITicketsRoyaltyInfoUpdated struct {
	Receiver       common.Address
	FeeBasisPoints *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRoyaltyInfoUpdated is a free log retrieval operation binding the contract event 0xae1d656a1268648b04ffa79c1416f05879338ae295aae3234d8db217356a1c62.
//
// Solidity: event RoyaltyInfoUpdated(address indexed receiver, uint96 feeBasisPoints)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterRoyaltyInfoUpdated(opts *bind.FilterOpts, receiver []common.Address) (*BOGOWITicketsRoyaltyInfoUpdatedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "RoyaltyInfoUpdated", receiverRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsRoyaltyInfoUpdatedIterator{contract: _BOGOWITickets.contract, event: "RoyaltyInfoUpdated", logs: logs, sub: sub}, nil
}

// WatchRoyaltyInfoUpdated is a free log subscription operation binding the contract event 0xae1d656a1268648b04ffa79c1416f05879338ae295aae3234d8db217356a1c62.
//
// Solidity: event RoyaltyInfoUpdated(address indexed receiver, uint96 feeBasisPoints)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchRoyaltyInfoUpdated(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsRoyaltyInfoUpdated, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "RoyaltyInfoUpdated", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsRoyaltyInfoUpdated)
				if err := _BOGOWITickets.contract.UnpackLog(event, "RoyaltyInfoUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoyaltyInfoUpdated is a log parse operation binding the contract event 0xae1d656a1268648b04ffa79c1416f05879338ae295aae3234d8db217356a1c62.
//
// Solidity: event RoyaltyInfoUpdated(address indexed receiver, uint96 feeBasisPoints)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseRoyaltyInfoUpdated(log types.Log) (*BOGOWITicketsRoyaltyInfoUpdated, error) {
	event := new(BOGOWITicketsRoyaltyInfoUpdated)
	if err := _BOGOWITickets.contract.UnpackLog(event, "RoyaltyInfoUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsTicketBurnedIterator is returned from FilterTicketBurned and is used to iterate over the raw logs and unpacked data for TicketBurned events raised by the BOGOWITickets contract.
type BOGOWITicketsTicketBurnedIterator struct {
	Event *BOGOWITicketsTicketBurned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsTicketBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsTicketBurned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsTicketBurned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsTicketBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsTicketBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsTicketBurned represents a TicketBurned event raised by the BOGOWITickets contract.
type BOGOWITicketsTicketBurned struct {
	TokenId *big.Int
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTicketBurned is a free log retrieval operation binding the contract event 0xf9017e896fd1920ef14b66376fe34ace88734084008a337b98e41568cf377b0b.
//
// Solidity: event TicketBurned(uint256 indexed tokenId, address indexed owner)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterTicketBurned(opts *bind.FilterOpts, tokenId []*big.Int, owner []common.Address) (*BOGOWITicketsTicketBurnedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "TicketBurned", tokenIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsTicketBurnedIterator{contract: _BOGOWITickets.contract, event: "TicketBurned", logs: logs, sub: sub}, nil
}

// WatchTicketBurned is a free log subscription operation binding the contract event 0xf9017e896fd1920ef14b66376fe34ace88734084008a337b98e41568cf377b0b.
//
// Solidity: event TicketBurned(uint256 indexed tokenId, address indexed owner)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchTicketBurned(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsTicketBurned, tokenId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "TicketBurned", tokenIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsTicketBurned)
				if err := _BOGOWITickets.contract.UnpackLog(event, "TicketBurned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTicketBurned is a log parse operation binding the contract event 0xf9017e896fd1920ef14b66376fe34ace88734084008a337b98e41568cf377b0b.
//
// Solidity: event TicketBurned(uint256 indexed tokenId, address indexed owner)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseTicketBurned(log types.Log) (*BOGOWITicketsTicketBurned, error) {
	event := new(BOGOWITicketsTicketBurned)
	if err := _BOGOWITickets.contract.UnpackLog(event, "TicketBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsTicketExpiredIterator is returned from FilterTicketExpired and is used to iterate over the raw logs and unpacked data for TicketExpired events raised by the BOGOWITickets contract.
type BOGOWITicketsTicketExpiredIterator struct {
	Event *BOGOWITicketsTicketExpired // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsTicketExpiredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsTicketExpired)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsTicketExpired)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsTicketExpiredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsTicketExpiredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsTicketExpired represents a TicketExpired event raised by the BOGOWITickets contract.
type BOGOWITicketsTicketExpired struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTicketExpired is a free log retrieval operation binding the contract event 0x876bd66db9cd47e8d01068b43a30a295d772a22173ebd1fdfa84f6d8b9e154a7.
//
// Solidity: event TicketExpired(uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterTicketExpired(opts *bind.FilterOpts, tokenId []*big.Int) (*BOGOWITicketsTicketExpiredIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "TicketExpired", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsTicketExpiredIterator{contract: _BOGOWITickets.contract, event: "TicketExpired", logs: logs, sub: sub}, nil
}

// WatchTicketExpired is a free log subscription operation binding the contract event 0x876bd66db9cd47e8d01068b43a30a295d772a22173ebd1fdfa84f6d8b9e154a7.
//
// Solidity: event TicketExpired(uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchTicketExpired(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsTicketExpired, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "TicketExpired", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsTicketExpired)
				if err := _BOGOWITickets.contract.UnpackLog(event, "TicketExpired", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTicketExpired is a log parse operation binding the contract event 0x876bd66db9cd47e8d01068b43a30a295d772a22173ebd1fdfa84f6d8b9e154a7.
//
// Solidity: event TicketExpired(uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseTicketExpired(log types.Log) (*BOGOWITicketsTicketExpired, error) {
	event := new(BOGOWITicketsTicketExpired)
	if err := _BOGOWITickets.contract.UnpackLog(event, "TicketExpired", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsTicketMintedIterator is returned from FilterTicketMinted and is used to iterate over the raw logs and unpacked data for TicketMinted events raised by the BOGOWITickets contract.
type BOGOWITicketsTicketMintedIterator struct {
	Event *BOGOWITicketsTicketMinted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsTicketMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsTicketMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsTicketMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsTicketMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsTicketMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsTicketMinted represents a TicketMinted event raised by the BOGOWITickets contract.
type BOGOWITicketsTicketMinted struct {
	TokenId           *big.Int
	BookingIdHash     [32]byte
	EventIdHash       [32]byte
	Buyer             common.Address
	RewardBasisPoints *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterTicketMinted is a free log retrieval operation binding the contract event 0xb07afce348c4f343297a2017d0123061f1c3eee2d5bfcd100f508b36789071a1.
//
// Solidity: event TicketMinted(uint256 indexed tokenId, bytes32 indexed bookingIdHash, bytes32 indexed eventIdHash, address buyer, uint256 rewardBasisPoints)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterTicketMinted(opts *bind.FilterOpts, tokenId []*big.Int, bookingIdHash [][32]byte, eventIdHash [][32]byte) (*BOGOWITicketsTicketMintedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var bookingIdHashRule []interface{}
	for _, bookingIdHashItem := range bookingIdHash {
		bookingIdHashRule = append(bookingIdHashRule, bookingIdHashItem)
	}
	var eventIdHashRule []interface{}
	for _, eventIdHashItem := range eventIdHash {
		eventIdHashRule = append(eventIdHashRule, eventIdHashItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "TicketMinted", tokenIdRule, bookingIdHashRule, eventIdHashRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsTicketMintedIterator{contract: _BOGOWITickets.contract, event: "TicketMinted", logs: logs, sub: sub}, nil
}

// WatchTicketMinted is a free log subscription operation binding the contract event 0xb07afce348c4f343297a2017d0123061f1c3eee2d5bfcd100f508b36789071a1.
//
// Solidity: event TicketMinted(uint256 indexed tokenId, bytes32 indexed bookingIdHash, bytes32 indexed eventIdHash, address buyer, uint256 rewardBasisPoints)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchTicketMinted(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsTicketMinted, tokenId []*big.Int, bookingIdHash [][32]byte, eventIdHash [][32]byte) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var bookingIdHashRule []interface{}
	for _, bookingIdHashItem := range bookingIdHash {
		bookingIdHashRule = append(bookingIdHashRule, bookingIdHashItem)
	}
	var eventIdHashRule []interface{}
	for _, eventIdHashItem := range eventIdHash {
		eventIdHashRule = append(eventIdHashRule, eventIdHashItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "TicketMinted", tokenIdRule, bookingIdHashRule, eventIdHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsTicketMinted)
				if err := _BOGOWITickets.contract.UnpackLog(event, "TicketMinted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTicketMinted is a log parse operation binding the contract event 0xb07afce348c4f343297a2017d0123061f1c3eee2d5bfcd100f508b36789071a1.
//
// Solidity: event TicketMinted(uint256 indexed tokenId, bytes32 indexed bookingIdHash, bytes32 indexed eventIdHash, address buyer, uint256 rewardBasisPoints)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseTicketMinted(log types.Log) (*BOGOWITicketsTicketMinted, error) {
	event := new(BOGOWITicketsTicketMinted)
	if err := _BOGOWITickets.contract.UnpackLog(event, "TicketMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsTicketRedeemedIterator is returned from FilterTicketRedeemed and is used to iterate over the raw logs and unpacked data for TicketRedeemed events raised by the BOGOWITickets contract.
type BOGOWITicketsTicketRedeemedIterator struct {
	Event *BOGOWITicketsTicketRedeemed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsTicketRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsTicketRedeemed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsTicketRedeemed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsTicketRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsTicketRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsTicketRedeemed represents a TicketRedeemed event raised by the BOGOWITickets contract.
type BOGOWITicketsTicketRedeemed struct {
	TokenId    *big.Int
	RedeemedBy common.Address
	Timestamp  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTicketRedeemed is a free log retrieval operation binding the contract event 0x577db72f216090dddc27ada8d4ad0fd6ab2c365a253821d43c50e4ba5b4f7184.
//
// Solidity: event TicketRedeemed(uint256 indexed tokenId, address indexed redeemedBy, uint256 timestamp)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterTicketRedeemed(opts *bind.FilterOpts, tokenId []*big.Int, redeemedBy []common.Address) (*BOGOWITicketsTicketRedeemedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var redeemedByRule []interface{}
	for _, redeemedByItem := range redeemedBy {
		redeemedByRule = append(redeemedByRule, redeemedByItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "TicketRedeemed", tokenIdRule, redeemedByRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsTicketRedeemedIterator{contract: _BOGOWITickets.contract, event: "TicketRedeemed", logs: logs, sub: sub}, nil
}

// WatchTicketRedeemed is a free log subscription operation binding the contract event 0x577db72f216090dddc27ada8d4ad0fd6ab2c365a253821d43c50e4ba5b4f7184.
//
// Solidity: event TicketRedeemed(uint256 indexed tokenId, address indexed redeemedBy, uint256 timestamp)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchTicketRedeemed(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsTicketRedeemed, tokenId []*big.Int, redeemedBy []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var redeemedByRule []interface{}
	for _, redeemedByItem := range redeemedBy {
		redeemedByRule = append(redeemedByRule, redeemedByItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "TicketRedeemed", tokenIdRule, redeemedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsTicketRedeemed)
				if err := _BOGOWITickets.contract.UnpackLog(event, "TicketRedeemed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTicketRedeemed is a log parse operation binding the contract event 0x577db72f216090dddc27ada8d4ad0fd6ab2c365a253821d43c50e4ba5b4f7184.
//
// Solidity: event TicketRedeemed(uint256 indexed tokenId, address indexed redeemedBy, uint256 timestamp)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseTicketRedeemed(log types.Log) (*BOGOWITicketsTicketRedeemed, error) {
	event := new(BOGOWITicketsTicketRedeemed)
	if err := _BOGOWITickets.contract.UnpackLog(event, "TicketRedeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BOGOWITickets contract.
type BOGOWITicketsTransferIterator struct {
	Event *BOGOWITicketsTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsTransfer represents a Transfer event raised by the BOGOWITickets contract.
type BOGOWITicketsTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*BOGOWITicketsTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsTransferIterator{contract: _BOGOWITickets.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsTransfer)
				if err := _BOGOWITickets.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseTransfer(log types.Log) (*BOGOWITicketsTransfer, error) {
	event := new(BOGOWITicketsTransfer)
	if err := _BOGOWITickets.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsTransferUnlockUpdatedIterator is returned from FilterTransferUnlockUpdated and is used to iterate over the raw logs and unpacked data for TransferUnlockUpdated events raised by the BOGOWITickets contract.
type BOGOWITicketsTransferUnlockUpdatedIterator struct {
	Event *BOGOWITicketsTransferUnlockUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsTransferUnlockUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsTransferUnlockUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsTransferUnlockUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsTransferUnlockUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsTransferUnlockUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsTransferUnlockUpdated represents a TransferUnlockUpdated event raised by the BOGOWITickets contract.
type BOGOWITicketsTransferUnlockUpdated struct {
	TokenId       *big.Int
	NewUnlockTime uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTransferUnlockUpdated is a free log retrieval operation binding the contract event 0xaf7074080ca22ecd9dd3dcbeb2e5106f25b5977d2ab0ce409b532693f7934c5d.
//
// Solidity: event TransferUnlockUpdated(uint256 indexed tokenId, uint64 newUnlockTime)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterTransferUnlockUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*BOGOWITicketsTransferUnlockUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "TransferUnlockUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsTransferUnlockUpdatedIterator{contract: _BOGOWITickets.contract, event: "TransferUnlockUpdated", logs: logs, sub: sub}, nil
}

// WatchTransferUnlockUpdated is a free log subscription operation binding the contract event 0xaf7074080ca22ecd9dd3dcbeb2e5106f25b5977d2ab0ce409b532693f7934c5d.
//
// Solidity: event TransferUnlockUpdated(uint256 indexed tokenId, uint64 newUnlockTime)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchTransferUnlockUpdated(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsTransferUnlockUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "TransferUnlockUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsTransferUnlockUpdated)
				if err := _BOGOWITickets.contract.UnpackLog(event, "TransferUnlockUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferUnlockUpdated is a log parse operation binding the contract event 0xaf7074080ca22ecd9dd3dcbeb2e5106f25b5977d2ab0ce409b532693f7934c5d.
//
// Solidity: event TransferUnlockUpdated(uint256 indexed tokenId, uint64 newUnlockTime)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseTransferUnlockUpdated(log types.Log) (*BOGOWITicketsTransferUnlockUpdated, error) {
	event := new(BOGOWITicketsTransferUnlockUpdated)
	if err := _BOGOWITickets.contract.UnpackLog(event, "TransferUnlockUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BOGOWITicketsUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the BOGOWITickets contract.
type BOGOWITicketsUnpausedIterator struct {
	Event *BOGOWITicketsUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BOGOWITicketsUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BOGOWITicketsUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BOGOWITicketsUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BOGOWITicketsUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BOGOWITicketsUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BOGOWITicketsUnpaused represents a Unpaused event raised by the BOGOWITickets contract.
type BOGOWITicketsUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_BOGOWITickets *BOGOWITicketsFilterer) FilterUnpaused(opts *bind.FilterOpts) (*BOGOWITicketsUnpausedIterator, error) {

	logs, sub, err := _BOGOWITickets.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &BOGOWITicketsUnpausedIterator{contract: _BOGOWITickets.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_BOGOWITickets *BOGOWITicketsFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *BOGOWITicketsUnpaused) (event.Subscription, error) {

	logs, sub, err := _BOGOWITickets.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BOGOWITicketsUnpaused)
				if err := _BOGOWITickets.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_BOGOWITickets *BOGOWITicketsFilterer) ParseUnpaused(log types.Log) (*BOGOWITicketsUnpaused, error) {
	event := new(BOGOWITicketsUnpaused)
	if err := _BOGOWITickets.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
