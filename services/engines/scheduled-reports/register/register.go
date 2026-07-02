package register

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/scheduled-reports/internal/api"
	"github.com/horizon/core/services/engines/scheduled-reports/internal/engine"
	"github.com/horizon/core/services/engines/scheduled-reports/internal/infrastructure"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	provider := infrastructure.NewPGDataProvider(pool)
	gen := engine.NewReportGenerator(provider)
	h := api.NewHandler(gen)
	mux.HandleFunc("POST /api/v1/reports/generate", h.GenerateReport)
	mux.HandleFunc("GET /api/v1/reports/sections", h.DefaultSections)
}
