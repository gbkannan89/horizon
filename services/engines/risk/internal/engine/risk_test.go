package engine

import (
	"testing"
)

func eq(t *testing.T, want, got int) {
	t.Helper()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func eqLevel(t *testing.T, want, got RiskLevel) {
	t.Helper()
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func eq64(t *testing.T, want, got int64) {
	t.Helper()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

// ---------------------------------------------------------------------------
// Unit: scorePct
// ---------------------------------------------------------------------------

func TestScorePct(t *testing.T) {
	e := &Engine{}

	t.Run("bands", func(t *testing.T) {
		threshold := 70.0
		eq(t, 85, e.scorePct(threshold*2, threshold))
		eq(t, 85, e.scorePct(threshold*2+1, threshold))
		eq(t, 70, e.scorePct(threshold*1.5, threshold))
		eq(t, 70, e.scorePct(threshold*2-0.01, threshold))
		eq(t, 55, e.scorePct(threshold*1.0, threshold))
		eq(t, 55, e.scorePct(threshold*1.5-0.01, threshold))
		eq(t, 35, e.scorePct(threshold*0.5, threshold))
		eq(t, 35, e.scorePct(threshold*1.0-0.01, threshold))
		eq(t, 15, e.scorePct(0, threshold))
		eq(t, 15, e.scorePct(threshold*0.5-0.01, threshold))
	})

	t.Run("zero threshold", func(t *testing.T) {
		// actual >= 0*2 → true, always hits first branch
		eq(t, 85, e.scorePct(0, 0))
		eq(t, 85, e.scorePct(100, 0))
	})
}

// ---------------------------------------------------------------------------
// Unit: scoreConc
// ---------------------------------------------------------------------------

func TestScoreConc(t *testing.T) {
	e := &Engine{}

	t.Run("bands", func(t *testing.T) {
		threshold := 15.0
		eq(t, 90, e.scoreConc(threshold*2, threshold))
		eq(t, 90, e.scoreConc(threshold*2+5, threshold))
		eq(t, 75, e.scoreConc(threshold*1.5, threshold))
		eq(t, 60, e.scoreConc(threshold*1.0, threshold))
		eq(t, 35, e.scoreConc(threshold*0.5, threshold))
		eq(t, 15, e.scoreConc(0, threshold))
		eq(t, 15, e.scoreConc(threshold*0.5-0.01, threshold))
	})

	t.Run("zero threshold", func(t *testing.T) {
		// actual >= 0*2 → true, always hits first branch
		eq(t, 90, e.scoreConc(0, 0))
		eq(t, 90, e.scoreConc(10, 0))
	})
}

// ---------------------------------------------------------------------------
// Unit: scoreInverse
// ---------------------------------------------------------------------------

func TestScoreInverse(t *testing.T) {
	e := &Engine{}

	t.Run("bands", func(t *testing.T) {
		target := 3.0
		eq(t, 15, e.scoreInverse(target*1.5, target))
		eq(t, 15, e.scoreInverse(target*2, target))
		eq(t, 35, e.scoreInverse(target*1.0, target))
		eq(t, 35, e.scoreInverse(target*1.5-0.01, target))
		eq(t, 55, e.scoreInverse(target*0.75, target))
		eq(t, 70, e.scoreInverse(target*0.5, target))
		eq(t, 85, e.scoreInverse(0, target))
		eq(t, 85, e.scoreInverse(target*0.5-0.01, target))
	})

	t.Run("zero target", func(t *testing.T) {
		eq(t, 50, e.scoreInverse(100, 0))
		eq(t, 50, e.scoreInverse(0, 0))
	})

	t.Run("negative target", func(t *testing.T) {
		eq(t, 50, e.scoreInverse(100, -1))
	})
}

// ---------------------------------------------------------------------------
// Unit: scorePctSimple
// ---------------------------------------------------------------------------

func TestScorePctSimple(t *testing.T) {
	e := &Engine{}
	eq(t, 80, e.scorePctSimple(40))
	eq(t, 80, e.scorePctSimple(50))
	eq(t, 65, e.scorePctSimple(30))
	eq(t, 65, e.scorePctSimple(39))
	eq(t, 45, e.scorePctSimple(15))
	eq(t, 45, e.scorePctSimple(29))
	eq(t, 25, e.scorePctSimple(5))
	eq(t, 25, e.scorePctSimple(14))
	eq(t, 10, e.scorePctSimple(0))
	eq(t, 10, e.scorePctSimple(4))
}

// ---------------------------------------------------------------------------
// Unit: toLevel
// ---------------------------------------------------------------------------

func TestToLevel(t *testing.T) {
	e := &Engine{}
	tests := []struct {
		score int
		want  RiskLevel
	}{
		{100, RLCritical}, {80, RLCritical}, {88, RLCritical},
		{79, RLElevated}, {60, RLElevated}, {65, RLElevated},
		{59, RLModerate}, {40, RLModerate}, {45, RLModerate},
		{39, RLLow}, {20, RLLow}, {25, RLLow},
		{19, RLMinimal}, {0, RLMinimal}, {-1, RLMinimal},
	}
	for _, tc := range tests {
		got := e.toLevel(tc.score)
		if tc.want != got {
			t.Errorf("toLevel(%d): want %s, got %s", tc.score, tc.want, got)
		}
	}
}

// ---------------------------------------------------------------------------
// Unit: crashImpact
// ---------------------------------------------------------------------------

func TestCrashImpact(t *testing.T) {
	e := &Engine{}
	eq64(t, -300_000, e.crashImpact(1_000_000, 0.30))
	eq64(t, 0, e.crashImpact(0, 0.50))
	eq64(t, -50_000, e.crashImpact(100_000, 0.50))
}

// ---------------------------------------------------------------------------
// Unit: runStressTests
// ---------------------------------------------------------------------------

func TestRunStressTests(t *testing.T) {
	e := &Engine{}
	inputs := Inputs{PortfolioValue: 1_000_000}
	results := e.runStressTests(inputs)

	if len(results) != 5 {
		t.Fatalf("expected 5 stress tests, got %d", len(results))
	}

	expected := []struct {
		scenario     string
		impact       int64
		recoveryDays int
	}{
		{"MarketCrash_30", -300_000, 730},
		{"ExtendedBearMarket", -450_000, 1095},
		{"InterestRateShock_3pct", -100_000, 180},
		{"InflationSpike_8pct", -150_000, 365},
		{"IncomeLoss_6mo", -500_000, 180},
	}
	for i, exp := range expected {
		if results[i].Scenario != exp.scenario {
			t.Errorf("stress[%d].Scenario: want %s, got %s", i, exp.scenario, results[i].Scenario)
		}
		if results[i].Impact != exp.impact {
			t.Errorf("stress[%d].Impact: want %d, got %d", i, exp.impact, results[i].Impact)
		}
		if results[i].RecoveryDays != exp.recoveryDays {
			t.Errorf("stress[%d].RecoveryDays: want %d, got %d", i, exp.recoveryDays, results[i].RecoveryDays)
		}
	}

	t.Run("zero portfolio", func(t *testing.T) {
		results := e.runStressTests(Inputs{PortfolioValue: 0})
		eq64(t, 0, results[0].Impact)
		eq64(t, 0, results[1].Impact)
		eq64(t, 0, results[2].Impact)
		eq64(t, 0, results[3].Impact)
		eq64(t, -500_000, results[4].Impact)
	})
}

// ---------------------------------------------------------------------------
// Unit: computeIndicators – verify all 18 indicators
// ---------------------------------------------------------------------------

func TestComputeIndicators(t *testing.T) {
	e := &Engine{}

	// Inputs where each indicator is at a distinct level.
	inputs := Inputs{
		SingleAssetConcPct:     30,   // scoreConc(30, 15)  → 90
		SectorConcPct:          45,   // scoreConc(45, 30)  → 75
		InstitutionConcPct:     40,   // scoreConc(40, 40)  → 60
		EquityExposureRatio:    35,   // scorePct(35, 70)   → 15
		CurrencyExposurePct:    10,   // scorePct(10, 20)   → 15
		FloatingRateLiabPct:    25,   // scorePct(25, 50)   → 15
		CashFlowCoverageMonths: 0,    // scoreInverse(0, 3) → 85
		LiquidityRatio:         0.75, // scoreInverse(0.75, 1.5) → 55
		DTIRatio:               20,   // scorePct(20, 40)   → 15
		CreditUtilizationPct:   25,   // scorePct(25, 50)   → 15
		RetirementFundingRatio: 40,   // scoreInverse(40, 80)→ 85
		WithdrawalRate:         2,    // scorePct(2, 4)     → 15
		GoalFundingGapPct:      10,   // scorePct(10, 20)   → 15
		IncomeVolatility:       1.0,  // min(1*50, 100)    → 50
		InsuranceGapLife:       5,    // scoreInverse(5,10) → 70
		InsuranceGapHealth:     50,   // scorePct(50, 50)   → 55
		SavingsConsistencyPct:  105,  // scoreInverse(105,70)→ 15
		AllocationDriftPct:     5,    // scorePct(5, 10)    → 15
	}

	indicators := e.computeIndicators(inputs)

	if len(indicators) != 18 {
		t.Fatalf("expected 18 indicators, got %d", len(indicators))
	}

	expected := []struct {
		name          string
		score         int
		level         RiskLevel
		warning, crit int
	}{
		{"SingleAssetConcentration", 90, RLCritical, 60, 80},
		{"SectorConcentration", 75, RLCritical, 60, 80},
		{"InstitutionConcentration", 60, RLCritical, 60, 80},
		{"EquityExposureRatio", 35, RLCritical, 60, 80},
		{"CurrencyExposure", 35, RLCritical, 60, 80},
		{"InterestRateExposure", 35, RLCritical, 60, 80},
		{"CashFlowCoverage", 85, RLCritical, 60, 80},
		{"LiquidityRatio", 70, RLCritical, 60, 80},
		{"DTIRatio", 35, RLCritical, 60, 80},
		{"CreditUtilization", 35, RLCritical, 60, 80},
		{"RetirementFundingRatio", 70, RLCritical, 60, 80},
		{"WithdrawalRate", 35, RLCritical, 60, 80},
		{"GoalFundingGapSeverity", 35, RLCritical, 60, 80},
		{"IncomeStability", 50, RLCritical, 60, 80},
		{"InsuranceGapLife", 70, RLCritical, 60, 80},
		{"InsuranceGapHealth", 55, RLCritical, 60, 80},
		{"SavingsConsistency", 15, RLCritical, 60, 80},
		{"AllocationDrift", 35, RLCritical, 60, 80},
	}

	for i, exp := range expected {
		if indicators[i].Name != exp.name {
			t.Errorf("indicator[%d].Name: want %s, got %s", i, exp.name, indicators[i].Name)
		}
		if indicators[i].Score != exp.score {
			t.Errorf("indicator[%d] %s.Score: want %d, got %d", i, exp.name, exp.score, indicators[i].Score)
		}
		if indicators[i].Level != exp.level {
			t.Errorf("indicator[%d] %s.Level: want %s, got %s", i, exp.name, exp.level, indicators[i].Level)
		}
		if indicators[i].Warning != exp.warning {
			t.Errorf("indicator[%d] %s.Warning: want %d, got %d", i, exp.name, exp.warning, indicators[i].Warning)
		}
		if indicators[i].Critical != exp.crit {
			t.Errorf("indicator[%d] %s.Critical: want %d, got %d", i, exp.name, exp.crit, indicators[i].Critical)
		}
	}
}

// ---------------------------------------------------------------------------
// Unit: computeCategoryScores
// ---------------------------------------------------------------------------

func TestComputeCategoryScores(t *testing.T) {
	e := &Engine{}
	indicators := []RiskIndicator{
		{"SingleAssetConcentration", 90, RLCritical, 60, 80},
		{"SectorConcentration", 75, RLCritical, 60, 80},
		{"EquityExposureRatio", 15, RLCritical, 60, 80},
		{"CashFlowCoverage", 85, RLCritical, 60, 80},
		{"DTIRatio", 15, RLCritical, 60, 80},
		{"RetirementFundingRatio", 85, RLCritical, 60, 80},
		{"GoalFundingGapSeverity", 15, RLCritical, 60, 80},
		{"InsuranceGapLife", 70, RLCritical, 60, 80},
		{"SavingsConsistency", 15, RLCritical, 60, 80},
	}

	cats := e.computeCategoryScores(indicators)

	// Verify all 8 categories are present.
	catMap := map[string]bool{}
	for _, c := range cats {
		catMap[c.Category] = true
	}
	expectedCats := []string{"Concentration", "Market", "Liquidity", "Credit", "Longevity", "LifeEvent", "Insurance", "Behavioural"}
	for _, cat := range expectedCats {
		if !catMap[cat] {
			t.Errorf("missing category: %s", cat)
		}
	}

	// Verify known averages.
	for _, c := range cats {
		switch c.Category {
		case "Concentration":
			if c.Score != 82 { // (90+75+?wait, InstitutionConcentration not included) ... actually only SA=90, Sector=75 → avg=82
				t.Errorf("Concentration avg: want 82, got %d", c.Score)
			}
		case "Market":
			if c.Score != 15 {
				t.Errorf("Market avg: want 15, got %d", c.Score)
			}
		case "Liquidity":
			if c.Score != 85 {
				t.Errorf("Liquidity avg: want 85, got %d", c.Score)
			}
		case "Credit":
			if c.Score != 15 {
				t.Errorf("Credit avg: want 15, got %d", c.Score)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Unit: computeComposite
// ---------------------------------------------------------------------------

func TestComputeComposite(t *testing.T) {
	e := &Engine{}

	t.Run("average of category scores", func(t *testing.T) {
		cats := []RiskCategoryScore{
			{"Concentration", 82, RLCritical},
			{"Market", 15, RLLow},
			{"Liquidity", 85, RLCritical},
			{"Credit", 15, RLLow},
			{"Longevity", 50, RLModerate},
			{"LifeEvent", 32, RLLow},
			{"Insurance", 62, RLElevated},
			{"Behavioural", 15, RLLow},
		}
		// (82+15+85+15+50+32+62+15) = 356 / 8 = 44
		eq(t, 44, e.computeComposite(cats))
	})

	t.Run("empty categories", func(t *testing.T) {
		eq(t, 0, e.computeComposite(nil))
		eq(t, 0, e.computeComposite([]RiskCategoryScore{}))
	})

	t.Run("single category", func(t *testing.T) {
		cats := []RiskCategoryScore{{"Concentration", 75, RLElevated}}
		eq(t, 75, e.computeComposite(cats))
	})
}

// ---------------------------------------------------------------------------
// Integration: Execute – full end-to-end
// ---------------------------------------------------------------------------

func TestExecute(t *testing.T) {
	e := NewEngine()

	inputs := Inputs{
		SingleAssetConcPct:     30,
		SectorConcPct:          45,
		InstitutionConcPct:     40,
		EquityExposureRatio:    35,
		CurrencyExposurePct:    10,
		FloatingRateLiabPct:    25,
		CashFlowCoverageMonths: 0,
		LiquidityRatio:         0.75,
		DTIRatio:               20,
		CreditUtilizationPct:   25,
		DepositConcentration:   0,
		InflationExposedPct:    0,
		RetirementFundingRatio: 40,
		WithdrawalRate:         2,
		GoalFundingGapPct:      10,
		IncomeVolatility:       1.0,
		InsuranceGapLife:       5,
		InsuranceGapHealth:     50,
		SavingsConsistencyPct:  105,
		AllocationDriftPct:     5,
		SequenceRiskScore:      0,
		PortfolioValue:         1_000_000,
	}

	output := e.Execute(inputs)

	// Verify top-level structure.
	if output.CompositeScore == 0 {
		t.Error("composite score should not be zero")
	}
	if output.CompositeLevel == "" {
		t.Error("composite level should be set")
	}
	if len(output.Indicators) != 18 {
		t.Errorf("expected 18 indicators, got %d", len(output.Indicators))
	}
	if len(output.CategoryScores) != 8 {
		t.Errorf("expected 8 category scores, got %d", len(output.CategoryScores))
	}
	if len(output.StressTests) != 5 {
		t.Errorf("expected 5 stress tests, got %d", len(output.StressTests))
	}
	if output.Confidence != "High" {
		t.Errorf("confidence should be High, got %s", output.Confidence)
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestExecute_EdgeCases(t *testing.T) {
	e := NewEngine()

	t.Run("all zeros", func(t *testing.T) {
		output := e.Execute(Inputs{})
		// Many inverse indicators default to 85, pushing composite into Moderate.
		if output.CompositeLevel == "" {
			t.Error("composite level should be set")
		}
		if len(output.Indicators) != 18 {
			t.Errorf("expected 18 indicators, got %d", len(output.Indicators))
		}
	})

	t.Run("maximum risk", func(t *testing.T) {
		inputs := Inputs{
			SingleAssetConcPct:     100,
			SectorConcPct:          100,
			InstitutionConcPct:     100,
			EquityExposureRatio:    100,
			CurrencyExposurePct:    100,
			FloatingRateLiabPct:    100,
			CashFlowCoverageMonths: 0,
			LiquidityRatio:         0,
			DTIRatio:               100,
			CreditUtilizationPct:   100,
			RetirementFundingRatio: 0,
			WithdrawalRate:         100,
			GoalFundingGapPct:      100,
			IncomeVolatility:       2.0,
			InsuranceGapLife:       0,
			InsuranceGapHealth:     0,
			SavingsConsistencyPct:  0,
			AllocationDriftPct:     100,
			PortfolioValue:         1_000_000,
		}
		output := e.Execute(inputs)
		if output.CompositeLevel != RLCritical {
			t.Errorf("max risk should yield Critical level, got %s", output.CompositeLevel)
		}
		if output.CompositeScore < 80 {
			t.Errorf("max risk composite should be >= 80, got %d", output.CompositeScore)
		}
	})

	t.Run("negative values", func(t *testing.T) {
		inputs := Inputs{
			SingleAssetConcPct: -1,
			LiquidityRatio:     -0.5,
		}
		output := e.Execute(inputs)
		// Should not panic, should produce valid output.
		if output == nil {
			t.Fatal("output should not be nil")
		}
		if len(output.Indicators) != 18 {
			t.Errorf("expected 18 indicators, got %d", len(output.Indicators))
		}
	})

	t.Run("large portfolio value", func(t *testing.T) {
		inputs := Inputs{PortfolioValue: 9_007_199_254_740_991} // near max int64
		output := e.Execute(inputs)
		if output == nil {
			t.Fatal("output should not be nil")
		}
		// Stress tests should handle large values without overflow.
		if output.StressTests[0].Impact > 0 {
			t.Errorf("crash impact should be negative, got %d", output.StressTests[0].Impact)
		}
	})

	t.Run("IncomeVolatility above 2.0 caps at 100", func(t *testing.T) {
		inputs := Inputs{IncomeVolatility: 5.0}
		indicators := e.computeIndicators(inputs)
		found := false
		for _, ind := range indicators {
			if ind.Name == "IncomeStability" {
				found = true
				if ind.Score != 100 {
					t.Errorf("IncomeStability with vol=5.0 should cap at 100, got %d", ind.Score)
				}
			}
		}
		if !found {
			t.Error("IncomeStability indicator not found")
		}
	})

	t.Run("Negative IncomeVolatility", func(t *testing.T) {
		inputs := Inputs{IncomeVolatility: -0.5}
		indicators := e.computeIndicators(inputs)
		for _, ind := range indicators {
			if ind.Name == "IncomeStability" {
				// math.Min(-25, 100) = -25 (float64), then int conversion
				if ind.Score != -25 {
					t.Errorf("IncomeStability with vol=-0.5 should be -25, got %d", ind.Score)
				}
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Verify the Execute output consistency: category score levels must match
// toLevel(category.Score)
// ---------------------------------------------------------------------------

func TestCategoryScoreLevelConsistency(t *testing.T) {
	e := NewEngine()
	inputs := Inputs{
		SingleAssetConcPct:     30,
		SectorConcPct:          45,
		InstitutionConcPct:     40,
		CashFlowCoverageMonths: 0,
		LiquidityRatio:         0.75,
		EquityExposureRatio:    35,
	}
	output := e.Execute(inputs)
	for _, cat := range output.CategoryScores {
		wantLevel := e.toLevel(cat.Score)
		if cat.Level != wantLevel {
			t.Errorf("category %s: score %d maps to level %s but got %s",
				cat.Category, cat.Score, wantLevel, cat.Level)
		}
	}
}
