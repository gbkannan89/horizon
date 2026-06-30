package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/timeline/internal/engine"
)

// DataProviders holds all interfaces for fetching timeline events.
type DataProviders struct {
	Financial  EventProvider
	Goal       EventProvider
	Account    EventProvider
	Asset      EventProvider
	Liability  EventProvider
	Portfolio  EventProvider
	Health     EventProvider
	Risk       EventProvider
	Rec        EventProvider
	Simulation EventProvider
	Optimize   EventProvider
	User       EventProvider
	Achieve    EventProvider
}

type EventProvider interface {
	GetEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error)
}

// Aggregator loads timeline data from all providers in parallel.
type Aggregator struct {
	providers DataProviders
	cache     Cache
}

type Cache interface {
	Get(ctx context.Context, key string) ([]engine.RawItem, bool)
	Set(ctx context.Context, key string, items []engine.RawItem)
	Invalidate(ctx context.Context, key string)
}

func New(providers DataProviders, cache Cache) *Aggregator {
	return &Aggregator{providers: providers, cache: cache}
}

// Aggregate fetches all timeline events from all providers in parallel.
func (a *Aggregator) Aggregate(ctx context.Context, userID string) ([]engine.RawItem, error) {
	cacheKey := "timeline:" + userID
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached, nil
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make(chan error, 14)
	allItems := make([]engine.RawItem, 0, 200)

	providers := []struct {
		name     string
		provider EventProvider
	}{
		{"financial", a.providers.Financial},
		{"goal", a.providers.Goal},
		{"account", a.providers.Account},
		{"asset", a.providers.Asset},
		{"liability", a.providers.Liability},
		{"portfolio", a.providers.Portfolio},
		{"health", a.providers.Health},
		{"risk", a.providers.Risk},
		{"rec", a.providers.Rec},
		{"simulation", a.providers.Simulation},
		{"optimize", a.providers.Optimize},
		{"user", a.providers.User},
		{"achieve", a.providers.Achieve},
	}

	for _, p := range providers {
		wg.Add(1)
		go func(pr EventProvider) {
			defer wg.Done()
			items, err := pr.GetEvents(ctx, userID, 100)
			if err != nil {
				errs <- err
				return
			}
			mu.Lock()
			allItems = append(allItems, items...)
			mu.Unlock()
		}(p.provider)
	}

	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil { return nil, e }
	}

	a.cache.Set(ctx, cacheKey, allItems)
	return allItems, nil
}

// Refresh forces a fresh load by clearing cache.
func (a *Aggregator) Refresh(ctx context.Context, userID string) ([]engine.RawItem, error) {
	a.cache.Invalidate(ctx, "timeline:"+userID)
	return a.Aggregate(ctx, userID)
}
