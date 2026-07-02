package http

type CreateRuleRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Category    string        `json:"category"`
	Priority    int           `json:"priority"`
	Enabled     bool          `json:"enabled"`
	Conditions  []ConditionDTO `json:"conditions"`
	Actions     []ActionDTO    `json:"actions"`
}

type ConditionDTO struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type ActionDTO struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type RuleResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
