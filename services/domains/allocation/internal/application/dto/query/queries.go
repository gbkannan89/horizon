package query

type GetByIDQuery struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type ListByGoalQuery struct {
	GoalID string `json:"goal_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type ListBySourceQuery struct {
	SourceID string `json:"source_id" validate:"required"`
	Cursor   string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type AllocView struct {
	AllocationID    string `json:"allocation_id"`
	GoalID          string `json:"goal_id"`
	FundingSourceID string `json:"funding_source_id"`
	AllocationType  string `json:"allocation_type"`
	Priority        int    `json:"priority"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	Health          string `json:"health"`
	ReservedAmount  int64  `json:"reserved_amount"`
	AllocatedAmount int64  `json:"allocated_amount"`
	CreatedAt       string `json:"created_at"`
}
type PaginatedResult struct {
	Allocations []AllocView `json:"allocations"`
	NextCursor  string      `json:"next_cursor,omitempty"`
	HasMore     bool        `json:"has_more"`
}
