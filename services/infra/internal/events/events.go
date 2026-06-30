package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type ServiceStarted struct{ BaseEvent; Service string }
func NewServiceStarted(id, svc string) ServiceStarted { return ServiceStarted{BaseEvent: NewBase("ServiceStarted", id), Service: svc} }
type ServiceHealthy struct{ BaseEvent }
func NewServiceHealthy(id string) ServiceHealthy { return ServiceHealthy{BaseEvent: NewBase("ServiceHealthy", id)} }
type ServiceUnhealthy struct{ BaseEvent; Reason string }
func NewServiceUnhealthy(id, reason string) ServiceUnhealthy { return ServiceUnhealthy{BaseEvent: NewBase("ServiceUnhealthy", id), Reason: reason} }
