package aggregator

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/horizon/core/services/experiences/planning/internal/engine"
)

// ---------------------------------------------------------------------------
// Mock providers
// ---------------------------------------------------------------------------

type mockGoalProvider struct {
	onTrack    int
	total      int
	fundingGap int64
	err        error
}

func (m *mockGoalProvider) GetGoalSummary(_ context.Context, _ string) (int, int, int64, error) {
	return m.onTrack, m.total, m.fundingGap, m.err
}

type mockAccountProvider struct {
	cash    int64
	income  int64
	expense int64
	err     error
}

func (m *mockAccountProvider) GetCashBalance(_ context.Context, _ string) (int64, int64, int64, error) {
	return m.cash, m.income, m.expense, m.err
}

type mockProjProvider struct {
	nw   int64
	inc  int64
	exp  int64
	pf   int64
	conf string
	err  error
}

func (m *mockProjProvider) GetPlanningProjection(_ context.Context, _ string) (int64, int64, int64, int64, string, error) {
	return m.nw, m.inc, m.exp, m.pf, m.conf, m.err
}

type mockRiskProvider struct {
	score int
	level string
	err   error
}

func (m *mockRiskProvider) GetRiskScore(_ context.Context, _ string) (int, string, error) {
	return m.score, m.level, m.err
}

type mockHealthProvider struct {
	score int
	grade string
	err   error
}

func (m *mockHealthProvider) GetHealthScore(_ context.Context, _ string) (int, string, error) {
	return m.score, m.grade, m.err
}

type mockRecProvider struct {
	has bool
	cnt int
	err error
}

func (m *mockRecProvider) HasRecommendations(_ context.Context, _ string) (bool, int, error) {
	return m.has, m.cnt, m.err
}

type mockOptProvider struct {
	has bool
	cnt int
	err error
}

func (m *mockOptProvider) HasOptimizations(_ context.Context, _ string) (bool, int, error) {
	return m.has, m.cnt, m.err
}

type mockSimProvider struct {
	has bool
	cnt int
	err error
}

func (m *mockSimProvider) HasSimulations(_ context.Context, _ string) (bool, int, error) {
	return m.has, m.cnt, m.err
}

type mockEventProvider struct {
	evtCnt int
	msCnt  int
	err    error
}

func (m *mockEventProvider) GetEventCount(_ context.Context, _ string) (int, int, error) {
	return m.evtCnt, m.msCnt, m.err
}

type mockScenarioProvider struct {
	scenarios []engine.Scenario
	err       error
}

func (m *mockScenarioProvider) GetScenarios(_ context.Context, _ string) ([]engine.Scenario, error) {
	return m.scenarios, m.err
}

type mockBudgetProvider struct {
	budgeted int64
	spent    int64
	rem      int64
	cats     []map[string]interface{}
	err      error
}

func (m *mockBudgetProvider) GetBudgetSummary(_ context.Context, _ string) (int64, int64, int64, []map[string]interface{}, error) {
	return m.budgeted, m.spent, m.rem, m.cats, m.err
}

type mockRecurringProvider struct {
	count int
	err   error
}

func (m *mockRecurringProvider) GetUpcomingCount(_ context.Context, _ string) (int, error) {
	return m.count, m.err
}

// ---------------------------------------------------------------------------
// Mock Cache
// ---------------------------------------------------------------------------

type mockCache struct {
	mu   sync.Mutex
	data map[string]interface{}
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string]interface{})}
}

func (m *mockCache) Get(_ context.Context, key string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *mockCache) Set(_ context.Context, key string, val interface{}) {
	m.mu.Lock()
	m.data[key] = val
	m.mu.Unlock()
}

func (m *mockCache) Invalidate(_ context.Context, key string) {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
}

// ---------------------------------------------------------------------------
// Default successful mock providers
// ---------------------------------------------------------------------------

func defaultProviders() DataProviders {
	return DataProviders{
		Goals: &mockGoalProvider{
			onTrack: 4, total: 6, fundingGap: 200000,
		},
		Accounts: &mockAccountProvider{
			cash: 50000, income: 150000, expense: 110000,
		},
		Projection: &mockProjProvider{
			nw: 5000000, inc: 1800000, exp: 1200000, pf: 3000000, conf: "Medium",
		},
		Risk: &mockRiskProvider{
			score: 3, level: "Moderate",
		},
		Health: &mockHealthProvider{
			score: 78, grade: "B+",
		},
		Recs: &mockRecProvider{
			has: true, cnt: 5,
		},
		Optimize: &mockOptProvider{
			has: true, cnt: 3,
		},
		Simulation: &mockSimProvider{
			has: true, cnt: 2,
		},
		Events: &mockEventProvider{
			evtCnt: 10, msCnt: 3,
		},
		Scenarios: &mockScenarioProvider{
			scenarios: []engine.Scenario{
				{ScenarioID: "s1", Name: "Optimistic"},
			},
		},
		Budget: &mockBudgetProvider{
			budgeted: 100000, spent: 75000, rem: 25000,
			cats: []map[string]interface{}{{"cat": "Food", "amount": float64(30000)}},
		},
		Recurring: &mockRecurringProvider{
			count: 4,
		},
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	a := New(defaultProviders(), newMockCache())
	if a == nil {
		t.Fatal("New() returned nil")
	}
}

func TestAggregate_Success(t *testing.T) {
	a := New(defaultProviders(), newMockCache())
	ctx := context.Background()

	in, err := a.Aggregate(ctx, "user-42", engine.PMOverview)
	if err != nil {
		t.Fatalf("Aggregate returned error: %v", err)
	}
	if in.UserID != "user-42" {
		t.Errorf("UserID = %q", in.UserID)
	}
	if in.Mode != engine.PMOverview {
		t.Errorf("Mode = %q", in.Mode)
	}

	// Goals
	if in.GoalsOnTrack != 4 || in.TotalGoals != 6 || in.FundingGap != 200000 {
		t.Errorf("goals: onTrack=%d total=%d gap=%d", in.GoalsOnTrack, in.TotalGoals, in.FundingGap)
	}
	// Accounts
	if in.CashBalance != 50000 || in.MonthlyIncome != 150000 || in.MonthlyExpenses != 110000 {
		t.Errorf("accounts: cash=%d inc=%d exp=%d", in.CashBalance, in.MonthlyIncome, in.MonthlyExpenses)
	}
	// Projection
	if in.NetWorthProjected != 5000000 || in.IncomeProjected != 1800000 || in.ExpensesProjected != 1200000 || in.PortfolioValue != 3000000 {
		t.Errorf("projection values mismatch")
	}
	if !in.HasProjection {
		t.Error("HasProjection should be true")
	}
	// Risk
	if in.RiskScore != 3 || in.RiskLevel != "Moderate" {
		t.Errorf("risk: score=%d level=%q", in.RiskScore, in.RiskLevel)
	}
	// Health
	if in.HealthScore != 78 || in.HealthGrade != "B+" {
		t.Errorf("health: score=%d grade=%q", in.HealthScore, in.HealthGrade)
	}
	// Recs
	if !in.HasRecs || in.RecCount != 5 {
		t.Errorf("recs: has=%v count=%d", in.HasRecs, in.RecCount)
	}
	// Opt
	if !in.HasOptimization || in.OptCount != 3 {
		t.Errorf("opt: has=%v count=%d", in.HasOptimization, in.OptCount)
	}
	// Sim
	if !in.HasSimulation || in.SimCount != 2 {
		t.Errorf("sim: has=%v count=%d", in.HasSimulation, in.SimCount)
	}
	// Events
	if in.EventCount != 10 || in.MilestoneCount != 3 {
		t.Errorf("events: evt=%d ms=%d", in.EventCount, in.MilestoneCount)
	}
	// Scenarios
	if len(in.Scenarios) != 1 || !in.HasAlternatives {
		t.Errorf("scenarios: len=%d hasAlt=%v", len(in.Scenarios), in.HasAlternatives)
	}
	// Budget
	if in.TotalBudgeted != 100000 || in.TotalSpent != 75000 || in.TotalRemaining != 25000 {
		t.Errorf("budget values mismatch")
	}
	if len(in.BudgetCategories) != 1 {
		t.Errorf("budget categories len = %d", len(in.BudgetCategories))
	}
	// Recurring
	if in.UpcomingRecurring != 4 {
		t.Errorf("recurring = %d", in.UpcomingRecurring)
	}
}

func TestAggregate_CacheHit(t *testing.T) {
	cache := newMockCache()
	a := New(defaultProviders(), cache)

	ctx := context.Background()

	// populate cache manually
	orig := &engine.Inputs{UserID: "cached-user", Mode: engine.PMOverview, GoalsOnTrack: 99}
	cache.Set(ctx, "plan:cached-user:Overview", orig)

	in, err := a.Aggregate(ctx, "cached-user", engine.PMOverview)
	if err != nil {
		t.Fatalf("Aggregate error: %v", err)
	}
	if in.GoalsOnTrack != 99 {
		t.Errorf("expected cached value 99, got %d", in.GoalsOnTrack)
	}
	if in.UserID != "cached-user" {
		t.Errorf("UserID = %q", in.UserID)
	}
}

func TestAggregate_ProviderError(t *testing.T) {
	p := defaultProviders()
	p.Goals = &mockGoalProvider{err: errors.New("goals down")}
	a := New(p, newMockCache())

	_, err := a.Aggregate(context.Background(), "user-x", engine.PMGoalPlanning)
	if err == nil {
		t.Fatal("expected error from goals provider")
	}
	if err.Error() != "goals down" {
		t.Errorf("error = %q", err.Error())
	}
}

func TestAggregate_ProjectionError(t *testing.T) {
	p := defaultProviders()
	p.Projection = &mockProjProvider{err: errors.New("projection unavailable")}
	a := New(p, newMockCache())

	_, err := a.Aggregate(context.Background(), "user-x", engine.PMRetirement)
	if err == nil {
		t.Fatal("expected error from projection provider")
	}
}

func TestAggregate_NilOptionalProviders(t *testing.T) {
	p := defaultProviders()
	p.Budget = nil
	p.Recurring = nil
	a := New(p, newMockCache())

	in, err := a.Aggregate(context.Background(), "user-nil", engine.PMOverview)
	if err != nil {
		t.Fatalf("Aggregate error: %v", err)
	}
	if in.TotalBudgeted != 0 || in.TotalSpent != 0 {
		t.Error("budget should be zero-valued when provider is nil")
	}
	if in.UpcomingRecurring != 0 {
		t.Error("recurring should be zero-valued when provider is nil")
	}
	// Non-optional fields should still be populated
	if in.GoalsOnTrack != 4 {
		t.Errorf("goals should still populate, got %d", in.GoalsOnTrack)
	}
}

func TestAggregate_OptionalBudgetErrorIgnored(t *testing.T) {
	p := defaultProviders()
	p.Budget = &mockBudgetProvider{err: errors.New("budget service down")}
	a := New(p, newMockCache())

	in, err := a.Aggregate(context.Background(), "user-budget-err", engine.PMDebtPlan)
	if err != nil {
		t.Fatalf("Aggregate should not fail on optional budget error: %v", err)
	}
	if in.TotalBudgeted != 0 {
		t.Error("budget should be zero on error")
	}
	// Other fields intact
	if in.RiskScore != 3 {
		t.Errorf("risk = %d", in.RiskScore)
	}
}

func TestAggregate_OptionalRecurringErrorIgnored(t *testing.T) {
	p := defaultProviders()
	p.Recurring = &mockRecurringProvider{err: errors.New("recurring down")}
	a := New(p, newMockCache())

	in, err := a.Aggregate(context.Background(), "user-rec-err", engine.PMInvestmentPlan)
	if err != nil {
		t.Fatalf("Aggregate should not fail on optional recurring error: %v", err)
	}
	if in.UpcomingRecurring != 0 {
		t.Error("recurring should be zero on error")
	}
}

func TestAggregate_CacheIsSetOnSuccess(t *testing.T) {
	cache := newMockCache()
	a := New(defaultProviders(), cache)

	ctx := context.Background()
	in, err := a.Aggregate(ctx, "cache-me", engine.PMOptimize)
	if err != nil {
		t.Fatalf("Aggregate error: %v", err)
	}
	if in == nil {
		t.Fatal("result is nil")
	}

	cached, ok := cache.Get(ctx, "plan:cache-me:Optimization")
	if !ok {
		t.Fatal("expected value in cache after aggregate")
	}
	casted, ok := cached.(*engine.Inputs)
	if !ok {
		t.Fatal("cached value is not *engine.Inputs")
	}
	if casted.UserID != "cache-me" {
		t.Errorf("cached UserID = %q", casted.UserID)
	}
}

func TestAggregate_ConcurrentSameUser(t *testing.T) {
	a := New(defaultProviders(), newMockCache())
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := a.Aggregate(ctx, "concurrent-user", engine.PMOverview)
			if err != nil {
				t.Errorf("concurrent aggregate error: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestRefresh_ClearsAllModes(t *testing.T) {
	cache := newMockCache()
	a := New(defaultProviders(), cache)

	ctx := context.Background()

	expectedModes := []engine.PlanningMode{
		engine.PMOverview, engine.PMGoalPlanning, engine.PMRetirement,
		engine.PMDebtPlan, engine.PMInvestmentPlan, engine.PMScenarioPlan, engine.PMOptimize,
	}

	// populate cache for all modes
	for _, mode := range expectedModes {
		cache.Set(ctx, "plan:refresh-user:"+string(mode), &engine.Inputs{Mode: mode})
	}
	// extra mode that should NOT be cleared
	cache.Set(ctx, "plan:refresh-user:EducationPlanning", &engine.Inputs{Mode: engine.PMEducationPlan})

	a.Refresh(ctx, "refresh-user")

	for _, mode := range expectedModes {
		_, ok := cache.Get(ctx, "plan:refresh-user:"+string(mode))
		if ok {
			t.Errorf("expected mode %q to be invalidated", mode)
		}
	}
	// EducationPlanning should not have been cleared
	_, ok := cache.Get(ctx, "plan:refresh-user:EducationPlanning")
	if !ok {
		t.Error("EducationPlanning should NOT have been invalidated by Refresh")
	}
}

func TestAggregate_ModePreservation(t *testing.T) {
	a := New(defaultProviders(), newMockCache())

	modes := []engine.PlanningMode{
		engine.PMOverview,
		engine.PMGoalPlanning,
		engine.PMRetirement,
		engine.PMDebtPlan,
		engine.PMInvestmentPlan,
		engine.PMScenarioPlan,
		engine.PMOptimize,
		engine.PMHistoricalComp,
	}

	for _, mode := range modes {
		in, err := a.Aggregate(context.Background(), "mode-test", mode)
		if err != nil {
			t.Fatalf("Aggregate mode %q error: %v", mode, err)
		}
		if in.Mode != mode {
			t.Errorf("expected mode %q, got %q", mode, in.Mode)
		}
	}
}

func TestAggregate_ScenariosEmptyNoAlternatives(t *testing.T) {
	p := defaultProviders()
	p.Scenarios = &mockScenarioProvider{scenarios: []engine.Scenario{}}
	a := New(p, newMockCache())

	in, err := a.Aggregate(context.Background(), "user-empty", engine.PMScenarioPlan)
	if err != nil {
		t.Fatalf("Aggregate error: %v", err)
	}
	if in.HasAlternatives {
		t.Error("HasAlternatives should be false with empty scenarios")
	}
	if len(in.Scenarios) != 0 {
		t.Errorf("expected 0 scenarios, got %d", len(in.Scenarios))
	}
}

func TestAggregate_ScenariosNilNoAlternatives(t *testing.T) {
	p := defaultProviders()
	p.Scenarios = &mockScenarioProvider{scenarios: nil}
	a := New(p, newMockCache())

	in, err := a.Aggregate(context.Background(), "user-nil-sc", engine.PMGoalPlanning)
	if err != nil {
		t.Fatalf("Aggregate error: %v", err)
	}
	if in.HasAlternatives {
		t.Error("HasAlternatives should be false with nil scenarios")
	}
}

func TestAggregate_FirstErrorWins(t *testing.T) {
	errGoal := errors.New("goal error")
	errRisk := errors.New("risk error")

	p := defaultProviders()
	p.Goals = &mockGoalProvider{err: errGoal}
	p.Risk = &mockRiskProvider{err: errRisk}
	a := New(p, newMockCache())

	_, err := a.Aggregate(context.Background(), "multi-err", engine.PMDebtPlan)
	if err == nil {
		t.Fatal("expected an error")
	}
	// Either of the two errors may surface depending on goroutine scheduling
	if err.Error() != "goal error" && err.Error() != "risk error" {
		t.Errorf("unexpected error %q", err)
	}
}
