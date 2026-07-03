package engine

import "testing"

func TestExecute(t *testing.T) {
	e := NewEngine()
	inputs := Inputs{
		SavingsRate: 20, CashFlowSurplus: 60000, DebtToIncomeRatio: 15,
		LiquidityRatio: 2.5, EmergencyFundMonths: 12, GoalProgressAvg: 90,
		OnTrackRatio: 0.8, InvestmentReturn: 15, DiversificationScore: 0.8,
		InsuranceCoverage: 95, RetirementFundedRatio: 95, AllocationAdherence: 90,
		SavingsConsistency: 90, PortfolioDrift: 1, BudgetAdherence: 95,
		PreviousOverallScore: 80,
	}

	output := e.Execute(inputs)
	if output.OverallScore < 80 {
		t.Errorf("expected high score for perfect inputs, got %d", output.OverallScore)
	}
	if output.ScoreGrade != SGExcellent {
		t.Errorf("expected Excellent grade, got %s", output.ScoreGrade)
	}
	if len(output.Dimensions) != 13 {
		t.Errorf("expected 13 dimensions, got %d", len(output.Dimensions))
	}
	if output.Trend != TImproving {
		t.Errorf("expected Improving trend (95->80), got %s", output.Trend)
	}
}

func TestExecute_PoorInputs(t *testing.T) {
	e := NewEngine()
	inputs := Inputs{
		SavingsRate: 0, CashFlowSurplus: -50000, DebtToIncomeRatio: 60,
		LiquidityRatio: 0.2, EmergencyFundMonths: 0, GoalProgressAvg: 5,
		InvestmentReturn: -5, DiversificationScore: 0, InsuranceCoverage: 0,
		RetirementFundedRatio: 0, AllocationAdherence: 0, SavingsConsistency: 0,
		PortfolioDrift: 30, BudgetAdherence: 0, PreviousOverallScore: 50,
	}

	output := e.Execute(inputs)
	if output.OverallScore > 40 {
		t.Errorf("expected low score for poor inputs, got %d", output.OverallScore)
	}
	if output.Trend != TDeclining {
		t.Errorf("expected Declining trend (20->50), got %s", output.Trend)
	}
}

func TestScoreSavings(t *testing.T) {
	e := NewEngine()
	tests := []struct{ rate, expected float64 }{
		{25, 100}, {20, 100}, {18, 80}, {12, 60}, {7, 40}, {2, 20}, {0, 20},
	}
	for _, tc := range tests {
		if got := e.scoreSavings(tc.rate); got != tc.expected {
			t.Errorf("scoreSavings(%.0f) = %.0f, want %.0f", tc.rate, got, tc.expected)
		}
	}
}

func TestScoreCashFlow(t *testing.T) {
	e := NewEngine()
	tests := []struct{ surplus int64; expected float64 }{
		{60000, 100}, {50001, 100}, {30000, 80}, {10000, 60}, {-10000, 40}, {-30000, 20},
	}
	for _, tc := range tests {
		if got := e.scoreCashFlow(tc.surplus); got != tc.expected {
			t.Errorf("scoreCashFlow(%d) = %.0f, want %.0f", tc.surplus, got, tc.expected)
		}
	}
}

func TestScoreDebt(t *testing.T) {
	e := NewEngine()
	tests := []struct{ dti, expected float64 }{
		{15, 100}, {25, 80}, {35, 60}, {45, 40}, {55, 20},
	}
	for _, tc := range tests {
		if got := e.scoreDebt(tc.dti); got != tc.expected {
			t.Errorf("scoreDebt(%.0f) = %.0f, want %.0f", tc.dti, got, tc.expected)
		}
	}
}

func TestScoreEmergency(t *testing.T) {
	e := NewEngine()
	tests := []struct{ months, expected float64 }{
		{12, 100}, {8, 80}, {4, 60}, {2, 40}, {0, 20},
	}
	for _, tc := range tests {
		if got := e.scoreEmergency(tc.months); got != tc.expected {
			t.Errorf("scoreEmergency(%.0f) = %.0f, want %.0f", tc.months, got, tc.expected)
		}
	}
}

func TestScoreGoalProgress(t *testing.T) {
	e := NewEngine()
	tests := []struct{ avg, expected float64 }{
		{90, 100}, {70, 80}, {50, 60}, {30, 40}, {10, 20},
	}
	for _, tc := range tests {
		if got := e.scoreGoalProgress(tc.avg); got != tc.expected {
			t.Errorf("scoreGoalProgress(%.0f) = %.0f, want %.0f", tc.avg, got, tc.expected)
		}
	}
}

func TestScoreBudget(t *testing.T) {
	e := NewEngine()
	tests := []struct{ adherence, expected float64 }{
		{95, 100}, {80, 80}, {60, 60}, {35, 40}, {10, 20},
	}
	for _, tc := range tests {
		if got := e.scoreBudget(tc.adherence); got != tc.expected {
			t.Errorf("scoreBudget(%.0f) = %.0f, want %.0f", tc.adherence, got, tc.expected)
		}
	}
}

func TestScorePortfolio(t *testing.T) {
	e := NewEngine()
	tests := []struct{ drift, expected float64 }{
		{1, 100}, {3, 80}, {7, 60}, {12, 40}, {20, 20},
	}
	for _, tc := range tests {
		if got := e.scorePortfolio(tc.drift); got != tc.expected {
			t.Errorf("scorePortfolio(%.0f) = %.0f, want %.0f", tc.drift, got, tc.expected)
		}
	}
}

func TestScoreBehaviour(t *testing.T) {
	e := NewEngine()
	if got := e.scoreBehaviour(80, 90); got != 85 {
		t.Errorf("scoreBehaviour(80,90) = %.0f, want 85", got)
	}
}

func TestComputeGrade(t *testing.T) {
	e := NewEngine()
	tests := []struct{ score int; grade ScoreGrade }{
		{100, SGExcellent}, {90, SGExcellent}, {85, SGGood}, {75, SGGood},
		{70, SGFair}, {60, SGFair}, {50, SGNeedsWork}, {40, SGNeedsWork},
		{30, SGCritical}, {0, SGCritical},
	}
	for _, tc := range tests {
		if got := e.computeGrade(tc.score); got != tc.grade {
			t.Errorf("computeGrade(%d) = %s, want %s", tc.score, got, tc.grade)
		}
	}
}

func TestComputeTrend(t *testing.T) {
	e := NewEngine()
	tests := []struct{ current, previous int; expected Trend }{
		{85, 80, TImproving}, {80, 80, TStable}, {75, 80, TDeclining},
		{85, 80, TImproving}, {90, 80, TImproving}, {70, 80, TDeclining},
	}
	for _, tc := range tests {
		if got := e.computeTrend(tc.current, tc.previous); got != tc.expected {
			t.Errorf("computeTrend(%d,%d) = %s, want %s", tc.current, tc.previous, got, tc.expected)
		}
	}
}

func TestScoreLiquidity(t *testing.T) {
	e := NewEngine()
	tests := []struct{ ratio, expected float64 }{
		{2.5, 100}, {1.8, 80}, {1.2, 60}, {0.7, 40}, {0.3, 20},
	}
	for _, tc := range tests {
		if got := e.scoreLiquidity(tc.ratio); got != tc.expected {
			t.Errorf("scoreLiquidity(%.1f) = %.0f, want %.0f", tc.ratio, got, tc.expected)
		}
	}
}

func TestScoreInvestment(t *testing.T) {
	e := NewEngine()
	tests := []struct{ ret, expected float64 }{
		{15, 100}, {10, 80}, {6, 60}, {2, 40}, {-1, 20},
	}
	for _, tc := range tests {
		if got := e.scoreInvestment(tc.ret); got != tc.expected {
			t.Errorf("scoreInvestment(%.0f) = %.0f, want %.0f", tc.ret, got, tc.expected)
		}
	}
}

func TestScoreDiversification(t *testing.T) {
	e := NewEngine()
	tests := []struct{ ds, expected float64 }{
		{0.8, 100}, {0.6, 80}, {0.4, 60}, {0.2, 40}, {0.05, 20},
	}
	for _, tc := range tests {
		if got := e.scoreDiversification(tc.ds); got != tc.expected {
			t.Errorf("scoreDiversification(%.1f) = %.0f, want %.0f", tc.ds, got, tc.expected)
		}
	}
}

func TestScoreInsurance(t *testing.T) {
	e := NewEngine()
	tests := []struct{ cov, expected float64 }{
		{95, 100}, {75, 80}, {60, 60}, {30, 40}, {10, 20},
	}
	for _, tc := range tests {
		if got := e.scoreInsurance(tc.cov); got != tc.expected {
			t.Errorf("scoreInsurance(%.0f) = %.0f, want %.0f", tc.cov, got, tc.expected)
		}
	}
}

func TestScoreRetirement(t *testing.T) {
	e := NewEngine()
	tests := []struct{ fr, expected float64 }{
		{95, 100}, {75, 80}, {55, 60}, {30, 40}, {10, 20},
	}
	for _, tc := range tests {
		if got := e.scoreRetirement(tc.fr); got != tc.expected {
			t.Errorf("scoreRetirement(%.0f) = %.0f, want %.0f", tc.fr, got, tc.expected)
		}
	}
}
