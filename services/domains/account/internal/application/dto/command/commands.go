package command

type CreateAccountCommand struct {
	OwnerID       string  `json:"owner_id" validate:"required"`
	AccountName   string  `json:"account_name" validate:"required"`
	AccountType   string  `json:"account_type" validate:"required"`
	Classification string `json:"classification" validate:"required"`
	Currency      string  `json:"currency" validate:"required"`
	OpenedDate    string  `json:"opened_date"`
	Visibility    string  `json:"visibility"`
	Ownership     string  `json:"ownership"`
	InstitutionID string  `json:"institution_id,omitempty"`
	HouseholdID   string  `json:"household_id,omitempty"`
	SubType       string  `json:"sub_type,omitempty"`
	CreditLimit   *int64  `json:"credit_limit,omitempty"`
	InterestRate  *float64 `json:"interest_rate,omitempty"`
	Country       string  `json:"country,omitempty"`
	SourceOfTruth string  `json:"source_of_truth,omitempty"`
	ExtRef        string  `json:"external_reference,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}
type ActivateAccountCommand struct{ AccountID string `json:"account_id" validate:"required"` }
type FreezeAccountCommand struct{ AccountID string `json:"account_id" validate:"required"` }
type UnfreezeAccountCommand struct{ AccountID string `json:"account_id" validate:"required"` }
type CloseAccountCommand struct {
	AccountID  string `json:"account_id" validate:"required"`
	ClosedDate string `json:"closed_date" validate:"required"`
	Reason     string `json:"reason,omitempty"`
}
type ArchiveAccountCommand struct{ AccountID string `json:"account_id" validate:"required"` }
type AccountResult struct {
	AccountID string `json:"account_id"`
	Status    string `json:"status"`
	Success   bool   `json:"success"`
}
