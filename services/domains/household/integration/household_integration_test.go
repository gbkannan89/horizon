package household_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/horizon/core/services/domains/household/internal/application"
	"github.com/horizon/core/services/domains/household/internal/application/dto/command"
	"github.com/horizon/core/services/domains/household/internal/application/dto/query"
	"github.com/horizon/core/services/domains/household/internal/domain"
	"github.com/horizon/core/services/internal/auth"
)

// ---------------------------------------------------------------------------
// Mock Repository
// ---------------------------------------------------------------------------

type mockRepo struct {
	mu         sync.Mutex
	households map[string]*domain.Household
}

func newMockRepo() *mockRepo {
	return &mockRepo{households: make(map[string]*domain.Household)}
}

func (r *mockRepo) Save(_ context.Context, h *domain.Household) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := h.ID()
	if id == "" {
		id = fmt.Sprintf("hh-%d", len(r.households)+1)
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
	}
	r.households[h.ID()] = h
	return nil
}

func (r *mockRepo) UpdateStatus(_ context.Context, id string, from, to domain.HouseholdStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.households[id]
	if !ok {
		return fmt.Errorf("household not found")
	}
	if h.Status() != from {
		return fmt.Errorf("status mismatch: got %s, want %s", h.Status(), from)
	}
	h.SetStatus(to)
	return nil
}

func (r *mockRepo) GetByID(_ context.Context, id string) (*domain.Household, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	h, ok := r.households[id]
	if !ok {
		return nil, fmt.Errorf("household not found")
	}
	return h, nil
}

func (r *mockRepo) ListByUser(_ context.Context, userID string, _ string, _ int) ([]*domain.Household, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	seen := make(map[string]bool)
	var result []*domain.Household
	for _, h := range r.households {
		if h.HeadOfHouseholdID() == userID {
			if !seen[h.ID()] {
				result = append(result, h)
				seen[h.ID()] = true
			}
		}
		for _, m := range h.Members() {
			if m.UserID == userID && !seen[h.ID()] {
				result = append(result, h)
				seen[h.ID()] = true
				break
			}
		}
	}
	return result, "", nil
}

func (r *mockRepo) ListByStatus(_ context.Context, status domain.HouseholdStatus, _ string, _ int) ([]*domain.Household, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Household
	for _, h := range r.households {
		if h.Status() == status {
			result = append(result, h)
		}
	}
	return result, "", nil
}

// ---------------------------------------------------------------------------
// Mock Publisher
// ---------------------------------------------------------------------------

type mockPublisher struct{}

func (m *mockPublisher) Publish(_ domain.DomainEvent) error { return nil }

// ---------------------------------------------------------------------------
// Response envelope types
// ---------------------------------------------------------------------------

type apiResponse struct {
	Success  bool              `json:"success"`
	Data     json.RawMessage   `json:"data,omitempty"`
	Error    map[string]string `json:"error,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ---------------------------------------------------------------------------
// HTTP handlers (test double matching register/householdHandler)
// ---------------------------------------------------------------------------

type handler struct {
	svc *application.HouseholdService
}

func newHandler(svc *application.HouseholdService) *handler {
	return &handler{svc: svc}
}

func (h *handler) createHousehold(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusCreated, okData(result))
}

func (h *handler) getHousehold(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) listHouseholds(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	cursor := r.URL.Query().Get("cursor")
	limit := 25
	result, err := h.svc.ListByUser(r.Context(), query.ListByUserQuery{UserID: userID, Cursor: cursor, Limit: limit})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) activateHousehold(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Activate(ctx, cmd)
	})
}

func (h *handler) pauseHousehold(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Pause(ctx, cmd)
	})
}

func (h *handler) resumeHousehold(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Resume(ctx, cmd)
	})
}

func (h *handler) dissolveHousehold(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Dissolve(ctx, cmd)
	})
}

func (h *handler) archiveHousehold(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
		return h.svc.Archive(ctx, cmd)
	})
}

type transitionFn func(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error)

func (h *handler) handleTransition(w http.ResponseWriter, r *http.Request, fn transitionFn) {
	id := r.PathValue("id")
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) addMember(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) acceptInvite(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) removeMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.PathValue("userId")
	requesterID := auth.UserIDFromRequest(r)
	cmd := command.RemoveMemberCommand{HouseholdID: id, RequesterID: requesterID, UserID: userID}
	result, err := h.svc.RemoveMember(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "MEMBER_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) updateMemberRole(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) linkAccount(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) unlinkAccount(w http.ResponseWriter, r *http.Request) {
	cmd := command.UnlinkAccountCommand{HouseholdID: r.PathValue("id"), AccountID: r.PathValue("accountId"), RequesterID: auth.UserIDFromRequest(r)}
	result, err := h.svc.UnlinkAccount(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "UNLINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) linkGoal(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) unlinkGoal(w http.ResponseWriter, r *http.Request) {
	cmd := command.UnlinkGoalCommand{HouseholdID: r.PathValue("id"), GoalID: r.PathValue("goalId"), RequesterID: auth.UserIDFromRequest(r)}
	result, err := h.svc.UnlinkGoal(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "UNLINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) linkBudget(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) unlinkBudget(w http.ResponseWriter, r *http.Request) {
	cmd := command.UnlinkBudgetCommand{HouseholdID: r.PathValue("id"), BudgetID: r.PathValue("budgetId"), RequesterID: auth.UserIDFromRequest(r)}
	result, err := h.svc.UnlinkBudget(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "UNLINK_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) addGoalContribution(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) getHouseholdSummary(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, okData(result))
}

// ---------------------------------------------------------------------------
// Response helpers
// ---------------------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success":  false,
		"error":    map[string]string{"code": code, "message": message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func okData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success":  true,
		"data":     data,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	}
}

// ---------------------------------------------------------------------------
// Test server setup
// ---------------------------------------------------------------------------

func newTestServer(t *testing.T) (*httptest.Server, *mockRepo) {
	t.Helper()
	repo := newMockRepo()
	pub := &mockPublisher{}
	svc := application.NewHouseholdService(repo, pub, time.Now)
	h := newHandler(svc)

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
	mux.HandleFunc("POST /api/v1/households/{id}/accounts", h.linkAccount)
	mux.HandleFunc("DELETE /api/v1/households/{id}/accounts/{accountId}", h.unlinkAccount)
	mux.HandleFunc("POST /api/v1/households/{id}/goals", h.linkGoal)
	mux.HandleFunc("DELETE /api/v1/households/{id}/goals/{goalId}", h.unlinkGoal)
	mux.HandleFunc("POST /api/v1/households/{id}/budgets", h.linkBudget)
	mux.HandleFunc("DELETE /api/v1/households/{id}/budgets/{budgetId}", h.unlinkBudget)
	mux.HandleFunc("POST /api/v1/households/{id}/goals/{goalId}/contributions", h.addGoalContribution)
	mux.HandleFunc("GET /api/v1/households/{id}/summary", h.getHouseholdSummary)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, repo
}

func doRequest(t *testing.T, srv *httptest.Server, method, path string, body []byte) *http.Response {
	t.Helper()
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, srv.URL+path, reqBody)
	if err != nil {
		t.Fatalf("creating %s %s: %v", method, path, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("executing %s %s: %v", method, path, err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func parseResponse(t *testing.T, resp *http.Response) apiResponse {
	t.Helper()
	var ar apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	return ar
}

func mustUnmarshalData(t *testing.T, raw json.RawMessage, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(raw, target); err != nil {
		t.Fatalf("unmarshalling data: %v", err)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected status %d, got %d", want, got)
	}
}

func assertSuccess(t *testing.T, ar apiResponse) {
	t.Helper()
	if !ar.Success {
		t.Fatal("expected success=true")
	}
	if ar.Error != nil {
		t.Fatalf("unexpected error: %v", ar.Error)
	}
	if ar.Data == nil {
		t.Fatal("expected data to be present")
	}
	if ar.Metadata == nil {
		t.Fatal("expected metadata to be present")
	}
	if _, ok := ar.Metadata["timestamp"]; !ok {
		t.Fatal("expected timestamp in metadata")
	}
}

func assertError(t *testing.T, ar apiResponse, expectedCode string) {
	t.Helper()
	if ar.Success {
		t.Fatal("expected success=false")
	}
	if ar.Data != nil {
		t.Fatal("expected no data on error")
	}
	if ar.Error == nil {
		t.Fatal("expected error object")
	}
	if ar.Error["code"] != expectedCode {
		t.Fatalf("expected error code %q, got %q", expectedCode, ar.Error["code"])
	}
}

// ---------------------------------------------------------------------------
// Tests: Full Household Lifecycle via HTTP
// ---------------------------------------------------------------------------

func TestHouseholdLifecycle_FullHTTP(t *testing.T) {
	srv, repo := newTestServer(t)
	userID := "lifecycle-head"
	adminID := "lifecycle-admin"
	memberID := "lifecycle-member"

	// 1. CREATE household (Draft)
	body := fmt.Sprintf(`{"name":"Lifecycle Home","household_type":"Couple","head_of_household_id":"%s","currency":"USD","country":"US","tags":["test","lifecycle"]}`, userID)
	resp := doRequest(t, srv, "POST", "/api/v1/households?user_id="+userID, []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var createResult command.HouseholdResult
	mustUnmarshalData(t, ar.Data, &createResult)
	if !createResult.Success {
		t.Fatal("expected success=true in result")
	}
	if createResult.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", createResult.Status)
	}
	householdID := createResult.HouseholdID
	if householdID == "" {
		t.Fatal("expected non-empty household_id")
	}

	// 2. GET (verify Draft)
	resp = doRequest(t, srv, "GET", "/api/v1/households/"+householdID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var view query.HouseholdDetailView
	mustUnmarshalData(t, ar.Data, &view)
	if view.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", view.Status)
	}
	if view.HouseholdID != householdID {
		t.Fatalf("expected household_id %s, got %s", householdID, view.HouseholdID)
	}
	if view.Name != "Lifecycle Home" {
		t.Fatalf("expected 'Lifecycle Home', got %s", view.Name)
	}
	if view.HouseholdType != "Couple" {
		t.Fatalf("expected Couple, got %s", view.HouseholdType)
	}
	if view.HeadOfHouseholdID != userID {
		t.Fatalf("expected head %s, got %s", userID, view.HeadOfHouseholdID)
	}
	if view.Currency != "USD" {
		t.Fatalf("expected USD, got %s", view.Currency)
	}
	if view.Country != "US" {
		t.Fatalf("expected US, got %s", view.Country)
	}
	if len(view.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(view.Tags))
	}
	if view.MemberCount == 0 {
		t.Fatal("expected at least 1 member")
	}

	// 3. ACTIVATE
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+householdID+"/activate?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var transition command.HouseholdResult
	mustUnmarshalData(t, ar.Data, &transition)
	if transition.Status != "Active" {
		t.Fatalf("expected Active, got %s", transition.Status)
	}

	// 4. GET (verify Active)
	resp = doRequest(t, srv, "GET", "/api/v1/households/"+householdID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Active" {
		t.Fatalf("expected Active, got %s", view.Status)
	}

	// 5. PAUSE
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+householdID+"/pause?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Paused" {
		t.Fatalf("expected Paused, got %s", transition.Status)
	}

	// 6. GET (verify Paused)
	resp = doRequest(t, srv, "GET", "/api/v1/households/"+householdID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Paused" {
		t.Fatalf("expected Paused, got %s", view.Status)
	}

	// 7. RESUME
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+householdID+"/resume?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Active" {
		t.Fatalf("expected Active, got %s", transition.Status)
	}

	// 8. DISSOLVE
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+householdID+"/dissolve?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Dissolved" {
		t.Fatalf("expected Dissolved, got %s", transition.Status)
	}

	// 9. ARCHIVE
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+householdID+"/archive?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Archived" {
		t.Fatalf("expected Archived, got %s", transition.Status)
	}

	// 10. Verify final state in repo
	repo.mu.Lock()
	archived, exists := repo.households[householdID]
	repo.mu.Unlock()
	if !exists {
		t.Fatal("household should exist in repo")
	}
	if string(archived.Status()) != "Archived" {
		t.Fatalf("expected Archived in repo, got %s", archived.Status())
	}

	// 11. Member management: invite -> accept
	// Reset: create a new active household
	body2 := fmt.Sprintf(`{"name":"Member Test","household_type":"Single","head_of_household_id":"%s","currency":"USD"}`, userID)
	resp = doRequest(t, srv, "POST", "/api/v1/households?user_id="+userID, []byte(body2))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)
	mustUnmarshalData(t, ar.Data, &createResult)
	memberHHID := createResult.HouseholdID

	// Activate it so we can invite
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/activate?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	// Invite member
	inviteBody := fmt.Sprintf(`{"user_id":"%s","role":"Member"}`, adminID)
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/members?user_id="+userID, []byte(inviteBody))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "GET", "/api/v1/households/"+memberHHID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	foundPending := false
	for _, m := range view.Members {
		if m.UserID == adminID && m.InviteStatus == "Pending" {
			foundPending = true
			break
		}
	}
	if !foundPending {
		t.Fatal("expected pending invite for admin user")
	}

	// Accept invite
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/members/"+adminID+"/accept?user_id="+adminID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "GET", "/api/v1/households/"+memberHHID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	foundAccepted := false
	for _, m := range view.Members {
		if m.UserID == adminID && m.InviteStatus == "Accepted" {
			foundAccepted = true
			break
		}
	}
	if !foundAccepted {
		t.Fatal("expected accepted invite for admin user")
	}

	// Remove member
	resp = doRequest(t, srv, "DELETE", "/api/v1/households/"+memberHHID+"/members/"+adminID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "GET", "/api/v1/households/"+memberHHID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	for _, m := range view.Members {
		if m.UserID == adminID {
			t.Fatal("admin should have been removed")
		}
	}

	// 12. Link/unlink account
	linkAccBody := `{"account_id":"acc-link-1","user_id":"` + userID + `"}`
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/accounts?user_id="+userID, []byte(linkAccBody))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "GET", "/api/v1/households/"+memberHHID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if len(view.LinkedAccounts) != 1 {
		t.Fatalf("expected 1 linked account, got %d", len(view.LinkedAccounts))
	}
	if view.LinkedAccounts[0].AccountID != "acc-link-1" {
		t.Fatalf("expected acc-link-1, got %s", view.LinkedAccounts[0].AccountID)
	}

	// Unlink
	resp = doRequest(t, srv, "DELETE", "/api/v1/households/"+memberHHID+"/accounts/acc-link-1?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "GET", "/api/v1/households/"+memberHHID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if len(view.LinkedAccounts) != 0 {
		t.Fatalf("expected 0 linked accounts, got %d", len(view.LinkedAccounts))
	}

	// 13. Link/unlink goal
	linkGoalBody := `{"goal_id":"goal-link-1","user_id":"` + userID + `"}`
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/goals?user_id="+userID, []byte(linkGoalBody))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "DELETE", "/api/v1/households/"+memberHHID+"/goals/goal-link-1?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	// 14. Link/unlink budget
	linkBudgetBody := `{"budget_id":"budget-link-1","user_id":"` + userID + `"}`
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/budgets?user_id="+userID, []byte(linkBudgetBody))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "DELETE", "/api/v1/households/"+memberHHID+"/budgets/budget-link-1?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	// 15. Role management
	// Add a member, then update their role
	inviteBody2 := fmt.Sprintf(`{"user_id":"%s","role":"Member"}`, memberID)
	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/members?user_id="+userID, []byte(inviteBody2))
	assertStatus(t, resp.StatusCode, http.StatusOK)

	resp = doRequest(t, srv, "POST", "/api/v1/households/"+memberHHID+"/members/"+memberID+"/accept?user_id="+memberID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	roleBody := `{"role":"Admin"}`
	resp = doRequest(t, srv, "PUT", "/api/v1/households/"+memberHHID+"/members/"+memberID+"/role?user_id="+userID, []byte(roleBody))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	resp = doRequest(t, srv, "GET", "/api/v1/households/"+memberHHID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	for _, m := range view.Members {
		if m.UserID == memberID && m.Role != "Admin" {
			t.Fatalf("expected member role Admin, got %s", m.Role)
		}
	}

	// 16. Get household summary
	resp = doRequest(t, srv, "GET", "/api/v1/households/"+memberHHID+"/summary?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var summary query.HouseholdFinancialSummary
	mustUnmarshalData(t, ar.Data, &summary)
	if summary.HouseholdID != memberHHID {
		t.Fatalf("expected household_id %s, got %s", memberHHID, summary.HouseholdID)
	}
}

// ---------------------------------------------------------------------------
// Tests: Create Household
// ---------------------------------------------------------------------------

func TestCreateHousehold_StoredInRepo(t *testing.T) {
	srv, repo := newTestServer(t)
	userID := "stored-head"

	body := fmt.Sprintf(`{"name":"Repo Check","household_type":"Family","head_of_household_id":"%s","currency":"EUR","country":"EU","tags":["check"]}`, userID)
	resp := doRequest(t, srv, "POST", "/api/v1/households?user_id="+userID, []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var result command.HouseholdResult
	mustUnmarshalData(t, ar.Data, &result)

	repo.mu.Lock()
	stored, exists := repo.households[result.HouseholdID]
	repo.mu.Unlock()
	if !exists {
		t.Fatal("household should be stored in repo")
	}
	if stored.Name() != "Repo Check" {
		t.Fatalf("expected 'Repo Check', got %s", stored.Name())
	}
	if stored.HeadOfHouseholdID() != userID {
		t.Fatalf("expected %s, got %s", userID, stored.HeadOfHouseholdID())
	}
	if stored.Currency() != "EUR" {
		t.Fatalf("expected EUR, got %s", stored.Currency())
	}
	if string(stored.HouseholdType()) != "Family" {
		t.Fatalf("expected Family, got %s", stored.HouseholdType())
	}
	if string(stored.Status()) != "Draft" {
		t.Fatalf("expected Draft, got %s", stored.Status())
	}
	if len(stored.Tags()) != 1 || stored.Tags()[0] != "check" {
		t.Fatalf("expected tags=[check], got %v", stored.Tags())
	}
}

func TestCreateHousehold_DefaultCurrency(t *testing.T) {
	srv, _ := newTestServer(t)
	userID := "currency-head"

	body := fmt.Sprintf(`{"name":"Currency Test","household_type":"Single","head_of_household_id":"%s","country":"IN"}`, userID)
	resp := doRequest(t, srv, "POST", "/api/v1/households?user_id="+userID, []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var result command.HouseholdResult
	mustUnmarshalData(t, ar.Data, &result)

	resp = doRequest(t, srv, "GET", "/api/v1/households/"+result.HouseholdID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var view query.HouseholdDetailView
	mustUnmarshalData(t, ar.Data, &view)
	if view.Currency != "INR" {
		t.Fatalf("expected default currency INR, got %s", view.Currency)
	}
}

// ---------------------------------------------------------------------------
// Tests: Error Scenarios
// ---------------------------------------------------------------------------

func TestCreateHousehold_InvalidJSON(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "POST", "/api/v1/households?user_id=err-user", []byte(`not json`))
	assertStatus(t, resp.StatusCode, http.StatusBadRequest)
	ar := parseResponse(t, resp)
	assertError(t, ar, "VALIDATION_ERROR")
}

func TestCreateHousehold_EmptyBody(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "POST", "/api/v1/households?user_id=err-user", []byte(`{}`))
	assertStatus(t, resp.StatusCode, http.StatusInternalServerError)
	ar := parseResponse(t, resp)
	assertError(t, ar, "CREATE_ERROR")
}

func TestGetHousehold_NotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/households/non-existent-id?user_id=err-user", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	ar := parseResponse(t, resp)
	assertError(t, ar, "NOT_FOUND")
}

func TestGetHouseholdSummary_NotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/households/non-existent-id/summary?user_id=err-user", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	ar := parseResponse(t, resp)
	assertError(t, ar, "NOT_FOUND")
}

// ---------------------------------------------------------------------------
// Tests: Invalid State Transitions via HTTP
// ---------------------------------------------------------------------------

func TestHousehold_InvalidTransitionsViaHTTP(t *testing.T) {
	tests := []struct {
		name       string
		setupState func(*domain.Household)
		endpoint   string
	}{
		{name: "activate_active", setupState: func(h *domain.Household) { domain.ActivateHousehold(h) }, endpoint: "activate"},
		{name: "pause_draft", setupState: func(h *domain.Household) {}, endpoint: "pause"},
		{name: "resume_draft", setupState: func(h *domain.Household) {}, endpoint: "resume"},
		{name: "resume_active", setupState: func(h *domain.Household) { domain.ActivateHousehold(h) }, endpoint: "resume"},
		{name: "dissolve_draft", setupState: func(h *domain.Household) {}, endpoint: "dissolve"},
		{name: "archive_draft", setupState: func(h *domain.Household) {}, endpoint: "archive"},
		{name: "archive_active", setupState: func(h *domain.Household) { domain.ActivateHousehold(h) }, endpoint: "archive"},
		{name: "pause_dissolved", setupState: func(h *domain.Household) { domain.ActivateHousehold(h); domain.DissolveHousehold(h) }, endpoint: "pause"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv, repo := newTestServer(t)
			hh := makeDomainHousehold(t)
			tc.setupState(hh)
			repo.mu.Lock()
			repo.households[hh.ID()] = hh
			repo.mu.Unlock()

			resp := doRequest(t, srv, "POST", "/api/v1/households/"+hh.ID()+"/"+tc.endpoint+"?user_id=user-head", nil)
			assertStatus(t, resp.StatusCode, http.StatusInternalServerError)
			ar := parseResponse(t, resp)
			assertError(t, ar, "TRANSITION_ERROR")
		})
	}
}

// ---------------------------------------------------------------------------
// Tests: List Households
// ---------------------------------------------------------------------------

func TestListHouseholds_Empty(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/households?user_id=empty-user", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var page query.PaginatedResult
	mustUnmarshalData(t, ar.Data, &page)
	if len(page.Households) != 0 {
		t.Fatalf("expected 0 households, got %d", len(page.Households))
	}
}

func TestListHouseholds_Multiple(t *testing.T) {
	srv, repo := newTestServer(t)
	userID := "list-user"

	svc := application.NewHouseholdService(repo, &mockPublisher{}, time.Now)
	for i := 0; i < 2; i++ {
		body := fmt.Sprintf(`{"name":"List HH %d","household_type":"Single","head_of_household_id":"%s","currency":"USD"}`, i+1, userID)
		req, _ := http.NewRequest("POST", "/api/v1/households?user_id="+userID, bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h := newHandler(svc)
		h.createHousehold(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("creating household %d: status %d", i, w.Code)
		}
	}

	resp := doRequest(t, srv, "GET", "/api/v1/households?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var page query.PaginatedResult
	mustUnmarshalData(t, ar.Data, &page)
	if len(page.Households) != 2 {
		t.Fatalf("expected 2 households, got %d", len(page.Households))
	}
}

// ---------------------------------------------------------------------------
// Tests: Response Format Consistency
// ---------------------------------------------------------------------------

func TestResponseFormat_SuccessEndpoints(t *testing.T) {
	srv, repo := newTestServer(t)
	hh := makeDomainHousehold(t)
	domain.ActivateHousehold(hh)
	domain.PauseHousehold(hh)
	repo.mu.Lock()
	repo.households[hh.ID()] = hh
	repo.mu.Unlock()

	tests := []struct {
		name   string
		method string
		path   string
		body   []byte
		code   int
	}{
		{name: "get household", method: "GET", path: "/api/v1/households/" + hh.ID() + "?user_id=user-head", code: 200},
		{name: "list households", method: "GET", path: "/api/v1/households?user_id=user-head", code: 200},
		{name: "resume", method: "POST", path: "/api/v1/households/" + hh.ID() + "/resume?user_id=user-head", code: 200},
		{name: "summary", method: "GET", path: "/api/v1/households/" + hh.ID() + "/summary?user_id=user-head", code: 200},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, srv, tc.method, tc.path, tc.body)
			assertStatus(t, resp.StatusCode, tc.code)
			ar := parseResponse(t, resp)
			assertSuccess(t, ar)
		})
	}
}

func TestResponseFormat_ErrorEndpoints(t *testing.T) {
	srv, _ := newTestServer(t)

	tests := []struct {
		name      string
		method    string
		path      string
		body      []byte
		code      int
		errorCode string
	}{
		{name: "not found", method: "GET", path: "/api/v1/households/does-not-exist?user_id=u", code: 404, errorCode: "NOT_FOUND"},
		{name: "bad json", method: "POST", path: "/api/v1/households?user_id=u", body: []byte(`{bad`), code: 400, errorCode: "VALIDATION_ERROR"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, srv, tc.method, tc.path, tc.body)
			assertStatus(t, resp.StatusCode, tc.code)
			ar := parseResponse(t, resp)
			assertError(t, ar, tc.errorCode)
		})
	}
}

// ---------------------------------------------------------------------------
// Tests: Accept Invite Forbidden
// ---------------------------------------------------------------------------

func TestAcceptInvite_Forbidden(t *testing.T) {
	srv, repo := newTestServer(t)
	hh := makeDomainHousehold(t)
	repo.mu.Lock()
	repo.households[hh.ID()] = hh
	repo.mu.Unlock()

	// Try to accept someone else's invite
	resp := doRequest(t, srv, "POST", "/api/v1/households/"+hh.ID()+"/members/user-head/accept?user_id=user-other", nil)
	assertStatus(t, resp.StatusCode, http.StatusForbidden)
	ar := parseResponse(t, resp)
	assertError(t, ar, "FORBIDDEN")
}

// ---------------------------------------------------------------------------
// Helper: create a domain Household with known ID
// ---------------------------------------------------------------------------

func makeDomainHousehold(t *testing.T) *domain.Household {
	t.Helper()
	headID := "user-head"
	members := []domain.HouseholdMember{
		{UserID: headID, Role: domain.RoleHead, AddedAt: time.Now().UTC().Format(time.RFC3339), InviteStatus: "Accepted"},
		{UserID: "user-admin", Role: domain.RoleAdmin, AddedAt: time.Now().UTC().Format(time.RFC3339), InviteStatus: "Accepted"},
		{UserID: "user-member", Role: domain.RoleMember, AddedAt: time.Now().UTC().Format(time.RFC3339), InviteStatus: "Accepted"},
	}
	now := time.Now().UTC()
	return domain.ReconstructFromDB(
		"test-hh-id", "Test Household", domain.HHCouple, headID, members,
		domain.HHStatusDraft, "USD", "US",
		0, 0, 0, domain.HHHealthy,
		[]string{}, "", nil, nil, nil, nil,
		now, now,
	)
}

// Ensure domain.Repository interface is satisfied at compile time
var _ domain.Repository = (*mockRepo)(nil)
