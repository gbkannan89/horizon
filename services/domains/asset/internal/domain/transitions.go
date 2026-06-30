package domain

import (
	"errors"
	"time"
)

// AcquireAsset transitions from Planned to Acquired.
func AcquireAsset(a *Asset, costBasis int64, acqDate time.Time) error {
	if a.status != StatusPlanned {
		return errors.New("only planned assets can be acquired")
	}
	a.status = StatusAcquired
	a.costBasis = costBasis
	a.acquisitionDate = &acqDate
	a.updatedAt = time.Now().UTC()
	return nil
}

// ActivateAsset transitions from Acquired to Active.
func ActivateAsset(a *Asset) error {
	if a.status != StatusAcquired {
		return errors.New("only acquired assets can be activated")
	}
	a.status = StatusActive
	a.updatedAt = time.Now().UTC()
	return nil
}

// PartiallyDispose records a partial sale.
func PartiallyDispose(a *Asset, qtySold float64, proceeds int64, remainingQty float64) error {
	if a.status != StatusActive {
		return errors.New("only active assets can be partially disposed")
	}
	if a.quantity != nil && qtySold > *a.quantity {
		return errors.New("cannot sell more than owned quantity")
	}
	if a.quantity != nil {
		ratio := 1.0 - (qtySold / *a.quantity)
		a.costBasis = int64(float64(a.costBasis) * ratio)
	}
	a.status = StatusPartiallyDisposed
	a.updatedAt = time.Now().UTC()
	return nil
}

// FullyDispose records full disposal.
func FullyDispose(a *Asset, proceeds int64, dispDate time.Time) error {
	if a.status != StatusActive && a.status != StatusPartiallyDisposed {
		return errors.New("only active or partially disposed assets can be fully disposed")
	}
	a.status = StatusFullyDisposed
	a.dispositionDate = &dispDate
	a.updatedAt = time.Now().UTC()
	return nil
}

// RevalueAsset updates the current valuation.
func RevalueAsset(a *Asset, newValue int64, method ValuationMethod) {
	a.currentValue = newValue
	a.valuationMethod = method
	a.valuationDate = time.Now().UTC()
	a.updatedAt = time.Now().UTC()
}

// SplitAsset adjusts quantity and cost basis per unit for a stock split.
func SplitAsset(a *Asset, oldQty, newQty float64) error {
	if newQty <= 0 {
		return errors.New("new quantity must be positive")
	}
	if a.quantity != nil {
		ratio := newQty / oldQty
		a.costBasis = int64(float64(a.costBasis) * (newQty / oldQty))
		a.quantity = &newQty
		if a.unitPrice > 0 {
			a.unitPrice = a.unitPrice / ratio
		}
	}
	a.updatedAt = time.Now().UTC()
	return nil
}

// ArchiveAsset transitions to Archived.
func ArchiveAsset(a *Asset) error {
	if a.status != StatusActive && a.status != StatusPartiallyDisposed {
		return errors.New("only active or partially disposed assets can be archived")
	}
	a.status = StatusArchived
	a.updatedAt = time.Now().UTC()
	return nil
}
