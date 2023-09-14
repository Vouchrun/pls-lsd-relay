package service

import (
	"math/big"
)

func (s *Service) syncExitElection() error {
	currentCycle, err := s.networkWithdrawContract.CurrentWithdrawCycle(nil)
	if err != nil {
		return err
	}

	if currentCycle.Uint64()-1 <= s.latestWithdrawCycleOfSyncExitElection {
		return nil
	}

	for i := s.latestBlockOfSyncBlock; i <= currentCycle.Uint64()-1; i++ {
		vals, err := s.networkWithdrawContract.GetEjectedValidatorsAtCycle(nil, big.NewInt(int64(i)))
		if err != nil {
			return err
		}
		if len(vals) == 0 {
			s.latestWithdrawCycleOfSyncExitElection = i
			continue
		}

		valList := make([]uint64, 0)
		for _, val := range vals {
			valList = append(valList, val.Uint64())
		}
		s.exitElections[i] = &ExitElection{
			WithdrawCycle:      i,
			ValidatorIndexList: valList,
		}
		s.latestWithdrawCycleOfSyncExitElection = i
	}

	return nil
}
