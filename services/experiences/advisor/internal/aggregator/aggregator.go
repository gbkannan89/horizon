package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/advisor/internal/engine"
)

type DataProviders struct {
	Dashboard   DashProvider
	Goals       GoalProvider
	Accounts    AcctProvider
	Portfolio   PfProvider
	Planning    PlanProvider
	Health      HealthProvider
	Risk        RiskProvider
	Projection  ProjProvider
	Recs        RecProvider
	Optimize    OptProvider
	Simulation  SimProvider
	Timeline    TimelineProvider
	Insights    InsightProvider
	Notifs      NotifProvider
}

type DashProvider interface {
	GetNetWorth(ctx context.Context, userID string) (netWorth int64, err error)
	GetCashBalance(ctx context.Context, userID string) (cash int64, err error)
	GetMonthlyFlow(ctx context.Context, userID string) (income, expenses int64, err error)
}
type GoalProvider interface {
	GetGoalProgress(ctx context.Context, userID string) (onTrack, atRisk, total int, fundingGap int64, err error)
}
type AcctProvider interface {
	GetAccountSummary(ctx context.Context, userID string) (count int, totalBalance int64, err error)
}
type PfProvider interface {
	GetPortfolioSummary(ctx context.Context, userID string) (value int64, totalReturn, returnPct float64, err error)
}
type PlanProvider interface{}
type HealthProvider interface {
	GetHealthScore(ctx context.Context, userID string) (int, string, error)
	GetHealthChange(ctx context.Context, userID string) (int, error)
}
type RiskProvider interface {
	GetRiskScore(ctx context.Context, userID string) (int, string, error)
}
type ProjProvider interface {
	GetProjectionSummary(ctx context.Context, userID string) (nwProjected int64, onTrack bool, confidence string, err error)
}
type RecProvider interface {
	GetRecommendationSummary(ctx context.Context, userID string) (has bool, count int, title string, err error)
}
type OptProvider interface {
	HasOptimizations(ctx context.Context, userID string) (bool, int, error)
}
type SimProvider interface {
	HasSimulations(ctx context.Context, userID string) (bool, int, error)
}
type TimelineProvider interface {
	GetEventCount(ctx context.Context, userID string) (int, error)
}
type InsightProvider interface {
	GetInsightSummary(ctx context.Context, userID string) (total, critical int, err error)
}
type NotifProvider interface {
	GetUnreadCount(ctx context.Context, userID string) (int, error)
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
	cacheKey := "adv:" + userID
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached.(*engine.Inputs), nil
	}

	inputs := &engine.Inputs{UserID: userID}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 18)

	wg.Add(1); go func() {
		defer wg.Done()
		nw, err := a.providers.Dashboard.GetNetWorth(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.NetWorth = nw; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		cash, err := a.providers.Dashboard.GetCashBalance(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.CashBalance = cash; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		inc, exp, err := a.providers.Dashboard.GetMonthlyFlow(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.MonthlyIncome = inc; inputs.MonthlyExp = exp; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		onTrack, atRisk, total, gap, err := a.providers.Goals.GetGoalProgress(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.OnTrack = onTrack; inputs.AtRisk = atRisk; inputs.TotalGoals = total; inputs.FundingGap = gap; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		count, balance, err := a.providers.Accounts.GetAccountSummary(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.TotalAccounts = count; inputs.TotalBalance = balance; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		val, ret, retPct, err := a.providers.Portfolio.GetPortfolioSummary(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.PortfolioValue = val; inputs.TotalReturn = ret; inputs.ReturnPct = retPct; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		score, grade, err := a.providers.Health.GetHealthScore(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HealthScore = score; inputs.HealthGrade = grade; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		chg, err := a.providers.Health.GetHealthChange(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HealthChg = chg; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		score, level, err := a.providers.Risk.GetRiskScore(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.RiskScore = score; inputs.RiskLevel = level; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		nwProj, onTrack, conf, err := a.providers.Projection.GetProjectionSummary(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.NWProjected = nwProj; inputs.ProjOnTrack = onTrack; inputs.ProjConf = conf; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, title, err := a.providers.Recs.GetRecommendationSummary(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasRecs = has; inputs.RecCount = cnt; inputs.RecTitle = title; mu.Unlock()
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
		cnt, err := a.providers.Timeline.GetEventCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.EventCount = cnt; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		total, crit, err := a.providers.Insights.GetInsightSummary(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.InsightTotal = total; inputs.InsightCrit = crit; mu.Unlock()
	}()
	wg.Add(1); go func() {
		defer wg.Done()
		cnt, err := a.providers.Notifs.GetUnreadCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.NotifUnread = cnt; mu.Unlock()
	}()

	wg.Wait()
	close(errs)
	for e := range errs { if e != nil { return nil, e } }

	a.cache.Set(ctx, cacheKey, inputs)
	return inputs, nil
}

func (a *Aggregator) Refresh(ctx context.Context, userID string) {
	a.cache.Invalidate(ctx, "adv:"+userID)
}
