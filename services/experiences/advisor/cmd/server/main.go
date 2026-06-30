package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/advisor/internal/aggregator"
	"github.com/horizon/core/services/experiences/advisor/internal/api"
	"github.com/horizon/core/services/experiences/advisor/internal/config"
	"github.com/horizon/core/services/experiences/advisor/internal/engine"
	"github.com/horizon/core/services/experiences/advisor/internal/infrastructure/cache"
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

	log.Printf("Starting Advisor Experience on %s", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultProviders() aggregator.DataProviders {
	return aggregator.DataProviders{
		Dashboard:   &noopDashProvider{},
		Goals:       &noopGoalProvider{},
		Accounts:    &noopAcctProvider{},
		Portfolio:   &noopPfProvider{},
		Planning:    &noopPlanProvider{},
		Health:      &noopHealthProvider{},
		Risk:        &noopRiskProvider{},
		Projection:  &noopProjProvider{},
		Recs:        &noopRecProvider{},
		Optimize:    &noopOptProvider{},
		Simulation:  &noopSimProvider{},
		Timeline:    &noopTimelineProvider{},
		Insights:    &noopInsightProvider{},
		Notifs:      &noopNotifProvider{},
	}
}

type noopDashProvider struct{}
func (n *noopDashProvider) GetNetWorth(ctx context.Context, userID string) (int64, error) { return 5000000, nil }
func (n *noopDashProvider) GetCashBalance(ctx context.Context, userID string) (int64, error) { return 500000, nil }
func (n *noopDashProvider) GetMonthlyFlow(ctx context.Context, userID string) (int64, int64, error) { return 150000, 100000, nil }

type noopGoalProvider struct{}
func (n *noopGoalProvider) GetGoalProgress(ctx context.Context, userID string) (int, int, int, int64, error) { return 3, 1, 5, 500000, nil }

type noopAcctProvider struct{}
func (n *noopAcctProvider) GetAccountSummary(ctx context.Context, userID string) (int, int64, error) { return 8, 1500000, nil }

type noopPfProvider struct{}
func (n *noopPfProvider) GetPortfolioSummary(ctx context.Context, userID string) (int64, float64, float64, error) { return 3500000, 75000, 2.3, nil }

type noopPlanProvider struct{}

type noopHealthProvider struct{}
func (n *noopHealthProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) { return 75, "Good", nil }
func (n *noopHealthProvider) GetHealthChange(ctx context.Context, userID string) (int, error) { return 3, nil }

type noopRiskProvider struct{}
func (n *noopRiskProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) { return 30, "Low", nil }

type noopProjProvider struct{}
func (n *noopProjProvider) GetProjectionSummary(ctx context.Context, userID string) (int64, bool, string, error) { return 8500000, true, "High", nil }

type noopRecProvider struct{}
func (n *noopRecProvider) GetRecommendationSummary(ctx context.Context, userID string) (bool, int, string, error) { return true, 3, "Increase Retirement Contribution", nil }

type noopOptProvider struct{}
func (n *noopOptProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) { return true, 2, nil }

type noopSimProvider struct{}
func (n *noopSimProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) { return true, 1, nil }

type noopTimelineProvider struct{}
func (n *noopTimelineProvider) GetEventCount(ctx context.Context, userID string) (int, error) { return 15, nil }

type noopInsightProvider struct{}
func (n *noopInsightProvider) GetInsightSummary(ctx context.Context, userID string) (int, int, error) { return 8, 0, nil }

type noopNotifProvider struct{}
func (n *noopNotifProvider) GetUnreadCount(ctx context.Context, userID string) (int, error) { return 3, nil }
