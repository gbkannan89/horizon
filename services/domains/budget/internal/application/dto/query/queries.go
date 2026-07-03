package query

type GetBudgetQuery struct {
	BudgetID string `json:"budget_id" validate:"required"`
}

type ListByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type ListByPeriodQuery struct {
	UserID    string `json:"user_id" validate:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

type ListByHouseholdQuery struct {
	HouseholdID string `json:"household_id" validate:"required"`
	Cursor      string `json:"cursor,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

type GetByCategoryQuery struct {
	UserID   string `json:"user_id" validate:"required"`
	Category string `json:"category" validate:"required"`
	Cursor   string `json:"cursor,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type BudgetVsActualQuery struct {
	BudgetID string `json:"budget_id" validate:"required"`
}

type BudgetView struct {
	BudgetID      string            `json:"budget_id"`
	HouseholdID   string            `json:"household_id,omitempty"`
	Name          string            `json:"name"`
	Period        string            `json:"period"`
	StartDate     string            `json:"start_date"`
	EndDate       string            `json:"end_date"`
	Status        string            `json:"status"`
	TotalBudgeted int64             `json:"total_budgeted"`
	TotalSpent    int64             `json:"total_spent"`
	TotalRemaining int64            `json:"total_remaining"`
	Currency      string            `json:"currency"`
	Categories    []CategoryView    `json:"categories"`
	Tags          []string          `json:"tags"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

type CategoryView struct {
	ID              string `json:"id"`
	Category        string `json:"category"`
	Subcategory     string `json:"subcategory,omitempty"`
	BudgetedAmount  int64  `json:"budgeted_amount"`
	SpentAmount     int64  `json:"spent_amount"`
	RemainingAmount int64  `json:"remaining_amount"`
	Rollover        bool   `json:"rollover"`
	SpentPct        float64 `json:"spent_pct"`
}

type PaginatedResult struct {
	Budgets   []BudgetView `json:"budgets"`
	NextCursor string       `json:"next_cursor,omitempty"`
	HasMore    bool         `json:"has_more"`
}
