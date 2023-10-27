package service

import "github.com/sirupsen/logrus"

func (s *Service) pruneBlocks() error {

	latestMerkleRootEpochStartBlock := uint64(0)
	if s.latestMerkleRootEpoch != 0 {
		latestMerkleRootEpochStartBlockRes, err := s.getEpochStartBlocknumberWithCheck(s.latestMerkleRootEpoch)
		if err != nil {
			return err
		}
		latestMerkleRootEpochStartBlock = latestMerkleRootEpochStartBlockRes
	}

	minHeight := s.latestDistributionPriorityFeeHeight
	if minHeight > s.latestDistributionWithdrawalHeight {
		minHeight = s.latestDistributionWithdrawalHeight
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
