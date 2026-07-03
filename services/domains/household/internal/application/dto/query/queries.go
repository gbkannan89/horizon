package query

type GetHouseholdQuery struct {
	HouseholdID string `json:"household_id" validate:"required"`
}

type ListByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type ListByStatusQuery struct {
	Status string `json:"status" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type HouseholdView struct {
	HouseholdID       string   `json:"household_id"`
	Name              string   `json:"name"`
	HouseholdType     string   `json:"household_type"`
	HeadOfHouseholdID string   `json:"head_of_household_id"`
	MemberCount       int      `json:"member_count"`
	Status            string   `json:"status"`
	Currency          string   `json:"currency"`
	Country           string   `json:"country"`
	TotalAssets       int64    `json:"total_assets"`
	TotalLiabilities  int64    `json:"total_liabilities"`
	TotalNetWorth     int64    `json:"total_net_worth"`
	Health            string   `json:"health"`
	Tags              []string `json:"tags"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	LinkedAccounts    []LinkedAccountView `json:"linked_accounts"`
	LinkedGoals       []LinkedGoalView    `json:"linked_goals"`
	LinkedBudgets     []LinkedBudgetView  `json:"linked_budgets"`
}

type LinkedAccountView struct {
	AccountID string `json:"account_id"`
	AddedBy   string `json:"added_by"`
	AddedAt   string `json:"added_at"`
}

type LinkedGoalView struct {
	GoalID  string `json:"goal_id"`
	AddedBy string `json:"added_by"`
	AddedAt string `json:"added_at"`
}

type LinkedBudgetView struct {
	BudgetID string `json:"budget_id"`
	AddedBy  string `json:"added_by"`
	AddedAt  string `json:"added_at"`
}

type GoalContributionView struct {
	GoalID string `json:"goal_id"`
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
	Date   string `json:"date"`
}

type HouseholdDetailView struct {
	HouseholdView
	Members           []HouseholdMemberView  `json:"members"`
	Notes             string                 `json:"notes"`
	GoalContributions []GoalContributionView `json:"goal_contributions"`
}

type HouseholdMemberView struct {
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	AddedAt      string `json:"added_at"`
	InviteStatus string `json:"invite_status"`
}

type PaginatedResult struct {
	Households []HouseholdView `json:"households"`
	NextCursor string          `json:"next_cursor,omitempty"`
	HasMore    bool            `json:"has_more"`
}

type HouseholdFinancialSummary struct {
	HouseholdID      string `json:"household_id"`
	TotalAssets      int64  `json:"total_assets"`
	TotalLiabilities int64  `json:"total_liabilities"`
	TotalNetWorth    int64  `json:"total_net_worth"`
	Currency         string `json:"currency"`
	LinkedAccounts   int    `json:"linked_accounts_count"`
	LinkedGoals      int    `json:"linked_goals_count"`
	LinkedBudgets    int    `json:"linked_budgets_count"`
	TotalContributions int64 `json:"total_contributions"`
}
