package domain

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseEvent struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) BaseEvent { return BaseEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseEvent) EventName() string { return b.Type }; func (b BaseEvent) EntityID() string { return b.ID }

type PortfolioCreated struct{ BaseEvent; PortfolioType PortfolioType }
func NewPortfolioCreated(p *Portfolio) PortfolioCreated { return PortfolioCreated{BaseEvent: NewBase("PortfolioCreated", p.ID()), PortfolioType: p.PortfolioType()} }
type PortfolioActivated struct{ BaseEvent }
func NewPortfolioActivated(p *Portfolio) PortfolioActivated { return PortfolioActivated{BaseEvent: NewBase("PortfolioActivated", p.ID())} }
type PortfolioRebalanced struct{ BaseEvent }
func NewPortfolioRebalanced(p *Portfolio) PortfolioRebalanced { return PortfolioRebalanced{BaseEvent: NewBase("PortfolioRebalanced", p.ID())} }
type PortfolioArchived struct{ BaseEvent }
func NewPortfolioArchived(p *Portfolio) PortfolioArchived { return PortfolioArchived{BaseEvent: NewBase("PortfolioArchived", p.ID())} }
type PortfolioHealthChanged struct{ BaseEvent; OldHealth, NewHealth PortfolioHealth }
func NewPortfolioHealthChanged(p *Portfolio, old PortfolioHealth) PortfolioHealthChanged { return PortfolioHealthChanged{BaseEvent: NewBase("PortfolioHealthChanged", p.ID()), OldHealth: old, NewHealth: p.Health()} }
type PortfolioValuationUpdated struct{ BaseEvent; OldValue, NewValue int64 }
func NewPortfolioValuationUpdated(p *Portfolio, old int64) PortfolioValuationUpdated { return PortfolioValuationUpdated{BaseEvent: NewBase("PortfolioValuationUpdated", p.ID()), OldValue: old, NewValue: p.CurrentValue()} }
type PortfolioDriftExceeded struct{ BaseEvent; DriftPercent float64 }
func NewPortfolioDriftExceeded(p *Portfolio, drift float64) PortfolioDriftExceeded { return PortfolioDriftExceeded{BaseEvent: NewBase("PortfolioDriftExceeded", p.ID()), DriftPercent: drift} }
