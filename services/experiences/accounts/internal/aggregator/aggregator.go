package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/accounts/internal/engine"
)

type DataProviders struct {
	Accounts   AccountProvider
	Trans      TransProvider
	Projection ProjProvider
	Risk       RiskProvider
	Health     HealthProvider
	Recs       RecProvider
	Optimize   OptProvider
	Simulation SimProvider
	Events     EventProvider
}

type AccountProvider interface {
	GetAccounts(ctx context.Context, userID string) ([]engine.AccountInput, error)
	GetAccountByID(ctx context.Context, userID, accountID string) (*engine.AccountInput, error)
}
type TransProvider interface {
	GetTransactions(ctx context.Context, accountID string, limit int) ([]engine.TransInput, error)
}
type ProjProvider interface {
	GetCashFlowProjection(ctx context.Context, userID string) (inflow, outflow, projected int64, err error)
}
type RiskProvider interface {
	GetRiskScore(ctx context.Context, userID string) (int, string, error)
}
type HealthProvider interface {
	GetHealthScore(ctx context.Context, userID string) (int, string, error)
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
	cacheKey := "accts:" + userID
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached.(*engine.Inputs), nil
	}

	inputs := &engine.Inputs{UserID: userID}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 12)

	wg.Add(1); go func() {
		defer wg.Done()
		accts, err := a.providers.Accounts.GetAccounts(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.Accounts = accts; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		inflow, outflow, projected, err := a.providers.Projection.GetCashFlowProjection(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.TotalInflow = inflow; inputs.TotalOutflow = outflow; inputs.ProjectedFlow = projected; mu.Unlock()
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
		mu.Lock(); inputs.HasOpts = has; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		has, cnt, err := a.providers.Simulation.HasSimulations(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.HasSims = has; mu.Unlock()
	}()

	wg.Add(1); go func() {
		defer wg.Done()
		cnt, err := a.providers.Events.GetEventCount(ctx, userID)
		if err != nil { errs <- err; return }
		mu.Lock(); inputs.EventCount = cnt; mu.Unlock()
	}()

	wg.Wait()
	close(errs)
	for e := range errs { if e != nil { return nil, e } }

	// Compute cash balance from accounts
	var cashBalance int64
	for _, a := range inputs.Accounts {
		cashBalance += a.CurrentBalance
	}
	inputs.CashBalance = cashBalance

	a.cache.Set(ctx, cacheKey, inputs)
	return inputs, nil
}

func (a *Aggregator) GetAccountDetail(ctx context.Context, userID, accountID string) (*engine.Inputs, error) {
	inputs, err := a.Aggregate(ctx, userID)
	if err != nil { return nil, err }

	// Find the specific account
	for _, a := range inputs.Accounts {
		if a.AccountID == accountID {
			inputs.Account = &a
			break
		}
	}

	// Fetch transactions
	txns, err := a.providers.Trans.GetTransactions(ctx, accountID, 50)
	if err == nil {
		inputs.Transactions = txns
	}

	return inputs, nil
}

func (a *Aggregator) Refresh(ctx context.Context, userID string) {
	a.cache.Invalidate(ctx, "accts:"+userID)
}
