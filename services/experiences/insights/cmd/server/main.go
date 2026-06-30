package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/insights/internal/aggregator"
	"github.com/horizon/core/services/experiences/insights/internal/api"
	"github.com/horizon/core/services/experiences/insights/internal/config"
	"github.com/horizon/core/services/experiences/insights/internal/engine"
	"github.com/horizon/core/services/experiences/insights/internal/infrastructure/cache"
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

	log.Printf("Starting Insights Experience on %s", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultProviders() aggregator.DataProviders {
	return aggregator.DataProviders{
		Health:     &noopHealthProvider{},
		Risk:       &noopRiskProvider{},
		Projection: &noopProjProvider{},
		Goals:      &noopGoalProvider{},
		Portfolio:  &noopPfProvider{},
		Recs:       &noopRecProvider{},
		Optimize:   &noopOptProvider{},
		Simulation: &noopSimProvider{},
		Events:     &noopEventProvider{},
		Accounts:   &noopAcctProvider{},
	}
}

type noopHealthProvider struct{}
func (n *noopHealthProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) { return 75, "Good", nil }
func (n *noopHealthProvider) GetHealthChange(ctx context.Context, userID string) (int, error) { return 2, nil }

type noopRiskProvider struct{}
func (n *noopRiskProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) { return 25, "Low", nil }

type noopProjProvider struct{}
func (n *noopProjProvider) GetSavingsRate(ctx context.Context, userID string) (float64, string, error) { return 25.0, "stable", nil }
func (n *noopProjProvider) GetNetWorthChange(ctx context.Context, userID string) (int64, int64, error) { return 50000, 1000000, nil }
func (n *noopProjProvider) GetCashFlowSurplus(ctx context.Context, userID string) (int64, error) { return 15000, nil }

type noopGoalProvider struct{}
func (n *noopGoalProvider) GetGoalProgress(ctx context.Context, userID string) (int, int, float64, error) { return 3, 5, 60.0, nil }
func (n *noopGoalProvider) GetMilestoneCount(ctx context.Context, userID string) (int, error) { return 2, nil }

type noopPfProvider struct{}
func (n *noopPfProvider) GetPortfolioReturn(ctx context.Context, userID string) (int64, float64, error) { return 5000000, 8.5, nil }

type noopRecProvider struct{}
func (n *noopRecProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) { return true, 3, nil }

type noopOptProvider struct{}
func (n *noopOptProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) { return true, 2, nil }

type noopSimProvider struct{}
func (n *noopSimProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) { return true, 1, nil }

type noopEventProvider struct{}
func (n *noopEventProvider) GetAchievementCount(ctx context.Context, userID string) (int, error) { return 5, nil }
func (n *noopEventProvider) GetSpendingAnomaly(ctx context.Context, userID string) (string, error) { return "", nil }

type noopAcctProvider struct{}
func (n *noopAcctProvider) GetAccountCount(ctx context.Context, userID string) (int, error) { return 3, nil }
