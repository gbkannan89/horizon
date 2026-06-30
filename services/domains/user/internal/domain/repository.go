package domain

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	UpdateStatus(ctx context.Context, userID string, fromStatus, toStatus UserStatus) error

	GetByID(ctx context.Context, userID string) (*User, error)
	ListByStatus(ctx context.Context, status UserStatus, cursor string, limit int) ([]*User, string, error)
	Search(ctx context.Context, query string, cursor string, limit int) ([]*User, string, error)
	GetConsents(ctx context.Context, userID string) ([]*ConsentRecord, error)
	GetHouseholdMemberships(ctx context.Context, userID string) ([]string, error)
}
