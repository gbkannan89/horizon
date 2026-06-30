package domain

import (
	"math"
	"time"
)

// C-01: OutstandingBalance = OriginalPrincipal − Sum(Payments) + AccruedInterest − Prepayments.
func CalculateOutstandingBalance(principal int64, totalPayments int64, totalInterest int64, totalPrepayments int64) int64 {
	b := principal - totalPayments + totalInterest - totalPrepayments
	if b < 0 { b = 0 }; return b
}

// C-02: Simple Interest Accrued.
func CalculateSimpleInterest(principal int64, rate float64, days int) int64 {
	return int64(float64(principal) * rate / 100.0 * float64(days) / 365.0)
}

// C-04: Principal Remaining.
func CalculatePrincipalRemaining(principal int64, totalPrincipalPaid int64) int64 {
	r := principal - totalPrincipalPaid; if r < 0 { r = 0 }; return r
}

// C-05: Interest Remaining = TotalInterest − InterestPaid.
func CalculateInterestRemaining(totalInterest, interestPaid int64) int64 {
	r := totalInterest - interestPaid; if r < 0 { r = 0 }; return r
}

// C-06: RepaymentProgress = PrincipalPaid / OriginalPrincipal * 100.
func CalculateRepaymentProgress(principalPaid, originalPrincipal int64) float64 {
	if originalPrincipal <= 0 { return 0 }
	return float64(principalPaid) / float64(originalPrincipal) * 100
}

// C-07: DebtUtilization = OutstandingBalance / OriginalPrincipal * 100.
func CalculateDebtUtilization(outstanding, original int64) float64 {
	if original <= 0 { return 0 }
	return float64(outstanding) / float64(original) * 100
}

// C-08: DebtBurden = TotalMonthlyPayments / MonthlyIncome * 100.
func CalculateDebtBurden(totalMonthlyPayment, monthlyIncome int64) float64 {
	if monthlyIncome <= 0 { return 0 }
	return float64(totalMonthlyPayment) / float64(monthlyIncome) * 100
}

// C-10: RemainingDuration in months.
func CalculateRemainingDuration(remainingInstallments int) int {
	return remainingInstallments
}

// C-12: LoanToValue = OutstandingBalance / CollateralValue * 100.
func CalculateLoanToValue(outstanding, collateralValue int64) float64 {
	if collateralValue <= 0 { return 0 }
	return float64(outstanding) / float64(collateralValue) * 100
}

// C-13: DebtServiceCoverage = NetIncome / TotalDebtService.
func CalculateDebtServiceCoverage(netIncome, totalDebtService int64) float64 {
	if totalDebtService <= 0 { return 0 }
	return float64(netIncome) / float64(totalDebtService)
}

// C-16: InterestBurden = TotalInterest / TotalPayments * 100.
func CalculateInterestBurden(totalInterest, totalPayments int64) float64 {
	if totalPayments <= 0 { return 0 }
	return float64(totalInterest) / float64(totalPayments) * 100
}

// C-19: DebtAging in days since last payment.
func CalculateDebtAging(lastPaymentDate time.Time, now time.Time) float64 {
	return now.Sub(lastPaymentDate).Hours() / 24
}

// C-03: Reducing Balance Interest.
func CalculateReducingBalanceInterest(outstanding int64, rate float64, months int) int64 {
	total := 0.0
	bal := float64(outstanding)
	mr := rate / 100.0 / 12.0
	for i := 0; i < months; i++ {
		interest := bal * mr
		total += interest
		bal += interest
	}
	return int64(total)
}

// C-11: EarlySettlementAmount = OutstandingBalance + Penalty + RemainingFees.
func CalculateEarlySettlementAmount(outstanding, penalty, remainingFees int64) int64 {
	return outstanding + penalty + remainingFees
}

// C-09: AverageInterestRate (balance-weighted).
func CalculateAverageInterestRate(balances []int64, rates []float64) float64 {
	totalBal := int64(0); weighted := 0.0
	for i, b := range balances { totalBal += b; weighted += float64(b) * rates[i] }
	if totalBal <= 0 { return 0 }
	return weighted / float64(totalBal)
}

// C-14: DebtExposure = OutstandingBalance.
func CalculateDebtExposure(outstanding int64) int64 { return outstanding }

// C-15: DebtConcentration = LiabilityOutstanding / TotalDebt * 100.
func CalculateDebtConcentration(liabilityOutstanding, totalDebt int64) float64 {
	if totalDebt <= 0 { return 0 }
	return float64(liabilityOutstanding) / float64(totalDebt) * 100
}

// C-17: RepaymentStability = OnTimePayments / TotalPayments * 100.
func CalculateRepaymentStability(ontime, total int) float64 {
	if total <= 0 { return 100 }
	return float64(ontime) / float64(total) * 100
}

// C-18: DebtVelocity = OutstandingChange / Period.
func CalculateDebtVelocity(outstandingOld, outstandingNew int64) int64 {
	return outstandingNew - outstandingOld
}

// C-20: DebtRisk = f(Health, Utilization, Burden, Aging).
func CalculateDebtRisk(health LiabilityHealth, utilization, burden float64, agingDays float64) float64 {
	risk := 0.0
	switch health {
	case LHHealthy: risk += 0.1
	case LHStable: risk += 0.25
	case LHWarning: risk += 0.5
	case LHCritical: risk += 0.75
	case LHDelinquentSt: risk += 0.9
	default: risk += 0.5
	}
	risk += (utilization / 100.0) * 0.2
	risk += (burden / 100.0) * 0.2
	risk += math.Min(agingDays/365.0, 1.0) * 0.1
	if risk > 1.0 { risk = 1.0 }
	return risk
}

// C-21: DebtQuality = 1 − DebtRisk.
func CalculateDebtQuality(risk float64) float64 {
	q := 1.0 - risk; if q < 0 { q = 0 }; return q
}
