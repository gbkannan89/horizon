package domain

import "time"

type DomainEvent interface{ EventName() string; EntityID() string }
type Publisher interface{ Publish(event DomainEvent) error }

type BaseDomainEvent struct{ Type, ID string; Time time.Time }

func NewBaseDomainEvent(evt, eid string) BaseDomainEvent {
	return BaseDomainEvent{Type: evt, ID: eid, Time: time.Now().UTC()}
}

func (b BaseDomainEvent) EventName() string { return b.Type }
func (b BaseDomainEvent) EntityID() string  { return b.ID }

type BudgetCreated struct{ BaseDomainEvent; Name string; Period BudgetPeriod }
func NewBudgetCreated(b *Budget) BudgetCreated {
	return BudgetCreated{
		BaseDomainEvent: NewBaseDomainEvent("BudgetCreated", b.ID()),
		Name: b.Name(), Period: b.Period(),
	}
}

type BudgetActivated struct{ BaseDomainEvent }
func NewBudgetActivated(b *Budget) BudgetActivated {
	return BudgetActivated{BaseDomainEvent: NewBaseDomainEvent("BudgetActivated", b.ID())}
}

type BudgetPaused struct{ BaseDomainEvent }
func NewBudgetPaused(b *Budget) BudgetPaused {
	return BudgetPaused{BaseDomainEvent: NewBaseDomainEvent("BudgetPaused", b.ID())}
}

type BudgetResumed struct{ BaseDomainEvent }
func NewBudgetResumed(b *Budget) BudgetResumed {
	return BudgetResumed{BaseDomainEvent: NewBaseDomainEvent("BudgetResumed", b.ID())}
}

type BudgetCompleted struct{ BaseDomainEvent }
func NewBudgetCompleted(b *Budget) BudgetCompleted {
	return BudgetCompleted{BaseDomainEvent: NewBaseDomainEvent("BudgetCompleted", b.ID())}
}

type BudgetArchived struct{ BaseDomainEvent }
func NewBudgetArchived(b *Budget) BudgetArchived {
	return BudgetArchived{BaseDomainEvent: NewBaseDomainEvent("BudgetArchived", b.ID())}
}

type BudgetOverspent struct{ BaseDomainEvent; Category string; OverspentAmount int64 }
func NewBudgetOverspent(b *Budget, category string, amount int64) BudgetOverspent {
	return BudgetOverspent{
		BaseDomainEvent: NewBaseDomainEvent("BudgetOverspent", b.ID()),
		Category: category, OverspentAmount: amount,
	}
}

type BudgetHealthChanged struct{ BaseDomainEvent; OldHealth, NewHealth BudgetHealth }
func NewBudgetHealthChanged(b *Budget, old BudgetHealth) BudgetHealthChanged {
	return BudgetHealthChanged{
		BaseDomainEvent: NewBaseDomainEvent("BudgetHealthChanged", b.ID()),
		OldHealth: old, NewHealth: b.determineHealth(),
	}
}

type BudgetCategoryUpdated struct{ BaseDomainEvent; CategoryID string; Budgeted, Spent int64 }
func NewBudgetCategoryUpdated(b *Budget, cat BudgetCategory) BudgetCategoryUpdated {
	return BudgetCategoryUpdated{
		BaseDomainEvent: NewBaseDomainEvent("BudgetCategoryUpdated", b.ID()),
		CategoryID: cat.ID, Budgeted: cat.BudgetedAmount, Spent: cat.SpentAmount,
	}
}
