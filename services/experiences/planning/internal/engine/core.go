package engine

import "fmt"

type CardType string

const (
	CTPlanOverview    CardType = "PlanOverview"
	CTGoalPlanning    CardType = "GoalPlanning"
	CTCashFlowPlan    CardType = "CashFlowPlanning"
	CTRetirementPlan  CardType = "RetirementPlanning"
	CTInvestmentPlan  CardType = "InvestmentPlanning"
	CTDebtPlan        CardType = "DebtPlanning"
	CTProjection      CardType = "Projection"
	CTRiskSum         CardType = "Risk"
	CTRecommend       CardType = "Recommendation"
	CTOptimization    CardType = "Optimization"
	CTSimulation      CardType = "Simulation"
	CTTimeline        CardType = "Timeline"
)

type Card struct {
	CardID   string      `json:"card_id"`
	CardType CardType    `json:"card_type"`
	Title    string      `json:"title"`
	Summary  string      `json:"summary"`
	Data     interface{} `json:"data,omitempty"`
	Priority int         `json:"priority"`
}

type PlanningMode string

const (
	PMOverview         PlanningMode = "Overview"
	PMGoalPlanning     PlanningMode = "GoalPlanning"
	PMRetirement       PlanningMode = "RetirementPlanning"
	PMEducationPlan    PlanningMode = "EducationPlanning"
	PMDebtPlan         PlanningMode = "DebtPlanning"
	PMInvestmentPlan   PlanningMode = "InvestmentPlanning"
	PMScenarioPlan     PlanningMode = "ScenarioPlanning"
	PMOptimize         PlanningMode = "Optimization"
	PMHistoricalComp   PlanningMode = "HistoricalComparison"
)

type PlanSummary struct {
	PlanID        string `json:"plan_id"`
	Name          string `json:"name"`
	IsBaseline    bool   `json:"is_baseline"`
	NetWorthEnd   int64  `json:"net_worth_at_end"`
	GoalsOnTrack  int    `json:"goals_on_track"`
	TotalGoals    int    `json:"total_goals"`
	FundingGap    int64  `json:"funding_gap"`
	RiskScore     int    `json:"risk_score"`
	HealthScore   int    `json:"health_score,omitempty"`
}

type ComparisonDelta struct {
	Metric      string `json:"metric"`
	Baseline    string `json:"baseline"`
	Alternative string `json:"alternative"`
	Delta       string `json:"delta"`
	Direction   string `json:"direction"`
}

type PlanComparison struct {
	PlanA  PlanSummary       `json:"plan_a"`
	PlanB  PlanSummary       `json:"plan_b"`
	Deltas []ComparisonDelta `json:"deltas"`
}

type Scenario struct {
	ScenarioID  string          `json:"scenario_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Assumptions AssumptionsData `json:"assumptions"`
	Projections ProjectionData  `json:"projections"`
}

type AssumptionsData struct {
	EquityReturn       float64 `json:"equity_return"`
	InflationRate      float64 `json:"inflation_rate"`
	ContributionGrowth float64 `json:"contribution_growth"`
	RetirementAge      int     `json:"retirement_age"`
}

type ProjectionData struct {
	NetWorthProjected int64   `json:"net_worth_projected"`
	IncomeProjected   int64   `json:"income_projected"`
	ExpensesProjected int64   `json:"expenses_projected"`
	PortfolioValue    int64   `json:"portfolio_value"`
	Confidence        string  `json:"confidence"`
}

type PlanningDashboard struct {
	Mode        PlanningMode `json:"mode"`
	Cards       []Card       `json:"cards"`
	CurrentPlan *PlanSummary `json:"current_plan,omitempty"`
	Scenarios   []Scenario   `json:"scenarios,omitempty"`
}

type Inputs struct {
	Mode              PlanningMode    `json:"mode"`
	NetWorthCurrent   int64           `json:"net_worth_current"`
	NetWorthProjected int64           `json:"net_worth_projected"`
	IncomeProjected   int64           `json:"income_projected"`
	ExpensesProjected int64           `json:"expenses_projected"`
	GoalsOnTrack      int             `json:"goals_on_track"`
	TotalGoals        int             `json:"total_goals"`
	FundingGap        int64           `json:"funding_gap"`
	RiskScore         int             `json:"risk_score"`
	RiskLevel         string          `json:"risk_level"`
	HealthScore       int             `json:"health_score"`
	HealthGrade       string          `json:"health_grade"`
	PortfolioValue    int64           `json:"portfolio_value"`
	CashBalance       int64           `json:"cash_balance"`
	MonthlyIncome     int64           `json:"monthly_income"`
	MonthlyExpenses   int64           `json:"monthly_expenses"`
	HasProjection     bool            `json:"has_projection"`
	HasSimulation     bool            `json:"has_simulation"`
	HasOptimization   bool            `json:"has_optimization"`
	HasRecs           bool            `json:"has_recommendations"`
	HasAlternatives   bool            `json:"has_alternatives"`
	RecCount          int             `json:"rec_count"`
	OptCount          int             `json:"opt_count"`
	SimCount          int             `json:"sim_count"`
	EventCount        int             `json:"event_count"`
	MilestoneCount    int             `json:"milestone_count"`
	UserID            string          `json:"user_id"`
	Assumptions       AssumptionsData `json:"assumptions"`
	Scenarios         []Scenario      `json:"scenarios"`
}

type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

func (c *Composer) BuildDashboard(inputs Inputs) *PlanningDashboard {
	plan := &PlanSummary{
		PlanID: "baseline", Name: "Current Plan", IsBaseline: true,
		NetWorthEnd: inputs.NetWorthProjected,
		GoalsOnTrack: inputs.GoalsOnTrack, TotalGoals: inputs.TotalGoals,
		FundingGap: inputs.FundingGap, RiskScore: inputs.RiskScore,
		HealthScore: inputs.HealthScore,
	}
	cards := c.buildCards(inputs)
	return &PlanningDashboard{
		Mode: inputs.Mode, Cards: cards, CurrentPlan: plan,
		Scenarios: inputs.Scenarios,
	}
}

func (c *Composer) BuildProjections(inputs Inputs) *ProjectionData {
	return &ProjectionData{
		NetWorthProjected: inputs.NetWorthProjected,
		IncomeProjected: inputs.IncomeProjected,
		ExpensesProjected: inputs.ExpensesProjected,
		PortfolioValue: inputs.PortfolioValue,
	}
}

func (c *Composer) BuildComparison(a, b PlanSummary) *PlanComparison {
	deltas := []ComparisonDelta{
		c.delta("Net Worth", a.NetWorthEnd, b.NetWorthEnd, true),
		c.deltaStr("Goals On Track", fmt.Sprint(a.GoalsOnTrack), fmt.Sprint(b.GoalsOnTrack), b.GoalsOnTrack > a.GoalsOnTrack),
		c.delta("Funding Gap", a.FundingGap, b.FundingGap, false),
		c.delta("Risk Score", int64(a.RiskScore), int64(b.RiskScore), false),
		c.delta("Health Score", int64(a.HealthScore), int64(b.HealthScore), true),
	}
	return &PlanComparison{PlanA: a, PlanB: b, Deltas: deltas}
}

func (c *Composer) delta(metric string, base, alt int64, higherBetter bool) ComparisonDelta {
	dir := "worse"
	if alt == base { dir = "unchanged" } else if higherBetter && alt > base { dir = "better" } else if !higherBetter && alt < base { dir = "better" }
	return ComparisonDelta{Metric: metric, Baseline: fmtMoney(base), Alternative: fmtMoney(alt), Delta: fmtMoney(alt - base), Direction: dir}
}

func (c *Composer) deltaStr(metric, base, alt string, better bool) ComparisonDelta {
	dir := "worse"
	if base == alt { dir = "unchanged" } else if better { dir = "better" }
	return ComparisonDelta{Metric: metric, Baseline: base, Alternative: alt, Delta: alt, Direction: dir}
}

func (c *Composer) buildCards(inputs Inputs) []Card {
	netWorthStr := fmtMoney(inputs.NetWorthProjected)
	cards := []Card{
		{CardID: "plan", CardType: CTPlanOverview, Title: "Financial Plan Overview",
			Summary: fmt.Sprintf("Net worth: %s", netWorthStr),
			Data:    map[string]interface{}{"net_worth": inputs.NetWorthProjected, "goals": fmt.Sprintf("%d/%d", inputs.GoalsOnTrack, inputs.TotalGoals)},
			Priority: 1},
		{CardID: "goal", CardType: CTGoalPlanning, Title: "Goal Planning",
			Summary: fmt.Sprintf("%d/%d goals on track", inputs.GoalsOnTrack, inputs.TotalGoals),
			Data:    map[string]interface{}{"on_track": inputs.GoalsOnTrack, "total": inputs.TotalGoals},
			Priority: 2},
		{CardID: "cf", CardType: CTCashFlowPlan, Title: "Cash Flow Planning",
			Summary: fmt.Sprintf("Income: %s | Expenses: %s", fmtMoney(inputs.MonthlyIncome), fmtMoney(inputs.MonthlyExpenses)),
			Priority: 3},
		{CardID: "ret", CardType: CTRetirementPlan, Title: "Retirement Planning",
			Summary: netWorthStr, Priority: 4},
		{CardID: "inv", CardType: CTInvestmentPlan, Title: "Investment Planning",
			Summary: fmt.Sprintf("Portfolio: %s", fmtMoney(inputs.PortfolioValue)), Priority: 5},
		{CardID: "debt", CardType: CTDebtPlan, Title: "Debt Planning",
			Summary: fmt.Sprintf("Funding gap: %s", fmtMoney(inputs.FundingGap)), Priority: 6},
		{CardID: "proj", CardType: CTProjection, Title: "Projection",
			Summary: netWorthStr, Data: map[string]interface{}{"projected": inputs.NetWorthProjected, "risk": inputs.RiskScore},
			Priority: 7},
	}
	cards = append(cards, Card{CardID: "risk", CardType: CTRiskSum, Title: "Risk Assessment",
		Summary: fmt.Sprintf("Score: %d — %s", inputs.RiskScore, inputs.RiskLevel), Priority: 8})
	if inputs.HasRecs {
		cards = append(cards, Card{CardID: "rec", CardType: CTRecommend, Title: "Recommendations",
			Summary: fmt.Sprintf("%d available", inputs.RecCount), Priority: 9})
	}
	if inputs.HasOptimization {
		cards = append(cards, Card{CardID: "opt", CardType: CTOptimization, Title: "Optimization",
			Summary: fmt.Sprintf("%d strategies", inputs.OptCount), Priority: 10})
	}
	if inputs.HasSimulation {
		cards = append(cards, Card{CardID: "sim", CardType: CTSimulation, Title: "Simulation",
			Summary: fmt.Sprintf("%d scenarios", inputs.SimCount), Priority: 11})
	}
	if inputs.EventCount > 0 {
		cards = append(cards, Card{CardID: "tl", CardType: CTTimeline, Title: "Timeline",
			Summary: fmt.Sprintf("%d events", inputs.EventCount), Priority: 12})
	}
	return cards
}

func CardSummary(count int) string { return fmt.Sprintf("%d items", count) }
func CardRiskSummary(score int, level string) string { return fmt.Sprintf("Score: %d — %s", score, level) }
func CardHealthSummary(score int, grade string) string { return fmt.Sprintf("Health: %d (%s)", score, grade) }
func CardEventSummary(count int) string { return fmt.Sprintf("%d events", count) }

func fmtMoney(v int64) string {
	if v >= 10000000 { return fmt.Sprintf("₹%.2fCr", float64(v)/10000000) }
	if v >= 100000 { return fmt.Sprintf("₹%.2fL", float64(v)/100000) }
	return fmt.Sprintf("₹%d", v)
}
