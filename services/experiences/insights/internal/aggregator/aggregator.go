package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/insights/internal/engine"
)

type DataProviders struct {
	Health     HealthProvider
	Risk       RiskProvider
	Projection ProjProvider
	Goals      GoalProvider
	Portfolio  PfProvider
	Recs       RecProvider
	Optimize   OptProvider
	Simulation SimProvider
	Events     EventProvider
	Accounts   AcctProvider
}

type HealthProvider interface {
	GetHealthScore(ctx context.Context, userID string) (score int, grade string, err error)
	GetHealthChange(ctx context.Context, userID string) (change int, err error)
}
type RiskProvider interface {
	GetRiskScore(ctx context.Context, userID string) (score int, level string, err error)
}
type ProjProvider interface {
	GetSavingsRate(ctx context.Context, userID string) (rate float64, trend string, err error)
	GetNetWorthChange(ctx context.Context, userID string) (change int64, netWorth int64, err error)
	GetCashFlowSurplus(ctx context.Context, userID string) (surplus int64, err error)
}
type GoalProvider interface {
	GetGoalProgress(ctx context.Context, userID string) (onTrack, total int, avgPct float64, err error)
	GetMilestoneCount(ctx context.Context, userID string) (int, error)
}
type PfProvider interface {
	GetPortfolioReturn(ctx context.Context, userID string) (value int64, ret float64, err error)
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
	GetAchievementCount(ctx context.Context, userID string) (int, error)
	GetSpendingAnomaly(ctx context.Context, userID string) (string, error)
}
type AcctProvider interface {
	GetAccountCount(ctx context.Context, userID string) (int, error)
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

func (a *Aggregator) Aggregate(ctx context.Context, userID string) (*engine.Inputs, error) {
	cacheKey := "ins:" + userID
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached.(*engine.Inputs), nil
	}

	inputs := &engine.Inputs{UserID: userID}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 14)

	wg.Add(1); go func() {
		defer wg.Done()
		score, grade, err := a.providers.Health.GetHealthScore(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HealthScore = score; inputs.HealthGrade = grade; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		change, err := a.providers.Health.GetHealthChange(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HealthChange = change; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		score, level, err := a.providers.Risk.GetRiskScore(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.RiskScore = score; inputs.RiskLevel = level; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		rate, trend, err := a.providers.Projection.GetSavingsRate(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.SavingsRate = rate; inputs.SavingsTrend = trend; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		change, nw, err := a.providers.Projection.GetNetWorthChange(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.NetWorthChange = change; inputs.NetWorth = nw; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		surplus, err := a.providers.Projection.GetCashFlowSurplus(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.CashFlowSurplus = surplus; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		onTrack, total, avgPct, err := a.providers.Goals.GetGoalProgress(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.GoalsOnTrack = onTrack; inputs.TotalGoals = total; inputs.GoalProgressAvg = avgPct; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		cnt, err := a.providers.Goals.GetMilestoneCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.MilestoneCount = cnt; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		val, ret, err := a.providers.Portfolio.GetPortfolioReturn(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.PortfolioValue = val; inputs.PortfolioReturn = ret; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, err := a.providers.Recs.HasRecommendations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasRecs = has; inputs.RecCount = cnt; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		has, _, err := a.providers.Optimize.HasOptimizations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasOpts = has; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		cnt, err := a.providers.Events.GetAchievementCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.AchievementCount = cnt; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		anomaly, err := a.providers.Events.GetSpendingAnomaly(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.SpendingAnomaly = anomaly; mu.Unlock()
	}()

	wg.Wait()
	close(errs)
	for e := range errs { if e != nil { return nil, e } }

	a.cache.Set(ctx, cacheKey, inputs)
	return inputs, nil
}

func (a *Aggregator) Refresh(ctx context.Context, userID string) {
	a.cache.Invalidate(ctx, "ins:"+userID)
}
