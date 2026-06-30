package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type InsightGeneratedEvt struct{ BaseEvent; Category string }
func NewInsightGenerated(id, cat string) InsightGeneratedEvt { return InsightGeneratedEvt{BaseEvent: NewBase("InsightGenerated", id), Category: cat} }
type InsightDismissed struct{ BaseEvent; Reason string }
func NewInsightDismissed(id, reason string) InsightDismissed { return InsightDismissed{BaseEvent: NewBase("InsightDismissed", id), Reason: reason} }
