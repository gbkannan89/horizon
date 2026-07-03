package domain

import "context"

type Repository interface {
	Save(ctx context.Context, household *Household) error
	UpdateStatus(ctx context.Context, id string, from, to HouseholdStatus) error
	GetByID(ctx context.Context, id string) (*Household, error)
	ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*Household, string, error)
	ListByStatus(ctx context.Context, status HouseholdStatus, cursor string, limit int) ([]*Household, string, error)
}
