package domain

import "context"

type Repository interface {
	Save(ctx context.Context, a *Allocation) error
	UpdateStatus(ctx context.Context, id string, from, to AllocationStatus) error
	GetByID(ctx context.Context, id string) (*Allocation, error)
	ListByGoal(ctx context.Context, goalID string, cursor string, limit int) ([]*Allocation, string, error)
	ListByFundingSource(ctx context.Context, sourceID string, cursor string, limit int) ([]*Allocation, string, error)
	ListByStatus(ctx context.Context, status AllocationStatus, cursor string, limit int) ([]*Allocation, string, error)
	GetActiveAllocations(ctx context.Context, goalID string) ([]*Allocation, error)
	GetTotalReserved(ctx context.Context, fundingSourceID string) (int64, error)
}
