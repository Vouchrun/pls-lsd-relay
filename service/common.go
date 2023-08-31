package service

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) getEpochStartBlocknumber(epoch uint64) (uint64, error) {
	eth2ValidatorBalanceSyncerStartSlot := utils.StartSlotOfEpoch(s.eth2Config, epoch)
	logrus.Debugf("getEpochStartBlocknumber: %d, epoch: %d ", eth2ValidatorBalanceSyncerStartSlot, epoch)

	retry := 0
	for {
		if retry > 10 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}

		targetBeaconBlock, exist, err := s.connection.GetBeaconBlock(eth2ValidatorBalanceSyncerStartSlot)
		if err != nil {
			return 0, err
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

// return (user reward, node reward, platform fee, totalWithdrawAmount) decimals 18
func (s *Service) getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight uint64) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, nil
}

// return (user reward, node reward, platform fee, totalWithdrawAmount) decimals 18
func (s *Service) getUserNodePlatformFromPriorityFee(latestDistributeHeight, targetEth1BlockHeight uint64) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, nil
}
