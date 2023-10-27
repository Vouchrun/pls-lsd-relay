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

	latestDistributionHeight, targetEth1BlockHeight, shouldGoNext, err := s.checkStateForDistributeWithdraw()
	if err != nil {
		return errors.Wrap(err, "distributeWithdrawals checkSyncState failed")
	}
	if !shouldGoNext {
		logrus.Debug("distributeWithdrawals should not go next")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"latestDistributionHeight": latestDistributionHeight,
		"targetEth1BlockHeight":    targetEth1BlockHeight,
		"latestBlockOfSyncBlock":   s.latestBlockOfSyncBlock,
	}).Debug("distributeWithdrawals")

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, _, err := s.getUserNodePlatformFromWithdrawals(latestDistributionHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromWithdrawals failed")
	}

	// -----2 cal newMaxClaimableWithdrawalIndex
	newMaxClaimableWithdrawalIndex, err := s.calMaxClaimableWithdrawalIndex(targetEth1BlockHeight, totalUserEthDeci)
	if err != nil {
		return errors.Wrap(err, "calMaxClaimableWithdrawIndex failed")
	}

	// -----3 send vote tx
	return s.sendDistributeTx(utils.DistributeTypeWithdrawals, big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawalIndex)))
}

// check sync and vote state
// return (latestDistributionHeight, targetEth1Blocknumber, shouldGoNext, err)
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

	latestDistributionHeight := s.latestDistributionWithdrawalHeight
	if err != nil {
		return 0, 0, false, err
	}
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

func (s *Service) calMaxClaimableWithdrawalIndex(targetEth1BlockHeight uint64, totalUserEthDeci decimal.Decimal) (uint64, error) {
	calOpts := s.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))
	maxClaimableWithdrawIndex, err := s.networkWithdrawalContract.MaxClaimableWithdrawalIndex(calOpts)
	if err != nil {
		return 0, err
	}
	// nextWithdrawIndex <= real value
	nextWithdrawIndex, err := s.networkWithdrawalContract.NextWithdrawalIndex(calOpts)
	if err != nil {
		return 0, err
	}
	totalWithdrawalShortages, err := s.networkWithdrawalContract.TotalWithdrawalShortages(calOpts)
	if err != nil {
		return 0, err
	}
	newMaxClaimableWithdrawalIndex := uint64(0)
	totalWithdrawalShortagesDeci := decimal.NewFromBigInt(totalWithdrawalShortages, 0)
	if totalWithdrawalShortagesDeci.LessThanOrEqual(totalUserEthDeci) {
		if nextWithdrawIndex.Uint64() >= 1 {
			newMaxClaimableWithdrawalIndex = nextWithdrawIndex.Uint64() - 1
		}
	} else {
		willMissingAmountDeci := totalWithdrawalShortagesDeci.Sub(totalUserEthDeci)
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

				latestUsersWaitAmountDeci = latestUsersWaitAmountDeci.Add(stakerWithdrawal.EthAmount)
				if latestUsersWaitAmountDeci.GreaterThan(willMissingAmountDeci) {
					if i >= 1 {
						newMaxClaimableWithdrawalIndex = i - 1
					}
					break
				}
			}
		}
	}
	if newMaxClaimableWithdrawalIndex < maxClaimableWithdrawIndex.Uint64() {
		newMaxClaimableWithdrawalIndex = maxClaimableWithdrawIndex.Uint64()
	}

	return newMaxClaimableWithdrawalIndex, nil
}

func (s *Service) sendDistributeTx(distributeType uint8, targetEth1BlockHeight, totalUserEth, totalNodeEth, totalPlatformEth, newMaxClaimableWithdrawalIndex *big.Int) error {
	err := s.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer s.connection.UnlockTxOpts()

	encodeBts, err := s.networkWithdrdawalAbi.Pack("distribute", distributeType, targetEth1BlockHeight,
		totalUserEth, totalNodeEth, totalPlatformEth, newMaxClaimableWithdrawalIndex)
	if err != nil {
		return err
	}

	proposalId := utils.ProposalId(s.networkWithdrawalAddress, encodeBts, targetEth1BlockHeight)

	// check voted
	hasVoted, err := s.networkProposalContract.HasVoted(nil, proposalId, s.keyPair.CommonAddress())
	if err != nil {
		return fmt.Errorf("networkProposalContract.HasVoted err: %s", err)
	}
	if hasVoted {
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"distributeType":                 distributeType,
		"targetEth1BlockHeight":          targetEth1BlockHeight,
		"totalUserEthDeci":               totalUserEth.String(),
		"totalNodeEthDeci":               totalNodeEth.String(),
		"totalPlatformEthDeci":           totalPlatformEth.String(),
		"newMaxClaimableWithdrawalIndex": newMaxClaimableWithdrawalIndex,
	}).Info("Will sendDistributeTx")

	tx, err := s.networkProposalContract.ExecProposal(s.connection.TxOpts(), s.networkWithdrawalAddress, encodeBts, targetEth1BlockHeight)
	if err != nil {
		return err
	}

	logrus.Infof("send Distribute tx hash: %s", tx.Hash().String())

	return s.waitProposalTxOk(tx.Hash(), proposalId)
}
