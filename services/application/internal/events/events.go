package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type RequestCompleted struct{ BaseEvent; Method string; Path string; Status int; DurationMs int64 }
func NewRequestCompleted(id, method, path string, status int, dur int64) RequestCompleted { return RequestCompleted{BaseEvent: NewBase("RequestCompleted", id), Method: method, Path: path, Status: status, DurationMs: dur} }
type AuthFailure struct{ BaseEvent; Reason string }
func NewAuthFailure(id, reason string) AuthFailure { return AuthFailure{BaseEvent: NewBase("AuthFailure", id), Reason: reason} }
