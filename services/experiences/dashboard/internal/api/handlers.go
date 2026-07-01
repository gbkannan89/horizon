package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/horizon/core/services/experiences/dashboard/internal/aggregator"
	"github.com/horizon/core/services/experiences/dashboard/internal/engine"
	"github.com/horizon/core/services/internal/auth"
)

type Handlers struct {
	aggregator *aggregator.DashboardAggregator
	composer   *engine.Composer
}

func New(agg *aggregator.DashboardAggregator, comp *engine.Composer) *Handlers {
	return &Handlers{aggregator: agg, composer: comp}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/dashboard/summary", h.GetSummary)
	mux.HandleFunc("GET /api/v1/dashboard/widgets", h.GetWidgets)
	mux.HandleFunc("GET /api/v1/dashboard/health", h.GetHealth)
	mux.HandleFunc("GET /api/v1/dashboard/recommendations", h.GetRecommendations)
	mux.HandleFunc("GET /api/v1/dashboard/milestones", h.GetMilestones)
}

// getDefaultUserID extracts from Authorization or query param.
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
	resp := map[string]interface{}{
		"success": false, "error": map[string]string{"code": code, "message": message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	}
	writeJSON(w, status, resp)
}

func (h *Handlers) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	mode := engine.DashboardMode(r.URL.Query().Get("mode"))
	if mode == "" { mode = engine.DMOverview }

	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", fmt.Sprintf("failed to load dashboard data: %v", err))
		return
	}

	output := h.composer.Compose(*inputs, mode)
	writeJSON(w, http.StatusOK, DashboardResponse{
		Success: true, Data: output,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	summary := h.composer.BuildSummary(*inputs)
	writeJSON(w, http.StatusOK, SummaryResponse{
		Success: true, Data: summary,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetWidgets(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	mode := engine.DashboardMode(r.URL.Query().Get("mode"))
	if mode == "" { mode = engine.DMOverview }

	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	output := h.composer.Compose(*inputs, mode)
	var all []engine.Widget
	all = append(all, output.Tier1...)
	all = append(all, output.Tier2...)
	all = append(all, output.Tier3...)
	if output.CriticalAlert != nil {
		all = append([]engine.Widget{*output.CriticalAlert}, all...)
	}
	writeJSON(w, http.StatusOK, WidgetsResponse{
		Success: true, Data: &WidgetList{Widgets: all, Count: len(all)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetHealth(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	widget := engine.Widget{
		WidgetID: "health", WidgetType: "health_score",
		Title: fmt.Sprintf("Health Score: %d (%s)", inputs.HealthScore, inputs.HealthGrade),
		Data:  inputs.HealthScore, Confidence: "High", Visible: true,
	}
	writeJSON(w, http.StatusOK, HealthResponse{
		Success: true, Data: &widget,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	var recs []engine.Widget
	if inputs.HasRecommendation {
		recs = append(recs, engine.Widget{
			WidgetID: "rec", WidgetType: "recommendation",
			Title: fmt.Sprintf("%d Recommendation(s) Available", inputs.RecCount),
			Data:  inputs.RecCount, Confidence: "High", Visible: true,
		})
	}
	writeJSON(w, http.StatusOK, RecommendationsResponse{
		Success: true, Data: &RecList{Recommendations: recs, Count: len(recs)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetMilestones(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	inputs, err := h.aggregator.Aggregate(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	milestones := []engine.Widget{
		{WidgetID: "networth", WidgetType: "milestone", Title: fmt.Sprintf("Net Worth: %s", fmtMoney(inputs.NetWorth)), Data: inputs.NetWorth, Visible: true},
	}
	if inputs.GoalCount > 0 {
		milestones = append(milestones, engine.Widget{
			WidgetID: "goals", WidgetType: "milestone",
			Title: fmt.Sprintf("%d/%d Goals On Track", inputs.GoalsOnTrack, inputs.GoalCount),
			Data:  map[string]int{"on_track": inputs.GoalsOnTrack, "total": inputs.GoalCount}, Visible: true,
		})
	}
	writeJSON(w, http.StatusOK, MilestonesResponse{
		Success: true, Data: &MilestoneList{Milestones: milestones, Count: len(milestones)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func fmtMoney(v int64) string {
	if v >= 10000000 { return fmt.Sprintf("₹%.2fCr", float64(v)/10000000) }
	if v >= 100000 { return fmt.Sprintf("₹%.2fL", float64(v)/100000) }
	return fmt.Sprintf("₹%d", v)
}
