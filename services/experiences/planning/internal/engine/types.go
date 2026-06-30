package engine

// CardType represents the types of planning workspace cards.
type CardType string
const (
	CTCurrentPlan     CardType = "CurrentPlan"
	CTFutureProjection CardType = "FutureProjection"
	CTAltStrategy     CardType = "AlternativeStrategy"
	CTSimResult       CardType = "SimulationResult"
	CTOptResult       CardType = "OptimizationResult"
	CTRiskImpact      CardType = "RiskImpact"
	CTGoalTradeoff    CardType = "GoalTradeoff"
	CTFundingGap      CardType = "FundingGap"
	CTDecisionSummary CardType = "DecisionSummary"
)

// PlanningCard represents a single card in the planning workspace.
type PlanningCard struct {
	CardID    string   `json:"card_id"`
	CardType  CardType `json:"card_type"`
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Detail    string   `json:"detail"`
	Priority  int      `json:"priority"`
}

// PlanningMode represents the focus mode.
type PlanningMode string
const (
	PMOverview          PlanningMode = "Overview"
	PMGoalPlanning      PlanningMode = "GoalPlanning"
	PMRetirement        PlanningMode = "RetirementPlanning"
	PMEducationPlanning PlanningMode = "EducationPlanning"
	PMDebtPlanning      PlanningMode = "DebtPlanning"
	PMInvestmentPlan    PlanningMode = "InvestmentPlanning"
	PMScenarioPlan      PlanningMode = "ScenarioPlanning"
	PMOptimization      PlanningMode = "Optimization"
	PMHistoricalCompare PlanningMode = "HistoricalComparison"
)

// ComparisonDelta represents the difference between two plans for a single metric.
type ComparisonDelta struct {
	Metric       string `json:"metric"`
	Baseline     string `json:"baseline"`
	Alternative  string `json:"alternative"`
	Delta        string `json:"delta"`
	Direction    string `json:"direction"` // better, worse, unchanged
}

// PlanSummary represents a saved plan.
type PlanSummary struct {
	PlanID         string `json:"plan_id"`
	Name           string `json:"name"`
	IsBaseline     bool   `json:"is_baseline"`
	NetWorthAtEnd  int64  `json:"net_worth_at_end"`
	GoalsOnTrack   int    `json:"goals_on_track"`
	TotalGoals     int    `json:"total_goals"`
	FundingGap     int64  `json:"funding_gap"`
	RiskScore      int    `json:"risk_score"`
}

// PlanningWorkspace is the complete planning view.
type PlanningWorkspace struct {
	Mode            PlanningMode      `json:"mode"`
	Cards           []PlanningCard    `json:"cards"`
	CurrentPlan     *PlanSummary      `json:"current_plan,omitempty"`
	AlternativePlans []PlanSummary    `json:"alternative_plans,omitempty"`
}

// PlanComparison is the side-by-side comparison result.
type PlanComparison struct {
	PlanA      PlanSummary       `json:"plan_a"`
	PlanB      PlanSummary       `json:"plan_b"`
	Deltas     []ComparisonDelta `json:"deltas"`
}

// Inputs contains all data consumed by the planning experience.
type Inputs struct {
	Mode             PlanningMode    `json:"mode"`
	NetWorthCurrent  int64           `json:"net_worth_current"`
	NetWorthProjected int64          `json:"net_worth_projected"`
	GoalsOnTrack     int             `json:"goals_on_track"`
	TotalGoals       int             `json:"total_goals"`
	FundingGapTotal  int64           `json:"funding_gap_total"`
	RiskScore        int             `json:"risk_score"`
	HasProjection    bool            `json:"has_projection"`
	HasSimulation    bool            `json:"has_simulation"`
	HasOptimization  bool            `json:"has_optimization"`
	HasAlternatives  bool            `json:"has_alternatives"`
	UserID           string          `json:"user_id"`
	Assumptions      AssumptionsData `json:"assumptions"`
}

// AssumptionsData represents the configurable planning assumptions.
type AssumptionsData struct {
	EquityReturn    float64 `json:"equity_return"`
	InflationRate   float64 `json:"inflation_rate"`
	ContributionGrowth float64 `json:"contribution_growth"`
	RetirementAge   int     `json:"retirement_age"`
}
