package feecalc

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// ApplyOverpaidCorrection implements the priority-fee "subtract overpaid amount"
// reconciliation. Returns the corrected totalUser/totalPlatform (totalNode is unchanged)
// or an error when the reconciliation fails. Pure; no I/O.
func ApplyOverpaidCorrection(totalUser, totalNode, totalPlatform, feePoolBalance, platformRate decimal.Decimal, fromBlock, toBlock uint64) (decimal.Decimal, decimal.Decimal, error) {
	overpaidAmount := totalUser.Add(totalNode).Add(totalPlatform).Sub(feePoolBalance)
	if overpaidAmount.GreaterThan(decimal.Zero) {
		platformFee := overpaidAmount.Mul(platformRate).Floor()
		userReward := overpaidAmount.Sub(platformFee)
		totalUser = totalUser.Sub(userReward)
		totalPlatform = totalPlatform.Sub(platformFee)

		if totalUser.LessThan(decimal.Zero) {
			return totalUser, totalPlatform,
				fmt.Errorf("total user eth less than zero, from block: %d to block: %d", fromBlock, toBlock)
		}
		if totalPlatform.LessThan(decimal.Zero) {
			return totalUser, totalPlatform,
				fmt.Errorf("total platform eth less than zero, from block: %d to block: %d", fromBlock, toBlock)
		}
	}

	if feePoolBalance.LessThan(totalUser.Add(totalNode).Add(totalPlatform)) {
		return totalUser, totalPlatform,
			fmt.Errorf("fee pool balance less than total user eth, node eth, platform eth, from block: %d to block: %d", fromBlock, toBlock)
	}

	return totalUser, totalPlatform, nil
}
