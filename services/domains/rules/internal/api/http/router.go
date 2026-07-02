package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handlers) {
	mux.HandleFunc("POST /api/v1/rules", h.CreateRule)
	mux.HandleFunc("GET /api/v1/rules/{id}", h.GetRule)
	mux.HandleFunc("GET /api/v1/rules", h.ListRules)
	mux.HandleFunc("DELETE /api/v1/rules/{id}", h.DeleteRule)
}
