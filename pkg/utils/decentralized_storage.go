package utils

import (
	"fmt"
	"strings"
)

func NodeRewardsFileNameAtEpoch(lsdToken string, chainID uint64, epoch uint64) string {
	return fmt.Sprintf("%s-rewards-%d-%d.json", strings.ToLower(lsdToken), chainID, epoch)
}

func NodeRewardsFileNameAtEpochOld(lsdToken string, epoch uint64) string {
	return fmt.Sprintf("%s-nodeRewards-%d.json", strings.ToLower(lsdToken), epoch)
}
