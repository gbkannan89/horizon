package query

type GetGoalQuery struct {
	GoalID string `json:"goal_id" validate:"required"`
}

type ListByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type ListByStatusQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Status string `json:"status" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type ListByImportanceQuery struct {
	UserID     string `json:"user_id" validate:"required"`
	Importance string `json:"importance" validate:"required"`
	Cursor     string `json:"cursor,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type ListByTypeQuery struct {
	UserID   string `json:"user_id" validate:"required"`
	GoalType string `json:"goal_type" validate:"required"`
	Cursor   string `json:"cursor,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type ListByHouseholdQuery struct {
	HouseholdID string `json:"household_id" validate:"required"`
	Cursor      string `json:"cursor,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

type GoalView struct {
	GoalID      string  `json:"goal_id"`
	UserID      string  `json:"user_id"`
	HouseholdID string  `json:"household_id,omitempty"`
	Name        string  `json:"name"`
	Importance  string  `json:"importance"`
	GoalType    string  `json:"goal_type"`
	Subtype     string  `json:"subtype,omitempty"`
	Priority    int     `json:"priority"`
	Progress    float64 `json:"progress"`
	Remaining   float64 `json:"remaining"`
	Status      string  `json:"status"`
	HasTargetDate bool  `json:"has_target_date"`
	CreatedAt   string  `json:"created_at"`
}

type PaginatedResult struct {
	Goals      []GoalView `json:"goals"`
	NextCursor string     `json:"next_cursor,omitempty"`
	HasMore    bool       `json:"has_more"`
}
