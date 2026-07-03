package register

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/packages/events"
	"github.com/horizon/core/services/experiences/notifications/internal/aggregator"
	"github.com/horizon/core/services/experiences/notifications/internal/api"
	"github.com/horizon/core/services/experiences/notifications/internal/engine"
	"github.com/horizon/core/services/experiences/notifications/internal/infrastructure/persistence"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool, pub events.Publisher) {
	p := persistence.NewNotifsPGProvider(pool)
	agg := aggregator.New(aggregator.DataProviders{
		Health: p, Risk: p, Projection: p, Goals: p, Recs: p,
		Optimize: p, Events: p, Prefs: p,
	}, genNoop{})
	stateRepo := persistence.NewPGStateRepository(pool)
	prefRepo := persistence.NewPGPreferenceRepository(pool)
	h := api.New(agg, engine.NewComposer(), stateRepo, prefRepo, pub)
	h.Register(mux)
}

type genNoop struct{}
func (genNoop) Get(_ context.Context, _ string) (interface{}, bool) { return nil, false }
func (genNoop) Set(_ context.Context, _ string, _ interface{})      {}
func (genNoop) Invalidate(_ context.Context, _ string)              {}
