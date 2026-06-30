package query

type GetLiabilityQuery struct{ LiabilityID string `json:"liability_id" validate:"required"` }
type ListByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}

type LiabView struct {
	LiabilityID         string  `json:"liability_id"`
	Name                string  `json:"name"`
	Classification      string  `json:"classification"`
	Currency            string  `json:"currency"`
	OriginalPrincipal   int64   `json:"original_principal"`
	OutstandingBalance  int64   `json:"outstanding_balance"`
	InterestRate        float64 `json:"interest_rate"`
	InterestModel       string  `json:"interest_model"`
	Status              string  `json:"status"`
	RemainingInstallments int   `json:"remaining_installments"`
	MaturityDate        string  `json:"maturity_date"`
	Health              string  `json:"health"`
	CreatedAt           string  `json:"created_at"`
}

type PaginatedResult struct {
	Liabilities []LiabView `json:"liabilities"`
	NextCursor  string     `json:"next_cursor,omitempty"`
	HasMore     bool       `json:"has_more"`
}
