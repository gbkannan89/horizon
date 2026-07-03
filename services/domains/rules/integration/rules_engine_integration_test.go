package rules_test

import (
	"bytes"
	"context"
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

	"github.com/horizon/core/services/domains/rules/internal/application"
	"github.com/horizon/core/services/domains/rules/internal/application/dto/command"
	"github.com/horizon/core/services/domains/rules/internal/application/dto/query"
	"github.com/horizon/core/services/domains/rules/internal/domain"
	"github.com/horizon/core/services/internal/auth"
)

// ---------------------------------------------------------------------------
// Engine types (replicated from services/engines/auto-categorize/internal/engine)
// so the test can evaluate rules without cross-module internal imports.
// ---------------------------------------------------------------------------

type engineOperator string

const (
	engineOpEquals       engineOperator = "equals"
	engineOpContains     engineOperator = "contains"
	engineOpGreaterThan  engineOperator = "greater_than"
	engineOpLessThan     engineOperator = "less_than"
	engineOpMatchesRegex engineOperator = "matches_regex"
	engineOpBetween      engineOperator = "between"
)

type engineActionType string

const (
	engineActionCategorize engineActionType = "categorize"
	engineActionTag        engineActionType = "tag"
	engineActionAlert      engineActionType = "alert"
	engineActionNotify     engineActionType = "notify"
	engineActionSkip       engineActionType = "skip"
	engineActionSplit      engineActionType = "split"
)

type engineCondition struct {
	Field    string         `json:"field"`
	Operator engineOperator `json:"operator"`
	Value    string         `json:"value"`
}

type engineAction struct {
	Type   engineActionType        `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type engineRule struct {
	ID         string
	UserID     string
	Name       string
	Category   string
	Priority   int
	Enabled    bool
	Conditions []engineCondition
	Actions    []engineAction
}

// engineEvaluate runs payload through rules and returns triggered actions.
// Replicates the logic in services/engines/auto-categorize/internal/engine.
func engineEvaluate(payload map[string]interface{}, rules []engineRule) ([]engineAction, error) {
	sortedRules := make([]engineRule, len(rules))
	copy(sortedRules, rules)
	sort.Slice(sortedRules, func(i, j int) bool {
		return sortedRules[i].Priority < sortedRules[j].Priority
	})

	var triggeredActions []engineAction

	for _, rule := range sortedRules {
		if !rule.Enabled {
			continue
		}

		ruleMatches := true
		for _, cond := range rule.Conditions {
			matched, err := engineEvaluateCondition(cond, payload)
			if err != nil {
				ruleMatches = false
				break
			}
			if !matched {
				ruleMatches = false
				break
			}
		}

		if ruleMatches {
			for _, a := range rule.Actions {
				if a.Params == nil {
					a.Params = map[string]interface{}{}
				}
				a.Params["rule_name"] = rule.Name
				a.Params["rule_id"] = rule.ID
			}
			triggeredActions = append(triggeredActions, rule.Actions...)
		}
	}

	return triggeredActions, nil
}

func engineEvaluateCondition(cond engineCondition, payload map[string]interface{}) (bool, error) {
	val, ok := payload[cond.Field]
	if !ok {
		return false, nil
	}

	strVal := fmt.Sprintf("%v", val)

	switch cond.Operator {
	case engineOpEquals:
		return strings.EqualFold(strVal, cond.Value), nil
	case engineOpContains:
		return strings.Contains(strings.ToLower(strVal), strings.ToLower(cond.Value)), nil
	case engineOpGreaterThan:
		payloadFloat, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return false, nil
		}
		condFloat, err := strconv.ParseFloat(cond.Value, 64)
		if err != nil {
			return false, nil
		}
		return payloadFloat > condFloat, nil
	case engineOpLessThan:
		payloadFloat, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return false, nil
		}
		condFloat, err := strconv.ParseFloat(cond.Value, 64)
		if err != nil {
			return false, nil
		}
		return payloadFloat < condFloat, nil
	case engineOpMatchesRegex:
		matched, err := regexp.MatchString(cond.Value, strVal)
		if err != nil {
			return false, err
		}
		return matched, nil
	case engineOpBetween:
		parts := strings.Split(cond.Value, ",")
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

	return false, fmt.Errorf("unknown operator: %s", cond.Operator)
}

// ---------------------------------------------------------------------------
// Mock Repository
// ---------------------------------------------------------------------------

type mockRepo struct {
	mu    sync.Mutex
	rules map[string]*domain.Rule
}

func newMockRepo() *mockRepo {
	return &mockRepo{rules: make(map[string]*domain.Rule)}
}

func (r *mockRepo) Save(_ context.Context, rule *domain.Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := rule.ID()
	if id == "" {
		id = fmt.Sprintf("rule-%d", len(r.rules)+1)
		now := time.Now().UTC().Format(time.RFC3339)
		rebuilt, err := domain.NewRule(
			id, rule.UserID(), rule.HouseholdID(), rule.Name(),
			rule.Description(), rule.Category(), rule.Priority(),
			rule.Enabled(), rule.Conditions(), rule.Actions(),
			now, now,
		)
		if err != nil {
			return err
		}
		*rule = *rebuilt
	}
	r.rules[rule.ID()] = rule
	return nil
}

func (r *mockRepo) GetByID(_ context.Context, id string) (*domain.Rule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rule, ok := r.rules[id]
	if !ok {
		return nil, fmt.Errorf("rule not found")
	}
	return rule, nil
}

func (r *mockRepo) ListByUser(_ context.Context, userID string) ([]*domain.Rule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Rule
	for _, rule := range r.rules {
		if rule.UserID() == userID {
			result = append(result, rule)
		}
	}
	return result, nil
}

func (r *mockRepo) ListByCategory(_ context.Context, category string) ([]*domain.Rule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Rule
	for _, rule := range r.rules {
		if rule.Category() == category {
			result = append(result, rule)
		}
	}
	return result, nil
}

func (r *mockRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, id)
	return nil
}

// ---------------------------------------------------------------------------
// Response envelope types
// ---------------------------------------------------------------------------

type apiResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// ---------------------------------------------------------------------------
// HTTP handlers (test double matching rules HTTP API)
// ---------------------------------------------------------------------------

type handler struct {
	svc *application.RuleService
}

func newHandler(svc *application.RuleService) *handler {
	return &handler{svc: svc}
}

func (h *handler) createRule(w http.ResponseWriter, r *http.Request) {
	var req createRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	conds := make([]domain.Condition, len(req.Conditions))
	for i, c := range req.Conditions {
		conds[i] = domain.Condition{Field: c.Field, Operator: domain.Operator(c.Operator), Value: c.Value}
	}
	acts := make([]domain.Action, len(req.Actions))
	for i, a := range req.Actions {
		acts[i] = domain.Action{Type: domain.ActionType(a.Type), Params: a.Params}
	}
	cmd := command.CreateRuleCommand{
		UserID:      auth.UserIDFromRequest(r),
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Priority:    req.Priority,
		Enabled:     req.Enabled,
		Conditions:  conds,
		Actions:     acts,
	}
	result, err := h.svc.Create(r.Context(), cmd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": result})
}

func (h *handler) getRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	view, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": view})
}

func (h *handler) listRules(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	rules, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rules == nil {
		rules = []query.RuleView{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": rules})
}

func (h *handler) deleteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

type createRuleRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Category    string             `json:"category"`
	Priority    int                `json:"priority"`
	Enabled     bool               `json:"enabled"`
	Conditions  []conditionDTO     `json:"conditions"`
	Actions     []actionDTO        `json:"actions"`
}

type conditionDTO struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type actionDTO struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

// ---------------------------------------------------------------------------
// Response helpers
// ---------------------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// ---------------------------------------------------------------------------
// Test server setup
// ---------------------------------------------------------------------------

func newTestServer(t *testing.T) (*httptest.Server, *mockRepo) {
	t.Helper()
	repo := newMockRepo()
	svc := application.NewRuleService(repo, time.Now)
	h := newHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/rules", h.createRule)
	mux.HandleFunc("GET /api/v1/rules/{id}", h.getRule)
	mux.HandleFunc("GET /api/v1/rules", h.listRules)
	mux.HandleFunc("DELETE /api/v1/rules/{id}", h.deleteRule)

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

// ---------------------------------------------------------------------------
// Tests: Rules HTTP API
// ---------------------------------------------------------------------------

func TestRulesAPI_CreateAndGet(t *testing.T) {
	srv, _ := newTestServer(t)
	userID := "rule-user"

	body := `{"name":"Groceries Rule","description":"Categorize grocery transactions","category":"Food","priority":1,"enabled":true,"conditions":[{"field":"description","operator":"contains","value":"grocery"}],"actions":[{"type":"categorize","params":{"category":"Food"}}]}`
	resp := doRequest(t, srv, "POST", "/api/v1/rules?user_id="+userID, []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)
	if !ar.Success {
		t.Fatalf("expected success, got error: %s", ar.Error)
	}

	var createResult command.RuleResult
	mustUnmarshalData(t, ar.Data, &createResult)
	if !createResult.Success {
		t.Fatal("expected success=true in result")
	}
	ruleID := createResult.RuleID
	if ruleID == "" {
		t.Fatal("expected non-empty rule_id")
	}

	// GET the created rule
	resp = doRequest(t, srv, "GET", "/api/v1/rules/"+ruleID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar = parseResponse(t, resp)
	if !ar.Success {
		t.Fatalf("expected success, got error: %s", ar.Error)
	}

	var view query.RuleView
	mustUnmarshalData(t, ar.Data, &view)
	if view.RuleID != ruleID {
		t.Fatalf("expected rule_id %s, got %s", ruleID, view.RuleID)
	}
	if view.Name != "Groceries Rule" {
		t.Fatalf("expected 'Groceries Rule', got %s", view.Name)
	}
	if view.Category != "Food" {
		t.Fatalf("expected Category Food, got %s", view.Category)
	}
	if view.Priority != 1 {
		t.Fatalf("expected Priority 1, got %d", view.Priority)
	}
	if !view.Enabled {
		t.Fatal("expected Enabled=true")
	}
	if len(view.Conditions) != 1 {
		t.Fatalf("expected 1 condition, got %d", len(view.Conditions))
	}
	if len(view.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(view.Actions))
	}
	if view.Actions[0].Type != domain.ActionCategorize {
		t.Fatalf("expected categorize action, got %s", view.Actions[0].Type)
	}
}

func TestRulesAPI_CreateAndList(t *testing.T) {
	srv, _ := newTestServer(t)
	userID := "list-user"

	for i := 0; i < 3; i++ {
		body := fmt.Sprintf(`{"name":"Rule %d","category":"Test","priority":%d,"enabled":true,"conditions":[{"field":"amount","operator":"greater_than","value":"%d"}],"actions":[{"type":"alert","params":{"severity":"low"}}]}`, i+1, i+1, i*100)
		resp := doRequest(t, srv, "POST", "/api/v1/rules?user_id="+userID, []byte(body))
		assertStatus(t, resp.StatusCode, http.StatusCreated)
	}

	resp := doRequest(t, srv, "GET", "/api/v1/rules?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)
	ar := parseResponse(t, resp)
	if !ar.Success {
		t.Fatalf("expected success, got error: %s", ar.Error)
	}

	var rules []query.RuleView
	mustUnmarshalData(t, ar.Data, &rules)
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
}

func TestRulesAPI_GetNotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/api/v1/rules/non-existent-id?user_id=u", nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
}

func TestRulesAPI_Delete(t *testing.T) {
	srv, _ := newTestServer(t)
	userID := "del-user"

	body := `{"name":"To Delete","category":"Temp","priority":5,"enabled":true,"conditions":[{"field":"amount","operator":"greater_than","value":"0"}],"actions":[{"type":"alert","params":{}}]}`
	resp := doRequest(t, srv, "POST", "/api/v1/rules?user_id="+userID, []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)

	var createResult command.RuleResult
	mustUnmarshalData(t, ar.Data, &createResult)
	ruleID := createResult.RuleID

	// Delete
	resp = doRequest(t, srv, "DELETE", "/api/v1/rules/"+ruleID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusOK)

	// Verify deleted
	resp = doRequest(t, srv, "GET", "/api/v1/rules/"+ruleID+"?user_id="+userID, nil)
	assertStatus(t, resp.StatusCode, http.StatusNotFound)
}

func TestRulesAPI_InvalidBody(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := doRequest(t, srv, "POST", "/api/v1/rules?user_id=u", []byte(`not json`))
	assertStatus(t, resp.StatusCode, http.StatusBadRequest)
}

func TestRulesAPI_EmptyName(t *testing.T) {
	srv, _ := newTestServer(t)
	body := `{"name":"","category":"Test","priority":1,"enabled":true,"conditions":[{"field":"amount","operator":"greater_than","value":"0"}],"actions":[{"type":"alert","params":{}}]}`
	resp := doRequest(t, srv, "POST", "/api/v1/rules?user_id=u", []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusInternalServerError)
}

// ---------------------------------------------------------------------------
// Tests: Engine Evaluation (auto-categorize flow)
// ---------------------------------------------------------------------------

func TestEngineEvaluate_NoRules(t *testing.T) {
	actions, err := engineEvaluate(map[string]interface{}{"amount": 100}, nil)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestEngineEvaluate_EmptyRules(t *testing.T) {
	actions, err := engineEvaluate(map[string]interface{}{"amount": 100}, []engineRule{})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestEngineEvaluate_SingleMatch(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-1", Name: "Big Purchase", Enabled: true, Priority: 1,
			Conditions: []engineCondition{
				{Field: "amount", Operator: engineOpGreaterThan, Value: "500"},
			},
			Actions: []engineAction{
				{Type: engineActionAlert, Params: map[string]interface{}{"severity": "high"}},
			},
		},
	}

	actions, err := engineEvaluate(map[string]interface{}{"amount": 1000}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != engineActionAlert {
		t.Errorf("expected Alert, got %s", actions[0].Type)
	}
	// Verify rule metadata was injected
	if actions[0].Params["rule_name"] != "Big Purchase" {
		t.Errorf("expected rule_name 'Big Purchase', got %v", actions[0].Params["rule_name"])
	}
	if actions[0].Params["rule_id"] != "rule-1" {
		t.Errorf("expected rule_id 'rule-1', got %v", actions[0].Params["rule_id"])
	}
}

func TestEngineEvaluate_NoMatch(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-1", Name: "Big Purchase", Enabled: true, Priority: 1,
			Conditions: []engineCondition{
				{Field: "amount", Operator: engineOpGreaterThan, Value: "500"},
			},
			Actions: []engineAction{{Type: engineActionAlert}},
		},
	}

	actions, err := engineEvaluate(map[string]interface{}{"amount": 100}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestEngineEvaluate_DisabledRule(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-1", Name: "Disabled Rule", Enabled: false, Priority: 1,
			Conditions: []engineCondition{{Field: "amount", Operator: engineOpGreaterThan, Value: "0"}},
			Actions:    []engineAction{{Type: engineActionAlert}},
		},
	}

	actions, err := engineEvaluate(map[string]interface{}{"amount": 100}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Error("expected 0 actions for disabled rule")
	}
}

func TestEngineEvaluate_PriorityOrder(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-low", Name: "Low Priority", Enabled: true, Priority: 10,
			Conditions: []engineCondition{{Field: "amount", Operator: engineOpGreaterThan, Value: "0"}},
			Actions:    []engineAction{{Type: engineActionCategorize, Params: map[string]interface{}{"category": "low"}}},
		},
		{
			ID: "rule-high", Name: "High Priority", Enabled: true, Priority: 1,
			Conditions: []engineCondition{{Field: "amount", Operator: engineOpGreaterThan, Value: "0"}},
			Actions:    []engineAction{{Type: engineActionCategorize, Params: map[string]interface{}{"category": "high"}}},
		},
	}

	actions, err := engineEvaluate(map[string]interface{}{"amount": 100}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}
	// High priority (priority 1) should be first
	if actions[0].Type != engineActionCategorize {
		t.Errorf("expected high priority action first")
	}
}

func TestEngineEvaluate_MultipleConditions(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-1", Name: "Groceries", Enabled: true, Priority: 1,
			Conditions: []engineCondition{
				{Field: "description", Operator: engineOpContains, Value: "grocery"},
				{Field: "amount", Operator: engineOpLessThan, Value: "5000"},
			},
			Actions: []engineAction{{Type: engineActionCategorize, Params: map[string]interface{}{"category": "Food"}}},
		},
	}

	// Both conditions match
	actions, err := engineEvaluate(map[string]interface{}{"description": "Walmart Grocery", "amount": 2500}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Error("expected 1 action when both conditions match")
	}

	// Only one condition matches
	actions, err = engineEvaluate(map[string]interface{}{"description": "Netflix", "amount": 2500}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Error("expected 0 actions when only one condition matches")
	}
}

func TestEngineEvaluate_AutoCategorizeWithRuleMetadata(t *testing.T) {
	rules := []engineRule{
		{
			ID: "cat-rule-1", Name: "Coffee Shop", Enabled: true, Priority: 1,
			Conditions: []engineCondition{
				{Field: "description", Operator: engineOpContains, Value: "starbucks"},
			},
			Actions: []engineAction{
				{Type: engineActionCategorize, Params: map[string]interface{}{"category": "Dining Out"}},
			},
		},
	}

	payload := map[string]interface{}{
		"description": "Starbucks Coffee",
		"amount":      12.50,
		"category":    "Uncategorized",
	}

	actions, err := engineEvaluate(payload, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}

	// The action should have categorize type and appropriate params
	if actions[0].Type != engineActionCategorize {
		t.Fatalf("expected categorize action, got %s", actions[0].Type)
	}
	if actions[0].Params["category"] != "Dining Out" {
		t.Fatalf("expected category 'Dining Out', got %v", actions[0].Params["category"])
	}
	if actions[0].Params["rule_name"] != "Coffee Shop" {
		t.Fatalf("expected rule_name 'Coffee Shop', got %v", actions[0].Params["rule_name"])
	}
}

func TestEngineEvaluate_MultipleRulesWithDifferentActions(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-1", Name: "Large Transaction", Enabled: true, Priority: 1,
			Conditions: []engineCondition{
				{Field: "amount", Operator: engineOpGreaterThan, Value: "10000"},
			},
			Actions: []engineAction{
				{Type: engineActionAlert, Params: map[string]interface{}{"severity": "high"}},
				{Type: engineActionNotify, Params: map[string]interface{}{"channel": "email"}},
			},
		},
		{
			ID: "rule-2", Name: "Small Transaction", Enabled: true, Priority: 2,
			Conditions: []engineCondition{
				{Field: "amount", Operator: engineOpLessThan, Value: "100"},
			},
			Actions: []engineAction{
				{Type: engineActionSkip, Params: map[string]interface{}{"reason": "insignificant"}},
			},
		},
	}

	// Large transaction match
	actions, err := engineEvaluate(map[string]interface{}{"amount": 50000}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions for large transaction, got %d", len(actions))
	}
	if actions[0].Type != engineActionAlert {
		t.Errorf("expected first action Alert, got %s", actions[0].Type)
	}
	if actions[1].Type != engineActionNotify {
		t.Errorf("expected second action Notify, got %s", actions[1].Type)
	}

	// Small transaction match
	actions, err = engineEvaluate(map[string]interface{}{"amount": 50}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action for small transaction, got %d", len(actions))
	}
	if actions[0].Type != engineActionSkip {
		t.Errorf("expected Skip action, got %s", actions[0].Type)
	}

	// No match
	actions, err = engineEvaluate(map[string]interface{}{"amount": 500}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("expected 0 actions for medium transaction, got %d", len(actions))
	}
}

func TestEngineEvaluate_MissingField(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-1", Name: "Category Check", Enabled: true, Priority: 1,
			Conditions: []engineCondition{
				{Field: "category", Operator: engineOpEquals, Value: "Food"},
			},
			Actions: []engineAction{{Type: engineActionTag, Params: map[string]interface{}{"tag": "edible"}}},
		},
	}

	// Missing field should not match
	actions, err := engineEvaluate(map[string]interface{}{"description": "Pizza"}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Error("expected 0 actions for missing field")
	}
}

func TestEngineEvaluate_EqualsOperator(t *testing.T) {
	matched, err := engineEvaluateCondition(engineCondition{Field: "type", Operator: engineOpEquals, Value: "income"}, map[string]interface{}{"type": "Income"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected match for case-insensitive equals")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "type", Operator: engineOpEquals, Value: "income"}, map[string]interface{}{"type": "expense"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected no match for different values")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "missing", Operator: engineOpEquals, Value: "x"}, map[string]interface{}{"type": "income"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected no match for missing field")
	}
}

func TestEngineEvaluate_ContainsOperator(t *testing.T) {
	matched, err := engineEvaluateCondition(engineCondition{Field: "description", Operator: engineOpContains, Value: "grocery"}, map[string]interface{}{"description": "Walmart Grocery Store"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected match for contains")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "description", Operator: engineOpContains, Value: "gym"}, map[string]interface{}{"description": "Netflix"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected no match")
	}
}

func TestEngineEvaluate_GreaterThanOperator(t *testing.T) {
	matched, err := engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpGreaterThan, Value: "500"}, map[string]interface{}{"amount": 1000})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected 1000 > 500")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpGreaterThan, Value: "500"}, map[string]interface{}{"amount": 300})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected 300 not > 500")
	}
}

func TestEngineEvaluate_LessThanOperator(t *testing.T) {
	matched, err := engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpLessThan, Value: "500"}, map[string]interface{}{"amount": 300})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected 300 < 500")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpLessThan, Value: "500"}, map[string]interface{}{"amount": 500})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected 500 not < 500")
	}
}

func TestEngineEvaluate_MatchesRegexOperator(t *testing.T) {
	matched, err := engineEvaluateCondition(engineCondition{Field: "email", Operator: engineOpMatchesRegex, Value: `^.*@example\.com$`}, map[string]interface{}{"email": "test@example.com"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected regex match")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "email", Operator: engineOpMatchesRegex, Value: `^.*@example\.com$`}, map[string]interface{}{"email": "test@other.com"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected regex no match")
	}

	// Bad regex should return error
	_, err = engineEvaluateCondition(engineCondition{Field: "email", Operator: engineOpMatchesRegex, Value: `[invalid`}, map[string]interface{}{"email": "test"})
	if err == nil {
		t.Error("expected error for bad regex")
	}
}

func TestEngineEvaluate_BetweenOperator(t *testing.T) {
	matched, err := engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpBetween, Value: "100,500"}, map[string]interface{}{"amount": 300})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected 300 between 100 and 500")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpBetween, Value: "100,500"}, map[string]interface{}{"amount": 50})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected 50 not between 100 and 500")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpBetween, Value: "100,500"}, map[string]interface{}{"amount": 100})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected 100 to match lower bound")
	}

	matched, err = engineEvaluateCondition(engineCondition{Field: "amount", Operator: engineOpBetween, Value: "100,500"}, map[string]interface{}{"amount": 500})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected 500 to match upper bound")
	}
}

func TestEngineEvaluate_UnknownOperator(t *testing.T) {
	_, err := engineEvaluateCondition(engineCondition{Field: "x", Operator: "unknown", Value: "y"}, map[string]interface{}{"x": "z"})
	if err == nil {
		t.Error("expected error for unknown operator")
	}
}

func TestEngineEvaluate_MixedEnabledAndDisabled(t *testing.T) {
	rules := []engineRule{
		{
			ID: "rule-enabled", Name: "Enabled Rule", Enabled: true, Priority: 1,
			Conditions: []engineCondition{{Field: "amount", Operator: engineOpGreaterThan, Value: "0"}},
			Actions:    []engineAction{{Type: engineActionCategorize, Params: map[string]interface{}{"category": "Matched"}}},
		},
		{
			ID: "rule-disabled", Name: "Disabled Rule", Enabled: false, Priority: 2,
			Conditions: []engineCondition{{Field: "amount", Operator: engineOpGreaterThan, Value: "0"}},
			Actions:    []engineAction{{Type: engineActionCategorize, Params: map[string]interface{}{"category": "ShouldNotMatch"}}},
		},
	}

	actions, err := engineEvaluate(map[string]interface{}{"amount": 100}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action from enabled rule, got %d", len(actions))
	}
	if actions[0].Params["category"] == "ShouldNotMatch" {
		t.Error("disabled rule should not have triggered")
	}
}

// ---------------------------------------------------------------------------
// Tests: Combined flow — create rules via API, then evaluate with engine
// ---------------------------------------------------------------------------

func TestEngineEvaluate_AfterAPICreation(t *testing.T) {
	srv, repo := newTestServer(t)
	userID := "combined-user"

	// Create a rule via API
	body := `{"name":"Dining Rule","category":"Food","priority":5,"enabled":true,"conditions":[{"field":"description","operator":"contains","value":"restaurant"}],"actions":[{"type":"categorize","params":{"category":"Dining"}}]}`
	resp := doRequest(t, srv, "POST", "/api/v1/rules?user_id="+userID, []byte(body))
	assertStatus(t, resp.StatusCode, http.StatusCreated)
	ar := parseResponse(t, resp)

	var createResult command.RuleResult
	mustUnmarshalData(t, ar.Data, &createResult)
	ruleID := createResult.RuleID

	// Fetch the rule from repo and convert to engineRule
	repo.mu.Lock()
	savedRule, ok := repo.rules[ruleID]
	repo.mu.Unlock()
	if !ok {
		t.Fatal("rule not found in repo")
	}

	er := domainRuleToEngineRule(savedRule)

	// Run evaluation using the engine
	payload := map[string]interface{}{
		"description": "Fancy Restaurant",
		"amount":      250.0,
		"category":    "Uncategorized",
	}

	actions, err := engineEvaluate(payload, []engineRule{er})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != engineActionCategorize {
		t.Fatalf("expected categorize action, got %s", actions[0].Type)
	}
	if actions[0].Params["category"] != "Dining" {
		t.Fatalf("expected category 'Dining', got %v", actions[0].Params["category"])
	}
}

func domainRuleToEngineRule(r *domain.Rule) engineRule {
	conds := make([]engineCondition, len(r.Conditions()))
	for i, c := range r.Conditions() {
		conds[i] = engineCondition{
			Field:    c.Field,
			Operator: engineOperator(c.Operator),
			Value:    c.Value,
		}
	}
	acts := make([]engineAction, len(r.Actions()))
	for i, a := range r.Actions() {
		acts[i] = engineAction{
			Type:   engineActionType(a.Type),
			Params: a.Params,
		}
	}
	return engineRule{
		ID:         r.ID(),
		Name:       r.Name(),
		Category:   r.Category(),
		Priority:   r.Priority(),
		Enabled:    r.Enabled(),
		Conditions: conds,
		Actions:    acts,
	}
}

// Ensure domain.Repository interface is satisfied at compile time
var _ domain.Repository = (*mockRepo)(nil)
