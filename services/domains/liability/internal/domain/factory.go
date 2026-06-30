package domain

import (
	"time"
)

// LiabilityFactory creates new Liability entities.
type LiabilityFactory struct{}

// NewLiabilityFactory returns a new LiabilityFactory.
func NewLiabilityFactory() *LiabilityFactory { return &LiabilityFactory{} }

// Create creates a new liability in Planned state.
func (f *LiabilityFactory) Create(
	id, name string,
	classification LiabilityClassification,
	ownerID, currency string,
	originalPrincipal int64,
	interestRate float64,
	interestModel InterestModel,
	repaymentProfile RepaymentProfile,
	repaymentModel RepaymentModel,
	installmentAmount int64,
	installmentFreq string,
	remainingInstallments int,
	maturityDate time.Time,
	sourceOfTruth string,
	liabilityType, householdID, institutionID, servicingAccountID, collateral, extRef string,
	tags []string, notes string,
) (*Liability, error) {

	if err := ValidateLiabilityName(name); err != nil { return nil, err }
	if err := ValidateClassification(string(classification)); err != nil { return nil, err }
	if err := ValidateInterestRate(interestRate); err != nil { return nil, err }
	if err := ValidateOriginalPrincipal(originalPrincipal); err != nil { return nil, err }
	if err := ValidateMaturityDate(maturityDate); err != nil { return nil, err }
	if tags == nil { tags = []string{} }

	now := time.Now().UTC()
	return &Liability{
		id: id, name: name, classification: classification, ownerID: ownerID,
		currency: currency, originalPrincipal: originalPrincipal, outstandingBalance: originalPrincipal,
		interestRate: interestRate, interestModel: interestModel,
		repaymentProfile: repaymentProfile, repaymentModel: repaymentModel,
		installmentAmount: installmentAmount, installmentFrequency: installmentFreq,
		remainingInstallments: remainingInstallments, maturityDate: maturityDate,
		liabilityHealth: LHStable, liabilityConfidence: LCUserVerifiedSt,
		status: LSPlanned, sourceOfTruth: sourceOfTruth, tags: tags,
		metadata: map[string]string{}, notes: notes, externalReference: extRef,
		createdAt: now, updatedAt: now,
	}, nil
}
