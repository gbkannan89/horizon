package engine

import (
	"fmt"
	"math"
	"sort"
)

type Inputs struct {
	MonthlyIncome      int64              `json:"monthly_income"`
	MonthlyExpenses    int64              `json:"monthly_expenses"`
	CurrentDebt        int64              `json:"current_debt"`
	PortfolioValue     int64              `json:"portfolio_value"`
	Goals              map[string]float64 `json:"goals"`
	GoalContributions  map[string]float64 `json:"goal_contributions"`
	InterestRates      map[string]float64 `json:"interest_rates"`
	RiskTolerance      float64            `json:"risk_tolerance"`
	RetirementTarget   int64              `json:"retirement_target"`
	RetirementCurrent  int64              `json:"retirement_current"`
	RetirementAge      int                `json:"retirement_age"`
	CurrentAge         int                `json:"current_age"`
	EquityReturn       float64            `json:"equity_return"`
	DebtReturn         float64            `json:"debt_return"`
	EmergencyMonths    float64            `json:"emergency_months"`
}

type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

func (e *Engine) Execute(inputs Inputs, objectives []Objective, constraints []Constraint) *OptimizationOutput {
	candidates := e.generateCandidates(inputs)
	evaluated := len(candidates)
	valid := e.filterConstraints(candidates, constraints)
	for i := range valid {
		valid[i].ObjectiveScores = e.scoreObjectives(valid[i], inputs, objectives)
		valid[i].CompositeScore = e.compositeScore(valid[i].ObjectiveScores, objectives)
	}
	if len(valid) == 0 {
		return &OptimizationOutput{Objectives: objectives, Constraints: constraints, CandidatesEvaluated: evaluated, CandidatesValid: 0, Status: OSFailed}
	}
	e.rankSolutions(valid)
	return &OptimizationOutput{
		Objectives: objectives, Constraints: constraints,
		SelectedSolution: valid[0], Alternatives: e.topAlternatives(valid),
		CandidatesEvaluated: evaluated, CandidatesValid: len(valid), Status: OSCompleted,
	}
}

func (e *Engine) generateCandidates(inputs Inputs) []Solution {
	var sols []Solution
	surplus := inputs.MonthlyIncome - inputs.MonthlyExpenses
	for alloc := 0.2; alloc <= 0.8; alloc += 0.1 {
		contrib := int64(float64(surplus) * alloc)
		sols = append(sols, Solution{
			Parameters: map[string]float64{"contribution": float64(contrib), "debt_allocation": (1.0 - alloc) * 100, "equity_pct": e.defaultEquity(inputs.RiskTolerance)},
			Description: fmt.Sprintf("%.0f%% to goals", alloc*100),
		})
	}
	for eq := 0.3; eq <= 0.9; eq += 0.1 {
		sols = append(sols, Solution{
			Parameters: map[string]float64{"contribution": float64(int64(float64(surplus)*0.5)), "debt_allocation": 50, "equity_pct": eq * 100},
			Description: fmt.Sprintf("Equity %.0f%%", eq*100),
		})
	}
	if len(sols) == 0 {
		sols = append(sols, Solution{Parameters: map[string]float64{"contribution": 0, "debt_allocation": 50, "equity_pct": 60}, Description: "Default"})
	}
	return sols
}

func (e *Engine) defaultEquity(rt float64) float64 {
	switch { case rt >= 0.8: return 90; case rt >= 0.6: return 75; case rt >= 0.4: return 60; case rt >= 0.2: return 40; default: return 20 }
}

func (e *Engine) filterConstraints(sols []Solution, cons []Constraint) []Solution {
	var v []Solution
	for _, s := range sols {
		pass := true
		for _, c := range cons {
			if c.Type == "hard" {
				switch c.Name {
				case "equity_max": if s.Parameters["equity_pct"] > c.Limit { pass = false }
				case "budget": if s.Parameters["contribution"] > c.Limit { pass = false }
				}
			}
		}
		if pass { v = append(v, s) }
	}
	return v
}

func (e *Engine) scoreObjectives(s Solution, inputs Inputs, objs []Objective) map[string]float64 {
	scores := make(map[string]float64)
	contrib := s.Parameters["contribution"]; eq := s.Parameters["equity_pct"]
	for _, o := range objs {
		switch o.Name {
		case "MaximizeNetWorth":
			scores[o.Name] = math.Min((float64(inputs.PortfolioValue)+contrib*12*20)/1e7, 1.0)
		case "MinimizeInterestPaid":
			scores[o.Name] = 1.0 - math.Min(float64(inputs.CurrentDebt)*0.12*(1.0-s.Parameters["debt_allocation"]/100.0*0.3)/1e6, 1.0)
		case "MaximizeRetirementCorpus":
			proj := float64(inputs.RetirementCurrent) + contrib*12*float64(inputs.RetirementAge-inputs.CurrentAge)
			t := float64(inputs.RetirementTarget)
			if t <= 0 { scores[o.Name] = 0.5 } else { scores[o.Name] = math.Min(proj/t, 1.0) }
		case "MinimizeRisk":
			scores[o.Name] = 1.0 - eq/100.0
		default:
			scores[o.Name] = 0.5
		}
	}
	return scores
}

func (e *Engine) compositeScore(scores map[string]float64, objs []Objective) float64 {
	t := 0.0; for _, o := range objs { if s, ok := scores[o.Name]; ok { t += s * o.Weight } }; return math.Round(t*100) / 100
}

func (e *Engine) rankSolutions(sols []Solution) {
	sort.Slice(sols, func(i, j int) bool { return sols[i].CompositeScore > sols[j].CompositeScore })
	for i := range sols { sols[i].Rank = i + 1 }
}

func (e *Engine) topAlternatives(sols []Solution) []Solution {
	if len(sols) <= 1 { return nil }
	m := 3; if len(sols)-1 < m { m = len(sols) - 1 }; return sols[1 : m+1]
}
