package domain

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type Base struct{ Type, ID string; Time time.Time }
func NewBase(evt, eid string) Base { return Base{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b Base) EventName() string { return b.Type }; func (b Base) EntityID() string { return b.ID }

type InstRegistered struct{ Base }
func NewInstRegistered(i *Institution) InstRegistered { return InstRegistered{Base: NewBase("InstitutionRegistered", i.ID())} }
type InstVerified struct{ Base }
func NewInstVerified(i *Institution) InstVerified { return InstVerified{Base: NewBase("InstitutionVerified", i.ID())} }
type InstActivated struct{ Base }
func NewInstActivated(i *Institution) InstActivated { return InstActivated{Base: NewBase("InstitutionActivated", i.ID())} }
type InstSuspended struct{ Base }
func NewInstSuspended(i *Institution) InstSuspended { return InstSuspended{Base: NewBase("InstitutionSuspended", i.ID())} }
type InstClosed struct{ Base }
func NewInstClosed(i *Institution) InstClosed { return InstClosed{Base: NewBase("InstitutionClosed", i.ID())} }
type InstArchived struct{ Base }
func NewInstArchived(i *Institution) InstArchived { return InstArchived{Base: NewBase("InstitutionArchived", i.ID())} }
type InstMetadataUpdated struct{ Base }
func NewInstMetadataUpdated(i *Institution) InstMetadataUpdated { return InstMetadataUpdated{Base: NewBase("InstitutionMetadataUpdated", i.ID())} }
type InstTrustChanged struct{ Base; OldTrust, NewTrust TrustLevel }
func NewInstTrustChanged(i *Institution, old TrustLevel) InstTrustChanged { return InstTrustChanged{Base: NewBase("InstitutionTrustChanged", i.ID()), OldTrust: old, NewTrust: i.TrustLevel()} }
type InstConnectivityChanged struct{ Base }
func NewInstConnectivityChanged(i *Institution) InstConnectivityChanged { return InstConnectivityChanged{Base: NewBase("InstitutionConnectivityChanged", i.ID())} }
type InstHealthChanged struct{ Base; OldHealth, NewHealth InstitutionHealth }
func NewInstHealthChanged(i *Institution, old InstitutionHealth) InstHealthChanged { return InstHealthChanged{Base: NewBase("InstitutionHealthChanged", i.ID()), OldHealth: old, NewHealth: i.Health()} }
type InstConfidenceChanged struct{ Base; OldConf, NewConf InstitutionConfidence }
func NewInstConfidenceChanged(i *Institution, old InstitutionConfidence) InstConfidenceChanged { return InstConfidenceChanged{Base: NewBase("InstitutionConfidenceChanged", i.ID()), OldConf: old, NewConf: i.Confidence()} }
type InstProductAdded struct{ Base; Product FinancialProduct }
func NewInstProductAdded(i *Institution, p FinancialProduct) InstProductAdded { return InstProductAdded{Base: NewBase("InstitutionProductAdded", i.ID()), Product: p} }
