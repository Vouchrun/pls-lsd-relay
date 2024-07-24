package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) submitBalances() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	beaconHead, err := s.connection.BeaconHead()
	if err != nil {
		return err
	}
	targetEpoch := (beaconHead.FinalizedEpoch / s.submitBalancesDuEpochs) * s.submitBalancesDuEpochs

	snapshotOnchain, err := s.networkBalancesContract.BalancesSnapshot(nil)
	if err != nil {
		return fmt.Errorf("networkBalancesContract.BalancesBlock err: %s", err)
	}

	targetBlock, err := s.getEpochStartBlocknumberWithCheck(targetEpoch)
	if err != nil {
		return err
	}

	// already update on this block, no need vote
	if targetBlock <= snapshotOnchain.Block.Uint64() {
		return nil
	}

	// wait sync block
	if targetBlock > s.latestBlockOfSyncBlock {
		return nil
	}

	s.log.WithFields(logrus.Fields{
		"targetEpoch":          targetEpoch,
		"targetBlock":          targetBlock,
		"balancesBlockOnChain": snapshotOnchain.Block.Uint64(),
	}).Debug("epochInfo")

	targetCallOpts := s.connection.CallOpts(big.NewInt(int64(targetBlock)))

	lsdTokenTotalSupply, err := s.lsdTokenContract.TotalSupply(targetCallOpts)
	if err != nil {
		return err
	}
	lsdTokenTotalSupplyDeci := decimal.NewFromBigInt(lsdTokenTotalSupply, 0)
	if lsdTokenTotalSupplyDeci.IsZero() {
		return nil
	}

	// deposit pool balance
	userDepositPoolBalance, err := s.userDepositContract.GetBalance(targetCallOpts)
	if err != nil {
		return err
	}
	userDepositPoolBalanceDeci := decimal.NewFromBigInt(userDepositPoolBalance, 0)

	targetValidators := s.GetValidatorDepositedListBeforeBlock(targetBlock)
	s.log.WithFields(logrus.Fields{
		"validatorDepositedList len": len(targetValidators),
	}).Debug("validatorDepositedList")

	// user eth from validators
	totalUserEthFromValidatorDeci := decimal.Zero
	for _, validator := range targetValidators {
		userAllEth, err := s.getUserEthInfoFromValidatorBalance(ctx, validator, targetEpoch)
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
		latestDistributeWithdrawalsHeight = big.NewInt(int64(s.startAtBlock))
	}
	userEthFromWithdrawDeci, _, _, _, err := s.getUserNodePlatformFromWithdrawals(latestDistributeWithdrawalsHeight.Uint64(), targetBlock)
	if err != nil {
		return err
	}

	// user eth from undistributed priority fee
	latestDistributePriorityFeeHeight, err := s.networkWithdrawContract.LatestDistributePriorityFeeHeight(targetCallOpts)
	if err != nil {
		return err
	}
	if latestDistributePriorityFeeHeight.Cmp(big.NewInt(0)) == 0 {
		latestDistributePriorityFeeHeight = big.NewInt(int64(s.startAtBlock))
	}
	userEthFromPriorityFeeDeci, _, _, _, err := s.getUserNodePlatformFromPriorityFee(latestDistributePriorityFeeHeight.Uint64(), targetBlock)
	if err != nil {
		return err
	}

	// ----final: total user eth = total user eth from validator + deposit pool balance + user undistributedWithdrawals +
	// 								+ user undistributed priority fee  - totalMissingAmountForWithdraw
	totalUserEthDeci := totalUserEthFromValidatorDeci.Add(userDepositPoolBalanceDeci).Add(userEthFromWithdrawDeci).
		Add(userEthFromPriorityFeeDeci).Sub(totalMissingAmountDeci)
	if totalUserEthDeci.BigInt().Cmp(lsdTokenTotalSupply) < 0 {
		s.log.WithFields(logrus.Fields{
			"old_totalUserEthDeci": totalUserEthDeci.StringFixed(0),
		}).Warn("adjust totalUserEthDeci to lsdTokenTotalSupply")
		totalUserEthDeci = decimal.NewFromBigInt(lsdTokenTotalSupply, 0)
	}

	// check exchange rate
	oldExchangeRate, err := s.networkBalancesContract.GetExchangeRate(targetCallOpts)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}
	oldExchangeRateDeci := decimal.NewFromBigInt(oldExchangeRate, 0)

	newExchangeRateDeci := totalUserEthDeci.Mul(decimal.NewFromInt(1e18)).Div(lsdTokenTotalSupplyDeci)
	rateChangeLimit, err := s.networkBalancesContract.RateChangeLimit(nil)
	if err != nil {
		return err
	}

	one18 := decimal.NewFromBigInt(big.NewInt(1), 18)
	rateChange := newExchangeRateDeci.Sub(oldExchangeRateDeci).Abs().Mul(one18).Div(oldExchangeRateDeci)
	rateInfoLog := s.log.WithFields(logrus.Fields{
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
		"rateChange":                        rateChange.StringFixed(0),
	})
	if rateChange.GreaterThan(decimal.NewFromBigInt(rateChangeLimit, 0)) {
		rateInfoLog.Error("exchangeRateInfo")
		return fmt.Errorf("exceed rate change limit %s, newExchangeRate %s, oldExchangeRate %s",
			rateChangeLimit.String(), newExchangeRateDeci.String(), oldExchangeRateDeci.String())
	}
	rateInfoLog.Info("exchangeRateInfo")

	return s.sendSubmitBalancesTx(big.NewInt(int64(targetBlock)), totalUserEthDeci.BigInt(), lsdTokenTotalSupply)

}

func (task *Service) getUserEthInfoFromValidatorBalance(ctx context.Context, validator *Validator, targetEpoch uint64) (decimal.Decimal, error) {
	validatorStatus, err := task.connection.GetValidatorStatus(ctx, types.BytesToValidatorPubkey(validator.Pubkey), &beacon.ValidatorStatusOptions{
		Epoch: &targetEpoch,
	})
	if err != nil {
		return decimal.Zero, err
	}

	status, err := mapValidatorStatus(&validatorStatus)
	if err != nil {
		return decimal.Zero, fmt.Errorf("unknown validator status: %d", status)
	}

	switch status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch:
		switch validator.NodeType {
		case utils.NodeTypeSolo:
			return decimal.Zero, nil
		case utils.NodeTypeTrust:
			return decimal.NewFromBigInt(big.NewInt(int64(task.manager.cfg.TrustNodeDepositAmount)), 18), nil
		default:
			// common node and trust node should not happen here
			return decimal.Zero, fmt.Errorf("unknown node type: %d", validator.NodeType)
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

		userDepositPlusReward, err := task.getUserDepositPlusReward(validator.NodeDepositAmountDeci, decimal.NewFromInt(int64(validatorStatus.Balance)).Mul(utils.GweiDeci))
		if err != nil {
			return decimal.Zero, errors.Wrap(err, "getUserDepositPlusReward failed")
		}
		return userDepositPlusReward, nil

	case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:
		return decimal.Zero, nil

	default:
		return decimal.Zero, fmt.Errorf("unknown validator status: %d", status)
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

		userRewardOfThisValidator, _, _ := utils.GetUserNodePlatformReward(s.nodeCommissionRate, s.platformCommissionRate, nodeDepositAmount, validatorTotalStakingReward)

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

	encodeBts, err := s.networkBalancesAbi.Pack("submitBalances", block, totalUserEth, lsdTokenTotalSupply)
	if err != nil {
		return err
	}

	proposalId := utils.ProposalId(s.networkBalancesAddress, encodeBts, block)

	// check voted
	hasVoted, err := s.networkProposalContract.HasVoted(nil, proposalId, s.connection.Keypair().CommonAddress())
	if err != nil {
		return fmt.Errorf("networkProposalContract.HasVoted err: %s", err)
	}
	if hasVoted {
		s.log.Info("already voted wait other voters")
		return nil
	}

	tx, err := s.networkProposalContract.ExecProposal(s.connection.TxOpts(), s.networkBalancesAddress, encodeBts, block)
	if err != nil {
		return err
	}

	s.log.Info("send submitBalances tx hash: ", tx.Hash().String())

	return s.waitProposalTxOk(tx.Hash(), proposalId)
}
