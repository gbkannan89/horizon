package engine

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Operator string

const (
	OpEquals       Operator = "equals"
	OpContains     Operator = "contains"
	OpGreaterThan  Operator = "greater_than"
	OpLessThan     Operator = "less_than"
	OpMatchesRegex Operator = "matches_regex"
	OpBetween      Operator = "between"
)

type ActionType string

const (
	ActionCategorize ActionType = "categorize"
	ActionTag        ActionType = "tag"
	ActionAlert      ActionType = "alert"
	ActionNotify     ActionType = "notify"
	ActionSkip       ActionType = "skip"
	ActionSplit      ActionType = "split"
)

type Condition struct {
	Field    string   `json:"field"`
	Operator Operator `json:"operator"`
	Value    string   `json:"value"`
}

type Action struct {
	Type   ActionType              `json:"type"`
	Params map[string]interface{}  `json:"params"`
}

type Rule struct {
	ID          string
	UserID      string
	Name        string
	Category    string
	Priority    int
	Enabled     bool
	Conditions  []Condition
	Actions     []Action
}

func Evaluate(payload map[string]interface{}, rules []Rule) ([]Action, error) {
	sortedRules := make([]Rule, len(rules))
	copy(sortedRules, rules)
	sort.Slice(sortedRules, func(i, j int) bool {
		return sortedRules[i].Priority < sortedRules[j].Priority
	})

	var triggeredActions []Action

	for _, rule := range sortedRules {
		if !rule.Enabled {
			continue
		}

		ruleMatches := true
		for _, cond := range rule.Conditions {
			matched, err := evaluateCondition(cond, payload)
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

func evaluateCondition(cond Condition, payload map[string]interface{}) (bool, error) {
	val, ok := payload[cond.Field]
	if !ok {
		return false, nil
	}

	strVal := fmt.Sprintf("%v", val)

	switch cond.Operator {
	case OpEquals:
		return strings.EqualFold(strVal, cond.Value), nil
	case OpContains:
		return strings.Contains(strings.ToLower(strVal), strings.ToLower(cond.Value)), nil
	case OpGreaterThan:
		payloadFloat, err := strconv.ParseFloat(strVal, 64)
		if err != nil { return false, nil }
		condFloat, err := strconv.ParseFloat(cond.Value, 64)
		if err != nil { return false, nil }
		return payloadFloat > condFloat, nil
	case OpLessThan:
		payloadFloat, err := strconv.ParseFloat(strVal, 64)
		if err != nil { return false, nil }
		condFloat, err := strconv.ParseFloat(cond.Value, 64)
		if err != nil { return false, nil }
		return payloadFloat < condFloat, nil
	case OpMatchesRegex:
		matched, err := regexp.MatchString(cond.Value, strVal)
		if err != nil { return false, err }
		return matched, nil
	case OpBetween:
		parts := strings.Split(cond.Value, ",")
		if len(parts) != 2 { return false, nil }
		min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		pf, err3 := strconv.ParseFloat(strVal, 64)
		if err1 != nil || err2 != nil || err3 != nil { return false, nil }
		return pf >= min && pf <= max, nil
	}

	return false, fmt.Errorf("unknown operator: %s", cond.Operator)
}
