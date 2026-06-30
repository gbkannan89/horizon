package engine

// RecStatus represents the recommendation lifecycle state.
type RecStatus string
const (
	RSDraft     RecStatus = "Draft"
	RSGenerated RecStatus = "Generated"
	RSActive    RecStatus = "Active"
	RSAccepted  RecStatus = "Accepted"
	RSRejected  RecStatus = "Rejected"
	RSExpired   RecStatus = "Expired"
	RSArchived  RecStatus = "Archived"
	RSHistorical RecStatus = "Historical"
)

// RecCategory represents the category of recommendation.
type RecCategory string
const (
	RCSavings              RecCategory = "Savings"
	RCDebt                 RecCategory = "Debt"
	RCInvestment           RecCategory = "Investment"
	RCGoalFunding          RecCategory = "GoalFunding"
	RCEmergencyFund        RecCategory = "EmergencyFund"
	RCInsurance            RecCategory = "Insurance"
	RCCashFlow             RecCategory = "CashFlow"
	RCExpenseReduction     RecCategory = "ExpenseReduction"
	RCIncomeGrowth         RecCategory = "IncomeGrowth"
	RCPortfolioRebalancing RecCategory = "PortfolioRebalancing"
	RCTax                  RecCategory = "Tax"
	RCRetirement           RecCategory = "Retirement"
	RCBehaviour            RecCategory = "Behaviour"
	RCFinancialDiscipline  RecCategory = "FinancialDiscipline"
)

// ScoringDimensions holds the 6 scoring factors.
type ScoringDimensions struct {
	Impact              float64 `json:"impact"`
	Urgency             float64 `json:"urgency"`
	Difficulty          float64 `json:"difficulty"`
	Confidence          float64 `json:"confidence"`
	GoalAlignment       float64 `json:"goal_alignment"`
	FinancialHealthImpact float64 `json:"financial_health_impact"`
}

// Recommendation is a single generated recommendation.
type Recommendation struct {
	ID                  string            `json:"id"`
	Category            RecCategory       `json:"category"`
	Priority            int               `json:"priority"`
	Score               float64           `json:"score"`
	ImpactRating        string            `json:"impact_rating"`
	UrgencyRating       string            `json:"urgency_rating"`
	DifficultyRating    string            `json:"difficulty_rating"`
	ConfidenceRating    string            `json:"confidence_rating"`
	Title               string            `json:"title"`
	Summary             string            `json:"summary"`
	Action              string            `json:"action"`
	ExpectedImprovement string            `json:"expected_improvement"`
	Evidence            []string          `json:"evidence"`
	Explanation         string            `json:"explanation"`
	Tradeoffs           []string          `json:"tradeoffs"`
	ExpirationTime      string            `json:"expiration_time"`
	Status              RecStatus         `json:"status"`
	Scores              ScoringDimensions `json:"scores"`
	CreatedAt           string            `json:"created_at,omitempty"`
}
