package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handlers) {
	mux.HandleFunc("POST /api/v1/recurring", h.CreateRecurring)
	mux.HandleFunc("GET /api/v1/recurring", h.ListByUser)
	mux.HandleFunc("GET /api/v1/recurring/{id}", h.GetByID)
	mux.HandleFunc("POST /api/v1/recurring/{id}/activate", h.Activate)
	mux.HandleFunc("POST /api/v1/recurring/{id}/pause", h.Pause)
	mux.HandleFunc("POST /api/v1/recurring/{id}/cancel", h.Cancel)
	mux.HandleFunc("POST /api/v1/recurring/{id}/archive", h.Archive)
}
