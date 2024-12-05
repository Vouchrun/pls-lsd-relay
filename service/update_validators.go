package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	ethpb "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	node_deposit "github.com/stafiprotocol/eth-lsd-relay/bindings/NodeDeposit"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) updateValidatorsFromNetwork() error {
	// 0. fetch new Nodes
	jobResult, err := s.connection.SubmitLatestCallJob(s.nodeDepositContract.NewGetNodesLengthMultiCall())
	if err != nil {
		return err
	}
	call := jobResult.Get()
	if call.Failed {
		return fmt.Errorf("nodeDepositContract.GetNodesLength failed: %w height: %d", call.Err, call.BlockNumber)
	}
	eth1LatestBlock := call.BlockNumber
	if eth1LatestBlock <= s.latestBlockOfUpdateValidator {
		return nil
	}
	opts := s.connection.CallOpts(big.NewInt(int64(eth1LatestBlock)))

	nodesLength := call.Outputs.(*node_deposit.GetNodesLengthMultiCallOutput).Length
	if nodesLength.Uint64() == 0 {
		return nil
	}

	s.log.WithFields(logrus.Fields{
		"eth1LatestBlock": eth1LatestBlock,
		"nodesLenOnChain": nodesLength.Int64(),
	}).Debug("updateValidatorsFromNetwork")

	if len(s.nodes) < int(nodesLength.Int64()) {
		nodesOnChain, err := s.nodeDepositContract.GetNodes(opts, big.NewInt(0), nodesLength)
		if err != nil {
			return fmt.Errorf("nodeDepositContract.GetNodes failed: %w", err)
		}
		newNodes := nodesOnChain[len(s.nodes):]
		for _, nodeAddress := range newNodes {
			nodeInfo, err := s.nodeDepositContract.NodeInfoOf(opts, nodeAddress)
			if err != nil {
				return err
			}
			pubkeys, err := s.nodeDepositContract.GetPubkeysOfNode(opts, nodeAddress)
			if err != nil {
				return err
			}
			newVals, err := s.fetchNewVals(opts, pubkeys)
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
	nodesPubkeyList, err := s.nodeDepositContract.GetPubkeysOfNodes(opts, lo.Keys(s.nodes))
	if err != nil {
		return errors.Wrapf(err, "get pubkeys of nodes, len: %d", len(s.nodes))
	}
	for addr, node := range s.nodes {
		pubkeys := nodesPubkeyList[addr]

		s.log.WithFields(logrus.Fields{
			"node":              node.NodeAddress,
			"pubkeysLenOnChain": len(pubkeys),
		}).Debug("updateValidatorsFromNetwork")

		if len(pubkeys) > int(node.PubkeyNumber) {
			newPubkeys := pubkeys[int(node.PubkeyNumber):]
			newVals, err := s.fetchNewVals(opts, newPubkeys)
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
	validValidatorPubkeys := make([][]byte, 0, len(s.validators))
	for _, val := range s.validators {
		if val.Status > utils.ValidatorStatusWithdrawUnmatch {
			continue
		}

		if val.Status == utils.ValidatorStatusStaked {
			continue
		}

		validValidatorPubkeys = append(validValidatorPubkeys, val.Pubkey)
	}
	pubkeyInfo, err := s.nodeDepositContract.GetPubkeyInfoList(opts, validValidatorPubkeys)
	if err != nil {
		return errors.Wrapf(err, "get pubkey info list, len: %d", len(validValidatorPubkeys))
	}
	for pubkeyStr, info := range pubkeyInfo {
		s.validators[pubkeyStr].Status = info.Status
	}

	s.latestBlockOfUpdateValidator = eth1LatestBlock
	return nil
}

func (s *Service) updateValidatorsFromBeacon() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
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
		if val.Status == 3 || val.Status > 4 {
			pubkeys = append(pubkeys, types.ValidatorPubkey(val.Pubkey))
		}
	}
	if len(pubkeys) == 0 {
		s.latestEpochOfUpdateValidator = finalEpoch
		return nil
	}

	validatorStatusMap, err := s.connection.GetValidatorStatuses(ctx, pubkeys, &beacon.ValidatorStatusOptions{
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
			validator.Status, err = mapValidatorStatus(&status)
			if err != nil {
				return fmt.Errorf("unsupported validator status %d", status.Status)
			}
			switch validator.Status {
			case utils.ValidatorStatusWaiting:
				validator.ValidatorIndex = status.Index
			case utils.ValidatorStatusActive, utils.ValidatorStatusActiveSlash,
				utils.ValidatorStatusExited, utils.ValidatorStatusExitedSlash,
				utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawableSlash,
				utils.ValidatorStatusWithdrawDone, utils.ValidatorStatusWithdrawDoneSlash:
				updateBaseInfo()
				updateBalance()
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

func mapValidatorStatus(status *beacon.ValidatorStatus) (uint8, error) {
	switch status.Status {
	case ethpb.ValidatorStatus_PENDING_INITIALIZED, ethpb.ValidatorStatus_PENDING_QUEUED: // pending
		return utils.ValidatorStatusWaiting, nil
	case ethpb.ValidatorStatus_ACTIVE_ONGOING, ethpb.ValidatorStatus_ACTIVE_EXITING, ethpb.ValidatorStatus_ACTIVE_SLASHED: // active
		if status.Slashed {
			return utils.ValidatorStatusActiveSlash, nil
		}
		return utils.ValidatorStatusActive, nil
	case ethpb.ValidatorStatus_EXITED_UNSLASHED, ethpb.ValidatorStatus_EXITED_SLASHED: // exited
		if status.Slashed {
			return utils.ValidatorStatusExitedSlash, nil
		}
		return utils.ValidatorStatusExited, nil
	case ethpb.ValidatorStatus_WITHDRAWAL_POSSIBLE: // withdrawable
		if status.Slashed {
			return utils.ValidatorStatusWithdrawableSlash, nil
		}
		return utils.ValidatorStatusWithdrawable, nil
	case ethpb.ValidatorStatus_WITHDRAWAL_DONE: // withdrawdone
		if status.Slashed {
			return utils.ValidatorStatusWithdrawDoneSlash, nil
		}
		return utils.ValidatorStatusWithdrawDone, nil
	default:
		return 0, fmt.Errorf("unsupported validator status %d", status.Status)
	}
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
			NodeDepositAmount:     new(big.Int).Div(pubkeyInfo.NodeDepositAmount, big.NewInt(1e9)).Uint64(), // convert wei to Gwei
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
