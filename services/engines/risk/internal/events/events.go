package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type RiskAssessmentUpdated struct{ BaseEvent; CompositeScore int; RiskLevel string }
func NewRiskAssessmentUpdated(id string, score int, level string) RiskAssessmentUpdated { return RiskAssessmentUpdated{BaseEvent: NewBase("RiskAssessmentUpdated", id), CompositeScore: score, RiskLevel: level} }
type RiskThresholdViolated struct{ BaseEvent; Indicator string; Score int; Threshold int }
func NewRiskThresholdViolated(id, indicator string, score, threshold int) RiskThresholdViolated { return RiskThresholdViolated{BaseEvent: NewBase("RiskThresholdViolated", id), Indicator: indicator, Score: score, Threshold: threshold} }
