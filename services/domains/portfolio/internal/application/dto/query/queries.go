package query

type GetPfQuery struct{ PortfolioID string `json:"portfolio_id" validate:"required"` }
type ListByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type ListByTypeQuery struct {
	UserID string `json:"user_id" validate:"required"`
	PfType string `json:"portfolio_type" validate:"required"`
	Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}

type PfView struct {
	PortfolioID   string `json:"portfolio_id"`
	Name          string `json:"name"`
	PortfolioType string `json:"portfolio_type"`
	BaseCurrency  string `json:"base_currency"`
	Status        string `json:"status"`
	CurrentValue  int64  `json:"current_value"`
	CostBasis     int64  `json:"cost_basis"`
	Health        string `json:"health"`
	RiskProfile   string `json:"risk_profile"`
	MemberCount   int    `json:"member_count"`
	CreatedAt     string `json:"created_at"`
}

type PaginatedResult struct {
	Portfolios []PfView `json:"portfolios"`
	NextCursor string   `json:"next_cursor,omitempty"`
	HasMore    bool     `json:"has_more"`
}
