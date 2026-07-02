package engine

import (
	"context"
	"fmt"
	"time"
)

type RuleLoader interface {
	LoadAlertRules(ctx context.Context, userID string) ([]Rule, error)
}

type AlertEngine struct {
	ruleLoader RuleLoader
}

func NewAlertEngine(loader RuleLoader) *AlertEngine {
	return &AlertEngine{ruleLoader: loader}
}

func (e *AlertEngine) EvaluateAlerts(ctx context.Context, req AlertRequest) ([]SmartAlert, error) {
	rules, err := e.ruleLoader.LoadAlertRules(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("load rules: %w", err)
	}
	if len(rules) == 0 {
		return []SmartAlert{}, nil
	}

	payload := map[string]interface{}{
		"description":  req.Description,
		"amount":       req.Amount,
		"category":     req.Category,
		"account_type": req.AccountType,
		"balance":      req.Balance,
	}
	for k, v := range req.Extra {
		payload[k] = v
	}

	actions, err := Evaluate(payload, rules)
	if err != nil {
		return nil, fmt.Errorf("evaluate: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	var alerts []SmartAlert
	for _, a := range actions {
		alert := SmartAlert{
			AlertID:      fmt.Sprintf("alert-%s-%d", now, len(alerts)),
			Severity:     "info",
			Category:     "Alert",
			SourceEngine: "smart-alert",
			Timestamp:    now,
		}
		if title, ok := a.Params["title"].(string); ok {
			alert.Title = title
		} else {
			alert.Title = "Smart Alert Triggered"
		}
		if summary, ok := a.Params["summary"].(string); ok {
			alert.Summary = summary
		} else {
			alert.Summary = fmt.Sprintf("Rule matched for transaction: %s", req.Description)
		}
		if sev, ok := a.Params["severity"].(string); ok {
			alert.Severity = sev
		}
		if cat, ok := a.Params["category"].(string); ok {
			alert.Category = cat
		}
		if name, ok := a.Params["rule_name"].(string); ok {
			alert.RuleName = name
		}
		if id, ok := a.Params["rule_id"].(string); ok {
			alert.RuleID = id
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}
