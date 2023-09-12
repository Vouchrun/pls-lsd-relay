package service

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/DepositContract"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/Erc20"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/LsdNetworkFactory"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/NetworkBalances"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/NetworkProposal"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/NetworkWithdraw"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/NodeDeposit"
	"github.com/stafiprotocol/eth-lsd-relay/bindings/UserDeposit"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

var (
	lsdNetworkFactoryAddressMainnet = common.HexToAddress("")
	lsdNetworkFactoryAddressTestnet = common.HexToAddress("")

	govDepositContractAddressMainnet = common.HexToAddress("")
	govDepositContractAddressTestnet = common.HexToAddress("")
)

type Service struct {
	stop                          chan struct{}
	eth1Endpoint                  string
	eth2Endpoint                  string
	keyPair                       *secp256k1.Keypair
	gasLimit                      *big.Int
	maxGasPrice                   *big.Int
	submitBalancesDuEpochs        uint64
	distributeWithdrawalsDuEpochs uint64
	distributePriorityFeeDuEpochs uint64

	// --- need init on start
	dev             bool
	eth1StartHeight uint64

	connection          *connection.Connection
	eth1Client          *ethclient.Client
	eth2Config          beacon.Eth2Config
	withdrawCredentials []byte
	domain              []byte // for eth2 sigs

	govDepositContractAddress common.Address
	lsdNetworkFactoryAddress  common.Address
	lsdTokenAddress           common.Address

	lsdNetworkFactoryContract *lsd_network_factory.LsdNetworkFactory
	nodeDepositContract       *node_deposit.NodeDeposit
	networkWithdrawContract   *network_withdraw.NetworkWithdraw
	govDepositContract        *deposit_contract.DepositContract
	networkProposalContract   *network_proposal.NetworkProposal
	networkBalancesContract   *network_balances.NetworkBalances
	lsdTokenContract          *erc20.Erc20
	userDeposit               *user_deposit.UserDeposit

	quenedHandlers []Handler

	dealedEth1Block    uint64
	networkCreateBlock uint64

	govDeposits map[string][][]byte // pubkey -> withdrawalCredentials

	validators map[string]*Validator // pubkey(hex.encodeToString) -> validator

	nodes map[common.Address]*Node // nodeAddress -> node

	stakerWithdrawals map[uint64]*StakerWithdrawal // withraw index => stakerWithdrawal
}

type Node struct {
	NodeAddress common.Address
	NodeType    uint8 // 1 light node 2 trust node
}
type Validator struct {
	Pubkey []byte

	NodeAddress       common.Address
	DepositSignature  []byte
	NodeDepositAmount decimal.Decimal
	DepositBlock      uint64
	ActiveEpoch       uint64
	EligibleEpoch     uint64
	ExitEpoch         uint64
	WithdrawableEpoch uint64
	NodeType          uint8  // 1 light node 2 trust node
	ValidatorIndex    uint64 // Notice!!!!!!: validator index is zero before status waiting

	Balance          uint64 // realtime balance
	EffectiveBalance uint64 // realtime effectiveBalance
	Status           uint8  // status details defined in pkg/utils/eth2.go
}

type StakerWithdrawal struct {
	WithdrawIndex uint64

	Address            string //hex with 0x prefix
	EthAmount          string
	BlockNumber        uint64
	ClaimedBlockNumber uint64
	Timestamp          uint64 // unstake tx timestamp
}

type Handler struct {
	method func() error
	name   string
}

func NewService(cfg *config.Config, keyPair *secp256k1.Keypair) (*Service, error) {
	if !common.IsHexAddress(cfg.Contracts.LsdTokenAddress) {
		return nil, fmt.Errorf("SsvTokenAddress contract address fmt err")
	}

	gasLimitDeci, err := decimal.NewFromString(cfg.GasLimit)
	if err != nil {
		return nil, err
	}

	if gasLimitDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("gas limit is zero")
	}
	maxGasPriceDeci, err := decimal.NewFromString(cfg.MaxGasPrice)
	if err != nil {
		return nil, err
	}
	if maxGasPriceDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("max gas price is zero")
	}

	eth1client, err := ethclient.Dial(cfg.Eth1Endpoint)
	if err != nil {
		return nil, err
	}
	s := &Service{
		stop:                          make(chan struct{}),
		eth1Endpoint:                  cfg.Eth1Endpoint,
		eth2Endpoint:                  cfg.Eth2Endpoint,
		eth1Client:                    eth1client,
		lsdTokenAddress:               common.HexToAddress(cfg.Contracts.LsdTokenAddress),
		keyPair:                       keyPair,
		gasLimit:                      gasLimitDeci.BigInt(),
		maxGasPrice:                   maxGasPriceDeci.BigInt(),
		submitBalancesDuEpochs:        225, // 1 day
		distributeWithdrawalsDuEpochs: 225,
		distributePriorityFeeDuEpochs: 225,
	}

	return s, nil
}

func (s *Service) Start() error {
	var err error
	s.connection, err = connection.NewConnection(s.eth1Endpoint, s.eth2Endpoint, s.keyPair,
		s.gasLimit, s.maxGasPrice)
	if err != nil {
		return err
	}

	chainId, err := s.eth1Client.ChainID(context.Background())
	if err != nil {
		return err
	}

	s.eth2Config, err = s.connection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}

	switch chainId.Uint64() {
	case 1: //mainnet
		s.dev = false
		if !bytes.Equal(s.eth2Config.GenesisForkVersion, params.MainnetConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
		s.lsdNetworkFactoryAddress = lsdNetworkFactoryAddressMainnet
		s.govDepositContractAddress = govDepositContractAddressMainnet

		domain, err := signing.ComputeDomain(
			params.MainnetConfig().DomainDeposit,
			params.MainnetConfig().GenesisForkVersion,
			params.MainnetConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		s.domain = domain

	case 11155111: // sepolia
		s.dev = true
		if !bytes.Equal(s.eth2Config.GenesisForkVersion, params.SepoliaConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
		s.lsdNetworkFactoryAddress = lsdNetworkFactoryAddressTestnet
		s.govDepositContractAddress = govDepositContractAddressTestnet

		domain, err := signing.ComputeDomain(
			params.SepoliaConfig().DomainDeposit,
			params.SepoliaConfig().GenesisForkVersion,
			params.SepoliaConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		s.domain = domain
	case 5: // goerli
		s.dev = true
		if !bytes.Equal(s.eth2Config.GenesisForkVersion, params.PraterConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
		s.lsdNetworkFactoryAddress = lsdNetworkFactoryAddressTestnet
		s.govDepositContractAddress = govDepositContractAddressTestnet

		domain, err := signing.ComputeDomain(
			params.PraterConfig().DomainDeposit,
			params.PraterConfig().GenesisForkVersion,
			params.PraterConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		s.domain = domain
	default:
		return fmt.Errorf("unsupport chainId: %d", chainId.Int64())
	}
	if err != nil {
		return err
	}

	// init dealed eth1 block
	latestBlockNumber, err := s.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}
	if latestBlockNumber > depositEventPreBlocks {
		s.dealedEth1Block = latestBlockNumber - depositEventPreBlocks
	}

	logrus.Info("init contracts...")
	err = s.initContract()
	if err != nil {
		return err
	}

	logrus.Info("start services...")
	s.appendHandlers(s.syncDepositInfo, s.updateValidatorsFromNetwork, s.updateValidatorsFromBeacon,
		s.voteWithdrawCredentials, s.submitBalances, s.distributeWithdrawals, s.distributePriorityFee)

	utils.SafeGo(s.voteService)

	return nil
}

func (s *Service) Stop() {
	close(s.stop)
}

func (s *Service) initContract() error {
	var err error
	s.govDepositContract, err = deposit_contract.NewDepositContract(s.govDepositContractAddress, s.eth1Client)
	if err != nil {
		return err
	}

	s.lsdNetworkFactoryContract, err = lsd_network_factory.NewLsdNetworkFactory(s.lsdNetworkFactoryAddress, s.eth1Client)
	if err != nil {
		return err
	}

	networkContracts, err := s.lsdNetworkFactoryContract.NetworkContractsOf(nil, s.lsdTokenAddress)
	if err != nil {
		return err
	}

	s.eth1StartHeight = networkContracts.Block.Uint64()

	s.nodeDepositContract, err = node_deposit.NewNodeDeposit(networkContracts.NodeDeposit, s.eth1Client)
	if err != nil {
		return err
	}
	s.networkWithdrawContract, err = network_withdraw.NewNetworkWithdraw(networkContracts.NodeDeposit, s.eth1Client)
	if err != nil {
		return err
	}

	s.networkProposalContract, err = network_proposal.NewNetworkProposal(networkContracts.NetworkProposal, s.eth1Client)
	if err != nil {
		return err
	}

	s.networkBalancesContract, err = network_balances.NewNetworkBalances(networkContracts.NetworkBalances, s.eth1Client)
	if err != nil {
		return err
	}
	s.lsdTokenContract, err = erc20.NewErc20(networkContracts.LsdToken, s.eth1Client)
	if err != nil {
		return err
	}

	s.networkCreateBlock = networkContracts.Block.Uint64()

	return nil
}

func (s *Service) voteService() {
	logrus.Info("start ssv service")
	retry := 0

Out:
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return
		}

		select {
		case <-s.stop:
			logrus.Info("task has stopped")
			return
		default:

			for _, handler := range s.quenedHandlers {
				funcName := handler.name
				logrus.Debugf("handler %s start.........", funcName)

				err := handler.method()
				if err != nil {
					logrus.Warnf("handler %s failed: %s, will retry.", funcName, err)
					time.Sleep(utils.RetryInterval * 4)
					retry++
					continue Out
				}
				logrus.Debugf("handler %s end.........", funcName)
			}

			retry = 0
		}

		time.Sleep(48 * time.Second) // 48 blocks
	}
}

func (s *Service) appendHandlers(handlers ...func() error) {
	for _, handler := range handlers {

		funcNameRaw := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

		splits := strings.Split(funcNameRaw, "/")
		funcName := splits[len(splits)-1]

		s.quenedHandlers = append(s.quenedHandlers, Handler{
			method: handler,
			name:   funcName,
		})
	}
}

func (s *Service) waitTxOk(txHash common.Hash) error {
	_, err := utils.WaitTxOkCommon(s.eth1Client, txHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetValidatorDepositedListBefore(block uint64) []*Validator {
	selectedValidator := make([]*Validator, 0)
	for _, v := range s.validators {
		if v.DepositBlock <= block {
			selectedValidator = append(selectedValidator, v)
		}
	}
	return selectedValidator
}
