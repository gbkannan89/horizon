package domain

type HouseholdStatus string

const (
	HHStatusDraft     HouseholdStatus = "Draft"
	HHStatusActive    HouseholdStatus = "Active"
	HHStatusPaused    HouseholdStatus = "Paused"
	HHStatusDissolved HouseholdStatus = "Dissolved"
	HHStatusArchived  HouseholdStatus = "Archived"
)

type HouseholdType string

const (
	HHSingle   HouseholdType = "Single"
	HHCouple   HouseholdType = "Couple"
	HHFamily   HouseholdType = "Family"
	HHRoommates HouseholdType = "Roommates"
	HHCustom   HouseholdType = "Custom"
)

type MemberRole string

const (
	RoleHead   MemberRole = "Head"
	RoleAdmin  MemberRole = "Admin"
	RoleMember MemberRole = "Member"
	RoleViewer MemberRole = "Viewer"
)

type HouseholdHealth string

const (
	HHHealthy HouseholdHealth = "Healthy"
	HHWarning HouseholdHealth = "Warning"
	HHCritical HouseholdHealth = "Critical"
)

type HouseholdMember struct {
	UserID    string     `json:"user_id"`
	Role      MemberRole `json:"role"`
	AddedAt   string     `json:"added_at"`
	InviteStatus string  `json:"invite_status"`
}

type LinkedAccount struct {
	AccountID string `json:"account_id"`
	AddedBy   string `json:"added_by"`
	AddedAt   string `json:"added_at"`
}

type LinkedGoal struct {
	GoalID  string `json:"goal_id"`
	AddedBy string `json:"added_by"`
	AddedAt string `json:"added_at"`
}

type LinkedBudget struct {
	BudgetID string `json:"budget_id"`
	AddedBy  string `json:"added_by"`
	AddedAt  string `json:"added_at"`
}

type GoalContribution struct {
	GoalID string `json:"goal_id"`
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
	Date   string `json:"date"`
}
