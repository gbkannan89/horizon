package domain

type GoalStatus string

const (
	StatusDraft        GoalStatus = "Draft"
	StatusActive       GoalStatus = "Active"
	StatusPaused       GoalStatus = "Paused"
	StatusCompleted    GoalStatus = "Completed"
	StatusMissed       GoalStatus = "Missed"
	StatusRestructured GoalStatus = "Restructured"
	StatusArchived     GoalStatus = "Archived"
)

type GoalImportance string

const (
	ImpMandatory  GoalImportance = "Mandatory"
	ImpEssential  GoalImportance = "Essential"
	ImpLifestyle  GoalImportance = "Lifestyle"
	ImpDream      GoalImportance = "Dream"
	ImpSpeculative GoalImportance = "Speculative"
)

var AllImportance = map[GoalImportance]int{
	ImpMandatory: 5, ImpEssential: 4, ImpLifestyle: 3, ImpDream: 2, ImpSpeculative: 1,
}

type GoalType string

const (
	TypeTimeBound GoalType = "TimeBound"
	TypeOpenEnded GoalType = "OpenEnded"
	TypeLifestyle GoalType = "Lifestyle"
)

type GoalSubtype string

const (
	SubSavings          GoalSubtype = "Savings"
	SubEducation        GoalSubtype = "Education"
	SubPurchase         GoalSubtype = "Purchase"
	SubRetirement       GoalSubtype = "Retirement"
	SubIncome           GoalSubtype = "Income"
	SubNetWorth         GoalSubtype = "NetWorth"
	SubDebtFreedom      GoalSubtype = "DebtFreedom"
	SubMonthlyCommitment GoalSubtype = "MonthlyCommitment"
	SubAnnualCommitment GoalSubtype = "AnnualCommitment"
)

type SuccessModel string

const (
	ModelTargetAmount         SuccessModel = "TargetAmount"
	ModelMonthlyPassiveIncome SuccessModel = "MonthlyPassiveIncome"
	ModelNetWorth             SuccessModel = "NetWorth"
	ModelDebtFree             SuccessModel = "DebtFree"
	ModelEmergencyFundMonths  SuccessModel = "EmergencyFundMonths"
	ModelCustomKPI            SuccessModel = "CustomKPI"
)

type RiskTolerance string

const (
	RiskConservative RiskTolerance = "Conservative"
	RiskModerate     RiskTolerance = "Moderate"
	RiskAggressive   RiskTolerance = "Aggressive"
)

var AllRiskTolerance = map[RiskTolerance]bool{
	RiskConservative: true, RiskModerate: true, RiskAggressive: true,
}

type SuccessCriteria struct {
	Model            SuccessModel `json:"model"`
	TargetValue      *float64     `json:"target_value,omitempty"`
	TargetMonths     *int         `json:"target_months,omitempty"`
	CustomDesc       *string      `json:"custom_description,omitempty"`
	CustomTargetVal  *float64     `json:"custom_target_value,omitempty"`
}

type ContributionSchedule struct {
	Amount    float64 `json:"amount"`
	Frequency string  `json:"frequency"`
}
