package domain

import (
	"fmt"
	"time"
)

type Allocation struct {
	id                string
	goalID            string
	fundingSourceID   string
	fundingSourceType FundingSourceType
	allocationType    AllocationType
	strategy          AllocationStrategy
	health            AllocationHealth
	confidence        AllocationConfidence
	fundingCommitment int64
	priority          int
	weight            *float64
	fixedAmount       *int64
	currency          string
	status            AllocationStatus
	effectiveDate     time.Time
	expirationDate    *time.Time
	reservedAmount    int64
	allocatedAmount   int64
	sourceOfTruth     string
	createdBy         string
	approvedBy        *string
	tags              []string
	metadata          map[string]string
	notes             string
	createdAt         time.Time
	updatedAt         time.Time
}

func (a *Allocation) ID() string { return a.id }
func (a *Allocation) GoalID() string { return a.goalID }
func (a *Allocation) FundingSourceID() string { return a.fundingSourceID }
func (a *Allocation) FundingSourceType() FundingSourceType { return a.fundingSourceType }
func (a *Allocation) AllocationType() AllocationType { return a.allocationType }
func (a *Allocation) Strategy() AllocationStrategy { return a.strategy }
func (a *Allocation) Health() AllocationHealth { return a.health }
func (a *Allocation) Confidence() AllocationConfidence { return a.confidence }
func (a *Allocation) FundingCommitment() int64 { return a.fundingCommitment }
func (a *Allocation) Priority() int { return a.priority }
func (a *Allocation) Weight() *float64 { return a.weight }
func (a *Allocation) FixedAmount() *int64 { return a.fixedAmount }
func (a *Allocation) Currency() string { return a.currency }
func (a *Allocation) Status() AllocationStatus { return a.status }
func (a *Allocation) EffectiveDate() time.Time { return a.effectiveDate }
func (a *Allocation) ExpirationDate() *time.Time { return a.expirationDate }
func (a *Allocation) ReservedAmount() int64 { return a.reservedAmount }
func (a *Allocation) AllocatedAmount() int64 { return a.allocatedAmount }
func (a *Allocation) CreatedBy() string { return a.createdBy }
func (a *Allocation) Tags() []string { return a.tags }
func (a *Allocation) CreatedAt() time.Time { return a.createdAt }
func (a *Allocation) UpdatedAt() time.Time { return a.updatedAt }

func (a *Allocation) SetStatus(s AllocationStatus) { a.status = s; a.updatedAt = time.Now().UTC() }

func (a *Allocation) CanTransitionTo(target AllocationStatus) error {
	m := map[AllocationStatus]map[AllocationStatus]bool{
		StDraft: {StPlanned: true, StCancelled: true},
		StPlanned: {StActive: true, StCancelled: true},
		StActive: {StPaused: true, StCompleted: true},
		StPaused: {StActive: true, StCancelled: true},
		StCompleted: {StArchived: true},
		StCancelled: {StArchived: true},
	}
	if al, ok := m[a.status]; ok && al[target] { return nil }
	if a.status == target { return fmt.Errorf("already in %s", target) }
	return fmt.Errorf("cannot transition from %s to %s", a.status, target)
}

func ReconstructFromDB(id, goalID, fundingSourceID string, allocType AllocationType, priority int,
	currency string, status AllocationStatus, reservedAmount, allocatedAmount int64,
	tags []string, createdAt, updatedAt time.Time) *Allocation {
	return &Allocation{
		id: id, goalID: goalID, fundingSourceID: fundingSourceID, allocationType: allocType,
		priority: priority, currency: currency, status: status,
		reservedAmount: reservedAmount, allocatedAmount: allocatedAmount,
		tags: tags, createdAt: createdAt, updatedAt: updatedAt,
	}
}