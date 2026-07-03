package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/planning/internal/engine"
)

type DataProviders struct {
	Goals       GoalProvider
	Accounts    AccountProvider
	Projection  ProjProvider
	Risk        RiskProvider
	Health      HealthProvider
	Recs        RecProvider
	Optimize    OptProvider
	Simulation  SimProvider
	Events      EventProvider
	Scenarios   ScenarioProvider
	Budget      BudgetProvider
	Recurring   RecurringProvider
}

type RecurringProvider interface {
	GetUpcomingCount(ctx context.Context, userID string) (int, error)
}

type GoalProvider interface {
	GetGoalSummary(ctx context.Context, userID string) (onTrack, total int, fundingGap int64, err error)
}
type AccountProvider interface {
	GetCashBalance(ctx context.Context, userID string) (int64, int64, int64, error)
}
type ProjProvider interface {
	GetPlanningProjection(ctx context.Context, userID string) (netWorth, income, expenses, portfolio int64, conf string, err error)
}
type RiskProvider interface {
	GetRiskScore(ctx context.Context, userID string) (score int, level string, err error)
}
type HealthProvider interface {
	GetHealthScore(ctx context.Context, userID string) (score int, grade string, err error)
}
type RecProvider interface {
	HasRecommendations(ctx context.Context, userID string) (bool, int, error)
}
type OptProvider interface {
	HasOptimizations(ctx context.Context, userID string) (bool, int, error)
}
type SimProvider interface {
	HasSimulations(ctx context.Context, userID string) (bool, int, error)
}
type EventProvider interface {
	GetEventCount(ctx context.Context, userID string) (int, int, error)
}
type ScenarioProvider interface {
	GetScenarios(ctx context.Context, userID string) ([]engine.Scenario, error)
}
type BudgetProvider interface {
	GetBudgetSummary(ctx context.Context, userID string) (totalBudgeted, totalSpent, totalRemaining int64, categories []map[string]interface{}, err error)
}

type Aggregator struct {
	providers DataProviders
	cache     Cache
}

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, val interface{})
	Invalidate(ctx context.Context, key string)
}

func New(providers DataProviders, cache Cache) *Aggregator {
	return &Aggregator{providers: providers, cache: cache}
}

func (a *Aggregator) Aggregate(ctx context.Context, userID string, mode engine.PlanningMode) (*engine.Inputs, error) {
	cacheKey := "plan:" + userID + ":" + string(mode)
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached.(*engine.Inputs), nil
	}

	inputs := &engine.Inputs{UserID: userID, Mode: mode}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 12)

	wg.Add(1); go func() {
		defer wg.Done()
		onTrack, total, gap, err := a.providers.Goals.GetGoalSummary(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.GoalsOnTrack = onTrack; inputs.TotalGoals = total; inputs.FundingGap = gap; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		cash, income, expenses, err := a.providers.Accounts.GetCashBalance(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.CashBalance = cash; inputs.MonthlyIncome = income; inputs.MonthlyExpenses = expenses; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		nw, inc, exp, pf, _, err := a.providers.Projection.GetPlanningProjection(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.NetWorthProjected = nw; inputs.IncomeProjected = inc; inputs.ExpensesProjected = exp; inputs.PortfolioValue = pf; inputs.HasProjection = true; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		score, level, err := a.providers.Risk.GetRiskScore(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.RiskScore = score; inputs.RiskLevel = level; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		score, grade, err := a.providers.Health.GetHealthScore(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HealthScore = score; inputs.HealthGrade = grade; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, err := a.providers.Recs.HasRecommendations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasRecs = has; inputs.RecCount = cnt; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, err := a.providers.Optimize.HasOptimizations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasOptimization = has; inputs.OptCount = cnt; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, err := a.providers.Simulation.HasSimulations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasSimulation = has; inputs.SimCount = cnt; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		evtCnt, msCnt, err := a.providers.Events.GetEventCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.EventCount = evtCnt; inputs.MilestoneCount = msCnt; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		scenarios, err := a.providers.Scenarios.GetScenarios(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.Scenarios = scenarios; inputs.HasAlternatives = len(scenarios) > 0; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if a.providers.Budget != nil {
			budgeted, spent, rem, cats, err := a.providers.Budget.GetBudgetSummary(ctx, userID)
			if err == nil {
				mu.Lock()
				inputs.TotalBudgeted = budgeted
				inputs.TotalSpent = spent
				inputs.TotalRemaining = rem
				inputs.BudgetCategories = cats
				mu.Unlock()
			}
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if a.providers.Recurring != nil {
			count, err := a.providers.Recurring.GetUpcomingCount(ctx, userID)
			if err == nil {
				mu.Lock(); inputs.UpcomingRecurring = count; mu.Unlock()
			}
		}
	}()

	wg.Wait()
	close(errs)
	for e := range errs { if e != nil { return nil, e } }

	a.cache.Set(ctx, cacheKey, inputs)
	return inputs, nil
}

func (a *Aggregator) Refresh(ctx context.Context, userID string) {
	for _, mode := range []engine.PlanningMode{engine.PMOverview, engine.PMGoalPlanning, engine.PMRetirement, engine.PMDebtPlan, engine.PMInvestmentPlan, engine.PMScenarioPlan, engine.PMOptimize} {
		a.cache.Invalidate(ctx, "plan:"+userID+":"+string(mode))
	}
}
