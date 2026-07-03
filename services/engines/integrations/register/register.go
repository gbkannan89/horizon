package register

import (
	"net/http"

	"github.com/horizon/core/services/engines/integrations/internal/api"
	"github.com/horizon/core/services/engines/integrations/internal/engine"
)

func RegisterRoutes(mux *http.ServeMux) {
	store := engine.NewIntegrationStore()
	h := api.NewHandler(store)
	mux.HandleFunc("POST /api/v1/integrations", h.Register)
	mux.HandleFunc("GET /api/v1/integrations", h.List)
	mux.HandleFunc("GET /api/v1/integrations/{id}", h.GetByID)
	mux.HandleFunc("DELETE /api/v1/integrations/{id}", h.Delete)
	mux.HandleFunc("POST /api/v1/integrations/{id}/toggle", h.Toggle)
}
