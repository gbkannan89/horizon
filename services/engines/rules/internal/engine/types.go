package engine

import (
	"time"

	"github.com/google/uuid"
)

type Operator string

const (
	OpEquals       Operator = "equals"
	OpContains     Operator = "contains"
	OpGreaterThan  Operator = "greater_than"
	OpLessThan     Operator = "less_than"
	OpMatchesRegex Operator = "matches_regex"
	OpBetween      Operator = "between" // assumes value is comma-separated "min,max"
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

// Condition represents a single logical evaluation criteria
type Condition struct {
	Field    string   `json:"field"`
	Operator Operator `json:"operator"`
	Value    string   `json:"value"` // Cast to required type internally for simplicity
}

// Action represents an effect applied if conditions meet
type Action struct {
	Type   ActionType             `json:"type"`
	Params map[string]interface{} `json:"params"`
}

// Rule represents the domain entity
type Rule struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	HouseholdID *uuid.UUID
	Name        string
	Description string
	Category    string
	Priority    int
	Enabled     bool
	Conditions  []Condition
	Actions     []Action
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
