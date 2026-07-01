package register

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/financial-event/internal/application"
	"github.com/horizon/core/services/domains/financial-event/internal/application/dto/query"
	"github.com/horizon/core/services/domains/financial-event/internal/domain"
	"github.com/horizon/core/services/domains/financial-event/internal/infrastructure/persistence"
	"github.com/horizon/core/services/internal/auth"
)

type transactionsHandler struct {
	svc *application.EventService
}

type noopPublisher struct{}

func (n *noopPublisher) Publish(event domain.DomainEvent) error { return nil }

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := persistence.NewFinancialEventRepository(pool)
	publisher := &noopPublisher{}
	svc := application.NewEventService(repo, publisher, time.Now)
	h := &transactionsHandler{svc: svc}

	mux.HandleFunc("GET /api/v1/transactions", h.listTransactions)
	mux.HandleFunc("GET /api/v1/transactions/{id}", h.getTransaction)
	mux.HandleFunc("GET /api/v1/transactions/search", h.searchTransactions)
	mux.HandleFunc("GET /api/v1/transactions/summary", h.transactionSummary)

	// Event write endpoints (mock MVP — Phase 2+ will implement real commands)
	mux.HandleFunc("POST /api/v1/events", h.createEvent)
	mux.HandleFunc("POST /api/v1/events/{id}/archive", h.archiveEvent)
}

func (h *transactionsHandler) listTransactions(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	limit := queryLimit(r, 25)
	cursor := r.URL.Query().Get("cursor")

	result, err := h.svc.ListByUser(r.Context(), query.ListEventsByUserQuery{
		UserID: userID, Cursor: cursor, Limit: limit,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"transactions": result.Events,
			"cursor":       result.NextCursor,
			"has_more":     result.HasMore,
			"total":        len(result.Events),
		},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *transactionsHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "transaction ID required")
		return
	}

	result, err := h.svc.GetEvent(r.Context(), query.GetEventQuery{EventID: id})
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "transaction not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"data":     result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *transactionsHandler) searchTransactions(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	q := r.URL.Query().Get("q")
	limit := queryLimit(r, 25)

	if q == "" {
		writeError(w, http.StatusBadRequest, "MISSING_QUERY", "search query 'q' is required")
		return
	}

	result, err := h.svc.ListByUser(r.Context(), query.ListEventsByUserQuery{
		UserID: userID, Limit: 100,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}

	var filtered []query.EventResult
	q = strings.ToLower(q)
	for _, evt := range result.Events {
		if strings.Contains(strings.ToLower(evt.Description), q) ||
			strings.Contains(strings.ToLower(evt.EventType), q) ||
			strings.Contains(strings.ToLower(evt.Source), q) ||
			strings.Contains(strings.ToLower(evt.Destination), q) {
			filtered = append(filtered, evt)
		}
	}
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"transactions": filtered,
			"total":        len(filtered),
		},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *transactionsHandler) transactionSummary(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	result, err := h.svc.ListByDateRange(r.Context(), query.ListEventsByDateRangeQuery{
		UserID: userID, Start: startOfMonth, End: now, Limit: 1000,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}

	var income, expenses int64
	var incomeCount, expenseCount int
	for _, e := range result.Events {
		if e.Amount > 0 {
			income += e.Amount
			incomeCount++
		} else {
			expenses += e.Amount
			expenseCount++
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"period_income":   income,
			"period_expenses": expenses,
			"net_flow":        income + expenses,
			"income_count":    incomeCount,
			"expense_count":   expenseCount,
			"total_count":     incomeCount + expenseCount,
		},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *transactionsHandler) createEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type     string  `json:"type"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       "evt-mock",
		"type":     req.Type,
		"amount":   req.Amount,
		"currency": req.Currency,
	})
}

func (h *transactionsHandler) archiveEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":   "archived",
		"event_id": id,
	})
}

func getUserID(r *http.Request) string {
	return auth.UserIDFromRequest(r)
}

func queryLimit(r *http.Request, def int) int {
	v := r.URL.Query().Get("limit")
	if v == "" { return def }
	n, err := strconv.Atoi(v)
	if err != nil { return def }
	if n <= 0 { return def }
	return n
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false, "error": map[string]string{"code": code, "message": message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}
