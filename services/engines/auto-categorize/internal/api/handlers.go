package api

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/engines/auto-categorize/internal/engine"
)

type Handler struct {
	eng *engine.CategorizeEngine
}

func NewHandler(eng *engine.CategorizeEngine) *Handler {
	return &Handler{eng: eng}
}

func (h *Handler) CategorizeTransaction(w http.ResponseWriter, r *http.Request) {
	var req engine.CategorizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.eng.Categorize(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
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
