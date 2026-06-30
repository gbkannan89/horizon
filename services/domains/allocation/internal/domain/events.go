package domain

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseDomainEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseDomainEvent { return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string { return b.ID }

type AllocationCreated struct{ BaseDomainEvent }
func NewAllocationCreated(a *Allocation) AllocationCreated { return AllocationCreated{BaseDomainEvent: NewBase("AllocationCreated", a.ID())} }
type AllocationApproved struct{ BaseDomainEvent }
func NewAllocationApproved(a *Allocation) AllocationApproved { return AllocationApproved{BaseDomainEvent: NewBase("AllocationApproved", a.ID())} }
type AllocationActivated struct{ BaseDomainEvent }
func NewAllocationActivated(a *Allocation) AllocationActivated { return AllocationActivated{BaseDomainEvent: NewBase("AllocationActivated", a.ID())} }
type AllocationPaused struct{ BaseDomainEvent }
func NewAllocationPaused(a *Allocation) AllocationPaused { return AllocationPaused{BaseDomainEvent: NewBase("AllocationPaused", a.ID())} }
type AllocationResumed struct{ BaseDomainEvent }
func NewAllocationResumed(a *Allocation) AllocationResumed { return AllocationResumed{BaseDomainEvent: NewBase("AllocationResumed", a.ID())} }
type AllocationCompleted struct{ BaseDomainEvent }
func NewAllocationCompleted(a *Allocation) AllocationCompleted { return AllocationCompleted{BaseDomainEvent: NewBase("AllocationCompleted", a.ID())} }
type AllocationCancelled struct{ BaseDomainEvent }
func NewAllocationCancelled(a *Allocation) AllocationCancelled { return AllocationCancelled{BaseDomainEvent: NewBase("AllocationCancelled", a.ID())} }
type AllocationArchived struct{ BaseDomainEvent }
func NewAllocationArchived(a *Allocation) AllocationArchived { return AllocationArchived{BaseDomainEvent: NewBase("AllocationArchived", a.ID())} }
type AllocationRecalculated struct{ BaseDomainEvent; ReservedAmount int64 }
func NewAllocationRecalculated(a *Allocation) AllocationRecalculated { return AllocationRecalculated{BaseDomainEvent: NewBase("AllocationRecalculated", a.ID()), ReservedAmount: a.ReservedAmount()} }
type AllocationConflictDetected struct{ BaseDomainEvent; ConflictingAllocationID string }
func NewAllocationConflictDetected(a *Allocation, conflictID string) AllocationConflictDetected { return AllocationConflictDetected{BaseDomainEvent: NewBase("AllocationConflictDetected", a.ID()), ConflictingAllocationID: conflictID} }
type AllocationUnderfunded struct{ BaseDomainEvent; Shortfall int64 }
func NewAllocationUnderfunded(a *Allocation, shortfall int64) AllocationUnderfunded { return AllocationUnderfunded{BaseDomainEvent: NewBase("AllocationUnderfunded", a.ID()), Shortfall: shortfall} }
type AllocationHealthChanged struct{ BaseDomainEvent; OldHealth, NewHealth AllocationHealth }
func NewAllocationHealthChanged(a *Allocation, old AllocationHealth) AllocationHealthChanged { return AllocationHealthChanged{BaseDomainEvent: NewBase("AllocationHealthChanged", a.ID()), OldHealth: old, NewHealth: a.Health()} }
