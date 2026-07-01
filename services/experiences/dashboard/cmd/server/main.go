package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/dashboard/internal/aggregator"
	"github.com/horizon/core/services/experiences/dashboard/internal/api"
	"github.com/horizon/core/services/experiences/dashboard/internal/config"
	"github.com/horizon/core/services/experiences/dashboard/internal/engine"
	"github.com/horizon/core/services/experiences/dashboard/internal/infrastructure/cache"
	"github.com/horizon/core/services/experiences/dashboard/internal/infrastructure/persistence"
)

func main() {
	cfg := config.Load()

	pool, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	c := cache.NewInMemory(5 * time.Minute)
	pgProvider := persistence.NewDashboardPGProvider(pool)
	providers := aggregator.DataProviders{
		Accounts:    pgProvider,
		Assets:      pgProvider,
		Liabilities: pgProvider,
		Goals:       pgProvider,
		Events:      pgProvider,
		Portfolio:   pgProvider,
		Health:      pgProvider,
		Risk:        pgProvider,
		Recs:        pgProvider,
		Projection:  pgProvider,
		Simulation:  pgProvider,
	}
	agg := aggregator.New(providers, c)
	composer := engine.NewComposer()
	handlers := api.New(agg, composer)
	router := api.NewRouter(handlers)
	srv := pkghttp.New(cfg.Addr(), router)

	log.Printf("Starting Dashboard Experience on %s (database: connected)", cfg.Addr())
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
