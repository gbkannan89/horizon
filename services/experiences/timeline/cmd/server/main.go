package main

import (
	"context"
	"log"
	"time"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/experiences/timeline/internal/aggregator"
	"github.com/horizon/core/services/experiences/timeline/internal/api"
	"github.com/horizon/core/services/experiences/timeline/internal/config"
	"github.com/horizon/core/services/experiences/timeline/internal/engine"
	"github.com/horizon/core/services/experiences/timeline/internal/infrastructure/cache"
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

	log.Printf("Starting Timeline Experience on %s", cfg.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultProviders() aggregator.DataProviders {
	return aggregator.DataProviders{
		Financial:  &noopProvider{name: "financial"},
		Goal:       &noopProvider{name: "goal"},
		Account:    &noopProvider{name: "account"},
		Asset:      &noopProvider{name: "asset"},
		Liability:  &noopProvider{name: "liability"},
		Portfolio:  &noopProvider{name: "portfolio"},
		Health:     &noopProvider{name: "health"},
		Risk:       &noopProvider{name: "risk"},
		Rec:        &noopProvider{name: "rec"},
		Simulation: &noopProvider{name: "simulation"},
		Optimize:   &noopProvider{name: "optimize"},
		User:       &noopProvider{name: "user"},
		Achieve:    &noopProvider{name: "achieve"},
	}
}

type noopProvider struct{ name string }

func (n *noopProvider) GetEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	return nil, nil
}
