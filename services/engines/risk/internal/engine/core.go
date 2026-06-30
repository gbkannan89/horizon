package engine

import "math"

// Inputs for risk assessment.
type Inputs struct {
	SingleAssetConcPct      float64 `json:"single_asset_conc_pct"`
	SectorConcPct           float64 `json:"sector_conc_pct"`
	InstitutionConcPct      float64 `json:"institution_conc_pct"`
	EquityExposureRatio     float64 `json:"equity_exposure_ratio"`
	CurrencyExposurePct     float64 `json:"currency_exposure_pct"`
	FloatingRateLiabPct     float64 `json:"floating_rate_liab_pct"`
	CashFlowCoverageMonths  float64 `json:"cash_flow_coverage_months"`
	LiquidityRatio          float64 `json:"liquidity_ratio"`
	DTIRatio                float64 `json:"dti_ratio"`
	CreditUtilizationPct    float64 `json:"credit_utilization_pct"`
	DepositConcentration    int64   `json:"deposit_concentration"`
	InflationExposedPct     float64 `json:"inflation_exposed_pct"`
	RetirementFundingRatio  float64 `json:"retirement_funding_ratio"`
	WithdrawalRate          float64 `json:"withdrawal_rate"`
	GoalFundingGapPct       float64 `json:"goal_funding_gap_pct"`
	IncomeVolatility        float64 `json:"income_volatility"`
	InsuranceGapLife        float64 `json:"insurance_gap_life"`
	InsuranceGapHealth      float64 `json:"insurance_gap_health"`
	SavingsConsistencyPct   float64 `json:"savings_consistency_pct"`
	AllocationDriftPct      float64 `json:"allocation_drift_pct"`
	SequenceRiskScore       float64 `json:"sequence_risk_score"`
	PortfolioValue          int64   `json:"portfolio_value"`
}

// Engine is the cross-cutting risk assessment engine.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

func (e *Engine) Execute(inputs Inputs) *RiskOutput {
	indicators := e.computeIndicators(inputs)
	catScores := e.computeCategoryScores(indicators)
	composite := e.computeComposite(catScores)
	stressTests := e.runStressTests(inputs)

	return &RiskOutput{
		CompositeScore: composite,
		CompositeLevel: e.toLevel(composite),
		Indicators:     indicators,
		CategoryScores: catScores,
		StressTests:    stressTests,
		Confidence:     "High",
	}
}

func (e *Engine) computeIndicators(inputs Inputs) []RiskIndicator {
	return []RiskIndicator{
		{"SingleAssetConcentration", e.scoreConc(inputs.SingleAssetConcPct, 15), RLCritical, 60, 80},
		{"SectorConcentration", e.scoreConc(inputs.SectorConcPct, 30), RLCritical, 60, 80},
		{"InstitutionConcentration", e.scoreConc(inputs.InstitutionConcPct, 40), RLCritical, 60, 80},
		{"EquityExposureRatio", e.scorePct(inputs.EquityExposureRatio, 70), RLCritical, 60, 80},
		{"CurrencyExposure", e.scorePct(inputs.CurrencyExposurePct, 20), RLCritical, 60, 80},
		{"InterestRateExposure", e.scorePct(inputs.FloatingRateLiabPct, 50), RLCritical, 60, 80},
		{"CashFlowCoverage", e.scoreInverse(inputs.CashFlowCoverageMonths, 3), RLCritical, 60, 80},
		{"LiquidityRatio", e.scoreInverse(inputs.LiquidityRatio, 1.5), RLCritical, 60, 80},
		{"DTIRatio", e.scorePct(inputs.DTIRatio, 40), RLCritical, 60, 80},
		{"CreditUtilization", e.scorePct(inputs.CreditUtilizationPct, 50), RLCritical, 60, 80},
		{"RetirementFundingRatio", e.scoreInverse(inputs.RetirementFundingRatio, 80), RLCritical, 60, 80},
		{"WithdrawalRate", e.scorePct(inputs.WithdrawalRate, 4), RLCritical, 60, 80},
		{"GoalFundingGapSeverity", e.scorePct(inputs.GoalFundingGapPct, 20), RLCritical, 60, 80},
		{"IncomeStability", int(math.Min(inputs.IncomeVolatility*50, 100)), RLCritical, 60, 80},
		{"InsuranceGapLife", e.scoreInverse(inputs.InsuranceGapLife, 10), RLCritical, 60, 80},
		{"InsuranceGapHealth", e.scorePct(100-inputs.InsuranceGapHealth, 50), RLCritical, 60, 80},
		{"SavingsConsistency", e.scoreInverse(inputs.SavingsConsistencyPct, 70), RLCritical, 60, 80},
		{"AllocationDrift", e.scorePct(inputs.AllocationDriftPct, 10), RLCritical, 60, 80},
	}
}

func (e *Engine) scorePct(actual, threshold float64) int {
	switch { case actual >= threshold*2: return 85; case actual >= threshold*1.5: return 70; case actual >= threshold: return 55; case actual >= threshold*0.5: return 35; default: return 15 }
}
func (e *Engine) scoreConc(actual, threshold float64) int {
	switch { case actual >= threshold*2: return 90; case actual >= threshold*1.5: return 75; case actual >= threshold: return 60; case actual >= threshold*0.5: return 35; default: return 15 }
}
func (e *Engine) scoreInverse(actual, target float64) int {
	if target <= 0 { return 50 }
	ratio := actual / target
	switch { case ratio >= 1.5: return 15; case ratio >= 1.0: return 35; case ratio >= 0.75: return 55; case ratio >= 0.5: return 70; default: return 85 }
}
func (e *Engine) scorePctSimple(pct float64) int {
	switch { case pct >= 40: return 80; case pct >= 30: return 65; case pct >= 15: return 45; case pct >= 5: return 25; default: return 10 }
}

func (e *Engine) computeCategoryScores(indicators []RiskIndicator) []RiskCategoryScore {
	cats := map[string][]int{}
	catMap := map[string]string{
		"SingleAssetConcentration": "Concentration", "SectorConcentration": "Concentration", "InstitutionConcentration": "Concentration",
		"EquityExposureRatio": "Market", "CurrencyExposure": "Market", "InterestRateExposure": "Market",
		"CashFlowCoverage": "Liquidity", "LiquidityRatio": "Liquidity",
		"DTIRatio": "Credit", "CreditUtilization": "Credit",
		"RetirementFundingRatio": "Longevity", "WithdrawalRate": "Longevity",
		"GoalFundingGapSeverity": "LifeEvent", "IncomeStability": "LifeEvent",
		"InsuranceGapLife": "Insurance", "InsuranceGapHealth": "Insurance",
		"SavingsConsistency": "Behavioural", "AllocationDrift": "Behavioural",
	}
	for _, ind := range indicators {
		cat := catMap[ind.Name]
		if _, ok := cats[cat]; !ok { cats[cat] = []int{} }
		cats[cat] = append(cats[cat], ind.Score)
	}
	var results []RiskCategoryScore
	for cat, scores := range cats {
		avg := 0; for _, s := range scores { avg += s }; avg /= len(scores)
		results = append(results, RiskCategoryScore{cat, avg, e.toLevel(avg)})
	}
	return results
}

func (e *Engine) computeComposite(cats []RiskCategoryScore) int {
	if len(cats) == 0 { return 0 }
	avg := 0; for _, c := range cats { avg += c.Score }; avg /= len(cats)
	return avg
}

func (e *Engine) toLevel(score int) RiskLevel {
	switch { case score >= 80: return RLCritical; case score >= 60: return RLElevated; case score >= 40: return RLModerate; case score >= 20: return RLLow; default: return RLMinimal }
}

func (e *Engine) runStressTests(inputs Inputs) []StressTestResult {
	return []StressTestResult{
		{"MarketCrash_30", e.crashImpact(inputs.PortfolioValue, 0.30), 730},
		{"ExtendedBearMarket", int64(float64(inputs.PortfolioValue) * -0.45), 1095},
		{"InterestRateShock_3pct", int64(float64(inputs.PortfolioValue) * -0.10), 180},
		{"InflationSpike_8pct", int64(float64(inputs.PortfolioValue) * -0.15), 365},
		{"IncomeLoss_6mo", int64(-500000), 180},
	}
}

func (e *Engine) crashImpact(portfolio int64, pct float64) int64 {
	return int64(-float64(portfolio) * pct)
}
