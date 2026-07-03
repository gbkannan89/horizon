package api

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/engines/currency/internal/engine"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) ListCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies := engine.GetCurrencies()
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": currencies})
}

func (h *Handler) Convert(w http.ResponseWriter, r *http.Request) {
	var req engine.ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.From == "" || req.To == "" {
		writeError(w, http.StatusBadRequest, "from and to currencies are required")
		return
	}

	result, err := engine.Convert(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": result})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
