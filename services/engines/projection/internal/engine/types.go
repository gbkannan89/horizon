package engine

// ProjectionType represents the type of projection computation.
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
const (
	ESRequested  EngineStatus = "Requested"
	ESQueued     EngineStatus = "Queued"
	ESRunning    EngineStatus = "Running"
	ESCompleted  EngineStatus = "Completed"
	ESFailed     EngineStatus = "Failed"
	ESExpired    EngineStatus = "Expired"
	ESSuperseded EngineStatus = "Superseded"
)

// Snapshot represents a single point in the projection timeline.
type Snapshot struct {
	Date      string `json:"date"`
	NetWorth  int64  `json:"net_worth,omitempty"`
	CashFlow  int64  `json:"cash_flow,omitempty"`
	GoalValue int64  `json:"goal_value,omitempty"`
	DebtValue int64  `json:"debt_value,omitempty"`
	InvValue  int64  `json:"investment_value,omitempty"`
	PfValue   int64  `json:"portfolio_value,omitempty"`
}

// Milestone represents a detected key event.
type Milestone struct {
	Date        string `json:"date"`
	MilestoneType string `json:"milestone_type"`
	Description string `json:"description"`
}

// ProjectionOutput is the result of a projection computation.
type ProjectionOutput struct {
	OutputID        string           `json:"output_id"`
	ProjectionType  ProjectionType   `json:"projection_type"`
	Version         int              `json:"version"`
	Status          EngineStatus     `json:"status"`
	Timeline        []Snapshot       `json:"timeline"`
	Milestones      []Milestone      `json:"milestones"`
	InputVersions   map[string]int   `json:"input_versions"`
	Assumptions     Assumptions      `json:"assumptions"`
	Confidence      string           `json:"confidence"`
}

// Assumptions holds the 12 configurable projection assumptions.
type Assumptions struct {
	InflationRate           float64 `json:"inflation_rate"`
	EquityReturn            float64 `json:"equity_return"`
	DebtReturn              float64 `json:"debt_return"`
	RealEstateAppreciation  float64 `json:"real_estate_appreciation"`
	SalaryGrowthRate        float64 `json:"salary_growth_rate"`
	ContributionGrowthRate  float64 `json:"contribution_growth_rate"`
	ExpenseGrowthRate       float64 `json:"expense_growth_rate"`
	LoanInterestRate        float64 `json:"loan_interest_rate"`
	TaxRate                 float64 `json:"tax_rate"`
	RetirementAge           int      `json:"retirement_age"`
	LifeExpectancy          int      `json:"life_expectancy"`
	EmergencyFundTargetMonths int    `json:"emergency_fund_target_months"`
}

// DefaultAssumptions returns the default assumption set.
func DefaultAssumptions() Assumptions {
	return Assumptions{
		InflationRate: 4.0, EquityReturn: 10.0, DebtReturn: 7.0,
		RealEstateAppreciation: 5.0, SalaryGrowthRate: 8.0,
		ContributionGrowthRate: 5.0, ExpenseGrowthRate: 4.0,
		LoanInterestRate: 9.0, TaxRate: 20.0, RetirementAge: 60,
		LifeExpectancy: 85, EmergencyFundTargetMonths: 6,
	}
}
