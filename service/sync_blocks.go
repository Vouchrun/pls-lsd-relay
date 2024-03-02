package service

import (
	"fmt"
	"time"

	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
	"golang.org/x/sync/errgroup"
)

var ErrExceedsValidatorUpdateBlock = fmt.Errorf("ErrExceedsValidatorUpdateBlock")
var ErrHandlerExit = fmt.Errorf("exit")
var ErrMissingEth1Block = fmt.Errorf("beacon chain missing eth1 block: %w", ErrHandlerExit)

// sync beacon and execution block info
func (s *Service) syncBlocks() error {
	beaconHead, err := s.connection.BeaconHead()
	if err != nil {
		return err
	}

	if beaconHead.FinalizedSlot <= s.latestSlotOfSyncBlock {
		return nil
	}
	latestSlotOfUpdateValidator := utils.EndSlotOfEpoch(s.eth2Config, s.latestEpochOfUpdateValidator)

	start := uint64(s.latestSlotOfSyncBlock + 1)
	end := beaconHead.FinalizedSlot
	if end > latestSlotOfUpdateValidator {
		end = latestSlotOfUpdateValidator
	}

	g := new(errgroup.Group)
	g.SetLimit(int(s.batchRequestBlocksNumber))

	for i := start; i <= end; i += s.batchRequestBlocksNumber {
		subStart := i
		subEnd := i + s.batchRequestBlocksNumber - 1
		if end < i+s.batchRequestBlocksNumber {
			subEnd = end
		}
		preLatestSyncBlock := s.latestBlockOfSyncBlock
		batchRequestStartTime := time.Now().Unix()

		blockReceiver := make([]*CachedBeaconBlock, s.batchRequestBlocksNumber)
		for j := subStart; j <= subEnd; j++ {
			// notice this
			slot := j
			g.Go(func() error {
				startTime := time.Now().Unix()
				beaconBlock, exist, err := s.manager.CacheBeaconBlock(slot)
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

				blockReceiver[slot-subStart] = beaconBlock

				saveTime := time.Now().Unix()
				s.log.Tracef("request block %d,start at %d, end at %d, save at: %d ", beaconBlock.ExecutionBlockNumber, startTime, endTime, saveTime)
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

		batchRequestWaitTime := time.Now().Unix()

		for _, beaconBlock := range blockReceiver {
			if beaconBlock == nil {
				continue
			}
			s.log.Tracef("save block: %d", beaconBlock.ExecutionBlockNumber)

			// update latest block
			if beaconBlock.ExecutionBlockNumber > s.latestBlockOfSyncBlock {
				if beaconBlock.ExecutionBlockNumber-s.latestBlockOfSyncBlock > 1 {
					// rpc error missing some blocks
					return fmt.Errorf("%w at slot: %d desired eth1 block: %d", ErrMissingEth1Block, beaconBlock.BeaconBlockId, s.latestBlockOfSyncBlock+1)
				}
				s.latestBlockOfSyncBlock = beaconBlock.ExecutionBlockNumber
			}
		}

		// update latest slot
		s.latestSlotOfSyncBlock = subEnd

		batchRequestEndTime := time.Now().Unix()
		s.log.Tracef("batch request block, start at: %d, wait at %d, end at %d", batchRequestStartTime, batchRequestWaitTime, batchRequestEndTime)
	}

	return nil
}
