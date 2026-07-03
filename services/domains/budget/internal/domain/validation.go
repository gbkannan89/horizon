package domain

import (
	"errors"
	"strings"
	"time"
)

func ValidateBudgetName(n string) error {
	t := strings.TrimSpace(n)
	if t == "" { return errors.New("budget name is required") }
	if len([]rune(t)) > 200 { return errors.New("budget name must be 200 characters or fewer") }
	return nil
}

func ValidateBudgetPeriod(p BudgetPeriod) error {
	switch p {
	case PeriodWeekly, PeriodMonthly, PeriodQuarterly, PeriodYearly, PeriodCustom:
		return nil
	}
	return errors.New("budget period is not recognized")
}

func ValidateCurrency(c string) error {
	if len(c) != 3 { return errors.New("currency must be a 3-letter ISO code") }
	return nil
}

func ValidateDateRange(start, end time.Time) error {
	if end.Before(start) { return errors.New("end date must be after start date") }
	diff := end.Sub(start)
	if diff > 365*24*time.Hour { return errors.New("budget period cannot exceed 1 year") }
	return nil
}

func ValidateCategoryName(c string) error {
	if strings.TrimSpace(c) == "" { return errors.New("category name is required") }
	return nil
}

func ValidateBudgetAmount(amount int64) error {
	if amount < 0 { return errors.New("budget amount must be non-negative") }
	return nil
}
