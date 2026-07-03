package register

import (
	"net/http"

	"github.com/horizon/core/services/engines/tax-planning/internal/api"
)

func RegisterRoutes(mux *http.ServeMux) {
	h := api.NewHandler()
	mux.HandleFunc("POST /api/v1/tax/analyze", h.Analyze)
}
