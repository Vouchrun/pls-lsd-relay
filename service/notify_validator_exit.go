package service

import (
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

var farfutureBlockHeight = uint64(1e11)

type ExitElection struct {
	ValidatorIndex    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex"`
	NotifyBlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_block_number"`
	NotifyTimestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_timestamp"`
	WithdrawCycle     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:withdraw_cycle"`

	ExitEpoch     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_epoch"`
	ExitTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_timestamp"`
}

func (s *Service) notifyValidatorExit() error {
	currentCycle, targetTimestamp := currentCycleAndStartTimestamp()
	preCycle := currentCycle - 1

	targetEpoch := utils.EpochAtTimestamp(s.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := s.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return err
	}
	targetCall := s.connection.CallOpts(big.NewInt(int64(targetBlockNumber)))

	ejectedValidator, err := s.networkWithdrawContract.GetEjectedValidatorsAtCycle(nil, big.NewInt(preCycle))
	if err != nil {
		return err
	}
	// return if already dealed
	if len(ejectedValidator) != 0 {
		logrus.Debugf("ejectedValidator %d at precycle %d", len(ejectedValidator), preCycle)
		return nil
	}

	totalMissingAmount, err := s.networkWithdrawContract.TotalMissingAmountForWithdraw(targetCall)
	if err != nil {
		return err
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

	// todo calc exited but not distributed amount
	exitButNotDistributedValidatorList := make([]*Validator, 0)
	if err != nil {
		return errors.Wrap(err, "GetValidatorListWithdrawableEpochAfter failed")
	}
	totalExitedButNotDistributedUserAmount := decimal.Zero
	notDistributeValidators := make(map[uint64]bool)
	for _, v := range exitButNotDistributedValidatorList {
		notDistributeValidators[v.ValidatorIndex] = true
		totalExitedButNotDistributedUserAmount = totalExitedButNotDistributedUserAmount.Add(utils.StandardEffectiveBalance.Sub(v.NodeDepositAmount))
	}

	// calc partial withdrawal not distributed amount
	latestDistributeWithdrawalHeight, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(targetCall)
	if err != nil {
		return err
	}
	// should exclude notDistributeValidators, as we has already calc
	userUndistributedWithdrawalsDeci, _, _, _, err := s.getUserNodePlatformFromWithdrawals(latestDistributeWithdrawalHeight.Uint64(), farfutureBlockHeight)
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

	// todo cal start cycle
	notExitElectionList := make([]*ExitElection, 0)
	if err != nil {
		return errors.Wrap(err, "GetAllNotExitElectionList failed")
	}
	startCycle := preCycle - 1
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

// utc 8:00
func currentCycleAndStartTimestamp() (int64, int64) {
	currentCycle := (time.Now().Unix()) / 86400
	targetTimestamp := currentCycle * 86400
	return currentCycle, targetTimestamp
}

func (s *Service) mustSelectValidatorsForExit(totalMissingAmount decimal.Decimal, targetEpoch uint64) ([]*big.Int, error) {
	return nil, nil
}
