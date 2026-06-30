package engine

import "math"

// Inputs for health score computation.
type Inputs struct {
	SavingsRate          float64 `json:"savings_rate"`
	CashFlowSurplus      int64   `json:"cash_flow_surplus"`
	DebtToIncomeRatio    float64 `json:"debt_to_income_ratio"`
	LiquidityRatio       float64 `json:"liquidity_ratio"`
	EmergencyFundMonths  float64 `json:"emergency_fund_months"`
	GoalProgressAvg      float64 `json:"goal_progress_avg"`
	OnTrackRatio         float64 `json:"on_track_ratio"`
	InvestmentReturn     float64 `json:"investment_return"`
	DiversificationScore float64 `json:"diversification_score"`
	InsuranceCoverage    float64 `json:"insurance_coverage"`
	RetirementFundedRatio float64 `json:"retirement_funded_ratio"`
	AllocationAdherence  float64 `json:"allocation_adherence"`
	SavingsConsistency   float64 `json:"savings_consistency"`
	PortfolioDrift       float64 `json:"portfolio_drift"`
	PreviousOverallScore int     `json:"previous_overall_score"`
}

// Engine computes financial health scores across 13 dimensions.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

func (e *Engine) Execute(inputs Inputs) *HealthScoreOutput {
	dims := e.computeDimensions(inputs)
	overall := e.computeOverall(dims)
	grade := e.computeGrade(overall)
	trend := e.computeTrend(overall, inputs.PreviousOverallScore)

	return &HealthScoreOutput{
		OverallScore:  overall,
		ScoreGrade:    grade,
		Trend:         trend,
		Dimensions:    dims,
		PreviousScore: inputs.PreviousOverallScore,
		Confidence:    "High",
	}
}

func (e *Engine) computeDimensions(inputs Inputs) []DimensionScore {
	return []DimensionScore{
		{"SavingsHealth", e.scoreSavings(inputs.SavingsRate), 15, 0},
		{"CashFlowHealth", e.scoreCashFlow(inputs.CashFlowSurplus), 10, 0},
		{"DebtHealth", e.scoreDebt(inputs.DebtToIncomeRatio), 15, 0},
		{"LiquidityHealth", e.scoreLiquidity(inputs.LiquidityRatio), 5, 0},
		{"EmergencyFundHealth", e.scoreEmergency(inputs.EmergencyFundMonths), 10, 0},
		{"GoalProgressHealth", e.scoreGoalProgress(inputs.GoalProgressAvg), 15, 0},
		{"InvestmentHealth", e.scoreInvestment(inputs.InvestmentReturn), 5, 0},
		{"DiversificationHealth", e.scoreDiversification(inputs.DiversificationScore), 5, 0},
		{"InsuranceHealth", e.scoreInsurance(inputs.InsuranceCoverage), 5, 0},
		{"RetirementHealth", e.scoreRetirement(inputs.RetirementFundedRatio), 10, 0},
		{"BehaviourHealth", e.scoreBehaviour(inputs.AllocationAdherence, inputs.SavingsConsistency), 5, 0},
		{"HouseholdHealth", 80, 5, 0},
		{"PortfolioHealth", e.scorePortfolio(inputs.PortfolioDrift), 5, 0},
	}
}

func (e *Engine) computeOverall(dims []DimensionScore) int {
	total := 0.0
	for i := range dims {
		dims[i].Contribution = dims[i].Score * dims[i].Weight / 100.0
		total += dims[i].Contribution
	}
	return int(math.Round(total))
}

func (e *Engine) computeGrade(score int) ScoreGrade {
	switch { case score >= 90: return SGExcellent; case score >= 75: return SGGood; case score >= 60: return SGFair; case score >= 40: return SGNeedsWork; default: return SGCritical }
}

func (e *Engine) computeTrend(current, previous int) Trend {
	diff := current - previous
	switch { case diff >= 5: return TImproving; case diff <= -5: return TDeclining; default: return TStable }
}

func (e *Engine) scoreSavings(rate float64) float64 {
	switch { case rate >= 20: return 100; case rate >= 15: return 80; case rate >= 10: return 60; case rate >= 5: return 40; default: return 20 }
}
func (e *Engine) scoreCashFlow(surplus int64) float64 {
	switch { case surplus > 50000: return 100; case surplus > 20000: return 80; case surplus > 0: return 60; case surplus > -20000: return 40; default: return 20 }
}
func (e *Engine) scoreDebt(dti float64) float64 {
	switch { case dti <= 20: return 100; case dti <= 30: return 80; case dti <= 40: return 60; case dti <= 50: return 40; default: return 20 }
}
func (e *Engine) scoreLiquidity(ratio float64) float64 {
	switch { case ratio >= 2: return 100; case ratio >= 1.5: return 80; case ratio >= 1: return 60; case ratio >= 0.5: return 40; default: return 20 }
}
func (e *Engine) scoreEmergency(months float64) float64 {
	switch { case months >= 12: return 100; case months >= 6: return 80; case months >= 3: return 60; case months >= 1: return 40; default: return 20 }
}
func (e *Engine) scoreGoalProgress(avg float64) float64 {
	switch { case avg >= 80: return 100; case avg >= 60: return 80; case avg >= 40: return 60; case avg >= 20: return 40; default: return 20 }
}
func (e *Engine) scoreInvestment(ret float64) float64 {
	switch { case ret >= 12: return 100; case ret >= 8: return 80; case ret >= 4: return 60; case ret >= 0: return 40; default: return 20 }
}
func (e *Engine) scoreDiversification(ds float64) float64 {
	switch { case ds >= 0.7: return 100; case ds >= 0.5: return 80; case ds >= 0.3: return 60; case ds >= 0.1: return 40; default: return 20 }
}
func (e *Engine) scoreInsurance(cov float64) float64 {
	switch { case cov >= 90: return 100; case cov >= 70: return 80; case cov >= 50: return 60; case cov >= 25: return 40; default: return 20 }
}
func (e *Engine) scoreRetirement(fr float64) float64 {
	switch { case fr >= 90: return 100; case fr >= 70: return 80; case fr >= 50: return 60; case fr >= 25: return 40; default: return 20 }
}
func (e *Engine) scoreBehaviour(adherence, consistency float64) float64 { return (adherence + consistency) / 2 }
func (e *Engine) scorePortfolio(drift float64) float64 {
	switch { case drift <= 2: return 100; case drift <= 5: return 80; case drift <= 10: return 60; case drift <= 15: return 40; default: return 20 }
}
