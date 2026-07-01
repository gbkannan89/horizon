package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/horizon/core/services/internal/auth"

	"github.com/horizon/core/services/experiences/timeline/internal/aggregator"
	"github.com/horizon/core/services/experiences/timeline/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
}

func New(agg *aggregator.Aggregator, comp *engine.Composer) *Handlers {
	return &Handlers{aggregator: agg, composer: comp}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/timeline", h.GetTimeline)
	mux.HandleFunc("GET /api/v1/timeline/recent", h.GetRecent)
	mux.HandleFunc("GET /api/v1/timeline/filter", h.GetFiltered)
	mux.HandleFunc("GET /api/v1/timeline/search", h.Search)
	mux.HandleFunc("GET /api/v1/timeline/{id}", h.GetByID)
}

func getDefaultUserID(r *http.Request) string {
	return auth.UserIDFromRequest(r)
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

func queryInt(qp string, def int) func(r *http.Request) int {
	return func(r *http.Request) int {
		v := r.URL.Query().Get(qp)
		if v == "" { return def }
		n, err := strconv.Atoi(v)
		if err != nil { return def }
		return n
	}
}

var limitFromQuery = queryInt("limit", 50)

func parseCategories(r *http.Request) []engine.EventCategory {
	raw := r.URL.Query().Get("categories")
	if raw == "" { return nil }
	parts := strings.Split(raw, ",")
	var cats []engine.EventCategory
	for _, p := range parts {
		cats = append(cats, engine.EventCategory(strings.TrimSpace(p)))
	}
	return cats
}

func parseFilters(r *http.Request) engine.FilterState {
	return engine.FilterState{
		Categories: parseCategories(r),
		EntityID:   r.URL.Query().Get("entity_id"),
		StartDate:  r.URL.Query().Get("start_date"),
		EndDate:    r.URL.Query().Get("end_date"),
		SearchText: r.URL.Query().Get("q"),
		Severity:   engine.Severity(r.URL.Query().Get("severity")),
		Goal:       r.URL.Query().Get("goal"),
		Account:    r.URL.Query().Get("account"),
		Asset:      r.URL.Query().Get("asset"),
		Liability:  r.URL.Query().Get("liability"),
		Portfolio:  r.URL.Query().Get("portfolio"),
		Rec:        r.URL.Query().Get("recommendation"),
	}
}

func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	view := engine.TimelineView(r.URL.Query().Get("view"))
	if view == "" { view = engine.TVThisMonth }

	items, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}

	output := h.composer.BuildFeed(engine.Inputs{
		Items: items, View: view,
		Filters: parseFilters(r),
		Cursor:  r.URL.Query().Get("cursor"),
		Limit:   limitFromQuery(r),
		UserID:  userID,
	})
	writeJSON(w, http.StatusOK, TimelineResponse{
		Success: true, Data: output,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRecent(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	items, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	recent := h.composer.GetRecent(items, limitFromQuery(r))
	writeJSON(w, http.StatusOK, ItemsResponse{
		Success: true, Data: &ItemList{Items: recent, Count: len(recent)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetFiltered(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	items, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}

	output := h.composer.BuildFeed(engine.Inputs{
		Items: items, View: engine.TVToday,
		Filters: parseFilters(r),
		Limit:   limitFromQuery(r),
		UserID:  userID,
	})
	writeJSON(w, http.StatusOK, TimelineResponse{
		Success: true, Data: output,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "MISSING_QUERY", "search query 'q' is required")
		return
	}

	items, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	results := h.composer.Search(items, q)
	writeJSON(w, http.StatusOK, ItemsResponse{
		Success: true, Data: &ItemList{Items: results, Count: len(results)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "timeline ID is required")
		return
	}

	userID := getDefaultUserID(r)
	items, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	item := h.composer.GetByID(items, id)
	if item == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "timeline item not found")
		return
	}
	writeJSON(w, http.StatusOK, ItemResponse{
		Success: true, Data: item,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}
