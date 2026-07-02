package command

import "github.com/horizon/core/services/domains/rules/internal/domain"

type CreateRuleCommand struct {
	UserID      string           `json:"user_id"`
	HouseholdID string           `json:"household_id,omitempty"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Category    string           `json:"category"`
	Priority    int              `json:"priority"`
	Enabled     bool             `json:"enabled"`
	Conditions  []domain.Condition `json:"conditions"`
	Actions     []domain.Action    `json:"actions"`
}

type UpdateRuleCommand struct {
	RuleID      string             `json:"rule_id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Category    string             `json:"category"`
	Priority    int                `json:"priority"`
	Enabled     bool               `json:"enabled"`
	Conditions  []domain.Condition `json:"conditions"`
	Actions     []domain.Action    `json:"actions"`
}

type RuleResult struct {
	RuleID  string `json:"rule_id"`
	Success bool   `json:"success"`
}
