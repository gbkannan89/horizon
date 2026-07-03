package api

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/engines/tax-planning/internal/engine"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Analyze(w http.ResponseWriter, r *http.Request) {
	var input engine.TaxInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result := input.Analyze()
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
