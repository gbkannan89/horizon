package engine

import (
	"sort"
)

// Engine is the central evaluator for rules
type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

// Evaluate runs a payload through a set of rules and returns all triggered Actions.
// Rules are evaluated in order of Priority (lowest number = highest priority).
// For this MVP, we assume implicit AND across all conditions within a single Rule.
func (e *Engine) Evaluate(payload map[string]interface{}, rules []Rule) ([]Action, error) {
	// Sort rules by priority ascending
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
				// If a condition errors out (e.g., bad regex), consider it a non-match
				ruleMatches = false
				break
			}
			if !matched {
				ruleMatches = false
				break
			}
		}

		if ruleMatches {
			triggeredActions = append(triggeredActions, rule.Actions...)
		}
	}

	return triggeredActions, nil
}
