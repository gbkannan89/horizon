package domain

import "context"

type Repository interface {
	Save(ctx context.Context, r *Rule) error
	GetByID(ctx context.Context, id string) (*Rule, error)
	ListByUser(ctx context.Context, userID string) ([]*Rule, error)
	ListByCategory(ctx context.Context, category string) ([]*Rule, error)
	Delete(ctx context.Context, id string) error
}
