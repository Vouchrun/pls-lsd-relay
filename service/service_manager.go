package service

import (
	"fmt"
	"math"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	xsync "github.com/puzpuzpuz/xsync/v3"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	lsd_network_factory "github.com/stafiprotocol/eth-lsd-relay/bindings/LsdNetworkFactory"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/local_store"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

type ServiceManager struct {
	stop       chan struct{}
	cfg        *config.Config
	connection *connection.CachedConnection
	srvs       *xsync.MapOf[string, *Service]
	localStore *local_store.LocalStore

	cachedBeaconBlock                  *xsync.MapOf[uint64, *CachedBeaconBlock] // beacon block id: (uint64) => beaconblock: (*CachedBeaconBlock)
	cachedBeaconBlockByExecBlockHeight *xsync.MapOf[uint64, *CachedBeaconBlock] // execution block height: (uint64) => beaconblock: (*CachedBeaconBlock)
	beaconBlockMutex                   *utils.KeyedMutex[uint64]
}

func NewServiceManager(cfg *config.Config, keyPair *secp256k1.Keypair) (*ServiceManager, error) {
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

	conn, err := connection.NewConnection(cfg.Eth1Endpoint, cfg.Eth2Endpoint, keyPair,
		gasLimitDeci.BigInt(), maxGasPriceDeci.BigInt())
	if err != nil {
		return nil, err
	}
	cachedConn, err := connection.NewCachedConnection(conn)
	if err != nil {
		return nil, err
	}
	if err = cachedConn.Start(); err != nil {
		return nil, err
	}
	localStore, err := local_store.NewLocalStore(cfg.BlockstoreFilePath)
	if err != nil {
		return nil, err
	}

	return &ServiceManager{
		stop:                               make(chan struct{}),
		cfg:                                cfg,
		connection:                         cachedConn,
		srvs:                               xsync.NewMapOf[string, *Service](),
		cachedBeaconBlock:                  xsync.NewMapOf[uint64, *CachedBeaconBlock](),
		cachedBeaconBlockByExecBlockHeight: xsync.NewMapOf[uint64, *CachedBeaconBlock](),
		beaconBlockMutex:                   &utils.KeyedMutex[uint64]{},
		localStore:                         localStore,
	}, nil
}

func (m *ServiceManager) Start() error {
	utils.SafeGoWithRestart(m.pruneCachedBeaconBlocksService)

	if !m.cfg.RunForEntrustedLsdNetwork {
		if _, err := m.newAndStartServiceFor(m.cfg.Contracts.LsdTokenAddress); err != nil {
			return err
		}
		return nil
	}

	// for entrusted lsd network
	err := retry.Do(m.syncEntrustedLsdTokens, retry.Attempts(3))
	if err != nil {
		return err
	}

	utils.SafeGo(m.startSyncService)

	return nil
}

func (m *ServiceManager) Stop() {
	close(m.stop)
	m.srvs.Range(func(key string, value *Service) bool {
		value.Stop()
		return true
	})
	m.connection.Stop()
}

func (m *ServiceManager) startSyncService() {
	logrus.Info("start listening new entrusted lsd token service")

	retry := 0

Out:
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return
		}

		select {
		case <-m.stop:
			logrus.Info("sync entrusted lsd token task has stopped")
			return
		default:

			err := m.syncEntrustedLsdTokens()
			if err != nil {
				logrus.Errorf("fail to sync entrusted token: %s", err.Error())
				time.Sleep(utils.RetryInterval * 4)
				retry++
				continue Out
			}

			retry = 0
		}

		time.Sleep(12 * time.Second)
	}
}

func (m *ServiceManager) syncEntrustedLsdTokens() error {
	lsdNetworkFactoryContract, err := lsd_network_factory.NewLsdNetworkFactory(common.HexToAddress(m.cfg.Contracts.LsdFactoryAddress), m.connection.Eth1Client())
	if err != nil {
		return err
	}
	tokens, err := lsdNetworkFactoryContract.GetEntrustedLsdTokens(nil)
	if err != nil {
		return err
	}
	tokenList := lo.Map(tokens, func(token common.Address, _ int) string { return token.String() })

	for _, token := range tokenList {
		if _, exist := m.srvs.Load(token); !exist {
			// add new entrusted lsd token
			if _, err := m.newAndStartServiceFor(token); err != nil {
				return err
			}
		}
	}

	m.srvs.Range(func(token string, srv *Service) bool {
		if !lo.Contains(tokenList, token) {
			// remove entrusted lsd token
			log := logrus.WithFields(logrus.Fields{
				"lsdToken": token,
			})
			log.Info("stopping service")
			srv.Stop()
			m.srvs.Delete(token)
			log.Warn("stopped service")
		}
		return true
	})

	return nil
}

func (m *ServiceManager) newAndStartServiceFor(lsdToken string) (*Service, error) {
	log := logrus.WithFields(logrus.Fields{
		"lsdToken": lsdToken,
	})
	log.Debug("new service instance")
	srvConfig := *m.cfg
	srvConfig.Contracts.LsdTokenAddress = lsdToken
	srv, err := NewService(&srvConfig, m, m.connection, m.localStore)
	if err != nil {
		return nil, fmt.Errorf("new service for lsd token %s err %s", lsdToken, err.Error())
	}
	if err = srv.Start(); err != nil {
		return nil, fmt.Errorf("start service for lsd token %s err %s", lsdToken, err.Error())
	}
	m.srvs.Store(lsdToken, srv)
	log.Warn("started service")
	return srv, nil
}

var notExistBeaconBlock = &CachedBeaconBlock{}

func (m *ServiceManager) CacheBeaconBlock(blockId uint64) (*CachedBeaconBlock, bool, error) {
	unlock := m.beaconBlockMutex.Lock(blockId)
	defer unlock()

	if block, ok := m.cachedBeaconBlock.Load(blockId); ok {
		if block == notExistBeaconBlock {
			return nil, false, nil
		}

		return block, true, nil
	}

	block, exist, err := m.connection.GetBeaconBlock(blockId)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		m.cachedBeaconBlock.Store(blockId, notExistBeaconBlock)
		return nil, false, nil
	}

	cachedBlock := CachedBeaconBlock{
		BeaconBlockId:        blockId,
		ExecutionBlockNumber: block.ExecutionBlockNumber,
		ProposerIndex:        block.ProposerIndex,
		Withdrawals:          make([]*CachedWithdrawal, 0, len(block.Withdrawals)),
	}
	for _, w := range block.Withdrawals {
		cachedBlock.Withdrawals = append(cachedBlock.Withdrawals, &CachedWithdrawal{
			ValidatorIndex: w.ValidatorIndex,
			Amount:         w.Amount,
		})
	}

	m.cachedBeaconBlockByExecBlockHeight.Store(block.ExecutionBlockNumber, &cachedBlock)
	m.cachedBeaconBlock.Store(blockId, &cachedBlock)
	return &cachedBlock, true, nil
}

func (m *ServiceManager) pruneCachedBeaconBlocksService() {
	for {
		m.pruneCachedBeaconBlocks()
		time.Sleep(time.Minute)
	}
}
func (m *ServiceManager) pruneCachedBeaconBlocks() {
	var minHeight uint64 = math.MaxUint64
	m.srvs.Range(func(key string, srv *Service) bool {
		if srv != nil &&
			!srv.waitFirstNodeStakeEvent &&
			srv.minExecutionBlockHeight > 0 &&
			srv.minExecutionBlockHeight < minHeight {
			minHeight = srv.minExecutionBlockHeight
		}
		return true
	})

	if minHeight == math.MaxUint64 {
		return
	}

	var maxClearableBeaconBlockId uint64 = 0
	m.cachedBeaconBlockByExecBlockHeight.Range(func(blockNumber uint64, b *CachedBeaconBlock) bool {
		if b != nil &&
			b.ExecutionBlockNumber < minHeight {
			if maxClearableBeaconBlockId < b.BeaconBlockId {
				maxClearableBeaconBlockId = b.BeaconBlockId
			}
			m.cachedBeaconBlockByExecBlockHeight.Delete(blockNumber)
		}
		return true
	})

	m.cachedBeaconBlock.Range(func(beaconBlockId uint64, b *CachedBeaconBlock) bool {
		if b != nil &&
			b.BeaconBlockId <= maxClearableBeaconBlockId {
			m.cachedBeaconBlock.Delete(beaconBlockId)
			m.beaconBlockMutex.Delete(beaconBlockId)
		}
		return true
	})
}
