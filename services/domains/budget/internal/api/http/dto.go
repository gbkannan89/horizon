package http

type CreateBudgetRequest struct {
	Name       string `json:"name"`
	Period     string `json:"period"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Currency   string `json:"currency"`
	Categories []CategoryRequest `json:"categories"`
}

type CategoryRequest struct {
	Category       string `json:"category"`
	BudgetedAmount int64  `json:"budgeted_amount"`
	Rollover       bool   `json:"rollover"`
}
