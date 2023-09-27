package service

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/sirupsen/logrus"
)

const (
	fetchEventBlockLimit      = uint64(4900)
	fetchEth1WaitBlockNumbers = uint64(2)
	depositEventPreBlocks     = 14400 // 2days
)

func (s *Service) syncEvents() error {
	latestBlockNumber, err := s.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}

	latestDistributeWithdrawalsHeight, err := s.networkWithdrawContract.LatestDistributeWithdrawalsHeight(nil)
	if err != nil {
		return err
	}
	s.latestDistributeWithdrawalsHeight = latestDistributeWithdrawalsHeight.Uint64()

	latestDistributePriorityFeeHeight, err := s.networkWithdrawContract.LatestDistributePriorityFeeHeight(nil)
	if err != nil {
		return err
	}
	s.latestDistributePriorityFeeHeight = latestDistributePriorityFeeHeight.Uint64()

	latestMerkleRootEpoch, err := s.networkWithdrawContract.LatestMerkleRootEpoch(nil)
	if err != nil {
		return err
	}
	s.latestMerkleRootEpoch = latestMerkleRootEpoch.Uint64()

	if latestBlockNumber > fetchEth1WaitBlockNumbers {
		latestBlockNumber -= fetchEth1WaitBlockNumbers
	}

	logrus.Debugf("latestBlockNumber: %d, latestBlockOfSyncDeposit: %d", latestBlockNumber, s.latestBlockOfSyncEvents)

	if latestBlockNumber <= uint64(s.latestBlockOfSyncEvents) {
		return nil
	}

	start := uint64(s.latestBlockOfSyncEvents + 1)
	end := latestBlockNumber

	for i := start; i <= end; i += fetchEventBlockLimit {
		subStart := i
		subEnd := i + fetchEventBlockLimit - 1
		if end < i+fetchEventBlockLimit {
			subEnd = end
		}

		err = s.fetchDepositContractEventsAndCache(subStart, subEnd)
		if err != nil {
			return err
		}

		err = s.fetchExitElectionEventAndCache(subStart, subEnd)
		if err != nil {
			return err
		}

		// update
		s.latestBlockOfSyncEvents = subEnd

		logrus.WithFields(logrus.Fields{
			"start": subStart,
			"end":   subEnd,
		}).Debug("syncDepositInfo already dealed blocks")
	}

	return nil
}

func (s *Service) fetchDepositContractEventsAndCache(start, end uint64) error {
	iterDeposited, err := s.govDepositContract.FilterDepositEvent(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	for iterDeposited.Next() {
		pubkeyStr := hex.EncodeToString(iterDeposited.Event.Pubkey)

		s.govDeposits[pubkeyStr] = append(s.govDeposits[pubkeyStr], iterDeposited.Event.WithdrawalCredentials)
	}

	return nil
}

func (s *Service) fetchExitElectionEventAndCache(start, end uint64) error {
	iter, err := s.networkWithdrawContract.FilterNotifyValidatorExit(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	for iter.Next() {
		cycle := iter.Event.WithdrawCycle.Uint64()

		valList := make([]uint64, 0)
		for _, val := range iter.Event.EjectedValidators {
			valList = append(valList, val.Uint64())
		}

		s.exitElections[cycle] = &ExitElection{
			WithdrawCycle:      cycle,
			ValidatorIndexList: valList,
		}
	}

	return nil
}
