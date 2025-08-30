package service

import (
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
		s.log.Debug("distributePriorityFee should not go next")
		return nil
	}

	s.log.WithFields(logrus.Fields{
		"latestDistributeHeight": latestDistributeHeight,
		"targetEth1BlockHeight":  targetEth1BlockHeight,
		"latestBlockOfSyncBlock": s.latestBlockOfSyncBlock,
	}).Debug("distributePriorityFee")

	// ----1 cal eth(from withdrawals) of user/node/platform
	log := s.log.WithFields(logrus.Fields{
		"distributePriorityFee": true,
	})
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, _, err := s.getUserNodePlatformFromPriorityFee(log, latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromPriorityFee failed")
	}

	// -----2 cal maxClaimableWithdrawIndex
	// find distribute withdrawals height as target block to cal this
	newMaxClaimableWithdrawIndex, err := s.calMaxClaimableWithdrawIndex(targetEth1BlockHeight, totalUserEthDeci)
	if err != nil {
		return errors.Wrap(err, "calMaxClaimableWithdrawIndex failed")
	}

	// -----3 send vote tx
	return s.sendDistributeTx(utils.DistributeTypePriorityFee, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)))
}

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, shouldGoNext, err)
func (s *Service) checkStateForDistributePriorityFee() (uint64, uint64, bool, error) {
	beaconHead, err := s.connection.BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	targetEpoch := (finalEpoch / s.distributePriorityFeeDuEpochs) * s.distributePriorityFeeDuEpochs
	targetEth1BlockHeight, err := s.getEpochStartBlocknumberWithCheck(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}
	s.log.Debugf("checkStateForDistributePriorityFee targetEth1Block: %d", targetEth1BlockHeight)

	latestDistributeHeight := s.latestDistributePriorityFeeHeight
	// init case
	if latestDistributeHeight == 0 {
		latestDistributeHeight = s.startAtBlock
	}

	if latestDistributeHeight >= targetEth1BlockHeight {
		s.log.Debugf("latestDistributeHeight: %d  targetEth1BlockHeight: %d", latestDistributeHeight, targetEth1BlockHeight)
		return 0, 0, false, nil
	}

	// wait sync block
	if targetEth1BlockHeight > s.latestBlockOfSyncBlock {
		return 0, 0, false, nil
	}

	return latestDistributeHeight, targetEth1BlockHeight, true, nil
}
