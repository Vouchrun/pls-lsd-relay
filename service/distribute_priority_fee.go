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
		return errors.Wrap(err, "distributeWithdrawals checkSyncState failed")
	}

	if !shouldGoNext {
		logrus.Debug("distributeWithdrawals should not go next")
		return nil
	}
	logrus.WithFields(logrus.Fields{
		"latestDistributeHeight": latestDistributeHeight,
		"targetEth1BlockHeight":  targetEth1BlockHeight,
	}).Debug("distributePriorityFee")

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, err := s.getUserNodePlatformFromPriorityFee(latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromPriorityFee failed")
	}

	// -----2 cal maxClaimableWithdrawIndex
	// todo find distribute withdrawals height as target block to cal this
	newMaxClaimableWithdrawIndex, err := s.calMaxClaimableWithdrawIndex(targetEth1BlockHeight, totalUserEthDeci)
	if err != nil {
		return errors.Wrap(err, "calMaxClaimableWithdrawIndex failed")
	}

	// check voted
	hasVoted, err := s.networkProposalContract.HasVoted(nil, utils.DistributeProposalId(utils.DistributeTypePriorityFee, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalAmountDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex))), s.keyPair.CommonAddress())
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
	return s.sendDistributeTx(utils.DistributeTypeWithdrawals, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)))
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
	targetEth1BlockHeight, err := s.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	logrus.Debugf("targetEth1Block %d", targetEth1BlockHeight)

	latestDistributeHeight, err := s.networkWithdrawContract.LatestDistributePriorityFeeHeight(nil)
	if err != nil {
		return 0, 0, false, err
	}
	// init case
	if latestDistributeHeight.Uint64() == 0 {
		latestDistributeHeight = big.NewInt(int64(s.networkCreateBlock))
	}

	if latestDistributeHeight.Uint64() >= targetEth1BlockHeight {
		logrus.Debug("latestDistributeHeight.Uint64() >= targetEth1BlockHeight")
		return 0, 0, false, nil
	}

	return latestDistributeHeight.Uint64(), targetEth1BlockHeight, true, nil
}
