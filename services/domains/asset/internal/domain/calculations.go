package domain

import (
	"math"
)

// C-01: Current Market Value.
func CalculateMarketValue(quantity, unitPrice float64) int64 {
	return int64(quantity * unitPrice)
}

// C-02: Book Value.
func CalculateBookValue(costBasis int64, accumulatedDepreciation int64, additionalInvestments int64) int64 {
	return costBasis - accumulatedDepreciation + additionalInvestments
}

// C-03: Unrealised Gain/Loss.
func CalculateUnrealisedGL(currentValue, costBasis int64) int64 {
	return currentValue - costBasis
}

// C-04: Realised Gain/Loss.
func CalculateRealisedGL(disposalProceeds, proRatedCostBasis int64) int64 {
	return disposalProceeds - proRatedCostBasis
}

// C-05: Adjusted Cost Basis (Post-Split).
func CalculateAdjustedCostBasis(costBasis int64, oldQty, newQty float64) int64 {
	if oldQty <= 0 {
		return costBasis
	}
	return int64(float64(costBasis) * (newQty / oldQty))
}

// C-06: Asset Allocation — returns classification percentage.
func CalculateAssetAllocation(assetValue int64, totalValue int64) float64 {
	if totalValue <= 0 {
		return 0
	}
	return float64(assetValue) / float64(totalValue) * 100
}

// C-07: Asset Concentration.
func CalculateAssetConcentration(assetValue, totalPortfolioValue int64) float64 {
	return CalculateAssetAllocation(assetValue, totalPortfolioValue)
}

// C-08: Liquidity Score (Immediate=1.0, Illiquid=0.0).
func CalculateLiquidityScore(lp LiquidityProfile) float64 {
	switch lp {
	case LQImmediate:
		return 1.0
	case LQShortTerm:
		return 0.85
	case LQMediumTerm:
		return 0.6
	case LQLongTerm:
		return 0.4
	case LQRestricted:
		return 0.2
	case LQLocked:
		return 0.1
	case LQIlliquid:
		return 0.0
	default:
		return 0.5
	}
}

// C-09: Valuation Change.
func CalculateValuationChange(currentValue, startValue int64) int64 {
	return currentValue - startValue
}

// C-11: Ownership Exposure.
func CalculateOwnershipExposure(currentValue int64, ownershipPct float64) int64 {
	return int64(float64(currentValue) * ownershipPct / 100.0)
}

// C-12: Currency Exposure.
func CalculateCurrencyExposure(currentValue int64, exchangeRate float64) int64 {
	return int64(float64(currentValue) * exchangeRate)
}

// C-14: Asset Risk Concentration.
func CalculateRiskConcentration(valueAtClassification, total int64) float64 {
	return CalculateAssetAllocation(valueAtClassification, total)
}

// C-16: Asset Growth.
func CalculateAssetGrowth(current, start float64) float64 {
	if start == 0 {
		return 0
	}
	return (current - start) / start * 100
}

// C-18: Asset Utilization.
func CalculateAssetUtilization(currentValue, maxHistoricalValue int64) float64 {
	if maxHistoricalValue <= 0 {
		return 100
	}
	return float64(currentValue) / float64(maxHistoricalValue) * 100
}

// C-19: Asset Volatility (simplified stddev from periodic returns).
func CalculateAssetVolatility(periodicReturns []float64) float64 {
	if len(periodicReturns) < 2 {
		return 0
	}
	mean := 0.0
	for _, r := range periodicReturns {
		mean += r
	}
	mean /= float64(len(periodicReturns))
	variance := 0.0
	for _, r := range periodicReturns {
		variance += (r - mean) * (r - mean)
	}
	variance /= float64(len(periodicReturns) - 1)
	return math.Sqrt(variance)
}

// C-13: Asset Diversification (HHI-based).
func CalculateDiversificationHHI(values []int64) float64 {
	total := int64(0)
	for _, v := range values {
		total += v
	}
	if total <= 0 {
		return 0
	}
	hhi := 0.0
	for _, v := range values {
		share := float64(v) / float64(total)
		hhi += share * share
	}
	return 1.0 - hhi
}

// C-15: Asset Quality.
func CalculateAssetQuality(confidence AssetConfidence, profile ValuationProfile, lp LiquidityProfile) float64 {
	score := 0.0
	switch confidence {
	case ACVerified:
		score += 0.4
	case ACMarketVerified:
		score += 0.35
	case ACInstitutionVerified:
		score += 0.3
	case ACUserVerified:
		score += 0.25
	case ACImported:
		score += 0.15
	default:
		score += 0.05
	}
	score += CalculateLiquidityScore(lp) * 0.3
	switch profile {
	case VPProfileMarketDriven, VPProfileExternalAppraisal:
		score += 0.3
	case VPProfileGovernmentAssessment, VPProfileRuleBased:
		score += 0.25
	case VPProfileManualAssessment, VPProfileIndexed:
		score += 0.15
	default:
		score += 0.05
	}
	if score > 1.0 {
		score = 1.0
	}
	return score
}
