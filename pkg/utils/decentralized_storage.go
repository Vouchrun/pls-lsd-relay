package utils

import (
	"fmt"
	"strings"
)

var nodeRewardsFileNameRaw string = "%s-nodeRewards-%d.json"

func NodeRewardsFileNameAtEpoch(lsdToken string, epoch uint64) string {
	return fmt.Sprintf(nodeRewardsFileNameRaw, strings.ToLower(lsdToken), epoch)
}

func NodeRewardsFileNameAtEpochOld(lsdToken string, epoch uint64) string {
	return fmt.Sprintf(nodeRewardsFileNameRaw, lsdToken, epoch)
}
