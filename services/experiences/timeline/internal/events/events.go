package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type TimelineViewed struct{ BaseEvent; View string }
func NewTimelineViewed(id, view string) TimelineViewed { return TimelineViewed{BaseEvent: NewBase("TimelineViewed", id), View: view} }
type TimelineEventSelected struct{ BaseEvent; EventType string }
func NewTimelineEventSelected(id, et string) TimelineEventSelected { return TimelineEventSelected{BaseEvent: NewBase("TimelineEventSelected", id), EventType: et} }
