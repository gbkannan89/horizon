package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/horizon/core/services/domains/household/register"
)

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

type apiResp struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   interface{}     `json:"error,omitempty"`
}

type householdResult struct {
	HouseholdID string `json:"household_id"`
	Status      string `json:"status"`
	Success     bool   `json:"success"`
}

type householdDetail struct {
	HouseholdID       string              `json:"household_id"`
	Name              string              `json:"name"`
	HouseholdType     string              `json:"household_type"`
	HeadOfHouseholdID string              `json:"head_of_household_id"`
	MemberCount       int                 `json:"member_count"`
	Status            string              `json:"status"`
	Currency          string              `json:"currency"`
	Country           string              `json:"country"`
	TotalAssets       int64               `json:"total_assets"`
	TotalLiabilities  int64               `json:"total_liabilities"`
	TotalNetWorth     int64               `json:"total_net_worth"`
	Health            string              `json:"health"`
	Tags              []string            `json:"tags"`
	Notes             string              `json:"notes"`
	Members           []memberView        `json:"members"`
	LinkedGoals       []linkedView        `json:"linked_goals"`
	LinkedBudgets     []linkedView        `json:"linked_budgets"`
	GoalContributions []contribView       `json:"goal_contributions"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
}

type memberView struct {
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	AddedAt      string `json:"added_at"`
	InviteStatus string `json:"invite_status"`
}

type linkedView struct {
	ID      string `json:"id,omitempty"`
	GoalID  string `json:"goal_id,omitempty"`
	BudgetID string `json:"budget_id,omitempty"`
	AddedBy string `json:"added_by"`
	AddedAt string `json:"added_at"`
}

type contribView struct {
	GoalID string `json:"goal_id"`
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
	Date   string `json:"date"`
}

type financialSummary struct {
	HouseholdID        string `json:"household_id"`
	TotalAssets        int64  `json:"total_assets"`
	TotalLiabilities   int64  `json:"total_liabilities"`
	TotalNetWorth      int64  `json:"total_net_worth"`
	Currency           string `json:"currency"`
	LinkedAccounts     int    `json:"linked_accounts_count"`
	LinkedGoals        int    `json:"linked_goals_count"`
	LinkedBudgets      int    `json:"linked_budgets_count"`
	TotalContributions int64  `json:"total_contributions"`
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()
	harness := register.NewE2ETestHarness()
	mux := http.NewServeMux()
	harness.RegisterRoutesOnMux(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &testServer{Server: srv}
}

func (ts *testServer) do(t *testing.T, method, path string, body interface{}) *http.Response {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req, err := http.NewRequest(method, ts.URL+path, &buf)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func (ts *testServer) get(t *testing.T, path string) *http.Response {
	return ts.do(t, http.MethodGet, path, nil)
}

func (ts *testServer) post(t *testing.T, path string, body interface{}) *http.Response {
	return ts.do(t, http.MethodPost, path, body)
}

func (ts *testServer) put(t *testing.T, path string, body interface{}) *http.Response {
	return ts.do(t, http.MethodPut, path, body)
}

func (ts *testServer) del(t *testing.T, path string) *http.Response {
	return ts.do(t, http.MethodDelete, path, nil)
}

func parseSuccess(t *testing.T, resp *http.Response) json.RawMessage {
	t.Helper()
	defer resp.Body.Close()
	var ar apiResp
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !ar.Success {
		t.Fatalf("API call failed (HTTP %d): %v", resp.StatusCode, ar.Error)
	}
	return ar.Data
}

func mustDecode[T any](t *testing.T, raw json.RawMessage, v *T) {
	t.Helper()
	if err := json.Unmarshal(raw, v); err != nil {
		t.Fatalf("decode %T: %v", v, err)
	}
}

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const (
	headUser   = "user-head"
	memberUser = "user-member"
	goalID     = "goal-family-vacation"
	budgetID   = "budget-monthly"
	hhName     = "The Sharma Family"
)

// ---------------------------------------------------------------------------
// E2E: Family Household Setup
// ---------------------------------------------------------------------------

func TestFamilyHouseholdE2E(t *testing.T) {
	ts := newTestServer(t)
	var hhID string

	// 1. Create household
	t.Run("01_create", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households?user_id="+headUser, map[string]interface{}{
			"name":                 hhName,
			"household_type":       "Family",
			"head_of_household_id": headUser,
			"currency":             "INR",
			"country":              "IN",
			"tags":                 []string{"family", "test"},
			"notes":                "Test household for E2E",
		})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}
		var res householdResult
		mustDecode(t, parseSuccess(t, resp), &res)
		if !res.Success {
			t.Fatal("expected success=true")
		}
		if res.HouseholdID == "" {
			t.Fatal("expected non-empty household_id")
		}
		if res.Status != "Draft" {
			t.Fatalf("expected Draft, got %s", res.Status)
		}
		hhID = res.HouseholdID
	})

	// 2. Get household by ID and verify fields
	t.Run("02_get", func(t *testing.T) {
		resp := ts.get(t, "/api/v1/households/"+hhID)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, resp), &det)

		if det.HouseholdID != hhID {
			t.Fatalf("ID mismatch: %s vs %s", hhID, det.HouseholdID)
		}
		if det.Name != hhName {
			t.Fatalf("name mismatch: %q vs %q", hhName, det.Name)
		}
		if det.HouseholdType != "Family" {
			t.Fatalf("type: expected Family, got %s", det.HouseholdType)
		}
		if det.Currency != "INR" {
			t.Fatalf("currency: expected INR, got %s", det.Currency)
		}
		if det.Country != "IN" {
			t.Fatalf("country: expected IN, got %s", det.Country)
		}
		if det.Status != "Draft" {
			t.Fatalf("status: expected Draft, got %s", det.Status)
		}
		if det.Health != "Healthy" {
			t.Fatalf("health: expected Healthy, got %s", det.Health)
		}
		if det.HeadOfHouseholdID != headUser {
			t.Fatalf("head: expected %s, got %s", headUser, det.HeadOfHouseholdID)
		}
		if det.MemberCount != 1 {
			t.Fatalf("member count: expected 1, got %d", det.MemberCount)
		}
		if len(det.Tags) != 2 || det.Tags[0] != "family" {
			t.Fatalf("tags: %v", det.Tags)
		}
		if det.Notes != "Test household for E2E" {
			t.Fatalf("notes: %s", det.Notes)
		}
		if det.TotalAssets != 0 || det.TotalLiabilities != 0 || det.TotalNetWorth != 0 {
			t.Fatalf("financials: %d/%d/%d", det.TotalAssets, det.TotalLiabilities, det.TotalNetWorth)
		}
		if len(det.Members) != 1 {
			t.Fatalf("expected 1 member, got %d", len(det.Members))
		}
		if det.Members[0].UserID != headUser {
			t.Fatalf("head member: %s", det.Members[0].UserID)
		}
		if det.Members[0].Role != "Head" {
			t.Fatalf("head role: %s", det.Members[0].Role)
		}
		if det.Members[0].InviteStatus != "accepted" {
			t.Fatalf("head status: %s", det.Members[0].InviteStatus)
		}
	})

	// 3. Activate
	t.Run("03_activate", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/activate?user_id="+headUser, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var res householdResult
		mustDecode(t, parseSuccess(t, resp), &res)
		if res.Status != "Active" {
			t.Fatalf("expected Active, got %s", res.Status)
		}
		// Verify via fetch
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if det.Status != "Active" {
			t.Fatalf("expected Active after verify, got %s", det.Status)
		}
	})

	// 4. Invite member
	t.Run("04_invite", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/members?user_id="+headUser, map[string]interface{}{
			"user_id": memberUser,
			"role":    "Member",
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var res householdResult
		mustDecode(t, parseSuccess(t, resp), &res)
		if !res.Success {
			t.Fatal("expected success=true")
		}
		// Verify member is pending
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if det.MemberCount != 2 {
			t.Fatalf("expected 2 members, got %d", det.MemberCount)
		}
		var found bool
		for _, m := range det.Members {
			if m.UserID == memberUser {
				found = true
				if m.Role != "Member" {
					t.Fatalf("role: expected Member, got %s", m.Role)
				}
				if m.InviteStatus != "Pending" {
					t.Fatalf("invite status: expected Pending, got %s", m.InviteStatus)
				}
				break
			}
		}
		if !found {
			t.Fatal("member not found after invite")
		}
	})

	// 5. Accept invite
	t.Run("05_accept", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/members/"+memberUser+"/accept?user_id="+memberUser, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		for _, m := range det.Members {
			if m.UserID == memberUser {
				if m.InviteStatus != "Accepted" {
					t.Fatalf("expected Accepted, got %s", m.InviteStatus)
				}
				return
			}
		}
		t.Fatal("member not found after accept")
	})

	// 6. Update role to Admin
	t.Run("06_update_role", func(t *testing.T) {
		resp := ts.put(t, "/api/v1/households/"+hhID+"/members/"+memberUser+"/role?user_id="+headUser, map[string]interface{}{
			"role": "Admin",
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		for _, m := range det.Members {
			if m.UserID == memberUser {
				if m.Role != "Admin" {
					t.Fatalf("expected Admin, got %s", m.Role)
				}
				return
			}
		}
		t.Fatal("member not found after role update")
	})

	// 7. Link goal
	t.Run("07_link_goal", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/goals?user_id="+headUser, map[string]interface{}{
			"goal_id": goalID,
			"user_id": headUser,
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if len(det.LinkedGoals) != 1 {
			t.Fatalf("expected 1 linked goal, got %d", len(det.LinkedGoals))
		}
	})

	// 8. Link budget
	t.Run("08_link_budget", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/budgets?user_id="+headUser, map[string]interface{}{
			"budget_id": budgetID,
			"user_id":   headUser,
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if len(det.LinkedBudgets) != 1 {
			t.Fatalf("expected 1 linked budget, got %d", len(det.LinkedBudgets))
		}
	})

	// 9. Add goal contribution
	t.Run("09_contribution", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/goals/"+goalID+"/contributions", map[string]interface{}{
			"user_id": headUser,
			"amount":  50000,
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if len(det.GoalContributions) != 1 {
			t.Fatalf("expected 1 contribution, got %d", len(det.GoalContributions))
		}
		if det.GoalContributions[0].Amount != 50000 {
			t.Fatalf("amount: expected 50000, got %d", det.GoalContributions[0].Amount)
		}
	})

	// 10. Summary
	t.Run("10_summary", func(t *testing.T) {
		resp := ts.get(t, "/api/v1/households/"+hhID+"/summary")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var s financialSummary
		mustDecode(t, parseSuccess(t, resp), &s)
		if s.LinkedGoals != 1 {
			t.Fatalf("expected 1 goal in summary, got %d", s.LinkedGoals)
		}
		if s.LinkedBudgets != 1 {
			t.Fatalf("expected 1 budget in summary, got %d", s.LinkedBudgets)
		}
		if s.LinkedAccounts != 0 {
			t.Fatalf("expected 0 accounts, got %d", s.LinkedAccounts)
		}
		if s.TotalContributions != 50000 {
			t.Fatalf("expected 50000 contributions, got %d", s.TotalContributions)
		}
	})

	// 11. Pause
	t.Run("11_pause", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/pause?user_id="+headUser, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var res householdResult
		mustDecode(t, parseSuccess(t, resp), &res)
		if res.Status != "Paused" {
			t.Fatalf("expected Paused, got %s", res.Status)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if det.Status != "Paused" {
			t.Fatalf("expected Paused after verify, got %s", det.Status)
		}
	})

	// 12. Resume
	t.Run("12_resume", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/resume?user_id="+headUser, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var res householdResult
		mustDecode(t, parseSuccess(t, resp), &res)
		if res.Status != "Active" {
			t.Fatalf("expected Active, got %s", res.Status)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if det.Status != "Active" {
			t.Fatalf("expected Active after verify, got %s", det.Status)
		}
	})

	// 13. Remove member
	t.Run("13_remove_member", func(t *testing.T) {
		resp := ts.del(t, "/api/v1/households/"+hhID+"/members/"+memberUser+"?user_id="+headUser)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if det.MemberCount != 1 {
			t.Fatalf("expected 1 member after removal, got %d", det.MemberCount)
		}
		for _, m := range det.Members {
			if m.UserID == memberUser {
				t.Fatal("member should have been removed")
			}
		}
	})

	// 14. Dissolve
	t.Run("14_dissolve", func(t *testing.T) {
		resp := ts.post(t, "/api/v1/households/"+hhID+"/dissolve?user_id="+headUser, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var res householdResult
		mustDecode(t, parseSuccess(t, resp), &res)
		if res.Status != "Dissolved" {
			t.Fatalf("expected Dissolved, got %s", res.Status)
		}
		var det householdDetail
		mustDecode(t, parseSuccess(t, ts.get(t, "/api/v1/households/"+hhID)), &det)
		if det.Status != "Dissolved" {
			t.Fatalf("expected Dissolved after verify, got %s", det.Status)
		}
	})
}
