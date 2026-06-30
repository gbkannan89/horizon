package engine

import (
	"fmt"
	"time"
)

type Engine struct {
	assumptions Assumptions
}

func NewEngine(assumptions Assumptions) *Engine {
	return &Engine{assumptions: assumptions}
}

func (e *Engine) Execute(inputs Inputs, projType ProjectionType) (*ProjectionOutput, error) {
	now := time.Now().UTC()
	frozen := inputs.freeze()
	versions := map[string]int{"input": 1}
	snapshots := e.buildTimeline(frozen, projType, now)
	milestones := e.detectMilestones(snapshots, projType)

	output := &ProjectionOutput{
		OutputID:       fmt.Sprintf("proj-%d", now.UnixMilli()),
		ProjectionType: projType,
		Version:        1,
		Status:         ESCompleted,
		Timeline:       snapshots,
		Milestones:     milestones,
		InputVersions:  versions,
		Assumptions:    e.assumptions,
		Confidence:     "High",
	}
	return output, nil
}

type Inputs struct {
	NetWorth    int64
	CashBalance int64
	GoalValues  map[string]int64
	GoalTargets map[string]int64
	AssetValues []int64
	Liabilities []int64
	Income      int64
	Expenses    int64
	Age         int
}

func (i Inputs) freeze() Inputs { return i }

func (e *Engine) buildTimeline(inputs Inputs, pt ProjectionType, start time.Time) []Snapshot {
	horizon := 30 * 365; snaps := []Snapshot{}; interval := 365
	for day := 0; day <= horizon; day += interval {
		d := start.AddDate(0, 0, day)
		s := Snapshot{Date: d.Format("2006-01-02")}
		switch pt {
		case PTNetWorth:
			s.NetWorth = e.projectNetWorth(inputs, day)
		case PTCashFlow:
			s.CashFlow = e.projectCashFlow(inputs, day)
		case PTGoalCompletion:
			s.GoalValue = e.projectGoalValue(inputs, day)
		case PTDebtPayoff:
			s.DebtValue = e.projectDebt(inputs, day)
		case PTInvestmentGrowth:
			s.InvValue = e.projectInvestments(inputs, day)
		case PTPortfolioValue:
			s.PfValue = e.projectPortfolio(inputs, day)
		case PTRetirement, PTFinancialIndependence:
			s.NetWorth = e.projectNetWorth(inputs, day)
		}
		snaps = append(snaps, s)
	}
	return snaps
}

func (e *Engine) projectNetWorth(inputs Inputs, day int) int64 {
	g := 1.0 + (e.assumptions.EquityReturn/100.0)*(float64(day)/365.0)
	return int64(float64(inputs.NetWorth) * g)
}
func (e *Engine) projectCashFlow(inputs Inputs, day int) int64 {
	y := float64(day) / 365.0
	ig := 1.0 + (e.assumptions.SalaryGrowthRate/100.0)*y
	eg := 1.0 + (e.assumptions.ExpenseGrowthRate/100.0)*y
	return int64(float64(inputs.Income)*ig) - int64(float64(inputs.Expenses)*eg)
}
func (e *Engine) projectGoalValue(inputs Inputs, day int) int64 {
	if len(inputs.GoalValues) == 0 { return 0 }; var total int64
	for _, v := range inputs.GoalValues {
		total += int64(float64(v) * (1.0 + (e.assumptions.EquityReturn/100.0)*(float64(day)/365.0)))
	}
	return total
}
func (e *Engine) projectDebt(inputs Inputs, day int) int64 {
	var total int64
	for _, v := range inputs.Liabilities {
		remaining := float64(v) - (float64(v)/120.0)*float64(day/30)
		if remaining < 0 { remaining = 0 }
		total += int64(remaining + remaining*(e.assumptions.LoanInterestRate/100.0)*(float64(day)/365.0))
	}
	return total
}
func (e *Engine) projectInvestments(inputs Inputs, day int) int64 {
	var total int64
	for _, v := range inputs.AssetValues {
		total += int64(float64(v) * (1.0 + (e.assumptions.EquityReturn/100.0)*(float64(day)/365.0)))
	}
	return total
}
func (e *Engine) projectPortfolio(inputs Inputs, day int) int64 { return e.projectInvestments(inputs, day) }

func (e *Engine) detectMilestones(snaps []Snapshot, pt ProjectionType) []Milestone {
	var ms []Milestone
	for i, s := range snaps {
		if i == 0 { continue }
		if pt == PTDebtPayoff && s.DebtValue <= 0 && snaps[i-1].DebtValue > 0 {
			ms = append(ms, Milestone{Date: s.Date, MilestoneType: "DebtFree", Description: "Projected debt-free date"})
		}
	}
	return ms
}
