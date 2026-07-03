package e2e_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// -------- Types (mirror services/engines/rules/internal/engine) --------

type operator string

const (
	opEquals       operator = "equals"
	opContains     operator = "contains"
	opGreaterThan  operator = "greater_than"
	opLessThan     operator = "less_than"
	opMatchesRegex operator = "matches_regex"
	opBetween      operator = "between"
)

type actionType string

const (
	actionCategorize actionType = "categorize"
	actionTag        actionType = "tag"
	actionAlert      actionType = "alert"
	actionNotify     actionType = "notify"
	actionSkip       actionType = "skip"
	actionSplit      actionType = "split"
)

type condition struct {
	Field    string   `json:"field"`
	Operator operator `json:"operator"`
	Value    string   `json:"value"`
}

type action struct {
	Type   actionType              `json:"type"`
	Params map[string]interface{}  `json:"params"`
}

type rule struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Category    string      `json:"category"`
	Priority    int         `json:"priority"`
	Enabled     bool        `json:"enabled"`
	Conditions  []condition `json:"conditions"`
	Actions     []action    `json:"actions"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

// -------- Engine Evaluate (mirror services/engines/rules/internal/engine) --------

func evaluate(payload map[string]interface{}, rules []rule) ([]action, error) {
	sorted := make([]rule, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority < sorted[j].Priority
	})

	var triggered []action

	for _, r := range sorted {
		if !r.Enabled {
			continue
		}

		match := true
		for _, c := range r.Conditions {
			ok, err := evaluateCondition(c, payload)
			if err != nil {
				match = false
				break
			}
			if !ok {
				match = false
				break
			}
		}

		if match {
			triggered = append(triggered, r.Actions...)
		}
	}

	return triggered, nil
}

func evaluateCondition(c condition, payload map[string]interface{}) (bool, error) {
	val, ok := payload[c.Field]
	if !ok {
		return false, nil
	}

	strVal := fmt.Sprintf("%v", val)

	switch c.Operator {
	case opEquals:
		return strings.EqualFold(strVal, c.Value), nil
	case opContains:
		return strings.Contains(strings.ToLower(strVal), strings.ToLower(c.Value)), nil
	case opGreaterThan:
		pf, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return false, nil
		}
		cf, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return false, nil
		}
		return pf > cf, nil
	case opLessThan:
		pf, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return false, nil
		}
		cf, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return false, nil
		}
		return pf < cf, nil
	case opMatchesRegex:
		matched, err := regexp.MatchString(c.Value, strVal)
		if err != nil {
			return false, err
		}
		return matched, nil
	case opBetween:
		parts := strings.Split(c.Value, ",")
		if len(parts) != 2 {
			return false, nil
		}
		min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		pf, err3 := strconv.ParseFloat(strVal, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return false, nil
		}
		return pf >= min && pf <= max, nil
	}

	return false, fmt.Errorf("unknown operator: %s", c.Operator)
}

// -------- In-memory rule store --------

type store struct {
	mu    sync.RWMutex
	rules map[string]*rule
}

func newStore() *store {
	return &store{rules: make(map[string]*rule)}
}

func (s *store) add(r *rule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[r.ID] = r
}

func (s *store) get(id string) *rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rules[id]
}

func (s *store) list(userID string) []*rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*rule
	for _, r := range s.rules {
		if r.UserID == userID {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt < out[j].CreatedAt
	})
	return out
}

func (s *store) delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rules, id)
}

// -------- HTTP server --------

type ruleStore interface {
	add(r *rule)
	get(id string) *rule
	list(userID string) []*rule
	delete(id string)
}

type ruleHandler struct {
	store ruleStore
	now   func() time.Time
}

func newHandler(s ruleStore) *ruleHandler {
	return &ruleHandler{store: s, now: time.Now}
}

type createRuleRequest struct {
	Name       string                 `json:"name"`
	Category   string                 `json:"category"`
	Priority   int                    `json:"priority"`
	Enabled    bool                   `json:"enabled"`
	Conditions []condition            `json:"conditions"`
	Actions    []action               `json:"actions"`
}

type ruleResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *ruleHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ruleResponse{Error: "invalid request body"})
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "default"
	}

	now := h.now().UTC().Format(time.RFC3339)
	rl := &rule{
		ID:         fmt.Sprintf("rule-%d", h.now().UnixNano()),
		UserID:     userID,
		Name:       req.Name,
		Category:   req.Category,
		Priority:   req.Priority,
		Enabled:    req.Enabled,
		Conditions: req.Conditions,
		Actions:    req.Actions,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	h.store.add(rl)
	writeJSON(w, http.StatusCreated, ruleResponse{
		Success: true,
		Data:    map[string]interface{}{"rule_id": rl.ID},
	})
}

func (h *ruleHandler) list(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "default"
	}

	rules := h.store.list(userID)
	if rules == nil {
		rules = []*rule{}
	}
	writeJSON(w, http.StatusOK, ruleResponse{Success: true, Data: rules})
}

func (h *ruleHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ruleResponse{Error: "missing id"})
		return
	}
	h.store.delete(id)
	writeJSON(w, http.StatusOK, ruleResponse{Success: true, Data: map[string]string{"status": "deleted"}})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func mustDecode[T any](t *testing.T, body io.Reader) *T {
	t.Helper()
	var v T
	if err := json.NewDecoder(body).Decode(&v); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return &v
}

// -------- Tests --------

func TestRulesEngineE2E_AutoCategorizeWorkflow(t *testing.T) {
	s := newStore()
	h := newHandler(s)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/rules", h.create)
	mux.HandleFunc("GET /api/v1/rules", h.list)
	mux.HandleFunc("DELETE /api/v1/rules/{id}", h.delete)

	server := httptest.NewServer(mux)
	defer server.Close()

	userID := "user-e2e-001"
	client := server.Client()
	base := server.URL + "/api/v1/rules"

	// -----------------------------------------------------------------------
	// Step 1: Create 3 rules
	// -----------------------------------------------------------------------

	// 1a. "Grocery" — description contains "grocery" → categorize as "Food"
	createAndAssert(t, client, base, userID, createRuleRequest{
		Name:     "Grocery",
		Category: "Food",
		Priority: 10,
		Enabled:  true,
		Conditions: []condition{
			{Field: "description", Operator: opContains, Value: "grocery"},
		},
		Actions: []action{
			{Type: actionCategorize, Params: map[string]interface{}{"category": "Food"}},
		},
	})

	// 1b. "Salary" — description contains "salary" AND amount > 50000 → categorize as "Income"
	createAndAssert(t, client, base, userID, createRuleRequest{
		Name:     "Salary",
		Category: "Income",
		Priority: 5,
		Enabled:  true,
		Conditions: []condition{
			{Field: "description", Operator: opContains, Value: "salary"},
			{Field: "amount", Operator: opGreaterThan, Value: "50000"},
		},
		Actions: []action{
			{Type: actionCategorize, Params: map[string]interface{}{"category": "Income"}},
		},
	})

	// 1c. "Large Purchase" — amount > 20000 → alert with severity "high"
	createAndAssert(t, client, base, userID, createRuleRequest{
		Name:     "Large Purchase",
		Category: "Alert",
		Priority: 1,
		Enabled:  true,
		Conditions: []condition{
			{Field: "amount", Operator: opGreaterThan, Value: "20000"},
		},
		Actions: []action{
			{Type: actionAlert, Params: map[string]interface{}{"severity": "high"}},
		},
	})

	// -----------------------------------------------------------------------
	// Step 2: List rules — verify all 3 returned
	// -----------------------------------------------------------------------
	rules := listRules(t, client, base, userID)
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	names := make(map[string]bool)
	for _, r := range rules {
		names[r.Name] = true
	}
	for _, n := range []string{"Grocery", "Salary", "Large Purchase"} {
		if !names[n] {
			t.Errorf("missing rule %q in list", n)
		}
	}

	// -----------------------------------------------------------------------
	// Step 3: Evaluate "Walmart Grocery" (amount 3500)
	// -----------------------------------------------------------------------
	actions, err := evaluate(map[string]interface{}{
		"description": "Walmart Grocery",
		"amount":      3500,
	}, toSlice(rules))
	if err != nil {
		t.Fatalf("evaluate grocery: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("grocery: expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != actionCategorize {
		t.Errorf("grocery: expected categorize action, got %s", actions[0].Type)
	}
	if actions[0].Params["category"] != "Food" {
		t.Errorf("grocery: expected category=Food, got %v", actions[0].Params["category"])
	}

	// -----------------------------------------------------------------------
	// Step 5-6: Evaluate "Monthly Salary" (amount 75000)
	//           Should match "Salary" (categorize) AND "Large Purchase" (alert)
	// -----------------------------------------------------------------------
	actions, err = evaluate(map[string]interface{}{
		"description": "Monthly Salary",
		"amount":      75000,
	}, toSlice(rules))
	if err != nil {
		t.Fatalf("evaluate salary: %v", err)
	}
	if len(actions) != 2 {
		t.Fatalf("salary: expected 2 actions, got %d. Actions: %+v", len(actions), actions)
	}

	// "Large Purchase" has priority 1, so its alert action should come first
	if actions[0].Type != actionAlert {
		t.Errorf("salary: expected first action to be alert, got %s", actions[0].Type)
	}
	if actions[1].Type != actionCategorize {
		t.Errorf("salary: expected second action to be categorize, got %s", actions[1].Type)
	}

	// -----------------------------------------------------------------------
	// Step 7: Disable the "Large Purchase" rule
	// -----------------------------------------------------------------------
	for _, r := range rules {
		if r.Name == "Large Purchase" {
			r.Enabled = false
			break
		}
	}

	// -----------------------------------------------------------------------
	// Step 8: Re-evaluate salary payload — should only match "Salary"
	// -----------------------------------------------------------------------
	actions, err = evaluate(map[string]interface{}{
		"description": "Monthly Salary",
		"amount":      75000,
	}, toSlice(rules))
	if err != nil {
		t.Fatalf("evaluate after disable: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("after disable: expected 1 action, got %d. Actions: %+v", len(actions), actions)
	}
	if actions[0].Type != actionCategorize {
		t.Errorf("after disable: expected categorize, got %s", actions[0].Type)
	}

	// -----------------------------------------------------------------------
	// Step 9: Delete the "Grocery" rule via HTTP
	// -----------------------------------------------------------------------
	for _, r := range rules {
		if r.Name == "Grocery" {
			req, _ := http.NewRequest("DELETE", base+"/"+r.ID, nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("delete grocery: %v", err)
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("delete grocery: status %d", resp.StatusCode)
			}
			break
		}
	}

	// -----------------------------------------------------------------------
	// Step 10: List rules — verify only 2 remain
	// -----------------------------------------------------------------------
	rules = listRules(t, client, base, userID)
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules after delete, got %d", len(rules))
	}
	for _, r := range rules {
		if r.Name == "Grocery" {
			t.Error("Grocery rule should have been deleted")
		}
	}
}

// -------- Helpers --------

func createAndAssert(t *testing.T, client *http.Client, base, userID string, req createRuleRequest) {
	t.Helper()
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", base+"?user_id="+userID, strings.NewReader(string(body)))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("create %s: %v", req.Name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create %s: status %d", req.Name, resp.StatusCode)
	}
}

func listRules(t *testing.T, client *http.Client, base, userID string) []*rule {
	t.Helper()
	resp, err := client.Get(base + "?user_id=" + userID)
	if err != nil {
		t.Fatalf("list rules: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list rules: status %d", resp.StatusCode)
	}
	var res ruleResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	data, err := json.Marshal(res.Data)
	if err != nil {
		t.Fatalf("marshal list data: %v", err)
	}
	var rules []*rule
	if err := json.Unmarshal(data, &rules); err != nil {
		t.Fatalf("unmarshal rules: %v", err)
	}
	return rules
}

func toSlice(rules []*rule) []rule {
	out := make([]rule, len(rules))
	for i, r := range rules {
		out[i] = *r
	}
	return out
}
