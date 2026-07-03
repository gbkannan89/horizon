package command

type CreateBudgetCommand struct {
	UserID      string `json:"user_id" validate:"required"`
	HouseholdID string `json:"household_id,omitempty"`
	Name        string `json:"name" validate:"required"`
	Period      string `json:"period" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
	Categories  []BudgetCategoryCmd `json:"categories,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
}

type BudgetCategoryCmd struct {
	Category       string `json:"category" validate:"required"`
	Subcategory    string `json:"subcategory,omitempty"`
	BudgetedAmount int64  `json:"budgeted_amount"`
	Rollover       bool   `json:"rollover"`
}

type BudgetIDCommand struct {
	BudgetID string `json:"budget_id" validate:"required"`
}

type UpdateCategoryCommand struct {
	BudgetID       string `json:"budget_id" validate:"required"`
	CategoryID     string `json:"category_id" validate:"required"`
	BudgetedAmount *int64 `json:"budgeted_amount,omitempty"`
	SpentAmount    *int64 `json:"spent_amount,omitempty"`
}

type BudgetResult struct {
	BudgetID string `json:"budget_id"`
	Status   string `json:"status"`
	Success  bool   `json:"success"`
}
