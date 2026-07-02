package engine

import (
	"context"
	"fmt"
)

type RuleLoader interface {
	LoadActiveRules(ctx context.Context, userID string) ([]Rule, error)
}

type CategorizeEngine struct {
	ruleLoader RuleLoader
}

func NewCategorizeEngine(loader RuleLoader) *CategorizeEngine {
	return &CategorizeEngine{ruleLoader: loader}
}

func (e *CategorizeEngine) Categorize(ctx context.Context, req CategorizeRequest) (*CategorizationResult, error) {
	rules, err := e.ruleLoader.LoadActiveRules(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("load rules: %w", err)
	}
	if len(rules) == 0 {
		return &CategorizationResult{
			TransactionID: req.TransactionID,
			Actions:       []string{},
		}, nil
	}

	payload := map[string]interface{}{
		"description": req.Description,
		"amount":      req.Amount,
		"category":    req.Category,
	}
	for k, v := range req.Extra {
		payload[k] = v
	}

	actions, err := Evaluate(payload, rules)
	if err != nil {
		return nil, fmt.Errorf("evaluate: %w", err)
	}

	result := &CategorizationResult{
		TransactionID: req.TransactionID,
		Actions:       make([]string, 0, len(actions)),
	}

	for _, a := range actions {
		result.Actions = append(result.Actions, string(a.Type))
		if a.Type == ActionCategorize {
			if cat, ok := a.Params["category"].(string); ok {
				result.SuggestedCategory = cat
			}
			if name, ok := a.Params["rule_name"].(string); ok {
				result.RuleName = name
			}
			if id, ok := a.Params["rule_id"].(string); ok {
				result.RuleID = id
			}
		}
	}

	return result, nil
}
