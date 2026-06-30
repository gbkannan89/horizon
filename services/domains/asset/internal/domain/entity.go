package domain

import (
	"fmt"
	"time"
)

// Asset represents an economic resource owned or controlled by the user.
type Asset struct {
	id                   string
	assetName            string
	classification       AssetClassification
	assetType            *string
	ownershipModel       OwnershipModel
	ownerID              string
	householdID          *string
	institutionID        *string
	accountID            *string
	currency             string
	acquisitionDate      *time.Time
	dispositionDate      *time.Time
	assetHealth          AssetHealth
	assetConfidence      AssetConfidence
	costBasis            int64
	currentValue         int64
	valuationProfile     ValuationProfile
	valuationMethod      ValuationMethod
	valuationDate        time.Time
	quantity             *float64
	unitPrice            float64
	liquidityProfile     LiquidityProfile
	ownershipPercentage  float64
	status               AssetStatus
	externalIdentifiers  map[string]string
	sourceOfTruth        string
	tags                 []string
	metadata             map[string]string
	notes                string
	createdAt            time.Time
	updatedAt            time.Time
}

func (a *Asset) ID() string                           { return a.id }
func (a *Asset) AssetName() string                    { return a.assetName }
func (a *Asset) Classification() AssetClassification   { return a.classification }
func (a *Asset) AssetType() *string                   { return a.assetType }
func (a *Asset) OwnershipModel() OwnershipModel        { return a.ownershipModel }
func (a *Asset) OwnerID() string                      { return a.ownerID }
func (a *Asset) Currency() string                     { return a.currency }
func (a *Asset) AcquisitionDate() *time.Time          { return a.acquisitionDate }
func (a *Asset) DispositionDate() *time.Time          { return a.dispositionDate }
func (a *Asset) AssetHealth() AssetHealth             { return a.assetHealth }
func (a *Asset) AssetConfidence() AssetConfidence     { return a.assetConfidence }
func (a *Asset) CostBasis() int64                     { return a.costBasis }
func (a *Asset) CurrentValue() int64                  { return a.currentValue }
func (a *Asset) ValuationProfile() ValuationProfile   { return a.valuationProfile }
func (a *Asset) ValuationMethod() ValuationMethod     { return a.valuationMethod }
func (a *Asset) ValuationDate() time.Time             { return a.valuationDate }
func (a *Asset) Quantity() *float64                   { return a.quantity }
func (a *Asset) UnitPrice() float64                   { return a.unitPrice }
func (a *Asset) LiquidityProfile() LiquidityProfile   { return a.liquidityProfile }
func (a *Asset) OwnershipPercentage() float64         { return a.ownershipPercentage }
func (a *Asset) Status() AssetStatus                  { return a.status }
func (a *Asset) SourceOfTruth() string                { return a.sourceOfTruth }
func (a *Asset) CreatedAt() time.Time                 { return a.createdAt }
func (a *Asset) UpdatedAt() time.Time                 { return a.updatedAt }

// SetStatus transitions the asset to a new lifecycle state.
func (a *Asset) SetStatus(s AssetStatus) { a.status = s; a.updatedAt = time.Now().UTC() }

// CanTransitionTo checks whether a state transition is valid.
func (a *Asset) CanTransitionTo(target AssetStatus) error {
	m := map[AssetStatus]map[AssetStatus]bool{
		StatusPlanned:         {StatusAcquired: true},
		StatusAcquired:        {StatusActive: true},
		StatusActive:          {StatusPartiallyDisposed: true, StatusFullyDisposed: true, StatusArchived: true},
		StatusPartiallyDisposed: {StatusActive: true, StatusFullyDisposed: true},
		StatusFullyDisposed:   {StatusHistorical: true},
		StatusArchived:        {StatusHistorical: true},
	}
	if al, ok := m[a.status]; ok && al[target] {
		return nil
	}
	if a.status == target {
		return fmt.Errorf("already in %s state", target)
	}
	if a.status == StatusArchived {
		return fmt.Errorf("archived assets cannot be reactivated")
	}
	return fmt.Errorf("cannot transition from %s to %s", a.status, target)
}

// SetCurrentValue updates the derived current value.
func (a *Asset) SetCurrentValue(v int64) { a.currentValue = v; a.updatedAt = time.Now().UTC() }

// SetCostBasis updates the cost basis.
func (a *Asset) SetCostBasis(b int64) { a.costBasis = b; a.updatedAt = time.Now().UTC() }

// SetQuantity updates the quantity.
func (a *Asset) SetQuantity(q float64) { a.quantity = &q; a.updatedAt = time.Now().UTC() }

// SetUnitPrice updates the unit price.
func (a *Asset) SetUnitPrice(p float64) { a.unitPrice = p; a.updatedAt = time.Now().UTC() }

func ReconstructFromDB(id, assetName string, cls AssetClassification, ownerID, currency string,
	costBasis int64, vp ValuationProfile, vm ValuationMethod, vd time.Time,
	qty *float64, up float64, lp LiquidityProfile, op float64, status AssetStatus,
	sot string, tags []string, createdAt time.Time) *Asset {
	return &Asset{
		id: id, assetName: assetName, classification: cls, ownerID: ownerID, currency: currency,
		costBasis: costBasis, valuationProfile: vp, valuationMethod: vm, valuationDate: vd,
		quantity: qty, unitPrice: up, liquidityProfile: lp, ownershipPercentage: op,
		status: status, sourceOfTruth: sot, tags: tags, createdAt: createdAt,
	}
}