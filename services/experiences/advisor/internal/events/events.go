package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type AdvisorOpened struct{ BaseEvent; Mode string }
func NewAdvisorOpened(id, mode string) AdvisorOpened { return AdvisorOpened{BaseEvent: NewBase("AdvisorOpened", id), Mode: mode} }
type AdvisorDecisionMade struct{ BaseEvent; Decision string }
func NewAdvisorDecisionMade(id, decision string) AdvisorDecisionMade { return AdvisorDecisionMade{BaseEvent: NewBase("AdvisorDecisionMade", id), Decision: decision} }
