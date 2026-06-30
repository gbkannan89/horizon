package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/horizon/core/services/experiences/advisor/internal/aggregator"
	"github.com/horizon/core/services/experiences/advisor/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
}

func New(agg *aggregator.Aggregator, comp *engine.Composer) *Handlers {
	return &Handlers{aggregator: agg, composer: comp}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/advisor", h.GetAdvisor)
	mux.HandleFunc("GET /api/v1/advisor/summary", h.GetSummary)
	mux.HandleFunc("GET /api/v1/advisor/context", h.GetContext)
	mux.HandleFunc("GET /api/v1/advisor/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/advisor/health", h.GetHealth)
	mux.HandleFunc("GET /api/v1/advisor/risk", h.GetRisk)
	mux.HandleFunc("GET /api/v1/advisor/recommendations", h.GetRecommendations)
	mux.HandleFunc("GET /api/v1/advisor/insights", h.GetInsights)
	mux.HandleFunc("GET /api/v1/advisor/timeline", h.GetTimeline)
	mux.HandleFunc("GET /api/v1/advisor/notifications", h.GetNotifications)
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

func parseMode(r *http.Request) engine.AdvisorMode {
	m := engine.AdvisorMode(r.URL.Query().Get("mode"))
	if m == "" { return engine.AMOverview }
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

func (h *Handlers) GetAdvisor(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	ws := h.composer.BuildWorkspace(*inputs, parseMode(r))
	writeJSON(w, http.StatusOK, WorkspaceResponse{Success: true, Data: ws, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetDashboard(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	ws := h.composer.BuildWorkspace(*inputs, parseMode(r))
	writeJSON(w, http.StatusOK, WorkspaceResponse{Success: true, Data: ws, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetContext(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	ctx := h.composer.BuildContext(*inputs)
	writeJSON(w, http.StatusOK, ContextResponse{Success: true, Data: ctx, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	summary := h.composer.BuildSummary(*inputs)
	writeJSON(w, http.StatusOK, SummaryResponse{Success: true, Data: summary, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetHealth(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "health", CardType: engine.CTHealth,
		Title: "Health Score", Summary: engine.CtxSummary(inputs.HealthScore, inputs.HealthGrade, inputs.HealthChg),
		Data: map[string]interface{}{"score": inputs.HealthScore, "grade": inputs.HealthGrade, "change": inputs.HealthChg}}
	writeJSON(w, http.StatusOK, CardListResponse{Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetRisk(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "risk", CardType: engine.CTRisk,
		Title: "Risk Assessment", Summary: engine.CtxRiskSummary(inputs.RiskScore, inputs.RiskLevel),
		Data: map[string]interface{}{"score": inputs.RiskScore, "level": inputs.RiskLevel}}
	writeJSON(w, http.StatusOK, CardListResponse{Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	cards := []engine.Card{}
	if inputs.HasRecs {
		cards = append(cards, engine.Card{CardID: "rec", CardType: engine.CTBestAction,
			Title: inputs.RecTitle, Summary: engine.CtxRecSummary(inputs.RecCount)})
	}
	writeJSON(w, http.StatusOK, CardListResponse{Success: true, Data: &CardList{Cards: cards, Count: len(cards)}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetInsights(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "insight", CardType: engine.CTInsight,
		Title: "Active Insights", Summary: engine.CtxInsightSummary(inputs.InsightTotal, inputs.InsightCrit),
		Data: map[string]interface{}{"total": inputs.InsightTotal, "critical": inputs.InsightCrit}}
	writeJSON(w, http.StatusOK, CardListResponse{Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "timeline", CardType: engine.CTTimeline,
		Title: "Recent Activity", Summary: engine.CtxEventSummary(inputs.EventCount),
		Data: map[string]interface{}{"count": inputs.EventCount}}
	writeJSON(w, http.StatusOK, CardListResponse{Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetNotifications(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "notif", CardType: engine.CTNotifications,
		Title: "Unread Notifications", Summary: engine.CtxNotifSummary(inputs.NotifUnread),
		Data: map[string]interface{}{"unread": inputs.NotifUnread}}
	writeJSON(w, http.StatusOK, CardListResponse{Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1}, Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)}})
}
