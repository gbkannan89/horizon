package register

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/smart-alert/internal/api"
	"github.com/horizon/core/services/engines/smart-alert/internal/engine"
	"github.com/horizon/core/services/engines/smart-alert/internal/infrastructure"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	loader := infrastructure.NewPGRuleLoader(pool)
	eng := engine.NewAlertEngine(loader)
	h := api.NewHandler(eng)
	mux.HandleFunc("POST /api/v1/smart-alert/evaluate", h.Evaluate)
}
