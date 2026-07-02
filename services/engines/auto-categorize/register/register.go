package register

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/auto-categorize/internal/api"
	"github.com/horizon/core/services/engines/auto-categorize/internal/engine"
	"github.com/horizon/core/services/engines/auto-categorize/internal/infrastructure"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	loader := infrastructure.NewPGRuleLoader(pool)
	eng := engine.NewCategorizeEngine(loader)
	h := api.NewHandler(eng)
	mux.HandleFunc("POST /api/v1/auto-categorize", h.CategorizeTransaction)
}
