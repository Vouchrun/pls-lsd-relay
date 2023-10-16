package service

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) distributePriorityFee() error {

	latestDistributeHeight, targetEth1BlockHeight, shouldGoNext, err := s.checkStateForDistributePriorityFee()
	if err != nil {
		return errors.Wrap(err, "distributePriorityFee checkSyncState failed")
	}
	if !shouldGoNext {
		logrus.Debug("distributePriorityFee should not go next")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"latestDistributeHeight": latestDistributeHeight,
		"targetEth1BlockHeight":  targetEth1BlockHeight,
		"latestBlockOfSyncBlock": s.latestBlockOfSyncBlock,
	}).Debug("distributePriorityFee")

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, _, err := s.getUserNodePlatformFromPriorityFee(latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromPriorityFee failed")
	}

	// -----2 cal maxClaimableWithdrawIndex
	// find distribute withdrawals height as target block to cal this
	newMaxClaimableWithdrawIndex, err := s.calMaxClaimableWithdrawIndex(targetEth1BlockHeight, totalUserEthDeci)
	if err != nil {
		return errors.Wrap(err, "calMaxClaimableWithdrawIndex failed")
	}

	proposalId := utils.DistributeProposalId(utils.DistributeTypePriorityFee, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)))
	// check voted
	hasVoted, err := s.networkProposalContract.HasVoted(nil, proposalId, s.keyPair.CommonAddress())
	if err != nil {
		return fmt.Errorf("networkProposalContract.HasVoted err: %s", err)
	}
	if hasVoted {
		logrus.Debug("distributePriorityFee voted")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"targetEth1BlockHeight":        targetEth1BlockHeight,
		"totalUserEthDeci":             totalUserEthDeci.String(),
		"totalNodeEthDeci":             totalNodeEthDeci.String(),
		"totalPlatformEthDeci":         totalPlatformEthDeci.String(),
		"newMaxClaimableWithdrawIndex": newMaxClaimableWithdrawIndex,
	}).Info("Will DistributePriorityFee")

	// -----3 send vote tx
	return s.sendDistributeTx(utils.DistributeTypePriorityFee, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)), proposalId)
}

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, shouldGoNext, err)
func (s *Service) checkStateForDistributePriorityFee() (uint64, uint64, bool, error) {
	beaconHead, err := s.connection.Eth2BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	targetEpoch := (finalEpoch / s.distributePriorityFeeDuEpochs) * s.distributePriorityFeeDuEpochs
	targetEth1BlockHeight, err := s.getEpochStartBlocknumberWithCheck(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}
	logrus.Debugf("checkStateForDistributePriorityFee targetEth1Block: %d", targetEth1BlockHeight)

	latestDistributeHeight := s.latestDistributePriorityFeeHeight
	// init case
	if latestDistributeHeight == 0 {
		latestDistributeHeight = s.networkCreateBlock
	}

	if latestDistributeHeight >= targetEth1BlockHeight {
		logrus.Debugf("latestDistributeHeight: %d  targetEth1BlockHeight: %d", latestDistributeHeight, targetEth1BlockHeight)
		return 0, 0, false, nil
	}

	// wait sync block
	if targetEth1BlockHeight > s.latestBlockOfSyncBlock {
		return 0, 0, false, nil
	}

	return latestDistributeHeight, targetEth1BlockHeight, true, nil
}
