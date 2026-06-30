package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type PlanCreated struct{ BaseEvent }
func NewPlanCreated(id string) PlanCreated { return PlanCreated{BaseEvent: NewBase("PlanCreated", id)} }
type PlanCompared struct{ BaseEvent }
func NewPlanCompared(id string) PlanCompared { return PlanCompared{BaseEvent: NewBase("PlanCompared", id)} }
