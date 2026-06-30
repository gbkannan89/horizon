package domain

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }
type BaseDomainEvent struct{ Type, ID string; Time time.Time }
func NewBaseDomainEvent(evt, eid string) BaseDomainEvent { return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()} }
func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string  { return b.ID }

type AccountCreated struct{ BaseDomainEvent; OwnerID string; AccountType string; Currency string }
func NewAccountCreated(a *Account) AccountCreated { return AccountCreated{BaseDomainEvent: NewBaseDomainEvent("AccountCreated", a.ID()), OwnerID: a.OwnerID(), AccountType: string(a.AccountType()), Currency: a.Currency()} }

type AccountActivated struct{ BaseDomainEvent }
func NewAccountActivated(a *Account) AccountActivated { return AccountActivated{BaseDomainEvent: NewBaseDomainEvent("AccountActivated", a.ID())} }

type AccountFrozen struct{ BaseDomainEvent }
func NewAccountFrozen(a *Account) AccountFrozen { return AccountFrozen{BaseDomainEvent: NewBaseDomainEvent("AccountFrozen", a.ID())} }

type AccountUnfrozen struct{ BaseDomainEvent }
func NewAccountUnfrozen(a *Account) AccountUnfrozen { return AccountUnfrozen{BaseDomainEvent: NewBaseDomainEvent("AccountUnfrozen", a.ID())} }

type AccountClosed struct{ BaseDomainEvent; Reason string }
func NewAccountClosed(a *Account, reason string) AccountClosed { return AccountClosed{BaseDomainEvent: NewBaseDomainEvent("AccountClosed", a.ID()), Reason: reason} }

type AccountArchived struct{ BaseDomainEvent }
func NewAccountArchived(a *Account) AccountArchived { return AccountArchived{BaseDomainEvent: NewBaseDomainEvent("AccountArchived", a.ID())} }

type AccountDormant struct{ BaseDomainEvent }
func NewAccountDormant(a *Account) AccountDormant { return AccountDormant{BaseDomainEvent: NewBaseDomainEvent("AccountDormant", a.ID())} }

type BalanceRecalculated struct{ BaseDomainEvent; CurrentBalance int64; AvailableBalance int64 }
func NewBalanceRecalculated(a *Account, cur, avail int64) BalanceRecalculated { return BalanceRecalculated{BaseDomainEvent: NewBaseDomainEvent("BalanceRecalculated", a.ID()), CurrentBalance: cur, AvailableBalance: avail} }

type AccountTypeChanged struct{ BaseDomainEvent; OldType, NewType AccountType }
func NewAccountTypeChanged(a *Account, old AccountType) AccountTypeChanged { return AccountTypeChanged{BaseDomainEvent: NewBaseDomainEvent("AccountTypeChanged", a.ID()), OldType: old, NewType: a.AccountType()} }

type AccountOwnershipChanged struct{ BaseDomainEvent; OldOwner, NewOwner string }
func NewAccountOwnershipChanged(a *Account, oldOwner string) AccountOwnershipChanged { return AccountOwnershipChanged{BaseDomainEvent: NewBaseDomainEvent("AccountOwnershipChanged", a.ID()), OldOwner: oldOwner, NewOwner: a.OwnerID()} }

type AccountVisibilityChanged struct{ BaseDomainEvent; OldVis, NewVis Visibility }
func NewAccountVisibilityChanged(a *Account, old Visibility) AccountVisibilityChanged { return AccountVisibilityChanged{BaseDomainEvent: NewBaseDomainEvent("AccountVisibilityChanged", a.ID()), OldVis: old, NewVis: a.Visibility()} }
