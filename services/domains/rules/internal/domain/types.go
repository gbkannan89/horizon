package domain

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
	id          string
	userID      string
	householdID string
	name        string
	description string
	category    string
	priority    int
	enabled     bool
	conditions  []Condition
	actions     []Action
	createdAt   string
	updatedAt   string
}

func (r *Rule) ID() string                    { return r.id }
func (r *Rule) UserID() string                { return r.userID }
func (r *Rule) HouseholdID() string           { return r.householdID }
func (r *Rule) Name() string                  { return r.name }
func (r *Rule) Description() string           { return r.description }
func (r *Rule) Category() string              { return r.category }
func (r *Rule) Priority() int                 { return r.priority }
func (r *Rule) Enabled() bool                 { return r.enabled }
func (r *Rule) Conditions() []Condition       { return r.conditions }
func (r *Rule) Actions() []Action             { return r.actions }
func (r *Rule) CreatedAt() string             { return r.createdAt }
func (r *Rule) UpdatedAt() string             { return r.updatedAt }

func NewRule(id, userID, householdID, name, description, category string,
	priority int, enabled bool, conditions []Condition, actions []Action,
	createdAt, updatedAt string,
) (*Rule, error) {
	if conditions == nil { conditions = []Condition{} }
	if actions == nil { actions = []Action{} }
	return &Rule{
		id: id, userID: userID, householdID: householdID,
		name: name, description: description, category: category,
		priority: priority, enabled: enabled,
		conditions: conditions, actions: actions,
		createdAt: createdAt, updatedAt: updatedAt,
	}, nil
}
