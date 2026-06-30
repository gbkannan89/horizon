package domain

import (
	"math"
	"time"
)

func CalculateNetWorth(totalAssets, totalCash, totalLiabilities int64) int64 { return totalAssets + totalCash - totalLiabilities }
func CalculateTotalValue(values []int64) int64 { var t int64; for _, v := range values { t += v }; return t }
func CalculateAllocationPct(classValue, totalValue int64) float64 { if totalValue <= 0 { return 0 }; return float64(classValue) / float64(totalValue) * 100 }
func CalculateDrift(current, target float64) float64 { return current - target }
func CalculateDiversificationScore(values []int64) float64 {
	total := int64(0); for _, v := range values { total += v }; if total <= 0 { return 0 }
	hhi := 0.0; for _, v := range values { s := float64(v) / float64(total); hhi += s * s }; return 1.0 - hhi
}
func CalculateConcentration(values []int64) float64 {
	total := int64(0); maxV := int64(0)
	for _, v := range values { total += v; if v > maxV { maxV = v } }; if total <= 0 { return 0 }; return float64(maxV) / float64(total) * 100
}
func CalculateLiquidityScore(liqVals []float64, values []int64) float64 {
	total := int64(0); weighted := 0.0
	for i, v := range values { total += v; weighted += liqVals[i] * float64(v) }; if total <= 0 { return 0 }; return weighted / float64(total)
}
func CalculateExposure(values []int64, pcts []float64) int64 {
	var total int64; for i, v := range values { total += int64(float64(v) * pcts[i] / 100.0) }; return total
}
func CalculatePeriodReturn(start, end int64) float64 { if start == 0 { return 0 }; return float64(end-start) / float64(start) * 100 }
func CalculateRebalanceNeed(drifts []float64) float64 { maxD := 0.0; for _, d := range drifts { if math.Abs(d) > maxD { maxD = math.Abs(d) } }; return maxD }
func CalculateRiskExposure(scores []float64, values []int64) float64 {
	total := int64(0); weighted := 0.0; for i, v := range values { total += v; weighted += scores[i] * float64(v) }; if total <= 0 { return 0 }; return weighted / float64(total)
}
func CalculateVolatility(returns []float64) float64 {
	if len(returns) < 2 { return 0 }; mean := 0.0
	for _, r := range returns { mean += r }; mean /= float64(len(returns))
	v := 0.0; for _, r := range returns { v += (r - mean) * (r - mean) }; return math.Sqrt(v / float64(len(returns)-1))
}
func CalculateSharpeRatio(portReturn, riskFreeRate, vol float64) float64 { if vol == 0 { return 0 }; return (portReturn - riskFreeRate) / vol }
func CalculateBeta(cov, mktVar float64) float64 { if mktVar == 0 { return 1 }; return cov / mktVar }
func CalculateAlpha(portReturn, rfr, beta, mktReturn float64) float64 { return portReturn - (rfr + beta*(mktReturn-rfr)) }
func CalculateAssetGrowth(cur, basis int64) float64 { if basis == 0 { return 0 }; return float64(cur-basis) / float64(basis) * 100 }
func CalculateAssetYield(income, avgValue int64) float64 { if avgValue == 0 { return 0 }; return float64(income) / float64(avgValue) * 100 }
func CalculateGoalMappingCoverage(mapped, total int64) float64 { if total <= 0 { return 0 }; return float64(mapped) / float64(total) * 100 }
func CalculateDrawdown(peak, trough int64) float64 { if peak <= 0 { return 0 }; return float64(peak-trough) / float64(peak) * 100 }
func CalculateRecoveryTime(peakDate, recoveryDate time.Time) float64 { return recoveryDate.Sub(peakDate).Hours() / 24 }
func CalculateWinRate(wins, total int) float64 { if total <= 0 { return 0 }; return float64(wins) / float64(total) * 100 }
