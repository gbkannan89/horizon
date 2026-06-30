package command

type CreatePortfolioCommand struct {
	OwnerID          string `json:"owner_id" validate:"required"`
	Name             string `json:"name" validate:"required"`
	PortfolioType    string `json:"portfolio_type" validate:"required"`
	BaseCurrency     string `json:"base_currency" validate:"required"`
	MembershipModel  string `json:"membership_model"`
	RiskProfile      string `json:"risk_profile"`
	BenchmarkProfile string `json:"benchmark_profile"`
	Benchmark        string `json:"benchmark,omitempty"`
	SourceOfTruth    string `json:"source_of_truth"`
	Tags             []string `json:"tags,omitempty"`
	Notes            string `json:"notes,omitempty"`
}

type IDCmd struct{ PortfolioID string `json:"portfolio_id" validate:"required"` }
type PfResult struct {
	PortfolioID string `json:"portfolio_id"`
	Status      string `json:"status"`
	Success     bool   `json:"success"`
}
