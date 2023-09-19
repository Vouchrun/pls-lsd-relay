package service

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
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
	g.SetLimit(int(s.BatchRequestBlocksNumber))

	for i := start; i <= end; i += s.BatchRequestBlocksNumber {
		subStart := i
		subEnd := i + s.BatchRequestBlocksNumber - 1
		if end < i+s.BatchRequestBlocksNumber {
			subEnd = end
		}
		preLatestSyncBlock := s.latestBlockOfSyncBlock
		batchRequestStartTime := time.Now().Unix()

		blockReciever := make([]*beacon.BeaconBlock, s.BatchRequestBlocksNumber)
		for j := subStart; j <= subEnd; j++ {
			// notice this
			slot := j
			g.Go(func() error {
				startTime := time.Now().Unix()
				beaconBlock, exist, err := s.connection.GetBeaconBlock(slot)
				if err != nil {
					return err
				}
				endTime := time.Now().Unix()

				// maybe not exist this slot
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

				blockReciever[slot-subStart] = &beaconBlock

				saveTime := time.Now().Unix()
				logrus.Tracef("request block %d,start at %d, end at %d, save at: %d ", beaconBlock.ExecutionBlockNumber, startTime, endTime, saveTime)
				return nil
			})
		}

		err = g.Wait()
		if err != nil {
			s.latestBlockOfSyncBlock = preLatestSyncBlock
			if err == ErrExceedsValidatorUpdateBlock {
				return nil
			}

			logrus.Warnf("sync block err: %s, will retry", err.Error())
			return err
		}

		batchRequestWaitTime := time.Now().Unix()

		s.cachedBeaconBlockMutex.Lock()
		for _, beaconBlock := range blockReciever {
			if beaconBlock == nil {
				continue
			}
			logrus.Tracef("save block: %d", beaconBlock.ExecutionBlockNumber)

			s.cachedBeaconBlock[beaconBlock.ExecutionBlockNumber] = beaconBlock
			// update latest block
			if beaconBlock.ExecutionBlockNumber > s.latestBlockOfSyncBlock {
				s.latestBlockOfSyncBlock = beaconBlock.ExecutionBlockNumber
			}
		}
		s.cachedBeaconBlockMutex.Unlock()

		// update latest slot
		s.latestSlotOfSyncBlock = subEnd

		batchRequestEndTime := time.Now().Unix()
		logrus.Tracef("batch request block, start at: %d, wait at %d, end at %d", batchRequestStartTime, batchRequestWaitTime, batchRequestEndTime)
	}

	return nil
}
