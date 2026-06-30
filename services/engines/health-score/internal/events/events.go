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

type HealthScoreCalculated struct{ BaseDomainEvent; OverallScore int; ScoreGrade string }
func NewHealthScoreCalculated(id string, score int, grade string) HealthScoreCalculated {
	return HealthScoreCalculated{BaseDomainEvent: NewBase("HealthScoreCalculated", id), OverallScore: score, ScoreGrade: grade}
}

type HealthScoreDeclined struct{ BaseDomainEvent; OldScore int; NewScore int }
func NewHealthScoreDeclined(id string, old, new int) HealthScoreDeclined {
	return HealthScoreDeclined{BaseDomainEvent: NewBase("HealthScoreDeclined", id), OldScore: old, NewScore: new}
}
