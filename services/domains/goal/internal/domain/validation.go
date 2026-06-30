package domain

import (
	"errors"
	"strings"
	"time"
)

func ValidateGoalName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.New("goal name is required")
	}
	if len([]rune(trimmed)) > 200 {
		return errors.New("goal name must be 200 characters or fewer")
	}
	return nil
}

func ValidatePriority(p int) error {
	if p < 1 {
		return errors.New("priority must be a positive integer")
	}
	return nil
}

func ValidateImportance(s string) error {
	if _, ok := AllImportance[GoalImportance(s)]; !ok {
		return errors.New("importance must be Mandatory, Essential, Lifestyle, Dream, or Speculative")
	}
	return nil
}

func ValidateTargetDate(td time.Time, created time.Time) error {
	if td.Before(created.Truncate(24 * time.Hour)) {
		return errors.New("target date must be in the future")
	}
	return nil
}

func ValidateSuccessCriteria(sc SuccessCriteria) error {
	switch sc.Model {
	case ModelTargetAmount, ModelMonthlyPassiveIncome, ModelNetWorth:
		if sc.TargetValue == nil || *sc.TargetValue <= 0 {
			return errors.New("target amount must be a positive value")
		}
	case ModelDebtFree:
	case ModelEmergencyFundMonths:
		if sc.TargetMonths == nil || *sc.TargetMonths <= 0 {
			return errors.New("target months must be a positive value")
		}
	case ModelCustomKPI:
		if sc.CustomDesc == nil || *sc.CustomDesc == "" {
			return errors.New("custom description is required for CustomKPI")
		}
		if sc.CustomTargetVal == nil || *sc.CustomTargetVal <= 0 {
			return errors.New("custom target value must be positive")
		}
	default:
		return errors.New("invalid success criteria model")
	}
	return nil
}

func ValidateRiskTolerance(s string) error {
	if s == "" {
		return nil
	}
	if !AllRiskTolerance[RiskTolerance(s)] {
		return errors.New("risk tolerance must be Conservative, Moderate, or Aggressive")
	}
	return nil
}
