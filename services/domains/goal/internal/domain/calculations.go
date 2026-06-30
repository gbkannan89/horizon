package domain

import (
	"math"
	"time"
)

func CalculateProgress(currentValue float64, sc SuccessCriteria) float64 {
	switch sc.Model {
	case ModelTargetAmount, ModelMonthlyPassiveIncome, ModelNetWorth:
		if sc.TargetValue == nil || *sc.TargetValue <= 0 {
			return 0
		}
		p := (currentValue / *sc.TargetValue) * 100
		if p > 100 {
			return 100
		}
		return p
	case ModelDebtFree:
		if currentValue >= 1 {
			return 100
		}
		return 0
	case ModelEmergencyFundMonths:
		if sc.TargetMonths == nil || *sc.TargetMonths <= 0 {
			return 0
		}
		if currentValue > float64(*sc.TargetMonths) {
			return 100
		}
		return (currentValue / float64(*sc.TargetMonths)) * 100
	case ModelCustomKPI:
		if sc.CustomTargetVal == nil || *sc.CustomTargetVal <= 0 {
			return 0
		}
		p := (currentValue / *sc.CustomTargetVal) * 100
		if p > 100 {
			return 100
		}
		return p
	}
	return 0
}

func CalculateRemaining(currentValue float64, sc SuccessCriteria) float64 {
	switch sc.Model {
	case ModelTargetAmount, ModelMonthlyPassiveIncome, ModelNetWorth:
		if sc.TargetValue == nil {
			return 0
		}
		return math.Max(0, *sc.TargetValue-currentValue)
	case ModelEmergencyFundMonths:
		if sc.TargetMonths == nil {
			return 0
		}
		return math.Max(0, float64(*sc.TargetMonths)-currentValue)
	case ModelDebtFree:
		if currentValue >= 1 {
			return 0
		}
		return currentValue
	case ModelCustomKPI:
		if sc.CustomTargetVal == nil {
			return 0
		}
		return math.Max(0, *sc.CustomTargetVal-currentValue)
	}
	return 0
}

func CalculateRequiredContribution(remaining float64, targetDate time.Time, now time.Time) float64 {
	remainingMonths := targetDate.Sub(now).Hours() / (30.44 * 24)
	if remainingMonths <= 0 {
		return remaining
	}
	return remaining / remainingMonths
}

func CalculateGoalHealth(progress, trajectory, deviation, timeRemainingRatio float64, importance GoalImportance) float64 {
	baseWeight := 1.0
	if w, ok := AllImportance[importance]; ok {
		baseWeight = float64(w) / 5.0
	}
	health := (progress/100.0)*0.3 + (1.0-math.Abs(deviation))*0.3 + timeRemainingRatio*0.2 + baseWeight*0.2
	if health < 0 {
		health = 0
	}
	if health > 1 {
		health = 1
	}
	_ = trajectory
	return health
}

func CalculateImportanceWeight(importance GoalImportance, priority int) float64 {
	base := 1.0
	if w, ok := AllImportance[importance]; ok {
		base = float64(w)
	}
	modifier := 1.0 / float64(priority)
	return base * modifier
}

func IsAtRisk(remaining, requiredContribution, timeRemaining float64, variabilityFactor float64) bool {
	if variabilityFactor < 0 {
		variabilityFactor = 0
	}
	return remaining > (requiredContribution * timeRemaining * (1.0 + variabilityFactor))
}
