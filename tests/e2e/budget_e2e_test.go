package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/horizon/core/services/domains/budget/register"
)

// ---------------------------------------------------------------------------
// Response types (defined locally to avoid importing internal DTO packages)
// ---------------------------------------------------------------------------

type apiResponse struct {
	Success  bool              `json:"success"`
	Data     json.RawMessage   `json:"data,omitempty"`
	Error    map[string]string `json:"error,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type budgetResult struct {
	BudgetID string `json:"budget_id"`
	Status   string `json:"status"`
	Success  bool   `json:"success"`
}

type categoryView struct {
	ID              string  `json:"id"`
	Category        string  `json:"category"`
	BudgetedAmount  int64   `json:"budgeted_amount"`
	SpentAmount     int64   `json:"spent_amount"`
	RemainingAmount int64   `json:"remaining_amount"`
	Rollover        bool    `json:"rollover"`
	SpentPct        float64 `json:"spent_pct"`
}

type budgetView struct {
	BudgetID       string         `json:"budget_id"`
	Name           string         `json:"name"`
	Period         string         `json:"period"`
	Status         string         `json:"status"`
	TotalBudgeted  int64          `json:"total_budgeted"`
	TotalSpent     int64          `json:"total_spent"`
	TotalRemaining int64          `json:"total_remaining"`
	Currency       string         `json:"currency"`
	Categories     []categoryView `json:"categories"`
	Tags           []string       `json:"tags"`
}

type paginatedResult struct {
	Budgets    []budgetView `json:"budgets"`
	NextCursor string       `json:"next_cursor,omitempty"`
	HasMore    bool         `json:"has_more"`
}

// ---------------------------------------------------------------------------
// HTTP helpers
// ---------------------------------------------------------------------------

func newBudgetTestServer(t *testing.T) (*httptest.Server, *register.E2ETestHarness) {
	t.Helper()
	harness := register.NewE2ETestHarness()
	mux := http.NewServeMux()
	harness.RegisterRoutesOnMux(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, harness
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

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected status %d, got %d", want, got)
	}
}

func assertSuccess(t *testing.T, ar apiResponse) {
	t.Helper()
	if !ar.Success {
		t.Fatalf("expected success=true, got error: %v", ar.Error)
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

func mustUnmarshal(t *testing.T, raw json.RawMessage, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(raw, target); err != nil {
		t.Fatalf("unmarshalling data: %v", err)
	}
}

// ---------------------------------------------------------------------------
// E2E Scenario 1: Monthly Budget Management
// ---------------------------------------------------------------------------

func TestE2E_MonthlyBudgetManagement(t *testing.T) {
	srv, harness := newBudgetTestServer(t)
	userID := "e2e-user"
	base := "?user_id=" + userID

	// ---- Step 1: Create a monthly budget with 3 categories ----
	createBody := `{
		"name":"Monthly Budget",
		"period":"Monthly",
		"start_date":"2026-07-01",
		"end_date":"2026-07-31",
		"currency":"INR",
		"categories":[
			{"category":"Food","budgeted_amount":15000,"rollover":false},
			{"category":"Transport","budgeted_amount":5000,"rollover":false},
			{"category":"Entertainment","budgeted_amount":3000,"rollover":false}
		]
	}`
	resp := doRequest(t, srv, "POST", "/api/v1/budgets"+base, []byte(createBody))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var createResult budgetResult
	mustUnmarshal(t, ar.Data, &createResult)
	if createResult.BudgetID == "" {
		t.Fatal("expected non-empty budget_id")
	}
	if createResult.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", createResult.Status)
	}
	if !createResult.Success {
		t.Fatal("expected success=true in result")
	}
	budgetID := createResult.BudgetID

	// ---- Step 2: Get the budget and verify categories ----
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var view budgetView
	mustUnmarshal(t, ar.Data, &view)

	if view.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", view.Status)
	}
	if view.Name != "Monthly Budget" {
		t.Fatalf("expected 'Monthly Budget', got %s", view.Name)
	}
	if view.Period != "Monthly" {
		t.Fatalf("expected Monthly, got %s", view.Period)
	}
	if view.Currency != "INR" {
		t.Fatalf("expected INR, got %s", view.Currency)
	}
	if view.TotalBudgeted != 23000 {
		t.Fatalf("expected total_budgeted=23000, got %d", view.TotalBudgeted)
	}
	if view.TotalRemaining != 23000 {
		t.Fatalf("expected total_remaining=23000, got %d", view.TotalRemaining)
	}
	if view.TotalSpent != 0 {
		t.Fatalf("expected total_spent=0, got %d", view.TotalSpent)
	}
	if len(view.Categories) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(view.Categories))
	}

	catMap := make(map[string]categoryView)
	for _, c := range view.Categories {
		catMap[c.Category] = c
	}

	food, ok := catMap["Food"]
	if !ok {
		t.Fatal("expected Food category")
	}
	if food.BudgetedAmount != 15000 {
		t.Fatalf("expected Food budgeted=15000, got %d", food.BudgetedAmount)
	}
	if food.SpentAmount != 0 {
		t.Fatalf("expected Food spent=0, got %d", food.SpentAmount)
	}
	if food.Rollover {
		t.Fatal("expected Food rollover=false")
	}

	transport, ok := catMap["Transport"]
	if !ok {
		t.Fatal("expected Transport category")
	}
	if transport.BudgetedAmount != 5000 {
		t.Fatalf("expected Transport budgeted=5000, got %d", transport.BudgetedAmount)
	}

	entertainment, ok := catMap["Entertainment"]
	if !ok {
		t.Fatal("expected Entertainment category")
	}
	if entertainment.BudgetedAmount != 3000 {
		t.Fatalf("expected Entertainment budgeted=3000, got %d", entertainment.BudgetedAmount)
	}

	// Save category IDs for later updates
	foodCatID := food.ID

	// ---- Step 3: Activate the budget ----
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/activate"+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var transitionResult budgetResult
	mustUnmarshal(t, ar.Data, &transitionResult)
	if transitionResult.Status != "Active" {
		t.Fatalf("expected Active, got %s", transitionResult.Status)
	}

	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshal(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Active" {
		t.Fatalf("expected Active after activation, got %s", view.Status)
	}

	// ---- Step 4: Update Food category spent amount ----
	if foodCatID == "" {
		t.Fatal("Food category ID not found")
	}

	updateBody := fmt.Sprintf(`{"category_id":"%s","spent_amount":8000}`, foodCatID)
	resp = doRequest(t, srv, "PUT", "/api/v1/budgets/"+budgetID+"/category"+base, []byte(updateBody))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var updateResult budgetResult
	mustUnmarshal(t, ar.Data, &updateResult)
	if !updateResult.Success {
		t.Fatal("expected success=true after category update")
	}

	// Verify via GET
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshal(t, parseResponse(t, resp).Data, &view)

	var updatedFood categoryView
	found := false
	for _, c := range view.Categories {
		if c.Category == "Food" {
			updatedFood = c
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Food category not found after update")
	}
	if updatedFood.SpentAmount != 8000 {
		t.Fatalf("expected Food spent=8000, got %d", updatedFood.SpentAmount)
	}
	if updatedFood.RemainingAmount != 7000 {
		t.Fatalf("expected Food remaining=7000, got %d", updatedFood.RemainingAmount)
	}

	// Verify via harness
	if got := harness.GetCategorySpent(budgetID, "Food"); got != 8000 {
		t.Fatalf("harness: expected Food spent=8000, got %d", got)
	}
	if got := harness.GetCategoryRemaining(budgetID, "Food"); got != 7000 {
		t.Fatalf("harness: expected Food remaining=7000, got %d", got)
	}

	// ---- Step 5: List budgets and verify count ----
	resp = doRequest(t, srv, "GET", "/api/v1/budgets"+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	var page paginatedResult
	mustUnmarshal(t, ar.Data, &page)
	if len(page.Budgets) != 1 {
		t.Fatalf("expected 1 budget, got %d", len(page.Budgets))
	}
	if page.Budgets[0].BudgetID != budgetID {
		t.Fatalf("expected budget_id %s, got %s", budgetID, page.Budgets[0].BudgetID)
	}
	if page.Budgets[0].Status != "Active" {
		t.Fatalf("expected Active status in list, got %s", page.Budgets[0].Status)
	}

	// ---- Step 6: Complete the budget ----
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/complete"+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	mustUnmarshal(t, ar.Data, &transitionResult)
	if transitionResult.Status != "Completed" {
		t.Fatalf("expected Completed, got %s", transitionResult.Status)
	}

	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshal(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Completed" {
		t.Fatalf("expected Completed after completion, got %s", view.Status)
	}

	// ---- Step 7: Archive the budget ----
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/archive"+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	assertSuccess(t, ar)

	mustUnmarshal(t, ar.Data, &transitionResult)
	if transitionResult.Status != "Archived" {
		t.Fatalf("expected Archived, got %s", transitionResult.Status)
	}

	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshal(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Archived" {
		t.Fatalf("expected Archived after archiving, got %s", view.Status)
	}

	// Final repo verification
	if !harness.BudgetExists(budgetID) {
		t.Fatal("budget should exist in repo after full lifecycle")
	}
	if harness.GetBudgetStatus(budgetID) != "Archived" {
		t.Fatalf("expected Archived in repo, got %s", harness.GetBudgetStatus(budgetID))
	}
}

// ---------------------------------------------------------------------------
// E2E Scenario 2: Budget Overspending
// ---------------------------------------------------------------------------

func TestE2E_BudgetOverspending(t *testing.T) {
	srv, harness := newBudgetTestServer(t)
	userID := "overspend-user"
	base := "?user_id=" + userID

	// ---- Step 1: Create a budget with a category budgeted 10000 ----
	createBody := `{
		"name":"Overspend Test",
		"period":"Monthly",
		"start_date":"2026-07-01",
		"end_date":"2026-07-31",
		"currency":"INR",
		"categories":[
			{"category":"Groceries","budgeted_amount":10000,"rollover":false}
		]
	}`
	resp := doRequest(t, srv, "POST", "/api/v1/budgets"+base, []byte(createBody))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	assertSuccess(t, ar)

	var createResult budgetResult
	mustUnmarshal(t, ar.Data, &createResult)
	budgetID := createResult.BudgetID
	if budgetID == "" {
		t.Fatal("expected non-empty budget_id")
	}

	// Verify creation
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	var view budgetView
	mustUnmarshal(t, parseResponse(t, resp).Data, &view)
	if len(view.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(view.Categories))
	}
	if view.Categories[0].BudgetedAmount != 10000 {
		t.Fatalf("expected budgeted=10000, got %d", view.Categories[0].BudgetedAmount)
	}
	if view.TotalBudgeted != 10000 {
		t.Fatalf("expected total_budgeted=10000, got %d", view.TotalBudgeted)
	}

	// ---- Step 2: Activate the budget ----
	resp = doRequest(t, srv, "POST", "/api/v1/budgets/"+budgetID+"/activate"+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	assertSuccess(t, parseResponse(t, resp))

	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshal(t, parseResponse(t, resp).Data, &view)
	if view.Status != "Active" {
		t.Fatalf("expected Active, got %s", view.Status)
	}

	// ---- Step 3: Update category with spent amount 15000 (overspent) ----
	catID := view.Categories[0].ID
	if catID == "" {
		t.Fatal("Groceries category ID not found")
	}

	updateBody := fmt.Sprintf(`{"category_id":"%s","spent_amount":15000}`, catID)
	resp = doRequest(t, srv, "PUT", "/api/v1/budgets/"+budgetID+"/category"+base, []byte(updateBody))
	assertStatus(t, resp.StatusCode, http.StatusOK)
	assertSuccess(t, parseResponse(t, resp))

	// ---- Step 4: Verify remaining = 0 (no rollover, overspent) ----
	resp = doRequest(t, srv, "GET", "/api/v1/budgets/"+budgetID+base, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	mustUnmarshal(t, parseResponse(t, resp).Data, &view)

	if len(view.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(view.Categories))
	}
	cat := view.Categories[0]
	if cat.SpentAmount != 15000 {
		t.Fatalf("expected spent=15000, got %d", cat.SpentAmount)
	}
	if cat.RemainingAmount != 0 {
		t.Fatalf("expected remaining=0 (no rollover for overspent), got %d", cat.RemainingAmount)
	}
	if cat.BudgetedAmount != 10000 {
		t.Fatalf("expected budgeted=10000, got %d", cat.BudgetedAmount)
	}

	// ---- Step 5: Check budget health shows overspent ----
	if !harness.BudgetExists(budgetID) {
		t.Fatal("budget should exist in repo")
	}

	totalBudgeted := harness.ComputeTotalBudgetedFromCategories(budgetID)
	totalSpent := harness.ComputeTotalSpentFromCategories(budgetID)
	if totalSpent <= totalBudgeted {
		t.Fatalf("expected total_spent (%d) > total_budgeted (%d) for overspent", totalSpent, totalBudgeted)
	}

	overspent := totalSpent - totalBudgeted
	if overspent != 5000 {
		t.Fatalf("expected overspent amount 5000, got %d", overspent)
	}

	if harness.GetCategoryRemaining(budgetID, "Groceries") != 0 {
		t.Fatalf("expected category remaining 0 in repo")
	}
	if harness.GetCategorySpent(budgetID, "Groceries") != 15000 {
		t.Fatalf("expected category spent 15000 in repo")
	}
}
