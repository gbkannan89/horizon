package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/financial-event/internal/application"
	"github.com/horizon/core/services/domains/financial-event/internal/application/dto/query"
	"github.com/horizon/core/services/domains/financial-event/internal/domain"
	"github.com/horizon/core/services/domains/financial-event/internal/infrastructure/persistence"
)

type server struct {
	svc  *application.EventService
	pool *pgxpool.Pool
}

func main() {
	dbURL := env("DATABASE_URL", "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable")
	port := env("PORT", "8060")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	repo := persistence.NewFinancialEventRepository(pool)
	publisher := &noopPublisher{}
	svc := application.NewEventService(repo, publisher, time.Now)

	s := &server{svc: svc, pool: pool}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", health("live"))
	mux.HandleFunc("GET /health/ready", health("ready"))
	mux.HandleFunc("GET /api/v1/transactions", s.listTransactions)
	mux.HandleFunc("GET /api/v1/transactions/{id}", s.getTransaction)
	mux.HandleFunc("GET /api/v1/transactions/search", s.searchTransactions)
	mux.HandleFunc("GET /api/v1/transactions/summary", s.transactionSummary)

	addr := ":" + port
	log.Printf("Starting Transactions (Financial Event) service on %s", addr)
	if err := http.ListenAndServe(addr, recovery(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

type noopPublisher struct{}

func (n *noopPublisher) Publish(event domain.DomainEvent) error { return nil }

func env(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}

func health(status string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": status})
	}
}

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"success": false, "error": map[string]string{"code": "PANIC", "message": "internal error"},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
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

func getUserID(r *http.Request) string {
	u := r.URL.Query().Get("user_id")
	if u != "" { return u }
	if a := r.Header.Get("Authorization"); strings.HasPrefix(a, "Bearer ") {
		if uid := strings.TrimPrefix(a, "Bearer "); uid != "" { return uid }
	}
	return "default"
}

func queryLimit(r *http.Request, def int) int {
	v := r.URL.Query().Get("limit")
	if v == "" { return def }
	n, err := strconv.Atoi(v)
	if err != nil { return def }
	if n <= 0 { return def }
	return n
}

// ---- Handlers ----

func (s *server) listTransactions(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	limit := queryLimit(r, 25)
	cursor := r.URL.Query().Get("cursor")

	result, err := s.svc.ListByUser(r.Context(), query.ListEventsByUserQuery{
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

func (s *server) getTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "transaction ID required")
		return
	}

	result, err := s.svc.GetEvent(r.Context(), query.GetEventQuery{EventID: id})
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

func (s *server) searchTransactions(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	q := r.URL.Query().Get("q")
	limit := queryLimit(r, 25)

	if q == "" {
		writeError(w, http.StatusBadRequest, "MISSING_QUERY", "search query 'q' is required")
		return
	}

	result, err := s.svc.ListByUser(r.Context(), query.ListEventsByUserQuery{
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

func (s *server) transactionSummary(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	result, err := s.svc.ListByDateRange(r.Context(), query.ListEventsByDateRangeQuery{
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
