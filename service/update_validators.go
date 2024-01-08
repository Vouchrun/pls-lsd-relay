package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	ethpb "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
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

	// 0. fetch new Nodes
	nodesLength, err := s.nodeDepositContract.GetNodesLength(call)
	if err != nil {
		return fmt.Errorf("nodeDepositContract.GetNodesLength failed: %w", err)
	}
	if nodesLength.Uint64() == 0 {
		return nil
	}

	nodesOnChain, err := s.nodeDepositContract.GetNodes(call, big.NewInt(0), nodesLength)
	if err != nil {
		return fmt.Errorf("nodeDepositContract.GetNodes failed: %w", err)
	}

	s.log.WithFields(logrus.Fields{
		"eth1LatestBlock": eth1LatestBlock,
		"nodesLenOnChain": len(nodesOnChain),
	}).Debug("updateValidatorsFromNetwork")

	if len(s.nodes) < len(nodesOnChain) {
		newNodes := nodesOnChain[len(s.nodes):]
		for _, nodeAddress := range newNodes {
			nodeInfo, err := s.nodeDepositContract.NodeInfoOf(call, nodeAddress)
			if err != nil {
				return err
			}
			pubkeys, err := s.nodeDepositContract.GetPubkeysOfNode(call, nodeAddress)
			if err != nil {
				return err
			}
			newVals, err := s.fetchNewVals(call, pubkeys)
			if err != nil {
				return errors.Wrapf(err, "new node fetchNewVals")
			}

			// cache validators
			for key, val := range newVals {
				s.validators[key] = val
			}
			// cache node
			s.nodes[nodeAddress] = &Node{
				NodeAddress:  nodeAddress,
				NodeType:     nodeInfo.NodeType,
				PubkeyNumber: uint64(len(newVals)),
			}
		}
	}

	// 1 fetch node's new pubkey
	for _, node := range s.nodes {
		pubkeys, err := s.nodeDepositContract.GetPubkeysOfNode(call, node.NodeAddress)
		if err != nil {
			return err
		}

		s.log.WithFields(logrus.Fields{
			"node":              node.NodeAddress,
			"pubkeysLenOnChain": len(pubkeys),
		}).Debug("updateValidatorsFromNetwork")

		if len(pubkeys) > int(node.PubkeyNumber) {
			newPubkeys := pubkeys[int(node.PubkeyNumber):]
			newVals, err := s.fetchNewVals(call, newPubkeys)
			if err != nil {
				return errors.Wrapf(err, "new pubkey fetchNewVals")
			}

			// cache validators
			for key, val := range newVals {
				s.validators[key] = val
			}
			// cache node
			node.PubkeyNumber += uint64(len(newVals))
		}
	}

	// 2. update validator status on network
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
	beaconHead, err := s.connection.BeaconHead()
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

	s.log.WithFields(logrus.Fields{
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
				validator.ExitEpoch = exitEpoch

				withdrawableEpoch := status.WithdrawableEpoch
				if withdrawableEpoch == math.MaxUint64 {
					withdrawableEpoch = 0
				}
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

	// cache validators by index
	s.validatorsByIndexMutex.Lock()
	for _, validator := range s.validators {
		if validator.ValidatorIndex > 0 {
			s.validatorsByIndex[validator.ValidatorIndex] = validator
		}
	}
	s.validatorsByIndexMutex.Unlock()

	s.latestEpochOfUpdateValidator = finalEpoch

	return nil
}

func (s *Service) fetchNewVals(call *bind.CallOpts, pubkeys [][]byte) (map[string]*Validator, error) {
	newVals := make(map[string]*Validator)
	for _, pubkey := range pubkeys {
		key := hex.EncodeToString(pubkey)
		if _, exist := s.validators[key]; exist {
			return nil, fmt.Errorf("validator %s duplicate", key)
		}

		pubkeyInfo, err := s.nodeDepositContract.PubkeyInfoOf(call, pubkey)
		if err != nil {
			return nil, err
		}

		nodeLocal, exist := s.nodes[pubkeyInfo.Owner]
		if !exist {
			nodeInfo, err := s.nodeDepositContract.NodeInfoOf(call, pubkeyInfo.Owner)
			if err != nil {
				return nil, err
			}

			node := Node{
				NodeAddress: pubkeyInfo.Owner,
				NodeType:    nodeInfo.NodeType,
			}

			s.nodes[node.NodeAddress] = &node

			nodeLocal = &node
		}

		filterBlock := pubkeyInfo.DepositBlock.Uint64()
		depositedIter, err := s.nodeDepositContract.FilterDeposited(&bind.FilterOpts{
			Start:   filterBlock,
			End:     &filterBlock,
			Context: context.Background(),
		})
		if err != nil {
			return nil, err
		}

		var depositSig []byte
		for depositedIter.Next() {
			if bytes.Equal(depositedIter.Event.Pubkey, pubkey) {
				depositSig = depositedIter.Event.ValidatorSignature
				break
			}
		}

		if len(depositSig) == 0 {
			return nil, fmt.Errorf("depositSignature empty, val pubkey: %s", key)
		}

		val := Validator{
			Pubkey:                pubkey,
			NodeAddress:           pubkeyInfo.Owner,
			DepositSignature:      depositSig,
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
		newVals[key] = &val
	}

	if len(pubkeys) != len(newVals) {
		return nil, fmt.Errorf("fetchNewVals, pubkeys length: %d not match newVals length: %d", len(pubkeys), len(newVals))
	}

	return newVals, nil
}
