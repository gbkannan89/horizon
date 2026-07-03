package engine

type TaxInput struct {
	GrossSalary       int64   `json:"gross_salary"`
	BusinessIncome    int64   `json:"business_income"`
	CapitalGains      int64   `json:"capital_gains"`
	RentalIncome      int64   `json:"rental_income"`
	OtherIncome       int64   `json:"other_income"`
	Section80C        int64   `json:"section_80c"`
	Section80D        int64   `json:"section_80d"`
	HRAReceived       int64   `json:"hra_received"`
	HRARentPaid       int64   `json:"hra_rent_paid"`
	HomeLoanInterest  int64   `json:"home_loan_interest"`
	NPSContribution   int64   `json:"nps_contribution"`
	HouseRentAllow    int64   `json:"house_rent_allowance,omitempty"`
}

type TaxRegime string
const (
	RegimeOld TaxRegime = "old"
	RegimeNew TaxRegime = "new"
)

type TaxSlab struct {
	MinIncome int64 `json:"min_income"`
	MaxIncome int64 `json:"max_income"`
	Rate      float64 `json:"rate"`
}

type TaxComputation struct {
	Regime        TaxRegime `json:"regime"`
	GrossIncome   int64     `json:"gross_income"`
	Deductions    int64     `json:"deductions"`
	TaxableIncome int64     `json:"taxable_income"`
	TaxAmount     int64     `json:"tax_amount"`
	CessAmount    int64     `json:"cess_amount"`
	TotalTax      int64     `json:"total_tax"`
	EffectiveRate float64   `json:"effective_rate"`
}

type TaxAnalysis struct {
	OldRegime TaxComputation `json:"old_regime"`
	NewRegime TaxComputation `json:"new_regime"`
	Recommended TaxRegime    `json:"recommended"`
	Savings     int64        `json:"savings"`
	Suggestions []string     `json:"suggestions,omitempty"`
}
