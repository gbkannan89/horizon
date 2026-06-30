package domain

import (
	"errors"
	"strings"
	"time"
)

// ValidateAssetName checks that the asset name is between 1 and 200 characters.
func ValidateAssetName(n string) error {
	t := strings.TrimSpace(n)
	if t == "" {
		return errors.New("asset name is required")
	}
	if len([]rune(t)) > 200 {
		return errors.New("asset name must be 200 characters or fewer")
	}
	return nil
}

// ValidateClassification checks that the classification is valid.
func ValidateClassification(c string) error {
	switch AssetClassification(c) {
	case ClsCash, ClsCashEquivalent, ClsEquity, ClsMutualFund, ClsETF, ClsBond,
		ClsRealEstate, ClsGold, ClsSilver, ClsCommodity, ClsVehicle,
		ClsBusinessOwnership, ClsRetirementAsset, ClsInsuranceCashValue,
		ClsCrypto, ClsCollectible, ClsDigitalAsset, ClsReceivable, ClsOther:
		return nil
	}
	return errors.New("asset classification is not recognized")
}

// ValidateOwnershipPercentage checks that the percentage is between 0 and 100.
func ValidateOwnershipPercentage(p float64) error {
	if p < 0 || p > 100 {
		return errors.New("ownership percentage must be between 0 and 100")
	}
	return nil
}

// ValidateQuantity checks that quantity is positive when set.
func ValidateQuantity(q *float64) error {
	if q != nil && *q <= 0 {
		return errors.New("quantity must be positive")
	}
	return nil
}

// ValidateCurrency checks that the currency code is valid.
func ValidateCurrency(c string) error {
	if len(c) != 3 {
		return errors.New("currency is not recognized")
	}
	return nil
}

// ValidateAcquisitionDate checks that acquisition date is in the past.
func ValidateAcquisitionDate(d *time.Time) error {
	if d != nil && d.After(time.Now()) {
		return errors.New("acquisition date cannot be in the future")
	}
	return nil
}
