package engine

import (
	"fmt"
	"math"
)

func (e *Engine) identifyGaps(inputs Inputs) []Recommendation {
	var recs []Recommendation
	for goal, gap := range inputs.GoalGaps {
		if gap < e.config.FundingGapThreshold { continue }
		recs = append(recs, Recommendation{
			Title: "Fund " + goal, Summary: "Gap of " + fmtMoney(gap),
			Action: "Increase contributions to " + goal, Category: RCGoalFunding,
			ExpectedImprovement: "Close funding gap for " + goal,
			Evidence: []string{"Gap: " + fmtMoney(gap), "Importance: " + inputs.GoalImportance[goal]},
		})
	}
	return recs
}

func (e *Engine) identifyOpportunities(inputs Inputs) []Recommendation {
	var recs []Recommendation
	for _, d := range inputs.Debts {
		if d.InterestRate < e.config.DebtInterestThreshold { continue }
		recs = append(recs, Recommendation{
			Title: "Optimize " + d.Name, Category: RCDebt,
			Summary: "High-interest debt at " + fmtFlt(d.InterestRate) + "%",
			Action: "Consider consolidation of " + d.Name,
			ExpectedImprovement: "Reduce interest payments",
			Evidence: []string{"Rate: " + fmtFlt(d.InterestRate) + "%", "Balance: " + fmtMoney(d.Balance)},
		})
	}
	if inputs.PortfolioDrift > e.config.DriftThresholdForRebalance {
		recs = append(recs, Recommendation{
			Title: "Rebalance portfolio", Category: RCPortfolioRebalancing,
			Summary: "Drift of " + fmtFlt(inputs.PortfolioDrift) + "%", Action: "Rebalance to target",
			ExpectedImprovement: "Reduce risk",
			Evidence: []string{"Drift: " + fmtFlt(inputs.PortfolioDrift) + "%"},
		})
	}
	if inputs.EmergencyMonths < e.config.EmergencyFundTargetMonths {
		recs = append(recs, Recommendation{
			Title: "Build emergency fund", Category: RCEmergencyFund,
			Summary: fmtFlt(inputs.EmergencyMonths) + " months coverage",
			Action: "Increase emergency contributions",
			ExpectedImprovement: "Improve resilience",
			Evidence: []string{"Coverage: " + fmtFlt(inputs.EmergencyMonths) + " months"},
		})
	}
	for cat, overspend := range inputs.OverspentCategories {
		if overspend > 0 {
			recs = append(recs, Recommendation{
				Title: "Reduce spending in " + cat, Category: RCCashFlow,
				Summary: "Overspent by " + fmtMoney(overspend),
				Action: "Adjust budget or reduce spending for " + cat,
				ExpectedImprovement: "Improve cash flow and budget adherence",
				Evidence: []string{"Overspend: " + fmtMoney(overspend), "Category: " + cat},
			})
		}
	}
	if inputs.CashFlowSurplus > 0 {
		recs = append(recs, Recommendation{
			Title: "Allocate surplus", Category: RCCashFlow,
			Summary: "Surplus of " + fmtMoney(inputs.CashFlowSurplus),
			Action: "Direct surplus to goals", ExpectedImprovement: "Accelerate progress",
			Evidence: []string{"Surplus: " + fmtMoney(inputs.CashFlowSurplus)},
		})
	}
	if inputs.SavingsRate < 20 {
		recs = append(recs, Recommendation{
			Title: "Increase savings rate", Category: RCSavings,
			Summary: "Rate: " + fmtFlt(inputs.SavingsRate) + "%",
			Action: "Increase savings by 5%", ExpectedImprovement: "Build wealth faster",
			Evidence: []string{"Rate: " + fmtFlt(inputs.SavingsRate) + "%"},
		})
	}
	return recs
}

func (e *Engine) scoreRecommendation(rec Recommendation, inputs Inputs) ScoringDimensions {
	s := ScoringDimensions{}
	switch rec.Category {
	case RCGoalFunding:
		s = ScoringDimensions{0.8, 0.7, 0.4, 0.8, 1.0, 0.7}
	case RCDebt:
		s = ScoringDimensions{0.7, 0.8, 0.5, 0.8, 0.6, 0.8}
	case RCPortfolioRebalancing:
		s = ScoringDimensions{0.6, 0.5, 0.3, 0.7, 0.5, 0.6}
	case RCEmergencyFund:
		s = ScoringDimensions{0.9, 0.9, 0.3, 0.9, 0.8, 1.0}
	case RCCashFlow:
		s = ScoringDimensions{0.6, 0.4, 0.2, 0.8, 0.7, 0.6}
	case RCSavings:
		s = ScoringDimensions{0.7, 0.5, 0.4, 0.7, 0.7, 0.7}
	default:
		s = ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}
	}
	return s
}

func (e *Engine) calculateComposite(s ScoringDimensions) float64 {
	score := s.Impact*0.3 + s.Urgency*0.2 + s.Difficulty*0.15 + s.Confidence*0.15 + s.GoalAlignment*0.1 + s.FinancialHealthImpact*0.1
	return math.Round(score*100) / 100
}

func (e *Engine) rankAndPrioritise(recs []Recommendation) []Recommendation {
	for i := 0; i < len(recs); i++ {
		for j := i + 1; j < len(recs); j++ {
			if recs[j].Score > recs[i].Score { recs[i], recs[j] = recs[j], recs[i] }
		}
	}
	for i := range recs { recs[i].Priority = i + 1 }
	return filterByThreshold(recs, e.config.MinRecommendationScore)
}

func filterByThreshold(recs []Recommendation, min float64) []Recommendation {
	var r []Recommendation
	for _, v := range recs { if v.Score >= min { r = append(r, v) } }
	return r
}

func (e *Engine) generateEvidence(rec Recommendation, inputs Inputs) []string { return rec.Evidence }
func (e *Engine) generateExplanation(rec Recommendation, inputs Inputs) string { return rec.Summary + ". " + rec.ExpectedImprovement + "." }

func fmtMoney(v int64) string { return fmt.Sprintf("₹%d", v) }
func fmtFlt(v float64) string { return fmt.Sprintf("%.1f", v) }
