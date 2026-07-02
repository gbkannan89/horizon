package engine

type CategorizationResult struct {
	TransactionID     string   `json:"transaction_id"`
	SuggestedCategory string   `json:"suggested_category"`
	RuleName          string   `json:"rule_name,omitempty"`
	RuleID            string   `json:"rule_id,omitempty"`
	Actions           []string `json:"actions"`
}

type CategorizeRequest struct {
	TransactionID string                 `json:"transaction_id"`
	UserID        string                 `json:"user_id"`
	Description   string                 `json:"description"`
	Amount        float64                `json:"amount"`
	Category      string                 `json:"category"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}
