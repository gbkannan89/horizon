package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/horizon/core/services/experiences/portfolio/internal/aggregator"
	"github.com/horizon/core/services/experiences/portfolio/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
}

func New(agg *aggregator.Aggregator, comp *engine.Composer) *Handlers {
	return &Handlers{aggregator: agg, composer: comp}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/portfolio/experience", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/portfolio/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/portfolio/allocation", h.GetAllocation)
	mux.HandleFunc("GET /api/v1/portfolio/performance", h.GetPerformance)
	mux.HandleFunc("GET /api/v1/portfolio/risk", h.GetRisk)
	mux.HandleFunc("GET /api/v1/portfolio/projection", h.GetProjection)
	mux.HandleFunc("GET /api/v1/portfolio/recommendations", h.GetRecommendations)
	mux.HandleFunc("GET /api/v1/portfolio/optimization", h.GetOptimization)
	mux.HandleFunc("GET /api/v1/portfolio/simulations", h.GetSimulations)
	mux.HandleFunc("GET /api/v1/portfolio/timeline", h.GetTimeline)
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

func (h *Handlers) GetDashboard(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, DashboardResponse{
		Success: true, Data: h.composer.BuildDashboard(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetAllocation(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, AllocationResponse{
		Success: true, Data: h.composer.BuildAllocation(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetPerformance(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, PerformanceResponse{
		Success: true, Data: h.composer.BuildPerformance(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRisk(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, RiskResponse{
		Success: true, Data: h.composer.BuildRisk(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetProjection(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, ProjectionResponse{
		Success: true, Data: h.composer.BuildProjection(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, RecsResponse{
		Success: true, Data: h.composer.BuildRecommendations(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetOptimization(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, OptResponse{
		Success: true, Data: h.composer.BuildOptimization(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetSimulations(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, SimResponse{
		Success: true, Data: h.composer.BuildSimulations(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, TimelineResponse{
		Success: true, Data: h.composer.BuildCardTimeline(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}
