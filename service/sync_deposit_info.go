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

func (s *Service) syncDepositInfo() error {
	latestBlockNumber, err := s.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}

	if latestBlockNumber > fetchEth1WaitBlockNumbers {
		latestBlockNumber -= fetchEth1WaitBlockNumbers
	}

	logrus.Debugf("latestBlockNumber: %d, dealedBlockNumber: %d", latestBlockNumber, s.dealedEth1Block)
	if latestBlockNumber <= uint64(s.dealedEth1Block) {
		return nil
	}

	start := uint64(s.dealedEth1Block + 1)
	end := latestBlockNumber

	for i := start; i <= end; i += fetchEventBlockLimit {
		subStart := i
		subEnd := i + fetchEventBlockLimit - 1
		if end < i+fetchEventBlockLimit {
			subEnd = end
		}

		err = s.fetchDepositContractEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"start": subStart,
			"end":   subEnd,
		}).Debug("already dealed blocks")
	}

	return nil
}

func (s *Service) fetchDepositContractEvents(start, end uint64) error {
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
