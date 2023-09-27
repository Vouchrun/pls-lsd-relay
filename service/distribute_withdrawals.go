package service

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) distributeWithdrawals() error {

	latestDistributeHeight, targetEth1BlockHeight, shouldGoNext, err := s.checkStateForDistributeWithdraw()
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
		"latestBlockOfSyncBlock": s.latestBlockOfSyncBlock,
	}).Debug("distributeWithdrawals")

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, _, err := s.getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromWithdrawals failed")
	}

	// -----2 cal maxClaimableWithdrawIndex
	newMaxClaimableWithdrawIndex, err := s.calMaxClaimableWithdrawIndex(targetEth1BlockHeight, totalUserEthDeci)
	if err != nil {
		return errors.Wrap(err, "calMaxClaimableWithdrawIndex failed")
	}

	// check voted
	hasVoted, err := s.networkProposalContract.HasVoted(nil, utils.DistributeProposalId(utils.DistributeTypeWithdrawals, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex))), s.keyPair.CommonAddress())
	if err != nil {
		return fmt.Errorf("networkProposalContract.HasVoted err: %s", err)
	}
	if hasVoted {
		logrus.Debug("networkProposalContract voted")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"targetEth1BlockHeight":        targetEth1BlockHeight,
		"totalUserEthDeci":             totalUserEthDeci.String(),
		"totalNodeEthDeci":             totalNodeEthDeci.String(),
		"totalPlatformEthDeci":         totalPlatformEthDeci.String(),
		"newMaxClaimableWithdrawIndex": newMaxClaimableWithdrawIndex,
	}).Info("Will DistributeWithdrawals")

	// -----3 send vote tx
	return s.sendDistributeTx(utils.DistributeTypeWithdrawals, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)))
}

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, shouldGoNext, err)
func (s *Service) checkStateForDistributeWithdraw() (uint64, uint64, bool, error) {
	beaconHead, err := s.connection.Eth2BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	targetEpoch := (finalEpoch / s.distributeWithdrawalsDuEpochs) * s.distributeWithdrawalsDuEpochs
	targetEth1BlockHeight, err := s.getEpochStartBlocknumberWithCheck(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	logrus.Debugf("targetEth1Block %d", targetEth1BlockHeight)

	latestDistributeHeight := s.latestDistributeWithdrawalsHeight
	if err != nil {
		return 0, 0, false, err
	}
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

func (s *Service) calMaxClaimableWithdrawIndex(targetEth1BlockHeight uint64, totalUserEthDeci decimal.Decimal) (uint64, error) {
	calOpts := s.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))
	maxClaimableWithdrawIndex, err := s.networkWithdrawContract.MaxClaimableWithdrawIndex(calOpts)
	if err != nil {
		return 0, err
	}
	// nextWithdrawIndex <= real value
	nextWithdrawIndex, err := s.networkWithdrawContract.NextWithdrawIndex(calOpts)
	if err != nil {
		return 0, err
	}
	totalMissingAmountForWithdraw, err := s.networkWithdrawContract.TotalMissingAmountForWithdraw(calOpts)
	if err != nil {
		return 0, err
	}
	newMaxClaimableWithdrawIndex := uint64(0)
	totalMissingAmountForWithdrawDeci := decimal.NewFromBigInt(totalMissingAmountForWithdraw, 0)
	if totalMissingAmountForWithdrawDeci.LessThanOrEqual(totalUserEthDeci) {
		if nextWithdrawIndex.Uint64() >= 1 {
			newMaxClaimableWithdrawIndex = nextWithdrawIndex.Uint64() - 1
		}
	} else {
		willMissingAmountDeci := totalMissingAmountForWithdrawDeci.Sub(totalUserEthDeci)
		if nextWithdrawIndex.Uint64() >= 1 {
			latestUsersWaitAmountDeci := decimal.Zero
			for i := nextWithdrawIndex.Uint64() - 1; i > maxClaimableWithdrawIndex.Uint64(); i-- {
				stakerWithdrawal, exist := s.stakerWithdrawals[i]
				if !exist {
					return 0, fmt.Errorf("stakerWithdrawal %d not exist", i)
				}

				// skip instantly withdrawal
				if stakerWithdrawal.ClaimedBlockNumber == stakerWithdrawal.BlockNumber {
					continue
				}

				ethAmountDeci, err := decimal.NewFromString(stakerWithdrawal.EthAmount)
				if err != nil {
					return 0, err
				}
				latestUsersWaitAmountDeci = latestUsersWaitAmountDeci.Add(ethAmountDeci)
				if latestUsersWaitAmountDeci.GreaterThan(willMissingAmountDeci) {
					if i >= 1 {
						newMaxClaimableWithdrawIndex = i - 1
					}
					break
				}
			}
		}
	}
	if newMaxClaimableWithdrawIndex < maxClaimableWithdrawIndex.Uint64() {
		newMaxClaimableWithdrawIndex = maxClaimableWithdrawIndex.Uint64()
	}

	return newMaxClaimableWithdrawIndex, nil
}

func (s *Service) sendDistributeTx(distributeType uint8, targetEth1BlockHeight, totalUserEth, totalNodeEth, totalPlatformEth, newMaxClaimableWithdrawIndex *big.Int) error {
	err := s.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer s.connection.UnlockTxOpts()

	tx, err := s.networkWithdrawContract.Distribute(s.connection.TxOpts(), distributeType, targetEth1BlockHeight,
		totalUserEth, totalNodeEth, totalPlatformEth, newMaxClaimableWithdrawIndex)
	if err != nil {
		return err
	}

	logrus.Infof("send Distribute tx hash: %s, type: %d", tx.Hash().String(), distributeType)

	return s.waitTxOk(tx.Hash())
}
