package http

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/domains/rules/internal/application"
	"github.com/horizon/core/services/domains/rules/internal/application/dto/command"
	"github.com/horizon/core/services/domains/rules/internal/application/dto/query"
	"github.com/horizon/core/services/domains/rules/internal/domain"
	"github.com/horizon/core/services/internal/auth"
)

type Handlers struct {
	svc *application.RuleService
}

func NewHandlers(svc *application.RuleService) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	conds := make([]domain.Condition, len(req.Conditions))
	for i, c := range req.Conditions {
		conds[i] = domain.Condition{Field: c.Field, Operator: domain.Operator(c.Operator), Value: c.Value}
	}
	acts := make([]domain.Action, len(req.Actions))
	for i, a := range req.Actions {
		acts[i] = domain.Action{Type: domain.ActionType(a.Type), Params: a.Params}
	}
	cmd := command.CreateRuleCommand{
		UserID: auth.UserIDFromRequest(r), Name: req.Name,
		Description: req.Description, Category: req.Category,
		Priority: req.Priority, Enabled: req.Enabled,
		Conditions: conds, Actions: acts,
	}
	result, err := h.svc.Create(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": result})
}

func (h *Handlers) GetRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	view, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": view})
}

func (h *Handlers) ListRules(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	rules, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rules == nil { rules = []query.RuleView{} }
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": rules})
}

func (h *Handlers) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
