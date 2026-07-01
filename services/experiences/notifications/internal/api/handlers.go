package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/horizon/core/services/internal/auth"
	"strconv"
	"sync"
	"time"

	"github.com/horizon/core/services/experiences/notifications/internal/aggregator"
	"github.com/horizon/core/services/experiences/notifications/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
	stateRepo  *engine.StateRepository
	mu         sync.RWMutex
	prefs      map[string][]engine.Preference
}

func New(agg *aggregator.Aggregator, comp *engine.Composer, repo *engine.StateRepository) *Handlers {
	return &Handlers{aggregator: agg, composer: comp, stateRepo: repo, prefs: make(map[string][]engine.Preference)}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/notifications", h.GetNotifications)
	mux.HandleFunc("GET /api/v1/notifications/unread", h.GetUnread)
	mux.HandleFunc("GET /api/v1/notifications/history", h.GetHistory)
	mux.HandleFunc("GET /api/v1/notifications/preferences", h.GetPreferences)
	mux.HandleFunc("PUT /api/v1/notifications/preferences", h.UpdatePreferences)
	mux.HandleFunc("GET /api/v1/notifications/search", h.Search)
	mux.HandleFunc("GET /api/v1/notifications/{id}", h.GetByID)
	mux.HandleFunc("POST /api/v1/notifications/{id}/read", h.MarkRead)
	mux.HandleFunc("POST /api/v1/notifications/{id}/archive", h.Archive)
	mux.HandleFunc("POST /api/v1/notifications/{id}/snooze", h.Snooze)
}

func getDefaultUserID(r *http.Request) string {
	return auth.UserIDFromRequest(r)
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" { return def }
	n, err := strconv.Atoi(v)
	if err != nil { return def }
	return n
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false, "error": map[string]string{"code": code, "message": message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) loadPrefs(userID string) map[engine.Category]engine.Preference {
	h.mu.RLock()
	prefs, ok := h.prefs[userID]
	h.mu.RUnlock()
	if !ok { return map[engine.Category]engine.Preference{} }
	m := make(map[engine.Category]engine.Preference)
	for _, p := range prefs { m[p.Category] = p }
	return m
}

func (h *Handlers) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	center := h.composer.BuildCenter(*inputs, h.loadPrefs(userID), h.stateRepo, queryInt(r, "limit", 50), r.URL.Query().Get("cursor"))
	writeJSON(w, http.StatusOK, CenterResponse{Success: true, Data: center, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetUnread(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	center := h.composer.BuildUnread(*inputs, h.loadPrefs(userID), h.stateRepo)
	writeJSON(w, http.StatusOK, CenterResponse{Success: true, Data: center, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	center := h.composer.BuildHistory(*inputs, h.loadPrefs(userID), h.stateRepo, queryInt(r, "limit", 50), r.URL.Query().Get("cursor"))
	writeJSON(w, http.StatusOK, CenterResponse{Success: true, Data: center, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	h.mu.RLock()
	prefs, ok := h.prefs[userID]
	h.mu.RUnlock()
	if !ok {
		inputs, err := h.aggregator.Aggregate(r.Context(), userID)
		if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
		prefs = inputs.Preferences
	}
	writeJSON(w, http.StatusOK, PrefListResponse{Success: true, Data: &PrefList{Preferences: prefs, Count: len(prefs)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	var req []engine.Preference
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid preference data")
		return
	}
	userID := getDefaultUserID(r)
	h.mu.Lock()
	h.prefs[userID] = req
	h.mu.Unlock()
	writeJSON(w, http.StatusOK, PrefListResponse{Success: true, Data: &PrefList{Preferences: req, Count: len(req)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" { writeError(w, http.StatusBadRequest, "MISSING_QUERY", "search query 'q' is required"); return }
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	results := h.composer.Search(*inputs, h.loadPrefs(userID), h.stateRepo, q)
	writeJSON(w, http.StatusOK, NotifListResponse{Success: true, Data: &NotifList{Notifications: results, Count: len(results)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "notification ID required"); return }
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	notif := h.composer.GetByID(*inputs, h.loadPrefs(userID), h.stateRepo, id)
	if notif == nil { writeError(w, http.StatusNotFound, "NOT_FOUND", "notification not found"); return }
	writeJSON(w, http.StatusOK, NotifResponse{Success: true, Data: notif, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) MarkRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "notification ID required"); return }
	h.stateRepo.SetState(id, engine.StateRead)
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "notif_id": id, "state": engine.StateRead})
}

func (h *Handlers) Archive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "notification ID required"); return }
	h.stateRepo.SetState(id, engine.StateArchived)
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "notif_id": id, "state": engine.StateArchived})
}

func (h *Handlers) Snooze(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "notification ID required"); return }
	var req struct {
		Until string `json:"until"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Until = time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	}
	h.stateRepo.Snooze(id, req.Until)
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "notif_id": id, "state": engine.StateSnoozed, "until": req.Until})
}
