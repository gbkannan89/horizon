package api

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/engines/smart-alert/internal/engine"
)

type Handler struct {
	eng *engine.AlertEngine
}

func NewHandler(eng *engine.AlertEngine) *Handler {
	return &Handler{eng: eng}
}

func (h *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var req engine.AlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	alerts, err := h.eng.EvaluateAlerts(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if alerts == nil {
		alerts = []engine.SmartAlert{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": alerts})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
