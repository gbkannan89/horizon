package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type RecommendationsGenerated struct{ BaseEvent; Count int }
func NewRecommendationsGenerated(id string, c int) RecommendationsGenerated { return RecommendationsGenerated{BaseEvent: NewBase("RecommendationsGenerated", id), Count: c} }
type RecommendationAccepted struct{ BaseEvent }
func NewRecommendationAccepted(id string) RecommendationAccepted { return RecommendationAccepted{BaseEvent: NewBase("RecommendationAccepted", id)} }
type RecommendationRejected struct{ BaseEvent }
func NewRecommendationRejected(id string) RecommendationRejected { return RecommendationRejected{BaseEvent: NewBase("RecommendationRejected", id)} }
