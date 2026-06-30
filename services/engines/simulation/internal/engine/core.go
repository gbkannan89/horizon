package engine

import (
	"fmt"
	"time"
)

// Inputs represents the baseline state for simulation.
type Inputs struct {
	NetWorth      int64   `json:"net_worth"`
	MonthlyIncome int64   `json:"monthly_income"`
	MonthlyEMI    int64   `json:"monthly_emi"`
	SavingsRate   float64 `json:"savings_rate"`
	EquityReturn  float64 `json:"equity_return"`
	InflationRate float64 `json:"inflation_rate"`
	GoalGap       int64   `json:"goal_gap"`
}

// Engine is the simulation computation engine.
type Engine struct{}

// NewEngine creates a new simulation engine.
func NewEngine() *Engine { return &Engine{} }

// Execute runs the simulation pipeline.
func (e *Engine) Execute(baseline Inputs, simType SimType, params map[string]interface{}) (*SimulationOutput, error) {
	scenario := e.applyScenario(baseline, simType, params)
	deltas := e.computeDeltas(baseline, scenario)
	summary := e.generateSummary(deltas)

	return &SimulationOutput{
		ScenarioID:   fmt.Sprintf("sim-%d", time.Now().UnixMilli()),
		SimType:      simType,
		Status:       SSCompleted,
		DeltaSummary: summary,
		Parameters:   params,
		Confidence:   "High",
	}, nil
}

func (e *Engine) applyScenario(baseline Inputs, st SimType, params map[string]interface{}) Inputs {
	s := baseline
	switch st {
	case STIncreaseSIP:
		if amt, ok := params["amount"].(float64); ok {
			s.MonthlyIncome += int64(amt)
		}
	case STDecreaseSIP:
		if amt, ok := params["amount"].(float64); ok {
			s.MonthlyIncome -= int64(amt)
		}
	case STEarlyLoanClosure, STExtraEMI:
		s.MonthlyEMI = 0
	case STSalaryIncrease:
		if amt, ok := params["amount"].(float64); ok {
			s.MonthlyIncome += int64(amt)
		}
	case STSalaryLoss:
		if pct, ok := params["percent"].(float64); ok {
			s.MonthlyIncome = int64(float64(s.MonthlyIncome) * (1.0 - pct/100.0))
		}
	case STMarketCrash:
		if pct, ok := params["percent"].(float64); ok {
			s.EquityReturn -= pct
		}
	case STMarketBoom:
		if pct, ok := params["percent"].(float64); ok {
			s.EquityReturn += pct
		}
	case STInflationChange:
		if pct, ok := params["percent"].(float64); ok {
			s.InflationRate = pct
		}
	case STExpenseReduction:
		if amt, ok := params["amount"].(float64); ok {
			s.MonthlyIncome += int64(amt)
		}
	case STInvestmentLumpSum:
		if amt, ok := params["amount"].(float64); ok {
			s.NetWorth += int64(amt)
		}
	}
	return s
}

func (e *Engine) computeDeltas(baseline, scenario Inputs) []DeltaEntry {
	var deltas []DeltaEntry
	deltas = append(deltas, e.delta("Net Worth", fmtMoney(baseline.NetWorth), fmtMoney(scenario.NetWorth)))
	deltas = append(deltas, e.delta("Monthly Income", fmtMoney(baseline.MonthlyIncome), fmtMoney(scenario.MonthlyIncome)))
	deltas = append(deltas, e.delta("Savings Rate", fmtFlt(baseline.SavingsRate), fmtFlt(scenario.SavingsRate)))
	deltas = append(deltas, e.delta("Monthly EMI", fmtMoney(baseline.MonthlyEMI), fmtMoney(scenario.MonthlyEMI)))
	return deltas
}

func (e *Engine) delta(metric, base, scenario string) DeltaEntry {
	d := DeltaEntry{Metric: metric, BaselineValue: base, ScenarioValue: scenario}
	if base == scenario { d.Direction = "NoChange"; d.Delta = "0" } else { d.Direction = "Change"; d.Delta = scenario }
	return d
}

func (e *Engine) generateSummary(deltas []DeltaEntry) []DeltaEntry { return deltas }

func fmtMoney(v int64) string { return fmt.Sprintf("₹%d", v) }
func fmtFlt(v float64) string { return fmt.Sprintf("%.1f", v) }
