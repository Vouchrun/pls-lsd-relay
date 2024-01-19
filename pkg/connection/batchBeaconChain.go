package connection

import (
	"time"

	"github.com/samber/lo"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	gomicrobee "github.com/stephennancekivell/go-micro-bee"
)

func (c *Connection) initBatchBeaconChain() {
	c.latestValicatorStatusMicrobeeSystem = gomicrobee.NewSystem[types.ValidatorPubkey, BeaconValidatorStatus](
		c.doValidatorStatusRequest,
		50,
		5*time.Second,
	)
}

func (c *Connection) SubmitFetchValidatorStatusJob(k types.ValidatorPubkey) (gomicrobee.JobResult[BeaconValidatorStatus], error) {
	return c.latestValicatorStatusMicrobeeSystem.Submit(k)
}

type BeaconValidatorStatus struct {
	beacon.ValidatorStatus
	BatchErr error
}

func (c *Connection) doValidatorStatusRequest(pubkeys []types.ValidatorPubkey) []BeaconValidatorStatus {
	validatorStatus, err := c.Eth2Client().GetValidatorStatuses(pubkeys, nil)
	if err != nil {
		return lo.Map(pubkeys, func(call types.ValidatorPubkey, index int) BeaconValidatorStatus {
			return BeaconValidatorStatus{BatchErr: err}
		})
	}

	return lo.Map(pubkeys, func(call types.ValidatorPubkey, index int) BeaconValidatorStatus {
		status, ok := validatorStatus[call]
		if !ok {
			return BeaconValidatorStatus{}
		}
		return BeaconValidatorStatus{ValidatorStatus: status}
	})
}
