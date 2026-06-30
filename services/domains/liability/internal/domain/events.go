package domain

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvt struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvt { return BaseEvt{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvt) EventName() string { return b.Type }
func (b BaseEvt) EntityID() string  { return b.ID }

type LiabilityCreated struct{ BaseEvt }
func NewLiabilityCreated(l *Liability) LiabilityCreated { return LiabilityCreated{BaseEvt: NewBase("LiabilityCreated", l.ID())} }
type LiabilityActivated struct{ BaseEvt }
func NewLiabilityActivated(l *Liability) LiabilityActivated { return LiabilityActivated{BaseEvt: NewBase("LiabilityActivated", l.ID())} }
type PaymentApplied struct{ BaseEvt; Amount int64 }
func NewPaymentApplied(l *Liability, amount int64) PaymentApplied { return PaymentApplied{BaseEvt: NewBase("PaymentApplied", l.ID()), Amount: amount} }
type PrepaymentApplied struct{ BaseEvt; Amount int64 }
func NewPrepaymentApplied(l *Liability, amount int64) PrepaymentApplied { return PrepaymentApplied{BaseEvt: NewBase("PrepaymentApplied", l.ID()), Amount: amount} }
type GracePeriodStarted struct{ BaseEvt }
func NewGracePeriodStarted(l *Liability) GracePeriodStarted { return GracePeriodStarted{BaseEvt: NewBase("GracePeriodStarted", l.ID())} }
type LiabilityDelinquent struct{ BaseEvt }
func NewLiabilityDelinquent(l *Liability) LiabilityDelinquent { return LiabilityDelinquent{BaseEvt: NewBase("LiabilityDelinquent", l.ID())} }
type LiabilityRestructured struct{ BaseEvt }
func NewLiabilityRestructured(l *Liability) LiabilityRestructured { return LiabilityRestructured{BaseEvt: NewBase("LiabilityRestructured", l.ID())} }
type LiabilitySettled struct{ BaseEvt }
func NewLiabilitySettled(l *Liability) LiabilitySettled { return LiabilitySettled{BaseEvt: NewBase("LiabilitySettled", l.ID())} }
type LiabilityWrittenOff struct{ BaseEvt }
func NewLiabilityWrittenOff(l *Liability) LiabilityWrittenOff { return LiabilityWrittenOff{BaseEvt: NewBase("LiabilityWrittenOff", l.ID())} }
type LiabilityArchived struct{ BaseEvt }
func NewLiabilityArchived(l *Liability) LiabilityArchived { return LiabilityArchived{BaseEvt: NewBase("LiabilityArchived", l.ID())} }
type InterestRateChanged struct{ BaseEvt; OldRate, NewRate float64 }
func NewInterestRateChanged(l *Liability, old float64) InterestRateChanged { return InterestRateChanged{BaseEvt: NewBase("InterestRateChanged", l.ID()), OldRate: old, NewRate: l.InterestRate()} }
type LiabilityHealthChanged struct{ BaseEvt; OldHealth, NewHealth LiabilityHealth }
func NewLiabilityHealthChanged(l *Liability, old LiabilityHealth) LiabilityHealthChanged { return LiabilityHealthChanged{BaseEvt: NewBase("LiabilityHealthChanged", l.ID()), OldHealth: old, NewHealth: l.LiabilityHealth()} }
