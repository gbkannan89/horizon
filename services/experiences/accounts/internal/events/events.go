package events

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type AccountListViewShown struct{ BaseEvent }
func NewAccountListViewShown(id string) AccountListViewShown { return AccountListViewShown{BaseEvent: NewBase("AccountListViewShown", id)} }
type AccountDetailViewShown struct{ BaseEvent; AccountID string }
func NewAccountDetailViewShown(id, aid string) AccountDetailViewShown { return AccountDetailViewShown{BaseEvent: NewBase("AccountDetailViewShown", id), AccountID: aid} }
type AccountCreatedEvt struct{ BaseEvent; AccountType string }
func NewAccountCreatedEvt(id, at string) AccountCreatedEvt { return AccountCreatedEvt{BaseEvent: NewBase("AccountCreated", id), AccountType: at} }
type AccountSearchPerformed struct{ BaseEvent; Query string; ResultCount int }
func NewAccountSearchPerformed(id, q string, rc int) AccountSearchPerformed { return AccountSearchPerformed{BaseEvent: NewBase("AccountSearchPerformed", id), Query: q, ResultCount: rc} }
