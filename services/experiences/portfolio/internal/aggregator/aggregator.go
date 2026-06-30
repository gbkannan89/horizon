package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/portfolio/internal/engine"
)

// DataProviders holds interfaces for fetching portfolio-related data.
type DataProviders struct {
	Portfolio   PortfolioProvider
	Assets      AssetProvider
	Allocation  AllocationProvider
	Performance PerfProvider
	Risk        RiskProvider
	Projection  ProjProvider
	Recs        RecProvider
	Optimize    OptProvider
	Simulation  SimProvider
	Events      EventProvider
}

type PortfolioProvider interface {
	GetPortfolioSummary(ctx context.Context, userID string) (value int64, totalReturn float64, returnPct float64, err error)
}
type AssetProvider interface {
	GetAllocations(ctx context.Context, userID string) ([]engine.AllocEntry, error)
}
type AllocationProvider interface{}
type PerfProvider interface {
	GetPerformance(ctx context.Context, userID string) (periodReturn, periodReturnPct, benchmarkRet float64, unrealizedGL, realizedGL int64, err error)
}
type RiskProvider interface {
	GetPortfolioRisk(ctx context.Context, userID string) (score int, level string, sharpe, vol, mdd, valueAtRisk float64, err error)
}
type ProjProvider interface {
	GetPortfolioProjection(ctx context.Context, userID string) (projectedVal float64, confidence string, horizonYears int, annualReturn float64, err error)
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
	GetEventCount(ctx context.Context, userID string) (int, error)
	GetMilestoneCount(ctx context.Context, userID string) (int, error)
}

// Aggregator loads portfolio data from all providers in parallel.
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
	cacheKey := "pf:" + userID
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached.(*engine.Inputs), nil
	}

	inputs := &engine.Inputs{UserID: userID}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 12)

	wg.Add(1); go func() {
		defer wg.Done()
		val, ret, retPct, err := a.providers.Portfolio.GetPortfolioSummary(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.PortfolioValue = val; inputs.TotalReturn = ret; inputs.TotalReturnPct = retPct; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		allocs, err := a.providers.Assets.GetAllocations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.Allocations = allocs; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		periodRet, periodRetPct, benchRet, unrealizedGL, realizedGL, err := a.providers.Performance.GetPerformance(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.PeriodReturn = periodRet; inputs.PeriodReturnPct = periodRetPct; inputs.BenchmarkRet = benchRet; inputs.UnrealizedGL = unrealizedGL; inputs.RealizedGL = realizedGL; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		score, level, sharpe, vol, mdd, valueAtRisk, err := a.providers.Risk.GetPortfolioRisk(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.RiskScore = score; inputs.RiskLevel = level; inputs.SharpeRatio = sharpe; inputs.Volatility = vol; inputs.MaxDrawdown = mdd; inputs.Var = valueAtRisk; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		projVal, conf, horizon, annRet, err := a.providers.Projection.GetPortfolioProjection(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.ProjectedVal = projVal; inputs.ProjectionConf = conf; inputs.HorizonYears = horizon; inputs.AnnualReturn = annRet; mu.Unlock()
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
		mu.Lock(); inputs.HasOpts = has; inputs.OptCount = cnt; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, err := a.providers.Simulation.HasSimulations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasSims = has; inputs.SimCount = cnt; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		cnt, err := a.providers.Events.GetEventCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.EventCount = cnt; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		cnt, err := a.providers.Events.GetMilestoneCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.MilestoneCount = cnt; mu.Unlock()
	}()

	wg.Wait()
	close(errs)
	for e := range errs { if e != nil { return nil, e } }

	a.cache.Set(ctx, cacheKey, inputs)
	return inputs, nil
}

func (a *Aggregator) Refresh(ctx context.Context, userID string) {
	a.cache.Invalidate(ctx, "pf:"+userID)
}
