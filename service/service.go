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

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	deposit_contract "github.com/stafiprotocol/eth-lsd-relay/bindings/DepositContract"
	erc20 "github.com/stafiprotocol/eth-lsd-relay/bindings/Erc20"
	fee_pool "github.com/stafiprotocol/eth-lsd-relay/bindings/FeePool"
	lsd_network_factory "github.com/stafiprotocol/eth-lsd-relay/bindings/LsdNetworkFactory"
	network_balances "github.com/stafiprotocol/eth-lsd-relay/bindings/NetworkBalances"
	network_proposal "github.com/stafiprotocol/eth-lsd-relay/bindings/NetworkProposal"
	network_withdraw "github.com/stafiprotocol/eth-lsd-relay/bindings/NetworkWithdraw"
	node_deposit "github.com/stafiprotocol/eth-lsd-relay/bindings/NodeDeposit"
	user_deposit "github.com/stafiprotocol/eth-lsd-relay/bindings/UserDeposit"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/local_store"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
	"github.com/web3-storage/go-w3s-client"
)

type Service struct {
	stop                chan struct{}
	startServiceOnce    sync.Once
	eth1Endpoint        string
	eth2Endpoint        string
	nodeRewardsFilePath string

	submitBalancesDuEpochs        uint64
	distributeWithdrawalsDuEpochs uint64
	distributePriorityFeeDuEpochs uint64
	merkleRootDuEpochs            uint64

	batchRequestBlocksNumber uint64

	connection          *connection.CachedConnection
	web3Client          w3s.Client
	eth2Config          beacon.Eth2Config
	withdrawCredentials []byte
	domain              []byte // for eth2 sigs

	lsdNetworkFactoryAddress common.Address
	lsdTokenAddress          common.Address
	feePoolAddress           common.Address
	networkWithdrawAddress   common.Address
	networkBalancesAddress   common.Address
	nodeDepositAddress       common.Address

	networkWithdrdawAbi abi.ABI
	networkBalancesAbi  abi.ABI
	nodeDepositAbi      abi.ABI

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

	latestSlotOfSyncBlock   uint64
	latestBlockOfSyncBlock  uint64
	waitFirstNodeStakeEvent bool
	localSyncedBlockHeight  uint64
	localStore              *local_store.LocalStore

	latestBlockOfSyncEvents      uint64
	latestBlockOfUpdateValidator uint64
	latestEpochOfUpdateValidator uint64
	startAtBlock                 uint64

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
	NodeAddress  common.Address
	NodeType     uint8 // 1 light node 2 trust node
	PubkeyNumber uint64
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

	Address            common.Address
	EthAmount          decimal.Decimal
	BlockNumber        uint64
	ClaimedBlockNumber uint64
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

func NewService(
	cfg *config.Config,
	conn *connection.CachedConnection,
	localStore *local_store.LocalStore,
) (*Service, error) {
	if !common.IsHexAddress(cfg.Contracts.LsdTokenAddress) {
		return nil, fmt.Errorf("LsdTokenAddress contract address fmt err")
	}
	if !common.IsHexAddress(cfg.Contracts.LsdFactoryAddress) {
		return nil, fmt.Errorf("LsdFactoryAddress contract address fmt err")
	}

	err := os.MkdirAll(cfg.LogFilePath, 0700)
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

	info, err := localStore.Read(cfg.Contracts.LsdTokenAddress)
	if err != nil {
		return nil, err
	}
	var localSyncedBlockHeight uint64 = 0
	if info != nil {
		localSyncedBlockHeight = info.SyncedHeight
	}

	s := &Service{
		stop:                     make(chan struct{}),
		connection:               conn,
		eth1Endpoint:             cfg.Eth1Endpoint,
		eth2Endpoint:             cfg.Eth2Endpoint,
		nodeRewardsFilePath:      cfg.LogFilePath,
		web3Client:               w3sClient,
		lsdTokenAddress:          common.HexToAddress(cfg.Contracts.LsdTokenAddress),
		lsdNetworkFactoryAddress: common.HexToAddress(cfg.Contracts.LsdFactoryAddress),
		batchRequestBlocksNumber: cfg.BatchRequestBlocksNumber,
		localSyncedBlockHeight:   localSyncedBlockHeight,
		localStore:               localStore,

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
	chainId, err := s.connection.ChainID()
	if err != nil {
		return err
	}

	s.eth2Config, err = s.connection.Eth2Config()
	if err != nil {
		return err
	}

	switch chainId.Uint64() {
	case 1: //mainnet
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
	case 17000: // holesky
		chainCfg := params.HoleskyConfig()
		if !bytes.Equal(s.eth2Config.GenesisForkVersion, chainCfg.GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
		domain, err := signing.ComputeDomain(
			chainCfg.DomainDeposit,
			chainCfg.GenesisForkVersion,
			chainCfg.ZeroHash[:],
		)
		if err != nil {
			return err
		}
		s.domain = domain
	case 5: // goerli
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
	s.latestBlockOfUpdateValidator = s.startAtBlock
	s.latestBlockOfSyncEvents = s.startAtBlock
	if err = s.initLatestBlockOfSyncBlock(); err != nil {
		return err
	}

	block, err := s.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(s.latestBlockOfSyncBlock)))
	if err != nil {
		return err
	}
	s.latestSlotOfSyncBlock = utils.SlotAtTimestamp(s.eth2Config, block.Time())
	logrus.Debugf("latestSlotOfSyncBlock: %d, latestBlockOfSyncBlock: %d", s.latestSlotOfSyncBlock, s.latestBlockOfSyncBlock)

	// init abi
	s.networkWithdrdawAbi, err = abi.JSON(strings.NewReader(network_withdraw.NetworkWithdrawABI))
	if err != nil {
		return err
	}
	s.networkBalancesAbi, err = abi.JSON(strings.NewReader(network_balances.NetworkBalancesABI))
	if err != nil {
		return err
	}
	s.nodeDepositAbi, err = abi.JSON(strings.NewReader(node_deposit.NodeDepositABI))
	if err != nil {
		return err
	}

	// start services
	logrus.Info("start services...")
	if s.waitFirstNodeStakeEvent {
		s.startSeekFirstNodeStakeEvent()
	} else {
		s.startHandlers()
	}

	return nil
}

func (s *Service) startSeekFirstNodeStakeEvent() {
	log := logrus.WithFields(logrus.Fields{
		"lsdToken": s.lsdTokenAddress.Hex(),
		"service":  "seekingFirstNodeStake",
	})
	utils.SafeGo(func() {
		log.Info("start service")
		defer log.Info("service stopped")
		for {
			select {
			case <-s.stop:
				return
			default:
				found, err := s.seekFirstNodeStakeEvent()
				if found {
					// first node stake event has been found
					log.Info("found first node stake event")
					return
				}
				if err != nil {
					log.WithFields(logrus.Fields{
						"err": err,
					}).Warn("seek first node stake event error")
				}
				utils.Sleep(s.stop, time.Minute*30)
			}
		}
	})
}
func (s *Service) seekFirstNodeStakeEvent() (bool, error) {
	latestBlock, err := s.connection.Eth1LatestBlock()
	if err != nil {
		return false, err
	}
	start := s.latestBlockOfSyncBlock + 1
	end := latestBlock

	for i := start; i <= end; i += fetchEventBlockLimit {
		subStart := i
		subEnd := i + fetchEventBlockLimit - 1
		if end < i+fetchEventBlockLimit {
			subEnd = end
		}

		iter, err := retry.DoWithData(func() (*node_deposit.NodeDepositDepositedIterator, error) {
			return s.nodeDepositContract.FilterDeposited(&bind.FilterOpts{
				Start:   subStart,
				End:     &subEnd,
				Context: context.Background(),
			})
		}, retry.Delay(time.Second*2), retry.Attempts(5))
		if err != nil {
			return false, err
		}

		hasEvent := iter.Next()
		iter.Close()
		if hasEvent {
			// found the first node stake event
			s.waitFirstNodeStakeEvent = false
			s.startAtBlock = iter.Event.Raw.BlockNumber - 1
			s.latestBlockOfSyncBlock = s.startAtBlock
			s.latestBlockOfUpdateValidator = s.latestBlockOfSyncBlock
			s.latestBlockOfSyncEvents = s.latestBlockOfSyncBlock

			block, err := s.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(s.latestBlockOfSyncBlock)))
			if err != nil {
				return false, err
			}
			s.latestSlotOfSyncBlock = utils.SlotAtTimestamp(s.eth2Config, block.Time())

			s.startHandlers()
			return true, nil
		}
	}

	if err = s.localStore.Update(local_store.Info{
		SyncedHeight: end,
		Address:      s.lsdTokenAddress.Hex(),
	}); err != nil {
		return false, err
	}
	s.latestBlockOfSyncBlock = end
	return false, nil
}

func (s *Service) startHandlers() {
	s.startServiceOnce.Do(func() {
		logrus.WithFields(logrus.Fields{
			"lsdToken":               s.lsdTokenAddress.Hex(),
			"latestBlockOfSyncBlock": s.latestBlockOfSyncBlock,
		}).Info("start voting handlers")

		s.appendSyncHandlers(s.syncBlocks, s.pruneBlocks)

		s.appendVoteHandlers(
			s.updateValidatorsFromNetwork, s.updateValidatorsFromBeacon, s.syncEvents,
			s.voteWithdrawCredentials,
			s.submitBalances, s.distributeWithdrawals, s.distributePriorityFee,
			s.setMerkleRoot, s.notifyValidatorExit)

		utils.SafeGo(s.syncService)
		utils.SafeGo(s.voteService)
	})
}

func (s *Service) Stop() {
	close(s.stop)
}

func (s *Service) initContract() error {
	var err error

	s.lsdNetworkFactoryContract, err = lsd_network_factory.NewLsdNetworkFactory(s.lsdNetworkFactoryAddress, s.connection.Eth1Client())
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

	s.govDepositContract, err = deposit_contract.NewDepositContract(ethDepositAddress, s.connection.Eth1Client())
	if err != nil {
		return err
	}

	s.nodeDepositContract, err = node_deposit.NewNodeDeposit(networkContracts.NodeDeposit, s.connection.Eth1Client())
	if err != nil {
		return err
	}

	s.userDepositContract, err = user_deposit.NewUserDeposit(networkContracts.UserDeposit, s.connection.Eth1Client())
	if err != nil {
		return err
	}
	s.networkWithdrawContract, err = network_withdraw.NewNetworkWithdraw(networkContracts.NetworkWithdraw, s.connection.Eth1Client())
	if err != nil {
		return err
	}

	s.networkProposalContract, err = network_proposal.NewNetworkProposal(networkContracts.NetworkProposal, s.connection.Eth1Client())
	if err != nil {
		return err
	}

	s.networkBalancesContract, err = network_balances.NewNetworkBalances(networkContracts.NetworkBalances, s.connection.Eth1Client())
	if err != nil {
		return err
	}
	s.lsdTokenContract, err = erc20.NewErc20(s.lsdTokenAddress, s.connection.Eth1Client())
	if err != nil {
		return err
	}

	s.feePoolContract, err = fee_pool.NewFeePool(networkContracts.FeePool, s.connection.Eth1Client())
	if err != nil {
		return err
	}

	s.networkWithdrawAddress = networkContracts.NetworkWithdraw
	s.networkBalancesAddress = networkContracts.NetworkBalances
	s.nodeDepositAddress = networkContracts.NodeDeposit

	s.startAtBlock = networkContracts.Block.Uint64()

	s.feePoolAddress = networkContracts.FeePool

	return nil
}

func (s *Service) initLatestBlockOfSyncBlock() error {
	s.latestBlockOfSyncBlock = math.MaxUint64
	checkAndUpdateLatestBlockOfSyncBlock := func(block uint64) {
		logrus.Debugf("checkAndUpdateLatestBlockOfSyncBlock block: %d", block)
		if block < s.latestBlockOfSyncBlock {
			s.latestBlockOfSyncBlock = block
		}
	}

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

	// latest block should less than LatestDistributeWithdrawalsHeight at cycle snapshot
	_, targetTimestamp, err := s.currentCycleAndStartTimestamp()
	if err != nil {
		return fmt.Errorf("currentCycleAndStartTimestamp failed: %w", err)
	}
	targetEpoch := utils.EpochAtTimestamp(s.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := s.getEpochStartBlocknumberWithCheck(targetEpoch)
	if err != nil {
		return err
	}
	targetCall := s.connection.CallOpts(big.NewInt(int64(targetBlockNumber)))
	latestDistributeWithdrawalHeight, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(targetCall)
	if err != nil {
		return err
	}
	checkAndUpdateLatestBlockOfSyncBlock(latestDistributeWithdrawalHeight.Uint64())

	// should greater network create block
	if s.latestBlockOfSyncBlock < s.startAtBlock {
		s.latestBlockOfSyncBlock = s.startAtBlock
		s.waitFirstNodeStakeEvent = true
	}
	// should be greater than local synced block height
	if s.latestBlockOfSyncBlock < s.localSyncedBlockHeight {
		s.latestBlockOfSyncBlock = s.localSyncedBlockHeight
		s.waitFirstNodeStakeEvent = true
	}

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

func (s *Service) GetValidatorDepositedListBeforeBlock(block uint64) []*Validator {
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

func (s *Service) notExitElectionListBefore(targetEpoch, willDealCycle uint64) []*ExitElection {
	els := make([]*ExitElection, 0)
	for cycle, e := range s.exitElections {
		if cycle >= willDealCycle {
			continue
		}
		for _, valIndex := range e.ValidatorIndexList {
			val, exist := s.getValidatorByIndex(valIndex)
			if exist {

				if val.ExitEpoch == 0 {
					els = append(els, e)
					break
				} else {
					if val.ExitEpoch > targetEpoch {
						els = append(els, e)
					}
				}
			}
		}
	}

	sort.SliceStable(els, func(i, j int) bool {
		return els[i].WithdrawCycle < els[j].WithdrawCycle
	})

	return els
}
