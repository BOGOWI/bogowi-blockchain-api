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

// RoleManagerMetaData contains all meta data concerning the RoleManager contract.
var RoleManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"ContractDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"ContractRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGrantedGlobally\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevokedGlobally\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUSINESS_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DAO_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DISTRIBUTOR_BACKEND_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TREASURY_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"}],\"name\":\"batchGrantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"}],\"name\":\"batchRevokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"checkRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"contractNames\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"deregisterContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegisteredContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"names\",\"type\":\"string[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"registerContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"registeredContracts\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RoleManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use RoleManagerMetaData.ABI instead.
var RoleManagerABI = RoleManagerMetaData.ABI

// RoleManager is an auto generated Go binding around an Ethereum contract.
type RoleManager struct {
	RoleManagerCaller     // Read-only binding to the contract
	RoleManagerTransactor // Write-only binding to the contract
	RoleManagerFilterer   // Log filterer for contract events
}

// RoleManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RoleManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RoleManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RoleManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoleManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RoleManagerSession struct {
	Contract     *RoleManager      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RoleManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RoleManagerCallerSession struct {
	Contract *RoleManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// RoleManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RoleManagerTransactorSession struct {
	Contract     *RoleManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// RoleManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RoleManagerRaw struct {
	Contract *RoleManager // Generic contract binding to access the raw methods on
}

// RoleManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RoleManagerCallerRaw struct {
	Contract *RoleManagerCaller // Generic read-only contract binding to access the raw methods on
}

// RoleManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RoleManagerTransactorRaw struct {
	Contract *RoleManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRoleManager creates a new instance of RoleManager, bound to a specific deployed contract.
func NewRoleManager(address common.Address, backend bind.ContractBackend) (*RoleManager, error) {
	contract, err := bindRoleManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RoleManager{RoleManagerCaller: RoleManagerCaller{contract: contract}, RoleManagerTransactor: RoleManagerTransactor{contract: contract}, RoleManagerFilterer: RoleManagerFilterer{contract: contract}}, nil
}

// NewRoleManagerCaller creates a new read-only instance of RoleManager, bound to a specific deployed contract.
func NewRoleManagerCaller(address common.Address, caller bind.ContractCaller) (*RoleManagerCaller, error) {
	contract, err := bindRoleManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RoleManagerCaller{contract: contract}, nil
}

// NewRoleManagerTransactor creates a new write-only instance of RoleManager, bound to a specific deployed contract.
func NewRoleManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RoleManagerTransactor, error) {
	contract, err := bindRoleManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RoleManagerTransactor{contract: contract}, nil
}

// NewRoleManagerFilterer creates a new log filterer instance of RoleManager, bound to a specific deployed contract.
func NewRoleManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RoleManagerFilterer, error) {
	contract, err := bindRoleManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RoleManagerFilterer{contract: contract}, nil
}

// bindRoleManager binds a generic wrapper to an already deployed contract.
func bindRoleManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RoleManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RoleManager *RoleManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RoleManager.Contract.RoleManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RoleManager *RoleManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoleManager.Contract.RoleManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RoleManager *RoleManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RoleManager.Contract.RoleManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RoleManager *RoleManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RoleManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RoleManager *RoleManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoleManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RoleManager *RoleManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RoleManager.Contract.contract.Transact(opts, method, params...)
}

// BUSINESSROLE is a free data retrieval call binding the contract method 0x3207247b.
//
// Solidity: function BUSINESS_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCaller) BUSINESSROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "BUSINESS_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BUSINESSROLE is a free data retrieval call binding the contract method 0x3207247b.
//
// Solidity: function BUSINESS_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerSession) BUSINESSROLE() ([32]byte, error) {
	return _RoleManager.Contract.BUSINESSROLE(&_RoleManager.CallOpts)
}

// BUSINESSROLE is a free data retrieval call binding the contract method 0x3207247b.
//
// Solidity: function BUSINESS_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) BUSINESSROLE() ([32]byte, error) {
	return _RoleManager.Contract.BUSINESSROLE(&_RoleManager.CallOpts)
}

// DAOROLE is a free data retrieval call binding the contract method 0xe9c26518.
//
// Solidity: function DAO_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCaller) DAOROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "DAO_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DAOROLE is a free data retrieval call binding the contract method 0xe9c26518.
//
// Solidity: function DAO_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerSession) DAOROLE() ([32]byte, error) {
	return _RoleManager.Contract.DAOROLE(&_RoleManager.CallOpts)
}

// DAOROLE is a free data retrieval call binding the contract method 0xe9c26518.
//
// Solidity: function DAO_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) DAOROLE() ([32]byte, error) {
	return _RoleManager.Contract.DAOROLE(&_RoleManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RoleManager.Contract.DEFAULTADMINROLE(&_RoleManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RoleManager.Contract.DEFAULTADMINROLE(&_RoleManager.CallOpts)
}

// DISTRIBUTORBACKENDROLE is a free data retrieval call binding the contract method 0xf68803c1.
//
// Solidity: function DISTRIBUTOR_BACKEND_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCaller) DISTRIBUTORBACKENDROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "DISTRIBUTOR_BACKEND_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DISTRIBUTORBACKENDROLE is a free data retrieval call binding the contract method 0xf68803c1.
//
// Solidity: function DISTRIBUTOR_BACKEND_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerSession) DISTRIBUTORBACKENDROLE() ([32]byte, error) {
	return _RoleManager.Contract.DISTRIBUTORBACKENDROLE(&_RoleManager.CallOpts)
}

// DISTRIBUTORBACKENDROLE is a free data retrieval call binding the contract method 0xf68803c1.
//
// Solidity: function DISTRIBUTOR_BACKEND_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) DISTRIBUTORBACKENDROLE() ([32]byte, error) {
	return _RoleManager.Contract.DISTRIBUTORBACKENDROLE(&_RoleManager.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerSession) MINTERROLE() ([32]byte, error) {
	return _RoleManager.Contract.MINTERROLE(&_RoleManager.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) MINTERROLE() ([32]byte, error) {
	return _RoleManager.Contract.MINTERROLE(&_RoleManager.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerSession) PAUSERROLE() ([32]byte, error) {
	return _RoleManager.Contract.PAUSERROLE(&_RoleManager.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) PAUSERROLE() ([32]byte, error) {
	return _RoleManager.Contract.PAUSERROLE(&_RoleManager.CallOpts)
}

// TREASURYROLE is a free data retrieval call binding the contract method 0xd11a57ec.
//
// Solidity: function TREASURY_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCaller) TREASURYROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "TREASURY_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TREASURYROLE is a free data retrieval call binding the contract method 0xd11a57ec.
//
// Solidity: function TREASURY_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerSession) TREASURYROLE() ([32]byte, error) {
	return _RoleManager.Contract.TREASURYROLE(&_RoleManager.CallOpts)
}

// TREASURYROLE is a free data retrieval call binding the contract method 0xd11a57ec.
//
// Solidity: function TREASURY_ROLE() view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) TREASURYROLE() ([32]byte, error) {
	return _RoleManager.Contract.TREASURYROLE(&_RoleManager.CallOpts)
}

// CheckRole is a free data retrieval call binding the contract method 0x12d9a6ad.
//
// Solidity: function checkRole(bytes32 role, address account) view returns(bool)
func (_RoleManager *RoleManagerCaller) CheckRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "checkRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckRole is a free data retrieval call binding the contract method 0x12d9a6ad.
//
// Solidity: function checkRole(bytes32 role, address account) view returns(bool)
func (_RoleManager *RoleManagerSession) CheckRole(role [32]byte, account common.Address) (bool, error) {
	return _RoleManager.Contract.CheckRole(&_RoleManager.CallOpts, role, account)
}

// CheckRole is a free data retrieval call binding the contract method 0x12d9a6ad.
//
// Solidity: function checkRole(bytes32 role, address account) view returns(bool)
func (_RoleManager *RoleManagerCallerSession) CheckRole(role [32]byte, account common.Address) (bool, error) {
	return _RoleManager.Contract.CheckRole(&_RoleManager.CallOpts, role, account)
}

// ContractNames is a free data retrieval call binding the contract method 0x1f80b872.
//
// Solidity: function contractNames(address ) view returns(string)
func (_RoleManager *RoleManagerCaller) ContractNames(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "contractNames", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractNames is a free data retrieval call binding the contract method 0x1f80b872.
//
// Solidity: function contractNames(address ) view returns(string)
func (_RoleManager *RoleManagerSession) ContractNames(arg0 common.Address) (string, error) {
	return _RoleManager.Contract.ContractNames(&_RoleManager.CallOpts, arg0)
}

// ContractNames is a free data retrieval call binding the contract method 0x1f80b872.
//
// Solidity: function contractNames(address ) view returns(string)
func (_RoleManager *RoleManagerCallerSession) ContractNames(arg0 common.Address) (string, error) {
	return _RoleManager.Contract.ContractNames(&_RoleManager.CallOpts, arg0)
}

// GetRegisteredContracts is a free data retrieval call binding the contract method 0x95877f32.
//
// Solidity: function getRegisteredContracts() pure returns(address[] addresses, string[] names)
func (_RoleManager *RoleManagerCaller) GetRegisteredContracts(opts *bind.CallOpts) (struct {
	Addresses []common.Address
	Names     []string
}, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "getRegisteredContracts")

	outstruct := new(struct {
		Addresses []common.Address
		Names     []string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Addresses = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Names = *abi.ConvertType(out[1], new([]string)).(*[]string)

	return *outstruct, err

}

// GetRegisteredContracts is a free data retrieval call binding the contract method 0x95877f32.
//
// Solidity: function getRegisteredContracts() pure returns(address[] addresses, string[] names)
func (_RoleManager *RoleManagerSession) GetRegisteredContracts() (struct {
	Addresses []common.Address
	Names     []string
}, error) {
	return _RoleManager.Contract.GetRegisteredContracts(&_RoleManager.CallOpts)
}

// GetRegisteredContracts is a free data retrieval call binding the contract method 0x95877f32.
//
// Solidity: function getRegisteredContracts() pure returns(address[] addresses, string[] names)
func (_RoleManager *RoleManagerCallerSession) GetRegisteredContracts() (struct {
	Addresses []common.Address
	Names     []string
}, error) {
	return _RoleManager.Contract.GetRegisteredContracts(&_RoleManager.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RoleManager *RoleManagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RoleManager *RoleManagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RoleManager.Contract.GetRoleAdmin(&_RoleManager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RoleManager *RoleManagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RoleManager.Contract.GetRoleAdmin(&_RoleManager.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RoleManager *RoleManagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RoleManager *RoleManagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RoleManager.Contract.HasRole(&_RoleManager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RoleManager *RoleManagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RoleManager.Contract.HasRole(&_RoleManager.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RoleManager *RoleManagerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RoleManager *RoleManagerSession) Paused() (bool, error) {
	return _RoleManager.Contract.Paused(&_RoleManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RoleManager *RoleManagerCallerSession) Paused() (bool, error) {
	return _RoleManager.Contract.Paused(&_RoleManager.CallOpts)
}

// RegisteredContracts is a free data retrieval call binding the contract method 0xa06617cd.
//
// Solidity: function registeredContracts(address ) view returns(bool)
func (_RoleManager *RoleManagerCaller) RegisteredContracts(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "registeredContracts", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RegisteredContracts is a free data retrieval call binding the contract method 0xa06617cd.
//
// Solidity: function registeredContracts(address ) view returns(bool)
func (_RoleManager *RoleManagerSession) RegisteredContracts(arg0 common.Address) (bool, error) {
	return _RoleManager.Contract.RegisteredContracts(&_RoleManager.CallOpts, arg0)
}

// RegisteredContracts is a free data retrieval call binding the contract method 0xa06617cd.
//
// Solidity: function registeredContracts(address ) view returns(bool)
func (_RoleManager *RoleManagerCallerSession) RegisteredContracts(arg0 common.Address) (bool, error) {
	return _RoleManager.Contract.RegisteredContracts(&_RoleManager.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RoleManager *RoleManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RoleManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RoleManager *RoleManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RoleManager.Contract.SupportsInterface(&_RoleManager.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RoleManager *RoleManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RoleManager.Contract.SupportsInterface(&_RoleManager.CallOpts, interfaceId)
}

// BatchGrantRole is a paid mutator transaction binding the contract method 0x46b5cb59.
//
// Solidity: function batchGrantRole(bytes32 role, address[] accounts) returns()
func (_RoleManager *RoleManagerTransactor) BatchGrantRole(opts *bind.TransactOpts, role [32]byte, accounts []common.Address) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "batchGrantRole", role, accounts)
}

// BatchGrantRole is a paid mutator transaction binding the contract method 0x46b5cb59.
//
// Solidity: function batchGrantRole(bytes32 role, address[] accounts) returns()
func (_RoleManager *RoleManagerSession) BatchGrantRole(role [32]byte, accounts []common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.BatchGrantRole(&_RoleManager.TransactOpts, role, accounts)
}

// BatchGrantRole is a paid mutator transaction binding the contract method 0x46b5cb59.
//
// Solidity: function batchGrantRole(bytes32 role, address[] accounts) returns()
func (_RoleManager *RoleManagerTransactorSession) BatchGrantRole(role [32]byte, accounts []common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.BatchGrantRole(&_RoleManager.TransactOpts, role, accounts)
}

// BatchRevokeRole is a paid mutator transaction binding the contract method 0x6f39feec.
//
// Solidity: function batchRevokeRole(bytes32 role, address[] accounts) returns()
func (_RoleManager *RoleManagerTransactor) BatchRevokeRole(opts *bind.TransactOpts, role [32]byte, accounts []common.Address) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "batchRevokeRole", role, accounts)
}

// BatchRevokeRole is a paid mutator transaction binding the contract method 0x6f39feec.
//
// Solidity: function batchRevokeRole(bytes32 role, address[] accounts) returns()
func (_RoleManager *RoleManagerSession) BatchRevokeRole(role [32]byte, accounts []common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.BatchRevokeRole(&_RoleManager.TransactOpts, role, accounts)
}

// BatchRevokeRole is a paid mutator transaction binding the contract method 0x6f39feec.
//
// Solidity: function batchRevokeRole(bytes32 role, address[] accounts) returns()
func (_RoleManager *RoleManagerTransactorSession) BatchRevokeRole(role [32]byte, accounts []common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.BatchRevokeRole(&_RoleManager.TransactOpts, role, accounts)
}

// DeregisterContract is a paid mutator transaction binding the contract method 0x64b626d8.
//
// Solidity: function deregisterContract(address contractAddress) returns()
func (_RoleManager *RoleManagerTransactor) DeregisterContract(opts *bind.TransactOpts, contractAddress common.Address) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "deregisterContract", contractAddress)
}

// DeregisterContract is a paid mutator transaction binding the contract method 0x64b626d8.
//
// Solidity: function deregisterContract(address contractAddress) returns()
func (_RoleManager *RoleManagerSession) DeregisterContract(contractAddress common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.DeregisterContract(&_RoleManager.TransactOpts, contractAddress)
}

// DeregisterContract is a paid mutator transaction binding the contract method 0x64b626d8.
//
// Solidity: function deregisterContract(address contractAddress) returns()
func (_RoleManager *RoleManagerTransactorSession) DeregisterContract(contractAddress common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.DeregisterContract(&_RoleManager.TransactOpts, contractAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RoleManager *RoleManagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RoleManager *RoleManagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.GrantRole(&_RoleManager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RoleManager *RoleManagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.GrantRole(&_RoleManager.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RoleManager *RoleManagerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RoleManager *RoleManagerSession) Pause() (*types.Transaction, error) {
	return _RoleManager.Contract.Pause(&_RoleManager.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RoleManager *RoleManagerTransactorSession) Pause() (*types.Transaction, error) {
	return _RoleManager.Contract.Pause(&_RoleManager.TransactOpts)
}

// RegisterContract is a paid mutator transaction binding the contract method 0xda23c0d9.
//
// Solidity: function registerContract(address contractAddress, string name) returns()
func (_RoleManager *RoleManagerTransactor) RegisterContract(opts *bind.TransactOpts, contractAddress common.Address, name string) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "registerContract", contractAddress, name)
}

// RegisterContract is a paid mutator transaction binding the contract method 0xda23c0d9.
//
// Solidity: function registerContract(address contractAddress, string name) returns()
func (_RoleManager *RoleManagerSession) RegisterContract(contractAddress common.Address, name string) (*types.Transaction, error) {
	return _RoleManager.Contract.RegisterContract(&_RoleManager.TransactOpts, contractAddress, name)
}

// RegisterContract is a paid mutator transaction binding the contract method 0xda23c0d9.
//
// Solidity: function registerContract(address contractAddress, string name) returns()
func (_RoleManager *RoleManagerTransactorSession) RegisterContract(contractAddress common.Address, name string) (*types.Transaction, error) {
	return _RoleManager.Contract.RegisterContract(&_RoleManager.TransactOpts, contractAddress, name)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RoleManager *RoleManagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RoleManager *RoleManagerSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.RenounceRole(&_RoleManager.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RoleManager *RoleManagerTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.RenounceRole(&_RoleManager.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RoleManager *RoleManagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RoleManager *RoleManagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.RevokeRole(&_RoleManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RoleManager *RoleManagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.RevokeRole(&_RoleManager.TransactOpts, role, account)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_RoleManager *RoleManagerTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_RoleManager *RoleManagerSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.TransferAdmin(&_RoleManager.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_RoleManager *RoleManagerTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _RoleManager.Contract.TransferAdmin(&_RoleManager.TransactOpts, newAdmin)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RoleManager *RoleManagerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoleManager.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RoleManager *RoleManagerSession) Unpause() (*types.Transaction, error) {
	return _RoleManager.Contract.Unpause(&_RoleManager.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RoleManager *RoleManagerTransactorSession) Unpause() (*types.Transaction, error) {
	return _RoleManager.Contract.Unpause(&_RoleManager.TransactOpts)
}

// RoleManagerContractDeregisteredIterator is returned from FilterContractDeregistered and is used to iterate over the raw logs and unpacked data for ContractDeregistered events raised by the RoleManager contract.
type RoleManagerContractDeregisteredIterator struct {
	Event *RoleManagerContractDeregistered // Event containing the contract specifics and raw log

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
func (it *RoleManagerContractDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerContractDeregistered)
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
		it.Event = new(RoleManagerContractDeregistered)
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
func (it *RoleManagerContractDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerContractDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerContractDeregistered represents a ContractDeregistered event raised by the RoleManager contract.
type RoleManagerContractDeregistered struct {
	ContractAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterContractDeregistered is a free log retrieval operation binding the contract event 0xc32d83d423bd05b8916851ed8856f32f6a435e4cd9224f3972a943075dd7727e.
//
// Solidity: event ContractDeregistered(address indexed contractAddress)
func (_RoleManager *RoleManagerFilterer) FilterContractDeregistered(opts *bind.FilterOpts, contractAddress []common.Address) (*RoleManagerContractDeregisteredIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "ContractDeregistered", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &RoleManagerContractDeregisteredIterator{contract: _RoleManager.contract, event: "ContractDeregistered", logs: logs, sub: sub}, nil
}

// WatchContractDeregistered is a free log subscription operation binding the contract event 0xc32d83d423bd05b8916851ed8856f32f6a435e4cd9224f3972a943075dd7727e.
//
// Solidity: event ContractDeregistered(address indexed contractAddress)
func (_RoleManager *RoleManagerFilterer) WatchContractDeregistered(opts *bind.WatchOpts, sink chan<- *RoleManagerContractDeregistered, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "ContractDeregistered", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerContractDeregistered)
				if err := _RoleManager.contract.UnpackLog(event, "ContractDeregistered", log); err != nil {
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

// ParseContractDeregistered is a log parse operation binding the contract event 0xc32d83d423bd05b8916851ed8856f32f6a435e4cd9224f3972a943075dd7727e.
//
// Solidity: event ContractDeregistered(address indexed contractAddress)
func (_RoleManager *RoleManagerFilterer) ParseContractDeregistered(log types.Log) (*RoleManagerContractDeregistered, error) {
	event := new(RoleManagerContractDeregistered)
	if err := _RoleManager.contract.UnpackLog(event, "ContractDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerContractRegisteredIterator is returned from FilterContractRegistered and is used to iterate over the raw logs and unpacked data for ContractRegistered events raised by the RoleManager contract.
type RoleManagerContractRegisteredIterator struct {
	Event *RoleManagerContractRegistered // Event containing the contract specifics and raw log

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
func (it *RoleManagerContractRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerContractRegistered)
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
		it.Event = new(RoleManagerContractRegistered)
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
func (it *RoleManagerContractRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerContractRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerContractRegistered represents a ContractRegistered event raised by the RoleManager contract.
type RoleManagerContractRegistered struct {
	ContractAddress common.Address
	Name            string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterContractRegistered is a free log retrieval operation binding the contract event 0x02331e3184c8c427ccd13836e43a51af710810d252021bce3e529d9a6d94cd98.
//
// Solidity: event ContractRegistered(address indexed contractAddress, string name)
func (_RoleManager *RoleManagerFilterer) FilterContractRegistered(opts *bind.FilterOpts, contractAddress []common.Address) (*RoleManagerContractRegisteredIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "ContractRegistered", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &RoleManagerContractRegisteredIterator{contract: _RoleManager.contract, event: "ContractRegistered", logs: logs, sub: sub}, nil
}

// WatchContractRegistered is a free log subscription operation binding the contract event 0x02331e3184c8c427ccd13836e43a51af710810d252021bce3e529d9a6d94cd98.
//
// Solidity: event ContractRegistered(address indexed contractAddress, string name)
func (_RoleManager *RoleManagerFilterer) WatchContractRegistered(opts *bind.WatchOpts, sink chan<- *RoleManagerContractRegistered, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "ContractRegistered", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerContractRegistered)
				if err := _RoleManager.contract.UnpackLog(event, "ContractRegistered", log); err != nil {
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

// ParseContractRegistered is a log parse operation binding the contract event 0x02331e3184c8c427ccd13836e43a51af710810d252021bce3e529d9a6d94cd98.
//
// Solidity: event ContractRegistered(address indexed contractAddress, string name)
func (_RoleManager *RoleManagerFilterer) ParseContractRegistered(log types.Log) (*RoleManagerContractRegistered, error) {
	event := new(RoleManagerContractRegistered)
	if err := _RoleManager.contract.UnpackLog(event, "ContractRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RoleManager contract.
type RoleManagerPausedIterator struct {
	Event *RoleManagerPaused // Event containing the contract specifics and raw log

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
func (it *RoleManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerPaused)
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
		it.Event = new(RoleManagerPaused)
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
func (it *RoleManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerPaused represents a Paused event raised by the RoleManager contract.
type RoleManagerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RoleManager *RoleManagerFilterer) FilterPaused(opts *bind.FilterOpts) (*RoleManagerPausedIterator, error) {

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RoleManagerPausedIterator{contract: _RoleManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RoleManager *RoleManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RoleManagerPaused) (event.Subscription, error) {

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerPaused)
				if err := _RoleManager.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_RoleManager *RoleManagerFilterer) ParsePaused(log types.Log) (*RoleManagerPaused, error) {
	event := new(RoleManagerPaused)
	if err := _RoleManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RoleManager contract.
type RoleManagerRoleAdminChangedIterator struct {
	Event *RoleManagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RoleManagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerRoleAdminChanged)
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
		it.Event = new(RoleManagerRoleAdminChanged)
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
func (it *RoleManagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerRoleAdminChanged represents a RoleAdminChanged event raised by the RoleManager contract.
type RoleManagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RoleManager *RoleManagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RoleManagerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RoleManagerRoleAdminChangedIterator{contract: _RoleManager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RoleManager *RoleManagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RoleManagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerRoleAdminChanged)
				if err := _RoleManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RoleManager *RoleManagerFilterer) ParseRoleAdminChanged(log types.Log) (*RoleManagerRoleAdminChanged, error) {
	event := new(RoleManagerRoleAdminChanged)
	if err := _RoleManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RoleManager contract.
type RoleManagerRoleGrantedIterator struct {
	Event *RoleManagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *RoleManagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerRoleGranted)
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
		it.Event = new(RoleManagerRoleGranted)
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
func (it *RoleManagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerRoleGranted represents a RoleGranted event raised by the RoleManager contract.
type RoleManagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RoleManagerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RoleManagerRoleGrantedIterator{contract: _RoleManager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RoleManagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerRoleGranted)
				if err := _RoleManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) ParseRoleGranted(log types.Log) (*RoleManagerRoleGranted, error) {
	event := new(RoleManagerRoleGranted)
	if err := _RoleManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerRoleGrantedGloballyIterator is returned from FilterRoleGrantedGlobally and is used to iterate over the raw logs and unpacked data for RoleGrantedGlobally events raised by the RoleManager contract.
type RoleManagerRoleGrantedGloballyIterator struct {
	Event *RoleManagerRoleGrantedGlobally // Event containing the contract specifics and raw log

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
func (it *RoleManagerRoleGrantedGloballyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerRoleGrantedGlobally)
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
		it.Event = new(RoleManagerRoleGrantedGlobally)
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
func (it *RoleManagerRoleGrantedGloballyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerRoleGrantedGloballyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerRoleGrantedGlobally represents a RoleGrantedGlobally event raised by the RoleManager contract.
type RoleManagerRoleGrantedGlobally struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGrantedGlobally is a free log retrieval operation binding the contract event 0x287e752593aea8255aeaa7ea93d3e06807273ef2426cde02d5572ee69f16e969.
//
// Solidity: event RoleGrantedGlobally(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) FilterRoleGrantedGlobally(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RoleManagerRoleGrantedGloballyIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "RoleGrantedGlobally", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RoleManagerRoleGrantedGloballyIterator{contract: _RoleManager.contract, event: "RoleGrantedGlobally", logs: logs, sub: sub}, nil
}

// WatchRoleGrantedGlobally is a free log subscription operation binding the contract event 0x287e752593aea8255aeaa7ea93d3e06807273ef2426cde02d5572ee69f16e969.
//
// Solidity: event RoleGrantedGlobally(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) WatchRoleGrantedGlobally(opts *bind.WatchOpts, sink chan<- *RoleManagerRoleGrantedGlobally, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "RoleGrantedGlobally", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerRoleGrantedGlobally)
				if err := _RoleManager.contract.UnpackLog(event, "RoleGrantedGlobally", log); err != nil {
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

// ParseRoleGrantedGlobally is a log parse operation binding the contract event 0x287e752593aea8255aeaa7ea93d3e06807273ef2426cde02d5572ee69f16e969.
//
// Solidity: event RoleGrantedGlobally(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) ParseRoleGrantedGlobally(log types.Log) (*RoleManagerRoleGrantedGlobally, error) {
	event := new(RoleManagerRoleGrantedGlobally)
	if err := _RoleManager.contract.UnpackLog(event, "RoleGrantedGlobally", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RoleManager contract.
type RoleManagerRoleRevokedIterator struct {
	Event *RoleManagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RoleManagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerRoleRevoked)
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
		it.Event = new(RoleManagerRoleRevoked)
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
func (it *RoleManagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerRoleRevoked represents a RoleRevoked event raised by the RoleManager contract.
type RoleManagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RoleManagerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RoleManagerRoleRevokedIterator{contract: _RoleManager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RoleManagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerRoleRevoked)
				if err := _RoleManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) ParseRoleRevoked(log types.Log) (*RoleManagerRoleRevoked, error) {
	event := new(RoleManagerRoleRevoked)
	if err := _RoleManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerRoleRevokedGloballyIterator is returned from FilterRoleRevokedGlobally and is used to iterate over the raw logs and unpacked data for RoleRevokedGlobally events raised by the RoleManager contract.
type RoleManagerRoleRevokedGloballyIterator struct {
	Event *RoleManagerRoleRevokedGlobally // Event containing the contract specifics and raw log

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
func (it *RoleManagerRoleRevokedGloballyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerRoleRevokedGlobally)
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
		it.Event = new(RoleManagerRoleRevokedGlobally)
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
func (it *RoleManagerRoleRevokedGloballyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerRoleRevokedGloballyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerRoleRevokedGlobally represents a RoleRevokedGlobally event raised by the RoleManager contract.
type RoleManagerRoleRevokedGlobally struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevokedGlobally is a free log retrieval operation binding the contract event 0xd8476cf9ac7354675e06297f0a3b1f57e056db7e4cefe100071b4ade09836415.
//
// Solidity: event RoleRevokedGlobally(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) FilterRoleRevokedGlobally(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RoleManagerRoleRevokedGloballyIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "RoleRevokedGlobally", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RoleManagerRoleRevokedGloballyIterator{contract: _RoleManager.contract, event: "RoleRevokedGlobally", logs: logs, sub: sub}, nil
}

// WatchRoleRevokedGlobally is a free log subscription operation binding the contract event 0xd8476cf9ac7354675e06297f0a3b1f57e056db7e4cefe100071b4ade09836415.
//
// Solidity: event RoleRevokedGlobally(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) WatchRoleRevokedGlobally(opts *bind.WatchOpts, sink chan<- *RoleManagerRoleRevokedGlobally, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "RoleRevokedGlobally", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerRoleRevokedGlobally)
				if err := _RoleManager.contract.UnpackLog(event, "RoleRevokedGlobally", log); err != nil {
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

// ParseRoleRevokedGlobally is a log parse operation binding the contract event 0xd8476cf9ac7354675e06297f0a3b1f57e056db7e4cefe100071b4ade09836415.
//
// Solidity: event RoleRevokedGlobally(bytes32 indexed role, address indexed account, address indexed sender)
func (_RoleManager *RoleManagerFilterer) ParseRoleRevokedGlobally(log types.Log) (*RoleManagerRoleRevokedGlobally, error) {
	event := new(RoleManagerRoleRevokedGlobally)
	if err := _RoleManager.contract.UnpackLog(event, "RoleRevokedGlobally", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RoleManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RoleManager contract.
type RoleManagerUnpausedIterator struct {
	Event *RoleManagerUnpaused // Event containing the contract specifics and raw log

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
func (it *RoleManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoleManagerUnpaused)
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
		it.Event = new(RoleManagerUnpaused)
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
func (it *RoleManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoleManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoleManagerUnpaused represents a Unpaused event raised by the RoleManager contract.
type RoleManagerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RoleManager *RoleManagerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RoleManagerUnpausedIterator, error) {

	logs, sub, err := _RoleManager.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RoleManagerUnpausedIterator{contract: _RoleManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RoleManager *RoleManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RoleManagerUnpaused) (event.Subscription, error) {

	logs, sub, err := _RoleManager.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoleManagerUnpaused)
				if err := _RoleManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_RoleManager *RoleManagerFilterer) ParseUnpaused(log types.Log) (*RoleManagerUnpaused, error) {
	event := new(RoleManagerUnpaused)
	if err := _RoleManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
