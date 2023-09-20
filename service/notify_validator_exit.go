package service

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) notifyValidatorExit() error {
	currentCycle, targetTimestamp, err := s.currentCycleAndStartTimestamp()
	if err != nil {
		return fmt.Errorf("currentCycleAndStartTimestamp failed: %w", err)
	}
	preCycle := currentCycle - 1

	targetEpoch := utils.EpochAtTimestamp(s.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := s.getEpochStartBlocknumberWithCheck(targetEpoch)
	if err != nil {
		return fmt.Errorf("getEpochStartBlocknumberWithCheck failed: %w", err)
	}

	// wait validator updated
	if targetEpoch > s.latestEpochOfUpdateValidator {
		logrus.Debugf("targetEpoch: %d  latestEpochOfUpdateValidator: %d", targetEpoch, s.latestEpochOfUpdateValidator)
		return nil
	}

	// wait sync block
	if targetBlockNumber > s.latestBlockOfSyncBlock {
		logrus.Debugf("targetBlockNumber: %d  latestBlockOfSyncBlock: %d", targetBlockNumber, s.latestBlockOfSyncBlock)
		return nil
	}

	targetCall := s.connection.CallOpts(big.NewInt(int64(targetBlockNumber)))

	ejectedValidator, err := s.networkWithdrawContract.GetEjectedValidatorsAtCycle(nil, big.NewInt(preCycle))
	if err != nil {
		return fmt.Errorf("GetEjectedValidatorsAtCycle failed: %w", err)
	}
	// return if already dealed
	if len(ejectedValidator) != 0 {
		logrus.Debugf("ejectedValidator %d at precycle %d", len(ejectedValidator), preCycle)
		return nil
	}

	totalMissingAmount, err := s.networkWithdrawContract.TotalMissingAmountForWithdraw(targetCall)
	if err != nil {
		return fmt.Errorf("TotalMissingAmountForWithdraw failed: %w, target block: %d", err, targetBlockNumber)
	}
	totalMissingAmountDeci := decimal.NewFromBigInt(totalMissingAmount, 0)

	// no need notify exit
	if totalMissingAmount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	userDepositBalance, err := s.userDepositContract.GetBalance(targetCall)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"currentCycle:":      currentCycle,
		"targetEpoch":        targetEpoch,
		"targetBlockNumber":  targetBlockNumber,
		"totalMissingAmount": totalMissingAmount,
		"userDepositBalance": userDepositBalance,
	}).Debug("notifyValidatorExit")

	// calc exited but not distributed amount
	exitButNotDistributedValidatorList := s.exitButNotDistributedValidatorList(targetEpoch)
	if err != nil {
		return errors.Wrap(err, "GetValidatorListWithdrawableEpochAfter failed")
	}
	totalExitedButNotDistributedUserAmount := decimal.Zero
	notDistributeValidators := make(map[uint64]bool)
	for _, v := range exitButNotDistributedValidatorList {
		notDistributeValidators[v.ValidatorIndex] = true
		totalExitedButNotDistributedUserAmount = totalExitedButNotDistributedUserAmount.Add(utils.StandardEffectiveBalanceDeci.Sub(v.NodeDepositAmountDeci))
	}

	// calc partial withdrawal not distributed amount
	latestDistributeWithdrawalHeight, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(targetCall)
	if err != nil {
		return err
	}
	// should exclude notDistributeValidators, as we has already calc
	userUndistributedWithdrawalsDeci, _, _, _, err := s.getUserNodePlatformFromWithdrawals(latestDistributeWithdrawalHeight.Uint64(), targetBlockNumber)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromWithdrawals failed")
	}

	totalPendingAmountDeci := totalExitedButNotDistributedUserAmount.Add(userUndistributedWithdrawalsDeci)
	// no need notify exit
	if totalMissingAmountDeci.LessThanOrEqual(totalPendingAmountDeci) {
		return nil
	}

	// final total missing amount
	finalTotalMissingAmountDeci := totalMissingAmountDeci.Sub(totalPendingAmountDeci)

	selectVals, err := s.mustSelectValidatorsForExit(finalTotalMissingAmountDeci, targetEpoch)
	if err != nil {
		return errors.Wrap(err, "selectValidatorsForExit failed")
	}
	if len(selectVals) == 0 {
		return fmt.Errorf("selectValidatorsForExit select zero vals, target epoch: %d", targetEpoch)
	}

	// cal start cycle
	startCycle := preCycle - 1
	notExitElectionList := s.notExitElectionList(uint64(currentCycle))
	if err != nil {
		return errors.Wrap(err, "GetAllNotExitElectionList failed")
	}
	if len(notExitElectionList) > 0 {
		startCycle = int64(notExitElectionList[0].WithdrawCycle)
	}
	logrus.WithFields(logrus.Fields{
		"startCycle": startCycle,
		"preCycle":   preCycle,
		"selectVal":  selectVals,
	}).Debug("will sendNotifyValidatorExitTx")

	// ---- send NotifyValidatorExit tx
	return s.sendNotifyExitTx(uint64(preCycle), uint64(startCycle), selectVals)
}

func (s *Service) sendNotifyExitTx(preCycle, startCycle uint64, selectVal []*big.Int) error {
	err := s.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer s.connection.UnlockTxOpts()
	tx, err := s.networkWithdrawContract.NotifyValidatorExit(s.connection.TxOpts(), big.NewInt(int64(preCycle)),
		big.NewInt(int64(startCycle)), selectVal)
	if err != nil {
		return err
	}

	logrus.Info("send NotifyValidatorExit tx hash: ", tx.Hash().String())

	return s.waitTxOk(tx.Hash())
}

func (s *Service) currentCycleAndStartTimestamp() (int64, int64, error) {
	currentCycle, err := s.networkWithdrawContract.CurrentWithdrawCycle(nil)
	if err != nil {
		return 0, 0, err
	}
	cycleSeconds, err := s.networkWithdrawContract.WithdrawCycleSeconds(nil)
	if err != nil {
		return 0, 0, err
	}

	targetTimestamp := currentCycle.Uint64() * cycleSeconds.Uint64()
	return int64(currentCycle.Uint64()), int64(targetTimestamp), nil
}

func (s *Service) mustSelectValidatorsForExit(totalMissingAmount decimal.Decimal, targetEpoch uint64) ([]*big.Int, error) {
	vals, err := s.getValidatorsOfTargetEpoch(targetEpoch)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].ActiveEpoch < vals[j].ActiveEpoch
	})

	return nil, nil
}
