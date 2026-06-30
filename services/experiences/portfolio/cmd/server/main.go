package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/portfolio/internal/aggregator"
	"github.com/horizon/core/services/experiences/portfolio/internal/api"
	"github.com/horizon/core/services/experiences/portfolio/internal/config"
	"github.com/horizon/core/services/experiences/portfolio/internal/engine"
	"github.com/horizon/core/services/experiences/portfolio/internal/infrastructure/cache"
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

	log.Printf("Starting Portfolio Experience on %s", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultProviders() aggregator.DataProviders {
	return aggregator.DataProviders{
		Portfolio:  &noopPfProvider{},
		Assets:     &noopAssetProvider{},
		Allocation: &noopAllocProvider{},
		Performance: &noopPerfProvider{},
		Risk:       &noopRiskProvider{},
		Projection: &noopProjProvider{},
		Recs:       &noopRecProvider{},
		Optimize:   &noopOptProvider{},
		Simulation: &noopSimProvider{},
		Events:     &noopEventProvider{},
	}
}

type noopPfProvider struct{}
func (n *noopPfProvider) GetPortfolioSummary(ctx context.Context, userID string) (int64, float64, float64, error) { return 0, 0, 0, nil }

type noopAssetProvider struct{}
func (n *noopAssetProvider) GetAllocations(ctx context.Context, userID string) ([]engine.AllocEntry, error) { return nil, nil }

type noopAllocProvider struct{}

type noopPerfProvider struct{}
func (n *noopPerfProvider) GetPerformance(ctx context.Context, userID string) (float64, float64, float64, int64, int64, error) { return 0, 0, 0, 0, 0, nil }

type noopRiskProvider struct{}
func (n *noopRiskProvider) GetPortfolioRisk(ctx context.Context, userID string) (int, string, float64, float64, float64, float64, error) { return 0, "", 0, 0, 0, 0, nil }

type noopProjProvider struct{}
func (n *noopProjProvider) GetPortfolioProjection(ctx context.Context, userID string) (float64, string, int, float64, error) { return 0, "", 0, 0, nil }

type noopRecProvider struct{}
func (n *noopRecProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) { return false, 0, nil }

type noopOptProvider struct{}
func (n *noopOptProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) { return false, 0, nil }

type noopSimProvider struct{}
func (n *noopSimProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) { return false, 0, nil }

type noopEventProvider struct{}
func (n *noopEventProvider) GetEventCount(ctx context.Context, userID string) (int, error) { return 0, nil }
func (n *noopEventProvider) GetMilestoneCount(ctx context.Context, userID string) (int, error) { return 0, nil }
