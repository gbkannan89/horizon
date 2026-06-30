package domain

import (
	"errors"
	"time"
)

type AllocationFactory struct{}

func NewAllocationFactory() *AllocationFactory { return &AllocationFactory{} }

func (f *AllocationFactory) Create(
	id, goalID, fundingSourceID string,
	fst FundingSourceType,
	allocType AllocationType,
	strategy AllocationStrategy,
	priority int,
	currency string,
	effectiveDate time.Time,
	weight *float64,
	fixedAmount *int64,
	expirationDate *time.Time,
	sourceOfTruth, createdBy string,
	tags []string,
	notes string,
) (*Allocation, error) {

	if err := ValidatePriority(priority); err != nil { return nil, err }
	if err := ValidateEffectiveDate(effectiveDate); err != nil { return nil, err }
	if allocType == ATFixedAmount && (fixedAmount == nil || *fixedAmount <= 0) {
		return nil, errors.New("fixed amount must be > 0 for FixedAmount type")
	}
	if allocType == ATPercentage && (weight == nil || *weight < 0 || *weight > 100) {
		return nil, errors.New("percentage must be 0-100 for Percentage type")
	}
	if expirationDate != nil {
		if err := ValidateExpirationDate(effectiveDate, *expirationDate); err != nil { return nil, err }
	}
	if createdBy == "Recommendation" {
		return nil, errors.New("recommendation allocations require approval")
	}

	now := time.Now().UTC()
	return &Allocation{
		id: id, goalID: goalID, fundingSourceID: fundingSourceID,
		fundingSourceType: fst, allocationType: allocType, strategy: strategy,
		health: AHHealthy, confidence: ACConfirmed, fundingCommitment: 0,
		priority: priority, weight: weight, fixedAmount: fixedAmount,
		currency: currency, status: StDraft, effectiveDate: effectiveDate,
		expirationDate: expirationDate, reservedAmount: 0, allocatedAmount: 0,
		sourceOfTruth: sourceOfTruth, createdBy: createdBy, tags: tags, notes: notes,
		createdAt: now, updatedAt: now,
	}, nil
}
