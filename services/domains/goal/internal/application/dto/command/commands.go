package command

type CreateGoalCommand struct {
	UserID         string   `json:"user_id" validate:"required"`
	Name           string   `json:"name" validate:"required"`
	Importance     string   `json:"importance" validate:"required"`
	GoalType       string   `json:"goal_type" validate:"required"`
	Subtype        string   `json:"subtype,omitempty"`
	SuccessModel   string   `json:"success_model" validate:"required"`
	TargetValue    *float64 `json:"target_value,omitempty"`
	TargetMonths   *int     `json:"target_months,omitempty"`
	CustomDesc     *string  `json:"custom_description,omitempty"`
	CustomTargetVal *float64 `json:"custom_target_value,omitempty"`
	Priority       int      `json:"priority" validate:"required"`
	TargetDate     *string  `json:"target_date,omitempty"`
	RiskTolerance  string   `json:"risk_tolerance,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	ParentGoalID   string   `json:"parent_goal_id,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

type ActivateGoalCommand struct {
	GoalID string `json:"goal_id" validate:"required"`
}

type PauseGoalCommand struct {
	GoalID string `json:"goal_id" validate:"required"`
	Reason string `json:"reason,omitempty"`
}

type ResumeGoalCommand struct {
	GoalID string `json:"goal_id" validate:"required"`
}

type RestructureGoalCommand struct {
	GoalID         string   `json:"goal_id" validate:"required"`
	SuccessModel   string   `json:"success_model" validate:"required"`
	TargetValue    *float64 `json:"target_value,omitempty"`
	TargetMonths   *int     `json:"target_months,omitempty"`
	CustomDesc     *string  `json:"custom_description,omitempty"`
	CustomTargetVal *float64 `json:"custom_target_value,omitempty"`
	TargetDate     *string  `json:"target_date,omitempty"`
	Priority       int      `json:"priority" validate:"required"`
}

type CompleteGoalCommand struct {
	GoalID string `json:"goal_id" validate:"required"`
}

type MissGoalCommand struct {
	GoalID string `json:"goal_id" validate:"required"`
}

type ArchiveGoalCommand struct {
	GoalID string `json:"goal_id" validate:"required"`
}

type UpdatePriorityCommand struct {
	GoalID      string `json:"goal_id" validate:"required"`
	NewPriority int    `json:"new_priority" validate:"required"`
}

type GoalResult struct {
	GoalID  string `json:"goal_id"`
	Status  string `json:"status"`
	Success bool   `json:"success"`
}
