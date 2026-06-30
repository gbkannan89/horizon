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
func (a *DashboardAggregator) Aggregate(ctx context.Context, userID string) (*engine.Inputs, error) {
	if cached, ok := a.cache.Get(ctx, userID); ok {
		return cached, nil
	}

	inputs := &engine.Inputs{UserID: userID}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 12)

	fetch := func(fn func() error) { defer wg.Done(); if err := fn(); err != nil { errs <- err } }

	wg.Add(1); go fetch(func() error {
		b, err := a.providers.Accounts.HasAccounts(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.HasAccounts = b; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		v, err := a.providers.Accounts.GetCashBalance(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.CashBalance = v; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		v, err := a.providers.Assets.GetTotalAssets(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.TotalAssets = v; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		v, err := a.providers.Liabilities.GetTotalLiabilities(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.TotalLiabilities = v; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		onTrack, atRisk, total, hasGoals, err := a.providers.Goals.GetGoalCounts(ctx, userID)
		if err != nil { return err }
		mu.Lock()
		inputs.GoalsOnTrack = onTrack
		inputs.GoalsAtRisk = atRisk
		inputs.GoalCount = total
		inputs.HasGoals = hasGoals
		mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		cnt, has, err := a.providers.Events.GetRecentCount(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.RecentEvents = cnt; inputs.HasEvents = has; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		v, err := a.providers.Portfolio.GetPortfolioValue(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.PortfolioValue = v; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		score, grade, err := a.providers.Health.GetHealthScore(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.HealthScore = score; inputs.HealthGrade = grade; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		score, level, err := a.providers.Risk.GetRiskScore(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.RiskScore = score; inputs.RiskLevel = level; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		has, cnt, err := a.providers.Recs.HasRecommendations(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.HasRecommendation = has; inputs.RecCount = cnt; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		income, expenses, err := a.providers.Projection.GetMonthlyIncomeExpenses(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.MonthlyIncome = income; inputs.MonthlyExpenses = expenses; mu.Unlock()
		return nil
	})

	wg.Add(1); go fetch(func() error {
		hasSim, err := a.providers.Simulation.HasSimulations(ctx, userID)
		if err != nil { return err }
		mu.Lock(); inputs.HasSimulation = hasSim; mu.Unlock()
		return nil
	})

	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil { return nil, e }
	}

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
