package engine

import "fmt"

// Engine composes the planning workspace from engine outputs.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// BuildWorkspace composes the planning workspace view.
func (e *Engine) BuildWorkspace(inputs Inputs) *PlanningWorkspace {
	cards := e.buildCards(inputs)

	current := &PlanSummary{
		PlanID: "baseline", Name: "Current Plan", IsBaseline: true,
		NetWorthAtEnd: inputs.NetWorthProjected, GoalsOnTrack: inputs.GoalsOnTrack,
		TotalGoals: inputs.TotalGoals, FundingGap: inputs.FundingGapTotal,
		RiskScore: inputs.RiskScore,
	}

	var alts []PlanSummary
	if inputs.HasAlternatives {
		alts = append(alts, PlanSummary{
			PlanID: "alt-1", Name: "Alternative Plan", IsBaseline: false,
			NetWorthAtEnd: inputs.NetWorthProjected + int64(float64(inputs.NetWorthProjected)*0.1),
			GoalsOnTrack: inputs.GoalsOnTrack, TotalGoals: inputs.TotalGoals,
			FundingGap: inputs.FundingGapTotal - 100000, RiskScore: inputs.RiskScore - 5,
		})
	}

	return &PlanningWorkspace{
		Mode: inputs.Mode, Cards: cards,
		CurrentPlan: current, AlternativePlans: alts,
	}
}

// ComparePlans produces a side-by-side comparison.
func (e *Engine) ComparePlans(a, b PlanSummary) *PlanComparison {
	deltas := []ComparisonDelta{
		{Metric: "Net Worth", Baseline: fmtMoney(a.NetWorthAtEnd), Alternative: fmtMoney(b.NetWorthAtEnd),
			Delta: fmtMoney(b.NetWorthAtEnd - a.NetWorthAtEnd), Direction: e.direction(b.NetWorthAtEnd > a.NetWorthAtEnd)},
		{Metric: "Goals On Track", Baseline: fmt.Sprint(a.GoalsOnTrack), Alternative: fmt.Sprint(b.GoalsOnTrack),
			Delta: fmt.Sprint(b.GoalsOnTrack - a.GoalsOnTrack), Direction: e.direction(b.GoalsOnTrack > a.GoalsOnTrack)},
		{Metric: "Funding Gap", Baseline: fmtMoney(a.FundingGap), Alternative: fmtMoney(b.FundingGap),
			Delta: fmtMoney(b.FundingGap - a.FundingGap), Direction: e.direction(b.FundingGap < a.FundingGap)},
		{Metric: "Risk Score", Baseline: fmt.Sprint(a.RiskScore), Alternative: fmt.Sprint(b.RiskScore),
			Delta: fmt.Sprint(b.RiskScore - a.RiskScore), Direction: e.direction(b.RiskScore < a.RiskScore)},
	}
	return &PlanComparison{PlanA: a, PlanB: b, Deltas: deltas}
}

func (e *Engine) buildCards(inputs Inputs) []PlanningCard {
	cards := []PlanningCard{
		{CardID: "cp", CardType: CTCurrentPlan, Title: "Current Plan",
			Summary: fmt.Sprintf("Projected net worth: %s", fmtMoney(inputs.NetWorthProjected)),
			Detail: fmt.Sprintf("%d/%d goals on track", inputs.GoalsOnTrack, inputs.TotalGoals), Priority: 1},
	}
	if inputs.HasProjection {
		cards = append(cards, PlanningCard{CardID: "fp", CardType: CTFutureProjection, Title: "Future Projection",
			Summary: fmt.Sprintf("Net worth at horizon: %s", fmtMoney(inputs.NetWorthProjected)), Priority: 2})
	}
	cards = append(cards, PlanningCard{CardID: "fg", CardType: CTFundingGap, Title: "Funding Gap",
		Summary: fmt.Sprintf("Total gap: %s across %d goals", fmtMoney(inputs.FundingGapTotal), inputs.TotalGoals), Priority: 3})

	return cards
}

func (e *Engine) direction(better bool) string {
	if better { return "better" }; return "worse"
}

func fmtMoney(v int64) string { return fmt.Sprintf("₹%d", v) }
