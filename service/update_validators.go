package service

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	ethpb "github.com/prysmaticlabs/prysm/v3/proto/eth/v1"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) updateValidatorsFromNetwork() error {
	eth1LatestBlock, err := s.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}
	if eth1LatestBlock <= s.latestBlockOfUpdateValidator {
		return nil
	}
	call := s.connection.CallOpts(big.NewInt(int64(eth1LatestBlock)))

	// 0. fetch new validators
	pubkeysLen, err := s.nodeDepositContract.GetPubkeysLength(call)
	if err != nil {
		return err
	}
	if len(s.validators) < int(pubkeysLen.Uint64()) {
		pubkeys, err := s.nodeDepositContract.GetPubkeys(call, big.NewInt(int64(len(s.validators))), pubkeysLen)
		if err != nil {
			return err
		}

		for _, pubkey := range pubkeys {
			key := hex.EncodeToString(pubkey)
			if _, exist := s.validators[key]; exist {
				return fmt.Errorf("validator %s duplicate", key)
			}

			pubkeyInfo, err := s.nodeDepositContract.PubkeyInfoOf(call, pubkey)
			if err != nil {
				return err
			}

			nodeLocal, exist := s.nodes[pubkeyInfo.Owner]
			if !exist {
				nodeInfo, err := s.nodeDepositContract.NodeInfoOf(call, pubkeyInfo.Owner)
				if err != nil {
					return err
				}

				node := Node{
					NodeAddress: pubkeyInfo.Owner,
					NodeType:    nodeInfo.NodeType,
				}

				s.nodes[node.NodeAddress] = &node
			}

			val := Validator{
				Pubkey:                pubkey,
				NodeAddress:           pubkeyInfo.Owner,
				DepositSignature:      pubkeyInfo.DepositSignature,
				NodeDepositAmountDeci: decimal.NewFromBigInt(pubkeyInfo.NodeDepositAmount, 0),
				NodeDepositAmount:     new(big.Int).Div(pubkeyInfo.NodeDepositAmount, big.NewInt(1e9)).Uint64(),
				DepositBlock:          pubkeyInfo.DepositBlock.Uint64(),
				ActiveEpoch:           0,
				EligibleEpoch:         0,
				ExitEpoch:             0,
				WithdrawableEpoch:     0,
				Balance:               0,
				EffectiveBalance:      0,
				NodeType:              nodeLocal.NodeType,
				Status:                pubkeyInfo.Status,
				ValidatorIndex:        0,
			}

			s.validators[key] = &val

		}
	}

	// 1. update validator status on network
	for _, val := range s.validators {
		if val.Status > utils.ValidatorStatusWithdrawUnmatch {
			continue
		}

		if val.Status == utils.ValidatorStatusStaked {
			continue
		}

		pubkeyInfo, err := s.nodeDepositContract.PubkeyInfoOf(call, val.Pubkey)
		if err != nil {
			return err
		}
		val.Status = pubkeyInfo.Status
	}

	s.latestBlockOfUpdateValidator = eth1LatestBlock
	return nil
}

func (s *Service) updateValidatorsFromBeacon() error {
	beaconHead, err := s.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch
	if finalEpoch <= s.latestEpochOfUpdateValidator {
		return nil
	}

	pubkeys := make([]types.ValidatorPubkey, 0)
	for _, val := range s.validators {
		pubkeys = append(pubkeys, types.ValidatorPubkey(val.Pubkey))
	}
	if len(pubkeys) == 0 {
		s.latestEpochOfUpdateValidator = finalEpoch
		return nil
	}

	validatorStatusMap, err := s.connection.GetValidatorStatuses(pubkeys, &beacon.ValidatorStatusOptions{
		Epoch: &finalEpoch,
	})
	if err != nil {
		return errors.Wrap(err, "syncValidatorLatestInfo GetValidatorStatuses failed")
	}

	logrus.WithFields(logrus.Fields{
		"validatorStatuses len": len(validatorStatusMap),
	}).Debug("validator statuses")

	for pubkey, status := range validatorStatusMap {
		pubkeyStr := pubkey.String()
		if status.Exists {
			// must exist here
			validator, exist := s.validators[pubkeyStr]
			if !exist {
				return fmt.Errorf("validator %s not exist", pubkeyStr)
			}

			updateBaseInfo := func() {
				// validator's info may be inited at any status
				validator.ActiveEpoch = status.ActivationEpoch
				validator.EligibleEpoch = status.ActivationEligibilityEpoch
				validator.ValidatorIndex = status.Index

				exitEpoch := status.ExitEpoch
				if exitEpoch == math.MaxUint64 {
					exitEpoch = 0
				}
				withdrawableEpoch := status.WithdrawableEpoch
				if withdrawableEpoch == math.MaxUint64 {
					withdrawableEpoch = 0
				}

				validator.ExitEpoch = exitEpoch
				validator.WithdrawableEpoch = withdrawableEpoch
			}

			updateBalance := func() {
				validator.Balance = status.Balance
				validator.EffectiveBalance = status.EffectiveBalance
			}

			switch status.Status {

			case ethpb.ValidatorStatus_PENDING_INITIALIZED, ethpb.ValidatorStatus_PENDING_QUEUED: // pending
				validator.Status = utils.ValidatorStatusWaiting
				validator.ValidatorIndex = status.Index

			case ethpb.ValidatorStatus_ACTIVE_ONGOING, ethpb.ValidatorStatus_ACTIVE_EXITING, ethpb.ValidatorStatus_ACTIVE_SLASHED: // active
				validator.Status = utils.ValidatorStatusActive
				if status.Slashed {
					validator.Status = utils.ValidatorStatusActiveSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_EXITED_UNSLASHED, ethpb.ValidatorStatus_EXITED_SLASHED: // exited
				validator.Status = utils.ValidatorStatusExited
				if status.Slashed {
					validator.Status = utils.ValidatorStatusExitedSlash
				}
				updateBaseInfo()
				updateBalance()
			case ethpb.ValidatorStatus_WITHDRAWAL_POSSIBLE: // withdrawable
				validator.Status = utils.ValidatorStatusWithdrawable
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawableSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_WITHDRAWAL_DONE: // withdrawdone
				validator.Status = utils.ValidatorStatusWithdrawDone
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawDoneSlash
				}
				updateBaseInfo()
				updateBalance()
			default:
				return fmt.Errorf("unsupported validator status %d", status.Status)
			}

		}
	}

	s.latestEpochOfUpdateValidator = finalEpoch

	return nil
}
