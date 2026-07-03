package register

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/household/internal/application"
	"github.com/horizon/core/services/domains/household/internal/application/dto/command"
	"github.com/horizon/core/services/domains/household/internal/application/dto/query"
	"github.com/horizon/core/services/domains/household/internal/domain"
	"github.com/horizon/core/services/domains/household/internal/infrastructure/persistence"
	"github.com/horizon/core/services/internal/auth"
)

type householdHandler struct {
	svc *application.HouseholdService
}

type noopPublisher struct{}

func (n *noopPublisher) Publish(event domain.DomainEvent) error { return nil }

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := persistence.NewHouseholdRepository(pool)
	publisher := &noopPublisher{}
	svc := application.NewHouseholdService(repo, publisher, time.Now)
	h := &householdHandler{svc: svc}

	mux.HandleFunc("POST /api/v1/households", h.createHousehold)
	mux.HandleFunc("GET /api/v1/households/{id}", h.getHousehold)
	mux.HandleFunc("GET /api/v1/households", h.listHouseholds)
	mux.HandleFunc("POST /api/v1/households/{id}/activate", h.activateHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/pause", h.pauseHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/resume", h.resumeHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/dissolve", h.dissolveHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/archive", h.archiveHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/members", h.addMember)
	mux.HandleFunc("POST /api/v1/households/{id}/members/{userId}/accept", h.acceptInvite)
	mux.HandleFunc("DELETE /api/v1/households/{id}/members/{userId}", h.removeMember)
	mux.HandleFunc("PUT /api/v1/households/{id}/members/{userId}/role", h.updateMemberRole)
	mux.HandleFunc("POST /api/v1/households/{id}/accounts", h.linkAccount)
	mux.HandleFunc("DELETE /api/v1/households/{id}/accounts/{accountId}", h.unlinkAccount)
	mux.HandleFunc("POST /api/v1/households/{id}/goals", h.linkGoal)
	mux.HandleFunc("DELETE /api/v1/households/{id}/goals/{goalId}", h.unlinkGoal)
	mux.HandleFunc("POST /api/v1/households/{id}/budgets", h.linkBudget)
	mux.HandleFunc("DELETE /api/v1/households/{id}/budgets/{budgetId}", h.unlinkBudget)
	mux.HandleFunc("POST /api/v1/households/{id}/goals/{goalId}/contributions", h.addGoalContribution)
	mux.HandleFunc("GET /api/v1/households/{id}/summary", h.getHouseholdSummary)
}

func (h *householdHandler) createHousehold(w http.ResponseWriter, r *http.Request) {
	var cmd command.CreateHouseholdCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	result, err := h.svc.Create(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "CREATE_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) getHousehold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Household ID required")
		return
	}

	result, err := h.svc.GetByID(r.Context(), query.GetHouseholdQuery{HouseholdID: id})
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Household not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) listHouseholds(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	cursor := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 { limit = 25 }

	result, err := h.svc.ListByUser(r.Context(), query.ListByUserQuery{
		UserID: userID, Cursor: cursor, Limit: limit,
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

func (h *householdHandler) activateHousehold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	handleTransition(w, r, id, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Activate(ctx, cmd)
	})
}

func (h *householdHandler) pauseHousehold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	handleTransition(w, r, id, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Pause(ctx, cmd)
	})
}

func (h *householdHandler) resumeHousehold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	handleTransition(w, r, id, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Resume(ctx, cmd)
	})
}

func (h *householdHandler) dissolveHousehold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	handleTransition(w, r, id, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Dissolve(ctx, cmd)
	})
}

func (h *householdHandler) archiveHousehold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	handleTransition(w, r, id, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Archive(ctx, cmd)
	})
}

func (h *householdHandler) addMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cmd command.AddMemberCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.HouseholdID = id
	cmd.RequesterID = auth.UserIDFromRequest(r)

	result, err := h.svc.AddMember(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "MEMBER_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) acceptInvite(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.PathValue("userId")
	requesterID := auth.UserIDFromRequest(r)

	if userID != requesterID {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "You can only accept your own invites")
		return
	}

	cmd := command.AcceptInviteCommand{HouseholdID: id, UserID: userID}
	result, err := h.svc.AcceptInvite(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "MEMBER_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) removeMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.PathValue("userId")
	requesterID := auth.UserIDFromRequest(r)
	cmd := command.RemoveMemberCommand{HouseholdID: id, RequesterID: requesterID, UserID: userID}

	result, err := h.svc.RemoveMember(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "MEMBER_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) updateMemberRole(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.PathValue("userId")
	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	requesterID := auth.UserIDFromRequest(r)
	cmd := command.UpdateMemberRoleCommand{HouseholdID: id, RequesterID: requesterID, UserID: userID, Role: req.Role}
	result, err := h.svc.UpdateMemberRole(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "MEMBER_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) linkAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cmd command.LinkAccountCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.HouseholdID = id
	cmd.RequesterID = auth.UserIDFromRequest(r)
	result, err := h.svc.LinkAccount(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) unlinkAccount(w http.ResponseWriter, r *http.Request) {
	cmd := command.UnlinkAccountCommand{HouseholdID: r.PathValue("id"), AccountID: r.PathValue("accountId"), RequesterID: auth.UserIDFromRequest(r)}
	result, err := h.svc.UnlinkAccount(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "UNLINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) linkGoal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cmd command.LinkGoalCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.HouseholdID = id
	cmd.RequesterID = auth.UserIDFromRequest(r)
	result, err := h.svc.LinkGoal(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) unlinkGoal(w http.ResponseWriter, r *http.Request) {
	cmd := command.UnlinkGoalCommand{HouseholdID: r.PathValue("id"), GoalID: r.PathValue("goalId"), RequesterID: auth.UserIDFromRequest(r)}
	result, err := h.svc.UnlinkGoal(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "UNLINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) linkBudget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cmd command.LinkBudgetCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.HouseholdID = id
	cmd.RequesterID = auth.UserIDFromRequest(r)
	result, err := h.svc.LinkBudget(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) unlinkBudget(w http.ResponseWriter, r *http.Request) {
	cmd := command.UnlinkBudgetCommand{HouseholdID: r.PathValue("id"), BudgetID: r.PathValue("budgetId"), RequesterID: auth.UserIDFromRequest(r)}
	result, err := h.svc.UnlinkBudget(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "UNLINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) addGoalContribution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	goalID := r.PathValue("goalId")
	var cmd command.AddGoalContributionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.HouseholdID = id
	cmd.GoalID = goalID
	result, err := h.svc.AddGoalContribution(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "CONTRIB_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *householdHandler) getHouseholdSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Household ID required")
		return
	}

	result, err := h.svc.GetHouseholdFinancialSummary(r.Context(), query.GetHouseholdQuery{HouseholdID: id})
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Household not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": result,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

type transitionFunc func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error)

func handleTransition(w http.ResponseWriter, r *http.Request, id string, fn transitionFunc) {
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Household ID required")
		return
	}
	userID := auth.UserIDFromRequest(r)
	result, err := fn(r.Context(), command.HouseholdIDCommand{HouseholdID: id, RequesterID: userID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "TRANSITION_ERROR", err.Error())
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

// ---------------------------------------------------------------------------
// E2E Test Harness — exported for external e2e tests (tests/e2e/)
// ---------------------------------------------------------------------------

// E2ETestHarness provides an in-memory test server for household E2E tests.
type E2ETestHarness struct {
	Mux *http.ServeMux
	Svc *application.HouseholdService
}

// NewE2ETestHarness creates a test harness with an in-memory repository.
func NewE2ETestHarness() *E2ETestHarness {
	repo := newInMemRepository()
	publisher := &noopPublisher{}
	svc := application.NewHouseholdService(repo, publisher, time.Now)
	h := &householdHandler{svc: svc}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/households", h.createHousehold)
	mux.HandleFunc("GET /api/v1/households/{id}", h.getHousehold)
	mux.HandleFunc("GET /api/v1/households", h.listHouseholds)
	mux.HandleFunc("POST /api/v1/households/{id}/activate", h.activateHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/pause", h.pauseHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/resume", h.resumeHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/dissolve", h.dissolveHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/archive", h.archiveHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/members", h.addMember)
	mux.HandleFunc("POST /api/v1/households/{id}/members/{userId}/accept", h.acceptInvite)
	mux.HandleFunc("DELETE /api/v1/households/{id}/members/{userId}", h.removeMember)
	mux.HandleFunc("PUT /api/v1/households/{id}/members/{userId}/role", h.updateMemberRole)
	mux.HandleFunc("POST /api/v1/households/{id}/goals", h.linkGoal)
	mux.HandleFunc("DELETE /api/v1/households/{id}/goals/{goalId}", h.unlinkGoal)
	mux.HandleFunc("POST /api/v1/households/{id}/budgets", h.linkBudget)
	mux.HandleFunc("DELETE /api/v1/households/{id}/budgets/{budgetId}", h.unlinkBudget)
	mux.HandleFunc("POST /api/v1/households/{id}/goals/{goalId}/contributions", h.addGoalContribution)
	mux.HandleFunc("GET /api/v1/households/{id}/summary", h.getHouseholdSummary)

	return &E2ETestHarness{Mux: mux, Svc: svc}
}

// RegisterRoutesOnMux registers the household routes on the provided mux.
func (h *E2ETestHarness) RegisterRoutesOnMux(mux *http.ServeMux) {
	// Routes are already registered in NewE2ETestHarness, but this method
	// copies them to another mux for testing with httptest.
	handler := &householdHandler{svc: h.Svc}
	mux.HandleFunc("POST /api/v1/households", handler.createHousehold)
	mux.HandleFunc("GET /api/v1/households/{id}", handler.getHousehold)
	mux.HandleFunc("GET /api/v1/households", handler.listHouseholds)
	mux.HandleFunc("POST /api/v1/households/{id}/activate", handler.activateHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/pause", handler.pauseHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/resume", handler.resumeHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/dissolve", handler.dissolveHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/archive", handler.archiveHousehold)
	mux.HandleFunc("POST /api/v1/households/{id}/members", handler.addMember)
	mux.HandleFunc("POST /api/v1/households/{id}/members/{userId}/accept", handler.acceptInvite)
	mux.HandleFunc("DELETE /api/v1/households/{id}/members/{userId}", handler.removeMember)
	mux.HandleFunc("PUT /api/v1/households/{id}/members/{userId}/role", handler.updateMemberRole)
	mux.HandleFunc("POST /api/v1/households/{id}/goals", handler.linkGoal)
	mux.HandleFunc("DELETE /api/v1/households/{id}/goals/{goalId}", handler.unlinkGoal)
	mux.HandleFunc("POST /api/v1/households/{id}/budgets", handler.linkBudget)
	mux.HandleFunc("DELETE /api/v1/households/{id}/budgets/{budgetId}", handler.unlinkBudget)
	mux.HandleFunc("POST /api/v1/households/{id}/goals/{goalId}/contributions", handler.addGoalContribution)
	mux.HandleFunc("GET /api/v1/households/{id}/summary", handler.getHouseholdSummary)
}

// inMemRepository provides a thread-safe in-memory store for testing.
type inMemRepository struct {
	mu         sync.Mutex
	households map[string]*domain.Household
	counter    int
}

func newInMemRepository() *inMemRepository {
	return &inMemRepository{households: make(map[string]*domain.Household)}
}

func (r *inMemRepository) Save(_ context.Context, h *domain.Household) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := h.ID()
	if id == "" {
		r.counter++
		id = fmt.Sprintf("HH%04d", r.counter)
		now := time.Now().UTC()
		rebuilt := domain.ReconstructFromDB(
			id, h.Name(), h.HouseholdType(), h.HeadOfHouseholdID(),
			h.Members(), h.Status(), h.Currency(), h.Country(),
			h.TotalAssets(), h.TotalLiabilities(), h.TotalNetWorth(),
			h.Health(), h.Tags(), h.Notes(),
			h.LinkedAccounts(), h.LinkedGoals(), h.LinkedBudgets(),
			h.GoalContributions(), now, now,
		)
		*h = *rebuilt
		r.households[id] = rebuilt
	} else {
		r.households[id] = domain.ReconstructFromDB(
			id, h.Name(), h.HouseholdType(), h.HeadOfHouseholdID(),
			h.Members(), h.Status(), h.Currency(), h.Country(),
			h.TotalAssets(), h.TotalLiabilities(), h.TotalNetWorth(),
			h.Health(), h.Tags(), h.Notes(),
			h.LinkedAccounts(), h.LinkedGoals(), h.LinkedBudgets(),
			h.GoalContributions(), h.CreatedAt(), h.UpdatedAt(),
		)
	}
	return nil
}

func (r *inMemRepository) GetByID(_ context.Context, id string) (*domain.Household, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.households[id]
	if !ok {
		return nil, fmt.Errorf("household %s not found", id)
	}
	return h, nil
}

func (r *inMemRepository) ListByUser(_ context.Context, userID, cursor string, limit int) ([]*domain.Household, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Household
	for _, h := range r.households {
		if h.HeadOfHouseholdID() == userID {
			result = append(result, h)
			continue
		}
		for _, m := range h.Members() {
			if m.UserID == userID && m.InviteStatus == "Accepted" {
				result = append(result, h)
				break
			}
		}
	}
	if limit <= 0 {
		limit = 25
	}
	if len(result) > limit {
		result = result[:limit]
	}
	return result, "", nil
}

func (r *inMemRepository) ListByStatus(_ context.Context, status domain.HouseholdStatus, cursor string, limit int) ([]*domain.Household, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Household
	for _, h := range r.households {
		if h.Status() == status {
			result = append(result, h)
		}
	}
	if limit <= 0 {
		limit = 25
	}
	if len(result) > limit {
		result = result[:limit]
	}
	return result, "", nil
}

func (r *inMemRepository) UpdateStatus(_ context.Context, id string, from, to domain.HouseholdStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.households[id]
	if !ok {
		return fmt.Errorf("household %s not found", id)
	}
	if h.Status() != from {
		return fmt.Errorf("status mismatch")
	}
	r.households[id] = domain.ReconstructFromDB(
		h.ID(), h.Name(), h.HouseholdType(), h.HeadOfHouseholdID(),
		h.Members(), to, h.Currency(), h.Country(),
		h.TotalAssets(), h.TotalLiabilities(), h.TotalNetWorth(),
		h.Health(), h.Tags(), h.Notes(),
		h.LinkedAccounts(), h.LinkedGoals(), h.LinkedBudgets(),
		h.GoalContributions(), h.CreatedAt(), time.Now().UTC(),
	)
	return nil
}

// Ensure inMemRepository satisfies domain.Repository.
var _ domain.Repository = (*inMemRepository)(nil)
