// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package network_withdrawal

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

// NetworkWithdrawalMetaData contains all meta data concerning the NetworkWithdrawal contract.
var NetworkWithdrawalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyClaimed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyDealtEpoch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyDealtHeight\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyNotifiedCycle\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyVoted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountNotZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountUnmatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BalanceNotEnough\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BlockNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CallerNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableDepositZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableRewardZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimableWithdrawalIndexOverflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CommissionRateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CycleNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DepositAmountLTMinAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EthAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMerkleProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidThreshold\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LengthNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LsdTokenAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeAlreadyRemoved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeNotClaimable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAuthorizedLsdToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotClaimable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotPubkeyOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotTrustNode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProposalExecFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyNumberOverLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PubkeyStatusUnmatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateChangeOverLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SecondsZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SoloNodeDepositAmountZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SoloNodeDepositDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubmitBalancesDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustNodeDepositDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UserDepositDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VotersDuplicate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VotersNotEnough\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VotersNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalIndexEmpty\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumINetworkWithdrawal.DistributionType\",\"name\":\"distributeType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dealtHeight\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"userAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nodeAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"platformAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxClaimableWithdrawIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"mvAmount\",\"type\":\"uint256\"}],\"name\":\"DistributeRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"EtherDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"claimableReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"claimableDeposit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumINetworkWithdrawal.ClaimType\",\"name\":\"claimType\",\"type\":\"uint8\"}],\"name\":\"NodeClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawalCycle\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ejectedStartWithdrawalCycle\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ejectedValidators\",\"type\":\"uint256[]\"}],\"name\":\"NotifyValidatorExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"dealtEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"nodeRewardsFileCid\",\"type\":\"string\"}],\"name\":\"SetMerkleRoot\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cycleSeconds\",\"type\":\"uint256\"}],\"name\":\"SetWithdrawalCycleSeconds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lsdTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"instantly\",\"type\":\"bool\"}],\"name\":\"Unstake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"withdrawIndexList\",\"type\":\"uint256[]\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"currentWithdrawalCycle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositEthAndUpdateTotalShortages\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumINetworkWithdrawal.DistributionType\",\"name\":\"_distributionType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_dealtHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_userAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nodeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_platformAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxClaimableWithdrawalIndex\",\"type\":\"uint256\"}],\"name\":\"distribute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ejectedStartCycle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ejectedValidatorsAtCycle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factoryAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factoryCommissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feePoolAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"cycle\",\"type\":\"uint256\"}],\"name\":\"getEjectedValidatorsAtCycle\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUnclaimedWithdrawalsOfUser\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_lsdTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_userDepositAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_networkProposalAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_networkBalancesAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_feePoolAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_factoryAddress\",\"type\":\"address\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestDistributionPriorityFeeHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestDistributionWithdrawalHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestMerkleRootEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lsdTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxClaimableWithdrawalIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkBalancesAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkProposalAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextWithdrawalIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_node\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_totalRewardAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalExitDepositAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"_merkleProof\",\"type\":\"bytes32[]\"},{\"internalType\":\"enumINetworkWithdrawal.ClaimType\",\"name\":\"_claimType\",\"type\":\"uint8\"}],\"name\":\"nodeClaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nodeClaimEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nodeCommissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nodeRewardsFileCid\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_withdrawCycle\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_ejectedStartCycle\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_validatorIndexList\",\"type\":\"uint256[]\"}],\"name\":\"notifyValidatorExit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"platformClaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"platformCommissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reinit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_factoryCommissionRate\",\"type\":\"uint256\"}],\"name\":\"setFactoryCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_dealtEpoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_merkleRoot\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_nodeRewardsFileCid\",\"type\":\"string\"}],\"name\":\"setMerkleRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_value\",\"type\":\"bool\"}],\"name\":\"setNodeClaimEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_platformCommissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nodeCommissionRate\",\"type\":\"uint256\"}],\"name\":\"setPlatformAndNodeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_withdrawalCycleSeconds\",\"type\":\"uint256\"}],\"name\":\"setWithdrawalCycleSeconds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalClaimedDepositOfNode\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalClaimedRewardOfNode\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalPlatformClaimedCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalPlatformCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalWithdrawalShortages\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lsdTokenAmount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"userDepositAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_withdrawalIndexList\",\"type\":\"uint256[]\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"withdrawalAtIndex\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawalCycleSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// NetworkWithdrawalABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkWithdrawalMetaData.ABI instead.
var NetworkWithdrawalABI = NetworkWithdrawalMetaData.ABI

// NetworkWithdrawal is an auto generated Go binding around an Ethereum contract.
type NetworkWithdrawal struct {
	NetworkWithdrawalCaller     // Read-only binding to the contract
	NetworkWithdrawalTransactor // Write-only binding to the contract
	NetworkWithdrawalFilterer   // Log filterer for contract events
}

// NetworkWithdrawalCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkWithdrawalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkWithdrawalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkWithdrawalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkWithdrawalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkWithdrawalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkWithdrawalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkWithdrawalSession struct {
	Contract     *NetworkWithdrawal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// NetworkWithdrawalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkWithdrawalCallerSession struct {
	Contract *NetworkWithdrawalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// NetworkWithdrawalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkWithdrawalTransactorSession struct {
	Contract     *NetworkWithdrawalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// NetworkWithdrawalRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkWithdrawalRaw struct {
	Contract *NetworkWithdrawal // Generic contract binding to access the raw methods on
}

// NetworkWithdrawalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkWithdrawalCallerRaw struct {
	Contract *NetworkWithdrawalCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkWithdrawalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkWithdrawalTransactorRaw struct {
	Contract *NetworkWithdrawalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkWithdrawal creates a new instance of NetworkWithdrawal, bound to a specific deployed contract.
func NewNetworkWithdrawal(address common.Address, backend bind.ContractBackend) (*NetworkWithdrawal, error) {
	contract, err := bindNetworkWithdrawal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawal{NetworkWithdrawalCaller: NetworkWithdrawalCaller{contract: contract}, NetworkWithdrawalTransactor: NetworkWithdrawalTransactor{contract: contract}, NetworkWithdrawalFilterer: NetworkWithdrawalFilterer{contract: contract}}, nil
}

// NewNetworkWithdrawalCaller creates a new read-only instance of NetworkWithdrawal, bound to a specific deployed contract.
func NewNetworkWithdrawalCaller(address common.Address, caller bind.ContractCaller) (*NetworkWithdrawalCaller, error) {
	contract, err := bindNetworkWithdrawal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalCaller{contract: contract}, nil
}

// NewNetworkWithdrawalTransactor creates a new write-only instance of NetworkWithdrawal, bound to a specific deployed contract.
func NewNetworkWithdrawalTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkWithdrawalTransactor, error) {
	contract, err := bindNetworkWithdrawal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalTransactor{contract: contract}, nil
}

// NewNetworkWithdrawalFilterer creates a new log filterer instance of NetworkWithdrawal, bound to a specific deployed contract.
func NewNetworkWithdrawalFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkWithdrawalFilterer, error) {
	contract, err := bindNetworkWithdrawal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalFilterer{contract: contract}, nil
}

// bindNetworkWithdrawal binds a generic wrapper to an already deployed contract.
func bindNetworkWithdrawal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NetworkWithdrawalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkWithdrawal *NetworkWithdrawalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkWithdrawal.Contract.NetworkWithdrawalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkWithdrawal *NetworkWithdrawalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.NetworkWithdrawalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkWithdrawal *NetworkWithdrawalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.NetworkWithdrawalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkWithdrawal *NetworkWithdrawalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkWithdrawal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkWithdrawal *NetworkWithdrawalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkWithdrawal *NetworkWithdrawalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.contract.Transact(opts, method, params...)
}

// CurrentWithdrawalCycle is a free data retrieval call binding the contract method 0x8c012ffb.
//
// Solidity: function currentWithdrawalCycle() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) CurrentWithdrawalCycle(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "currentWithdrawalCycle")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentWithdrawalCycle is a free data retrieval call binding the contract method 0x8c012ffb.
//
// Solidity: function currentWithdrawalCycle() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) CurrentWithdrawalCycle() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.CurrentWithdrawalCycle(&_NetworkWithdrawal.CallOpts)
}

// CurrentWithdrawalCycle is a free data retrieval call binding the contract method 0x8c012ffb.
//
// Solidity: function currentWithdrawalCycle() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) CurrentWithdrawalCycle() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.CurrentWithdrawalCycle(&_NetworkWithdrawal.CallOpts)
}

// EjectedStartCycle is a free data retrieval call binding the contract method 0x8a699828.
//
// Solidity: function ejectedStartCycle() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) EjectedStartCycle(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "ejectedStartCycle")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EjectedStartCycle is a free data retrieval call binding the contract method 0x8a699828.
//
// Solidity: function ejectedStartCycle() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) EjectedStartCycle() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.EjectedStartCycle(&_NetworkWithdrawal.CallOpts)
}

// EjectedStartCycle is a free data retrieval call binding the contract method 0x8a699828.
//
// Solidity: function ejectedStartCycle() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) EjectedStartCycle() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.EjectedStartCycle(&_NetworkWithdrawal.CallOpts)
}

// EjectedValidatorsAtCycle is a free data retrieval call binding the contract method 0x261a792d.
//
// Solidity: function ejectedValidatorsAtCycle(uint256 , uint256 ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) EjectedValidatorsAtCycle(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "ejectedValidatorsAtCycle", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EjectedValidatorsAtCycle is a free data retrieval call binding the contract method 0x261a792d.
//
// Solidity: function ejectedValidatorsAtCycle(uint256 , uint256 ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) EjectedValidatorsAtCycle(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _NetworkWithdrawal.Contract.EjectedValidatorsAtCycle(&_NetworkWithdrawal.CallOpts, arg0, arg1)
}

// EjectedValidatorsAtCycle is a free data retrieval call binding the contract method 0x261a792d.
//
// Solidity: function ejectedValidatorsAtCycle(uint256 , uint256 ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) EjectedValidatorsAtCycle(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _NetworkWithdrawal.Contract.EjectedValidatorsAtCycle(&_NetworkWithdrawal.CallOpts, arg0, arg1)
}

// FactoryAddress is a free data retrieval call binding the contract method 0x966dae0e.
//
// Solidity: function factoryAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) FactoryAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "factoryAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FactoryAddress is a free data retrieval call binding the contract method 0x966dae0e.
//
// Solidity: function factoryAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalSession) FactoryAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.FactoryAddress(&_NetworkWithdrawal.CallOpts)
}

// FactoryAddress is a free data retrieval call binding the contract method 0x966dae0e.
//
// Solidity: function factoryAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) FactoryAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.FactoryAddress(&_NetworkWithdrawal.CallOpts)
}

// FactoryCommissionRate is a free data retrieval call binding the contract method 0xc2156b4b.
//
// Solidity: function factoryCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) FactoryCommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "factoryCommissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FactoryCommissionRate is a free data retrieval call binding the contract method 0xc2156b4b.
//
// Solidity: function factoryCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) FactoryCommissionRate() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.FactoryCommissionRate(&_NetworkWithdrawal.CallOpts)
}

// FactoryCommissionRate is a free data retrieval call binding the contract method 0xc2156b4b.
//
// Solidity: function factoryCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) FactoryCommissionRate() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.FactoryCommissionRate(&_NetworkWithdrawal.CallOpts)
}

// FeePoolAddress is a free data retrieval call binding the contract method 0x4319ebe4.
//
// Solidity: function feePoolAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) FeePoolAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "feePoolAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeePoolAddress is a free data retrieval call binding the contract method 0x4319ebe4.
//
// Solidity: function feePoolAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalSession) FeePoolAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.FeePoolAddress(&_NetworkWithdrawal.CallOpts)
}

// FeePoolAddress is a free data retrieval call binding the contract method 0x4319ebe4.
//
// Solidity: function feePoolAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) FeePoolAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.FeePoolAddress(&_NetworkWithdrawal.CallOpts)
}

// GetEjectedValidatorsAtCycle is a free data retrieval call binding the contract method 0x2c0f4166.
//
// Solidity: function getEjectedValidatorsAtCycle(uint256 cycle) view returns(uint256[])
func (_NetworkWithdrawal *NetworkWithdrawalCaller) GetEjectedValidatorsAtCycle(opts *bind.CallOpts, cycle *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "getEjectedValidatorsAtCycle", cycle)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetEjectedValidatorsAtCycle is a free data retrieval call binding the contract method 0x2c0f4166.
//
// Solidity: function getEjectedValidatorsAtCycle(uint256 cycle) view returns(uint256[])
func (_NetworkWithdrawal *NetworkWithdrawalSession) GetEjectedValidatorsAtCycle(cycle *big.Int) ([]*big.Int, error) {
	return _NetworkWithdrawal.Contract.GetEjectedValidatorsAtCycle(&_NetworkWithdrawal.CallOpts, cycle)
}

// GetEjectedValidatorsAtCycle is a free data retrieval call binding the contract method 0x2c0f4166.
//
// Solidity: function getEjectedValidatorsAtCycle(uint256 cycle) view returns(uint256[])
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) GetEjectedValidatorsAtCycle(cycle *big.Int) ([]*big.Int, error) {
	return _NetworkWithdrawal.Contract.GetEjectedValidatorsAtCycle(&_NetworkWithdrawal.CallOpts, cycle)
}

// GetUnclaimedWithdrawalsOfUser is a free data retrieval call binding the contract method 0xfd6b5a49.
//
// Solidity: function getUnclaimedWithdrawalsOfUser(address user) view returns(uint256[])
func (_NetworkWithdrawal *NetworkWithdrawalCaller) GetUnclaimedWithdrawalsOfUser(opts *bind.CallOpts, user common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "getUnclaimedWithdrawalsOfUser", user)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUnclaimedWithdrawalsOfUser is a free data retrieval call binding the contract method 0xfd6b5a49.
//
// Solidity: function getUnclaimedWithdrawalsOfUser(address user) view returns(uint256[])
func (_NetworkWithdrawal *NetworkWithdrawalSession) GetUnclaimedWithdrawalsOfUser(user common.Address) ([]*big.Int, error) {
	return _NetworkWithdrawal.Contract.GetUnclaimedWithdrawalsOfUser(&_NetworkWithdrawal.CallOpts, user)
}

// GetUnclaimedWithdrawalsOfUser is a free data retrieval call binding the contract method 0xfd6b5a49.
//
// Solidity: function getUnclaimedWithdrawalsOfUser(address user) view returns(uint256[])
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) GetUnclaimedWithdrawalsOfUser(user common.Address) ([]*big.Int, error) {
	return _NetworkWithdrawal.Contract.GetUnclaimedWithdrawalsOfUser(&_NetworkWithdrawal.CallOpts, user)
}

// LatestDistributionPriorityFeeHeight is a free data retrieval call binding the contract method 0xb360a0f6.
//
// Solidity: function latestDistributionPriorityFeeHeight() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) LatestDistributionPriorityFeeHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "latestDistributionPriorityFeeHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestDistributionPriorityFeeHeight is a free data retrieval call binding the contract method 0xb360a0f6.
//
// Solidity: function latestDistributionPriorityFeeHeight() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) LatestDistributionPriorityFeeHeight() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.LatestDistributionPriorityFeeHeight(&_NetworkWithdrawal.CallOpts)
}

// LatestDistributionPriorityFeeHeight is a free data retrieval call binding the contract method 0xb360a0f6.
//
// Solidity: function latestDistributionPriorityFeeHeight() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) LatestDistributionPriorityFeeHeight() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.LatestDistributionPriorityFeeHeight(&_NetworkWithdrawal.CallOpts)
}

// LatestDistributionWithdrawalHeight is a free data retrieval call binding the contract method 0xf5cb8807.
//
// Solidity: function latestDistributionWithdrawalHeight() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) LatestDistributionWithdrawalHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "latestDistributionWithdrawalHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestDistributionWithdrawalHeight is a free data retrieval call binding the contract method 0xf5cb8807.
//
// Solidity: function latestDistributionWithdrawalHeight() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) LatestDistributionWithdrawalHeight() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.LatestDistributionWithdrawalHeight(&_NetworkWithdrawal.CallOpts)
}

// LatestDistributionWithdrawalHeight is a free data retrieval call binding the contract method 0xf5cb8807.
//
// Solidity: function latestDistributionWithdrawalHeight() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) LatestDistributionWithdrawalHeight() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.LatestDistributionWithdrawalHeight(&_NetworkWithdrawal.CallOpts)
}

// LatestMerkleRootEpoch is a free data retrieval call binding the contract method 0xb5ca7410.
//
// Solidity: function latestMerkleRootEpoch() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) LatestMerkleRootEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "latestMerkleRootEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestMerkleRootEpoch is a free data retrieval call binding the contract method 0xb5ca7410.
//
// Solidity: function latestMerkleRootEpoch() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) LatestMerkleRootEpoch() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.LatestMerkleRootEpoch(&_NetworkWithdrawal.CallOpts)
}

// LatestMerkleRootEpoch is a free data retrieval call binding the contract method 0xb5ca7410.
//
// Solidity: function latestMerkleRootEpoch() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) LatestMerkleRootEpoch() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.LatestMerkleRootEpoch(&_NetworkWithdrawal.CallOpts)
}

// LsdTokenAddress is a free data retrieval call binding the contract method 0x87505b9d.
//
// Solidity: function lsdTokenAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) LsdTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "lsdTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LsdTokenAddress is a free data retrieval call binding the contract method 0x87505b9d.
//
// Solidity: function lsdTokenAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalSession) LsdTokenAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.LsdTokenAddress(&_NetworkWithdrawal.CallOpts)
}

// LsdTokenAddress is a free data retrieval call binding the contract method 0x87505b9d.
//
// Solidity: function lsdTokenAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) LsdTokenAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.LsdTokenAddress(&_NetworkWithdrawal.CallOpts)
}

// MaxClaimableWithdrawalIndex is a free data retrieval call binding the contract method 0x3330641b.
//
// Solidity: function maxClaimableWithdrawalIndex() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) MaxClaimableWithdrawalIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "maxClaimableWithdrawalIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxClaimableWithdrawalIndex is a free data retrieval call binding the contract method 0x3330641b.
//
// Solidity: function maxClaimableWithdrawalIndex() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) MaxClaimableWithdrawalIndex() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.MaxClaimableWithdrawalIndex(&_NetworkWithdrawal.CallOpts)
}

// MaxClaimableWithdrawalIndex is a free data retrieval call binding the contract method 0x3330641b.
//
// Solidity: function maxClaimableWithdrawalIndex() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) MaxClaimableWithdrawalIndex() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.MaxClaimableWithdrawalIndex(&_NetworkWithdrawal.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_NetworkWithdrawal *NetworkWithdrawalSession) MerkleRoot() ([32]byte, error) {
	return _NetworkWithdrawal.Contract.MerkleRoot(&_NetworkWithdrawal.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) MerkleRoot() ([32]byte, error) {
	return _NetworkWithdrawal.Contract.MerkleRoot(&_NetworkWithdrawal.CallOpts)
}

// NetworkBalancesAddress is a free data retrieval call binding the contract method 0x38fcf092.
//
// Solidity: function networkBalancesAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) NetworkBalancesAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "networkBalancesAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NetworkBalancesAddress is a free data retrieval call binding the contract method 0x38fcf092.
//
// Solidity: function networkBalancesAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalSession) NetworkBalancesAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.NetworkBalancesAddress(&_NetworkWithdrawal.CallOpts)
}

// NetworkBalancesAddress is a free data retrieval call binding the contract method 0x38fcf092.
//
// Solidity: function networkBalancesAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) NetworkBalancesAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.NetworkBalancesAddress(&_NetworkWithdrawal.CallOpts)
}

// NetworkProposalAddress is a free data retrieval call binding the contract method 0xb4701c09.
//
// Solidity: function networkProposalAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) NetworkProposalAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "networkProposalAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NetworkProposalAddress is a free data retrieval call binding the contract method 0xb4701c09.
//
// Solidity: function networkProposalAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalSession) NetworkProposalAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.NetworkProposalAddress(&_NetworkWithdrawal.CallOpts)
}

// NetworkProposalAddress is a free data retrieval call binding the contract method 0xb4701c09.
//
// Solidity: function networkProposalAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) NetworkProposalAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.NetworkProposalAddress(&_NetworkWithdrawal.CallOpts)
}

// NextWithdrawalIndex is a free data retrieval call binding the contract method 0xbba9282e.
//
// Solidity: function nextWithdrawalIndex() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) NextWithdrawalIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "nextWithdrawalIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextWithdrawalIndex is a free data retrieval call binding the contract method 0xbba9282e.
//
// Solidity: function nextWithdrawalIndex() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) NextWithdrawalIndex() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.NextWithdrawalIndex(&_NetworkWithdrawal.CallOpts)
}

// NextWithdrawalIndex is a free data retrieval call binding the contract method 0xbba9282e.
//
// Solidity: function nextWithdrawalIndex() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) NextWithdrawalIndex() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.NextWithdrawalIndex(&_NetworkWithdrawal.CallOpts)
}

// NodeClaimEnabled is a free data retrieval call binding the contract method 0xd3638c7e.
//
// Solidity: function nodeClaimEnabled() view returns(bool)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) NodeClaimEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "nodeClaimEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// NodeClaimEnabled is a free data retrieval call binding the contract method 0xd3638c7e.
//
// Solidity: function nodeClaimEnabled() view returns(bool)
func (_NetworkWithdrawal *NetworkWithdrawalSession) NodeClaimEnabled() (bool, error) {
	return _NetworkWithdrawal.Contract.NodeClaimEnabled(&_NetworkWithdrawal.CallOpts)
}

// NodeClaimEnabled is a free data retrieval call binding the contract method 0xd3638c7e.
//
// Solidity: function nodeClaimEnabled() view returns(bool)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) NodeClaimEnabled() (bool, error) {
	return _NetworkWithdrawal.Contract.NodeClaimEnabled(&_NetworkWithdrawal.CallOpts)
}

// NodeCommissionRate is a free data retrieval call binding the contract method 0x4636e4e5.
//
// Solidity: function nodeCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) NodeCommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "nodeCommissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NodeCommissionRate is a free data retrieval call binding the contract method 0x4636e4e5.
//
// Solidity: function nodeCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) NodeCommissionRate() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.NodeCommissionRate(&_NetworkWithdrawal.CallOpts)
}

// NodeCommissionRate is a free data retrieval call binding the contract method 0x4636e4e5.
//
// Solidity: function nodeCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) NodeCommissionRate() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.NodeCommissionRate(&_NetworkWithdrawal.CallOpts)
}

// NodeRewardsFileCid is a free data retrieval call binding the contract method 0xd57dc824.
//
// Solidity: function nodeRewardsFileCid() view returns(string)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) NodeRewardsFileCid(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "nodeRewardsFileCid")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NodeRewardsFileCid is a free data retrieval call binding the contract method 0xd57dc824.
//
// Solidity: function nodeRewardsFileCid() view returns(string)
func (_NetworkWithdrawal *NetworkWithdrawalSession) NodeRewardsFileCid() (string, error) {
	return _NetworkWithdrawal.Contract.NodeRewardsFileCid(&_NetworkWithdrawal.CallOpts)
}

// NodeRewardsFileCid is a free data retrieval call binding the contract method 0xd57dc824.
//
// Solidity: function nodeRewardsFileCid() view returns(string)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) NodeRewardsFileCid() (string, error) {
	return _NetworkWithdrawal.Contract.NodeRewardsFileCid(&_NetworkWithdrawal.CallOpts)
}

// PlatformCommissionRate is a free data retrieval call binding the contract method 0x1da4dd0d.
//
// Solidity: function platformCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) PlatformCommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "platformCommissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PlatformCommissionRate is a free data retrieval call binding the contract method 0x1da4dd0d.
//
// Solidity: function platformCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) PlatformCommissionRate() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.PlatformCommissionRate(&_NetworkWithdrawal.CallOpts)
}

// PlatformCommissionRate is a free data retrieval call binding the contract method 0x1da4dd0d.
//
// Solidity: function platformCommissionRate() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) PlatformCommissionRate() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.PlatformCommissionRate(&_NetworkWithdrawal.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_NetworkWithdrawal *NetworkWithdrawalSession) ProxiableUUID() ([32]byte, error) {
	return _NetworkWithdrawal.Contract.ProxiableUUID(&_NetworkWithdrawal.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) ProxiableUUID() ([32]byte, error) {
	return _NetworkWithdrawal.Contract.ProxiableUUID(&_NetworkWithdrawal.CallOpts)
}

// TotalClaimedDepositOfNode is a free data retrieval call binding the contract method 0x6c570dc1.
//
// Solidity: function totalClaimedDepositOfNode(address ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) TotalClaimedDepositOfNode(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "totalClaimedDepositOfNode", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalClaimedDepositOfNode is a free data retrieval call binding the contract method 0x6c570dc1.
//
// Solidity: function totalClaimedDepositOfNode(address ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) TotalClaimedDepositOfNode(arg0 common.Address) (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalClaimedDepositOfNode(&_NetworkWithdrawal.CallOpts, arg0)
}

// TotalClaimedDepositOfNode is a free data retrieval call binding the contract method 0x6c570dc1.
//
// Solidity: function totalClaimedDepositOfNode(address ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) TotalClaimedDepositOfNode(arg0 common.Address) (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalClaimedDepositOfNode(&_NetworkWithdrawal.CallOpts, arg0)
}

// TotalClaimedRewardOfNode is a free data retrieval call binding the contract method 0xbb2d840c.
//
// Solidity: function totalClaimedRewardOfNode(address ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) TotalClaimedRewardOfNode(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "totalClaimedRewardOfNode", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalClaimedRewardOfNode is a free data retrieval call binding the contract method 0xbb2d840c.
//
// Solidity: function totalClaimedRewardOfNode(address ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) TotalClaimedRewardOfNode(arg0 common.Address) (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalClaimedRewardOfNode(&_NetworkWithdrawal.CallOpts, arg0)
}

// TotalClaimedRewardOfNode is a free data retrieval call binding the contract method 0xbb2d840c.
//
// Solidity: function totalClaimedRewardOfNode(address ) view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) TotalClaimedRewardOfNode(arg0 common.Address) (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalClaimedRewardOfNode(&_NetworkWithdrawal.CallOpts, arg0)
}

// TotalPlatformClaimedCommission is a free data retrieval call binding the contract method 0xe2f84d36.
//
// Solidity: function totalPlatformClaimedCommission() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) TotalPlatformClaimedCommission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "totalPlatformClaimedCommission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalPlatformClaimedCommission is a free data retrieval call binding the contract method 0xe2f84d36.
//
// Solidity: function totalPlatformClaimedCommission() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) TotalPlatformClaimedCommission() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalPlatformClaimedCommission(&_NetworkWithdrawal.CallOpts)
}

// TotalPlatformClaimedCommission is a free data retrieval call binding the contract method 0xe2f84d36.
//
// Solidity: function totalPlatformClaimedCommission() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) TotalPlatformClaimedCommission() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalPlatformClaimedCommission(&_NetworkWithdrawal.CallOpts)
}

// TotalPlatformCommission is a free data retrieval call binding the contract method 0xfef25c0d.
//
// Solidity: function totalPlatformCommission() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) TotalPlatformCommission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "totalPlatformCommission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalPlatformCommission is a free data retrieval call binding the contract method 0xfef25c0d.
//
// Solidity: function totalPlatformCommission() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) TotalPlatformCommission() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalPlatformCommission(&_NetworkWithdrawal.CallOpts)
}

// TotalPlatformCommission is a free data retrieval call binding the contract method 0xfef25c0d.
//
// Solidity: function totalPlatformCommission() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) TotalPlatformCommission() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalPlatformCommission(&_NetworkWithdrawal.CallOpts)
}

// TotalWithdrawalShortages is a free data retrieval call binding the contract method 0xd69d8431.
//
// Solidity: function totalWithdrawalShortages() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) TotalWithdrawalShortages(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "totalWithdrawalShortages")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalWithdrawalShortages is a free data retrieval call binding the contract method 0xd69d8431.
//
// Solidity: function totalWithdrawalShortages() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) TotalWithdrawalShortages() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalWithdrawalShortages(&_NetworkWithdrawal.CallOpts)
}

// TotalWithdrawalShortages is a free data retrieval call binding the contract method 0xd69d8431.
//
// Solidity: function totalWithdrawalShortages() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) TotalWithdrawalShortages() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.TotalWithdrawalShortages(&_NetworkWithdrawal.CallOpts)
}

// UserDepositAddress is a free data retrieval call binding the contract method 0x46773830.
//
// Solidity: function userDepositAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) UserDepositAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "userDepositAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserDepositAddress is a free data retrieval call binding the contract method 0x46773830.
//
// Solidity: function userDepositAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalSession) UserDepositAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.UserDepositAddress(&_NetworkWithdrawal.CallOpts)
}

// UserDepositAddress is a free data retrieval call binding the contract method 0x46773830.
//
// Solidity: function userDepositAddress() view returns(address)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) UserDepositAddress() (common.Address, error) {
	return _NetworkWithdrawal.Contract.UserDepositAddress(&_NetworkWithdrawal.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) Version(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_NetworkWithdrawal *NetworkWithdrawalSession) Version() (uint8, error) {
	return _NetworkWithdrawal.Contract.Version(&_NetworkWithdrawal.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) Version() (uint8, error) {
	return _NetworkWithdrawal.Contract.Version(&_NetworkWithdrawal.CallOpts)
}

// WithdrawalAtIndex is a free data retrieval call binding the contract method 0xa8e1b8ef.
//
// Solidity: function withdrawalAtIndex(uint256 ) view returns(address _address, uint256 _amount)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) WithdrawalAtIndex(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Address common.Address
	Amount  *big.Int
}, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "withdrawalAtIndex", arg0)

	outstruct := new(struct {
		Address common.Address
		Amount  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Address = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// WithdrawalAtIndex is a free data retrieval call binding the contract method 0xa8e1b8ef.
//
// Solidity: function withdrawalAtIndex(uint256 ) view returns(address _address, uint256 _amount)
func (_NetworkWithdrawal *NetworkWithdrawalSession) WithdrawalAtIndex(arg0 *big.Int) (struct {
	Address common.Address
	Amount  *big.Int
}, error) {
	return _NetworkWithdrawal.Contract.WithdrawalAtIndex(&_NetworkWithdrawal.CallOpts, arg0)
}

// WithdrawalAtIndex is a free data retrieval call binding the contract method 0xa8e1b8ef.
//
// Solidity: function withdrawalAtIndex(uint256 ) view returns(address _address, uint256 _amount)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) WithdrawalAtIndex(arg0 *big.Int) (struct {
	Address common.Address
	Amount  *big.Int
}, error) {
	return _NetworkWithdrawal.Contract.WithdrawalAtIndex(&_NetworkWithdrawal.CallOpts, arg0)
}

// WithdrawalCycleSeconds is a free data retrieval call binding the contract method 0x9ca9f888.
//
// Solidity: function withdrawalCycleSeconds() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCaller) WithdrawalCycleSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkWithdrawal.contract.Call(opts, &out, "withdrawalCycleSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalCycleSeconds is a free data retrieval call binding the contract method 0x9ca9f888.
//
// Solidity: function withdrawalCycleSeconds() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalSession) WithdrawalCycleSeconds() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.WithdrawalCycleSeconds(&_NetworkWithdrawal.CallOpts)
}

// WithdrawalCycleSeconds is a free data retrieval call binding the contract method 0x9ca9f888.
//
// Solidity: function withdrawalCycleSeconds() view returns(uint256)
func (_NetworkWithdrawal *NetworkWithdrawalCallerSession) WithdrawalCycleSeconds() (*big.Int, error) {
	return _NetworkWithdrawal.Contract.WithdrawalCycleSeconds(&_NetworkWithdrawal.CallOpts)
}

// DepositEth is a paid mutator transaction binding the contract method 0x439370b1.
//
// Solidity: function depositEth() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) DepositEth(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "depositEth")
}

// DepositEth is a paid mutator transaction binding the contract method 0x439370b1.
//
// Solidity: function depositEth() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) DepositEth() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.DepositEth(&_NetworkWithdrawal.TransactOpts)
}

// DepositEth is a paid mutator transaction binding the contract method 0x439370b1.
//
// Solidity: function depositEth() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) DepositEth() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.DepositEth(&_NetworkWithdrawal.TransactOpts)
}

// DepositEthAndUpdateTotalShortages is a paid mutator transaction binding the contract method 0xa59914c6.
//
// Solidity: function depositEthAndUpdateTotalShortages() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) DepositEthAndUpdateTotalShortages(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "depositEthAndUpdateTotalShortages")
}

// DepositEthAndUpdateTotalShortages is a paid mutator transaction binding the contract method 0xa59914c6.
//
// Solidity: function depositEthAndUpdateTotalShortages() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) DepositEthAndUpdateTotalShortages() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.DepositEthAndUpdateTotalShortages(&_NetworkWithdrawal.TransactOpts)
}

// DepositEthAndUpdateTotalShortages is a paid mutator transaction binding the contract method 0xa59914c6.
//
// Solidity: function depositEthAndUpdateTotalShortages() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) DepositEthAndUpdateTotalShortages() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.DepositEthAndUpdateTotalShortages(&_NetworkWithdrawal.TransactOpts)
}

// Distribute is a paid mutator transaction binding the contract method 0xc980ba89.
//
// Solidity: function distribute(uint8 _distributionType, uint256 _dealtHeight, uint256 _userAmount, uint256 _nodeAmount, uint256 _platformAmount, uint256 _maxClaimableWithdrawalIndex) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) Distribute(opts *bind.TransactOpts, _distributionType uint8, _dealtHeight *big.Int, _userAmount *big.Int, _nodeAmount *big.Int, _platformAmount *big.Int, _maxClaimableWithdrawalIndex *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "distribute", _distributionType, _dealtHeight, _userAmount, _nodeAmount, _platformAmount, _maxClaimableWithdrawalIndex)
}

// Distribute is a paid mutator transaction binding the contract method 0xc980ba89.
//
// Solidity: function distribute(uint8 _distributionType, uint256 _dealtHeight, uint256 _userAmount, uint256 _nodeAmount, uint256 _platformAmount, uint256 _maxClaimableWithdrawalIndex) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) Distribute(_distributionType uint8, _dealtHeight *big.Int, _userAmount *big.Int, _nodeAmount *big.Int, _platformAmount *big.Int, _maxClaimableWithdrawalIndex *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Distribute(&_NetworkWithdrawal.TransactOpts, _distributionType, _dealtHeight, _userAmount, _nodeAmount, _platformAmount, _maxClaimableWithdrawalIndex)
}

// Distribute is a paid mutator transaction binding the contract method 0xc980ba89.
//
// Solidity: function distribute(uint8 _distributionType, uint256 _dealtHeight, uint256 _userAmount, uint256 _nodeAmount, uint256 _platformAmount, uint256 _maxClaimableWithdrawalIndex) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) Distribute(_distributionType uint8, _dealtHeight *big.Int, _userAmount *big.Int, _nodeAmount *big.Int, _platformAmount *big.Int, _maxClaimableWithdrawalIndex *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Distribute(&_NetworkWithdrawal.TransactOpts, _distributionType, _dealtHeight, _userAmount, _nodeAmount, _platformAmount, _maxClaimableWithdrawalIndex)
}

// Init is a paid mutator transaction binding the contract method 0x99e133f9.
//
// Solidity: function init(address _lsdTokenAddress, address _userDepositAddress, address _networkProposalAddress, address _networkBalancesAddress, address _feePoolAddress, address _factoryAddress) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) Init(opts *bind.TransactOpts, _lsdTokenAddress common.Address, _userDepositAddress common.Address, _networkProposalAddress common.Address, _networkBalancesAddress common.Address, _feePoolAddress common.Address, _factoryAddress common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "init", _lsdTokenAddress, _userDepositAddress, _networkProposalAddress, _networkBalancesAddress, _feePoolAddress, _factoryAddress)
}

// Init is a paid mutator transaction binding the contract method 0x99e133f9.
//
// Solidity: function init(address _lsdTokenAddress, address _userDepositAddress, address _networkProposalAddress, address _networkBalancesAddress, address _feePoolAddress, address _factoryAddress) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) Init(_lsdTokenAddress common.Address, _userDepositAddress common.Address, _networkProposalAddress common.Address, _networkBalancesAddress common.Address, _feePoolAddress common.Address, _factoryAddress common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Init(&_NetworkWithdrawal.TransactOpts, _lsdTokenAddress, _userDepositAddress, _networkProposalAddress, _networkBalancesAddress, _feePoolAddress, _factoryAddress)
}

// Init is a paid mutator transaction binding the contract method 0x99e133f9.
//
// Solidity: function init(address _lsdTokenAddress, address _userDepositAddress, address _networkProposalAddress, address _networkBalancesAddress, address _feePoolAddress, address _factoryAddress) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) Init(_lsdTokenAddress common.Address, _userDepositAddress common.Address, _networkProposalAddress common.Address, _networkBalancesAddress common.Address, _feePoolAddress common.Address, _factoryAddress common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Init(&_NetworkWithdrawal.TransactOpts, _lsdTokenAddress, _userDepositAddress, _networkProposalAddress, _networkBalancesAddress, _feePoolAddress, _factoryAddress)
}

// NodeClaim is a paid mutator transaction binding the contract method 0xfdf435e9.
//
// Solidity: function nodeClaim(uint256 _index, address _node, uint256 _totalRewardAmount, uint256 _totalExitDepositAmount, bytes32[] _merkleProof, uint8 _claimType) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) NodeClaim(opts *bind.TransactOpts, _index *big.Int, _node common.Address, _totalRewardAmount *big.Int, _totalExitDepositAmount *big.Int, _merkleProof [][32]byte, _claimType uint8) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "nodeClaim", _index, _node, _totalRewardAmount, _totalExitDepositAmount, _merkleProof, _claimType)
}

// NodeClaim is a paid mutator transaction binding the contract method 0xfdf435e9.
//
// Solidity: function nodeClaim(uint256 _index, address _node, uint256 _totalRewardAmount, uint256 _totalExitDepositAmount, bytes32[] _merkleProof, uint8 _claimType) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) NodeClaim(_index *big.Int, _node common.Address, _totalRewardAmount *big.Int, _totalExitDepositAmount *big.Int, _merkleProof [][32]byte, _claimType uint8) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.NodeClaim(&_NetworkWithdrawal.TransactOpts, _index, _node, _totalRewardAmount, _totalExitDepositAmount, _merkleProof, _claimType)
}

// NodeClaim is a paid mutator transaction binding the contract method 0xfdf435e9.
//
// Solidity: function nodeClaim(uint256 _index, address _node, uint256 _totalRewardAmount, uint256 _totalExitDepositAmount, bytes32[] _merkleProof, uint8 _claimType) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) NodeClaim(_index *big.Int, _node common.Address, _totalRewardAmount *big.Int, _totalExitDepositAmount *big.Int, _merkleProof [][32]byte, _claimType uint8) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.NodeClaim(&_NetworkWithdrawal.TransactOpts, _index, _node, _totalRewardAmount, _totalExitDepositAmount, _merkleProof, _claimType)
}

// NotifyValidatorExit is a paid mutator transaction binding the contract method 0x1e0f4aae.
//
// Solidity: function notifyValidatorExit(uint256 _withdrawCycle, uint256 _ejectedStartCycle, uint256[] _validatorIndexList) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) NotifyValidatorExit(opts *bind.TransactOpts, _withdrawCycle *big.Int, _ejectedStartCycle *big.Int, _validatorIndexList []*big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "notifyValidatorExit", _withdrawCycle, _ejectedStartCycle, _validatorIndexList)
}

// NotifyValidatorExit is a paid mutator transaction binding the contract method 0x1e0f4aae.
//
// Solidity: function notifyValidatorExit(uint256 _withdrawCycle, uint256 _ejectedStartCycle, uint256[] _validatorIndexList) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) NotifyValidatorExit(_withdrawCycle *big.Int, _ejectedStartCycle *big.Int, _validatorIndexList []*big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.NotifyValidatorExit(&_NetworkWithdrawal.TransactOpts, _withdrawCycle, _ejectedStartCycle, _validatorIndexList)
}

// NotifyValidatorExit is a paid mutator transaction binding the contract method 0x1e0f4aae.
//
// Solidity: function notifyValidatorExit(uint256 _withdrawCycle, uint256 _ejectedStartCycle, uint256[] _validatorIndexList) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) NotifyValidatorExit(_withdrawCycle *big.Int, _ejectedStartCycle *big.Int, _validatorIndexList []*big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.NotifyValidatorExit(&_NetworkWithdrawal.TransactOpts, _withdrawCycle, _ejectedStartCycle, _validatorIndexList)
}

// PlatformClaim is a paid mutator transaction binding the contract method 0xaaf82770.
//
// Solidity: function platformClaim(address _recipient) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) PlatformClaim(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "platformClaim", _recipient)
}

// PlatformClaim is a paid mutator transaction binding the contract method 0xaaf82770.
//
// Solidity: function platformClaim(address _recipient) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) PlatformClaim(_recipient common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.PlatformClaim(&_NetworkWithdrawal.TransactOpts, _recipient)
}

// PlatformClaim is a paid mutator transaction binding the contract method 0xaaf82770.
//
// Solidity: function platformClaim(address _recipient) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) PlatformClaim(_recipient common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.PlatformClaim(&_NetworkWithdrawal.TransactOpts, _recipient)
}

// Reinit is a paid mutator transaction binding the contract method 0xc482ceaf.
//
// Solidity: function reinit() returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) Reinit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "reinit")
}

// Reinit is a paid mutator transaction binding the contract method 0xc482ceaf.
//
// Solidity: function reinit() returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) Reinit() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Reinit(&_NetworkWithdrawal.TransactOpts)
}

// Reinit is a paid mutator transaction binding the contract method 0xc482ceaf.
//
// Solidity: function reinit() returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) Reinit() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Reinit(&_NetworkWithdrawal.TransactOpts)
}

// SetFactoryCommissionRate is a paid mutator transaction binding the contract method 0xfd005077.
//
// Solidity: function setFactoryCommissionRate(uint256 _factoryCommissionRate) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) SetFactoryCommissionRate(opts *bind.TransactOpts, _factoryCommissionRate *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "setFactoryCommissionRate", _factoryCommissionRate)
}

// SetFactoryCommissionRate is a paid mutator transaction binding the contract method 0xfd005077.
//
// Solidity: function setFactoryCommissionRate(uint256 _factoryCommissionRate) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) SetFactoryCommissionRate(_factoryCommissionRate *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetFactoryCommissionRate(&_NetworkWithdrawal.TransactOpts, _factoryCommissionRate)
}

// SetFactoryCommissionRate is a paid mutator transaction binding the contract method 0xfd005077.
//
// Solidity: function setFactoryCommissionRate(uint256 _factoryCommissionRate) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) SetFactoryCommissionRate(_factoryCommissionRate *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetFactoryCommissionRate(&_NetworkWithdrawal.TransactOpts, _factoryCommissionRate)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x12b81931.
//
// Solidity: function setMerkleRoot(uint256 _dealtEpoch, bytes32 _merkleRoot, string _nodeRewardsFileCid) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) SetMerkleRoot(opts *bind.TransactOpts, _dealtEpoch *big.Int, _merkleRoot [32]byte, _nodeRewardsFileCid string) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "setMerkleRoot", _dealtEpoch, _merkleRoot, _nodeRewardsFileCid)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x12b81931.
//
// Solidity: function setMerkleRoot(uint256 _dealtEpoch, bytes32 _merkleRoot, string _nodeRewardsFileCid) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) SetMerkleRoot(_dealtEpoch *big.Int, _merkleRoot [32]byte, _nodeRewardsFileCid string) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetMerkleRoot(&_NetworkWithdrawal.TransactOpts, _dealtEpoch, _merkleRoot, _nodeRewardsFileCid)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x12b81931.
//
// Solidity: function setMerkleRoot(uint256 _dealtEpoch, bytes32 _merkleRoot, string _nodeRewardsFileCid) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) SetMerkleRoot(_dealtEpoch *big.Int, _merkleRoot [32]byte, _nodeRewardsFileCid string) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetMerkleRoot(&_NetworkWithdrawal.TransactOpts, _dealtEpoch, _merkleRoot, _nodeRewardsFileCid)
}

// SetNodeClaimEnabled is a paid mutator transaction binding the contract method 0xf1583c08.
//
// Solidity: function setNodeClaimEnabled(bool _value) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) SetNodeClaimEnabled(opts *bind.TransactOpts, _value bool) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "setNodeClaimEnabled", _value)
}

// SetNodeClaimEnabled is a paid mutator transaction binding the contract method 0xf1583c08.
//
// Solidity: function setNodeClaimEnabled(bool _value) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) SetNodeClaimEnabled(_value bool) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetNodeClaimEnabled(&_NetworkWithdrawal.TransactOpts, _value)
}

// SetNodeClaimEnabled is a paid mutator transaction binding the contract method 0xf1583c08.
//
// Solidity: function setNodeClaimEnabled(bool _value) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) SetNodeClaimEnabled(_value bool) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetNodeClaimEnabled(&_NetworkWithdrawal.TransactOpts, _value)
}

// SetPlatformAndNodeCommissionRate is a paid mutator transaction binding the contract method 0x7a1a934d.
//
// Solidity: function setPlatformAndNodeCommissionRate(uint256 _platformCommissionRate, uint256 _nodeCommissionRate) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) SetPlatformAndNodeCommissionRate(opts *bind.TransactOpts, _platformCommissionRate *big.Int, _nodeCommissionRate *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "setPlatformAndNodeCommissionRate", _platformCommissionRate, _nodeCommissionRate)
}

// SetPlatformAndNodeCommissionRate is a paid mutator transaction binding the contract method 0x7a1a934d.
//
// Solidity: function setPlatformAndNodeCommissionRate(uint256 _platformCommissionRate, uint256 _nodeCommissionRate) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) SetPlatformAndNodeCommissionRate(_platformCommissionRate *big.Int, _nodeCommissionRate *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetPlatformAndNodeCommissionRate(&_NetworkWithdrawal.TransactOpts, _platformCommissionRate, _nodeCommissionRate)
}

// SetPlatformAndNodeCommissionRate is a paid mutator transaction binding the contract method 0x7a1a934d.
//
// Solidity: function setPlatformAndNodeCommissionRate(uint256 _platformCommissionRate, uint256 _nodeCommissionRate) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) SetPlatformAndNodeCommissionRate(_platformCommissionRate *big.Int, _nodeCommissionRate *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetPlatformAndNodeCommissionRate(&_NetworkWithdrawal.TransactOpts, _platformCommissionRate, _nodeCommissionRate)
}

// SetWithdrawalCycleSeconds is a paid mutator transaction binding the contract method 0xb9a65a6c.
//
// Solidity: function setWithdrawalCycleSeconds(uint256 _withdrawalCycleSeconds) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) SetWithdrawalCycleSeconds(opts *bind.TransactOpts, _withdrawalCycleSeconds *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "setWithdrawalCycleSeconds", _withdrawalCycleSeconds)
}

// SetWithdrawalCycleSeconds is a paid mutator transaction binding the contract method 0xb9a65a6c.
//
// Solidity: function setWithdrawalCycleSeconds(uint256 _withdrawalCycleSeconds) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) SetWithdrawalCycleSeconds(_withdrawalCycleSeconds *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetWithdrawalCycleSeconds(&_NetworkWithdrawal.TransactOpts, _withdrawalCycleSeconds)
}

// SetWithdrawalCycleSeconds is a paid mutator transaction binding the contract method 0xb9a65a6c.
//
// Solidity: function setWithdrawalCycleSeconds(uint256 _withdrawalCycleSeconds) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) SetWithdrawalCycleSeconds(_withdrawalCycleSeconds *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.SetWithdrawalCycleSeconds(&_NetworkWithdrawal.TransactOpts, _withdrawalCycleSeconds)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _lsdTokenAmount) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) Unstake(opts *bind.TransactOpts, _lsdTokenAmount *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "unstake", _lsdTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _lsdTokenAmount) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) Unstake(_lsdTokenAmount *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Unstake(&_NetworkWithdrawal.TransactOpts, _lsdTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _lsdTokenAmount) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) Unstake(_lsdTokenAmount *big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Unstake(&_NetworkWithdrawal.TransactOpts, _lsdTokenAmount)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.UpgradeTo(&_NetworkWithdrawal.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.UpgradeTo(&_NetworkWithdrawal.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.UpgradeToAndCall(&_NetworkWithdrawal.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.UpgradeToAndCall(&_NetworkWithdrawal.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x983d95ce.
//
// Solidity: function withdraw(uint256[] _withdrawalIndexList) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) Withdraw(opts *bind.TransactOpts, _withdrawalIndexList []*big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.Transact(opts, "withdraw", _withdrawalIndexList)
}

// Withdraw is a paid mutator transaction binding the contract method 0x983d95ce.
//
// Solidity: function withdraw(uint256[] _withdrawalIndexList) returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) Withdraw(_withdrawalIndexList []*big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Withdraw(&_NetworkWithdrawal.TransactOpts, _withdrawalIndexList)
}

// Withdraw is a paid mutator transaction binding the contract method 0x983d95ce.
//
// Solidity: function withdraw(uint256[] _withdrawalIndexList) returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) Withdraw(_withdrawalIndexList []*big.Int) (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Withdraw(&_NetworkWithdrawal.TransactOpts, _withdrawalIndexList)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkWithdrawal.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalSession) Receive() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Receive(&_NetworkWithdrawal.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkWithdrawal *NetworkWithdrawalTransactorSession) Receive() (*types.Transaction, error) {
	return _NetworkWithdrawal.Contract.Receive(&_NetworkWithdrawal.TransactOpts)
}

// NetworkWithdrawalAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalAdminChangedIterator struct {
	Event *NetworkWithdrawalAdminChanged // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalAdminChanged)
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
		it.Event = new(NetworkWithdrawalAdminChanged)
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
func (it *NetworkWithdrawalAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalAdminChanged represents a AdminChanged event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*NetworkWithdrawalAdminChangedIterator, error) {

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalAdminChangedIterator{contract: _NetworkWithdrawal.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalAdminChanged) (event.Subscription, error) {

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalAdminChanged)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseAdminChanged(log types.Log) (*NetworkWithdrawalAdminChanged, error) {
	event := new(NetworkWithdrawalAdminChanged)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalBeaconUpgradedIterator struct {
	Event *NetworkWithdrawalBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalBeaconUpgraded)
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
		it.Event = new(NetworkWithdrawalBeaconUpgraded)
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
func (it *NetworkWithdrawalBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalBeaconUpgraded represents a BeaconUpgraded event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*NetworkWithdrawalBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalBeaconUpgradedIterator{contract: _NetworkWithdrawal.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalBeaconUpgraded)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseBeaconUpgraded(log types.Log) (*NetworkWithdrawalBeaconUpgraded, error) {
	event := new(NetworkWithdrawalBeaconUpgraded)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalDistributeRewardsIterator is returned from FilterDistributeRewards and is used to iterate over the raw logs and unpacked data for DistributeRewards events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalDistributeRewardsIterator struct {
	Event *NetworkWithdrawalDistributeRewards // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalDistributeRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalDistributeRewards)
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
		it.Event = new(NetworkWithdrawalDistributeRewards)
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
func (it *NetworkWithdrawalDistributeRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalDistributeRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalDistributeRewards represents a DistributeRewards event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalDistributeRewards struct {
	DistributeType            uint8
	DealtHeight               *big.Int
	UserAmount                *big.Int
	NodeAmount                *big.Int
	PlatformAmount            *big.Int
	MaxClaimableWithdrawIndex *big.Int
	MvAmount                  *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterDistributeRewards is a free log retrieval operation binding the contract event 0xf10021cf129ec9c5003084ae330dba6d0bd143c47a2677c4d68939a19423206b.
//
// Solidity: event DistributeRewards(uint8 distributeType, uint256 dealtHeight, uint256 userAmount, uint256 nodeAmount, uint256 platformAmount, uint256 maxClaimableWithdrawIndex, uint256 mvAmount)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterDistributeRewards(opts *bind.FilterOpts) (*NetworkWithdrawalDistributeRewardsIterator, error) {

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "DistributeRewards")
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalDistributeRewardsIterator{contract: _NetworkWithdrawal.contract, event: "DistributeRewards", logs: logs, sub: sub}, nil
}

// WatchDistributeRewards is a free log subscription operation binding the contract event 0xf10021cf129ec9c5003084ae330dba6d0bd143c47a2677c4d68939a19423206b.
//
// Solidity: event DistributeRewards(uint8 distributeType, uint256 dealtHeight, uint256 userAmount, uint256 nodeAmount, uint256 platformAmount, uint256 maxClaimableWithdrawIndex, uint256 mvAmount)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchDistributeRewards(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalDistributeRewards) (event.Subscription, error) {

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "DistributeRewards")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalDistributeRewards)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "DistributeRewards", log); err != nil {
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

// ParseDistributeRewards is a log parse operation binding the contract event 0xf10021cf129ec9c5003084ae330dba6d0bd143c47a2677c4d68939a19423206b.
//
// Solidity: event DistributeRewards(uint8 distributeType, uint256 dealtHeight, uint256 userAmount, uint256 nodeAmount, uint256 platformAmount, uint256 maxClaimableWithdrawIndex, uint256 mvAmount)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseDistributeRewards(log types.Log) (*NetworkWithdrawalDistributeRewards, error) {
	event := new(NetworkWithdrawalDistributeRewards)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "DistributeRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalEtherDepositedIterator is returned from FilterEtherDeposited and is used to iterate over the raw logs and unpacked data for EtherDeposited events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalEtherDepositedIterator struct {
	Event *NetworkWithdrawalEtherDeposited // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalEtherDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalEtherDeposited)
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
		it.Event = new(NetworkWithdrawalEtherDeposited)
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
func (it *NetworkWithdrawalEtherDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalEtherDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalEtherDeposited represents a EtherDeposited event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalEtherDeposited struct {
	From   common.Address
	Amount *big.Int
	Time   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherDeposited is a free log retrieval operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterEtherDeposited(opts *bind.FilterOpts, from []common.Address) (*NetworkWithdrawalEtherDepositedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "EtherDeposited", fromRule)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalEtherDepositedIterator{contract: _NetworkWithdrawal.contract, event: "EtherDeposited", logs: logs, sub: sub}, nil
}

// WatchEtherDeposited is a free log subscription operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchEtherDeposited(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalEtherDeposited, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "EtherDeposited", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalEtherDeposited)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "EtherDeposited", log); err != nil {
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

// ParseEtherDeposited is a log parse operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseEtherDeposited(log types.Log) (*NetworkWithdrawalEtherDeposited, error) {
	event := new(NetworkWithdrawalEtherDeposited)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "EtherDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalInitializedIterator struct {
	Event *NetworkWithdrawalInitialized // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalInitialized)
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
		it.Event = new(NetworkWithdrawalInitialized)
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
func (it *NetworkWithdrawalInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalInitialized represents a Initialized event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterInitialized(opts *bind.FilterOpts) (*NetworkWithdrawalInitializedIterator, error) {

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalInitializedIterator{contract: _NetworkWithdrawal.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalInitialized) (event.Subscription, error) {

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalInitialized)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseInitialized(log types.Log) (*NetworkWithdrawalInitialized, error) {
	event := new(NetworkWithdrawalInitialized)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalNodeClaimedIterator is returned from FilterNodeClaimed and is used to iterate over the raw logs and unpacked data for NodeClaimed events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalNodeClaimedIterator struct {
	Event *NetworkWithdrawalNodeClaimed // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalNodeClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalNodeClaimed)
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
		it.Event = new(NetworkWithdrawalNodeClaimed)
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
func (it *NetworkWithdrawalNodeClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalNodeClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalNodeClaimed represents a NodeClaimed event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalNodeClaimed struct {
	Index            *big.Int
	Account          common.Address
	ClaimableReward  *big.Int
	ClaimableDeposit *big.Int
	ClaimType        uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNodeClaimed is a free log retrieval operation binding the contract event 0x659e842f0209726671f562e8d7ee3a25d2cef92c2b85dcd268af93ef5d13582c.
//
// Solidity: event NodeClaimed(uint256 index, address account, uint256 claimableReward, uint256 claimableDeposit, uint8 claimType)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterNodeClaimed(opts *bind.FilterOpts) (*NetworkWithdrawalNodeClaimedIterator, error) {

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "NodeClaimed")
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalNodeClaimedIterator{contract: _NetworkWithdrawal.contract, event: "NodeClaimed", logs: logs, sub: sub}, nil
}

// WatchNodeClaimed is a free log subscription operation binding the contract event 0x659e842f0209726671f562e8d7ee3a25d2cef92c2b85dcd268af93ef5d13582c.
//
// Solidity: event NodeClaimed(uint256 index, address account, uint256 claimableReward, uint256 claimableDeposit, uint8 claimType)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchNodeClaimed(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalNodeClaimed) (event.Subscription, error) {

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "NodeClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalNodeClaimed)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "NodeClaimed", log); err != nil {
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

// ParseNodeClaimed is a log parse operation binding the contract event 0x659e842f0209726671f562e8d7ee3a25d2cef92c2b85dcd268af93ef5d13582c.
//
// Solidity: event NodeClaimed(uint256 index, address account, uint256 claimableReward, uint256 claimableDeposit, uint8 claimType)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseNodeClaimed(log types.Log) (*NetworkWithdrawalNodeClaimed, error) {
	event := new(NetworkWithdrawalNodeClaimed)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "NodeClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalNotifyValidatorExitIterator is returned from FilterNotifyValidatorExit and is used to iterate over the raw logs and unpacked data for NotifyValidatorExit events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalNotifyValidatorExitIterator struct {
	Event *NetworkWithdrawalNotifyValidatorExit // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalNotifyValidatorExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalNotifyValidatorExit)
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
		it.Event = new(NetworkWithdrawalNotifyValidatorExit)
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
func (it *NetworkWithdrawalNotifyValidatorExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalNotifyValidatorExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalNotifyValidatorExit represents a NotifyValidatorExit event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalNotifyValidatorExit struct {
	WithdrawalCycle             *big.Int
	EjectedStartWithdrawalCycle *big.Int
	EjectedValidators           []*big.Int
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterNotifyValidatorExit is a free log retrieval operation binding the contract event 0xb83477449e27b4bab4f28c938d033b953557d6a1b9b4469a43d229f78ed6e55c.
//
// Solidity: event NotifyValidatorExit(uint256 withdrawalCycle, uint256 ejectedStartWithdrawalCycle, uint256[] ejectedValidators)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterNotifyValidatorExit(opts *bind.FilterOpts) (*NetworkWithdrawalNotifyValidatorExitIterator, error) {

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "NotifyValidatorExit")
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalNotifyValidatorExitIterator{contract: _NetworkWithdrawal.contract, event: "NotifyValidatorExit", logs: logs, sub: sub}, nil
}

// WatchNotifyValidatorExit is a free log subscription operation binding the contract event 0xb83477449e27b4bab4f28c938d033b953557d6a1b9b4469a43d229f78ed6e55c.
//
// Solidity: event NotifyValidatorExit(uint256 withdrawalCycle, uint256 ejectedStartWithdrawalCycle, uint256[] ejectedValidators)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchNotifyValidatorExit(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalNotifyValidatorExit) (event.Subscription, error) {

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "NotifyValidatorExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalNotifyValidatorExit)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "NotifyValidatorExit", log); err != nil {
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

// ParseNotifyValidatorExit is a log parse operation binding the contract event 0xb83477449e27b4bab4f28c938d033b953557d6a1b9b4469a43d229f78ed6e55c.
//
// Solidity: event NotifyValidatorExit(uint256 withdrawalCycle, uint256 ejectedStartWithdrawalCycle, uint256[] ejectedValidators)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseNotifyValidatorExit(log types.Log) (*NetworkWithdrawalNotifyValidatorExit, error) {
	event := new(NetworkWithdrawalNotifyValidatorExit)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "NotifyValidatorExit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalSetMerkleRootIterator is returned from FilterSetMerkleRoot and is used to iterate over the raw logs and unpacked data for SetMerkleRoot events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalSetMerkleRootIterator struct {
	Event *NetworkWithdrawalSetMerkleRoot // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalSetMerkleRootIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalSetMerkleRoot)
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
		it.Event = new(NetworkWithdrawalSetMerkleRoot)
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
func (it *NetworkWithdrawalSetMerkleRootIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalSetMerkleRootIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalSetMerkleRoot represents a SetMerkleRoot event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalSetMerkleRoot struct {
	DealtEpoch         *big.Int
	MerkleRoot         [32]byte
	NodeRewardsFileCid string
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterSetMerkleRoot is a free log retrieval operation binding the contract event 0xec43b2424d0504da473794ad49016df3e06fb0d772bb403d724c9e5d53d8afb8.
//
// Solidity: event SetMerkleRoot(uint256 indexed dealtEpoch, bytes32 merkleRoot, string nodeRewardsFileCid)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterSetMerkleRoot(opts *bind.FilterOpts, dealtEpoch []*big.Int) (*NetworkWithdrawalSetMerkleRootIterator, error) {

	var dealtEpochRule []interface{}
	for _, dealtEpochItem := range dealtEpoch {
		dealtEpochRule = append(dealtEpochRule, dealtEpochItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "SetMerkleRoot", dealtEpochRule)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalSetMerkleRootIterator{contract: _NetworkWithdrawal.contract, event: "SetMerkleRoot", logs: logs, sub: sub}, nil
}

// WatchSetMerkleRoot is a free log subscription operation binding the contract event 0xec43b2424d0504da473794ad49016df3e06fb0d772bb403d724c9e5d53d8afb8.
//
// Solidity: event SetMerkleRoot(uint256 indexed dealtEpoch, bytes32 merkleRoot, string nodeRewardsFileCid)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchSetMerkleRoot(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalSetMerkleRoot, dealtEpoch []*big.Int) (event.Subscription, error) {

	var dealtEpochRule []interface{}
	for _, dealtEpochItem := range dealtEpoch {
		dealtEpochRule = append(dealtEpochRule, dealtEpochItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "SetMerkleRoot", dealtEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalSetMerkleRoot)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "SetMerkleRoot", log); err != nil {
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

// ParseSetMerkleRoot is a log parse operation binding the contract event 0xec43b2424d0504da473794ad49016df3e06fb0d772bb403d724c9e5d53d8afb8.
//
// Solidity: event SetMerkleRoot(uint256 indexed dealtEpoch, bytes32 merkleRoot, string nodeRewardsFileCid)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseSetMerkleRoot(log types.Log) (*NetworkWithdrawalSetMerkleRoot, error) {
	event := new(NetworkWithdrawalSetMerkleRoot)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "SetMerkleRoot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalSetWithdrawalCycleSecondsIterator is returned from FilterSetWithdrawalCycleSeconds and is used to iterate over the raw logs and unpacked data for SetWithdrawalCycleSeconds events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalSetWithdrawalCycleSecondsIterator struct {
	Event *NetworkWithdrawalSetWithdrawalCycleSeconds // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalSetWithdrawalCycleSecondsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalSetWithdrawalCycleSeconds)
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
		it.Event = new(NetworkWithdrawalSetWithdrawalCycleSeconds)
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
func (it *NetworkWithdrawalSetWithdrawalCycleSecondsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalSetWithdrawalCycleSecondsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalSetWithdrawalCycleSeconds represents a SetWithdrawalCycleSeconds event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalSetWithdrawalCycleSeconds struct {
	CycleSeconds *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSetWithdrawalCycleSeconds is a free log retrieval operation binding the contract event 0x104b9e6ae581a335baf19e149870062416733a2a2fa10ed5c2bdb82bf4267478.
//
// Solidity: event SetWithdrawalCycleSeconds(uint256 cycleSeconds)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterSetWithdrawalCycleSeconds(opts *bind.FilterOpts) (*NetworkWithdrawalSetWithdrawalCycleSecondsIterator, error) {

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "SetWithdrawalCycleSeconds")
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalSetWithdrawalCycleSecondsIterator{contract: _NetworkWithdrawal.contract, event: "SetWithdrawalCycleSeconds", logs: logs, sub: sub}, nil
}

// WatchSetWithdrawalCycleSeconds is a free log subscription operation binding the contract event 0x104b9e6ae581a335baf19e149870062416733a2a2fa10ed5c2bdb82bf4267478.
//
// Solidity: event SetWithdrawalCycleSeconds(uint256 cycleSeconds)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchSetWithdrawalCycleSeconds(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalSetWithdrawalCycleSeconds) (event.Subscription, error) {

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "SetWithdrawalCycleSeconds")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalSetWithdrawalCycleSeconds)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "SetWithdrawalCycleSeconds", log); err != nil {
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

// ParseSetWithdrawalCycleSeconds is a log parse operation binding the contract event 0x104b9e6ae581a335baf19e149870062416733a2a2fa10ed5c2bdb82bf4267478.
//
// Solidity: event SetWithdrawalCycleSeconds(uint256 cycleSeconds)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseSetWithdrawalCycleSeconds(log types.Log) (*NetworkWithdrawalSetWithdrawalCycleSeconds, error) {
	event := new(NetworkWithdrawalSetWithdrawalCycleSeconds)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "SetWithdrawalCycleSeconds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalUnstakeIterator is returned from FilterUnstake and is used to iterate over the raw logs and unpacked data for Unstake events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalUnstakeIterator struct {
	Event *NetworkWithdrawalUnstake // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalUnstakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalUnstake)
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
		it.Event = new(NetworkWithdrawalUnstake)
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
func (it *NetworkWithdrawalUnstakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalUnstakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalUnstake represents a Unstake event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalUnstake struct {
	From           common.Address
	LsdTokenAmount *big.Int
	EthAmount      *big.Int
	WithdrawIndex  *big.Int
	Instantly      bool
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUnstake is a free log retrieval operation binding the contract event 0xc7ccdcb2d25f572c6814e377dbb34ea4318a4b7d3cd890f5cfad699d75327c7c.
//
// Solidity: event Unstake(address indexed from, uint256 lsdTokenAmount, uint256 ethAmount, uint256 withdrawIndex, bool instantly)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterUnstake(opts *bind.FilterOpts, from []common.Address) (*NetworkWithdrawalUnstakeIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "Unstake", fromRule)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalUnstakeIterator{contract: _NetworkWithdrawal.contract, event: "Unstake", logs: logs, sub: sub}, nil
}

// WatchUnstake is a free log subscription operation binding the contract event 0xc7ccdcb2d25f572c6814e377dbb34ea4318a4b7d3cd890f5cfad699d75327c7c.
//
// Solidity: event Unstake(address indexed from, uint256 lsdTokenAmount, uint256 ethAmount, uint256 withdrawIndex, bool instantly)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchUnstake(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalUnstake, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "Unstake", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalUnstake)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "Unstake", log); err != nil {
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

// ParseUnstake is a log parse operation binding the contract event 0xc7ccdcb2d25f572c6814e377dbb34ea4318a4b7d3cd890f5cfad699d75327c7c.
//
// Solidity: event Unstake(address indexed from, uint256 lsdTokenAmount, uint256 ethAmount, uint256 withdrawIndex, bool instantly)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseUnstake(log types.Log) (*NetworkWithdrawalUnstake, error) {
	event := new(NetworkWithdrawalUnstake)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "Unstake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalUpgradedIterator struct {
	Event *NetworkWithdrawalUpgraded // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalUpgraded)
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
		it.Event = new(NetworkWithdrawalUpgraded)
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
func (it *NetworkWithdrawalUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalUpgraded represents a Upgraded event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*NetworkWithdrawalUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalUpgradedIterator{contract: _NetworkWithdrawal.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalUpgraded)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseUpgraded(log types.Log) (*NetworkWithdrawalUpgraded, error) {
	event := new(NetworkWithdrawalUpgraded)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkWithdrawalWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the NetworkWithdrawal contract.
type NetworkWithdrawalWithdrawIterator struct {
	Event *NetworkWithdrawalWithdraw // Event containing the contract specifics and raw log

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
func (it *NetworkWithdrawalWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkWithdrawalWithdraw)
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
		it.Event = new(NetworkWithdrawalWithdraw)
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
func (it *NetworkWithdrawalWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkWithdrawalWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkWithdrawalWithdraw represents a Withdraw event raised by the NetworkWithdrawal contract.
type NetworkWithdrawalWithdraw struct {
	From              common.Address
	WithdrawIndexList []*big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x67e9df8b3c7743c9f1b625ba4f2b4e601206dbd46ed5c33c85a1242e4d23a2d1.
//
// Solidity: event Withdraw(address indexed from, uint256[] withdrawIndexList)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) FilterWithdraw(opts *bind.FilterOpts, from []common.Address) (*NetworkWithdrawalWithdrawIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.FilterLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return &NetworkWithdrawalWithdrawIterator{contract: _NetworkWithdrawal.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x67e9df8b3c7743c9f1b625ba4f2b4e601206dbd46ed5c33c85a1242e4d23a2d1.
//
// Solidity: event Withdraw(address indexed from, uint256[] withdrawIndexList)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *NetworkWithdrawalWithdraw, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkWithdrawal.contract.WatchLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkWithdrawalWithdraw)
				if err := _NetworkWithdrawal.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x67e9df8b3c7743c9f1b625ba4f2b4e601206dbd46ed5c33c85a1242e4d23a2d1.
//
// Solidity: event Withdraw(address indexed from, uint256[] withdrawIndexList)
func (_NetworkWithdrawal *NetworkWithdrawalFilterer) ParseWithdraw(log types.Log) (*NetworkWithdrawalWithdraw, error) {
	event := new(NetworkWithdrawalWithdraw)
	if err := _NetworkWithdrawal.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
