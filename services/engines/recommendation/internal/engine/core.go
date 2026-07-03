package engine

// Engine is the recommendation generation engine.
type Engine struct {
	config Config
}

// NewEngine creates a new recommendation engine.
func NewEngine(cfg Config) *Engine {
	return &Engine{config: cfg}
}

// Execute runs the 11-step recommendation pipeline.
func (e *Engine) Execute(inputs Inputs) ([]Recommendation, error) {
	// Step 1-2: Freeze inputs
	frozen := inputs

	// Step 3-4: Identify gaps and opportunities
	candidates := e.identifyGaps(frozen)
	candidates = append(candidates, e.identifyOpportunities(frozen)...)

	// Step 5-6: Score, rank, and prioritise
	for i := range candidates {
		candidates[i].Scores = e.scoreRecommendation(candidates[i], frozen)
		candidates[i].Score = e.calculateComposite(candidates[i].Scores)
	}
	candidates = e.rankAndPrioritise(candidates)

	// Step 7-8: Generate evidence and explanations
	for i := range candidates {
		candidates[i].Evidence = e.generateEvidence(candidates[i], frozen)
		candidates[i].Explanation = e.generateExplanation(candidates[i], frozen)
		candidates[i].Status = RSGenerated
	}

	if len(candidates) > e.config.MaxTotalRecommendations {
		candidates = candidates[:e.config.MaxTotalRecommendations]
	}
	return candidates, nil
}

// Inputs holds all domain data consumed by the engine.
type Inputs struct {
	GoalGaps        map[string]int64  `json:"goal_gaps"`        // goalName -> gap
	GoalImportance  map[string]string `json:"goal_importance"`  // goalName -> Mandatory/Essential/...
	Debts           []DebtInfo        `json:"debts"`
	PortfolioDrift  float64           `json:"portfolio_drift"`
	EmergencyMonths float64           `json:"emergency_months"`
	CashFlowSurplus int64             `json:"cash_flow_surplus"`
	SavingsRate     float64           `json:"savings_rate"`
	InsuranceGap    float64           `json:"insurance_gap"`
	HasTaxOpportunity bool            `json:"has_tax_opportunity"`
	OverspentCategories map[string]int64 `json:"overspent_categories"`
}

// DebtInfo holds debt information for analysis.
type DebtInfo struct {
	Name         string  `json:"name"`
	Balance      int64   `json:"balance"`
	InterestRate float64 `json:"interest_rate"`
}
