package engine

import (
	"testing"
)

func defaultInputs() Inputs {
	return Inputs{
		NetWorth:      5_000_000,
		MonthlyIncome: 150_000,
		MonthlyEMI:    25_000,
		SavingsRate:   0.3,
		EquityReturn:  12.0,
		InflationRate: 6.0,
		GoalGap:       500_000,
	}
}

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
}

func TestExecute_BaselineOnly(t *testing.T) {
	e := NewEngine()
	in := defaultInputs()
	params := map[string]interface{}{}

	out, err := e.Execute(in, STSalaryIncrease, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != SSCompleted {
		t.Errorf("expected status Completed, got %s", out.Status)
	}
	if out.SimType != STSalaryIncrease {
		t.Errorf("expected sim type SalaryIncrease, got %s", out.SimType)
	}
	if out.ScenarioID == "" {
		t.Error("expected non-empty scenario ID")
	}
	if out.Confidence != "High" {
		t.Errorf("expected High confidence, got %s", out.Confidence)
	}
}

func TestExecute_StatusConstants(t *testing.T) {
	if SSDraft != "Draft" || SSRunning != "Running" || SSCompleted != "Completed" || SSFailed != "Failed" || SSExpired != "Expired" || SSArchived != "Archived" || SSHistorical != "Historical" {
		t.Error("one or more SimStatus constants changed unexpectedly")
	}
}

// ---------------------------------------------------------------------------
// All 20 simulation types
// ---------------------------------------------------------------------------

func TestSimType_STIncreaseSIP(t *testing.T) {
	e := NewEngine()
	in := defaultInputs()
	params := map[string]interface{}{"amount": float64(10_000)}

	out, err := e.Execute(in, STIncreaseSIP, params)
	if err != nil {
		t.Fatal(err)
	}
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹160000" {
		t.Errorf("expected scenario income ₹160000, got %s", income.ScenarioValue)
	}
}

func TestSimType_STDecreaseSIP(t *testing.T) {
	e := NewEngine()
	in := defaultInputs()
	params := map[string]interface{}{"amount": float64(5_000)}

	out, _ := e.Execute(in, STDecreaseSIP, params)
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹145000" {
		t.Errorf("expected ₹145000, got %s", income.ScenarioValue)
	}
}

func TestSimType_STEarlyLoanClosure(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{}
	out, _ := e.Execute(defaultInputs(), STEarlyLoanClosure, params)
	emi := out.DeltaSummary[3]
	if emi.ScenarioValue != "₹0" {
		t.Errorf("expected ₹0 EMI after closure, got %s", emi.ScenarioValue)
	}
}

func TestSimType_STExtraEMI(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STExtraEMI, map[string]interface{}{})
	emi := out.DeltaSummary[3]
	if emi.ScenarioValue != "₹0" {
		t.Errorf("expected ₹0 EMI (extra EMI pays off), got %s", emi.ScenarioValue)
	}
}

func TestSimType_STSalaryIncrease(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"amount": float64(25_000)}
	out, _ := e.Execute(defaultInputs(), STSalaryIncrease, params)
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹175000" {
		t.Errorf("expected ₹175000, got %s", income.ScenarioValue)
	}
}

func TestSimType_STSalaryLoss(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		want    string
	}{
		{"ten_percent", 10.0, "₹135000"},
		{"fifty_percent", 50.0, "₹75000"},
		{"zero_percent", 0.0, "₹150000"},
		{"hundred_percent", 100.0, "₹0"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEngine()
			params := map[string]interface{}{"percent": float64(tc.percent)}
			out, _ := e.Execute(defaultInputs(), STSalaryLoss, params)
			income := out.DeltaSummary[1]
			if income.ScenarioValue != tc.want {
				t.Errorf("expected %s, got %s", tc.want, income.ScenarioValue)
			}
		})
	}
}

func TestSimType_STJobChange(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STJobChange, map[string]interface{}{})
	// no-op: values should be identical to baseline
	verifyNoChange(t, out)
}

func TestSimType_STMarriage(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STMarriage, map[string]interface{}{})
	verifyNoChange(t, out)
}

func TestSimType_STChild(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STChild, map[string]interface{}{})
	verifyNoChange(t, out)
}

func TestSimType_STHousePurchase(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STHousePurchase, map[string]interface{}{})
	verifyNoChange(t, out)
}

func TestSimType_STVehiclePurchase(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STVehiclePurchase, map[string]interface{}{})
	verifyNoChange(t, out)
}

func TestSimType_STInvestmentLumpSum(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"amount": float64(1_000_000)}
	out, _ := e.Execute(defaultInputs(), STInvestmentLumpSum, params)
	nw := out.DeltaSummary[0]
	if nw.ScenarioValue != "₹6000000" {
		t.Errorf("expected ₹6000000, got %s", nw.ScenarioValue)
	}
}

func TestSimType_STExpenseReduction(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"amount": float64(8_000)}
	out, _ := e.Execute(defaultInputs(), STExpenseReduction, params)
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹158000" {
		t.Errorf("expected ₹158000, got %s", income.ScenarioValue)
	}
}

func TestSimType_STGoalDelay(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STGoalDelay, map[string]interface{}{})
	verifyNoChange(t, out)
}

func TestSimType_STGoalPriorityChange(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STGoalPriorityChange, map[string]interface{}{})
	verifyNoChange(t, out)
}

func TestSimType_STEmergency(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STEmergency, map[string]interface{}{})
	verifyNoChange(t, out)
}

func TestSimType_STMarketCrash(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"percent": float64(8.0)}
	out, _ := e.Execute(defaultInputs(), STMarketCrash, params)
	// EquityReturn: 12.0 - 8.0 = 4.0 — but the delta entry is for "Savings Rate" field,
	// and the scenario EquityReturn is only used internally, not surfaced in deltas directly.
	// We verify the returned deltas still contain the correct structure.
	if len(out.DeltaSummary) != 4 {
		t.Fatalf("expected 4 delta entries, got %d", len(out.DeltaSummary))
	}
	// Net worth and monthly income should be unchanged
	nw := out.DeltaSummary[0]
	if nw.Direction != "NoChange" {
		t.Errorf("net worth should be NoChange on crash (no amount param), got %s", nw.Direction)
	}
}

func TestSimType_STMarketBoom(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"percent": float64(5.0)}
	out, _ := e.Execute(defaultInputs(), STMarketBoom, params)
	if len(out.DeltaSummary) != 4 {
		t.Fatalf("expected 4 delta entries, got %d", len(out.DeltaSummary))
	}
	nw := out.DeltaSummary[0]
	if nw.Direction != "NoChange" {
		t.Errorf("net worth should be NoChange on boom (no amount param), got %s", nw.Direction)
	}
}

func TestSimType_STInflationChange(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"percent": float64(8.5)}
	out, _ := e.Execute(defaultInputs(), STInflationChange, params)
	if len(out.DeltaSummary) != 4 {
		t.Fatalf("expected 4 delta entries, got %d", len(out.DeltaSummary))
	}
}

func TestSimType_STRetirementAgeChange(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STRetirementAgeChange, map[string]interface{}{})
	verifyNoChange(t, out)
}

// ---------------------------------------------------------------------------
// Edge cases: missing / wrong-typed / nil params
// ---------------------------------------------------------------------------

func TestExecute_MissingAmount(t *testing.T) {
	e := NewEngine()
	// amount missing — type assertion fails, no mutation
	out, _ := e.Execute(defaultInputs(), STIncreaseSIP, map[string]interface{}{})
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹150000" {
		t.Errorf("expected no change ₹150000, got %s", income.ScenarioValue)
	}
}

func TestExecute_MissingPercent(t *testing.T) {
	e := NewEngine()
	// percent missing for SalaryLoss — no mutation
	out, _ := e.Execute(defaultInputs(), STSalaryLoss, map[string]interface{}{})
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹150000" {
		t.Errorf("expected no change ₹150000, got %s", income.ScenarioValue)
	}
}

func TestExecute_WrongTypeAmount(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"amount": "not a number"}
	out, _ := e.Execute(defaultInputs(), STSalaryIncrease, params)
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹150000" {
		t.Errorf("expected no change ₹150000, got %s", income.ScenarioValue)
	}
}

func TestExecute_WrongTypePercent(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"percent": "100"}
	out, _ := e.Execute(defaultInputs(), STMarketCrash, params)
	// string "100" fails float64 type assertion -> no-op
	sr := out.DeltaSummary[2]
	if sr.ScenarioValue != "0.3" {
		t.Errorf("expected savings rate 0.3 unchanged, got %s", sr.ScenarioValue)
	}
}

func TestExecute_NilParams(t *testing.T) {
	e := NewEngine()
	in := defaultInputs()
	out, err := e.Execute(in, STSalaryIncrease, nil)
	if err != nil {
		t.Fatal(err)
	}
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹150000" {
		t.Errorf("expected no change with nil params, got %s", income.ScenarioValue)
	}
}

// ---------------------------------------------------------------------------
// Edge cases: negative / zero / extreme values
// ---------------------------------------------------------------------------

func TestExecute_NegativeAmount(t *testing.T) {
	e := NewEngine()
	// Negative SIP increase should actually decrease income
	out, _ := e.Execute(defaultInputs(), STIncreaseSIP, map[string]interface{}{"amount": float64(-10_000)})
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹140000" {
		t.Errorf("expected ₹140000 (negative amount reduces income), got %s", income.ScenarioValue)
	}
}

func TestExecute_NegativePercentSalaryLoss(t *testing.T) {
	e := NewEngine()
	// Negative percent should increase income (1 - (-25)/100 = 1.25)
	out, _ := e.Execute(defaultInputs(), STSalaryLoss, map[string]interface{}{"percent": float64(-25.0)})
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹187500" {
		t.Errorf("expected ₹187500 (negative loss = gain), got %s", income.ScenarioValue)
	}
}

func TestExecute_NegativePercentMarketCrash(t *testing.T) {
	e := NewEngine()
	// Negative crash percent increases equity return
	out, _ := e.Execute(defaultInputs(), STMarketCrash, map[string]interface{}{"percent": float64(-3.0)})
	// EquityReturn becomes 12 - (-3) = 15 — no visible delta field change
	// But we verify the call doesn't panic
	if out.Status != SSCompleted {
		t.Errorf("expected completed, got %s", out.Status)
	}
}

func TestExecute_ZeroAmount(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STInvestmentLumpSum, map[string]interface{}{"amount": float64(0)})
	nw := out.DeltaSummary[0]
	if nw.ScenarioValue != "₹5000000" {
		t.Errorf("expected ₹5000000 for zero lump sum, got %s", nw.ScenarioValue)
	}
}

func TestExecute_ZeroPercentSalaryLoss(t *testing.T) {
	e := NewEngine()
	out, _ := e.Execute(defaultInputs(), STSalaryLoss, map[string]interface{}{"percent": float64(0)})
	income := out.DeltaSummary[1]
	if income.ScenarioValue != "₹150000" {
		t.Errorf("expected ₹150000 unchanged, got %s", income.ScenarioValue)
	}
}

func TestExecute_HugeAmount(t *testing.T) {
	e := NewEngine()
	params := map[string]interface{}{"amount": float64(1 << 40)} // ~1 trillion
	out, _ := e.Execute(defaultInputs(), STIncreaseSIP, params)
	income := out.DeltaSummary[1]
	want := "₹1249628096" // 150000 + 1099511627776 = 1099511777776... wait
	// 150000 + 1099511627776 = 1099511777776
	want = "₹1099511777776"
	if income.ScenarioValue != want {
		t.Errorf("expected large income %s, got %s", want, income.ScenarioValue)
	}
}

// ---------------------------------------------------------------------------
// Multiple sims on the same engine instance
// ---------------------------------------------------------------------------

func TestExecute_ReuseEngine(t *testing.T) {
	e := NewEngine()

	out1, _ := e.Execute(defaultInputs(), STIncreaseSIP, map[string]interface{}{"amount": float64(10_000)})
	out2, _ := e.Execute(defaultInputs(), STDecreaseSIP, map[string]interface{}{"amount": float64(10_000)})

	if out1.DeltaSummary[1].ScenarioValue != "₹160000" {
		t.Errorf("first sim expected ₹160000, got %s", out1.DeltaSummary[1].ScenarioValue)
	}
	if out2.DeltaSummary[1].ScenarioValue != "₹140000" {
		t.Errorf("second sim expected ₹140000, got %s", out2.DeltaSummary[1].ScenarioValue)
	}
}

// ---------------------------------------------------------------------------
// Delta computation
// ---------------------------------------------------------------------------

func TestComputeDeltas_NoChange(t *testing.T) {
	e := NewEngine()
	base := defaultInputs()
	deltas := e.computeDeltas(base, base)

	for i, d := range deltas {
		if d.Direction != "NoChange" {
			t.Errorf("entry %d (%s) expected NoChange, got %s", i, d.Metric, d.Direction)
		}
		if d.Delta != "0" {
			t.Errorf("entry %d (%s) expected delta 0, got %s", i, d.Metric, d.Delta)
		}
		if d.BaselineValue != d.ScenarioValue {
			t.Errorf("entry %d (%s) baseline %s != scenario %s", i, d.Metric, d.BaselineValue, d.ScenarioValue)
		}
	}
}

func TestComputeDeltas_AllDifferent(t *testing.T) {
	e := NewEngine()
	base := defaultInputs()
	scenario := defaultInputs()
	scenario.NetWorth = 6_000_000
	scenario.MonthlyIncome = 200_000
	scenario.SavingsRate = 0.5
	scenario.MonthlyEMI = 0

	deltas := e.computeDeltas(base, scenario)
	if len(deltas) != 4 {
		t.Fatalf("expected 4 deltas, got %d", len(deltas))
	}
	// 0: Net Worth
	if deltas[0].BaselineValue != "₹5000000" || deltas[0].ScenarioValue != "₹6000000" {
		t.Errorf("net worth delta mismatch: %s -> %s", deltas[0].BaselineValue, deltas[0].ScenarioValue)
	}
	// 1: Monthly Income
	if deltas[1].BaselineValue != "₹150000" || deltas[1].ScenarioValue != "₹200000" {
		t.Errorf("income delta mismatch: %s -> %s", deltas[1].BaselineValue, deltas[1].ScenarioValue)
	}
	// 2: Savings Rate
	if deltas[2].BaselineValue != "0.3" || deltas[2].ScenarioValue != "0.5" {
		t.Errorf("savings rate delta mismatch: %s -> %s", deltas[2].BaselineValue, deltas[2].ScenarioValue)
	}
	// 3: Monthly EMI
	if deltas[3].BaselineValue != "₹25000" || deltas[3].ScenarioValue != "₹0" {
		t.Errorf("EMI delta mismatch: %s -> %s", deltas[3].BaselineValue, deltas[3].ScenarioValue)
	}

	for _, d := range deltas {
		if d.Direction != "Change" {
			t.Errorf("expected Change for %s, got %s", d.Metric, d.Direction)
		}
	}
}

func TestComputeDeltas_Ordering(t *testing.T) {
	e := NewEngine()
	deltas := e.computeDeltas(defaultInputs(), defaultInputs())
	expectedOrder := []string{"Net Worth", "Monthly Income", "Savings Rate", "Monthly EMI"}
	for i, d := range deltas {
		if d.Metric != expectedOrder[i] {
			t.Errorf("position %d expected %s, got %s", i, expectedOrder[i], d.Metric)
		}
	}
}

// ---------------------------------------------------------------------------
// Delta entry helper
// ---------------------------------------------------------------------------

func TestDelta(t *testing.T) {
	e := NewEngine()

	t.Run("equal", func(t *testing.T) {
		d := e.delta("Metric", "₹100", "₹100")
		if d.Direction != "NoChange" || d.Delta != "0" {
			t.Errorf("expected NoChange/0, got %s/%s", d.Direction, d.Delta)
		}
	})

	t.Run("different", func(t *testing.T) {
		d := e.delta("Metric", "₹100", "₹200")
		if d.Direction != "Change" || d.Delta != "₹200" {
			t.Errorf("expected Change/₹200, got %s/%s", d.Direction, d.Delta)
		}
	})
}

// ---------------------------------------------------------------------------
// generateSummary
// ---------------------------------------------------------------------------

func TestGenerateSummary(t *testing.T) {
	e := NewEngine()
	deltas := []DeltaEntry{
		{Metric: "A", BaselineValue: "₹100", ScenarioValue: "₹200", Delta: "₹200", Direction: "Change"},
	}
	result := e.generateSummary(deltas)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0] != deltas[0] {
		t.Error("generateSummary should return the same slice")
	}
}

// ---------------------------------------------------------------------------
// SimStatus constants
// ---------------------------------------------------------------------------

func TestSimStatusValues(t *testing.T) {
	if SSDraft != "Draft" { t.Error("SSDraft != Draft") }
	if SSRunning != "Running" { t.Error("SSRunning != Running") }
	if SSCompleted != "Completed" { t.Error("SSCompleted != Completed") }
	if SSFailed != "Failed" { t.Error("SSFailed != Failed") }
	if SSExpired != "Expired" { t.Error("SSExpired != Expired") }
	if SSArchived != "Archived" { t.Error("SSArchived != Archived") }
	if SSHistorical != "Historical" { t.Error("SSHistorical != Historical") }
}

// ---------------------------------------------------------------------------
// SimType constants
// ---------------------------------------------------------------------------

func TestSimTypeValues(t *testing.T) {
	if STIncreaseSIP != "IncreaseSIP" { t.Error("STIncreaseSIP != IncreaseSIP") }
	if STDecreaseSIP != "DecreaseSIP" { t.Error("STDecreaseSIP != DecreaseSIP") }
	if STEarlyLoanClosure != "EarlyLoanClosure" { t.Error("STEarlyLoanClosure != EarlyLoanClosure") }
	if STExtraEMI != "ExtraEMI" { t.Error("STExtraEMI != ExtraEMI") }
	if STSalaryIncrease != "SalaryIncrease" { t.Error("STSalaryIncrease != SalaryIncrease") }
	if STSalaryLoss != "SalaryLoss" { t.Error("STSalaryLoss != SalaryLoss") }
	if STJobChange != "JobChange" { t.Error("STJobChange != JobChange") }
	if STMarriage != "Marriage" { t.Error("STMarriage != Marriage") }
	if STChild != "Child" { t.Error("STChild != Child") }
	if STHousePurchase != "HousePurchase" { t.Error("STHousePurchase != HousePurchase") }
	if STVehiclePurchase != "VehiclePurchase" { t.Error("STVehiclePurchase != VehiclePurchase") }
	if STInvestmentLumpSum != "InvestmentLumpSum" { t.Error("STInvestmentLumpSum != InvestmentLumpSum") }
	if STExpenseReduction != "ExpenseReduction" { t.Error("STExpenseReduction != ExpenseReduction") }
	if STGoalDelay != "GoalDelay" { t.Error("STGoalDelay != GoalDelay") }
	if STGoalPriorityChange != "GoalPriorityChange" { t.Error("STGoalPriorityChange != GoalPriorityChange") }
	if STEmergency != "Emergency" { t.Error("STEmergency != Emergency") }
	if STMarketCrash != "MarketCrash" { t.Error("STMarketCrash != MarketCrash") }
	if STMarketBoom != "MarketBoom" { t.Error("STMarketBoom != MarketBoom") }
	if STInflationChange != "InflationChange" { t.Error("STInflationChange != InflationChange") }
	if STRetirementAgeChange != "RetirementAgeChange" { t.Error("STRetirementAgeChange != RetirementAgeChange") }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func verifyNoChange(t *testing.T, out *SimulationOutput) {
	t.Helper()
	for _, d := range out.DeltaSummary {
		if d.Direction != "NoChange" {
			t.Errorf("expected NoChange for %s, got %s", d.Metric, d.Direction)
		}
		if d.BaselineValue != d.ScenarioValue {
			t.Errorf("%s: baseline %s != scenario %s", d.Metric, d.BaselineValue, d.ScenarioValue)
		}
	}
}
