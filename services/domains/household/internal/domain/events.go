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

type HouseholdCreated struct{ BaseDomainEvent; Name string; HeadID string }
func NewHouseholdCreated(h *Household) HouseholdCreated {
	return HouseholdCreated{
		BaseDomainEvent: NewBaseDomainEvent("HouseholdCreated", h.ID()),
		Name: h.Name(), HeadID: h.HeadOfHouseholdID(),
	}
}

type HouseholdActivated struct{ BaseDomainEvent }
func NewHouseholdActivated(h *Household) HouseholdActivated {
	return HouseholdActivated{BaseDomainEvent: NewBaseDomainEvent("HouseholdActivated", h.ID())}
}

type HouseholdPaused struct{ BaseDomainEvent }
func NewHouseholdPaused(h *Household) HouseholdPaused {
	return HouseholdPaused{BaseDomainEvent: NewBaseDomainEvent("HouseholdPaused", h.ID())}
}

type HouseholdResumed struct{ BaseDomainEvent }
func NewHouseholdResumed(h *Household) HouseholdResumed {
	return HouseholdResumed{BaseDomainEvent: NewBaseDomainEvent("HouseholdResumed", h.ID())}
}

type HouseholdDissolved struct{ BaseDomainEvent }
func NewHouseholdDissolved(h *Household) HouseholdDissolved {
	return HouseholdDissolved{BaseDomainEvent: NewBaseDomainEvent("HouseholdDissolved", h.ID())}
}

type HouseholdArchived struct{ BaseDomainEvent }
func NewHouseholdArchived(h *Household) HouseholdArchived {
	return HouseholdArchived{BaseDomainEvent: NewBaseDomainEvent("HouseholdArchived", h.ID())}
}

type MemberAdded struct{ BaseDomainEvent; UserID string; Role MemberRole }
func NewMemberAdded(h *Household, userID string, role MemberRole) MemberAdded {
	return MemberAdded{
		BaseDomainEvent: NewBaseDomainEvent("MemberAdded", h.ID()),
		UserID: userID, Role: role,
	}
}

type MemberRemoved struct{ BaseDomainEvent; UserID string }
func NewMemberRemoved(h *Household, userID string) MemberRemoved {
	return MemberRemoved{
		BaseDomainEvent: NewBaseDomainEvent("MemberRemoved", h.ID()),
		UserID: userID,
	}
}

type MemberRoleChanged struct{ BaseDomainEvent; UserID string; OldRole, NewRole MemberRole }
func NewMemberRoleChanged(h *Household, userID string, oldRole, newRole MemberRole) MemberRoleChanged {
	return MemberRoleChanged{
		BaseDomainEvent: NewBaseDomainEvent("MemberRoleChanged", h.ID()),
		UserID: userID, OldRole: oldRole, NewRole: newRole,
	}
}

type HouseholdHealthChanged struct{ BaseDomainEvent; OldHealth, NewHealth HouseholdHealth }
func NewHouseholdHealthChanged(h *Household, old HouseholdHealth) HouseholdHealthChanged {
	return HouseholdHealthChanged{
		BaseDomainEvent: NewBaseDomainEvent("HouseholdHealthChanged", h.ID()),
		OldHealth: old, NewHealth: h.Health(),
	}
}

type AccountLinked struct{ BaseDomainEvent; AccountID, AddedBy string }
func NewAccountLinked(h *Household, accountID, addedBy string) AccountLinked {
	return AccountLinked{
		BaseDomainEvent: NewBaseDomainEvent("AccountLinked", h.ID()),
		AccountID: accountID, AddedBy: addedBy,
	}
}

type AccountUnlinked struct{ BaseDomainEvent; AccountID string }
func NewAccountUnlinked(h *Household, accountID string) AccountUnlinked {
	return AccountUnlinked{
		BaseDomainEvent: NewBaseDomainEvent("AccountUnlinked", h.ID()),
		AccountID: accountID,
	}
}

type GoalLinked struct{ BaseDomainEvent; GoalID, AddedBy string }
func NewGoalLinked(h *Household, goalID, addedBy string) GoalLinked {
	return GoalLinked{
		BaseDomainEvent: NewBaseDomainEvent("GoalLinked", h.ID()),
		GoalID: goalID, AddedBy: addedBy,
	}
}

type GoalUnlinked struct{ BaseDomainEvent; GoalID string }
func NewGoalUnlinked(h *Household, goalID string) GoalUnlinked {
	return GoalUnlinked{
		BaseDomainEvent: NewBaseDomainEvent("GoalUnlinked", h.ID()),
		GoalID: goalID,
	}
}

type BudgetLinked struct{ BaseDomainEvent; BudgetID, AddedBy string }
func NewBudgetLinked(h *Household, budgetID, addedBy string) BudgetLinked {
	return BudgetLinked{
		BaseDomainEvent: NewBaseDomainEvent("BudgetLinked", h.ID()),
		BudgetID: budgetID, AddedBy: addedBy,
	}
}

type BudgetUnlinked struct{ BaseDomainEvent; BudgetID string }
func NewBudgetUnlinked(h *Household, budgetID string) BudgetUnlinked {
	return BudgetUnlinked{
		BaseDomainEvent: NewBaseDomainEvent("BudgetUnlinked", h.ID()),
		BudgetID: budgetID,
	}
}

type GoalContributionAdded struct{ BaseDomainEvent; GoalID, UserID string; Amount int64 }
func NewGoalContributionAdded(h *Household, contrib GoalContribution) GoalContributionAdded {
	return GoalContributionAdded{
		BaseDomainEvent: NewBaseDomainEvent("GoalContributionAdded", h.ID()),
		GoalID: contrib.GoalID, UserID: contrib.UserID, Amount: contrib.Amount,
	}
}
