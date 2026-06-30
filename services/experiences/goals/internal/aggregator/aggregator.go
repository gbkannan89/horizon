package aggregator

import (
	"context"
	"sync"

	"github.com/horizon/core/services/experiences/goals/internal/engine"
)

// DataProviders holds interfaces for fetching goal-related data.
type DataProviders struct {
	Goals       GoalProvider
	Allocation  AllocationProvider
	Projection  ProjectionProvider
	Recs        RecProvider
	Optimize    OptProvider
	Simulation  SimProvider
	Risk        RiskProvider
	Health      HealthProvider
	Events      EventProvider
}

type GoalProvider interface {
	GetGoals(ctx context.Context, userID string) ([]engine.GoalData, error)
	GetGoalByID(ctx context.Context, userID, goalID string) (*engine.GoalData, error)
}

type AllocationProvider interface {
	GetMonthlyContribution(ctx context.Context, goalID string) (float64, error)
}

type ProjectionProvider interface {
	GetGoalProjection(ctx context.Context, goalID string) (projectedDate string, projectedValue float64, onTrack bool, err error)
}

type RecProvider interface {
	GetGoalRecommendations(ctx context.Context, goalID string) ([]engine.RecItem, error)
}

type OptProvider interface {
	GetGoalOptimizations(ctx context.Context, goalID string) ([]engine.OptItem, error)
}

type SimProvider interface {
	HasSimulations(ctx context.Context, goalID string) (bool, error)
}

type RiskProvider interface {
	HasRisk(ctx context.Context, goalID string) (bool, error)
}

type HealthProvider interface{}

type EventProvider interface {
	GetGoalEvents(ctx context.Context, goalID string) ([]engine.GoalEvent, error)
	GetGoalMilestones(ctx context.Context, goalID string) ([]engine.Milestone, error)
}

// Aggregator loads goal data from all providers in parallel.
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

func (a *Aggregator) GetGoals(ctx context.Context, userID string) ([]engine.GoalData, error) {
	cacheKey := "goals:list:" + userID
	if cached, ok := a.cache.Get(ctx, cacheKey); ok {
		return cached.([]engine.GoalData), nil
	}

	goals, err := a.providers.Goals.GetGoals(ctx, userID)
	if err != nil { return nil, err }

	// Enrich each goal with data from other providers
	var wg sync.WaitGroup
	for i := range goals {
		wg.Add(1)
		go func(g *engine.GoalData) {
			defer wg.Done()
			var mu sync.Mutex
			var inner sync.WaitGroup

			inner.Add(1)
			go func() {
				defer inner.Done()
				if amt, err := a.providers.Allocation.GetMonthlyContribution(ctx, g.GoalID); err == nil {
					mu.Lock(); g.MonthlyContribution = amt; mu.Unlock()
				}
			}()

			inner.Add(1)
			go func() {
				defer inner.Done()
				if date, val, onTrack, err := a.providers.Projection.GetGoalProjection(ctx, g.GoalID); err == nil {
					mu.Lock(); g.ProjectedDate = date; g.ProjectedValue = val; g.OnTrack = onTrack; mu.Unlock()
				}
			}()

			inner.Add(1)
			go func() {
				defer inner.Done()
				if recs, err := a.providers.Recs.GetGoalRecommendations(ctx, g.GoalID); err == nil {
					mu.Lock(); g.HasRecommendation = len(recs) > 0; g.RecCount = len(recs); mu.Unlock()
				}
			}()

			inner.Add(1)
			go func() {
				defer inner.Done()
				if opts, err := a.providers.Optimize.GetGoalOptimizations(ctx, g.GoalID); err == nil {
					mu.Lock(); g.HasOptimization = len(opts) > 0; g.OptCount = len(opts); mu.Unlock()
				}
			}()

			inner.Add(1)
			go func() {
				defer inner.Done()
				if has, err := a.providers.Risk.HasRisk(ctx, g.GoalID); err == nil {
					mu.Lock(); g.HasRisk = has; mu.Unlock()
				}
			}()

			inner.Add(1)
			go func() {
				defer inner.Done()
				if evts, err := a.providers.Events.GetGoalEvents(ctx, g.GoalID); err == nil {
					mu.Lock(); g.EventCount = len(evts); mu.Unlock()
				}
			}()

			inner.Add(1)
			go func() {
				defer inner.Done()
				if mss, err := a.providers.Events.GetGoalMilestones(ctx, g.GoalID); err == nil {
					mu.Lock(); g.MilestoneCount = len(mss); mu.Unlock()
				}
			}()

			inner.Wait()
		}(&goals[i])
	}
	wg.Wait()

	a.cache.Set(ctx, cacheKey, goals)
	return goals, nil
}

// GetGoalDetail returns a single enriched goal.
func (a *Aggregator) GetGoalDetail(ctx context.Context, userID, goalID string) (*engine.GoalData, error) {
	goal, err := a.providers.Goals.GetGoalByID(ctx, userID, goalID)
	if err != nil { return nil, err }

	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if amt, err := a.providers.Allocation.GetMonthlyContribution(ctx, goalID); err == nil {
			mu.Lock(); goal.MonthlyContribution = amt; mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if date, val, onTrack, err := a.providers.Projection.GetGoalProjection(ctx, goalID); err == nil {
			mu.Lock(); goal.ProjectedDate = date; goal.ProjectedValue = val; goal.OnTrack = onTrack; mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if recs, err := a.providers.Recs.GetGoalRecommendations(ctx, goalID); err == nil {
			mu.Lock(); goal.HasRecommendation = len(recs) > 0; goal.RecCount = len(recs); mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if opts, err := a.providers.Optimize.GetGoalOptimizations(ctx, goalID); err == nil {
			mu.Lock(); goal.HasOptimization = len(opts) > 0; goal.OptCount = len(opts); mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if has, err := a.providers.Risk.HasRisk(ctx, goalID); err == nil {
			mu.Lock(); goal.HasRisk = has; mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if evts, err := a.providers.Events.GetGoalEvents(ctx, goalID); err == nil {
			mu.Lock(); goal.EventCount = len(evts); mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if mss, err := a.providers.Events.GetGoalMilestones(ctx, goalID); err == nil {
			mu.Lock(); goal.MilestoneCount = len(mss); mu.Unlock()
		}
	}()

	wg.Wait()
	return goal, nil
}

// GetGoalData returns the raw goal data for a single goal (used by detail endpoints).
func (a *Aggregator) GetGoalData(ctx context.Context, userID, goalID string) (*engine.GoalData, error) {
	return a.GetGoalDetail(ctx, userID, goalID)
}

// GetRecs returns recommendations for a goal.
func (a *Aggregator) GetRecs(ctx context.Context, goalID string) ([]engine.RecItem, error) {
	return a.providers.Recs.GetGoalRecommendations(ctx, goalID)
}

// GetOpts returns optimizations for a goal.
func (a *Aggregator) GetOpts(ctx context.Context, goalID string) ([]engine.OptItem, error) {
	return a.providers.Optimize.GetGoalOptimizations(ctx, goalID)
}

// GetEvents returns timeline events for a goal.
func (a *Aggregator) GetEvents(ctx context.Context, goalID string) ([]engine.GoalEvent, error) {
	return a.providers.Events.GetGoalEvents(ctx, goalID)
}

// GetMilestones returns milestones for a goal.
func (a *Aggregator) GetMilestones(ctx context.Context, goalID string) ([]engine.Milestone, error) {
	return a.providers.Events.GetGoalMilestones(ctx, goalID)
}

// Refresh clears cache for a user.
func (a *Aggregator) Refresh(ctx context.Context, userID string) {
	a.cache.Invalidate(ctx, "goals:list:"+userID)
}
