package domain

import (
	"fmt"
	"time"
)

// Liability represents a financial obligation owed by the user.
type Liability struct {
	id                    string
	name                  string
	classification        LiabilityClassification
	liabilityType         string
	ownerID               string
	householdID           *string
	institutionID         *string
	servicingAccountID    *string
	currency              string
	originalPrincipal     int64
	outstandingBalance    int64
	liabilityHealth       LiabilityHealth
	liabilityConfidence   LiabilityConfidence
	interestRate          float64
	interestModel         InterestModel
	repaymentProfile      RepaymentProfile
	repaymentModel        RepaymentModel
	installmentAmount     int64
	installmentFrequency  string
	remainingInstallments int
	maturityDate          time.Time
	collateral            string
	status                LiabilityStatus
	sourceOfTruth         string
	tags                  []string
	metadata              map[string]string
	notes                 string
	externalReference     string
	createdAt             time.Time
	updatedAt             time.Time
}

func (l *Liability) ID() string                        { return l.id }
func (l *Liability) Name() string                      { return l.name }
func (l *Liability) Classification() LiabilityClassification { return l.classification }
func (l *Liability) OwnerID() string                   { return l.ownerID }
func (l *Liability) Currency() string                  { return l.currency }
func (l *Liability) OriginalPrincipal() int64          { return l.originalPrincipal }
func (l *Liability) OutstandingBalance() int64         { return l.outstandingBalance }
func (l *Liability) InterestRate() float64             { return l.interestRate }
func (l *Liability) InterestModel() InterestModel      { return l.interestModel }
func (l *Liability) RepaymentProfile() RepaymentProfile { return l.repaymentProfile }
func (l *Liability) RepaymentModel() RepaymentModel    { return l.repaymentModel }
func (l *Liability) InstallmentAmount() int64          { return l.installmentAmount }
func (l *Liability) RemainingInstallments() int        { return l.remainingInstallments }
func (l *Liability) MaturityDate() time.Time           { return l.maturityDate }
func (l *Liability) Status() LiabilityStatus           { return l.status }
func (l *Liability) LiabilityHealth() LiabilityHealth  { return l.liabilityHealth }
func (l *Liability) LiabilityConfidence() LiabilityConfidence { return l.liabilityConfidence }
func (l *Liability) CreatedAt() time.Time              { return l.createdAt }
func (l *Liability) UpdatedAt() time.Time              { return l.updatedAt }

// SetStatus transitions the liability to a new state.
func (l *Liability) SetStatus(s LiabilityStatus) { l.status = s; l.updatedAt = time.Now().UTC() }

// SetOutstandingBalance updates the derived balance.
func (l *Liability) SetOutstandingBalance(b int64) { l.outstandingBalance = b; l.updatedAt = time.Now().UTC() }

// SetRemainingInstallments updates the count.
func (l *Liability) SetRemainingInstallments(n int) { l.remainingInstallments = n; l.updatedAt = time.Now().UTC() }

// CanTransitionTo checks validity of a state change.
func (l *Liability) CanTransitionTo(target LiabilityStatus) error {
	m := map[LiabilityStatus]map[LiabilityStatus]bool{
		LSPlanned:    {LSApproved: true},
		LSApproved:   {LSActive: true},
		LSActive:     {LSGracePeriod: true, LSDelinquent: true, LSRestructured: true, LSSettled: true},
		LSGracePeriod: {LSActive: true, LSDelinquent: true},
		LSDelinquent:  {LSRestructured: true, LSWrittenOff: true},
		LSRestructured: {LSActive: true},
		LSSettled:     {LSHistorical: true, LSArchived: true},
		LSWrittenOff:  {LSArchived: true},
		LSHistorical:  {LSArchived: true},
	}
	if al, ok := m[l.status]; ok && al[target] {
		return nil
	}
	if l.status == target {
		return fmt.Errorf("already in %s state", target)
	}
	if l.status == LSArchived {
		return fmt.Errorf("archived liabilities cannot transition")
	}
	return fmt.Errorf("cannot transition from %s to %s", l.status, target)
}

func ReconstructFromDB(id, name string, cls LiabilityClassification, ownerID, currency string,
	originalPrincipal, interestRate float64, interestModel InterestModel, repaymentModel RepaymentModel,
	remainingMonths int, maturityDate time.Time, lh LiabilityHealth, lc LiabilityConfidence,
	status LiabilityStatus, sot string, tags []string, createdAt time.Time) *Liability {
	return &Liability{
		id: id, name: name, classification: cls, ownerID: ownerID, currency: currency,
		originalPrincipal: int64(originalPrincipal), interestRate: interestRate,
		interestModel: interestModel, repaymentModel: repaymentModel,
		remainingInstallments: remainingMonths, maturityDate: maturityDate,
		liabilityHealth: lh, liabilityConfidence: lc, status: status, sourceOfTruth: sot,
		tags: tags, createdAt: createdAt,
	}
}