package engine

type BalanceType string
const (BTCurrent BalanceType = "current"; BTAvailable BalanceType = "available"; BTCleared BalanceType = "cleared"; BTUncleared BalanceType = "uncleared"; BTReserved BalanceType = "reserved"; BTAllocated BalanceType = "allocated"; BTSpendable BalanceType = "spendable"; BTProjected BalanceType = "projected")

type AccountCardData struct {
	AccountID        string `json:"account_id"`; AccountName string `json:"account_name"`
	AccountType      string `json:"account_type"`; Classification string `json:"classification"`
	Status           string `json:"status"`; Currency string `json:"currency"`
	CurrentBalance   int64  `json:"current_balance"`; AvailableBalance int64 `json:"available_balance"`
	InstitutionName  string `json:"institution_name,omitempty"`; LastActivityDate string `json:"last_activity_date,omitempty"`
	HasImport        bool   `json:"has_import"`
}

type BalanceEntry struct {
	Type BalanceType `json:"type"`; Value int64 `json:"value"`
	Label string `json:"label"`; Description string `json:"description"`
}

type AccountDetailData struct {
	AccountID string `json:"account_id"`; AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`; Classification string `json:"classification"`
	Status string `json:"status"`; Currency string `json:"currency"`
	Balances []BalanceEntry `json:"balances"`
	InstitutionName string `json:"institution_name,omitempty"`; ImportStatus string `json:"import_status,omitempty"`
	LastSyncTime string `json:"last_sync_time,omitempty"`; Transactions []TransPreview `json:"transactions"`
	OpenedDate string `json:"opened_date"`; ClosedDate string `json:"closed_date,omitempty"`
	LiquidityProfile string `json:"liquidity_profile"`; AccountHealth string `json:"account_health"`
	CreditLimit *int64 `json:"credit_limit,omitempty"`; InterestRate *float64 `json:"interest_rate,omitempty"`
	Visibility string `json:"visibility"`; Tags []string `json:"tags"`; Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}

type TransPreview struct {
	EventID string `json:"event_id"`; EventType string `json:"event_type"`
	Amount int64 `json:"amount"`; Currency string `json:"currency"`
	Description string `json:"description"`; Category string `json:"category"`
	EventDate string `json:"event_date"`; RunningBalance int64 `json:"running_balance"`
}

type AccountsList struct {
	Accounts []AccountCardData `json:"accounts"`; Total int `json:"total"`
	Cursor string `json:"cursor,omitempty"`; HasMore bool `json:"has_more"`
	Summary *AccountsSummary `json:"summary,omitempty"`
}

type AccountsSummary struct {
	TotalBalance int64 `json:"total_balance"`
	CountByType map[string]int `json:"count_by_type"`
	CountByStatus map[string]int `json:"count_by_status"`
}

type AccountInput struct {
	AccountID string `json:"account_id"`; AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`; Classification string `json:"classification"`
	Status string `json:"status"`; Currency string `json:"currency"`
	CurrentBalance int64 `json:"current_balance"`; AvailableBalance int64 `json:"available_balance"`
	ClearedBalance int64 `json:"cleared_balance"`; UnclearedBalance int64 `json:"uncleared_balance"`
	ReservedBalance int64 `json:"reserved_balance"`; AllocatedBalance int64 `json:"allocated_balance"`
	SpendableBalance int64 `json:"spendable_balance"`; ProjectedBalance *int64 `json:"projected_balance,omitempty"`
	InstitutionName string `json:"institution_name,omitempty"`; InstitutionID string `json:"institution_id,omitempty"`
	ImportStatus string `json:"import_status,omitempty"`; LastSyncTime string `json:"last_sync_time,omitempty"`
	HasImport bool `json:"has_import"`; OpenedDate string `json:"opened_date"`
	ClosedDate string `json:"closed_date,omitempty"`; LiquidityProfile string `json:"liquidity_profile"`
	AccountHealth string `json:"account_health"`; CreditLimit *int64 `json:"credit_limit,omitempty"`
	InterestRate *float64 `json:"interest_rate,omitempty"`; Visibility string `json:"visibility"`
	Tags []string `json:"tags"`; Notes string `json:"notes"`; CreatedAt string `json:"created_at"`
}

type TransInput struct {
	EventID string `json:"event_id"`; EventType string `json:"event_type"`
	Amount int64 `json:"amount"`; Currency string `json:"currency"`
	Description string `json:"description"`; Category string `json:"category"`
	EventDate string `json:"event_date"`; RunningBalance int64 `json:"running_balance"`
}

type FilterOpts struct {
	AccountTypes []string `json:"account_types,omitempty"`
	Statuses     []string `json:"statuses,omitempty"`
	Currencies   []string `json:"currencies,omitempty"`
}
