package service

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
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
	"github.com/stafiprotocol/eth-lsd-relay/bindings/FeePool"
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
	"github.com/web3-storage/go-w3s-client"
)

type Service struct {
	stop                chan struct{}
	eth1Endpoint        string
	eth2Endpoint        string
	nodeRewardsFilePath string
	keyPair             *secp256k1.Keypair
	gasLimit            *big.Int
	maxGasPrice         *big.Int

	submitBalancesDuEpochs        uint64
	distributeWithdrawalsDuEpochs uint64
	distributePriorityFeeDuEpochs uint64
	merkleRootDuEpochs            uint64

	batchRequestBlocksNumber uint64

	// --- need init on start
	dev bool

	connection          *connection.Connection
	eth1Client          *ethclient.Client
	web3Client          w3s.Client
	eth2Config          beacon.Eth2Config
	withdrawCredentials []byte
	domain              []byte // for eth2 sigs

	lsdNetworkFactoryAddress common.Address
	lsdTokenAddress          common.Address
	feePoolAddress           common.Address

	lsdNetworkFactoryContract *lsd_network_factory.LsdNetworkFactory
	nodeDepositContract       *node_deposit.NodeDeposit
	networkWithdrawContract   *network_withdraw.NetworkWithdraw
	govDepositContract        *deposit_contract.DepositContract
	networkProposalContract   *network_proposal.NetworkProposal
	networkBalancesContract   *network_balances.NetworkBalances
	lsdTokenContract          *erc20.Erc20
	userDepositContract       *user_deposit.UserDeposit
	feePoolContract           *fee_pool.FeePool

	nodeCommissionRate     decimal.Decimal
	platfromCommissionRate decimal.Decimal

	quenedVoteHandlers []Handler
	quenedSyncHandlers []Handler

	latestSlotOfSyncBlock        uint64
	latestBlockOfSyncBlock       uint64
	latestBlockOfSyncEvents      uint64
	latestBlockOfUpdateValidator uint64
	latestEpochOfUpdateValidator uint64
	networkCreateBlock           uint64

	cycleSeconds                      uint64
	latestDistributeWithdrawalsHeight uint64
	latestDistributePriorityFeeHeight uint64
	latestMerkleRootEpoch             uint64

	govDeposits map[string][][]byte // pubkey(hex.encodeToString) -> withdrawalCredentials

	validators             map[string]*Validator // pubkey(hex.encodeToString) -> validator
	validatorsByIndex      map[uint64]*Validator // validator index -> validator
	validatorsByIndexMutex sync.RWMutex

	nodes map[common.Address]*Node // nodeAddress -> node

	stakerWithdrawals map[uint64]*StakerWithdrawal // withraw index => stakerWithdrawal

	cachedBeaconBlock      map[uint64]*CachedBeaconBlock // executionBlockNumber => beaconblock
	cachedBeaconBlockMutex sync.RWMutex

	exitElections map[uint64]*ExitElection // cycle -> exitElection
}

type Node struct {
	NodeAddress common.Address
	NodeType    uint8 // 1 light node 2 trust node
}
type Validator struct {
	Pubkey []byte

	NodeAddress           common.Address
	DepositSignature      []byte
	NodeDepositAmountDeci decimal.Decimal // decimals 18
	NodeDepositAmount     uint64          //decimals 9
	DepositBlock          uint64
	ActiveEpoch           uint64
	EligibleEpoch         uint64
	ExitEpoch             uint64
	WithdrawableEpoch     uint64
	NodeType              uint8  // 1 light node 2 trust node
	ValidatorIndex        uint64 // Notice!!!!!!: validator index is zero before status waiting

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

type ExitElection struct {
	WithdrawCycle      uint64
	ValidatorIndexList []uint64
}

type CachedBeaconBlock struct {
	ProposerIndex uint64
	Withdrawals   []*CachedWithdrawal

	// execute layer
	FeeRecipient common.Address
	Transactions []*CachedTransaction
	PriorityFee  *big.Int // may be nil if not pool validator
}

type CachedTransaction struct {
	// big endian
	Recipient []byte
	// big endian
	Amount []byte
}

type CachedWithdrawal struct {
	ValidatorIndex uint64
	Amount         uint64
}

type Handler struct {
	method func() error
	name   string
}

func NewService(cfg *config.Config, keyPair *secp256k1.Keypair) (*Service, error) {
	if !common.IsHexAddress(cfg.Contracts.LsdTokenAddress) {
		return nil, fmt.Errorf("LsdTokenAddress contract address fmt err")
	}
	if !common.IsHexAddress(cfg.Contracts.LsdFactoryAddress) {
		return nil, fmt.Errorf("LsdFactoryAddress contract address fmt err")
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

	err = os.MkdirAll(cfg.LogFilePath, 0700)
	if err != nil {
		return nil, fmt.Errorf("LogFilePath %w", err)
	}

	if isDir, err := utils.IsDir(cfg.LogFilePath); err != nil {
		return nil, err
	} else if !isDir {
		return nil, fmt.Errorf("logFilePath %s is not dir", cfg.LogFilePath)
	}
	if len(cfg.Web3StorageApiToken) == 0 {
		return nil, fmt.Errorf("web3StorageApiToken empty")
	}
	w3sClient, err := w3s.NewClient(w3s.WithToken(cfg.Web3StorageApiToken))
	if err != nil {
		return nil, fmt.Errorf("error creating new Web3.Storage client: %w", err)
	}
	if cfg.BatchRequestBlocksNumber == 0 {
		return nil, fmt.Errorf("BatchRequestBlocksNumber is zero")
	}

	s := &Service{
		stop:                     make(chan struct{}),
		eth1Endpoint:             cfg.Eth1Endpoint,
		eth2Endpoint:             cfg.Eth2Endpoint,
		nodeRewardsFilePath:      cfg.LogFilePath,
		eth1Client:               eth1client,
		web3Client:               w3sClient,
		lsdTokenAddress:          common.HexToAddress(cfg.Contracts.LsdTokenAddress),
		lsdNetworkFactoryAddress: common.HexToAddress(cfg.Contracts.LsdFactoryAddress),
		keyPair:                  keyPair,
		gasLimit:                 gasLimitDeci.BigInt(),
		maxGasPrice:              maxGasPriceDeci.BigInt(),
		batchRequestBlocksNumber: cfg.BatchRequestBlocksNumber,

		govDeposits:       make(map[string][][]byte),
		validators:        make(map[string]*Validator),
		validatorsByIndex: make(map[uint64]*Validator),
		nodes:             make(map[common.Address]*Node),
		stakerWithdrawals: make(map[uint64]*StakerWithdrawal),
		cachedBeaconBlock: make(map[uint64]*CachedBeaconBlock),
		exitElections:     make(map[uint64]*ExitElection),
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

	logrus.Info("init contracts...")
	err = s.initContract()
	if err != nil {
		return err
	}
	logrus.Debug("init contracts end")

	credentials, err := s.nodeDepositContract.WithdrawCredentials(nil)
	if err != nil {
		return err
	}
	s.withdrawCredentials = credentials

	// get updateBalances epochs
	updateBalancesEpochs, err := s.networkBalancesContract.UpdateBalancesEpochs(nil)
	if err != nil {
		return err
	}
	if updateBalancesEpochs.Uint64() == 0 {
		return fmt.Errorf("updateBalancesEpochs is zero")
	}

	s.submitBalancesDuEpochs = updateBalancesEpochs.Uint64()
	s.distributeWithdrawalsDuEpochs = updateBalancesEpochs.Uint64()
	s.distributePriorityFeeDuEpochs = updateBalancesEpochs.Uint64()
	s.merkleRootDuEpochs = updateBalancesEpochs.Uint64()

	// init commission
	nodeCommissionRate, err := s.networkWithdrawContract.NodeCommissionRate(nil)
	if err != nil {
		return err
	}
	s.nodeCommissionRate = decimal.NewFromBigInt(nodeCommissionRate, 0).Div(decimal.NewFromInt(1e18))
	platformCommissionRate, err := s.networkWithdrawContract.PlatformCommissionRate(nil)
	if err != nil {
		return err
	}
	s.platfromCommissionRate = decimal.NewFromBigInt(platformCommissionRate, 0).Div(decimal.NewFromInt(1e18))

	logrus.Infof("nodeCommission rate: %s, platformCommission rate: %s", s.nodeCommissionRate.String(), s.platfromCommissionRate.String())

	// init cycle seconds
	cycleSeconds, err := s.networkWithdrawContract.WithdrawCycleSeconds(nil)
	if err != nil {
		return err
	}
	s.cycleSeconds = cycleSeconds.Uint64()

	// init latest block and slot number
	s.latestBlockOfUpdateValidator = s.networkCreateBlock
	s.latestBlockOfSyncEvents = s.networkCreateBlock

	s.latestBlockOfSyncBlock = math.MaxUint64
	checkAndUpdateLatestBlockOfSyncBlock := func(block uint64) {
		logrus.Debugf("checkAndUpdateLatestBlockOfSyncBlock block: %d", block)
		if block < s.latestBlockOfSyncBlock {
			s.latestBlockOfSyncBlock = block
		}
	}
	latestDistributeWithdrawalsHeight, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(nil)
	if err != nil {
		return fmt.Errorf("LatestDistributeWithdrawalsHeight %w", err)
	}
	checkAndUpdateLatestBlockOfSyncBlock(latestDistributeWithdrawalsHeight.Uint64())

	latestDistributePriorityFeeHeight, err := s.networkWithdrawContract.LatestDistributePriorityFeeHeight(nil)
	if err != nil {
		return err
	}
	checkAndUpdateLatestBlockOfSyncBlock(latestDistributePriorityFeeHeight.Uint64())

	merkleRootEpoch, err := s.networkWithdrawContract.LatestMerkleRootEpoch(nil)
	if err != nil {
		return err
	}
	if merkleRootEpoch.Uint64() > 0 {
		epochBlockNumber, err := s.getEpochStartBlocknumberWithCheck(merkleRootEpoch.Uint64())
		if err != nil {
			return err
		}
		checkAndUpdateLatestBlockOfSyncBlock(epochBlockNumber)
	} else {
		checkAndUpdateLatestBlockOfSyncBlock(0)
	}
	if s.latestBlockOfSyncBlock < s.networkCreateBlock {
		s.latestBlockOfSyncBlock = s.networkCreateBlock
	}

	block, err := s.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(s.latestBlockOfSyncBlock)))
	if err != nil {
		return err
	}
	s.latestSlotOfSyncBlock = utils.SlotAtTimestamp(s.eth2Config, block.Time())
	logrus.Debugf("latestSlotOfSyncBlock: %d, latestBlockOfSyncBlock: %d", s.latestSlotOfSyncBlock, s.latestBlockOfSyncBlock)

	// start services
	logrus.Info("start services...")
	s.appendSyncHandlers(s.syncBlocks, s.pruneBlocks)

	s.appendVoteHandlers(
		s.updateValidatorsFromNetwork, s.updateValidatorsFromBeacon, s.syncEvents,
		s.voteWithdrawCredentials,
		s.submitBalances, s.distributeWithdrawals, s.distributePriorityFee,
		s.setMerkleRoot, s.notifyValidatorExit)

	utils.SafeGo(s.syncService)
	utils.SafeGo(s.voteService)

	return nil
}

func (s *Service) Stop() {
	close(s.stop)
}

func (s *Service) initContract() error {
	var err error

	s.lsdNetworkFactoryContract, err = lsd_network_factory.NewLsdNetworkFactory(s.lsdNetworkFactoryAddress, s.eth1Client)
	if err != nil {
		return err
	}

	networkContracts, err := s.lsdNetworkFactoryContract.NetworkContractsOfLsdToken(nil, s.lsdTokenAddress)
	if err != nil {
		return err
	}
	logrus.Infof("networkContracts: %+v", networkContracts)

	ethDepositAddress, err := s.lsdNetworkFactoryContract.EthDepositAddress(nil)
	if err != nil {
		return err
	}

	s.govDepositContract, err = deposit_contract.NewDepositContract(ethDepositAddress, s.eth1Client)
	if err != nil {
		return err
	}

	s.nodeDepositContract, err = node_deposit.NewNodeDeposit(networkContracts.NodeDeposit, s.eth1Client)
	if err != nil {
		return err
	}

	s.userDepositContract, err = user_deposit.NewUserDeposit(networkContracts.UserDeposit, s.eth1Client)
	if err != nil {
		return err
	}
	s.networkWithdrawContract, err = network_withdraw.NewNetworkWithdraw(networkContracts.NetworkWithdraw, s.eth1Client)
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

	s.feePoolContract, err = fee_pool.NewFeePool(networkContracts.FeePool, s.eth1Client)
	if err != nil {
		return err
	}

	s.networkCreateBlock = networkContracts.Block.Uint64()

	s.feePoolAddress = networkContracts.FeePool

	return nil
}

func (s *Service) voteService() {
	logrus.Info("start vote service")
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

			for _, handler := range s.quenedVoteHandlers {
				funcName := handler.name
				logrus.Debugf("handler %s start...", funcName)

				err := handler.method()
				if err != nil {
					logrus.Warnf("handler %s failed: %s, will retry.", funcName, err)
					time.Sleep(utils.RetryInterval * 4)
					retry++
					continue Out
				}
				logrus.Debugf("handler %s end", funcName)
			}

			retry = 0
		}

		time.Sleep(12 * time.Second)
	}
}

func (s *Service) syncService() {
	logrus.Info("start sync service")
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

			for _, handler := range s.quenedSyncHandlers {
				funcName := handler.name
				logrus.Debugf("handler %s start...", funcName)

				err := handler.method()
				if err != nil {
					logrus.Warnf("handler %s failed: %s, will retry.", funcName, err)
					time.Sleep(utils.RetryInterval * 4)
					retry++
					continue Out
				}
				logrus.Debugf("handler %s end", funcName)
			}

			retry = 0
		}

		time.Sleep(12 * time.Second)
	}
}

func (s *Service) appendVoteHandlers(handlers ...func() error) {
	for _, handler := range handlers {

		funcNameRaw := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

		splits := strings.Split(funcNameRaw, "/")
		funcName := splits[len(splits)-1]

		s.quenedVoteHandlers = append(s.quenedVoteHandlers, Handler{
			method: handler,
			name:   funcName,
		})
	}
}

func (s *Service) appendSyncHandlers(handlers ...func() error) {
	for _, handler := range handlers {

		funcNameRaw := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

		splits := strings.Split(funcNameRaw, "/")
		funcName := splits[len(splits)-1]

		s.quenedSyncHandlers = append(s.quenedSyncHandlers, Handler{
			method: handler,
			name:   funcName,
		})
	}
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

func (s *Service) getBeaconBlock(eth1BlcokNumber uint64) (*CachedBeaconBlock, error) {
	s.cachedBeaconBlockMutex.RLock()
	defer s.cachedBeaconBlockMutex.RUnlock()

	block, exist := s.cachedBeaconBlock[eth1BlcokNumber]
	if !exist {
		return nil, fmt.Errorf("block %d not cached", eth1BlcokNumber)
	}
	return block, nil

}

func (s *Service) getValidatorByIndex(valIndex uint64) (*Validator, bool) {
	s.validatorsByIndexMutex.RLock()
	defer s.validatorsByIndexMutex.RUnlock()

	v, exist := s.validatorsByIndex[valIndex]
	return v, exist
}

func (s *Service) exitButNotDistributedValidatorList(epoch uint64) []*Validator {
	vals := make([]*Validator, 0)
	for _, v := range s.validators {
		if v.ExitEpoch != 0 && v.ExitEpoch <= epoch && v.WithdrawableEpoch > epoch {
			vals = append(vals, v)
		}
	}
	return vals
}

func (s *Service) notExitElectionListBefore(willDealCycle uint64) []*ExitElection {
	els := make([]*ExitElection, 0)
	for cycle, e := range s.exitElections {
		if cycle >= willDealCycle {
			continue
		}
		for _, valIndex := range e.ValidatorIndexList {
			val, exist := s.getValidatorByIndex(valIndex)
			if exist && val.ExitEpoch == 0 {
				els = append(els, e)
				break
			}
		}
	}

	sort.SliceStable(els, func(i, j int) bool {
		return els[i].WithdrawCycle < els[j].WithdrawCycle
	})

	return els
}
