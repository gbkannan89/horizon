package register

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/timeline/internal/aggregator"
	"github.com/horizon/core/services/experiences/timeline/internal/api"
	"github.com/horizon/core/services/experiences/timeline/internal/engine"
	"github.com/horizon/core/services/experiences/timeline/internal/infrastructure/persistence"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	pg := persistence.NewTimelinePGProvider(pool)
	agg := aggregator.New(aggregator.DataProviders{
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
	}, tlNoop{})
	h := api.New(agg, engine.NewComposer())
	h.Register(mux)
}

type tlNoop struct{}
func (tlNoop) Get(_ context.Context, _ string) ([]engine.RawItem, bool) { return nil, false }
func (tlNoop) Set(_ context.Context, _ string, _ []engine.RawItem)       {}
func (tlNoop) Invalidate(_ context.Context, _ string)                    {}
