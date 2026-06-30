package events

import (
	"time"

	"github.com/horizon/core/services/engines/projection/internal/engine"
)

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }

type BaseDomainEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseDomainEvent { return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string  { return b.ID }

type ProjectionGenerated struct{ BaseDomainEvent; ProjectionType engine.ProjectionType; Version int }
func NewProjectionGenerated(id string, pt engine.ProjectionType, ver int) ProjectionGenerated {
	return ProjectionGenerated{BaseDomainEvent: NewBase("ProjectionGenerated", id), ProjectionType: pt, Version: ver}
}

type ProjectionUpdated struct{ BaseDomainEvent; ProjectionType engine.ProjectionType; Version int }
func NewProjectionUpdated(id string, pt engine.ProjectionType, ver int) ProjectionUpdated {
	return ProjectionUpdated{BaseDomainEvent: NewBase("ProjectionUpdated", id), ProjectionType: pt, Version: ver}
}

type ProjectionExpired struct{ BaseDomainEvent; ProjectionType engine.ProjectionType }
func NewProjectionExpired(id string, pt engine.ProjectionType) ProjectionExpired {
	return ProjectionExpired{BaseDomainEvent: NewBase("ProjectionExpired", id), ProjectionType: pt}
}
