package engine

import (
	"context"
	"testing"

	"github.com/horizon/core/services/ai/provider"
)

// mockAIProvider wraps provider.MockProvider to satisfy the full AIProvider interface
// (MockProvider is missing StreamChat).
type mockAIProvider struct {
	*provider.MockProvider
}

func newMockAI() *mockAIProvider {
	return &mockAIProvider{MockProvider: provider.NewMock()}
}

func (m *mockAIProvider) StreamChat(ctx context.Context, req provider.ChatRequest) (<-chan string, error) {
	ch := make(chan string)
	close(ch)
	return ch, nil
}

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

func TestNewComposer(t *testing.T) {
	t.Run("without AI provider", func(t *testing.T) {
		c := NewComposer(nil)
		if c == nil {
			t.Fatal("NewComposer(nil) returned nil")
		}
		if c.ai != nil {
			t.Error("expected nil ai provider")
		}
	})

	t.Run("with AI provider", func(t *testing.T) {
		mock := newMockAI()
		c := NewComposer(mock)
		if c == nil {
			t.Fatal("NewComposer(mock) returned nil")
		}
		if c.ai == nil {
			t.Error("expected non-nil ai provider")
		}
	})
}

func TestGenerate_RuleBased_Health(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("health improving", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HealthScore: 65, HealthChange: 5, HealthGrade: "B"})
		got := findByID(ins, "ins-health-up")
		if got == nil {
			t.Fatal("expected ins-health-up")
		}
		if got.Category != ICHealth {
			t.Errorf("category = %q, want %q", got.Category, ICHealth)
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q, want %q", got.Priority, PriorityLow)
		}
		if got.Confidence != "High" {
			t.Errorf("confidence = %q", got.Confidence)
		}
		if got.SourceEngine != "HealthScore" {
			t.Errorf("source = %q", got.SourceEngine)
		}
		if got.Severity != "positive" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("health declining", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HealthScore: 55, HealthChange: -10, HealthGrade: "C"})
		got := findByID(ins, "ins-health-down")
		if got == nil {
			t.Fatal("expected ins-health-down")
		}
		if got.Category != ICHealth {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityHigh {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "warning" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("health critical", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HealthScore: 35, HealthGrade: "F"})
		got := findByID(ins, "ins-health-critical")
		if got == nil {
			t.Fatal("expected ins-health-critical")
		}
		if got.Priority != PriorityCritical {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "critical" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("health no insights when score is zero", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HealthScore: 0})
		if findByID(ins, "ins-health-up") != nil {
			t.Error("expected no health insight for score 0")
		}
		if findByID(ins, "ins-health-down") != nil {
			t.Error("expected no health decline insight for score 0")
		}
		if findByID(ins, "ins-health-critical") != nil {
			t.Error("expected no health critical insight for score 0")
		}
	})

	t.Run("health critical and declining both fire", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HealthScore: 30, HealthChange: -5, HealthGrade: "F"})
		if findByID(ins, "ins-health-down") == nil {
			t.Error("expected ins-health-down")
		}
		if findByID(ins, "ins-health-critical") == nil {
			t.Error("expected ins-health-critical")
		}
	})
}

func TestGenerate_RuleBased_Risk(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("risk critical >= 80", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 85, RiskLevel: "Very High"})
		got := findByID(ins, "ins-risk-critical")
		if got == nil {
			t.Fatal("expected ins-risk-critical")
		}
		if got.Category != ICRisk {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityCritical {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "critical" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("risk elevated >= 60", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 72, RiskLevel: "High"})
		got := findByID(ins, "ins-risk-high")
		if got == nil {
			t.Fatal("expected ins-risk-high")
		}
		if got.Priority != PriorityHigh {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "warning" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("risk low < 30", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 25, RiskLevel: "Low"})
		got := findByID(ins, "ins-risk-low")
		if got == nil {
			t.Fatal("expected ins-risk-low")
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "positive" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("risk no insight for mid-range", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 45, RiskLevel: "Moderate"})
		if findByID(ins, "ins-risk-critical") != nil {
			t.Error("unexpected risk critical")
		}
		if findByID(ins, "ins-risk-high") != nil {
			t.Error("unexpected risk high")
		}
		if findByID(ins, "ins-risk-low") != nil {
			t.Error("unexpected risk low")
		}
	})

	t.Run("risk boundary 60 exactly", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 60, RiskLevel: "High"})
		if findByID(ins, "ins-risk-high") == nil {
			t.Error("expected ins-risk-high at boundary 60")
		}
	})

	t.Run("risk boundary 80 exactly", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 80, RiskLevel: "Very High"})
		if findByID(ins, "ins-risk-critical") == nil {
			t.Error("expected ins-risk-critical at boundary 80")
		}
	})

	t.Run("risk boundary 29 (below low threshold)", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 29, RiskLevel: "Low"})
		if findByID(ins, "ins-risk-low") == nil {
			t.Error("expected ins-risk-low at 29")
		}
	})
}

func TestGenerate_RuleBased_Savings(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("savings improving trend", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{SavingsRate: 25.5, SavingsTrend: "improving"})
		got := findByID(ins, "ins-savings")
		if got == nil {
			t.Fatal("expected ins-savings")
		}
		if got.Category != ICSavings {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityMedium {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "improving" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("savings declining trend", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{SavingsRate: 10.0, SavingsTrend: "declining"})
		got := findByID(ins, "ins-savings")
		if got == nil {
			t.Fatal("expected ins-savings")
		}
		if got.Severity != "declining" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("savings stable trend", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{SavingsRate: 15.0, SavingsTrend: "stable"})
		got := findByID(ins, "ins-savings")
		if got == nil {
			t.Fatal("expected ins-savings")
		}
		if got.Severity != "stable" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("savings no insight when rate is zero", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{SavingsRate: 0})
		if findByID(ins, "ins-savings") != nil {
			t.Error("expected no savings insight for rate 0")
		}
	})
}

func TestGenerate_RuleBased_CashFlow(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("cash flow surplus", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{CashFlowSurplus: 50000})
		got := findByID(ins, "ins-cf-surplus")
		if got == nil {
			t.Fatal("expected ins-cf-surplus")
		}
		if got.Category != ICCashFlow {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "positive" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("cash flow deficit", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{CashFlowSurplus: -25000})
		got := findByID(ins, "ins-cf-deficit")
		if got == nil {
			t.Fatal("expected ins-cf-deficit")
		}
		if got.Priority != PriorityHigh {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "warning" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("cash flow zero", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{CashFlowSurplus: 0})
		if findByID(ins, "ins-cf-surplus") != nil {
			t.Error("unexpected surplus insight")
		}
		if findByID(ins, "ins-cf-deficit") != nil {
			t.Error("unexpected deficit insight")
		}
	})
}

func TestGenerate_RuleBased_Goals(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("all goals on track - ratio >= 1.0", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{TotalGoals: 5, GoalsOnTrack: 5})
		got := findByID(ins, "ins-goals-all")
		if got == nil {
			t.Fatal("expected ins-goals-all")
		}
		if got.Category != ICGoal {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "achievement" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("goals at risk - ratio < 0.5", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{TotalGoals: 6, GoalsOnTrack: 2})
		got := findByID(ins, "ins-goals-risk")
		if got == nil {
			t.Fatal("expected ins-goals-risk")
		}
		if got.Priority != PriorityHigh {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "warning" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("goals mid-range no insight", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{TotalGoals: 4, GoalsOnTrack: 3})
		if findByID(ins, "ins-goals-all") != nil {
			t.Error("unexpected goals-all")
		}
		if findByID(ins, "ins-goals-risk") != nil {
			t.Error("unexpected goals-risk")
		}
	})

	t.Run("goals no insight when total is zero", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{TotalGoals: 0, GoalsOnTrack: 0})
		if findByID(ins, "ins-goals-all") != nil {
			t.Error("unexpected goals-all when no goals")
		}
		if findByID(ins, "ins-goals-risk") != nil {
			t.Error("unexpected goals-risk when no goals")
		}
	})

	t.Run("goals boundary ratio 1.0", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{TotalGoals: 3, GoalsOnTrack: 3})
		if findByID(ins, "ins-goals-all") == nil {
			t.Error("expected ins-goals-all at ratio 1.0")
		}
	})

	t.Run("goals boundary ratio just below 0.5", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{TotalGoals: 5, GoalsOnTrack: 2})
		if findByID(ins, "ins-goals-risk") == nil {
			t.Error("expected ins-goals-risk at ratio 0.4")
		}
	})
}

func TestGenerate_RuleBased_NetWorth(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("net worth growing", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{NetWorth: 10000000, NetWorthChange: 500000})
		got := findByID(ins, "ins-nw-up")
		if got == nil {
			t.Fatal("expected ins-nw-up")
		}
		if got.Category != ICNetWorth {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "positive" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("net worth declined", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{NetWorth: 5000000, NetWorthChange: -200000})
		got := findByID(ins, "ins-nw-down")
		if got == nil {
			t.Fatal("expected ins-nw-down")
		}
		if got.Priority != PriorityMedium {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "warning" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("net worth unchanged", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{NetWorth: 1000000, NetWorthChange: 0})
		if findByID(ins, "ins-nw-up") != nil {
			t.Error("unexpected nw-up")
		}
		if findByID(ins, "ins-nw-down") != nil {
			t.Error("unexpected nw-down")
		}
	})
}

func TestGenerate_RuleBased_Portfolio(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("portfolio positive return", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{PortfolioValue: 3000000, PortfolioReturn: 12.5})
		got := findByID(ins, "ins-pf-up")
		if got == nil {
			t.Fatal("expected ins-pf-up")
		}
		if got.Category != ICPortfolio {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Confidence != "Medium" {
			t.Errorf("confidence = %q", got.Confidence)
		}
		if got.Severity != "positive" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("portfolio no insight when return is zero or negative", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{PortfolioValue: 1000000, PortfolioReturn: 0})
		if findByID(ins, "ins-pf-up") != nil {
			t.Error("unexpected pf-up for zero return")
		}
	})

	t.Run("portfolio no insight when return is negative", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{PortfolioReturn: -5.0})
		if findByID(ins, "ins-pf-up") != nil {
			t.Error("unexpected pf-up for negative return")
		}
	})
}

func TestGenerate_RuleBased_Achievements(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("achievements present", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{AchievementCount: 3})
		got := findByID(ins, "ins-ach")
		if got == nil {
			t.Fatal("expected ins-ach")
		}
		if got.Category != ICAchievement {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "achievement" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("achievements zero", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{AchievementCount: 0})
		if findByID(ins, "ins-ach") != nil {
			t.Error("unexpected achievement insight")
		}
	})
}

func TestGenerate_RuleBased_Behaviour(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("spending anomaly present", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{SpendingAnomaly: "Unusual spending in dining category"})
		got := findByID(ins, "ins-anomaly")
		if got == nil {
			t.Fatal("expected ins-anomaly")
		}
		if got.Category != ICBehaviour {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityMedium {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "info" {
			t.Errorf("severity = %q", got.Severity)
		}
		if got.Summary != "Unusual spending in dining category" {
			t.Errorf("summary = %q", got.Summary)
		}
	})

	t.Run("spending anomaly empty", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{SpendingAnomaly: ""})
		if findByID(ins, "ins-anomaly") != nil {
			t.Error("unexpected anomaly insight")
		}
	})
}

func TestGenerate_RuleBased_Recommendations(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("recommendations available", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HasRecs: true, RecCount: 5})
		got := findByID(ins, "ins-recs")
		if got == nil {
			t.Fatal("expected ins-recs")
		}
		if got.Category != ICRecommend {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityMedium {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.SourceEngine != "Recommendation" {
			t.Errorf("source = %q", got.SourceEngine)
		}
	})

	t.Run("recommendations not available", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HasRecs: false})
		if findByID(ins, "ins-recs") != nil {
			t.Error("unexpected recs insight")
		}
	})
}

func TestGenerate_RuleBased_Optimizations(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("optimizations available", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HasOpts: true})
		got := findByID(ins, "ins-opts")
		if got == nil {
			t.Fatal("expected ins-opts")
		}
		if got.Category != ICOptimization {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityMedium {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.SourceEngine != "Optimization" {
			t.Errorf("source = %q", got.SourceEngine)
		}
	})

	t.Run("optimizations not available", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{HasOpts: false})
		if findByID(ins, "ins-opts") != nil {
			t.Error("unexpected opts insight")
		}
	})
}

func TestGenerate_RuleBased_Milestones(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("milestones present", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{MilestoneCount: 2})
		got := findByID(ins, "ins-ms")
		if got == nil {
			t.Fatal("expected ins-ms")
		}
		if got.Category != ICMilestone {
			t.Errorf("category = %q", got.Category)
		}
		if got.Priority != PriorityLow {
			t.Errorf("priority = %q", got.Priority)
		}
		if got.Severity != "achievement" {
			t.Errorf("severity = %q", got.Severity)
		}
	})

	t.Run("milestones zero", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{MilestoneCount: 0})
		if findByID(ins, "ins-ms") != nil {
			t.Error("unexpected milestone insight")
		}
	})
}

func TestGenerate_AIInsights(t *testing.T) {
	mock := newMockAI()
	c := NewComposer(mock)
	ctx := context.Background()

	ins := c.generate(ctx, Inputs{
		HealthScore:      75,
		SavingsRate:      20,
		NetWorth:         5000000,
		SpendingAnomaly:  "Increased dining out",
	})

	ai0 := findByID(ins, "ins-ai-0")
	if ai0 == nil {
		t.Fatal("expected ins-ai-0 from mock provider")
	}
	if ai0.SourceEngine != "AI" {
		t.Errorf("source = %q", ai0.SourceEngine)
	}
	if ai0.Category != ICBehaviour {
		t.Errorf("ai-0 category = %q, want Behaviour", ai0.Category)
	}
	if ai0.Priority != PriorityMedium {
		t.Errorf("ai-0 priority = %q", ai0.Priority)
	}
	if ai0.Title != "High Discretionary Spending" {
		t.Errorf("ai-0 title = %q", ai0.Title)
	}
	if ai0.Confidence != "High" {
		t.Errorf("ai-0 confidence = %q", ai0.Confidence)
	}

	ai1 := findByID(ins, "ins-ai-1")
	if ai1 == nil {
		t.Fatal("expected ins-ai-1 from mock provider")
	}
	if ai1.Category != ICSavings {
		t.Errorf("ai-1 category = %q, want Savings", ai1.Category)
	}
	if ai1.Title != "Potential to Increase Savings" {
		t.Errorf("ai-1 title = %q", ai1.Title)
	}
	if ai1.Confidence != "Medium" {
		t.Errorf("ai-1 confidence = %q", ai1.Confidence)
	}
}

func TestGenerate_AIInsights_NilProvider(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	ins := c.generate(ctx, Inputs{
		HealthScore: 75,
		SavingsRate: 20,
	})

	for _, in := range ins {
		if in.SourceEngine == "AI" {
			t.Errorf("unexpected AI insight with nil provider: %s", in.InsightID)
		}
	}
}

func TestGenerate_CombinedScenario(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	ins := c.generate(ctx, Inputs{
		HealthScore:      35,
		HealthChange:     -3,
		HealthGrade:      "F",
		RiskScore:        88,
		RiskLevel:        "Critical",
		SavingsRate:      12.0,
		SavingsTrend:     "declining",
		CashFlowSurplus:  -15000,
		TotalGoals:       4,
		GoalsOnTrack:     1,
		NetWorth:         2000000,
		NetWorthChange:   -100000,
		PortfolioValue:   1500000,
		PortfolioReturn:  8.2,
		AchievementCount: 5,
		SpendingAnomaly:  "Unusual transaction pattern",
		HasRecs:          true,
		RecCount:         3,
		HasOpts:          true,
		MilestoneCount:   2,
	})

	expectedIDs := []string{
		"ins-health-down", "ins-health-critical",
		"ins-risk-critical",
		"ins-savings",
		"ins-cf-deficit",
		"ins-goals-risk",
		"ins-nw-down",
		"ins-pf-up",
		"ins-ach",
		"ins-anomaly",
		"ins-recs",
		"ins-opts",
		"ins-ms",
	}

	for _, id := range expectedIDs {
		if findByID(ins, id) == nil {
			t.Errorf("expected insight %s in combined scenario", id)
		}
	}
}

func TestGenerate_MinimalInputs(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	t.Run("risk score 0 triggers low risk insight", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{})
		if got := findByID(ins, "ins-risk-low"); got == nil {
			t.Fatal("expected ins-risk-low for RiskScore=0")
		}
	})

	t.Run("mid-range risk with no other signals", func(t *testing.T) {
		ins := c.generate(ctx, Inputs{RiskScore: 45, RiskLevel: "Moderate"})
		if len(ins) != 0 {
			t.Errorf("expected 0 insights for mid-range risk and no other signals, got %d", len(ins))
		}
	})
}

func TestBuildDashboard(t *testing.T) {
	mock := newMockAI()
	c := NewComposer(mock)
	ctx := context.Background()

	inputs := Inputs{
		HealthScore:      30,
		HealthChange:     -5,
		HealthGrade:      "F",
		RiskScore:        85,
		RiskLevel:        "Critical",
		SavingsRate:      15.0,
		SavingsTrend:     "declining",
		CashFlowSurplus:  -20000,
		TotalGoals:       5,
		GoalsOnTrack:     2,
		NetWorthChange:   100000,
		AchievementCount: 2,
		SpendingAnomaly:  "Alert",
		HasRecs:          true,
		RecCount:         3,
		HasOpts:          true,
		MilestoneCount:   1,
	}

	dash := c.BuildDashboard(ctx, inputs)

	if dash.Total == 0 {
		t.Fatal("expected non-zero insights")
	}

	expectedCategories := []InsightCategory{
		ICHealth, ICRisk, ICSavings, ICCashFlow,
		ICGoal, ICNetWorth, ICAchievement, ICBehaviour,
		ICRecommend, ICOptimization, ICMilestone,
	}
	for _, cat := range expectedCategories {
		if _, ok := dash.ByCategory[cat]; !ok {
			t.Errorf("missing category %q in ByCategory", cat)
		}
	}

	if _, ok := dash.ByPriority[PriorityCritical]; !ok {
		t.Error("missing PriorityCritical in ByPriority")
	}
	if _, ok := dash.ByPriority[PriorityHigh]; !ok {
		t.Error("missing PriorityHigh in ByPriority")
	}

	if len(dash.Critical) == 0 {
		t.Error("expected critical insights")
	}
	if len(dash.High) == 0 {
		t.Error("expected high insights")
	}

	allCount := dash.ByPriority[PriorityCritical] + dash.ByPriority[PriorityHigh] +
		dash.ByPriority[PriorityMedium] + dash.ByPriority[PriorityLow] +
		dash.ByPriority[PriorityInfo]
	if dash.Total != allCount {
		t.Errorf("total %d != sum of byPriority %d", dash.Total, allCount)
	}

	var catSum int
	for _, v := range dash.ByCategory {
		catSum += v
	}
	if dash.Total != catSum {
		t.Errorf("total %d != sum of byCategory %d", dash.Total, catSum)
	}
}

func TestBuildDashboard_AIProvider(t *testing.T) {
	mock := newMockAI()
	c := NewComposer(mock)
	ctx := context.Background()

	dash := c.BuildDashboard(ctx, Inputs{
		HealthScore: 80,
		SavingsRate: 25,
		NetWorth:    10000000,
	})

	if dash.Total == 0 {
		t.Fatal("expected insights")
	}
	if _, ok := dash.ByCategory[ICBehaviour]; !ok {
		t.Error("expected AI Behaviour insight in dashboard")
	}
	if _, ok := dash.ByCategory[ICSavings]; !ok {
		t.Error("expected AI Savings insight in dashboard")
	}
}

func TestBuildFeed(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	inputs := Inputs{
		HealthScore:   65,
		HealthChange:  5,
		HealthGrade:   "B",
		SavingsRate:   20,
		AchievementCount: 3,
		MilestoneCount:   2,
	}

	t.Run("default max", func(t *testing.T) {
		feed := c.BuildFeed(ctx, inputs, 0, "")
		if feed.Total != len(feed.Insights) {
			t.Errorf("total %d != len %d", feed.Total, len(feed.Insights))
		}
		if feed.Cursor != "" {
			t.Errorf("expected empty cursor, got %q", feed.Cursor)
		}
	})

	t.Run("max less than total", func(t *testing.T) {
		feed := c.BuildFeed(ctx, inputs, 1, "")
		if len(feed.Insights) > 1 {
			t.Errorf("expected at most 1 insight, got %d", len(feed.Insights))
		}
		if !feed.HasMore {
			t.Error("expected HasMore true when limited")
		}
		if feed.Cursor == "" {
			t.Error("expected non-empty cursor when hasMore")
		}
	})

	t.Run("cursor points to last item", func(t *testing.T) {
		feed := c.BuildFeed(ctx, inputs, 1, "")
		if feed.Cursor != "" {
			if feed.Cursor != feed.Insights[len(feed.Insights)-1].InsightID {
				t.Errorf("cursor %q should match last insight ID", feed.Cursor)
			}
		}
	})
}

func TestBuildFeed_WithAI(t *testing.T) {
	mock := newMockAI()
	c := NewComposer(mock)
	ctx := context.Background()

	feed := c.BuildFeed(ctx, Inputs{HealthScore: 70}, 10, "")
	if feed.Total == 0 {
		t.Fatal("expected insights in feed")
	}

	var hasAI bool
	for _, ins := range feed.Insights {
		if ins.SourceEngine == "AI" {
			hasAI = true
			break
		}
	}
	if !hasAI {
		t.Error("expected AI-generated insights in feed")
	}
}

func TestFilterByCategory(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	inputs := Inputs{
		HealthScore:      50,
		HealthChange:     10,
		HealthGrade:      "C",
		RiskScore:        70,
		RiskLevel:        "High",
		SavingsRate:      20,
		CashFlowSurplus:  50000,
		AchievementCount: 2,
	}

	t.Run("filter by Health", func(t *testing.T) {
		filtered := c.FilterByCategory(ctx, inputs, ICHealth)
		if len(filtered) == 0 {
			t.Fatal("expected health insights")
		}
		for _, ins := range filtered {
			if ins.Category != ICHealth {
				t.Errorf("expected all health, got %q", ins.Category)
			}
		}
	})

	t.Run("filter by none", func(t *testing.T) {
		filtered := c.FilterByCategory(ctx, inputs, ICDebt)
		if len(filtered) != 0 {
			t.Errorf("expected 0 debt insights, got %d", len(filtered))
		}
	})

	t.Run("filter by Risk", func(t *testing.T) {
		filtered := c.FilterByCategory(ctx, inputs, ICRisk)
		if len(filtered) == 0 {
			t.Fatal("expected risk insights")
		}
		for _, ins := range filtered {
			if ins.Category != ICRisk {
				t.Errorf("expected all risk, got %q", ins.Category)
			}
		}
	})
}

func TestFilterByPriority(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	inputs := Inputs{
		HealthScore:    35,
		HealthGrade:    "F",
		RiskScore:      85,
		RiskLevel:      "Critical",
		SavingsRate:    15,
		AchievementCount: 1,
		MilestoneCount:   1,
	}

	t.Run("filter by Critical", func(t *testing.T) {
		filtered := c.FilterByPriority(ctx, inputs, PriorityCritical)
		if len(filtered) == 0 {
			t.Fatal("expected critical insights")
		}
		for _, ins := range filtered {
			if ins.Priority != PriorityCritical {
				t.Errorf("expected all critical, got %q", ins.Priority)
			}
		}
	})

	t.Run("filter by Low", func(t *testing.T) {
		filtered := c.FilterByPriority(ctx, inputs, PriorityLow)
		if len(filtered) == 0 {
			t.Fatal("expected low insights")
		}
		for _, ins := range filtered {
			if ins.Priority != PriorityLow {
				t.Errorf("expected all low, got %q", ins.Priority)
			}
		}
	})

	t.Run("filter by non-existent", func(t *testing.T) {
		filtered := c.FilterByPriority(ctx, Inputs{}, PriorityInfo)
		if len(filtered) != 0 {
			t.Errorf("expected 0 info insights, got %d", len(filtered))
		}
	})
}

func TestSearch(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	inputs := Inputs{
		HealthScore:      50,
		HealthChange:     5,
		HealthGrade:      "B",
		CashFlowSurplus:  -10000,
		AchievementCount: 3,
	}

	t.Run("search by title", func(t *testing.T) {
		matched := c.Search(ctx, inputs, "Cash Flow")
		if len(matched) == 0 {
			t.Fatal("expected match for 'Cash Flow'")
		}
	})

	t.Run("search by summary", func(t *testing.T) {
		matched := c.Search(ctx, inputs, "shortfall")
		if len(matched) == 0 {
			t.Fatal("expected match for 'shortfall' in summary")
		}
	})

	t.Run("search case insensitive", func(t *testing.T) {
		matched := c.Search(ctx, inputs, "cash flow")
		if len(matched) == 0 {
			t.Fatal("expected case-insensitive match")
		}
	})

	t.Run("search no matches", func(t *testing.T) {
		matched := c.Search(ctx, inputs, "zzzznotfound")
		if len(matched) != 0 {
			t.Errorf("expected 0 matches, got %d", len(matched))
		}
	})

	t.Run("search empty query", func(t *testing.T) {
		matched := c.Search(ctx, inputs, "")
		if len(matched) == 0 {
			t.Error("expected all insights for empty query")
		}
		all := c.generate(ctx, inputs)
		if len(matched) != len(all) {
			t.Errorf("expected %d matches for empty query, got %d", len(all), len(matched))
		}
	})
}

func TestGetByID(t *testing.T) {
	c := NewComposer(nil)
	ctx := context.Background()

	inputs := Inputs{
		HealthScore:  70,
		HealthChange: 10,
		HealthGrade:  "A",
		SavingsRate:  30,
	}

	t.Run("found", func(t *testing.T) {
		got := c.GetByID(ctx, inputs, "ins-health-up")
		if got == nil {
			t.Fatal("expected to find ins-health-up")
		}
		if got.InsightID != "ins-health-up" {
			t.Errorf("id = %q", got.InsightID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		got := c.GetByID(ctx, inputs, "ins-nonexistent")
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestAddMeta(t *testing.T) {
	c := NewComposer(nil)
	base := Insight{InsightID: "test"}

	t.Run("fills defaults", func(t *testing.T) {
		got := c.addMeta(base, Inputs{}, "warning")
		if got.Timestamp == "" {
			t.Error("timestamp should not be empty")
		}
		if got.Severity != "warning" {
			t.Errorf("severity = %q", got.Severity)
		}
		if got.Confidence != "Medium" {
			t.Errorf("confidence = %q", got.Confidence)
		}
	})

	t.Run("preserves existing values", func(t *testing.T) {
		existing := Insight{
			InsightID:  "test2",
			Timestamp:  "2025-01-01T00:00:00Z",
			Severity:   "critical",
			Confidence: "High",
		}
		got := c.addMeta(existing, Inputs{}, "warning")
		if got.Timestamp != "2025-01-01T00:00:00Z" {
			t.Errorf("timestamp overwritten to %q", got.Timestamp)
		}
		if got.Severity != "critical" {
			t.Errorf("severity overwritten to %q", got.Severity)
		}
		if got.Confidence != "High" {
			t.Errorf("confidence overwritten to %q", got.Confidence)
		}
	})
}

func TestInsightCategoryConstants(t *testing.T) {
	expected := []InsightCategory{
		ICOpportunity, ICWarning, ICAchievement, ICMilestone,
		ICForecast, ICOptimization, ICRecommend, ICRisk,
		ICHealth, ICPortfolio, ICGoal, ICCashFlow,
		ICNetWorth, ICSavings, ICBehaviour, ICDebt,
	}
	if len(expected) != 16 {
		t.Errorf("expected 16 categories, got %d", len(expected))
	}
}

func TestPriorityConstants(t *testing.T) {
	priorities := []Priority{PriorityCritical, PriorityHigh, PriorityMedium, PriorityLow, PriorityInfo}
	if len(priorities) != 5 {
		t.Errorf("expected 5 priorities, got %d", len(priorities))
	}
}

func findByID(insights []Insight, id string) *Insight {
	for i := range insights {
		if insights[i].InsightID == id {
			return &insights[i]
		}
	}
	return nil
}
