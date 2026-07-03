package register

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/horizon/core/services/domains/budget/internal/application"
	"github.com/horizon/core/services/domains/budget/internal/application/dto/command"
	"github.com/horizon/core/services/domains/budget/internal/application/dto/query"
	"github.com/horizon/core/services/domains/budget/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock repository
// ---------------------------------------------------------------------------

type mockRepo struct {
	mu          sync.Mutex
	budgets     map[string]*domain.Budget
	failMethods map[string]bool
}

func newMockRepo() *mockRepo {
	return &mockRepo{budgets: make(map[string]*domain.Budget), failMethods: make(map[string]bool)}
}

func (r *mockRepo) setFail(method string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.failMethods[method] = true
}

func (r *mockRepo) shouldFail(method string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.failMethods[method]
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

func (r *mockRepo) ListByUser(_ context.Context, userID string, cursor string, limit int) ([]*domain.Budget, string, error) {
	r.mu.Lock()
	shouldFail := r.failMethods["ListByUser"]
	r.mu.Unlock()
	if shouldFail {
		return nil, "", fmt.Errorf("db error")
	}

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

func (r *mockRepo) ListByPeriod(_ context.Context, userID string, start, end time.Time, cursor string, limit int) ([]*domain.Budget, string, error) {
	return nil, "", nil
}

func (r *mockRepo) GetByCategory(_ context.Context, userID, category string, cursor string, limit int) ([]*domain.Budget, string, error) {
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

func (r *mockRepo) ListByHousehold(_ context.Context, householdID string, cursor string, limit int) ([]*domain.Budget, string, error) {
	r.mu.Lock()
	shouldFail := r.failMethods["ListByHousehold"]
	r.mu.Unlock()
	if shouldFail {
		return nil, "", fmt.Errorf("db error")
	}

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
// Mock publisher
// ---------------------------------------------------------------------------

type mockPublisher struct{}

func (m *mockPublisher) Publish(_ domain.DomainEvent) error { return nil }

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func newTestHarness(t *testing.T) (*http.ServeMux, *mockRepo, *budgetHandler) {
	t.Helper()
	repo := newMockRepo()
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
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
	return mux, repo, h
}

func newTestHarnessWithBudget(t *testing.T, budget *domain.Budget) (*http.ServeMux, *mockRepo, *budgetHandler) {
	t.Helper()
	mux, repo, h := newTestHarness(t)
	repo.mu.Lock()
	repo.budgets[budget.ID()] = budget
	repo.mu.Unlock()
	return mux, repo, h
}

func execRequest(mux *http.ServeMux, method, target string, body []byte) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, target, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func execRequestWithUser(mux *http.ServeMux, method, target, userID string, body []byte) *httptest.ResponseRecorder {
	url := target
	if userID != "" {
		sep := "?"
		for _, c := range target {
			if c == '?' {
				sep = "&"
				break
			}
		}
		url = target + sep + "user_id=" + userID
	}
	return execRequest(mux, method, url, body)
}

func execRequestWithHouseholdRole(mux *http.ServeMux, method, target, householdID, role string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	if body != nil {
		req.Body = http.NoBody
		req = httptest.NewRequest(method, target, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	if role != "" {
		req.Header.Set("X-Household-Role-"+householdID, role)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

// responseBody is a generic container for JSON response parsing.
type responseBody struct {
	Success  bool                   `json:"success"`
	Data     json.RawMessage        `json:"data,omitempty"`
	Error    map[string]string      `json:"error,omitempty"`
	Metadata map[string]string      `json:"metadata,omitempty"`
}

func decodeResponse(t *testing.T, body *bytes.Buffer) responseBody {
	t.Helper()
	var resp responseBody
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

func makeTestBudget(t *testing.T) *domain.Budget {
	t.Helper()
	f := domain.NewBudgetFactory()
	now := time.Now()
	b, err := f.Create("budget-1", "user-1", "", "Test Budget", domain.PeriodMonthly, now, now.AddDate(0, 1, 0), "INR", []domain.BudgetCategory{
		{ID: "cat-1", Category: "Food", BudgetedAmount: 10000},
		{ID: "cat-2", Category: "Transport", BudgetedAmount: 5000},
	}, nil)
	if err != nil {
		t.Fatalf("makeTestBudget: %v", err)
	}
	return b
}

// ---------------------------------------------------------------------------
// Tests: Create Budget
// ---------------------------------------------------------------------------

func TestCreateBudget_Success(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	body := `{"name":"Monthly Groceries","period":"Monthly","start_date":"2025-01-01","end_date":"2025-01-31","currency":"INR","categories":[{"category":"Food","budgeted_amount":15000}],"tags":["groceries"]}`
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets", "user-1", []byte(body))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	if resp.Data == nil {
		t.Fatal("expected data in response")
	}
	var result command.BudgetResult
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if result.BudgetID == "" {
		t.Error("expected non-empty budget_id")
	}
	if result.Status != "Draft" {
		t.Errorf("expected Draft, got %s", result.Status)
	}
	if !result.Success {
		t.Error("expected success=true on result")
	}
	// Verify default user ID
	w2 := execRequest(mux, "GET", "/api/v1/budgets/"+result.BudgetID+"?user_id=user-1", nil)
	if w2.Code != 200 {
		t.Errorf("expected 200 on follow-up get, got %d", w2.Code)
	}
}

func TestCreateBudget_DefaultUser(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	body := `{"name":"Test 2","period":"Monthly","start_date":"2025-01-01","end_date":"2025-01-31","currency":"INR"}`
	w := execRequest(mux, "POST", "/api/v1/budgets?user_id=user-2", []byte(body))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
}

func TestCreateBudget_InvalidBody(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets", "user-1", []byte(`not json`))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Success {
		t.Fatal("expected success=false")
	}
	if resp.Error["code"] != "VALIDATION_ERROR" {
		t.Errorf("expected VALIDATION_ERROR, got %v", resp.Error["code"])
	}
}

func TestCreateBudget_EmptyBody(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets", "user-1", []byte(`{}`))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "CREATE_ERROR" {
		t.Errorf("expected CREATE_ERROR, got %v", resp.Error["code"])
	}
}

// ---------------------------------------------------------------------------
// Tests: Get Budget
// ---------------------------------------------------------------------------

func TestGetBudget_Success(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "GET", "/api/v1/budgets/"+b.ID(), "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	var view query.BudgetView
	if err := json.Unmarshal(resp.Data, &view); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if view.BudgetID != b.ID() {
		t.Errorf("expected budget_id %s, got %s", b.ID(), view.BudgetID)
	}
	if view.Name != "Test Budget" {
		t.Errorf("expected 'Test Budget', got '%s'", view.Name)
	}
	if view.Status != "Draft" {
		t.Errorf("expected Draft, got %s", view.Status)
	}
	if len(view.Categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(view.Categories))
	}
}

func TestGetBudget_MissingID(t *testing.T) {
	_, _, h := newTestHarness(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.getBudget(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "MISSING_ID" {
		t.Errorf("expected MISSING_ID, got %v", resp.Error["code"])
	}
}

func TestGetBudget_NotFound(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	w := execRequestWithUser(mux, "GET", "/api/v1/budgets/non-existent", "user-1", nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND, got %v", resp.Error["code"])
	}
}

// ---------------------------------------------------------------------------
// Tests: List Budgets
// ---------------------------------------------------------------------------

func TestListBudgets_Success(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "GET", "/api/v1/budgets", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	var page query.PaginatedResult
	if err := json.Unmarshal(resp.Data, &page); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if len(page.Budgets) != 1 {
		t.Errorf("expected 1 budget, got %d", len(page.Budgets))
	}
	if page.Budgets[0].BudgetID != b.ID() {
		t.Errorf("expected budget_id %s, got %s", b.ID(), page.Budgets[0].BudgetID)
	}
}

func TestListBudgets_Empty(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	w := execRequestWithUser(mux, "GET", "/api/v1/budgets", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	var page query.PaginatedResult
	if err := json.Unmarshal(resp.Data, &page); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if len(page.Budgets) != 0 {
		t.Errorf("expected 0 budgets, got %d", len(page.Budgets))
	}
}

func TestListBudgets_WithPagination(t *testing.T) {
	repo := newMockRepo()
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}

	for i := 0; i < 3; i++ {
		_, err := svc.Create(context.Background(), command.CreateBudgetCommand{
			UserID:    "user-1",
			Name:      fmt.Sprintf("Budget %d", i+1),
			Period:    "Monthly",
			StartDate: time.Now().Format("2006-01-02"),
			EndDate:   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			Currency:  "INR",
		})
		if err != nil {
			t.Fatalf("svc.Create %d: %v", i, err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/budgets", h.listBudgets)

	w := execRequestWithUser(mux, "GET", "/api/v1/budgets?limit=2", "user-1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	var page query.PaginatedResult
	if err := json.Unmarshal(resp.Data, &page); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if len(page.Budgets) < 2 {
		t.Errorf("expected at least 2 budgets, got %d", len(page.Budgets))
	}
}

// ---------------------------------------------------------------------------
// Tests: List By Household
// ---------------------------------------------------------------------------

func TestListByHousehold_Success(t *testing.T) {
	f := domain.NewBudgetFactory()
	now := time.Now()
	b, err := f.Create("bh-1", "user-1", "hh-1", "Household Budget",
		domain.PeriodMonthly, now, now.AddDate(0, 1, 0), "INR", nil, nil)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	repo := newMockRepo()
	repo.mu.Lock()
	repo.budgets["bh-1"] = b
	repo.mu.Unlock()
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/households/{id}/budgets", h.listByHousehold)

	w := execRequestWithHouseholdRole(mux, "GET", "/api/v1/households/hh-1/budgets", "hh-1", "admin", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	var page query.PaginatedResult
	if err := json.Unmarshal(resp.Data, &page); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if len(page.Budgets) != 1 {
		t.Errorf("expected 1 budget, got %d", len(page.Budgets))
	}
	if page.Budgets[0].BudgetID != "bh-1" {
		t.Errorf("expected bh-1, got %s", page.Budgets[0].BudgetID)
	}
}

func TestListByHousehold_MissingID(t *testing.T) {
	_, _, h := newTestHarness(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.listByHousehold(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "MISSING_ID" {
		t.Errorf("expected MISSING_ID, got %v", resp.Error["code"])
	}
}

func TestListByHousehold_Forbidden(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	w := execRequest(mux, "GET", "/api/v1/households/hh-1/budgets", nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "FORBIDDEN" {
		t.Errorf("expected FORBIDDEN, got %v", resp.Error["code"])
	}
}

func TestListByHousehold_Empty(t *testing.T) {
	repo := newMockRepo()
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/households/{id}/budgets", h.listByHousehold)

	w := execRequestWithHouseholdRole(mux, "GET", "/api/v1/households/hh-empty/budgets", "hh-empty", "viewer", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	var page query.PaginatedResult
	if err := json.Unmarshal(resp.Data, &page); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if len(page.Budgets) != 0 {
		t.Errorf("expected 0 budgets, got %d", len(page.Budgets))
	}
}

// ---------------------------------------------------------------------------
// Tests: State transitions (activate, pause, resume, complete, archive)
// ---------------------------------------------------------------------------

func TestActivateBudget_Success(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/activate", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	var result command.BudgetResult
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.Status != "Active" {
		t.Errorf("expected Active, got %s", result.Status)
	}
}

func TestActivateBudget_MissingID(t *testing.T) {
	_, _, h := newTestHarness(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	h.activateBudget(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "MISSING_ID" {
		t.Errorf("expected MISSING_ID, got %v", resp.Error["code"])
	}
}

func TestActivateBudget_InvalidTransition(t *testing.T) {
	b := makeTestBudget(t)
	// Activate first
	_ = domain.ActivateBudget(b)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/activate", "user-1", nil)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "TRANSITION_ERROR" {
		t.Errorf("expected TRANSITION_ERROR, got %v", resp.Error["code"])
	}
}

func TestPauseBudget_Success(t *testing.T) {
	b := makeTestBudget(t)
	_ = domain.ActivateBudget(b)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/pause", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result command.BudgetResult
	json.Unmarshal(decodeResponse(t, w.Body).Data, &result)
	if result.Status != "Paused" {
		t.Errorf("expected Paused, got %s", result.Status)
	}
}

func TestPauseBudget_MissingID(t *testing.T) {
	_, _, h := newTestHarness(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	h.pauseBudget(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestResumeBudget_Success(t *testing.T) {
	b := makeTestBudget(t)
	_ = domain.ActivateBudget(b)
	_ = domain.PauseBudget(b)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/resume", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result command.BudgetResult
	json.Unmarshal(decodeResponse(t, w.Body).Data, &result)
	if result.Status != "Active" {
		t.Errorf("expected Active, got %s", result.Status)
	}
}

func TestResumeBudget_InvalidTransition(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/resume", "user-1", nil)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestCompleteBudget_Success(t *testing.T) {
	b := makeTestBudget(t)
	_ = domain.ActivateBudget(b)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/complete", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result command.BudgetResult
	json.Unmarshal(decodeResponse(t, w.Body).Data, &result)
	if result.Status != "Completed" {
		t.Errorf("expected Completed, got %s", result.Status)
	}
}

func TestCompleteBudget_InvalidTransition(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/complete", "user-1", nil)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestArchiveBudget_Success(t *testing.T) {
	b := makeTestBudget(t)
	_ = domain.ActivateBudget(b)
	_ = domain.CompleteBudget(b)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/archive", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result command.BudgetResult
	json.Unmarshal(decodeResponse(t, w.Body).Data, &result)
	if result.Status != "Archived" {
		t.Errorf("expected Archived, got %s", result.Status)
	}
}

func TestArchiveBudget_InvalidTransition(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "POST", "/api/v1/budgets/"+b.ID()+"/archive", "user-1", nil)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Tests: Budget vs Actual
// ---------------------------------------------------------------------------

func TestBudgetVsActual_Success(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "GET", "/api/v1/budgets/"+b.ID()+"/vs-actual", "user-1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
}

func TestBudgetVsActual_MissingID(t *testing.T) {
	_, _, h := newTestHarness(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.budgetVsActual(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "MISSING_ID" {
		t.Errorf("expected MISSING_ID, got %v", resp.Error["code"])
	}
}

func TestBudgetVsActual_NotFound(t *testing.T) {
	mux, _, _ := newTestHarness(t)
	w := execRequestWithUser(mux, "GET", "/api/v1/budgets/non-existent/vs-actual", "user-1", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND, got %v", resp.Error["code"])
	}
}

// ---------------------------------------------------------------------------
// Tests: Update Category
// ---------------------------------------------------------------------------

func TestUpdateCategory_Success(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	body := `{"category_id":"cat-1","budgeted_amount":12000}`
	w := execRequestWithUser(mux, "PUT", "/api/v1/budgets/"+b.ID()+"/category", "user-1", []byte(body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	var result command.BudgetResult
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.BudgetID != b.ID() {
		t.Errorf("expected budget_id %s, got %s", b.ID(), result.BudgetID)
	}
}

func TestUpdateCategory_InvalidBody(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	w := execRequestWithUser(mux, "PUT", "/api/v1/budgets/"+b.ID()+"/category", "user-1", []byte(`not json`))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "VALIDATION_ERROR" {
		t.Errorf("expected VALIDATION_ERROR, got %v", resp.Error["code"])
	}
}

func TestUpdateCategory_MissingID(t *testing.T) {
	_, _, h := newTestHarness(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/", nil)
	h.updateCategory(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateCategory_SpendAmount(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	body := `{"category_id":"cat-1","budgeted_amount":10000,"spent_amount":2500}`
	w := execRequestWithUser(mux, "PUT", "/api/v1/budgets/"+b.ID()+"/category", "user-1", []byte(body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if !resp.Success {
		t.Fatal("expected success=true")
	}
}

// ---------------------------------------------------------------------------
// Tests: Helper functions (writeJSON, writeError, okData)
// ---------------------------------------------------------------------------

func TestWriteJSON_SetsContentType(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"hello": "world"})
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", w.Header().Get("Content-Type"))
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["hello"] != "world" {
		t.Errorf("expected 'world', got '%s'", body["hello"])
	}
}

func TestWriteJSON_WithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusCreated, map[string]int{"count": 42})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestWriteError_Format(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Success {
		t.Error("expected success=false")
	}
	if resp.Error == nil {
		t.Fatal("expected error map")
	}
	if resp.Error["code"] != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND, got %v", resp.Error["code"])
	}
	if resp.Error["message"] != "Resource not found" {
		t.Errorf("expected 'Resource not found', got '%s'", resp.Error["message"])
	}
	if resp.Metadata == nil {
		t.Fatal("expected metadata")
	}
	if _, ok := resp.Metadata["timestamp"]; !ok {
		t.Error("expected timestamp in metadata")
	}
}

func TestWriteError_VariousCodes(t *testing.T) {
	tests := []struct {
		status  int
		code    string
		message string
	}{
		{http.StatusBadRequest, "VALIDATION_ERROR", "bad input"},
		{http.StatusForbidden, "FORBIDDEN", "not allowed"},
		{http.StatusInternalServerError, "DB_ERROR", "connection failed"},
	}
	for _, tc := range tests {
		t.Run(tc.code, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeError(w, tc.status, tc.code, tc.message)
			if w.Code != tc.status {
				t.Errorf("expected %d, got %d", tc.status, w.Code)
			}
			resp := decodeResponse(t, w.Body)
			if resp.Error["code"] != tc.code {
				t.Errorf("expected %s, got %s", tc.code, resp.Error["code"])
			}
			if resp.Error["message"] != tc.message {
				t.Errorf("expected %s, got %s", tc.message, resp.Error["message"])
			}
		})
	}
}

func TestOkData_Structure(t *testing.T) {
	data := okData("hello")
	if data["success"] != true {
		t.Error("expected success=true")
	}
	if data["data"] != "hello" {
		t.Errorf("expected 'hello', got '%v'", data["data"])
	}
	meta, ok := data["metadata"].(map[string]string)
	if !ok {
		t.Fatal("expected metadata map")
	}
	if _, exists := meta["timestamp"]; !exists {
		t.Error("expected timestamp in metadata")
	}
}

func TestOkData_WithStruct(t *testing.T) {
	type testStruct struct{ ID string; Value int }
	d := okData(testStruct{ID: "abc", Value: 42})
	inner, ok := d["data"].(testStruct)
	if !ok {
		t.Fatal("expected testStruct data")
	}
	if inner.ID != "abc" {
		t.Errorf("expected abc, got %s", inner.ID)
	}
}

// ---------------------------------------------------------------------------
// Tests: Response format consistency across all endpoints
// ---------------------------------------------------------------------------

func TestResponseFormat_AllEndpoints(t *testing.T) {
	b := makeTestBudget(t)
	mux, repo, _ := newTestHarnessWithBudget(t, b)

	// Activate the budget for state-based endpoints
	_ = domain.ActivateBudget(b)
	repo.mu.Lock()
	repo.budgets[b.ID()] = b
	repo.mu.Unlock()

	type endpoint struct {
		method string
		path   string
		body   []byte
		code   int
		userID string
	}

	endpoints := []endpoint{
		// Success cases
		{method: "GET", path: "/api/v1/budgets/" + b.ID(), code: 200, userID: "user-1"},
		{method: "GET", path: "/api/v1/budgets", code: 200, userID: "user-1"},
		{method: "POST", path: "/api/v1/budgets/" + b.ID() + "/pause", code: 200, userID: "user-1"},
		// Error cases
		{method: "GET", path: "/api/v1/budgets/bad-id", code: 404, userID: "user-1"},
		{method: "POST", path: "/api/v1/budgets/" + b.ID() + "/activate", code: 500, userID: "user-1"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := execRequestWithUser(mux, ep.method, ep.path, ep.userID, ep.body)
			if w.Code != ep.code {
				t.Fatalf("expected %d, got %d", ep.code, w.Code)
			}
			resp := decodeResponse(t, w.Body)
			if resp.Metadata == nil {
				t.Error("missing metadata")
			} else if _, ok := resp.Metadata["timestamp"]; !ok {
				t.Error("metadata missing timestamp")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Tests: Full lifecycle integration
// ---------------------------------------------------------------------------

func TestBudgetFullLifecycle(t *testing.T) {
	f := domain.NewBudgetFactory()
	now := time.Now()
	b, err := f.Create("lifecycle-1", "user-life", "", "Lifecycle Budget",
		domain.PeriodMonthly, now, now.AddDate(0, 1, 0), "INR", nil, nil)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	repo := newMockRepo()
	repo.mu.Lock()
	repo.budgets["lifecycle-1"] = b
	repo.mu.Unlock()
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/budgets/{id}/activate", h.activateBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/pause", h.pauseBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/resume", h.resumeBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/complete", h.completeBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/archive", h.archiveBudget)

	activate := func() *command.BudgetResult {
		w := execRequestWithUser(mux, "POST", "/api/v1/budgets/lifecycle-1/activate", "user-life", nil)
		if w.Code != 200 {
			t.Fatalf("activate: %s", w.Body.String())
		}
		var r command.BudgetResult
		json.Unmarshal(decodeResponse(t, w.Body).Data, &r)
		return &r
	}
	pause := func() *command.BudgetResult {
		w := execRequestWithUser(mux, "POST", "/api/v1/budgets/lifecycle-1/pause", "user-life", nil)
		if w.Code != 200 {
			t.Fatalf("pause: %s", w.Body.String())
		}
		var r command.BudgetResult
		json.Unmarshal(decodeResponse(t, w.Body).Data, &r)
		return &r
	}
	resume := func() *command.BudgetResult {
		w := execRequestWithUser(mux, "POST", "/api/v1/budgets/lifecycle-1/resume", "user-life", nil)
		if w.Code != 200 {
			t.Fatalf("resume: %s", w.Body.String())
		}
		var r command.BudgetResult
		json.Unmarshal(decodeResponse(t, w.Body).Data, &r)
		return &r
	}
	complete := func() *command.BudgetResult {
		w := execRequestWithUser(mux, "POST", "/api/v1/budgets/lifecycle-1/complete", "user-life", nil)
		if w.Code != 200 {
			t.Fatalf("complete: %s", w.Body.String())
		}
		var r command.BudgetResult
		json.Unmarshal(decodeResponse(t, w.Body).Data, &r)
		return &r
	}
	archive := func() *command.BudgetResult {
		w := execRequestWithUser(mux, "POST", "/api/v1/budgets/lifecycle-1/archive", "user-life", nil)
		if w.Code != 200 {
			t.Fatalf("archive: %s", w.Body.String())
		}
		var r command.BudgetResult
		json.Unmarshal(decodeResponse(t, w.Body).Data, &r)
		return &r
	}

	// Draft -> Active -> Paused -> Active -> Completed -> Archived
	if r := activate(); r.Status != "Active" {
		t.Errorf("expected Active, got %s", r.Status)
	}
	if r := pause(); r.Status != "Paused" {
		t.Errorf("expected Paused, got %s", r.Status)
	}
	if r := resume(); r.Status != "Active" {
		t.Errorf("expected Active, got %s", r.Status)
	}
	if r := complete(); r.Status != "Completed" {
		t.Errorf("expected Completed, got %s", r.Status)
	}
	if r := archive(); r.Status != "Archived" {
		t.Errorf("expected Archived, got %s", r.Status)
	}
}

// ---------------------------------------------------------------------------
// Tests: Missing ID across all path-value handlers
// ---------------------------------------------------------------------------

func TestMissingID_AllHandlers(t *testing.T) {
	_, _, h := newTestHarness(t)

	type handlerCase struct {
		name string
		fn   func(http.ResponseWriter, *http.Request)
	}

	cases := []handlerCase{
		{"getBudget", h.getBudget},
		{"activateBudget", h.activateBudget},
		{"pauseBudget", h.pauseBudget},
		{"resumeBudget", h.resumeBudget},
		{"completeBudget", h.completeBudget},
		{"archiveBudget", h.archiveBudget},
		{"budgetVsActual", h.budgetVsActual},
		{"updateCategory", h.updateCategory},
		{"listByHousehold", h.listByHousehold},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			c.fn(w, r)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Tests: Auth extraction (user_id query param)
// ---------------------------------------------------------------------------

func TestCreateBudget_SetsUserID(t *testing.T) {
	repo := newMockRepo()
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/budgets", h.createBudget)

	body := `{"name":"Auth Test","period":"Monthly","start_date":"2025-01-01","end_date":"2025-01-31","currency":"INR"}`
	w := execRequest(mux, "POST", "/api/v1/budgets?user_id=auth-user", []byte(body))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// Verify the saved budget has the correct user ID
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for id, b := range repo.budgets {
		if b.Name() == "Auth Test" {
			if b.UserID() != "auth-user" {
				t.Errorf("expected auth-user, got %s", b.UserID())
			}
			_ = id
			return
		}
	}
	t.Error("budget not found in repo")
}

func TestListBudgets_RepoError(t *testing.T) {
	repo := newMockRepo()
	repo.setFail("ListByUser")
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/budgets", h.listBudgets)

	w := execRequestWithUser(mux, "GET", "/api/v1/budgets", "user-1", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "DB_ERROR" {
		t.Errorf("expected DB_ERROR, got %v", resp.Error["code"])
	}
}

func TestListByHousehold_RepoError(t *testing.T) {
	repo := newMockRepo()
	repo.setFail("ListByHousehold")
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/households/{id}/budgets", h.listByHousehold)

	w := execRequestWithHouseholdRole(mux, "GET", "/api/v1/households/hh-1/budgets", "hh-1", "admin", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "DB_ERROR" {
		t.Errorf("expected DB_ERROR, got %v", resp.Error["code"])
	}
}

func TestUpdateCategory_CategoryNotFound(t *testing.T) {
	b := makeTestBudget(t)
	mux, _, _ := newTestHarnessWithBudget(t, b)
	body := `{"category_id":"nonexistent","budgeted_amount":5000}`
	w := execRequestWithUser(mux, "PUT", "/api/v1/budgets/"+b.ID()+"/category", "user-1", []byte(body))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
	resp := decodeResponse(t, w.Body)
	if resp.Error["code"] != "UPDATE_ERROR" {
		t.Errorf("expected UPDATE_ERROR, got %v", resp.Error["code"])
	}
}

func TestDefaultUserID(t *testing.T) {
	repo := newMockRepo()
	pub := &mockPublisher{}
	svc := application.NewBudgetService(repo, pub, time.Now)
	h := &budgetHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/budgets", h.createBudget)

	body := `{"name":"No Auth","period":"Monthly","start_date":"2025-01-01","end_date":"2025-01-31","currency":"INR"}`
	w := execRequest(mux, "POST", "/api/v1/budgets", []byte(body))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, b := range repo.budgets {
		if b.Name() == "No Auth" {
			if b.UserID() != "default" {
				t.Errorf("expected 'default', got '%s'", b.UserID())
			}
			return
		}
	}
	t.Error("budget not found")
}
