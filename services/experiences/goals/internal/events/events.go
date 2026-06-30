package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type GoalOpened struct{ BaseEvent; GoalID string }
func NewGoalOpened(id, gid string) GoalOpened { return GoalOpened{BaseEvent: NewBase("GoalOpened", id), GoalID: gid} }
