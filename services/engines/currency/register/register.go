package register

import (
	"net/http"

	"github.com/horizon/core/services/engines/currency/internal/api"
)

func RegisterRoutes(mux *http.ServeMux) {
	h := api.NewHandler()
	mux.HandleFunc("GET /api/v1/currencies", h.ListCurrencies)
	mux.HandleFunc("POST /api/v1/currencies/convert", h.Convert)
}
