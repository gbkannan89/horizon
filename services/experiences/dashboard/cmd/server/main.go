package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/dashboard/internal/aggregator"
	"github.com/horizon/core/services/experiences/dashboard/internal/api"
	"github.com/horizon/core/services/experiences/dashboard/internal/config"
	"github.com/horizon/core/services/experiences/dashboard/internal/engine"
	"github.com/horizon/core/services/experiences/dashboard/internal/infrastructure/cache"
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

	log.Printf("Starting Dashboard Experience on %s", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultProviders() aggregator.DataProviders {
	return aggregator.DataProviders{
		Accounts:    &defaultAccountProvider{},
		Assets:      &defaultAssetProvider{},
		Liabilities: &defaultLiabilityProvider{},
		Goals:       &defaultGoalProvider{},
		Events:      &defaultEventProvider{},
		Portfolio:   &defaultPortfolioProvider{},
		Health:      &defaultHealthProvider{},
		Risk:        &defaultRiskProvider{},
		Recs:        &defaultRecProvider{},
		Projection:  &defaultProjectionProvider{},
		Simulation:  &defaultSimulationProvider{},
	}
}

type defaultAccountProvider struct{}
func (d *defaultAccountProvider) GetCashBalance(ctx context.Context, userID string) (int64, error) { return 0, nil }
func (d *defaultAccountProvider) HasAccounts(ctx context.Context, userID string) (bool, error) { return false, nil }

type defaultAssetProvider struct{}
func (d *defaultAssetProvider) GetTotalAssets(ctx context.Context, userID string) (int64, error) { return 0, nil }

type defaultLiabilityProvider struct{}
func (d *defaultLiabilityProvider) GetTotalLiabilities(ctx context.Context, userID string) (int64, error) { return 0, nil }

type defaultGoalProvider struct{}
func (d *defaultGoalProvider) GetGoalCounts(ctx context.Context, userID string) (onTrack, atRisk, total int, hasGoals bool, err error) { return 0, 0, 0, false, nil }

type defaultEventProvider struct{}
func (d *defaultEventProvider) GetRecentCount(ctx context.Context, userID string) (int, bool, error) { return 0, false, nil }

type defaultPortfolioProvider struct{}
func (d *defaultPortfolioProvider) GetPortfolioValue(ctx context.Context, userID string) (int64, error) { return 0, nil }

type defaultHealthProvider struct{}
func (d *defaultHealthProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) { return 75, "Good", nil }

type defaultRiskProvider struct{}
func (d *defaultRiskProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) { return 25, "Low", nil }

type defaultRecProvider struct{}
func (d *defaultRecProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) { return false, 0, nil }

type defaultProjectionProvider struct{}
func (d *defaultProjectionProvider) GetMonthlyIncomeExpenses(ctx context.Context, userID string) (int64, int64, error) { return 0, 0, nil }

type defaultSimulationProvider struct{}
func (d *defaultSimulationProvider) HasSimulations(ctx context.Context, userID string) (bool, error) { return false, nil }
