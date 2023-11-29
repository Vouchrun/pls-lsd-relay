// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fee_pool

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

// FeePoolMetaData contains all meta data concerning the FeePool contract.
var FeePoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyClaimed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyDealedEpoch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyDealedHeight\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyNotifiedCycle\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyVoted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountNotZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountUnmatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BalanceNotEnough\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BlockNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CallerNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableDepositZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableRewardZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableWithdrawIndexOverflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CommissionRateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CycleNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DepositAmountGTMaxAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DepositAmountLTMinAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EthAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMerkleProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidThreshold\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LengthNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LsdTokenAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeAlreadyRemoved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeNotClaimable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAuthorizedLsdToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotClaimable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotPubkeyOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotTrustNode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProposalExecFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyNumberOverLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyStatusUnmatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateChangeOverLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SoloNodeDepositAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SoloNodeDepositDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubmitBalancesDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"min\",\"type\":\"uint256\"}],\"name\":\"TooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustNodeDepositDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UserDepositDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VotersDuplicate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VotersNotEnough\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VotersNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawIndexEmpty\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"EtherWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_networkWithdrawAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_networkProposalAddress\",\"type\":\"address\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkProposalAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkWithdrawAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reinit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// FeePoolABI is the input ABI used to generate the binding from.
// Deprecated: Use FeePoolMetaData.ABI instead.
var FeePoolABI = FeePoolMetaData.ABI

// FeePool is an auto generated Go binding around an Ethereum contract.
type FeePool struct {
	FeePoolCaller     // Read-only binding to the contract
	FeePoolTransactor // Write-only binding to the contract
	FeePoolFilterer   // Log filterer for contract events
}

// FeePoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeePoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeePoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeePoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeePoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeePoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeePoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeePoolSession struct {
	Contract     *FeePool          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeePoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeePoolCallerSession struct {
	Contract *FeePoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// FeePoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeePoolTransactorSession struct {
	Contract     *FeePoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// FeePoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeePoolRaw struct {
	Contract *FeePool // Generic contract binding to access the raw methods on
}

// FeePoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeePoolCallerRaw struct {
	Contract *FeePoolCaller // Generic read-only contract binding to access the raw methods on
}

// FeePoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeePoolTransactorRaw struct {
	Contract *FeePoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeePool creates a new instance of FeePool, bound to a specific deployed contract.
func NewFeePool(address common.Address, backend bind.ContractBackend) (*FeePool, error) {
	contract, err := bindFeePool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeePool{FeePoolCaller: FeePoolCaller{contract: contract}, FeePoolTransactor: FeePoolTransactor{contract: contract}, FeePoolFilterer: FeePoolFilterer{contract: contract}}, nil
}

// NewFeePoolCaller creates a new read-only instance of FeePool, bound to a specific deployed contract.
func NewFeePoolCaller(address common.Address, caller bind.ContractCaller) (*FeePoolCaller, error) {
	contract, err := bindFeePool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeePoolCaller{contract: contract}, nil
}

// NewFeePoolTransactor creates a new write-only instance of FeePool, bound to a specific deployed contract.
func NewFeePoolTransactor(address common.Address, transactor bind.ContractTransactor) (*FeePoolTransactor, error) {
	contract, err := bindFeePool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeePoolTransactor{contract: contract}, nil
}

// NewFeePoolFilterer creates a new log filterer instance of FeePool, bound to a specific deployed contract.
func NewFeePoolFilterer(address common.Address, filterer bind.ContractFilterer) (*FeePoolFilterer, error) {
	contract, err := bindFeePool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeePoolFilterer{contract: contract}, nil
}

// bindFeePool binds a generic wrapper to an already deployed contract.
func bindFeePool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeePoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeePool *FeePoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeePool.Contract.FeePoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeePool *FeePoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeePool.Contract.FeePoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeePool *FeePoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeePool.Contract.FeePoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeePool *FeePoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeePool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeePool *FeePoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeePool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeePool *FeePoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeePool.Contract.contract.Transact(opts, method, params...)
}

// NetworkProposalAddress is a free data retrieval call binding the contract method 0xb4701c09.
//
// Solidity: function networkProposalAddress() view returns(address)
func (_FeePool *FeePoolCaller) NetworkProposalAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeePool.contract.Call(opts, &out, "networkProposalAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NetworkProposalAddress is a free data retrieval call binding the contract method 0xb4701c09.
//
// Solidity: function networkProposalAddress() view returns(address)
func (_FeePool *FeePoolSession) NetworkProposalAddress() (common.Address, error) {
	return _FeePool.Contract.NetworkProposalAddress(&_FeePool.CallOpts)
}

// NetworkProposalAddress is a free data retrieval call binding the contract method 0xb4701c09.
//
// Solidity: function networkProposalAddress() view returns(address)
func (_FeePool *FeePoolCallerSession) NetworkProposalAddress() (common.Address, error) {
	return _FeePool.Contract.NetworkProposalAddress(&_FeePool.CallOpts)
}

// NetworkWithdrawAddress is a free data retrieval call binding the contract method 0x1587ef1b.
//
// Solidity: function networkWithdrawAddress() view returns(address)
func (_FeePool *FeePoolCaller) NetworkWithdrawAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeePool.contract.Call(opts, &out, "networkWithdrawAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NetworkWithdrawAddress is a free data retrieval call binding the contract method 0x1587ef1b.
//
// Solidity: function networkWithdrawAddress() view returns(address)
func (_FeePool *FeePoolSession) NetworkWithdrawAddress() (common.Address, error) {
	return _FeePool.Contract.NetworkWithdrawAddress(&_FeePool.CallOpts)
}

// NetworkWithdrawAddress is a free data retrieval call binding the contract method 0x1587ef1b.
//
// Solidity: function networkWithdrawAddress() view returns(address)
func (_FeePool *FeePoolCallerSession) NetworkWithdrawAddress() (common.Address, error) {
	return _FeePool.Contract.NetworkWithdrawAddress(&_FeePool.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FeePool *FeePoolCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FeePool.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FeePool *FeePoolSession) ProxiableUUID() ([32]byte, error) {
	return _FeePool.Contract.ProxiableUUID(&_FeePool.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FeePool *FeePoolCallerSession) ProxiableUUID() ([32]byte, error) {
	return _FeePool.Contract.ProxiableUUID(&_FeePool.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_FeePool *FeePoolCaller) Version(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FeePool.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_FeePool *FeePoolSession) Version() (uint8, error) {
	return _FeePool.Contract.Version(&_FeePool.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_FeePool *FeePoolCallerSession) Version() (uint8, error) {
	return _FeePool.Contract.Version(&_FeePool.CallOpts)
}

// Init is a paid mutator transaction binding the contract method 0xf09a4016.
//
// Solidity: function init(address _networkWithdrawAddress, address _networkProposalAddress) returns()
func (_FeePool *FeePoolTransactor) Init(opts *bind.TransactOpts, _networkWithdrawAddress common.Address, _networkProposalAddress common.Address) (*types.Transaction, error) {
	return _FeePool.contract.Transact(opts, "init", _networkWithdrawAddress, _networkProposalAddress)
}

// Init is a paid mutator transaction binding the contract method 0xf09a4016.
//
// Solidity: function init(address _networkWithdrawAddress, address _networkProposalAddress) returns()
func (_FeePool *FeePoolSession) Init(_networkWithdrawAddress common.Address, _networkProposalAddress common.Address) (*types.Transaction, error) {
	return _FeePool.Contract.Init(&_FeePool.TransactOpts, _networkWithdrawAddress, _networkProposalAddress)
}

// Init is a paid mutator transaction binding the contract method 0xf09a4016.
//
// Solidity: function init(address _networkWithdrawAddress, address _networkProposalAddress) returns()
func (_FeePool *FeePoolTransactorSession) Init(_networkWithdrawAddress common.Address, _networkProposalAddress common.Address) (*types.Transaction, error) {
	return _FeePool.Contract.Init(&_FeePool.TransactOpts, _networkWithdrawAddress, _networkProposalAddress)
}

// Reinit is a paid mutator transaction binding the contract method 0xc482ceaf.
//
// Solidity: function reinit() returns()
func (_FeePool *FeePoolTransactor) Reinit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeePool.contract.Transact(opts, "reinit")
}

// Reinit is a paid mutator transaction binding the contract method 0xc482ceaf.
//
// Solidity: function reinit() returns()
func (_FeePool *FeePoolSession) Reinit() (*types.Transaction, error) {
	return _FeePool.Contract.Reinit(&_FeePool.TransactOpts)
}

// Reinit is a paid mutator transaction binding the contract method 0xc482ceaf.
//
// Solidity: function reinit() returns()
func (_FeePool *FeePoolTransactorSession) Reinit() (*types.Transaction, error) {
	return _FeePool.Contract.Reinit(&_FeePool.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FeePool *FeePoolTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _FeePool.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FeePool *FeePoolSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FeePool.Contract.UpgradeTo(&_FeePool.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FeePool *FeePoolTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FeePool.Contract.UpgradeTo(&_FeePool.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FeePool *FeePoolTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FeePool.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FeePool *FeePoolSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FeePool.Contract.UpgradeToAndCall(&_FeePool.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FeePool *FeePoolTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FeePool.Contract.UpgradeToAndCall(&_FeePool.TransactOpts, newImplementation, data)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0x3bed33ce.
//
// Solidity: function withdrawEther(uint256 _amount) returns()
func (_FeePool *FeePoolTransactor) WithdrawEther(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _FeePool.contract.Transact(opts, "withdrawEther", _amount)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0x3bed33ce.
//
// Solidity: function withdrawEther(uint256 _amount) returns()
func (_FeePool *FeePoolSession) WithdrawEther(_amount *big.Int) (*types.Transaction, error) {
	return _FeePool.Contract.WithdrawEther(&_FeePool.TransactOpts, _amount)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0x3bed33ce.
//
// Solidity: function withdrawEther(uint256 _amount) returns()
func (_FeePool *FeePoolTransactorSession) WithdrawEther(_amount *big.Int) (*types.Transaction, error) {
	return _FeePool.Contract.WithdrawEther(&_FeePool.TransactOpts, _amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FeePool *FeePoolTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeePool.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FeePool *FeePoolSession) Receive() (*types.Transaction, error) {
	return _FeePool.Contract.Receive(&_FeePool.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FeePool *FeePoolTransactorSession) Receive() (*types.Transaction, error) {
	return _FeePool.Contract.Receive(&_FeePool.TransactOpts)
}

// FeePoolAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the FeePool contract.
type FeePoolAdminChangedIterator struct {
	Event *FeePoolAdminChanged // Event containing the contract specifics and raw log

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
func (it *FeePoolAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeePoolAdminChanged)
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
		it.Event = new(FeePoolAdminChanged)
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
func (it *FeePoolAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeePoolAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeePoolAdminChanged represents a AdminChanged event raised by the FeePool contract.
type FeePoolAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FeePool *FeePoolFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*FeePoolAdminChangedIterator, error) {

	logs, sub, err := _FeePool.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &FeePoolAdminChangedIterator{contract: _FeePool.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FeePool *FeePoolFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *FeePoolAdminChanged) (event.Subscription, error) {

	logs, sub, err := _FeePool.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeePoolAdminChanged)
				if err := _FeePool.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FeePool *FeePoolFilterer) ParseAdminChanged(log types.Log) (*FeePoolAdminChanged, error) {
	event := new(FeePoolAdminChanged)
	if err := _FeePool.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeePoolBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the FeePool contract.
type FeePoolBeaconUpgradedIterator struct {
	Event *FeePoolBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *FeePoolBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeePoolBeaconUpgraded)
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
		it.Event = new(FeePoolBeaconUpgraded)
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
func (it *FeePoolBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeePoolBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeePoolBeaconUpgraded represents a BeaconUpgraded event raised by the FeePool contract.
type FeePoolBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FeePool *FeePoolFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*FeePoolBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FeePool.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &FeePoolBeaconUpgradedIterator{contract: _FeePool.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FeePool *FeePoolFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *FeePoolBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FeePool.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeePoolBeaconUpgraded)
				if err := _FeePool.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FeePool *FeePoolFilterer) ParseBeaconUpgraded(log types.Log) (*FeePoolBeaconUpgraded, error) {
	event := new(FeePoolBeaconUpgraded)
	if err := _FeePool.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeePoolEtherWithdrawnIterator is returned from FilterEtherWithdrawn and is used to iterate over the raw logs and unpacked data for EtherWithdrawn events raised by the FeePool contract.
type FeePoolEtherWithdrawnIterator struct {
	Event *FeePoolEtherWithdrawn // Event containing the contract specifics and raw log

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
func (it *FeePoolEtherWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeePoolEtherWithdrawn)
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
		it.Event = new(FeePoolEtherWithdrawn)
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
func (it *FeePoolEtherWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeePoolEtherWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeePoolEtherWithdrawn represents a EtherWithdrawn event raised by the FeePool contract.
type FeePoolEtherWithdrawn struct {
	Amount *big.Int
	Time   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdrawn is a free log retrieval operation binding the contract event 0xfe8f3b8f6ceb26dd77f4305e42b5199c054cd056585d04be113a00c6ea75acb1.
//
// Solidity: event EtherWithdrawn(uint256 amount, uint256 time)
func (_FeePool *FeePoolFilterer) FilterEtherWithdrawn(opts *bind.FilterOpts) (*FeePoolEtherWithdrawnIterator, error) {

	logs, sub, err := _FeePool.contract.FilterLogs(opts, "EtherWithdrawn")
	if err != nil {
		return nil, err
	}
	return &FeePoolEtherWithdrawnIterator{contract: _FeePool.contract, event: "EtherWithdrawn", logs: logs, sub: sub}, nil
}

// WatchEtherWithdrawn is a free log subscription operation binding the contract event 0xfe8f3b8f6ceb26dd77f4305e42b5199c054cd056585d04be113a00c6ea75acb1.
//
// Solidity: event EtherWithdrawn(uint256 amount, uint256 time)
func (_FeePool *FeePoolFilterer) WatchEtherWithdrawn(opts *bind.WatchOpts, sink chan<- *FeePoolEtherWithdrawn) (event.Subscription, error) {

	logs, sub, err := _FeePool.contract.WatchLogs(opts, "EtherWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeePoolEtherWithdrawn)
				if err := _FeePool.contract.UnpackLog(event, "EtherWithdrawn", log); err != nil {
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

// ParseEtherWithdrawn is a log parse operation binding the contract event 0xfe8f3b8f6ceb26dd77f4305e42b5199c054cd056585d04be113a00c6ea75acb1.
//
// Solidity: event EtherWithdrawn(uint256 amount, uint256 time)
func (_FeePool *FeePoolFilterer) ParseEtherWithdrawn(log types.Log) (*FeePoolEtherWithdrawn, error) {
	event := new(FeePoolEtherWithdrawn)
	if err := _FeePool.contract.UnpackLog(event, "EtherWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeePoolInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FeePool contract.
type FeePoolInitializedIterator struct {
	Event *FeePoolInitialized // Event containing the contract specifics and raw log

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
func (it *FeePoolInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeePoolInitialized)
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
		it.Event = new(FeePoolInitialized)
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
func (it *FeePoolInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeePoolInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeePoolInitialized represents a Initialized event raised by the FeePool contract.
type FeePoolInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FeePool *FeePoolFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeePoolInitializedIterator, error) {

	logs, sub, err := _FeePool.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeePoolInitializedIterator{contract: _FeePool.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FeePool *FeePoolFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeePoolInitialized) (event.Subscription, error) {

	logs, sub, err := _FeePool.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeePoolInitialized)
				if err := _FeePool.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FeePool *FeePoolFilterer) ParseInitialized(log types.Log) (*FeePoolInitialized, error) {
	event := new(FeePoolInitialized)
	if err := _FeePool.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeePoolUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FeePool contract.
type FeePoolUpgradedIterator struct {
	Event *FeePoolUpgraded // Event containing the contract specifics and raw log

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
func (it *FeePoolUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeePoolUpgraded)
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
		it.Event = new(FeePoolUpgraded)
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
func (it *FeePoolUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeePoolUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeePoolUpgraded represents a Upgraded event raised by the FeePool contract.
type FeePoolUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FeePool *FeePoolFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FeePoolUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FeePool.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FeePoolUpgradedIterator{contract: _FeePool.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FeePool *FeePoolFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FeePoolUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FeePool.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeePoolUpgraded)
				if err := _FeePool.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FeePool *FeePoolFilterer) ParseUpgraded(log types.Log) (*FeePoolUpgraded, error) {
	event := new(FeePoolUpgraded)
	if err := _FeePool.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
