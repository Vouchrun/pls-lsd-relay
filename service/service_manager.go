package service

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	xsync "github.com/puzpuzpuz/xsync/v3"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	lsd_network_factory "github.com/stafiprotocol/eth-lsd-relay/bindings/LsdNetworkFactory"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

type ServiceManager struct {
	stop       chan struct{}
	connection *connection.Connection
	srvs       *xsync.MapOf[string, *Service]
}

func NewServiceManager(cfg *config.Config, keyPair *secp256k1.Keypair) (*ServiceManager, error) {
	return &ServiceManager{
		stop: make(chan struct{}),
		srvs: xsync.NewMapOf[string, *Service](),
	}, nil
}

func (m *ServiceManager) Start() error {
	err := retry.Do(m.syncEntrustedLsdTokens, retry.Attempts(utils.RetryLimit))
	if err != nil {
		return err
	}

	utils.SafeGo(m.startSyncService)

	return nil
	// step1. get all entrusted lsd tokens
	// blockNumber, err := m.connection.Eth1LatestBlock()
	// if err != nil {
	// 	return err
	// }

	// call sync entrusted lsd token manually
	// lsdNetworkFactoryContract, err := lsd_network_factory.NewLsdNetworkFactory(common.Address{}, m.connection.Eth1Client())
	// if err != nil {
	// 	return err
	// }
	// tokens, err := lsdNetworkFactoryContract.GetEntrustedLsdTokens(nil)
	// if err != nil {
	// 	return err
	// }

	// // step2. start all entrusted lsd tokens
	// for _, token := range tokens {
	// 	// start
	// 	m.srvs.Store(token.String(), &Service{})
	// }

	// step3. start sync entrusted lsd tokens
}

func (m *ServiceManager) Stop() {
	close(m.stop)
}

func (m *ServiceManager) startSyncService() {
	logrus.Info("start listening new entrusted lsd token service")

	for {
		select {
		case <-m.stop:
			logrus.Info("sync entrusted lsd token task has stopped")
			return
		default:
			// sync new entrusted lsd tokens
			err := retry.Do(m.syncEntrustedLsdTokens, retry.Attempts(utils.RetryLimit))
			if err != nil {
				logrus.Error("task has stopped")
				utils.ShutdownRequestChannel <- struct{}{}
				return
			}

			time.Sleep(12 * time.Second)
		}
	}
}

func (m *ServiceManager) syncEntrustedLsdTokens() error {
	lsdNetworkFactoryContract, err := lsd_network_factory.NewLsdNetworkFactory(common.Address{}, m.connection.Eth1Client())
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
			m.srvs.Store(token, &Service{}) // fixme
		}
	}

	m.srvs.Range(func(token string, srv *Service) bool {
		if !lo.Contains(tokenList, token) {
			// remove entrusted lsd token
			srv.Stop()
			m.srvs.Delete(token)
		}
		return true
	})

	return nil
}
