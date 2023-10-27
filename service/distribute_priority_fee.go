package service

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) distributePriorityFee() error {

	latestDistributionHeight, targetEth1BlockHeight, shouldGoNext, err := s.checkStateForDistributePriorityFee()
	if err != nil {
		return errors.Wrap(err, "distributePriorityFee checkSyncState failed")
	}
	if !shouldGoNext {
		logrus.Debug("distributePriorityFee should not go next")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"latestDistributionHeight": latestDistributionHeight,
		"targetEth1BlockHeight":    targetEth1BlockHeight,
		"latestBlockOfSyncBlock":   s.latestBlockOfSyncBlock,
	}).Debug("distributePriorityFee")

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, _, err := s.getUserNodePlatformFromPriorityFee(latestDistributionHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromPriorityFee failed")
	}

	// -----2 cal maxClaimableWithdrawalIndex
	// find distribute withdrawals height as target block to cal this
	newMaxClaimableWithdrawalIndex, err := s.calMaxClaimableWithdrawalIndex(targetEth1BlockHeight, totalUserEthDeci)
	if err != nil {
		return errors.Wrap(err, "calMaxClaimableWithdrawalIndex failed")
	}

	// -----3 send vote tx
	return s.sendDistributeTx(utils.DistributeTypePriorityFee, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawalIndex)))
}

// check sync and vote state
// return (latestDistributionHeight, targetEth1Blocknumber, shouldGoNext, err)
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

	latestDistributionHeight := s.latestDistributionPriorityFeeHeight
	// init case
	if latestDistributionHeight == 0 {
		latestDistributionHeight = s.networkCreateBlock
	}

	if latestDistributionHeight >= targetEth1BlockHeight {
		logrus.Debugf("latestDistributionHeight: %d  targetEth1BlockHeight: %d", latestDistributionHeight, targetEth1BlockHeight)
		return 0, 0, false, nil
	}

	// wait sync block
	if targetEth1BlockHeight > s.latestBlockOfSyncBlock {
		return 0, 0, false, nil
	}

	return latestDistributionHeight, targetEth1BlockHeight, true, nil
}
