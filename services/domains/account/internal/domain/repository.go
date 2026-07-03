package domain

import "context"

type Repository interface {
	Save(ctx context.Context, account *Account) error
	UpdateStatus(ctx context.Context, id string, from, to AccountStatus) error
	GetByID(ctx context.Context, id string) (*Account, error)
	ListByUser(ctx context.Context, ownerID string, cursor string, limit int) ([]*Account, string, error)
	ListByType(ctx context.Context, ownerID string, at AccountType, cursor string, limit int) ([]*Account, string, error)
	ListByStatus(ctx context.Context, ownerID string, status AccountStatus, cursor string, limit int) ([]*Account, string, error)
	ListByInstitution(ctx context.Context, institutionID string, cursor string, limit int) ([]*Account, string, error)
	ListByHousehold(ctx context.Context, householdID string, cursor string, limit int) ([]*Account, string, error)
}
