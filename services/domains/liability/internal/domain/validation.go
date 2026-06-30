package domain

import (
	"errors"
	"strings"
	"time"
)

func ValidateLiabilityName(n string) error {
	t := strings.TrimSpace(n)
	if t == "" { return errors.New("liability name is required") }
	if len([]rune(t)) > 200 { return errors.New("liability name must be 200 characters or fewer") }
	return nil
}

func ValidateClassification(c string) error {
	switch LiabilityClassification(c) {
	case LCMortgage, LCHomeLoan, LCPersonalLoan, LCVehicleLoan, LCEducationLoan,
		LCCreditCard, LCLineOfCredit, LCBusinessLoan, LCGoldLoan, LCMarginLoan,
		LCTaxLiability, LCFamilyLoan, LCEmployerLoan, LCBuyNowPayLater, LCOverdraft, LCOther:
		return nil
	}
	return errors.New("liability classification is not recognized")
}

func ValidateInterestRate(r float64) error {
	if r < 0 { return errors.New("interest rate must be non-negative") }
	return nil
}

func ValidateOriginalPrincipal(p int64) error {
	if p <= 0 { return errors.New("original principal must be positive") }
	return nil
}

func ValidateMaturityDate(md time.Time) error {
	if md.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("maturity date must be in the future")
	}
	return nil
}
