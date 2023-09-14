package service

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

var ErrExceedsValidatorUpdateBlock = fmt.Errorf("ErrExceedsValidatorUpdateBlock")

func (s *Service) syncBlocks() error {
	beaconHead, err := s.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	if beaconHead.FinalizedSlot <= s.latestSlotOfSyncBlock {
		return nil
	}

	start := uint64(s.latestSlotOfSyncBlock + 1)
	end := beaconHead.FinalizedSlot

	g := new(errgroup.Group)
	groupLimit := 16
	g.SetLimit(groupLimit)

	for i := start; i <= end; i++ {
		subStart := i
		subEnd := i + uint64(groupLimit) - 1
		if end < i+uint64(groupLimit) {
			subEnd = end
		}
		preLatestSyncBlock := s.latestBlockOfSyncBlock

		for j := subStart; j <= subEnd; j++ {
			g.Go(func() error {
				beaconBlock, exist, err := s.connection.GetBeaconBlock(j)
				if err != nil {
					return err
				}
				if !exist {
					return nil
				}
				// wait validator updated
				if beaconBlock.ExecutionBlockNumber > s.latestBlockOfUpdateValidator {
					return ErrExceedsValidatorUpdateBlock
				}
				_, isPoolVal := s.getValidatorByIndex(beaconBlock.ProposerIndex)
				if isPoolVal {
					fee, err := s.connection.GetELRewardForBlock(beaconBlock.ExecutionBlockNumber)
					if err != nil {
						return err
					}
					beaconBlock.PriorityFee = fee
				}

				s.cachedBeaconBlockMutex.Lock()
				s.cachedBeaconBlock[beaconBlock.ExecutionBlockNumber] = &beaconBlock
				s.cachedBeaconBlockMutex.Unlock()

				if beaconBlock.ExecutionBlockNumber > s.latestBlockOfSyncBlock {
					s.latestBlockOfSyncBlock = beaconBlock.ExecutionBlockNumber
				}

				return nil
			})
		}

		err = g.Wait()
		if err != nil {
			s.latestBlockOfSyncBlock = preLatestSyncBlock

			if err == ErrExceedsValidatorUpdateBlock {
				return nil
			}
			return err
		}

		s.latestSlotOfSyncBlock = subEnd

	}

	return nil
}
