package register

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/dashboard/internal/aggregator"
	"github.com/horizon/core/services/experiences/dashboard/internal/api"
	"github.com/horizon/core/services/experiences/dashboard/internal/engine"
	"github.com/horizon/core/services/experiences/dashboard/internal/infrastructure/persistence"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	p := persistence.NewDashboardPGProvider(pool)
	agg := aggregator.New(aggregator.DataProviders{
		Accounts: p, Assets: p, Liabilities: p, Goals: p, Events: p,
		Portfolio: p, Health: p, Risk: p, Recs: p, Projection: p, Simulation: p,
	}, dashNoop{})
	h := api.New(agg, engine.NewComposer())
	h.Register(mux)
}

type dashNoop struct{}

func (dashNoop) Get(_ context.Context, _ string) (*engine.Inputs, bool) { return nil, false }
func (dashNoop) Set(_ context.Context, _ string, _ *engine.Inputs)      {}
