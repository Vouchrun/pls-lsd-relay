package service

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

func (s *Service) waitProposalTxOk(txHash common.Hash, proposalId [32]byte) error {
	_, err := s.connection.Eth1Client().WaitTxOkCommon(txHash)
	if err != nil {
		p, err := s.networkProposalContract.Proposals(nil, proposalId)
		if err != nil {
			return err
		}
		// proposal is executed
		if p.Status == 2 {
			return nil
		}
		return err
	}

	return nil
}

func (s *Service) waitProposalsTxOk(txHash common.Hash, proposalIds [][32]byte) error {
	_, err := s.connection.Eth1Client().WaitTxOkCommon(txHash)
	if err != nil {
		allProposalsExecuted := true
		for _, proposalId := range proposalIds {
			p, err := s.networkProposalContract.Proposals(nil, proposalId)
			if err != nil {
				return err
			}
			// proposal is executed
			if p.Status != 2 {
				allProposalsExecuted = false
				break
			}
		}
		if allProposalsExecuted {
			return nil
		}

		return err
	}

	return nil
}

func (s *Service) getEpochStartBlocknumberWithCheck(epoch uint64) (uint64, error) {
	s.cacheEpochToBlockIDMutex.Lock()
	defer s.cacheEpochToBlockIDMutex.Unlock()

	if blockID, ok := s.cacheEpochToBlockID.Get(epoch); ok {
		return blockID, nil
	}

	targetBlock, err := s.getEpochStartBlocknumber(epoch)
	if err != nil {
		return 0, err
	}

	if targetBlock < s.startAtBlock {
		targetBlock = s.startAtBlock + 1
	}
	s.cacheEpochToBlockID.Add(epoch, targetBlock)
	return targetBlock, nil
}

func (s *Service) getEpochStartBlocknumber(epoch uint64) (uint64, error) {
	eth2ValidatorBalanceSyncerStartSlot := utils.StartSlotOfEpoch(s.eth2Config, epoch)
	retry := 0
	for {
		if retry > 10 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}

		targetBeaconBlock, exist, err := s.connection.GetBeaconBlock(eth2ValidatorBalanceSyncerStartSlot)
		if err != nil {
			return 0, fmt.Errorf("fail to get beacon block[%d]: %w", eth2ValidatorBalanceSyncerStartSlot, err)
		}
		// we will use next slot if not exist
		if !exist {
			eth2ValidatorBalanceSyncerStartSlot++
			retry++
			continue
		}
		if targetBeaconBlock.ExecutionBlockNumber == 0 {
			return 0, fmt.Errorf("beacon slot %d executionBlockNumber is zero", eth2ValidatorBalanceSyncerStartSlot)
		}
		return targetBeaconBlock.ExecutionBlockNumber, nil
	}
}

// return (user reward, node reward, platform fee, nodeRewardMap) decimals 18
func (s *Service) getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight uint64) (decimal.Decimal, decimal.Decimal, decimal.Decimal, NodeNewRewardsMap, error) {
	totalUserEthDeci := decimal.Zero
	totalNodeEthDeci := decimal.Zero
	totalPlatformEthDeci := decimal.Zero
	nodeNewRewardsMap := make(NodeNewRewardsMap)

	for i := latestDistributeHeight + 1; i <= targetEth1BlockHeight; i++ {
		block, err := s.getBeaconBlock(i)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
		}

		for _, w := range block.Withdrawals {
			val, exist := s.getValidatorByIndex(w.ValidatorIndex)
			if !exist {
				continue
			}

			totalReward := uint64(0)
			userDeposit := uint64(0)
			nodeDeposit := uint64(0)

			switch {

			case w.Amount < utils.MaxPartialWithdrawalAmount: // partial withdrawal
				totalReward = w.Amount

			case w.Amount >= utils.MaxPartialWithdrawalAmount && w.Amount < utils.StandardEffectiveBalance: // slash
				totalReward = 0

				userDeposit = utils.StandardEffectiveBalance - val.NodeDepositAmount
				if userDeposit > w.Amount {
					userDeposit = w.Amount
					nodeDeposit = 0
				} else {
					nodeDeposit = w.Amount - userDeposit
				}

			case w.Amount >= utils.StandardEffectiveBalance: // full withdrawal
				totalReward = w.Amount - utils.StandardEffectiveBalance

				userDeposit = utils.StandardEffectiveBalance - val.NodeDepositAmount
				nodeDeposit = val.NodeDepositAmount

			default:
				return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("unknown withdrawal's amount %d, valIndex: %d", w.Amount, val.ValidatorIndex)
			}

			// distribute reward
			userRewardDeci, nodeRewardDeci, platformFeeDeci := utils.GetUserNodePlatformReward(s.nodeCommissionRate, s.platformCommissionRate, val.NodeDepositAmountDeci, decimal.NewFromInt(int64(totalReward)).Mul(utils.GweiDeci))
			userDepositDeci := decimal.NewFromInt(int64(userDeposit)).Mul(utils.GweiDeci)
			nodeDepositDeci := decimal.NewFromInt(int64(nodeDeposit)).Mul(utils.GweiDeci)

			// cal node reward
			nodeNewReward, exist := nodeNewRewardsMap[val.NodeAddress]
			if exist {
				nodeNewReward.TotalRewardAmount = nodeNewReward.TotalRewardAmount.Add(nodeRewardDeci)
				nodeNewReward.TotalExitDepositAmount = nodeNewReward.TotalExitDepositAmount.Add(nodeDepositDeci)
			} else {
				n := NodeNewReward{
					Address:                val.NodeAddress.String(),
					TotalRewardAmount:      nodeRewardDeci,
					TotalExitDepositAmount: nodeDepositDeci,
				}
				nodeNewRewardsMap[val.NodeAddress] = &n
			}

			// cal total vals
			totalUserEthDeci = totalUserEthDeci.Add(userRewardDeci).Add(userDepositDeci)
			totalNodeEthDeci = totalNodeEthDeci.Add(nodeRewardDeci).Add(nodeDepositDeci)
			totalPlatformEthDeci = totalPlatformEthDeci.Add(platformFeeDeci)
		}
	}

	return totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, nodeNewRewardsMap, nil
}

func WalkTrace(seekFn func(tx *connection.TxTrace) bool, amount decimal.Decimal, trace connection.TxTrace) decimal.Decimal {
	value := decimal.NewFromBigInt(trace.Value.ToInt(), 0)

	if seekFn(&trace) {
		amount = amount.Add(value)
	}

	if trace.Calls != nil {
		for _, call := range trace.Calls {
			amount = WalkTrace(seekFn, amount, call)
		}
	}

	return amount
}

func (s *Service) cacheFeePoolBalances(log *logrus.Entry, fromBlock, toBlock uint64) error {
	log.Debug("start cache fee pool balances")
	defer func() {
		log.Debug("end cache fee pool balances")
	}()

	for i := fromBlock; i <= toBlock; i += s.batchQueryBalanceBlockNumbers {
		end := i + s.batchQueryBalanceBlockNumbers - 1
		if end > toBlock {
			end = toBlock
		}
		blocks := make([]uint64, 0, s.batchQueryBalanceBlockNumbers)
		for j := i; j <= end; j++ {
			if _, ok := s.feePoolBalances.Load(uint64(j)); !ok {
				blocks = append(blocks, uint64(j))
			}
		}
		if len(blocks) > 0 {
			feePoolBalances, err := s.connection.Eth1Client().(*connection.Eth1Client).BatchBalancesAtBlocks(context.Background(), s.feePoolAddress, blocks)
			if err != nil {
				return fmt.Errorf("fail to batch query fee pool balances: %w", err)
			}
			for block, balance := range feePoolBalances {
				s.feePoolBalances.Store(block, balance)
			}
		}
	}
	return nil
}

func (s *Service) getFeePoolBalance(blockNumber uint64) (*big.Int, error) {
	balance, ok := s.feePoolBalances.Load(blockNumber)
	if !ok {
		return nil, fmt.Errorf("fee pool balance not found for block %d", blockNumber)
	}
	return balance.(*big.Int), nil
}

// return (user reward, node reward, platform fee) decimals 18
func (s *Service) getUserNodePlatformFromPriorityFee(log *logrus.Entry, latestDistributeHeight, targetEth1BlockHeight uint64) (decimal.Decimal, decimal.Decimal, decimal.Decimal, NodeNewRewardsMap, error) {
	ctx := context.Background()
	totalUserEthDeci := decimal.Zero
	totalNodeEthDeci := decimal.Zero
	totalPlatformEthDeci := decimal.Zero
	nodeNewRewardsMap := make(NodeNewRewardsMap)

	log = log.WithFields(logrus.Fields{
		"fromBlock":                          latestDistributeHeight + 1,
		"targetBlock":                        targetEth1BlockHeight,
		"getUserNodePlatformFromPriorityFee": true,
	})

	log.Debug("start getUserNodePlatformFromPriorityFee")
	defer func() {
		log.Debug("end getUserNodePlatformFromPriorityFee")
	}()

	if err := s.cacheFeePoolBalances(log, latestDistributeHeight, targetEth1BlockHeight); err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
	}

	log.Debug("start filter all withdrawn events")
	// filter all withdrawn events
	withdrawals := make(map[uint64]*big.Int)
	for i := latestDistributeHeight + 1; i <= targetEth1BlockHeight; i += s.eventFilterMaxSpanBlocks {
		end := i + s.eventFilterMaxSpanBlocks - 1
		if end > targetEth1BlockHeight {
			end = targetEth1BlockHeight
		}
		withdrawIter, err := s.feePoolContract.FilterEtherWithdrawn(&bind.FilterOpts{
			Start:   i,
			End:     &end,
			Context: context.Background(),
		})
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("filter ether withdrawn failed: %w", err)
		}
		for withdrawIter.Next() {
			block := withdrawIter.Event.Raw.BlockNumber
			if _, ok := withdrawals[block]; !ok {
				withdrawals[block] = big.NewInt(0)
			}
			withdrawals[block] = new(big.Int).Add(withdrawals[block], withdrawIter.Event.Amount)
		}
	}
	log.Debug("end filter all withdrawn events")

	for i := latestDistributeHeight + 1; i <= targetEth1BlockHeight; i++ {
		// report progress in every 30 blocks
		if (i-latestDistributeHeight)%30 == 0 {
			log.WithFields(logrus.Fields{
				"block":    i,
				"progress": float64(i-latestDistributeHeight) / float64(targetEth1BlockHeight-latestDistributeHeight),
			}).Debug("report progress")
		}

		block, err := s.getBeaconBlock(i)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
		}

		// cal priority fee at this block
		feePoolPreBalance, err := s.getFeePoolBalance(i - 1)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
		}
		feePoolCurBalance, err := s.getFeePoolBalance(i)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
		}

		decreaseAmount, ok := withdrawals[i]
		if !ok {
			decreaseAmount = big.NewInt(0)
		}
		totalFeePoolCurBalance := new(big.Int).Add(feePoolCurBalance, decreaseAmount)
		if totalFeePoolCurBalance.Cmp(feePoolPreBalance) < 0 {
			return decimal.Zero, decimal.Zero, decimal.Zero, nil, fmt.Errorf("should not happened here when cal priority fee, block: %d", i)
		}
		feeAmountAtThisBlock := decimal.NewFromBigInt(new(big.Int).Sub(totalFeePoolCurBalance, feePoolPreBalance), 0)

		var userRewardDeci, nodeRewardDeci, platformFeeDeci decimal.Decimal
		val, _ := s.getValidatorByIndex(block.ProposerIndex)
		if val == nil {
			if feeAmountAtThisBlock.GreaterThan(decimal.Zero) {
				nodeRewardDeci = decimal.Zero
				platformFeeDeci = feeAmountAtThisBlock.Mul(s.platformCommissionRate).Floor()
				userRewardDeci = feeAmountAtThisBlock.Sub(platformFeeDeci.Add(nodeRewardDeci))
				log.WithFields(logrus.Fields{
					"block":  i,
					"amount": feeAmountAtThisBlock.DivRound(decimal.NewFromInt(1e18), 18).StringFixed(18),
				}).Debug("found transferFee")
			}
		} else {
			// get transfered fee from trace call
			trace, err := s.connection.Eth1Client().Debug_TraceBlockByNumber(ctx, big.NewInt(int64(i)), connection.Tracer{Tracer: "callTracer"})
			if err != nil {
				return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
			}
			transferFee := decimal.Zero
			seekFn := func(tx *connection.TxTrace) bool {
				return utils.In(s.transferFeeAddresses, strings.ToLower(tx.From)) &&
					strings.EqualFold(tx.To, s.feePoolAddress.String())
			}
			for _, tx := range trace {
				amount := WalkTrace(seekFn, decimal.Zero, tx.Result)
				if amount.GreaterThan(decimal.Zero) {
					transferFee = transferFee.Add(amount)
					log.WithFields(logrus.Fields{
						"block":  i,
						"txHash": tx.TxHash.Hex(),
						"amount": amount.DivRound(decimal.NewFromInt(1e18), 18).StringFixed(18),
					}).Debug("found transferFee")
				}
			}
			if transferFee.GreaterThan(decimal.Zero) {
				_platformFeeDeci := transferFee.Mul(s.platformCommissionRate).Floor()
				_userRewardDeci := transferFee.Sub(_platformFeeDeci)
				userRewardDeci = userRewardDeci.Add(_userRewardDeci)
				platformFeeDeci = platformFeeDeci.Add(_platformFeeDeci)
			}

			// only distribute tip fee to node
			tipFee := feeAmountAtThisBlock.Sub(transferFee)
			if tipFee.GreaterThan(decimal.Zero) {
				// cal rewards
				_userRewardDeci, _nodeRewardDeci, _platformFeeDeci := utils.GetUserNodePlatformReward(s.nodeCommissionRate, s.platformCommissionRate, val.NodeDepositAmountDeci, tipFee)
				userRewardDeci = userRewardDeci.Add(_userRewardDeci)
				nodeRewardDeci = nodeRewardDeci.Add(_nodeRewardDeci)
				platformFeeDeci = platformFeeDeci.Add(_platformFeeDeci)

				// cal node reward
				nodeNewReward, exist := nodeNewRewardsMap[val.NodeAddress]
				if exist {
					nodeNewReward.TotalRewardAmount = nodeNewReward.TotalRewardAmount.Add(nodeRewardDeci)
				} else {
					n := NodeNewReward{
						Address:                val.NodeAddress.String(),
						TotalRewardAmount:      nodeRewardDeci,
						TotalExitDepositAmount: decimal.Zero,
					}
					nodeNewRewardsMap[val.NodeAddress] = &n
				}
			}
		}

		// cal total vals
		totalUserEthDeci = totalUserEthDeci.Add(userRewardDeci)
		totalNodeEthDeci = totalNodeEthDeci.Add(nodeRewardDeci)
		totalPlatformEthDeci = totalPlatformEthDeci.Add(platformFeeDeci)
	}
	log.WithFields(logrus.Fields{
		"progress": float64(1),
	}).Debug("report progress: finished")

	{
		// hotfix: distribute blocked transfer fee
		feePoolBalance, err := s.getFeePoolBalance(targetEth1BlockHeight)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, nil, err
		}
		feePoolBalanceDeci := decimal.NewFromBigInt(feePoolBalance, 0)
		blockedTransferFeeDeci := feePoolBalanceDeci.Sub(totalUserEthDeci.Add(totalNodeEthDeci).Add(totalPlatformEthDeci))
		if blockedTransferFeeDeci.GreaterThan(decimal.Zero) {
			maxAmountPerEra := decimal.NewFromInt(int64(s.manager.cfg.DistributeBlockedTransferFeePerEra)).Mul(utils.EtherDeci)
			currentDistributeAmount := blockedTransferFeeDeci
			if currentDistributeAmount.GreaterThan(maxAmountPerEra) {
				currentDistributeAmount = maxAmountPerEra
			}

			platformFeeDeci := currentDistributeAmount.Mul(s.platformCommissionRate).Floor()
			userRewardDeci := currentDistributeAmount.Sub(platformFeeDeci)
			log.WithFields(logrus.Fields{
				"block":  targetEth1BlockHeight,
				"amount": currentDistributeAmount.DivRound(decimal.NewFromInt(1e18), 18).StringFixed(18),
			}).Debug("distribute blocked transferee fee")

			totalUserEthDeci = totalUserEthDeci.Add(userRewardDeci)
			totalPlatformEthDeci = totalPlatformEthDeci.Add(platformFeeDeci)
		}
	}

	return totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, nodeNewRewardsMap, nil
}

// include withdrawals fee
func (s *Service) getNodeNewRewardsBetween(latestDistributeHeight, targetEth1BlockHeight uint64) (NodeNewRewardsMap, error) {
	_, _, _, nodeNewRewardsMapFromWithdrawals, err := s.getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return nil, err
	}
	log := s.log.WithFields(logrus.Fields{
		"getNodeNewRewardsBetween": true,
	})
	_, _, _, nodeNewRewardsMapFromPriorityFee, err := s.getUserNodePlatformFromPriorityFee(log, latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return nil, err
	}

	finalNodeRewardsMap := make(NodeNewRewardsMap, 0)
	for _, node := range nodeNewRewardsMapFromWithdrawals {
		address := common.HexToAddress(node.Address)
		f, exist := finalNodeRewardsMap[address]
		if exist {
			f.TotalRewardAmount = f.TotalRewardAmount.Add(node.TotalRewardAmount)
			f.TotalExitDepositAmount = f.TotalExitDepositAmount.Add(node.TotalExitDepositAmount)
		} else {
			finalNodeRewardsMap[address] = &NodeNewReward{
				Address:                node.Address,
				TotalRewardAmount:      node.TotalRewardAmount,
				TotalExitDepositAmount: node.TotalExitDepositAmount,
			}
		}
	}

	for _, node := range nodeNewRewardsMapFromPriorityFee {
		address := common.HexToAddress(node.Address)
		f, exist := finalNodeRewardsMap[address]
		if exist {
			f.TotalRewardAmount = f.TotalRewardAmount.Add(node.TotalRewardAmount)
			f.TotalExitDepositAmount = f.TotalExitDepositAmount.Add(node.TotalExitDepositAmount)
		} else {
			finalNodeRewardsMap[address] = &NodeNewReward{
				Address:                node.Address,
				TotalRewardAmount:      node.TotalRewardAmount,
				TotalExitDepositAmount: node.TotalExitDepositAmount,
			}
		}
	}

	return finalNodeRewardsMap, nil
}

func (s *Service) getValidatorsOfTargetEpoch(ctx context.Context, targetEpoch uint64) ([]*Validator, error) {
	vals := make([]*Validator, 0)
	pubkeys := make([]types.ValidatorPubkey, 0)
	for _, val := range s.validators {
		if val.Status == 3 || val.Status > 4 {
			pubkeys = append(pubkeys, types.ValidatorPubkey(val.Pubkey))
		}
	}
	if len(pubkeys) == 0 {
		return nil, nil
	}

	validatorStatusMap, err := s.connection.GetValidatorStatuses(ctx, pubkeys, &beacon.ValidatorStatusOptions{
		Epoch: &targetEpoch,
	})
	if err != nil {
		return nil, errors.Wrap(err, "syncValidatorLatestInfo GetValidatorStatuses failed")
	}

	s.log.WithFields(logrus.Fields{
		"validatorStatuses len": len(validatorStatusMap),
	}).Debug("validator statuses")

	for pubkey, status := range validatorStatusMap {
		pubkeyStr := pubkey.String()
		if status.Exists {
			// must exist here
			valOfLatest, exist := s.validators[pubkeyStr]
			if !exist {
				return nil, fmt.Errorf("validator %s not exist", pubkeyStr)
			}

			// copy value
			validator := *valOfLatest

			updateBaseInfo := func() {
				// validator's info may be inited at any status
				validator.ActiveEpoch = status.ActivationEpoch
				validator.EligibleEpoch = status.ActivationEligibilityEpoch
				validator.ValidatorIndex = status.Index

				exitEpoch := status.ExitEpoch
				if exitEpoch == math.MaxUint64 {
					exitEpoch = 0
				}
				withdrawableEpoch := status.WithdrawableEpoch
				if withdrawableEpoch == math.MaxUint64 {
					withdrawableEpoch = 0
				}

				validator.ExitEpoch = exitEpoch
				validator.WithdrawableEpoch = withdrawableEpoch
			}

			updateBalance := func() {
				validator.Balance = status.Balance
				validator.EffectiveBalance = status.EffectiveBalance
			}

			switch status.Status {

			case ethpb.ValidatorStatus_PENDING_INITIALIZED, ethpb.ValidatorStatus_PENDING_QUEUED: // pending
				validator.Status = utils.ValidatorStatusWaiting
				validator.ValidatorIndex = status.Index

			case ethpb.ValidatorStatus_ACTIVE_ONGOING, ethpb.ValidatorStatus_ACTIVE_EXITING, ethpb.ValidatorStatus_ACTIVE_SLASHED: // active
				validator.Status = utils.ValidatorStatusActive
				if status.Slashed {
					validator.Status = utils.ValidatorStatusActiveSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_EXITED_UNSLASHED, ethpb.ValidatorStatus_EXITED_SLASHED: // exited
				validator.Status = utils.ValidatorStatusExited
				if status.Slashed {
					validator.Status = utils.ValidatorStatusExitedSlash
				}
				updateBaseInfo()
				updateBalance()
			case ethpb.ValidatorStatus_WITHDRAWAL_POSSIBLE: // withdrawable
				validator.Status = utils.ValidatorStatusWithdrawable
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawableSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_WITHDRAWAL_DONE: // withdrawdone
				validator.Status = utils.ValidatorStatusWithdrawDone
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawDoneSlash
				}
				updateBaseInfo()
				updateBalance()
			default:
				return nil, fmt.Errorf("unsupported validator status %d", status.Status)
			}

			if validator.ActiveEpoch < targetEpoch {
				vals = append(vals, &validator)
			}
		}
	}

	return vals, nil
}

func (s *Service) exitButNotFullWithdrawedValidatorListAtEpoch(ctx context.Context, epoch uint64) ([]*Validator, error) {
	vals := make([]*Validator, 0)

	pubkeys := make([]types.ValidatorPubkey, 0)
	for _, val := range s.validators {
		// skip not already actived vals
		if val.ActiveEpoch == 0 || val.ActiveEpoch > epoch {
			continue
		}
		if val.Status == 3 || val.Status > 4 {
			pubkeys = append(pubkeys, types.ValidatorPubkey(val.Pubkey))
		}
	}
	if len(pubkeys) == 0 {
		return nil, nil
	}

	validatorStatusMap, err := s.connection.GetValidatorStatuses(ctx, pubkeys, &beacon.ValidatorStatusOptions{
		Epoch: &epoch,
	})
	if err != nil {
		return nil, errors.Wrap(err, "exitButNotFullWithdrawedValidatorListAtEpoch GetValidatorStatuses failed")
	}

	s.log.WithFields(logrus.Fields{
		"validatorStatuses len": len(validatorStatusMap),
	}).Debug("validator statuses")

	for pubkey, status := range validatorStatusMap {
		pubkeyStr := pubkey.String()
		if !status.Exists {
			return nil, fmt.Errorf("validator %s status not exist on beacon", pubkeyStr)
		}

		// must exist here
		validatorCached, exist := s.validators[pubkeyStr]
		if !exist {
			return nil, fmt.Errorf("validator %s not exist in cached", pubkeyStr)
		}

		// copy val
		validator := *validatorCached

		updateBaseInfo := func() {
			// validator's info may be inited at any status
			validator.ActiveEpoch = status.ActivationEpoch
			validator.EligibleEpoch = status.ActivationEligibilityEpoch
			validator.ValidatorIndex = status.Index

			exitEpoch := status.ExitEpoch
			if exitEpoch == math.MaxUint64 {
				exitEpoch = 0
			}
			validator.ExitEpoch = exitEpoch

			withdrawableEpoch := status.WithdrawableEpoch
			if withdrawableEpoch == math.MaxUint64 {
				withdrawableEpoch = 0
			}
			validator.WithdrawableEpoch = withdrawableEpoch
		}

		updateBalance := func() {
			validator.Balance = status.Balance
			validator.EffectiveBalance = status.EffectiveBalance
		}

		switch status.Status {

		case ethpb.ValidatorStatus_ACTIVE_ONGOING, ethpb.ValidatorStatus_ACTIVE_EXITING, ethpb.ValidatorStatus_ACTIVE_SLASHED: // active
			validator.Status = utils.ValidatorStatusActive
			if status.Slashed {
				validator.Status = utils.ValidatorStatusActiveSlash
			}
			updateBaseInfo()
			updateBalance()

		case ethpb.ValidatorStatus_EXITED_UNSLASHED, ethpb.ValidatorStatus_EXITED_SLASHED: // exited
			validator.Status = utils.ValidatorStatusExited
			if status.Slashed {
				validator.Status = utils.ValidatorStatusExitedSlash
			}
			updateBaseInfo()
			updateBalance()
		case ethpb.ValidatorStatus_WITHDRAWAL_POSSIBLE: // withdrawable
			validator.Status = utils.ValidatorStatusWithdrawable
			if status.Slashed {
				validator.Status = utils.ValidatorStatusWithdrawableSlash
			}
			updateBaseInfo()
			updateBalance()

		case ethpb.ValidatorStatus_WITHDRAWAL_DONE: // withdrawdone
			validator.Status = utils.ValidatorStatusWithdrawDone
			if status.Slashed {
				validator.Status = utils.ValidatorStatusWithdrawDoneSlash
			}
			updateBaseInfo()
			updateBalance()
		default:
			return nil, fmt.Errorf("exitButNotFullWithdrawedValidatorListAtEpoch, unsupported validator: %s status %d", pubkeyStr, status.Status)
		}

		// filter exited and not full withdrawed
		if validator.ExitEpoch > 0 && validator.Balance > 0 {
			vals = append(vals, &validator)
		}

	}

	return vals, nil
}
