package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/notifications/internal/aggregator"
	"github.com/horizon/core/services/experiences/notifications/internal/api"
	"github.com/horizon/core/services/experiences/notifications/internal/config"
	"github.com/horizon/core/services/experiences/notifications/internal/engine"
	"github.com/horizon/core/services/experiences/notifications/internal/infrastructure/cache"
	"github.com/horizon/core/services/experiences/notifications/internal/infrastructure/persistence"
)

func main() {
	cfg := config.Load()

	pool, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	c := cache.NewInMemory(5 * time.Minute)
	pgProvider := persistence.NewNotifsPGProvider(pool)

	providers := aggregator.DataProviders{
		Health:     pgProvider,
		Risk:       pgProvider,
		Projection: pgProvider,
		Goals:      pgProvider,
		Recs:       pgProvider,
		Optimize:   pgProvider,
		Events:     pgProvider,
		Prefs:      pgProvider,
	}
	agg := aggregator.New(providers, c)
	composer := engine.NewComposer()
	stateRepo := engine.NewStateRepo()
	handlers := api.New(agg, composer, stateRepo)
	router := api.NewRouter(handlers)
	srv := pkghttp.New(cfg.Addr(), router)

	log.Printf("Starting Notification Experience on %s (database: connected)", cfg.Addr())
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
