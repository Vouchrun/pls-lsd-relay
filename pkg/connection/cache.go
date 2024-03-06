package connection

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

type CachedConnection struct {
	*Connection
	stop chan struct{}

	// cache data
	beaconHead               beacon.BeaconHead
	beaconHeadErr            error
	eth1LatestBlockNumber    uint64
	eth1LatestBlockNumberErr error

	chainId              *big.Int
	eth2Config           *beacon.Eth2Config
	validatorStatusCache sync.Map
}

func NewCachedConnection(conn *Connection) (*CachedConnection, error) {
	cc := CachedConnection{
		Connection:           conn,
		stop:                 make(chan struct{}),
		validatorStatusCache: sync.Map{},
	}
	return &cc, nil
}

func (c *CachedConnection) Start() error {
	err := c.cacheChainID()
	if err != nil {
		return err
	}

	if err = utils.ExecuteFns(
		c.syncBeaconHead,
		c.syncEth1LatestBlockNumber,
	); err != nil {
		return err
	}

	utils.SafeGo(c.syncBeaconHeadService)
	utils.SafeGo(c.syncEth1LatestBlockService)

	return nil
}

func (c *CachedConnection) Stop() {
	close(c.stop)
}

func (c *CachedConnection) BeaconHead() (beacon.BeaconHead, error) {
	return c.beaconHead, c.beaconHeadErr
}

func (c *CachedConnection) ChainID() (*big.Int, error) {
	return c.chainId, nil
}

func (c *CachedConnection) Eth2Config() (beacon.Eth2Config, error) {
	if c.eth2Config == nil {
		cfg, err := retry.DoWithData(c.GetEth2Config,
			retry.Delay(time.Second*2), retry.Attempts(150))
		if err != nil {
			return beacon.Eth2Config{}, err
		}
		c.eth2Config = &cfg
	}
	return *c.eth2Config, nil
}

// internal jobs

func (c *CachedConnection) syncBeaconHeadService() {
	for {
		select {
		case <-c.stop:
			return
		default:
			if err := c.syncBeaconHead(); err != nil {
				logrus.Errorf("connection cache: fail to sync beacon head: %s", err.Error())
			}
		}
		time.Sleep(12 * time.Second)
	}
}

func (c *CachedConnection) syncBeaconHead() error {
	c.beaconHead, c.beaconHeadErr = retry.DoWithData(c.GetBeaconHead,
		retry.Delay(time.Second*2), retry.Attempts(5))
	return c.beaconHeadErr
}

func (c *CachedConnection) Eth1LatestBlock() (uint64, error) {
	return c.eth1LatestBlockNumber, c.eth1LatestBlockNumberErr
}

func (c *CachedConnection) syncEth1LatestBlockService() {
	for {
		select {
		case <-c.stop:
			return
		default:
			if err := c.syncEth1LatestBlockNumber(); err != nil {
				logrus.Errorf("connection cache: fail to sync eth1 latest block number: %s", err.Error())
			}
		}
		time.Sleep(12 * time.Second)
	}
}

// LatestBlock returns the latest block from the current chain
func (c *CachedConnection) syncEth1LatestBlockNumber() error {
	c.eth1LatestBlockNumber, c.eth1LatestBlockNumberErr = retry.DoWithData(c.Connection.Eth1LatestBlock,
		retry.Delay(time.Second*2), retry.Attempts(5))
	return c.eth1LatestBlockNumberErr
}

func (c *CachedConnection) cacheChainID() (err error) {
	c.chainId, err = retry.DoWithData(func() (*big.Int, error) { return c.eth1Client.ChainID(context.Background()) },
		retry.Delay(time.Second*2), retry.Attempts(150))
	return
}

func (c *CachedConnection) GetValidatorStatus(ctx context.Context, pubkey types.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (validatorStatus beacon.ValidatorStatus, err error) {
	cacheKey := ""
	if opts != nil {
		if opts.Epoch != nil {
			cacheKey += fmt.Sprintf("%d", *opts.Epoch)
		}
		if opts.Slot != nil {
			cacheKey += fmt.Sprintf("_%d", *opts.Slot)
		}
	}
	if cacheKey != "" {
		cacheKey += "_" + pubkey.Hex()
		status, ok := c.validatorStatusCache.Load(cacheKey)
		if ok {
			return status.(beacon.ValidatorStatus), nil
		}
	}

	status, err := c.Connection.GetValidatorStatus(ctx, pubkey, opts)
	if err != nil {
		return status, err
	}

	if cacheKey != "" {
		c.validatorStatusCache.Store(cacheKey, status)
	}

	return status, nil
}
