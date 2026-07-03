package query

type GetAccountQuery struct{ AccountID string `json:"account_id" validate:"required"` }
type ListByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type ListByTypeQuery struct {
	UserID    string `json:"user_id" validate:"required"`
	AcctType  string `json:"account_type" validate:"required"`
	Cursor    string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type ListByStatusQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Status string `json:"status" validate:"required"`
	Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type ListByHouseholdQuery struct {
	HouseholdID string `json:"household_id" validate:"required"`
	Cursor      string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type AccountView struct {
	AccountID      string `json:"account_id"`
	AccountName    string `json:"account_name"`
	AccountType    string `json:"account_type"`
	Classification string `json:"classification"`
	Currency       string `json:"currency"`
	Status         string `json:"status"`
	OwnerID        string `json:"owner_id"`
	HouseholdID    string `json:"household_id,omitempty"`
	Visibility     string `json:"visibility"`
	Liquidity      string `json:"liquidity_profile"`
	Health         string `json:"account_health"`
	CreatedAt      string `json:"created_at"`
}
type PaginatedResult struct {
	Accounts   []AccountView `json:"accounts"`
	NextCursor string        `json:"next_cursor,omitempty"`
	HasMore    bool          `json:"has_more"`
}
