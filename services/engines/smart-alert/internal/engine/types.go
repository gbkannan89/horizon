package engine

type SmartAlert struct {
	AlertID      string `json:"alert_id"`
	Title        string `json:"title"`
	Summary      string `json:"summary"`
	Severity     string `json:"severity"`
	Category     string `json:"category"`
	RuleName     string `json:"rule_name,omitempty"`
	RuleID       string `json:"rule_id,omitempty"`
	SourceEngine string `json:"source_engine"`
	Timestamp     string `json:"timestamp"`
}

type AlertRequest struct {
	UserID      string                 `json:"user_id"`
	Description string                 `json:"description"`
	Amount      float64                `json:"amount"`
	Category    string                 `json:"category"`
	AccountType string                 `json:"account_type,omitempty"`
	Balance     float64                `json:"balance,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}
