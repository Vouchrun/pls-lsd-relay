package service

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func (s *Service) pruneBlocks() error {

	latestMerkleRootEpochStartBlock := uint64(0)
	if s.latestMerkleRootEpoch != 0 {
		latestMerkleRootEpochStartBlockRes, err := s.getEpochStartBlocknumberWithCheck(s.latestMerkleRootEpoch)
		if err != nil {
			return err
		}
		latestMerkleRootEpochStartBlock = latestMerkleRootEpochStartBlockRes
	}

	minHeight := s.latestDistributePriorityFeeHeight
	if minHeight > s.latestDistributeWithdrawalsHeight {
		minHeight = s.latestDistributeWithdrawalsHeight
	}
	if minHeight > latestMerkleRootEpochStartBlock {
		minHeight = latestMerkleRootEpochStartBlock
	}

	_, targetTimestamp, err := s.currentCycleAndStartTimestamp()
	if err != nil {
		return fmt.Errorf("currentCycleAndStartTimestamp failed: %w", err)
	}
	targetEpoch := utils.EpochAtTimestamp(s.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := s.getEpochStartBlocknumberWithCheck(targetEpoch)
	if err != nil {
		return err
	}
	targetCall := s.connection.CallOpts(big.NewInt(int64(targetBlockNumber)))
	latestDistributeWithdrawalHeightOnCycleSnapshot, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(targetCall)
	if err != nil {
		return err
	}
	if minHeight > latestDistributeWithdrawalHeightOnCycleSnapshot.Uint64() {
		minHeight = s.latestDistributeWithdrawalsHeight
	}

	if minHeight == 0 {
		return nil
	}

	s.cachedBeaconBlockMutex.Lock()
	defer s.cachedBeaconBlockMutex.Unlock()

	maxHeight := uint64(0)
	for blockNumber := range s.cachedBeaconBlock {
		if blockNumber < minHeight {
			delete(s.cachedBeaconBlock, blockNumber)
		}
		if blockNumber > maxHeight {
			maxHeight = blockNumber
		}
	}

	logrus.Debugf("cachedBlocks minHeight: %d, maxHeight: %d", minHeight, maxHeight)

	return nil
}
