package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/planning/internal/aggregator"
	"github.com/horizon/core/services/experiences/planning/internal/api"
	"github.com/horizon/core/services/experiences/planning/internal/config"
	"github.com/horizon/core/services/experiences/planning/internal/engine"
	"github.com/horizon/core/services/experiences/planning/internal/infrastructure/cache"
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

	log.Printf("Starting Planning Experience on %s", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultProviders() aggregator.DataProviders {
	return aggregator.DataProviders{
		Goals:       &noopGoalProvider{},
		Accounts:    &noopAccountProvider{},
		Projection:  &noopProjProvider{},
		Risk:        &noopRiskProvider{},
		Health:      &noopHealthProvider{},
		Recs:        &noopRecProvider{},
		Optimize:    &noopOptProvider{},
		Simulation:  &noopSimProvider{},
		Events:      &noopEventProvider{},
		Scenarios:   &noopScenarioProvider{},
	}
}

type noopGoalProvider struct{}
func (n *noopGoalProvider) GetGoalSummary(ctx context.Context, userID string) (int, int, int64, error) { return 0, 0, 0, nil }

type noopAccountProvider struct{}
func (n *noopAccountProvider) GetCashBalance(ctx context.Context, userID string) (int64, int64, int64, error) { return 0, 0, 0, nil }

type noopProjProvider struct{}
func (n *noopProjProvider) GetPlanningProjection(ctx context.Context, userID string) (int64, int64, int64, int64, string, error) { return 0, 0, 0, 0, "", nil }

type noopRiskProvider struct{}
func (n *noopRiskProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) { return 25, "Low", nil }

type noopHealthProvider struct{}
func (n *noopHealthProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) { return 75, "Good", nil }

type noopRecProvider struct{}
func (n *noopRecProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) { return false, 0, nil }

type noopOptProvider struct{}
func (n *noopOptProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) { return false, 0, nil }

type noopSimProvider struct{}
func (n *noopSimProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) { return false, 0, nil }

type noopEventProvider struct{}
func (n *noopEventProvider) GetEventCount(ctx context.Context, userID string) (int, int, error) { return 0, 0, nil }

type noopScenarioProvider struct{}
func (n *noopScenarioProvider) GetScenarios(ctx context.Context, userID string) ([]engine.Scenario, error) { return nil, nil }
