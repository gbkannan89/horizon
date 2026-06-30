package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/notifications/internal/engine"
)

type DataProviders struct {
	Health     HealthProvider
	Risk       RiskProvider
	Projection ProjProvider
	Goals      GoalProvider
	Recs       RecProvider
	Optimize   OptProvider
	Events     EventProvider
	Prefs      PrefProvider
}

type HealthProvider interface {
	GetHealthScore(ctx context.Context, userID string) (int, string, error)
	GetHealthChange(ctx context.Context, userID string) (int, error)
}
type RiskProvider interface {
	GetRiskScore(ctx context.Context, userID string) (int, string, error)
}
type ProjProvider interface {
	GetCashFlowSurplus(ctx context.Context, userID string) (int64, error)
	GetNetWorthChange(ctx context.Context, userID string) (int64, error)
}
type GoalProvider interface {
	GetGoalProgress(ctx context.Context, userID string) (onTrack, total, atRisk int, err error)
}
type RecProvider interface {
	HasRecommendations(ctx context.Context, userID string) (bool, int, error)
}
type OptProvider interface {
	HasOptimizations(ctx context.Context, userID string) (bool, error)
}
type EventProvider interface {
	GetAchievementCount(ctx context.Context, userID string) (int, error)
}
type PrefProvider interface {
	GetPreferences(ctx context.Context, userID string) ([]engine.Preference, error)
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
	cacheKey := "notif:" + userID
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached.(*engine.Inputs), nil
	}

	inputs := &engine.Inputs{UserID: userID}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 12)

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
		onTrack, total, atRisk, err := a.providers.Goals.GetGoalProgress(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.GoalOnTrack = onTrack; inputs.TotalGoals = total; inputs.GoalRisk = atRisk; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		surplus, err := a.providers.Projection.GetCashFlowSurplus(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.CashSurplus = surplus; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		change, err := a.providers.Projection.GetNetWorthChange(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.NetWorthChg = change; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, err := a.providers.Recs.HasRecommendations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasRecs = has; inputs.RecCount = cnt; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		has, err := a.providers.Optimize.HasOptimizations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasOpts = has; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		cnt, err := a.providers.Events.GetAchievementCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.AchieveCnt = cnt; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		prefs, err := a.providers.Prefs.GetPreferences(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.Preferences = prefs; mu.Unlock()
	}()

	wg.Wait()
	close(errs)
	for e := range errs { if e != nil { return nil, e } }

	a.cache.Set(ctx, cacheKey, inputs)
	return inputs, nil
}

func (a *Aggregator) Refresh(ctx context.Context, userID string) {
	a.cache.Invalidate(ctx, "notif:"+userID)
}
