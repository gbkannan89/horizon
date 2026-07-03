package engine

import (
	"math"
	"testing"
)

func defaultConfig() Config {
	return Config{
		MaxRecommendationsPerCategory: 3,
		MaxTotalRecommendations:       20,
		MinRecommendationScore:        0.3,
		DriftThresholdForRebalance:    5.0,
		FundingGapThreshold:           100000,
		DebtInterestThreshold:         12.0,
		EmergencyFundTargetMonths:     6.0,
	}
}

func TestIdentifyGaps_AboveThreshold(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		GoalGaps:       map[string]int64{"Retirement": 500000, "Vacation": 50000},
		GoalImportance: map[string]string{"Retirement": "Mandatory", "Vacation": "Optional"},
	}
	recs := e.identifyGaps(inputs)
	if len(recs) != 1 {
		t.Fatalf("expected 1 gap recommendation (Vacation below threshold), got %d", len(recs))
	}
	if recs[0].Category != RCGoalFunding {
		t.Errorf("expected RCGoalFunding category, got %s", recs[0].Category)
	}
	if recs[0].Title != "Fund Retirement" {
		t.Errorf("expected 'Fund Retirement', got '%s'", recs[0].Title)
	}
	if len(recs[0].Evidence) != 2 {
		t.Errorf("expected 2 evidence items, got %d", len(recs[0].Evidence))
	}
}

func TestIdentifyGaps_Empty(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := e.identifyGaps(Inputs{GoalGaps: nil})
	if len(recs) != 0 {
		t.Errorf("expected 0 recs for nil gaps, got %d", len(recs))
	}
	recs = e.identifyGaps(Inputs{GoalGaps: map[string]int64{}})
	if len(recs) != 0 {
		t.Errorf("expected 0 recs for empty gaps, got %d", len(recs))
	}
}

func TestIdentifyGaps_BelowThreshold(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		GoalGaps:       map[string]int64{"PettyCash": 500},
		GoalImportance: map[string]string{"PettyCash": "Low"},
	}
	recs := e.identifyGaps(inputs)
	if len(recs) != 0 {
		t.Errorf("expected 0 recs below threshold, got %d", len(recs))
	}
}

func TestIdentifyGaps_ExactThreshold(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		GoalGaps:       map[string]int64{"College": 100000},
		GoalImportance: map[string]string{"College": "Mandatory"},
	}
	recs := e.identifyGaps(inputs)
	if len(recs) != 1 {
		t.Errorf("expected 1 rec at exact threshold, got %d", len(recs))
	}
}

func TestIdentifyOpportunities_DebtAboveThreshold(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		Debts: []DebtInfo{
			{Name: "CreditCard", Balance: 200000, InterestRate: 18.0},
			{Name: "LowInterestLoan", Balance: 500000, InterestRate: 8.0},
		},
	}
	recs := e.identifyOpportunities(inputs)
	found := false
	for _, r := range recs {
		if r.Category == RCDebt && r.Title == "Optimize CreditCard" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected debt recommendation for CreditCard (18%%) above threshold")
	}
	for _, r := range recs {
		if r.Category == RCDebt && r.Title == "Optimize LowInterestLoan" {
			t.Errorf("unexpected recommendation for debt below threshold")
		}
	}
}

func TestIdentifyOpportunities_PortfolioRebalance(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := e.identifyOpportunities(Inputs{PortfolioDrift: 8.0})
	found := false
	for _, r := range recs {
		if r.Category == RCPortfolioRebalancing {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected rebalance recommendation for drift 8.0%% > 5.0%%")
	}

	recs = e.identifyOpportunities(Inputs{PortfolioDrift: 3.0})
	for _, r := range recs {
		if r.Category == RCPortfolioRebalancing {
			t.Errorf("unexpected rebalance recommendation for drift 3.0%% < 5.0%%")
		}
	}
}

func TestIdentifyOpportunities_EmergencyFund(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := e.identifyOpportunities(Inputs{EmergencyMonths: 2.0})
	found := false
	for _, r := range recs {
		if r.Category == RCEmergencyFund {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected emergency fund recommendation for 2 months < 6")
	}

	recs = e.identifyOpportunities(Inputs{EmergencyMonths: 6.0})
	for _, r := range recs {
		if r.Category == RCEmergencyFund {
			t.Errorf("unexpected emergency fund recommendation at exactly 6 months")
		}
	}
}

func TestIdentifyOpportunities_OverspentCategories(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		OverspentCategories: map[string]int64{"Dining": 15000, "Entertainment": 5000, "Utilities": 0},
	}
	recs := e.identifyOpportunities(inputs)
	count := 0
	for _, r := range recs {
		if r.Category == RCCashFlow {
			count++
		}
	}
	if count < 2 {
		t.Errorf("expected at least 2 CashFlow overspend recs (Dining, Entertainment), got %d", count)
	}
	for _, r := range recs {
		if r.Category == RCCashFlow && r.Title == "Reduce spending in Dining" {
			if r.Summary != "Overspent by ₹15000" {
				t.Errorf("unexpected summary: %s", r.Summary)
			}
		}
	}
}

func TestIdentifyOpportunities_CashFlowSurplus(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := e.identifyOpportunities(Inputs{CashFlowSurplus: 30000})
	found := false
	for _, r := range recs {
		if r.Category == RCCashFlow && r.Title == "Allocate surplus" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Allocate surplus' recommendation for positive cash flow")
	}

	recs = e.identifyOpportunities(Inputs{CashFlowSurplus: 0})
	for _, r := range recs {
		if r.Category == RCCashFlow && r.Title == "Allocate surplus" {
			t.Errorf("unexpected surplus recommendation for zero surplus")
		}
	}
}

func TestIdentifyOpportunities_SavingsRate(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := e.identifyOpportunities(Inputs{SavingsRate: 15})
	found := false
	for _, r := range recs {
		if r.Category == RCSavings {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected savings recommendation for rate 15%% < 20%%")
	}

	recs = e.identifyOpportunities(Inputs{SavingsRate: 20})
	for _, r := range recs {
		if r.Category == RCSavings {
			t.Errorf("unexpected savings recommendation for rate exactly 20%%")
		}
	}

	recs = e.identifyOpportunities(Inputs{SavingsRate: 25})
	for _, r := range recs {
		if r.Category == RCSavings {
			t.Errorf("unexpected savings recommendation for rate 25%% >= 20%%")
		}
	}
}

func TestIdentifyOpportunities_EmptyInputs(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := e.identifyOpportunities(Inputs{
		PortfolioDrift:  0,
		EmergencyMonths: 6,
		SavingsRate:     20,
	})
	if len(recs) != 0 {
		t.Errorf("expected 0 recs for non-triggering inputs, got %d", len(recs))
	}
}

func TestIdentifyOpportunities_DefaultZeroValuesTrigger(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := e.identifyOpportunities(Inputs{})
	if len(recs) != 2 {
		t.Errorf("expected 2 recs (emergency fund + savings) from zero inputs, got %d", len(recs))
	}
}

func TestIdentifyOpportunities_AllTriggersSimultaneous(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		Debts:              []DebtInfo{{Name: "Card", Balance: 100000, InterestRate: 24.0}},
		PortfolioDrift:     10.0,
		EmergencyMonths:    1.0,
		OverspentCategories: map[string]int64{"Shopping": 20000},
		CashFlowSurplus:    50000,
		SavingsRate:        5,
	}
	recs := e.identifyOpportunities(inputs)
	categories := map[RecCategory]bool{}
	for _, r := range recs {
		categories[r.Category] = true
	}
	expected := []RecCategory{RCDebt, RCPortfolioRebalancing, RCEmergencyFund, RCCashFlow, RCSavings}
	for _, cat := range expected {
		if !categories[cat] {
			t.Errorf("missing category %s in combined opportunities", cat)
		}
	}
}

func TestScoreRecommendation_AllCategories(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{}

	tests := []struct {
		name     string
		category RecCategory
		want     ScoringDimensions
	}{
		{"GoalFunding", RCGoalFunding, ScoringDimensions{0.8, 0.7, 0.4, 0.8, 1.0, 0.7}},
		{"Debt", RCDebt, ScoringDimensions{0.7, 0.8, 0.5, 0.8, 0.6, 0.8}},
		{"PortfolioRebalancing", RCPortfolioRebalancing, ScoringDimensions{0.6, 0.5, 0.3, 0.7, 0.5, 0.6}},
		{"EmergencyFund", RCEmergencyFund, ScoringDimensions{0.9, 0.9, 0.3, 0.9, 0.8, 1.0}},
		{"CashFlow", RCCashFlow, ScoringDimensions{0.6, 0.4, 0.2, 0.8, 0.7, 0.6}},
		{"Savings", RCSavings, ScoringDimensions{0.7, 0.5, 0.4, 0.7, 0.7, 0.7}},
		{"Insurance", RCInsurance, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
		{"ExpenseReduction", RCExpenseReduction, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
		{"IncomeGrowth", RCIncomeGrowth, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
		{"Tax", RCTax, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
		{"Retirement", RCRetirement, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
		{"Behaviour", RCBehaviour, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
		{"FinancialDiscipline", RCFinancialDiscipline, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
		{"Investment", RCInvestment, ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := e.scoreRecommendation(Recommendation{Category: tc.category}, inputs)
			if got != tc.want {
				t.Errorf("scoreRecommendation(%s) = %+v, want %+v", tc.category, got, tc.want)
			}
		})
	}
}

func TestCalculateComposite_KnownValues(t *testing.T) {
	e := &Engine{}

	tests := []struct {
		name   string
		scores ScoringDimensions
		want   float64
	}{
		{"EmergencyFund high", ScoringDimensions{0.9, 0.9, 0.3, 0.9, 0.8, 1.0}, 0.81},
		{"GoalFunding", ScoringDimensions{0.8, 0.7, 0.4, 0.8, 1.0, 0.7}, 0.73},
		{"Debt", ScoringDimensions{0.7, 0.8, 0.5, 0.8, 0.6, 0.8}, 0.71},
		{"Savings", ScoringDimensions{0.7, 0.5, 0.4, 0.7, 0.7, 0.7}, 0.61},
		{"PortfolioRebalancing", ScoringDimensions{0.6, 0.5, 0.3, 0.7, 0.5, 0.6}, 0.54},
		{"CashFlow", ScoringDimensions{0.6, 0.4, 0.2, 0.8, 0.7, 0.6}, 0.54},
		{"Default", ScoringDimensions{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}, 0.50},
		{"Max values", ScoringDimensions{1.0, 1.0, 1.0, 1.0, 1.0, 1.0}, 1.00},
		{"All zeros", ScoringDimensions{0, 0, 0, 0, 0, 0}, 0.00},
		{"Only Impact", ScoringDimensions{1.0, 0, 0, 0, 0, 0}, 0.30},
		{"Only Urgency", ScoringDimensions{0, 1.0, 0, 0, 0, 0}, 0.20},
		{"Only Difficulty", ScoringDimensions{0, 0, 1.0, 0, 0, 0}, 0.15},
		{"Only Confidence", ScoringDimensions{0, 0, 0, 1.0, 0, 0}, 0.15},
		{"Only GoalAlignment", ScoringDimensions{0, 0, 0, 0, 1.0, 0}, 0.10},
		{"Only FinancialHealth", ScoringDimensions{0, 0, 0, 0, 0, 1.0}, 0.10},
		{"Rounding check", ScoringDimensions{0.333, 0.333, 0.333, 0.333, 0.333, 0.333}, 0.33},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := e.calculateComposite(tc.scores)
			if got != tc.want {
				t.Errorf("calculateComposite(%+v) = %.2f, want %.2f", tc.scores, got, tc.want)
			}
		})
	}
}

func TestCalculateComposite_ExtremeValues(t *testing.T) {
	e := &Engine{}
	got := e.calculateComposite(ScoringDimensions{1e308, 0, 0, 0, 0, 0})
	if !math.IsInf(got, 1) {
		t.Logf("extreme large value produces +Inf as expected: %v", got)
	}
	got = e.calculateComposite(ScoringDimensions{-1.0, -1.0, -1.0, -1.0, -1.0, -1.0})
	if got != -1.00 {
		t.Errorf("expected -1.00 for all -1 values, got %.2f", got)
	}
}

func TestRankAndPrioritise_SortsDescending(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := []Recommendation{
		{Score: 0.3, Category: RCSavings},
		{Score: 0.9, Category: RCEmergencyFund},
		{Score: 0.5, Category: RCDebt},
	}
	result := e.rankAndPrioritise(recs)
	if len(result) != 3 {
		t.Fatalf("expected 3 recs, got %d", len(result))
	}
	if result[0].Score != 0.9 || result[1].Score != 0.5 || result[2].Score != 0.3 {
		t.Errorf("recs not sorted descending: %+v", result)
	}
}

func TestRankAndPrioritise_PriorityAssignment(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := []Recommendation{
		{Score: 0.3},
		{Score: 0.9},
		{Score: 0.5},
	}
	result := e.rankAndPrioritise(recs)
	for i, r := range result {
		if r.Priority != i+1 {
			t.Errorf("expected priority %d, got %d for rec at index %d", i+1, r.Priority, i)
		}
	}
}

func TestRankAndPrioritise_FiltersBelowThreshold(t *testing.T) {
	cfg := defaultConfig()
	cfg.MinRecommendationScore = 0.5
	e := &Engine{config: cfg}
	recs := []Recommendation{
		{Score: 0.9},
		{Score: 0.3},
		{Score: 0.7},
	}
	result := e.rankAndPrioritise(recs)
	if len(result) != 2 {
		t.Fatalf("expected 2 recs (score >= 0.5), got %d", len(result))
	}
	if result[0].Score != 0.9 || result[1].Score != 0.7 {
		t.Errorf("unexpected ordering: %+v", result)
	}
}

func TestRankAndPrioritise_AllBelowThreshold(t *testing.T) {
	e := &Engine{config: Config{MinRecommendationScore: 0.8}}
	recs := []Recommendation{
		{Score: 0.3},
		{Score: 0.5},
		{Score: 0.7},
	}
	result := e.rankAndPrioritise(recs)
	if len(result) != 0 {
		t.Errorf("expected 0 recs, all below threshold, got %d", len(result))
	}
}

func TestRankAndPrioritise_EmptyInput(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	result := e.rankAndPrioritise(nil)
	if len(result) != 0 {
		t.Errorf("expected 0 recs for nil input, got %d", len(result))
	}
	result = e.rankAndPrioritise([]Recommendation{})
	if len(result) != 0 {
		t.Errorf("expected 0 recs for empty input, got %d", len(result))
	}
}

func TestRankAndPrioritise_EqualScores(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs := []Recommendation{
		{Score: 0.5, Category: RCSavings},
		{Score: 0.5, Category: RCDebt},
	}
	result := e.rankAndPrioritise(recs)
	if len(result) != 2 {
		t.Fatalf("expected 2 recs, got %d", len(result))
	}
	if result[0].Priority != 1 || result[1].Priority != 2 {
		t.Errorf("priorities not sequential for equal scores: %+v", result)
	}
}

func TestFilterByThreshold(t *testing.T) {
	tests := []struct {
		name   string
		recs   []Recommendation
		min    float64
		wantN  int
	}{
		{"all above", []Recommendation{{Score: 0.8}, {Score: 0.6}}, 0.5, 2},
		{"some below", []Recommendation{{Score: 0.8}, {Score: 0.2}}, 0.5, 1},
		{"all below", []Recommendation{{Score: 0.1}, {Score: 0.2}}, 0.5, 0},
		{"edge at threshold", []Recommendation{{Score: 0.5}}, 0.5, 1},
		{"nil input", nil, 0.5, 0},
		{"zero min", []Recommendation{{Score: 0}, {Score: 0.1}}, 0, 2},
		{"negative min", []Recommendation{{Score: -0.1}, {Score: 0.1}}, 0, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filterByThreshold(tc.recs, tc.min)
			if len(got) != tc.wantN {
				t.Errorf("filterByThreshold returned %d recs, want %d", len(got), tc.wantN)
			}
		})
	}
}

func TestGenerateEvidence(t *testing.T) {
	e := &Engine{}
	rec := Recommendation{Evidence: []string{"item1", "item2"}}
	got := e.generateEvidence(rec, Inputs{})
	if len(got) != 2 {
		t.Errorf("expected 2 evidence items, got %d", len(got))
	}
	if got[0] != "item1" || got[1] != "item2" {
		t.Errorf("evidence mismatch: %+v", got)
	}
}

func TestGenerateExplanation(t *testing.T) {
	e := &Engine{}
	rec := Recommendation{
		Summary:             "Test summary",
		ExpectedImprovement: "Test improvement",
	}
	got := e.generateExplanation(rec, Inputs{})
	want := "Test summary. Test improvement."
	if got != want {
		t.Errorf("generateExplanation = '%s', want '%s'", got, want)
	}
}

func TestExecute_FullPipeline(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		GoalGaps:        map[string]int64{"Retirement": 500000},
		GoalImportance:  map[string]string{"Retirement": "Mandatory"},
		Debts:           []DebtInfo{{Name: "Card", Balance: 100000, InterestRate: 24.0}},
		PortfolioDrift:  10.0,
		EmergencyMonths: 2.0,
		CashFlowSurplus: 50000,
		SavingsRate:     10,
	}

	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(recs) == 0 {
		t.Fatal("expected at least 1 recommendation")
	}

	for i, r := range recs {
		if r.Status != RSGenerated {
			t.Errorf("rec[%d] status = %s, want Generated", i, r.Status)
		}
		if r.Score <= 0 {
			t.Errorf("rec[%d] score = %.2f, want > 0", i, r.Score)
		}
		if r.Scores == (ScoringDimensions{}) {
			t.Errorf("rec[%d] has zero ScoringDimensions", i)
		}
		if len(r.Evidence) == 0 {
			t.Errorf("rec[%d] has no evidence", i)
		}
		if r.Explanation == "" {
			t.Errorf("rec[%d] has no explanation", i)
		}
		if r.Priority != i+1 {
			t.Errorf("rec[%d] priority = %d, want %d", i, r.Priority, i+1)
		}
	}

	if recs[0].Score < recs[len(recs)-1].Score {
		t.Errorf("recs not sorted descending by score")
	}
}

func TestExecute_MaxTotalTruncation(t *testing.T) {
	cfg := defaultConfig()
	cfg.MaxTotalRecommendations = 2
	e := &Engine{config: cfg}

	inputs := Inputs{
		GoalGaps: map[string]int64{
			"G1": 200000, "G2": 200000, "G3": 200000,
		},
		GoalImportance: map[string]string{
			"G1": "Mandatory", "G2": "Mandatory", "G3": "Mandatory",
		},
	}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if len(recs) > 2 {
		t.Errorf("expected at most 2 recs (MaxTotalRecommendations=2), got %d", len(recs))
	}
}

func TestExecute_EmptyInputs(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	recs, err := e.Execute(Inputs{})
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if len(recs) < 1 {
		t.Errorf("expected at least 1 rec (emergency fund/savings from zero defaults), got %d", len(recs))
	}
}

func TestExecute_OnlyGoalGaps(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		GoalGaps:        map[string]int64{"College": 300000},
		GoalImportance:  map[string]string{"College": "Essential"},
		EmergencyMonths: 6,
		SavingsRate:     20,
	}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 rec, got %d", len(recs))
	}
	if recs[0].Category != RCGoalFunding {
		t.Errorf("expected RCGoalFunding, got %s", recs[0].Category)
	}
}

func TestExecute_AllSurplusNoGaps(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{CashFlowSurplus: 100000}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	found := false
	for _, r := range recs {
		if r.Category == RCCashFlow && r.Title == "Allocate surplus" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Allocate surplus' recommendation")
	}
}

func TestExecute_NegativeSurplus(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{CashFlowSurplus: -50000, SavingsRate: 30, EmergencyMonths: 12, PortfolioDrift: 2}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	for _, r := range recs {
		if r.Category == RCCashFlow && r.Title == "Allocate surplus" {
			t.Errorf("unexpected surplus recommendation for negative surplus")
		}
	}
}

func TestExecute_OnlyHighInterestDebt(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		Debts:           []DebtInfo{{Name: "Loan", Balance: 1000000, InterestRate: 15.0}},
		EmergencyMonths: 6,
		SavingsRate:     20,
	}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 rec, got %d", len(recs))
	}
	if recs[0].Category != RCDebt {
		t.Errorf("expected RCDebt, got %s", recs[0].Category)
	}
}

func TestExecute_AllRecommendationsBelowMinScore(t *testing.T) {
	cfg := defaultConfig()
	cfg.MinRecommendationScore = 0.99
	e := &Engine{config: cfg}
	inputs := Inputs{
		CashFlowSurplus: 100000,
		SavingsRate:     5,
	}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("expected 0 recs (all below min score 0.99), got %d", len(recs))
	}
}

func TestExecute_CategoriesPresent(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		GoalGaps:           map[string]int64{"Ret": 500000},
		GoalImportance:     map[string]string{"Ret": "High"},
		Debts:              []DebtInfo{{Name: "Card", Balance: 100000, InterestRate: 24.0}},
		PortfolioDrift:     10.0,
		EmergencyMonths:    1.0,
		OverspentCategories: map[string]int64{"Dining": 20000},
		CashFlowSurplus:    60000,
		SavingsRate:        5,
	}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	present := map[RecCategory]bool{}
	for _, r := range recs {
		present[r.Category] = true
	}

	expected := []RecCategory{RCGoalFunding, RCDebt, RCPortfolioRebalancing, RCEmergencyFund, RCCashFlow, RCSavings}
	for _, cat := range expected {
		if !present[cat] {
			t.Errorf("missing category %s in pipeline output", cat)
		}
	}
}

func TestNewEngine(t *testing.T) {
	cfg := defaultConfig()
	e := NewEngine(cfg)
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
	if e.config != cfg {
		t.Errorf("config not set correctly")
	}
}

func TestFmtMoney(t *testing.T) {
	tests := []struct {
		val  int64
		want string
	}{
		{0, "₹0"},
		{100, "₹100"},
		{-500, "₹-500"},
		{1000000, "₹1000000"},
	}
	for _, tc := range tests {
		got := fmtMoney(tc.val)
		if got != tc.want {
			t.Errorf("fmtMoney(%d) = '%s', want '%s'", tc.val, got, tc.want)
		}
	}
}

func TestFmtFlt(t *testing.T) {
	tests := []struct {
		val  float64
		want string
	}{
		{0.0, "0.0"},
		{5.5, "5.5"},
		{12.345, "12.3"},
		{-3.14, "-3.1"},
	}
	for _, tc := range tests {
		got := fmtFlt(tc.val)
		if got != tc.want {
			t.Errorf("fmtFlt(%.3f) = '%s', want '%s'", tc.val, got, tc.want)
		}
	}
}

func TestExecute_ResultStability(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		GoalGaps:        map[string]int64{"Ret": 500000},
		GoalImportance:  map[string]string{"Ret": "High"},
		Debts:           []DebtInfo{{Name: "Card", Balance: 100000, InterestRate: 24.0}},
		PortfolioDrift:  10.0,
		EmergencyMonths: 1.0,
		CashFlowSurplus: 60000,
		SavingsRate:     5,
	}

	recs1, _ := e.Execute(inputs)
	recs2, _ := e.Execute(inputs)
	if len(recs1) != len(recs2) {
		t.Fatalf("result length differs between runs: %d vs %d", len(recs1), len(recs2))
	}
	for i := range recs1 {
		if recs1[i].Score != recs2[i].Score {
			t.Errorf("rec[%d] score differs between runs: %.2f vs %.2f", i, recs1[i].Score, recs2[i].Score)
		}
		if recs1[i].Category != recs2[i].Category {
			t.Errorf("rec[%d] category differs between runs: %s vs %s", i, recs1[i].Category, recs2[i].Category)
		}
	}
}

func TestExecute_DebtBelowThresholdNoRec(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		Debts: []DebtInfo{
			{Name: "CheapLoan", Balance: 1000000, InterestRate: 7.0},
		},
	}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	for _, r := range recs {
		if r.Category == RCDebt {
			t.Errorf("unexpected debt recommendation for 7%% interest (below 12%% threshold)")
		}
	}
}

func TestExecute_DriftBelowThresholdNoRec(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{PortfolioDrift: 3.0}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	for _, r := range recs {
		if r.Category == RCPortfolioRebalancing {
			t.Errorf("unexpected rebalance recommendation for 3%% drift (below 5%% threshold)")
		}
	}
}

func TestExecute_EmergencyAtTargetNoRec(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{EmergencyMonths: 6.0}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	for _, r := range recs {
		if r.Category == RCEmergencyFund {
			t.Errorf("unexpected emergency fund rec at exactly 6 months (at target)")
		}
	}
}

func TestExecute_MultipleIdenticalOverspentCategories(t *testing.T) {
	e := &Engine{config: defaultConfig()}
	inputs := Inputs{
		OverspentCategories: map[string]int64{"Food": 5000, "Transport": 3000, "Entertainment": 7000},
	}
	recs, err := e.Execute(inputs)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	cashFlowCount := 0
	for _, r := range recs {
		if r.Category == RCCashFlow {
			cashFlowCount++
		}
	}
	if cashFlowCount < 3 {
		t.Errorf("expected 3 CashFlow overspend recs, got %d", cashFlowCount)
	}
}
