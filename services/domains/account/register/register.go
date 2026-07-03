package register

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/account/internal/application"
	"github.com/horizon/core/services/domains/account/internal/application/dto/query"
	"github.com/horizon/core/services/domains/account/internal/domain"
	"github.com/horizon/core/services/domains/account/internal/infrastructure/persistence"
	"github.com/horizon/core/services/internal/auth"
)

type accountHandler struct {
	svc *application.AccountService
}

type noopPublisher struct{}

func (n *noopPublisher) Publish(event domain.DomainEvent) error { return nil }

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := persistence.NewAccountRepository(pool)
	publisher := &noopPublisher{}
	svc := application.NewAccountService(repo, publisher, time.Now)
	h := &accountHandler{svc: svc}

	mux.HandleFunc("GET /api/v1/households/{id}/accounts", h.listHouseholdAccounts)
}

func (h *accountHandler) listHouseholdAccounts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Household ID required")
		return
	}

	role := auth.HouseholdRoleFromRequest(r, id)
	if role == "" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "You do not have access to this household")
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 { limit = 25 }

	result, err := h.svc.ListByHousehold(r.Context(), query.ListByHouseholdQuery{
		HouseholdID: id, Cursor: cursor, Limit: limit,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
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
