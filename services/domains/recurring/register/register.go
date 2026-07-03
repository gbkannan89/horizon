package register

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/packages/events"
	apihttp "github.com/horizon/core/services/domains/recurring/internal/api/http"
	"github.com/horizon/core/services/domains/recurring/internal/application"
	"github.com/horizon/core/services/domains/recurring/internal/infrastructure/persistence"
)

type noopPublisher struct{}

func (n *noopPublisher) Publish(ctx context.Context, envelope events.Envelope) error { return nil }

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := persistence.NewPostgresRepository(pool)
	pub := &noopPublisher{}
	svc := application.NewService(repo, pub)
	handlers := apihttp.NewHandlers(svc)
	apihttp.RegisterRoutes(mux, handlers)
}
