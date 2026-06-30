package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/goals/internal/aggregator"
	"github.com/horizon/core/services/experiences/goals/internal/api"
	"github.com/horizon/core/services/experiences/goals/internal/config"
	"github.com/horizon/core/services/experiences/goals/internal/engine"
	"github.com/horizon/core/services/experiences/goals/internal/infrastructure/cache"
)

func main() {
	cfg := config.Load()
	c := cache.NewInMemory(5 * time.Minute)
	providers := defaultProviders()
	agg := aggregator.New(providers, c)
	composer := engine.NewComposer()
	handlers := api.New(agg, composer)
	router := api.NewRouter(handlers)
	srv := pkghttp.New(cfg.Addr(), router)

	log.Printf("Starting Goals Experience on %s", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultProviders() aggregator.DataProviders {
	return aggregator.DataProviders{
		Goals:      &noopGoalProvider{},
		Allocation: &noopAllocationProvider{},
		Projection: &noopProjectionProvider{},
		Recs:       &noopRecProvider{},
		Optimize:   &noopOptProvider{},
		Simulation: &noopSimProvider{},
		Risk:       &noopRiskProvider{},
		Health:     &noopHealthProvider{},
		Events:     &noopEventProvider{},
	}
}

type noopGoalProvider struct{}
func (n *noopGoalProvider) GetGoals(ctx context.Context, userID string) ([]engine.GoalData, error) { return nil, nil }
func (n *noopGoalProvider) GetGoalByID(ctx context.Context, userID, goalID string) (*engine.GoalData, error) { return nil, nil }

type noopAllocationProvider struct{}
func (n *noopAllocationProvider) GetMonthlyContribution(ctx context.Context, goalID string) (float64, error) { return 0, nil }

type noopProjectionProvider struct{}
func (n *noopProjectionProvider) GetGoalProjection(ctx context.Context, goalID string) (string, float64, bool, error) { return "", 0, false, nil }

type noopRecProvider struct{}
func (n *noopRecProvider) GetGoalRecommendations(ctx context.Context, goalID string) ([]engine.RecItem, error) { return nil, nil }

type noopOptProvider struct{}
func (n *noopOptProvider) GetGoalOptimizations(ctx context.Context, goalID string) ([]engine.OptItem, error) { return nil, nil }

type noopSimProvider struct{}
func (n *noopSimProvider) HasSimulations(ctx context.Context, goalID string) (bool, error) { return false, nil }

type noopRiskProvider struct{}
func (n *noopRiskProvider) HasRisk(ctx context.Context, goalID string) (bool, error) { return false, nil }

type noopHealthProvider struct{}

type noopEventProvider struct{}
func (n *noopEventProvider) GetGoalEvents(ctx context.Context, goalID string) ([]engine.GoalEvent, error) { return nil, nil }
func (n *noopEventProvider) GetGoalMilestones(ctx context.Context, goalID string) ([]engine.Milestone, error) { return nil, nil }
