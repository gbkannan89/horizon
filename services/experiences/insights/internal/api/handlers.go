package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/horizon/core/services/experiences/insights/internal/aggregator"
	"github.com/horizon/core/services/experiences/insights/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
}

func New(agg *aggregator.Aggregator, comp *engine.Composer) *Handlers {
	return &Handlers{aggregator: agg, composer: comp}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/insights", h.GetInsights)
	mux.HandleFunc("GET /api/v1/insights/summary", h.GetSummary)
	mux.HandleFunc("GET /api/v1/insights/opportunities", h.GetOpportunities)
	mux.HandleFunc("GET /api/v1/insights/warnings", h.GetWarnings)
	mux.HandleFunc("GET /api/v1/insights/achievements", h.GetAchievements)
	mux.HandleFunc("GET /api/v1/insights/forecast", h.GetForecast)
	mux.HandleFunc("GET /api/v1/insights/trends", h.GetTrends)
	mux.HandleFunc("GET /api/v1/insights/timeline", h.GetTimeline)
	mux.HandleFunc("GET /api/v1/insights/search", h.Search)
	mux.HandleFunc("GET /api/v1/insights/{id}", h.GetByID)
}

func getDefaultUserID(r *http.Request) string {
	u := r.URL.Query().Get("user_id")
	if u != "" { return u }
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		if uid := strings.TrimPrefix(auth, "Bearer "); uid != "" { return uid }
	}
	return "default"
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

func (h *Handlers) GetInsights(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	limit := queryInt(r, "limit", 20)
	feed := h.composer.BuildFeed(*inputs, limit, r.URL.Query().Get("cursor"))
	writeJSON(w, http.StatusOK, FeedResponse{Success: true, Data: feed, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	dash := h.composer.BuildDashboard(*inputs)
	writeJSON(w, http.StatusOK, DashboardResponse{Success: true, Data: dash, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetOpportunities(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	all := h.composer.FilterByCategory(*inputs, engine.ICOpportunity)
	writeJSON(w, http.StatusOK, InsightListResponse{Success: true, Data: &InsightList{Insights: all, Count: len(all)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetWarnings(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	all := h.composer.FilterByCategory(*inputs, engine.ICWarning)
	all2 := h.composer.FilterByCategory(*inputs, engine.ICRisk)
	all = append(all, all2...)
	writeJSON(w, http.StatusOK, InsightListResponse{Success: true, Data: &InsightList{Insights: all, Count: len(all)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetAchievements(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	all := h.composer.FilterByCategory(*inputs, engine.ICAchievement)
	all2 := h.composer.FilterByCategory(*inputs, engine.ICMilestone)
	all = append(all, all2...)
	writeJSON(w, http.StatusOK, InsightListResponse{Success: true, Data: &InsightList{Insights: all, Count: len(all)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetForecast(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	all := h.composer.FilterByCategory(*inputs, engine.ICForecast)
	all2 := h.composer.FilterByCategory(*inputs, engine.ICCashFlow)
	all = append(all, all2...)
	writeJSON(w, http.StatusOK, InsightListResponse{Success: true, Data: &InsightList{Insights: all, Count: len(all)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetTrends(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	all := h.composer.FilterByCategory(*inputs, engine.ICSavings)
	all2 := h.composer.FilterByCategory(*inputs, engine.ICNetWorth)
	all3 := h.composer.FilterByCategory(*inputs, engine.ICPortfolio)
	all = append(all, all2...)
	all = append(all, all3...)
	writeJSON(w, http.StatusOK, InsightListResponse{Success: true, Data: &InsightList{Insights: all, Count: len(all)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	limit := queryInt(r, "limit", 10)
	feed := h.composer.BuildFeed(*inputs, limit, "")
	writeJSON(w, http.StatusOK, FeedResponse{Success: true, Data: feed, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" { writeError(w, http.StatusBadRequest, "MISSING_QUERY", "search query 'q' is required"); return }
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	results := h.composer.Search(*inputs, q)
	writeJSON(w, http.StatusOK, InsightListResponse{Success: true, Data: &InsightList{Insights: results, Count: len(results)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "insight ID required"); return }
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	insight := h.composer.GetByID(*inputs, id)
	if insight == nil { writeError(w, http.StatusNotFound, "NOT_FOUND", "insight not found"); return }
	writeJSON(w, http.StatusOK, InsightResponse{Success: true, Data: insight, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}
