package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/dashboard/internal/engine"
)

// DataProviders holds all interfaces for fetching domain/engine data.
type DataProviders struct {
	Accounts    AccountProvider
	Assets      AssetProvider
	Liabilities LiabilityProvider
	Goals       GoalProvider
	Events      EventProvider
	Portfolio   PortfolioProvider
	Health      HealthProvider
	Risk        RiskProvider
	Recs        RecommendationProvider
	Projection  ProjectionProvider
	Simulation  SimulationProvider
	Optimization OptimizationProvider
}

type AccountProvider interface {
	GetCashBalance(ctx context.Context, userID string) (int64, error)
	HasAccounts(ctx context.Context, userID string) (bool, error)
}
type AssetProvider interface {
	GetTotalAssets(ctx context.Context, userID string) (int64, error)
}
type LiabilityProvider interface {
	GetTotalLiabilities(ctx context.Context, userID string) (int64, error)
}
type GoalProvider interface {
	GetGoalCounts(ctx context.Context, userID string) (onTrack, atRisk, total int, hasGoals bool, err error)
}
type EventProvider interface {
	GetRecentCount(ctx context.Context, userID string) (int, bool, error)
}
type PortfolioProvider interface {
	GetPortfolioValue(ctx context.Context, userID string) (int64, error)
}
type HealthProvider interface {
	GetHealthScore(ctx context.Context, userID string) (score int, grade string, err error)
}
type RiskProvider interface {
	GetRiskScore(ctx context.Context, userID string) (score int, level string, err error)
}
type RecommendationProvider interface {
	HasRecommendations(ctx context.Context, userID string) (bool, int, error)
}
type ProjectionProvider interface {
	GetMonthlyIncomeExpenses(ctx context.Context, userID string) (income, expenses int64, err error)
}
type SimulationProvider interface {
	HasSimulations(ctx context.Context, userID string) (bool, error)
}
type OptimizationProvider interface{}

// DashboardAggregator loads data from all providers in parallel.
type DashboardAggregator struct {
	providers DataProviders
	cache     Cache
}

type Cache interface {
	Get(ctx context.Context, key string) (*engine.Inputs, bool)
	Set(ctx context.Context, key string, inputs *engine.Inputs)
}

func New(providers DataProviders, cache Cache) *DashboardAggregator {
	return &DashboardAggregator{providers: providers, cache: cache}
}

// Aggregate fetches all dashboard data in parallel goroutines.
// Errors from individual providers are silently skipped (best-effort).
func (a *DashboardAggregator) Aggregate(ctx context.Context, userID string) (*engine.Inputs, error) {
	if cached, ok := a.cache.Get(ctx, userID); ok {
		return cached, nil
	}

	inputs := &engine.Inputs{UserID: userID}
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1); go func() {
		defer wg.Done()
		if b, err := a.providers.Accounts.HasAccounts(ctx, userID); err == nil {
			mu.Lock(); inputs.HasAccounts = b; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if v, err := a.providers.Accounts.GetCashBalance(ctx, userID); err == nil {
			mu.Lock(); inputs.CashBalance = v; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if v, err := a.providers.Assets.GetTotalAssets(ctx, userID); err == nil {
			mu.Lock(); inputs.TotalAssets = v; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if v, err := a.providers.Liabilities.GetTotalLiabilities(ctx, userID); err == nil {
			mu.Lock(); inputs.TotalLiabilities = v; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if onTrack, atRisk, total, hasGoals, err := a.providers.Goals.GetGoalCounts(ctx, userID); err == nil {
			mu.Lock()
			inputs.GoalsOnTrack = onTrack
			inputs.GoalsAtRisk = atRisk
			inputs.GoalCount = total
			inputs.HasGoals = hasGoals
			mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if cnt, has, err := a.providers.Events.GetRecentCount(ctx, userID); err == nil {
			mu.Lock(); inputs.RecentEvents = cnt; inputs.HasEvents = has; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if v, err := a.providers.Portfolio.GetPortfolioValue(ctx, userID); err == nil {
			mu.Lock(); inputs.PortfolioValue = v; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if score, grade, err := a.providers.Health.GetHealthScore(ctx, userID); err == nil {
			mu.Lock(); inputs.HealthScore = score; inputs.HealthGrade = grade; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if score, level, err := a.providers.Risk.GetRiskScore(ctx, userID); err == nil {
			mu.Lock(); inputs.RiskScore = score; inputs.RiskLevel = level; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if has, cnt, err := a.providers.Recs.HasRecommendations(ctx, userID); err == nil {
			mu.Lock(); inputs.HasRecommendation = has; inputs.RecCount = cnt; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if income, expenses, err := a.providers.Projection.GetMonthlyIncomeExpenses(ctx, userID); err == nil {
			mu.Lock(); inputs.MonthlyIncome = income; inputs.MonthlyExpenses = expenses; mu.Unlock()
		}
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		if hasSim, err := a.providers.Simulation.HasSimulations(ctx, userID); err == nil {
			mu.Lock(); inputs.HasSimulation = hasSim; mu.Unlock()
		}
	}()

	wg.Wait()

	inputs.NetWorth = inputs.TotalAssets + inputs.CashBalance - inputs.TotalLiabilities
	if inputs.CashBalance > 0 || inputs.TotalAssets > 0 {
		inputs.IsFirstTimeUser = false
	}
	if inputs.GoalCount == 0 && !inputs.HasAccounts {
		inputs.IsFirstTimeUser = true
	}

	a.cache.Set(ctx, userID, inputs)
	return inputs, nil
}

// Refresh forces a fresh aggregation by clearing cache then re-aggregating.
func (a *DashboardAggregator) Refresh(ctx context.Context, userID string) (*engine.Inputs, error) {
	a.cache.Set(ctx, userID, nil)
	return a.Aggregate(ctx, userID)
}
