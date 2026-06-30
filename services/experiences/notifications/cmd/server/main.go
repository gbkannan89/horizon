package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/notifications/internal/aggregator"
	"github.com/horizon/core/services/experiences/notifications/internal/api"
	"github.com/horizon/core/services/experiences/notifications/internal/config"
	"github.com/horizon/core/services/experiences/notifications/internal/engine"
	"github.com/horizon/core/services/experiences/notifications/internal/infrastructure/cache"
)

func main() {
	cfg := config.Load()
	c := cache.NewInMemory(5 * time.Minute)
	providers := defaultProviders()
	agg := aggregator.New(providers, c)
	composer := engine.NewComposer()
	stateRepo := engine.NewStateRepo()
	handlers := api.New(agg, composer, stateRepo)
	router := api.NewRouter(handlers)
	srv := pkghttp.New(cfg.Addr(), router)

	log.Printf("Starting Notification Experience on %s", cfg.Addr())
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
		Recs:       &noopRecProvider{},
		Optimize:   &noopOptProvider{},
		Events:     &noopEventProvider{},
		Prefs:      &noopPrefProvider{},
	}
}

type noopHealthProvider struct{}
func (n *noopHealthProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) { return 75, "Good", nil }
func (n *noopHealthProvider) GetHealthChange(ctx context.Context, userID string) (int, error) { return 2, nil }

type noopRiskProvider struct{}
func (n *noopRiskProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) { return 25, "Low", nil }

type noopProjProvider struct{}
func (n *noopProjProvider) GetCashFlowSurplus(ctx context.Context, userID string) (int64, error) { return 15000, nil }
func (n *noopProjProvider) GetNetWorthChange(ctx context.Context, userID string) (int64, error) { return 50000, nil }

type noopGoalProvider struct{}
func (n *noopGoalProvider) GetGoalProgress(ctx context.Context, userID string) (int, int, int, error) { return 3, 5, 0, nil }

type noopRecProvider struct{}
func (n *noopRecProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) { return true, 3, nil }

type noopOptProvider struct{}
func (n *noopOptProvider) HasOptimizations(ctx context.Context, userID string) (bool, error) { return true, nil }

type noopEventProvider struct{}
func (n *noopEventProvider) GetAchievementCount(ctx context.Context, userID string) (int, error) { return 5, nil }

type noopPrefProvider struct{}
func (n *noopPrefProvider) GetPreferences(ctx context.Context, userID string) ([]engine.Preference, error) {
	return []engine.Preference{}, nil
}
