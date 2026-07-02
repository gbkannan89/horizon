package register

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	apihttp "github.com/horizon/core/services/domains/rules/internal/api/http"
	"github.com/horizon/core/services/domains/rules/internal/application"
	"github.com/horizon/core/services/domains/rules/internal/infrastructure/persistence"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := persistence.NewPostgresRuleRepository(pool)
	svc := application.NewRuleService(repo, time.Now)
	handlers := apihttp.NewHandlers(svc)
	apihttp.RegisterRoutes(mux, handlers)
}

var _ = context.Background
