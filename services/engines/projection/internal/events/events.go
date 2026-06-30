package events

import (
	"time"
)

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }
func (b BaseEvent) EntityID() string { return b.ID }

type ProjectionGenerated struct{ BaseEvent; ProjType string; Version int }
func NewProjectionGenerated(id, projType string, v int) ProjectionGenerated { return ProjectionGenerated{BaseEvent: NewBase("ProjectionGenerated", id), ProjType: projType, Version: v} }
type ProjectionExpired struct{ BaseEvent }
func NewProjectionExpired(id string) ProjectionExpired { return ProjectionExpired{BaseEvent: NewBase("ProjectionExpired", id)} }
type ProjectionFailed struct{ BaseEvent; Reason string }
func NewProjectionFailed(id, reason string) ProjectionFailed { return ProjectionFailed{BaseEvent: NewBase("ProjectionFailed", id), Reason: reason} }
