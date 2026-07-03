package api

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/engines/integrations/internal/engine"
)

type Handler struct {
	store *engine.IntegrationStore
}

func NewHandler(store *engine.IntegrationStore) *Handler { return &Handler{store: store} }

func userID(r *http.Request) string {
	if u := r.URL.Query().Get("user_id"); u != "" { return u }
	return "default"
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string `json:"name"`
		Provider   string `json:"provider"`
		APIKey     string `json:"api_key,omitempty"`
		WebhookURL string `json:"webhook_url,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	integ, err := h.store.Register(userID(r), req.Name, req.Provider, req.APIKey, req.WebhookURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": integ})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	integrations, err := h.store.ListByUser(userID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if integrations == nil { integrations = []*engine.Integration{} }
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": integrations})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	integ, err := h.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": integ})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) Toggle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	integ, err := h.store.Toggle(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": integ})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
