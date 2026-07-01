package register

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/portfolio/internal/aggregator"
	"github.com/horizon/core/services/experiences/portfolio/internal/api"
	"github.com/horizon/core/services/experiences/portfolio/internal/engine"
	"github.com/horizon/core/services/experiences/portfolio/internal/infrastructure/persistence"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	p := persistence.NewPortfolioPGProvider(pool)
	agg := aggregator.New(aggregator.DataProviders{
		Portfolio: p, Assets: p, Allocation: p, Performance: p, Risk: p,
		Projection: p, Recs: p, Optimize: p, Simulation: p, Events: p,
	}, genNoop{})
	h := api.New(agg, engine.NewComposer())
	h.Register(mux)
}

type genNoop struct{}
func (genNoop) Get(_ context.Context, _ string) (interface{}, bool) { return nil, false }
func (genNoop) Set(_ context.Context, _ string, _ interface{})      {}
func (genNoop) Invalidate(_ context.Context, _ string)              {}
