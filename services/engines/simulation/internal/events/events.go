package events

import (
	"time"

	"github.com/horizon/core/services/engines/simulation/internal/engine"
)

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }

type BaseDomainEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseDomainEvent { return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string  { return b.ID }

type SimulationCreated struct{ BaseDomainEvent; SimType engine.SimType }
func NewSimulationCreated(id string, st engine.SimType) SimulationCreated {
	return SimulationCreated{BaseDomainEvent: NewBase("SimulationCreated", id), SimType: st}
}

type SimulationCompletedEvt struct{ BaseDomainEvent; SimType engine.SimType }
func NewSimulationCompletedEvt(id string, st engine.SimType) SimulationCompletedEvt {
	return SimulationCompletedEvt{BaseDomainEvent: NewBase("SimulationCompleted", id), SimType: st}
}
