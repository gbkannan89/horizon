package api

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/engines/scheduled-reports/internal/engine"
)

type Handler struct {
	gen *engine.ReportGenerator
}

func NewHandler(gen *engine.ReportGenerator) *Handler {
	return &Handler{gen: gen}
}

func (h *Handler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		Name     string `json:"name"`
		Schedule string `json:"schedule,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	schedule := engine.ReportSchedule{
		ID:        "on-demand",
		UserID:    req.UserID,
		Name:      req.Name,
		Frequency: engine.FreqMonthly,
		Sections:  h.gen.DefaultSections(),
		Enabled:   true,
	}

	report, err := h.gen.Generate(r.Context(), schedule)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": report})
}

func (h *Handler) DefaultSections(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"sections": h.gen.DefaultSections(),
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
