package domain

import (
	"errors"
	"time"
)

type Household struct {
	id               string
	name             string
	householdType    HouseholdType
	headOfHouseholdID string
	members          []HouseholdMember
	status           HouseholdStatus
	currency         string
	country          string
	totalAssets      int64
	totalLiabilities int64
	totalNetWorth    int64
	health           HouseholdHealth
	metadata         map[string]string
	tags             []string
	notes            string
	linkedAccounts    []LinkedAccount
	linkedGoals       []LinkedGoal
	linkedBudgets     []LinkedBudget
	goalContributions []GoalContribution
	createdAt        time.Time
	updatedAt        time.Time
}

func (h *Household) ID() string                       { return h.id }
func (h *Household) Name() string                     { return h.name }
func (h *Household) HouseholdType() HouseholdType       { return h.householdType }
func (h *Household) HeadOfHouseholdID() string         { return h.headOfHouseholdID }
func (h *Household) Members() []HouseholdMember        { return h.members }
func (h *Household) Status() HouseholdStatus           { return h.status }
func (h *Household) Currency() string                  { return h.currency }
func (h *Household) Country() string                   { return h.country }
func (h *Household) TotalAssets() int64                { return h.totalAssets }
func (h *Household) TotalLiabilities() int64           { return h.totalLiabilities }
func (h *Household) TotalNetWorth() int64              { return h.totalNetWorth }
func (h *Household) Health() HouseholdHealth           { return h.health }
func (h *Household) Metadata() map[string]string       { return h.metadata }
func (h *Household) Tags() []string                    { return h.tags }
func (h *Household) Notes() string                     { return h.notes }
func (h *Household) LinkedAccounts() []LinkedAccount   { return h.linkedAccounts }
func (h *Household) LinkedGoals() []LinkedGoal         { return h.linkedGoals }
func (h *Household) LinkedBudgets() []LinkedBudget     { return h.linkedBudgets }
func (h *Household) GoalContributions() []GoalContribution { return h.goalContributions }
func (h *Household) CreatedAt() time.Time              { return h.createdAt }
func (h *Household) UpdatedAt() time.Time              { return h.updatedAt }

func (h *Household) SetStatus(s HouseholdStatus) {
	h.status = s
	h.updatedAt = time.Now().UTC()
}

func (h *Household) SetHealth(health HouseholdHealth) {
	h.health = health
	h.updatedAt = time.Now().UTC()
}

func (h *Household) InviteMember(inviterID, userID string, role MemberRole) error {
	if !h.CanManageMembers(inviterID) {
		return errors.New("user does not have permission to manage members")
	}
	if h.headOfHouseholdID == userID {
		return errors.New("user is already the head of household")
	}
	for _, m := range h.members {
		if m.UserID == userID {
			return errors.New("user is already a member or invited")
		}
	}
	h.members = append(h.members, HouseholdMember{
		UserID:       userID,
		Role:         role,
		AddedAt:      time.Now().UTC().Format(time.RFC3339),
		InviteStatus: "Pending",
	})
	h.updatedAt = time.Now().UTC()
	return nil
}

func (h *Household) AcceptInvite(userID string) error {
	for i, m := range h.members {
		if m.UserID == userID {
			if m.InviteStatus == "Accepted" {
				return errors.New("invite already accepted")
			}
			h.members[i].InviteStatus = "Accepted"
			h.updatedAt = time.Now().UTC()
			return nil
		}
	}
	return errors.New("user not invited to household")
}

func (h *Household) RemoveMember(removerID, userID string) error {
	if !h.CanManageMembers(removerID) && removerID != userID {
		return errors.New("user does not have permission to remove members")
	}
	if h.headOfHouseholdID == userID {
		return errors.New("cannot remove the head of household")
	}
	if removerID == userID {
		// allow self removal
	}

	var filtered []HouseholdMember
	found := false
	for _, m := range h.members {
		if m.UserID == userID {
			found = true
			continue
		}
		filtered = append(filtered, m)
	}
	if !found {
		return errors.New("user is not a member of the household")
	}
	h.members = filtered
	h.updatedAt = time.Now().UTC()
	return nil
}

func (h *Household) UpdateMemberRole(updaterID, userID string, newRole MemberRole) error {
	if !h.CanManageMembers(updaterID) {
		return errors.New("user does not have permission to manage members")
	}
	if h.headOfHouseholdID == userID {
		return errors.New("cannot change the role of the head of household")
	}
	for i, m := range h.members {
		if m.UserID == userID {
			h.members[i].Role = newRole
			h.updatedAt = time.Now().UTC()
			return nil
		}
	}
	return errors.New("user is not a member of the household")
}

func (h *Household) SetFinancials(assets, liabilities, netWorth int64) {
	h.totalAssets = assets
	h.totalLiabilities = liabilities
	h.totalNetWorth = netWorth
	h.updatedAt = time.Now().UTC()
}

func (h *Household) LinkAccount(userID string, acc LinkedAccount) error {
	if !h.CanLink(userID) {
		return errors.New("user does not have permission to link accounts")
	}
	h.linkedAccounts = append(h.linkedAccounts, acc)
	h.updatedAt = time.Now().UTC()
	return nil
}
func (h *Household) UnlinkAccount(userID, accountID string) error {
	if !h.CanLink(userID) {
		return errors.New("user does not have permission to unlink accounts")
	}
	var filtered []LinkedAccount
	for _, a := range h.linkedAccounts {
		if a.AccountID != accountID { filtered = append(filtered, a) }
	}
	h.linkedAccounts = filtered
	h.updatedAt = time.Now().UTC()
	return nil
}

func (h *Household) LinkGoal(userID string, goal LinkedGoal) error {
	if !h.CanLink(userID) {
		return errors.New("user does not have permission to link goals")
	}
	h.linkedGoals = append(h.linkedGoals, goal)
	h.updatedAt = time.Now().UTC()
	return nil
}
func (h *Household) UnlinkGoal(userID, goalID string) error {
	if !h.CanLink(userID) {
		return errors.New("user does not have permission to unlink goals")
	}
	var filtered []LinkedGoal
	for _, g := range h.linkedGoals {
		if g.GoalID != goalID { filtered = append(filtered, g) }
	}
	h.linkedGoals = filtered
	h.updatedAt = time.Now().UTC()
	return nil
}

func (h *Household) LinkBudget(userID string, budget LinkedBudget) error {
	if !h.CanLink(userID) {
		return errors.New("user does not have permission to link budgets")
	}
	h.linkedBudgets = append(h.linkedBudgets, budget)
	h.updatedAt = time.Now().UTC()
	return nil
}
func (h *Household) UnlinkBudget(userID, budgetID string) error {
	if !h.CanLink(userID) {
		return errors.New("user does not have permission to unlink budgets")
	}
	var filtered []LinkedBudget
	for _, b := range h.linkedBudgets {
		if b.BudgetID != budgetID { filtered = append(filtered, b) }
	}
	h.linkedBudgets = filtered
	h.updatedAt = time.Now().UTC()
	return nil
}

func (h *Household) AddGoalContribution(contrib GoalContribution) {
	h.goalContributions = append(h.goalContributions, contrib)
	h.updatedAt = time.Now().UTC()
}

func ReconstructFromDB(
	id, name string,
	hhType HouseholdType,
	headID string,
	members []HouseholdMember,
	status HouseholdStatus,
	currency, country string,
	totalAssets, totalLiabilities, totalNetWorth int64,
	health HouseholdHealth,
	tags []string,
	notes string,
	linkedAccounts []LinkedAccount,
	linkedGoals []LinkedGoal,
	linkedBudgets []LinkedBudget,
	goalContributions []GoalContribution,
	createdAt, updatedAt time.Time,
) *Household {
	if tags == nil { tags = []string{} }
	if linkedAccounts == nil { linkedAccounts = []LinkedAccount{} }
	if linkedGoals == nil { linkedGoals = []LinkedGoal{} }
	if linkedBudgets == nil { linkedBudgets = []LinkedBudget{} }
	if goalContributions == nil { goalContributions = []GoalContribution{} }
	return &Household{
		id: id, name: name, householdType: hhType,
		headOfHouseholdID: headID, members: members, status: status,
		currency: currency, country: country,
		totalAssets: totalAssets, totalLiabilities: totalLiabilities,
		totalNetWorth: totalNetWorth, health: health,
		tags: tags, notes: notes, metadata: map[string]string{},
		linkedAccounts: linkedAccounts, linkedGoals: linkedGoals,
		linkedBudgets: linkedBudgets, goalContributions: goalContributions,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}
