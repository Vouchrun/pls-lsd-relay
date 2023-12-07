package connection

import (
	"context"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

type CachedConnection struct {
	*Connection
	stop chan struct{}

	// cache data
	beaconHead    beacon.BeaconHead
	beaconHeadErr error

	chainId    *big.Int
	eth2Config *beacon.Eth2Config

	beaconBlockCache *lru.Cache[uint64, beacon.BeaconBlock]
	beaconBlockMutex *utils.KeyedMutex[uint64]
}

func NewCachedConnection(conn *Connection) (*CachedConnection, error) {
	beaconBlockCache, err := lru.New[uint64, beacon.BeaconBlock](1024 * 1000)
	if err != nil {
		return nil, err
	}
	cc := CachedConnection{
		Connection:       conn,
		beaconBlockCache: beaconBlockCache,
	}
	return &cc, nil
}

func (c *CachedConnection) Start() error {
	err := c.cacheChainID()
	if err != nil {
		return err
	}

	c.syncBeaconHead()
	utils.SafeGo(c.syncBeaconHeadService)

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
		cfg, err := retry.DoWithData(c.Eth2Client().GetEth2Config,
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
			c.syncBeaconHead()
		}
		time.Sleep(12 * time.Second)
	}
}

func (c *CachedConnection) syncBeaconHead() {
	c.beaconHead, c.beaconHeadErr = retry.DoWithData(c.eth2Client.GetBeaconHead,
		retry.Delay(time.Second*2), retry.Attempts(150))
}

func (c *CachedConnection) cacheChainID() (err error) {
	c.chainId, err = retry.DoWithData(func() (*big.Int, error) { return c.eth1Client.ChainID(context.Background()) },
		retry.Delay(time.Second*2), retry.Attempts(150))
	return
}

func (c *CachedConnection) GetBeaconBlock(blockId uint64) (beacon.BeaconBlock, bool, error) {
	// lock block for concurrent
	unlock := c.beaconBlockMutex.Lock(blockId)
	defer unlock()

	// load from cache
	block, ok := c.beaconBlockCache.Get(blockId)
	if ok {
		return block, true, nil
	}

	block, ok, err := c.Connection.GetBeaconBlock(blockId)
	if err != nil {
		return beacon.BeaconBlock{}, ok, err
	}
	if ok {
		c.beaconBlockCache.Add(blockId, block)
	}

	return block, ok, nil
}
