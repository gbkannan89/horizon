package engine

import (
	"testing"
)

func TestFmtMoney(t *testing.T) {
	tests := []struct {
		name string
		v    int64
		want string
	}{
		{"crore", 12300000, "₹1.23Cr"},
		{"crore_exact", 10000000, "₹1.00Cr"},
		{"lakh", 500000, "₹5.00L"},
		{"lakh_frac", 750000, "₹7.50L"},
		{"thousand", 99000, "₹99000"},
		{"zero", 0, "₹0"},
		{"negative", -5000, "₹-5000"},
		{"small", 99, "₹99"},
		{"boundary_lakh", 99999, "₹99999"},
		{"boundary_lakh_plus", 100000, "₹1.00L"},
		{"boundary_crore_minus", 9999999, "₹100.00L"},
		{"boundary_crore", 10000000, "₹1.00Cr"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmtMoney(tt.v); got != tt.want {
				t.Errorf("fmtMoney(%d) = %q, want %q", tt.v, got, tt.want)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("CardSummary", func(t *testing.T) {
		if got := CardSummary(5); got != "5 items" {
			t.Errorf("CardSummary = %q", got)
		}
		if got := CardSummary(0); got != "0 items" {
			t.Errorf("CardSummary 0 = %q", got)
		}
	})

	t.Run("CardRiskSummary", func(t *testing.T) {
		if got := CardRiskSummary(3, "Moderate"); got != "Score: 3 — Moderate" {
			t.Errorf("CardRiskSummary = %q", got)
		}
	})

	t.Run("CardHealthSummary", func(t *testing.T) {
		if got := CardHealthSummary(85, "A"); got != "Health: 85 (A)" {
			t.Errorf("CardHealthSummary = %q", got)
		}
	})

	t.Run("CardEventSummary", func(t *testing.T) {
		if got := CardEventSummary(12); got != "12 events" {
			t.Errorf("CardEventSummary = %q", got)
		}
		if got := CardEventSummary(1); got != "1 events" {
			t.Errorf("CardEventSummary 1 = %q", got)
		}
	})
}

func TestNewComposer(t *testing.T) {
	c := NewComposer()
	if c == nil {
		t.Fatal("NewComposer() returned nil")
	}
}

func TestBuildDashboard(t *testing.T) {
	c := NewComposer()

	t.Run("overview mode minimal", func(t *testing.T) {
		in := Inputs{
			Mode:              PMOverview,
			NetWorthProjected: 5000000,
			GoalsOnTrack:      3,
			TotalGoals:        5,
			FundingGap:        200000,
			RiskScore:         4,
			RiskLevel:         "High",
			HealthScore:       72,
		}
		dash := c.BuildDashboard(in)
		if dash.Mode != PMOverview {
			t.Errorf("Mode = %q, want %q", dash.Mode, PMOverview)
		}
		if dash.CurrentPlan == nil {
			t.Fatal("CurrentPlan is nil")
		}
		if dash.CurrentPlan.PlanID != "baseline" {
			t.Errorf("PlanID = %q", "baseline")
		}
		if dash.CurrentPlan.NetWorthEnd != 5000000 {
			t.Errorf("NetWorthEnd = %d", dash.CurrentPlan.NetWorthEnd)
		}
		if len(dash.Cards) == 0 {
			t.Error("expected at least some cards")
		}
	})

	t.Run("all optional cards enabled", func(t *testing.T) {
		in := Inputs{
			Mode:              PMGoalPlanning,
			NetWorthProjected: 10000000,
			GoalsOnTrack:      5,
			TotalGoals:        5,
			IncomeProjected:   800000,
			ExpensesProjected: 600000,
			PortfolioValue:    3000000,
			FundingGap:        0,
			RiskScore:         2,
			RiskLevel:         "Low",
			MonthlyIncome:     100000,
			MonthlyExpenses:   80000,
			HasRecs:           true,
			RecCount:          3,
			HasOptimization:   true,
			OptCount:          5,
			HasSimulation:     true,
			SimCount:          4,
			EventCount:        7,
			UpcomingRecurring: 2,
		}
		dash := c.BuildDashboard(in)
		cardTypes := make(map[CardType]bool)
		for _, card := range dash.Cards {
			cardTypes[card.CardType] = true
		}
		expected := []CardType{
			CTPlanOverview, CTGoalPlanning, CTCashFlowPlan, CTRetirementPlan,
			CTInvestmentPlan, CTDebtPlan, CTProjection, CTRiskSum,
			CTRecommend, CTOptimization, CTSimulation, CTTimeline, CTRecurring,
		}
		for _, ct := range expected {
			if !cardTypes[ct] {
				t.Errorf("missing card type %q", ct)
			}
		}
	})

	t.Run("no optional cards when flags are false", func(t *testing.T) {
		in := Inputs{
			Mode:              PMRetirement,
			NetWorthProjected: 0,
			HasRecs:           false,
			HasOptimization:   false,
			HasSimulation:     false,
			EventCount:        0,
			UpcomingRecurring: 0,
		}
		dash := c.BuildDashboard(in)
		for _, card := range dash.Cards {
			switch card.CardType {
			case CTRecommend, CTOptimization, CTSimulation, CTTimeline, CTRecurring:
				t.Errorf("unexpected card %q with zero/disabled flags", card.CardType)
			}
		}
	})

	t.Run("scenarios forwarded", func(t *testing.T) {
		scenarios := []Scenario{
			{ScenarioID: "s1", Name: "Optimistic"},
			{ScenarioID: "s2", Name: "Pessimistic"},
		}
		in := Inputs{
			Mode:              PMScenarioPlan,
			NetWorthProjected: 0,
			Scenarios:         scenarios,
		}
		dash := c.BuildDashboard(in)
		if len(dash.Scenarios) != 2 {
			t.Fatalf("expected 2 scenarios, got %d", len(dash.Scenarios))
		}
		if dash.Scenarios[0].ScenarioID != "s1" {
			t.Errorf("first scenario id = %q", dash.Scenarios[0].ScenarioID)
		}
	})

	t.Run("cards have correct priority ordering", func(t *testing.T) {
		in := Inputs{
			Mode:              PMDebtPlan,
			NetWorthProjected: 100000,
			GoalsOnTrack:      2,
			TotalGoals:        2,
			HasRecs:           true,
			RecCount:          1,
		}
		dash := c.BuildDashboard(in)
		for i := 1; i < len(dash.Cards); i++ {
			if dash.Cards[i].Priority <= dash.Cards[i-1].Priority {
				t.Errorf("card %q (priority %d) should have higher priority than %q (priority %d)",
					dash.Cards[i].Title, dash.Cards[i].Priority,
					dash.Cards[i-1].Title, dash.Cards[i-1].Priority)
			}
		}
	})

	t.Run("plan overview card data", func(t *testing.T) {
		in := Inputs{
			Mode:              PMInvestmentPlan,
			NetWorthProjected: 25000000,
			GoalsOnTrack:      4,
			TotalGoals:        6,
		}
		dash := c.BuildDashboard(in)
		var planCard *Card
		for i, card := range dash.Cards {
			if card.CardType == CTPlanOverview {
				planCard = &dash.Cards[i]
				break
			}
		}
		if planCard == nil {
			t.Fatal("missing PlanOverview card")
		}
		data, ok := planCard.Data.(map[string]interface{})
		if !ok {
			t.Fatal("PlanOverview data is not map")
		}
		nw, _ := data["net_worth"].(int64)
		if nw != 25000000 {
			t.Errorf("net_worth = %v", data["net_worth"])
		}
		goals, _ := data["goals"].(string)
		if goals != "4/6" {
			t.Errorf("goals = %q", goals)
		}
	})
}

func TestBuildProjections(t *testing.T) {
	c := NewComposer()

	in := Inputs{
		NetWorthProjected: 5000000,
		IncomeProjected:   1200000,
		ExpensesProjected: 900000,
		PortfolioValue:    3000000,
	}

	p := c.BuildProjections(in)
	if p.NetWorthProjected != 5000000 {
		t.Errorf("NetWorthProjected = %d", p.NetWorthProjected)
	}
	if p.IncomeProjected != 1200000 {
		t.Errorf("IncomeProjected = %d", p.IncomeProjected)
	}
	if p.ExpensesProjected != 900000 {
		t.Errorf("ExpensesProjected = %d", p.ExpensesProjected)
	}
	if p.PortfolioValue != 3000000 {
		t.Errorf("PortfolioValue = %d", p.PortfolioValue)
	}
	if p.Confidence != "" {
		t.Errorf("Confidence should be empty, got %q", p.Confidence)
	}
}

func TestBuildComparison(t *testing.T) {
	c := NewComposer()

	t.Run("higher is better deltas", func(t *testing.T) {
		a := PlanSummary{NetWorthEnd: 1000000, GoalsOnTrack: 3, TotalGoals: 5, FundingGap: 200000, RiskScore: 5, HealthScore: 70}
		b := PlanSummary{NetWorthEnd: 1500000, GoalsOnTrack: 4, TotalGoals: 5, FundingGap: 100000, RiskScore: 3, HealthScore: 80}
		comp := c.BuildComparison(a, b)

		if comp.PlanA != a || comp.PlanB != b {
			t.Error("plans not preserved")
		}
		if len(comp.Deltas) != 5 {
			t.Fatalf("expected 5 deltas, got %d", len(comp.Deltas))
		}

		check := func(metric, dir string) {
			for _, d := range comp.Deltas {
				if d.Metric == metric {
					if d.Direction != dir {
						t.Errorf("delta %q direction = %q, want %q", metric, d.Direction, dir)
					}
					return
				}
			}
			t.Errorf("delta metric %q not found", metric)
		}
		check("Net Worth", "better")
		check("Goals On Track", "better")
		check("Funding Gap", "better")
		check("Risk Score", "better")
		check("Health Score", "better")
	})

	t.Run("lower is better deltas", func(t *testing.T) {
		a := PlanSummary{NetWorthEnd: 1000000, FundingGap: 500000, RiskScore: 8}
		b := PlanSummary{NetWorthEnd: 500000, FundingGap: 600000, RiskScore: 5}
		comp := c.BuildComparison(a, b)

		for _, d := range comp.Deltas {
			switch d.Metric {
			case "Net Worth":
				if d.Direction != "worse" {
					t.Errorf("Net Worth should be worse, got %q", d.Direction)
				}
			case "Funding Gap":
				if d.Direction != "worse" {
					t.Errorf("Funding Gap should be worse (higher gap), got %q", d.Direction)
				}
			case "Risk Score":
				if d.Direction != "better" {
					t.Errorf("Risk Score should be better (lower score), got %q", d.Direction)
				}
			}
		}
	})

	t.Run("unchanged", func(t *testing.T) {
		a := PlanSummary{NetWorthEnd: 1000000, GoalsOnTrack: 3, FundingGap: 200000, RiskScore: 5, HealthScore: 70}
		b := PlanSummary{NetWorthEnd: 1000000, GoalsOnTrack: 3, FundingGap: 200000, RiskScore: 5, HealthScore: 70}
		comp := c.BuildComparison(a, b)

		check := func(metric, dir string) {
			for _, d := range comp.Deltas {
				if d.Metric == metric {
					if d.Direction != dir {
						t.Errorf("delta %q direction = %q, want %q", metric, d.Direction, dir)
					}
					return
				}
			}
			t.Errorf("delta metric %q not found", metric)
		}
		check("Net Worth", "unchanged")
		check("Goals On Track", "unchanged")
		check("Funding Gap", "unchanged")
		check("Risk Score", "unchanged")
		check("Health Score", "unchanged")
	})

	t.Run("delta values formatting", func(t *testing.T) {
		a := PlanSummary{NetWorthEnd: 1000000}
		b := PlanSummary{NetWorthEnd: 2500000}
		comp := c.BuildComparison(a, b)

		var nwDelta *ComparisonDelta
		for _, d := range comp.Deltas {
			if d.Metric == "Net Worth" {
				nwDelta = &d
				break
			}
		}
		if nwDelta == nil {
			t.Fatal("missing Net Worth delta")
		}
		if nwDelta.Baseline != "₹10.00L" {
			t.Errorf("baseline = %q", nwDelta.Baseline)
		}
		if nwDelta.Alternative != "₹25.00L" {
			t.Errorf("alternative = %q", nwDelta.Alternative)
		}
		if nwDelta.Delta != "₹15.00L" {
			t.Errorf("delta = %q", nwDelta.Delta)
		}
	})
}

func TestDelta(t *testing.T) {
	c := NewComposer()

	t.Run("higher better", func(t *testing.T) {
		d := c.delta("test", 100, 200, true)
		if d.Direction != "better" {
			t.Errorf("expected better, got %q", d.Direction)
		}
		if d.Delta != "₹100" {
			t.Errorf("delta = %q", d.Delta)
		}
	})

	t.Run("higher better unchanged", func(t *testing.T) {
		d := c.delta("test", 100, 100, true)
		if d.Direction != "unchanged" {
			t.Errorf("expected unchanged, got %q", d.Direction)
		}
	})

	t.Run("higher better worse", func(t *testing.T) {
		d := c.delta("test", 200, 100, true)
		if d.Direction != "worse" {
			t.Errorf("expected worse, got %q", d.Direction)
		}
	})

	t.Run("lower better", func(t *testing.T) {
		d := c.delta("test", 200, 100, false)
		if d.Direction != "better" {
			t.Errorf("expected better, got %q", d.Direction)
		}
	})

	t.Run("lower better unchanged", func(t *testing.T) {
		d := c.delta("test", 100, 100, false)
		if d.Direction != "unchanged" {
			t.Errorf("expected unchanged, got %q", d.Direction)
		}
	})

	t.Run("lower better worse", func(t *testing.T) {
		d := c.delta("test", 100, 200, false)
		if d.Direction != "worse" {
			t.Errorf("expected worse, got %q", d.Direction)
		}
	})
}

func TestDeltaStr(t *testing.T) {
	c := NewComposer()

	t.Run("better", func(t *testing.T) {
		d := c.deltaStr("test", "3", "5", true)
		if d.Direction != "better" {
			t.Errorf("expected better, got %q", d.Direction)
		}
		if d.Alternative != d.Delta {
			t.Error("delta should equal alt for string deltas")
		}
	})

	t.Run("unchanged", func(t *testing.T) {
		d := c.deltaStr("test", "5", "5", true)
		if d.Direction != "unchanged" {
			t.Errorf("expected unchanged, got %q", d.Direction)
		}
	})

	t.Run("better_when_flag_true_and_alt_differs", func(t *testing.T) {
		d := c.deltaStr("test", "5", "3", true)
		if d.Direction != "better" {
			t.Errorf("expected better, got %q", d.Direction)
		}
	})

	t.Run("worse_when_flag_false", func(t *testing.T) {
		d := c.deltaStr("test", "5", "3", false)
		if d.Direction != "worse" {
			t.Errorf("expected worse, got %q", d.Direction)
		}
	})
}
