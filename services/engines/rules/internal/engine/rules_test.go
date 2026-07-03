package engine

import (
	"testing"

	"github.com/google/uuid"
)

func TestEvaluate_NoRules(t *testing.T) {
	e := NewEngine()
	actions, err := e.Evaluate(map[string]interface{}{"amount": 100}, nil)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestEvaluate_SingleMatch(t *testing.T) {
	e := NewEngine()
	rules := []Rule{
		{
			ID: uuid.New(), Name: "Big Purchase", Enabled: true, Priority: 1,
			Conditions: []Condition{
				{Field: "amount", Operator: OpGreaterThan, Value: "500"},
			},
			Actions: []Action{
				{Type: ActionAlert, Params: map[string]interface{}{"severity": "high"}},
			},
		},
	}

	actions, err := e.Evaluate(map[string]interface{}{"amount": 1000}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type != ActionAlert {
		t.Errorf("expected Alert, got %s", actions[0].Type)
	}
}

func TestEvaluate_NoMatch(t *testing.T) {
	e := NewEngine()
	rules := []Rule{
		{
			ID: uuid.New(), Name: "Big Purchase", Enabled: true, Priority: 1,
			Conditions: []Condition{
				{Field: "amount", Operator: OpGreaterThan, Value: "500"},
			},
			Actions: []Action{{Type: ActionAlert}},
		},
	}

	actions, err := e.Evaluate(map[string]interface{}{"amount": 100}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestEvaluate_DisabledRule(t *testing.T) {
	e := NewEngine()
	rules := []Rule{
		{
			ID: uuid.New(), Name: "Disabled Rule", Enabled: false, Priority: 1,
			Conditions: []Condition{{Field: "amount", Operator: OpGreaterThan, Value: "0"}},
			Actions:    []Action{{Type: ActionAlert}},
		},
	}

	actions, err := e.Evaluate(map[string]interface{}{"amount": 100}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Error("expected 0 actions for disabled rule")
	}
}

func TestEvaluate_PriorityOrder(t *testing.T) {
	e := NewEngine()
	rules := []Rule{
		{
			ID: uuid.New(), Name: "Low Priority", Enabled: true, Priority: 10,
			Conditions: []Condition{{Field: "amount", Operator: OpGreaterThan, Value: "0"}},
			Actions:    []Action{{Type: ActionCategorize, Params: map[string]interface{}{"category": "low"}}},
		},
		{
			ID: uuid.New(), Name: "High Priority", Enabled: true, Priority: 1,
			Conditions: []Condition{{Field: "amount", Operator: OpGreaterThan, Value: "0"}},
			Actions:    []Action{{Type: ActionCategorize, Params: map[string]interface{}{"category": "high"}}},
		},
	}

	actions, err := e.Evaluate(map[string]interface{}{"amount": 100}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}
	// High priority (priority 1) should be first
	if actions[0].Type != ActionCategorize {
		t.Errorf("expected high priority action first")
	}
}

func TestEvaluate_MultipleConditions(t *testing.T) {
	e := NewEngine()
	rules := []Rule{
		{
			ID: uuid.New(), Name: "Groceries", Enabled: true, Priority: 1,
			Conditions: []Condition{
				{Field: "description", Operator: OpContains, Value: "grocery"},
				{Field: "amount", Operator: OpLessThan, Value: "5000"},
			},
			Actions: []Action{{Type: ActionCategorize, Params: map[string]interface{}{"category": "Food"}}},
		},
	}

	// Both conditions match
	actions, err := e.Evaluate(map[string]interface{}{"description": "Walmart Grocery", "amount": 2500}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 1 {
		t.Error("expected 1 action when both conditions match")
	}

	// Only one condition matches
	actions, err = e.Evaluate(map[string]interface{}{"description": "Netflix", "amount": 2500}, rules)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if len(actions) != 0 {
		t.Error("expected 0 actions when only one condition matches")
	}
}

func TestEvaluateCondition_Equals(t *testing.T) {
	matched, err := evaluateCondition(Condition{Field: "type", Operator: OpEquals, Value: "income"}, map[string]interface{}{"type": "Income"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected match for case-insensitive equals")
	}

	matched, err = evaluateCondition(Condition{Field: "type", Operator: OpEquals, Value: "income"}, map[string]interface{}{"type": "expense"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected no match for different values")
	}

	matched, err = evaluateCondition(Condition{Field: "missing", Operator: OpEquals, Value: "x"}, map[string]interface{}{"type": "income"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected no match for missing field")
	}
}

func TestEvaluateCondition_Contains(t *testing.T) {
	matched, err := evaluateCondition(Condition{Field: "description", Operator: OpContains, Value: "grocery"}, map[string]interface{}{"description": "Walmart Grocery Store"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !matched {
		t.Error("expected match for contains")
	}

	matched, err = evaluateCondition(Condition{Field: "description", Operator: OpContains, Value: "gym"}, map[string]interface{}{"description": "Netflix"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if matched {
		t.Error("expected no match")
	}
}

func TestEvaluateCondition_GreaterThan(t *testing.T) {
	matched, err := evaluateCondition(Condition{Field: "amount", Operator: OpGreaterThan, Value: "500"}, map[string]interface{}{"amount": 1000})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if !matched { t.Error("expected 1000 > 500") }

	matched, err = evaluateCondition(Condition{Field: "amount", Operator: OpGreaterThan, Value: "500"}, map[string]interface{}{"amount": 300})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if matched { t.Error("expected 300 not > 500") }
}

func TestEvaluateCondition_LessThan(t *testing.T) {
	matched, err := evaluateCondition(Condition{Field: "amount", Operator: OpLessThan, Value: "500"}, map[string]interface{}{"amount": 300})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if !matched { t.Error("expected 300 < 500") }
}

func TestEvaluateCondition_MatchesRegex(t *testing.T) {
	matched, err := evaluateCondition(Condition{Field: "email", Operator: OpMatchesRegex, Value: `^.*@example\.com$`}, map[string]interface{}{"email": "test@example.com"})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if !matched { t.Error("expected regex match") }

	matched, err = evaluateCondition(Condition{Field: "email", Operator: OpMatchesRegex, Value: `^.*@example\.com$`}, map[string]interface{}{"email": "test@other.com"})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if matched { t.Error("expected regex no match") }

	// Bad regex
	_, err = evaluateCondition(Condition{Field: "email", Operator: OpMatchesRegex, Value: `[invalid`}, map[string]interface{}{"email": "test"})
	if err == nil { t.Error("expected error for bad regex") }
}

func TestEvaluateCondition_Between(t *testing.T) {
	matched, err := evaluateCondition(Condition{Field: "amount", Operator: OpBetween, Value: "100,500"}, map[string]interface{}{"amount": 300})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if !matched { t.Error("expected 300 between 100 and 500") }

	matched, err = evaluateCondition(Condition{Field: "amount", Operator: OpBetween, Value: "100,500"}, map[string]interface{}{"amount": 50})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if matched { t.Error("expected 50 not between 100 and 500") }

	matched, err = evaluateCondition(Condition{Field: "amount", Operator: OpBetween, Value: "100,500"}, map[string]interface{}{"amount": 500})
	if err != nil { t.Fatalf("evaluate: %v", err) }
	if !matched { t.Error("expected 500 to match upper bound") }
}

func TestEvaluateCondition_UnknownOperator(t *testing.T) {
	_, err := evaluateCondition(Condition{Field: "x", Operator: "unknown", Value: "y"}, map[string]interface{}{"x": "z"})
	if err == nil {
		t.Error("expected error for unknown operator")
	}
}

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Fatal("expected non-nil engine")
	}
}
