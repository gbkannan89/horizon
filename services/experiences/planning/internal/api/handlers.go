package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/horizon/core/services/internal/auth"
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
	mux.HandleFunc("GET /api/v1/planning/budget", h.GetBudget)
	mux.HandleFunc("GET /api/v1/planning/retirement", h.GetRetirement)
	mux.HandleFunc("GET /api/v1/planning/emergency-fund", h.GetEmergencyFund)
	mux.HandleFunc("GET /api/v1/planning/debt-payoff", h.GetDebtPayoff)
	mux.HandleFunc("GET /api/v1/planning/investment", h.GetInvestment)
}

func getDefaultUserID(r *http.Request) string {
	return auth.UserIDFromRequest(r)
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

func (h *Handlers) GetBudget(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{
		"total_budget": 109000, "total_spent": 107500, "surplus": 1500,
		"categories": []map[string]interface{}{
			{"category": "Housing", "planned": 45000, "actual": 43500, "remaining": 1500},
			{"category": "Food", "planned": 15000, "actual": 16200, "remaining": -1200},
			{"category": "Transport", "planned": 8000, "actual": 7200, "remaining": 800},
			{"category": "Utilities", "planned": 6000, "actual": 5800, "remaining": 200},
			{"category": "Entertainment", "planned": 5000, "actual": 4800, "remaining": 200},
			{"category": "Savings", "planned": 30000, "actual": 30000, "remaining": 0},
		},
	}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetRetirement(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{
		"current_corpus": 38500000, "target_corpus": 50000000, "progress_pct": 77,
		"monthly_contribution": 30000, "recommended_contribution": 45000,
		"projected_retirement_age": 58, "on_track": true, "confidence": "Medium",
	}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetEmergencyFund(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{
		"current_savings": 450000, "target_amount": 600000, "months_covered": 4.5,
		"target_months": 6, "monthly_expenses": 100000, "progress_pct": 75, "status": "In Progress",
	}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetDebtPayoff(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{
		"total_debt": 4545000, "total_monthly": 90000, "debt_to_income_ratio": 28.0,
		"debts": []map[string]interface{}{
			{"name": "Home Loan", "principal": 4500000, "interest_rate": 8.5, "monthly_emi": 45000, "remaining_months": 180},
			{"name": "Credit Card", "principal": 45000, "interest_rate": 42.0, "monthly_emi": 45000, "remaining_months": 1},
		}, "debt_free_date": "2032-06-01",
	}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetInvestment(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{
		"portfolio_value": 4580000, "total_return": 12.5, "risk_level": "Moderate", "rebalance_needed": true,
		"allocations": []map[string]interface{}{
			{"type": "Equity", "current": 45, "target": 50, "drift": -5},
			{"type": "Debt", "current": 30, "target": 25, "drift": 5},
			{"type": "Gold", "current": 10, "target": 10, "drift": 0},
			{"type": "Cash", "current": 5, "target": 5, "drift": 0},
		},
	}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}
