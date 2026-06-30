package events

import (
	"time"

	"github.com/horizon/core/services/engines/recommendation/internal/engine"
)

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }

type BaseDomainEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseDomainEvent { return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string  { return b.ID }

type RecommendationsGenerated struct{ BaseDomainEvent; Count int }
func NewRecommendationsGenerated(id string, count int) RecommendationsGenerated {
	return RecommendationsGenerated{BaseDomainEvent: NewBase("RecommendationsGenerated", id), Count: count}
}

type RecommendationAcceptedEvt struct{ BaseDomainEvent; Category engine.RecCategory }
func NewRecommendationAcceptedEvt(id string, cat engine.RecCategory) RecommendationAcceptedEvt {
	return RecommendationAcceptedEvt{BaseDomainEvent: NewBase("RecommendationAccepted", id), Category: cat}
}
