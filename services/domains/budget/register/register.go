package register

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/budget/internal/application"
	"github.com/horizon/core/services/domains/budget/internal/application/dto/command"
	"github.com/horizon/core/services/domains/budget/internal/application/dto/query"
	"github.com/horizon/core/services/domains/budget/internal/domain"
	"github.com/horizon/core/services/domains/budget/internal/infrastructure/persistence"
	"github.com/horizon/core/services/internal/auth"
)

type budgetHandler struct {
	svc *application.BudgetService
}

type noopPublisher struct{}

func (n *noopPublisher) Publish(event domain.DomainEvent) error { return nil }

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := persistence.NewBudgetRepository(pool)
	publisher := &noopPublisher{}
	svc := application.NewBudgetService(repo, publisher, time.Now)
	h := &budgetHandler{svc: svc}

	mux.HandleFunc("POST /api/v1/budgets", h.createBudget)
	mux.HandleFunc("GET /api/v1/budgets/{id}", h.getBudget)
	mux.HandleFunc("GET /api/v1/budgets", h.listBudgets)
	mux.HandleFunc("POST /api/v1/budgets/{id}/activate", h.activateBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/pause", h.pauseBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/resume", h.resumeBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/complete", h.completeBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/archive", h.archiveBudget)
	mux.HandleFunc("PUT /api/v1/budgets/{id}/category", h.updateCategory)
	mux.HandleFunc("GET /api/v1/budgets/{id}/vs-actual", h.budgetVsActual)
	mux.HandleFunc("GET /api/v1/households/{id}/budgets", h.listByHousehold)
}

func (h *budgetHandler) createBudget(w http.ResponseWriter, r *http.Request) {
	var cmd command.CreateBudgetCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.UserID = auth.UserIDFromRequest(r)
	result, err := h.svc.Create(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "CREATE_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, okData(result))
}

func (h *budgetHandler) getBudget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "Budget ID required"); return }
	result, err := h.svc.GetByID(r.Context(), query.GetBudgetQuery{BudgetID: id})
	if err != nil { writeError(w, http.StatusNotFound, "NOT_FOUND", "Budget not found"); return }
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *budgetHandler) listBudgets(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	cursor := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 { limit = 25 }
	result, err := h.svc.ListByUser(r.Context(), query.ListByUserQuery{UserID: userID, Cursor: cursor, Limit: limit})
	if err != nil { writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *budgetHandler) listByHousehold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "Household ID required"); return }

	role := auth.HouseholdRoleFromRequest(r, id)
	if role == "" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "You do not have access to this household")
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 { limit = 25 }
	result, err := h.svc.ListByHousehold(r.Context(), query.ListByHouseholdQuery{HouseholdID: id, Cursor: cursor, Limit: limit})
	if err != nil { writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *budgetHandler) activateBudget(w http.ResponseWriter, r *http.Request) {
	handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Activate(ctx, cmd)
	})
}
func (h *budgetHandler) pauseBudget(w http.ResponseWriter, r *http.Request) {
	handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Pause(ctx, cmd)
	})
}
func (h *budgetHandler) resumeBudget(w http.ResponseWriter, r *http.Request) {
	handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Resume(ctx, cmd)
	})
}
func (h *budgetHandler) completeBudget(w http.ResponseWriter, r *http.Request) {
	handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Complete(ctx, cmd)
	})
}
func (h *budgetHandler) archiveBudget(w http.ResponseWriter, r *http.Request) {
	handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Archive(ctx, cmd)
	})
}

func (h *budgetHandler) budgetVsActual(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "Budget ID required"); return }
	result, err := h.svc.GetByID(r.Context(), query.GetBudgetQuery{BudgetID: id})
	if err != nil { writeError(w, http.StatusNotFound, "NOT_FOUND", "Budget not found"); return }
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *budgetHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cmd command.UpdateCategoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.BudgetID = id
	result, err := h.svc.UpdateCategory(r.Context(), cmd)
	if err != nil { writeError(w, http.StatusInternalServerError, "UPDATE_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, okData(result))
}

type transitionFn func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error)

func handleTransition(w http.ResponseWriter, r *http.Request, fn transitionFn) {
	id := r.PathValue("id")
	if id == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "Budget ID required"); return }
	result, err := fn(r.Context(), command.BudgetIDCommand{BudgetID: id})
	if err != nil { writeError(w, http.StatusInternalServerError, "TRANSITION_ERROR", err.Error()); return }
	writeJSON(w, http.StatusOK, okData(result))
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

func okData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true, "data": data,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	}
}
