package command

type CreateHouseholdCommand struct {
	Name             string   `json:"name" validate:"required"`
	HouseholdType    string   `json:"household_type" validate:"required"`
	HeadOfHouseholdID string  `json:"head_of_household_id" validate:"required"`
	Currency         string   `json:"currency" validate:"required"`
	Country          string   `json:"country"`
	Tags             []string `json:"tags,omitempty"`
	Notes            string   `json:"notes,omitempty"`
}

type HouseholdIDCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	RequesterID string `json:"requester_id"`
}

type AddMemberCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	RequesterID string `json:"requester_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
	Role        string `json:"role" validate:"required"`
}

type AcceptInviteCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
}

type RemoveMemberCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	RequesterID string `json:"requester_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
}

type UpdateMemberRoleCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	RequesterID string `json:"requester_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
	Role        string `json:"role" validate:"required"`
}

type HouseholdResult struct {
	HouseholdID string `json:"household_id"`
	Status      string `json:"status"`
	Success     bool   `json:"success"`
}

type LinkAccountCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	AccountID   string `json:"account_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
	RequesterID string `json:"requester_id"`
}

type UnlinkAccountCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	AccountID   string `json:"account_id" validate:"required"`
	RequesterID string `json:"requester_id"`
}

type LinkGoalCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	GoalID      string `json:"goal_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
	RequesterID string `json:"requester_id"`
}

type UnlinkGoalCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	GoalID      string `json:"goal_id" validate:"required"`
	RequesterID string `json:"requester_id"`
}

type LinkBudgetCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	BudgetID    string `json:"budget_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
	RequesterID string `json:"requester_id"`
}

type UnlinkBudgetCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	BudgetID    string `json:"budget_id" validate:"required"`
	RequesterID string `json:"requester_id"`
}

type AddGoalContributionCommand struct {
	HouseholdID string `json:"household_id" validate:"required"`
	GoalID      string `json:"goal_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
	Amount      int64  `json:"amount" validate:"required"`
	RequesterID string `json:"requester_id"`
}
