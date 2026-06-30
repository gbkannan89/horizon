package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type NotificationViewed struct{ BaseEvent }
func NewNotificationViewed(id string) NotificationViewed { return NotificationViewed{BaseEvent: NewBase("NotificationViewed", id)} }
type NotificationDismissed struct{ BaseEvent }
func NewNotificationDismissed(id string) NotificationDismissed { return NotificationDismissed{BaseEvent: NewBase("NotificationDismissed", id)} }
