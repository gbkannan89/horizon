package engine

// ProjectionType represents the type of projection.
type ProjectionType string
const (
	PTCashFlow            ProjectionType = "CashFlow"
	PTNetWorth            ProjectionType = "NetWorth"
	PTGoalCompletion      ProjectionType = "GoalCompletion"
	PTDebtPayoff          ProjectionType = "DebtPayoff"
	PTInvestmentGrowth    ProjectionType = "InvestmentGrowth"
	PTRetirement          ProjectionType = "Retirement"
	PTFinancialIndependence ProjectionType = "FinancialIndependence"
	PTPortfolioValue      ProjectionType = "PortfolioValue"
)

// EngineStatus represents the lifecycle state.
type EngineStatus string
const (ESRequested EngineStatus = "Requested"; ESQueued EngineStatus = "Queued"; ESRunning EngineStatus = "Running"; ESCompleted EngineStatus = "Completed"; ESFailed EngineStatus = "Failed"; ESExpired EngineStatus = "Expired"; ESSuperseded EngineStatus = "Superseded")

// Inputs holds all data needed for a projection.
type Inputs struct {
	NetWorth        int64             `json:"net_worth"`
	TotalAssets     int64             `json:"total_assets"`
	TotalLiabilities int64            `json:"total_liabilities"`
	MonthlyIncome   int64             `json:"monthly_income"`
	MonthlyExpenses int64             `json:"monthly_expenses"`
	CashBalance     int64             `json:"cash_balance"`
	GoalValues      map[string]int64  `json:"goal_values"`
	GoalTargets     map[string]int64  `json:"goal_targets"`
	AssetValues     []int64           `json:"asset_values"`
	Liabilities     []int64           `json:"liabilities"`
}

// Snapshot represents a single point in the projection timeline.
type Snapshot struct {
	Date            string `json:"date"`
	NetWorth        int64  `json:"net_worth,omitempty"`
	NetCashFlow     int64  `json:"net_cash_flow,omitempty"`
	CashFlow        int64  `json:"cash_flow,omitempty"`
	GoalValue       int64  `json:"goal_value,omitempty"`
	DebtValue       int64  `json:"debt_value,omitempty"`
	InvValue        int64  `json:"investment_value,omitempty"`
	PfValue         int64  `json:"portfolio_value,omitempty"`
	TotalAssets     int64  `json:"total_assets,omitempty"`
	TotalLiabilities int64 `json:"total_liabilities,omitempty"`
}

// Milestone represents a key event detected in the projection.
type Milestone struct {
	Date        string `json:"date"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// SummaryMetrics contains computed summary statistics.
type SummaryMetrics struct {
	FinalNetWorth       int64   `json:"final_net_worth,omitempty"`
	TotalGrowth         int64   `json:"total_growth,omitempty"`
	AnnualizedReturn    float64 `json:"annualized_return,omitempty"`
	AvgMonthlyCashFlow  int64   `json:"avg_monthly_cash_flow,omitempty"`
	TotalCashFlow       int64   `json:"total_cash_flow,omitempty"`
	FinalDebtBalance    int64   `json:"final_debt_balance,omitempty"`
	DebtFreeDate        string  `json:"debt_free_date,omitempty"`
	GoalCompletionRate  float64 `json:"goal_completion_rate,omitempty"`
	RetirementReadiness string  `json:"retirement_readiness,omitempty"`
}

// Assumptions holds configurable projection parameters.
type Assumptions struct {
	InflationRate          float64 `json:"inflation_rate"`
	EquityReturn           float64 `json:"equity_return"`
	DebtReturn             float64 `json:"debt_return"`
	SalaryGrowthRate       float64 `json:"salary_growth_rate"`
	ExpenseGrowthRate      float64 `json:"expense_growth_rate"`
	ContributionGrowthRate float64 `json:"contribution_growth_rate"`
	LoanInterestRate       float64 `json:"loan_interest_rate"`
	TaxRate                float64 `json:"tax_rate"`
	RetirementAge          int     `json:"retirement_age"`
	LifeExpectancy         int     `json:"life_expectancy"`
}

// DefaultAssumptions returns sensible defaults.
func DefaultAssumptions() Assumptions {
	return Assumptions{
		InflationRate: 4.0, EquityReturn: 10.0, DebtReturn: 7.0,
		SalaryGrowthRate: 8.0, ExpenseGrowthRate: 4.0, ContributionGrowthRate: 5.0,
		LoanInterestRate: 9.0, TaxRate: 20.0, RetirementAge: 60, LifeExpectancy: 85,
	}
}

// ProjectionOutput is the complete result of a projection computation.
type ProjectionOutput struct {
	OutputID        string          `json:"output_id"`
	ProjectionType  ProjectionType  `json:"projection_type"`
	Version         int             `json:"version"`
	Status          EngineStatus    `json:"status"`
	Timeline        []Snapshot      `json:"timeline"`
	Milestones      []Milestone     `json:"milestones"`
	SummaryMetrics  SummaryMetrics  `json:"summary_metrics"`
	Assumptions     Assumptions     `json:"assumptions"`
	ConfidenceLevel string          `json:"confidence_level"`
}
