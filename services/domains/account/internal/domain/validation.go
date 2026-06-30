package domain

import (
	"errors"
	"strings"
	"time"
)

func ValidateAccountName(n string) error {
	t := strings.TrimSpace(n)
	if t == "" { return errors.New("account name is required") }
	if len([]rune(t)) > 200 { return errors.New("account name must be 200 characters or fewer") }
	return nil
}

func ValidateAccountType(t string) error {
	if _, ok := typeClassifications[AccountType(t)]; !ok {
		return errors.New("account type is not recognized")
	}
	return nil
}

func ValidateCurrency(c string) error {
	if len(c) != 3 { return errors.New("currency is not recognized") }
	return nil
}

func ValidateClassificationConsistency(at AccountType, cls AccountClassification) error {
	allowed, ok := typeClassifications[at]
	if !ok { return errors.New("unknown account type") }
	for _, a := range allowed {
		if a == cls { return nil }
	}
	return errors.New("classification is inconsistent with account type")
}

func ValidateCreditLimit(at AccountType, limit *int64) error {
	if limit != nil {
		switch at {
		case ATCreditCard, ATLoan, ATMortgage, ATOverdraft:
		default:
			return errors.New("credit limit is only valid for credit-type accounts")
		}
		if *limit < 0 { return errors.New("credit limit must be non-negative") }
	}
	return nil
}

func ValidateInterestRate(rate *float64) error {
	if rate != nil && (*rate < 0 || *rate > 100) {
		return errors.New("interest rate must be between 0 and 100")
	}
	return nil
}

func ValidateClosedDate(status AccountStatus, closedDate *time.Time) error {
	if status == StatusClosed && closedDate == nil {
		return errors.New("closed date is required when status is Closed")
	}
	if status != StatusClosed && closedDate != nil {
		return errors.New("closed date must not be set when account is not Closed")
	}
	return nil
}
