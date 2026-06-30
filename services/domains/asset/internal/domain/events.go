package domain

import "time"

// DomainEvent is the interface for domain events.
type DomainEvent interface{ EventName() string; EntityID() string }

// Publisher publishes domain events.
type Publisher interface{ Publish(event DomainEvent) error }

// Base provides common event fields.
type Base struct{ Type, ID string; Time time.Time }

// NewBase creates a new Base event.
func NewBase(evt, eid string) Base { return Base{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b Base) EventName() string { return b.Type }
func (b Base) EntityID() string  { return b.ID }

type AssetCreated struct{ Base; Classification AssetClassification }
func NewAssetCreated(a *Asset) AssetCreated { return AssetCreated{Base: NewBase("AssetCreated", a.ID()), Classification: a.Classification()} }

type AssetActivated struct{ Base }
func NewAssetActivated(a *Asset) AssetActivated { return AssetActivated{Base: NewBase("AssetActivated", a.ID())} }

type AssetPartiallyDisposed struct{ Base; QtySold float64; Proceeds int64 }
func NewAssetPartiallyDisposed(a *Asset, qty float64, proceeds int64) AssetPartiallyDisposed { return AssetPartiallyDisposed{Base: NewBase("AssetPartiallyDisposed", a.ID()), QtySold: qty, Proceeds: proceeds} }

type AssetFullyDisposed struct{ Base; Proceeds int64; RealisedGainLoss int64 }
func NewAssetFullyDisposed(a *Asset, proceeds int64, gainLoss int64) AssetFullyDisposed { return AssetFullyDisposed{Base: NewBase("AssetFullyDisposed", a.ID()), Proceeds: proceeds, RealisedGainLoss: gainLoss} }

type AssetRevalued struct{ Base; OldValue, NewValue int64; Method ValuationMethod }
func NewAssetRevalued(a *Asset, oldVal int64) AssetRevalued { return AssetRevalued{Base: NewBase("AssetRevalued", a.ID()), OldValue: oldVal, NewValue: a.CurrentValue(), Method: a.ValuationMethod()} }

type AssetSplit struct{ Base; OldQty, NewQty float64; AdjustedCostBasis int64 }
func NewAssetSplit(a *Asset, oldQty, newQty float64) AssetSplit { return AssetSplit{Base: NewBase("AssetSplit", a.ID()), OldQty: oldQty, NewQty: newQty, AdjustedCostBasis: a.CostBasis()} }

type AssetArchived struct{ Base }
func NewAssetArchived(a *Asset) AssetArchived { return AssetArchived{Base: NewBase("AssetArchived", a.ID())} }

type AssetTransferred struct{ Base; FromAccount, ToAccount string }
func NewAssetTransferred(a *Asset, from, to string) AssetTransferred { return AssetTransferred{Base: NewBase("AssetTransferred", a.ID()), FromAccount: from, ToAccount: to} }

type CostBasisAdjusted struct{ Base; OldCostBasis, NewCostBasis int64; Reason string }
func NewCostBasisAdjusted(a *Asset, oldBasis int64, reason string) CostBasisAdjusted { return CostBasisAdjusted{Base: NewBase("CostBasisAdjusted", a.ID()), OldCostBasis: oldBasis, NewCostBasis: a.CostBasis(), Reason: reason} }

type OwnershipChanged struct{ Base; OldOwner, NewOwner string; OldPct, NewPct float64 }
func NewOwnershipChanged(a *Asset, oldOwner string, oldPct float64) OwnershipChanged { return OwnershipChanged{Base: NewBase("OwnershipChanged", a.ID()), OldOwner: oldOwner, NewOwner: a.OwnerID(), OldPct: oldPct, NewPct: a.OwnershipPercentage()} }
