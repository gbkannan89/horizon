package events

import (
	"time"

	"github.com/horizon/core/services/engines/risk/internal/engine"
)

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }

type BaseDomainEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseDomainEvent { return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string  { return b.ID }

type RiskAssessmentGenerated struct{ BaseDomainEvent; CompositeScore int; RiskLevel engine.RiskLevel }
func NewRiskAssessmentGenerated(id string, score int, level engine.RiskLevel) RiskAssessmentGenerated {
	return RiskAssessmentGenerated{BaseDomainEvent: NewBase("RiskAssessmentGenerated", id), CompositeScore: score, RiskLevel: level}
}

type RiskThresholdExceeded struct{ BaseDomainEvent; Indicator string; Score int; Threshold int }
func NewRiskThresholdExceeded(id, indicator string, score, threshold int) RiskThresholdExceeded {
	return RiskThresholdExceeded{BaseDomainEvent: NewBase("RiskThresholdExceeded", id), Indicator: indicator, Score: score, Threshold: threshold}
}
