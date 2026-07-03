package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeHeadMember(userID string) HouseholdMember {
	return HouseholdMember{
		UserID:       userID,
		Role:         RoleHead,
		AddedAt:      time.Now().UTC().Format(time.RFC3339),
		InviteStatus: "Accepted",
	}
}

func makeAcceptedMember(userID string, role MemberRole) HouseholdMember {
	return HouseholdMember{
		UserID:       userID,
		Role:         role,
		AddedAt:      time.Now().UTC().Format(time.RFC3339),
		InviteStatus: "Accepted",
	}
}

func makePendingMember(userID string, role MemberRole) HouseholdMember {
	return HouseholdMember{
		UserID:       userID,
		Role:         role,
		AddedAt:      time.Now().UTC().Format(time.RFC3339),
		InviteStatus: "Pending",
	}
}

// ---------------------------------------------------------------------------
// Factory
// ---------------------------------------------------------------------------

func TestHouseholdFactory_Create(t *testing.T) {
	t.Parallel()

	factory := NewHouseholdFactory()
	headID := "user-head"
	headMember := makeHeadMember(headID)

	t.Run("valid creation", func(t *testing.T) {
		h, err := factory.Create("hh-1", "My Home", HHCouple, headID,
			[]HouseholdMember{headMember}, "USD", "US", []string{"important"}, "some notes")
		require.NoError(t, err)
		assert.Equal(t, "hh-1", h.ID())
		assert.Equal(t, "My Home", h.Name())
		assert.Equal(t, HHCouple, h.HouseholdType())
		assert.Equal(t, headID, h.HeadOfHouseholdID())
		assert.Equal(t, HHStatusDraft, h.Status())
		assert.Equal(t, "USD", h.Currency())
		assert.Equal(t, "US", h.Country())
		assert.Equal(t, HHHealthy, h.Health())
		assert.Equal(t, []string{"important"}, h.Tags())
		assert.Equal(t, "some notes", h.Notes())
		assert.NotZero(t, h.CreatedAt())
		assert.NotZero(t, h.UpdatedAt())
	})

	t.Run("fails on empty name", func(t *testing.T) {
		_, err := factory.Create("hh-2", "", HHCouple, headID,
			[]HouseholdMember{headMember}, "USD", "US", nil, "")
		assert.EqualError(t, err, "household name is required")
	})

	t.Run("fails on name too long", func(t *testing.T) {
		longName := strings.Repeat("a", 201)
		_, err := factory.Create("hh-3", longName, HHCouple, headID,
			[]HouseholdMember{headMember}, "USD", "US", nil, "")
		assert.EqualError(t, err, "household name must be 200 characters or fewer")
	})

	t.Run("fails on invalid household type", func(t *testing.T) {
		_, err := factory.Create("hh-4", "test", HouseholdType("Invalid"), headID,
			[]HouseholdMember{headMember}, "USD", "US", nil, "")
		assert.EqualError(t, err, "household type is not recognized")
	})

	t.Run("fails on invalid currency", func(t *testing.T) {
		_, err := factory.Create("hh-5", "test", HHSingle, headID,
			[]HouseholdMember{headMember}, "US", "US", nil, "")
		assert.EqualError(t, err, "currency must be a 3-letter ISO code")
	})

	t.Run("fails on empty members", func(t *testing.T) {
		_, err := factory.Create("hh-6", "test", HHCouple, headID,
			[]HouseholdMember{}, "USD", "US", nil, "")
		assert.EqualError(t, err, "household must have at least one member")
	})

	t.Run("fails on too many members", func(t *testing.T) {
		members := make([]HouseholdMember, 51)
		for i := range members {
			members[i] = makeHeadMember(headID)
		}
		_, err := factory.Create("hh-7", "test", HHCouple, headID,
			members, "USD", "US", nil, "")
		assert.EqualError(t, err, "household cannot have more than 50 members")
	})

	t.Run("fails when head not in members", func(t *testing.T) {
		other := makeAcceptedMember("user-other", RoleMember)
		_, err := factory.Create("hh-8", "test", HHSingle, headID,
			[]HouseholdMember{other}, "USD", "US", nil, "")
		assert.EqualError(t, err, "head of household must be a member")
	})

	t.Run("fails when head does not have Head role", func(t *testing.T) {
		wrongRole := makeAcceptedMember(headID, RoleMember)
		_, err := factory.Create("hh-9", "test", HHSingle, headID,
			[]HouseholdMember{wrongRole}, "USD", "US", nil, "")
		assert.EqualError(t, err, "head of household must have Head role")
	})

	t.Run("fails on invalid member role", func(t *testing.T) {
		invalid := HouseholdMember{UserID: "user-x", Role: MemberRole("Bogus")}
		_, err := factory.Create("hh-10", "test", HHCouple, headID,
			[]HouseholdMember{headMember, invalid}, "USD", "US", nil, "")
		assert.EqualError(t, err, "member role is not recognized")
	})

	t.Run("nil tags becomes empty slice", func(t *testing.T) {
		h, err := factory.Create("hh-11", "test", HHSingle, headID,
			[]HouseholdMember{headMember}, "USD", "US", nil, "")
		require.NoError(t, err)
		assert.NotNil(t, h.Tags())
		assert.Empty(t, h.Tags())
	})
}

// ---------------------------------------------------------------------------
// State Transitions
// ---------------------------------------------------------------------------

func TestStateTransitions(t *testing.T) {
	t.Parallel()

	t.Run("ActivateHousehold", func(t *testing.T) {
		t.Run("valid from Draft", func(t *testing.T) {
			h := testDraftHousehold()
			err := ActivateHousehold(h)
			assert.NoError(t, err)
			assert.Equal(t, HHStatusActive, h.Status())
		})

		t.Run("invalid from Active", func(t *testing.T) {
			h := testActiveHousehold()
			err := ActivateHousehold(h)
			assert.EqualError(t, err, "only draft households can be activated")
		})

		t.Run("invalid from Paused", func(t *testing.T) {
			h := testPausedHousehold()
			err := ActivateHousehold(h)
			assert.EqualError(t, err, "only draft households can be activated")
		})

		t.Run("invalid from Dissolved", func(t *testing.T) {
			h := testDissolvedHousehold()
			err := ActivateHousehold(h)
			assert.EqualError(t, err, "only draft households can be activated")
		})

		t.Run("invalid from Archived", func(t *testing.T) {
			h := testArchivedHousehold()
			err := ActivateHousehold(h)
			assert.EqualError(t, err, "only draft households can be activated")
		})
	})

	t.Run("PauseHousehold", func(t *testing.T) {
		t.Run("valid from Active", func(t *testing.T) {
			h := testActiveHousehold()
			err := PauseHousehold(h)
			assert.NoError(t, err)
			assert.Equal(t, HHStatusPaused, h.Status())
		})

		t.Run("invalid from Draft", func(t *testing.T) {
			h := testDraftHousehold()
			err := PauseHousehold(h)
			assert.EqualError(t, err, "only active households can be paused")
		})

		t.Run("invalid from Paused", func(t *testing.T) {
			h := testPausedHousehold()
			err := PauseHousehold(h)
			assert.EqualError(t, err, "only active households can be paused")
		})

		t.Run("invalid from Dissolved", func(t *testing.T) {
			h := testDissolvedHousehold()
			err := PauseHousehold(h)
			assert.EqualError(t, err, "only active households can be paused")
		})
	})

	t.Run("ResumeHousehold", func(t *testing.T) {
		t.Run("valid from Paused", func(t *testing.T) {
			h := testPausedHousehold()
			err := ResumeHousehold(h)
			assert.NoError(t, err)
			assert.Equal(t, HHStatusActive, h.Status())
		})

		t.Run("invalid from Draft", func(t *testing.T) {
			h := testDraftHousehold()
			err := ResumeHousehold(h)
			assert.EqualError(t, err, "only paused households can be resumed")
		})

		t.Run("invalid from Active", func(t *testing.T) {
			h := testActiveHousehold()
			err := ResumeHousehold(h)
			assert.EqualError(t, err, "only paused households can be resumed")
		})

		t.Run("invalid from Dissolved", func(t *testing.T) {
			h := testDissolvedHousehold()
			err := ResumeHousehold(h)
			assert.EqualError(t, err, "only paused households can be resumed")
		})

		t.Run("invalid from Archived", func(t *testing.T) {
			h := testArchivedHousehold()
			err := ResumeHousehold(h)
			assert.EqualError(t, err, "only paused households can be resumed")
		})
	})

	t.Run("DissolveHousehold", func(t *testing.T) {
		t.Run("valid from Active", func(t *testing.T) {
			h := testActiveHousehold()
			err := DissolveHousehold(h)
			assert.NoError(t, err)
			assert.Equal(t, HHStatusDissolved, h.Status())
		})

		t.Run("valid from Paused", func(t *testing.T) {
			h := testPausedHousehold()
			err := DissolveHousehold(h)
			assert.NoError(t, err)
			assert.Equal(t, HHStatusDissolved, h.Status())
		})

		t.Run("invalid from Draft", func(t *testing.T) {
			h := testDraftHousehold()
			err := DissolveHousehold(h)
			assert.EqualError(t, err, "only active or paused households can be dissolved")
		})

		t.Run("invalid from Dissolved", func(t *testing.T) {
			h := testDissolvedHousehold()
			err := DissolveHousehold(h)
			assert.EqualError(t, err, "only active or paused households can be dissolved")
		})

		t.Run("invalid from Archived", func(t *testing.T) {
			h := testArchivedHousehold()
			err := DissolveHousehold(h)
			assert.EqualError(t, err, "only active or paused households can be dissolved")
		})
	})

	t.Run("ArchiveHousehold", func(t *testing.T) {
		t.Run("valid from Dissolved", func(t *testing.T) {
			h := testDissolvedHousehold()
			err := ArchiveHousehold(h)
			assert.NoError(t, err)
			assert.Equal(t, HHStatusArchived, h.Status())
		})

		t.Run("invalid from Draft", func(t *testing.T) {
			h := testDraftHousehold()
			err := ArchiveHousehold(h)
			assert.EqualError(t, err, "only dissolved households can be archived")
		})

		t.Run("invalid from Active", func(t *testing.T) {
			h := testActiveHousehold()
			err := ArchiveHousehold(h)
			assert.EqualError(t, err, "only dissolved households can be archived")
		})

		t.Run("invalid from Paused", func(t *testing.T) {
			h := testPausedHousehold()
			err := ArchiveHousehold(h)
			assert.EqualError(t, err, "only dissolved households can be archived")
		})

		t.Run("invalid from Archived", func(t *testing.T) {
			h := testArchivedHousehold()
			err := ArchiveHousehold(h)
			assert.EqualError(t, err, "only dissolved households can be archived")
		})
	})
}

// ---------------------------------------------------------------------------
// Member Management
// ---------------------------------------------------------------------------

func TestMemberManagement(t *testing.T) {
	t.Parallel()

	headID := "user-head"
	adminID := "user-admin"
	memberID := "user-member"
	viewerID := "user-viewer"
	outsiderID := "user-outsider"

	t.Run("InviteMember", func(t *testing.T) {
		t.Run("head can invite", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(headID, "new-user", RoleMember)
			assert.NoError(t, err)
			assert.Len(t, h.Members(), 5) // 4 initial + 1 new
		})

		t.Run("admin can invite", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(adminID, "new-user", RoleMember)
			assert.NoError(t, err)
		})

		t.Run("member cannot invite", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(memberID, "new-user", RoleMember)
			assert.EqualError(t, err, "user does not have permission to manage members")
		})

		t.Run("viewer cannot invite", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(viewerID, "new-user", RoleMember)
			assert.EqualError(t, err, "user does not have permission to manage members")
		})

		t.Run("outsider cannot invite", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(outsiderID, "new-user", RoleMember)
			assert.EqualError(t, err, "user does not have permission to manage members")
		})

		t.Run("cannot invite head of household", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(headID, headID, RoleMember)
			assert.EqualError(t, err, "user is already the head of household")
		})

		t.Run("cannot invite existing member", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(headID, adminID, RoleMember)
			assert.EqualError(t, err, "user is already a member or invited")
		})

		t.Run("invite sets pending status", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.InviteMember(headID, "fresh-user", RoleMember)
			require.NoError(t, err)
			for _, m := range h.Members() {
				if m.UserID == "fresh-user" {
					assert.Equal(t, "Pending", m.InviteStatus)
					assert.Equal(t, RoleMember, m.Role)
					return
				}
			}
			t.Fatal("invited user not found")
		})
	})

	t.Run("AcceptInvite", func(t *testing.T) {
		t.Run("accepts pending invite", func(t *testing.T) {
			h := testDraftHousehold()
			h.InviteMember(headID, "pending-user", RoleMember) // nolint:errcheck
			err := h.AcceptInvite("pending-user")
			assert.NoError(t, err)
			for _, m := range h.Members() {
				if m.UserID == "pending-user" {
					assert.Equal(t, "Accepted", m.InviteStatus)
					return
				}
			}
		})

		t.Run("fails if already accepted (head is already accepted)", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.AcceptInvite(headID)
			assert.EqualError(t, err, "invite already accepted")
		})

		t.Run("fails for non-invited user", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.AcceptInvite("unknown")
			assert.EqualError(t, err, "user not invited to household")
		})

		t.Run("fails on already accepted invite", func(t *testing.T) {
			h := testActiveHousehold()
			// Add a pending member, accept, then accept again
			h.InviteMember(headID, "duel", RoleMember) // nolint:errcheck
			require.NoError(t, h.AcceptInvite("duel"))
			err := h.AcceptInvite("duel")
			assert.EqualError(t, err, "invite already accepted")
		})
	})

	t.Run("RemoveMember", func(t *testing.T) {
		t.Run("head can remove member", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(headID, memberID)
			assert.NoError(t, err)
			assert.Len(t, h.Members(), 3)
		})

		t.Run("admin can remove member", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(adminID, memberID)
			assert.NoError(t, err)
		})

		t.Run("member cannot remove others", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(memberID, viewerID)
			assert.EqualError(t, err, "user does not have permission to remove members")
		})

		t.Run("viewer cannot remove others", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(viewerID, memberID)
			assert.EqualError(t, err, "user does not have permission to remove members")
		})

		t.Run("self-removal allowed for member", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(memberID, memberID)
			assert.NoError(t, err)
			assert.Len(t, h.Members(), 3)
		})

		t.Run("self-removal allowed for viewer", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(viewerID, viewerID)
			assert.NoError(t, err)
		})

		t.Run("cannot remove head of household", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(headID, headID)
			assert.EqualError(t, err, "cannot remove the head of household")
		})

		t.Run("cannot remove head by admin", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(adminID, headID)
			assert.EqualError(t, err, "cannot remove the head of household")
		})

		t.Run("fails for non-member target", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.RemoveMember(headID, "ghost")
			assert.EqualError(t, err, "user is not a member of the household")
		})
	})

	t.Run("UpdateMemberRole", func(t *testing.T) {
		t.Run("head can change member role", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.UpdateMemberRole(headID, memberID, RoleAdmin)
			assert.NoError(t, err)
			for _, m := range h.Members() {
				if m.UserID == memberID {
					assert.Equal(t, RoleAdmin, m.Role)
					return
				}
			}
		})

		t.Run("admin can change member role", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.UpdateMemberRole(adminID, viewerID, RoleMember)
			assert.NoError(t, err)
		})

		t.Run("member cannot change roles", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.UpdateMemberRole(memberID, viewerID, RoleAdmin)
			assert.EqualError(t, err, "user does not have permission to manage members")
		})

		t.Run("viewer cannot change roles", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.UpdateMemberRole(viewerID, memberID, RoleAdmin)
			assert.EqualError(t, err, "user does not have permission to manage members")
		})

		t.Run("cannot change head role", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.UpdateMemberRole(headID, headID, RoleMember)
			assert.EqualError(t, err, "cannot change the role of the head of household")
		})

		t.Run("cannot change head role by admin", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.UpdateMemberRole(adminID, headID, RoleMember)
			assert.EqualError(t, err, "cannot change the role of the head of household")
		})

		t.Run("fails for non-existent member", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.UpdateMemberRole(headID, "nobody", RoleMember)
			assert.EqualError(t, err, "user is not a member of the household")
		})
	})
}

// ---------------------------------------------------------------------------
// Link / Unlink Accounts, Goals, Budgets
// ---------------------------------------------------------------------------

func TestLinkUnlink(t *testing.T) {
	t.Parallel()

	headID := "user-head"
	adminID := "user-admin"
	memberID := "user-member"
	viewerID := "user-viewer"
	outsiderID := "user-outsider"

	acc := LinkedAccount{AccountID: "acc-1", AddedBy: headID, AddedAt: time.Now().UTC().Format(time.RFC3339)}
	goal := LinkedGoal{GoalID: "goal-1", AddedBy: headID, AddedAt: time.Now().UTC().Format(time.RFC3339)}
	budget := LinkedBudget{BudgetID: "budget-1", AddedBy: headID, AddedAt: time.Now().UTC().Format(time.RFC3339)}

	t.Run("LinkAccount", func(t *testing.T) {
		t.Run("head can link", func(t *testing.T) {
			h := testActiveHousehold()
			require.NoError(t, h.LinkAccount(headID, acc))
			assert.Len(t, h.LinkedAccounts(), 1)
		})
		t.Run("admin can link", func(t *testing.T) {
			h := testActiveHousehold()
			require.NoError(t, h.LinkAccount(adminID, acc))
		})
		t.Run("member cannot link", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.LinkAccount(memberID, acc)
			assert.EqualError(t, err, "user does not have permission to link accounts")
		})
		t.Run("viewer cannot link", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.LinkAccount(viewerID, acc)
			assert.EqualError(t, err, "user does not have permission to link accounts")
		})
		t.Run("outsider cannot link", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.LinkAccount(outsiderID, acc)
			assert.EqualError(t, err, "user does not have permission to link accounts")
		})
	})

	t.Run("UnlinkAccount", func(t *testing.T) {
		setup := func() *Household {
			h := testActiveHousehold()
			h.LinkAccount(headID, acc) // nolint:errcheck
			return h
		}
		t.Run("head can unlink", func(t *testing.T) {
			h := setup()
			require.NoError(t, h.UnlinkAccount(headID, "acc-1"))
			assert.Empty(t, h.LinkedAccounts())
		})
		t.Run("admin can unlink", func(t *testing.T) {
			h := setup()
			require.NoError(t, h.UnlinkAccount(adminID, "acc-1"))
		})
		t.Run("member cannot unlink", func(t *testing.T) {
			h := setup()
			err := h.UnlinkAccount(memberID, "acc-1")
			assert.EqualError(t, err, "user does not have permission to unlink accounts")
		})
		t.Run("unlink non-existent keeps existing accounts", func(t *testing.T) {
			h := setup()
			require.NoError(t, h.UnlinkAccount(headID, "non-existent"))
			assert.Len(t, h.LinkedAccounts(), 1) // non-matching ID is not filtered out
		})
	})

	t.Run("LinkGoal", func(t *testing.T) {
		t.Run("head can link", func(t *testing.T) {
			h := testActiveHousehold()
			require.NoError(t, h.LinkGoal(headID, goal))
			assert.Len(t, h.LinkedGoals(), 1)
		})
		t.Run("admin can link", func(t *testing.T) {
			h := testActiveHousehold()
			require.NoError(t, h.LinkGoal(adminID, goal))
		})
		t.Run("member cannot link", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.LinkGoal(memberID, goal)
			assert.EqualError(t, err, "user does not have permission to link goals")
		})
	})

	t.Run("UnlinkGoal", func(t *testing.T) {
		setup := func() *Household {
			h := testActiveHousehold()
			h.LinkGoal(headID, goal) // nolint:errcheck
			return h
		}
		t.Run("head can unlink", func(t *testing.T) {
			h := setup()
			require.NoError(t, h.UnlinkGoal(headID, "goal-1"))
			assert.Empty(t, h.LinkedGoals())
		})
		t.Run("admin can unlink", func(t *testing.T) {
			h := setup()
			require.NoError(t, h.UnlinkGoal(adminID, "goal-1"))
		})
		t.Run("member cannot unlink", func(t *testing.T) {
			h := setup()
			err := h.UnlinkGoal(memberID, "goal-1")
			assert.EqualError(t, err, "user does not have permission to unlink goals")
		})
	})

	t.Run("LinkBudget", func(t *testing.T) {
		t.Run("head can link", func(t *testing.T) {
			h := testActiveHousehold()
			require.NoError(t, h.LinkBudget(headID, budget))
			assert.Len(t, h.LinkedBudgets(), 1)
		})
		t.Run("admin can link", func(t *testing.T) {
			h := testActiveHousehold()
			require.NoError(t, h.LinkBudget(adminID, budget))
		})
		t.Run("member cannot link", func(t *testing.T) {
			h := testActiveHousehold()
			err := h.LinkBudget(memberID, budget)
			assert.EqualError(t, err, "user does not have permission to link budgets")
		})
	})

	t.Run("UnlinkBudget", func(t *testing.T) {
		setup := func() *Household {
			h := testActiveHousehold()
			h.LinkBudget(headID, budget) // nolint:errcheck
			return h
		}
		t.Run("head can unlink", func(t *testing.T) {
			h := setup()
			require.NoError(t, h.UnlinkBudget(headID, "budget-1"))
			assert.Empty(t, h.LinkedBudgets())
		})
		t.Run("admin can unlink", func(t *testing.T) {
			h := setup()
			require.NoError(t, h.UnlinkBudget(adminID, "budget-1"))
		})
		t.Run("member cannot unlink", func(t *testing.T) {
			h := setup()
			err := h.UnlinkBudget(memberID, "budget-1")
			assert.EqualError(t, err, "user does not have permission to unlink budgets")
		})
	})

	t.Run("AddGoalContribution", func(t *testing.T) {
		h := testActiveHousehold()
		contrib := GoalContribution{GoalID: "goal-1", UserID: headID, Amount: 1000, Date: "2026-01-01"}
		h.AddGoalContribution(contrib)
		assert.Len(t, h.GoalContributions(), 1)
		assert.Equal(t, int64(1000), h.GoalContributions()[0].Amount)

		h.AddGoalContribution(GoalContribution{GoalID: "goal-1", UserID: memberID, Amount: 500, Date: "2026-02-01"})
		assert.Len(t, h.GoalContributions(), 2)
	})
}

// ---------------------------------------------------------------------------
// SetFinancials / SetHealth / SetStatus
// ---------------------------------------------------------------------------

func TestSetMethods(t *testing.T) {
	t.Parallel()

	t.Run("SetFinancials", func(t *testing.T) {
		h := testDraftHousehold()
		h.SetFinancials(10000, 3000, 7000)
		assert.Equal(t, int64(10000), h.TotalAssets())
		assert.Equal(t, int64(3000), h.TotalLiabilities())
		assert.Equal(t, int64(7000), h.TotalNetWorth())
		assert.False(t, h.UpdatedAt().IsZero())
	})

	t.Run("SetHealth", func(t *testing.T) {
		h := testDraftHousehold()
		h.SetHealth(HHWarning)
		assert.Equal(t, HHWarning, h.Health())
	})

	t.Run("SetStatus", func(t *testing.T) {
		h := testDraftHousehold()
		h.SetStatus(HHStatusActive)
		assert.Equal(t, HHStatusActive, h.Status())
	})
}

// ---------------------------------------------------------------------------
// Permission checks
// ---------------------------------------------------------------------------

func TestPermissions(t *testing.T) {
	t.Parallel()

	headID := "user-head"
	adminID := "user-admin"
	memberID := "user-member"
	viewerID := "user-viewer"
	outsiderID := "user-outsider"
	pendingID := "user-pending"

	t.Run("HasAccess", func(t *testing.T) {
		h := testActiveHousehold()
		assert.True(t, h.HasAccess(headID), "head")
		assert.True(t, h.HasAccess(adminID), "admin")
		assert.True(t, h.HasAccess(memberID), "member")
		assert.True(t, h.HasAccess(viewerID), "viewer")
		assert.False(t, h.HasAccess(outsiderID), "outsider")
		assert.False(t, h.HasAccess(pendingID), "pending user")
	})

	t.Run("CanView", func(t *testing.T) {
		h := testActiveHousehold()
		assert.True(t, h.CanView(headID))
		assert.True(t, h.CanView(adminID))
		assert.False(t, h.CanView(outsiderID))
	})

	t.Run("CanEdit", func(t *testing.T) {
		h := testActiveHousehold()
		assert.True(t, h.CanEdit(headID), "head can edit")
		assert.True(t, h.CanEdit(adminID), "admin can edit")
		assert.False(t, h.CanEdit(memberID), "member cannot edit")
		assert.False(t, h.CanEdit(viewerID), "viewer cannot edit")
		assert.False(t, h.CanEdit(outsiderID), "outsider cannot edit")
		assert.False(t, h.CanEdit(pendingID), "pending user cannot edit")
	})

	t.Run("CanTransition delegates to CanEdit", func(t *testing.T) {
		h := testActiveHousehold()
		assert.True(t, h.CanTransition(headID))
		assert.False(t, h.CanTransition(memberID))
	})

	t.Run("CanLink delegates to CanEdit", func(t *testing.T) {
		h := testActiveHousehold()
		assert.True(t, h.CanLink(headID))
		assert.False(t, h.CanLink(memberID))
	})

	t.Run("CanManageMembers delegates to CanEdit", func(t *testing.T) {
		h := testActiveHousehold()
		assert.True(t, h.CanManageMembers(headID))
		assert.False(t, h.CanManageMembers(viewerID))
	})

	t.Run("GetMemberRole", func(t *testing.T) {
		h := testActiveHousehold()
		role, ok := h.GetMemberRole(headID)
		assert.True(t, ok)
		assert.Equal(t, RoleHead, role)

		role, ok = h.GetMemberRole(adminID)
		assert.True(t, ok)
		assert.Equal(t, RoleAdmin, role)

		role, ok = h.GetMemberRole(memberID)
		assert.True(t, ok)
		assert.Equal(t, RoleMember, role)

		role, ok = h.GetMemberRole(viewerID)
		assert.True(t, ok)
		assert.Equal(t, RoleViewer, role)

		_, ok = h.GetMemberRole(outsiderID)
		assert.False(t, ok)

		_, ok = h.GetMemberRole(pendingID)
		assert.False(t, ok)
	})
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

func TestValidation(t *testing.T) {
	t.Parallel()

	t.Run("ValidateHouseholdName", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
		}{
			{"valid", "My Household", ""},
			{"empty", "", "household name is required"},
			{"spaces only", "   ", "household name is required"},
			{"too long", strings.Repeat("a", 201), "household name must be 200 characters or fewer"},
			{"exactly 200", strings.Repeat("b", 200), ""},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateHouseholdName(tc.input)
				if tc.want == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.want)
				}
			})
		}
	})

	t.Run("ValidateHouseholdType", func(t *testing.T) {
		tests := []struct {
			name  string
			input HouseholdType
			want  string
		}{
			{"Single", HHSingle, ""},
			{"Couple", HHCouple, ""},
			{"Family", HHFamily, ""},
			{"Roommates", HHRoommates, ""},
			{"Custom", HHCustom, ""},
			{"invalid", HouseholdType("Unknown"), "household type is not recognized"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateHouseholdType(tc.input)
				if tc.want == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.want)
				}
			})
		}
	})

	t.Run("ValidateMemberRole", func(t *testing.T) {
		tests := []struct {
			name  string
			input MemberRole
			want  string
		}{
			{"Head", RoleHead, ""},
			{"Admin", RoleAdmin, ""},
			{"Member", RoleMember, ""},
			{"Viewer", RoleViewer, ""},
			{"invalid", MemberRole("SuperAdmin"), "member role is not recognized"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateMemberRole(tc.input)
				if tc.want == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.want)
				}
			})
		}
	})

	t.Run("ValidateCurrency", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
		}{
			{"USD", "USD", ""},
			{"INR", "INR", ""},
			{"empty", "", "currency must be a 3-letter ISO code"},
			{"two letters", "US", "currency must be a 3-letter ISO code"},
			{"four letters", "USDD", "currency must be a 3-letter ISO code"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateCurrency(tc.input)
				if tc.want == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.want)
				}
			})
		}
	})

	t.Run("ValidateMemberCount", func(t *testing.T) {
		tests := []struct {
			name  string
			input int
			want  string
		}{
			{"one member", 1, ""},
			{"50 members", 50, ""},
			{"zero members", 0, "household must have at least one member"},
			{"51 members", 51, "household cannot have more than 50 members"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				members := make([]HouseholdMember, tc.input)
				err := ValidateMemberCount(members)
				if tc.want == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.want)
				}
			})
		}
	})

	t.Run("ValidateHeadOfHousehold", func(t *testing.T) {
		head := makeHeadMember("u1")
		member := makeAcceptedMember("u2", RoleMember)
		admin := makeAcceptedMember("u3", RoleAdmin)

		tests := []struct {
			name    string
			members []HouseholdMember
			headID  string
			want    string
		}{
			{"valid head", []HouseholdMember{head, member}, "u1", ""},
			{"head not found", []HouseholdMember{member}, "u1", "head of household must be a member"},
			{"head has wrong role", []HouseholdMember{member}, "u2", "head of household must have Head role"},
			{"head is admin (wrong role)", []HouseholdMember{head, admin}, "u3", "head of household must have Head role"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateHeadOfHousehold(tc.members, tc.headID)
				if tc.want == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.want)
				}
			})
		}
	})
}

// ---------------------------------------------------------------------------
// Calculations
// ---------------------------------------------------------------------------

func TestCalculations(t *testing.T) {
	t.Parallel()

	t.Run("CalculateTotalNetWorth", func(t *testing.T) {
		tests := []struct {
			name      string
			assets    int64
			liabs     int64
			expected  int64
		}{
			{"positive", 10000, 3000, 7000},
			{"zero", 5000, 5000, 0},
			{"negative", 2000, 8000, -6000},
			{"all zero", 0, 0, 0},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expected, CalculateTotalNetWorth(tc.assets, tc.liabs))
			})
		}
	})

	t.Run("DetermineHouseholdHealth", func(t *testing.T) {
		tests := []struct {
			name       string
			netWorth   int64
			status     HouseholdStatus
			ratio      float64
			expected   HouseholdHealth
		}{
			{"healthy", 1000, HHStatusActive, 0.8, HHHealthy},
			{"critical negative net worth", -1, HHStatusActive, 0.8, HHCritical},
			{"critical dissolved", 1000, HHStatusDissolved, 0.8, HHCritical},
			{"critical archived", 1000, HHStatusArchived, 0.8, HHCritical},
			{"warning low ratio", 1000, HHStatusActive, 0.4, HHWarning},
			{"warning zero net worth", 0, HHStatusActive, 0.8, HHWarning},
			{"negative net worth supersedes low ratio", -1, HHStatusActive, 0.1, HHCritical},
			{"dissolved with negative net worth", -100, HHStatusDissolved, 0.9, HHCritical},
			{"paused healthy", 100, HHStatusPaused, 0.6, HHHealthy},
			{"paused zero net worth", 0, HHStatusPaused, 0.6, HHWarning},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got := DetermineHouseholdHealth(tc.netWorth, tc.status, tc.ratio)
				assert.Equal(t, tc.expected, got)
			})
		}
	})

	t.Run("CalculateActiveMemberRatio", func(t *testing.T) {
		tests := []struct {
			name          string
			total         int
			active        int
			expected      float64
		}{
			{"all active", 4, 4, 1.0},
			{"half active", 4, 2, 0.5},
			{"none active", 5, 0, 0},
			{"zero total", 0, 0, 0},
			{"one active of many", 10, 1, 0.1},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expected, CalculateActiveMemberRatio(tc.total, tc.active))
			})
		}
	})

	t.Run("DeriveDefaultCurrency", func(t *testing.T) {
		tests := []struct {
			country  string
			expected string
		}{
			{"IN", "INR"},
			{"India", "INR"},
			{"US", "USD"},
			{"USA", "USD"},
			{"United States", "USD"},
			{"GB", "GBP"},
			{"UK", "GBP"},
			{"United Kingdom", "GBP"},
			{"EU", "EUR"},
			{"Eurozone", "EUR"},
			{"AE", "AED"},
			{"UAE", "AED"},
			{"SG", "SGD"},
			{"Singapore", "SGD"},
			{"unknown", "INR"},
			{"", "INR"},
		}
		for _, tc := range tests {
			t.Run(tc.country, func(t *testing.T) {
				assert.Equal(t, tc.expected, DeriveDefaultCurrency(tc.country))
			})
		}
	})
}

// ---------------------------------------------------------------------------
// ReconstructFromDB
// ---------------------------------------------------------------------------

func TestReconstructFromDB(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	headID := "user-head"
	members := []HouseholdMember{
		{UserID: headID, Role: RoleHead, AddedAt: now.Format(time.RFC3339), InviteStatus: "Accepted"},
		{UserID: "user-admin", Role: RoleAdmin, AddedAt: now.Format(time.RFC3339), InviteStatus: "Accepted"},
	}
	accounts := []LinkedAccount{
		{AccountID: "acc-1", AddedBy: headID, AddedAt: now.Format(time.RFC3339)},
	}
	goals := []LinkedGoal{
		{GoalID: "goal-1", AddedBy: headID, AddedAt: now.Format(time.RFC3339)},
	}
	budgets := []LinkedBudget{
		{BudgetID: "budget-1", AddedBy: headID, AddedAt: now.Format(time.RFC3339)},
	}
	contribs := []GoalContribution{
		{GoalID: "goal-1", UserID: headID, Amount: 500, Date: "2026-01-01"},
	}

	h := ReconstructFromDB(
		"hh-recon", "Reconstructed Home", HHFamily, headID, members,
		HHStatusActive, "EUR", "EU",
		50000, 15000, 35000, HHHealthy,
		[]string{"tag1", "tag2"},
		accounts, goals, budgets, contribs,
		now, now,
	)

	assert.Equal(t, "hh-recon", h.ID())
	assert.Equal(t, "Reconstructed Home", h.Name())
	assert.Equal(t, HHFamily, h.HouseholdType())
	assert.Equal(t, headID, h.HeadOfHouseholdID())
	assert.Len(t, h.Members(), 2)
	assert.Equal(t, HHStatusActive, h.Status())
	assert.Equal(t, "EUR", h.Currency())
	assert.Equal(t, "EU", h.Country())
	assert.Equal(t, int64(50000), h.TotalAssets())
	assert.Equal(t, int64(15000), h.TotalLiabilities())
	assert.Equal(t, int64(35000), h.TotalNetWorth())
	assert.Equal(t, HHHealthy, h.Health())
	assert.Equal(t, []string{"tag1", "tag2"}, h.Tags())
	assert.Len(t, h.LinkedAccounts(), 1)
	assert.Len(t, h.LinkedGoals(), 1)
	assert.Len(t, h.LinkedBudgets(), 1)
	assert.Len(t, h.GoalContributions(), 1)
	assert.Equal(t, now.Unix(), h.CreatedAt().Unix())
	assert.Equal(t, now.Unix(), h.UpdatedAt().Unix())

	t.Run("nil slices become empty", func(t *testing.T) {
		h2 := ReconstructFromDB(
			"hh-nil", "Nil Test", HHSingle, headID, members,
			HHStatusDraft, "USD", "US",
			0, 0, 0, HHHealthy,
			nil, nil, nil, nil, nil,
			now, now,
		)
		assert.NotNil(t, h2.Tags())
		assert.Empty(t, h2.Tags())
		assert.NotNil(t, h2.LinkedAccounts())
		assert.Empty(t, h2.LinkedAccounts())
		assert.NotNil(t, h2.LinkedGoals())
		assert.Empty(t, h2.LinkedGoals())
		assert.NotNil(t, h2.LinkedBudgets())
		assert.Empty(t, h2.LinkedBudgets())
		assert.NotNil(t, h2.GoalContributions())
		assert.Empty(t, h2.GoalContributions())
	})
}

// ---------------------------------------------------------------------------
// Helpers — construct Households with known state
// ---------------------------------------------------------------------------

func testDraftHousehold() *Household {
	headID := "user-head"
	return ReconstructFromDB(
		"hh-id", "Test Household", HHCouple, headID,
		[]HouseholdMember{
			makeHeadMember(headID),
			makeAcceptedMember("user-admin", RoleAdmin),
			makeAcceptedMember("user-member", RoleMember),
			makeAcceptedMember("user-viewer", RoleViewer),
		},
		HHStatusDraft, "USD", "US",
		0, 0, 0, HHHealthy,
		[]string{}, nil, nil, nil, nil,
		time.Now().UTC(), time.Now().UTC(),
	)
}

func testActiveHousehold() *Household {
	headID := "user-head"
	return ReconstructFromDB(
		"hh-id", "Active Household", HHCouple, headID,
		[]HouseholdMember{
			makeHeadMember(headID),
			makeAcceptedMember("user-admin", RoleAdmin),
			makeAcceptedMember("user-member", RoleMember),
			makeAcceptedMember("user-viewer", RoleViewer),
		},
		HHStatusActive, "USD", "US",
		10000, 3000, 7000, HHHealthy,
		[]string{}, nil, nil, nil, nil,
		time.Now().UTC(), time.Now().UTC(),
	)
}

func testPausedHousehold() *Household {
	headID := "user-head"
	return ReconstructFromDB(
		"hh-id", "Paused Household", HHCouple, headID,
		[]HouseholdMember{
			makeHeadMember(headID),
		},
		HHStatusPaused, "USD", "US",
		5000, 1000, 4000, HHHealthy,
		[]string{}, nil, nil, nil, nil,
		time.Now().UTC(), time.Now().UTC(),
	)
}

func testDissolvedHousehold() *Household {
	headID := "user-head"
	return ReconstructFromDB(
		"hh-id", "Dissolved Household", HHSingle, headID,
		[]HouseholdMember{
			makeHeadMember(headID),
		},
		HHStatusDissolved, "USD", "US",
		0, 0, 0, HHCritical,
		[]string{}, nil, nil, nil, nil,
		time.Now().UTC(), time.Now().UTC(),
	)
}

func testArchivedHousehold() *Household {
	headID := "user-head"
	return ReconstructFromDB(
		"hh-id", "Archived Household", HHSingle, headID,
		[]HouseholdMember{
			makeHeadMember(headID),
		},
		HHStatusArchived, "USD", "US",
		0, 0, 0, HHCritical,
		[]string{}, nil, nil, nil, nil,
		time.Now().UTC(), time.Now().UTC(),
	)
}
