package domain

import "context"

type Repository interface {
	Save(ctx context.Context, p *Portfolio) error
	UpdateStatus(ctx context.Context, id string, from, to PortfolioStatus) error
	GetByID(ctx context.Context, id string) (*Portfolio, error)
	ListByUser(ctx context.Context, ownerID string, cursor string, limit int) ([]*Portfolio, string, error)
	ListByType(ctx context.Context, ownerID string, pt PortfolioType, cursor string, limit int) ([]*Portfolio, string, error)
	GetPortfolioHistory(ctx context.Context, id string, cursor string, limit int) ([]*Portfolio, string, error)
}
