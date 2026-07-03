package domain

type BudgetStatus string

const (
	BStatusDraft     BudgetStatus = "Draft"
	BStatusActive    BudgetStatus = "Active"
	BStatusPaused    BudgetStatus = "Paused"
	BStatusCompleted BudgetStatus = "Completed"
	BStatusArchived  BudgetStatus = "Archived"
)

type BudgetPeriod string

const (
	PeriodWeekly  BudgetPeriod = "Weekly"
	PeriodMonthly BudgetPeriod = "Monthly"
	PeriodQuarterly BudgetPeriod = "Quarterly"
	PeriodYearly  BudgetPeriod = "Yearly"
	PeriodCustom  BudgetPeriod = "Custom"
)

type BudgetHealth string

const (
	HealthOnTrack   BudgetHealth = "OnTrack"
	HealthOverspent BudgetHealth = "Overspent"
	HealthUnderTrack BudgetHealth = "UnderTrack"
	HealthCritical  BudgetHealth = "Critical"
)

type BudgetCategory struct {
	ID             string `json:"id"`
	BudgetID       string `json:"budget_id"`
	Category       string `json:"category"`
	Subcategory    string `json:"subcategory,omitempty"`
	BudgetedAmount int64  `json:"budgeted_amount"`
	SpentAmount    int64  `json:"spent_amount"`
	RemainingAmount int64 `json:"remaining_amount"`
	Rollover       bool   `json:"rollover"`
}
