package query

import "github.com/horizon/core/services/domains/rules/internal/domain"

type RuleView struct {
	RuleID      string             `json:"rule_id"`
	UserID      string             `json:"user_id"`
	HouseholdID string             `json:"household_id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Category    string             `json:"category"`
	Priority    int                `json:"priority"`
	Enabled     bool               `json:"enabled"`
	Conditions  []domain.Condition `json:"conditions"`
	Actions     []domain.Action    `json:"actions"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}
