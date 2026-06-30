package events

import (
	"time"
)

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }

type BaseDomainEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseDomainEvent { return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string  { return b.ID }

type OptimizationGenerated struct{ BaseDomainEvent; CandidatesEvaluated int; CandidatesValid int }
func NewOptimizationGenerated(id string, evaluated, valid int) OptimizationGenerated {
	return OptimizationGenerated{BaseDomainEvent: NewBase("OptimizationGenerated", id), CandidatesEvaluated: evaluated, CandidatesValid: valid}
}

type OptimizationCompleted struct{ BaseDomainEvent; SelectedScore float64 }
func NewOptimizationCompleted(id string, score float64) OptimizationCompleted {
	return OptimizationCompleted{BaseDomainEvent: NewBase("OptimizationCompleted", id), SelectedScore: score}
}
