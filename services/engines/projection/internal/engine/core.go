package engine

import (
	"fmt"
	"math"
	"time"
)

// Engine performs deterministic financial projections.
type Engine struct {
	assumptions Assumptions
}

// NewEngine creates a projection engine with given assumptions.
func NewEngine(assumptions Assumptions) *Engine {
	return &Engine{assumptions: assumptions}
}

// Execute runs the projection pipeline for a given type.
func (e *Engine) Execute(inputs Inputs, projType ProjectionType) *ProjectionOutput {
	now := time.Now().UTC()
	horizon := 30 * 365 // 30 years
	interval := 365     // annual snapshots

	var snaps []Snapshot
	for day := 0; day <= horizon; day += interval {
		d := now.AddDate(0, 0, day)
		s := Snapshot{Date: d.Format("2006-01-02")}

		switch projType {
		case PTCashFlow:
			s.NetCashFlow = e.projectCashFlow(inputs, day)
		case PTNetWorth:
			s.NetWorth = e.projectNetWorth(inputs, day)
			s.TotalAssets = inputs.TotalAssets + int64(float64(inputs.TotalAssets)*e.assumptions.EquityReturn/100.0*float64(day)/365.0)
			s.TotalLiabilities = e.projectDebt(inputs, day)
		case PTGoalCompletion:
			for goal, current := range inputs.GoalValues {
				target := inputs.GoalTargets[goal]
				growth := current + int64(float64(current)*e.assumptions.EquityReturn/100.0*float64(day)/365.0)
				if growth > target { growth = target }
				s.GoalValue += growth
			}
		case PTDebtPayoff:
			s.DebtValue = e.projectDebt(inputs, day)
		case PTInvestmentGrowth:
			for _, v := range inputs.AssetValues {
				s.InvValue += int64(float64(v) * (1.0 + e.assumptions.EquityReturn/100.0*float64(day)/365.0))
			}
		case PTRetirement:
			s.NetWorth = e.projectNetWorth(inputs, day)
			s.CashFlow = e.projectCashFlow(inputs, day)
		case PTFinancialIndependence:
			s.InvValue = e.projectInvestments(inputs, day)
			s.CashFlow = e.projectCashFlow(inputs, day)
		case PTPortfolioValue:
			s.PfValue = e.projectPortfolio(inputs, day)
		}

		snaps = append(snaps, s)
	}

	ms := e.detectMilestones(snaps, projType)
	summary := e.computeSummary(snaps, projType)

	return &ProjectionOutput{
		ProjectionType:  projType,
		Version:         1,
		Status:          ESCompleted,
		Timeline:        snaps,
		Milestones:      ms,
		SummaryMetrics:  summary,
		Assumptions:     e.assumptions,
		ConfidenceLevel: e.computeConfidence(inputs),
	}
}

func (e *Engine) projectNetWorth(inputs Inputs, day int) int64 {
	years := float64(day) / 365.0
	growth := 1.0 + e.assumptions.EquityReturn/100.0*years
	return int64(float64(inputs.NetWorth) * growth)
}

func (e *Engine) projectCashFlow(inputs Inputs, day int) int64 {
	years := float64(day) / 365.0
	income := (float64(inputs.MonthlyIncome) + float64(inputs.RecurringIncome)) * (1.0 + e.assumptions.SalaryGrowthRate/100.0*years)
	baseExp := inputs.MonthlyExpenses + inputs.RecurringExpenses
	if inputs.BudgetExpenses > 0 { baseExp = inputs.BudgetExpenses + inputs.RecurringExpenses }
	expenses := float64(baseExp) * (1.0 + e.assumptions.ExpenseGrowthRate/100.0*years)
	return int64(income - expenses)
}

func (e *Engine) projectDebt(inputs Inputs, day int) int64 {
	var total int64
	for _, v := range inputs.Liabilities {
		remaining := float64(v) - (float64(v)/120.0)*float64(day/30)
		if remaining < 0 { remaining = 0 }
		interest := remaining * e.assumptions.LoanInterestRate / 100.0 * float64(day) / 365.0
		total += int64(remaining + interest)
	}
	return total
}

func (e *Engine) projectInvestments(inputs Inputs, day int) int64 {
	var total int64
	for _, v := range inputs.AssetValues {
		total += int64(float64(v) * (1.0 + e.assumptions.EquityReturn/100.0*float64(day)/365.0))
	}
	return total
}

func (e *Engine) projectPortfolio(inputs Inputs, day int) int64 {
	return e.projectInvestments(inputs, day)
}

func (e *Engine) detectMilestones(snaps []Snapshot, pt ProjectionType) []Milestone {
	var ms []Milestone
	for i, s := range snaps {
		if i == 0 { continue }
		if pt == PTDebtPayoff && s.DebtValue <= 0 && snaps[i-1].DebtValue > 0 {
			ms = append(ms, Milestone{Date: s.Date, Type: "DebtFree", Description: "Projected debt-free date"})
		}
		if pt == PTNetWorth {
			for _, m := range []int64{1000000, 5000000, 10000000, 50000000, 100000000} {
				if s.NetWorth >= m && snaps[i-1].NetWorth < m {
					ms = append(ms, Milestone{Date: s.Date, Type: "NetWorthMilestone", Description: fmt.Sprintf("Net worth crosses ₹%d", m)})
				}
			}
		}
	}
	return ms
}

func (e *Engine) computeSummary(snaps []Snapshot, pt ProjectionType) SummaryMetrics {
	if len(snaps) == 0 { return SummaryMetrics{} }
	first := snaps[0]
	last := snaps[len(snaps)-1]

	sm := SummaryMetrics{}
	switch pt {
	case PTNetWorth:
		sm.FinalNetWorth = last.NetWorth
		sm.TotalGrowth = last.NetWorth - first.NetWorth
		sm.AnnualizedReturn = e.calcReturn(first.NetWorth, last.NetWorth, len(snaps))
	case PTCashFlow:
		sm.AvgMonthlyCashFlow = last.CashFlow
		sm.TotalCashFlow = last.CashFlow * 12
	case PTDebtPayoff:
		sm.FinalDebtBalance = last.DebtValue
		sm.DebtFreeDate = ""
		for _, s := range snaps {
			if s.DebtValue <= 0 { sm.DebtFreeDate = s.Date; break }
		}
	case PTGoalCompletion:
		sm.GoalCompletionRate = math.Min(float64(last.GoalValue)/100.0, 1.0)
	case PTInvestmentGrowth:
		sm.TotalGrowth = last.InvValue
	case PTRetirement:
		sm.FinalNetWorth = last.NetWorth
		sm.RetirementReadiness = e.assessRetirement(last.NetWorth)
	}
	return sm
}

func (e *Engine) calcReturn(start, end int64, years int) float64 {
	if start == 0 || years == 0 { return 0 }
	return (float64(end)/float64(start) - 1.0) * 100.0
}

func (e *Engine) assessRetirement(netWorth int64) string {
	target := int64(50000000) // ₹5Cr target
	if netWorth >= target { return "OnTrack" }
	if netWorth >= target/2 { return "Progressing" }
	return "NeedsAttention"
}

func (e *Engine) computeConfidence(inputs Inputs) string {
	if len(inputs.GoalValues) > 0 && len(inputs.AssetValues) > 0 {
		return "High"
	}
	if len(inputs.AssetValues) > 0 {
		return "Medium"
	}
	return "Low"
}
