package register

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/insights/internal/aggregator"
	"github.com/horizon/core/services/experiences/insights/internal/api"
	"github.com/horizon/core/services/experiences/insights/internal/engine"
	"github.com/horizon/core/services/experiences/insights/internal/infrastructure/persistence"
	"github.com/horizon/core/services/ai/provider"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool, ai provider.AIProvider) {
	p := persistence.NewInsightsPGProvider(pool)
	agg := aggregator.New(aggregator.DataProviders{
		Health: p, Risk: p, Projection: p, Goals: p, Portfolio: p,
		Recs: p, Optimize: p, Simulation: p, Events: p, Accounts: p,
	}, genNoop{})
	h := api.New(agg, engine.NewComposer(ai))
	h.Register(mux)
}

type genNoop struct{}
func (genNoop) Get(_ context.Context, _ string) (interface{}, bool) { return nil, false }
func (genNoop) Set(_ context.Context, _ string, _ interface{})      {}
func (genNoop) Invalidate(_ context.Context, _ string)              {}
