package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/timeline/internal/aggregator"
	"github.com/horizon/core/services/experiences/timeline/internal/api"
	"github.com/horizon/core/services/experiences/timeline/internal/config"
	"github.com/horizon/core/services/experiences/timeline/internal/engine"
	"github.com/horizon/core/services/experiences/timeline/internal/infrastructure/cache"
	"github.com/horizon/core/services/experiences/timeline/internal/infrastructure/persistence"
)

func main() {
	cfg := config.Load()

	pool, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	c := cache.NewInMemory(5 * time.Minute)
	pg := persistence.NewTimelinePGProvider(pool)

	providers := aggregator.DataProviders{
		Financial:  &persistence.FinEventsProvider{P: pg},
		Goal:       &persistence.GoalEventsProvider{P: pg},
		Account:    &persistence.AcctEventsProvider{P: pg},
		Asset:      &persistence.AssetEventsProvider{P: pg},
		Liability:  &persistence.LiabEventsProvider{P: pg},
		Portfolio:  &persistence.PfEventsProvider{P: pg},
		Health:     &persistence.HealthEventsProvider{P: pg},
		Risk:       &persistence.RiskEventsProvider{P: pg},
		Rec:        &persistence.RecEventsProvider{P: pg},
		Simulation: &persistence.SimEventsProvider{P: pg},
		Optimize:   &persistence.OptEventsProvider{P: pg},
		User:       &persistence.UserEventsProvider{},
		Achieve:    &persistence.AchieveEventsProvider{},
	}
	agg := aggregator.New(providers, c)
	composer := engine.NewComposer()
	handlers := api.New(agg, composer)
	router := api.NewRouter(handlers)
	srv := pkghttp.New(cfg.Addr(), router)

	log.Printf("Starting Timeline Experience on %s (database: connected)", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func connectDB(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
