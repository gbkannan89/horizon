package command

type CreateLiabilityCommand struct {
	OwnerID              string  `json:"owner_id" validate:"required"`
	Name                 string  `json:"name" validate:"required"`
	Classification       string  `json:"classification" validate:"required"`
	Currency             string  `json:"currency" validate:"required"`
	OriginalPrincipal    int64   `json:"original_principal" validate:"required"`
	InterestRate         float64 `json:"interest_rate" validate:"required"`
	InterestModel        string  `json:"interest_model" validate:"required"`
	RepaymentProfile     string  `json:"repayment_profile"`
	RepaymentModel       string  `json:"repayment_model"`
	InstallmentAmount    int64   `json:"installment_amount"`
	InstallmentFreq      string  `json:"installment_frequency"`
	RemainingInstallments int    `json:"remaining_installments"`
	MaturityDate         string  `json:"maturity_date" validate:"required"`
	SourceOfTruth        string  `json:"source_of_truth"`
	LiabilityType        string  `json:"liability_type,omitempty"`
	HouseholdID          string  `json:"household_id,omitempty"`
	InstitutionID        string  `json:"institution_id,omitempty"`
	ServicingAccountID   string  `json:"servicing_account_id,omitempty"`
	Collateral           string  `json:"collateral,omitempty"`
	ExtRef               string  `json:"external_reference,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	Notes                string  `json:"notes,omitempty"`
}

type PaymentCommand struct {
	LiabilityID string `json:"liability_id" validate:"required"`
	Amount      int64  `json:"amount" validate:"required"`
}

type PrepaymentCommand struct {
	LiabilityID string `json:"liability_id" validate:"required"`
	Amount      int64  `json:"amount" validate:"required"`
}

type IDCmd struct{ LiabilityID string `json:"liability_id" validate:"required"` }
type LiabResult struct {
	LiabilityID string `json:"liability_id"`
	Status      string `json:"status"`
	Success     bool   `json:"success"`
}
