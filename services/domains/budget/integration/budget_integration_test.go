package budget_test

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

	"github.com/horizon/core/services/domains/budget/internal/application"
	"github.com/horizon/core/services/domains/budget/internal/application/dto/command"
	"github.com/horizon/core/services/domains/budget/internal/application/dto/query"
	"github.com/horizon/core/services/domains/budget/internal/domain"
	"github.com/horizon/core/services/internal/auth"
)

// ---------------------------------------------------------------------------
// Mock Repository
// ---------------------------------------------------------------------------

type mockRepo struct {
	mu      sync.Mutex
	budgets map[string]*domain.Budget
}

func newMockRepo() *mockRepo {
	return &mockRepo{budgets: make(map[string]*domain.Budget)}
}

func (r *mockRepo) Save(_ context.Context, b *domain.Budget) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := b.ID()
	if id == "" {
		id = fmt.Sprintf("budget-%d", len(r.budgets)+1)
		now := time.Now().UTC()
		rebuilt := domain.ReconstructFromDB(
			id, b.UserID(), b.HouseholdID(), b.Name(), b.Period(),
			b.StartDate(), b.EndDate(), b.Status(),
			b.TotalBudgeted(), b.TotalSpent(), b.TotalRemaining(),
			b.Currency(), b.Categories(), b.Tags(), now, now,
		)
		*b = *rebuilt
	}
	r.budgets[b.ID()] = b
	return nil
}

func (r *mockRepo) UpdateStatus(_ context.Context, id string, from, to domain.BudgetStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.budgets[id]
	if !ok {
		return fmt.Errorf("budget not found")
	}
	if b.Status() != from {
		return fmt.Errorf("status mismatch")
	}
	b.SetStatus(to)
	return nil
}

func (r *mockRepo) GetByID(_ context.Context, id string) (*domain.Budget, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.budgets[id]
	if !ok {
		return nil, fmt.Errorf("budget not found")
	}
	return b, nil
}

func (r *mockRepo) ListByUser(_ context.Context, userID string, _ string, _ int) ([]*domain.Budget, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Budget
	for _, b := range r.budgets {
		if b.UserID() == userID {
			result = append(result, b)
		}
	}
	return result, "", nil
}

func (r *mockRepo) ListByPeriod(_ context.Context, _ string, _, _ time.Time, _ string, _ int) ([]*domain.Budget, string, error) {
	return nil, "", nil
}

func (r *mockRepo) GetByCategory(_ context.Context, _, _ string, _ string, _ int) ([]*domain.Budget, string, error) {
	return nil, "", nil
}

func (r *mockRepo) GetBudgetVsActual(_ context.Context, budgetID string) ([]domain.BudgetCategory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.budgets[budgetID]
	if !ok {
		return nil, fmt.Errorf("budget not found")
	}
	return b.Categories(), nil
}

func (r *mockRepo) ListByHousehold(_ context.Context, householdID string, _ string, _ int) ([]*domain.Budget, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Budget
	for _, b := range r.budgets {
		if b.HouseholdID() == householdID {
			result = append(result, b)
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
// HTTP handlers (test double mirroring register/budgetHandler)
// ---------------------------------------------------------------------------

type handler struct {
	svc *application.BudgetService
}

func newHandler(svc *application.BudgetService) *handler {
	return &handler{svc: svc}
}

func (h *handler) createBudget(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) getBudget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Budget ID required")
		return
	}
	result, err := h.svc.GetByID(r.Context(), query.GetBudgetQuery{BudgetID: id})
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Budget not found")
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) listBudgets(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) activateBudget(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Activate(ctx, cmd)
	})
}

func (h *handler) pauseBudget(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Pause(ctx, cmd)
	})
}

func (h *handler) resumeBudget(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Resume(ctx, cmd)
	})
}

func (h *handler) completeBudget(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Complete(ctx, cmd)
	})
}

func (h *handler) archiveBudget(w http.ResponseWriter, r *http.Request) {
	h.handleTransition(w, r, func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
		return h.svc.Archive(ctx, cmd)
	})
}

type transitionFn func(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error)

func (h *handler) handleTransition(w http.ResponseWriter, r *http.Request, fn transitionFn) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Budget ID required")
		return
	}
	result, err := fn(r.Context(), command.BudgetIDCommand{BudgetID: id})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "TRANSITION_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) updateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cmd command.UpdateCategoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	cmd.BudgetID = id
	result, err := h.svc.UpdateCategory(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "UPDATE_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) budgetVsActual(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Budget ID required")
		return
	}
	result, err := h.svc.GetByID(r.Context(), query.GetBudgetQuery{BudgetID: id})
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Budget not found")
		return
	}
	writeJSON(w, http.StatusOK, okData(result))
}

func (h *handler) listByHousehold(w http.ResponseWriter, r *http.Request) {
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
	limit := 25
	result, err := h.svc.ListByHousehold(r.Context(), query.ListByHouseholdQuery{HouseholdID: id, Cursor: cursor, Limit: limit})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
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
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := newHandler(svc)

	mux := http.NewServeMux()
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
// Tests: Complete Budget Lifecycle via HTTP
// ---------------------------------------------------------------------------

func TestBudgetLifecycle_FullHTTP(t *testing.T) {
	srv, repo := newTestServer(t)
	userID := "lifecycle-user"

	// 1. CREATE
	body := `{"name":"Lifecycle Budget","period":"Monthly","start_date":"2025-06-01","end_date":"2025-06-30","currency":"INR","categories":[{"category":"Food","budgeted_amount":10000},{"category":"Transport","budgeted_amount":5000}]}`
	resp := doRequest(t, srv, "POST", "/api/v1/budgets?user_id="+userID, []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var createResult command.BudgetResult
	mustUnmarshalData(t, ar.Data, &createResult)
	if !createResult.Success {
		t.Fatal("expected success=true in result")
	}
	if createResult.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", createResult.Status)
	}
	budgetID := createResult.BudgetID
	if budgetID == "" {
		t.Fatal("expected non-empty budget_id")
	}

	// 2. GET (verify Draft)
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var view query.BudgetView
	mustUnmarshalData(t, ar.Data, &view)
	if view.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", view.Status)
	}
	if view.BudgetID != budgetID {
		t.Fatalf("expected budget_id %s, got %s", budgetID, view.BudgetID)
	}
	if view.Name != "Lifecycle Budget" {
		t.Fatalf("expected 'Lifecycle Budget', got %s", view.Name)
	}
	if view.Period != "Monthly" {
		t.Fatalf("expected Monthly, got %s", view.Period)
	}
	if view.Currency != "INR" {
		t.Fatalf("expected INR, got %s", view.Currency)
	}
	if len(view.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(view.Categories))
	}
	if view.TotalBudgeted != 15000 {
		t.Fatalf("expected total_budgeted=15000, got %d", view.TotalBudgeted)
	}

	// 3. ACTIVATE
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/activate?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var transition command.BudgetResult
	mustUnmarshalData(t, ar.Data, &transition)
	if transition.Status != "Active" {
		t.Fatalf("expected Active, got %s", transition.Status)
	}

	// 4. GET (verify Active)
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Active" {
		t.Fatalf("expected Active, got %s", view.Status)
	}

	// 5. PAUSE
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/pause?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Paused" {
		t.Fatalf("expected Paused, got %s", transition.Status)
	}

	// 6. GET (verify Paused)
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Paused" {
		t.Fatalf("expected Paused, got %s", view.Status)
	}

	// 7. RESUME
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/resume?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Active" {
		t.Fatalf("expected Active, got %s", transition.Status)
	}

	// 8. COMPLETE
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/complete?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Completed" {
		t.Fatalf("expected Completed, got %s", transition.Status)
	}

	// 9. GET (verify Completed)
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Completed" {
		t.Fatalf("expected Completed, got %s", view.Status)
	}

	// 10. ARCHIVE
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/archive?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &transition)
	if transition.Status != "Archived" {
		t.Fatalf("expected Archived, got %s", transition.Status)
	}

	// 11. GET (verify Archived)
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshalData(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Archived" {
		t.Fatalf("expected Archived, got %s", view.Status)
	}

	// 12. Verify the repo has the budget with final archived state
	repo.mu.Lock()
	archived, exists := repo.budgets[budgetID]
	repo.mu.Unlock()
	if !exists {
		t.Fatal("budget should exist in repo")
	}
	if archived.Status() != domain.BStatusArchived {
		t.Fatalf("expected Archived in repo, got %s", archived.Status())
	}
}

// ---------------------------------------------------------------------------
// Tests: Create Budget
// ---------------------------------------------------------------------------

func TestCreateBudget_WithMultipleCategories(t *testing.T) {
	srv, _ := newTestServer(t)
	body := `{"name":"Category Test","period":"Monthly","start_date":"2025-07-01","end_date":"2025-07-31","currency":"USD","categories":[
		{"category":"Food","budgeted_amount":20000,"rollover":true},
		{"category":"Transport","budgeted_amount":8000},
		{"category":"Entertainment","budgeted_amount":5000,"rollover":false}
	],"tags":["test","categories"]}`
	resp := doRequest(t, srv, "POST", "/api/v1/budgets?user_id=cat-user", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var result command.BudgetResult
	mustUnmarshalData(t, ar.Data, &result)
	if result.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", result.Status)
	}

	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+result.BudgetID+"?user_id=cat-user", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var view query.BudgetView
	mustUnmarshalData(t, ar.Data, &view)
	if len(view.Categories) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(view.Categories))
	}
	catMap := make(map[string]query.CategoryView)
	for _, c := range view.Categories {
		catMap[c.Category] = c
	}

	food, ok := catMap["Food"]
	if !ok {
		t.Fatal("expected Food category")
	}
	if food.BudgetedAmount != 20000 {
		t.Fatalf("expected Food budget 20000, got %d", food.BudgetedAmount)
	}
	if !food.Rollover {
		t.Fatal("expected Food rollover=true")
	}

	transport, ok := catMap["Transport"]
	if !ok {
		t.Fatal("expected Transport category")
	}
	if transport.BudgetedAmount != 8000 {
		t.Fatalf("expected Transport budget 8000, got %d", transport.BudgetedAmount)
	}
	if transport.Rollover {
		t.Fatal("expected Transport rollover=false")
	}

	entertainment, ok := catMap["Entertainment"]
	if !ok {
		t.Fatal("expected Entertainment category")
	}
	if entertainment.BudgetedAmount != 5000 {
		t.Fatalf("expected Entertainment budget 5000, got %d", entertainment.BudgetedAmount)
	}

	if view.TotalBudgeted != 33000 {
		t.Fatalf("expected total_budgeted=33000, got %d", view.TotalBudgeted)
	}
	if view.TotalRemaining != 33000 {
		t.Fatalf("expected total_remaining=33000, got %d", view.TotalRemaining)
	}
	if len(view.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(view.Tags))
	}
}

func TestCreateBudget_DefaultValues(t *testing.T) {
	srv, _ := newTestServer(t)
	body := `{"name":"Defaults","period":"Monthly","start_date":"2025-08-01","end_date":"2025-08-31"}`
	resp := doRequest(t, srv, "POST", "/api/v1/budgets", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var result command.BudgetResult
	mustUnmarshalData(t, ar.Data, &result)
	if result.BudgetID == "" {
		t.Fatal("expected budget_id")
	}

	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+result.BudgetID+"?user_id=default", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var view query.BudgetView
	mustUnmarshalData(t, ar.Data, &view)
	if view.Currency != "INR" {
		t.Fatalf("expected default currency INR, got %s", view.Currency)
	}
}

// ---------------------------------------------------------------------------
// Tests: Error Scenarios
// ---------------------------------------------------------------------------

func TestCreateBudget_InvalidJSON(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "POST", "/api/v1/budgets?user_id=err-user", []byte(`not json`))
	assertStatus(t, resp.StatusCode, http.StatusBadRequest)
	ar := parseResponse(t, resp)
	assertError(t, ar, "VALIDATION_ERROR")
}

func TestCreateBudget_EmptyBody(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "POST", "/api/v1/budgets?user_id=err-user", []byte(`{}`))
	assertStatus(t, resp.StatusCode, http.StatusInternalServerError)
	ar := parseResponse(t, resp)
	assertError(t, ar, "CREATE_ERROR")
}

func TestCreateBudget_EmptyName(t *testing.T) {
	srv, _ := newTestServer(t)
	body := `{"name":"","period":"Monthly","start_date":"2025-01-01","end_date":"2025-01-31","currency":"INR"}`
	resp := doRequest(t, srv, "POST", "/api/v1/budgets?user_id=err-user", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusInternalServerError)
	ar := parseResponse(t, resp)
	assertError(t, ar, "CREATE_ERROR")
}



func TestGetBudget_NotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/budgets/non-existent-id?user_id=err-user", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	ar := parseResponse(t, resp)
	assertError(t, ar, "NOT_FOUND")
}

func TestBudgetVsActual_NotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/budgets/bad-id/vs-actual?user_id=err-user", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
	ar := parseResponse(t, resp)
	assertError(t, ar, "NOT_FOUND")
}

func TestUpdateCategory_InvalidBody(t *testing.T) {
	srv, repo := newTestServer(t)

	// Create a budget first
	budget := makeDomainBudget(t)
	repo.mu.Lock()
	repo.budgets[budget.ID()] = budget
	repo.mu.Unlock()

	resp := doRequest(t, srv, "PUT", "/api/v1/budgets/"+budget.ID()+"/category?user_id=u", []byte(`not json`))
	assertStatus(t, resp.StatusCode, http.StatusBadRequest)
	ar := parseResponse(t, resp)
	assertError(t, ar, "VALIDATION_ERROR")
}

func TestUpdateCategory_CategoryNotFound(t *testing.T) {
	srv, repo := newTestServer(t)
	budget := makeDomainBudget(t)
	repo.mu.Lock()
	repo.budgets[budget.ID()] = budget
	repo.mu.Unlock()

	body := `{"category_id":"nonexistent","budgeted_amount":5000}`
	resp := doRequest(t, srv, "PUT", "/api/v1/budgets/"+budget.ID()+"/category?user_id=u", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusInternalServerError)
	ar := parseResponse(t, resp)
	assertError(t, ar, "UPDATE_ERROR")
}

// ---------------------------------------------------------------------------
// Tests: Invalid State Transitions (table-driven)
// ---------------------------------------------------------------------------

func TestBudget_InvalidTransitionsViaHTTP(t *testing.T) {
	tests := []struct {
		name       string
		setupState func(*domain.Budget)
		endpoint   string // relative to /api/v1/budgets/{id}/
	}{
		{name: "activate_active", setupState: func(b *domain.Budget) { domain.ActivateBudget(b) }, endpoint: "activate"},
		{name: "pause_draft", setupState: func(b *domain.Budget) {}, endpoint: "pause"},
		{name: "resume_draft", setupState: func(b *domain.Budget) {}, endpoint: "resume"},
		{name: "resume_active", setupState: func(b *domain.Budget) { domain.ActivateBudget(b) }, endpoint: "resume"},
		{name: "complete_draft", setupState: func(b *domain.Budget) {}, endpoint: "complete"},
		{name: "archive_draft", setupState: func(b *domain.Budget) {}, endpoint: "archive"},
		{name: "archive_active", setupState: func(b *domain.Budget) { domain.ActivateBudget(b) }, endpoint: "archive"},
		{name: "pause_completed", setupState: func(b *domain.Budget) { domain.ActivateBudget(b); domain.CompleteBudget(b) }, endpoint: "pause"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv, repo := newTestServer(t)
			budget := makeDomainBudget(t)
			tc.setupState(budget)
			repo.mu.Lock()
			repo.budgets[budget.ID()] = budget
			repo.mu.Unlock()

			resp := doRequest(t, srv, "POST", "/api/v1/budgets/"+budget.ID()+"/"+tc.endpoint+"?user_id=u", nil)
			assertStatus(t, resp.StatusCode, http.StatusInternalServerError)
			ar := parseResponse(t, resp)
			assertError(t, ar, "TRANSITION_ERROR")
		})
	}
}

// ---------------------------------------------------------------------------
// Tests: List Budgets
// ---------------------------------------------------------------------------

func TestListBudgets_Empty(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/budgets?user_id=empty-user", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var page query.PaginatedResult
	mustUnmarshalData(t, ar.Data, &page)
	if len(page.Budgets) != 0 {
		t.Fatalf("expected 0 budgets, got %d", len(page.Budgets))
	}
}

func TestListBudgets_Multiple(t *testing.T) {
	srv, repo := newTestServer(t)
	userID := "list-user"

	// Create two budgets via the service
	f := domain.NewBudgetFactory()
	now := time.Now()
	for i := 0; i < 2; i++ {
		b, err := f.Create(
			fmt.Sprintf("list-budget-%d", i+1), userID, "",
			fmt.Sprintf("List Budget %d", i+1),
			domain.PeriodMonthly, now, now.AddDate(0, 1, 0),
			"INR", nil, nil,
		)
		if err != nil {
			t.Fatalf("creating budget %d: %v", i, err)
		}
		repo.mu.Lock()
		repo.budgets[b.ID()] = b
		repo.mu.Unlock()
	}

	resp := doRequest(t, srv, "GET", "/api/v1/budgets?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var page query.PaginatedResult
	mustUnmarshalData(t, ar.Data, &page)
	if len(page.Budgets) != 2 {
		t.Fatalf("expected 2 budgets, got %d", len(page.Budgets))
	}
}

// ---------------------------------------------------------------------------
// Tests: List by Household
// ---------------------------------------------------------------------------

func TestListByHousehold_WithRole(t *testing.T) {
	srv, repo := newTestServer(t)
	hhID := "hh-integration"

	f := domain.NewBudgetFactory()
	now := time.Now()
	b, err := f.Create("hh-budget-1", "u1", hhID, "HH Budget", domain.PeriodMonthly, now, now.AddDate(0, 1, 0), "INR", nil, nil)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	repo.mu.Lock()
	repo.budgets[b.ID()] = b
	repo.mu.Unlock()

	req, _ := http.NewRequest("GET", srv.URL+"/api/v1/households/"+hhID+"/budgets", nil)
	req.Header.Set("X-Household-Role-"+hhID, "admin")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var page query.PaginatedResult
	mustUnmarshalData(t, ar.Data, &page)
	if len(page.Budgets) != 1 {
		t.Fatalf("expected 1 budget, got %d", len(page.Budgets))
	}
	if page.Budgets[0].BudgetID != b.ID() {
		t.Fatalf("expected budget_id %s, got %s", b.ID(), page.Budgets[0].BudgetID)
	}
}

func TestListByHousehold_Forbidden(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/households/hh-no-access/budgets", nil)
	assertStatus(t, resp.StatusCode, http.StatusForbidden)
	ar := parseResponse(t, resp)
	assertError(t, ar, "FORBIDDEN")
}

func TestListByHousehold_Empty(t *testing.T) {
	srv, _ := newTestServer(t)
	hhID := "hh-empty"
	req, _ := http.NewRequest("GET", srv.URL+"/api/v1/households/"+hhID+"/budgets", nil)
	req.Header.Set("X-Household-Role-"+hhID, "viewer")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var page query.PaginatedResult
	mustUnmarshalData(t, ar.Data, &page)
	if len(page.Budgets) != 0 {
		t.Fatalf("expected 0 budgets, got %d", len(page.Budgets))
	}
}

// ---------------------------------------------------------------------------
// Tests: Update Category
// ---------------------------------------------------------------------------

func TestUpdateCategory_BudgetedAmount(t *testing.T) {
	srv, repo := newTestServer(t)
	budget := makeDomainBudget(t)
	repo.mu.Lock()
	repo.budgets[budget.ID()] = budget
	repo.mu.Unlock()

	catID := budget.Categories()[0].ID
	body := fmt.Sprintf(`{"category_id":"%s","budgeted_amount":12000}`, catID)
	resp := doRequest(t, srv, "PUT", "/api/v1/budgets/"+budget.ID()+"/category?user_id=u", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var result command.BudgetResult
	mustUnmarshalData(t, ar.Data, &result)
	if result.BudgetID != budget.ID() {
		t.Fatalf("expected budget_id %s, got %s", budget.ID(), result.BudgetID)
	}

	// Verify the update persisted in the repo
	repo.mu.Lock()
	updated := repo.budgets[budget.ID()]
	repo.mu.Unlock()
	if updated.Categories()[0].BudgetedAmount != 12000 {
		t.Fatalf("expected budgeted_amount 12000, got %d", updated.Categories()[0].BudgetedAmount)
	}
}

func TestUpdateCategory_SpentAmount(t *testing.T) {
	srv, repo := newTestServer(t)
	budget := makeDomainBudget(t)
	repo.mu.Lock()
	repo.budgets[budget.ID()] = budget
	repo.mu.Unlock()

	catID := budget.Categories()[0].ID
	body := fmt.Sprintf(`{"category_id":"%s","budgeted_amount":10000,"spent_amount":2500}`, catID)
	resp := doRequest(t, srv, "PUT", "/api/v1/budgets/"+budget.ID()+"/category?user_id=u", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	repo.mu.Lock()
	updated := repo.budgets[budget.ID()]
	repo.mu.Unlock()
	cat := updated.Categories()[0]
	if cat.SpentAmount != 2500 {
		t.Fatalf("expected spent_amount 2500, got %d", cat.SpentAmount)
	}
	// 10000 - 2500 = 7500
	if cat.RemainingAmount != 7500 {
		t.Fatalf("expected remaining_amount 7500, got %d", cat.RemainingAmount)
	}
}

// ---------------------------------------------------------------------------
// Tests: Budget vs Actual
// ---------------------------------------------------------------------------

func TestBudgetVsActual_Success(t *testing.T) {
	srv, repo := newTestServer(t)
	budget := makeDomainBudget(t)
	repo.mu.Lock()
	repo.budgets[budget.ID()] = budget
	repo.mu.Unlock()

	resp := doRequest(t, srv, "GET", "/api/v1/budgets/"+budget.ID()+"/vs-actual?user_id=u", nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)
}

// ---------------------------------------------------------------------------
// Tests: Response Format Consistency
// ---------------------------------------------------------------------------

func TestResponseFormat_SuccessEndpoints(t *testing.T) {
	srv, repo := newTestServer(t)
	budget := makeDomainBudget(t)
	domain.ActivateBudget(budget)
	domain.PauseBudget(budget)
	repo.mu.Lock()
	repo.budgets[budget.ID()] = budget
	repo.mu.Unlock()

	tests := []struct {
		name   string
		method string
		path   string
		body   []byte
		code   int
	}{
		{name: "get budget", method: "GET", path: "/api/v1/budgets/" + budget.ID() + "?user_id=u", code: 200},
		{name: "list budgets", method: "GET", path: "/api/v1/budgets?user_id=u", code: 200},
		{name: "resume", method: "POST", path: "/api/v1/budgets/" + budget.ID() + "/resume?user_id=u", code: 200},
		{name: "vs-actual", method: "GET", path: "/api/v1/budgets/" + budget.ID() + "/vs-actual?user_id=u", code: 200},
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
		name        string
		method      string
		path        string
		body        []byte
		code        int
		errorCode   string
	}{
		{name: "not found", method: "GET", path: "/api/v1/budgets/does-not-exist?user_id=u", code: 404, errorCode: "NOT_FOUND"},
		{name: "bad json", method: "POST", path: "/api/v1/budgets?user_id=u", body: []byte(`{bad`), code: 400, errorCode: "VALIDATION_ERROR"},
		{name: "forbidden", method: "GET", path: "/api/v1/households/hh-x/budgets", code: 403, errorCode: "FORBIDDEN"},
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
// Test: Full Create via HTTP then verify in repo
// ---------------------------------------------------------------------------

func TestCreateBudget_StoredInRepo(t *testing.T) {
	srv, repo := newTestServer(t)
	body := `{"name":"Repo Check","period":"Monthly","start_date":"2025-09-01","end_date":"2025-09-30","currency":"EUR","categories":[{"category":"Food","budgeted_amount":10000}],"tags":["check"]}`
	resp := doRequest(t, srv, "POST", "/api/v1/budgets?user_id=repo-user", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var result command.BudgetResult
	mustUnmarshalData(t, ar.Data, &result)

	repo.mu.Lock()
	stored, exists := repo.budgets[result.BudgetID]
	repo.mu.Unlock()
	if !exists {
		t.Fatal("budget should be stored in repo")
	}
	if stored.Name() != "Repo Check" {
		t.Fatalf("expected 'Repo Check', got %s", stored.Name())
	}
	if stored.UserID() != "repo-user" {
		t.Fatalf("expected repo-user, got %s", stored.UserID())
	}
	if stored.Currency() != "EUR" {
		t.Fatalf("expected EUR, got %s", stored.Currency())
	}
	if stored.Status() != domain.BStatusDraft {
		t.Fatalf("expected Draft, got %s", stored.Status())
	}
	if stored.TotalBudgeted() != 10000 {
		t.Fatalf("expected 10000, got %d", stored.TotalBudgeted())
	}
	if len(stored.Categories()) != 1 {
		t.Fatalf("expected 1 category, got %d", len(stored.Categories()))
	}
	if len(stored.Tags()) != 1 || stored.Tags()[0] != "check" {
		t.Fatalf("expected tags=[check], got %v", stored.Tags())
	}
}

// ---------------------------------------------------------------------------
// Helper: create a domain Budget with known ID and categories
// ---------------------------------------------------------------------------

func makeDomainBudget(t *testing.T) *domain.Budget {
	t.Helper()
	f := domain.NewBudgetFactory()
	now := time.Now()
	b, err := f.Create(
		"test-budget-id", "test-user", "", "Test Budget",
		domain.PeriodMonthly, now, now.AddDate(0, 1, 0), "INR",
		[]domain.BudgetCategory{
			{ID: "cat-food", Category: "Food", BudgetedAmount: 10000},
			{ID: "cat-transport", Category: "Transport", BudgetedAmount: 5000},
		},
		nil,
	)
	if err != nil {
		t.Fatalf("makeDomainBudget: %v", err)
	}
	return b
}

// Ensure domain.Repository interface is satisfied at compile time
var _ domain.Repository = (*mockRepo)(nil)
