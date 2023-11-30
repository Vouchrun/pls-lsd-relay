package connection

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

type CachedConnection struct {
	*Connection
	stop chan struct{}

	beaconHead    beacon.BeaconHead
	beaconHeadErr error
}

func (c *CachedConnection) Start() {
	utils.SafeGo(c.syncBeaconHead)
}

func (c *CachedConnection) Stop() {
	close(c.stop)
}

func (c *CachedConnection) BeaconHead() (beacon.BeaconHead, error) {
	return c.eth2Client.GetBeaconHead()
}

// internal jobs

func (c *CachedConnection) syncBeaconHead() {
	for {
		select {
		case <-c.stop:
			return
		default:
			c.beaconHead, c.beaconHeadErr = retry.DoWithData(c.eth2Client.GetBeaconHead,
				retry.Delay(time.Second*2), retry.Attempts(150))
		}
		time.Sleep(12 * time.Second)
	}
}
