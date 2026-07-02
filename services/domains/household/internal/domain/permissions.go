package domain

// HasAccess returns true if the user has at least Viewer access to the household.
func (h *Household) HasAccess(userID string) bool {
	if h.headOfHouseholdID == userID {
		return true
	}
	for _, m := range h.members {
		if m.UserID == userID && m.InviteStatus == "Accepted" {
			return true
		}
	}
	return false
}

// CanView returns true if the user has access to view household data.
func (h *Household) CanView(userID string) bool {
	return h.HasAccess(userID)
}

// CanEdit returns true if the user is the Head or an Admin.
func (h *Household) CanEdit(userID string) bool {
	if h.headOfHouseholdID == userID {
		return true
	}
	for _, m := range h.members {
		if m.UserID == userID && m.InviteStatus == "Accepted" && m.Role == RoleAdmin {
			return true
		}
	}
	return false
}

// CanTransition returns true if the user can change household state.
func (h *Household) CanTransition(userID string) bool {
	return h.CanEdit(userID)
}

// CanLink returns true if the user can link/unlink accounts, goals, budgets.
func (h *Household) CanLink(userID string) bool {
	return h.CanEdit(userID)
}

// CanManageMembers returns true if the user is the Head or an Admin.
func (h *Household) CanManageMembers(userID string) bool {
	return h.CanEdit(userID)
}

// GetMemberRole returns the role of the user in the household.
func (h *Household) GetMemberRole(userID string) (MemberRole, bool) {
	if h.headOfHouseholdID == userID {
		return RoleHead, true
	}
	for _, m := range h.members {
		if m.UserID == userID && m.InviteStatus == "Accepted" {
			return m.Role, true
		}
	}
	return "", false
}
