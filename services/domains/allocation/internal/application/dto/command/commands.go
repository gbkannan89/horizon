package command

type CreateAllocationCommand struct {
	GoalID          string  `json:"goal_id" validate:"required"`
	FundingSourceID string  `json:"funding_source_id" validate:"required"`
	AllocationType  string  `json:"allocation_type" validate:"required"`
	Strategy        string  `json:"strategy"`
	Priority        int     `json:"priority"`
	Currency        string  `json:"currency" validate:"required"`
	EffectiveDate   string  `json:"effective_date" validate:"required"`
	Weight          *float64 `json:"weight,omitempty"`
	FixedAmount     *int64  `json:"fixed_amount,omitempty"`
	ExpirationDate  *string `json:"expiration_date,omitempty"`
	SourceOfTruth   string  `json:"source_of_truth"`
	CreatedBy       string  `json:"created_by"`
	Tags            []string `json:"tags,omitempty"`
	Notes           string  `json:"notes,omitempty"`
}
type ApproveAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type PlanAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type ActivateAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type PauseAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type ResumeAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type CompleteAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type CancelAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type ArchiveAllocationCommand struct{ AllocationID string `json:"allocation_id" validate:"required"` }
type AllocResult struct {
	AllocationID string `json:"allocation_id"`
	Status       string `json:"status"`
	Success      bool   `json:"success"`
}
