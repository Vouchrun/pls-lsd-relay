package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

var (
	oneEther     = decimal.NewFromInt(1e18)
	platformRate = decimal.RequireFromString("0.05")
)

// pls builds a decimal of the given whole-PLS amount scaled to 18 decimals.
func pls(amount string) decimal.Decimal {
	return decimal.RequireFromString(amount).Mul(oneEther)
}

// legacy is a copy of the original inline reconciliation block from
// service/common.go, used to prove the extracted helper is equivalent.
func legacy(totalUser, totalNode, totalPlatform, feePoolBalance, rate decimal.Decimal, fromBlock, toBlock uint64) (decimal.Decimal, decimal.Decimal, error) {
	overpaidAmount := totalUser.Add(totalNode).Add(totalPlatform).Sub(feePoolBalance)
	if overpaidAmount.GreaterThan(decimal.Zero) {
		platformFee := overpaidAmount.Mul(rate).Floor()
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

// 1. Incident reproduction: the live failure (25 distribution cycles in one range).
func TestReconcileIncidentReproduction(t *testing.T) {
	_, _, err := applyOverpaidCorrection(
		pls("20736409"), pls("1494660"), pls("1170056"),
		pls("136503.63"), platformRate,
		27596634, 27654906)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "total user eth less than zero") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 2. Normal cycle: feePoolBalance >= totals, so overpaid <= 0 and nothing changes.
func TestReconcileNormalCycle(t *testing.T) {
	totalUser, totalNode, totalPlatform := pls("100"), pls("50"), pls("10")
	feePoolBalance := pls("200")

	gotUser, gotPlatform, err := applyOverpaidCorrection(
		totalUser, totalNode, totalPlatform, feePoolBalance, platformRate, 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotUser.Equal(totalUser) {
		t.Fatalf("totalUser changed: got %s want %s", gotUser, totalUser)
	}
	if !gotPlatform.Equal(totalPlatform) {
		t.Fatalf("totalPlatform changed: got %s want %s", gotPlatform, totalPlatform)
	}
}

// 3. Small overpaid within both shares: totals reduced, both stay non-negative.
func TestReconcileSmallOverpaid(t *testing.T) {
	totalUser, totalNode, totalPlatform := pls("100"), pls("50"), pls("10")
	feePoolBalance := pls("140") // overpaid = 160 - 140 = 20

	gotUser, gotPlatform, err := applyOverpaidCorrection(
		totalUser, totalNode, totalPlatform, feePoolBalance, platformRate, 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := pls("81"); !gotUser.Equal(want) { // 100 - (20 - floor(20*0.05))
		t.Fatalf("totalUser: got %s want %s", gotUser, want)
	}
	if want := pls("9"); !gotPlatform.Equal(want) { // 10 - floor(20*0.05)
		t.Fatalf("totalPlatform: got %s want %s", gotPlatform, want)
	}
}

// 4. The correction equalizes the totals to feePoolBalance when overpaid > 0
// (userReward + platformFee == overpaid), so the trailing guard passes.
func TestReconcileBalanceEqualsFeePoolAfterCorrection(t *testing.T) {
	totalUser, totalNode, totalPlatform := pls("100"), pls("50"), pls("10")
	feePoolBalance := pls("140")

	gotUser, gotPlatform, err := applyOverpaidCorrection(
		totalUser, totalNode, totalPlatform, feePoolBalance, platformRate, 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := gotUser.Add(totalNode).Add(gotPlatform)
	if !sum.Equal(feePoolBalance) {
		t.Fatalf("corrected totals must equal feePoolBalance: got %s want %s", sum, feePoolBalance)
	}
}

// 5. Equivalence with the legacy inline math on all passing cases.
func TestReconcileEquivalenceLegacy(t *testing.T) {
	rate := decimal.RequireFromString("0.125")
	cases := []struct {
		name                          string
		user, node, platform, balance decimal.Decimal
	}{
		{"normal", pls("1000"), pls("200"), pls("50"), pls("1300")},
		{"exact balance", pls("1000"), pls("200"), pls("50"), pls("1250")},
		{"small overpaid", pls("1000"), pls("200"), pls("50"), pls("1240")},
		{"zero totals", decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero},
		{"fractional", pls("10.5"), pls("3.25"), pls("1.75"), pls("15.1")},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotUser, gotPlatform, gotErr := applyOverpaidCorrection(
				c.user, c.node, c.platform, c.balance, rate, 100, 200)
			wantUser, wantPlatform, wantErr := legacy(
				c.user, c.node, c.platform, c.balance, rate, 100, 200)

			if (gotErr == nil) != (wantErr == nil) {
				t.Fatalf("error mismatch: got %v want %v", gotErr, wantErr)
			}
			if gotErr != nil && gotErr.Error() != wantErr.Error() {
				t.Fatalf("error text mismatch: got %q want %q", gotErr.Error(), wantErr.Error())
			}
			if !gotUser.Equal(wantUser) {
				t.Fatalf("totalUser mismatch: got %s want %s", gotUser, wantUser)
			}
			if !gotPlatform.Equal(wantPlatform) {
				t.Fatalf("totalPlatform mismatch: got %s want %s", gotPlatform, wantPlatform)
			}
		})
	}
}
