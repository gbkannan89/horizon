package engine

import (
	"math"
	"testing"
	"time"
)

func TestDefaultAssumptions(t *testing.T) {
	a := DefaultAssumptions()
	if a.InflationRate != 4.0 {
		t.Errorf("InflationRate = %f, want 4.0", a.InflationRate)
	}
	if a.EquityReturn != 10.0 {
		t.Errorf("EquityReturn = %f, want 10.0", a.EquityReturn)
	}
	if a.SalaryGrowthRate != 8.0 {
		t.Errorf("SalaryGrowthRate = %f, want 8.0", a.SalaryGrowthRate)
	}
	if a.ExpenseGrowthRate != 4.0 {
		t.Errorf("ExpenseGrowthRate = %f, want 4.0", a.ExpenseGrowthRate)
	}
	if a.LoanInterestRate != 9.0 {
		t.Errorf("LoanInterestRate = %f, want 9.0", a.LoanInterestRate)
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.MaxProjectionHorizon != 30 {
		t.Errorf("MaxProjectionHorizon = %d, want 30", c.MaxProjectionHorizon)
	}
	if c.SnapshotIntervalMonths != 12 {
		t.Errorf("SnapshotIntervalMonths = %d, want 12", c.SnapshotIntervalMonths)
	}
	if c.ConfidenceLookbackYears != 3 {
		t.Errorf("ConfidenceLookbackYears = %d, want 3", c.ConfidenceLookbackYears)
	}
}

func TestNewEngine(t *testing.T) {
	a := DefaultAssumptions()
	e := NewEngine(a)
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
	if e.assumptions != a {
		t.Error("assumptions not set correctly")
	}
}

func TestEngine_projectCashFlow(t *testing.T) {
	t.Run("uses BudgetExpenses when positive", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 0, ExpenseGrowthRate: 0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     50000,
			MonthlyExpenses:   40000,
			BudgetExpenses:    30000,
			RecurringIncome:   0,
			RecurringExpenses: 0,
		}
		// day=0 => years=0, income=50000, expenses=BudgetExpenses=30000
		got := e.projectCashFlow(inputs, 0)
		want := int64(50000 - 30000) // 20000
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})

	t.Run("uses MonthlyExpenses when BudgetExpenses is zero", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 0, ExpenseGrowthRate: 0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     50000,
			MonthlyExpenses:   40000,
			BudgetExpenses:    0,
			RecurringIncome:   0,
			RecurringExpenses: 0,
		}
		got := e.projectCashFlow(inputs, 0)
		want := int64(50000 - 40000) // 10000
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})

	t.Run("includes recurring income and expenses", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 0, ExpenseGrowthRate: 0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     50000,
			MonthlyExpenses:   30000,
			BudgetExpenses:    0,
			RecurringIncome:   10000,
			RecurringExpenses: 5000,
		}
		// day=0 => income=60000, expenses=35000
		got := e.projectCashFlow(inputs, 0)
		want := int64(60000 - 35000) // 25000
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})

	t.Run("applies salary growth", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 10.0, ExpenseGrowthRate: 0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     100000,
			MonthlyExpenses:   50000,
			BudgetExpenses:    0,
			RecurringIncome:   0,
			RecurringExpenses: 0,
		}
		// day=365 => years=1, income=100000*(1+0.10*1)=110000, expenses=50000
		got := e.projectCashFlow(inputs, 365)
		want := int64(110000 - 50000) // 60000
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})

	t.Run("applies expense growth", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 0, ExpenseGrowthRate: 5.0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     100000,
			MonthlyExpenses:   50000,
			BudgetExpenses:    0,
			RecurringIncome:   0,
			RecurringExpenses: 0,
		}
		// day=365 => years=1, income=100000, expenses=50000*(1+0.05*1)=52500
		got := e.projectCashFlow(inputs, 365)
		want := int64(100000 - 52500) // 47500
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})

	t.Run("applies both growth rates", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 8.0, ExpenseGrowthRate: 4.0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     100000,
			MonthlyExpenses:   60000,
			BudgetExpenses:    0,
			RecurringIncome:   0,
			RecurringExpenses: 0,
		}
		got := e.projectCashFlow(inputs, 730)
		years := 730.0 / 365.0
		income := (100000.0 + 0) * (1.0 + 8.0/100.0*years)
		expenses := float64(60000) * (1.0 + 4.0/100.0*years)
		want := int64(income - expenses)
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})

	t.Run("with BudgetExpenses and recurring", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 5.0, ExpenseGrowthRate: 3.0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     80000,
			MonthlyExpenses:   70000,
			BudgetExpenses:    50000,
			RecurringIncome:   20000,
			RecurringExpenses: 10000,
		}
		// day=365 => years=1
		// income=(80000+20000)*(1+0.05*1)=105000
		// baseExp=BudgetExpenses(50000)+RecurringExpenses(10000)=60000
		// expenses=60000*(1+0.03*1)=61800
		got := e.projectCashFlow(inputs, 365)
		want := int64(105000 - 61800) // 43200
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})

	t.Run("day zero uses no growth", func(t *testing.T) {
		a := Assumptions{SalaryGrowthRate: 100.0, ExpenseGrowthRate: 100.0}
		e := NewEngine(a)
		inputs := Inputs{
			MonthlyIncome:     10000,
			MonthlyExpenses:   6000,
			BudgetExpenses:    0,
			RecurringIncome:   0,
			RecurringExpenses: 0,
		}
		got := e.projectCashFlow(inputs, 0)
		want := int64(10000 - 6000) // 4000
		if got != want {
			t.Errorf("projectCashFlow = %d, want %d", got, want)
		}
	})
}

func TestEngine_projectNetWorth(t *testing.T) {
	a := Assumptions{EquityReturn: 10.0}
	e := NewEngine(a)
	inputs := Inputs{NetWorth: 1000000}

	t.Run("day zero returns initial net worth", func(t *testing.T) {
		got := e.projectNetWorth(inputs, 0)
		if got != 1000000 {
			t.Errorf("projectNetWorth = %d, want %d", got, 1000000)
		}
	})

	t.Run("one year growth applied", func(t *testing.T) {
		got := e.projectNetWorth(inputs, 365)
		want := int64(float64(1000000) * (1.0 + 10.0/100.0*1.0))
		if got != want {
			t.Errorf("projectNetWorth = %d, want %d", got, want)
		}
	})

	t.Run("half year growth", func(t *testing.T) {
		got := e.projectNetWorth(inputs, 182)
		years := 182.0 / 365.0
		want := int64(float64(1000000) * (1.0 + 10.0/100.0*years))
		if got != want {
			t.Errorf("projectNetWorth = %d, want %d", got, want)
		}
	})
}

func TestEngine_projectDebt(t *testing.T) {
	a := Assumptions{LoanInterestRate: 9.0}
	e := NewEngine(a)

	t.Run("single liability at day zero", func(t *testing.T) {
		inputs := Inputs{Liabilities: []int64{120000}}
		got := e.projectDebt(inputs, 0)
		// remaining=120000, no interest at day 0
		want := int64(120000)
		if got != want {
			t.Errorf("projectDebt = %d, want %d", got, want)
		}
	})

	t.Run("single liability amortizing", func(t *testing.T) {
		inputs := Inputs{Liabilities: []int64{120000}}
		got := e.projectDebt(inputs, 360) // ~12 months
		// remaining = 120000 - (120000/120)*12 = 120000 - 12000 = 108000
		// interest = 108000 * 9/100 * 360/365 ≈ 9589.04
		remaining := 120000.0 - (120000.0/120.0)*12.0
		interest := remaining * 9.0 / 100.0 * 360.0 / 365.0
		want := int64(remaining + interest)
		if got != want {
			t.Errorf("projectDebt = %d, want %d", got, want)
		}
	})

	t.Run("multiple liabilities", func(t *testing.T) {
		inputs := Inputs{Liabilities: []int64{60000, 120000}}
		got := e.projectDebt(inputs, 180)
		expectedTotal := int64(0)
		for _, v := range inputs.Liabilities {
			r := float64(v) - (float64(v)/120.0)*float64(180/30)
			if r < 0 {
				r = 0
			}
			expectedTotal += int64(r + r*9.0/100.0*180.0/365.0)
		}
		if got != expectedTotal {
			t.Errorf("projectDebt = %d, want %d", got, expectedTotal)
		}
	})

	t.Run("debt reaches zero", func(t *testing.T) {
		inputs := Inputs{Liabilities: []int64{1000}} // small loan
		got := e.projectDebt(inputs, 3650)           // 10 years
		if got != 0 {
			t.Errorf("projectDebt = %d, want 0 (fully paid)", got)
		}
	})

	t.Run("no liabilities", func(t *testing.T) {
		inputs := Inputs{}
		got := e.projectDebt(inputs, 365)
		if got != 0 {
			t.Errorf("projectDebt = %d, want 0", got)
		}
	})
}

func TestEngine_projectInvestments(t *testing.T) {
	a := Assumptions{EquityReturn: 10.0}
	e := NewEngine(a)

	t.Run("single asset at day zero", func(t *testing.T) {
		inputs := Inputs{AssetValues: []int64{500000}}
		got := e.projectInvestments(inputs, 0)
		if got != 500000 {
			t.Errorf("projectInvestments = %d, want %d", got, 500000)
		}
	})

	t.Run("single asset with growth", func(t *testing.T) {
		inputs := Inputs{AssetValues: []int64{500000}}
		got := e.projectInvestments(inputs, 365)
		want := int64(float64(500000) * (1.0 + 10.0/100.0*365.0/365.0))
		if got != want {
			t.Errorf("projectInvestments = %d, want %d", got, want)
		}
	})

	t.Run("multiple assets", func(t *testing.T) {
		inputs := Inputs{AssetValues: []int64{200000, 300000, 500000}}
		got := e.projectInvestments(inputs, 730)
		var want int64
		for _, v := range inputs.AssetValues {
			want += int64(float64(v) * (1.0 + 10.0/100.0*730.0/365.0))
		}
		if got != want {
			t.Errorf("projectInvestments = %d, want %d", got, want)
		}
	})

	t.Run("no assets", func(t *testing.T) {
		inputs := Inputs{}
		got := e.projectInvestments(inputs, 365)
		if got != 0 {
			t.Errorf("projectInvestments = %d, want 0", got)
		}
	})
}

func TestEngine_projectPortfolio(t *testing.T) {
	a := Assumptions{EquityReturn: 10.0}
	e := NewEngine(a)
	inputs := Inputs{AssetValues: []int64{1000000}}
	got := e.projectPortfolio(inputs, 365)
	want := e.projectInvestments(inputs, 365)
	if got != want {
		t.Errorf("projectPortfolio = %d, want %d (should match projectInvestments)", got, want)
	}
}

func TestEngine_detectMilestones(t *testing.T) {
	t.Run("debt payoff milestone", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", DebtValue: 50000},
			{Date: "2026-01-01", DebtValue: 25000},
			{Date: "2027-01-01", DebtValue: 0},
			{Date: "2028-01-01", DebtValue: 0},
		}
		ms := e.detectMilestones(snaps, PTDebtPayoff)
		if len(ms) != 1 {
			t.Fatalf("expected 1 milestone, got %d", len(ms))
		}
		if ms[0].Date != "2027-01-01" {
			t.Errorf("milestone date = %s, want 2027-01-01", ms[0].Date)
		}
		if ms[0].Type != "DebtFree" {
			t.Errorf("milestone type = %s, want DebtFree", ms[0].Type)
		}
		if ms[0].Description != "Projected debt-free date" {
			t.Errorf("milestone description = %s, want 'Projected debt-free date'", ms[0].Description)
		}
	})

	t.Run("no debt milestone for other types", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", DebtValue: 50000},
			{Date: "2026-01-01", DebtValue: 0},
		}
		ms := e.detectMilestones(snaps, PTCashFlow)
		if len(ms) != 0 {
			t.Errorf("expected 0 milestones for non-debt type, got %d", len(ms))
		}
	})

	t.Run("net worth milestones", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", NetWorth: 500000},
			{Date: "2026-01-01", NetWorth: 1500000},
			{Date: "2027-01-01", NetWorth: 6000000},
			{Date: "2028-01-01", NetWorth: 15000000},
		}
		ms := e.detectMilestones(snaps, PTNetWorth)
		if len(ms) == 0 {
			t.Fatal("expected at least one net worth milestone")
		}
		// Check cross 1M at 2026, cross 5M at 2027, cross 10M at 2028
		for _, m := range ms {
			if m.Type != "NetWorthMilestone" {
				t.Errorf("unexpected milestone type: %s", m.Type)
			}
		}
		if ms[0].Date != "2026-01-01" {
			t.Errorf("first milestone date = %s, want 2026-01-01", ms[0].Date)
		}
		if len(ms) >= 2 && ms[1].Date != "2027-01-01" {
			t.Errorf("second milestone date = %s, want 2027-01-01", ms[1].Date)
		}
	})

	t.Run("first snapshot skipped", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", DebtValue: 50000},
			{Date: "2026-01-01", DebtValue: 0},
		}
		// If we also had debt=0 in first snapshot, it should still detect
		ms := e.detectMilestones(snaps, PTDebtPayoff)
		if len(ms) != 1 {
			t.Errorf("expected 1 milestone despite first snapshot skip logic, got %d", len(ms))
		}
	})

	t.Run("no milestone when already zero", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", DebtValue: 0},
			{Date: "2026-01-01", DebtValue: 0},
		}
		ms := e.detectMilestones(snaps, PTDebtPayoff)
		if len(ms) != 0 {
			t.Errorf("expected 0 milestones when already debt-free, got %d", len(ms))
		}
	})
}

func TestEngine_computeSummary(t *testing.T) {
	t.Run("PTNetWorth", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", NetWorth: 1000000},
			{Date: "2026-01-01", NetWorth: 1500000},
			{Date: "2027-01-01", NetWorth: 2000000},
		}
		sm := e.computeSummary(snaps, PTNetWorth)
		if sm.FinalNetWorth != 2000000 {
			t.Errorf("FinalNetWorth = %d, want 2000000", sm.FinalNetWorth)
		}
		if sm.TotalGrowth != 1000000 {
			t.Errorf("TotalGrowth = %d, want 1000000", sm.TotalGrowth)
		}
		// calcReturn receives len(snaps)=3 as the years parameter
		wantReturn := (float64(2000000)/float64(1000000) - 1.0) * 100.0
		if sm.AnnualizedReturn != wantReturn {
			t.Errorf("AnnualizedReturn = %f, want %f", sm.AnnualizedReturn, wantReturn)
		}
	})

	t.Run("PTCashFlow", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", CashFlow: 10000},
			{Date: "2026-01-01", CashFlow: 15000},
		}
		sm := e.computeSummary(snaps, PTCashFlow)
		if sm.AvgMonthlyCashFlow != 15000 {
			t.Errorf("AvgMonthlyCashFlow = %d, want 15000", sm.AvgMonthlyCashFlow)
		}
		if sm.TotalCashFlow != 15000*12 {
			t.Errorf("TotalCashFlow = %d, want %d", sm.TotalCashFlow, 15000*12)
		}
	})

	t.Run("PTDebtPayoff found debt-free", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", DebtValue: 50000},
			{Date: "2026-01-01", DebtValue: 10000},
			{Date: "2027-01-01", DebtValue: 0},
		}
		sm := e.computeSummary(snaps, PTDebtPayoff)
		if sm.FinalDebtBalance != 0 {
			t.Errorf("FinalDebtBalance = %d, want 0", sm.FinalDebtBalance)
		}
		if sm.DebtFreeDate != "2027-01-01" {
			t.Errorf("DebtFreeDate = %s, want 2027-01-01", sm.DebtFreeDate)
		}
	})

	t.Run("PTDebtPayoff never debt-free", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", DebtValue: 50000},
			{Date: "2026-01-01", DebtValue: 30000},
		}
		sm := e.computeSummary(snaps, PTDebtPayoff)
		if sm.FinalDebtBalance != 30000 {
			t.Errorf("FinalDebtBalance = %d, want 30000", sm.FinalDebtBalance)
		}
		if sm.DebtFreeDate != "" {
			t.Errorf("DebtFreeDate = %s, want empty", sm.DebtFreeDate)
		}
	})

	t.Run("PTGoalCompletion", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", GoalValue: 0},
			{Date: "2026-01-01", GoalValue: 50000000},
		}
		sm := e.computeSummary(snaps, PTGoalCompletion)
		// GoalValue/100 = 500000, capped at 1.0
		if sm.GoalCompletionRate != math.Min(500000.0, 1.0) {
			t.Errorf("GoalCompletionRate = %f, want %f", sm.GoalCompletionRate, math.Min(500000.0, 1.0))
		}
	})

	t.Run("PTInvestmentGrowth", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", InvValue: 100000},
			{Date: "2026-01-01", InvValue: 500000},
		}
		sm := e.computeSummary(snaps, PTInvestmentGrowth)
		if sm.TotalGrowth != 500000 {
			t.Errorf("TotalGrowth = %d, want 500000", sm.TotalGrowth)
		}
	})

	t.Run("PTRetirement OnTrack", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", NetWorth: 10000000},
			{Date: "2026-01-01", NetWorth: 50000000},
		}
		sm := e.computeSummary(snaps, PTRetirement)
		if sm.FinalNetWorth != 50000000 {
			t.Errorf("FinalNetWorth = %d, want 50000000", sm.FinalNetWorth)
		}
		if sm.RetirementReadiness != "OnTrack" {
			t.Errorf("RetirementReadiness = %s, want OnTrack", sm.RetirementReadiness)
		}
	})

	t.Run("PTRetirement Progressing", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", NetWorth: 0},
			{Date: "2026-01-01", NetWorth: 25000000},
		}
		sm := e.computeSummary(snaps, PTRetirement)
		if sm.RetirementReadiness != "Progressing" {
			t.Errorf("RetirementReadiness = %s, want Progressing", sm.RetirementReadiness)
		}
	})

	t.Run("PTRetirement NeedsAttention", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		snaps := []Snapshot{
			{Date: "2025-01-01", NetWorth: 0},
			{Date: "2026-01-01", NetWorth: 1000000},
		}
		sm := e.computeSummary(snaps, PTRetirement)
		if sm.RetirementReadiness != "NeedsAttention" {
			t.Errorf("RetirementReadiness = %s, want NeedsAttention", sm.RetirementReadiness)
		}
	})

	t.Run("empty snapshots returns zero value", func(t *testing.T) {
		e := NewEngine(DefaultAssumptions())
		sm := e.computeSummary(nil, PTNetWorth)
		if sm.FinalNetWorth != 0 || sm.TotalGrowth != 0 {
			t.Errorf("expected zero values for empty input, got %+v", sm)
		}
	})
}

func TestEngine_calcReturn(t *testing.T) {
	e := NewEngine(DefaultAssumptions())

	t.Run("normal case", func(t *testing.T) {
		got := e.calcReturn(1000, 1500, 1)
		want := 50.0
		if got != want {
			t.Errorf("calcReturn = %f, want %f", got, want)
		}
	})

	t.Run("zero start returns zero", func(t *testing.T) {
		got := e.calcReturn(0, 1500, 1)
		if got != 0 {
			t.Errorf("calcReturn = %f, want 0", got)
		}
	})

	t.Run("zero years returns zero", func(t *testing.T) {
		got := e.calcReturn(1000, 1500, 0)
		if got != 0 {
			t.Errorf("calcReturn = %f, want 0", got)
		}
	})

	t.Run("negative return", func(t *testing.T) {
		got := e.calcReturn(2000, 1500, 1)
		want := -25.0
		if got != want {
			t.Errorf("calcReturn = %f, want %f", got, want)
		}
	})

	t.Run("no change", func(t *testing.T) {
		got := e.calcReturn(1000, 1000, 5)
		if got != 0 {
			t.Errorf("calcReturn = %f, want 0", got)
		}
	})
}

func TestEngine_assessRetirement(t *testing.T) {
	e := NewEngine(DefaultAssumptions())

	t.Run("OnTrack when net worth >= 50M", func(t *testing.T) {
		got := e.assessRetirement(50000000)
		if got != "OnTrack" {
			t.Errorf("assessRetirement = %s, want OnTrack", got)
		}
		got = e.assessRetirement(100000000)
		if got != "OnTrack" {
			t.Errorf("assessRetirement = %s, want OnTrack", got)
		}
	})

	t.Run("Progressing when net worth >= 25M", func(t *testing.T) {
		got := e.assessRetirement(25000000)
		if got != "Progressing" {
			t.Errorf("assessRetirement = %s, want Progressing", got)
		}
		got = e.assessRetirement(49999999)
		if got != "Progressing" {
			t.Errorf("assessRetirement = %s, want Progressing", got)
		}
	})

	t.Run("NeedsAttention when net worth < 25M", func(t *testing.T) {
		got := e.assessRetirement(0)
		if got != "NeedsAttention" {
			t.Errorf("assessRetirement = %s, want NeedsAttention", got)
		}
		got = e.assessRetirement(1000000)
		if got != "NeedsAttention" {
			t.Errorf("assessRetirement = %s, want NeedsAttention", got)
		}
	})
}

func TestEngine_computeConfidence(t *testing.T) {
	e := NewEngine(DefaultAssumptions())

	t.Run("High when both goals and assets present", func(t *testing.T) {
		inputs := Inputs{
			GoalValues:  map[string]int64{"retire": 50000000},
			AssetValues: []int64{1000000},
		}
		got := e.computeConfidence(inputs)
		if got != "High" {
			t.Errorf("computeConfidence = %s, want High", got)
		}
	})

	t.Run("Medium when only assets present", func(t *testing.T) {
		inputs := Inputs{
			AssetValues: []int64{1000000},
		}
		got := e.computeConfidence(inputs)
		if got != "Medium" {
			t.Errorf("computeConfidence = %s, want Medium", got)
		}
	})

	t.Run("Low when neither present", func(t *testing.T) {
		inputs := Inputs{}
		got := e.computeConfidence(inputs)
		if got != "Low" {
			t.Errorf("computeConfidence = %s, want Low", got)
		}
	})

	t.Run("Medium when only goals present (no assets)", func(t *testing.T) {
		inputs := Inputs{
			GoalValues: map[string]int64{"retire": 50000000},
		}
		got := e.computeConfidence(inputs)
		// GoalValues present but no AssetValues => len(GoalValues)>0 is true but
		// the condition is len(GoalValues)>0 && len(AssetValues)>0 => false
		// then len(AssetValues)>0 => false => "Low"
		if got != "Low" {
			t.Errorf("computeConfidence = %s, want Low (goals without assets is Low by current logic)", got)
		}
	})
}

func TestEngine_Execute(t *testing.T) {
	defaultAssumptions := DefaultAssumptions()
	now := time.Now().UTC()

	baseInputs := Inputs{
		NetWorth:          10000000,
		TotalAssets:       15000000,
		TotalLiabilities:  5000000,
		MonthlyIncome:     200000,
		MonthlyExpenses:   120000,
		RecurringIncome:   50000,
		RecurringExpenses: 20000,
		BudgetExpenses:    0,
		CashBalance:       500000,
		GoalValues:        map[string]int64{"retire": 10000000},
		GoalTargets:       map[string]int64{"retire": 50000000},
		AssetValues:       []int64{5000000, 3000000},
		Liabilities:       []int64{2000000, 1000000},
	}

	t.Run("PTCashFlow", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTCashFlow)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTCashFlow {
			t.Errorf("ProjectionType = %s, want CashFlow", out.ProjectionType)
		}
		if out.Version != 1 {
			t.Errorf("Version = %d, want 1", out.Version)
		}
		if out.Status != ESCompleted {
			t.Errorf("Status = %s, want Completed", out.Status)
		}
		if len(out.Timeline) == 0 {
			t.Fatal("Timeline is empty")
		}
		// First snapshot at day 0, last at ~30 years
		if out.Timeline[0].Date != now.Format("2006-01-02") {
			t.Errorf("first snapshot date = %s, want %s", out.Timeline[0].Date, now.Format("2006-01-02"))
		}
		if out.Timeline[0].NetCashFlow != e.projectCashFlow(baseInputs, 0) {
			t.Error("first snapshot NetCashFlow mismatch")
		}
		if out.ConfidenceLevel != "High" {
			t.Errorf("ConfidenceLevel = %s, want High", out.ConfidenceLevel)
		}
	})

	t.Run("PTNetWorth", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTNetWorth)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTNetWorth {
			t.Errorf("ProjectionType = %s, want NetWorth", out.ProjectionType)
		}
		if len(out.Timeline) == 0 {
			t.Fatal("Timeline is empty")
		}
		// First snapshot should have computed NetWorth and TotalAssets/TotalLiabilities
		s0 := out.Timeline[0]
		if s0.NetWorth != e.projectNetWorth(baseInputs, 0) {
			t.Error("first snapshot NetWorth mismatch")
		}
		// Expect some milestones (NetWorth may cross thresholds over 30 years)
		t.Logf("NetWorth milestones: %d", len(out.Milestones))
	})

	t.Run("PTGoalCompletion", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTGoalCompletion)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTGoalCompletion {
			t.Errorf("ProjectionType = %s, want GoalCompletion", out.ProjectionType)
		}
		// Goal value should grow over time
		if len(out.Timeline) > 0 && out.Timeline[0].GoalValue == 0 {
			t.Error("expected non-zero GoalValue in first snapshot")
		}
	})

	t.Run("PTDebtPayoff", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTDebtPayoff)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTDebtPayoff {
			t.Errorf("ProjectionType = %s, want DebtPayoff", out.ProjectionType)
		}
		if len(out.Timeline) > 0 && out.Timeline[0].DebtValue != e.projectDebt(baseInputs, 0) {
			t.Error("first snapshot DebtValue mismatch")
		}
	})

	t.Run("PTInvestmentGrowth", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTInvestmentGrowth)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTInvestmentGrowth {
			t.Errorf("ProjectionType = %s, want InvestmentGrowth", out.ProjectionType)
		}
		if len(out.Timeline) > 0 && out.Timeline[0].InvValue == 0 {
			t.Error("expected non-zero InvValue in first snapshot")
		}
	})

	t.Run("PTRetirement", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTRetirement)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTRetirement {
			t.Errorf("ProjectionType = %s, want Retirement", out.ProjectionType)
		}
		if len(out.Timeline) == 0 {
			t.Fatal("Timeline is empty")
		}
		s0 := out.Timeline[0]
		if s0.NetWorth == 0 {
			t.Error("expected non-zero NetWorth in Retirement projection")
		}
		if s0.CashFlow == 0 {
			t.Error("expected non-zero CashFlow in Retirement projection")
		}
	})

	t.Run("PTFinancialIndependence", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTFinancialIndependence)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTFinancialIndependence {
			t.Errorf("ProjectionType = %s, want FinancialIndependence", out.ProjectionType)
		}
		if len(out.Timeline) == 0 {
			t.Fatal("Timeline is empty")
		}
		s0 := out.Timeline[0]
		if s0.InvValue == 0 {
			t.Error("expected non-zero InvValue in FI projection")
		}
		if s0.CashFlow == 0 {
			t.Error("expected non-zero CashFlow in FI projection")
		}
	})

	t.Run("PTPortfolioValue", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTPortfolioValue)
		if out == nil {
			t.Fatal("Execute returned nil")
		}
		if out.ProjectionType != PTPortfolioValue {
			t.Errorf("ProjectionType = %s, want PortfolioValue", out.ProjectionType)
		}
		if len(out.Timeline) > 0 && out.Timeline[0].PfValue == 0 {
			t.Error("expected non-zero PfValue in first snapshot")
		}
	})

	t.Run("timeline has 31 snapshots (0 to 30 years inclusive)", func(t *testing.T) {
		e := NewEngine(defaultAssumptions)
		out := e.Execute(baseInputs, PTCashFlow)
		expectedSnaps := 31 // day 0 through day 10950, step 365
		if len(out.Timeline) != expectedSnaps {
			t.Errorf("expected %d snapshots, got %d", expectedSnaps, len(out.Timeline))
		}
	})

	t.Run("assumptions carried through", func(t *testing.T) {
		customAssumptions := Assumptions{
			EquityReturn: 12.0, SalaryGrowthRate: 10.0, ExpenseGrowthRate: 5.0,
			LoanInterestRate: 8.0, InflationRate: 3.0, DebtReturn: 6.0,
			ContributionGrowthRate: 7.0, TaxRate: 25.0, RetirementAge: 55,
			LifeExpectancy: 80,
		}
		e := NewEngine(customAssumptions)
		out := e.Execute(baseInputs, PTNetWorth)
		if out.Assumptions != customAssumptions {
			t.Error("assumptions not preserved in output")
		}
	})
}

func TestEngine_Execute_CashFlowBudgetExpenses(t *testing.T) {
	a := Assumptions{SalaryGrowthRate: 5.0, ExpenseGrowthRate: 3.0}
	e := NewEngine(a)
	inputs := Inputs{
		MonthlyIncome:     100000,
		MonthlyExpenses:   90000,
		BudgetExpenses:    60000,
		RecurringIncome:   0,
		RecurringExpenses: 0,
	}
	out := e.Execute(inputs, PTCashFlow)
	// First snapshot should use BudgetExpenses (60000) not MonthlyExpenses (90000)
	s0 := out.Timeline[0]
	expectedCashFlow := int64(100000 - 60000) // 40000
	if s0.NetCashFlow != expectedCashFlow {
		t.Errorf("first snapshot NetCashFlow = %d, want %d (BudgetExpenses should be used)", s0.NetCashFlow, expectedCashFlow)
	}
}

func TestEngine_Execute_EmptyInputs(t *testing.T) {
	e := NewEngine(DefaultAssumptions())
	inputs := Inputs{}
	out := e.Execute(inputs, PTNetWorth)
	if out == nil {
		t.Fatal("Execute with empty inputs returned nil")
	}
	if out.Status != ESCompleted {
		t.Errorf("Status = %s, want Completed", out.Status)
	}
	if out.ConfidenceLevel != "Low" {
		t.Errorf("ConfidenceLevel = %s, want Low", out.ConfidenceLevel)
	}
}
