package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func ValidateDisplayName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.New("display name is required")
	}
	if len([]rune(trimmed)) > 100 {
		return errors.New("display name must be 100 characters or fewer")
	}
	return nil
}

func ValidateLegalName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("legal name is required")
	}
	return nil
}

func ValidateDateOfBirth(dob *time.Time) error {
	if dob == nil {
		return nil
	}
	age := time.Since(*dob).Hours() / (365.25 * 24)
	if age < 0 || age > 150 {
		return errors.New("date of birth is invalid")
	}
	return nil
}

var iso3166Map = map[string]bool{
	"IN": true, "US": true, "GB": true, "DE": true, "FR": true, "JP": true,
	"AU": true, "CA": true, "SG": true, "AE": true, "NZ": true, "ZA": true,
}

func ValidateCountry(code string) error {
	code = strings.ToUpper(code)
	if !iso3166Map[code] {
		return errors.New("country is not recognized")
	}
	return nil
}

var iso4217Map = map[string]bool{
	"INR": true, "USD": true, "EUR": true, "GBP": true, "JPY": true,
	"CHF": true, "AUD": true, "CAD": true, "SGD": true, "HKD": true, "CNY": true, "NZD": true,
}

func ValidateBaseCurrency(code string) error {
	code = strings.ToUpper(code)
	if !iso4217Map[code] {
		return errors.New("base currency is not recognized")
	}
	return nil
}

var localeRegexp = regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)

func ValidateLocale(locale string) error {
	if !localeRegexp.MatchString(locale) {
		return errors.New("locale is not recognized")
	}
	return nil
}

var ianaTimezones = map[string]bool{
	"Asia/Kolkata": true, "America/New_York": true, "America/Chicago": true,
	"America/Denver": true, "America/Los_Angeles": true, "Europe/London": true,
	"Europe/Paris": true, "Europe/Berlin": true, "Asia/Tokyo": true,
	"Asia/Singapore": true, "Australia/Sydney": true, "Pacific/Auckland": true,
	"UTC": true,
}

func ValidateTimezone(tz string) error {
	if !ianaTimezones[tz] {
		return errors.New("timezone is not recognized")
	}
	return nil
}

func ValidateMinorUser(userType UserType, dob *time.Time) error {
	if userType != TypeMinor || dob == nil {
		return nil
	}
	age := time.Since(*dob).Hours() / (365.25 * 24)
	if age >= 18 {
		return errors.New("minor user type requires age below majority")
	}
	return nil
}

func ValidateSOT(sot SOT) error {
	switch sot {
	case SOTSystem, SOTManual, SOTImport, SOTMigration:
		return nil
	}
	return fmt.Errorf("invalid source of truth: %s", sot)
}
