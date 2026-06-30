package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/horizon/core/services/experiences/planning/internal/aggregator"
	"github.com/horizon/core/services/experiences/planning/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
}

func New(agg *aggregator.Aggregator, comp *engine.Composer) *Handlers {
	return &Handlers{aggregator: agg, composer: comp}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/planning", h.GetPlanning)
	mux.HandleFunc("GET /api/v1/planning/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/planning/projections", h.GetProjections)
	mux.HandleFunc("GET /api/v1/planning/scenarios", h.GetScenarios)
	mux.HandleFunc("GET /api/v1/planning/compare", h.GetCompare)
	mux.HandleFunc("GET /api/v1/planning/recommendations", h.GetRecommendations)
	mux.HandleFunc("GET /api/v1/planning/optimizations", h.GetOptimizations)
	mux.HandleFunc("GET /api/v1/planning/risk", h.GetRisk)
	mux.HandleFunc("GET /api/v1/planning/health", h.GetHealth)
	mux.HandleFunc("GET /api/v1/planning/timeline", h.GetTimeline)
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

func parseMode(r *http.Request) engine.PlanningMode {
	m := engine.PlanningMode(r.URL.Query().Get("mode"))
	if m == "" { return engine.PMOverview }
	return m
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

func (h *Handlers) GetPlanning(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, DashboardResponse{
		Success: true, Data: h.composer.BuildDashboard(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetDashboard(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, DashboardResponse{
		Success: true, Data: h.composer.BuildDashboard(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetProjections(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, ProjectionResponse{
		Success: true, Data: h.composer.BuildProjections(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetScenarios(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), engine.PMScenarioPlan)
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, ScenarioListResponse{
		Success: true, Data: &ScenarioList{Scenarios: inputs.Scenarios, Total: len(inputs.Scenarios)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetCompare(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), engine.PMScenarioPlan)
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }

	if len(inputs.Scenarios) < 1 {
		writeError(w, http.StatusNotFound, "NO_SCENARIOS", "no scenarios available for comparison")
		return
	}

	baseline := engine.PlanSummary{
		PlanID: "baseline", Name: "Current Plan", IsBaseline: true,
		NetWorthEnd: inputs.NetWorthProjected, GoalsOnTrack: inputs.GoalsOnTrack,
		TotalGoals: inputs.TotalGoals, FundingGap: inputs.FundingGap,
		RiskScore: inputs.RiskScore, HealthScore: inputs.HealthScore,
	}

	alt := inputs.Scenarios[0]
	altPlan := engine.PlanSummary{
		PlanID: alt.ScenarioID, Name: alt.Name, IsBaseline: false,
		NetWorthEnd: alt.Projections.NetWorthProjected,
		GoalsOnTrack: inputs.GoalsOnTrack, TotalGoals: inputs.TotalGoals,
		FundingGap: inputs.FundingGap, RiskScore: inputs.RiskScore,
		HealthScore: inputs.HealthScore,
	}

	writeJSON(w, http.StatusOK, ComparisonResponse{
		Success: true, Data: h.composer.BuildComparison(baseline, altPlan),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	cards := []engine.Card{}
	if inputs.HasRecs {
		cards = append(cards, engine.Card{CardID: "rec", CardType: engine.CTRecommend,
			Title: "Active Recommendations", Summary: engine.CardSummary(inputs.RecCount)})
	}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: cards, Count: len(cards)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetOptimizations(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	cards := []engine.Card{}
	if inputs.HasOptimization {
		cards = append(cards, engine.Card{CardID: "opt", CardType: engine.CTOptimization,
			Title: "Optimization Strategies", Summary: engine.CardSummary(inputs.OptCount)})
	}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: cards, Count: len(cards)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRisk(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "risk", CardType: engine.CTRiskSum,
		Title: "Risk Assessment",
		Summary: engine.CardRiskSummary(inputs.RiskScore, inputs.RiskLevel),
		Data:    map[string]interface{}{"score": inputs.RiskScore, "level": inputs.RiskLevel}}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetHealth(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "health", CardType: engine.CTPlanOverview,
		Title: "Health Score",
		Summary: engine.CardHealthSummary(inputs.HealthScore, inputs.HealthGrade),
		Data:    map[string]interface{}{"score": inputs.HealthScore, "grade": inputs.HealthGrade}}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r), parseMode(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "tl", CardType: engine.CTTimeline,
		Title: "Recent Events",
		Summary: engine.CardEventSummary(inputs.EventCount),
		Data:    map[string]interface{}{"events": inputs.EventCount, "milestones": inputs.MilestoneCount}}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}
