package domain

import "errors"

func ValidatePortfolioName(n string) error {
	if len(n) == 0 { return errors.New("portfolio name is required") }
	if len(n) > 200 { return errors.New("portfolio name must be 200 characters or fewer") }
	return nil
}

func ValidateCurrency(c string) error {
	if len(c) != 3 { return errors.New("currency is not recognized") }
	return nil
}
