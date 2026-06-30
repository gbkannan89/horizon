package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func ValidateEventType(s string) error {
	if !ValidEventType(s) {
		return fmt.Errorf("event type is not recognized: %s", s)
	}
	return nil
}

func ValidateCurrency(code string) error {
	if len(code) != 3 {
		return fmt.Errorf("currency is not recognized: %s", code)
	}
	code = strings.ToUpper(code)
	supported := map[string]bool{"INR": true, "USD": true, "EUR": true, "GBP": true,
		"JPY": true, "CHF": true, "AUD": true, "CAD": true, "SGD": true, "HKD": true, "CNY": true, "NZD": true}
	if !supported[code] {
		return fmt.Errorf("currency is not recognized: %s", code)
	}
	return nil
}

func ValidateEventDate(ed time.Time) error {
	if ed.IsZero() {
		return errors.New("event date is required")
	}
	return nil
}

func ValidateEffectiveDate(eff, eventDate time.Time, maxDaysBefore int) error {
	if eff.IsZero() {
		return errors.New("effective date is required")
	}
	if eff.After(eventDate) {
		return errors.New("effective date may not be in the future")
	}
	maxPast := eventDate.AddDate(0, 0, -maxDaysBefore)
	if eff.Before(maxPast) {
		return fmt.Errorf("effective date may not be more than %d days before the event date", maxDaysBefore)
	}
	return nil
}

func ValidateAmount(amount int64) error {
	if amount == 0 {
		return errors.New("event amount must be non-zero")
	}
	return nil
}
