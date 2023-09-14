package service

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) submitBalances() error {
	beaconHead, err := s.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	targetEpoch := (beaconHead.FinalizedEpoch / s.submitBalancesDuEpochs) * s.submitBalancesDuEpochs

	balancesBlockOnChain, err := s.networkBalancesContract.BalancesBlock(nil)
	if err != nil {
		return fmt.Errorf("networkBalancesContract.BalancesBlock err: %s", err)
	}

	targetBlock, err := s.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"targetEpoch":          targetEpoch,
		"targetBlock":          targetBlock,
		"balancesBlockOnChain": balancesBlockOnChain.String(),
	}).Debug("epocheInfo")

	// already update on this block, no need vote
	if targetBlock <= balancesBlockOnChain.Uint64() {
		return nil
	}

	targetCallOpts := s.connection.CallOpts(big.NewInt(int64(targetBlock)))

	lsdTokenTotalSupply, err := s.lsdTokenContract.TotalSupply(targetCallOpts)
	if err != nil {
		return err
	}
	lsdTokenTotalSupplyDeci := decimal.NewFromBigInt(lsdTokenTotalSupply, 0)

	// deposit pool balance
	userDepositPoolBalance, err := s.userDepositContract.GetBalance(targetCallOpts)
	if err != nil {
		return err
	}
	userDepositPoolBalanceDeci := decimal.NewFromBigInt(userDepositPoolBalance, 0)

	targetValidators := s.GetValidatorDepositedListBefore(targetBlock)
	logrus.WithFields(logrus.Fields{
		"validatorDepositedList len": len(targetValidators),
	}).Debug("validatorDepositedList")

	// user eth from validators
	totalUserEthFromValidatorDeci := decimal.Zero
	for _, validator := range targetValidators {
		userAllEth, err := s.getUserEthInfoFromValidatorBalance(validator, targetEpoch)
		if err != nil {
			return err
		}
		totalUserEthFromValidatorDeci = totalUserEthFromValidatorDeci.Add(userAllEth)
	}

	// total missing amount for withdraw
	totalMissingAmount, err := s.networkWithdrawContract.TotalMissingAmountForWithdraw(targetCallOpts)
	if err != nil {
		return err
	}
	totalMissingAmountDeci := decimal.NewFromBigInt(totalMissingAmount, 0)

	// user eth from undistributed withdrawals
	latestDistributeWithdrawalsHeight, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(targetCallOpts)
	if err != nil {
		return err
	}
	if latestDistributeWithdrawalsHeight.Cmp(big.NewInt(0)) == 0 {
		latestDistributeWithdrawalsHeight = big.NewInt(int64(s.networkCreateBlock))
	}
	userEthFromWithdrawDeci, _, _, _, err := s.getUserNodePlatformFromWithdrawals(latestDistributeWithdrawalsHeight.Uint64(), targetBlock)
	if err != nil {
		return err
	}

	// user eth from undistributed priorityfee
	latestDistributePriorityFeeHeight, err := s.networkWithdrawContract.LatestDistributePriorityFeeHeight(targetCallOpts)
	if err != nil {
		return err
	}
	if latestDistributePriorityFeeHeight.Cmp(big.NewInt(0)) == 0 {
		latestDistributePriorityFeeHeight = big.NewInt(int64(s.networkCreateBlock))
	}
	userEthFromPriorityFeeDeci, _, _, _, err := s.getUserNodePlatformFromPriorityFee(latestDistributePriorityFeeHeight.Uint64(), targetBlock)
	if err != nil {
		return err
	}

	// ----final: total user eth = total user eth from validator + deposit pool balance + user undistributedWithdrawals +
	// 								+ user undistributed priority fee  - totalMissingAmountForWithdraw
	totalUserEthDeci := totalUserEthFromValidatorDeci.Add(userDepositPoolBalanceDeci).Add(userEthFromWithdrawDeci).
		Add(userEthFromPriorityFeeDeci).Sub(totalMissingAmountDeci)

	// check exchange rate
	oldExchangeRate, err := s.networkBalancesContract.GetExchangeRate(targetCallOpts)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}
	oldExchangeRateDeci := decimal.NewFromBigInt(oldExchangeRate, 0)
	newExchangeRateDeci := totalUserEthDeci.Mul(decimal.NewFromInt(1e18)).Div(lsdTokenTotalSupplyDeci)
	if newExchangeRateDeci.LessThanOrEqual(oldExchangeRateDeci) {
		logrus.WithFields(logrus.Fields{
			"newExchangeRate": newExchangeRateDeci.StringFixed(0),
			"oldExchangeRate": oldExchangeRate.String(),
		}).Warn("new exchangeRate less than old")
		return nil
	}
	var maxRateChangeDeci = decimal.NewFromInt(11e14) //0.0011
	if newExchangeRateDeci.GreaterThan(oldExchangeRateDeci.Add(maxRateChangeDeci)) {
		return fmt.Errorf("newExchangeRate %s too big than oldExchangeRate %s", newExchangeRateDeci.String(), oldExchangeRateDeci.String())
	}

	hasVoted, err := s.networkProposalContract.HasVoted(nil, utils.SubmitBalancesProposalId(big.NewInt(int64(targetBlock)),
		totalUserEthDeci.BigInt(), lsdTokenTotalSupply), s.keyPair.CommonAddress())
	if err != nil {
		return err
	}
	if hasVoted {
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"targetBlockNumber":                 targetBlock,
		"targetEpoch":                       targetEpoch,
		"totalUserEthFromValidator":         totalUserEthFromValidatorDeci.StringFixed(0),
		"userDepositPoolBalanceDeci":        userDepositPoolBalanceDeci.StringFixed(0),
		"userUndistributedWithdrawalsDeci":  userEthFromWithdrawDeci.StringFixed(0),
		"userUndistributedPriorityFeeDeci":  userEthFromPriorityFeeDeci.StringFixed(0),
		"totalMissingAmountForWithdrawDeci": totalMissingAmountDeci.StringFixed(0),
		"totalUserEth":                      totalUserEthDeci.StringFixed(0),
		"lsdTokenTotalSupply":               lsdTokenTotalSupplyDeci.StringFixed(0),
		"newExchangeRate":                   newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate":                   oldExchangeRateDeci.StringFixed(0),
	}).Info("exchangeRateInfo")

	return s.sendSubmitBalancesTx(big.NewInt(int64(targetBlock)), totalUserEthDeci.BigInt(), lsdTokenTotalSupply)

}

func (task *Service) getUserEthInfoFromValidatorBalance(validator *Validator, targetEpoch uint64) (decimal.Decimal, error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch:
		switch validator.NodeType {
		case utils.NodeTypeLight:
			return decimal.Zero, nil
		case utils.NodeTypeTrust:
			return utils.StandardTrustNodeFakeDepositBalance, nil
		default:
			// common node and trust node should not happen here
			return decimal.Zero, fmt.Errorf("unknow node type: %d", validator.NodeType)
		}

	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		userDepositBalance := utils.StandardEffectiveBalanceDeci.Sub(validator.NodeDepositAmountDeci)
		return userDepositBalance, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
		utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:

		userDepositBalance := utils.StandardEffectiveBalanceDeci.Sub(validator.NodeDepositAmountDeci)
		// case: activeEpoch 155747 > targetEpoch 155700
		if validator.ActiveEpoch > targetEpoch {
			return userDepositBalance, nil
		}

		validatorStatus, err := task.connection.GetValidatorStatus(types.BytesToValidatorPubkey(validator.Pubkey), &beacon.ValidatorStatusOptions{
			Epoch: &targetEpoch,
		})
		if err != nil {
			return decimal.Zero, err
		}

		userDepositPlusReward, err := task.getUserDepositPlusReward(validator.NodeDepositAmountDeci, decimal.NewFromInt(int64(validatorStatus.Balance)).Mul(utils.GweiDeci))
		if err != nil {
			return decimal.Zero, errors.Wrap(err, "getUserDepositPlusReward failed")
		}
		return userDepositPlusReward, nil

	case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:
		return decimal.Zero, nil

	default:
		return decimal.Zero, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}

func (s *Service) getUserDepositPlusReward(nodeDepositAmount, validatorBalance decimal.Decimal) (decimal.Decimal, error) {
	userDepositAmount := utils.StandardEffectiveBalanceDeci.Sub(nodeDepositAmount)

	switch {
	case validatorBalance.IsZero(): //withdrawdone case
		return decimal.Zero, nil
	case validatorBalance.GreaterThan(decimal.Zero) && validatorBalance.LessThan(utils.StandardEffectiveBalanceDeci):
		loss := utils.StandardEffectiveBalanceDeci.Sub(validatorBalance)
		if loss.LessThan(nodeDepositAmount) {
			return userDepositAmount, nil
		} else {
			return validatorBalance, nil
		}
	case validatorBalance.Equal(utils.StandardEffectiveBalanceDeci):
		return userDepositAmount, nil
	case validatorBalance.GreaterThan(utils.StandardEffectiveBalanceDeci):
		// total staking reward
		validatorTotalStakingReward := validatorBalance.Sub(utils.StandardEffectiveBalanceDeci)

		userRewardOfThisValidator, _, _ := utils.GetUserNodePlatformReward(nodeDepositAmount, validatorTotalStakingReward)

		return userDepositAmount.Add(userRewardOfThisValidator), nil
	default:
		// should not happen here
		return decimal.Zero, fmt.Errorf("unknown balance")
	}
}

func (s *Service) sendSubmitBalancesTx(block, totalUserEth, lsdTokenTotalSupply *big.Int) error {
	err := s.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer s.connection.UnlockTxOpts()

	tx, err := s.networkBalancesContract.SubmitBalances(
		s.connection.TxOpts(),
		block,
		totalUserEth,
		lsdTokenTotalSupply)
	if err != nil {
		return err
	}

	logrus.Info("send submitBalances tx hash: ", tx.Hash().String())

	return s.waitTxOk(tx.Hash())
}
