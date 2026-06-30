package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/horizon/core/services/experiences/accounts/internal/aggregator"
	"github.com/horizon/core/services/experiences/accounts/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
	acctEngine *engine.Engine
}

func New(agg *aggregator.Aggregator, comp *engine.Composer, ae *engine.Engine) *Handlers {
	return &Handlers{aggregator: agg, composer: comp, acctEngine: ae}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/accounts/experience", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/accounts/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/accounts/summary", h.GetSummary)
	mux.HandleFunc("GET /api/v1/accounts/balances", h.GetBalances)
	mux.HandleFunc("GET /api/v1/accounts/cashflow", h.GetCashFlow)
	mux.HandleFunc("GET /api/v1/accounts/health", h.GetHealth)
	mux.HandleFunc("GET /api/v1/accounts/recommendations", h.GetRecommendations)
	mux.HandleFunc("GET /api/v1/accounts/projections", h.GetProjections)
	mux.HandleFunc("GET /api/v1/accounts/timeline", h.GetTimeline)
	mux.HandleFunc("GET /api/v1/accounts/{id}", h.GetAccountDetail)
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

func parseFilters(r *http.Request) *engine.FilterOpts {
	types := r.URL.Query()["type"]
	statuses := r.URL.Query()["status"]
	currencies := r.URL.Query()["currency"]
	if len(types) == 0 && len(statuses) == 0 && len(currencies) == 0 { return nil }
	return &engine.FilterOpts{AccountTypes: types, Statuses: statuses, Currencies: currencies}
}

func queryLimit(r *http.Request, def int) int {
	v := r.URL.Query().Get("limit")
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

func (h *Handlers) GetDashboard(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, DashboardResponse{
		Success: true, Data: h.composer.BuildDashboard(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	list := h.acctEngine.BuildAccountsList(*inputs, parseFilters(r), r.URL.Query().Get("cursor"), queryLimit(r, 25))
	writeJSON(w, http.StatusOK, ListResponse{
		Success: true, Data: list,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetBalances(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, BalanceResponse{
		Success: true, Data: h.composer.BuildBalanceSummary(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetCashFlow(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, CashFlowResponse{
		Success: true, Data: h.composer.BuildCashFlowSummary(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetHealth(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, HealthResponse{
		Success: true, Data: h.composer.BuildHealthSummary(*inputs),
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	cards := []engine.Card{}
	if inputs.HasRecs {
		cards = append(cards, engine.Card{CardID: "rec", CardType: engine.CTRecommendation, Title: "Recommendations",
			Summary: engine.SummaryCount(inputs.RecCount)})
	}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: cards, Count: len(cards)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetProjections(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "proj", CardType: engine.CTProjection,
		Title: "Cash Flow Projection",
		Summary: engine.SummaryMoney(inputs.ProjectedFlow),
		Data: map[string]interface{}{"projected": inputs.ProjectedFlow}}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	inputs, err := h.aggregator.Aggregate(r.Context(), getDefaultUserID(r))
	if err != nil { writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error()); return }
	card := engine.Card{CardID: "tl", CardType: engine.CTTimeline,
		Title: "Recent Activity",
		Summary: engine.SummaryCount(inputs.EventCount)}
	writeJSON(w, http.StatusOK, CardListResponse{
		Success: true, Data: &CardList{Cards: []engine.Card{card}, Count: 1},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetAccountDetail(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("id")
	if accountID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "account ID required"); return }

	inputs, err := h.aggregator.GetAccountDetail(r.Context(), getDefaultUserID(r), accountID)
	if err != nil { writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error()); return }
	if inputs.Account == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "account not found")
		return
	}
	detail := h.acctEngine.BuildAccountDetail(*inputs.Account, inputs.Transactions)
	writeJSON(w, http.StatusOK, DetailResponse{
		Success: true, Data: detail,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}
