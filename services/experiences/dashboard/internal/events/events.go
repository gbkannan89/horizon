package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type DashboardRefreshed struct{ BaseEvent; State string }
func NewDashboardRefreshed(id, state string) DashboardRefreshed { return DashboardRefreshed{BaseEvent: NewBase("DashboardRefreshed", id), State: state} }
type DashboardAlert struct{ BaseEvent; Severity string; Widget string; Message string }
func NewDashboardAlert(id, severity, widget, msg string) DashboardAlert { return DashboardAlert{BaseEvent: NewBase("DashboardAlert", id), Severity: severity, Widget: widget, Message: msg} }
