package service

import "github.com/sirupsen/logrus"

func (s *Service) pruneBlocks() error {
	latestDistributeWithdrawalsHeight, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(nil)
	if err != nil {
		return err
	}

	latestDistributePriorityFeeHeight, err := s.networkWithdrawContract.LatestDistributePriorityFeeHeight(nil)
	if err != nil {
		return err
	}

	latestMerkleRootEpoch, err := s.networkWithdrawContract.LatestMerkleRootEpoch(nil)
	if err != nil {
		return err
	}

	latestMerkleRootEpochStartBlock := uint64(0)
	if latestMerkleRootEpoch.Uint64() != 0 {
		latestMerkleRootEpochStartBlockRes, err := s.getEpochStartBlocknumberWithCheck(latestMerkleRootEpoch.Uint64())
		if err != nil {
			return err
		}
		latestMerkleRootEpochStartBlock = latestMerkleRootEpochStartBlockRes
	}

	minHeight := latestDistributePriorityFeeHeight.Uint64()
	if minHeight > latestDistributeWithdrawalsHeight.Uint64() {
		minHeight = latestDistributeWithdrawalsHeight.Uint64()
	}
	if minHeight > latestMerkleRootEpochStartBlock {
		minHeight = latestMerkleRootEpochStartBlock
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
